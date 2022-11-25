package models

import (
	"fmt"
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"math"
	"time"
)

func GetSLATabvalues(page *entities.SLATabEntity) (entities.SLATabresponsesEntity, bool, error, string) {
	logger.Log.Println("In side GetSLATabvalues")
	t := entities.SLATabresponsesEntity{}
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()

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
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}
	issladue, err := dataAccess.FetchSLADueRow(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	if err != nil {
		logger.Log.Println("error is ----->", err)

		return t, false, err, "Something Went Wrong"
	}
	if issladue > 0 {


		// ------------- 15.03.2022 -------------------------------------------------
		res, err := dataAccess.Getrecorddetails(page.RecordID)
		if err != nil {
			logger.Log.Println(err)
			return t, false, err, "Something Went Wrong"
		}
		returnValue, _, _, _ := SLACriteriaRespResl(page.ClientID, page.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)
		currentTime := time.Now()
		//strDate := currentTime.Format("2006-01-02")

		zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
		datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)
		strDate, strtsec := GetStringDate(datetime)

		logger.Log.Println("YYYY-MM-DD : ", currentTime.Format("2006-01-02"))
		starttime := int64(0)
		endtime := int64(0)
		isHolliday := bool(false)
		dayweemstarttime := int64(0)
		dayweemendtime := int64(0)
		logger.Log.Println("strDate  value is   ---------------------------->", strDate)
		if returnValue.Supportgroupspecific == 1 {

			grpID, err := dataAccess.GetCurrentSupportGRP(page.RecordID)
			if err != nil {
				logger.Log.Println("error is ----->", err)

				return t, false, err, "Something Went Wrong"
			}
			starttime, endtime, _, isHolliday, _, _ = GetSupportGroupHoliday(page.ClientID, page.Mstorgnhirarchyid, grpID, strDate)
			dayweemstarttime, dayweemendtime, _, _, _, _, _ = GetSupportGroupDayOfWeek(page.ClientID, page.Mstorgnhirarchyid, grpID, strDate)
		} else {
			starttime, endtime, _, isHolliday, _, _ = GetClientHoliday(page.ClientID, page.Mstorgnhirarchyid, strDate)
			dayweemstarttime, dayweemendtime, _, _, _, _ = GetClientDayOfWeek(page.ClientID, page.Mstorgnhirarchyid, strDate)
		}

		logger.Log.Println("starttime, endtime, _, isHolliday  value is  ---------------------------->", starttime, endtime, isHolliday)
		logger.Log.Println("dayweemstarttime, dayweemendtime  value is  ---------------------------->", dayweemstarttime, dayweemendtime)


		var resultTime int64
		if strtsec > dayweemstarttime {
			dayweemstarttime = strtsec
		}

		resultTime = calculateWorkingTime(dayweemstarttime, dayweemendtime, starttime, endtime, isHolliday)
		logger.Log.Println("resultTime  value is  ---------------------------->", resultTime)

		// --------------- 15.03.2022 -------------------------------------------------



		resolutionvalues, err1 := dataAccess.GetResolutiondetails(page,resultTime)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}

		responsevalues, err1 := dataAccess.GetResponsedetails(page,resultTime)
		if err1 != nil {
			return t, false, err1, "Something Went Wrong"
		}
		t.Resolutionetails = resolutionvalues
		t.Responsedetails = responsevalues
	}
	holidayvalues, err1 := dataAccess.GetHolidaydetails(page)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	t.Holidaydetails = holidayvalues

	return t, true, err, ""
}

// func GetSLAResolution(page *entities.SLATabEntity) (entities.SLAMeterEntity, bool, error, string) {

// 	logger.Log.Println("GetSLAResolution page.ClientID no ---->", page.ClientID)
// 	logger.Log.Println("GetSLAResolution page.Mstorgnhirarchyid no ---->", page.Mstorgnhirarchyid)
// 	logger.Log.Println("GetSLAResolution page.RecordtypeID no ---->", page.RecordtypeID)
// 	logger.Log.Println("GetSLAResolution page.WorkingcatID no ---->", page.WorkingcatID)
// 	logger.Log.Println("GetSLAResolution page.PriorityID no ---->", page.PriorityID)

// 	currentTime := time.Now().UTC()
// 	t := entities.SLAMeterEntity{}
// 	//today := utility.TimeParse(currentTime.Format("2006-01-02 15:04:05"), "")
// 	// Get SLA total response and resolution time

// 	zonediff, _, _, _ := utility.Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
// 	today := utility.AddSubSecondsToDate(currentTime, zonediff.UTCdiff)

// 	returnValue, _, _, _ := utility.SLACriteriaRespResl(page.ClientID, page.Mstorgnhirarchyid, page.RecordtypeID, page.WorkingcatID, page.PriorityID)
// 	logger.Log.Println(returnValue.Responsetimeinsec)
// 	logger.Log.Println(returnValue.Resolutiontimeinsec)

