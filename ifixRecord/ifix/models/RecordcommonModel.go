package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"strconv"

	"log"
	"strings"
	"time"
)
func InsertBulkRecordTermvalues(pz *entities.BulkRecordcommonEntity) (int64, bool, error, string) {
	logger.Log.Println("In side InsertBulkRecordTermvalues", pz)
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	var err error
	dataAccess := dao.DbConn{DB: db}
	var tz entities.RecordcommonEntity
	logger.Log.Println("In side InsertBulkRecordTermvalues", pz.Termseq)

	tz.Termseq = pz.Termseq
	tz.Termvalue = pz.Termvalue
	tz.Userid = pz.Userid
	tz.Usergroupid = pz.Usergroupid
	for i := 0; i < len(pz.RecordID); i++ {

		clientid, orgid, recordstageid, err := dataAccess.Getclientorgvalues(pz.RecordID[i])
		//logger.Log.Println("Term name is  -----2222222222222222222222222222222------>", termnm)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		tz.RecordID = pz.RecordID[i]
		tz.ClientID = clientid
		tz.Mstorgnhirarchyid = orgid
		tz.RecordstageID = recordstageid
		hashmap, err := dataAccess.Termsequance(tz.ClientID, tz.Mstorgnhirarchyid)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		if len(hashmap) == 0 {
			return 0, false, err, "Something Went Wrong"
		}
		if tz.TermID == 0 {
			id := hashmap[tz.Termseq]
			tz.TermID = id
		}
		termnm, err := dataAccess.Gettermnamebyid(tz.TermID, tz.ClientID, tz.Mstorgnhirarchyid)
		//logger.Log.Println("Term name is  -----2222222222222222222222222222222------>", termnm)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		_, err = dataAccess.InsertRecordTermvalues(&tz)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
	}
	return 0, true, err, ""
}
/*func InsertRecordTermvalues(tz *entities.RecordcommonEntity) (int64, bool, error, string) {
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println(tz)
	id, status, err, msg := InsertRecordTermvalueswithdb(tz, db)
	if err != nil {
		logger.Log.Println("Error in DBConnection")
		return 0, false, err, "Something Went Wrong"
	}

	return id, status, err, msg
}*/

func InsertRecordTermvalues(tz *entities.RecordcommonEntity) (int64, bool, error, string) {
	logger.Log.Println("In side Mstslastatemodel")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	log.Println("database connection failure", err)
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
	//childids, err := dataAccess.Getchildrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	childids, err := dataAccess.GetchildrecordidsForCommon(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	//logger.Log.Println("Child ticket size is ----------->", len(childids))
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	hashmap, err := dataAccess.Termsequance(tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if len(hashmap) == 0 {
		return 0, false, err, "Something Went Wrong"
	}
	if tz.TermID == 0 {
		id := hashmap[tz.Termseq]
		tz.TermID = id
	}
	//logger.Log.Println("Termid value is ==============================>", tz.TermID)
	//logger.Log.Println("Term seuance value is ==============================>", tz.Termseq)

	releationID, err := dataAccess.Checktermreleation(tz.TermID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	logger.Log.Println("Child ticket size is ----------->", releationID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	termnm, err := dataAccess.Gettermnamebyid(tz.TermID, tz.ClientID, tz.Mstorgnhirarchyid)
	//logger.Log.Println("Term name is  -----2222222222222222222222222222222------>", termnm)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	grplevel, err := dataAccess.GetgrplevelID(tz.ClientID, tz.Mstorgnhirarchyid, tz.Usergroupid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	//logger.Log.Println("Term id  is in InsertRecordTermvalues ----------->", tz.TermID)

	id, err := dataAccess.InsertRecordTermvalues(tz)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	if tz.Termseq == 29 {
		fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if fcount > 0 {
			//Update Stage TBL For Impact
			err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 17, "The value is :: "+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

		}

	}

	Username, err := dataAccess.GetUsername(tz.Userid)
	if err != nil {
		logger.Log.Println(err)
		return 0, false, err, "Something Went Wrong"
	}

	err = dataAccess.UpdateUserInfoWithoutTrn(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
	if err != nil {
		logger.Log.Println(err)
		//return 0, false, err, "Something Went Wrong"
	}
	//logger.Log.Println("childids values is  ----------->", childids)
	//logger.Log.Println("releationID values is  ----------->", releationID)
	//logger.Log.Println("grplevel values is  ----------->", grplevel)

	if len(childids) > 0 && releationID > 0 && grplevel > 1 {
		for k := 0; k < len(childids); k++ {
			//	logger.Log.Println("222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222  ----------->")
			_, err := dataAccess.InsertRecordTermvaluesForchilds(tz, childids[k])
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

			err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

			if tz.Termseq == 11 {
				go CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
			}
			if tz.Termseq == 1 {
				go FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
			}
			err = dataAccess.UpdateUserInfoWithoutTrn(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
			if err != nil {
				logger.Log.Println(err)
				//return 0, false, err, "Something Went Wrong"
			}

			if tz.Termseq == 29 {
				fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
				if fcount > 0 {
					err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 17, "The value is :: "+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}

				}

			}

		}
	} // child for loop

	// For Stask ticket type ......
	typeseq, err := dataAccess.Getrecordtypeseq(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	// For Stask ticket type ......
	if len(childids) > 0 && releationID > 0 && grplevel == 0 { // For SR ticket type...
		if typeseq == 2 || typeseq == 4 {
			for k := 0; k < len(childids); k++ {
				//	logger.Log.Println("333333333333333333333333333333333333333333333333333333333333333333333333333333333333  ----------->", typeseq)
				_, err := dataAccess.InsertRecordTermvaluesForchilds(tz, childids[k])
				if err != nil {
					return 0, false, err, "Something Went Wrong"
				}

				err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
				if err != nil {
					return 0, false, err, "Something Went Wrong"
				}

				if tz.Termseq == 29 {
					fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
					if fcount > 0 {
						err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 17, "The value is :: "+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
						if err != nil {
							return 0, false, err, "Something Went Wrong"
						}

					}

				}

				if tz.Termseq == 11 {
					go CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
				}
				if tz.Termseq == 1 {
					go FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
				}
				err = dataAccess.UpdateUserInfoWithoutTrn(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
				if err != nil {
					logger.Log.Println(err)
					//return 0, false, err, "Something Went Wrong"
				}

			}
		}

	} // child for loop

	if typeseq == 3 || typeseq == 5 {
		_, Seq, _, err := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		//3,8,11,14
		if Seq != 3 && Seq != 8 && Seq != 11 && Seq != 14 {
			parentid, err := dataAccess.Getparentrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
			//logger.Log.Println("Child ticket size is ----------->", len(parentid))
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

			if len(parentid) > 0 && releationID > 0 {
				for k := 0; k < len(parentid); k++ {
					_, err := dataAccess.InsertRecordTermvaluesForchilds(tz, parentid[k])
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}

					err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}

					if tz.Termseq == 29 {
						fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k])
						if fcount > 0 {
							err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], 17, "The value is ::"+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
							if err != nil {
								return 0, false, err, "Something Went Wrong"
							}

						}

					}

					if tz.Termseq == 11 {
						go CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k])
					}
					if tz.Termseq == 1 {
						go FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k])
					}

				}
			} // parent for loop
		} // Seq if condition end here...
	} // typeseq if condition
	// For Stask ticket type ............

	if id > 0 {
		if tz.Termseq == 11 {
			CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		}

		if tz.Termseq == 1 {
			FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		}
		if tz.Termseq == 29 {
			fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if fcount > 0 {
				//Update Stage TBL For Impact

				err = dataAccess.UpdateFollowupcountFromCommon(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, fcount)
				if err != nil {
					log.Println("error is ----->", err)
				}

				//Update Stage TBL For Impact
				go FollowupCountEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, fcount)
			}

		}

		if tz.Termseq == 30 {
			obcount, err := dataAccess.GetOutboundcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if err != nil {
				log.Println("error is ----->", err)
			}
			err = dataAccess.UpdateOutboundcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, obcount)
			if err != nil {
				log.Println("error is ----->", err)
			}
		}

		// ebonding code start here 10.06.2022
		if tz.Termseq == 1 || tz.Termseq == 11 {
			workingdifftypeID, workingdiffID, err := dataAccess.GetWorkingCategoryDetails(tz.RecordID)
			if err != nil {
				//return "", 0, errors.New("Something Went Wrong")
				logger.Log.Println("ebonding error is ----->", err)
			}

			ebondingSeq, err := dataAccess.GetEbondingSeq(tz.ClientID, tz.Mstorgnhirarchyid, workingdifftypeID, workingdiffID)
			if err != nil {
				//return "", 0, errors.New("Something Went Wrong")
				logger.Log.Println("ebonding error is ----->", err)
			}

			ticketID, recordtypeid, err := dataAccess.GetRecordCodeANDTypeID(tz.RecordID)
			if err != nil {
				//return "", 0, errors.New("Something Went Wrong")
				logger.Log.Println("ebonding error is ----->", err)
			}

			logger.Log.Println("ebonding seq number is =====> ", ebondingSeq)
			if ebondingSeq > 0 {

				if tz.Termseq == 1 {
					ebonding := entities.EbondingRecordEntity{}
					ebonding.ClientID = tz.ClientID
					ebonding.MstorgnhirarchyID = tz.Mstorgnhirarchyid
					ebonding.RecorddifftypeID = 2
					ebonding.RecorddiffID = recordtypeid
					ebonding.RecordID = tz.RecordID
					ebonding.RecordStagedID = tz.RecordstageID
					ebonding.RecordCode = ticketID
					ebonding.EbondingSeq = ebondingSeq
					ebonding.EbondingModuleSeq = 4
					ebonding.UploadedFileName = tz.Termdescription
					ebonding.OriginalFileName = tz.Termvalue
					go Ebonding(&ebonding)
				} else if tz.Termseq == 11 {
					ebonding := entities.EbondingRecordEntity{}
					ebonding.ClientID = tz.ClientID
					ebonding.MstorgnhirarchyID = tz.Mstorgnhirarchyid
					ebonding.RecorddifftypeID = 2
					ebonding.RecorddiffID = recordtypeid
					ebonding.RecordID = tz.RecordID
					ebonding.RecordStagedID = tz.RecordstageID
					ebonding.RecordCode = ticketID
					ebonding.EbondingSeq = ebondingSeq
					ebonding.EbondingModuleSeq = 3
					ebonding.Worknote = tz.Termvalue
					go Ebonding(&ebonding)
				}

			}

		}

		//ebonding code end here 10.06.2022


		return id, true, err, ""
	} else {
		return 0, false, err, "Something Went Wrong"
	}

}

