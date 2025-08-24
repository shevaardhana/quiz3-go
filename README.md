# Library API

Aplikasi RESTful API untuk mengelola **kategori** dan **buku** menggunakan Go, Gin, dan PostgreSQL.  
API dilengkapi dengan Basic Auth untuk autentikasi.

---

## **Fitur**

- CRUD Kategori (`/api/categories`)  
- CRUD Buku (`/api/books`)  
- Menampilkan buku berdasarkan kategori (`/api/categories/:id/books`)  
- Validasi input (misal `releaseYear`, `totalPage`)  
- Konversi otomatis `thickness` dari `totalPage`  
- Middleware Basic Auth untuk proteksi endpoint  

---

## **Instalasi**

1. Clone repository
```bash
git clone https://github.com/username/library-api.git
cd library-api

## **Create Database with Table**
CREATE DATABASE db_library;

CREATE TABLE buku (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    release_year INTEGER,
    price INTEGER,
    total_page INTEGER,
    thickness VARCHAR(100),
    category_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP,
    modified_by VARCHAR(100)
);

CREATE TABLE kategori (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP,
    modified_by VARCHAR(100)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP,
    modified_by VARCHAR(100)
);

INSERT INTO users (username, password, created_at, created_by)
VALUES ('sheva', 'sheva123', NOW(), 'system');

## Edit file database/db.go sesuai konfigurasi DB Anda (host, port, user, password, dbname).

## **Instalasi**
go mod tidy

## **menjalankan aplikasi**
go run main.go

## **path API**
| Method | Path                        | Kegunaan                              |
| ------ | --------------------------- | ------------------------------------- |
| GET    | `/api/categories`           | Menampilkan seluruh kategori          |
| POST   | `/api/categories`           | Menambahkan kategori                  |
| GET    | `/api/categories/:id`       | Menampilkan detail kategori           |
| PUT    | `/api/categories/:id`       | Update kategori                       |
| DELETE | `/api/categories/:id`       | Menghapus kategori                    |
| GET    | `/api/categories/:id/books` | Menampilkan buku berdasarkan kategori |


| Method | Path             | Kegunaan                 |
| ------ | ---------------- | ------------------------ |
| GET    | `/api/books`     | Menampilkan seluruh buku |
| POST   | `/api/books`     | Menambahkan buku         |
| GET    | `/api/books/:id` | Menampilkan detail buku  |
| PUT    | `/api/books/:id` | Update buku              |
| DELETE | `/api/books/:id` | Menghapus buku           |


