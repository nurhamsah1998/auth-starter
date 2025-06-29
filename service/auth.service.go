package service

import (
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurhamsah1998/auth-starter/internal"
	"github.com/nurhamsah1998/auth-starter/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthService struct{}
	/// form validasi body payload yang dikirim oleh client
	FormRegister struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"required,numeric,min=10,max=15"`
		Password    string `json:"password" validate:"required,min=8,max=100"`
	}
	FormLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=100"`
	}
	FormActivation struct {
		Activation string `json:"activation" validate:"required"`
	}
	FormResetPassword struct {
		NewPassword    string `json:"new_password" validate:"required,min=8,max=100"`
		ReTypePassword string `json:"retype_password" validate:"required,min=8,max=100"`
	}
	FormForgotPassword struct {
		Email string `json:"email" validate:"required,email"`
	}
)

// / service handler untuk menginject servis ke controller
func AuthHandler() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormRegister{}
	/// validasi format json
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}
	/// validasi body payload yang dikirim oleh client (frontend)
	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"message": errors,
			"error":   true,
		})
	}

	user.Email = bodyPayload.Email
	/// hashing password
	pwdHash, errPwdH := bcrypt.GenerateFromPassword([]byte(bodyPayload.Password), 10)
	if errPwdH != nil {
		return errors.New("failed hashing password")
	}
	user.Password = string(pwdHash)
	user.Profile.Name = bodyPayload.Name
	user.Profile.PhoneNumber = bodyPayload.PhoneNumber
	/// proses insert ke database
	res := internal.DB.Create(&user)
	/// jika terjadi error ketika proses insert
	/// contoh : error unique username
	if res.RowsAffected == 0 {
		return errors.New(res.Error.Error())
	}

	/// membuat token JWT untuk aktivasi akun
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":              user.ID,
		"email":           user.Email,
		"code_activation": internal.KeyGenerate(10),
		/// token aktivasi akan expired/kadaluarsa dalam 168 jam kedepan setelah berhasil register
		"exp": time.Now().Add(168 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("ACTIVATION_TOKEN"))) /// <--- secret key untuk token aktivasi (mengambil dari file .env)
	/// jika terjadi error ketika proses pembuatan token
	if err != nil {
		return errors.New("failed create token activation")
	}
	/// proses insert token aktivasi ke tabel users pada bagian kolom "Activation",
	/// dengan menggunakan metode update.
	internal.DB.Model(&user).Update("Activation", tokenString)

	return c.Status(201).JSON(fiber.Map{"message": "Register successfully. We send a unique code for activation", "data": bodyPayload})
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormLogin{}
	/// validasi format json
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}
	/// validasi body payload yang dikirim oleh client (frontend)
	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"message": errors,
			"error":   true,
		})
	}
	/// proses pencarian data by email
	res := internal.DB.Preload("Profile").Find(&user, "email = ?", bodyPayload.Email)
	/// jika pencarian data by email di rables users tidak ditemukan
	if res.RowsAffected == 0 {
		return errors.New("invalid credential")
	}
	/// melakukan compare/perbandingan antara password yang dikirim client
	/// dan yang ada di database
	errPwd := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyPayload.Password))
	/// error ketika password tidak sama
	if errPwd != nil {
		return errors.New("invalid credential")
	}
	/// error ketika client yang sudah daftar,
	/// tapi belum melakukan activasi mencoba untuk login.
	if user.Activation != "" {
		return errors.New("invalid credential")
	}
	/// proses pembuatan akses token untuk client
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		/// token JWT akan expired/kadaluarsa dalam 168 jam kedepan setelah berhasil login
		"exp": time.Now().Add(168 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN"))) /// <--- secret key untuk token login (mengambil dari file .env)
	data := fiber.Map{
		"token":   tokenString,
		"id":      user.ID,
		"email":   user.Email,
		"profile": user.Profile,
	}
	return c.Status(200).JSON(fiber.Map{"message": "Login successfully", "data": data})
}

