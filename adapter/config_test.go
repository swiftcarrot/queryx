package adapter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swiftcarrot/queryx/schema"
	"github.com/swiftcarrot/queryx/types"
)

func TestNewPostgreSQLConfig(t *testing.T) {
	url := "postgres://postgres:postgres@localhost:5432/queryx_test?sslmode=disable"
	config := NewConfig(&schema.Config{Adapter: "postgresql", URL: types.StringOrEnv{Value: url}})
	require.Equal(t, url, config.URL)
	require.Equal(t, "postgres://postgres:postgres@localhost:5432?sslmode=disable", config.URL2)
	require.Equal(t, "queryx_test", config.Database)
}

func TestNewMySQLConfig(t *testing.T) {
	url := "root@tcp(localhost:3306)/queryx_test?parseTime=true&loc=Asia%2FShanghai"
	config := NewConfig(&schema.Config{Adapter: "mysql", URL: types.StringOrEnv{Value: url}})
	require.Equal(t, url, config.URL)
	require.Equal(t, "root@tcp(localhost:3306)/", config.URL2)
	require.Equal(t, "queryx_test", config.Database)
}

func TestNewSQLiteConfig(t *testing.T) {
	url := "file:test.sqlite3"
	config := NewConfig(&schema.Config{Adapter: "sqlite", URL: types.StringOrEnv{Value: url}})
	require.Equal(t, url, config.URL)
	require.Equal(t, "", config.URL2)
	require.Equal(t, "test.sqlite3", config.Database)
}
