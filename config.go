package happening

import (
	goconfig "github.com/msbranco/goconfig"
	"reflect"
)

type Config struct {
	Daemon      bool   `ini:"daemonize"`
	LogFile     string `ini:"log_file"`
	LogLevel    string `ini:"log_level"`
	Pidfile     string `ini:"pidfile"`
	StoragePath string `ini:"storage_path"`
}

func NewConfig() *Config {
	return &Config{
		Daemon:      DEFAULT_DAEMON_MODE,
		LogLevel:    DEFAULT_LOG_LEVEL,
		LogFile:     DEFAULT_LOG_FILE,
		Pidfile:     DEFAULT_PID_FILE,
		StoragePath: DEFAULT_STORAGE_PATH,
	}
}

func (c *Config) FromFile(path string, section string) error {
	return loadConfigFromFile(path, c, section)
}

func loadConfigFromFile(path string, obj interface{}, section string) error {
	iniConfig, err := goconfig.ReadConfigFile(path)
	if err != nil {
		return err
	}

	config := reflect.ValueOf(obj).Elem()
	configType := config.Type()

	for i := 0; i < config.NumField(); i++ {
		structField := config.Field(i)
		fieldTag := configType.Field(i).Tag.Get("ini")

		switch {
		case structField.Type().Kind() == reflect.Bool:
			config_value, err := iniConfig.GetBool(section, fieldTag)
			if err == nil {
				structField.SetBool(config_value)
			}
		case structField.Type().Kind() == reflect.String:
			config_value, err := iniConfig.GetString(section, fieldTag)
			if err == nil {
				structField.SetString(config_value)
			}
		case structField.Type().Kind() == reflect.Int:
			config_value, err := iniConfig.GetInt64(section, fieldTag)
			if err == nil {
				structField.SetInt(config_value)
			}
		}
	}

	return nil
}

// A bit verbose, and not that dry, but could not find
// more clever for now.
func (c *Config) UpdateFromCmdline(cmdline *Cmdline) {
	if *cmdline.DaemonMode != DEFAULT_DAEMON_MODE {
		c.Daemon = *cmdline.DaemonMode
	}

	if *cmdline.LogLevel != DEFAULT_LOG_LEVEL {
		c.LogLevel = *cmdline.LogLevel
	}
}
