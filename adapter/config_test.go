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
	require.Equal(t, "postgres://postgres:password@localhost:5432/queryx_test?sslmode=disable", c.TSFormat())
}

func TestNewMySQLConfig(t *testing.T) {
	u := "mysql://root:@localhost:3306/queryx_test"
	// url := "root@tcp(localhost:3306)/queryx_test?parseTime=true&loc=Asia%2FShanghai"
	c := NewConfig(&schema.Config{URL: types.StringOrEnv{Value: u}})
	require.Equal(t, "mysql", c.Adapter)
	require.Equal(t, "root", c.Username)
	require.Equal(t, "", c.Password)
	require.Equal(t, "localhost", c.Host)
	require.Equal(t, "3306", c.Port)
	require.Equal(t, "queryx_test", c.Database)
	require.Equal(t, "root:@tcp(localhost:3306)/queryx_test?parseTime=true", c.URL)
	require.Equal(t, "root:@tcp(localhost:3306)/?parseTime=true", c.URL2)
	require.Equal(t, "mysql://root:@localhost:3306/queryx_test?parseTime=true", c.TSFormat())

}

func TestNewSQLiteConfig(t *testing.T) {
	u := "sqlite:test.sqlite3"
	c := NewConfig(&schema.Config{URL: types.StringOrEnv{Value: u}})
	require.Equal(t, "sqlite", c.Adapter)
	require.Equal(t, "file:test.sqlite3", c.URL)
	require.Equal(t, "test.sqlite3", c.Database)
}
