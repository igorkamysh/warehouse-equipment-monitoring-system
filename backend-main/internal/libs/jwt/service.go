package jwt

import (
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type serviceJWT struct{}

func NewService() *serviceJWT {
	return &serviceJWT{}
}

type jwtData struct {
	UserId      int64
	PhoneNumber string
	JobPosition entities.UserJob
	Exp         int64
}

func (s *serviceJWT) GenerateToken(user entities.User, secret string, tokenTTL time.Duration) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["phone_number"] = user.PhoneNumber
	claims["job_position"] = user.JobPosition
	claims["exp"] = time.Now().Add(tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "create JWT token")
	}
	return tokenString, nil
}

func GetDataFromClaims(claims map[string]interface{}) (jwtData, error) {
	var data jwtData

	userId, ok := claims["id"].(float64)
	if !ok {
		return data, errors.New("")
	}
	data.UserId = int64(userId)

	phoneNumber, ok := claims["phone_number"].(string)
	if !ok {
		return data, errors.New("failed get phone number from claims")
	}
	data.PhoneNumber = phoneNumber

	userJob, ok := claims["job_position"]
	if !ok {
		return data, errors.New("failed get job position from claims")
	}

	job, ok := userJob.(entities.UserJob)
	if !ok {
		return data, errors.New("")
	}
	data.JobPosition = job

	exp, ok := claims["exp"].(float64)
	if !ok {
		return data, errors.New("failed to parse `exp` from claims")
	}

	data.Exp = int64(exp)
	return data, nil
}
