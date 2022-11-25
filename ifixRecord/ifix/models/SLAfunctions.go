package models

import (
	"fmt"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"math"
	"time"
)

func GetStringDate(somedate time.Time) (string, int64) {
	strDate := somedate.Format("2006-01-02")
	newDate := TimeParse(strDate+" 00:00:00", "")
	currentTime := somedate.Unix() - newDate.Unix()
	return strDate, currentTime

}

func GetSLATimeForClient(clientid int64, mstorgnhirarchyid int64, somedate time.Time, supportgroupspecific int64, supportgroupid int64) (time.Time, string) {
	p := logger.Log.Println
	p("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	strDate, currentTime := GetStringDate(somedate)
	p(strDate)
	p(currentTime)
	starttime := int64(0)
	endtime := int64(0)
	isHolliday := bool(false)
	dayweemstarttime := int64(0)
	dayweemendtime := int64(0)

	if supportgroupspecific == 1 {
		starttime, endtime, _, isHolliday, _, _ = GetSupportGroupHoliday(clientid, mstorgnhirarchyid, supportgroupid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _, _ = GetSupportGroupDayOfWeek(clientid, mstorgnhirarchyid, supportgroupid, strDate)
	} else {
		starttime, endtime, _, isHolliday, _, _ = GetClientHoliday(clientid, mstorgnhirarchyid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _ = GetClientDayOfWeek(clientid, mstorgnhirarchyid, strDate)
	}
	// p("dateofholiday >> ", reqdate)
	// p("isHolliday >> ", isweekday)
	if isHolliday == true {
		if currentTime < dayweemstarttime && dayweemstarttime != starttime {
			return AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), dayweemstarttime), "1st point"
		} else if (currentTime >= dayweemstarttime && currentTime <= dayweemendtime) && (dayweemstarttime <= starttime) && (currentTime <= starttime) {
			return AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), currentTime), "2nd point"
		} else if (currentTime >= dayweemstarttime && currentTime <= dayweemendtime) && (currentTime >= starttime && currentTime <= endtime) && (endtime <= dayweemendtime) {
			return AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), endtime+1), "3rd point"
		} else if (currentTime >= dayweemstarttime && currentTime <= dayweemendtime) && (currentTime >= endtime) {
			return AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), currentTime), "4th point"
		} else {
			nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
			resultDate, msg := GetSLATimeForClient(clientid, mstorgnhirarchyid, nextDay, supportgroupspecific, supportgroupid)
			return resultDate, msg
		}
	} else {
		if currentTime < dayweemstarttime {
			return AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), dayweemstarttime), "1st point no holiday"
		} else if currentTime >= dayweemstarttime && currentTime <= dayweemendtime {
			return AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), currentTime), "2st point no holiday"
		} else {
			nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
			resultDate, msg := GetSLATimeForClient(clientid, mstorgnhirarchyid, nextDay, supportgroupspecific, supportgroupid)
			return resultDate, msg
		}
	}
}

// func SLADueTimeCalculation(therecordid int64, resolutiondone int64, fromclientuserid int64, slastateid int64, rcorddatetime string, clientid int64, mstorgnhirarchyid int64, recordtypeid int64, workingdiffid int64, recordpriorityid int64) {
// 	var mstslaentityid = int64(1)
// 	var finaltherecordid = therecordid
// 	// if therecordid == 0 {
// 	// finaltherecordid = therecordid
// 	// }
// 	finalresdone := int64(0)
// 	if resolutiondone != 0 {
// 		finalresdone = resolutiondone
// 	}

// 	//select working category,priority,record type for the record

// 	logger.Log.Println("clientid >>>>>", clientid)
// 	logger.Log.Println("mstorgnhirarchyid >>>>>>", mstorgnhirarchyid)
// 	logger.Log.Println("recordtypeid >>>>>>>>", recordtypeid)
// 	logger.Log.Println("workingdiffid >>>>>>>>", workingdiffid)
// 	logger.Log.Println("recordpriorityid >>>>>>>", recordpriorityid)

// 	returnValue, _, _, _ := SLACriteriaRespResl(clientid, mstorgnhirarchyid, recordtypeid, workingdiffid, recordpriorityid)
// 	logger.Log.Println(returnValue.Clientid)
// 	logger.Log.Println(returnValue.Supportgroupspecific)
// 	logger.Log.Println(returnValue.Slaid)
// 	logger.Log.Println(returnValue.Id)
// 	logger.Log.Println(returnValue.Responsetimeinsec)
// 	logger.Log.Println(returnValue.Resolutiontimeinsec)
// 	p := logger.Log.Println
// 	//p("????????????? ", GetWeekDay("2021-03-13"))
// 	// Specific to Client or Organisation Unit
// 	currentTime := time.Now()
// 	var dt = currentTime.Format("2006-01-02 15:04:05")
// 	logger.Log.Println("dt value is -------------->", dt)
// 	//today := AddSubSecondsToDate(TimeParse(dt, ""), 2600)
// 	today := TimeParse(rcorddatetime, "")
// 	resolutiontimeinsec := returnValue.Resolutiontimeinsec

// 	SLAResponseStartTime := time.Now()
// 	SLAResponseEndTime := time.Now()
// 	SLAResponseEndTimeInt := int64(0)
// 	msg := ""

// 	// check history table
// 	trnslarecords, _, _, _ := GetTrnslaentityhistory(returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid)
// 	if trnslarecords.Slastartstopindicator == 3 {
// 		today = TimeParse(trnslarecords.Recorddatetime, "")
// 		slarecords, _, _, _ := GetMstsladue(returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid)
// 		logger.Log.Println("In time second >>>", slarecords.Remainingtime)
// 		logger.Log.Println("in percentage >>>>> ", slarecords.Completepercent)
// 		resolutiontimeinsec = slarecords.Remainingtime

// 		// if resume store previous resolution data
// 		SLAResponseStartTime = TimeParse(slarecords.Startdatetimeresponse, "")
// 		SLAResponseEndTime = TimeParse(slarecords.Duedatetimeresponse, "")
// 		SLAResponseEndTimeInt = slarecords.Duedatetimeresponseint
// 	} else {

