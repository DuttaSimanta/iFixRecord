package models

import (
	"fmt"
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"strings"
	"time"
	"database/sql"
)

//var lock = &sync.Mutex{}

func TimeParse(somedate string, layout string) time.Time {
	if layout == "" {
		layout = "2006-01-02T15:04:05.000Z"
	}
	somedate = strings.Replace(somedate, " ", "T", 1) + ".000Z"
	t, _ := time.Parse(layout, somedate)
	return t
}

func DateTimeFormat(somedate time.Time, format string) string {
	if format == "" {
		return somedate.Format("2006-01-02 15:04:05")
	} else {
		return somedate.Format(format)
	}
}

func GetWeekDay(somedate string) int64 {
	date := TimeParse(somedate+" 00:00:00", "")
	dt := date.Weekday()
	today := strings.ToUpper(dt.String())
	if today == "SUNDAY" {
		return 1
	} else if today == "MONDAY" {
		return 2
	} else if today == "TUESDAY" {
		return 3
	} else if today == "WEDNESDAY" {
		return 4
	} else if today == "THURSDAY" {
		return 5
	} else if today == "FRIDAY" {
		return 6
	} else {
		return 7
	}
}

//This function is for add or subtract seconds to date
func AddSubSecondsToDate(someDate time.Time, seconds int64) time.Time {
	timein := someDate.Add(time.Second * time.Duration(seconds))
	return timein
}

//This function is for subttract date to date
func SubtractDateToDate(someDate time.Time, anotherDate time.Time) int64 {
	diff := int64((someDate.Sub(anotherDate) / time.Second))
	return diff
}

func CalculateBill(price int, no int) int {
	var totalPrice = price * no
	return totalPrice
}

//Multiple Return Value
// Call area, perimeter := rectProps(10.8, 5.6)
// area, _ := rectProps(10.8, 5.6) // perimeter is discarded

func rectProps(length, width float64) (float64, float64) {
	var area = length * width
	var perimeter = (length + width) * 2
	return area, perimeter
}

/* Named Return - no return variable required

func rectProps(length, width float64)(area, perimeter float64) {
    area = length * width
    perimeter = (length + width) * 2
    return //no explicit return value
}

*/

// to check - https://golangbot.com/go-packages/
// use of label

func CheckSlaFromWFTransition() bool {

	return false
}

func SLADueCheckforOUHoliday() bool {

	return false
}

func SLADueCheckforOUOfficeHours() bool {

	return false
}

func SLAPrelDueCheck(clientid int, ouid int, recordid int) bool {
	// The Parameters - Need to be discussed?
	// Will a SLA be always Client/OU or Support Group Specific?
	// Check which SLA is associated for the Record - How to check this - Need to be discussed?
	// Check SLA is specific to Support Group
	// If not - Check SLA Availability and Office Hour for Client
	// Else - get the Support Group
	//   Check SLA Availability and Office Hour for Specific Support Group
	// Can it be SG Timezone - Need to Discuss?
	//

	return false
}

/*

1.   Selection for Date Time
2.   Check SLA is within OU or SG
3.   OU/SG - TimeZone wise Time Selection
4.   Check Holiday - Start Time/End Time
5.   Check Office Hour - Start Time/End Time
4.1. for 4, Add One Day - Same Check Recursion
4.2  If false,
5.1. Office Hour - Add Start Hour/End Hour to Date
5.2. Check Current Time in that - Mark as Sla Duration Start
5.3. Add Time for Resposponse SLA End Time
6.   Check the Response SLA End Date Time for Holiday and Office Hour as 4
7.   Mark the SLA END Time in a Office Hour
8.   Return SLA Start Time and End Time

*/

/*

Whole Process
--At the time of record retrieval
1. Get SLA Fullfillment Criteria based on Ticket Type, Working Category, Priority
1.1 (To be done) effect on timezone - will work at last

2. Make SLA Call For All SLA Due Time
3. Check Transition and Make All History and Other Records
4. Check ResumePause and Make History and Other Records
5. Return Clock related all information
6. Stop the SLA

*/

