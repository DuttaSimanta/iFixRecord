package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

func GetAdditionalFields(page *entities.AdditionalfieldRequestEntity) ([]entities.AdditionalFieldEntity, bool, error, string) {
	logger.Log.Println("In side GetAdditionalFields")
	t := []entities.AdditionalFieldEntity{}
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
	values, err := dataAccess.GetAdditionalFields(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetAdditionalFieldsBYTypeCat(page *entities.AdditionalfieldRequestEntity) ([]entities.AdditionalFieldEntity, bool, error, string) {
	logger.Log.Println("In side GetAdditionalFields")
	t := []entities.AdditionalFieldEntity{}
	var err error
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
	if len(page.RecordCatSet) > 0 {
		fieldIds := []int64{}
		res := []entities.AdditionalFieldEntity{}
		for _, v := range page.RecordCatSet {
			var data = &entities.AdditionalfieldRequestEntity{}
			data.Clientid = page.Clientid
			data.Mstorgnhirarchyid = page.Mstorgnhirarchyid
			var diff = entities.AdditionalFieldDiffEntity{}
			diff.Mstdifferentiationtypeid = page.RecordTypeDiffTypeID
			diff.Mstdifferentiationid = page.RecordTypeDiffID
			data.Mstdifferentiationset = append(data.Mstdifferentiationset, diff)
			data.Mstdifferentiationset = append(data.Mstdifferentiationset, v)
			values, err := dataAccess.GetAdditionalFields(data)
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
			res = append(res, values...)
		}
		resp := []entities.AdditionalFieldEntity{}
		for _, vv := range res {
			if itemExists(fieldIds, vv.FieldID) == false {
				resp = append(resp, vv)
				fieldIds = append(fieldIds, vv.FieldID)
			}
		}
		return resp, true, err, ""
	} else {
		return t, false, err, "Something Went Wrong"
	}
}
func itemExists(arr []int64, item int64) bool {

	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return true
		}
	}

	return false
}
