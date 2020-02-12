package main

import (
	"log"
	"net/http"
)

func main() {
	// "Signin" and "Welcome" are the handlers
	http.HandleFunc("/token", GetToken)
	http.HandleFunc("/refresh", RefreshToken)
	http.HandleFunc("/welcome", Welcome)

	//http.HandleFunc("/validate", ValidateMiddleware(Welcome))

	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		authorizationHeader := req.Header.Get("authorization")
// 		if authorizationHeader != "" {
// 			bearerToken := strings.Split(authorizationHeader, " ")
// 			if len(bearerToken) == 2 {
// 				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
// 					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 						return nil, fmt.Errorf("There was an error")
// 					}
// 					return []byte("secret"), nil
// 				})
// 				if error != nil {
// 					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
// 					return
// 				}
// 				if token.Valid {
// 					context.Set(req, "decoded", token.Claims)
// 					next(w, req)
// 				} else {
// 					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
// 				}
// 			}
// 		} else {
// 			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
// 		}
// 	})
// }
