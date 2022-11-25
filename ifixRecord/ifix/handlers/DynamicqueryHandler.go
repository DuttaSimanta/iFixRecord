package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"io/ioutil"
	"net/http"
)

func ThrowGetDynamicQueryResultV1Response(successMessage string, recordResponseData []interface{}, w http.ResponseWriter, success bool) {
	var response = entities.QueryResponeData{}
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
func ThrowGetDynamicQueryResultResponse(successMessage string, recordResponseData entities.QueryDetailsEntity, w http.ResponseWriter, success bool) {
	var response = entities.QueryResultResponeData{}
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

func ThrowGetDynamicQueryCountResultResponse(successMessage string, recordResponseData map[string]interface{}, w http.ResponseWriter, success bool) {
	var response = entities.QueryCountResultResponeData{}
	response.Status = success
	response.Message = successMessage
	response.Response = recordResponseData
	jsonResponse, jsonError := json.Marshal(response)
	//var result map[string]interface{}
	//json.Unmarshal(jsonResponse, &result)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)
}
func GetDynamicQueryResult(w http.ResponseWriter, req *http.Request) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &result)
	response, success, _, msg := models.GetDynamicQueryResultTemp(result)
	ThrowGetDynamicQueryResultResponse(msg, response, w, success)

}
func GetDynamicQueryResultV1(w http.ResponseWriter, req *http.Request) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &result)
	response, success, _, msg := models.GetDynamicQueryResult(result)
	//log.Println(response)
	ThrowGetDynamicQueryResultV1Response(msg, response, w, success)

}
func GetDynamicCountQueryResult(w http.ResponseWriter, req *http.Request) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &result)
	response, success, _, msg := models.GetDynamicQueryCountResultTemp(result)
	result["total"] = response
	ThrowGetDynamicQueryCountResultResponse(msg, result, w, success)

}

func RecordGridResult(w http.ResponseWriter, req *http.Request) {
	var reqData map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		json.Unmarshal(body, &reqData)
		response, success, _, msg := models.RecordGridResult(reqData)
		ThrowGetDynamicQueryCountResultResponse(msg, response, w, success)
	}
}

func RecordGridResultOnly(w http.ResponseWriter, req *http.Request) {
	var reqData map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		json.Unmarshal(body, &reqData)
		response, success, _, msg := models.RecordGridResultOnly(reqData)
		ThrowGetDynamicQueryCountResultResponse(msg, response, w, success)
	}
}

func RecordFilterAdd(w http.ResponseWriter, req *http.Request) {
	var reqData map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		json.Unmarshal(body, &reqData)
		response, success, _, msg := models.RecordFilterAdd(reqData)
		ThrowGetDynamicQueryCountResultResponse(msg, response, w, success)
	}
}

func RecordFilterList(w http.ResponseWriter, req *http.Request) {
	var reqData map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		json.Unmarshal(body, &reqData)
		response, success, _, msg := models.RecordFilterList(reqData)
		ThrowGetDynamicQueryCountResultResponse(msg, response, w, success)
	}
}

func RecordFilterDelete(w http.ResponseWriter, req *http.Request) {
	var reqData map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		json.Unmarshal(body, &reqData)
		response, success, _, msg := models.RecordFilterDelete(reqData)
		ThrowGetDynamicQueryCountResultResponse(msg, response, w, success)
	}
}



func RecordFilterUpdate(w http.ResponseWriter, req *http.Request) {
	var reqData map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		json.Unmarshal(body, &reqData)
		response, success, _, msg := models.RecordFilterUpdate(reqData)
		ThrowGetDynamicQueryCountResultResponse(msg, response, w, success)
	}
}