// 		SLAResponseStartTime, msg = GetSLATimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, today)
// 		p("SLAResponseStartTime @@@@@@@@@@2 ", SLAResponseStartTime)
// 		p("SLAResponseStartTime @@@@@@@@@@2 ", msg)

// 		//End Time Response
// 		responseTimeDate := AddSubSecondsToDate(SLAResponseStartTime, returnValue.Responsetimeinsec)
// 		p("responseTimeDate @@@@@@@@@@2 ", responseTimeDate)

// 		SLAResponseEndTime, _ = GetSLAEndTimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, SLAResponseStartTime, returnValue.Responsetimeinsec)
// 		p("SLAResponseEndTime @@@@@@@@@@2 ", SLAResponseEndTime)
// 		p("SLAResponseEndTime @@@@@@@@@@2 ")
// 		// calculate date time to int
// 		SLAResponseEndTimeInt = SLAResponseEndTime.Unix()
// 	}

// 	SLAResolutioneStartTime, msg1 := GetSLATimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, today)
// 	p("SLAResolutioneStartTime @@@@@@@@@@2 ", SLAResolutioneStartTime)
// 	p("SLAResolutioneStartTime @@@@@@@@@@2 ", msg1)

// 	// End Time
// 	resolutionTimeDate := AddSubSecondsToDate(SLAResolutioneStartTime, resolutiontimeinsec)
// 	p("resolutionTimeDate ^^^^^^^^^^^^^^^^^^^^^^ ", resolutionTimeDate)
// 	SLAResolutioneEndTime, msg3 := GetSLAEndTimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, SLAResolutioneStartTime, resolutiontimeinsec)
// 	p("SLAResolutioneEndTime @@@@@@@@@@2 ", SLAResolutioneEndTime)
// 	p("SLAResolutioneEndTime @@@@@@@@@@2 ", msg3)

// 	SLAResolutioneEndTimeInt := SLAResolutioneEndTime.Unix()
// 	p("SLAResolutioneEndTimeInt >>>>>>>>>>>>>>>>>>>>> ", SLAResolutioneEndTimeInt)

// 	p("SLAResponseEndTimeInt >>>>>>>>>>>>>>>>>>>>> ", SLAResponseEndTimeInt)
// 	InsertSLARecord(returnValue.Clientid, returnValue.Mstorgnhirarchyid, mstslaentityid, finaltherecordid, rcorddatetime, slastateid, fromclientuserid,
// 		SLAResponseStartTime, SLAResolutioneStartTime, SLAResponseEndTime, SLAResolutioneEndTime, SLAResolutioneEndTimeInt, SLAResponseEndTimeInt, finalresdone)

// 	//t := entities.SLAResponseEntity{}

// 	// t.SLAResponseStartTime = SLAResponseStartTime
// 	// t.SLAResolutioneStartTime = SLAResolutioneStartTime
// 	// t.ResponseTimeDate = responseTimeDate
// 	// t.SLAResponseEndTime = SLAResponseEndTime
// 	// t.ResolutionTimeDate = resolutionTimeDate
// 	// t.SLAResolutioneEndTime = SLAResolutioneEndTime

// 	//return t, nil
// }

