package adapter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swiftcarrot/queryx/schema"
	"github.com/swiftcarrot/queryx/types"
)

func TestNewPostgreSQLConfig(t *testing.T) {
	u := "postgresql://postgres:password@localhost:5432/queryx_test?sslmode=disable"
	c := NewConfig(&schema.Config{URL: types.StringOrEnv{Value: u}})
	require.Equal(t, "postgresql", c.Adapter)
	require.Equal(t, "postgres", c.Username)
	require.Equal(t, "password", c.Password)
	require.Equal(t, "localhost", c.Host)
	require.Equal(t, "5432", c.Port)
	require.Equal(t, "queryx_test", c.Database)
	require.Equal(t, "sslmode=disable", c.Options.Encode())
	require.Equal(t, "postgres://postgres:password@localhost:5432/queryx_test?sslmode=disable", c.URL)
	require.Equal(t, "postgres://postgres:password@localhost:5432/?sslmode=disable", c.URL2)
}

func TestNewMySQLConfig(t *testing.T) {
	url := "mysql://root:@localhost:3306/queryx_test"
	// url := "root@tcp(localhost:3306)/queryx_test?parseTime=true&loc=Asia%2FShanghai"
	c := NewConfig(&schema.Config{URL: types.StringOrEnv{Value: url}})
	require.Equal(t, "mysql", c.Adapter)
	require.Equal(t, "root", c.Username)
	require.Equal(t, "", c.Password)
	require.Equal(t, "localhost", c.Host)
	require.Equal(t, "3306", c.Port)
	require.Equal(t, "queryx_test", c.Database)

	require.Equal(t, "root@tcp(localhost:3306)/queryx_test?parseTime=true", c.URL)
	require.Equal(t, "root@tcp(localhost:3306)/?parseTime=true", c.URL2)
}

func TestNewSQLiteConfig(t *testing.T) {
	url := "sqlite:test.sqlite3"
	config := NewConfig(&schema.Config{Adapter: "sqlite", URL: types.StringOrEnv{Value: url}})
	require.Equal(t, url, config.URL)
	require.Equal(t, "", config.URL2)
	require.Equal(t, "test.sqlite3", config.Database)
}