/*

2. Explanation
2.1. Check the Time at record creation - Starting hour this or when it will start? - assumption - starting hour is record creation time
2.2. Check the SLA from No 1 - SG or OU Specific?
2.3. As per 2.2, SG or OU - check the day holiday range
2.4. If not Holiday go for Office Hour - SG or OU -->> 2.8
2.5. If the current time is within Holiday range - start time will be end of holiday range -if that is office hour
2.6. If the current time is before Holday range - start time will be now - if that is office hour
2.7. If the current time is not within holiday range - goto 2.8
2.8. If the current time is between Office Hour - mark the start time
2.9. Get End Time calculated from response and resolution SLA
2.10. Make the same calculation as start time for end time marking

What are the conditions - one can face for start time -

A.   The date for the time is holiday
A.1. The time can be before holday start or after holiday end
A.2. If the condition is A.1 check the hour is within office hour
     - if yes, mark the start hour
     - if no, Goto Next Day and Check for A.
A.3. If the time within holiday start and holiday end
     - check after holiday - office hour exists or not
     - if yes - mark start time
     - if no - Goto Next Day and Check for A.
B.   The date for the time is not holiday
B.1. The time can be before or after office hour
     - if before - mark office hour start as start hour
     - if after - Check for A.
   - if no mark the start hour.

*/

func InsertTrnslaentityhistory(tz entities.TrnslaentityhistoryEntity) (int64, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, err, "Something Went Wrong"
	// }

	t := int64(0)
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	lastInsertedId, err1 := dataAccess.InsertTrnslaentityhistory(&tz)
	fmt.Println(lastInsertedId)
	if err1 != nil {
		return 0, err1, "Something Went Wrong"
	}
	return lastInsertedId, nil, ""
}

func InsertMstsladue(tz entities.MstsladueEntity) (int64, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, err, "Something Went Wrong2222222222222222222222222222"
	// }
	t := int64(0)
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, err, "Something Went Wrong.."
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	lastInsertedId, err1 := dataAccess.InsertMstsladue(&tz)
	// p(t)
	if err1 != nil {
		return 0, err1, "Something Went Wrong2222222222222222222222222222222222"
	}
	return lastInsertedId, nil, ""
}

func GetMstsladue(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.MstsladueEntity, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }
	t := entities.MstsladueEntity{}
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.GetMstsladue(clientid, mstorgnhirarchyid, therecordid)
	// p(t)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	return t, true, nil, ""
}

func GetTrnslaentityhistory(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.TrnslaentityhistoryEntity, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()

	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("GetTrnslaentityhistory Get Statement Prepare Error", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	t := entities.TrnslaentityhistoryEntity{}
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.GetTrnslaentityhistory(clientid, mstorgnhirarchyid, therecordid)
	// p(t)
	if err1 != nil {
		logger.Log.Println("GetTrnslaentityhistory Get Statement Prepare Error", err1)
		return t, false, err1, "Something Went Wrong"
	}
	return t, true, nil, ""
}

func GetSupportGroupId(clientid int64, mstorgnhirarchyid int64, mstslafullfillmentcriteriaid int64) (int64, bool, error, string) {
	// p := fmt.Println
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }

	t := int64(0)
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	var err error
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.GetSupportGroupId(clientid, mstorgnhirarchyid, mstslafullfillmentcriteriaid)
	// p(t)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}

	return t, true, err, ""
}