func SLADueTimeCalculation(therecordid int64, resolutiondone int64, fromclientuserid int64, slastateid int64, rcorddatetime string, clientid int64, mstorgnhirarchyid int64, recordtypeid int64, workingdiffid int64, recordpriorityid int64, flag string, supportgroupid int64) {
	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>22222222222222222222222222222222222222222222222222222222222222222222222222>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.", supportgroupid, therecordid)
	var mstslaentityid = int64(1)
	var finaltherecordid = therecordid
	// if therecordid == 0 {
	// finaltherecordid = therecordid
	// }
	finalresdone := int64(0)
	if resolutiondone != 0 {
		finalresdone = resolutiondone
	}

	//select working category,priority,record type for the record
	logger.Log.Println("Flag value is  >>>>>", flag)
	logger.Log.Println("Record date time is  >>>>>", rcorddatetime)

	logger.Log.Println("clientid >>>>>", clientid)
	logger.Log.Println("mstorgnhirarchyid >>>>>>", mstorgnhirarchyid)
	logger.Log.Println("recordtypeid >>>>>>>>", recordtypeid)
	logger.Log.Println("workingdiffid >>>>>>>>", workingdiffid)
	logger.Log.Println("recordpriorityid >>>>>>>", recordpriorityid)

	returnValue, _, _, _ := SLACriteriaRespResl(clientid, mstorgnhirarchyid, recordtypeid, workingdiffid, recordpriorityid)
	if returnValue.Id == 0 {
		logger.Log.Println("****************************** Not data found in SLACriteriaRespResl *********************** ")
		return
	}
	logger.Log.Println("returnValue.Responsetimeinsec ----------->", returnValue.Responsetimeinsec)
	logger.Log.Println("returnValue.Resolutiontimeinsec------------->", returnValue.Resolutiontimeinsec)
	logger.Log.Println("returnValue.Supportgroupspecific------------->", returnValue.Supportgroupspecific)
	p := logger.Log.Println
	//p("????????????? ", GetWeekDay("2021-03-13"))
	// Specific to Client or Organisation Unit
	currentTime := time.Now()
	var dt = currentTime.Format("2006-01-02 15:04:05")
	logger.Log.Println("dt value is -------------->", dt)
	//today := AddSubSecondsToDate(TimeParse(dt, ""), 2600)
	today := TimeParse(rcorddatetime, "")
	if today.IsZero() {
		logger.Log.Println("****************************** Incorrect date capture *********************** ", today)
		return
	}
	p(today)
	resolutiontimeinsec := returnValue.Resolutiontimeinsec
	logger.Log.Println("resolutiontimeinsec  resolutiontimeinsec  resolutiontimeinsec  resolutiontimeinsec    >>>", resolutiontimeinsec)

	SLAResponseStartTime := time.Now()
	//SLAResolutioneStartTime := time.Now()
	SLAResolutioneStartTime := TimeParse(rcorddatetime, "")
	if SLAResolutioneStartTime.IsZero() {
		logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResolutioneStartTime)
		return
	}
	SLAResponseEndTime := time.Now()
	SLAResponseEndTimeInt := int64(0)
	msg := ""
	msg2 := ""
	p(msg, msg2)
	// check history table
	trnslarecords, _, _, _ := GetTrnslaentityhistory(returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid)
	/*if trnslarecords.Id == 0 {
		logger.Log.Println("****************************** Not data found in GetTrnslaentityhistory *********************** ")
		return
	}*/
	logger.Log.Println("trnslarecords.Slastartstopindicator  111111111111111111111111111  >>>", trnslarecords.Slastartstopindicator)
	var totalpushTime int64
	slarecords, _, _, _ := GetMstsladue(returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid)
	if slarecords.Id == 0 {
		logger.Log.Println("****************************** data not found in mstsladue *********************** ")
		//return
	}
	logger.Log.Println("In time second >>>", slarecords.Remainingtime)
	logger.Log.Println("in percentage >>>>> ", slarecords.Completepercent)
	totalpushTime = slarecords.PushTime

	if trnslarecords.Slastartstopindicator == 3 {
		logger.Log.Println("trnslarecords.Slastartstopindicator  333333333333333333333333333  == 3    >>>")

		//resolutiontimeinsec = slarecords.Remainingtime
		today = TimeParse(trnslarecords.Recorddatetime, "")
		if today.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", today)
			return
		}

		if flag == "SC" {
			resolutiontimeinsec = returnValue.Resolutiontimeinsec
			SLAResolutioneStartTime = TimeParse(slarecords.Startdatetimeresolution, "")

		} else if flag != "P" {
			logger.Log.Println("flag != P >>>>>>>>>>>>>>>>>>>>>>>>...", flag)
			resolutiontimeinsec = slarecords.Remainingtime
			SLAResolutioneStartTime = TimeParse(trnslarecords.Recorddatetime, "")
		} else {
			logger.Log.Println("Not equals P value is  ----------------------------------------------    >>>")
			SLAResolutioneStartTime = TimeParse(rcorddatetime, "")
		}
		if SLAResolutioneStartTime.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResolutioneStartTime)
			return
		}

		// if resume store previous resolution data
		SLAResponseStartTime = TimeParse(slarecords.Startdatetimeresponse, "")
		SLAResponseEndTime = TimeParse(slarecords.Duedatetimeresponse, "")
		SLAResponseEndTimeInt = slarecords.Duedatetimeresponseint
		if SLAResponseStartTime.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResponseStartTime)
			return
		}
		if SLAResponseEndTime.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResponseEndTime)
			return
		}

		//SLAResolutioneStartTime = TimeParse(trnslarecords.Recorddatetime, "")
		p(">>>>>>>>>>>>>>>>>>>>>> PushTime >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		//pushRecord, _, _, _ := GetTrnslaentityhistorytype2(returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid)
//		pushRecord, _, _, _ := GetTrnslaentityhistorytype2(returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid, trnslarecords.Id)
		pushRecord, isPushTimeCalculated, _, _, _ := CheckTrnslaentityhistorytype(returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid, trnslarecords.Id)//changes done by josim as per santanu da

		logger.Log.Println("pushRecord value is  ----------------------------------------------    >>>", pushRecord,isPushTimeCalculated)
		if pushRecord.Id != 0 && !isPushTimeCalculated {
			pushDateTime := TimeParse(pushRecord.Recorddatetime, "")
			if pushDateTime.IsZero() {
				logger.Log.Println("****************************** Incorrect date capture *********************** ", pushDateTime)
				return
			}
			//changes supportgrp here  .. .......Santanu,simanta
//=========================================To get last grpid=====================================================
			if returnValue.Supportgroupspecific == 1 {
				requestid, _, err, _ := Getrequestid(therecordid)
				if err != nil {
					logger.Log.Println("****************************** Not Getting Requestid *********************** ", requestid)
				}
				lastgroupid, _, err, _ := GetLastSupportgrp(requestid, supportgroupid)
				if lastgroupid > 0 {
					supportgroupid = lastgroupid
				}
				logger.Log.Println("****************************** currentgrpid,lastgrpid *********************** ", supportgroupid, lastgroupid)
			}
//=========================================To get last grpid=====================================================


			pushTime := CalculateWorkingHourBetweenTwoDates(returnValue.Clientid, returnValue.Mstorgnhirarchyid, pushDateTime, today, int64(0), returnValue.Supportgroupspecific, supportgroupid)
			upRes, _, _ := UpdatePushTimeInHistory(trnslarecords.Id, pushTime)
			p("Update push time result >>>>>>>>> ", upRes)
			totalpushTime = pushTime + slarecords.PushTime

		}
		if flag == "P" {
			p(">>>>>>>>>>>>>>>>>>>>  22222222222222222222222222222222222 >>>>>>>>>>>>>>>>>>> ", totalpushTime)
			resolutiontimeinsec = resolutiontimeinsec + totalpushTime
		} else if flag == "SC" && totalpushTime > 0 {
			resolutiontimeinsec = resolutiontimeinsec + totalpushTime
		}
		//p(">>>>>>>>>>>>>>>>>>>>  total push time >>>>>>>>>>>>>>>>>>> ", totalpushTime)
	} else {
		if trnslarecords.Slastartstopindicator == 5 {
			today = TimeParse(trnslarecords.Recorddatetime, "")
			if today.IsZero() {
				logger.Log.Println("****************************** Incorrect date capture *********************** ", today)
				return
			}
		}
		SLAResponseStartTime, msg = GetSLATimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, today, returnValue.Supportgroupspecific, supportgroupid)
		if SLAResponseStartTime.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResponseStartTime)
			return
		}
		p("SLAResponseStartTime @@@@@@@@@@2 ", msg)

		//End Time Response
		responseTimeDate := AddSubSecondsToDate(SLAResponseStartTime, returnValue.Responsetimeinsec)
		p("responseTimeDate @@@@@@@@@@2 ", responseTimeDate)

		//SLAResponseEndTime, msg2 = GetSLAEndTimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, SLAResponseStartTime, returnValue.Responsetimeinsec)
		var responseTimeInSec = returnValue.Responsetimeinsec
		if flag == "P" {
			responseTimeInSec = returnValue.Responsetimeinsec + slarecords.PushTime
		}

		SLAResponseEndTime, msg2 = GetSLAEndTimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, SLAResponseStartTime, responseTimeInSec, returnValue.Supportgroupspecific, supportgroupid)
		if SLAResponseEndTime.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResponseEndTime)
			return
		}
		p("SLAResponseEndTime @@@@@@@@@@2 ", msg2)
		// calculate date time to int
		SLAResponseEndTimeInt = SLAResponseEndTime.Unix()

		SLAResolutioneStartTime, _ = GetSLATimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, today, returnValue.Supportgroupspecific, supportgroupid)
		if SLAResolutioneStartTime.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResolutioneStartTime)
			return
		}
		p("SLAResolutioneStartTime @@@@@@@@@@2 ", SLAResolutioneStartTime)
		p("slarecords.PushTime @@@@@@@@@@2 ", slarecords.PushTime)
		if flag == "P" {
			p("slarecords.Remainingtime   if condition                   @@@@@@@@@@2 ", slarecords.Remainingtime)
			//resolutiontimeinsec = slarecords.Remainingtime
			p("resolutiontimeinsec        if condition        @@@@@@@@@@2 ", resolutiontimeinsec)
			resolutiontimeinsec = resolutiontimeinsec + slarecords.PushTime
		}
	}

	p("**************flag *********************> ", flag)

	p("resolutiontimeinsec @@@@@@@@@@2 ", resolutiontimeinsec)
	SLAResolutioneEndTime, msg3 := GetSLAEndTimeForClient(returnValue.Clientid, returnValue.Mstorgnhirarchyid, SLAResolutioneStartTime, resolutiontimeinsec, returnValue.Supportgroupspecific, supportgroupid)
	if SLAResolutioneEndTime.IsZero() {
		logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResolutioneEndTime)
		return
	}
	p("SLAResolutioneEndTime @@@@@@@@@@2 ", SLAResolutioneEndTime)
	p("SLAResolutioneEndTime @@@@@@@@@@2 ", msg3)

	SLAResolutioneEndTimeInt := SLAResolutioneEndTime.Unix()
	//p("SLAResolutioneEndTimeInt >>>>>>>>>>>>>>>>>>>>> ", SLAResolutioneEndTimeInt)

	//p("SLAResponseEndTimeInt >>>>>>>>>>>>>>>>>>>>> ", SLAResponseEndTimeInt)

	slatabentity := entities.SLATabEntity{ClientID: clientid, Mstorgnhirarchyid: mstorgnhirarchyid, RecordID: finaltherecordid, RecordtypeID: recordtypeid, WorkingcatID: workingdiffid, PriorityID: recordpriorityid}

	InsertSLARecord(returnValue.Clientid, returnValue.Mstorgnhirarchyid, mstslaentityid, finaltherecordid, rcorddatetime, slastateid, fromclientuserid,
		SLAResponseStartTime, SLAResolutioneStartTime, SLAResponseEndTime, SLAResolutioneEndTime, SLAResolutioneEndTimeInt, SLAResponseEndTimeInt, finalresdone, trnslarecords.Slastartstopindicator, totalpushTime, flag, slatabentity, returnValue.Supportgroupspecific, supportgroupid)

}

