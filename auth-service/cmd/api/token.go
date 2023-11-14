package main

func (app *Config) GenerateToken(email string) (string, error) {

	_, tokenString, err := app.Token.Encode(map[string]interface{}{"user_email": email})

	if err != nil {
		return "", err
	}

	return tokenString, err

}
