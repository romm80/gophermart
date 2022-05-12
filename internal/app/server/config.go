package server

type Config struct {
	Host           string `env:"RUN_ADDRESS"`
	DB             string `env:"DATABASE_URI"`
	Accrual        string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Key            []byte
	WorkerPoolSize int
}

var CFG Config
