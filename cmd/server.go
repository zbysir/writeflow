package cmd

import (
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"github.com/zbysir/writeflow/internal/apiservice"
	"github.com/zbysir/writeflow/internal/pkg/config"
	"github.com/zbysir/writeflow/internal/pkg/db"
	"github.com/zbysir/writeflow/internal/pkg/signal"
	"github.com/zbysir/writeflow/internal/repo"
	"github.com/zbysir/writeflow/pkg/modules/llm"
	"strings"
)

type ApiParams struct {
	Address string `json:"address"`
	Secret  string `json:"secret"`
	OpenAI  OpenAI `json:"openai"`
	PGDB    PGDB   `json:"pgdb"`
}

type OpenAI struct {
	APIKey string `json:"apikey"`
}

type PGDB struct {
	Debug    bool   `json:"debug"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func Api() *cobra.Command {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cmd := &cobra.Command{
		Use:   "api",
		Short: "api start a api service",
		//Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := config.Get[ApiParams](v)
			if err != nil {
				return err
			}

			if p.Secret == "" {
				if !config.IsDebug() {
					p.Secret = funk.RandomString(8)
				}
			}

			//gin.SetMode(gin.ReleaseMode)
			//log.Infof("config: %+v", p)

			kvDb, err := db.NewKvDb("./data").Open("db", "default")
			if err != nil {
				return err
			}

			flowRepo := repo.NewBoltDBFlow(kvDb)
			sysRepo := repo.NewBoltDBSystem(kvDb)

			openAIClient := openai.NewClient(p.OpenAI.APIKey)

			documentRepo, err := repo.NewPGStorage(repo.DbConfig{
				Debug:    false,
				Host:     p.PGDB.Host,
				Port:     p.PGDB.Port,
				User:     p.PGDB.User,
				Password: p.PGDB.Password,
				DBName:   p.PGDB.DBName,
			}, llm.NewMarkDoneSplit(2048), llm.NewOpenAIEmbedding(openAIClient))
			if err != nil {
				return err
			}

			service, err := apiservice.NewApiService(apiservice.Config{Secret: p.Secret, ListenAddress: p.Address}, flowRepo, sysRepo, documentRepo)
			if err != nil {
				return err
			}

			ctx, c := signal.NewContext()
			defer c()
			err = service.Run(ctx, p.Address)
			if err != nil {
				return err
			}
			return nil
		},
	}

	config.DeclareFlag(v, cmd, "address", "a", ":9433", "service listen address")
	config.DeclareFlag(v, cmd, "secret", "c", "", "secret for web ui")
	config.DeclareFlag(v, cmd, "pgdb.password", "", "123456", "db password")
	config.DeclareFlag(v, cmd, "pgdb.host", "", "localhost", "db password")
	config.DeclareFlag(v, cmd, "pgdb.dbname", "", "writeflow", "db password")
	config.DeclareFlag(v, cmd, "pgdb.user", "", "postgres", "db password")
	config.DeclareFlag(v, cmd, "pgdb.port", "", "5432", "db password")
	config.DeclareFlag(v, cmd, "openai.apikey", "", "", "db password")

	return cmd
}
