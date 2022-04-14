package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() error {
	home, _ := os.UserHomeDir()

	cfgFile := "config.yml"
	viper.SetConfigName(cfgFile)
	viper.SetConfigType("yaml")

	viper.AddConfigPath(home)

	cfgPath := filepath.Join(home, ".config", "gh-changelog")
	viper.AddConfigPath(cfgPath)

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgPath, 0750); err != nil {
			return fmt.Errorf("failed to create config directory: %s", err)
		}
	}

	cfgFilePath := filepath.Join(cfgPath, cfgFile)

	if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
		_, err := os.Create(filepath.Clean(cfgFilePath))
		if err != nil {
			return fmt.Errorf("failed to initialise %s: %s", cfgFilePath, err)
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %s", err)
	}

	setDefaults()

	return nil
}

func setDefaults() {
	viper.SetDefault("fileName", "CHANGELOG.md")
	viper.SetDefault("excludedLabels", []string{"maintenance"})

	sections := make(map[string][]string)
	sections["Changed"] = []string{"backwards-incompatible"}
	sections["Added"] = []string{"feature", "enhancement"}
	sections["Fixed"] = []string{"bug", "bugfix", "documentation"}

	viper.SetDefault("sections", sections)
}
