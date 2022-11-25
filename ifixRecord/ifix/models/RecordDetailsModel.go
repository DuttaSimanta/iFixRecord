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

func SaveChildRecord(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side SaveChildRecord")
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
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	if len(page.ChildIDS) > 0 {
		var itselfCheck int64
		for _, v := range page.ChildIDS {
			if page.ParentID == v {
				itselfCheck = 1
			}
		}
		if itselfCheck == 1 {
			return t, false, err, "Ticket should NOT be attached as a child ticket to itself."
		}

	} else {
		return t, false, err, "Child should not be null"
	}
	duprec, err2 := dataAccess.CheckDuplicateChildRecord(page)
	if err2 != nil {
		return t, false, err2, "Something Went Wrong"
	}
	parentaschild, err2 := dataAccess.CheckParentAsChildRecord(page)
	if err2 != nil {
		return t, false, err2, "Something Went Wrong"
	}
	childasparent, err2 := dataAccess.CheckChildAsParentRecord(page)
	if err2 != nil {
		return t, false, err2, "Something Went Wrong"
	}
	if duprec > 0 {
		return t, false, err2, "Child ticket can only have one parent"
	} else if childasparent > 0 {
		return t, false, err2, "Parent tickets cannot be attached as Child to another parent"
	} else if parentaschild > 0 {
		return t, false, err2, "Parent tickets cannot be attached as Child to another parent"
	} else {
		recordCode, err1 := dataAccess.GetRecordIDByID(page.ParentID)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}
		values, err1 := dataAccess.SaveChildRecord(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		if len(page.ChildIDS) > 0 {
			tx, err := db.Begin()
			if err != nil {
				logger.Log.Println("Transaction creation error in Add child Record", err)
				return t, false, err, "Something Went Wrong"
			}
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
				//db.Close()
				return t, false, err, "Something Went Wrong"
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Log.Println("response body ------> ", err)
				tx.Rollback()
				//db.Close()
				return t, false, err, "Something Went Wrong"
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
				//db.Close()
				return t, false, err, "Something Went Wrong"
			}
			for _, v := range page.ChildIDS {
				childCode, err1 := dataAccess.GetRecordIDByID(v)
				if err1 != nil {
					tx.Rollback()
					//db.Close()
					return t, false, err1, "Something Went Wrong"
				}
				var logData = "(Parent Ticket ID:" + recordCode + ", Child Ticket ID:" + childCode + ")"
				err = dao.InsertActivityLogs(tx, page.Clientid, page.Mstorgnhirarchyid, page.ParentID, 9, logData, page.Userid, page.GroupID)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					return t, false, err, "Something Went Wrong"
				}
				_, err = Childrecordstatusupdation(tx, page.Clientid, page.Mstorgnhirarchyid, v, stateID, page.Userid, page.GroupID, db)
				if err != nil {
					log.Println("Child Record status updation failed", err)
					tx.Rollback()
					//	db.Close()
					return t, false, err, "Something Went Wrong"
				}
			}
			err = tx.Commit()
			if err != nil {
				log.Println("DB commit is failed", err)
				tx.Rollback()
				//db.Close()
				return t, false, err, "Something Went Wrong"
			}
		}
		//	db.Close()
		return values, true, err, "Child Record has been Added"
	}

}

func RemoveChildRecord(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side RemoveChildRecord")
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
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	recordCode, err1 := dataAccess.GetRecordIDByID(page.ParentID)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	values, err1 := dataAccess.RemoveChildRecord(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Println("Transaction creation error in Add child Record", err)
		return t, false, err, "Something Went Wrong"
	}
	for _, v := range page.ChildIDS {
		childCode, err1 := dataAccess.GetRecordIDByID(v)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}
		var logData = "(Parent Ticket ID:" + recordCode + ", Child Ticket ID:" + childCode + ")"
		err = dao.InsertActivityLogs(tx, page.Clientid, page.Mstorgnhirarchyid, page.ParentID, 11, logData, page.Userid, page.GroupID)
		if err != nil {
			log.Println("error is ----->", err)
			tx.Rollback()
			return t, false, err, "Something Went Wrong"
		}
		err = dao.InsertActivityLogs(tx, page.Clientid, page.Mstorgnhirarchyid, v, 11, logData, page.Userid, page.GroupID)
		if err != nil {
			log.Println("error is ----->", err)
			tx.Rollback()
			return t, false, err, "Something Went Wrong"
		}
	}
	// =================   08.05.2022 ===========================================
	reqbd := &entities.ParentchildEntity{}
	reqbd.Transactionid = page.ChildIDS[0]
	reqbd.Userid = page.Userid
	reqbd.Createdgroupid = page.GroupID
	reqbd.Usergroupid = page.GroupID
	postBody, _ := json.Marshal(reqbd)

	logger.Log.Println("Record status request body -->", reqbd)

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(dbconfig.MASTER_URL+"/detachchildticket", "application/json", responseBody)
	if err != nil {
		logger.Log.Println("An Error Occured --->", err)
		tx.Rollback()
		return t, false, err, "Something Went Wrong"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Println("response body ------> ", err)
		tx.Rollback()
		return t, false, err, "Something Went Wrong"
	}
	sb := string(body)
	wfres := entities.WorkflowResponse{}
	json.Unmarshal([]byte(sb), &wfres)

	logger.Log.Println("Record status response message -->", wfres.Success)
	logger.Log.Println("Record status response error message -->", wfres.Message)
	// =================  08.03.2022 ============================================
	err = tx.Commit()
	if err != nil {
		log.Println("DB commit is failed", err)
		tx.Rollback()
		return t, false, err, "Something Went Wrong"
	}
	//defer db.Close()
	return values, true, err, "Child Record has been removed"

}

