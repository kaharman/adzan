<div align="center">
  <img src="assets/logo.svg" alt="Adzan Logo" width="128" height="128" />

  # 🕌 Adzan API

  **High-Performance Indonesian Prayer Time Engine &amp; API (REST &amp; gRPC)**  
  *Jadwal waktu shalat presisi khusus wilayah Indonesia 🇮🇩 selaras dengan data hisab resmi Bimas Islam Kementerian Agama RI.*

  [![Go Report Card](https://goreportcard.com/badge/github.com/kaharman/adzan)](https://goreportcard.com/report/github.com/kaharman/adzan)
  [![License: MIT](https://img.shields.io/badge/License-MIT-1c6b5b.svg)](LICENSE)
  [![Docker GHCR](https://img.shields.io/badge/Docker-ghcr.io%2Fkaharman%2Fadzan-2b7d6d?logo=docker&logoColor=white)](https://github.com/kaharman/adzan/pkgs/container/adzan)
  [![Web App](https://img.shields.io/badge/Web%20App-adzan.kaharman.my.id-2d81af)](https://adzan.kaharman.my.id)

  <br />

  [🌐 Web Demo](https://adzan.kaharman.my.id) • [📖 Skenario](#-skenario-penggunaan) • [💾 Data Source](#-pilihan-data-source) • [🚀 Cara Menjalankan](#-cara-menjalankan) • [📡 Endpoint API](#-dokumentasi-endpoint-rest-api)
</div>

---

**ADZAN** merupakan REST & gRPC API yang menyediakan jadwal waktu shalat presisi khusus untuk seluruh wilayah 🇮🇩 Indonesia.

Berbeda dari API waktu shalat pada umumnya yang menghitung posisi matahari secara *real-time* menggunakan rumus hisab/algoritma astronomi kompleks, **ADZAN menggunakan metode berbasis data tabel (lookup table berulang/perpetual)** yang disadur langsung dari jadwal hisab resmi **Bimas Islam Kementerian Agama Republik Indonesia**.

---

## ⚡ Mengapa Berbasis Data Tabel?

Sebagian besar API jadwal shalat mengandalkan kalkulasi trigonometri/astronomi saat request datang (*on-the-fly*). Pendekatan tersebut membutuhkan beban CPU lebih tinggi dan sering kali menghasilkan selisih 1–3 menit dari jadwal resmi yang digunakan oleh masjid-masjid di Indonesia.

Di wilayah tropis seperti Indonesia, variasi siklus waktu shalat pada tanggal dan bulan yang sama antar-tahun sangat kecil (rentang perbedaannya dapat diabaikan). Dengan menyimpan jadwal 366 hari per kota dalam bentuk data tabel:
* **Performa Tinggi & Latensi Sangat Rendah:** Query beroperasi dalam kompleksitas $O(1)$ langsung dari memory cache atau file CSV lokal tanpa kalkulasi rumit.
* **Hemat Sumber Daya:** Sangat ringan dan efisien, cocok di-deploy pada server berspesifikasi kecil, micro-instance, hingga single-board computer (SBC).
* **Konsisten:** Waktu shalat yang disajikan selaras dengan acuan Kementerian Agama RI yang berlaku di masyarakat.

---

## 💡 Skenario Penggunaan

ADZAN dirancang fleksibel untuk dijadikan penyedia data (*backend provider*) berbagai integrasi:
* 🌐 **Website & Portal:** Menampilkan widget jadwal shalat harian atau tabel bulanan di portal berita, blog, atau web instansi.
* 📱 **Aplikasi Mobile:** Sumber data backend untuk aplikasi Android / iOS, lock screen notification, dan aplikasi pengingat ibadah.
* 🖥️ **Desktop & Smartphone Widgets:** Widget layar utama smartphone atau sistem operasi desktop.
* 🕌 **Perangkat Jadwal Shalat Masjid:** Digunakan oleh display jam digital masjid, display running text LED, Smart TV masjid, hingga perangkat otomasi audio adzan berbasis mini-PC / Raspberry Pi / ESP32.

---

## 💾 Pilihan Data Source

ADZAN mendukung dua pilihan arsitektur penyimpanan data yang dapat diatur melalui konfigurasi:
1. **CSV (Standalone / File-Based):**
   - Data dibaca dari file CSV per kota (`data/praytimes/{city_code}.csv`) dan daftar kota (`data/cities.csv`).
   - Dilengkapi *in-memory cache* dan *lazy loading* otomatis.
   - **Tanpa dependensi luar:** Aplikasi dapat langsung berjalan mandiri tanpa memerlukan server database/Redis.
2. **Redis (In-Memory Key-Value):**
   - Data disimpan dalam Redis Hash berulang bulanan (`praytime:{city_code}:{MM}`).
   - Cocok untuk deployment skala besar atau arsitektur microservices terdistribusi.

---

## 🚀 Cara Menjalankan

### Opsi A: Menjalankan Langsung (Tanpa Docker / Bare Metal)

Pastikan sistem Anda telah terpasang **Go (versi 1.16+)**.

1. **Clone repositori:**
   ```bash
   git clone https://github.com/kaharman/adzan.git
   cd adzan
   ```

2. **Siapkan file konfigurasi:**
   Salin template konfigurasi bawaan:
   ```bash
   cp config.yaml.example config.yaml
   ```
   *Buka `config.yaml` dan sesuaikan pengaturan seperti `data_source` (`csv` atau `redis`), `port`, serta `server.base_path` jika memerlukan prefix URL.*

3. **Jalankan aplikasi:**
   ```bash
   # Menjalankan langsung via Go
   go run app/main.go

   # Atau compile binary terlebih dahulu
   go build -o adzan-service app/main.go
   ./adzan-service
   ```
   Aplikasi akan mulai berjalan pada port yang ditentukan (default `4321`).

---

### Opsi B: Menjalankan Menggunakan Docker (GitHub Packages)

Docker image resmi otomatis di-build multi-arsitektur (`linux/amd64` dan `linux/arm64`) dan dipublikasikan ke **GitHub Container Registry (GHCR)**.

1. **Pull image terbaru:**
   ```bash
   docker pull ghcr.io/kaharman/adzan:latest
   ```

2. **Jalankan kontainer (Mode CSV bawaan):**
   ```bash
   docker run -d --name adzan-service \
     -p 4321:4321 \
     ghcr.io/kaharman/adzan:latest
   ```

3. **Jalankan kontainer dengan volume mount data kustom:**
   ```bash
   docker run -d --name adzan-service \
     -p 4321:4321 \
     -v $(pwd)/data:/data \
     -e DATA_SOURCE=csv \
     -e CSV_PATH=/data \
     ghcr.io/kaharman/adzan:latest
   ```

4. **Jalankan kontainer dengan data source Redis:**
   ```bash
   docker run -d --name adzan-service \
     -p 4321:4321 \
     -e DATA_SOURCE=redis \
     -e REDIS_HOST=host.docker.internal \
     -e REDIS_PORT=6379 \
     ghcr.io/kaharman/adzan:latest
   ```

---

## 📖 Dokumentasi & Contoh API

Dokumentasi OpenAPI / Swagger UI interaktif tersedia langsung di:
```text
http://localhost:4321/openapi/
```
*(Atau `http://localhost:4321{base_path}/openapi/` jika menggunakan `server.base_path`)*

### Contoh Endpoint HTTP:

* **Jadwal Hari Ini (Kota Jakarta):**
  ```bash
  curl -s "http://localhost:4321/today?city_code=3171"
  ```

* **Jadwal Tanggal Spesifik:**
  ```bash
  curl -s "http://localhost:4321/daily?city_code=3171&date=2026-11-02"
  ```

* **Jadwal 1 Bulan Penuh:**
  ```bash
  curl -s "http://localhost:4321/monthly?city_code=3171&month=11&year=2026"
  ```

* **Jadwal 1 Tahun:**
  ```bash
  curl -s "http://localhost:4321/yearly?city_code=3171&year=2026"
  ```

* **Pencarian Kota:**
  ```bash
  curl -s "http://localhost:4321/city/search?city=bandung"
  ```

* **Detail Kota berdasarkan Kode / Koordinat:**
  ```bash
  curl -s "http://localhost:4321/city?city_code=3171"
  curl -s "http://localhost:4321/city?latitude=-6.170089&longitude=106.831050"
  ```

---

## 🌐 Web Interface

Tersedia antarmuka web modern dengan jam analog presisi, hitung mundur menuju waktu shalat berikutnya, kalender Hijriah, pratinjau jadwal imsakiyah, dan live core API explorer:
* **Live Web App**: [https://adzan.kaharman.my.id](https://adzan.kaharman.my.id)

---

## ⚠️ Disclaimer

Data waktu shalat yang disajikan dalam API ini dihimpun dan disesuaikan langsung dari data publik **Bimas Islam Kementerian Agama Republik Indonesia**. Penyusun/pengembang berusaha menjaga keakuratan pemrosesan data, namun tidak bertanggung jawab atas potensi ketidaksesuaian atau kesalahan data yang bersumber langsung dari data primer penyedia aslinya. Pengguna disarankan untuk tetap melakukan kroscek dengan otoritas setempat bila diperlukan.

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE).

