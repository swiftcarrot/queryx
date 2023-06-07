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
// use in `sql.Open("adapter", config.ConnectionURL())`
// TODO: add examples for different adapters for reference
func (c *Config) ConnectionURL(withDatabase bool) string {
	rawURL := ""
	if c.URL.EnvKey != "" {
		rawURL = os.Getenv(c.URL.EnvKey)
	} else {
		// TODO: construct connection url
		//switch c.Adapter {
		//case "postgresql":
		//	rawURL = fmt.Sprintf("postgres://%+v:%+v@%+v:%+v/%+v?sslmode=disable",c.User,c.Password,c.Host,c.Port,c.Database)
		//case "mysql":
		//	rawURL =fmt.Sprintf("%+v:%+v@tcp(%+v:%+v)/%+v?parseTime=true",c.User,c.Password,c.Host,c.Port,c.Database)
		//case "sqlite":
		//	rawURL = fmt.Sprintf("file:%+v",c.Host)
		//}
		rawURL = c.URL.Value
	}

	u, _ := url.Parse(rawURL)

	if !withDatabase {
		// TODO: populate config from connection url
		if c.Adapter == "mysql" {
			i := strings.LastIndexAny(rawURL, "/")
			i2 := strings.Index(rawURL, "?")
			c.Database = rawURL[i+1 : i2]
			rawURL = strings.Replace(rawURL, rawURL[i+1:i2], "", -1)
			return rawURL
		} else if c.Adapter == "sqlite" {
			database := strings.Split(rawURL, ":")[1]
			c.Database = database
		} else {
			c.Database = u.Path[1:]
			u.Path = ""
		}

		return u.String()
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
	u, _ := url.Parse(rawURL)
	if c.Adapter == "mysql" {
		i := strings.LastIndexAny(rawURL, "/")
		i2 := strings.Index(rawURL, "?")
		c.Database = rawURL[i+1 : i2]
	} else {
		c.Database = u.Path[1:]
		u.Path = ""
	}
	return c.Database
}
