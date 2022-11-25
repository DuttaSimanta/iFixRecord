package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
)

func Recordadditionalfieldupdate(data *entities.RecordcategoryupdateEntity) (int64, bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// dbcon, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("Error in DBConnection in side Recordcategoryupdate")
	// 	return 0, false, err, "Something Went Wrong"
	// }
	//defer dbcon.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Println("Transaction creation error in Recordcategoryupdate", err)
		//dbcon.Close()
		return 0, false, err, "Something Went Wrong"
	}
	dataAccess := dao.DbConn{DB: db}
	for i := 0; i < len(data.Additionalfields); i++ {
		var a = data.Additionalfields[i]
		if len(a.Val) > 0 {
			oldvalue, err := dataAccess.Getadditionaloldvalue(a.Termsid, data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			err = dao.UpdateRecordAdditional(tx, data.ClientID, data.Mstorgnhirarchyid, a.ID, a.Termsid, a.Val, data.RecordID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			fildsname, err := dao.Gettermnamebyid(tx, a.Termsid, data.ClientID, data.Mstorgnhirarchyid)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			err = dao.InsertActivityLogs(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 5, fildsname+" has been changed from "+oldvalue+" to "+a.Val, data.UserID, data.UsergroupID)
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}
		}

	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Println("Error is --->", err)
		tx.Rollback()
		//dbcon.Close()
		return 0, false, err, "Something Went Wrong"
	}

	//dbcon.Close()
	return 0, true, err, "Additional fields updated successfully."
}
