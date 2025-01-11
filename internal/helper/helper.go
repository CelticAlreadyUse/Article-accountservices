package helper

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

func Hashpassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func GenerateToken(userID int64) (strToken string, err error) {
	duration, err := time.ParseDuration(config.JWTExp().String())
	if err != nil {
		return "", err
	}
	claims := &model.CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	strToken, err = token.SignedString([]byte(config.JWTSigningKey()))
	if err != nil {
		return "", err
	}
	return
}
func DecodeToken(token string, claim *model.CustomClaims) (err error) {
	jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSigningKey()), nil
	})
	return
}
func ValidateToken(tokenString string, config model.ConfigJWT) (*model.CustomClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.SigningKey), nil
    })
    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %v", err)
    }
    if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, fmt.Errorf("invalid token")
}
func IsTokenExpired(expTime *jwt.NumericDate) bool {
    if expTime == nil {
        return true
    }
    return time.Now().After(expTime.Time)
}
func GenerateOTP()string{
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // Random 6 digit OTP

}
func CheckDataSame(data,send string)error{
	if data == send{
		return nil
	}else{
		return errors.New("the data is not the same")
	}
}

func SendEmail(to, subject, body string)error{
	mail := gomail.NewMessage()
	mail.SetHeader("From", "wahyusantosokanisius@gmial.com") // Ganti dengan email Anda
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)

	// Konfigurasi SMTP
	port, _ := strconv.Atoi(config.SMTPPort()) 
	dialer := gomail.NewDialer(config.SMTPHost(), port, config.SMTPEmail(), config.SMTPPasswrod())

	// Kirim email
	if err := dialer.DialAndSend(mail); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}