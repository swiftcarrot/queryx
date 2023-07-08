package schema

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/swiftcarrot/queryx/types"
)

func TestMySQLConnectionURL(t *testing.T) {
	dsn := "mysql:mysql@tcp(localhost:3306)/queryx_test?param=value"
	config := &Config{Adapter: "mysql", URL: types.StringOrEnv{Value: dsn}}

	require.Equal(t, "mysql:mysql@tcp(localhost:3306)/", config.ConnectionURL(false))
	require.Equal(t, dsn, config.ConnectionURL(true))
	require.Equal(t, "queryx_test", config.Database)
}
