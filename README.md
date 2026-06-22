# Fleet Management System - Transjakarta Technical Test

Sistem manajemen armada (*Fleet Management*) backend yang dirancang menggunakan bahasa pemrograman **Golang**, menerapkan pola **Clean Architecture**, dan dikemas secara rapi menggunakan **Docker**.

## Fitur Utama
1. **MQTT Subscriber**: Menerima koordinat GPS armada secara *real-time* via protokol MQTT.
2. **PostgreSQL Storage**: Menyimpan jejak riwayat pergerakan kendaraan ke dalam database rasional.
3. **REST API**: Menyediakan *endpoint* untuk melihat lokasi terakhir dan riwayat perjalanan armada.
4. **Geofencing & RabbitMQ**: Otomatis mendeteksi jika kendaraan memasuki radius 50 meter dari titik tertentu (Monas), dan seketika membunyikan peringatan (*alert*) melalui RabbitMQ.

## Teknologi yang Digunakan
- **Bahasa**: Golang (1.24+)
- **Framework Web**: Gin Gonic
- **Database**: PostgreSQL 16 (diakses menggunakan pgxpool)
- **Message Broker (MQTT)**: Eclipse Mosquitto
- **Message Queue (AMQP)**: RabbitMQ 3 Management
- **Infrastruktur**: Docker & Docker Compose

## Arsitektur Sistem (All-in-One)
Aplikasi didesain untuk menghindari *over-engineering*. Meski fitur-fiturnya (API, MQTT Sub, RabbitMQ Worker) seolah-olah butuh *microservices* terpisah, sistem ini dirakit secara padu dalam **1 Aplikasi Tunggal** yang berjalan harmonis secara bersamaan di latar belakang (*background*) memanfaatkan *Goroutine* bawaan Golang.

---

## Cara Menjalankan Aplikasi

### 1. Prasyarat (*Prerequisites*)
Pastikan di komputer Anda sudah terinstal perangkat lunak berikut:
- **Docker Desktop** (atau Docker Engine + Docker Compose)
- **Golang** (Versi 1.24 atau ke atas)
- **Postman** (untuk pengujian API)

### 2. Memulai Infrastruktur (Otomatis)
Aplikasi backend, beserta database, dan *message broker* telah dikonfigurasi dalam satu file Compose.
Buka terminal di dalam folder proyek ini, lalu jalankan:

```bash
docker-compose up -d --build
```

*Sistem akan otomatis mengunduh image, melakukan kompilasi build Golang, menjalankan kontainer, serta membuat tabel database (Auto-Migration).*

### 3. Cek Status Container
Untuk memastikan semuanya berjalan normal:
```bash
docker ps
```
Anda seharusnya melihat 4 kontainer berjalan: `fleet_api`, `fleet_postgres`, `fleet_rabbitmq`, dan `fleet_mosquitto`.

Untuk melihat aktivitas (log) dari otak aplikasi kita secara *live*:
```bash
docker logs -f fleet_api
```

### 4. Menjalankan Simulasi Kendaraan (Mock Publisher)
Kami telah menyediakan skrip *bot* khusus untuk menyimulasikan bus yang sedang bergerak menuju pusat *Geofence*.
Buka tab terminal/CMD **baru**, dan jalankan:

```bash
go run cmd/publisher/main.go
```
*Perhatikan terminal `docker logs -f fleet_api` Anda! Anda akan melihat secara visual sistem menerima data MQTT, mengeksekusi Rumus Haversine, lalu memicu Alarm RabbitMQ.*

---

## Pengujian REST API

Anda bisa melakukan uji coba pemanggilan data melalui Postman atau browser:

### 1. Mengecek Lokasi Terkini Kendaraan
**GET** `http://localhost:3000/vehicles/B1234XYZ/location`

### 2. Melihat Riwayat Lengkap Kendaraan
**GET** `http://localhost:3000/vehicles/B1234XYZ/history?start=1700000000&end=2000000000`

### 3. Dashboard Kelinci (RabbitMQ UI)
Buka browser dan akses antarmuka administrator RabbitMQ:
**URL:** `http://localhost:15672`  
**Username:** `guest` | **Password:** `guest`  
*(Cek tab "Queues" untuk melihat lonjakan grafik pesan `geofence_alerts` masuk!)*
