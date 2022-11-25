package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
	"strconv"
)

func GetAssetDetailsByAssetID(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side GetAssetDetailsByAssetID model")
	var response map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	attributes, err := dataAccess.GetAssetAttributesByAseetID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}

	assetdata, err := dataAccess.GetOnlyAssetIDByID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	response = make(map[string]interface{}, 3)
	response["assetattributes"] = attributes
	response["total"] = len(assetdata)
	for _, v := range assetdata {
		page["id"] = v["id"]
		page["mstdifferentiationtypeid"] = v["mstdifferentiationtypeid"]
		assetattr, err := dataAccess.GetAssetAttributesByAssetID(page)
		if err != nil {
			return response, false, err, "Something Went Wrong"
		}
		v["attributes"] = assetattr
	}

	response["assetvales"] = assetdata
	return response, true, err, ""
}
func GetRecordAssetByID(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side GetRecordAssetByID model")
	var response map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	attributes, err := dataAccess.GetAssetAttributesByTypeID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	assetcount, err := dataAccess.CountAssetByRecordID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	assetdata, err := dataAccess.GetRecordAssetByID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	response = make(map[string]interface{}, 3)
	response["assetattributes"] = attributes
	response["total"] = assetcount
	for _, v := range assetdata {
		page["id"] = v["id"]
		assetattr, err := dataAccess.GetAssetAttributesByAssetID(page)
		if err != nil {
			return response, false, err, "Something Went Wrong"
		}
		v["attributes"] = assetattr
	}

	response["assetvales"] = assetdata
	return response, true, err, ""
}

func GetAllAssetTypeNDetailsByRecordID(page map[string]interface{}) ([]map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side GetAllAssetTypeNDetailsByRecordID model")
	var response []map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	assettype, err := dataAccess.GetAssetTypeIDByRecordID(page)
	for _, v := range assettype {
		page["mstdifferentiationtypeid"] = v["id"]
		attributes, err := dataAccess.GetAssetAttributesByTypeID(page)
		if err != nil {
			return response, false, err, "Something Went Wrong"
		}
		v["attributes"] = attributes
		page["mstdifferentiationtypeid"] = v["id"]
		assetdata, err := dataAccess.GetAllRecordAssetByID(page)
		if err != nil {
			return response, false, err, "Something Went Wrong"
		}
		for _, x := range assetdata {
			page["id"] = x["id"]
			assetattr, err := dataAccess.GetAssetAttributesByAssetID(page)
			if err != nil {
				return response, false, err, "Something Went Wrong"
			}
			x["attributes"] = assetattr
		}

		v["assets"] = assetdata

		response = append(response, v)
	}
	return response, true, err, ""
}
func GetAssetHistroyByAssetID(page *entities.FetchAssetHistoryRequest) (map[string]interface{}, bool, error, string) {
	logger.Log.Println("GetAssetHistroyByAssetID Req====>", page)
	var response = make(map[string]interface{})
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	assetHistory, err := dataAccess.GetAssetHistroyByAssetID(page)
	if err != nil {
		//tx.Rollback()
		logger.Log.Println(err)
		return response, false, err, "Something Went Wrong"
	}
	response = assetHistory

	return response, true, nil, ""
}
func InsertRecordAsset(page *entities.InsertRecordAssetEntity) (map[string]interface{}, bool, error, string) {
	var response = make(map[string]interface{})
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	var keys []int64
	var values []string
	for k, v := range page.AssetDetails {
		var key int64
		if n, err := strconv.Atoi(k); err == nil {
			//fmt.Println(n + 1)
			key = int64(n)
		}
		keys = append(keys, int64(key))
		values = append(values, v.(string))
		log.Println("Asset Keys===>", k)
		log.Println("AssetValues===>", v)
	}
	log.Println("Asset Keys===>", keys)

	dataAccess := dao.DbConn{DB: db}

	tx, err := db.Begin()
	if err != nil {
		logger.Log.Println("Transaction creation error in Add asset with Record", err)
		tx.Rollback()
		//db.Close()
		return response, false, err, "Something Went Wrong"
	}
	assetID, err := dataAccess.GetAssetIdToInsert(tx, page.Clientid, page.Mstorgnhirarchyid)
	if err != nil {
		tx.Rollback()
		logger.Log.Println(err)
		return response, false, err, "Something Went Wrong"
	}
	logger.Log.Println("Asset ID===>", assetID)
	var trnAsset entities.TrnAsset
	trnAsset.ClientId = page.Clientid
	trnAsset.MstOrgnHirarchyId = page.Mstorgnhirarchyid
	trnAsset.MstDifftypeid = page.MstDiffTypeID
	trnAsset.AssetId = assetID
	trnAsset.AdditionalAttr = ""
	trnAsset.ActiveFlag = 1
	trnAsset.DeleteFlag = 0
	lastInsertedTrnAssetId, insertTrnAssetError := dataAccess.InsertTrnAsset(tx, &trnAsset)
	if insertTrnAssetError != nil {
		tx.Rollback()
		logger.Log.Println(insertTrnAssetError)
		return response, false, err, "Something Went Wrong"
	}
	response["assetid"] = lastInsertedTrnAssetId
	page.Assetid = lastInsertedTrnAssetId
	for i := 0; i < len(keys); i++ {
		var mapAssetDiff entities.MapAssetDifferentiation

		//log.Printf("%s\n", text)
		mapAssetDiff.Clientid = page.Clientid
		mapAssetDiff.Mstorgnhirarchyid = page.Mstorgnhirarchyid
		mapAssetDiff.Mstdifferentiationtypeid = page.MstDiffTypeID
		mapAssetDiff.Mstdifferentiationid = keys[i]
		mapAssetDiff.Trnassetid = lastInsertedTrnAssetId
		mapAssetDiff.Value = values[i]
		mapAssetDiff.Deleteflg = 0
		mapAssetDiff.Activeflg = 1
		//	mapAssetDiff.AuditTransactionId = 1
		mapAssetDiffError := dataAccess.InsertMapAssetDiff(tx, &mapAssetDiff)
		if mapAssetDiffError != nil {
			logger.Log.Println(mapAssetDiffError)
			tx.Rollback()
			return response, false, err, "Something Went Wrong"
		}

		//coloumnCount++

	}
	page.Assetid = lastInsertedTrnAssetId
	response["assetcode"] = assetID
	_, err1 := dataAccess.InsertRecordAsset(tx, page, assetID)
	if err1 != nil {
		return response, false, err, "Something Went Wrong"
	}

	err = tx.Commit()
	if err != nil {
		log.Println("DB commit is failed", err)
		tx.Rollback()
		//db.Close()
		return response, false, err, "Something Went Wrong"
	}
	//db.Close()
	return response, true, nil, "Asset has been added and associated with record successfully."
}

