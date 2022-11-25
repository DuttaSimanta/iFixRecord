package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

/*func GetRecordcreatedata(page *entities.RecordcreaterequestEntity) (entities.RecordcreateEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordcreatedatamodel")
	t := entities.RecordcreateEntity{}
	if mutexutility.MutexLocked(lock) == false {
		lock.Lock()
		defer lock.Unlock()
	}
	db, err := dbconfig.ConnectMySqlDb()
	//defer db.Close()
	if err != nil {
		logger.Log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}
	dataAccess := dao.DbConn{DB: db}
	addReq := entities.AdditionalfieldRequestEntity{}
	addReq.Clientid = page.Clientid
	addReq.Mstorgnhirarchyid = page.Mstorgnhirarchyid
	aa := entities.AdditionalFieldDiffEntity{}
	aa.Mstdifferentiationid = page.Recorddiffid
	aa.Mstdifferentiationtypeid = page.Recorddifftypeid
	bb := []entities.AdditionalFieldDiffEntity{}
	bb = append(bb, aa)
	catvalues, err := dataAccess.GetAllRecordcategories(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	dirvalues, err := dataAccess.Getdirection(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	termslist, err := dataAccess.GetAllRecordtermslist(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	statusvalues, err := dataAccess.GetRecordstatusdata(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	assetchkval, err := dataAccess.CheckAssetCount(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	workcatval, err := dataAccess.GetWorkingCatLabel(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	logger.Log.Println("############################################################################dirvalues --->", dirvalues)
	if dirvalues == 1 {
		impactvalues, err := dataAccess.GetRecordimpact(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		urgencyvalues, err := dataAccess.GetRecordurgency(page)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		t.Recordimpact = impactvalues
		t.Recordurgency = urgencyvalues
	}
		logger.Log.Println("############################################################################dirvalues --->", dirvalues)
	t.WorkingCatLabelID = workcatval.WorkingCatLabelID
	t.AssetAttached = assetchkval.AssetAttached
	t.Recordstatus = statusvalues
	t.Businessmatrixdirection = dirvalues
	t.Recordterms = termslist
	t.Recordcatpos = -1
	if len(catvalues) > 0 {
		var tempCat int64
		var tempPos int64
		var seqNos int64
		catDiffTypes := []entities.RecordcategorydetailsEntity{}
		catDiffs := []entities.RecordcategorydetailsEntity{}
		catDiffType := entities.RecordcategorydetailsEntity{}
		for k, v := range catvalues {
			//logger.Log.Println("tempCat ---------------------------------->", tempCat)
			//logger.Log.Println("v.Typeid ------------------------------->", v.Typeid)
			if tempCat != v.Typeid {
				if k != 0 {
					seqNos++
					catDiffType.Child = catDiffs
					if len(catDiffs) == 1 {
						catDiffType.IsDisabled = true
					} else {
						catDiffType.IsDisabled = false
					}

					if catDiffType.IsDisabled == false && t.Recordcatpos == -1 {
						t.Recordcatpos = tempPos
					}
					catDiffType.Sequanceno = seqNos
					catDiffTypes = append(catDiffTypes, catDiffType)
					catDiffs = []entities.RecordcategorydetailsEntity{}
					tempPos++
				}

				tempCat = v.Typeid
				catDiffType.ID = v.Typeid
				catDiffType.Title = v.Typename
				catDiffType.Sequanceno = v.Typeseq
			}
			catDiff := entities.RecordcategorydetailsEntity{}
			catDiff.ID = v.ID
			catDiff.Title = v.Name
			catDiff.Sequanceno = v.Seqno
			//logger.Log.Println("catDiff values --->", catDiff)
			catDiffs = append(catDiffs, catDiff)
			//logger.Log.Println("catDiffs values --->", catDiffs)
		}
		seqNos++
		if len(catDiffs) <= 1 {
			catDiffType.IsDisabled = true
		} else {
			catDiffType.IsDisabled = false
		}
		if !catDiffType.IsDisabled && t.Recordcatpos == -1 {
			t.Recordcatpos = tempPos
		}
		catDiffType.Sequanceno = seqNos
		catDiffType.Child = catDiffs
		catDiffTypes = append(catDiffTypes, catDiffType)
		//logger.Log.Println("catDiffTypes values ------------------->", catDiffTypes)
		if t.Recordcatpos == -1 {
			t.Recordcatpos = 0
		} else {
			for k, _ := range catDiffTypes {
				if k > int(t.Recordcatpos) {
					//logger.Log.Println(k, int(t.Recordcatpos))
					catDiffTypes[k].Child = []entities.RecordcategorydetailsEntity{}
				}
			}
		}
		for _, vv := range catDiffTypes {
			if vv.IsDisabled == true && len(vv.Child) > 0 {
				aa := entities.AdditionalFieldDiffEntity{}
				aa.Mstdifferentiationid = vv.Child[0].ID
				aa.Mstdifferentiationtypeid = vv.ID
				bb = append(bb, aa)
			}
		}
		addReq.Mstdifferentiationset = bb
		addfields, err := dataAccess.GetAdditionalFieldsByDiffId(&addReq)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		t.AdditionalFields = addfields
		t.Recordcategory = catDiffTypes

		//logger.Log.Println("category value is in if --->", t)
		return t, true, err, ""
	} else {
		//logger.Log.Println("category value is in else --->", t)
		return t, true, err, ""
	}

}*/

