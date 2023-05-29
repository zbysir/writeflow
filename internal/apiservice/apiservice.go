package apiservice

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zbysir/writeflow/internal/pkg/auth"
	"github.com/zbysir/writeflow/internal/pkg/config"
	"github.com/zbysir/writeflow/internal/pkg/httpsrv"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/internal/repo"
	"net/http"
)

type Config struct {
	Secret string // 单机部署，输入 secret 就能使用
}
type ApiService struct {
	config Config

	flowRepo repo.Flow
}

func NewApiService(config Config, flowRepo repo.Flow) *ApiService {
	return &ApiService{config: config, flowRepo: flowRepo}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               // 请求方法
		origin := c.Request.Header.Get("Origin") // 请求头部

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*") // 这是允许访问所有域
		}
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, Cookie")
		c.Header("Access-Control-Allow-Methods", "OPTIONS,GET,PUT,POST,DELETE")
		c.Header("Access-Control-Allow-Credentials", "true") //  跨域请求是否需要带cookie信息 默认设置为true

		//放行所有 OPTIONS 方法
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
			log.Infof("3 %v", err)

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

// localhost:9090/api/file/tree
func (a *ApiService) Run(ctx context.Context, addr string) (err error) {
	if !config.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(Cors())

	var gateway = r.Group("/").Use(ErrorHandler())

	gateway.Use(Auth(a.config.Secret)).GET("/ws/:key", func(c *gin.Context) {
		//key := c.Param("key")
		//conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		//if err != nil {
		//	c.Error(err)
		//	return
		//}
		//a.hub.Add(key, conn)
	})

	var api = r.Group("/api").Use(ErrorHandler(), Cors())
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
		// TODO Get domain from query host
		c.SetCookie("token", t, 7*24*3600, "", "localhost:9433", false, true)
		c.JSON(200, "ok")
	})

	apiAuth := api.Use(Auth(a.config.Secret))

	a.RegisterRepo(apiAuth)

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
