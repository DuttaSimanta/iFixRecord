package models

import (
	"database/sql"
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
	"time"
)

func UpdateResponseValueinStagetbl(tx *sql.Tx, db *sql.DB, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, recorddiffseq int64, Usergroupid int64, Supportgroupspecific int64, recordtypeSeq int64,previousseq int64) error {

	logger.Log.Println("In side UpdateResponseValueinStagetbl")
	zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
	dataAccess := dao.DbConn{DB: db}
	if recordtypeSeq == 1 || recordtypeSeq == 3 {
		value, err := dataAccess.GetFirstResponseValue(ClientID, Mstorgnhirarchyid, RecordID)
		if err != nil {
			logger.Log.Println("error is ----->", err)
			return err
		}
		if value == "NA" {
			createdt, err := dataAccess.GetRecordcreatedate(ClientID, Mstorgnhirarchyid, RecordID)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
			currentTime := time.Now().UTC()
			//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
			strtDay := AddSubSecondsToDate(TimeParse(createdt, ""), zonediff.UTCdiff)
			today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)

			logger.Log.Println("ClientID is ---***********************-->", ClientID)
			logger.Log.Println("Mstorgnhirarchyid is ---*******************-->", Mstorgnhirarchyid)
			logger.Log.Println("strtDay is --****************************--->", strtDay)
			logger.Log.Println("today is ---***************************************-->", today)
			logger.Log.Println("Supportgroupspecific is ----->", Supportgroupspecific)
			logger.Log.Println("Usergroupid is ----->", Usergroupid)

			responsetiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
			log.Println("responsetiimetaken  value is ----->", responsetiimetaken)

			err = dataAccess.UpdateFirstResponse(ClientID, Mstorgnhirarchyid, RecordID, responsetiimetaken)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
		} else {
			if previousseq == 1 || previousseq == 10  {
				if recorddiffseq == 2 {
				reopendt, err := dataAccess.GetReopendate(ClientID, Mstorgnhirarchyid, RecordID, previousseq)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					return err
				}
				if len(reopendt) > 0 {
					currentTime := time.Now().UTC()
					today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
					strtDay := AddSubSecondsToDate(TimeParse(reopendt, ""), zonediff.UTCdiff)
					responsetiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
					
					logger.Log.Println("Start time taken value is --------->",strtDay)
					logger.Log.Println("Today time taken value is ---***********************-->", today)
					logger.Log.Println("responsetiimetaken  value is ----->", responsetiimetaken)

					err = dataAccess.UpdateLatestResponse(ClientID, Mstorgnhirarchyid, RecordID, responsetiimetaken)
					if err != nil {
						logger.Log.Println("error is ----->", err)
						return err
					}
				}
			    }	
			}

		}
	}else if recordtypeSeq == 5 {
		err := dataAccess.UpdateFirstResponse(ClientID, Mstorgnhirarchyid, RecordID, 0)
		if err != nil {
			logger.Log.Println("error is ----->", err)
			return err
		}
	}
	return nil

}