func GetAdditionalInfoBasedonCategory(page *entities.RecordcreaterequestEntity) (entities.RecordcatchildNEstimatedEfforEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordprioritydata")
	t := entities.RecordcatchildNEstimatedEfforEntity{}
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

	// typevalues, err := dataAccess.GetRecordprioritydata(page)

	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }
	estimatedefforts, slaCompliances, changeTypes, err1 := dataAccess.GetEstimateEffort(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	//t.ChildCategories = typevalues
	t.EstimatedEfforts = estimatedefforts
	t.SlaCompliance = slaCompliances
	t.ChangeType = changeTypes
	return t, true, err, ""
}

func GetRecordcreatedata(page *entities.RecordcreaterequestEntity) (entities.RecordcreateEntity, bool, error, string) {

	logger.Log.Println("In side GetRecordcreatedatamodel")

	t := entities.RecordcreateEntity{}

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

	addReq := entities.AdditionalfieldRequestEntity{}

	addReq.Clientid = page.Clientid

	addReq.Mstorgnhirarchyid = page.Mstorgnhirarchyid

	aa := entities.AdditionalFieldDiffEntity{}

	aa.Mstdifferentiationid = page.Recorddiffid

	aa.Mstdifferentiationtypeid = page.Recorddifftypeid

	bb := []entities.AdditionalFieldDiffEntity{}

	bb = append(bb, aa)

	catvalues, err := dataAccess.GetAllRecordcategories(page)

	if err != nil {

		return t, false, err, "Something Went Wrong"

	}

	dirvalues, err := dataAccess.Getdirection(page)

	if err != nil {

		return t, false, err, "Something Went Wrong"

	}

	termslist, err := dataAccess.GetAllRecordtermslist(page)

	if err != nil {

		return t, false, err, "Something Went Wrong"

	}

	statusvalues, err := dataAccess.GetRecordstatusdata(page)

	if err != nil {

		return t, false, err, "Something Went Wrong"

	}

	assetchkval, err := dataAccess.CheckAssetCount(page)

	if err != nil {

		return t, false, err, "Something Went Wrong"

	}

	workcatval, err := dataAccess.GetWorkingCatLabel(page)

	if err != nil {

		return t, false, err, "Something Went Wrong"

	}

	if dirvalues == 1 {

		impactvalues, err := dataAccess.GetRecordimpact(page)

		if err != nil {

			return t, false, err, "Something Went Wrong"

		}

		urgencyvalues, err := dataAccess.GetRecordurgency(page)

		if err != nil {

			return t, false, err, "Something Went Wrong"

		}

		t.Recordimpact = impactvalues

		t.Recordurgency = urgencyvalues

	}

	//            logger.Log.Println("catvalues length --->", len(catvalues))

	t.WorkingCatLabelID = workcatval.WorkingCatLabelID

	t.AssetAttached = assetchkval.AssetAttached

	t.Recordstatus = statusvalues

	t.Businessmatrixdirection = dirvalues

	t.Recordterms = termslist

	t.Recordcatpos = -1

	if len(catvalues) > 0 {

		var tempCat int64

		var tempPos int64

		var seqNos int64

		catDiffTypes := []entities.RecordcategorydetailsEntity{}

		catDiffs := []entities.RecordcategorydetailsEntity{}

		catDiffType := entities.RecordcategorydetailsEntity{}

		for k, v := range catvalues {

			//logger.Log.Println("tempCat ---------------------------------->", tempCat)

			//logger.Log.Println("v.Typeid ------------------------------->", v.Typeid)

			if tempCat != v.Typeid {

				if k != 0 {

					seqNos++

					catDiffType.Child = catDiffs

					if len(catDiffs) == 1 {

						catDiffType.IsDisabled = true

					} else {

						catDiffType.IsDisabled = false

					}

					if catDiffType.IsDisabled == false && t.Recordcatpos == -1 {

						t.Recordcatpos = tempPos

					}

					catDiffType.Sequanceno = seqNos

					catDiffTypes = append(catDiffTypes, catDiffType)

					catDiffs = []entities.RecordcategorydetailsEntity{}

					tempPos++

				}

				tempCat = v.Typeid

				catDiffType.ID = v.Typeid

				catDiffType.Title = v.Typename

				catDiffType.Sequanceno = v.Typeseq

			}

			catDiff := entities.RecordcategorydetailsEntity{}

			catDiff.ID = v.ID

			catDiff.Title = v.Name

			catDiff.Sequanceno = v.Seqno

			//logger.Log.Println("catDiff values --->", catDiff)

			catDiffs = append(catDiffs, catDiff)

			//logger.Log.Println("catDiffs values --->", catDiffs)

		}

		seqNos++

		if len(catDiffs) <= 1 {

			catDiffType.IsDisabled = true

		} else {

			catDiffType.IsDisabled = false

		}

		if !catDiffType.IsDisabled && t.Recordcatpos == -1 {

			t.Recordcatpos = tempPos

		}

		catDiffType.Sequanceno = seqNos

		catDiffType.Child = catDiffs

		catDiffTypes = append(catDiffTypes, catDiffType)

		//logger.Log.Println("catDiffTypes values ------------------->", catDiffTypes)

		if t.Recordcatpos == -1 {

			t.Recordcatpos = 0

		} else {

			for k, _ := range catDiffTypes {

				if k > int(t.Recordcatpos) {

					//logger.Log.Println(k, int(t.Recordcatpos))

					catDiffTypes[k].Child = []entities.RecordcategorydetailsEntity{}

				}

			}

		}

		for _, vv := range catDiffTypes {

			if vv.Sequanceno == 1 && len(vv.Child) > 0 {

				aa := entities.AdditionalFieldDiffEntity{}

				aa.Mstdifferentiationid = vv.Child[0].ID

				aa.Mstdifferentiationtypeid = vv.ID

				bb = append(bb, aa)

			}

		}

		addReq.Mstdifferentiationset = bb

		addfields, err := dataAccess.GetAdditionalFieldsByDiffId(&addReq)

		if err != nil {

			return t, false, err, "Something Went Wrong"

		}

		t.AdditionalFields = addfields

		t.Recordcategory = catDiffTypes

		//logger.Log.Println("category value is in if --->", t)

		return t, true, err, ""

	} else {

		//logger.Log.Println("category value is in else --->", t)

		return t, true, err, ""

	}

}

func GetRecordcatchild(page *entities.RecordcreaterequestEntity) ([]entities.RecordcatchildEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordcatchild")
	t := []entities.RecordcatchildEntity{}
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

	typevalues, err := dataAccess.GetRecordcatchild(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	return typevalues, true, err, ""
}

func GetRecordprioritydata(page *entities.RecordcreaterequestEntity) (entities.RecordcatchildNEstimatedEfforEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordprioritydata")
	t := entities.RecordcatchildNEstimatedEfforEntity{}
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

	typevalues, err := dataAccess.GetRecordprioritydata(page)

	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	estimatedefforts, slaCompliances, changeTypes, err1 := dataAccess.GetEstimateEffort(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	t.ChildCategories = typevalues
	t.EstimatedEfforts = estimatedefforts
	t.SlaCompliance = slaCompliances
	t.ChangeType = changeTypes
	return t, true, err, ""
}
func GetRecordtypedata(page *entities.RecordcreaterequestEntity) ([]entities.RecordtypedetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetRecordtypedatamodel")
	t := []entities.RecordtypedetailsEntity{}
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

	typevalues, err := dataAccess.GetAllRecordtypes(page)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	return typevalues, true, err, ""
}
