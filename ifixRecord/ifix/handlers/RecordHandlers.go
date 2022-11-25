package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"ifixRecord/ifix/validators"
	"net/http"
)

//ThrowMstClientResponse function is used to throw success response in JSON format
func ThrowMstClientResponse(successMessage string, recordnumber string, w http.ResponseWriter, ID int64) {
	var response = entities.RecordResponeData{}
	response.Success = true
	response.Message = successMessage
	response.Response = recordnumber
	response.ID = ID
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// AddRecordData function is used to handle record data insert operation
func AddRecordData(w http.ResponseWriter, req *http.Request) {
	var recordData = entities.RecordEntity{}
	var responseError []string
	jsonError := recordData.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		responseError = validators.ValidateAddRecordData(recordData)
		logger.Log.Println(responseError)

		if len(responseError) == 0 {
			recordnumber, ID, modelResponseError := models.CreateRecordModel(&recordData)
			logger.Log.Println(modelResponseError)
			if modelResponseError != nil {
				entities.ThrowJSONResponse(entities.DbErrorResponse(modelResponseError.Error()), w)
			} else {
				//data, success, _, msg := models.GetAllAssetdifferentiation(&data)
				ThrowMstClientResponse("record data inserted successfully.", recordnumber, w, ID)
			}
		} else {
			entities.ThrowErrorResponse(responseError, w)
		}

		//models.APIcall()
	}
}
