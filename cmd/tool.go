package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zbysir/writeflow/internal/pkg/config"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"github.com/zbysir/writeflow/pkg/writeflow"
)

type toolParams struct {
	Address string `json:"address"`
	Secret  string `json:"secret"`
}

func Tool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tool",
		Short: "",
	}

	cmd.AddCommand(checkPlugin())
	return cmd
}

func checkPlugin() *cobra.Command {
	v := viper.New()

	type checkPluginParams struct {
		Url  string `json:"url"`
		Path string `json:"path"`
	}

	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Check plugin is available",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := config.Get[checkPluginParams](v)
			if err != nil {
				return err
			}

			// log.Infof("config: %+v", p)
			wf := writeflow.NewWriteFlow()
			if p.Url != "" {
				pm := writeflow.NewGoPkgPluginManager(nil)
				s, err := pm.Load(p.Url)
				if err != nil {
					return err
				}

				err = s.Register(wf)
				if err != nil {
					return err
				}
			} else if p.Path != "" {
				sf := writeflow.NewSysFs(p.Path)
				s := writeflow.NewGoPkgPlugin(sf, p.Path)
				err = s.Register(wf)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("url or path must be set")
			}

			log.Infof("plugin is available")
			return nil
		},
	}

	config.DeclareFlag(v, cmd, "path", "p", "", "plugin path")
	config.DeclareFlag(v, cmd, "url", "u", "", "plugin url")

	return cmd
}
