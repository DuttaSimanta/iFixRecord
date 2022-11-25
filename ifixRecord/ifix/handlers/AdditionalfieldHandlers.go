package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"log"
	"net/http"
)

func ThrowAdditionalFieldResponse(successMessage string, recordResponseData []entities.AdditionalFieldEntity, w http.ResponseWriter, success bool) {
	var response = entities.AdditionalFieldResponeData{}
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

func GetAdditionalFields(w http.ResponseWriter, req *http.Request) {
	var data = entities.AdditionalfieldRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, _ := models.GetAdditionalFields(&data)
		//log.Print(data)
		ThrowAdditionalFieldResponse("", data, w, success)
	}
}

func GetAdditionalFieldsBYTypeCat(w http.ResponseWriter, req *http.Request) {
	var data = entities.AdditionalfieldRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, _ := models.GetAdditionalFieldsBYTypeCat(&data)
		//log.Print(data)
		ThrowAdditionalFieldResponse("", data, w, success)
	}
}
