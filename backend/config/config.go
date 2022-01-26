package config

type CompleteConfig struct {
	DB DB
	Logger Logger
	Web Web
}

func GetConfig() CompleteConfig {
	return CompleteConfig{
		DB:     GetDB(),
		Logger: GetLogger(),
		Web:    GetWeb(),
	}
}
