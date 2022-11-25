package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"log"
	"net/http"
)

func ThrowRecordParentDetailsResponse(successMessage string, recordResponseData int64, w http.ResponseWriter, success bool) {
	var response = entities.RecordDetailsParentResponeData{}
	response.Status = success
	response.Message = successMessage
	response.Details = recordResponseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowParentRecordDetailsResponse(successMessage string, recordResponseData []entities.RecordDetailsEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordDetailsResponeData{}
	response.Status = success
	response.Message = successMessage
	response.Details = recordResponseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func Addparentfromchild(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.SaveParentRecord(&data)
		//log.Print(data)
		ThrowRecordParentDetailsResponse(msg, data, w, success)
	}
}

func GetParentRecordDetailsByNo(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetParentRecordDetailsByNo(&data)
		//log.Print(data)
		ThrowParentRecordDetailsResponse(msg, data, w, success)
	}
}

func ChildRecordSearchCriteria(w http.ResponseWriter, req *http.Request) {
	var data = entities.ChildRecordSearchEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.ChildRecordSearchCriteria(&data)
		//log.Print(data)
		ThrowRecordDetailsResponse(msg, data, w, success)
	}
}