/*func GetChildRecordsBYParentID(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetChildRecordsBYParentID")
	t := []entities.RecordDetailsEntity{}
	db, err := dbconfig.ConnectMySqlDb()
	if err != nil {
		logger.Log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}
	defer db.Close()
	dataAccess := dao.DbConn{DB: db}

	values, err := dataAccess.GetChildRecordsBYParentID(page)
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
}*/

func GetChildRecordsBYParentID(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetChildRecordsBYParentID")
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
	typeseq, err := dataAccess.Getrecordtypeseq(page.Clientid, page.Mstorgnhirarchyid, page.ParentID)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	if typeseq == 1 {
		t, err = dataAccess.GetChildRecordsBYParentID(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
	} else {
		t, err = dataAccess.GetStaskChildRecordsBYParentID(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
	}

	for i, v := range t {
		page.RecordStageID = v.RecordStageID
		workflowid, err := dataAccess.GetProcessidbyworkingcatid(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		t[i].WorkFlowDetails = workflowid
	}
	return t, true, err, ""
}

func GetRecordDetails(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordDetails")
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

	orgnID, err := dataAccess.GetOrgnIDbyrecordID(page.Recordid)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	page.Mstorgnhirarchyid = orgnID
	values, err := dataAccess.GetRecordDetails(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	for i, v := range values {
		page.Recordid = v.ID
		page.RecordStageID = v.RecordStageID
		page.TermsSeq = 24
		respCount, err := dataAccess.GetResolRespBreachCount(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		if respCount > 0 {
			values[i].IsRespBreach = true
		} else {
			values[i].IsRespBreach = false
		}

		page.TermsSeq = 25
		reslCount, err := dataAccess.GetResolRespBreachCount(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		if reslCount > 0 {
			values[i].IsReslBreach = true
		} else {
			values[i].IsReslBreach = false
		}
		workflowid, err := dataAccess.GetProcessidbyworkingcatid(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		values[i].WorkFlowDetails = workflowid
		// Access Denied logic
		checkflg, err := dataAccess.CheckOrgnIDbyuserID(page.Clientid, orgnID, page.Userid)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		if len(checkflg) > 0 {
			values[i].Haspermission = true
		} else {
			values[i].Haspermission = false
		}
		// Access Denied logic
	}
	return values, true, err, ""
}

func GetRecordDetailsByNo(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordDetailsByNo")
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
	orgnID, err := dataAccess.GetOrgnIDbyrecordCode(page.RecordNo)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	checkflg, err := dataAccess.CheckOrgnIDbyuserID(page.Clientid, orgnID, page.Userid)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	//if len(checkflg) > 0 {
	values, err := dataAccess.GetRecordDetailsByNo(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	for i, v := range values {
		page.RecordStageID = v.RecordStageID
		page.Recordid = v.ID
		//workflowid, err := dataAccess.GetProcessidbyworkingcatidNew(page, values[0].Mstorgnhirarchyid)
		workflowid, err := dataAccess.GetProcessidbyworkingcatidNew(page, values[0].Mstorgnhirarchyid)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		values[i].WorkFlowDetails = workflowid
		if len(checkflg) > 0 {
			values[i].Haspermission = true
		} else {
			values[i].Haspermission = false
		}
	}

	return values, true, err, ""
	// } else {
	// 	return t, false, err, "Access Denied.."
	// }

}

func GetRecordDetailsByNoForlinkrecord(page *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordDetailsByNo")
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

	values, err := dataAccess.GetRecordDetailsByNoForlinkrecord(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	for i, v := range values {
		page.RecordStageID = v.RecordStageID
		page.Recordid = v.ID
		//workflowid, err := dataAccess.GetProcessidbyworkingcatidNew(page, values[0].Mstorgnhirarchyid)
		workflowid, err := dataAccess.GetProcessidbyworkingcatidNew(page, values[0].Mstorgnhirarchyid)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		values[i].WorkFlowDetails = workflowid
	}
	return values, true, err, ""
}

func GetRecordCatDetails(page *entities.RecordDetailsRequestEntity) (entities.RecordCatDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordCatDetails Model")
	t := entities.RecordCatDetailsEntity{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	//	defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	catvalues, err := dataAccess.GetAllTypeWiseCategories(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	dirvalues, err := dataAccess.GetRecordDirection(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	fieldvalues, err := dataAccess.GetAllRecordFields(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	workflowid, err := dataAccess.GetProcessidbyworkingcatid(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	assetchkval, err := dataAccess.CheckAssetCountForDetails(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	t.AssetAttached = assetchkval.AssetAttached
	reqStat := &entities.RecordcreaterequestEntity{
		Clientid:           0,
		Mstorgnhirarchyid:  0,
		Recorddifftypeid:   0,
		Recorddiffid:       0,
		Recorddiffparentid: 0,
		Recordtypeid:       0,
		Recordimpactid:     0,
		Recordurgencyid:    0,
		Recordcatid:        0,
	}
	reqStat.Clientid = page.Clientid
	reqStat.Mstorgnhirarchyid = page.Mstorgnhirarchyid
	reqStat.Recorddifftypeid = page.RecordDiffTypeid
	reqStat.Recorddiffid = page.RecordDiffid
	statusvalues, err := dataAccess.GetRecordstatusdata(reqStat)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	t.RecordCreateStatus = statusvalues
	if dirvalues == 1 {
		impactvalues, err := dataAccess.GetRecordWiseImpact(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		urgencyvalues, err := dataAccess.GetRecordWiseUrgency(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		if len(impactvalues) > 0 {
			t.Recordimpact = impactvalues
		} else {
			t.Recordimpact = []entities.RecordcatchildEntity{}
		}
		if len(urgencyvalues) > 0 {
			t.Recordurgency = urgencyvalues
		} else {
			t.Recordurgency = []entities.RecordcatchildEntity{}
		}

	} else {
		t.Recordimpact = []entities.RecordcatchildEntity{}
		t.Recordurgency = []entities.RecordcatchildEntity{}
	}

	//t.WorkingCatLabelID = workcatval.WorkingCatLabelID
	//t.AssetAttached = assetchkval.AssetAttached
	//t.Recordstatus = statusvalues
	t.Businessmatrixdirection = dirvalues
	t.RecordFields = fieldvalues
	t.WorkFlowDetails = workflowid

	//t.Recordterms = termslist
	t.Recordcatpos = -1
	if len(catvalues) > 0 {
		var tempCat int64
		var tempPos int64
		var seqNos int64
		var parentID int64
		var allowFlag int64
		parentID = 0
		allowFlag = 0
		catDiffTypes := []entities.RecordcatEntity{}
		catDiffs := []entities.RecordcatEntity{}
		catDiffType := entities.RecordcatEntity{}
		for k, v := range catvalues {

			if tempCat != v.Typeid {
				if k != 0 {
					allowFlag = parentID
					seqNos++
					catDiffType.Child = catDiffs
					if len(catDiffs) > 0 && seqNos != 1 {
						logger.Log.Println("seq No ==>", seqNos)
						logger.Log.Println("catDiffs ==>", catDiffs)
						catDiffType.IsDisabled = false
					} else {
						catDiffType.IsDisabled = true
					}
					if catDiffType.IsDisabled == false && t.Recordcatpos == -1 {
						t.Recordcatpos = tempPos
					}
					catDiffType.Sequanceno = seqNos
					catDiffTypes = append(catDiffTypes, catDiffType)
					catDiffs = []entities.RecordcatEntity{}
					tempPos++
				}

				tempCat = v.Typeid
				catDiffType.ID = v.Typeid
				catDiffType.Title = v.Typename
				catDiffType.Sequanceno = v.Typeseq
			}
			if parentID == v.ParentID || allowFlag == v.ParentID {
				log.Println(parentID, v.ID)
				catDiff := entities.RecordcatEntity{}
				catDiff.ID = v.ID
				catDiff.Title = v.Name
				catDiff.Sequanceno = v.Seqno
				catDiff.IsSelected = v.Selected
				if v.Selected == 1 {
					parentID = v.ID
				}
				//logger.Log.Println("catDiff values --->", catDiff)
				catDiffs = append(catDiffs, catDiff)
				//logger.Log.Println("catDiffs values --->", catDiffs)
			}
		}
		seqNos++
		catDiffType.Sequanceno = seqNos
		catDiffType.Child = catDiffs
		catDiffTypes = append(catDiffTypes, catDiffType)
		/*for k, _ := range catDiffTypes {
			if k > int(t.Recordcatpos) {
				logger.Log.Println(k, int(t.Recordcatpos))
				catDiffTypes[k].Child = []entities.RecordcatEntity{}
			}
		}*/
		if len(catDiffTypes) > 0 {
			t.Recordcategory = catDiffTypes
		} else {
			t.Recordcategory = []entities.RecordcatEntity{}
		}
		for i := 0; i < len(t.Recordcategory); i++ {
			if t.Recordcategory[i].Sequanceno == 5 {
				for j := 0; j < len(t.Recordcategory[i].Child); j++ {
					if t.Recordcategory[i].Child[j].IsSelected == 1 {
						req := entities.RecordcreaterequestEntity{}
						req.Clientid = page.Clientid
						req.Mstorgnhirarchyid = page.Mstorgnhirarchyid
						req.Recordcatid = t.Recordcategory[i].Child[j].ID
						t.EstimatedEfforts, t.SlaCompliances, t.ChangeTypes, err = dataAccess.GetEstimateEffort(&req)
						if err != nil {
							return t, false, err, "Something Went Wrong"
						}
					}

				}

			}

		}
		logger.Log.Println("category value is in if --->", t)
		return t, true, err, ""
	} else {
		logger.Log.Println("category value is in else --->", t)
		return t, true, err, ""
	}

}

func GetRecordCatByLastID(page *entities.ChildRecordSearchEntity) (entities.RecordCatDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordCatByLastID Model")
	t := entities.RecordCatDetailsEntity{}
	finalCatDetails := []entities.RecordcatEntity{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	//	defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	dirRequest := &entities.RecordDetailsRequestEntity{}
	dirRequest.Clientid = page.Clientid
	dirRequest.Mstorgnhirarchyid = page.Mstorgnhirarchyid
	dirRequest.RecordDiffTypeid = page.RecordDiffTypeid
	dirRequest.RecordDiffid = page.RecordDiffid
	dirvalues, err := dataAccess.GetRecordDirection(dirRequest)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	req := entities.RecordcreaterequestEntity{}
	req.Clientid = page.Clientid
	req.Mstorgnhirarchyid = page.Mstorgnhirarchyid
	req.Recordcatid = page.CategoryID
	estimatedefforts, _, _, err := dataAccess.GetEstimateEffort(&req)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	for {
		catDetails, err := dataAccess.GetRecordCatByLastID(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		if len(catDetails) == 0 {
			break
		}
		var totsel int64
		tempCatDetails := []entities.RecordcatEntity{}
		for _, v := range catDetails {
			tempd := entities.RecordcatEntity{}
			tempd.ID = v.ID
			tempd.IsSelected = v.Selected
			tempd.Sequanceno = v.Seqno
			tempd.Title = v.Name
			totsel = totsel + v.Selected
			if tempd.IsSelected == 1 {
				page.CategoryID = v.ParentID
			}
			tempd.Child = []entities.RecordcatEntity{}
			tempCatDetails = append(tempCatDetails, tempd)
		}
		if totsel == 0 {
			break
		} else {
			tempdtype := entities.RecordcatEntity{}
			tempdtype.ID = catDetails[0].Typeid
			tempdtype.Title = catDetails[0].Typename
			tempdtype.Sequanceno = catDetails[0].Typeseq
			tempdtype.IsSelected = 1
			if len(tempCatDetails) > 1 {
				tempdtype.IsDisabled = false
			} else {
				tempdtype.IsDisabled = false
			}
			tempdtype.Child = tempCatDetails
			finalCatDetails = append(finalCatDetails, entities.RecordcatEntity{})
			copy(finalCatDetails[1:], finalCatDetails)
			finalCatDetails[0] = tempdtype
		}

	}
	for i, _ := range finalCatDetails {
		finalCatDetails[i].Sequanceno = int64(i) + 1
	}
	t.Businessmatrixdirection = dirvalues
	t.EstimatedEfforts = estimatedefforts
	t.Recordcategory = finalCatDetails
	return t, true, err, ""
}

func GetRecordAccessPermissionByNo(page *entities.RecordDetailsRequestEntity) (entities.RecordAccessEntity, bool, error, string) {
	t := entities.RecordAccessEntity{}
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
	orgnID, err := dataAccess.GetOrgnIDbyrecordCode(page.RecordNo)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	checkflg, err := dataAccess.CheckOrgnIDbyuserID(page.Clientid, orgnID, page.Userid)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	t.Clientid = page.Clientid
	t.Mstorgnhirarchyid = orgnID
	t.RecordNo = page.RecordNo
	if len(checkflg) > 0 {
		t.Haspermission = true
	} else {
		t.Haspermission = false
	}
	return t, true, err, ""

}
