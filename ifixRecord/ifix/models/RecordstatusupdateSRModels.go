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
	"net/http"

	"log"
	"time"
)

func UpdateSRrecordstatus(tz *entities.RecordstatusEntity, db *sql.DB, recordtypeSeq int64) (int64, bool, error, string) {
	tx, err := db.Begin()
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
	var childids []int64
	logger.Log.Println("tz.Issrrequestor value  is --++++++++++++++++++++++++++++++++++++++++++++++++++++++--------->", tz.Issrrequestor)
	logger.Log.Println("tz.Changestatus value  is --++++++++++++++++++++++++++++++++++++++++++++++++++++++--------->", tz.Changestatus)
	if tz.Issrrequestor == 0 {
		childids, err = dataAccess.Getchildrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, 2, diffid)
		if err != nil {
			log.Println("database connection failure", err)
			return 0, false, err, "Something Went Wrong"
		}
	} else {
		childids, err = dataAccess.GetchildrecordidForPendingVendorActions(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, 2, diffid)
		if err != nil {
			log.Println("database connection failure", err)
			return 0, false, err, "Something Went Wrong"
		}
	}

	ID, statusID, statusSeq, err := SRrecordstatusupdation(tx, tz, db, recordtypeSeq)
	logger.Log.Println("STATUS ID  is ----------->", statusID, statusSeq)
	if err != nil {
		log.Println("Parent Record status updation failed", err)
		tx.Rollback()
		//db.Close()
		return 0, false, err, "Something Went Wrong"
	}

	if ID == -1 {
		return ID, true, err, ""
	}

	logger.Log.Println("Child records length is --FOR SR --------->", len(childids))
	logger.Log.Println("Change status value is  ----                  SR status --------------------------------------->", tz.RecordID, childids)
	if tz.Changestatus == 0 {
		for i := 0; i < len(childids); i++ {
			//change in 19.10.2022
			_, currentseqno, _, err := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, childids[i])
			if err != nil {
				log.Println("Child Record status updation failed", err)
				tx.Rollback()
				//db.Close()
				return 0, false, err, "Something Went Wrong"
			}
			//change in 19.10.2022
			logger.Log.Println("Current Sequance value is  --------------------+++++++++++++++++++++++++++++++++++++++++++>", currentseqno)
			if currentseqno != 8 && currentseqno != 11 && currentseqno != 14 {
			logger.Log.Println(" In Side If condition Current Sequance value is  --------------------+++++++++++++++++++++++++++++++++++++++++++>", currentseqno, childids[i])
			_, err = ChildSRrecordstatusupdation(tx, tz.ClientID, tz.Mstorgnhirarchyid, childids[i], statusID, tz.UserID, tz.Usergroupid, db, tz.RecordID, statusSeq, recordtypeSeq)
			if err != nil {
				log.Println("Child Record status updation failed", err)
				tx.Rollback()
				//db.Close()
				return 0, false, err, "Something Went Wrong"
			   }
		      }
		}
	}

	if ID > 0 {
		err = tx.Commit()
		if err != nil {
			log.Println("DB commit is failed", err)
			tx.Rollback()
			return 0, false, err, "Something Went Wrong"
		}
	}

	var workflowflag bool
	var errormsg string
	if len(childids) > 0 && tz.Changestatus == 0 {
		reqbd := &entities.ParentchildEntity{}
		reqbd.Parentid = tz.RecordID
		reqbd.Childids = childids
		reqbd.Userid = tz.UserID
		reqbd.Isupdate = false
		reqbd.Createdgroupid = tz.Usergroupid
		postBody, _ := json.Marshal(reqbd)

		logger.Log.Println("Record status request body -->", reqbd)

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
	logger.Log.Println("ID value is -------- -->", ID)
	if ID > 0 {
		//Email Notification Start Here
		logger.Log.Println("1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
		statusID, Seq, _, _ := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		logger.Log.Println(statusID)
		go StatusChangeEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, statusID, Seq)
		return ID, true, err, ""
	}

	return 0, false, err, "Something Went Wrong"
}

