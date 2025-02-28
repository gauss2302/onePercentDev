package config

import "fmt"

type ServerConfig struct {
	Port string `mapstructure:"port" description:"Port on which the server will listen"`
	Host string `mapstructure:"host" description:"Host address for the server"`
}

func (s *ServerConfig) Validate() error {
	if s.Port == "" {
		return fmt.Errorf("server.port is required")
	}
	if s.Host == "" {
		return fmt.Errorf("server.host is required")
	}
	return nil
}

func (s *ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
