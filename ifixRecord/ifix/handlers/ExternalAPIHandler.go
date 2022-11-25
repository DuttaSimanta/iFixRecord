package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"log"
	"net/http"
)
func GetifixRecordDockerStatus(w http.ResponseWriter, req *http.Request) {
	// var data = entities.ReportDownloadListEntity{}
	// jsonError := data.FromJSON(req.Body)
	// if jsonError != nil {
	// 	log.Print(jsonError)
	// 	logger.Log.Println(jsonError)
	// 	entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	// } else {
	// 	data, success, _, msg := models.GetDownloadList(&data)
	// 	ThrowGetDownloadListResponse(msg, data, w, success)
	// }
	type response struct {
		Success bool `json:"success"`
	}
	var res response
	res.Success = true
	jsonResponse, jsonError := json.Marshal(res)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
func UpdateExternalRecordAttachment(w http.ResponseWriter, req *http.Request) {
    var data = entities.FileAttacmentToRecordEntity{}
    jsonError := data.FromJSON(req.Body)
    if jsonError != nil {
        log.Print(jsonError)
        logger.Log.Println(jsonError)
        entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
    } else {
        var token = req.Header.Get("Authorization")
        success := models.CheckToken(token, data.Userid)
        if !success {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        success, _, data := models.UpdateExternalRecordAttachment(&data)
        ThrowExternalRecordCreateAllResponse("", data, w, success)
    }
}

//==============================================
func GetExternalRecordDetailsNo(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntityAPI{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		var token = req.Header.Get("Authorization")
		success := models.CheckToken(token, data.Userid)
		if !success {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, success, _, msg := models.GetExternalRecordDetailsByNo(&data)
		ThrowExternalRecordAllResponse(msg, data, w, success)
	}
}

func GetExternalRecordDetailsbyDate(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordDetailsRequestEntityAPI{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		var token = req.Header.Get("Authorization")
		success := models.CheckToken(token, data.Userid)
		if !success {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data, success, _, msg := models.GetExternalRecordDetailsByDate(&data)
		ThrowtokenExternalRecordAllResponse(msg, data, w, success)
	}
}
func ThrowtokenExternalRecordAllResponse(successMessage string, responseData []entities.RecordDetailsEntityAPI, w http.ResponseWriter, success bool) {
	var response = entities.TokenRecordDetailsEntityResponse{}
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
func ThrowExternalRecordAllResponse(successMessage string, responseData entities.RecordDetailsEntityAPI, w http.ResponseWriter, success bool) {
	var response = entities.RecordDetailsEntityResponse{}
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

func ThrowExternalRecordCreateAllResponse(successMessage string, responseData string, w http.ResponseWriter, success bool) {
	var response = entities.ExternalRecordEntityResponse{}
	response.Success = success
	response.Message = responseData
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func CreateExternalRecord(w http.ResponseWriter, req *http.Request) {
	var data = entities.ExternalCreateRecord{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		var token = req.Header.Get("Authorization")
		success := models.CheckToken(token, data.Userid)
		if !success {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		success, _, data := models.ExternalRecordCreate(&data)
		ThrowExternalRecordCreateAllResponse("", data, w, success)
	}
}

func CreateExternalRecordStatusUpdate(w http.ResponseWriter, req *http.Request) {
logger.Log.Println("============================================================== CreateExternalRecordStatusUpdate STARTED =============================================================:")
	var data = entities.ExternalCreateRecord{}
	jsonError := data.FromJSON(req.Body)
	postBody, _ := json.Marshal(data)
	logger.Log.Println("PAYLOAD->>>>>>>>>>>>>>>>>>>>>>>>>>>>", string(postBody))
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		var token = req.Header.Get("Authorization")
		success := models.CheckToken(token, data.Userid)
		if !success {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		success, _, data := models.ExternalRecordStatusUpdate(&data)
				logger.Log.Println("============================================================== CreateExternalRecordStatusUpdate ENDED =============================================================")

		ThrowExternalRecordCreateAllResponse("", data, w, success)
	}
}

func CreateExternalRecordGrpUpdate(w http.ResponseWriter, req *http.Request) {
	var data = entities.ExternalCreateRecord{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		var token = req.Header.Get("Authorization")
		success := models.CheckToken(token, data.Userid)
		if !success {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		success, _, data := models.ExternalRecordGrpUpdate(&data)
		ThrowExternalRecordCreateAllResponse("", data, w, success)
	}
}

func CreateExternalRecordUserUpdate(w http.ResponseWriter, req *http.Request) {
	var data = entities.ExternalCreateRecord{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		var token = req.Header.Get("Authorization")
		success := models.CheckToken(token, data.Userid)
		if !success {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		success, _, data := models.ExternalRecordUserUpdate(&data)
		ThrowExternalRecordCreateAllResponse("", data, w, success)
	}
}

func ExternalRecordCommentUpdate(w http.ResponseWriter, req *http.Request) {
	var data = entities.ExternalCreateRecord{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		var token = req.Header.Get("Authorization")
		success := models.CheckToken(token, data.Userid)
		if !success {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		success, _, data := models.ExternalRecordInternalCommentUpdate(&data)
		ThrowExternalRecordCreateAllResponse("", data, w, success)
	}
}
