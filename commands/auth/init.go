package auth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qernal/cli-qernal/charm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type config struct {
	Token string `yaml:"token"`
}

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log in to your Qernal account",
		Long: `log in to your Qernal account by searching for the QERNAL_TOKEN environment variable first. 

 the order in which values are searched for:

1. **QERNAL_TOKEN environment variable:** If set, this is used as the token.
2. **$HOME/.qernal/config.yaml file:** If the environment variable is not found, the CLI checks for the token in this file.
3. **User input:** If neither of the above is found, the user is prompted to enter their Qernal token.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := getQernalToken()
			if err != nil {
				return err
			}

			return saveConfig(token)
		},
	}

	cfgPath = filepath.Join(os.Getenv("HOME"), ".qernal", "config.yaml")
)

func getQernalToken() (string, error) {
	// 1. Check environment variable
	if token := os.Getenv("QERNAL_TOKEN"); token != "" {
		fmt.Println(charm.SuccessStyle.Render("configuring CLI using environment variable ✅"))

		return token, nil
	}

	// 2. Check config file
	if token, err := readConfig(cfgPath); err == nil {
		fmt.Println(charm.SuccessStyle.Render(fmt.Sprintf("Using token from %s.", cfgPath)))
		return token, nil
	} else if os.IsNotExist(err) {
		// File doesn't exist, continue to prompt user
		token, err := charm.GetSensitiveInput("clientid@clientsecret", ".....")
		if err != nil {
			fmt.Println(charm.ErrorStyle.Render(fmt.Sprintf("error retrieving input %s", err.Error())))
			return "", err
		}
		return token, nil

	}
	token, err := charm.GetSensitiveInput("Enter your token", ".....")
	if err != nil {
		fmt.Println(charm.ErrorStyle.Render(fmt.Sprintf("error retrieving input %s", err.Error())))
		return "", err
	}
	return token, nil
}

func readConfig(cfgPath string) (string, error) {
	viper.SetConfigFile(cfgPath)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return "", fmt.Errorf("error reading config file, %s", err)
	}

	// Unmarshal the config into a struct
	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		return "", fmt.Errorf("unable to decode into struct, %v", err)
	}

	return cfg.Token, nil
}

func saveConfig(token string) error {
	cfg := &config{Token: token}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cfgPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0644)
}
