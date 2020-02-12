package main

import (
	"fmt"

	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"strings"

	"github.com/dgrijalva/jwt-go"
)

var testUsers = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GetToken(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var creds Credentials
	inputJSON := request.Body
	fmt.Println(" inputJSON ", inputJSON)

	err := json.NewDecoder(strings.NewReader(request.Body)).Decode(&creds)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not autherized'}", StatusCode: http.StatusBadRequest}, nil
	}

	expectedPassword, ok := testUsers[creds.Username]

	if !ok || expectedPassword != creds.Password {
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not autherized'}", StatusCode: http.StatusUnauthorized}, nil
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not autherized'}", StatusCode: http.StatusInternalServerError}, nil
	}

	// claims := token.Claims.(jwt.MapClaims)
	// /* Set token claims */
	// claims["admin"] = true
	// claims["name"] = "Ado Kukic"
	// claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	//TODO this should
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })

	m := map[string][]string{
		"token": {tokenString},
		//	"Expires": {expirationTime},
	}
	h := http.Header(m)
	fmt.Println(h.Get("token"))
	/* Finally, write the token to the browser window */

	return events.APIGatewayProxyResponse{ // Success HTTP response
		Body:       string(tokenString),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(GetToken)
}
