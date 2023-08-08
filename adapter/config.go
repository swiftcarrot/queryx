package adapter

import (
	"net/url"
	"os"
	"strings"

	"github.com/swiftcarrot/queryx/schema"
)

type Config struct {
	Adapter     string
	Environment string
	// full connection url
	URL string
	// connection url without database
	URL2 string
	// database, for sqlite it is sqlite file path
	Database string
}

func NewConfig(cfg *schema.Config) *Config {
	rawURL := ""
	if cfg.URL.EnvKey != "" {
		rawURL = os.Getenv(cfg.URL.EnvKey)
	} else {
		rawURL = cfg.URL.Value
	}

	c := &Config{
		Adapter:     cfg.Adapter,
		Environment: cfg.Environment,
		URL:         rawURL,
	}

	switch cfg.Adapter {
	case "postgresql":
		u, _ := url.Parse(rawURL)
		c.Database = u.Path[1:]
		u.Path = ""
		c.URL2 = u.String()
	case "mysql":
		parts := strings.Split(rawURL, "/")
		c.Database = strings.Split(parts[1], "?")[0]
		c.URL2 = parts[0] + "/"
	case "sqlite":
		c.Database = strings.Split(rawURL, ":")[1]
	}

	return c
}
