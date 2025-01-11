package repository

import (
	"context"
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/redis/go-redis/v9"
)
type otpRepository struct {
	redisClient *redis.Client
	ctx         context.Context
}
func NewOTPRepository(redisClient *redis.Client,ctx context.Context) model.OTPRepository {
	return &otpRepository{
		redisClient: redisClient,
		ctx:         ctx,
	}
}
func (r *otpRepository) StoreOTP(email string, otp string, ttl time.Duration) error {
	key := "otp:" + email
	return r.redisClient.Set(r.ctx, key, otp, ttl).Err()
}
func (r *otpRepository) ValidateOTP(data model.OTPRequestValidate) (bool, error) {
	key := "otp:" + data.Email
	storedOTP, err := r.redisClient.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return false, nil 
	} else if err != nil {
		return false, err
	}
	return storedOTP == data.OTPCode, nil
}
