package models

import (
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"math"
	"time"
        "log"
)

func UpdateStagingTableDetails(page *entities.SLATabEntity) (entities.SLAMeterEntity, bool, error, string) {
	p := logger.Log.Println
	// p := fmt.Println
	// p("GetSLAResolution page.ClientID no ---->", page.ClientID)
	// p("GetSLAResolution page.Mstorgnhirarchyid no ---->", page.Mstorgnhirarchyid)
	// p("GetSLAResolution page.RecordtypeID no ---->", page.RecordtypeID)
	// p("GetSLAResolution page.WorkingcatID no ---->", page.WorkingcatID)
	// p("GetSLAResolution page.PriorityID no ---->", page.PriorityID)
	t := entities.SLAMeterEntity{}
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, true, nil, ""
		}
		db = dbcon
	}

	currentTime := time.Now().UTC()

	var runningTimewh = int64(0)
	var runningTimeincludeleave = int64(0)
      //  var supportgroupid int64
	/*_, diffseq, _, err := Getcurrentsatusid(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	if err != nil {
		logger.Log.Println(err)

	}
	if diffseq != 0 {
		currentgrp, _, err := GetRecordCurrentGrpID(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		if err != nil {
			logger.Log.Println(err)
		}
		if currentgrp != 0 {
			supportgroupid = currentgrp
		} else {
			logger.Log.Println("NOT GETTING GRPID")
		}

	}*/



	//supportgroupid, _ = FetchCurrentGrpID(page.RecordID)
	supportgroupid, _, err := GetRecordCurrentGrpID(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	if err != nil {
		logger.Log.Println(err)
	}
	p("Current support grp id is >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", supportgroupid)
	zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
	today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
	todayCurrentTime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)

	returnValue, _, _, _ := SLACriteriaRespResl(page.ClientID, page.Mstorgnhirarchyid, page.RecordtypeID, page.WorkingcatID, page.PriorityID)

	p(returnValue.Responsetimeinsec)
	p(returnValue.Resolutiontimeinsec)

	slarecords, _, _, _ := GetMstsladue(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	initialStartTimeResolution,prvpushtime, _, _ := GetFirstStartTimeResolution(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, db)
	logger.Log.Println("initialStartTimeResolution >>>> ", initialStartTimeResolution)
	// p("sla due record ", slarecords)
	startdatetimeresponse := slarecords.Startdatetimeresponse
	duedatetimeresponse := slarecords.Duedatetimeresponse

	startdatetimeresolution := slarecords.Startdatetimeresolution
	duedatetimeresolution := slarecords.Duedatetimeresolution

	// slarecords.Responseremainingtime
	// slarecords.Responsepercentage

	respRemainingTime := slarecords.Responseremainingtime
	respPercent := slarecords.Responsepercentage
	// if respPercent >= 100 {
	// 	respPercent = 100
	// }
	trnslarecords, _, _, _ := GetTrnslaentityhistory(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	if trnslarecords.Slastartstopindicator == 2 {
		// p("In time second >>>", slarecords.Remainingtime)
		// p("in percentage >>>>> ", slarecords.Completepercent)
		today = TimeParse(trnslarecords.Recorddatetime, "")

	} else if trnslarecords.Slastartstopindicator == 4 {
		_, _, _ = UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, slarecords.Remainingtime, slarecords.Completepercent, respRemainingTime, respPercent)
		// fmt.Println("is update >>>", uResult)
		t.Remainresolutiontime = slarecords.Remainingtime
		t.Resolutionpercent = slarecords.Completepercent
		t.Remainresponsetime = respRemainingTime
		t.Responsepercent = respPercent
		t.RecordID = page.RecordID
		return t, true, nil, ""
	}
	// For Resolution meter
	availableTime := int64(0)
	result := int64(0)
	if today.Unix() > TimeParse(duedatetimeresolution, "").Unix() {
		result, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, TimeParse(duedatetimeresolution, ""), startdatetimeresolution, today.Format("2006-01-02 15:04:05"), availableTime, returnValue.Supportgroupspecific, supportgroupid)
		result = -result
	} else {
		result, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, today, startdatetimeresolution, duedatetimeresolution, availableTime, returnValue.Supportgroupspecific, supportgroupid)
	}
	// p("slarecords.PushTime                 >>>", slarecords.PushTime)
	//result = result - slarecords.PushTime
	doneSec := (returnValue.Resolutiontimeinsec - result)
	percent := (float64(doneSec) / float64(returnValue.Resolutiontimeinsec)) * 100

	//For Response meter
	if slarecords.Isresponsecomplete == 0 {
		availableTime1 := int64(0)
		result1 := int64(0)
		if today.Unix() > TimeParse(duedatetimeresponse, "").Unix() {
			result1, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, TimeParse(duedatetimeresponse, ""), startdatetimeresponse, today.Format("2006-01-02 15:04:05"), availableTime1, returnValue.Supportgroupspecific, supportgroupid)
			result1 = -result1
		} else {
			result1, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, today, startdatetimeresponse, duedatetimeresponse, availableTime1, returnValue.Supportgroupspecific, supportgroupid)
		}
		doneSec1 := (returnValue.Responsetimeinsec - result1)
		percent1 := (float64(doneSec1) / float64(returnValue.Responsetimeinsec)) * 100
		respRemainingTime = result1
		respPercent = percent1

	}
	// p("resolution complete flag ", slarecords.Isresolutioncomplete)
	// p("response complete flag ", slarecords.Isresponsecomplete)
	// New condition for upgrade SLA, login if response completed then never violate during upgrade or downgrade
	if math.Signbit(float64(respRemainingTime)) && slarecords.Isresponsecomplete == 0 {
		UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 0)
	}
	if math.Signbit(float64(result)) {
		UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 1)
	}
	if trnslarecords.Slastartstopindicator != 2 {
		_, _, _ = UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, result, percent, respRemainingTime, respPercent)
		// p("is update >>>", uResult)
		t.Remainresolutiontime = result
		t.Resolutionpercent = percent
		t.Remainresponsetime = respRemainingTime
		t.Responsepercent = respPercent
		t.RecordID = page.RecordID

	} else {
		t.Remainresolutiontime = slarecords.Remainingtime
		t.Resolutionpercent = slarecords.Completepercent
		t.Remainresponsetime = respRemainingTime
		t.Responsepercent = respPercent
		t.RecordID = page.RecordID
	}

	// p("Remaining resolution time >>>>>>>>>>>>>> ", t.Remainresolutiontime)
	// fmt.Println("t.Remaining response time >>>>>>>>>>>>>> ", t.Remainresponsetime)
	// fmt.Println("t.Remaining response percentage >>>>>>>>>>>>>> ", t.Responsepercent)
	// fmt.Println("SLA push time >>>>>>>>>>>>>> ", slarecords.PushTime)

	var respoverduetime int64
	var respoverdueperc float64
	var resooverduetime int64
	var resooverdueperc float64
	if math.Signbit(float64(t.Remainresolutiontime)) {
		resooverduetime = t.Remainresolutiontime
	} else {
		resooverduetime = 0
	}
	if math.Signbit(float64(t.Remainresponsetime)) {
		respoverduetime = t.Remainresponsetime
	} else {
		respoverduetime = 0
	}
	// To calculate overdue percentage
	if t.Resolutionpercent > float64(100) {
		resooverdueperc = (t.Resolutionpercent - float64(100))
	} else {
		resooverdueperc = 0
	}
	if t.Responsepercent > float64(100) {
		// fmt.Println("t.Remaining response percentage >>>>>>>>>>>>>> overdue")
		respoverdueperc = (t.Responsepercent - float64(100))
	} else {
		respoverdueperc = 0
	}
	totalpushTime := slarecords.PushTime
	if trnslarecords.Slastartstopindicator == 2 {
		// p("returnValue.Clientid, returnValue.Mstorgnhirarchyid, therecordid, trnslarecords.Id >>> ", returnValue.Clientid, returnValue.Mstorgnhirarchyid, page.RecordID, trnslarecords.Id)
		pushRecord, _, _, _ := GetTrnslaentityhistorytype2(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, trnslarecords.Id)
		logger.Log.Println("pushRecord value is  ----------------------------------------------    >>>", pushRecord)
		if pushRecord.Id != 0 {
			pushDateTime := TimeParse(pushRecord.Recorddatetime, "")
			if pushDateTime.IsZero() {
				logger.Log.Println("****************************** Incorrect date capture *********************** ", pushDateTime)
			}
			pushTime := CalculateWorkingHourBetweenTwoDates(returnValue.Clientid, returnValue.Mstorgnhirarchyid, pushDateTime, todayCurrentTime, int64(0), returnValue.Supportgroupspecific, supportgroupid)
			// p("push time result >>>>>>>>> ", pushTime)
			totalpushTime = pushTime + slarecords.PushTime

		}
	}
	// p("slarecords.PushTime >>>>>>>>>>>>> ", totalpushTime)
	// p("111111111111111111 startdatetimeresolution >>>>>>>>>>>>>> ", startdatetimeresolution)
	// p("page.ClientID, page.Mstorgnhirarchyid, supportgroupid >>> ", page.ClientID, page.Mstorgnhirarchyid, supportgroupid)
	// p("returnValue.Supportgroupspecific >>> ", returnValue.Supportgroupspecific)
	// Calculate the sla running time in working hour
	//runningTimewh = CalculateWorkingHourBetweenTwoDates(page.ClientID, page.Mstorgnhirarchyid, TimeParse(startdatetimeresolution, ""), today, int64(0), returnValue.Supportgroupspecific, supportgroupid)
	// p("business aging >>> ", runningTimewh)
	//runningTimeincludeleave = SubtractDateToDate(today, TimeParse(startdatetimeresolution, ""))
	// p("calender aging >>> ", runningTimeincludeleave)
	// p("222222222222222222222 >>>>>>>>>>>>>> ")

	log.Println(prvpushtime)

	count, err := GetMstSLADueRowcount(page.RecordID)
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
	}
	if count > 1 {
		PrevIdletime, err := GetMstSLADuePrevPushTime(page.RecordID, slarecords.Id)
		if err != nil {
			logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
		}
		totalpushTime = totalpushTime + PrevIdletime
	}

	runningTimewh = CalculateWorkingHourBetweenTwoDates(page.ClientID, page.Mstorgnhirarchyid, TimeParse(initialStartTimeResolution, ""), todayCurrentTime, int64(0), returnValue.Supportgroupspecific, supportgroupid)
	runningTimeincludeleave = SubtractDateToDate(todayCurrentTime, TimeParse(initialStartTimeResolution, ""))


	recordfullfilldetails := entities.SLARcordFullfillUpdate{}
	recordfullfilldetails.Responseslameterpercentage = t.Responsepercent
	recordfullfilldetails.Resolutionslameterpercentage = t.Resolutionpercent
	recordfullfilldetails.Businessaging = runningTimewh
	recordfullfilldetails.Calendaraging = runningTimeincludeleave
	recordfullfilldetails.Actualeffort = 0
	recordfullfilldetails.Slaidletime = totalpushTime
	recordfullfilldetails.Respoverduetime = respoverduetime
	recordfullfilldetails.Resooverduetime = resooverduetime
	recordfullfilldetails.Respoverdueperc = respoverdueperc
	recordfullfilldetails.Resooverdueperc = resooverdueperc
	recordfullfilldetails.ClientID = page.ClientID
	recordfullfilldetails.Mstorgnhirarchyid = page.Mstorgnhirarchyid
	recordfullfilldetails.RecordID = page.RecordID

	_, _, uerr := UpdateRecordFullFillDetails(recordfullfilldetails)
	if uerr != nil {

		p("Recoed full fill details update >>>>>>>>>>>>>>>>> ", uerr)
	}
	return t, true, nil, ""
}