// 	slarecords, _, _, _ := utility.GetMstsladue(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
// 	if slarecords.Id == 0 {
// 		return t, true, nil, ""
// 	}
// 	logger.Log.Println("sla due record ", slarecords)
// 	startdatetimeresponse := slarecords.Startdatetimeresponse
// 	duedatetimeresponse := slarecords.Duedatetimeresponse

// 	startdatetimeresolution := slarecords.Startdatetimeresolution
// 	duedatetimeresolution := slarecords.Duedatetimeresolution

// 	// slarecords.Responseremainingtime
// 	// slarecords.Responsepercentage

// 	respRemainingTime := slarecords.Responseremainingtime
// 	respPercent := slarecords.Responsepercentage
// 	if respPercent >= 100 {
// 		respPercent = 100
// 	}

// 	logger.Log.Println("slarecords entity  >>>", slarecords)
// 	logger.Log.Println("respRemainingTime   >>>", respRemainingTime)

// 	trnslarecords, _, _, _ := utility.GetTrnslaentityhistory(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
// 	logger.Log.Println("In time second >>>", slarecords.Remainingtime)

// 	if trnslarecords.Slastartstopindicator == 2 {
// 		logger.Log.Println("In time second >>>", slarecords.Remainingtime)
// 		logger.Log.Println("in percentage >>>>> ", slarecords.Completepercent)
// 		today = utility.TimeParse(trnslarecords.Recorddatetime, "")

// 	} else if trnslarecords.Slastartstopindicator == 4 {
// 		uResult, _, _ := utility.UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, slarecords.Remainingtime, slarecords.Completepercent, respRemainingTime, respPercent)
// 		fmt.Println("is update >>>", uResult)
// 		t.Remainresolutiontime = slarecords.Remainingtime
// 		t.Resolutionpercent = slarecords.Completepercent
// 		t.Remainresponsetime = respRemainingTime
// 		t.Responsepercent = respPercent

// 		return t, true, nil, ""
// 	}

// 	// For Resolution meter
// 	availableTime := int64(0)
// 	result := int64(0)
// 	if today.Unix() > utility.TimeParse(duedatetimeresolution, "").Unix() {
// 		result, _ = utility.GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, utility.TimeParse(duedatetimeresolution, ""), startdatetimeresolution, today.Format("2006-01-02 15:04:05"), availableTime)
// 		result = -result
// 	} else {
// 		result, _ = utility.GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, today, startdatetimeresolution, duedatetimeresolution, availableTime)
// 	}
// 	doneSec := (returnValue.Resolutiontimeinsec - result)
// 	percent := (float64(doneSec) / float64(returnValue.Resolutiontimeinsec)) * 100
// 	if percent >= float64(100) {
// 		percent = float64(100)
// 	}
// 	//For Response meter
// 	if slarecords.Isresponsecomplete == 0 {
// 		availableTime1 := int64(0)
// 		result1 := int64(0)
// 		if today.Unix() > utility.TimeParse(duedatetimeresponse, "").Unix() {
// 			result1, _ = utility.GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, utility.TimeParse(duedatetimeresponse, ""), startdatetimeresponse, today.Format("2006-01-02 15:04:05"), availableTime1)
// 			result1 = -result1
// 		} else {
// 			result1, _ = utility.GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, today, startdatetimeresponse, duedatetimeresponse, availableTime1)
// 		}
// 		doneSec1 := (returnValue.Responsetimeinsec - result1)
// 		percent1 := (float64(doneSec1) / float64(returnValue.Responsetimeinsec)) * 100
// 		logger.Log.Println("percent1   ------------------------>", percent1)
// 		respRemainingTime = result1
// 		if percent1 >= float64(100) {
// 			respPercent = float64(100)
// 		} else {
// 			respPercent = percent1
// 		}

// 	}

// 	logger.Log.Println("resolution complete flag ", slarecords.Isresolutioncomplete)
// 	logger.Log.Println("response complete flag ", slarecords.Isresponsecomplete)
// 	if math.Signbit(float64(respRemainingTime)) {
// 		utility.UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 0)
// 	}
// 	if math.Signbit(float64(result)) {
// 		utility.UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 1)
// 	}
// 	uResult, _, _ := utility.UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, result, percent, respRemainingTime, respPercent)

// 	logger.Log.Println("respRemainingTime   >>>", respRemainingTime)
// 	logger.Log.Println("is update >>>", uResult)
// 	t.Remainresolutiontime = result
// 	t.Resolutionpercent = percent
// 	t.Remainresponsetime = respRemainingTime
// 	t.Responsepercent = respPercent

// 	return t, true, nil, ""
// }

