package configs

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type envConfig struct {
	MongoDBURI  string
	MONGODBName string
	Port        string
	JWTSecretKey string
}

var envCfg envConfig

func init() {
	dir, err := os.Getwd()
	if err != nil{
		panic(err)
	}

	envPath := filepath.Join(dir, ".env")

	if err := godotenv.Load(envPath); err != nil{
		panic(err)
	}

	envCfg = envConfig{
		MongoDBURI: os.Getenv("MONGODB_URI"),
		MONGODBName: os.Getenv("MONGODB_NAME"),
		Port: os.Getenv("PORT"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

func GetConfig() *envConfig{
	return &envCfg
}