package handlers

import (
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/models"
	"net/http"
)

// ThrowRecordcommonAllResponse function is used to throw success response of All data in JSON format
func ThrowRecordcommonAllResponse(successMessage string, responseData []entities.RecordcommonresponseEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordcommonAllResponse{}
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

// ThrowRecordcommonIntResponse function is used to throw success response of integer data in JSON format
func ThrowRecordcommonIntResponse(successMessage string, responseData int64, w http.ResponseWriter, success bool) {
	var response = entities.RecordcommonResponseInt{}
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

// ThrowTermnamesAllResponse function is used to throw success response of All data in JSON format
func ThrowTermnamesAllResponse(successMessage string, responseData []entities.RecordTermnamesEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordTermnamesAllResponse{}
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

func ThrowCountAllResponse(successMessage string, data entities.RecordcountEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordcountAllResponse{}
	response.Success = success
	response.Message = successMessage
	response.Prioritycount = data.Prioritycount
	response.Followupcount = data.Followupcount
	response.Reopencount = data.Reopencount
	response.Pendingvendoractioncount = data.Pendingvendoractioncount
	response.Outboundcount = data.Outboundcount
	response.Aging = data.Aging
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowRecentrecordAllResponse(successMessage string, responseData []entities.RecentrecordEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecentrecordAllResponse{}
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

func ThrowRecordlogsAllResponse(successMessage string, responseData []entities.RecordlogsEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordlogsAllResponse{}
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

func ThrowFrequentrecordAllResponse(successMessage string, responseData []entities.FrequentRecordEntity, w http.ResponseWriter, success bool) {
	var response = entities.FrequentrecordsAllResponse{}
	response.Success = success
	response.Message = successMessage
	response.Details = responseData

	jsonResponse, jsonError := json.Marshal(response)
	//logger.Log.Println("----------------->", jsonResponse)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowParentrecordAllResponse(successMessage string, responseData []entities.ParentticketEntity, w http.ResponseWriter, success bool) {
	var response = entities.ParentticketEntityAllResponse{}
	response.Success = success
	response.Message = successMessage
	response.Details = responseData

	jsonResponse, jsonError := json.Marshal(response)
	//logger.Log.Println("----------------->", jsonResponse)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowActivityAllResponse(successMessage string, responseData []entities.Recordactivitymst, w http.ResponseWriter, success bool) {
	var response = entities.ActivitynamesAllResponse{}
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

func ThrowNewActivitylogsAllResponse(successMessage string, responseData []entities.NewActivitylogsEntity, w http.ResponseWriter, success bool) {
	var response = entities.NewActivityAllResponse{}
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

func ThrowActivitySearchAllResponse(successMessage string, responseData []entities.NewActivitylogsEntity, w http.ResponseWriter, success bool) {
	var response = entities.ActivitySearchAllResponse{}
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

func ThrowPendingstatusAllResponse(successMessage string, responseData []entities.Pendingstatustermvalue, w http.ResponseWriter, success bool) {
	var response = entities.PendingstatusAllResponse{}
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

func ThrowAttachmentAllResponse(successMessage string, responseData []entities.RecordAttachmentfiles, w http.ResponseWriter, success bool) {
	var response = entities.RecordattachmentAllResponse{}
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

func ThrowDocumentAllResponse(successMessage string, w http.ResponseWriter, success bool) {
	var response = entities.DocumentupdateAllResponse{}
	response.Success = success
	response.Message = successMessage
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowcustomervisibleAllResponse(successMessage string, responseData []entities.Customervisiblecomment, w http.ResponseWriter, success bool) {
	var response = entities.CustomervisiblecommentAllResponse{}
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

func ThrowDeleteattchmentAllResponse(successMessage string, w http.ResponseWriter, success bool) {
	var response = entities.DocumentupdateAllResponse{}
	response.Success = success
	response.Message = successMessage
	jsonResponse, jsonError := json.Marshal(response)
	if jsonError != nil {
		logger.Log.Fatal("Internel Server Error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func ThrowTermvaluesbyseqAllResponse(successMessage string, responseData []entities.Recordtermseqvalue, w http.ResponseWriter, success bool) {
	var response = entities.RecordtermseqvalueAllResponse{}
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
func InsertBulkRecordTermvalues(w http.ResponseWriter, req *http.Request) {
	var data = entities.BulkRecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		//models.APIcall()
		data, success, _, msg := models.InsertBulkRecordTermvalues(&data)
		ThrowRecordcommonIntResponse(msg, data, w, success)
	}
}
func InsertTermvalue(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		//models.APIcall()
		data, success, _, msg := models.InsertRecordTermvalues(&data)
		ThrowRecordcommonIntResponse(msg, data, w, success)
	}
}

func InsertMultipleTermvalue(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordmultiplecommonEntity{}
	jsonError := data.FrommultipleJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.InsertMultipleTermvalues(&data)
		ThrowRecordcommonIntResponse(msg, data, w, success)
	}
}

func GetAllTermvalues(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetAllcommontermvalues(&data)
		ThrowRecordcommonAllResponse(msg, data, w, success)
	}
}

func GetTermvaluesbytermid(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetTermvalueagainsttermid(&data)
		ThrowRecordcommonAllResponse(msg, data, w, success)
	}
}

func GetTermnames(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetTermnames(&data)
		ThrowTermnamesAllResponse(msg, data, w, success)
	}
}

func GetTermnamesbystate(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonstateEntity{}
	jsonError := data.FromstateJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetTermnamesbystate(&data)
		ThrowTermnamesAllResponse(msg, data, w, success)
	}
}

//int64, int64, bool, error, string
func GetRecordcount(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.Getrecordcount(&data)
		ThrowCountAllResponse(msg, data, w, success)
	}
}

func Getrecentrecords(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.Getrecentrecord(&data)
		ThrowRecentrecordAllResponse(msg, data, w, success)
	}
}

func GetRecordlogs(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetRecordlogs(&data)
		ThrowRecordlogsAllResponse(msg, data, w, success)
	}
}

func Getfrequentissues(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.Getfrequentissues(&data)
		logger.Log.Println(data)
		ThrowFrequentrecordAllResponse(msg, data, w, success)
	}
}

func GetParentrecord(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetParentrecord(&data)
		logger.Log.Println(data)
		ThrowParentrecordAllResponse(msg, data, w, success)
	}
}

func GetActivitymstnames(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError =GetActivitymstnames validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetActivitymstnames(&data)
		ThrowActivityAllResponse(msg, data, w, success)
	}
}

func GetNewActivitylogs(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetNewActivitylogs(&data)
		ThrowNewActivitylogsAllResponse(msg, data, w, success)
	}
}

func Activitylogsearch(w http.ResponseWriter, req *http.Request) {
	var data = entities.Activitylogsearchcriteria{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.Activitysearchresults(&data)
		ThrowActivitySearchAllResponse(msg, data, w, success)
	}
}

func GetPendingstatustermvalue(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetPendingstatustermvalue(&data)
		ThrowPendingstatusAllResponse(msg, data, w, success)
	}
}

func GetAttachmentfiles(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetAttachmentfiles(&data)
		ThrowAttachmentAllResponse(msg, data, w, success)
	}
}

