package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	RestPort   string `required:"true" envconfig:"REST_PORT"`
	GrpcPort   string `required:"true" envconfig:"GRPC_PORT"`
	MongoDBUri string `required:"true" envconfig:"MONGO_DB_URI"`
}

var Global Env = Env{}

func Load() error {
	if err := envconfig.Process("", &Global); err != nil {
		return err
	}
  return nil
}