func InsertMultipleTermvalues(tz *entities.RecordmultiplecommonEntity) (int64, bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// dbcon, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	//dbcon.Close()
	// 	logger.Log.Println("Error in DBConnection in side InsertMultipleTermvalues")
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
		logger.Log.Println("Transaction creation error in InsertMultipleTermvalues", err)
		return 0, false, err, "Something Went Wrong"
	}
	dataAccess := dao.DbConn{DB: db}
	hashmap, err := dataAccess.TermsIDs(tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	Username, err := dataAccess.GetUsername(tz.Userid)
	if err != nil {
		logger.Log.Println(err)
		tx.Rollback()
		//dbcon.Close()
		return 0, false, err, "Something Went Wrong"
	}

	for i := 0; i < len(tz.Details); i++ {
		logger.Log.Println("Details value is ---------------->", &tz.Details[i])
		v := &tz.Details[i]
		if len(v.Insertedvalue) > 0 {
			_, err := dao.InsertMultipleRecordTermvalues(tx, &tz.Details[i], tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.RecordstageID, tz.ForuserID, tz.Userid, tz.Usergroupid)
			if err != nil {
				tx.Rollback()
				//	dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}
			termnm, err := dao.Gettermnamebyid(tx, tz.Details[i].ID, tz.ClientID, tz.Mstorgnhirarchyid)
			if err != nil {
				tx.Rollback()
				//	dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}
			err = dao.InsertActivityLogsfromterms(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 100, termnm+" :: "+tz.Details[i].Insertedvalue, tz.Userid, tz.Usergroupid, tz.Details[i].ID)
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				//	dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			termseq := hashmap[tz.Details[i].ID]
			if termseq == 33 {
				err = dataAccess.UpdateClosureCode(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 35 {
				err = dataAccess.UpdateClosureComment(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 17 {
				err = dataAccess.UpdateClosureComment(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 4 {
				err = dataAccess.UpdatePendingvendorname(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 5 {
				err = dataAccess.UpdatePendingvendorticketID(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 6 {
				err = dataAccess.UpdateResolutioncode(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 8 {
				err = dataAccess.UpdateResolutioncomment(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			// New Addition on 28.04.2022
			if termseq == 23 {
				err = dataAccess.UpdateResponseBreachCode(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 24 {
				err = dataAccess.UpdateResponseBreachComment(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 25 {
				err = dataAccess.UpdateResolutionBreachCode(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			if termseq == 26 {
				err = dataAccess.UpdateResolutionBreachComment(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Details[i].Insertedvalue)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//	dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			// New Addition on 28.04.2022

			err = dataAccess.UpdateUserInfo(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//	dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

		} // If condition end here ...
	} // Details end here...

	//dataAccess := dao.DbConn{DB: dbcon}
	//childids, err := dataAccess.Getchildrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	childids, err := dataAccess.GetchildrecordidsForCommon(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	logger.Log.Println("Child ticket size is ----------->", len(childids))

	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if len(tz.Details) > 0 {
		//rec := &tz.Details[0]

		if len(childids) > 0 {
			for k := 0; k < len(childids); k++ {
				for i := 0; i < len(tz.Details); i++ {
					logger.Log.Println("Details vaInsertMultipleRecordTermvalueslue is ---------------->", &tz.Details[i])
					v := &tz.Details[i]
					if len(v.Insertedvalue) > 0 {
						releationID, err := dataAccess.Checktermreleation(tz.Details[i].ID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
						logger.Log.Println("Releation value  is  ----------->", releationID)
						if err != nil {
							return 0, false, err, "Something Went Wrong"
						}
						if releationID > 0 {
							_, err := dao.InsertMultipleRecordTermvalues(tx, &tz.Details[i], tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.RecordstageID, tz.ForuserID, tz.Userid, tz.Usergroupid)
							if err != nil {
								tx.Rollback()
								//	dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							termnm, err := dao.Gettermnamebyid(tx, tz.Details[i].ID, tz.ClientID, tz.Mstorgnhirarchyid)
							if err != nil {
								tx.Rollback()
								//	dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}
							err = dao.InsertActivityLogsfromterms(tx, tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 100, termnm+" :: "+tz.Details[i].Insertedvalue, tz.Userid, tz.Usergroupid, tz.Details[i].ID)
							if err != nil {
								log.Println("error is ----->", err)
								tx.Rollback()
								//	dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							termseq := hashmap[tz.Details[i].ID]
							if termseq == 33 {
								err = dataAccess.UpdateClosureCode(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 35 {
								err = dataAccess.UpdateClosureComment(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 17 {
								err = dataAccess.UpdateClosureComment(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 4 {
								err = dataAccess.UpdatePendingvendorname(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 5 {
								err = dataAccess.UpdatePendingvendorticketID(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 6 {
								err = dataAccess.UpdateResolutioncode(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 8 {
								err = dataAccess.UpdateResolutioncomment(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							// New Addition on 28.04.2022
							if termseq == 23 {
								err = dataAccess.UpdateResponseBreachCode(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 24 {
								err = dataAccess.UpdateResponseBreachComment(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 25 {
								err = dataAccess.UpdateResolutionBreachCode(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							if termseq == 26 {
								err = dataAccess.UpdateResolutionBreachComment(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Details[i].Insertedvalue)
								if err != nil {
									log.Println("error is ----->", err)
									tx.Rollback()
									//	dbcon.Close()
									return 0, false, err, "Something Went Wrong"
								}
							}

							// New Addition on 28.04.2022

							err = dataAccess.UpdateUserInfo(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], tz.Userid, Username)
							if err != nil {
								logger.Log.Println(err)
								tx.Rollback()
								//	dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}
						}
					}

				} // details end here...

			}
		} // child id if condition .......

		// For Stask logic ...........

		typeseq, err := dataAccess.Getrecordtypeseq(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		if typeseq == 3 {
			_, Seq, _, err := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}
			//3,8,11,14
			if Seq != 3 && Seq != 8 && Seq != 11 && Seq != 14 {

				parentid, err := dataAccess.Getparentrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
				logger.Log.Println("Parent ticket size is ----------->", len(parentid))
				if err != nil {
					return 0, false, err, "Something Went Wrong"
				}

				if len(parentid) > 0 {
					for k := 0; k < len(parentid); k++ {
						for i := 0; i < len(tz.Details); i++ {
							logger.Log.Println("Details vaInsertMultipleRecordTermvalueslue is ---------------->", &tz.Details[i])
							v := &tz.Details[i]
							if len(v.Insertedvalue) > 0 {
								releationID, err := dataAccess.Checktermreleation(tz.Details[i].ID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
								logger.Log.Println("Releation value  is ----------->", releationID)
								if err != nil {
									return 0, false, err, "Something Went Wrong"
								}
								if releationID > 0 {
									_, err := dao.InsertMultipleRecordTermvalues(tx, &tz.Details[i], tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.RecordstageID, tz.ForuserID, tz.Userid, tz.Usergroupid)
									if err != nil {
										tx.Rollback()
										//	dbcon.Close()
										return 0, false, err, "Something Went Wrong"
									}

									termnm, err := dao.Gettermnamebyid(tx, tz.Details[i].ID, tz.ClientID, tz.Mstorgnhirarchyid)
									if err != nil {
										tx.Rollback()
										//	dbcon.Close()
										return 0, false, err, "Something Went Wrong"
									}
									err = dao.InsertActivityLogsfromterms(tx, tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], 100, termnm+" :: "+tz.Details[i].Insertedvalue, tz.Userid, tz.Usergroupid, tz.Details[i].ID)
									if err != nil {
										log.Println("error is ----->", err)
										tx.Rollback()
										//	dbcon.Close()
										return 0, false, err, "Something Went Wrong"
									}

									termseq := hashmap[tz.Details[i].ID]
									if termseq == 33 {
										err = dataAccess.UpdateClosureCode(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 35 {
										err = dataAccess.UpdateClosureComment(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 17 {
										err = dataAccess.UpdateClosureComment(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 4 {
										err = dataAccess.UpdatePendingvendorname(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 5 {
										err = dataAccess.UpdatePendingvendorticketID(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 6 {
										err = dataAccess.UpdateResolutioncode(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 8 {
										err = dataAccess.UpdateResolutioncomment(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									// New Addition on 28.04.2022
									if termseq == 23 {
										err = dataAccess.UpdateResponseBreachCode(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 24 {
										err = dataAccess.UpdateResponseBreachComment(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 25 {
										err = dataAccess.UpdateResolutionBreachCode(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									if termseq == 26 {
										err = dataAccess.UpdateResolutionBreachComment(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], tz.Details[i].Insertedvalue)
										if err != nil {
											log.Println("error is ----->", err)
											tx.Rollback()
											//	dbcon.Close()
											return 0, false, err, "Something Went Wrong"
										}
									}

									// New Addition on 28.04.2022

								}
							}
						} // details end here...

					}
				} // child id if condition ......
			} //Seq if condition end here ...
		} // typeseq if condition

		// For Stask logic ...........

	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		//	dbcon.Close()
		return 0, false, err, "Something Went Wrong"
	}
	return 0, true, err, ""
}

func GetAllcommontermvalues(page *entities.RecordcommonEntity) ([]entities.RecordcommonresponseEntity, bool, error, string) {
	logger.Log.Println("In side GetAllcommontermvalues")
	t := []entities.RecordcommonresponseEntity{}
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
	values, err1 := dataAccess.GetAllcommontermvalues(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetTermvalueagainsttermid(page *entities.RecordcommonEntity) ([]entities.RecordcommonresponseEntity, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	t := []entities.RecordcommonresponseEntity{}
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
	values, err1 := dataAccess.GetTermvalueagainsttermid(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetTermnames(page *entities.RecordcommonEntity) ([]entities.RecordTermnamesEntity, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	t := []entities.RecordTermnamesEntity{}
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
	values, err1 := dataAccess.GetTermnames(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetTermnamesbystate(page *entities.RecordcommonstateEntity) ([]entities.RecordTermnamesEntity, bool, error, string) {
	logger.Log.Println("In side GetTermnamesbystate")
	t := []entities.RecordTermnamesEntity{}
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
	values, err1 := dataAccess.GetTermnamesbystate(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func Getrecordcount(tz *entities.RecordcommonEntity) (entities.RecordcountEntity, bool, error, string) {
	t := entities.RecordcountEntity{}
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
	pcount, err := dataAccess.GetPrioritycount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	pvacount, err := dataAccess.GetPendingvendorcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	rcount, err := dataAccess.GetReopencount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	/*aging, err := dataAccess.GetAging(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	if aging > 0 {
		err = dataAccess.UpdateAging(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, (aging * 24 * 60 * 60))
		if err != nil {
			log.Println("error is ----->", err)
		}
	}*/


	_, recordseq, _, err := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	recordcreatedate, err := dataAccess.GetRecordcreatedate(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	var enddate string
	var aging int64
	zonediff, _, _, _ := Getutcdiff(tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		logger.Log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}

	if len(recordcreatedate) > 0 {
		fromdatetime := AddSubSecondsToDate(TimeParse(recordcreatedate, ""), zonediff.UTCdiff)
		var t1 int64
		if recordseq == 3 {
			enddate, err = dataAccess.GetReopendate(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 3)
			if err != nil {
				logger.Log.Println("database connection failure", err)
				return t, false, err, "Something Went Wrong"
			}
			todatetime1 := AddSubSecondsToDate(TimeParse(enddate, ""), zonediff.UTCdiff)
			t1 = SubtractDateToDate(todatetime1, fromdatetime)
		} else {
			currentTime := time.Now()
			todatetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
			t1 = SubtractDateToDate(todatetime, fromdatetime)
		}
		logger.Log.Println("Daycount value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", tz.RecordID, t1)
		if t1 > 0 {
			aging = (t1 / (24 * 60 * 60))
			if aging > 0 {
				err = dataAccess.UpdateAging(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, (aging * 24 * 60 * 60))
				if err != nil {
					log.Println("error is ----->", err)
				}
			}
		}
	}




	typeseq, err := dataAccess.Getrecordtypeseq(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	var fcount int64
	var obcount int64

	if typeseq == 2 {
		// parentid, err := dataAccess.Getparentrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
		// logger.Log.Println("Child ticket size is ----------->", len(parentid))
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		//if len(parentid) > 0 {
		fcount, err = dataAccess.GetFollowupcountForSR(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}

		obcount, err = dataAccess.GetOutboundcountForSR(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		//}

	} else {
		fcount, err = dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}

		obcount, err = dataAccess.GetOutboundcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
	}

	t.Followupcount = fcount
	t.Outboundcount = obcount
	t.Pendingvendoractioncount = pvacount
	t.Prioritycount = (pcount - 1)
	t.Reopencount = rcount
	t.Aging = aging

	return t, true, err, ""
}

func Getrecentrecord(tz *entities.RecordcommonEntity) ([]entities.RecentrecordEntity, bool, error, string) {
	t := []entities.RecentrecordEntity{}
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
	orgntype, err1 := dataAccess.GetOrgnType(tz.ClientID, tz.Mstorgnhirarchyid)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	grplevel, err := dataAccess.GetgrplevelID(tz.ClientID, tz.Mstorgnhirarchyid, tz.Usergroupid)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	logger.Log.Println("GrpID value is -------------->", tz.Usergroupid)
	util, err1 := dataAccess.Gettimediff(tz.ClientID, tz.Mstorgnhirarchyid)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	if grplevel == 1 {
		t, err = dataAccess.GetRecentrecords(tz, util.Timediff)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
	} else {
		t, err = dataAccess.GetRecentrecordsForResolver(tz, util.Timediff, orgntype)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
	}

	return t, true, err, ""
}

func GetRecordlogs(tz *entities.RecordcommonEntity) ([]entities.RecordlogsEntity, bool, error, string) {
	t := []entities.RecordlogsEntity{}
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
	values, err := dataAccess.GetRecordlogs(tz)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	return values, true, err, ""
}

func Getfrequentissues(tz *entities.RecordcommonEntity) ([]entities.FrequentRecordEntity, bool, error, string) {
	t := []entities.FrequentRecordEntity{}
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

	grplevel, err := dataAccess.GetgrplevelID(tz.ClientID, tz.Mstorgnhirarchyid, tz.Usergroupid)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	logger.Log.Println("GrpID value is -------------->", tz.Usergroupid)

	if grplevel == 1 {
		values, err := dataAccess.Getfrequentissues(tz)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		for i := 0; i < len(values); i++ {
			val := entities.FrequentRecordEntity{}
			var rows = values[i]
			res := strings.Split(rows.Parentcatname, "->")
			res1 := "<b>" + res[len(res)-2] + "</b> - " + res[len(res)-1]
			val.Mstorgnhirarchyid = rows.Mstorgnhirarchyid
			val.LastlevelID = rows.LastlevelID
			val.Count = rows.Count
			val.ParentcatID = rows.ParentcatID
			val.Parentcatname = res1
			val.Recorddifftypeid = rows.Recorddifftypeid
			val.Recorddiffid = rows.Recorddiffid
			val.Seq = rows.Seq
			t = append(t, val)
		}
	} else {
		values, err := dataAccess.Getfrequentissuesresolver(tz)
		if err != nil {
			return t, false, err, "Something Went Wrong"
		}
		for i := 0; i < len(values); i++ {
			/*val := entities.FrequentRecordEntity{}
			var rows = values[i]
			res := strings.Split(rows.Parentcatname, "->")
			res1 := "<b>" + res[0] + " -> " + res[len(res)-2] + " - </b>" + res[len(res)-1]
			val.Mstorgnhirarchyid = rows.Mstorgnhirarchyid
			val.LastlevelID = rows.LastlevelID
			val.Count = rows.Count
			val.ParentcatID = rows.ParentcatID
			val.Parentcatname = res1
			val.Recorddifftypeid = rows.Recorddifftypeid
			val.Recorddiffid = rows.Recorddiffid
			val.Seq = rows.Seq
			t = append(t, val)*/

			val := entities.FrequentRecordEntity{}
			var rows = values[i]
			parentcategoryids, parentcategorynames, err := dataAccess.GetCategoryInfo(rows.LastlevelID)
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
			res := strings.Split(parentcategorynames, "->")
			if len(res) > 0 {
				res1 := "<b>" + res[0] + " -> " + res[len(res)-2] + " - </b>" + res[len(res)-1]
				val.Mstorgnhirarchyid = rows.Mstorgnhirarchyid
				val.LastlevelID = rows.LastlevelID
				val.Count = rows.Count
				val.ParentcatID = parentcategoryids
				val.Parentcatname = res1
				val.Recorddifftypeid = rows.Recorddifftypeid
				val.Recorddiffid = rows.Recorddiffid
				val.Seq = rows.Seq
				t = append(t, val)
			}
		}
	}

	return t, true, err, ""
}

func GetParentrecord(tz *entities.RecordcommonEntity) ([]entities.ParentticketEntity, bool, error, string) {
	t := []entities.ParentticketEntity{}
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

	//var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}

	dataAccess := dao.DbConn{DB: db}
	values, err := dataAccess.GetParentrecord(tz)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetActivitymstnames(page *entities.RecordcommonEntity) ([]entities.Recordactivitymst, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	t := []entities.Recordactivitymst{}
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
	values, err1 := dataAccess.GetActivitymstnames(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetNewActivitylogs(page *entities.RecordcommonEntity) ([]entities.NewActivitylogsEntity, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	t := []entities.NewActivitylogsEntity{}
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
	values, err1 := dataAccess.GetNewActivitylogs(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	//rows.Scan(&value.ID, &value.RecordID, &value.Logvalue, &value.Createddate, &value.Name, &value.Activitydesc, &value.Supportgroupname, &value.Termname, &value.Status)

	util, err1 := dataAccess.Gettimediff(page.ClientID, page.Mstorgnhirarchyid)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	var str string
	if util.Timediff > 0 {
		for i := 0; i < len(values); i++ {
			//logger.Log.Println("In side activity logs for loop str value -------------------------------->", str)
			//logger.Log.Println("In side activity logs for loop array  value -------------------------------->", values[i].Activitydesc)
			if str != values[i].Logvalue {
				v := entities.NewActivitylogsEntity{}
				v.ID = values[i].ID
				v.RecordID = values[i].RecordID
				v.Logvalue = values[i].Logvalue
				v.Createddate = values[i].Createddate
				v.Name = values[i].Name
				v.Activitydesc = values[i].Activitydesc
				v.Supportgroupname = values[i].Supportgroupname
				v.Termname = values[i].Termname
				v.Status = values[i].Status
				v.Showcreatedate = dao.Convertdate(values[i].Createddate, util.Timediff)
				str = values[i].Logvalue
				t = append(t, v)
			}

		}
	}

	return t, true, err, ""
}

func Activitysearchresults(page *entities.Activitylogsearchcriteria) ([]entities.NewActivitylogsEntity, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	t := []entities.NewActivitylogsEntity{}
	values := []entities.NewActivitylogsEntity{}
	values1 := []entities.NewActivitylogsEntity{}
	//var sequance string
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
	util, err1 := dataAccess.Gettimediff(page.ClientID, page.Mstorgnhirarchyid)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	for i := 0; i < len(page.Searchfilter); i++ {
		if page.Searchfilter[i].Seq == 100 {
			values, err = dataAccess.Searchtermlogs(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Searchfilter[i].ID, values, util.Timediff)
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
			logger.Log.Println("values1 is --------------->", t)
			//t = append(values)
		} else {
			//sequance = sequance + "," + strconv.FormatInt(page.Searchfilter[i].Seq, 10)
			values1, err = dataAccess.Searchnormallogs(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Searchfilter[i].Seq, values1, util.Timediff)
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
		}

	}
	t = append(values, values1...)
	return t, true, err, ""
}

func GetTermnamesbysequance(page *entities.RecordcommonEntity) ([]entities.RecordTermnamesEntity, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	t := []entities.RecordTermnamesEntity{}
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
	for i := 0; i < len(page.Sequance); i++ {
		values, err1 := dataAccess.GetTermnamesbysequance(page, page.Sequance[i])
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}
		t = append(t, values)
	}

	return t, true, err, ""
}

func GetPendingstatustermvalue(page *entities.RecordcommonEntity) ([]entities.Pendingstatustermvalue, bool, error, string) {
	logger.Log.Println("In side GetPendingstatustermvalue")
	t := []entities.Pendingstatustermvalue{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //	defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
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
	values, err1 := dataAccess.GetPendingstatustermvalue(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetAttachmentfiles(page *entities.RecordcommonEntity) ([]entities.RecordAttachmentfiles, bool, error, string) {
	logger.Log.Println("In side GetAttachmentfiles")
	t := []entities.RecordAttachmentfiles{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //	defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }

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
	values, err1 := dataAccess.GetAttachmentfiles(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	util, err1 := dataAccess.Gettimediff(page.ClientID, page.Mstorgnhirarchyid)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	if util.Timediff > 0 {
		for i := 0; i < len(values); i++ {
			if len(values[i].Originalname) > 0 {
				v := entities.RecordAttachmentfiles{}
				v.ID = values[i].ID
				v.RecordID = values[i].RecordID
				v.ClientID = values[i].ClientID
				v.Mstorgnhirarchyid = values[i].Mstorgnhirarchyid
				v.RecordtermID = values[i].RecordtermID
				v.Originalname = values[i].Originalname
				v.Uploadname = values[i].Uploadname
				v.Createdate = values[i].Createdate
				v.Supportgrouplevelid = values[i].Supportgrouplevelid
				v.Createdbyid = values[i].Createdbyid
				v.Createdgrpid = values[i].Createdgrpid
				v.Name = values[i].Name
				v.RecorduserID = values[i].RecorduserID
				v.RecordusergrpID = values[i].RecordusergrpID
				v.RecordoriginaluserID = values[i].RecordoriginaluserID
				v.RecordoriginalusergrpID = values[i].RecordoriginalusergrpID

				v.Showcreatedate = dao.Convertdate(values[i].Createdate, util.Timediff)
				t = append(t, v)
			}

		}
	}

	return t, true, err, ""
}

func Updatedocumentcount(page *entities.RecordcommonEntity) (bool, error, string) {
	logger.Log.Println("In side GetAttachmentfiles")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	err1 := dataAccess.Updatedocumentcount(page.ClientID, page.Mstorgnhirarchyid, page.ID)
	if err1 != nil {
		return false, err1, "Something Went Wrong"
	}

	return true, err, "Document count successfully updated"
}

/*func Customervisiblecomment(page *entities.RecordcommonEntity) ([]entities.Customervisiblecomment, bool, error, string) {
	logger.Log.Println("In side GetAttachmentfiles")
	t := []entities.Customervisiblecomment{}
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
	var daycount int64
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
	values, err := dataAccess.GetLastRecordcomment(page)
	if err != nil {
		logger.Log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}
	//logger.Log.Println("Last comment date is >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", values[0].Createddate)
	currentTime := time.Now()
	zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
	todatetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	if len(values) > 0 {
		fromdatetime := AddSubSecondsToDate(TimeParse(values[0].Createddate, ""), zonediff.UTCdiff)
		daycount = CalculateWorkingDays(page.ClientID, page.Mstorgnhirarchyid, fromdatetime, todatetime, 0)
		if daycount <= 0 {
			values[0].Daycount = 0
		} else {
			values[0].Daycount = (daycount - 1)
		}

		//values[0].Daycount = (daycount - 1)
		return values, true, err, ""
	} else {
		recordcreatedate, err := dataAccess.GetRecordcreatedate(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		if err != nil {
			logger.Log.Println("database connection failure", err)
			return t, false, err, "Something Went Wrong"
		}
		fromdatetime := AddSubSecondsToDate(TimeParse(recordcreatedate, ""), zonediff.UTCdiff)
		daycount = CalculateWorkingDays(page.ClientID, page.Mstorgnhirarchyid, fromdatetime, todatetime, 0)
		value := entities.Customervisiblecomment{}
		if daycount <= 0 {
			value.Daycount = 0
		} else {
			value.Daycount = (daycount - 1)
		}
		//value.Daycount = (daycount - 1)
		value.ClientID = page.ClientID
		value.Mstorgnhirarchyid = page.Mstorgnhirarchyid
		err = dataAccess.UpdateWorknotedaycount(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, (daycount - 1))
		if err != nil {
			log.Println("error is ----->", err)
		}
		t = append(t, value)
		return t, true, err, ""
	}

}*/

/*func Customervisiblecomment(page *entities.RecordcommonEntity) ([]entities.Customervisiblecomment, bool, error, string) {
	logger.Log.Println("In side GetAttachmentfiles")
	t := []entities.Customervisiblecomment{}
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

	var daycount int64
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
	values, err := dataAccess.GetLastRecordcomment(page)
	if err != nil {
		logger.Log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}
	currentTime := time.Now()
	zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
	todatetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	res, err := dataAccess.Getrecorddetails(page.RecordID)
	//res.SupportgroupId = page.Usergroupid
	if err != nil {
		logger.Log.Println(err)
		return t, false, err, "Something Went Wrong"
	}
	returnValue, _, _, _ := SLACriteriaRespResl(page.ClientID, page.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)
	if len(values) > 0 {

		//logger.Log.Println("Last comment date is >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", values[0].Createddate)
		fromdatetime := AddSubSecondsToDate(TimeParse(values[0].Createddate, ""), zonediff.UTCdiff)
		//logger.Log.Println("fromdatetime value  is >>>>>>>>>>>>>>>>>11>>>>>>>>>>>>>>>", fromdatetime)
		//logger.Log.Println("todatetime value  is >>>>>>>>>>>>>>>>>11>>>>>>>>>>>>>>>", todatetime)
		//logger.Log.Println("Supportgroupspecific value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", returnValue.Supportgroupspecific)
		if returnValue.Supportgroupspecific == 1 {
			daycount = CalculateWorkingDays(page.ClientID, page.Mstorgnhirarchyid, fromdatetime, todatetime, 0, returnValue.Supportgroupspecific, page.Usergroupid)
			if daycount <= 0 {
				values[0].Daycount = 0
			} else {
				values[0].Daycount = (daycount - 1)
			}
		} else {
			daycount = CalculateWorkingDays(page.ClientID, page.Mstorgnhirarchyid, fromdatetime, todatetime, 0, 0, page.Usergroupid)
			if daycount <= 0 {
				values[0].Daycount = 0
			} else {
				values[0].Daycount = (daycount - 1)
			}
		}

		//values[0].Daycount = (daycount - 1)
		logger.Log.Println("Daycount value  is >>>>>>>>>>>>>>>>>11>>>>>>>>>>>>>>>", values[0].Daycount)
		return values, true, err, ""
	} else {
		//logger.Log.Println("data value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		recordcreatedate, err := dataAccess.GetRecordcreatedate(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		if err != nil {
			logger.Log.Println("database connection failure", err)
			return t, false, err, "Something Went Wrong"
		}
		//logger.Log.Println("recordcreatedate value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", recordcreatedate)
		value := entities.Customervisiblecomment{}
		fromdatetime := AddSubSecondsToDate(TimeParse(recordcreatedate, ""), zonediff.UTCdiff)
		//logger.Log.Println("fromdatetime value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", fromdatetime)
		//logger.Log.Println("todatetime value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", todatetime)
		//logger.Log.Println("Supportgroupspecific value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", returnValue.Supportgroupspecific)
		if returnValue.Supportgroupspecific == 1 {
			daycount = CalculateWorkingDays(page.ClientID, page.Mstorgnhirarchyid, fromdatetime, todatetime, 0, returnValue.Supportgroupspecific, page.Usergroupid)
			if daycount <= 0 {
				value.Daycount = 0
			} else {
				value.Daycount = (daycount - 1)
			}
		} else {
			daycount = CalculateWorkingDays(page.ClientID, page.Mstorgnhirarchyid, fromdatetime, todatetime, 0, 0, page.Usergroupid)
			if daycount <= 0 {
				value.Daycount = 0
			} else {
				value.Daycount = (daycount - 1)
			}
		}

		//value.Daycount = (daycount - 1)
		value.ClientID = page.ClientID
		value.Mstorgnhirarchyid = page.Mstorgnhirarchyid
		logger.Log.Println("todatetime value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", daycount)
		var finaldaycount int64
		if daycount == 0 {
			finaldaycount = daycount
		} else {
			finaldaycount = daycount
		}
		err = dataAccess.UpdateWorknotedaycount(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, finaldaycount)
		if err != nil {
			log.Println("error is ----->", err)
		}
		t = append(t, value)
		return t, true, err, ""
	}

}*/

func Customervisiblecomment(page *entities.RecordcommonEntity) ([]entities.Customervisiblecomment, bool, error, string) {
	logger.Log.Println("In side GetAttachmentfiles")
	t := []entities.Customervisiblecomment{}
	var daycount int64
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
	values, err := dataAccess.GetLastRecordcomment(page)
	if err != nil {
		logger.Log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}
	currentTime := time.Now()
	zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
	todatetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	if len(values) > 0 {
		if len(values[0].Createddate) > 0 {
			fromdatetime := AddSubSecondsToDate(TimeParse(values[0].Createddate, ""), zonediff.UTCdiff)
			worknotcalendarsec := SubtractDateToDate(todatetime, fromdatetime)
			logger.Log.Println("Daycount value  is >>>>>>>>>>>>>>>>>11>>>>>>>>>>>>>>>", page.RecordID, worknotcalendarsec)
			if worknotcalendarsec > 0 {
				daycount = (worknotcalendarsec / (24 * 60 * 60))
				values[0].Daycount = daycount
			}
		}
		logger.Log.Println("Daycount value  is >>>>>>>>>>>>>>>>>11>>>>>>>>>>>>>>>", page.RecordID, values[0].Daycount)
		return values, true, err, ""
	} else {
		_, recordseq, _, err := dataAccess.Getcurrentsatusid(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		var enddate string
		recordcreatedate, err := dataAccess.GetRecordcreatedate(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		if err != nil {
			logger.Log.Println("database connection failure", err)
			return t, false, err, "Something Went Wrong"
		}

		value := entities.Customervisiblecomment{}
		if len(recordcreatedate) > 0 {
			fromdatetime := AddSubSecondsToDate(TimeParse(recordcreatedate, ""), zonediff.UTCdiff)
			var worknotcalendarsec1 int64
			if recordseq == 3 {
				enddate, err = dataAccess.GetReopendate(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 3)
				if err != nil {
					logger.Log.Println("database connection failure", err)
					return t, false, err, "Something Went Wrong"
				}
				todatetime1 := AddSubSecondsToDate(TimeParse(enddate, ""), zonediff.UTCdiff)
				worknotcalendarsec1 = SubtractDateToDate(todatetime1, fromdatetime)
			} else {
				worknotcalendarsec1 = SubtractDateToDate(todatetime, fromdatetime)
			}
			logger.Log.Println("Daycount value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", page.RecordID, worknotcalendarsec1)
			if worknotcalendarsec1 > 0 {
				daycount = (worknotcalendarsec1 / (24 * 60 * 60))
			}
		}
		value.Daycount = daycount
		value.ClientID = page.ClientID
		value.Mstorgnhirarchyid = page.Mstorgnhirarchyid
		logger.Log.Println("todatetime value  is >>>>>>>>>>>>>>>>>22>>>>>>>>>>>>>>>", page.RecordID, daycount)
		err = dataAccess.UpdateWorknotedaycount(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, daycount)
		if err != nil {
			log.Println("error is ----->", err)
		}
		t = append(t, value)
		return t, true, err, ""
	}

}

func DropAttachmentfiles(page *entities.RecordAttachmentfiles) (bool, error, string) {
	logger.Log.Println("In side DropAttachmentfiles")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// dbcon, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("Error in DBConnection in side InsertMultipleTermvalues")
	// 	return false, err, "Something Went Wrong"
	// }
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	tx, err := db.Begin()
	if err != nil {
		//	dbcon.Close()
		logger.Log.Println("Transaction creation error in InsertMultipleTermvalues", err)
		return false, err, "Something Went Wrong"
	}

	err = dao.Updateattachfiles(tx, page.ID)
	if err != nil {
		tx.Rollback()
		//	dbcon.Close()
		logger.Log.Println("Trnordertarcking update failure", err)
		return false, err, "Something Went Wrong"
	}
	err = dao.Deletefromactivitylogs(tx, page.RecordID, page.Originalname, page.Createdbyid, page.Createdgrpid, page.RecordtermID)
	if err != nil {
		tx.Rollback()
		//	dbcon.Close()
		logger.Log.Println("Activitylogs update failure", err)
		return false, err, "Something Went Wrong"
	}

	err = dao.InsertActivityLogs(tx, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 13, "Attachment name :: "+page.Originalname, page.Userid, page.Usergroupid)
	if err != nil {
		log.Println("error is ----->", err)
		tx.Rollback()
		//	dbcon.Close()
		return false, err, "Something Went Wrong"
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		//	dbcon.Close()
		return false, err, "Something Went Wrong"
	}
	//Email Start Here
	go FileDeleteEmail(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	//Email End Here
	//dbcon.Close()
	return true, err, "Attachment Deleted Successfully"
}

func GetTermvaluebysequance(page *entities.RecordcommonEntity) ([]entities.Recordtermseqvalue, bool, error, string) {
	logger.Log.Println("In side GetAttachmentfiles")
	t := []entities.Recordtermseqvalue{}
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
	values, err1 := dataAccess.GetTermvaluebysequance(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	util, err1 := dataAccess.Gettimediff(page.ClientID, page.Mstorgnhirarchyid)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	if util.Timediff > 0 {
		for i := 0; i < len(values); i++ {
			v := entities.Recordtermseqvalue{}
			v.Recordtermvalue = values[i].Recordtermvalue
			v.Name = values[i].Name
			v.Createddate = values[i].Createddate
			v.Showcreatedate = dao.Convertdate(values[i].Createddate, util.Timediff)
			t = append(t, v)
		}
	}

	return t, true, err, ""
}

func Parentchildcollaborationlogs(page *entities.Activitylogsearchcriteria) ([]entities.NewActivitylogsEntity, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	//m := make(map[string][]entities.NewActivitylogsEntity)
	t := []entities.NewActivitylogsEntity{}
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
	util, err1 := dataAccess.Gettimediff(page.ClientID, page.Mstorgnhirarchyid)

	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	diffid, err := dataAccess.Getrecordtypediffid(page.RecordID, page.ClientID, page.Mstorgnhirarchyid)
	if err != nil {
		log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}
	childids, err := dataAccess.Getchildrecordidswithcode(page.RecordID, page.ClientID, page.Mstorgnhirarchyid, 2, diffid, page.Recordcode)
	if err != nil {
		log.Println("database connection failure", err)
		return t, false, err, "Something Went Wrong"
	}
	for k := 0; k < len(childids); k++ {
		for i := 0; i < len(page.Searchbyseq); i++ {
			t, err = dataAccess.Searchtermlogsbysequance(page.ClientID, page.Mstorgnhirarchyid, childids[k].ID, page.Searchbyseq[i], t, util.Timediff, childids[k].Code)
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
			logger.Log.Println("values1 is --------------->", t)
		}
	}

	return t, true, err, ""
}

func GetRecordID(page *entities.RecordcommonEntity) (int64, bool, error, string) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	var ID int64
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return ID, false, err, "Something Went Wrong"
	// }
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return ID, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	ID, err1 := dataAccess.GetRecordID(page.ClientID, page.Mstorgnhirarchyid, page.Recordno)
	if err1 != nil {
		return ID, false, err1, "Something Went Wrong"
	}

	return ID, true, err, ""
}

func GetTabTermnames(page *entities.RecordcommonEntity) (entities.RecordTabTermsEntity, bool, error, string) {
	logger.Log.Println("In side GetTabTermnames")
	t := entities.RecordTabTermsEntity{}
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
	if page.Recordstatusid == 0 {
		values, err1 := dataAccess.GetScheduleTabTermnames(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		values1, err1 := dataAccess.GetPlanTabTermnames(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		t.ScheduleTab = values
		t.PlanTab = values1
	} else {
		values, err1 := dataAccess.GetScheduleTabTermnameswithStatus(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		values1, err1 := dataAccess.GetPlanTabTermnameswithStatus(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		t.ScheduleTab = values
		t.PlanTab = values1
	}

	return t, true, err, ""
}

func GetTabTermvalues(page *entities.RecordcommonEntity) (entities.RecordTabTermsEntity, bool, error, string) {
	logger.Log.Println("In side GetTabTermnames")
	t := entities.RecordTabTermsEntity{}
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
	if page.Recordstatusid == 0 {
		values, err1 := dataAccess.GetScheduleTabTermvalues(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		values1, err1 := dataAccess.GetPlanTabTermvalues(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}
		t.ScheduleTab = values
		t.PlanTab = values1
	} else {
		values, err1 := dataAccess.GetScheduleTabTermvalueswithStatus(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		values1, err1 := dataAccess.GetPlanTabTermvalueswithStatus(page)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}
		t.ScheduleTab = values
		t.PlanTab = values1
	}

	return t, true, err, ""
}

func Removelinkrecord(page *entities.LinkRecordEntity) (bool, error, string) {
	logger.Log.Println("In side GetAttachmentfiles")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	err1 := dataAccess.RemoverecordLink(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.LinkrecordID)
	if err1 != nil {
		return false, err1, "Something Went Wrong"
	}

	return true, err, "Record deleted successfully"
}

func GetLinkRecordsByID(page *entities.LinkRecordEntity) ([]entities.LinkRecordDetailsEntity, bool, error, string) {
	logger.Log.Println("In side GetTabTermnames")
	t := []entities.LinkRecordDetailsEntity{}
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
	values, err1 := dataAccess.GetLinkRecordsByID(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func SaveLinkRecordsByID(page *entities.LinkRecordEntity) (int64, bool, error, string) {
	logger.Log.Println("In side GetTabTermnames")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return 0, false, err, "Something Went Wrong"
	// }

	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	checkflag, err := dataAccess.CheckLinkRecordFlag(page.ClientID, page.Mstorgnhirarchyid, page.LinkrecordID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if checkflag == 0 {
		values, err1 := dataAccess.SaveLinkRecordsByID(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.LinkrecordID)
		if err1 != nil {
			return 0, false, err1, "Something Went Wrong"
		}

		err = dataAccess.InsertActivityLogs(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 16, "Child/Linked Recordno is :: "+page.Linkrecordno, page.Userid, page.Usergroupid, 0)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}

		return values, true, err, ""
	} else {
		return 0, false, err, "This records already linked with other record."
	}

}

func GetParentRecordInfoByID(page *entities.RecordcommonEntity) (entities.ParentRecordInfoEntity, bool, error, string) {
	logger.Log.Println("In side GetTabTermnames")
	t := entities.ParentRecordInfoEntity{}
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
	values, err1 := dataAccess.GetParentRecordInfo(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return values, true, err, ""
}

func GetSLADuetimeCalculation(page *entities.RecordcommonEntity) (bool, error, string) {
	logger.Log.Println("In side GetSLADuetimeCalculation")
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }

	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	recorddiffid, _, _, err := dataAccess.Getcurrentsatusid(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	seq, err1 := dataAccess.Getemeterseqno(page.ClientID, page.Mstorgnhirarchyid, recorddiffid, 2)
	if err1 != nil {
		logger.Log.Println(err1)
		return false, err, "Something Went Wrong"
	}
	if seq > 0 {

		res, err := dataAccess.Getrecorddetails(page.RecordID)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		grpID, err := dataAccess.FetchCurrentGrpID(page.RecordID)
		if err != nil {
			return false, err, "Something Went Wrong"
		}

		returnValue, _, _, _ := SLACriteriaRespResl(page.ClientID, page.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)
		if returnValue.Supportgroupspecific == 1 {
			count, err := dataAccess.GetSupportgrpdayofweekcount(page.ClientID, page.Mstorgnhirarchyid, grpID)
			if err != nil {
				return false, err, "Something Went Wrong"
			}
			if count < 7 {
				return false, err, "Day Of Week Not Properly Configured.Please Check."
			}
		} else {
			count, err := dataAccess.GetOrganizationdayofweekcount(page.ClientID, page.Mstorgnhirarchyid)
			if err != nil {
				return false, err, "Something Went Wrong"
			}
			if count < 7 {
				return false, err, "Day Of Week Not Properly Configured.Please Check."
			}
		}
		currentTime := time.Now()
		zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
		datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
		SLADueTimeCalculation(page.RecordID, 0, 1, 3, datetime.Format("2006-01-02 15:04:05"), page.ClientID, page.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "SC", grpID)

		t := entities.SLATabEntity{}
		t.ClientID = page.ClientID
		t.Mstorgnhirarchyid = page.Mstorgnhirarchyid
		t.RecordID = page.RecordID
		sladtls, _, err, _ := GetSLATabvalues(&t)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		logger.Log.Println(sladtls)
		err = dataAccess.UpdateALLSLAFields(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, sladtls.Responsedetails.Responseduetime, sladtls.Responsedetails.Responseclockstatus, sladtls.Resolutionetails.Resolutionduetime, sladtls.Resolutionetails.Resolutionclockstatus)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}

	}

	return true, err, ""
}

func GetParentrecordForIM(tz *entities.RecordcommonEntity) ([]entities.ParentticketEntity, bool, error, string) {
	t := []entities.ParentticketEntity{}
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
	values, err := dataAccess.GetParentrecordForIM(tz)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	return values, true, err, ""
}

/*func UpdateVendorTickeID(tz *entities.RecordcommonEntity) (int64, bool, error, string) {

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	hashmap, err := dataAccess.Termsequance(tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if len(hashmap) == 0 {
		return 0, false, err, "Something Went Wrong"
	}
	if tz.TermID == 0 {
		id := hashmap[tz.Termseq]
		tz.TermID = id
	}
	termnm, err := dataAccess.Gettermnamebyid(tz.TermID, tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	err = dataAccess.UpdateRecordTermvalues(tz)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 100, "Update "+termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	err = dataAccess.UpdateRecordFulldetailsVendorTicketID(tz)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	return 0, true, err, "Data updated successfully.."
}*/

func UpdateVendorTickeID(tz *entities.RecordcommonEntity) (int64, bool, error, string) {

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	hashmap, err := dataAccess.Termsequance(tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if len(hashmap) == 0 {
		return 0, false, err, "Something Went Wrong"
	}
	if tz.TermID == 0 {
		id := hashmap[tz.Termseq]
		tz.TermID = id
	}
	termnm, err := dataAccess.Gettermnamebyid(tz.TermID, tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	id, err := dataAccess.GetVendorTicketID(tz)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if id > 0 {
		err = dataAccess.UpdateRecordTermvalues(tz)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
	} else {
		_, err := dataAccess.InsertRecordTermvalues(tz)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
	}

	err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 100, "Update "+termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	err = dataAccess.UpdateRecordFulldetailsVendorTicketID(tz)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	return 0, true, err, "Data updated successfully.."
}
func UpdateSlaDueTime(tz *entities.RecordNoEntity) (int64, bool, error, string) {
	logger.Log.Println("In side Updaterecordstatus----------000000000000000000000000000---------------------------->", tz)
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	recordinfo, err := dataAccess.GetRecordidsToSlaUpdate(tz.RecordNo)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	for _, v := range recordinfo {
		logger.Log.Println("Recorid:", v.RecordID)
		res, err := dataAccess.Getrecorddetails(v.RecordID)
		if err != nil {
			logger.Log.Println(err)
			return 0, false, err, "Something Went Wrong"
		}
		Grpid, err := dataAccess.FetchFirstGrpID(v.RecordID)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}

//		fmt.Println("Grpid:", Grpid)
		zonediff, _, _, _ := Getutcdiff(v.ClientID, v.Mstorgnhirarchyid)
		datetime, _ := time.Parse("2006-01-02 15:04:05", v.Datetime)
		datetime = AddSubSecondsToDate(datetime, zonediff.UTCdiff)
		logger.Log.Println("DATETIME:", datetime)
		logger.Log.Println("Org,Record,date,grp:",v.Mstorgnhirarchyid,v.RecordID,datetime.Format("2006-01-02 15:04:05"), Grpid)
		SLADueTimeCalculation(v.RecordID, 0, 1, 3, datetime.Format("2006-01-02 15:04:05"), v.ClientID, v.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "", Grpid)
	}
	return 0, true, err, ""
}

