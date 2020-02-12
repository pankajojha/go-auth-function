package main

import (
	"fmt"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Welcome(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	headers := http.Header(request.MultiValueHeaders)
	tknStr := headers.Get("token")

	inputJSON := request.Body
	fmt.Println(" inputJSON ", inputJSON, headers, tknStr)

	if tknStr == "" {
		// If the structure of the body is wrong, return an HTTP error
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not StatusUnauthorized'}", StatusCode: http.StatusUnauthorized}, nil
	}

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not StatusUnauthorized'}", StatusCode: http.StatusUnauthorized}, nil
		}
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not StatusUnauthorized'}", StatusCode: http.StatusBadRequest}, nil
	}
	if !tkn.Valid {
		return events.APIGatewayProxyResponse{Body: "{status:403, success:false, reason : 'You are not StatusUnauthorized'}", StatusCode: http.StatusUnauthorized}, nil
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))

	return events.APIGatewayProxyResponse{ // Success HTTP response
		Body:       "Welcome %s!" + string(claims.Username),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Welcome)
}
