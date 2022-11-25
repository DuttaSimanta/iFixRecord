package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"log"
	"net/http"
)

func GetAdditionalInfoBasedonCategory(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcreaterequestEntity{}
	logger.Log.Println("In GetAdditionalInfoBasedonCategory Handler")
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, _ := models.GetAdditionalInfoBasedonCategory(&data)
		//log.Print(data)
		ThrowPriorityResponse("", data, w, success)
	}
}
//ThrowRecordcreateResponse function is used to throw success response in JSON format
func ThrowRecordcreateResponse(successMessage string, recordResponseData entities.RecordcreateEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordcreateResponeData{}
	response.Status = success
	response.Message = successMessage
	response.Response = recordResponseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

//ThrowRecordcreateResponse function is used to throw success response in JSON format
func ThrowRecordtypeResponse(successMessage string, recordResponseData []entities.RecordtypedetailsEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordtypeResponeData{}
	response.Status = success
	response.Message = successMessage
	response.Response = recordResponseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

//ThrowRecordcreateResponse function is used to throw success response in JSON format
func ThrowRecordcatchildResponse(successMessage string, recordResponseData []entities.RecordcatchildEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordcatchildResponeData{}
	response.Status = success
	response.Message = successMessage
	response.Response = recordResponseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowPriorityResponse(successMessage string, recordResponseData entities.RecordcatchildNEstimatedEfforEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordPriotiryResponeData{}
	response.Status = success
	response.Message = successMessage
	response.Response = recordResponseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetRecordcreatedata(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcreaterequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, _ := models.GetRecordcreatedata(&data)
		//log.Print(data)
		ThrowRecordcreateResponse("", data, w, success)
	}
}

func GetRecordcatchild(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcreaterequestEntity{}
	logger.Log.Println("In Handler")
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetRecordcatchild(&data)
		//log.Print(data)
		ThrowRecordcatchildResponse(msg, data, w, success)
	}
}

func GetRecordtypedata(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcreaterequestEntity{}
	logger.Log.Println("In Handler")
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, _ := models.GetRecordtypedata(&data)
		//log.Print(data)
		ThrowRecordtypeResponse("", data, w, success)
	}
}

func GetRecordprioritydata(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcreaterequestEntity{}
	logger.Log.Println("In GetRecordprioritydata Handler")
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, _ := models.GetRecordprioritydata(&data)
		//log.Print(data)
		ThrowPriorityResponse("", data, w, success)
	}
}
