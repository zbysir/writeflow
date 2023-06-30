package apiservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zbysir/writeflow/internal/pkg/auth"
	"github.com/zbysir/writeflow/internal/pkg/config"
	"github.com/zbysir/writeflow/internal/pkg/http_file_server"
	"github.com/zbysir/writeflow/internal/pkg/httpsrv"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/internal/repo"
	"github.com/zbysir/writeflow/internal/usecase"
	"github.com/zbysir/writeflow/pkg/writeflowui"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type Config struct {
	Secret        string // 单机部署，输入 secret 就能使用
	ListenAddress string
}
type ApiService struct {
	config Config

	flowRepo    repo.Flow
	flowUsecase *usecase.Flow
}

func NewApiService(config Config, flowRepo repo.Flow, sysRepo repo.System) (*ApiService, error) {
	flow, err := usecase.NewFlow(flowRepo, sysRepo)
	if err != nil {
		return nil, err
	}
	return &ApiService{config: config, flowRepo: flowRepo, flowUsecase: flow}, nil
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*") // 这是允许访问所有域
		}
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, Cookie")
		c.Header("Access-Control-Allow-Methods", "OPTIONS,GET,PUT,POST,DELETE")
		c.Header("Access-Control-Allow-Credentials", "true") //  跨域请求是否需要带 cookie 信息 默认设置为 true

		// 放行所有 OPTIONS 方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
			c.Abort()
			return
		}

		c.Next()
	}
}

var AuthErr = errors.New("need login")

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, e := range c.Errors {
			err := e.Err
			code := 400
			if errors.Is(err, AuthErr) {
				code = 401
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"code": code,
				"msg":  err.Error(),
			})

			return
		}
	}

}

func Auth(secret string) gin.HandlerFunc {
	if secret == "" {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		t, _ := c.Cookie("token")
		if t == "" {
			c.Error(AuthErr)
			c.Abort()
			return
		}
		if !auth.CheckToken(secret, t) {
			c.Error(AuthErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (a *ApiService) Run(ctx context.Context, addr string) (err error) {
	if !config.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(Cors())

	r.NoRoute(func(c *gin.Context) {
		proto := "http"
		host := c.Request.Host

		referer := c.Request.Referer()
		if referer != "" {
			parse, err := url.Parse(referer)
			if err != nil {
				log.Errorf("parse referer err: %+v", err)
			}
			host = parse.Host
			proto = parse.Scheme
		}

		apiHost := fmt.Sprintf("%s://%s", proto, host)
		wsHost := ""
		if proto == "https" {
			wsHost = fmt.Sprintf("wss://%s", host)
		} else {
			wsHost = fmt.Sprintf("ws://%s", host)
		}

		d := writeflowui.UIFs(writeflowui.UIConfig{
			ApiHost: apiHost,
			WsHost:  wsHost,
		})

		// try file
		p := strings.TrimPrefix(c.Request.URL.Path, "/")
		f, err := d.Open(p)
		if err != nil {
			f, err = d.Open(filepath.Join(p, "index.html"))
			if err != nil {
				// fallback to index.html
				c.Request.URL.Path = ""
			} else {
				f.Close()
			}
		} else {
			f.Close()
		}

		http_file_server.WrapEtagHandler(http.FileServer(http.FS(d))).ServeHTTP(c.Writer, c.Request)
	})

	var api = r.Group("/api").Use(ErrorHandler(), Cors())

	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	api.GET("/ws/:topic", func(c *gin.Context) {
		topic := c.Param("topic")
		if topic == "" {
			c.Error(errors.New("need topic"))
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.Error(err)
			return
		}
		a.flowUsecase.AddWs(topic, conn)
	})

	api.POST("/auth", func(c *gin.Context) {
		// 创建 token
		var p struct {
			Secret string `json:"secret"`
		}
		err = c.BindJSON(&p)
		if err != nil {
			c.Error(err)
			return
		}

		// 如果是空，则验证 token
		if p.Secret == "" {
			t, _ := c.Cookie("token")
			if t == "" {
				c.Error(AuthErr)
				return
			}
			if !auth.CheckToken(a.config.Secret, t) {
				c.Error(AuthErr)
				return
			}

			c.JSON(200, "ok")
			return
		}
		log.Infof("a.config.Secret, %v %v", a.config.Secret, p.Secret)
		if p.Secret != a.config.Secret {
			c.Error(AuthErr)
			return
		}
		t := auth.CreateToken(p.Secret)
		c.SetCookie("token", t, 7*24*3600, "", c.Request.Host, false, true)
		c.JSON(200, "ok")
	})

	apiAuth := api.Use(Auth(a.config.Secret))

	a.RegisterFlow(apiAuth)

	s, err := httpsrv.NewService(addr)
	if err != nil {
		return
	}
	s.Handler("/", r.Handler().ServeHTTP)
	err = s.Start(ctx)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
			log.Infof("http service shutdown")
		} else {
			return err
		}
	}
	return nil
}