func GetSupportGroupHoliday(clientid int64, mstorgnhirarchyid int64, supportGroupId int64, today string) (int64, int64, string, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return 0, 0, "", false, err, "Something Went Wrong"
	// }

	starttime := int64(0)
	endtime := int64(0)
	var dateofholiday string
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return 0, 0, "", false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	starttime, endtime, dateofholiday, err1 := dataAccess.GetSupportGroupHoliday(clientid, mstorgnhirarchyid, supportGroupId, today)
	if err1 != nil {
		return 0, 0, "", false, err1, "Something Went Wrong"
	}
	if starttime == 0 && endtime == 0 {
		return starttime, endtime, dateofholiday, false, nil, ""
	} else {
		// p(AddSubSecondsToDate(TimeParse(dateofholiday, ""), starttime))
		return starttime, endtime, dateofholiday, true, nil, ""
	}
}

func GetClientHoliday(clientid int64, mstorgnhirarchyid int64, today string) (int64, int64, string, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return 0, 0, "", false, err, "Something Went Wrong"
	// }
	starttime := int64(0)
	endtime := int64(0)
	var dateofholiday string
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return 0, 0, "", false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	starttime, endtime, dateofholiday, err1 := dataAccess.GetClientHoliday(clientid, mstorgnhirarchyid, today)
	if err1 != nil {
		return 0, 0, "", false, err1, "Something Went Wrong"
	}
	if starttime == 0 && endtime == 0 {
		return starttime, endtime, dateofholiday, false, nil, ""
	} else {
		// p(AddSubSecondsToDate(TimeParse(dateofholiday, ""), starttime))
		return starttime, endtime, dateofholiday, true, nil, ""
	}
}

func GetClientDayOfWeek(clientid int64, mstorgnhirarchyid int64, today string) (int64, int64, string, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return 0, 0, "", false, err, "Something Went Wrong"
	// }
	starttime := int64(0)
	endtime := int64(0)
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return 0, 0, "", false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	dayofweekid := GetWeekDay(today)
	starttime, endtime, err1 := dataAccess.GetClientDayOfWeek(clientid, mstorgnhirarchyid, dayofweekid)
	if err1 != nil {
		return 0, 0, today, false, err1, "Something Went Wrong"
	}
	if starttime == 0 && endtime == 0 {
		return starttime, endtime, today, false, nil, ""
	} else {
		return starttime, endtime, today, true, nil, ""
	}
}

func GetSupportGroupDayOfWeek(clientid int64, mstorgnhirarchyid int64, supportGroupId int64, today string) (int64, int64, int64, string, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return 0, 0, 0, "", false, err, "Something Went Wrong"
	// }
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return 0, 0, 0, "", false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)

	dataAccess := dao.DbConn{DB: db}

	dayofweekid := GetWeekDay(today)
	starttime, endtime, nextdayforward, err1 := dataAccess.GetSupportGroupDayOfWeek(clientid, mstorgnhirarchyid, supportGroupId, dayofweekid)
	if err1 != nil {
		return 0, 0, 0, today, false, err1, "Something Went Wrong"
	}
	if starttime == 0 && endtime == 0 {
		return starttime, endtime, nextdayforward, today, false, nil, ""
	} else {
		return starttime, endtime, nextdayforward, today, true, nil, ""
	}
}

func SLACriteriaRespResl(clientid int64, orgnunitid int64, recordtypeid int64, recordworkingcategory int64, recordpriority int64) (entities.MstslafullfillmentcriteriaEntity, bool, error, string) {
	t := entities.MstslafullfillmentcriteriaEntity{}
	// p("1111")
	var indata = entities.MstslafullfillmentcriteriaEntity{}
	indata.Clientid = clientid
	indata.Mstorgnhirarchyid = orgnunitid
	indata.Mstrecorddifferentiationtickettypeid = recordtypeid
	indata.Mstrecorddifferentiationpriorityid = recordpriority
	indata.Mstrecorddifferentiationworkingcatid = recordworkingcategory
	// p("4444")
	//logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }

	// p("3333")
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}
	var err error
	t, err2 := dataAccess.GetSpecificMstslafullfillmentcriteriaWithoutWCat(&indata)
	if err2 != nil {
		logger.Log.Println(" sla id value is ---------------------------->", err2)
		return t, false, err2, "Something Went Wrong"
	}
	logger.Log.Println(" sla id value is ----------1111111111111111------------------>", t)
	return t, true, err, ""
}