func GetSLAEndTimeForClient(clientid int64, mstorgnhirarchyid int64, somedate time.Time, second int64, supportgroupspecific int64, supportgroupid int64) (time.Time, string) {
	p := logger.Log.Println
	p("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	strDate, currentTime := GetStringDate(somedate)
	p("somedate value is -------------------->", somedate)
	//p(currentTime)
	starttime := int64(0)
	endtime := int64(0)
	isHolliday := bool(false)
	dayweemstarttime := int64(0)
	dayweemendtime := int64(0)

	p("supportgroupspecific   value is ---------111111111111----------->", supportgroupspecific)
	if supportgroupspecific == 1 {
		starttime, endtime, _, isHolliday, _, _ = GetSupportGroupHoliday(clientid, mstorgnhirarchyid, supportgroupid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _, _ = GetSupportGroupDayOfWeek(clientid, mstorgnhirarchyid, supportgroupid, strDate)
	} else {
		starttime, endtime, _, isHolliday, _, _ = GetClientHoliday(clientid, mstorgnhirarchyid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _ = GetClientDayOfWeek(clientid, mstorgnhirarchyid, strDate)
	}
	p("starttime, endtime, _, isHolliday  value is ---------111111111111----------->", starttime, endtime, isHolliday)
	p("dayweemstarttime, dayweemendtime    value is ---------111111111111----------->", dayweemstarttime, dayweemendtime)

	var workingHour = int64(0)
	currentTime1 := time.Now()
	zonediff, _, _, _ := Getutcdiff(clientid, mstorgnhirarchyid)
	datetime := AddSubSecondsToDate(currentTime1, zonediff.UTCdiff)
	today, timeinsec := GetStringDate(datetime)
	p("-----------today ---------------- ", today, timeinsec)

	if currentTime > dayweemstarttime { //strDate == today &&
		dayweemstarttime = currentTime
	}

	workingHour = calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
	p("-----------eorking hour ---------------- ", workingHour, second)
	//p("-----------dayweemstarttime 1111111111111111111111 ---------------- ", dayweemstarttime)
	if workingHour < 0 {
		workingHour = 0
	}
	if workingHour >= second {
		return AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), dayweemstarttime+second), "working hour grater than totalhour"
	} else {
		nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
		second = second - workingHour
		resultDate, msg := GetSLAEndTimeForClient(clientid, mstorgnhirarchyid, nextDay, second, supportgroupspecific, supportgroupid)
		return resultDate, msg
	}
}

// func GetSLARemainingTimeForClient(clientid int64, mstorgnhirarchyid int64, id int64, somedate time.Time, startdatetimeresponse string, duedatetimeresponse string, availableTime int64) (int64, string) {
// 	p := logger.Log.Println
// 	p("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ ", duedatetimeresponse)
// 	strDate, currentTime := GetStringDate(somedate)
// 	p("date start >>>>>>>>>>> ", strDate)
// 	p("date start >>>>>>>>>>> ", currentTime)

// 	starttime, endtime, _, isHolliday, _, _ := GetClientHoliday(clientid, mstorgnhirarchyid, strDate)
// 	p("holiday starttime >> ", starttime)
// 	p("holiday endtime >> ", endtime)
// 	p("holiday isHolliday >> ", isHolliday)
// 	dayweemstarttime, dayweemendtime, _, _, _, _ := GetClientDayOfWeek(clientid, mstorgnhirarchyid, strDate)
// 	p("dayweemstarttime >> ", dayweemstarttime)
// 	p("dayweemendtime >> ", dayweemendtime)

// 	// startRespDate, startRespTime := GetStringDate(TimeParse(startdatetimeresponse, ""))
// 	// p(">>>>>>>>> ", startRespDate)
// 	// p(">>>>>>>>", startRespTime)
// 	endRespDate, endtRespTime := GetStringDate(TimeParse(duedatetimeresponse, ""))
// 	p(">>>>>>>>> ", endRespDate)
// 	p(">>>>>>>>", endtRespTime)

// 	p("****************************** Condition Start **************************")
// 	if strDate == endRespDate {
// 		dayweemendtime = endtRespTime
// 		resultTime := calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
// 		p("result time >>>>>>>>>>>>>", resultTime, strDate)
// 		return availableTime + resultTime, "END"
// 	} else {
// 		var availableTime1 = int64(0)
// 		if dayweemstarttime < currentTime && currentTime <= dayweemendtime {
// 			dayweemstarttime = currentTime
// 			availableTime1 = calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
// 		} else if dayweemendtime < currentTime {
// 			availableTime1 = 0
// 		} else {
// 			availableTime1 = calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
// 		}
// 		p("result time >>>>>>>>>>>>>", availableTime1, strDate)
// 		nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
// 		resultTime, msg := GetSLARemainingTimeForClient(clientid, mstorgnhirarchyid, id, nextDay, startdatetimeresponse, duedatetimeresponse, availableTime1)
// 		return availableTime + resultTime, msg
// 	}

//}

func GetSLARemainingTimeForClient(clientid int64, mstorgnhirarchyid int64, id int64, somedate time.Time, startdatetimeresponse string, duedatetimeresponse string, availableTime int64, supportgroupspecific int64, supportgroupid int64) (int64, string) {
	p := logger.Log.Println
	p("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ supportgroupspecific supportgroupspecific $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ ", supportgroupspecific)
	if somedate.IsZero() {
		p("****************************** Incorrect date capture in GetSLARemainingTimeForClient *********************** ", somedate)
		return 0, ""
	}
	// p("current date time -------------------- ", somedate)
	strDate, currentTime := GetStringDate(somedate)
	// p("date start >>>>>>>>>>> ", strDate)
	// p("date start >>>>>>>>>>> ", currentTime)
	starttime := int64(0)
	endtime := int64(0)
	isHolliday := bool(false)
	dayweemstarttime := int64(0)
	dayweemendtime := int64(0)

	if supportgroupspecific == 1 {
		starttime, endtime, _, isHolliday, _, _ = GetSupportGroupHoliday(clientid, mstorgnhirarchyid, supportgroupid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _, _ = GetSupportGroupDayOfWeek(clientid, mstorgnhirarchyid, supportgroupid, strDate)
	} else {
		starttime, endtime, _, isHolliday, _, _ = GetClientHoliday(clientid, mstorgnhirarchyid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _ = GetClientDayOfWeek(clientid, mstorgnhirarchyid, strDate)
	}

	// startRespDate, startRespTime := GetStringDate(TimeParse(startdatetimeresponse, ""))
	p("dayweemstarttime, dayweemendtime             >>>>>>>>> ", dayweemstarttime, dayweemendtime)
	// p(">>>>>>>>", startRespTime)
	//p("duedatetimeresponse >>>>>>>>> ", duedatetimeresponse)
	endRespDate, endtRespTime := GetStringDate(TimeParse(duedatetimeresponse, ""))
	//p(">>>>>>>>> ", endRespDate)
	//p(">>>>>>>>", endtRespTime)

	p("strDate **************************", strDate)
	p("strDate **************************", endRespDate)
	if strDate == endRespDate {
		if dayweemstarttime < currentTime {
			dayweemstarttime = currentTime
		}
		dayweemendtime = endtRespTime
		resultTime := calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
		//p("result time >>>>>>>>>>>>>", resultTime, strDate)
		return availableTime + resultTime, "END"
	} else {
		var availableTime1 = int64(0)
		if dayweemstarttime < currentTime && currentTime <= dayweemendtime {
			dayweemstarttime = currentTime
			availableTime1 = calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
		} else if dayweemendtime < currentTime {
			availableTime1 = 0
		} else {
			availableTime1 = calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
		}
		p("result time >>>>>>>>>>>>>", availableTime1, strDate)
		nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
		resultTime, msg := GetSLARemainingTimeForClient(clientid, mstorgnhirarchyid, id, nextDay, startdatetimeresponse, duedatetimeresponse, availableTime1, supportgroupspecific, supportgroupid)
		return availableTime + resultTime, msg
	}

}

func InsertSLARecord(Clientid int64, Mstorgnhirarchyid int64, mstslaentityid int64, finaltherecordid int64, rcorddatetime string, slastateid int64, fromclientuserid int64,
	SLAResponseStartTime time.Time, SLAResolutioneStartTime time.Time, SLAResponseEndTime time.Time, SLAResolutioneEndTime time.Time, SLAResolutioneEndTimeInt int64, SLAResponseEndTimeInt int64, finalresdone int64, indicator int64, totalpushTime int64, flag string, slatabentity entities.SLATabEntity, supportgroupspecific int64, supportgroupid int64) {

	// if SLAResponseStartTime.IsZero() {
	// 	logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResponseStartTime)
	// 	return
	// }
	// if SLAResolutioneStartTime.IsZero() {
	// 	logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResolutioneStartTime)
	// 	return
	// }
	// if SLAResponseEndTime.IsZero() {
	// 	logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResponseEndTime)
	// 	return
	// }
	// if SLAResolutioneEndTime.IsZero() {
	// 	logger.Log.Println("****************************** Incorrect date capture *********************** ", SLAResolutioneEndTime)
	// 	return
	// }
	// Insert transaction history data
	var transactionHist = entities.TrnslaentityhistoryEntity{}
	var sladue = entities.MstsladueEntity{}
	slarecords, _, _, _ := GetMstsladue(Clientid, Mstorgnhirarchyid, finaltherecordid)
	logger.Log.Println("sla due record ", slarecords)
	slastartstopindicator := int64(1)

	if slarecords.Id != 0 && indicator != 5 {
		slastartstopindicator = 3
	}

	if flag == "P" {
		slastartstopindicator = indicator
	}

	p := logger.Log.Println
	// TrnslaentityhistoryEntity entity
	transactionHist.Clientid = Clientid
	transactionHist.Mstorgnhirarchyid = Mstorgnhirarchyid
	transactionHist.Mstslaentityid = mstslaentityid
	transactionHist.Therecordid = finaltherecordid
	transactionHist.Recorddatetime = rcorddatetime
	// p("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// p(TimeParse(rcorddatetime, "2006-01-02 15:04:05").Unix())
	transactionHist.Recorddatetoint = TimeParse(rcorddatetime, "").Unix()
	transactionHist.Donotupdatesladue = 1
	transactionHist.Recordtimetoint = 0
	transactionHist.Mstslastateid = slastateid
	transactionHist.Commentonrecord = "Start of SLA"
	transactionHist.Slastartstopindicator = slastartstopindicator
	transactionHist.Fromclientuserid = fromclientuserid
	transactionid, _, _ := InsertTrnslaentityhistory(transactionHist)

	// mstsladie entity
	sladue.Clientid = Clientid
	sladue.Mstorgnhirarchyid = Mstorgnhirarchyid
	sladue.Mstslaentityid = mstslaentityid
	sladue.Therecordid = finaltherecordid
	sladue.Latestone = 1
	sladue.Startdatetimeresponse = DateTimeFormat(SLAResponseStartTime, "")
	sladue.Startdatetimeresolution = DateTimeFormat(SLAResolutioneStartTime, "")
	sladue.Duedatetimeresponse = DateTimeFormat(SLAResponseEndTime, "")
	sladue.Duedatetimeresolution = DateTimeFormat(SLAResolutioneEndTime, "")
	sladue.DuedatetimeresolutionInt = SLAResolutioneEndTimeInt
	sladue.Duedatetimeresponseint = SLAResponseEndTimeInt
	sladue.Remainingtime = 0
	sladue.Completepercent = 0
	sladue.Duedatetimetominute = 0
	sladue.Resoltiondone = finalresdone
	sladue.Trnslaentityhistoryid = transactionid

	if finalresdone != 0 {
		sladue.Resolutiondatetime = DateTimeFormat(time.Now(), "")
	} else {
		sladue.Resolutiondatetime = DateTimeFormat(time.Now(), "")
	}
	sladue.Lastupdatedattime = DateTimeFormat(time.Now(), "")
	if slarecords.Id != 0 && indicator != 5 {
		//p("111111111111111111111111111111111111111--------if part------------------------>")
		uResult, _, _ := UpdateMstsladue(Clientid, Mstorgnhirarchyid, finaltherecordid, DateTimeFormat(SLAResponseEndTime, ""), DateTimeFormat(SLAResolutioneEndTime, ""), SLAResolutioneEndTimeInt, SLAResponseEndTimeInt, totalpushTime)
		p("***** SLA dao update *********** ", uResult)
	} else {
		//p("111111111111111111111111111111111111111--------else part------------------------>")
		sladueid, _, _ := InsertMstsladue(sladue)
		p("***** SLA dao ID *********** ", sladueid)
	}

	// call from SLATabModel.go
	if flag == "P" {
		UpdateSLARemainingTime(slatabentity, indicator, supportgroupspecific, supportgroupid)
	}
}

//-------------------------------------

// func calculateWorkingTime(dayweemstarttime int64, dayweemendtime int64, starttime int64, endtime int64, isHolliday bool) int64 {
// 	var workingHour = int64(0)
// 	if isHolliday == true {
// 		logger.Log.Println(dayweemstarttime, starttime)
// 		logger.Log.Println(dayweemendtime, endtime)
// 		if (dayweemstarttime >= starttime) && (dayweemendtime <= endtime) {
// 			workingHour = 0
// 		} else if (dayweemstarttime >= starttime) && (dayweemendtime >= endtime) {
// 			workingHour = (dayweemstarttime - starttime) + (dayweemendtime - endtime)
// 		} else {
// 			workingHour = (starttime - dayweemstarttime) + (dayweemendtime - endtime)
// 		}
// 	} else {
// 		workingHour = dayweemendtime - dayweemstarttime
// 	}
// 	return workingHour
// }

func calculateWorkingTime(dayweemstarttime int64, dayweemendtime int64, starttime int64, endtime int64, isHolliday bool) int64 {
	logger.Log.Println("start time >>>>>>>>>>>>>>>>>>>>>> ", dayweemstarttime, starttime, endtime)
	var workingHour = int64(0)
	if isHolliday == true {
		logger.Log.Println(dayweemstarttime, starttime)
		logger.Log.Println(dayweemendtime, endtime)
		if (dayweemstarttime >= starttime) && (dayweemendtime <= endtime) {
			workingHour = 0
		} else if (dayweemstarttime >= starttime) && (dayweemendtime >= endtime) {
			// workingHour = (dayweemstarttime - starttime) + (dayweemendtime - endtime)
			workingHour = dayweemendtime - endtime
		} else if (dayweemstarttime <= starttime) && (dayweemendtime <= endtime) {
			workingHour = starttime - dayweemstarttime
		} else {
			workingHour = (starttime - dayweemstarttime) + (dayweemendtime - endtime)
		}
	} else {
		workingHour = dayweemendtime - dayweemstarttime
	}
	return workingHour
}

/* func checkslagroup(Clientid int64, Mstorgnhirarchyid int64, Id int64) {
	currentTime := time.Now().UTC()
	today := TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
	slarecords, _, _, _ := GetMstsladue(Clientid, Mstorgnhirarchyid, Id)
	logger.Log.Println("sla due record ", slarecords)
	startdatetimeresolution := slarecords.Startdatetimeresolution
	duedatetimeresolution := slarecords.Duedatetimeresolution
	availableTime := int64(0)
	// startdatetimeresponse := slarecords.Startdatetimeresponse
	//  duedatetimeresponse := slarecords.Duedatetimeresponse
	result, msg := GetSLARemainingTimeForClient(Clientid, Mstorgnhirarchyid, Id, today, startdatetimeresolution, duedatetimeresolution, availableTime)
	logger.Log.Println(result)
	logger.Log.Println(msg)
} */

func CalculateWorkingHourBetweenTwoDates(clientid int64, mstorgnhirarchyid int64, startDate time.Time, endDate time.Time, second int64, supportgroupspecific int64, supportgroupid int64) int64 {
	strDate, currentTime := GetStringDate(startDate)
	strEndDate, strendTime := GetStringDate(endDate)
	starttime := int64(0)
	endtime := int64(0)
	isHolliday := bool(false)
	dayweemstarttime := int64(0)
	dayweemendtime := int64(0)

	if supportgroupspecific == 1 {
		starttime, endtime, _, isHolliday, _, _ = GetSupportGroupHoliday(clientid, mstorgnhirarchyid, supportgroupid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _, _ = GetSupportGroupDayOfWeek(clientid, mstorgnhirarchyid, supportgroupid, strDate)
	} else {
		starttime, endtime, _, isHolliday, _, _ = GetClientHoliday(clientid, mstorgnhirarchyid, strDate)
		dayweemstarttime, dayweemendtime, _, _, _, _ = GetClientDayOfWeek(clientid, mstorgnhirarchyid, strDate)
	}
	var workingHour = int64(0)
	if currentTime > dayweemstarttime {
		dayweemstarttime = currentTime
	}
	if strEndDate == strDate && strendTime < dayweemendtime {
		dayweemendtime = strendTime
	}
	fmt.Println(">>>>>>>>>>>>>. startDate >>>>>>>>>>>>>>", startDate)
	fmt.Println(">>>>>>>>>>>>>. endDate >>>>>>>>>>>>>>", endDate)
	fmt.Println(">>>>>>>>>>>>>. strDate, currentTime >>>>>>>>>>>>>>", strDate, currentTime)
	fmt.Println(">>>>>>>>>>>>>. strDate, currentTime >>>>>>>>>>>>>>", strEndDate, strendTime)
	fmt.Println(">>>>>>>>>>>>>. dayweemstarttime >>>>>>>>>>>>>>", dayweemstarttime)
	fmt.Println(">>>>>>>>>>>>>. dayweemendtime >>>>>>>>>>>>>>", dayweemendtime)

	workingHour = calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
	if strDate == strEndDate {
		return workingHour + second
	} else {
		nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
		second = second + workingHour
		resultDate := CalculateWorkingHourBetweenTwoDates(clientid, mstorgnhirarchyid, nextDay, endDate, second, supportgroupspecific, supportgroupid)
		return resultDate
	}
}

func calculateRemainingTime(ClientID int64, Mstorgnhirarchyid int64, duedatetimeresolution string, RecordID int64, startdatetimeresolution string, supportgroupspecific int64, supportgroupid int64) int64 {
	currentTime := time.Now().UTC()
	zonediff, _, _, _ := Getutcdiff(ClientID, Mstorgnhirarchyid)
	today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	availableTime := int64(0)
	result := int64(0)
	if today.Unix() > TimeParse(duedatetimeresolution, "").Unix() {
		result, _ = GetSLARemainingTimeForClient(ClientID, Mstorgnhirarchyid, RecordID, TimeParse(duedatetimeresolution, ""), startdatetimeresolution, today.Format("2006-01-02 15:04:05"), availableTime, supportgroupspecific, supportgroupid)
		result = -result
	} else {
		result, _ = GetSLARemainingTimeForClient(ClientID, Mstorgnhirarchyid, RecordID, today, startdatetimeresolution, duedatetimeresolution, availableTime, supportgroupspecific, supportgroupid)
	}
	return result
}

func UpdateSLARemainingTime(page entities.SLATabEntity, indicator int64, supportgroupspecific int64, supportgroupid int64) {
	logger.Log.Println("************************************************ UpdateSLARemainingTime ******************************************************")
	logger.Log.Println("GetSLAResolution page.ClientID no ---->", page.ClientID)
	logger.Log.Println("GetSLAResolution page.Mstorgnhirarchyid no ---->", page.Mstorgnhirarchyid)
	logger.Log.Println("GetSLAResolution page.RecordtypeID no ---->", page.RecordtypeID)
	logger.Log.Println("GetSLAResolution page.WorkingcatID no ---->", page.WorkingcatID)
	logger.Log.Println("GetSLAResolution page.PriorityID no ---->", page.PriorityID)

	currentTime := time.Now().UTC()
	//today := utility.TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
	// Get SLA total response and resolution time

	zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
	today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	if indicator == 2 {
		lastPushHistory, _, _, _ := GetTrnslaentityhistoryLastPushTime(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		today = TimeParse(lastPushHistory.Recorddatetime, "")
		if today.IsZero() {
			logger.Log.Println("****************************** Incorrect date capture *********************** ", today)
			return
		}
	}
	returnValue, _, _, _ := SLACriteriaRespResl(page.ClientID, page.Mstorgnhirarchyid, page.RecordtypeID, page.WorkingcatID, page.PriorityID)

	logger.Log.Println(returnValue.Responsetimeinsec)
	logger.Log.Println(returnValue.Resolutiontimeinsec)

	slarecords, _, _, _ := GetMstsladue(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	logger.Log.Println("sla due record ", slarecords)
	startdatetimeresponse := slarecords.Startdatetimeresponse
	duedatetimeresponse := slarecords.Duedatetimeresponse

	startdatetimeresolution := slarecords.Startdatetimeresolution
	duedatetimeresolution := slarecords.Duedatetimeresolution

	logger.Log.Println("startdatetimeresolution value is  ---->", startdatetimeresolution)
	logger.Log.Println("duedatetimeresponse value is  ---->", duedatetimeresponse)
	// slarecords.Responseremainingtime
	// slarecords.Responsepercentage

	respRemainingTime := slarecords.Responseremainingtime
	respPercent := slarecords.Responsepercentage
	// if respPercent >= 100 {
	// 	respPercent = 100
	// }

	trnslarecords, _, _, _ := GetTrnslaentityhistory(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	logger.Log.Println("In time second >>>", slarecords.Remainingtime)
	if trnslarecords.Slastartstopindicator == 4 {
		uResult, _, _ := UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, slarecords.Remainingtime, slarecords.Completepercent, respRemainingTime, respPercent)
		logger.Log.Println("is update >>>", uResult)
	} else {
		// For Resolution meter
		availableTime := int64(0)
		result := int64(0)
		if today.Unix() > TimeParse(duedatetimeresolution, "").Unix() {
			logger.Log.Println("If part  >>>", today)
			result, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, TimeParse(duedatetimeresolution, ""), startdatetimeresolution, today.Format("2006-01-02 15:04:05"), availableTime, supportgroupspecific, supportgroupid)
			result = -result
		} else {
			logger.Log.Println("Else part  >>>", today)
			result, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, today, startdatetimeresolution, duedatetimeresolution, availableTime, supportgroupspecific, supportgroupid)
		}
		logger.Log.Println("slarecords.PushTime           >>>", slarecords.PushTime)
		//result = result - slarecords.PushTime
		doneSec := (returnValue.Resolutiontimeinsec - result)
		percent := (float64(doneSec) / float64(returnValue.Resolutiontimeinsec)) * 100

		//For Response meter
		if slarecords.Isresponsecomplete == 0 {
			availableTime1 := int64(0)
			result1 := int64(0)
			if today.Unix() > TimeParse(duedatetimeresponse, "").Unix() {
				result1, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, TimeParse(duedatetimeresponse, ""), startdatetimeresponse, today.Format("2006-01-02 15:04:05"), availableTime1, supportgroupspecific, supportgroupid)
				result1 = -result1
			} else {
				result1, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, today, startdatetimeresponse, duedatetimeresponse, availableTime1, supportgroupspecific, supportgroupid)
			}
			doneSec1 := (returnValue.Responsetimeinsec - result1)
			percent1 := (float64(doneSec1) / float64(returnValue.Responsetimeinsec)) * 100
			respRemainingTime = result1
			respPercent = percent1
		}
		logger.Log.Println("resolution complete flag ", slarecords.Isresolutioncomplete)
		logger.Log.Println("response complete flag ", slarecords.Isresponsecomplete)
		// New condition for upgrade SLA, login if response completed then never violate during upgrade or downgrade
		if math.Signbit(float64(respRemainingTime)) && slarecords.Isresponsecomplete != 0 {
			UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 0)
		}
		if math.Signbit(float64(result)) {
			UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 1)
		}
		logger.Log.Println("result value is   >>>", result)
		uResult, _, _ := UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, result, percent, respRemainingTime, respPercent)
		logger.Log.Println("is update >>>", uResult)
	}
}

