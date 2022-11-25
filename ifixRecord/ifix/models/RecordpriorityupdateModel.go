package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
	"time"
)

func Updatepriority(tz *entities.RecordpriorityEntity) (int64, bool, error, string) {

	logger.Log.Println("In side model")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return 0, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	hashmap, err := dataAccess.Getdiffdtls(tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	tx, err := db.Begin()
	if err != nil {
		//db.Close()
		logger.Log.Println("Transaction creation error.", err)
		return 0, false, err, "Something Went Wrong"
	}

	stageID, err := dao.InsertRecordStage(tx, tz)
	if err != nil {
		tx.Rollback()
		//db.Close()
		return 0, false, err, "Something Went Wrong"
	}
	if stageID > 0 {
		pid, _, err := dao.GetlatestDiffID(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 5)
		if err != nil {
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}
		prepirotynm, err := dao.GetPreviouspriorityname(tx, pid)
		if err != nil {
			logger.Log.Println(err)
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		err = dao.Updateoldpriorityflag(tx, pid)
		if err != nil {
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		id, err := dao.InsertRecordMapDifferrtiation(tx, tz, stageID, 0)
		if err != nil {
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		curpirotynm, err := dao.Getpriorityname(tx, tz.Recorddiffid)
		logger.Log.Println("current ")
		if err != nil {
			logger.Log.Println(err)
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		//Update Stage TBL For Priority

		//err = dao.UpdateStagePriority(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Recorddiffid, curpirotynm)
		err = dataAccess.UpdateStagePriority(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Recorddiffid, curpirotynm)
		if err != nil {
			log.Println("error is ----->", err)
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		//Update Stage TBL For Priority

		// Newly addition in 18.07.2021
		//typeID, _, err := dao.GetlatestDiffID(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 2)
		typeID, err := dataAccess.GetRecordTypeID(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		impactID, urgencyID, err := dao.GetImpactUrgencydtls(tx, tz.ClientID, tz.Mstorgnhirarchyid, typeID, tz.Recorddiffid)
		if err != nil {
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		// Impact ID
		preimpactID, _, err := dao.GetlatestDiffID(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 7)
		if err != nil {
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		if preimpactID > 0 {
			// err = dao.Updateoldpriorityflag(tx, preimpactID)
			// if err != nil {
			// 	tx.Rollback()
			// 	//	db.Close()
			// 	return 0, false, err, "Something Went Wrong"
			// }

			err = dao.UpdateoldpriorityflagNew(tx, tz.RecordID, 7)
			if err != nil {
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			err := dao.InsertTrnRecordMapDifferrtiation(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, stageID, 7, impactID, 0, "")
			if err != nil {
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			//Update Stage TBL For Impact

			//err = dao.UpdateStageImpact(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, impactID, hashmap[impactID])
			err = dataAccess.UpdateStageImpact(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, impactID, hashmap[impactID])
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			//Update Stage TBL For Impact

		}

		//Urgency ID
		preurgencyID, _, err := dao.GetlatestDiffID(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 8)
		if err != nil {
			tx.Rollback()
			//	db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		if preurgencyID > 0 {
			// err = dao.Updateoldpriorityflag(tx, preurgencyID)
			// if err != nil {
			// 	tx.Rollback()
			// 	//	db.Close()
			// 	return 0, false, err, "Something Went Wrong"
			// }

			err = dao.UpdateoldpriorityflagNew(tx, tz.RecordID, 8)
			if err != nil {
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			err := dao.InsertTrnRecordMapDifferrtiation(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, stageID, 8, urgencyID, 0, "")
			if err != nil {
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			//Update Stage TBL For Impact

			//err = dao.UpdateStageUrgency(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, urgencyID, hashmap[urgencyID])
			err = dataAccess.UpdateStageUrgency(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, urgencyID, hashmap[urgencyID])
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			//Update Stage TBL For Impact
		}
		//activity log entry here

		err = dao.InsertActivityLogs(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 10, "From "+prepirotynm+" To "+curpirotynm, tz.Userid, tz.Usergroupid)
		if err != nil {
			log.Println("error is ----->", err)
			tx.Rollback()
			//	db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		Username, err := dataAccess.GetUsername(tz.Userid)
		if err != nil {
			logger.Log.Println(err)
			tx.Rollback()
			//	db.Close()
			return 0, false, err, "Something Went Wrong"
		}

		//activity log entry end
		logger.Log.Println("Priority change id is --------------------------------------------------->", id)
		if id > 0 {
			err = tx.Commit()
			if err != nil {
				logger.Log.Println("commit error --------------------------------------------------->", err)
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			err = dataAccess.UpdateUserInfo(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}

			// For SLA calling
			historyrecord, err := dataAccess.GetLatesttrnhistoryAll(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			currentTime := time.Now()
			zonediff, _, _, _ := Getutcdiff(tz.ClientID, tz.Mstorgnhirarchyid)
			datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)

			createdt, err := dataAccess.Getrecordcreatedate(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if err != nil {
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}
			reopendt, err := dataAccess.Getrecordreopencreatedate(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if err != nil {
				tx.Rollback()
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}
			if len(createdt) > 0 {
				logger.Log.Println("createdt  1111111111111111---------------------------------------------------->", createdt)
				datetime1 := AddSubSecondsToDate(TimeParse(createdt, ""), zonediff.UTCdiff)
				// value assign
				histrn := entities.TrnslaentityhistoryEntity{}
				histrn.Clientid = tz.ClientID
				histrn.Mstorgnhirarchyid = tz.Mstorgnhirarchyid
				histrn.Therecordid = tz.RecordID
				histrn.Mstslaentityid = historyrecord.Mstslaentityid
				histrn.Donotupdatesladue = historyrecord.Donotupdatesladue
				histrn.Mstslastateid = historyrecord.Mstslastateid
				histrn.Commentonrecord = historyrecord.Commentonrecord
				histrn.Fromclientuserid = historyrecord.Fromclientuserid
				histrn.Recorddatetime = datetime.Format("2006-01-02 15:04:05")
				histrn.Recorddatetoint = TimeParse(datetime.Format("2006-01-02 15:04:05"), "").Unix()
				histrn.Slastartstopindicator = historyrecord.Slastartstopindicator
				trnid, err := dataAccess.InsertTrnslaentityhistory(&histrn)
				if err != nil {
					logger.Log.Println(err)
					//	db.Close()
					return 0, false, err, "Something Went Wrong"
				}
				logger.Log.Println("history table id---------------------------------------------------->", trnid)
				res, err := dataAccess.Getrecorddetails(tz.RecordID)
				if err != nil {
					return 0, false, err, "Something Went Wrong"
				}
				grpID, err := dataAccess.FetchCurrentGrpID(tz.RecordID)
				if err != nil {
					return 0, false, err, "Something Went Wrong"
				}

				returnValue, _, _, _ := SLACriteriaRespResl(tz.ClientID, tz.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)
				if returnValue.Supportgroupspecific == 1 {
					count, err := dataAccess.GetSupportgrpdayofweekcount(tz.ClientID, tz.Mstorgnhirarchyid, grpID)
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}
					if count < 7 {
						return 0, false, err, "Day Of Week Not Properly Configured.Please Check."
					}
				} else {
					count, err := dataAccess.GetOrganizationdayofweekcount(tz.ClientID, tz.Mstorgnhirarchyid)
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}
					if count < 7 {
						return 0, false, err, "Day Of Week Not Properly Configured.Please Check."
					}
				}

				if len(reopendt) > 0 {
					datetime2 := AddSubSecondsToDate(TimeParse(reopendt, ""), zonediff.UTCdiff)
					SLADueTimeCalculation(tz.RecordID, 0, 1, 3, datetime2.Format("2006-01-02 15:04:05"), tz.ClientID, tz.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "P", grpID)
				} else {
					SLADueTimeCalculation(tz.RecordID, 0, 1, 3, datetime1.Format("2006-01-02 15:04:05"), tz.ClientID, tz.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "P", grpID)

				}

				t := entities.SLATabEntity{}
				t.ClientID = tz.ClientID
				t.Mstorgnhirarchyid = tz.Mstorgnhirarchyid
				t.RecordID = tz.RecordID
				sladtls, _, err, _ := GetSLATabvalues(&t)
				if err != nil {
					logger.Log.Println(err)
				}
				err = dataAccess.UpdateSLAFields(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, sladtls.Responsedetails.Responseduetime, sladtls.Responsedetails.Responseclockstatus, sladtls.Resolutionetails.Resolutionduetime, sladtls.Resolutionetails.Resolutionclockstatus)
				if err != nil {
					logger.Log.Println(err)
				}

				//For Email Notification Start
				//	go PriorityChangeEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Recorddiffid)
				//For Email Notification End
				//	db.Close()
				return stageID, true, err, "Priority update successfully done."
			} else {
				//	db.Close()
				return 0, false, err, "Something Went Wrong"
			}
			//End SLA calling
		}
	}
	tx.Rollback()
	//db.Close()
	return 0, false, err, "Something Went Wrong"
}

//curl -v localhost:8083/updatepriority -d '{"clientid":2,"mstorgnhirarchyid":2,"recorddifftypeid":5,"recorddiffid":471,"recordid":255,"userid":13,"usergroupid":12,"originaluserid":13,"originalusergroupid":12,"recordname":"test","recordesc":"test1111"}'

//{"success":true,"message":"Priority update successfully done.","stageid":261}
