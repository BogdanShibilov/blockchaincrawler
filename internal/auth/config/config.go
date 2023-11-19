package config

type Config struct {
	Server    Server    `yaml:"server"`
	Auth      Auth      `yaml:"auth"`
	Transport Transport `yaml:"transport"`
}

type Server struct {
	Port string `yaml:"port"`
}

// AccesTokenDuration and RefreshTokenDuration is in minutes
type Auth struct {
	JwtSecretKey         string `yaml:"jwtSecretKey"`
	AccessTokenDuration  int    `yaml:"accessTokenDuration"`
	RefreshTokenDuration int    `yaml:"refreshTokenDuration"`
}

type Transport struct {
	UserTransport UserTransport `yaml:"userTransport"`
}

type UserTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