func (s *AuthService) Activation(c *fiber.Ctx) error {
	user := model.User{}
	activation := c.Params("token_activation")
	bodyPayload := FormActivation{}
	/// validasi format json
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}
	/// validasi body payload yang dikirim oleh client (frontend)
	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"message": errors,
			"error":   true,
		})
	}

	token, err := jwt.Parse(activation, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("ACTIVATION_TOKEN")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return errors.New(err.Error())
	}
	/// proses validasi user dan kode aktivasi
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id := claims["id"]
		code_activation := claims["code_activation"]
		res := internal.DB.First(&user, "id = ?", id)
		/// respon error ketika user/client tidak ditemukan
		if res.RowsAffected == 0 {
			return errors.New("user not found")
		}
		/// respon error ketika kode aktivasi yang dikirim client,
		/// tidak sama dengan kode aktivasi yang di database
		if code_activation != bodyPayload.Activation {
			return errors.New("invalid code activation")
		}
		/// respon error ketika client mencoba untuk aktivasi lagi setelah berhasil
		if user.Activation == "" {
			return errors.New("your account is already activated")

		}
	} else {
		return errors.New(err.Error())
	}
	/// proses update kolom "Activation" pada tabel users,
	/// set ke nil, menandakan client sudah melakukan aktivasi
	internal.DB.Model(&user).Update("Activation", nil)
	return c.Status(200).JSON(fiber.Map{"message": "Successfully activated account", "error": false})
}

func (s *AuthService) ResetPassword(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormResetPassword{}
	/// mengambil token dari param url
	resetPwdToken := c.Params("reset_pwd_token")
	/// validasi format json
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}

	/// proses validasi client
	token, err := jwt.Parse(resetPwdToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("RESET_PASSWORD_TOKEN")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return errors.New(err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email := claims["email"]
		res := internal.DB.First(&user, "email = ?", email)
		/// respon error ketika user/client tidak ditemukan
		if res.RowsAffected == 0 {
			return errors.New("cannot reset password")
		}
	} else {
		return errors.New(err.Error())
	}
	/// validasi body payload yang dikirim oleh client (frontend)
	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"message": errors,
			"error":   true,
		})
	}
	/// respon error jika password baru dan retype password tidak sama
	if bodyPayload.NewPassword != bodyPayload.ReTypePassword {
		return errors.New("password not match")
	}

	/// hashing password baru
	pwdHash, errPwdH := bcrypt.GenerateFromPassword([]byte(bodyPayload.NewPassword), 10)
	if errPwdH != nil {
		return errors.New("failed hashing password")
	}
	/// proses update kolom "Password" pada tabel users berdasarkan,
	/// email yang dikirim user ketika forgot password
	internal.DB.Model(&user).Update("Password", string(pwdHash))
	return c.Status(200).JSON(fiber.Map{"message": "Reset password successfully", "error": false})
}

func (s *AuthService) ForgotPassword(c *fiber.Ctx) error {
	user := model.User{}
	bodyPayload := FormForgotPassword{}
	/// validasi format json
	if err := c.BodyParser(&bodyPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body", "error": true})
	}
	/// validasi body payload yang dikirim oleh client (frontend)
	if err := internal.ClassValidate.Struct(bodyPayload); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		return c.Status(422).JSON(fiber.Map{
			"message": errors,
			"error":   true,
		})
	}
	/// pencarian client/user berdasarkan email yang dikirim client/user
	res := internal.DB.First(&user, "email = ?", bodyPayload.Email)
	/// respon error ketika user/client tidak ditemukan
	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}

	/// proses pembuatan akses token untuk reset password
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		/// token reset password akan expired/kadaluarsa dalam 168 jam kedepan,
		//  setelah berhasil mengirimkan email
		"exp": time.Now().Add(168 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("RESET_PASSWORD_TOKEN"))) /// <--- secret key untuk reset password (mengambil dari file .env)
	///
	/// pada block dibawah ini bisa diimplementasikan dengan service pihak ketiga,
	/// untuk mengirimkan link reset password.
	//
	//
	// {masukan kodemu disini}
	//
	//
	/// contoh link yang akan diakses oleh user/client dari sisi frontend,
	/// ketika mendapatkan link reset password dari service pihak ketiga tadi : https://domainmu.com/auth/reset-password/{token-reset-password}.
	/// pada tampilan/UI kurang lebih bisa berisi form new password dan retype new password,
	/// ketika user/client click tombol reset password, maka akan hit api ke https://domainmu.com/api/auth/reset-password/{token-reset-password},
	/// untuk proses reset passwordnya
	///
	/// NOTE : Link diatas hanya contoh, bisa dicustom sendiri sesuai kebutuhan
	///
	return c.Status(200).JSON(fiber.Map{"message": "Successfully sending link to your email for reset password", "error": false, "data": tokenString})
}
