package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
func ENV() string {
	return viper.GetString("env");
}
func PORT() string {
	return viper.GetString("port");
}
func MySQLDBHost() string {
	return viper.GetString("mysql.dbhost")
}
func MySQLDBPort() string {
	return viper.GetString("mysql.dbport");
}
func MySQLDBUser() string {
	return viper.GetString("mysql.dbuser")
}
func MySQLDBPass() string {
	return viper.GetString("mysql.dbpass")
}
func MySQLDBName() string {
	return viper.GetString("mysql.dbname")
}
func JWTSigningKey() string {
	return viper.GetString("jwt.signing_key")
}
func JWTExp() time.Duration {
	return viper.GetDuration("jwt.exp")
}
func EMAIExp()time.Time{
	return viper.GetTime("email_verification.expTime")
}
func SMTPHost()string{
	return viper.GetString("smtp.host")
}
func SMTPPort()string{
	return viper.GetString("smtp.port")
}
func SMTPName()string{
	return viper.GetString("smtp.name")
}
func SMTPEmail()string{
	return viper.GetString("smtp.email")
}
func SMTPPasswrod()string{
	return viper.GetString("smtp.password")
}
func RedisHost()string{
	return viper.GetString("redis.host")
}
func RedisPass()string{
	return viper.GetString("redis.pass")
}
func RedisDB()string{
	return viper.GetString("redis.db")
}