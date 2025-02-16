package config

import (
	"flag"
	"os"
)

var (
	DatabaseDSN string
)

func ParseFlags() {
	flag.StringVar(&DatabaseDSN, "d", "postgres://postgres:password@db:5432/shop", "postgres connection url")
	flag.Parse()

	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	if envDatabaseDSN != "" {
		DatabaseDSN = envDatabaseDSN
	}
}
