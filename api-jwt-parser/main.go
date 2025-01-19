package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

// mySingingkey is the secret key used to sign the JWT, we need it here to verify the JWT signature and auhtorize the user to access the resource
var mySigningKey = []byte(os.Getenv("SECRET"))

// handler function for the home page
func homePage(w http.ResponseWriter, r *http.Request) {
	// write the response to the user
	w.Write([]byte("Super Secret Information"))
}

// isAuthorized is a middleware function that will check if the user is authorized to access the resource or not
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the "Token" header
		tokenString := r.Header.Get("Token")
		if tokenString == "" {
			http.Error(w, "Token is missing", http.StatusForbidden)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		})

		// Handle parsing and validation errors
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Verify custom claims (audience and issuer)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Check audience (aud)
		expectedAud := "parsergeneratorservices.io"
		aud, ok := claims["aud"].(string)
		if !ok || aud != expectedAud {
			http.Error(w, "Invalid audience", http.StatusUnauthorized)
			return
		}

		// Check issuer (iss)
		expectedIss := "ahmad.io"
		iss, ok := claims["iss"].(string)
		if !ok || iss != expectedIss {
			http.Error(w, "Invalid issuer", http.StatusUnauthorized)
			return
		}

		// If everything is valid, proceed to the next handler
		endpoint(w, r)
	})
}

// this function is the handler for the home page, it will be called when the user access the home page
func handleRequests() {
	// handler for the route "/" and it will call the middleware function isAuthorized to see if we can display the home page or not
	// handle will allow the default mux to handle the request the isAuthorized function will return a type of http.Handler so we can use it here
	// registers this handler for the route within our default mux
	http.Handle("/", isAuthorized(homePage))
	// listen and serve the server at port 9001 and fatal so the server will not continue if there is an error and log the error if no error then nothing will be logged
	// use default mux to handle the request
	// limitations is that we cant have nested routes and
	log.Fatal(http.ListenAndServe(":9001", nil))
}

// calling handleRequest function to call the handleRequest function
func main() {
	fmt.Printf("Starting server at port 9001\n")
	handleRequests()

}