func UpdateResolutionValueinStagetbl(tx1 *sql.Tx, db *sql.DB, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, recorddiffseq int64, Usergroupid int64, Supportgroupspecific int64, recordtypeSeq int64) error {
	logger.Log.Println("In side UpdateResolutionValueinStagetbl")
	zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
	dataAccess := dao.DbConn{DB: db}
	if recordtypeSeq == 1 || recordtypeSeq == 3 {
		//dataAccess := dao.DbConn{DB: db}
		slaid, err := dataAccess.GetSLAdataexist(ClientID, Mstorgnhirarchyid, RecordID)
		if err != nil {
			logger.Log.Println(err)
			return err
		}
		if slaid > 0 {
			value, err := dataAccess.GetFirstResolutionValue(ClientID, Mstorgnhirarchyid, RecordID)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}

			if value == "NA" {
				createdt, err := dataAccess.GetRecordcreatedate(ClientID, Mstorgnhirarchyid, RecordID)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					return err
				}
				currentTime := time.Now().UTC()
				today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
				strtDay := AddSubSecondsToDate(TimeParse(createdt, ""), zonediff.UTCdiff)
				resolutiontiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
				//log.Println("responsetiimetaken  value is ----->", resolutiontiimetaken)
				logger.Log.Println("Resolution Start date  value is ----->", strtDay)
				logger.Log.Println("Resolution Todate   value is ----->", today)
				logger.Log.Println("Resolution Time taken  value is ----->", resolutiontiimetaken)
				//code change by kaustubh
			idleTime, idleTimeErr := dataAccess.GetIdleTime(ClientID, Mstorgnhirarchyid, RecordID)
			if idleTimeErr != nil {
				logger.Log.Println("idleTimeErr is ----->", idleTimeErr)
				return idleTimeErr
			}
			resolutionTimeTakenWithoutIdleTime := resolutiontiimetaken - idleTime
			err = dataAccess.UpdateFirstResolution(ClientID, Mstorgnhirarchyid, RecordID, resolutiontiimetaken, resolutionTimeTakenWithoutIdleTime)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}

			err = dataAccess.Updateresotimeexcludeidletime(ClientID, Mstorgnhirarchyid, RecordID, idleTime)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					return err
				}

			} else {
				if recorddiffseq == 3 {
					var seqno int64
					if recordtypeSeq == 1 {
						seqno = 10
					} else {
						seqno = 1
					}

					reopendt, err := dataAccess.GetReopendate(ClientID, Mstorgnhirarchyid, RecordID, seqno)
					if err != nil {
						logger.Log.Println("error is ----->", err)
						return err
					}
					if len(reopendt) > 0 {
						currentTime := time.Now().UTC()
						//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
						today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
						strtDay := AddSubSecondsToDate(TimeParse(reopendt, ""), zonediff.UTCdiff)
						resolutiontiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
						//log.Println("responsetiimetaken  value is ----->", resolutiontiimetaken)
						
						logger.Log.Println("Resolution Start date  value is ----->", strtDay)
				logger.Log.Println("Resolution Todate   value is ----->", today)
				logger.Log.Println("Resolution Time taken  value is ----->", resolutiontiimetaken)					
						//code change by Kaustubh
										//idleTime, idleTimeErr := dataAccess.GetIdleTime(ClientID, Mstorgnhirarchyid, RecordID)

						idleTime, idleTimeErr := dataAccess.GetIdleTime(ClientID, Mstorgnhirarchyid, RecordID)
						if idleTimeErr != nil {
							logger.Log.Println("idleTimeErr is ----->", idleTimeErr)
							return idleTimeErr
						}
						resolutionTimeTakenWithoutIdleTime := resolutiontiimetaken - idleTime
						err = dataAccess.UpdateLatestResolution(ClientID, Mstorgnhirarchyid, RecordID, resolutiontiimetaken, resolutionTimeTakenWithoutIdleTime)
						if err != nil {
							logger.Log.Println("error is ----->", err)
							return err
						}

						err = dataAccess.Updateresotimeexcludeidletime(ClientID, Mstorgnhirarchyid, RecordID, idleTime)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					return err
				}
					}
				}
			}
			go ExecuteStagingDetailsUpdate(RecordID)
		}
	}else if recordtypeSeq == 5 {
		err := dataAccess.UpdateFirstResolution(ClientID, Mstorgnhirarchyid, RecordID, 0,0)
		if err != nil {
			logger.Log.Println("error is ----->", err)
			return err
          }
	}


	return nil
}

func UpdateUserreplytimetakenValueinStagetbl(tx1 *sql.Tx, db *sql.DB, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, recorddiffseq int64, Usergroupid int64, Supportgroupspecific int64, recordtypeSeq int64) error {
	logger.Log.Println("In side UpdateUserreplytimetakenValueinStagetbl")
	zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
	if recordtypeSeq == 1 || recordtypeSeq == 3 {
		dataAccess := dao.DbConn{DB: db}
		previousdt, err := dataAccess.GetPreviousstatusdate(ClientID, Mstorgnhirarchyid, RecordID)
		logger.Log.Println("previousdt  -----------previousdt previousdt previousdt previousdt previousdt previousdt previousdt previousdt-- ----->", previousdt)
		if err != nil {
			logger.Log.Println("error is ----->", err)
			return err
		}
		if len(previousdt) > 0 {
			currentTime := time.Now().UTC()
			//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
			today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
			strtDay := AddSubSecondsToDate(TimeParse(previousdt, ""), zonediff.UTCdiff)
			replytimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
			err = dataAccess.UpdateUserreplydatetime(ClientID, Mstorgnhirarchyid, RecordID, replytimetaken)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
		}

	}
	return nil
}