func SLAConfigurationForResponseDueTime(therecordid int) int {

	return 0
}

func SLAConfigurationForResolutionDueTime(therecordid int) int {

	return 0
}

func IsSLAStartTimeWithinOUOfficeHour() bool {

	return false
}

func IsSLAStartTimeWithinSGOfficeHour() bool {

	return false
}

func SLADueCheckforSGHoliday(supgroupid int) bool {

	return false
}

func SLADueCheckforSGOfficeHours(supgroupid int) bool {

	return false
}

func IsSLASupportGroupSpecific(clientid int, ouid int, slaId int) bool {
	//Check whether SLA is support group specific or not

	return false
}

func GetSupportGroupforTheSLA(clientid int, ouid int, slaid int) bool {
	//return specific support group id

	return false
}

func SetSLATransactionHistory() bool {
	// Insert Record in SLA Transaction History

	return false
}

func SetSLADue() bool {
	// Insert Record in SLA Due

	return false
}

func CheckSLAConfiguration() bool {

	return false
}

func GetSLALatestStatus() bool {
	// Will get SLA Start Stop Indicator

	return false
}

func isTheSLAViolated() bool {
	// SLA Violation Rule Check

	return false
}

func GetUpgradeIndicationSLA() bool {

	return false
}

func GetDowngradeIndicationSLA() bool {

	return false
}

func GetOUTimeZone() bool {
	//OU - Organisation Unit

	return false
}

func GetSGTimeZone() bool {
	//SG - Support Group

	return false
}

// func UpdateRemainingPercent(clientid int64, mstorgnhirarchyid int64, therecordid int64, remainingtime int64, completepercent float64) (bool, error, string) {
// 	db, err := ConnectMySqlDb()
// 	if err != nil {
// 		return false, err, "Something Went Wrong"
// 	}
// 	dataAccess := dao.DbConn{DB: db}

// 	t, err1 := dataAccess.UpdateRemainingPercent(clientid, mstorgnhirarchyid, therecordid, remainingtime, completepercent)
// 	// p(t)
// 	if err1 != nil {
// 		return false, err1, "Something Went Wrong"
// 	}
// 	return t, nil, ""
// }

func UpdateRemainingPercent(clientid int64, mstorgnhirarchyid int64, therecordid int64, remainingtime int64, completepercent float64, responseremainingtime int64, responsepercentage float64) (bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return false, err, "Something Went Wrong"
	// }
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.UpdateRemainingPercent(clientid, mstorgnhirarchyid, therecordid, remainingtime, completepercent, responseremainingtime, responsepercentage)
	// p(t)
	if err1 != nil {
		return false, err1, "Something Went Wrong"
	}
	return t, nil, ""
}

func UpdateMstsladue(clientid int64, mstorgnhirarchyid int64, therecordid int64, duedatetimeresponse string, duedatetimeresolution string, duedatetimeresolutionint int64, duedatetimeresponseint int64, totalpushTime int64) (bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
	// 	return false, err, "Something Went Wrong1111111111111111111"
	// }

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.UpdateMstsladue(clientid, mstorgnhirarchyid, therecordid, duedatetimeresponse, duedatetimeresolution, duedatetimeresolutionint, duedatetimeresponseint, totalpushTime)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return false, err1, "Something Went Wrong11111111111111111111111111111"
	}
	return t, nil, ""
}

func UpdateResponseEndFlag(clientid int64, mstorgnhirarchyid int64, therecordid int64) (bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
	// 	return false, err, "Something Went Wrong"
	// }
	currentTime := time.Now().UTC()
	responseCompleteTime := currentTime.Format("2006-01-02 15:04:05")
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.UpdateResponseEndFlag(clientid, mstorgnhirarchyid, therecordid, responseCompleteTime)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return false, err1, "Something Went Wrong"
	}
	return t, nil, ""
}

