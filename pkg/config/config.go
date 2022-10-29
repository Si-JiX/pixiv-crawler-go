package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var VarsFile *VarsConfig

var Vars = struct {
	Host              string `mapstructure:"host"`
	PixivRefreshToken string `mapstructure:"pixiv_refresh_token"`
	PixivToken        string `mapstructure:"pixiv_token"`
	UserID            int    `mapstructure:"pixiv_user_id"`
	ThreadMax         int    `mapstructure:"thread_max"`
	VersionName       string `mapstructure:"version_name"`
}{}

type VarsConfig struct {
	Vipers     *viper.Viper
	UpdateFile bool
}

func (is *VarsConfig) LoadConfig() {
	if err := is.Vipers.ReadInConfig(); err != nil {
		fmt.Printf("Read config-settings.json Error:[%v]\n", err)
		is.UpdateFile = false
	} else {
		is.UpdateFile = true
	}
}
func (is *VarsConfig) VarsUnmarshal() {
	if err := is.Vipers.Unmarshal(&Vars); err != nil {
		fmt.Println(err)
	}
}
func (is *VarsConfig) SaveConfig() {
	is.LoadConfig()
	if !is.UpdateFile {
		if err := is.Vipers.SafeWriteConfig(); err != nil {
			fmt.Println("Safe write config file failed,please check the permission or create config-settings.json file manually.")
			fmt.Println("Detailed error message as follows:", err)
		} else {
			fmt.Println("safe write config file success!")
		}
	} else {
		if err := is.Vipers.WriteConfig(); err != nil {
			fmt.Println("Update config file failed,please check the permission.")
		}
	}
	is.VarsUnmarshal()
}

func VarsConfigInit() {
	VarsFile = &VarsConfig{Vipers: viper.New()}
	VarsFile.Vipers.SetConfigName("config-settings")
	VarsFile.Vipers.SetConfigType("json")
	VarsFile.Vipers.AddConfigPath(".")
	// path to look for the config file in
	VarsFile.Vipers.SetDefault("host", "https://app-api.pixiv.net")
	VarsFile.Vipers.SetDefault("thread_max", 16)
	VarsFile.Vipers.Set("version_name", "1.3.8")
	VarsFile.Vipers.SetDefault("PIXIV_TOKEN", "")
	VarsFile.Vipers.SetDefault("PIXIV_REFRESH_TOKEN", "")
	VarsFile.SaveConfig()
}
