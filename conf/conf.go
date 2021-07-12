package conf

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

type AppConfig struct {
	debug  bool
	base   string
	level  uint8 // log level
	Server struct {
		Grpc struct {
			Host string
			Port string
		}
		Rest struct {
			Host string
			Port string
		}
	}
	Client struct {
		Redis struct {
			Host     string
			Port     string
			User     string
			Password string
		}
		Postgres struct {
			Host     string
			Port     string
			User     string
			Password string
			Database string
		}
	}
	JWT struct {
		Secret        string
		RefreshSecret string
		TokenDur      time.Duration
	}
}

type appConfig interface {
	load() error
	get() *AppConfig
	Debug(bool)
	Mode() bool
	SetPath(string)
	Level(uint8)
}

var config appConfig

func init() {
	config = new(AppConfig)
}

func GetAppConfig() *AppConfig {
	return config.get()
}

func (c *AppConfig) get() *AppConfig {
	return c
}

func (c *AppConfig) Debug(debug bool) {
	fmt.Printf("setting up debug mode : %v \n", aurora.Yellow(debug))
	c.debug = debug
}

func (c *AppConfig) Mode() bool {
	return c.debug
}

func (c *AppConfig) SetPath(path string) {

	c.base = path

	if err := config.load(); err != nil {
		panic(fmt.Errorf("cannot unmarshal config %s", err))
	}
}

func (c *AppConfig) load() error {

	viper.AddConfigPath(c.base)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found: %v", aurora.Red(err)))
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("cannot unmarshal configs %v", err)
	}

	return nil
}

func (c *AppConfig) Level(level uint8) {
	c.level = level
}