/*func CalculateWorkingDays(clientid int64, mstorgnhirarchyid int64, startDate time.Time, endDate time.Time, dayCount int64) int64 {

	strDate, _ := GetStringDate(startDate)
	strendDate, _ := GetStringDate(endDate)
	dayofweekid := GetWeekDay(strDate)
	logger.Log.Println("Day number >>>>>>>> ", dayofweekid)
	starttime, endtime, err1 := GetClientDayOfWeekNew(clientid, mstorgnhirarchyid, dayofweekid)
	if err1 != nil {
	}
	if starttime != 0 || endtime != 0 {
		dayCount = dayCount + 1
	}
	logger.Log.Println("start date >>>>>>>> ", strDate)
	logger.Log.Println("end date >>>>>>>> ", strendDate)
	if strDate == strendDate {
		logger.Log.Println("****************** End ************************")
		return dayCount
	} else {
		nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
		result := CalculateWorkingDays(clientid, mstorgnhirarchyid, nextDay, endDate, dayCount)
		return result
	}

}*/

func CalculateWorkingDays(clientid int64, mstorgnhirarchyid int64, startDate time.Time, endDate time.Time, dayCount int64, supportgroupspecific int64, supportgroupid int64) int64 {

	strDate, _ := GetStringDate(startDate)
	strendDate, _ := GetStringDate(endDate)
	dayofweekid := GetWeekDay(strDate)
	starttime := int64(0)
	endtime := int64(0)
	logger.Log.Println("Day number >>>>>>>> ", dayofweekid)

	logger.Log.Println("startDate value is  >>>22222222222222222222222222222>>>>>>>>>>>>> ", startDate)

	logger.Log.Println("endDate value is   >>>>>>>>>>>>>3333333333333333333333>>>>>>>>>>>>>>>>>> ", endDate)

	if supportgroupspecific == 1 {
		starttime, endtime, _, _, _, _, _ = GetSupportGroupDayOfWeek(clientid, mstorgnhirarchyid, supportgroupid, strDate)
	} else {
		starttime, endtime, _, _, _, _ = GetClientDayOfWeek(clientid, mstorgnhirarchyid, strDate)
	}

	/* starttime, endtime, err1 := GetClientDayOfWeekNew(clientid, mstorgnhirarchyid, dayofweekid)
	   if err1 != nil {
	   } */
	if starttime != 0 || endtime != 0 {
		dayCount = dayCount + 1
	}
	logger.Log.Println("start date >>>>>>>> ", strDate)
	logger.Log.Println("end date >>>>>>>> ", strendDate)
	if strDate == strendDate {
		logger.Log.Println("****************** End ************************")
		return dayCount
	} else {
		nextDay := AddSubSecondsToDate(TimeParse(strDate+" 00:00:00", ""), 86400)
		result := CalculateWorkingDays(clientid, mstorgnhirarchyid, nextDay, endDate, dayCount, supportgroupspecific, supportgroupid)
		return result
	}

}