// To get UTC difference
func Getutcdiff(clientid int64, mstorgnhirarchyid int64) (entities.ZoneEntity, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }
	t := entities.ZoneEntity{}

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.Getutcdiff(clientid, mstorgnhirarchyid)
	// fmt.Println(t)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	return t, true, nil, ""
}

func UpdateViolateFlag(clientid int64, mstorgnhirarchyid int64, therecordid int64, flag int) (bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println("UpdateViolateFlag ------------>", err)
	// 	return false, err, "Something Went Wrong"
	// }
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.UpdateViolateFlag(clientid, mstorgnhirarchyid, therecordid, flag)
	if err1 != nil {
		logger.Log.Println("UpdateViolateFlag ------------>", err1)
		return false, err1, "Something Went Wrong"
	}
	return t, nil, ""
}

func UpdateRessolutionEndFlag(clientid int64, mstorgnhirarchyid int64, therecordid int64) (bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// //defer db.Close()
	// if err != nil {
	// 	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
	// 	return false, err, "Something Went Wrong"
	// }
	currentTime := time.Now().UTC()
	completeTime := currentTime.Format("2006-01-02 15:04:05")
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.UpdateRessolutionEndFlag(clientid, mstorgnhirarchyid, therecordid, completeTime)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return false, err1, "Something Went Wrong"
	}
	return t, nil, ""
}

func GetLastPushRow(allValue []entities.TrnslaentityhistoryEntity) entities.TrnslaentityhistoryEntity {
	result := entities.TrnslaentityhistoryEntity{}
	for i := 0; i < len(allValue); i++ {
		// fmt.Println(allValue[i])
		if allValue[i].Slastartstopindicator == 2 {
			result = allValue[i]
		} else {
			break
		}

	}
	return result
}

func GetTrnslaentityhistorytype2(clientid int64, mstorgnhirarchyid int64, therecordid int64, trnId int64) (entities.TrnslaentityhistoryEntity, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("GetTrnslaentityhistorytype2 Get Statement Prepare Error", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	t := entities.TrnslaentityhistoryEntity{}
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	result, err1 := dataAccess.GetTrnslaentityhistorytype2(clientid, mstorgnhirarchyid, therecordid, trnId)
	// p(t)
	if err1 != nil {
		logger.Log.Println("GetTrnslaentityhistorytype2 Get Statement Prepare Error", err1)
		return t, false, err1, "Something Went Wrong"
	}
	t = GetLastPushRow(result)
	return t, true, nil, ""
}
func CheckTrnslaentityhistorytype(clientid int64, mstorgnhirarchyid int64, therecordid int64, trnId int64) (entities.TrnslaentityhistoryEntity, bool, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("GetTrnslaentityhistorytype2 Get Statement Prepare Error", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	t := entities.TrnslaentityhistoryEntity{}
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	result, err1 := dataAccess.CheckTrnslaentityhistorytype(clientid, mstorgnhirarchyid, therecordid, trnId)
	// p(t)
	if err1 != nil {
		logger.Log.Println("GetTrnslaentityhistorytype2 Get Statement Prepare Error", err1)
		return t, false, false, err1, "Something Went Wrong"
	}
	t = GetLastPushRow(result)
	isPushTimeCalculated := CheckPushtimeCalculatedOrNot(result)
	return t, isPushTimeCalculated, true, nil, ""
}
func CheckPushtimeCalculatedOrNot(allValue []entities.TrnslaentityhistoryEntity) bool {
	var result bool = false
	for i := 1; i < len(allValue); i++ {
		// fmt.Println(allValue[i])
		if allValue[i].PushTime > 0 {
			result = true
		} else {
			break
		}
	}
	return result
}
func UpdatePushTimeInHistory(historyId int64, pushtime int64) (bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
	// 	return false, err, "Something Went Wrong"
	// }
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.UpdatePushTimeInHistory(historyId, pushtime)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return false, err1, "Something Went Wrong"
	}
	return t, nil, ""
}

// This method will return the starttime and endtime of the week day
func GetClientDayOfWeekNew(clientid int64, mstorgnhirarchyid int64, dayofweekid int64) (int64, int64, error) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
	// 	return 0, 0, err
	// }
	starttime := int64(0)
	endtime := int64(0)
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return 0, 0, err
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	var getclientWeekDay = `SELECT starttimeinteger, endtimeinteger FROM mstclientdayofweek WHERE clientid = ? AND mstorgnhirarchyid = ?  AND dayofweekid = ? AND deleteflg = 0 AND activeflg = 1`
	rows, err := db.Query(getclientWeekDay, clientid, mstorgnhirarchyid, dayofweekid)
	defer rows.Close()
	if err != nil {
		return 0, 0, err
	}
	for rows.Next() {
		err := rows.Scan(&starttime, &endtime)
		if err != nil {
		}
	}
	return starttime, endtime, nil
}

