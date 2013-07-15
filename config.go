package happening

import (
	goconfig "github.com/msbranco/goconfig"
	"reflect"
)

type Config struct {
	Daemon   		bool   `ini:"daemonize"`
	LogFile  		string `ini:"log_file"`
	LogLevel 		string `ini:"log_level"`
	Pidfile  		string `ini:"pidfile"`
	StoragePath		string `ini:"storage_path"`
}

func NewConfig() *Config {
	return &Config{
		Daemon:   false,
		LogLevel: "INFO",
		LogFile:  "/var/log/happening.log",
		Pidfile:  "/var/run/happening.pid",
		StoragePath: "/var/lib/happening/",
	}
}

func (c *Config) FromFile(path string, section string) error {
	return loadConfigFromFile(path, c, section)
}

func loadConfigFromFile(path string, obj interface{}, section string) error {
	ini_config, err := goconfig.ReadConfigFile(path)
	if err != nil {
		return err
	}

	config := reflect.ValueOf(obj).Elem()
	config_type := config.Type()

	for i := 0; i < config.NumField(); i++ {
		struct_field := config.Field(i)
		field_tag := config_type.Field(i).Tag.Get("ini")

		switch {
		case struct_field.Type().Kind() == reflect.Bool:
			config_value, err := ini_config.GetBool(section, field_tag)
			if err == nil {
				struct_field.SetBool(config_value)
			}
		case struct_field.Type().Kind() == reflect.String:
			config_value, err := ini_config.GetString(section, field_tag)
			if err == nil {
				struct_field.SetString(config_value)
			}
		case struct_field.Type().Kind() == reflect.Int:
			config_value, err := ini_config.GetInt64(section, field_tag)
			if err == nil {
				struct_field.SetInt(config_value)
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
