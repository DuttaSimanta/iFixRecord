package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

func GetDynamicQueryResult(page map[string]interface{}) ([]interface{}, bool, error, string) {
	logger.Log.Println("In side GetDynamicQueryResult model")
	t := []interface{}{}
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

	values, err := dataAccess.GetQueryNParams(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	if len(values) == 0 {
		return t, false, err, "Query not configured!"
	}
	results, err := dataAccess.GetQueryResuls(page, values[0])
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	return results, true, err, ""
}
func GetDynamicQueryResultTemp(page map[string]interface{}) (entities.QueryDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetDynamicQueryResult model")
	t := entities.QueryDetailsEntity{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //defer db.Close()
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

	values, err := dataAccess.GetQueryNParams(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	if len(values) == 0 {
		return t, false, err, "Query not configured!"
	}
	timediff, _, err := dataAccess.Gettimezonediff(page["clientid"], page["mstorgnhirarchyid"])
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	results, err := dataAccess.GetQueryResulsTemp(page, values[0], timediff)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	page["querytype"] = 1
	values1, err := dataAccess.GetQueryNParams(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	if len(values1) == 0 {
		return t, false, err, "Query not configured!"
	}
	res, err := dataAccess.GetQueryResulsCountTemp(page, values1[0])
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	for i, v := range results {
		page["recordid"] = v.ID
		page["recordstageid"] = v.StageID
		cat, err := dataAccess.GetCategoryNameByRecordID(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		results[i].Categories = cat
	}
	t.Values = results
	t.Total = res.ID
	return t, true, err, ""
}

func GetDynamicQueryCountResultTemp(page map[string]interface{}) (int64, bool, error, string) {
	logger.Log.Println("In side GetDynamicQueryResult model")
	var t int64
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //defer db.Close()
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

	values, err := dataAccess.GetQueryNParams(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	if len(values) == 0 {
		return t, false, err, "Query not configured!"
	}
	results, err := dataAccess.GetQueryResulsCountTemp(page, values[0])
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	t = results.ID
	return t, true, err, ""
}

func RecordGridResult(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	var t map[string]interface{}
	var response map[string]interface{}
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
/*	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
*/
if dbSlave == nil {
		dbcon, err := ConnectMySqlDbSlave()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		dbSlave = dbcon
	}
	dataAccess := dao.DbConn{DB: dbSlave}

	// New Addition
	_, err1 := dataAccess.GetOrgnTypeForGrid(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	//if orgntype == 2 {
	getOrgnID, err1 := dataAccess.GetOrgnIDForGrid(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	logger.Log.Println("getOrgnID ------------->", getOrgnID)
	//}
	//End

	fConfig, err := dataAccess.GetQueryNParams(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	if len(fConfig) == 0 {
		return t, false, err, "Query not configured!"
	}

	if fConfig[0].QueryType == 2 {
		result, err := dataAccess.RecordGridResult(page, fConfig[0], getOrgnID)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		response = make(map[string]interface{}, 3)
		response["result"] = result
	} else {
		response = make(map[string]interface{}, 2)
	}
	total, err := dataAccess.RecordGridCount(page, fConfig[0], getOrgnID)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	response["querytype"] = fConfig[0].QueryType
	response["menuid"] = page["menuid"]
	response["total"] = total
	return response, true, nil, ""
}

func RecordGridResultOnly(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	var t map[string]interface{}
	var response map[string]interface{}
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
/*	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
*/
if dbSlave == nil {
		dbcon, err := ConnectMySqlDbSlave()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		dbSlave = dbcon
	}
	dataAccess := dao.DbConn{DB: dbSlave}
	response = make(map[string]interface{}, 2)
	result, err := dataAccess.RecordGridResultOnly(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	response["result"] = result

	total, err := dataAccess.RecordGridCountOnly(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	response["total"] = total
	return response, true, nil, ""
}

func RecordFilterAdd(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	var t map[string]interface{}
	var response map[string]interface{}
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

	countNmae, err := dataAccess.RecordFilterNameCheck(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	if countNmae > 0 {
		return t, false, err, "Filter Name is already exists."
	} else {
		response = make(map[string]interface{}, 1)
		inserid, err := dataAccess.RecordFilterAdd(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		response["id"] = inserid
	}

	return response, true, nil, ""
}

func RecordFilterList(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	var t map[string]interface{}
	var response map[string]interface{}
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
	response = make(map[string]interface{}, 2)
	filtercount, err := dataAccess.RecordFilterCount(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	response["total"] = filtercount

	result, err := dataAccess.RecordFilterList(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	response["result"] = result

	return response, true, nil, ""
}

func RecordFilterDelete(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	var t map[string]interface{}
	var response map[string]interface{}
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

	result, err := dataAccess.RecordFilterDelete(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	response = make(map[string]interface{}, 1)
	response["result"] = result
	return response, true, nil, ""
}



func RecordFilterUpdate(page map[string]interface{}) (map[string]interface{}, bool, error, string) {
	var t map[string]interface{}
	var response map[string]interface{}
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

	// countNmae, err := dataAccess.RecordFilterNameCheckForDuplicate(page)
	// if err != nil {
	//  return t, false, err, "Something Went Wrong"
	// }
	// if countNmae > 0 {
	//  return t, false, err, "Filter Name is already exists."
	// } else {
	response = make(map[string]interface{}, 1)
	inserid, err := dataAccess.RecordFilterUpdate(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	response["id"] = inserid
	// }

	return response, true, nil, ""
}
