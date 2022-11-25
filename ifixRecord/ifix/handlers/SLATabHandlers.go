package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"net/http"
)

func ThrowSLATabAllResponse(successMessage string, responseData entities.SLATabresponsesEntity, w http.ResponseWriter, success bool) {
	var response = entities.SLATabAllResponse{}
	response.Success = success
	response.Message = successMessage
	response.Details = responseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowSLAmeterlResponse(successMessage string, responseData entities.SLAMeterEntity, w http.ResponseWriter, success bool) {
	var response = entities.SLAMeterAllResponse{}
	response.Success = success
	response.Message = successMessage
	response.Details = responseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetSLATabvalues(w http.ResponseWriter, req *http.Request) {
	var data = entities.SLATabEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetSLATabvalues(&data)
		ThrowSLATabAllResponse(msg, data, w, success)
	}
}

func GetSLAResolution(w http.ResponseWriter, req *http.Request) {
	var data = entities.SLATabEntity{}
	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.", jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetSLAResolution(&data)
		logger.Log.Println(msg)
		ThrowSLAmeterlResponse(msg, data, w, success)
	}
}

func SLAcalculation(w http.ResponseWriter, req *http.Request) {
	var data = entities.SLAEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		//data, success, _, msg := models.GetSLAResolution(&data)
		//logger.Log.Println(msg)
		//ThrowSLAmeterlResponse(msg, data, w, success)
	}
}