func ExecuteStagingDetailsUpdate(recordId int64) {

	p := logger.Log.Println
	p("inside ExecuteRemainingUpdate &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&  Started &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")

	slaSchedule, success, _ := GetrecordFromsladue(recordId)
	//p("slaSchedule data >>>> ", slaSchedule, success)
	if success {
		var ids []int64
		var recordids []int64
		for _, elem := range slaSchedule {
			// p(elem.ID)
			ids = append(ids, elem.ID)
			recordids = append(recordids, elem.RecordID) //This is only for printing in log done by josim
		}
		p("Recordids and sladueid going to update: ", recordids, len(recordids))
		// p(ids)
		_, err, _ := UpdateExecuteFlag(ids, 1)
		if err != nil {
			logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
		}
		for v, elem := range slaSchedule {
			p("************************ Executed Id ************************* ", elem.ID)
			page, success1, err := GetCategoryIds(elem.ClientID, elem.Mstorgnhirarchyid, elem.RecordID)
			p("page data >>>> ", page)
			if err != nil {
				logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
			}
			if success1 && page.ClientID != 0 {
				_, _, err, _ = UpdateStagingTableDetails(&page)
				if err != nil {
					logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
				}
			}
			//break
			logger.Log.Println("updated :", v)
		}
		_, err, _ = UpdateExecuteFlag(ids, 0)
		if err != nil {
			logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
		}
	}

	p("inside ExecuteRemainingUpdate &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&  Ended &&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
}
