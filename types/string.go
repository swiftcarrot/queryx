package types

type String struct{}

type StringOrEnv struct {
	Value   string
	EnvKey  string
	Default string
}
