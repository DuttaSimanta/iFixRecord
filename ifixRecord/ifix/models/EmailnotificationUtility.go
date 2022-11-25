package models

import (
	"bytes"
	"encoding/json"
	"ifixRecord/ifix/dbconfig"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"net/http"
)

func FollowupCountEmail(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, Count int64) {
	logger.Log.Println("Follow up count is   ---------------------------->", Count)
	reqbd := &entities.EmailnotificationTemplate{}
	reqbd.Clientid = ClientID
	reqbd.Mstorgnhirarchyid = Mstorgnhirarchyid
	reqbd.Recordid = RecordID
	reqbd.Channeltype = 1
	reqbd.Eventnotificationid = 5
	reqbd.Followupcount = Count

	postBody, _ := json.Marshal(reqbd)

	logger.Log.Println("Record priority request body -->", reqbd)

	responseBody := bytes.NewBuffer(postBody)
	http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// resp, err := http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// if err != nil {
	// 	logger.Log.Println("An Error Occured --->", err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Log.Println("response body ------> ", err)
	// }
	// sb := string(body)
	// logger.Log.Println("response body ------> ", sb)
}

func PriorityChangeEmail(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, Recorddiffid int64) {
	reqbd := &entities.EmailnotificationTemplate{}
	reqbd.Clientid = ClientID
	reqbd.Mstorgnhirarchyid = Mstorgnhirarchyid
	reqbd.Recordid = RecordID
	reqbd.Eventnotificationid = 7
	reqbd.Channeltype = 1
	reqbd.Priorityid = Recorddiffid

	postBody, _ := json.Marshal(reqbd)

	logger.Log.Println("Record priority request body -->", reqbd)

	responseBody := bytes.NewBuffer(postBody)
	http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// resp, err := http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// if err != nil {
	// 	logger.Log.Println("An Error Occured --->", err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Log.Println("response body ------> ", err)
	// }
	// sb := string(body)
	// logger.Log.Println("response body ------> ", sb)
}

func StatusChangeEmail(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, StatusID int64,Statusseq int64) {

	reqbd := &entities.EmailnotificationTemplate{}
	reqbd.Clientid = ClientID
	reqbd.Mstorgnhirarchyid = Mstorgnhirarchyid
	reqbd.Recordid = RecordID
	reqbd.Eventnotificationid = 1
	reqbd.Channeltype = 1
	reqbd.Statusid = StatusID
	/*if StatusID == 18 {
		reqbd.TermsSeq = 78
	} else if StatusID == 19 {
		reqbd.TermsSeq = 7
	} else if StatusID == 20 {
		reqbd.TermsSeq = 9
	} else if StatusID == 21 {
		reqbd.TermsSeq = 10
	} else {
		reqbd.TermsSeq = 3
	}*/


	if Statusseq == 4 { //Pending for Vendor Action
		reqbd.TermsSeq = 78
	} else if Statusseq == 5 { //Pending User Action
		reqbd.TermsSeq = 12
	} else if Statusseq == 6 { //Pending Hold by Client
		reqbd.TermsSeq = 83
	} else if Statusseq == 7 { //Pending for Testing / Observation
		reqbd.TermsSeq = 84
	} else {
		reqbd.TermsSeq = 3
	}

	postBody, _ := json.Marshal(reqbd)

	logger.Log.Println("Record Status request body -->", reqbd)
	logger.Log.Println("11111111111")
	responseBody := bytes.NewBuffer(postBody)
	logger.Log.Println("22222222222")
	http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	logger.Log.Println("333333333333")

	//resp, err := http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// if err != nil {
	// 	logger.Log.Println("An Error Occured --->", err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Log.Println("response body ------> ", err)
	// }
	// sb := string(body)
	// logger.Log.Println("response body ------> ", sb)

}

func CustomerVisibleWorkNotesEmail(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) {
	reqbd := &entities.EmailnotificationTemplate{}
	reqbd.Clientid = ClientID
	reqbd.Mstorgnhirarchyid = Mstorgnhirarchyid
	reqbd.Recordid = RecordID
	reqbd.Eventnotificationid = 2
	reqbd.Channeltype = 1

	postBody, _ := json.Marshal(reqbd)

	logger.Log.Println("Record priority request body -->", reqbd)

	responseBody := bytes.NewBuffer(postBody)
	http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// resp, err := http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// if err != nil {
	// 	logger.Log.Println("An Error Occured --->", err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Log.Println("response body ------> ", err)
	// }
	// sb := string(body)
	// logger.Log.Println("response body ------> ", sb)
}

func FileAttachmentEmail(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) {
	reqbd := &entities.EmailnotificationTemplate{}
	reqbd.Clientid = ClientID
	reqbd.Mstorgnhirarchyid = Mstorgnhirarchyid
	reqbd.Recordid = RecordID
	reqbd.Eventnotificationid = 3
	reqbd.Channeltype = 1

	postBody, _ := json.Marshal(reqbd)

	logger.Log.Println("Record priority request body -->", reqbd)

	responseBody := bytes.NewBuffer(postBody)
	http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)

	// resp, err := http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// if err != nil {
	// 	logger.Log.Println("An Error Occured --->", err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Log.Println("response body ------> ", err)
	// }
	// sb := string(body)
	// logger.Log.Println("response body ------> ", sb)
}

func FileDeleteEmail(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) {
	reqbd := &entities.EmailnotificationTemplate{}
	reqbd.Clientid = ClientID
	reqbd.Mstorgnhirarchyid = Mstorgnhirarchyid
	reqbd.Recordid = RecordID
	reqbd.Eventnotificationid = 4
	reqbd.Channeltype = 1

	postBody, _ := json.Marshal(reqbd)

	logger.Log.Println("Record priority request body -->", reqbd)

	responseBody := bytes.NewBuffer(postBody)
	http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// resp, err := http.Post(dbconfig.NOTIFICATION_URL+"/sendnotification", "application/json", responseBody)
	// if err != nil {
	// 	logger.Log.Println("An Error Occured --->", err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.Log.Println("response body ------> ", err)
	// }
	// sb := string(body)
	// logger.Log.Println("response body ------> ", sb)
}
