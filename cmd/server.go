package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
	"github.com/zbysir/writeflow/internal/apiservice"
	"github.com/zbysir/writeflow/internal/pkg/config"
	"github.com/zbysir/writeflow/internal/pkg/db"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/internal/pkg/signal"
	"github.com/zbysir/writeflow/internal/repo"
)

type ApiParams struct {
	Address string `json:"address"`
	Secret  string `json:"secret"`
}

func Api() *cobra.Command {
	v := viper.New()
	v.AutomaticEnv()

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
			log.Infof("config: %+v", p)

			kvDb, err := db.NewKvDb("./data").Open("db", "default")
			if err != nil {
				return err
			}

			flowRepo := repo.NewBoltDBFlow(kvDb)

			service := apiservice.NewApiService(apiservice.Config{Secret: p.Secret}, flowRepo)

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

	return cmd
}
