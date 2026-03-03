# Silsilah — Platform Manajemen Keluarga & Silsilah Mandiri

**Silsilah** adalah platform modern yang dirancang untuk membantu Anda mencatat, merapikan, dan memvisualisasikan garis keturunan keluarga lintas generasi. Dengan antarmuka yang premium, aman, dan kolaboratif, Silsilah memastikan warisan digital keluarga Anda tetap terjaga.

## ✨ Fitur Utama (Fase Aktif)

- **Dashboard Keluarga**: Kelola banyak keluarga dalam satu akun.
- **SSO BangLipai Portal**: Otentikasi aman menggunakan Logto.
- **Premium UI/UX**: Antarmuka berbasis SvelteKit + Tailwind 4 dengan dukungan Dark/Light mode otomatis.
- **Manajemen Anggota**: CRUD anggota keluarga, relasi, dan pengelolaan role (Owner, Admin, Editor, Viewer).
- **Visualisasi Pohon**: (Dalam Pengembangan) Visualisasi diagram silsilah interaktif.

## 🛠️ Tech Stack

### Backend (BFF)

- **Language**: Go 1.23+
- **Framework**: [Fiber v2](https://gofiber.io/)
- **Database**: PostgreSQL with [GORM](https://gorm.io/)
- **Auth**: [Logto SSO](https://logto.io/) Integration
- **Dependency Injection**: [Uber FX](https://github.com/uber-go/fx)

### Frontend

- **Framework**: [SvelteKit](https://kit.svelte.dev/) (Svelte 5)
- **Styling**: [Tailwind CSS v4](https://tailwindcss.com/)
- **Components**: [shadcn-svelte](https://shadcn-svelte.com/)
- **Icons**: [Lucide Svelte](https://lucide.dev/)

## 🚀 Persiapan Pengembangan

### 1. Prasyarat

- Go 1.23+
- Node.js 20+ & pnpm
- PostgreSQL
- Logto Instance (Endpoint & Credentials)

### 2. Backend Setup

1. Masuk ke folder backend: `cd backend`
2. Salin template profil: `cp config.example.json config.json` (Sesuaikan isinya)
3. Jalankan migrasi dan server:
   ```bash
   go run main.go
   ```

### 3. Frontend Setup

1. Masuk ke folder frontend: `cd frontend`
2. Install dependensi: `pnpm install`
3. Buat file `.env`:
   ```env
   PUBLIC_BACKEND_URL=http://localhost:8080
   INTERNAL_BACKEND_URL=http://localhost:8080
   ```
4. Jalankan mode development:
   ```bash
   pnpm dev
   ```

## 📂 Struktur Proyek

```text
.
├── backend/            # Go Fiber Application (BFF & API)
│   ├── app/            # Source code (Controller, Service, Domain)
│   ├── cmd/            # Entry point
│   └── utils/          # Helpers & Config
├── frontend/           # SvelteKit Application
│   ├── src/
│   │   ├── lib/        # Shared components & API client
│   │   └── routes/     # SvelteKit Pages & Layout Groups
│   └── tailwind.config.js
└── README.md
```

## 🔒 Keamanan

Silsilah menggunakan pola **Backend-for-Frontend (BFF)**. Frontend tidak menyimpan token akses secara langsung; semua sesi dikelola melalui cookie server-side yang aman (`httpOnly`, `secure`) yang divalidasi oleh backend Go terhadap Logto.

---

Dikembangkan dengan ❤️ oleh **Kukuh**.
