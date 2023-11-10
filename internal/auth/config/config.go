package config

type Config struct {
	GrpcServer GrpcServer `yaml:"grpcServer"`
	Auth       Auth       `yaml:"auth"`
	Transport  Transport  `yaml:"transport"`
}

type GrpcServer struct {
	Port string `yaml:"port"`
}

type Auth struct {
	JwtSecretKey string `yaml:"jwtSecretKey"`
}

type Transport struct {
	UserGrpc UserGrpcTransport `yaml:"userGrpc"`
}

type UserGrpcTransport struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
