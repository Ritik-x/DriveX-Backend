package config

import "os"

type Config struct {
		AppEnv string
	Port   string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret string


	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSS3Bucket        string
}

func Load() Config {
	return Config{
AppEnv : os.Getenv("APP_ENV"),
Port:   os.Getenv("PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		JWTSecret: os.Getenv("JWT_SECRET"),


		AWSRegion:          os.Getenv("AWS_REGION"),
AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
AWSS3Bucket:        os.Getenv("AWS_S3_BUCKET"),
	}

}