func GetSLAResolution(page *entities.SLATabEntity) (entities.SLAMeterEntity, bool, error, string) {

	logger.Log.Println("GetSLAResolution page.ClientID no ---->", page.ClientID)
	logger.Log.Println("GetSLAResolution page.Mstorgnhirarchyid no ---->", page.Mstorgnhirarchyid)
	logger.Log.Println("GetSLAResolution page.RecordtypeID no ---->", page.RecordtypeID)
	logger.Log.Println("GetSLAResolution page.WorkingcatID no ---->", page.WorkingcatID)
	logger.Log.Println("GetSLAResolution page.PriorityID no ---->", page.PriorityID)
	logger.Log.Println("GetSLAResolution page.RecordID no ---->", page.RecordID)
	logger.Log.Println("GetSLAResolution page.SupportgroupId --11111111111111111111111-->", page.SupportgroupId)

	currentTime := time.Now().UTC()
	t := entities.SLAMeterEntity{}

	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	log.Println("database connection failure", err)
	// 	return t, true, nil, ""
	// }
	//defer db.Close()

	// lock.Lock()
	// defer lock.Unlock()

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, true, nil, ""
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in SLA Model ---------------------------->", db.Stats().OpenConnections)

	dataAccess := dao.DbConn{DB: db}

	grpID, _ := dataAccess.FetchCurrentGrpID(page.RecordID)
	// if err != nil {
	// 	return false, err, "Something Went Wrong"
	// }

	supportgroupid := grpID

	slaid, err := dataAccess.GetSLAdataexist(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	if err != nil {
		logger.Log.Println(err)
		return t, true, nil, ""
	}

	if slaid > 0 {

		zonediff, _, _, _ := Getutcdiff(page.ClientID, page.Mstorgnhirarchyid)
		today := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)

		returnValue, _, _, _ := SLACriteriaRespResl(page.ClientID, page.Mstorgnhirarchyid, page.RecordtypeID, page.WorkingcatID, page.PriorityID)
		if returnValue.Id == 0 {
			logger.Log.Println("******************************** No Record found from SLACriteriaRespResl *************************************")
			return t, true, nil, ""
		}

		logger.Log.Println(returnValue.Responsetimeinsec)
		logger.Log.Println(returnValue.Resolutiontimeinsec)

		slarecords, _, _, _ := GetMstsladue(page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
		if slarecords.Id == 0 {
			logger.Log.Println("******************************** No Record found from GetMstsladue *************************************")
			return t, true, nil, ""
		}
		logger.Log.Println("sla due record ", slarecords)
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
		if trnslarecords.Id == 0 {
			logger.Log.Println("******************************** No Record found from GetTrnslaentityhistory *************************************")
			return t, true, nil, ""
		}
		logger.Log.Println("In time second >>>", slarecords.Remainingtime)
		if trnslarecords.Slastartstopindicator == 2 {
			logger.Log.Println("In time second >>>", slarecords.Remainingtime)
			logger.Log.Println("in percentage >>>>> ", slarecords.Completepercent)
			today = TimeParse(trnslarecords.Recorddatetime, "")
			if today.IsZero() {
				logger.Log.Println("****************************** Incorrect date capture *********************** ", today)
				return t, true, nil, ""
			}

		} else if trnslarecords.Slastartstopindicator == 4 {
			uResult, _, _ := UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, slarecords.Remainingtime, slarecords.Completepercent, respRemainingTime, respPercent)
			fmt.Println("is update >>>", uResult)
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
		logger.Log.Println("today.Unix()           ---->", today.Unix())
		logger.Log.Println("TimeParse(duedatetimeresolution).Unix()           ---->", TimeParse(duedatetimeresolution, "").Unix())

		if today.Unix() > TimeParse(duedatetimeresolution, "").Unix() {
			logger.Log.Println("today.Unix() > TimeParse(duedatetimeresolution, ).Unix()           ---->")
			result, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, TimeParse(duedatetimeresolution, ""), startdatetimeresolution, today.Format("2006-01-02 15:04:05"), availableTime, returnValue.Supportgroupspecific, supportgroupid)
			result = -result
		} else {
			result, _ = GetSLARemainingTimeForClient(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, today, startdatetimeresolution, duedatetimeresolution, availableTime, returnValue.Supportgroupspecific, supportgroupid)
		}
		logger.Log.Println("slarecords.PushTime                 >>>", slarecords.PushTime)
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
		logger.Log.Println("resolution complete flag ", slarecords.Isresolutioncomplete)
		logger.Log.Println("response complete flag ", slarecords.Isresponsecomplete)
		// New condition for upgrade SLA, login if response completed then never violate during upgrade or downgrade
		if math.Signbit(float64(respRemainingTime)) && slarecords.Isresponsecomplete == 0 {
			UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 0)
		}
		if math.Signbit(float64(result)) {
			UpdateViolateFlag(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, 1)
		}
		if trnslarecords.Slastartstopindicator != 2 {
			uResult, _, _ := UpdateRemainingPercent(page.ClientID, page.Mstorgnhirarchyid, page.RecordID, result, percent, respRemainingTime, respPercent)
			logger.Log.Println("is update >>>", uResult)
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
	}
	return t, true, nil, ""
}
