package models

import (
	"bytes"
	"encoding/json"
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/dbconfig"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"io/ioutil"
	"log"
	"net/http"
)

func SaveParentRecord(page *entities.RecordDetailsRequestEntity) (int64, bool, error, string) {
	logger.Log.Println("In side SaveChildRecord")
	var insertedID int64
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return insertedID, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return insertedID, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	if page.ParentID == page.ChildID {
		return insertedID, false, err, "Ticket should NOT be attached as a child ticket to itself."
	}

	if page.ParentID == 0 {
		return insertedID, false, err, "Parent should not be null"
	}

	duprec, err2 := dataAccess.CheckChildRecordIsExist(page)
	if err2 != nil {
		return insertedID, false, err2, "Something Went Wrong"
	}
	parentaschild, err2 := dataAccess.CheckParentAsChildRecord1(page)
	if err2 != nil {
		return insertedID, false, err2, "Something Went Wrong"
	}
	if duprec > 0 {
		return insertedID, false, err2, "This child record is already present."
	} else if parentaschild > 0 {
		return insertedID, false, err2, "Parent record cannot be attached as Child to another parent"
	} else {
		parentrecordCode, err1 := dataAccess.GetRecordIDByID(page.ParentID)
		if err1 != nil {
			return insertedID, false, err1, "Something Went Wrong"
		}
		childrecordCode, err1 := dataAccess.GetRecordIDByID(page.ChildID)
		if err1 != nil {
			return insertedID, false, err1, "Something Went Wrong"
		}
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Println("Transaction creation error in Add child Record", err)
			return insertedID, false, err, "Something Went Wrong"
		}
		values, err1 := dao.SaveParentRecord(tx, page)
		if err1 != nil {
			tx.Rollback()
			//db.Close()
			return insertedID, false, err1, "Something Went Wrong"
		}

		if values > 0 {
			var logData = "(Parent Ticket ID:" + parentrecordCode + ", Child Ticket ID:" + childrecordCode + ")"
			err = dao.InsertActivityLogs(tx, page.Clientid, page.Mstorgnhirarchyid, page.ParentID, 15, logData, page.Userid, page.GroupID)
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				//	db.Close()
				return insertedID, false, err, "Something Went Wrong"
			}
			page.ChildIDS = append(page.ChildIDS, page.ChildID)
			reqbd := &entities.ParentchildEntity{}
			reqbd.Parentid = page.ParentID
			reqbd.Childids = page.ChildIDS
			reqbd.Userid = page.Userid
			reqbd.Createdgroupid = page.GroupID
			reqbd.IsAttaching = 1
			postBody, _ := json.Marshal(reqbd)

			logger.Log.Println("Record status request body -->", reqbd)

			responseBody := bytes.NewBuffer(postBody)
			resp, err := http.Post(dbconfig.MASTER_URL+"/updatechildstatus", "application/json", responseBody)
			if err != nil {
				logger.Log.Println("An Error Occured --->", err)
				tx.Rollback()
				//	db.Close()
				return insertedID, false, err, "SomethinIDg Went Wrong"
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Log.Println("response body ------> ", err)
				tx.Rollback()
				//	db.Close()
				return insertedID, false, err, "Something Went Wrong"
			}
			sb := string(body)
			wfres := entities.WorkflowResponse{}
			json.Unmarshal([]byte(sb), &wfres)

			logger.Log.Println("Record status response message -->", wfres.Success)
			logger.Log.Println("Record status response error message -->", wfres.Message)

			stateID, err := dataAccess.GetrecordlateststateID(page.ParentID, page.Clientid, page.Mstorgnhirarchyid)
			if err != nil {
				log.Println("Find child record status error", err)
				tx.Rollback()
				//	db.Close()
				return insertedID, false, err, "Something Went Wrong"
			}
			_, err = Childrecordstatusupdation(tx, page.Clientid, page.Mstorgnhirarchyid, page.ChildID, stateID, page.Userid, page.GroupID, db)
			if err != nil {
				log.Println("Child Record status updation failed", err)
				tx.Rollback()
				//	db.Close()
				return insertedID, false, err, "Something Went Wrong"
			}

			err = tx.Commit()
			if err != nil {
				log.Println("DB commit is failed", err)
				tx.Rollback()
				//	db.Close()
				return insertedID, false, err, "Something Went Wrong"
			}
			//db.Close()
			return values, true, err, "Parent Record has been Added"
		}

	}
	return insertedID, false, err, "Something Went Wrong"

}

func GetParentRecordDetailsByNo(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetParentRecordDetailsByNo")
	t := []entities.RecordDetailsEntity{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	checkparent, err := dataAccess.CheckParentOrNot(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	id, _, err := dataAccess.GetIdAgainstRecordNo(page.Clientid, page.Mstorgnhirarchyid, page.RecordNo)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	_, currentseqno, _, err := dataAccess.Getcurrentsatusid(page.Clientid, page.Mstorgnhirarchyid, id)
	if err != nil {
		logger.Log.Println(err)
		return t, false, err, "Something Went Wrong"
	}

	if checkparent == 0 {
		return t, false, err, "This Parent record can not to be attached with child record."
	} else if currentseqno == 1 || currentseqno == 10 {
		return t, false, err, "This Parent record can not to be attached with child record."
	} else {
		values, err := dataAccess.GetRecordDetailsByNo(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		for i, v := range values {
			page.RecordStageID = v.RecordStageID
			workflowid, err := dataAccess.GetProcessidbyworkingcatid(page)
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
			values[i].WorkFlowDetails = workflowid
		}
		return values, true, err, ""
	}

}

func ChildRecordSearchCriteria(page *entities.ChildRecordSearchEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetParentRecordDetailsByNo")
	t := []entities.RecordDetailsEntity{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	if len(page.RecordNo) > 0 {
		id, _, err := dataAccess.GetIdAgainstRecordNo(page.Clientid, page.Mstorgnhirarchyid, page.RecordNo)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}

		_, currentseqno, _, err := dataAccess.Getcurrentsatusid(page.Clientid, page.Mstorgnhirarchyid, id)
		if err != nil {
			logger.Log.Println(err)
			return t, false, err, "Something Went Wrong"
		}

		if currentseqno == 3 || currentseqno == 8 || currentseqno == 11 {
			return t, false, err, "You are not allowed to attach Cancelled, Resolved or Closed Ticket as Child"
		}

		recorddiffID, err := dataAccess.Getrecordtypediffid(id, page.Clientid, page.Mstorgnhirarchyid)
		if err != nil {
			logger.Log.Println(err)
			return t, false, err, "Something Went Wrong"
		}

		if recorddiffID != page.RecordDiffid {
			return t, false, err, "Ticket Type is not matched."
		}
	}

	values, err := dataAccess.GetRecordDetailsByOthers(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	for i, v := range values {
		page.RecordStageID = v.RecordStageID
		workflowid, err := dataAccess.GetProcessidbyworkingcatidforsearch(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		values[i].WorkFlowDetails = workflowid
	}
	return values, true, err, ""

}
