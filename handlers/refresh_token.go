package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func RefreshToken(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route

	//tknStr := request.MultiValueHeaders["token"]
	headers := http.Header(request.MultiValueHeaders)
	tknStr := headers.Get("token")

	inputJSON := request.Body
	fmt.Println(" inputJSON ", inputJSON, headers, tknStr)

	if tknStr == "" {
		// If the structure of the body is wrong, return an HTTP error
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not autherized'}", StatusCode: http.StatusBadRequest}, nil
	}

	//	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not autherized'}", StatusCode: http.StatusUnauthorized}, nil
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not autherized'}", StatusCode: http.StatusUnauthorized}, nil
		}
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not autherized'}", StatusCode: http.StatusUnauthorized}, nil
	}
	// (END) The code uptil this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not StatusBadRequest'}", StatusCode: http.StatusBadRequest}, nil
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not StatusInternalServerError'}", StatusCode: http.StatusInternalServerError}, nil
	}

	// Set the new token as the users `session_token` cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })

	m := map[string][]string{
		"token": {tokenString},
		//	"Expires": {expirationTime},
	}

	h := http.Header(m)
	fmt.Println(h.Get("token"))

	return events.APIGatewayProxyResponse{ // Success HTTP response
		Body:       string(tokenString),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(RefreshToken)
}
