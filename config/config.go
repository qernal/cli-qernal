package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (

	// represents the current  config
	Current Config
)

type ErrNoTokenFound struct{}

func (m *ErrNoTokenFound) Error() string {
	return "no token found"
}

type Config struct {
	Token string `yaml:"token"`
}

func Load() {

	token, found := os.LookupEnv("QERNAL_TOKEN")

	if found && token != "" {
		Current.Token = token
		return
	}

	// lookup token in config
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.SetConfigFile(fmt.Sprintf("%s/%s/%s", home, ".qernal", "config.yaml"))

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Print(fmt.Errorf("unable to read qernal config: %w", err).Error())
	}

	token = viper.Get("token").(string)
	if Current.Token != "" {
		Current.Token = token
		return
	}

	// TODO: prompt user for token
	fmt.Println("no token found")
}