func Updatedocumentcount(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		success, _, msg := models.Updatedocumentcount(&data)
		ThrowDocumentAllResponse(msg, w, success)
	}
}

func GetTermnamesbyseq(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetTermnamesbysequance(&data)
		ThrowTermnamesAllResponse(msg, data, w, success)
	}
}

func Customervisiblecomment(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.Customervisiblecomment(&data)
		ThrowcustomervisibleAllResponse(msg, data, w, success)
	}
}

func Deleteattchfile(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordAttachmentfiles{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		success, _, msg := models.DropAttachmentfiles(&data)
		ThrowDeleteattchmentAllResponse(msg, w, success)
	}
}

func GetTermvaluebyseq(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.GetTermvaluebysequance(&data)
		ThrowTermvaluesbyseqAllResponse(msg, data, w, success)
	}
}

func Parentchildcollaborationlogs(w http.ResponseWriter, req *http.Request) {
	var data = entities.Activitylogsearchcriteria{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		data, success, _, msg := models.Parentchildcollaborationlogs(&data)
		ThrowParentChildSearchAllResponse(msg, data, w, success)
	}
}

func ThrowParentChildSearchAllResponse(successMessage string, responseData []entities.NewActivitylogsEntity, w http.ResponseWriter, success bool) {
	var response = entities.ParentchildSearchAllResponse{}
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

func GetRecordID(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetRecordID(&data)
		ThrowRecordcommonIntResponse(msg, data, w, success)
	}
}

