//***************************//
// Package Name: Router
// Date Of Creation: 09/01/2021
// Authour Name: Moloy Mondal
// History: N/A
// Synopsis: Roueter defination
// Functions: NewRouter,PostMiddleware,ThrowBlankResponse
// Global Variable: N/A
// Version: 1.0.0
//***************************//

package router

import (
	"ifixRecord/ifix/entities"
	"log"
	"net/http"
)

// NewRouter method is used to map with path and handler
func NewRouter() {
	for _, route := range routes {
		if route.Method == "POST"{
			http.Handle("/recordapi"+route.Pattern, PostMiddleware(http.HandlerFunc(route.HandlerFunc)))
		} else {
			http.Handle("/recordapi"+route.Pattern, GetMiddleware(http.HandlerFunc(route.HandlerFunc)))
		}
	}
}

// PostMiddleware method is used to handle post method. No other method will not applied
func PostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// POST method checking. No other method will be allowed within this function
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		//w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type, Authorization,Auth")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,X-Requested-With,Auth")
		w.Header().Set("Access-Control-Expose-Headers", "Autho")
		w.Header().Set("Cache-Control", "no-cache,no-store")
		if req.Method != "POST" {
			log.Println(req.Method + " is called in " + req.URL.Path)
			entities.ThrowJSONResponse(entities.NotPostMethodResponse(), w)
			return
		}
		next.ServeHTTP(w, req)
	})
}

// GetMiddleware method is used to handle post method. No other method will not applied
func GetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// POST method checking. No other method will be allowed within this function
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type, Authorization,Auth")
		w.Header().Set("Access-Control-Expose-Headers", "Autho")
		w.Header().Set("Cache-Control", "no-cache,no-store")
		if req.Method != "GET" {
			log.Println(req.Method + " is called in " + req.URL.Path)
			entities.ThrowJSONResponse(entities.NotPostMethodResponse(), w)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func ThrowBlankResponse(w http.ResponseWriter, req *http.Request) {
	entities.ThrowJSONResponse(entities.BlankPathCheckResponse(), w)
}
