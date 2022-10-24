package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var Vars = struct {
	Host             string `mapstructure:"HOST"`
	PixapiReTokenKey string `mapstructure:"pixapi_re_token_key"`
	PixapiTokenKey   string `mapstructure:"pixapi_token_key"`
	PixapiTokenTime  string `mapstructure:"pixapi_token_time"`
	Threadmax        int    `mapstructure:"thread_max"`
	Vision           string `mapstructure:"vision"`
}{}

func VarsConfigInit() {
	var viper_init = viper.New()
	viper_init.SetConfigName("config-settings")
	viper_init.SetConfigType("json")
	viper_init.AddConfigPath(".")
	// path to look for the config file in
	viper_init.SetDefault("ThreadMax", 16)
	viper_init.SetDefault("vision", "1.0.9")
	viper_init.SetDefault("PIXAPI_TOKEN_KEY", "")
	viper_init.SetDefault("PIXAPI_RE_TOKEN_KEY", "")
	viper_init.SetDefault("PIXAPI_TOKEN_TIME", time.Now())
	viper_init.SetDefault("HOST", "https://app-api.pixiv.net")

	if err := viper_init.ReadInConfig(); err != nil {
		fmt.Printf("Read config-settings.json Error:[%v]\n", err)
		if err = viper_init.SafeWriteConfig(); err != nil {
			fmt.Println("Safe write config file failed,please check the permission or create config-settings.json file manually.")
			fmt.Println("Detailed error message as follows:", err)
		} else {
			fmt.Println("safe write config file success!")
		}
	} else {
		if err = viper_init.WriteConfig(); err != nil {
			fmt.Println("Update config file failed,please check the permission.")
		}
	}
	if err := viper_init.Unmarshal(&Vars); err != nil {
		panic(err)
	}
}