func UpdateFollowuptimetakenValueinStagetbl(tx1 *sql.Tx, db *sql.DB, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, Usergroupid int64, Supportgroupspecific int64, recordtypeSeq int64) error {
	logger.Log.Println("In side UpdateFollowuptimetakenValueinStagetbl")
	zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
	if recordtypeSeq == 1 || recordtypeSeq == 3 {
		dataAccess := dao.DbConn{DB: db}
		previousdt, err := dataAccess.GetPreviousstatusdate(ClientID, Mstorgnhirarchyid, RecordID)
		if err != nil {
			log.Println("error is ----->", err)
			return err
		}
		if len(previousdt) > 0 {
			currentTime := time.Now().UTC()
			//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
			today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
			strtDay := AddSubSecondsToDate(TimeParse(previousdt, ""), zonediff.UTCdiff)
			followuptimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
			log.Println("followuptimetaken  value is ----->", followuptimetaken)

			err = dataAccess.UpdateFollowuptimetaken(ClientID, Mstorgnhirarchyid, RecordID, followuptimetaken)
			if err != nil {
				log.Println("error is ----->", err)
				return err
			}
		}
	}
	return nil
}

func UpdateCloseValueinStagetbl(tx1 *sql.Tx, db *sql.DB, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, recorddiffseq int64, Usergroupid int64, Supportgroupspecific int64, recordtypeSeq int64,previousseq int64) error {

	// For close status update response resolution meter for converted ticket
	logger.Log.Println("In side UpdateCloseValueinStagetbl ---------------->", recorddiffseq, recordtypeSeq)
	zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
	dataAccess := dao.DbConn{DB: db}
      if previousseq != 3 {
	if recordtypeSeq == 1 || recordtypeSeq == 3 {

		issladue, err := dataAccess.FetchSLADueRow(ClientID, Mstorgnhirarchyid, RecordID)
		if err != nil {
			logger.Log.Println("error is ----->", err)

			return err
		}

		isresponsecomplete, isresolutioncomplete, err := dataAccess.FetchResponseResolutionCompleteValue(ClientID, Mstorgnhirarchyid, RecordID)
		if err != nil {
			logger.Log.Println("error is ----->", err)

			return err
		}
		logger.Log.Println("In side UpdateCloseValueinStagetbl -------isresolutioncomplete  --------->", isresponsecomplete, isresolutioncomplete)
		value, err := dataAccess.GetFirstResponseValue(ClientID, Mstorgnhirarchyid, RecordID)
		if err != nil {
			logger.Log.Println("error is ----->", err)

			return err
		}
		if value == "NA" {
			createdt, err := dataAccess.GetRecordcreatedate(ClientID, Mstorgnhirarchyid, RecordID)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
			currentTime := time.Now().UTC()
			//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
			today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
			strtDay := AddSubSecondsToDate(TimeParse(createdt, ""), zonediff.UTCdiff)
			responsetiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
			log.Println("responsetiimetaken  value is ----->", responsetiimetaken)

			err = dataAccess.UpdateFirstResponse(ClientID, Mstorgnhirarchyid, RecordID, responsetiimetaken)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
		} else {
			if recorddiffseq == 1 || recorddiffseq == 10 {
				reopendt, err := dataAccess.GetReopendate(ClientID, Mstorgnhirarchyid, RecordID, recorddiffseq)
				if err != nil {
					logger.Log.Println("error is ----->", err)

					return err
				}
				if len(reopendt) > 0 {
					currentTime := time.Now().UTC()
					//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
					today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
					strtDay := AddSubSecondsToDate(TimeParse(reopendt, ""), zonediff.UTCdiff)
					responsetiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
					log.Println("responsetiimetaken  value is ----->", responsetiimetaken)
					
					err = dataAccess.UpdateLatestResponse(ClientID, Mstorgnhirarchyid, RecordID, responsetiimetaken)
					if err != nil {
						log.Println("error is ----->", err)

						return err
					}
				}
			}

		}
		logger.Log.Println("In side UpdateCloseValueinStagetbl -------issladue  --------->", issladue)
		if isresponsecomplete == 0 && issladue > 0 {
			currentTime := time.Now()
			flag, err := dataAccess.UpdateResponseEndFlag(ClientID, Mstorgnhirarchyid, RecordID, currentTime.Format("2006-01-02 15:04:05"))
			logger.Log.Println(flag)
			if err != nil {
				logger.Log.Println(err)
				return err
			}

		}

		if isresolutioncomplete == 0 && issladue > 0 {
			_, err, _ := UpdateRessolutionEndFlag(ClientID, Mstorgnhirarchyid, RecordID)
			if err != nil {
				logger.Log.Println(err)
				return err
			}

			currentTime := time.Now()
			zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
			datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
			histrn := entities.TrnslaentityhistoryEntity{}
			histrn.Clientid = ClientID
			histrn.Mstorgnhirarchyid = Mstorgnhirarchyid
			histrn.Therecordid = RecordID
			histrn.Recorddatetime = datetime.Format("2006-01-02 15:04:05")
			histrn.Recorddatetoint = TimeParse(datetime.Format("2006-01-02 15:04:05"), "").Unix()
			histrn.Slastartstopindicator = 4

			trnid, err := dataAccess.InsertTrnslaentityhistory(&histrn)
			if err != nil {
				logger.Log.Println(err)
				return err
			}
			logger.Log.Println("history table id---->", trnid)

		}

		value, err = dataAccess.GetFirstResolutionValue(ClientID, Mstorgnhirarchyid, RecordID)
		if err != nil {
			logger.Log.Println("error is ----->", err)

			return err
		}

		if value == "NA" {
			createdt, err := dataAccess.GetRecordcreatedate(ClientID, Mstorgnhirarchyid, RecordID)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
			currentTime := time.Now().UTC()
			//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
			today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
			strtDay := AddSubSecondsToDate(TimeParse(createdt, ""), zonediff.UTCdiff)
			resolutiontiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
			log.Println("responsetiimetaken  value is ----->", resolutiontiimetaken)
			//code change by kaustubh
			idleTime, idleTimeErr := dataAccess.GetIdleTime(ClientID, Mstorgnhirarchyid, RecordID)
			if idleTimeErr != nil {
				logger.Log.Println("idleTimeErr is ----->", idleTimeErr)
				return idleTimeErr
			}
			resolutionTimeTakenWithoutIdleTime := resolutiontiimetaken - idleTime
			err = dataAccess.UpdateFirstResolution(ClientID, Mstorgnhirarchyid, RecordID, resolutiontiimetaken, resolutionTimeTakenWithoutIdleTime)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
			/*err = dataAccess.UpdateFirstResolution(ClientID, Mstorgnhirarchyid, RecordID, resolutiontiimetaken)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}*/
		} else {
			reopendt, err := dataAccess.GetReopendate(ClientID, Mstorgnhirarchyid, RecordID, recorddiffseq)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return err
			}
			if len(reopendt) > 0 {
				currentTime := time.Now().UTC()
				//today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
				today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
				strtDay := AddSubSecondsToDate(TimeParse(reopendt, ""), zonediff.UTCdiff)
				resolutiontiimetaken := CalculateWorkingHourBetweenTwoDates(ClientID, Mstorgnhirarchyid, strtDay, today, 0, Supportgroupspecific, Usergroupid)
				log.Println("responsetiimetaken  value is ----->", resolutiontiimetaken)
				//code change by Kaustubh
								idleTime, idleTimeErr := dataAccess.GetIdleTime(ClientID, Mstorgnhirarchyid, RecordID)

				if idleTimeErr != nil {
					logger.Log.Println("idleTimeErr is ----->", idleTimeErr)
					return idleTimeErr
				}
				resolutionTimeTakenWithoutIdleTime := resolutiontiimetaken - idleTime
				err = dataAccess.UpdateLatestResolution(ClientID, Mstorgnhirarchyid, RecordID, resolutiontiimetaken, resolutionTimeTakenWithoutIdleTime)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					return err
				}
				/*err = dataAccess.UpdateLatestResolution(ClientID, Mstorgnhirarchyid, RecordID, resolutiontiimetaken)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					return err
				}*/
			}
		}
	}
   }
	// For close status update response resolution meter for converted ticket

	err := dataAccess.UpdateCloseddate(ClientID, Mstorgnhirarchyid, RecordID)
	if err != nil {
		logger.Log.Println("error is ----->", err)

		return err
	}
	return nil
}