func SRrecordstatusupdation(tx *sql.Tx, tz *entities.RecordstatusEntity, db *sql.DB, recordtypeSeq int64) (int64, int64, int64, error) {
	logger.Log.Println("In side Mstslastatemodel")
	dataAccess := dao.DbConn{DB: db}
	// ---------------- New Addition 19.08.2021 --------------------------------------

	//check status priority
	taskflag, err := dataAccess.GetcurrentTaskflag(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		logger.Log.Println(err)
		return 0, 0, 0, err
	}
	// ---------------------------------------------------------------------------------

	recorddiffid, recorddiffseq, currentstatusname, err := dataAccess.Getrecorddiffidbystateid(tz.ClientID, tz.Mstorgnhirarchyid, tz.ReordstatusID)
	if err != nil {
		logger.Log.Println(err)
		return 0, 0, 0, err
	}
	logger.Log.Println("Task flag value is ------------------------------------->", taskflag)
	logger.Log.Println("Task flag value is ------------------------------------->", recorddiffseq)
	if recorddiffid > 0 {
		//if taskflag == 0 || recorddiffseq == 3 || recorddiffseq == 14 || recorddiffseq == 9 || recorddiffseq == 10 || recorddiffseq == 28 || recorddiffseq == 26 || recorddiffseq == 11 || recorddiffseq == 16 || recorddiffseq == 4 || recorddiffseq == 5 {
		if taskflag == 0 || recorddiffseq > 0 {
			laststageID, err := dataAccess.GetMaxstageID(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if err != nil {
				logger.Log.Println(err)
				return 0, 0, 0, err
			}
			previousstatus, _, previousstatusname, err := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
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

			isApproveworkflow, err := dataAccess.GetIsapprovalworkflow(tz.RecordID)
			if err != nil {
				return 0, 0, 0, err
			}
			log.Println("isApproveworkflow is ----->", isApproveworkflow)

			if recordtypeSeq == 2 && recorddiffseq == 15 || recorddiffseq == 17 { //|| isApproveworkflow == 1

				err = dao.UpdateApproveflag(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}
			}

			if recordtypeSeq == 4 && recorddiffseq == 26 { // || isApproveworkflow == 1
				err = dao.UpdateApproveflag(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
				if err != nil {
					log.Println("error is ----->", err)
					tx.Rollback()
					return 0, 0, 0, err
				}
			}

			// res, err := dataAccess.Getrecorddetails(tz.RecordID)
			// if err != nil {
			// 	logger.Log.Println(err)
			// 	return 0, 0, 0, err
			// }
			//returnValue, _, _, _ := SLACriteriaRespResl(tz.ClientID, tz.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)

			if id > 0 {

				//activity log entry here
				if previousstatus != recorddiffid {
					/*err = dao.InsertActivityLogs(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 4, "From "+previousstatusname+" To "+currentstatusname, tz.UserID, tz.Usergroupid)
					if err != nil {
						log.Println("error is ----->", err)
						return 0, 0, 0, err
					}*/
					logger.Log.Println("Record type sequance is ------------------------->", recordtypeSeq)

					if recordtypeSeq == 2 {
						currentlogseq, err := dataAccess.GetLatestActivitylogSeq(tz.RecordID)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							return 0, 0, 0, err
						}
						logger.Log.Println("Current activity log sequance is ------------------------->", currentlogseq)

						if currentlogseq != 6 {
							err = dao.InsertActivityLogs(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 4, "From "+previousstatusname+" To "+currentstatusname, tz.UserID, tz.Usergroupid)
							if err != nil {
								log.Println("error is ----->", err)
								return 0, 0, 0, err
							}
						}

					} else {
						err = dao.InsertActivityLogs(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 4, "From "+previousstatusname+" To "+currentstatusname, tz.UserID, tz.Usergroupid)
						if err != nil {
							log.Println("error is ----->", err)
							return 0, 0, 0, err
						}
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

					// if recorddiffseq == 16 {
					// 	err := UpdateResponseValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
					// 	if err != nil {
					// 		logger.Log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, 0, 0, err
					// 	}

					// }

					if recorddiffseq == 3 {
						err = UpdateStageResolverInfoForSR(db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.UserID, tz.Usergroupid)
						if err != nil {
							logger.Log.Println(err)
							tx.Rollback()
							return 0, 0, 0, err
						}
					}
					// if recorddiffseq == 4 {
					// 	//
					// 	err = dao.UpdateFollowupcount(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
					// 	if err != nil {
					// 		log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		//db.Close()
					// 		return 0, 0, 0, err
					// 	}
					// }

					// if recorddiffseq == 5 {
					// 	err = dao.UpdatePendinguserAction(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
					// 	if err != nil {
					// 		log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		//db.Close()
					// 		return 0, 0, 0, err
					// 	}
					// }

					// if recorddiffseq == 8 {
					// 	err := UpdateCloseValueinStagetbl(tx, db, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, recorddiffseq, tz.Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
					// 	if err != nil {
					// 		logger.Log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, 0, 0, err
					// 	}
					// }

					// if recorddiffseq == 10 {
					// 	err = dao.UpdateReopenCount(tx, tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
					// 	if err != nil {
					// 		log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		//db.Close()
					// 		return 0, 0, 0, err
					// 	}

					// 	// For SR Calling update response datetime
					// }

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

					//Update Stage TBL For Status
				}

				return id, recorddiffid, recorddiffseq, err
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

func ChildSRrecordstatusupdation(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, Recordid int64, statusID int64, UserID int64, Usergroupid int64, db *sql.DB, ParentRecordID int64, statusSeq int64, parentrecordtypeSeq int64) (int64, error) {
	logger.Log.Println("statusID  is -----------------category change         --------------------------------------------------------------->", statusID)
	logger.Log.Println("statusSeq  is ------------------------------- category change     ------------------------------------------------->", statusSeq)
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
		return 0, err
	}

	//New Addition
	ChildrecordtypeID, err := dataAccess.FetchRecordtypeID(ClientID, Mstorgnhirarchyid, Recordid, 2)
	if err != nil {
		logger.Log.Println(err)
		return 0, err
	}

	ParentrcordtypeID, err := dataAccess.FetchRecordtypeID(ClientID, Mstorgnhirarchyid, ParentRecordID, 2)
	if err != nil {
		logger.Log.Println(err)
		return 0, err
	}
	//End Addition
	var changestatusID int64
	logger.Log.Println("parentrecordtypeSeq  is -------------------------------------------------------------------------------->", parentrecordtypeSeq)
	logger.Log.Println("statusSeq  is -------------------------------------------------------------------------------->", statusSeq)
	if parentrecordtypeSeq == 4 || parentrecordtypeSeq == 2 {
		if statusSeq == 11 {
			_, currentseqno, _, err := dataAccess.Getcurrentsatusid(ClientID, Mstorgnhirarchyid, Recordid)
			if err != nil {
				logger.Log.Println(err)
				return 0, err
			}
			logger.Log.Println("currentseqno  is -------------------------------------------------------------------------------->", currentseqno)
			if currentseqno == 3 {
				changestatusID, err = dataAccess.FetchDifferentiationIDBySeq(ClientID, Mstorgnhirarchyid, 3, 8)
				if err != nil {
					logger.Log.Println(err)
					return 0, err
				}
			} else {
				changestatusID, err = dataAccess.FetchDifferentiationIDBySeq(ClientID, Mstorgnhirarchyid, 3, 11)
				if err != nil {
					logger.Log.Println(err)
					return 0, err
				}
			}
		}
	}

	var recorddiffid int64
	var seqno int64
	var currentstatusname string
	logger.Log.Println("changestatusID  is -------------------------------------------------------------------------------->", changestatusID)
	if changestatusID > 0 {
		recorddiffid, seqno, currentstatusname, err = dataAccess.FetchDifferentiationDetailsByID(changestatusID)
		if err != nil {
			logger.Log.Println(err)
			return 0, err
		}
	} else {
		logger.Log.Println("ClientID value is ----------------------1111111111111111111111------------------------->", ClientID)
		logger.Log.Println("Mstorgnhirarchyid value is ----------------------------------------------->", Mstorgnhirarchyid)
		logger.Log.Println("statusID value is ----------------------------------------------->", statusID)
		logger.Log.Println("ParentrcordtypeID value is ----------------------------------------------->", ParentrcordtypeID)
		logger.Log.Println("ChildrecordtypeID value is ----------------------------------------------->", ChildrecordtypeID)
		recorddiffid, seqno, currentstatusname, err = dataAccess.GetrecorddiffidforstaskNew(ClientID, Mstorgnhirarchyid, statusID, ParentrcordtypeID, ChildrecordtypeID)
		if err != nil {
			logger.Log.Println(err)
			return 0, err
		}
	}
	//logger.Log.Println("recorddiffid, seqno, currentstatusname  is -------------------------------------------------------------------------------->", recorddiffid, seqno, currentstatusname)
	// recorddiffid, seqno, currentstatusname, err = dataAccess.GetrecorddiffidforstaskNew(ClientID, Mstorgnhirarchyid, statusID, ParentrcordtypeID, ChildrecordtypeID)
	// if err != nil {
	// 	logger.Log.Println(err)
	// 	return 0, err
	// }
	if recorddiffid > 0 {

		previousstatus, currentseqno, previousstatusname, err := dataAccess.Getcurrentsatusid(ClientID, Mstorgnhirarchyid, Recordid)
		log.Println(currentseqno)
		if err != nil {
			logger.Log.Println(err)
			return 0, err
		}
		if currentseqno == 8 || currentseqno == 11 {
			return 0, err
		} else {

			laststageID, err := dataAccess.GetMaxstageID(ClientID, Mstorgnhirarchyid, Recordid)
			if err != nil {
				logger.Log.Println(err)
				return 0, err
			}
			err = dao.Updatepreviousstatus(tx, Recordid, ClientID, Mstorgnhirarchyid)
			if err != nil {
				logger.Log.Println(err)
				return 0, err
			}
			// End New logic added in 14.06.2021
			id, err := dao.Updaterecordstatus(tx, ClientID, Mstorgnhirarchyid, Recordid, recorddiffid, laststageID)
			if err != nil {
				logger.Log.Println(err)
				return 0, err
			}

			if id > 0 {
				//activity log entry here

				logger.Log.Println("Previous Status ID value is ----------------------------------------------->", previousstatus)
				logger.Log.Println("recorddiffid Status ID value is ----------------------------------------------->", recorddiffid)
				logger.Log.Println("Previous Status name value is ----------------------------------------------->", previousstatusname)
				logger.Log.Println("currentstatusname Status name value is ----------------------------------------------->", currentstatusname)

				res, err := dataAccess.Getrecorddetails(Recordid)
				if err != nil {
					logger.Log.Println(err)
					return 0, err
				}
				returnValue, _, _, _ := SLACriteriaRespResl(ClientID, Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)

				if previousstatus != recorddiffid {
					err = dao.InsertActivityLogs(tx, ClientID, Mstorgnhirarchyid, Recordid, 4, "From "+previousstatusname+" To "+currentstatusname, UserID, Usergroupid)
					if err != nil {
						log.Println("error is ----->", err)
						return 0, err
					}

					//Update Stage TBL For Status

					err = dataAccess.UpdateStageStatus(ClientID, Mstorgnhirarchyid, Recordid, recorddiffid, currentstatusname)
					if err != nil {
						log.Println("error is ----->", err)
						tx.Rollback()
						//db.Close()
						return 0, err
					}

					// if seqno == 1 {
					// 	err := UpdateResponseValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, seqno, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
					// 	if err != nil {
					// 		logger.Log.Println("error is ----->", err)
					// 		tx.Rollback()
					// 		return 0, err
					// 	}

					// }

					if seqno == 2 {
						err := UpdateResponseValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, seqno, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq,currentseqno)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							tx.Rollback()
							return 0, err
						}

					}

					if seqno == 3 {

						err := UpdateResolutionValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, seqno, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							tx.Rollback()
							return 0, err
						}

						err = UpdateStageResolverInfo(tx, db, ClientID, Mstorgnhirarchyid, Recordid, UserID, Usergroupid, returnValue.Supportgroupspecific, seqno, recordtypeSeq)
						if err != nil {
							logger.Log.Println(err)
							tx.Rollback()
							return 0, err
						}
					}

					if seqno == 4 {

						err = dataAccess.UpdateFollowupcount(ClientID, Mstorgnhirarchyid, Recordid)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							//db.Close()
							return 0, err
						}
					}

					if seqno == 2 && currentseqno == 4 {
						err := UpdateFollowuptimetakenValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							tx.Rollback()
							return 0, err
						}

					}

					if seqno == 5 {
						err = dataAccess.UpdatePendinguserAction(ClientID, Mstorgnhirarchyid, Recordid)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							//db.Close()
							return 0, err
						}
					}

					if seqno == 8 {

						err := UpdateCloseValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, seqno, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq,currentseqno)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							tx.Rollback()
							return 0, err
						}
					}

					if seqno == 9 {
						err := UpdateUserreplytimetakenValueinStagetbl(tx, db, ClientID, Mstorgnhirarchyid, Recordid, seqno, Usergroupid, returnValue.Supportgroupspecific, recordtypeSeq)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							tx.Rollback()
							return 0, err
						}
					}

					if currentseqno == 3 && seqno == 1 {
						/*err = dataAccess.UpdateReopenCount(ClientID, Mstorgnhirarchyid, Recordid)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							//db.Close()
							return 0, err
						}*/


						var reopentimetaken int64
						var reopentimetakencalendar int64
						resolvedt, err := dataAccess.GetReopendate(ClientID, Mstorgnhirarchyid, Recordid, 3)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							return 0, err
						}
						if len(resolvedt) > 0 {
							zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
							currentTime := time.Now().UTC()
							today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
							strtDay := AddSubSecondsToDate(TimeParse(resolvedt, ""), zonediff.UTCdiff)
							logger.Log.Println("reopentimetaken is ---***********************-->", ClientID, Mstorgnhirarchyid, strtDay, today, 0, Usergroupid, returnValue.Supportgroupspecific)
							reopentimetakencalendar = SubtractDateToDate(today, strtDay)
							reopentimetaken = CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, returnValue.Supportgroupspecific, Usergroupid)
							logger.Log.Println("reopentimetaken  value is ----->", reopentimetaken, reopentimetakencalendar)
						}

						err = dataAccess.UpdateReopenCount1(ClientID, Mstorgnhirarchyid, Recordid, reopentimetaken, reopentimetakencalendar)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							return 0, err
						}
					}

					//Update Stage TBL For Status

				} else {
					//Update Stage TBL For Status

					err = dataAccess.UpdateStageStatus(ClientID, Mstorgnhirarchyid, Recordid, recorddiffid, currentstatusname)
					if err != nil {
						log.Println("error is ----->", err)
						tx.Rollback()
						return 0, err
					}

					//Update Stage TBL For Status
				}

				Username, err := dataAccess.GetUsername(UserID)
				if err != nil {
					logger.Log.Println(err)
					tx.Rollback()
					return 0, err
				}

				err = dataAccess.UpdateUserInfo(ClientID, Mstorgnhirarchyid, Recordid, UserID, Username)
				if err != nil {
					logger.Log.Println(err)
					tx.Rollback()
					return 0, err
				}

				//activity log entry end

				//For Response meter status checking

				seq, err := dataAccess.Getemeterseqno(ClientID, Mstorgnhirarchyid, recorddiffid, 1)
				if err != nil {
					logger.Log.Println(err)
					return 0, err
				}
				logger.Log.Println("Responsemeter sequance no ---->", seq)
				if seq > 0 {
					res, err := dataAccess.Getrecorddetails(Recordid)
					res.SupportgroupId = Usergroupid
					if err != nil {
						logger.Log.Println(err)
						return 0, err
					}
					slaid, err := dataAccess.GetSLAdataexist(ClientID, Mstorgnhirarchyid, Recordid)
					if err != nil {
						logger.Log.Println(err)
						return 0, err
					}
					if slaid > 0 {
						GetSLAResolution(&res)
						if seq == 4 {
							flag, err := dataAccess.UpdateResponseEndFlag(ClientID, Mstorgnhirarchyid, Recordid, currentTime.Format("2006-01-02 15:04:05"))
							logger.Log.Println(flag)
							if err != nil {
								logger.Log.Println(err)
								return 0, err
							}
						}
					}

				}

				if recordtypeSeq == 3 {
					seq, err1 := dataAccess.Getemeterseqno(ClientID, Mstorgnhirarchyid, recorddiffid, 2)
					if err1 != nil {
						logger.Log.Println(err1)
						return 0, err
					}
					logger.Log.Println("Resolutionmeter sequance no ---->", seq)
					if seq > 0 {
						res, err := dataAccess.Getrecorddetails(Recordid)
						res.SupportgroupId = Usergroupid
						if err != nil {
							logger.Log.Println(err)
							return 0, err
						}
						historyrecord, err := dataAccess.GetLatesttrnhistory(ClientID, Mstorgnhirarchyid, Recordid)
						if err != nil {
							logger.Log.Println(err)
							return 0, err
						}
						if seq != 1 {
							// Change in 15.05.2021 -----------------------------

							histrn := entities.TrnslaentityhistoryEntity{}
							histrn.Clientid = ClientID
							histrn.Mstorgnhirarchyid = Mstorgnhirarchyid
							histrn.Therecordid = Recordid
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
								return 0, err
							}
							logger.Log.Println("history table id---->", trnid)

						}
						if seq == 2 {
							_, _, err, _ = GetSLAResolution(&res)
							if err != nil {
								logger.Log.Println(err)
								return 0, err
							}
						}
						if seq == 4 {
							_, err, _ := UpdateRessolutionEndFlag(ClientID, Mstorgnhirarchyid, Recordid)
							if err != nil {
								logger.Log.Println(err)
								return 0, err
							}
							_, _, err, _ = GetSLAResolution(&res)
							if err != nil {
								logger.Log.Println(err)
								return 0, err
							}
						}
						logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
						if seq == 1 || seq == 5 || seq == 3 {

							if currentseqno == 3 && seqno == 1 { // For Reopen Case
								resolvestateID, err := dataAccess.FetchResolvedStateID(ClientID, Mstorgnhirarchyid, currentseqno)
								if err != nil {
									return 0, err
								}
								grpID, err := dataAccess.FetchCurrentGrpIDForReopen(Recordid, resolvestateID)
								if err != nil {
									return 0, err
								}
								logger.Log.Println("Reopen GrpID is >>>>>>>>>>>>>>>>>>>>>>>>>>>>>", grpID)
								returnValue, _, _, _ := SLACriteriaRespResl(ClientID, Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)
								if returnValue.Supportgroupspecific == 1 {
									count, err := dataAccess.GetSupportgrpdayofweekcount(ClientID, Mstorgnhirarchyid, grpID)
									if err != nil {
										return 0, err
									}
									if count < 7 {
										return 0, errors.New("Day Of Week Not Properly Configured.Please Check.")
									}
								} else {
									count, err := dataAccess.GetOrganizationdayofweekcount(ClientID, Mstorgnhirarchyid)
									if err != nil {
										return 0, err
									}
									if count < 7 {
										return 0, errors.New("Day Of Week Not Properly Configured.Please Check.")
									}
								}
								SLADueTimeCalculation(Recordid, 0, 1, 3, datetime.Format("2006-01-02 15:04:05"), ClientID, Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "", grpID)
							} else {
								grpID, err := dataAccess.FetchCurrentGrpID(Recordid)
								if err != nil {
									return 0, err
								}
								if returnValue.Supportgroupspecific == 1 {
									count, err := dataAccess.GetSupportgrpdayofweekcount(ClientID, Mstorgnhirarchyid, grpID)
									if err != nil {
										return 0, err
									}
									if count < 7 {
										return 0, errors.New("Day Of Week Not Properly Configured.Please Check.")
									}
								} else {
									count, err := dataAccess.GetOrganizationdayofweekcount(ClientID, Mstorgnhirarchyid)
									if err != nil {
										return 0, err
									}
									if count < 7 {
										return 0, errors.New("Day Of Week Not Properly Configured.Please Check.")
									}
								}
								SLADueTimeCalculation(Recordid, 0, 1, 3, datetime.Format("2006-01-02 15:04:05"), ClientID, Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "", grpID)
							}

						}

						t := entities.SLATabEntity{}
						t.ClientID = ClientID
						t.Mstorgnhirarchyid = Mstorgnhirarchyid
						t.RecordID = Recordid
						sladtls, _, err, _ := GetSLATabvalues(&t)
						logger.Log.Println("------------------666666666666666666666666666666", sladtls)
						if err != nil {
							logger.Log.Println("------------------666666666666666666666666666666", err)
						}
						logger.Log.Println("------------------666666666666666666666666666666", err)
						err = dataAccess.UpdateSLAFields(ClientID, Mstorgnhirarchyid, Recordid, sladtls.Responsedetails.Responseduetime, sladtls.Responsedetails.Responseclockstatus, sladtls.Resolutionetails.Resolutionduetime, sladtls.Resolutionetails.Resolutionclockstatus)
						if err != nil {
							logger.Log.Println(err)
						}

					}
				}
				//End
				return id, err
			} else {
				return 0, err
			}
		}
	} else {
		return 0, err
	}

}
