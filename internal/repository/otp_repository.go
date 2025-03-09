package repository
import (
	"context"
	"time"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)
type otpRepository struct {
	redisClient *redis.Client
	ctx         context.Context
}
func NewOTPRepository(redisClient *redis.Client, ctx context.Context) model.OTPRepository {
	return &otpRepository{
		redisClient: redisClient,
		ctx:         ctx,
	}
}
func (r *otpRepository) StoreOTP(email string, otp string, ttl time.Duration) error {
	key := "otp_emailverified:" + email
	return r.redisClient.Set(r.ctx, key, otp, ttl).Err()
}
func (r *otpRepository) ValidateOTP(data model.OTPRequestValidate) (bool, error) {
	key := "otp_emailverified:" + data.Email
	storedOTP, err := r.redisClient.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return storedOTP == data.OTPCode, nil
}
func (r *otpRepository) StoredOTPPass(email string, otp string, ttl time.Duration) error {
	key := "otp_pass:" + email
	return r.redisClient.Set(r.ctx, key, otp, ttl).Err()
}
func (r *otpRepository) GenerateTokenPass(email string) (string, error) {
	passToken := uuid.New().String()
	r.redisClient.Set(context.Background(), "reset:"+passToken, email, 10*time.Minute)
	return passToken, nil
}
func (r *otpRepository) ValidateOTPPass(validate model.OTPRequestValidate) (bool, error) {
	key := "otp_pass:" + validate.Email
	storedOTP, err := r.redisClient.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return storedOTP == validate.OTPCode, nil
}
