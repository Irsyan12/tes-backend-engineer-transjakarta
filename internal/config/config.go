package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	MQTTBroker     string
	RabbitMQURL    string
	GeofenceLat    float64
	GeofenceLng    float64
	GeofenceRadius float64
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	rmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	lat, _ := strconv.ParseFloat(os.Getenv("GEOFENCE_LAT"), 64)
	lng, _ := strconv.ParseFloat(os.Getenv("GEOFENCE_LNG"), 64)
	rad, _ := strconv.ParseFloat(os.Getenv("GEOFENCE_RADIUS"), 64)

	return &Config{
		AppPort:        port,
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		MQTTBroker:     os.Getenv("MQTT_BROKER"),
		RabbitMQURL:    rmqURL,
		GeofenceLat:    lat,
		GeofenceLng:    lng,
		GeofenceRadius: rad,
	}
}