func GetTrnslaentityhistoryLastPushTime(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.TrnslaentityhistoryLastPushEntity, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("GetTrnslaentityhistory Get Statement Prepare Error", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	t := entities.TrnslaentityhistoryLastPushEntity{}
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.GetTrnslaentityhistoryLastPushTime(clientid, mstorgnhirarchyid, therecordid)
	// p(t)
	if err1 != nil {
		logger.Log.Println("GetTrnslaentityhistory Get Statement Prepare Error", err1)
		return t, false, err1, "Something Went Wrong"
	}
	return t, true, nil, ""
}

func GetrecordFromsladue(recordId int64) ([]entities.SLASchedulerEntity, bool, error) {
	t := []entities.SLASchedulerEntity{}
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	logger.Log.Println("inside GetrecordFromsladue")
	t, err1 := dataAccess.GetrecordFromsladue(recordId)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return t, false, err1
	}
	return t, true, nil
}

func UpdateExecuteFlag(id []int64, flag int64) (bool, error, string) {
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return false, err, ""
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	t, err1 := dataAccess.UpdateExecuteFlag(id, flag)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return false, err1, "Something Went Wrong"
	}
	return t, nil, ""
}

// Santanu Dhabal

func GetCategoryIds(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.SLATabEntity, bool, error) {

	t := entities.SLATabEntity{}
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.GetCategoryIds(clientid, mstorgnhirarchyid, therecordid)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return t, false, err1
	}
	return t, true, nil
}

func FetchCurrentGrpID(ID int64) (int64, error) {
	// logger.Log.Println("inside FetchCurrentGrpID")

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return 0, err
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	// fmt.Println("UpdateRecordFullFillDetails >>>>>>>>>>>>>> ", details)
	groupId, err1 := dataAccess.FetchCurrentGrpID(ID)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return 0, err1
	}
	return groupId, err1
}

// Santanu Dhabal
// This method will update details of a record
// 11.08.2021

func UpdateRecordFullFillDetails(details entities.SLARcordFullfillUpdate) (string, bool, error) {
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return "", false, err
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	// fmt.Println("UpdateRecordFullFillDetails >>>>>>>>>>>>>> ", details)
	_, err1 := dataAccess.UpdateRecordFullFillDetails(details)
	// p(t)
	if err1 != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
		return "", false, err1
	}
	return "", true, nil
}


func GetFirstStartTimeResolution(clientid int64, mstorgnhirarchyid int64, therecordid int64, db *sql.DB) (string,int64, bool, error) {
    // lock.Lock()
    // defer lock.Unlock()
    // db, err := dbconfig.ConnectMySqlDb()
    // t := entities.MstsladueEntity{}
    // // defer db.Close()
    // if err != nil {
    //  logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
    //  return t, false, err, "Something Went Wrong"
    // }
    dataAccess := dao.DbConn{DB: db}

    t, t1,err1 := dataAccess.GetFirstStartTimeResolution(clientid, mstorgnhirarchyid, therecordid)
    // p(t)
    if err1 != nil {
        logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err1)
        return t,t1, false, err1
    }
    return t,t1, true, nil
}




