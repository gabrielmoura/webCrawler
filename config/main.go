package config

import (
	"errors"
	"flag"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"github.com/gabrielmoura/go/pkg/ternary"
	"github.com/spf13/viper"
	"regexp"
	"strings"
)

var (
	MaxConcurrency = flag.Int("maxConcurrency", 10, "Max number of concurrent requests")
	MaxDepth       = flag.Int("maxDepth", 2, "Max depth to crawl")
	enabledConfig  = flag.Bool("config", false, "Enable config file")
	enableProxy    = flag.Bool("proxy", false, "Enable Proxy")
	proxyURL       = flag.String("proxyURL", "http://localhost:4444", "Proxy URL")
	inicialURL     = flag.String("url", "https://www.uol.com.br", "URL inicial")
	cacheMode      = flag.Bool("mem", false, "Cache mode")
	// tlds list of Top-Level Domains
	tlds        = flag.String("tlds", "", "TLDs to filter EX: com,br,org")
	postgresURI = flag.String("postgresURI", "postgres://postgres:Strong@P4ssword@localhost/crawler", "Postgres URI")
)

func splitComma(txt string) []string {
	if txt == "" {
		return []string{}
	}
	return strings.Split(txt, ",")
}

type Config struct {
	MaxConcurrency int          `mapstructure:"MAX_CONCURRENCY"`
	MaxDepth       int          `mapstructure:"MAX_DEPTH"`
	PostgresURI    string       `mapstructure:"POSTGRES_URI"`
	AppName        string       `mapstructure:"APP_NAME"`
	TimeFormat     string       `mapstructure:"TIME_FORMAT"`
	TimeZone       string       `mapstructure:"TIME_ZONE"`
	InicialURL     string       `mapstructure:"URL"`
	Cache          *CacheConfig `mapstructure:"CACHE"`
	Proxy          *Proxy       `mapstructure:"PROXY"`
	Filter         *Filter      `mapstructure:"FILTER"`
}
type CacheConfig struct {
	DBDir string `mapstructure:"DB_DIR"`
	Mode  string `mapstructure:"MODE"` // "mem" or "disc'
}
type Proxy struct {
	Enabled  bool   `mapstructure:"ENABLED"`
	ProxyURL string `mapstructure:"PROXY_URL"`
}
type Filter struct {
	Tlds        []string `mapstructure:"TLDS"`
	IgnoreLocal bool     `mapstructure:"IGNORE_LOCAL"`
}

var Conf *Config

func loadByFlag() error {
	cfg := &Config{

		AppName:        "WebCrawler",
		TimeFormat:     "02-Jan-2006",
		TimeZone:       "America/Sao_Paulo",
		MaxConcurrency: *MaxConcurrency,
		MaxDepth:       *MaxDepth,
		PostgresURI:    *postgresURI,
		InicialURL:     *inicialURL,
		Cache: &CacheConfig{
			DBDir: "/tmp/WebCrawler",
			Mode:  ternary.Ternary(*cacheMode, "mem", "disc"),
		},
		Proxy: &Proxy{
			Enabled:  *enableProxy,
			ProxyURL: *proxyURL,
		},
		Filter: &Filter{
			Tlds:        splitComma(*tlds),
			IgnoreLocal: false,
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
	vip.SetDefault("APP_NAME", "WebCrawler")
	vip.SetDefault("TIME_FORMAT", "02-Jan-2006")
	vip.SetDefault("TIME_ZONE", "America/Sao_Paulo")
	vip.SetDefault("DB_DIR", "/tmp/WebCrawler")
	vip.SetDefault("POSTGRES_URI", "postgres://postgres:Strong@P4ssword@localhost/crawler")

	vip.SetDefault("PROXY.ENABLED", false)
	vip.SetDefault("PROXY.PROXY_URL", "http://localhost:4444")

	vip.SetDefault("FILTER.TLDS", []string{})
	vip.SetDefault("FILTER.IGNORE_LOCAL", false)

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
	if vip.IsSet("FILTER.TLDS") {
		vip.Set("FILTER.TLDS", splitComma(vip.GetString("FILTER.TLDS")))
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
		log.Logger.Info("Carregando configurações do arquivo")
		return loadByConfigFile()
	} else {
		log.Logger.Info("Carregando configurações por flag")
		return loadByFlag()
	}
}
