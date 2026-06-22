package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	mqttclient "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	opts := mqttclient.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("fleet-management-mock-publisher")

	client := mqttclient.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT: %v", token.Error())
	}
	defer client.Disconnect(250)

	log.Println("Mock Publisher connected to", broker)

	vehicleID := "B1234XYZ"
	topic := "/fleet/vehicle/" + vehicleID + "/location"

	// SIMULASI PERGERAKAN MOBIL
	// Kita set titik awal mobil (Latitude) berada sedikit di luar area Geofence
	// Pusat Geofence kita ada di: -6.2088 (seperti di .env)
	lat := -6.2095 // Lebih kecil dari -6.2088, berarti di luarnya
	lng := 106.8456

	for {
		// 1. Bungkus data GPS saat ini ke dalam JSON
		payload := map[string]interface{}{
			"vehicle_id": vehicleID,
			"latitude":   lat,
			"longitude":  lng,
			"timestamp":  time.Now().Unix(),
		}

		body, _ := json.Marshal(payload)

		// 2. Tembakkan (Publish) data tersebut ke Menara Radio (MQTT)
		token := client.Publish(topic, 1, false, body)
		token.Wait()

		log.Printf("Mengirim lokasi GPS ke %s: %s", topic, string(body))

		// 3. Gerakkan mobil sedikit maju mendekati pusat Geofence
		lat += 0.0001
		
		// Jika mobil sudah menabrak pusat geofence (-6.2088),
		// Kembalikan posisi mobil ke titik awal (luar geofence) untuk mengulang simulasi
		if lat > -6.2088 {
			lat = -6.2095 
		}

		// 4. Tunggu 2 detik sebelum mengirim lokasi GPS berikutnya
		time.Sleep(2 * time.Second)
	}
}
