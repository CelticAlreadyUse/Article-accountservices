package helper

import (
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


func Hashpassword(password string)(string,error){
	bytes,err := bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes),err
}
func CheckPasswword(password, hash string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}

func GenerateToken(userID int64)(strToken string,err error){
	expiredAt :=time.Now().UTC().Add(config.JWTExp())
	strToken,err = jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"exp":expiredAt.Unix(),
		"user_id" :userID,
	}).SignedString([]byte(config.JWTSigningKey()))
		return
}

func DecodeToken(token string,claim *model.CustomClaims)(err error){
	jwt.ParseWithClaims(token,claim,func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSigningKey()),nil
	})
	return
}
