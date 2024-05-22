package config

import (
	"errors"
	"flag"
	"github.com/spf13/viper"
	"log"
	"regexp"
)

var (
	MaxConcurrency = flag.Int("maxConcurrency", 10, "Max number of concurrent requests")
	MaxDepth       = flag.Int("maxDepth", 2, "Max depth to crawl")
	enabledConfig  = flag.Bool("config", false, "Enable config file")
	enabledI2P     = flag.Bool("i2p", false, "Enable I2P")

	inicialURL = flag.String("url", "https://www.uol.com.br", "URL inicial")
)

type Config struct {
	MaxConcurrency int
	MaxDepth       int
	MongoURI       string
	AppName        string  `mapstructure:"APP_NAME"`
	TimeFormat     string  `mapstructure:"TIME_FORMAT"`
	TimeZone       string  `mapstructure:"TIME_ZONE"`
	DBDir          string  `mapstructure:"DB_DIR"`
	I2PCfg         *I2PCfg `mapstructure:"I2P_CFG"`
	InicialURL     string  `mapstructure:"URL"`
}
type I2PCfg struct {
	Enabled         bool   `mapstructure:"ENABLED"`
	HttpHostAndPort string `mapstructure:"HTTP_HOST_AND_PORT"`
	Host            string `mapstructure:"HOST"`
	Url             string `mapstructure:"URL"`
	HttpsUrl        string `mapstructure:"HTTPS_URL"`
	SAMAddr         string `mapstructure:"SAM_ADDR"`
	KeyPath         string `mapstructure:"KEY_PATH"`
}

var Conf *Config

func loadByFlag() error {
	cfg := &Config{

		AppName:        "DavServer",
		TimeFormat:     "02-Jan-2006",
		TimeZone:       "America/Sao_Paulo",
		MaxConcurrency: *MaxConcurrency,
		MaxDepth:       *MaxDepth,
		MongoURI:       "mongodb://root:Strong%40P4word@localhost:27017",
		DBDir:          "/tmp/badgerDB",
		InicialURL:     *inicialURL,
		I2PCfg: &I2PCfg{
			Enabled:         *enabledI2P,
			HttpHostAndPort: "127.0.0.1:7672",
			Host:            "",
			Url:             "127.0.0.1:7672",
			HttpsUrl:        "",
			SAMAddr:         "127.0.0.1:7656",
			KeyPath:         "./",
		},
	}
	// Atualiza a variável global Conf
	Conf = cfg
	return nil
}
func loadByConfigFile() error {
	var cfg Config
	vip := viper.New()

	// Definindo valores padrão
	vip.SetDefault("APP_NAME", "DavServer")
	vip.SetDefault("TIME_FORMAT", "02-Jan-2006")
	vip.SetDefault("TIME_ZONE", "America/Sao_Paulo")
	vip.SetDefault("DB_DIR", "/tmp/badgerDB")

	vip.SetDefault("I2P_CFG.ENABLED", false)
	vip.SetDefault("I2P_CFG.SAM_ADDR", "127.0.0.1:7656")
	vip.SetDefault("I2P_CFG.HTTP_HOST_AND_PORT", "127.0.0.1:7672")
	vip.SetDefault("I2P_CFG.KEY_PATH", "./")

	// Lendo o arquivo de configuração conf.yml
	vip.SetConfigName("conf")
	vip.SetConfigType("yml")
	vip.AddConfigPath(".")
	vip.AddConfigPath("/opt/crw")
	vip.AddConfigPath("/etc/crw")

	// Lendo as configurações do arquivo conf.yml
	if err := vip.ReadInConfig(); err != nil {
		// Se o arquivo conf.yml não for encontrado, continue sem erro
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}
	// Se APP_NAME não estiver definido no padrão, retorne um erro
	if vip.IsSet("APP_NAME") {
		regex := regexp.MustCompile("^[A-Za-z0-9]+$")
		if !regex.MatchString(vip.GetString("APP_NAME")) {
			return errors.New("APP_NAME só pode conter letras e números")
		}
	}

	// Atribua as configurações ao cfg
	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}

	// Atualiza a variável global Conf
	Conf = &cfg

	return nil
}
func LoadConfig() error {
	flag.Parse()
	if *enabledConfig {
		log.Printf("Carregando configurações do arquivo")
		return loadByConfigFile()
	} else {
		log.Printf("Carregando configurações por flag")
		return loadByFlag()
	}
}
