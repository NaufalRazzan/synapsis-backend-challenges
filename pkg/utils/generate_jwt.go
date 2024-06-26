package utils

import (
	"synapsis-backend-test/configs"
	"synapsis-backend-test/models"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(user models.Users) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.User_id,
		"email":    user.Email,
		"password": user.Password,
	})

	signedToken, err := token.SignedString([]byte(configs.GetConfig().JWTSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
