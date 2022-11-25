package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"io/ioutil"
	"net/http"
)

func GetMiscDataByRecordID(w http.ResponseWriter, req *http.Request) {
	var request map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Log.Println(err)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	}
	//defer req.Body.Close()
	x0 := json.Unmarshal
	x0(body, &request)
	response, success, _, msg := models.GetMiscDataByRecordID(request)
	ThrowMiscReponse(msg, response, w, success)

}

func ThrowMiscReponse(successMessage string, responseData map[string]interface{}, w http.ResponseWriter, success bool) {
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