func UpdateRecordAsset(page *entities.UpdateRecordAssetEntity) (map[string]interface{}, bool, error, string) {
	var response map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	tx, err := db.Begin()
	if err != nil {
		logger.Log.Println("Transaction creation error in Add asset with Record", err)
		tx.Rollback()
		//db.Close()
		return response, false, err, "Something Went Wrong"
	}
	util, err2 := dataAccess.Gettimediff(page.Clientid, page.Mstorgnhirarchyid)
	if err2 != nil {
		return response, false, err2, "Something Went Wrong"
	}
	_, err1 := dataAccess.UpdateRecordAsset(tx, page, util.Timediff)
	if err1 != nil {
		return response, false, err, "Something Went Wrong"
	}
	err = tx.Commit()
	if err != nil {
		log.Println("DB commit is failed", err)
		tx.Rollback()
		//db.Close()
		return response, false, err, "Something Went Wrong"
	}
	//db.Close()
	return response, true, nil, "Asset has been updated and associated with record successfully."
}

func AddAssetWithRecord(page *entities.RecordAssetRequestEntity) (map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side GetRecordAssetByID model")
	var response map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	recordCode, err1 := dataAccess.GetRecordIDByID(page.Recordid)
	if err1 != nil {
		return response, false, err1, "Something Went Wrong"
	}
	_, err1 = dataAccess.AddAssetWithRecord(page)
	if err1 != nil {
		return response, false, err, "Something Went Wrong"
	}
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Println("Transaction creation error in Add asset with Record", err)
		tx.Rollback()
		//db.Close()
		return response, false, err, "Something Went Wrong"
	}
	for _, v := range page.AssetID {
		assetIDS, err1 := dataAccess.GetAssetIDByID(v)
		if err1 != nil {
			return response, false, err1, "Something Went Wrong"
		}
		var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"
		err = dao.InsertActivityLogs(tx, page.Clientid, page.Mstorgnhirarchyid, page.Recordid, 7, logData, page.Userid, page.GroupID)
		if err != nil {
			log.Println("error is ----->", err)
			tx.Rollback()
			return response, false, err, "Something Went Wrong"
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println("DB commit is failed", err)
		tx.Rollback()
		//db.Close()
		return response, false, err, "Something Went Wrong"
	}
	//db.Close()
	return response, true, err, "Asset has been associated with record successfully."
}

func DeleteAssetFromRecord(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side DeleteAssetFromRecord model")
	var response map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	recordCode, err := dataAccess.GetRecordIDByIDInterface(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}

	assetIDS, err := dataAccess.GetAssetIDByIDInterface(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"
	_, err = dataAccess.DeleteAssetFromRecord(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	page["logval"] = logData
	_, err = dataAccess.InsertActivityLogForRecordAsset(page, 8)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	return response, true, err, "Asset has been disassociated with record successfully."
}

func GetAssetTypeByRecordID(page map[string]interface{}) ([]map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side GetAssetTypeByRecordID model")
	var response []map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	response, err = dataAccess.GetAssetTypeByRecordID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	if response == nil {
		response = make([]map[string]interface{}, 0)
		return response, true, err, ""
	} else {
		return response, true, err, ""
	}

}

func GetAssetfieldSpecificDataBYRecordID(page *entities.AssetAttrNameValRequestEntity) ([]map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side GetAssetfieldSpecificDataBYRecordID model")
	var response []map[string]interface{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return response, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return response, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	results, err := dataAccess.GetAssetfieldSpecificDataBYRecordID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	if len(results) > 0 {
		var tmpastId int64
		tmpRow := make(map[string]interface{}, len(page.AssetFieldsNames)+1)
		for _, v := range results {
			if tmpastId != v.Id {
				if tmpastId > 0 {
					response = append(response, tmpRow)
					tmpRow = make(map[string]interface{}, len(page.AssetFieldsNames)+1)
				}
				tmpastId = v.Id
			}
			tmpRow["id"] = tmpastId
			tmpRow[v.Name] = v.Value
		}
		response = append(response, tmpRow)
	}

	if len(response) == 0 {
		response = make([]map[string]interface{}, 0)
		return response, true, err, ""
	} else {
		return response, true, err, ""
	}

}
