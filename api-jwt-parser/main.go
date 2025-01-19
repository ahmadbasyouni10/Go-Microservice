package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// mySingingkey is the secret key used to sign the JWT, we need it here to verify the JWT signature and auhtorize the user to access the resource
var mySigningKey = []byte(os.Getenv("SECRET"))

// handler function for the home page
func homePage(w http.ResponseWriter, r *http.Request) {
	// write the response to the user
	fmt.Fprintf(w, "Super Secret Information")
}

// isAuthorized is a middleware function that will check if the user is authorized to access the resource or not
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
		}
	})
}

// this function is the handler for the home page, it will be called when the user access the home page
func handleRequests() {
	// handler for the route "/" and it will call the middleware function isAuthorized to see if we can display the home page or not
	// handle will allow the default mux to handle the request the isAuthorized function will return a type of http.Handler so we can use it here
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
	fmt.Printf("Server started\n")

}
