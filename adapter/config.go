package adapter

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/swiftcarrot/queryx/schema"
)

type Config struct {
	Adapter string
	// full connection url
	URL string
	// connection url without database
	URL2 string
	// database, for sqlite it is sqlite file path
	Database string
	Host     string
	Port     string
	Username string
	Password string
	Options  url.Values
}

func NewConfigFromURL(rawURL string) *Config {
	c := &Config{}

	u, _ := url.Parse(rawURL)

	c.Adapter = u.Scheme
	c.Username = u.User.Username()
	c.Password, _ = u.User.Password()
	c.Host, c.Port, _ = net.SplitHostPort(u.Host)
	c.Database = strings.TrimPrefix(u.Path, "/")
	c.Options, _ = url.ParseQuery(u.RawQuery)

	return c
}

func NewConfig(cfg *schema.Config) *Config {
	rawURL := ""
	if cfg.URL.EnvKey != "" {
		rawURL = os.Getenv(cfg.URL.EnvKey)
	} else {
		rawURL = cfg.URL.Value
	}

	c := NewConfigFromURL(rawURL)
	c.URL = c.GoFormat()

	db := c.Database
	c.Database = ""
	c.URL2 = c.GoFormat()
	c.Database = db

	return c
}

func (c *Config) GoFormat() string {
	var u string

	switch c.Adapter {
	case "postgresql":
		u = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
	case "mysql":
		c.Options.Set("parseTime", "true")
		u = fmt.Sprintf("%s@tcp(%s:%s)/%s", c.Username, c.Host, c.Port, c.Database)
	case "sqlite":
		return ""
	default:
		return ""
	}

	options := c.Options.Encode()
	if options != "" {
		u = u + "?" + options
	}

	return u
}

func (c *Config) TSFormat() string {
	var u string

	switch c.Adapter {
	case "postgresql":
		u = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
	case "mysql":
		u = fmt.Sprintf("%s@tcp(%s:%s)/%s", c.Username, c.Host, c.Port, c.Database)
	case "sqlite":
		return ""
	default:
		return ""
	}

	options := c.Options.Encode()
	if options != "" {
		u = u + "?" + options
	}

	return u
}
