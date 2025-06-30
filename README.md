# Auth Starter Go-Fiber
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

### Struktur folder 📂
```
📦 auth-starter
├── 📁 config
│   └── config.go
├── 📁 constant
│   └── response.go
├── 📁 database
│   └── connection.go
├── 📁 handler
│   ├── auth.handler.go
│   └── user.handler.go
├── 📁 middleware
│   └── jwt.middleware.go
├── 📁 model
│   ├── profile.go
│   └── user.go
├── 📁 repository
│   ├── profile.repository.go
│   └── user.repository.go
├── 📁 router
│   └── router.go
├── 📁 service
│   ├── auth.service.go
│   └── user.service.go
├── .env
├── go.mod
├── go.sum
└── main.go

```


### API docs 🌐
untuk API dokumentasi sudah ada, foldernya bernama doc_api. untuk membuka dokumentasi gunakan tool [Bruno](https://www.usebruno.com/)
