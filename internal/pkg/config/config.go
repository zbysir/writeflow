package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func DeclareFlag(v *viper.Viper, c *cobra.Command, name string, shorthand string, defaultVal any, usage string) {
	flags := c.PersistentFlags()

	switch defaultVal := defaultVal.(type) {
	case string:
		flags.StringP(name, shorthand, defaultVal, usage)
	}

	err := v.BindPFlag(name, flags.Lookup(name))
	if err != nil {
		panic(err)
	}
}

func GetAll(v *viper.Viper) map[string]interface{} {
	return v.AllSettings()
}

func Get[T any](v *viper.Viper) (T, error) {
	var t T
	err := v.Unmarshal(&t, func(config *mapstructure.DecoderConfig) {
		config.TagName = "json"
	})
	if err != nil {
		return t, err
	}
	return t, nil
}

func IsDebug() bool {
	s, ok := os.LookupEnv("DEBUG")
	if ok && s != "false" {
		return true
	}

	return false
}
