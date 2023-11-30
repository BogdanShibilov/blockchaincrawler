package config

type Config struct {
	Http      Http      `yaml:"http"`
	Transport Transport `yaml:"transport"`
}

type Http struct {
	Port string `yaml:"port"`
}

type Transport struct {
	BlockInfoTransport BlockInfoTransport `yaml:"blockInfoTransport"`
	AuthTransport      AuthTransport      `yaml:"authTransport"`
}

type BlockInfoTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AuthTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
