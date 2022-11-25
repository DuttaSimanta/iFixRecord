package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"net/http"
)

func ThrowPriorityupdateAllResponse(stageID int64, success bool, msg string, w http.ResponseWriter) {
	var response = entities.RecordPriorityResponeEntity{}
	response.Success = success
	response.Message = msg
	response.StageID = stageID
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func Updaterecordpriority(w http.ResponseWriter, req *http.Request) {
	var recordData = entities.RecordpriorityEntity{}
	jsonError := recordData.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		stageID, success, err, msg := models.Updatepriority(&recordData)
		logger.Log.Println(err)
		if err != nil {
			entities.ThrowJSONResponse(entities.DbErrorResponse(err.Error()), w)
		} else {
			//data, success, _, msg := models.GetAllAssetdifferentiation(&data)
			ThrowPriorityupdateAllResponse(stageID, success, msg, w)
		}

	}
}

func Updaterecordcategory(w http.ResponseWriter, req *http.Request) {
	var recordData = entities.RecordcategoryupdateEntity{}
	jsonError := recordData.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		stageID, success, err, msg := models.Recordcategoryupdate(&recordData)
		logger.Log.Println(err)
		if err != nil {
			entities.ThrowJSONResponse(entities.DbErrorResponse(err.Error()), w)
		} else {
			//data, success, _, msg := models.GetAllAssetdifferentiation(&data)
			ThrowPriorityupdateAllResponse(stageID, success, msg, w)
		}

	}
}

func Updateadditionalfields(w http.ResponseWriter, req *http.Request) {
	var recordData = entities.RecordcategoryupdateEntity{}
	jsonError := recordData.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		stageID, success, err, msg := models.Recordadditionalfieldupdate(&recordData)
		logger.Log.Println(err)
		if err != nil {
			entities.ThrowJSONResponse(entities.DbErrorResponse(err.Error()), w)
		} else {
			//data, success, _, msg := models.GetAllAssetdifferentiation(&data)
			ThrowPriorityupdateAllResponse(stageID, success, msg, w)
		}

	}
}
