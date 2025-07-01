# Auth Starter Go-Fiber
![aut_starter_small](https://github.com/user-attachments/assets/ed9fda79-60e3-49af-b9d5-113566a29380)

Ini adalah starter auth backend menggunakan [Golang](https://go.dev/), [Fiber](https://gofiber.io/) dan [Gorm](https://gorm.io/)(postgre).
repo ini tidak sengaja terbuat karena sebelumnya akan dipakai untuk server PPDB sekolah.
Fitur yang tersedia didalam repo ini masih cukup basic, bisa kalian custom sendiri sesuai dengan kebutuhan atau bisa kalian jadikan untuk materi pembelajaran.
Starter ini sedikit mengadopsi penggunaan backend javascript seperti [Express JS](https://expressjs.com/) dan [Nest JS](https://nestjs.com/) dengan harapan akan mempermudah kalian jika memiliki background javascript/typescript

### Fitur 📑

- Form validation menggunakan [package validator](https://github.com/go-playground/validator).
- [Password hashing](https://pkg.go.dev/golang.org/x/crypto/bcrypt).
- Register.
  - name
  - email
  - phone_number
  - password
- Login.
  - email
  - password
- Aktivation kode setelah register.
  - setelah register akan generate token yang berisi kode aktivasi yang nantinya akan dikirim lewat email. Untuk repo ini masih belum diintegrasikan dengan service pihak ketiga
- Forgot password.
- Reset password.
- Refresh token
- Guard (token verify).
  - middleware untuk memverifikasi token ketika client mengakses resourse tertentu

### Base struktur folder 📂

```
📦 auth-starter
├── 📁 config
│   └── config.go
├── 📁 controller
│   └── auth.controller.go
│   └── profile.controller.go
├── 📁 doc_api
│   └── 📁 auth
│   │    └── Activation.bru
│   │    └── folder.bru
│   │    └── health.bru
│   │    └── Login.bru
│   │    └── refresh token.bru
│   │    └── register.bru
│   └── 📁 Profile
│   │    └── folder.bru
│   │    └── profile.bru
│   └── bruno.json
├── 📁 internal
│   └── 📁 middleware
│   │    └── guard.go
│   └── 📁 model
│   │    └── user.go
│   │    └── profile.go
│   └── db.go
│   └── key_generator.go
│   └── validation.go
├── 📁 server
│   └── router.go
├── 📁 service
│   └── 📁 auth
│   │    └── activation.service.go
│   │    └── forgot_password.service.go
│   │    └── handler.service.go
│   │    └── login.service.go
│   │    └── refresh_token.service.go
│   │    └── register.service.go
│   │    └── reset_password.service.go
│   └── health.service.go
│   └── profile.service.go
├── THIRD_PARTY_LICENSES/
│   └── golang-x-crypto.LICENSE
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── license
└── main.go
└── README.MD

```

### API docs 🌐

untuk API dokumentasi sudah ada, foldernya bernama doc_api. untuk tool dokumentasi saya menggunakan [Bruno](https://www.usebruno.com/)

### Cara install 💻

- clone repository

```
git clone https://github.com/nurhamsah1998/auth-starter.git
```

- install

```
go mod tidy
```

- jalankan server

```
go run main.go
```

- Rekomendasi pakai [Air](https://github.com/air-verse/air) seperti nodemon di javascript
- file .env

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=root
DB_NAME=authDB
DB_SSLMODE=disable
ACCESS_TOKEN=SJHuyh76HYuj43Derf09MNbv8Jnhu7
REFRESH_TOKEN=Nhu76BGhjd8uji8HNB56tgtyh
ACTIVATION_TOKEN=POlk8Iu8jNh7yhG32ZsdxRTmjhN
RESET_PASSWORD_TOKEN=jhsTGWmdk8Yhe6Cfdr5cDfe
```

## Lisensi Pihak Ketiga 📚

Proyek ini menggunakan pustaka pihak ketiga yang dilisensikan di bawah MIT, BSD, dan Apache 2.0.
Semua hak cipta tetap milik pemilik aslinya.

- Daftar pustaka dan jenis lisensi: [docs/LICENSE_REPORT.txt](./docs/LICENSE_REPORT.txt)
- Teks lisensi lengkap: [docs/licenses/](./docs/licenses/)