func GetMstSLADuePrevPushTime(recordID int64, id int64) (int64, error) {
	// logger.Log.Println("In side GetMstSLADueRowcount")
	logger.Log.Println("In side GetMstSLADueRowcount----->", recordID)
	var pushtime int64
	// var grpname string
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return pushtime, err
		}
		db = dbcon
	}
	var sql = "select sum(pushtime) from mstsladue where therecordid=? and id<? and deleteflg=0"
	rows, err := db.Query(sql, recordID, id)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("getStageID Get Statement Prepare Error", err)
		return pushtime, err
	}
	for rows.Next() {
		err = rows.Scan(&pushtime)
		logger.Log.Println("getStageID rows.next() Error", err)
	}
	return pushtime, nil
}
func GetMstSLADueRowcount(recordID int64) (int64, error) {
	// logger.Log.Println("In side GetMstSLADueRowcount")
	logger.Log.Println("In side GetMstSLADueRowcount----->", recordID)
	var count int64
	// var grpname string
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return count, err
		}
		db = dbcon
	}
	var sql = "select count(id) from mstsladue where therecordid=? and deleteflg=0"
	rows, err := db.Query(sql, recordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("getStageID Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("getStageID rows.next() Error", err)
	}
	return count, nil
}



func Getcurrentsatusid(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, int64, string, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var recorddiffid int64
	var name string
	var seqno int64
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return recorddiffid, seqno, name, err
		}
		db = dbcon
	}
	var latestsatusid = "SELECT a.recorddiffid,b.seqno,b.name FROM maprecordtorecorddifferentiation a,mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=3 AND a.recordid=? AND a.islatest=1 AND a.recorddiffid=b.id and b.seqno in(3,5,13);"
	rows, err := db.Query(latestsatusid, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getcurrentsatusid Get Statement Prepare Error", err)
		return recorddiffid, seqno, name, err
	}
	for rows.Next() {
		rows.Scan(&recorddiffid, &seqno, &name)

	}
	return recorddiffid, seqno, name, nil
}
func GetRecordCurrentGrpID(ClientID int64, OrgnIDTypeID int64, recordID int64) (int64, string, error) {
	logger.Log.Println("In side getStageID")
	logger.Log.Println("In side getStageID----->", recordID)
	var grpID int64
	var grpname string
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return grpID, grpname, err
		}
		db = dbcon
	}
	var sql = "select assignedgroupid,assignedgroup from recordfulldetails where clientid=? and mstorgnhirarchyid=? and recordid=? and deleteflg=0"
	rows, err := db.Query(sql, ClientID, OrgnIDTypeID, recordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("getStageID Get Statement Prepare Error", err)
		return grpID, grpname, err
	}
	for rows.Next() {
		err = rows.Scan(&grpID, &grpname)
		logger.Log.Println("getStageID rows.next() Error", err)
	}
	return grpID, grpname, nil
}





//=========================================To get last grpid=====================================================
func GetLastSupportgrp(requestid int64, supportgroupid int64) (int64, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }
	var t int64 = 0
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.GetLastSupportgrp(requestid, supportgroupid)
	// p(t)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	return t, true, nil, ""
}
func Getrequestid(therecordid int64) (int64, bool, error, string) {
	// logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := ConnectMySqlDb()
	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }
	var t int64 = 0
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection in side CreateRecordModel")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	logger.Log.Println("Database Open Connection Count is in Record Model ---------------------------->", db.Stats().OpenConnections)
	dataAccess := dao.DbConn{DB: db}

	t, err1 := dataAccess.Getrequestid(therecordid)
	// p(t)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	return t, true, nil, ""
}
//=========================================To get last grpid=====================================================

