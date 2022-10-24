package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var Vars = struct {
	Host              string `mapstructure:"HOST"`
	PixivRefreshToken string `mapstructure:"PIXIV_REFRESH_TOKEN"`
	PixivToken        string `mapstructure:"PIXIV_TOKEN"`
	PixapiTokenTime   string `mapstructure:"pixapi_token_time"`
	ThreadMax         int    `mapstructure:"thread_max"`
	Vision            string `mapstructure:"version_name"`
}{}

var Vipers = viper.New()

func VarsConfigInit() {
	Vipers.SetConfigName("config-settings")
	Vipers.SetConfigType("json")
	Vipers.AddConfigPath(".")
	// path to look for the config file in
	Vipers.SetDefault("ThreadMax", 16)
	Vipers.SetDefault("vision", "1.0.9")
	Vipers.SetDefault("PIXAPI_TOKEN_KEY", "")
	Vipers.SetDefault("PIXAPI_RE_TOKEN_KEY", "")
	Vipers.SetDefault("PIXAPI_TOKEN_TIME", time.Now())
	Vipers.SetDefault("HOST", "https://app-api.pixiv.net")

	if err := Vipers.ReadInConfig(); err != nil {
		fmt.Printf("Read config-settings.json Error:[%v]\n", err)
		if err = Vipers.SafeWriteConfig(); err != nil {
			fmt.Println("Safe write config file failed,please check the permission or create config-settings.json file manually.")
			fmt.Println("Detailed error message as follows:", err)
		} else {
			fmt.Println("safe write config file success!")
		}
	} else {
		if err = Vipers.WriteConfig(); err != nil {
			fmt.Println("Update config file failed,please check the permission.")
		}
	}
	if err := Vipers.Unmarshal(&Vars); err != nil {
		panic(err)
	}
}
