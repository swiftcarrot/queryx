// Code generated by queryx, DO NOT EDIT.

package queryx

import ()

type Config struct {
	URL string
}

func NewConfig(env string) *Config {
	switch env {
	case "development":
		return &Config{
			URL: fixURL("postgresql://postgres:postgres@localhost:5432/test?sslmode=disable"),
		}
	case "test":
		return &Config{
			URL: fixURL(getenv("DATABASE_URL")),
		}
	}
	return nil
}
func fixURL(rawURL string) string {
	return rawURL
}
