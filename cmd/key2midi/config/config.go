package config

import "github.com/spf13/viper"

func Load() {
	viper.SetConfigType("json")
	viper.AddConfigPath("./")

	if err := viper.MergeInConfig(); err != nil {
		panic(err)
	}
}