func UpdateStageResolverInfo(tx1 *sql.Tx, db *sql.DB, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, UserID int64, Usergroupid int64, Supportgroupspecific int64, recorddiffseq int64, recordtypeSeq int64) error {
	logger.Log.Println("In side UpdateStageResolverInfo")
	dataAccess := dao.DbConn{DB: db}

	if recordtypeSeq == 1 || recordtypeSeq == 3 {
		currentTime := time.Now()
		zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
		datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
		logger.Log.Println("In side UpdateStageResolverInfo==================================================================>", ClientID, Mstorgnhirarchyid, datetime, Supportgroupspecific, Usergroupid)

		if Supportgroupspecific == 1 {
			count, err := dataAccess.GetSupportgrpdayofweekcount(ClientID, Mstorgnhirarchyid, Usergroupid)
			if err != nil {
				return err
			}
			if count == 7 {
				time, _ := GetSLAEndTimeForClient(ClientID, Mstorgnhirarchyid, datetime, 259200, Supportgroupspecific, Usergroupid)

				count, err := dataAccess.FetchAutoCloseRecordCount(RecordID)
				if err != nil {
					logger.Log.Println(err)
					return err
				}
				if count > 0 {
					err = dataAccess.DeleteFromClosureTable(RecordID)
					if err != nil {
						logger.Log.Println(err)
						return err
					}
				}

				err = dataAccess.InsertRecordClosure(ClientID, Mstorgnhirarchyid, RecordID, recorddiffseq, time.Format("2006-01-02 15:04:05"))
				if err != nil {
					logger.Log.Println(err)
					return err
				}
			}
		} else {
			count, err := dataAccess.GetOrganizationdayofweekcount(ClientID, Mstorgnhirarchyid)
			if err != nil {
				return err
			}
			if count == 7 {
				time, _ := GetSLAEndTimeForClient(ClientID, Mstorgnhirarchyid, datetime, 259200, Supportgroupspecific, Usergroupid)

				count, err := dataAccess.FetchAutoCloseRecordCount(RecordID)
				if err != nil {
					logger.Log.Println(err)
					return err
				}
				if count > 0 {
					err = dataAccess.DeleteFromClosureTable(RecordID)
					if err != nil {
						logger.Log.Println(err)
						return err
					}
				}

				err = dataAccess.InsertRecordClosure(ClientID, Mstorgnhirarchyid, RecordID, recorddiffseq, time.Format("2006-01-02 15:04:05"))
				if err != nil {
					logger.Log.Println(err)
					return err
				}
			}
		}

	}
	originalInfo, err := dataAccess.GetOriginalInfo(UserID)
	if err != nil {
		logger.Log.Println(err)
		return err
	}

	grpname, err := dataAccess.GetGrpname(Usergroupid)

	if err != nil {
		logger.Log.Println(err)
		return err
	}

	err = dataAccess.UpdateStageResolver(ClientID, Mstorgnhirarchyid, RecordID, UserID, originalInfo.Orgcreatorname, Usergroupid, grpname)
	if err != nil {
		logger.Log.Println("error is ----->", err)
		return err
	}
	return nil
}



func UpdateStageResolverInfoForSR(db *sql.DB, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, UserID int64, Usergroupid int64) error {
	logger.Log.Println("In side UpdateStageResolverInfoForSR")
	dataAccess := dao.DbConn{DB: db}
	originalInfo, err := dataAccess.GetOriginalInfo(UserID)
	if err != nil {
		logger.Log.Println(err)
		return err
	}

	grpname, err := dataAccess.GetGrpname(Usergroupid)

	if err != nil {
		logger.Log.Println(err)
		return err
	}

	err = dataAccess.UpdateStageResolverForSR(ClientID, Mstorgnhirarchyid, RecordID, UserID, originalInfo.Orgcreatorname, Usergroupid, grpname)
	if err != nil {
		logger.Log.Println("error is ----->", err)
		return err
	}

	count, err := dataAccess.FetchAutoCloseRecordCount(RecordID)
	if err != nil {
		logger.Log.Println(err)
		return err
	}
	if count > 0 {
		err = dataAccess.DeleteFromClosureTable(RecordID)
		if err != nil {
			logger.Log.Println(err)
			return err
		}
	}

	err = dataAccess.InsertRecordClosureForSR(ClientID, Mstorgnhirarchyid, RecordID, 3)
	if err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}
