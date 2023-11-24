package main

import "time"

func (app *Config) GenerateToken(email string) (string, error) {

	_, tokenString, err := app.Token.Encode(map[string]interface{}{
		"user_email": email,
		"exp":        time.Now().Add(time.Minute * 15).Unix()})

	if err != nil {
		return "", err
	}

	return tokenString, err

}
