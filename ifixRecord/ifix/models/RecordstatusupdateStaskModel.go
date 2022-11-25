package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/dbconfig"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"

	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func UpdateStaskrecordstatus(tz *entities.RecordstatusEntity, db *sql.DB, recordtypeSeq int64) (int64, bool, error, string) {

	tx, err := db.Begin()
	logger.Log.Println(tx)
	if err != nil {
		logger.Log.Println("Transaction creation error in Updaterecordstatus", err)
		return 0, false, err, "Something Went Wrong"
	}
	dataAccess := dao.DbConn{DB: db}
	diffid, err := dataAccess.Getrecordtypediffid(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		log.Println("database connection failure", err)
		return 0, false, err, "Something Went Wrong"
	}
	parentID, err := dataAccess.Getparentrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, 2, diffid)
	if err != nil {
		log.Println("database connection failure", err)
		return 0, false, err, "Something Went Wrong"
	}
	ID, statusID, statusSeq, err := Staskrecordstatusupdation(tx, tz, db, recordtypeSeq, parentID[0])
	logger.Log.Println("STATUS ID  is ----------->", statusID)
	if err != nil {
		log.Println("Parent Record status updation failed", err)
		tx.Rollback()
		//db.Close()
		return 0, false, err, "Something Went Wrong"
	}

	logger.Log.Println("Child records length is --FOR SR --------->", len(parentID))
	logger.Log.Println("Change status value is  ----                  SR TASK             status --------------------------------------->", tz.RecordID, tz.Changestatus, parentID)
	//for i := 0; i < len(parentID); i++ {
	var previousstatus int64
	var recorddiffid int64
	if len(parentID) > 0 {
		_, previousstatus, recorddiffid, err = ParentSRrecordstatusupdation(tx, tz.ClientID, tz.Mstorgnhirarchyid, parentID[0], statusID, tz.UserID, tz.Usergroupid, db, diffid, tz, statusSeq, recordtypeSeq)
		if err != nil {
			log.Println("Child Record status updation failed", err)
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}
	}

	//}

	logger.Log.Println("ID value is -------- -->", ID)
	if ID > 0 {
		err = tx.Commit()
		if err != nil {
			log.Println("DB commit is failed", err)
			tx.Rollback()
			//db.Close()
			return 0, false, err, "Something Went Wrong"
		}
	}

	var workflowflag bool
	var errormsg string
	logger.Log.Println("previousstatus value is ----------------------------------------------------------- -->", previousstatus)
	logger.Log.Println("recorddiffid value is ----------------------------------------------------------- -->", recorddiffid)

	if statusSeq == 11 {
		err := dataAccess.Updatechildrecordflag(tz.ClientID, tz.Mstorgnhirarchyid, parentID[0], tz.RecordID)
		if err != nil {
			logger.Log.Println("database connection failure", err)

		}
	}

	if len(parentID) > 0 && previousstatus != recorddiffid && recorddiffid != -1 {
		reqbd := &entities.ParentchildEntity{}
		reqbd.Parentid = tz.RecordID
		reqbd.Childids = parentID
		reqbd.Userid = tz.UserID
		reqbd.Isupdate = false
		reqbd.Createdgroupid = tz.Usergroupid
		postBody, _ := json.Marshal(reqbd)

		logger.Log.Println("Record status request body -->", reqbd,string(postBody))

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(dbconfig.MASTER_URL+"/updatetaskstatus", "application/json", responseBody)
		if err != nil {
			logger.Log.Println("An Error Occured --->", err)
			return 0, false, err, "Something Went Wrong"
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Println("response body ------> ", err)
			return 0, false, err, "Something Went Wrong"
		}
		sb := string(body)
		wfres := entities.WorkflowResponse{}
		json.Unmarshal([]byte(sb), &wfres)
		workflowflag = wfres.Success
		errormsg = wfres.Message
		logger.Log.Println("Record status response message -->", workflowflag)
		logger.Log.Println("Record status response error message -->", errormsg)
	}
	if ID > 0 {
		//Email Notification Start Here

		logger.Log.Println("1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
		statusID, _, _, _ := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		logger.Log.Println(statusID)

		// for i := 0; i < len(parentID); i++ {
		// 	pstatusID, _, _, _ := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, parentID[i])
		// 	logger.Log.Println("Parent status id for mail notification --------------------->", pstatusID)
		// 	StatusChangeEmail(tz.ClientID, tz.Mstorgnhirarchyid, parentID[i], pstatusID)
		// }

		logger.Log.Println("1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
		//Email Notification End Here
		//db.Close()
		return ID, true, err, ""
	}
	logger.Log.Println("222222222222222222222222222222222222222222222222222222222222222222222222")
	return 0, false, err, "Something Went Wrong"
}

func Staskrecordstatusupdation(tx *sql.Tx, tz *entities.RecordstatusEntity, db *sql.DB, recordtypeSeq int64, ParentID int64) (int64, int64, int64, error) {
	logger.Log.Println("In side Mstslastatemodel")
	currentTime := time.Now()
	zonediff, _, _, _ := Getutcdiff(tz.ClientID, tz.Mstorgnhirarchyid)
	datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	log.Println(datetime)
	dataAccess := dao.DbConn{DB: db}
	logger.Log.Println("STATUS TASK  ID  is -------------------------------------------------------------------------------->", tz.ReordstatusID)

	isApproveflag, err := dataAccess.GetIsapprovalFlag(ParentID)
	if err != nil {
		return 0, 0, 0, err
	}
	var statusID int64
	if isApproveflag == 1 && recordtypeSeq == 5 {
		_, statusseq, _, err := dataAccess.Getrecorddiffidbystateid(tz.ClientID, tz.Mstorgnhirarchyid, tz.ReordstatusID)
		logger.Log.Println("statusseq  is -------------------------------------------------------------------------------->", statusseq)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}

		if statusseq == 20 {
			statusID, err = dataAccess.FetchDifferentiationIDBySeq(tz.ClientID, tz.Mstorgnhirarchyid, 3, 29)
			if err != nil {
				logger.Log.Println(err)
				return 0, 0, 0, err
			}

		}
	}
	var recorddiffid int64
	var recorddiffseq int64
	var currentstatusname string
	if statusID > 0 {
		recorddiffid, recorddiffseq, currentstatusname, err = dataAccess.FetchDifferentiationDetailsByID(statusID)
		//logger.Log.Println("recorddiffid, recorddiffseq, currentstatusname  is -------------------------------------------------------------------------------->", recorddiffid, recorddiffseq, currentstatusname)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}
	} else {
		recorddiffid, recorddiffseq, currentstatusname, err = dataAccess.Getrecorddiffidbystateid(tz.ClientID, tz.Mstorgnhirarchyid, tz.ReordstatusID)
		//logger.Log.Println("recorddiffid, recorddiffseq, currentstatusname  is -------------------------------------------------------------------------------->", recorddiffid, recorddiffseq, currentstatusname)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}
	}

	if recorddiffid > 0 {
		laststageID, err := dataAccess.GetMaxstageID(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}
		previousstatus, previousseq, previousstatusname, err := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}

		err = dao.Updatepreviousstatus(tx, tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}
		id, err := dao.Updaterecordstatus(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffid, laststageID)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}

		if id > 0 {

			//activity log entry here
			logger.Log.Println("STASK Previous status value is ------------------------------------------------->", previousstatusname)
			logger.Log.Println("STASK Current status value is ------------------------------------------------->", currentstatusname)
			logger.Log.Println(err)
			if previousstatus != recorddiffid {
				err = dao.InsertActivityLogs(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 4, "From "+previousstatusname+" To "+currentstatusname, tz.UserID, tz.Usergroupid)
				if err != nil {
					log.Println("error is ----->", err)
					return 0, 0, 0, err
				}
				//Update Stage TBL For Status

				err = dataAccess.UpdateStageStatus(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffid, currentstatusname)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//db.Close()
					return 0, 0, 0, err
				}

				//Update Stage TBL For Status
			} else {
				//Update Stage TBL For Status

				err = dataAccess.UpdateStageStatus(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffid, currentstatusname)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//db.Close()
					return 0, 0, 0, err
				}
				//Update Stage TBL For Status
			}

			res, err := dataAccess.Getrecorddetails(tz.RecordID)
			if err != nil {
				logger.Log.Println(err)
				return 0, 0, 0, err
			}
			returnValue, _, _, _ := SLACriteriaRespResl(tz.ClientID, tz.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)

			// if recorddiffseq == 1 {
			// 	err := UpdateResponseValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
			// 	if err != nil {
			// 		logger.Log.Println("error is ----->", err)
			// 		tx.Rollback()
			// 		return 0, 0, 0, err
			// 	}

			// }

			if recorddiffseq == 2 {
				err := UpdateResponseValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq,previousseq)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}

			}

			if recorddiffseq == 3 {

				err := UpdateResolutionValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}

				err = UpdateStageResolverInfo(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.UserID, tz.Usergroupid, returnValue.Supportgroupspecific, recorddiffseq, recordtypeSeq)
				if err != nil {
					logger.Log.Println(err)
					tx.Rollback()
					return 0, 0, 0, err
				}

				/*err = UpdateCloseValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}*/
			}

			if recorddiffseq == 4 {

				err = dataAccess.UpdateFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}
			}

			if recorddiffseq == 2 && previousseq == 4 {
				err := UpdateFollowuptimetakenValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}

			}

			if recorddiffseq == 5 {
				err = dataAccess.UpdatePendinguserAction(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}
			}

			if recorddiffseq == 8 {

				err := UpdateCloseValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq,previousseq)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}
			}

			if recorddiffseq == 9 {
				err := UpdateUserreplytimetakenValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}
			}

			if recorddiffseq == 10 {
				/*err = dataAccess.UpdateReopenCount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					//db.Close()
					return 0, 0, 0, err
				}*/

				var reopentimetaken int64
				var reopentimetakencalendar int64
				resolvedt, err := dataAccess.GetReopendate(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 3)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					return 0, 0, 0, err
				}
				if len(resolvedt) > 0 {
					zonediff, _, _, _ := Getutcdiff(tz.ClientID, tz.Mstorgnhirarchyid)
					currentTime := time.Now().UTC()
					today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
					strtDay := AddSubSecondsToDate(TimeParse(resolvedt, ""), zonediff.UTCdiff)
					logger.Log.Println("reopentimetaken is ---***********************-->", tz.ClientID, tz.Mstorgnhirarchyid, strtDay, today, 0, tz.Usergroupid, returnValue.Supportgroupspecific)
					reopentimetakencalendar = SubtractDateToDate(today, strtDay)
					reopentimetaken = CalculateWorkingHourBetweenTwoDates(tz.ClientID, tz.Mstorgnhirarchyid, strtDay, today, 0, returnValue.Supportgroupspecific, tz.Usergroupid)

					logger.Log.Println("reopentimetaken  value is ----->", reopentimetaken, reopentimetakencalendar)
				}
				err = dataAccess.UpdateReopenCount1(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, reopentimetaken, reopentimetakencalendar)
				if err != nil {
					log.Println("error is ----->", err)
					dataAccess.InsertErrorByRecordID(tz.ClientID, tz.Mstorgnhirarchyid, "Status Update-323", tz.RecordID, err.Error())
					tx.Rollback()
					return 0, 0, 0, err
				}
			}

			Username, err := dataAccess.GetUsername(tz.UserID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return 0, 0, 0, err
			}

			err = dataAccess.UpdateUserInfo(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.UserID, Username)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return 0, 0, 0, err
			}

			seq, err := dataAccess.Getemeterseqno(tz.ClientID, tz.Mstorgnhirarchyid, recorddiffid, 1)
			if err != nil {
				logger.Log.Println(err)
				return 0, 0, 0, err
			}
			logger.Log.Println("Responsemeter sequance no ---->", seq)
			if seq > 0 {
				res, err := dataAccess.Getrecorddetails(tz.RecordID)
				res.SupportgroupId = tz.Usergroupid
				if err != nil {
					logger.Log.Println(err)
					return 0, 0, 0, err
				}
				slaid, err := dataAccess.GetSLAdataexist(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
				if err != nil {
					logger.Log.Println(err)
					return 0, 0, 0, err
				}
				if slaid > 0 {
					GetSLAResolution(&res)
					if seq == 4 {
						flag, err := dataAccess.UpdateResponseEndFlag(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, currentTime.Format("2006-01-02 15:04:05"))
						logger.Log.Println(flag)
						if err != nil {
							logger.Log.Println(err)
							return 0, 0, 0, err
						}
					}
				}

			}
			if recordtypeSeq == 3 {
				seq, err1 := dataAccess.Getemeterseqno(tz.ClientID, tz.Mstorgnhirarchyid, recorddiffid, 2)
				if err1 != nil {
					logger.Log.Println(err1)
					return 0, 0, 0, err
				}
				logger.Log.Println("Resolutionmeter sequance no ---->", seq)
				if seq > 0 {
					res, err := dataAccess.Getrecorddetails(tz.RecordID)
					res.SupportgroupId = tz.Usergroupid
					if err != nil {
						logger.Log.Println(err)
						return 0, 0, 0, err
					}
					historyrecord, err := dataAccess.GetLatesttrnhistory(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
					if err != nil {
						logger.Log.Println(err)
						return 0, 0, 0, err
					}
					if seq != 1 {
						// Change in 15.05.2021 -----------------------------

						histrn := entities.TrnslaentityhistoryEntity{}
						histrn.Clientid = tz.ClientID
						histrn.Mstorgnhirarchyid = tz.Mstorgnhirarchyid
						histrn.Therecordid = tz.RecordID
						if historyrecord.Slastartstopindicator == 2 && seq == 2 {
							var dt = historyrecord.Recorddatetime
							var dttime = historyrecord.Recorddatetoint
							histrn.Recorddatetime = dt
							histrn.Recorddatetoint = dttime
						} else {
							histrn.Recorddatetime = datetime.Format("2006-01-02 15:04:05")
							histrn.Recorddatetoint = TimeParse(datetime.Format("2006-01-02 15:04:05"), "").Unix()
						}
						histrn.Slastartstopindicator = seq

						trnid, err := dataAccess.InsertTrnslaentityhistory(&histrn)
						if err != nil {
							logger.Log.Println(err)
							return 0, 0, 0, err
						}
						logger.Log.Println("history table id---->", trnid)

					}
					if seq == 2 {
						_, _, err, _ = GetSLAResolution(&res)
						if err != nil {
							logger.Log.Println(err)
							return 0, 0, 0, err
						}
					}
					if seq == 4 {
						_, err, _ := UpdateRessolutionEndFlag(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
						if err != nil {
							logger.Log.Println(err)
							return 0, 0, 0, err
						}
						_, _, err, _ = GetSLAResolution(&res)
						if err != nil {
							logger.Log.Println(err)
							return 0, 0, 0, err
						}
					}
					logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
					if seq == 1 || seq == 5 || seq == 3 || seq == 20 {
						grpID, err := dataAccess.FetchCurrentGrpID(tz.RecordID)
						if err != nil {
							return 0, 0, 0, err
						}
						returnValue, _, _, _ := SLACriteriaRespResl(tz.ClientID, tz.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)
						if returnValue.Supportgroupspecific == 1 {
							count, err := dataAccess.GetSupportgrpdayofweekcount(tz.ClientID, tz.Mstorgnhirarchyid, grpID)
							if err != nil {
								return 0, 0, 0, err
							}
							if count < 7 {
								return 0, 0, 0, errors.New("Day Of Week Not Properly Configured.Please Check.")
							}
						} else {
							count, err := dataAccess.GetOrganizationdayofweekcount(tz.ClientID, tz.Mstorgnhirarchyid)
							if err != nil {
								return 0, 0, 0, err
							}
							if count < 7 {
								return 0, 0, 0, errors.New("Day Of Week Not Properly Configured.Please Check.")
							}
						}
						SLADueTimeCalculation(tz.RecordID, 0, 1, 3, datetime.Format("2006-01-02 15:04:05"), tz.ClientID, tz.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "", grpID)
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

				}
			}

			return id, recorddiffid, recorddiffseq, err
		} else {
			return 0, 0, 0, err
		}
	} else {
		return 0, 0, 0, err
	}

}

func ParentSRrecordstatusupdation(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, Recordid int64, statusID int64, UserID int64, Usergroupid int64, db *sql.DB, TypeID int64, tz *entities.RecordstatusEntity, statusSeq int64, recordtypeSeq int64) (int64, int64, int64, error) {
	logger.Log.Println("In side Mstslastatemodel")
	//Condition For SLA
	currentTime := time.Now()
	zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
	datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	log.Println(datetime)

	dataAccess := dao.DbConn{DB: db}
	recordtypeSeq, err := dataAccess.Getrecordtypeseq(ClientID, Mstorgnhirarchyid, Recordid)
	log.Println(recordtypeSeq)
	if err != nil {
		logger.Log.Println(err)
		return 0, 0, 0, err
	}

	if statusSeq == 8 { //   || statusSeq == 11
		return 0, 0, -1, err
	} else {
		//check status priority
		pstatusID, priority, err := dataAccess.GetStatusPriority(ClientID, Mstorgnhirarchyid, Recordid, TypeID, tz.RecordID)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}

		updatepriority, err := dataAccess.GetUpdateStatusPriority(ClientID, Mstorgnhirarchyid, TypeID, statusID)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}

		logger.Log.Println("Priority status id & priority id values are ------------------------>", pstatusID, priority) // 2340,5
		logger.Log.Println("updatepriority priority id values are ------------------------>", updatepriority)            //0
		logger.Log.Println("Parameter staus id value is   ------------------------>", statusID)                          //2342

		if priority > 0 && priority < updatepriority || priority == updatepriority || updatepriority == 0 {

			statusID = pstatusID
		}
		logger.Log.Println("After priority set staus id value is   ------------------------>", statusID)

		// check status priority
		var recorddiffid int64
		var seqno int64
		var currentstatusname string

		_, currentseq, _, err := dataAccess.Getrecorddifferation(ClientID, Mstorgnhirarchyid, statusSeq)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}
		isapproveworkflow, err := dataAccess.GetIsapprovalworkflow(Recordid)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}

		//New Addition
		ChildrecordtypeID, err := dataAccess.FetchRecordtypeID(ClientID, Mstorgnhirarchyid, tz.RecordID, 2)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}

		ParentrcordtypeID, err := dataAccess.FetchRecordtypeID(ClientID, Mstorgnhirarchyid, Recordid, 2)
		if err != nil {
			logger.Log.Println(err)
			return 0, 0, 0, err
		}
		//End Addition
		logger.Log.Println("currentseq value is ------------------------>", currentseq)
		logger.Log.Println("isapproveworkflow value is  ------------------------>", isapproveworkflow)

		if currentseq == 1 && isapproveworkflow == 0 {
			if priority < updatepriority || priority == updatepriority {
				//recorddiffid, seqno, currentstatusname, err = dataAccess.Getrecorddiffidforstask(ClientID, Mstorgnhirarchyid, statusID)
				recorddiffid, seqno, currentstatusname, err = dataAccess.GetrecorddiffidforstaskNewFromStask(ClientID, Mstorgnhirarchyid, statusID, ParentrcordtypeID, ChildrecordtypeID)
				if err != nil {
					logger.Log.Println(err)
					return 0, 0, 0, err
				}
			} else {
				recorddiffid, seqno, currentstatusname, err = dataAccess.Getrecorddifferation(ClientID, Mstorgnhirarchyid, 17)
				if err != nil {
					logger.Log.Println(err)
					return 0, 0, 0, err
				}
			}

		} else {
			//recorddiffid, seqno, currentstatusname, err = dataAccess.Getrecorddiffidforstask(ClientID, Mstorgnhirarchyid, statusID)
			recorddiffid, seqno, currentstatusname, err = dataAccess.GetrecorddiffidforstaskNewFromStask(ClientID, Mstorgnhirarchyid, statusID, ParentrcordtypeID, ChildrecordtypeID)
			if err != nil {
				logger.Log.Println(err)
				return 0, 0, 0, err
			}
		}
		logger.Log.Println("recorddiffid value is ------------------------>", recorddiffid)
		logger.Log.Println("seqno value is  ------------------------>", seqno)
		logger.Log.Println("currentstatusname value is  ------------------------>", currentstatusname)

		if recorddiffid > 0 {

			previousstatus, currentseqno, previousstatusname, err := dataAccess.Getcurrentsatusid(ClientID, Mstorgnhirarchyid, Recordid)
			log.Println(currentseqno)
			if err != nil {
				logger.Log.Println(err)
				return 0, 0, 0, err
			}
			if previousstatus != recorddiffid {
				laststageID, err := dataAccess.GetMaxstageID(ClientID, Mstorgnhirarchyid, Recordid)
				if err != nil {
					logger.Log.Println(err)
					return 0, 0, 0, err
				}
				err = dao.Updatepreviousstatus(tx, Recordid, ClientID, Mstorgnhirarchyid)
				if err != nil {
					logger.Log.Println(err)
					return 0, 0, 0, err
				}
				// End New logic added in 14.06.2021
				id, err := dao.Updaterecordstatus(tx, ClientID, Mstorgnhirarchyid, Recordid, recorddiffid, laststageID)
				if err != nil {
					logger.Log.Println(err)
					return 0, 0, 0, err
				}

				err = dao.UpdateTaskflag(tx, ClientID, Mstorgnhirarchyid, Recordid, recorddiffid)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}
				logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.++++++++++++++++++++++++++++++++++++++++++++++++++++", id)
				if id > 0 {
					//activity log entry here
					if previousstatus != recorddiffid {
						/*err = dao.InsertActivityLogs(tx, ClientID, Mstorgnhirarchyid, Recordid, 4, "From "+previousstatusname+" To "+currentstatusname, UserID, Usergroupid)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							return 0, 0, 0, err
						}*/

						logger.Log.Println("Record type sequance is ------------------------->", recordtypeSeq)

						if recordtypeSeq == 2 {
							currentlogseq, err := dataAccess.GetLatestActivitylogSeq(Recordid)
							if err != nil {
								log.Println("error is ----->", err)
								tx.Rollback()
								return 0, 0, 0, err
							}
							logger.Log.Println("Current activity log sequance is ------------------------->", currentlogseq)
							if currentlogseq != 6 {
								err = dao.InsertActivityLogs(tx, ClientID, Mstorgnhirarchyid, Recordid, 4, "From "+previousstatusname+" To "+currentstatusname, UserID, Usergroupid)
								if err != nil {
									logger.Log.Println("error is ----->", err)
									return 0, 0, 0, err
								}
							}

						} else {
							err = dao.InsertActivityLogs(tx, ClientID, Mstorgnhirarchyid, Recordid, 4, "From "+previousstatusname+" To "+currentstatusname, UserID, Usergroupid)
							if err != nil {
								logger.Log.Println("error is ----->", err)
								return 0, 0, 0, err
							}
						}

						//Update Stage TBL For Status

						err = dataAccess.UpdateStageStatus(ClientID, Mstorgnhirarchyid, Recordid, recorddiffid, currentstatusname)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							tx.Rollback()
							return 0, 0, 0, err
						}
						//Update Stage TBL For Status

					} else {
						//Update Stage TBL For Status

						err = dataAccess.UpdateStageStatus(ClientID, Mstorgnhirarchyid, Recordid, recorddiffid, currentstatusname)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							tx.Rollback()
							//db.Close()
							return 0, 0, 0, err
						}

						//Update Stage TBL For Status
					}

					// res, err := dataAccess.Getrecorddetails(Recordid)
					// if err != nil {
					// 	logger.Log.Println(err)
					// 	return 0, 0, 0, err
					// }
					// returnValue, _, _, _ := SLACriteriaRespResl(ClientID, Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)

					// if seqno == 16 {
					// 	err := UpdateResponseValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, seqno, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
					// 	if err != nil {
					// 		logger.Log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, 0, 0, err
					// 	}

					// }

					
					if seqno == 3 {
						err = UpdateStageResolverInfoForSR(db, ClientID, Mstorgnhirarchyid, Recordid, UserID, Usergroupid)
						if err != nil {
							logger.Log.Println(err)
							tx.Rollback()
							return 0, 0, 0, err
						}
					}
					// if seqno == 4 {

					// 	err = dao.UpdateFollowupcount(tx, ClientID, Mstorgnhirarchyid, Recordid)
					// 	if err != nil {
					// 		log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, 0, 0, err
					// 	}
					// }

					// if seqno == 5 {
					// 	err = dao.UpdatePendinguserAction(tx, ClientID, Mstorgnhirarchyid, Recordid)
					// 	if err != nil {
					// 		log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, 0, 0, err
					// 	}
					// }

					// if seqno == 8 {
					// 	err := UpdateCloseValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, seqno, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
					// 	if err != nil {
					// 		logger.Log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, 0, 0, err
					// 	}
					// }

					// if seqno == 10 {
					// 	err = dao.UpdateReopenCount(tx, ClientID, Mstorgnhirarchyid, Recordid)
					// 	if err != nil {
					// 		log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, 0, 0, err
					// 	}
					// }

					Username, err := dataAccess.GetUsername(UserID)
					if err != nil {
						logger.Log.Println(err)
						tx.Rollback()
						return 0, 0, 0, err
					}

					err = dataAccess.UpdateUserInfo(ClientID, Mstorgnhirarchyid, Recordid, UserID, Username)
					if err != nil {
						logger.Log.Println(err)
						tx.Rollback()
						return 0, 0, 0, err
					}

					logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.++++++++++++++++++++++++++++++++++++++++++++++++++++", id)
					//activity log entry end

					//End
					return id, previousstatus, recorddiffid, err
				} else {
					return 0, 0, 0, err
				}

			} else {
				return 0, 0, 0, err
			}
		} else {
			return 0, 0, 0, err
		}
	}

}
