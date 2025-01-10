package redis

import "github.com/redis/go-redis/v9"


func InitRedisClient() *redis.Client{
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Ganti sesuai konfigurasi Redis Anda
		Password: "",               // Kosongkan jika tidak ada password
		DB:       0,                // Gunakan database Redis default (0)
	})
}