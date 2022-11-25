package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
)

func GetMiscDataByRecordID(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	logger.Log.Println("In side GetMiscDataByRecordID model")
	var response map[string]interface{}
	t := []entities.RecorddifferentionSingle{}
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
	terms, err := dataAccess.GetFileTermByrecordTypeID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	status, err := dataAccess.GetFirstStatusByrecordTypeID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	priority, err := dataAccess.GetPriorityByrecordTypeNCatID(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	catvalues, err := GetRecorddifferentiationbyparent(page, t)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	workcat, err := dataAccess.GetWorkCatLabel(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	assetcount, err := dataAccess.GetAssetCountByRecordType(page)
	if err != nil {
		return response, false, err, "Something Went Wrong"
	}
	response = make(map[string]interface{}, 6)
	response["terms"] = terms
	response["status"] = status
	response["category"] = catvalues
	response["priority"] = priority
	response["workcat"] = workcat
	response["assetcount"] = assetcount
	return response, true, err, ""
}
func GetRecorddifferentiationbyparent(page map[string]interface{}, output []entities.RecorddifferentionSingle) ([]entities.RecorddifferentionSingle, error) {
	logger.Log.Println("In side Recorddifferentiationmodel")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return output, err
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return output, err
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	values, err1 := dataAccess.GetRecorddifferentiationbyparent(page)
	if err1 != nil {
		return output, err1
	}
	if len(values) > 0 {
		output = append(output, values[0])
		log.Print(output)
		if values[0].Parentid > 0 {
			page["recordcatdiffid"] = values[0].Parentid
			return GetRecorddifferentiationbyparent(page, output)
		} else {
			return output, err
		}
	} else {
		return output, err
	}
}