func GetTabTermnames(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetTabTermnames(&data)
		ThrowTabTermnamesAllResponse(msg, data, w, success)
	}
}

func GetTabTermvalues(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetTabTermvalues(&data)
		ThrowTabTermnamesAllResponse(msg, data, w, success)
	}
}

func ThrowTabTermnamesAllResponse(successMessage string, responseData entities.RecordTabTermsEntity, w http.ResponseWriter, success bool) {
	var response = entities.RecordTabTermnamesAllResponse{}
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

func Removelinkrecord(w http.ResponseWriter, req *http.Request) {
	var data = entities.LinkRecordEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		success, _, msg := models.Removelinkrecord(&data)
		ThrowDocumentAllResponse(msg, w, success)
	}
}

func GetLinkRecorddetails(w http.ResponseWriter, req *http.Request) {
	var data = entities.LinkRecordEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetLinkRecordsByID(&data)
		ThrowLinkRecorddetailsAllResponse(msg, data, w, success)
	}
}

func ThrowLinkRecorddetailsAllResponse(successMessage string, responseData []entities.LinkRecordDetailsEntity, w http.ResponseWriter, success bool) {
	var response = entities.LinkRecordDetailsAllResponse{}
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

func Savelinkrecord(w http.ResponseWriter, req *http.Request) {
	var data = entities.LinkRecordEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.SaveLinkRecordsByID(&data)
		ThrowLinkRecordAllResponse(msg, data, w, success)
	}
}

func ThrowLinkRecordAllResponse(successMessage string, responseData int64, w http.ResponseWriter, success bool) {
	var response = entities.LinkRecordSaveAllResponse{}
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

func GetParentRecordInfo(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetParentRecordInfoByID(&data)
		ThrowParentRecordInfoAllResponse(msg, data, w, success)
	}
}

func ThrowParentRecordInfoAllResponse(successMessage string, responseData entities.ParentRecordInfoEntity, w http.ResponseWriter, success bool) {
	var response = entities.ParentRecordInfoAllResponse{}
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


func GetSLADuetimeCalculation(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		success, _, msg := models.GetSLADuetimeCalculation(&data)
		ThrowDocumentAllResponse(msg, w, success)
	}
}



func GetParentrecordForIM(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.GetParentrecordForIM(&data)
		logger.Log.Println(data)
		ThrowParentrecordAllResponse(msg, data, w, success)
	}
}


func UpdateVendorTickeID(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordcommonEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		//models.APIcall()
		data, success, _, msg := models.UpdateVendorTickeID(&data)
		ThrowRecordcommonIntResponse(msg, data, w, success)
	}
}
func UpdateSlaDueTime(w http.ResponseWriter, req *http.Request) {
	var data = entities.RecordNoEntity{}
	jsonError := data.FromJSON(req.Body)

	if jsonError != nil {
		logger.Log.Println(jsonError)
		entities.ThrowJSONResponse(entities.JSONParseErrorResponse(), w)
	} else {
		//responseError = validators.ValidateAddMstClient(data)

		//if(len(responseError)==0){
		data, success, _, msg := models.UpdateSlaDueTime(&data)
		ThrowRecordstatusIntResponse(msg, data, w, success)
	}
}
