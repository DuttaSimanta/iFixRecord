//***************************//
// Package entities
// Date Of Creation: 19/12/2020
// Authour Name: Moloy Mondal
// History: N/A
// Synopsis: This file is used to define commonly used api response format and common response as output.
// Functions: BlankPathCheckResponse,NotPostMethodResponse, InternalServerErrorResponse, JSONParseErrorResponse,
// DbConErrorResponse, DbErrorResponse, ThrowJSONResponse, ThrowErrorResponse
// BlankPathCheckResponse() Parameter:  N/A
// NotPostMethodResponse() Parameter:  N/A
// InternalServerErrorResponse() Parameter:  N/A
// JSONParseErrorResponse() Parameter:  N/A
// DbConErrorResponse Parameter:  N/A
// DbErrorResponse Parameter:  N/A
// ThrowJSONResponse Parameter: (<Structure APIResponse>,<http.ResponseWriter>)
// ThrowErrorResponse Parameter: (<[]string>,<http.ResponseWriter>)
// Global Variable: N/A
// Version: 1.0.0
//***************************//

package entities

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIResponse Structure used to handle http response using json
type APIResponse struct {
	Status   bool   `json:"success"`
	Message  string `json:"message"`
	Response string `json:"response"`
}

// ErrorResponse Structure used to handle error  response using json
type ErrorResponse struct {
	Status   bool     `json:"success"`
	Message  string   `json:"message"`
	Response []string `json:"response"`
}

// BlankPathCheckResponse function is used to return blank path response
func BlankPathCheckResponse() APIResponse {
	var response = APIResponse{}
	response.Status = false
	response.Message = "404 not found."
	log.Println("Blank request called")
	return response
}

// NotPostMethodResponse function is used to return not post method response
func NotPostMethodResponse() APIResponse {
	var response = APIResponse{}
	response.Status = false
	response.Message = "405 method not allowed."
	return response
}

// InternalServerErrorResponse function is used to return Internal server error response
func InternalServerErrorResponse() APIResponse {
	var response = APIResponse{}
	response.Status = false
	response.Message = "Internal Server Error."
	return response
}

// JSONParseErrorResponse function is used to return JSON parse error response
func JSONParseErrorResponse() APIResponse {
	var response = APIResponse{}
	response.Status = false
	response.Message = "501 JSON parse Error."
	return response
}

// DbConErrorResponse function is used to return database connection error response
func DbConErrorResponse() APIResponse {
	var response = APIResponse{}
	response.Status = false
	response.Message = "502 DB Connection Error."
	return response
}

// DbErrorResponse function is used to return database Insertion error response
func DbErrorResponse(message string) APIResponse {
	var response = APIResponse{}
	response.Status = false
	response.Message = message
	return response
}

// ThrowJSONResponse function is used to throw response in JSON format
func ThrowJSONResponse(response APIResponse, w http.ResponseWriter) {
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// ThrowErrorResponse function is used to throw response in JSON format
func ThrowErrorResponse(responseErr []string, w http.ResponseWriter) {
	var response = ErrorResponse{}
	response.Status = false
	response.Message = "201 Operational Error."
	response.Response = responseErr
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
