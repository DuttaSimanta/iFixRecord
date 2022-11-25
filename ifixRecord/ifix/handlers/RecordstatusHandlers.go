package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"net/http"
)

// ThrowRecordstatusIntResponse function is used to throw success response of integer data in JSON format
func ThrowRecordstatusIntResponse(successMessage string, responseData int64, w http.ResponseWriter, success bool) {
	var response = entities.RecordstatusResponeData{}
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

func Updaterecordstatusvalue(w http.ResponseWriter, req *http.Request) {
	logger.Log.Println("======================================================  updaterecordstatus api STARTED======================================")

	var data = entities.RecordstatusEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.RecordStatusUpdation(&data)
			logger.Log.Println("======================================================  updaterecordstatus api ENDED======================================")

		ThrowRecordstatusIntResponse(msg, data, w, success)
	}
}
