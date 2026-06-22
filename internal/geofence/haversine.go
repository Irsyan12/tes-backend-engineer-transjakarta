package geofence

import "math"

// earthRadiusMeters adalah konstanta jari-jari bumi.
// Karena bumi itu bulat, kita butuh ini agar perhitungannya akurat di peta melengkung.
const earthRadiusMeters = 6371000

// CalculateDistance adalah otak dari Geofence kita (Pagar Gaib).
// Fungsi ini menggunakan "Rumus Haversine" untuk menghitung jarak lurus (garis burung)
// antara koordinat mobil (lat1, lon1) dengan pusat geofence (lat2, lon2).
// Hasil yang dikembalikan adalah dalam satuan Meter.
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Mengubah derajat GPS (Latitude/Longitude) ke format Radian untuk perhitungan matematika
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	// Ini murni penerapan rumus matematika Haversine
	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Pow(math.Sin(dLon/2), 2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Mengkalikan hasil sudut dengan jari-jari bumi untuk mendapatkan hasil akhir dalam Meter
	return earthRadiusMeters * c
}
