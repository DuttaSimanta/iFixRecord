package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"io/ioutil"
	"log"
	"net/http"
)

func GetAssetfieldSpecificDataBYRecordID(w http.ResponseWriter, req *http.Request) {
	var data = entities.AssetAttrNameValRequestEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		response, success, _, msg := models.GetAssetfieldSpecificDataBYRecordID(&data)
		ThrowAssetArrReponse(msg, response, w, success)
	}
}
func GetAssetTypeByRecordID(w http.ResponseWriter, req *http.Request) {
	var request map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &request)
	response, success, _, msg := models.GetAssetTypeByRecordID(request)
	ThrowAssetArrReponse(msg, response, w, success)

}

func GetRecordAssetByID(w http.ResponseWriter, req *http.Request) {
	var request map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &request)
	response, success, _, msg := models.GetRecordAssetByID(request)
	ThrowAssetReponse(msg, response, w, success)

}

func GetAssetDetailsByAssetID(w http.ResponseWriter, req *http.Request) {
	var request map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &request)
	response, success, _, msg := models.GetAssetDetailsByAssetID(request)
	ThrowAssetReponse(msg, response, w, success)

}
func UpdateRecordAsset(w http.ResponseWriter, req *http.Request) {
	var request = entities.UpdateRecordAssetEntity{}
	jsonError := request.FromJSON(req.Body)
	if jsonError != nil {
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
		logger.Log.Println(jsonError)

	} else {
		response, success, _, msg := models.UpdateRecordAsset(&request)
		ThrowAssetReponse(msg, response, w, success)
	}
}
func GetAssetHistroyByAssetID(w http.ResponseWriter, req *http.Request) {
	var request = entities.FetchAssetHistoryRequest{}
	jsonError := request.FromJSON(req.Body)
	if jsonError != nil {
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
		logger.Log.Println(jsonError)

	} else {
		response, success, _, msg := models.GetAssetHistroyByAssetID(&request)
		ThrowAssetReponse(msg, response, w, success)
	}
}

func InsertRecordAsset(w http.ResponseWriter, req *http.Request) {
	var request = entities.InsertRecordAssetEntity{}
	jsonError := request.FromJSON(req.Body)
	if jsonError != nil {
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
		logger.Log.Println(jsonError)

	} else {
		response, success, _, msg := models.InsertRecordAsset(&request)
		ThrowAssetReponse(msg, response, w, success)
	}
}


func AddAssetWithRecord(w http.ResponseWriter, req *http.Request) {
	var request = entities.RecordAssetRequestEntity{}
	jsonError := request.FromJSON(req.Body)
	if jsonError != nil {
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
		logger.Log.Println(jsonError)

	} else {
		response, success, _, msg := models.AddAssetWithRecord(&request)
		ThrowAssetReponse(msg, response, w, success)
	}
}

func DeleteAssetFromRecord(w http.ResponseWriter, req *http.Request) {
	var request map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &request)
	response, success, _, msg := models.DeleteAssetFromRecord(request)
	ThrowAssetReponse(msg, response, w, success)
}

func GetAllAssetTypeNDetailsByRecordID(w http.ResponseWriter, req *http.Request) {
	var request map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &request)
	response, success, _, msg := models.GetAllAssetTypeNDetailsByRecordID(request)
	ThrowAssetArrReponse(msg, response, w, success)

}

func ThrowAssetReponse(successMessage string, responseData map[string]interface{}, w http.ResponseWriter, success bool) {
	response := make(map[string]interface{}, 3)
	response["details"] = responseData
	response["message"] = successMessage
	response["success"] = success
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowAssetArrReponse(successMessage string, responseData []map[string]interface{}, w http.ResponseWriter, success bool) {
	response := make(map[string]interface{}, 3)
	response["details"] = responseData
	response["message"] = successMessage
	response["success"] = success
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
