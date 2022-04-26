package server

type Config struct {
	Host    string `env:"RUN_ADDRESS"`
	DB      string `env:"DATABASE_URI"`
	Accrual string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Key     []byte
}

var CFG Config
