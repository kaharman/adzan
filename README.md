# Adzan API

Adzan API menyediakan jadwal waktu shalat di seluruh wilayah Indonesia (berdasarkan hisab Bimas Islam Kemenag RI).

Mendukung dua mode data source:
1. **CSV (Standalone / File-based):** Menggunakan file CSV per kota (`data/praytimes/{city_code}.csv`) tanpa membutuhkan Redis server.
2. **Redis (In-Memory Key-Value):** Menggunakan Redis Hash perpetual bulanan `praytime:{city_code}:{MM}`.

---

## 🚀 Menjalankan Menggunakan Docker (GitHub Packages)

Docker image otomatis di-build dan dipublikasikan ke **GitHub Container Registry (GHCR)**:

```bash
# Pull image terbaru
docker pull ghcr.io/kaharman/adzan:latest
```

### Menjalankan Mode CSV (Standalone)
Mount folder data lokal Anda ke `/data` di dalam container:

```bash
docker run -d --name adzan-service \
  -p 4321:4321 \
  -v $(pwd)/data:/data \
  -e DATA_SOURCE=csv \
  -e CSV_PATH=/data \
  ghcr.io/kaharman/adzan:latest
```

### Menjalankan Mode Redis
```bash
docker run -d --name adzan-service \
  -p 4321:4321 \
  -e DATA_SOURCE=redis \
  -e REDIS_HOST=host.docker.internal \
  -e REDIS_PORT=6379 \
  ghcr.io/kaharman/adzan:latest
```

---

## 🛠️ Build Image Sendiri Secara Lokal

```bash
docker build -t adzan:latest .
```

---

## 📡 Endpoint API Contoh

* **Jadwal Harian Kota Jakarta:**
  ```bash
  curl -s "http://localhost:4321/daily?city_code=3171"
  ```
* **Jadwal Bulanan:**
  ```bash
  curl -s "http://localhost:4321/monthly?city_code=3171&month=11&year=2026"
  ```
* **Pencarian Kota:**
  ```bash
  curl -s "http://localhost:4321/city?city=bandung"
  ```
