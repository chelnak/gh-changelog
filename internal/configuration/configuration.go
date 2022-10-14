// Package configuration contains a number of methods that are used
// to provide configuration to the wider application. It uses viper
// to pull config from either the environment or a config file then
// unmarhsals the config into the configuration struct. The configuration struct
// is made available to the application via a package level variable
// called Config.
package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var Config configuration

type configuration struct {
	FileName                string              `mapstructure:"file_name" yaml:"file_name" json:"fileName"`
	ExcludedLabels          []string            `mapstructure:"excluded_labels" yaml:"excluded_labels" json:"excludedLabels"`
	Sections                map[string][]string `mapstructure:"sections" yaml:"sections" json:"sections"`
	SkipEntriesWithoutLabel bool                `mapstructure:"skip_entries_without_label" yaml:"skip_entries_without_label" json:"skipEntriesWithoutLabel"`
	ShowUnreleased          bool                `mapstructure:"show_unreleased" yaml:"show_unreleased" json:"showUnreleased"`
	CheckForUpdates         bool                `mapstructure:"check_for_updates" yaml:"check_for_updates" json:"checkForUpdates"`
	Logger                  string              `mapstructure:"logger" yaml:"logger" json:"logger"`
}

type writeOptions struct {
	data      string
	lexerName string
	noColor   bool
	writer    io.Writer
}

func prettyWrite(opts writeOptions) error {
	lexer := lexers.Get(opts.lexerName)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)

	style := styles.Get("native")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal16m")

	if opts.noColor {
		formatter = formatters.Get("noop")
	}

	iterator, err := lexer.Tokenise(nil, opts.data)
	if err != nil {
		return err
	}

	return formatter.Format(opts.writer, style, iterator)
}

func (c *configuration) PrintJSON(noColor bool, writer io.Writer) error {
	b, err := json.MarshalIndent(c, "", "  ")
	b = append(b, '\n')
	if err != nil {
		return err
	}

	opts := writeOptions{
		data:      string(b),
		lexerName: "json",
		noColor:   noColor,
		writer:    writer,
	}

	return prettyWrite(opts)
}

func (c *configuration) PrintYAML(noColor bool, writer io.Writer) error {
	b, err := yaml.Marshal(c)
	y := []byte("---\n")
	y = append(y, b...)
	if err != nil {
		return err
	}

	opts := writeOptions{
		data:      string(y),
		lexerName: "yaml",
		noColor:   noColor,
		writer:    writer,
	}

	return prettyWrite(opts)
}

func InitConfig() error {
	home, _ := os.UserHomeDir()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	cfgPath := filepath.Join(home, ".config", "gh-changelog")
	viper.AddConfigPath(cfgPath)

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgPath, 0750); err != nil {
			return fmt.Errorf("failed to create config directory: %s", err)
		}
	}

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		err := viper.SafeWriteConfig()
		if err != nil {
			return fmt.Errorf("failed to write config: %s", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("changelog")

	err := viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %s", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("file_name", "CHANGELOG.md")
	viper.SetDefault("excluded_labels", []string{"maintenance", "dependencies"})

	sections := make(map[string][]string)
	sections["changed"] = []string{"backwards-incompatible"}
	sections["added"] = []string{"feature", "enhancement"}
	sections["fixed"] = []string{"bug", "bugfix", "documentation"}

	viper.SetDefault("sections", sections)

	viper.SetDefault("skip_entries_without_label", false)

	viper.SetDefault("show_unreleased", true)

	viper.SetDefault("check_for_updates", true)

	viper.SetDefault("no_color", false)

	viper.SetDefault("logger", "spinner")
}
