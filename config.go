package happening

import (
    goconfig "github.com/msbranco/goconfig"
    "reflect"
)

type Config struct {
    Pidfile     string `ini:"pidfile"`
}


func NewConfig() *Config {
    return &Config{
        Pidfile: "/var/run/happening.pid",
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
