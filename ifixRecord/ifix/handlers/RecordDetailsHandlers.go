package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"log"
	"net/http"
)

func ThrowRecordDetailsResponse(successMessage string, recordResponseData []entities.RecordDetailsEntity, w http.ResponseWriter, success bool) {
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

func ThrowRecordCatDetailsResponse(successMessage string, recordResponseData entities.RecordCatDetailsEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordCatDetailsResponeData{}
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

func GetRecordDetails(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, _ := models.GetRecordDetails(&data)
		//log.Print(data)
		ThrowRecordDetailsResponse("", data, w, success)
	}
}
func GetRecordDetailsByNo(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetRecordDetailsByNo(&data)
		//log.Print(data)
		ThrowRecordDetailsResponse(msg, data, w, success)
	}
}

func GetRecordDetailsByNoForlinkrecord(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetRecordDetailsByNoForlinkrecord(&data)
		//log.Print(data)
		ThrowRecordDetailsResponse(msg, data, w, success)
	}
}

func SaveChildRecord(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.SaveChildRecord(&data)
		//log.Print(data)
		ThrowRecordDetailsResponse(msg, data, w, success)
	}
}

func RemoveChildRecord(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.RemoveChildRecord(&data)
		//log.Print(data)
		ThrowRecordDetailsResponse(msg, data, w, success)
	}
}

func GetChildRecordsBYParentID(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetChildRecordsBYParentID(&data)
		//log.Print(data)
		ThrowRecordDetailsResponse(msg, data, w, success)
	}
}

func GetRecordCatDetails(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetRecordCatDetails(&data)
		//log.Print(data)
		ThrowRecordCatDetailsResponse(msg, data, w, success)
	}
}

func GetCatByLastID(w http.ResponseWriter, req *http.Request) {
	var data = entities.ChildRecordSearchEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetRecordCatByLastID(&data)
		//log.Print(data)
		ThrowRecordCatDetailsResponse(msg, data, w, success)
	}
}

func GetRecordAccesspermissionByNo(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetRecordAccessPermissionByNo(&data)
		//log.Print(data)
		ThrowRecordAccessDetailsResponse(msg, data, w, success)
	}
}

func ThrowRecordAccessDetailsResponse(successMessage string, recordResponseData entities.RecordAccessEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordAccessDetailsResponeData{}
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
