package mqtt

import (
	"context"
	"encoding/json"
	"log"

	"fleet-management/internal/config"
	"fleet-management/internal/geofence"
	"fleet-management/internal/rabbitmq"
	"fleet-management/internal/service"

	mqttclient "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber struct {
	client   mqttclient.Client
	service  service.LocationService
	config   *config.Config
	producer *rabbitmq.Producer
}

func NewSubscriber(cfg *config.Config, svc service.LocationService, prod *rabbitmq.Producer) *Subscriber {
	opts := mqttclient.NewClientOptions()
	opts.AddBroker(cfg.MQTTBroker)
	opts.SetClientID("fleet-management-subscriber")

	opts.OnConnect = func(c mqttclient.Client) {
		log.Println("MQTT Connected Successfully to", cfg.MQTTBroker)
	}
	opts.OnConnectionLost = func(c mqttclient.Client, err error) {
		log.Printf("MQTT Connection lost: %v", err)
	}

	client := mqttclient.NewClient(opts)

	return &Subscriber{
		client:   client,
		service:  svc,
		config:   cfg,
		producer: prod,
	}
}

func (s *Subscriber) Start() {
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT: %v", token.Error())
	}

	topic := "/fleet/vehicle/+/location"

	token := s.client.Subscribe(topic, 1, s.messageHandler)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	log.Printf("Subscribed to MQTT topic: %s", topic)
}

func (s *Subscriber) Stop() {
	s.client.Disconnect(250)
	log.Println("MQTT Disconnected")
}

func (s *Subscriber) messageHandler(client mqttclient.Client, msg mqttclient.Message) {
	payload := msg.Payload()
	log.Printf("[MQTT] Pesan Masuk dari topik %s: %s", msg.Topic(), string(payload))

	var data struct {
		VehicleID string  `json:"vehicle_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timestamp int64   `json:"timestamp"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		log.Printf("[MQTT] Error - Payload JSON tidak valid: %v", err)
		return
	}

	if data.VehicleID == "" || data.Timestamp == 0 {
		log.Printf("[MQTT] Error - Data tidak lengkap")
		return
	}

	// 1. Simpan ke Database
	err := s.service.Save(context.Background(), data.VehicleID, data.Latitude, data.Longitude, data.Timestamp)
	if err != nil {
		log.Printf("[MQTT] Error - Gagal menyimpan ke DB: %v", err)
		return
	}

	// 2. Geofence Check (Haversine Formula)
	distance := geofence.CalculateDistance(
		s.config.GeofenceLat,
		s.config.GeofenceLng,
		data.Latitude,
		data.Longitude,
	)

	// Jika jarak <= radius (contoh: 50 meter)
	if distance <= s.config.GeofenceRadius {
		eventPayload := map[string]interface{}{
			"vehicle_id": data.VehicleID,
			"event":      "geofence_entry",
			"location": map[string]float64{
				"latitude":  data.Latitude,
				"longitude": data.Longitude,
			},
			"timestamp": data.Timestamp,
		}

		err = s.producer.PublishGeofenceEvent(context.Background(), eventPayload)
		if err != nil {
			log.Printf("[RABBITMQ] Error publishing geofence event: %v", err)
		} else {
			log.Printf("[RABBITMQ] 🚨 Geofence event published for vehicle %s", data.VehicleID)
		}
	}
}
