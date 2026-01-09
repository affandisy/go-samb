## Gudang SAMB

Manajemen Gudang dengan menggunakan Golang + React

## Fitur Utama
- **Master Data**: Manajemen Supplier, Customer, Product, dan Warehouse.
- **Transaksi Masuk (Trx In)**: Menambah stok gudang.
- **Transaksi Keluar (Trx Out)**: Mengurangi stok dengan **Validasi Stok Otomatis** (Mencegah input jika stok kurang).
- **Laporan Stok**: Melihat posisi stok terkini (Hanya menampilkan barang yang memiliki stok).

## Teknologi yang Digunakan
### Backend
- **Language**: Go (Golang) v1.25+
- **Framework**: Echo v4
- **Database Driver**: lib/pq (PostgreSQL Driver)
- **Architecture**: Layered Architecture (Handler -> Service -> Repository)
- **Library**: Godotenv (Config)

### Frontend
- **Framework**: React (via Vite)
- **UI Library**: Bootstrap 5 (React-Bootstrap)
- **HTTP Client**: Axios
- **Routing**: React Router Dom

### Database
- **PostgreSQL**

---

## Prasyarat
Pastikan sudah terinstall:
1. **Go** (v1.24+)
2. **Node.js** (v18+) & **NPM**
3. **PostgreSQL**

---

## ⚙️ Cara Menjalankan (Local Development)

### 1. Setup Database
Buat database baru di PostgreSQL (misal: `samb`), lalu jalankan query SQL yang ada di file `migrations/00001_schema.up.sql` untuk membuat tabel dan data awal.

### 2. Setup Backend
```bash
cd backend

# Buat file .env sesuai konfigurasi database Anda
cp .env.example .env

# Edit file .env, sesuaikan DB_HOST, DB_USER, DB_PASSWORD, dll.
# Pastikan SERVER_PORT=8080

# Download dependencies
go mod tidy

# Jalankan Server
go run cmd/main.go
```

### 3. Setup Frontend
cd frontend

#### Install dependencies
npm install

#### Jalankan Frontend
npm run dev

## API Reference (cURL Testing)
### Master Data
```bash
# Get All Warehouses
curl -X GET http://localhost:8080/api/warehouses

# Get All Products
curl -X GET http://localhost:8080/api/products

# Get All Suppliers
curl -X GET http://localhost:8080/api/suppliers

# Get All Customers
curl -X GET http://localhost:8080/api/customers
```
### Transaksi Barang Masuk
```bash
curl -X POST http://localhost:8080/api/trx-in \
  -H "Content-Type: application/json" \
  -d '{
    "trx_in_no": "IN-2026-001",
    "whs_idf": 1,
    "trx_in_date": "2026-01-10",
    "trx_in_supp_idf": 1,
    "trx_in_notes": "Penerimaan awal tahun",
    "details": [
      {
        "trx_in_d_product_idf": 1,
        "trx_in_d_qty_dus": 50,
        "trx_in_d_qty_pcs": 12
      },
      {
        "trx_in_d_product_idf": 2,
        "trx_in_d_qty_dus": 100,
        "trx_in_d_qty_pcs": 0
      }
    ]
  }'
```
### Transaksi Barang Keluar
```bash
curl -X POST http://localhost:8080/api/trx-out \
  -H "Content-Type: application/json" \
  -d '{
    "trx_out_no": "OUT-2026-001",
    "whs_idf": 1,
    "trx_out_date": "2026-01-11",
    "trx_out_cust_idf": 1,
    "trx_out_notes": "Pengiriman reguler ke Customer A",
    "details": [
      {
        "trx_out_d_product_idf": 1,
        "trx_out_d_qty_dus": 5,
        "trx_out_d_qty_pcs": 0
      }
    ]
  }'
```

### Laporan Stock
```bash
curl -X GET http://localhost:8080/api/stock-report
```

### Postman Test

[![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)](https://documenter.getpostman.com/view/49261332/2sBXVfhqNF)