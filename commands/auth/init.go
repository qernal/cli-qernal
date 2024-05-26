package auth

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/qernal/cli-qernal/charm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Qernalconfig struct {
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
			token, err := GetQernalToken()
			if err != nil {
				return err
			}

			return saveConfig(token)
		},
	}

	cfgPath = filepath.Join(os.Getenv("HOME"), ".qernal", "config.yaml")
)

func GetQernalToken() (string, error) {
	// 1. Check environment variable
	if token := os.Getenv("QERNAL_TOKEN"); token != "" {
		fmt.Println(charm.SuccessStyle.Render("configuring CLI using environment variable ✅"))

		return token, nil
	}

	// 2. Check config file
	if config, err := readConfig(cfgPath); err == nil {
		if err := validatePermissions(cfgPath); err != nil {
			fmt.Println(charm.WarningStyle.Render(err.Error())) // Use custom style
		}
		return config.Token, nil
	} else if os.IsNotExist(err) {
		// File doesn't exist, continue to prompt user
		token, err := charm.GetSensitiveInput("clientid@clientsecret", "")
		if err != nil {
			fmt.Println(charm.ErrorStyle.Render(fmt.Sprintf("error retrieving input %s", err.Error())))
			return "", err
		}
		return token, nil

	}
	token, err := charm.GetSensitiveInput("Enter your token", "")
	if err != nil {
		fmt.Println(charm.ErrorStyle.Render(fmt.Sprintf("error retrieving input %s", err.Error())))
		return "", err
	}
	return token, nil
}

func readConfig(cfgPath string) (Qernalconfig, error) {
	viper.SetConfigFile(cfgPath)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return Qernalconfig{}, fmt.Errorf("error reading config file, %s", err)
	}

	// Unmarshal the config into a struct
	var cfg Qernalconfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return Qernalconfig{}, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return cfg, nil
}

// TODO: use viper to only save update values
func saveConfig(token string) error {
	cfg := &Qernalconfig{Token: token}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cfgPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0600)
}

func validatePermissions(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Check if owner has read and write access, others don't have any access
	if fileInfo.Mode()&os.ModePerm != 0600 {
		_, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		return fmt.Errorf(
			"WARNING: Qernal configuration file is readable by others, set the permissions to 600 on file %s\nYou can run 'chmod 600 %s' to fix this.",
			filePath,
			filePath,
		)
	}
	return nil
}
