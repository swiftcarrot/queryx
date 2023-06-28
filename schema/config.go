package schema

import (
	"net/url"
	"os"
	"strings"

	"github.com/swiftcarrot/queryx/types"
)

type Config struct {
	Adapter     string
	Environment string
	URL         types.StringOrEnv
	Database    string
	Host        string
	Port        string
	User        string
	Password    string
	Encoding    string
	Pool        int
	Timeout     int
	Socket      string
	Options     map[string]string
	RawOptions  string
}

// return database connection url from config
// `sql.Open("adapter", config.ConnectionURL(true))`
func (c *Config) ConnectionURL(withDatabase bool) string {
	rawURL := ""
	if c.URL.EnvKey != "" {
		rawURL = os.Getenv(c.URL.EnvKey)
	} else {
		rawURL = c.URL.Value
	}

	if !withDatabase {
		if c.Adapter == "postgresql" {
			u, _ := url.Parse(rawURL)
			c.Database = u.Path[1:]
			u.Path = ""
			return u.String()
		} else if c.Adapter == "mysql" {
			parts := strings.Split(rawURL, "/")
			c.Database = strings.Split(parts[1], "?")[0]
			return parts[0] + "/"
		} else if c.Adapter == "sqlite" {
			database := strings.Split(rawURL, ":")[1]
			c.Database = database
		}
	}

	return rawURL
}

func (c *Config) GetDatabaseName() string {
	rawURL := ""
	if c.URL.EnvKey != "" {
		rawURL = os.Getenv(c.URL.EnvKey)
	} else if c.URL.Value != "" {
		rawURL = c.URL.Value
	}

	if c.Adapter == "postgresql" {
		u, _ := url.Parse(rawURL)
		c.Database = u.Path[1:]
	} else if c.Adapter == "mysql" {
		parts := strings.Split(rawURL, "/")
		c.Database = strings.Split(parts[1], "?")[0]
	} else if c.Adapter == "sqlite" {
	}

	return c.Database
}
