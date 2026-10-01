package config

type Config struct {
	Name         string `yaml:"name" default:"e2engine.sqlite" validate:"min(3)"`
	Driver       string `yaml:"driver" default:"sqlite" validate:"oneof(sqlite)"`
	MaxOpenConns int    `yaml:"max_open_conns" default:"1" validate:"min(1)"`
	MaxIdleConns int    `yaml:"max_idle_conns" default:"1" validate:"min(1)"`
}
