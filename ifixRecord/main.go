//***************************//
// Package Name: Main
// Date Of Creation: 11/01/2021
// Authour Name: Simanta Dutta
// History: N/A
// Synopsis: This is Entry Point of API.
// Functions: main
// Global Variable: N/A
// Version: 1.0.0
//***************************//
package main

import (
	"ifixRecord/ifix/logger"
	router "ifixRecord/ifix/routers"
	"net/http"
)

// main function is uded to call the API function and to start the http server for API
func main() {
//	logger.Log.Panicln("STopped Record docke")
	router.NewRouter()
	logger.Log.Println("Server started at 8083 port")
	logger.Log.Fatal(http.ListenAndServe(":8083", nil))
}
