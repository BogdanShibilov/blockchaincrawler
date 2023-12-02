package config

type Config struct {
	Http      Http      `yaml:"http"`
	Transport Transport `yaml:"transport"`
	Jwt       Jwt       `yaml:"jwt"`
}

type Http struct {
	Port string `yaml:"port"`
}

type Transport struct {
	BlockInfoTransport BlockInfoTransport `yaml:"blockInfoTransport"`
	AuthTransport      AuthTransport      `yaml:"authTransport"`
	UserTransport      UserTransport      `yaml:"userTransport"`
}

type BlockInfoTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AuthTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type UserTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Jwt struct {
	Secret string `yaml:"secret"`
}
