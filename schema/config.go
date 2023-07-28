package schema

import (
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
