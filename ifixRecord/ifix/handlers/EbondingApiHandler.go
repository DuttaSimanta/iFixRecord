package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"log"
	"net/http"
)

func ThrowEbondingApiIntResponse(successMessage string, w http.ResponseWriter, success bool) {
	var response = entities.EbondingApiRespone{}
	response.Success = success
	response.Message = successMessage
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func EbondingTicket(w http.ResponseWriter, req *http.Request) {
	var data = entities.EbondingRecordEntity{}
	jsonError := data.FromJSON(req.Body)
	if jsonError != nil {
		log.Print(jsonError)
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		// var token = req.Header.Get("Authorization")
		// success := models.CheckToken(token, data.Userid)
		// if !success {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	return
		// }
		success, _, msg := models.EbondingTicket(&data)
		ThrowEbondingApiIntResponse(msg, w, success)
	}
}
