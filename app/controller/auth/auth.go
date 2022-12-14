package auth

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"time"
	"wb/app/model"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contract interface {
	Login(username, password string) string
	Register(username, password, fullname string)
	Logout(sessionToken string)
}

type Auth struct {
	redisClient *redis.Client
	dbConn      *gorm.DB
}

func NewAuth(redisClient *redis.Client, dbConn *gorm.DB) Contract {
	return &Auth{
		redisClient: redisClient,
		dbConn:      dbConn,
	}
}

func (a *Auth) Login(username, password string) string {

	auth := model.Auth{
		Username: username,
	}

	auth.FindByUsername(a.dbConn)

	hashedPassword := getSHA384Password(password, auth.Salt)

	if hashedPassword != auth.Password {
		return ""
	}

	sessionToken := uuid.NewString()
	a.redisClient.Set(context.TODO(), sessionToken, auth.ID, 30*time.Minute)

	return sessionToken
}

func (a *Auth) Register(username, password, fullname string) {
	salt := uuid.NewString()
	hashedPassword := getSHA384Password(password, salt)
	auth := &model.Auth{
		Username: username,
		Password: hashedPassword,
		Salt:     salt,
		Fullname: fullname,
	}

	auth.Create(a.dbConn)
}

func (a *Auth) Logout(sessionToken string) {
	a.redisClient.Del(context.TODO(), sessionToken)
}

func getSHA384Password(pass string, salt string) string {
	hasher := sha512.New384()
	hasher.Write([]byte(salt))
	hasher.Write([]byte("\xf3\xf4\xe6")) // Magic
	hasher.Write([]byte(pass))
	hasher.Write([]byte("\xdb\xf9\xb0")) // More magic
	return hex.EncodeToString(hasher.Sum(nil))
}
