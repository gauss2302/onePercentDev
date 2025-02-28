package config

import "fmt"

type PostgresConfig struct {
	URL string `mapstructure:"url" description:"Database connection URL"`
}

func (p *PostgresConfig) Validate() error {
	if p.URL == "" {
		return fmt.Errorf("postgres.url is required")
	}
	return nil
}

func (p *PostgresConfig) GetDSN() string {
	return p.URL
}
