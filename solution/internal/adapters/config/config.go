package config

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"os"
	"solution/internal/adapters/logger"
)

type Config struct {
	Database *pgxpool.Pool
	Redis    *redis.Client
	Minio    *minio.Client
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("failed to read config: %v", err)
	}
}

func Configure() *Config {
	initConfig()

	logger.New(
		viper.GetBool("settings.debug"),
		viper.GetString("settings.timezone"),
	)
	logger.Log.Debugf("Debug mode: %t", viper.GetBool("settings.debug"))

	// Initialize database
	logger.Log.Info("Initializing database...")

	logger.Log.Debug("Configuring postgres connection string")
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		viper.GetString("service.database.ssl-mode"),
		viper.GetString("settings.timezone"),
	)

	logger.Log.Debug("Configuring database")
	pgxConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		logger.Log.Panicf("Failed to parse config: %v", err)
		os.Exit(1)
	}

	pgxConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	logger.Log.Debugf("dsn: %s", dsn)
	logger.Log.Debug("Connecting to postgres...")
	database, errConnect := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if errConnect != nil {
		logger.Log.Panicf("Failed to connect to postgres: %v", errConnect)
		os.Exit(1)
	} else {
		logger.Log.Info("Connected to postgres")
	}

	if err := otelpgx.RecordStats(database); err != nil {
		logger.Log.Panicf("Unable to record database stats: %v", err)
	}

	logger.Log.Info("Database initialized")

	logger.Log.Info("Initializing redis...")
	redisAddress := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})
	logger.Log.Info("Redis initialized")

	logger.Log.Info("Initializing minio...")
	endpoint := os.Getenv("MINIO_HOST")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	logger.Log.Info("Minio initialized")

	if err != nil {
		logger.Log.Panicf("Failed to connect to minio: %v", err)
		os.Exit(1)
	}
	logger.Log.Info("Minio initialized")

	return &Config{
		Database: database,
		Redis:    redisClient,
		Minio:    minioClient,
	}
}
