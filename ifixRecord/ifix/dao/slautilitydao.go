package dao

import (
	//"database/sql"
	"fmt"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"strings"
)

var getMstslafullfillmentcriteria = `SELECT id as Id, clientid as Clientid, mstorgnhirarchyid as Mstorgnhirarchyid, 
slaid as Slaid, mstrecorddifferentiationtickettypeid as Mstrecorddifferentiationtickettypeid, 
mstrecorddifferentiationpriorityid as Mstrecorddifferentiationpriorityid, mstrecorddifferentiationworkingcatid as Mstrecorddifferentiationworkingcatid, 
responsetimeinhr as Responsetimeinhr, responsetimeinmin as Responsetimeinmin, responsetimeinsec as Responsetimeinsec, 
resolutiontimeinhr as Resolutiontimeinhr, resolutiontimeinmin as Resolutiontimeinmin, resolutiontimeinsec as Resolutiontimeinsec, 
supportgroupspecific as Supportgroupspecific, activeflg as Activeflg
FROM mstslafullfillmentcriteria
WHERE clientid = ? AND mstorgnhirarchyid = ?  AND mstrecorddifferentiationtickettypeid =? 
AND mstrecorddifferentiationpriorityid =? AND mstrecorddifferentiationworkingcatid =?
AND deleteflg =0 and activeflg=1`

var getsupportGroupid = `SELECT mstclientsupportgroupid FROM mstslaresponsiblesupportgroup WHERE clientid = ? AND mstorgnhirarchyid = ? AND mstslafullfillmentcriteriaid = ? AND activeflg = 1 AND deleteflg = 0`

var getclientsupportgroupholyday = `SELECT starttimeinteger, endtimeinteger, FROM_UNIXTIME(dateofholiday) as dateofholiday FROM mstclientsupportgroupholiday WHERE clientid = ? AND mstorgnhirarchyid = ? AND mstclientsupportgroupid = ? AND dateofholiday = round(UNIX_TIMESTAMP(?)) AND deleteflg = 0 AND activeflg = 1`

var getclientholyday = `SELECT starttimeinteger, endtimeinteger, FROM_UNIXTIME(dateofholiday) as dateofholiday FROM mstclientholiday WHERE clientid = ? AND mstorgnhirarchyid = ? AND dateofholiday = round(UNIX_TIMESTAMP(?)) AND deleteflg = 0 AND activeflg = 1`

var getclientWeekDay = `SELECT starttimeinteger, endtimeinteger FROM mstclientdayofweek WHERE clientid = ? AND mstorgnhirarchyid = ?  AND dayofweekid = ? AND deleteflg = 0 AND activeflg = 1`

var getsupportgroupWeekDay = `SELECT starttimeinteger, endtimeinteger, nextdayforward FROM mstclientsupportgroupdayofweek WHERE clientid = ? AND mstorgnhirarchyid = ? AND mstclientsupportgroupid = ?  AND dayofweekid = ? AND deleteflg = 0 AND activeflg = 1`

var insertTrnslaentityhistory = "INSERT INTO trnslaentityhistory (clientid, mstorgnhirarchyid, mstslaentityid, therecordid, recorddatetime, recorddatetoint, donotupdatesladue, recordtimetoint, mstslastateid, commentonrecord, slastartstopindicator, fromclientuserid) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

var getTrnslaentityhistory = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, recorddatetime, recorddatetoint, donotupdatesladue, recordtimetoint, mstslastateid, commentonrecord, slastartstopindicator, fromclientuserid FROM trnslaentityhistory WHERE clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 ORDER BY id DESC LIMIT 1`

//var updateRemainingPercent = `UPDATE mstsladue SET remainingtime = ?, completepercent = ? WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

//var updateMstsladue = `UPDATE mstsladue SET duedatetimeresponse = ?, duedatetimeresolution = ?, duedatetimeresolutionint = ?, duedatetimeresponseint = ?,pushtime=? WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

var updateMstsladue = `UPDATE mstsladue SET duedatetimeresponse = ?, duedatetimeresolution = ?, duedatetimeresolutionint = ?, duedatetimeresponseint = ?,pushtime=? WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ? order by id desc limit 1`

var insertMstsladue = "INSERT INTO mstsladue (clientid, mstorgnhirarchyid, mstslaentityid, therecordid, latestone, startdatetimeresponse, startdatetimeresolution, duedatetimeresponse, duedatetimeresolution, duedatetimetominute, resoltiondone, resolutiondatetime, lastupdatedattime, trnslaentityhistoryid, duedatetimeresolutionint, duedatetimeresponseint,remainingtime, completepercent, responseremainingtime, responsepercentage) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

var updateRemainingPercent = `UPDATE mstsladue SET remainingtime = ?, completepercent = ?, responseremainingtime = ?, responsepercentage = ? WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

//var getMstsladue = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, latestone, startdatetimeresponse, startdatetimeresolution, duedatetimeresponse, duedatetimeresolution, duedatetimeresolutionint, duedatetimeresponseint, resoltiondone, resolutiondatetime, trnslaentityhistoryid, remainingtime, completepercent, responseremainingtime, responsepercentage, isresponsecomplete FROM mstsladue WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ? AND activeflg = 1 AND deleteflg = 0 ORDER BY id DESC LIMIT 1`

var updateResponseEndFlag = `UPDATE mstsladue SET isresponsecomplete = 1, responseCompleteTime = ? WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

var getutcdiff = `SELECT b.utcdiff FROM mstslatimezone as a, zone as b WHERE a.clientid = ? AND a.mstorgnhirarchyid = ? AND a.msttimezoneid = b.zone_id AND a.deleteflg = 0 AND a.activeflg = 1`

//var getMstsladue = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, latestone, startdatetimeresponse, startdatetimeresolution, duedatetimeresponse, duedatetimeresolution, duedatetimeresolutionint, duedatetimeresponseint, resoltiondone, resolutiondatetime, trnslaentityhistoryid, remainingtime, completepercent, responseremainingtime, responsepercentage, isresponsecomplete, isresolutioncomplete FROM mstsladue WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ? AND activeflg = 1 AND deleteflg = 0 ORDER BY id DESC LIMIT 1`

//var updateRespViolateFlag = `UPDATE mstsladue SET isresponsecomplete = 2 WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

//var updateResolViolateFlag = `UPDATE mstsladue SET isresolutioncomplete = 2 WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

var updateRespViolateFlag = `UPDATE mstsladue SET isresponseviolation = 1 WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

var updateResolViolateFlag = `UPDATE mstsladue SET isresolutionviolation = 1 WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

var updateRessolutionEndFlag = `UPDATE mstsladue SET isresolutioncomplete = 1, resolutionCompleteTime = ? WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ?`

var getMstslafullfillmentcriteriawithoutworkingcategory = `SELECT id as Id, clientid as Clientid, mstorgnhirarchyid as Mstorgnhirarchyid,
slaid as Slaid, mstrecorddifferentiationtickettypeid as Mstrecorddifferentiationtickettypeid,
mstrecorddifferentiationpriorityid as Mstrecorddifferentiationpriorityid, mstrecorddifferentiationworkingcatid as Mstrecorddifferentiationworkingcatid,
responsetimeinhr as Responsetimeinhr, responsetimeinmin as Responsetimeinmin, responsetimeinsec as Responsetimeinsec,
resolutiontimeinhr as Resolutiontimeinhr, resolutiontimeinmin as Resolutiontimeinmin, resolutiontimeinsec as Resolutiontimeinsec,
supportgroupspecific as Supportgroupspecific, activeflg as Activeflg
FROM mstslafullfillmentcriteria
WHERE clientid = ? AND mstorgnhirarchyid = ?  AND mstrecorddifferentiationtickettypeid =?
AND mstrecorddifferentiationpriorityid =? AND deleteflg =0 and activeflg=1`

var getMstsladue = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, latestone, startdatetimeresponse, startdatetimeresolution, duedatetimeresponse, duedatetimeresolution, duedatetimeresolutionint, duedatetimeresponseint, resoltiondone, resolutiondatetime, trnslaentityhistoryid, remainingtime, completepercent, responseremainingtime, responsepercentage, isresponsecomplete, isresolutioncomplete, pushtime FROM mstsladue WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ? AND activeflg = 1 AND deleteflg = 0 ORDER BY id DESC LIMIT 1`

//var getTrnslaentityhistorytype2 = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, recorddatetime, recorddatetoint, donotupdatesladue, recordtimetoint, mstslastateid, commentonrecord, slastartstopindicator, fromclientuserid FROM trnslaentityhistory WHERE clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 AND slastartstopindicator = 2 ORDER BY id DESC LIMIT 1`


//var getTrnslaentityhistorytype2 = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, recorddatetime, recorddatetoint, donotupdatesladue, recordtimetoint, mstslastateid, commentonrecord, slastartstopindicator, fromclientuserid FROM trnslaentityhistory  WHERE clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 AND slastartstopindicator = 2 AND id <= ? AND id > (SELECT id FROM trnslaentityhistory b WHERE id < ? and slastartstopindicator=1 and b.therecordid=trnslaentityhistory.therecordid order BY id DESC LIMIT 1) ORDER BY id DESC`
var getTrnslaentityhistorytype2 = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, recorddatetime, recorddatetoint, donotupdatesladue, recordtimetoint, mstslastateid, commentonrecord, slastartstopindicator, fromclientuserid FROM trnslaentityhistory WHERE clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 AND slastartstopindicator = 2 AND id <= ? AND id > (SELECT id FROM trnslaentityhistory WHERE id > ? and slastartstopindicator=1 ORDER BY id DESC LIMIT 1) ORDER BY id DESC`

var checkTrnslaentityhistorytype = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, recorddatetime,
 recorddatetoint, donotupdatesladue, recordtimetoint, mstslastateid, commentonrecord, slastartstopindicator,
  fromclientuserid,pushtime from trnslaentityhistory where  clientid=? and mstorgnhirarchyid=? and therecordid=? and 
  id>=(select max(id) from trnslaentityhistory where therecordid= ? and slastartstopindicator=2 and id>(select max(id) 
  from trnslaentityhistory where therecordid=? and slastartstopindicator=1));`

var updatePushTimeInHistory = `UPDATE trnslaentityhistory SET pushtime = ? WHERE id = ?`

//var getTrnslaentityhistorytype2 = `SELECT id, clientid, mstorgnhirarchyid, mstslaentityid, therecordid, recorddatetime, recorddatetoint, donotupdatesladue, recordtimetoint, mstslastateid, commentonrecord, slastartstopindicator, fromclientuserid FROM trnslaentityhistory WHERE clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 AND slastartstopindicator = 2 AND id <= ?   ORDER BY id DESC`

var getrecordFromsladue = `SELECT id, clientid, mstorgnhirarchyid, therecordid FROM mstsladue WHERE therecordid = ?`

var updateExecuteFlag = `UPDATE mstsladue SET isexecute = ? WHERE id IN(%s)`

var getCategoryIds = `SELECT
a.clientid,
a.mstorgnhirarchyid,
 a.recordid,
 a.ticket_type_id,
 c.priority_id,
 d.Working_category
	  
FROM
 (SELECT
	 clientid,
		 mstorgnhirarchyid,
		 recordid,
		 recorddiffid AS ticket_type_id
 FROM
	 maprecordtorecorddifferentiation m
 WHERE
	 recorddifftypeid = 2 AND recordid = ?
		 AND islatest = 1) a,
 (SELECT
		n.clientid,
		 n.mstorgnhirarchyid,
		 n.recordid,
		 n.recorddiffid AS status_id
		
 FROM
	 maprecordtorecorddifferentiation n
 WHERE
		 n.recorddifftypeid = 3
		 AND n.recordid = ?
		 AND n.islatest = 1) b,
 (SELECT
	 o.clientid,
		 o.mstorgnhirarchyid,
		 o.recordid,
		 o.recorddiffid AS priority_id
 FROM
	 maprecordtorecorddifferentiation o
 WHERE
		 o.recorddifftypeid = 5
		 AND o.recordid = ?
		 AND o.islatest = 1) c,
(SELECT
	 clientid,
		 mstorgnhirarchyid,
		 recordid,
		 recorddiffid AS Working_category
 FROM
	 maprecordtorecorddifferentiation m
 WHERE
	 isworking = 1 AND recordid = ?
		 AND islatest = 1) d
  

WHERE
	 a.recordid = b.recordid
	 AND a.recordid = c.recordid
	 AND a.recordid = d.recordid
	 AND a.clientid = ?
	 AND a.mstorgnhirarchyid = ?
	 AND a.recordid = ?`

func (dbc DbConn) GetSpecificMstslafullfillmentcriteriaWithoutWCat(page *entities.MstslafullfillmentcriteriaEntity) (entities.MstslafullfillmentcriteriaEntity, error) {
	value := entities.MstslafullfillmentcriteriaEntity{}
	rows, err := dbc.DB.Query(getMstslafullfillmentcriteriawithoutworkingcategory, page.Clientid, page.Mstorgnhirarchyid, page.Mstrecorddifferentiationtickettypeid,
		page.Mstrecorddifferentiationpriorityid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return value, err
	}
	for rows.Next() {
		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Slaid,
			&value.Mstrecorddifferentiationtickettypeid, &value.Mstrecorddifferentiationpriorityid,
			&value.Mstrecorddifferentiationworkingcatid, &value.Responsetimeinhr, &value.Responsetimeinmin,
			&value.Responsetimeinsec, &value.Resolutiontimeinhr, &value.Resolutiontimeinmin, &value.Resolutiontimeinsec,
			&value.Supportgroupspecific, &value.Activeflg)
		if err != nil {
			logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		}
		// p("dao 4444")
	}
	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", value)
	return value, nil
}

func (dbc DbConn) UpdateRessolutionEndFlag(clientid int64, mstorgnhirarchyid int64, therecordid int64, resolutionCompleteTime string) (bool, error) {
	p := logger.Log.Println

	rows, err := dbc.DB.Query(updateRessolutionEndFlag, resolutionCompleteTime, therecordid, clientid, mstorgnhirarchyid)
	defer rows.Close()
	p(rows)
	if err != nil {
		logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return false, err
	}
	// p("dao 2222")
	return true, nil
}

func (dbc DbConn) UpdateViolateFlag(clientid int64, mstorgnhirarchyid int64, therecordid int64, flag int) (bool, error) {
	p := logger.Log.Println
	sql := ""
	if flag == 0 {
		sql = updateRespViolateFlag
	} else {
		sql = updateResolViolateFlag
	}
	rows, err := dbc.DB.Query(sql, therecordid, clientid, mstorgnhirarchyid)
	defer rows.Close()
	p(rows)
	if err != nil {
		logger.Log.Println("UpdateViolateFlag ------------>", err)
		return false, err
	}
	return true, nil
}

func (dbc DbConn) GetMstsladue(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.MstsladueEntity, error) {
	value := entities.MstsladueEntity{}
	rows, err := dbc.DB.Query(getMstsladue, therecordid, clientid, mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
		return value, err
	}
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Mstslaentityid,
			&value.Therecordid, &value.Latestone,
			&value.Startdatetimeresponse, &value.Startdatetimeresolution, &value.Duedatetimeresponse,
			&value.Duedatetimeresolution, &value.DuedatetimeresolutionInt, &value.Duedatetimeresponseint, &value.Resoltiondone, &value.Resolutiondatetime,
			&value.Trnslaentityhistoryid, &value.Remainingtime, &value.Completepercent, &value.Responseremainingtime, &value.Responsepercentage, &value.Isresponsecomplete, &value.Isresolutioncomplete, &value.PushTime)
		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
	}
	return value, nil
}

func (dbc DbConn) Getutcdiff(clientid int64, mstorgnhirarchyid int64) (entities.ZoneEntity, error) {
	// p := fmt.Println

	value := entities.ZoneEntity{}
	// p("Client ID ")

	rows, err := dbc.DB.Query(getutcdiff, clientid, mstorgnhirarchyid)
	defer rows.Close()
	// p("dao 1111")
	// p(rows)
	if err != nil {
		//logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return value, err
	}
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan(&value.UTCdiff)
		// p("dao 3333")
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
		// p("dao 4444")
	}
	return value, nil
}

func (dbc DbConn) UpdateResponseEndFlag(clientid int64, mstorgnhirarchyid int64, therecordid int64, responseCompleteTime string) (bool, error) {
	p := fmt.Println

	rows, err := dbc.DB.Query(updateResponseEndFlag, responseCompleteTime, therecordid, clientid, mstorgnhirarchyid)
	defer rows.Close()
	// p("dao 1111")
	p(rows)
	if err != nil {
		//logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return false, err
	}
	// p("dao 2222")
	return true, nil
}

func (dbc DbConn) UpdateRemainingPercent(clientid int64, mstorgnhirarchyid int64, therecordid int64, remainingtime int64, completepercent float64, responseremainingtime int64, responsepercentage float64) (bool, error) {
	logger.Log.Println("UpdateRemainingPercent query is ------------->", updateRemainingPercent)
	logger.Log.Println("UpdateRemainingPercent parameter is ------------->", remainingtime, completepercent, responseremainingtime, responsepercentage, therecordid, clientid, mstorgnhirarchyid)
	rows, err := dbc.DB.Query(updateRemainingPercent, remainingtime, completepercent, responseremainingtime, responsepercentage, therecordid, clientid, mstorgnhirarchyid)
	defer rows.Close()
	// p("dao 1111")
	//p(rows)
	if err != nil {
		//logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return false, err
	}
	// p("dao 2222")
	return true, nil
}

func (dbc DbConn) UpdateMstsladue(clientid int64, mstorgnhirarchyid int64, therecordid int64, duedatetimeresponse string, duedatetimeresolution string, duedatetimeresolutionint int64, duedatetimeresponseint int64, totalpushTime int64) (bool, error) {
	p := fmt.Println
	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.", duedatetimeresolutionint)
	rows, err := dbc.DB.Query(updateMstsladue, duedatetimeresponse, duedatetimeresolution, duedatetimeresolutionint, duedatetimeresponseint, totalpushTime, therecordid, clientid, mstorgnhirarchyid)
	defer rows.Close()
	// p("dao 1111")
	p(rows)
	if err != nil {
		logger.Log.Println("UpdateMstsladue Get Statement Prepare Error", err)
		return false, err
	}
	// p("dao 2222")
	return true, nil
}

// func (dbc DbConn) GetMstsladue(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.MstsladueEntity, error) {
// 	value := entities.MstsladueEntity{}
// 	rows, err := dbc.DB.Query(getMstsladue, therecordid, clientid, mstorgnhirarchyid)
// 	defer rows.Close()
// 	if err != nil {
// 		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
// 		return value, err
// 	}
// 	// p("dao 2222")

// 	for rows.Next() {
// 		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Mstslaentityid,
// 			&value.Therecordid, &value.Latestone,
// 			&value.Startdatetimeresponse, &value.Startdatetimeresolution, &value.Duedatetimeresponse,
// 			&value.Duedatetimeresolution, &value.DuedatetimeresolutionInt, &value.Duedatetimeresponseint, &value.Resoltiondone, &value.Resolutiondatetime,
// 			&value.Trnslaentityhistoryid, &value.Remainingtime, &value.Completepercent, &value.Responseremainingtime, &value.Responsepercentage, &value.Isresponsecomplete)
// 		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
// 	}
// 	return value, nil
// }

func (dbc DbConn) GetTrnslaentityhistory(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.TrnslaentityhistoryEntity, error) {
	// p := fmt.Println

	value := entities.TrnslaentityhistoryEntity{}
	// p("Client ID ")

	rows, err := dbc.DB.Query(getTrnslaentityhistory, clientid, mstorgnhirarchyid, therecordid)
	defer rows.Close()
	// p("dao 1111")
	// p(rows)
	if err != nil {
		logger.Log.Println("GetTrnslaentityhistory Get Statement Prepare Error", err)
		return value, err
	}
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Mstslaentityid,
			&value.Therecordid, &value.Recorddatetime,
			&value.Recorddatetoint, &value.Donotupdatesladue, &value.Recordtimetoint,
			&value.Mstslastateid, &value.Commentonrecord, &value.Slastartstopindicator, &value.Fromclientuserid)
		// p("dao 3333")
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
		// p("dao 4444")
	}
	return value, nil
}

func (dbc DbConn) InsertTrnslaentityhistory(tz *entities.TrnslaentityhistoryEntity) (int64, error) {
	stmt, err := dbc.DB.Prepare(insertTrnslaentityhistory)
	defer stmt.Close()
	if err != nil {
		logger.Log.Println("history table id----------22222222222222222222222222222222------------>", err)
		return 0, err
	}
	res, err := stmt.Exec(tz.Clientid, tz.Mstorgnhirarchyid, tz.Mstslaentityid, tz.Therecordid, tz.Recorddatetime, tz.Recorddatetoint, tz.Donotupdatesladue, tz.Recordtimetoint, tz.Mstslastateid, tz.Commentonrecord, tz.Slastartstopindicator, tz.Fromclientuserid)
	if err != nil {
		logger.Log.Println("history table id----------22222222222222222222222222222222------------>", err)
		return 0, err
	}
	lastInsertedId, err := res.LastInsertId()
	return lastInsertedId, nil
}

// func (dbc DbConn) InsertMstsladue(tz *entities.MstsladueEntity) (int64, error) {
// 	// fmt.Println(tz.Clientid)
// 	stmt, err := dbc.DB.Prepare(insertMstsladue)
// 	defer stmt.Close()
// 	if err != nil {
// 		return 0, err
// 	}
// 	res, err := stmt.Exec(tz.Clientid, tz.Mstorgnhirarchyid, tz.Mstslaentityid, tz.Therecordid, tz.Latestone, tz.Startdatetimeresponse, tz.Startdatetimeresolution, tz.Duedatetimeresponse, tz.Duedatetimeresolution, tz.Duedatetimetominute, tz.Resoltiondone, tz.Resolutiondatetime, tz.Lastupdatedattime, tz.Trnslaentityhistoryid, tz.DuedatetimeresolutionInt, tz.Duedatetimeresponseint, tz.Remainingtime, tz.Completepercent)
// 	if err != nil {
// 		// fmt.Println("-----------------------------------------")
// 		fmt.Println(err)
// 		return 0, err
// 	}
// 	lastInsertedId, err := res.LastInsertId()
// 	return lastInsertedId, nil
// }

func (dbc DbConn) InsertMstsladue(tz *entities.MstsladueEntity) (int64, error) {
	stmt, err := dbc.DB.Prepare(insertMstsladue)
	defer stmt.Close()
	if err != nil {
		logger.Log.Println("statement error ------------->", err)
		return 0, err
	}
	res, err := stmt.Exec(tz.Clientid, tz.Mstorgnhirarchyid, tz.Mstslaentityid, tz.Therecordid, tz.Latestone, tz.Startdatetimeresponse, tz.Startdatetimeresolution, tz.Duedatetimeresponse, tz.Duedatetimeresolution, tz.Duedatetimetominute, tz.Resoltiondone, tz.Resolutiondatetime, tz.Lastupdatedattime, tz.Trnslaentityhistoryid, tz.DuedatetimeresolutionInt, tz.Duedatetimeresponseint, tz.Remainingtime, tz.Completepercent, tz.Responseremainingtime, tz.Responsepercentage)
	if err != nil {
		// fmt.Println("-----------------------------------------")
		logger.Log.Println(err)
		return 0, err
	}
	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		// fmt.Println("-----------------------------------------")
		logger.Log.Println(err)
		return 0, err
	}
	return lastInsertedId, nil
}

func (dbc DbConn) GetSpecificMstslafullfillmentcriteria(page *entities.MstslafullfillmentcriteriaEntity) (entities.MstslafullfillmentcriteriaEntity, error) {
	value := entities.MstslafullfillmentcriteriaEntity{}
	rows, err := dbc.DB.Query(getMstslafullfillmentcriteria, page.Clientid, page.Mstorgnhirarchyid, page.Mstrecorddifferentiationtickettypeid,
		page.Mstrecorddifferentiationpriorityid, page.Mstrecorddifferentiationworkingcatid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return value, err
	}
	for rows.Next() {
		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Slaid,
			&value.Mstrecorddifferentiationtickettypeid, &value.Mstrecorddifferentiationpriorityid,
			&value.Mstrecorddifferentiationworkingcatid, &value.Responsetimeinhr, &value.Responsetimeinmin,
			&value.Responsetimeinsec, &value.Resolutiontimeinhr, &value.Resolutiontimeinmin, &value.Resolutiontimeinsec,
			&value.Supportgroupspecific, &value.Activeflg)
		logger.Log.Println("GetSpecificMstslafullfillmentcriteria ------------------------>", err)
	}
	return value, nil
}

func (dbc DbConn) GetSupportGroupId(clientid int64, mstorgnhirarchyid int64, mstslafullfillmentcriteriaid int64) (int64, error) {
	p := fmt.Println
	supportGroupId := int64(0)
	// p("clientid", clientid)
	// p("mstorgnhirarchyid", mstorgnhirarchyid)
	// p("mstslafullfillmentcriteriaid", mstslafullfillmentcriteriaid)
	rows, err := dbc.DB.Query(getsupportGroupid, clientid, mstorgnhirarchyid, mstslafullfillmentcriteriaid)
	defer rows.Close()
	// p("dao 1111")
	p(rows)
	if err != nil {
		//logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return supportGroupId, err
	}
	// p("dao 2222")
	// p(">>>>>>", rows)
	for rows.Next() {
		err := rows.Scan(&supportGroupId)
		// p("dao 3333")
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
		// p("dao 4444")
	}
	return supportGroupId, nil
}

func (dbc DbConn) GetSupportGroupHoliday(clientid int64, mstorgnhirarchyid int64, supportGroupId int64, today string) (int64, int64, string, error) {
	p := fmt.Println
	p("************** GetSupportGroupHoliday *******************")
	starttime := int64(0)
	endtime := int64(0)
	var dateofholiday string
	rows, err := dbc.DB.Query(getclientsupportgroupholyday, clientid, mstorgnhirarchyid, supportGroupId, today)
	defer rows.Close()
	if err != nil {
		//logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return 0, 0, "", err
	}
	for rows.Next() {
		err := rows.Scan(&starttime, &endtime, &dateofholiday)
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
	}
	return starttime, endtime, dateofholiday, nil
}

func (dbc DbConn) GetClientHoliday(clientid int64, mstorgnhirarchyid int64, today string) (int64, int64, string, error) {
	p := fmt.Println
	p("************** GetClientHoliday *******************")
	starttime := int64(0)
	endtime := int64(0)
	var dateofholiday string
	rows, err := dbc.DB.Query(getclientholyday, clientid, mstorgnhirarchyid, today)
	defer rows.Close()
	if err != nil {
		//logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return 0, 0, "", err
	}
	for rows.Next() {
		err := rows.Scan(&starttime, &endtime, &dateofholiday)
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
	}
	return starttime, endtime, dateofholiday, nil
}

// This method will return the starttime and endtime of the week day
func (dbc DbConn) GetClientDayOfWeek(clientid int64, mstorgnhirarchyid int64, dayofweekid int64) (int64, int64, error) {
	//p := fmt.Println
	logger.Log.Println("************** GetClientDayOfWeek *******************", dayofweekid)
	starttime := int64(0)
	endtime := int64(0)
	// print(dayofweekid)
	rows, err := dbc.DB.Query(getclientWeekDay, clientid, mstorgnhirarchyid, dayofweekid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("************** GetClientDayOfWeek *******************", err)
		return 0, 0, err
	}
	for rows.Next() {
		err := rows.Scan(&starttime, &endtime)
		if err != nil {
			logger.Log.Println("************** GetClientDayOfWeek *******************", err)
		}
	}
	return starttime, endtime, nil
}

// This method will return the starttime and endtime of the week day for support group
func (dbc DbConn) GetSupportGroupDayOfWeek(clientid int64, mstorgnhirarchyid int64, supportGroupId int64, dayofweekid int64) (int64, int64, int64, error) {
	p := fmt.Println
	p("************** GetSupportGroupDayOfWeek *******************")
	starttime := int64(0)
	endtime := int64(0)
	nextdayforward := int64(0)
	// print(dayofweekid)
	rows, err := dbc.DB.Query(getsupportgroupWeekDay, clientid, mstorgnhirarchyid, supportGroupId, dayofweekid)
	defer rows.Close()
	if err != nil {
		//logger.Log.Println("GetAllMstslafullfillmentcriteria Get Statement Prepare Error", err)
		return 0, 0, 0, err
	}
	for rows.Next() {
		err := rows.Scan(&starttime, &endtime, &nextdayforward)
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
	}
	return starttime, endtime, nextdayforward, nil
}

// func (dbc DbConn) GetTrnslaentityhistorytype2(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.TrnslaentityhistoryEntity, error) {
// 	// p := fmt.Println

// 	value := entities.TrnslaentityhistoryEntity{}
// 	// p("Client ID ")

// 	rows, err := dbc.DB.Query(getTrnslaentityhistorytype2, clientid, mstorgnhirarchyid, therecordid)
// 	defer rows.Close()
// 	// p("dao 1111")
// 	// p(rows)
// 	if err != nil {
// 		logger.Log.Println("GetTrnslaentityhistory Get Statement Prepare Error", err)
// 		return value, err
// 	}
// 	// p("dao 2222")

// 	for rows.Next() {
// 		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Mstslaentityid,
// 			&value.Therecordid, &value.Recorddatetime,
// 			&value.Recorddatetoint, &value.Donotupdatesladue, &value.Recordtimetoint,
// 			&value.Mstslastateid, &value.Commentonrecord, &value.Slastartstopindicator, &value.Fromclientuserid)
// 		// p("dao 3333")
// 		if err != nil {
// 			panic(err) // Error related to the iteration of rows
// 		}
// 		// p("dao 4444")
// 	}
// 	return value, nil
// }


/*func (dbc DbConn) GetTrnslaentityhistorytype2(clientid int64, mstorgnhirarchyid int64, therecordid int64, trnId int64) ([]entities.TrnslaentityhistoryEntity, error) {
	// p := fmt.Println

	value := entities.TrnslaentityhistoryEntity{}
	allValue := []entities.TrnslaentityhistoryEntity{}
	// p("Client ID ")

	rows, err := dbc.DB.Query(getTrnslaentityhistorytype2, clientid, mstorgnhirarchyid, therecordid, trnId)
	defer rows.Close()
	// p("dao 1111")
	// p(rows)
	if err != nil {
		logger.Log.Println("GetTrnslaentityhistory Get Statement Prepare Error", err)
		return allValue, err
	}
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Mstslaentityid,
			&value.Therecordid, &value.Recorddatetime,
			&value.Recorddatetoint, &value.Donotupdatesladue, &value.Recordtimetoint,
			&value.Mstslastateid, &value.Commentonrecord, &value.Slastartstopindicator, &value.Fromclientuserid)
		// p("dao 3333")
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
		allValue = append(allValue, entities.TrnslaentityhistoryEntity{Id: value.Id, Clientid: value.Clientid, Mstorgnhirarchyid: value.Mstorgnhirarchyid, Mstslaentityid: value.Mstslaentityid,
			Therecordid: value.Therecordid, Recorddatetime: value.Recorddatetime,
			Recorddatetoint: value.Recorddatetoint, Donotupdatesladue: value.Donotupdatesladue, Recordtimetoint: value.Recordtimetoint,
			Mstslastateid: value.Mstslastateid, Commentonrecord: value.Commentonrecord, Slastartstopindicator: value.Slastartstopindicator, Fromclientuserid: value.Fromclientuserid})
		// p("dao 4444")
	}
	return allValue, nil
}*/


func (dbc DbConn) GetTrnslaentityhistorytype2(clientid int64, mstorgnhirarchyid int64, therecordid int64, trnId int64) ([]entities.TrnslaentityhistoryEntity, error) {
    // p := fmt.Println

    value := entities.TrnslaentityhistoryEntity{}
    allValue := []entities.TrnslaentityhistoryEntity{}
    // p("Client ID ")

    rows, err := dbc.DB.Query(getTrnslaentityhistorytype2, clientid, mstorgnhirarchyid, therecordid, trnId, trnId)
    defer rows.Close()
    // p("dao 1111")
    // p(rows)
    if err != nil {
        logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetTrnslaentityhistory Get Statement Prepare Error", err)
        return allValue, err
    }
    // p("dao 2222")

    for rows.Next() {
        err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Mstslaentityid,
            &value.Therecordid, &value.Recorddatetime,
            &value.Recorddatetoint, &value.Donotupdatesladue, &value.Recordtimetoint,
            &value.Mstslastateid, &value.Commentonrecord, &value.Slastartstopindicator, &value.Fromclientuserid)
        // p("dao 3333")
        if err != nil {
            // panic(err) // Error related to the iteration of rows
            logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>Scan error in GetTrnslaentityhistorytype2", err)
        }
        allValue = append(allValue, entities.TrnslaentityhistoryEntity{Id: value.Id, Clientid: value.Clientid, Mstorgnhirarchyid: value.Mstorgnhirarchyid, Mstslaentityid: value.Mstslaentityid,
            Therecordid: value.Therecordid, Recorddatetime: value.Recorddatetime,
            Recorddatetoint: value.Recorddatetoint, Donotupdatesladue: value.Donotupdatesladue, Recordtimetoint: value.Recordtimetoint,
            Mstslastateid: value.Mstslastateid, Commentonrecord: value.Commentonrecord, Slastartstopindicator: value.Slastartstopindicator, Fromclientuserid: value.Fromclientuserid})
        // p("dao 4444")
    }
    return allValue, nil
}
func (dbc DbConn) CheckTrnslaentityhistorytype(clientid int64, mstorgnhirarchyid int64, therecordid int64, trnId int64) ([]entities.TrnslaentityhistoryEntity, error) {
	// p := fmt.Println

	// value := entities.TrnslaentityhistoryEntity{}
	allValue := []entities.TrnslaentityhistoryEntity{}
	// p("Client ID ")

	rows, err := dbc.DB.Query(checkTrnslaentityhistorytype, clientid, mstorgnhirarchyid, therecordid, therecordid, therecordid)
	defer rows.Close()
	// p("dao 1111")
	// p(rows)
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetTrnslaentityhistory Get Statement Prepare Error", err)
		return allValue, err
	}
	// p("dao 2222")

	for rows.Next() {
		value := entities.TrnslaentityhistoryEntity{}

		err := rows.Scan(&value.Id, &value.Clientid, &value.Mstorgnhirarchyid, &value.Mstslaentityid,
			&value.Therecordid, &value.Recorddatetime,
			&value.Recorddatetoint, &value.Donotupdatesladue, &value.Recordtimetoint,
			&value.Mstslastateid, &value.Commentonrecord, &value.Slastartstopindicator, &value.Fromclientuserid, &value.PushTime)
		// p("dao 3333")
		if err != nil {
			// panic(err) // Error related to the iteration of rows
			logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>Scan error in GetTrnslaentityhistorytype2", err)
		}
		// allValue = append(allValue, entities.TrnslaentityhistoryEntity{Id: value.Id, Clientid: value.Clientid, Mstorgnhirarchyid: value.Mstorgnhirarchyid, Mstslaentityid: value.Mstslaentityid,
		// 	Therecordid: value.Therecordid, Recorddatetime: value.Recorddatetime,
		// 	Recorddatetoint: value.Recorddatetoint, Donotupdatesladue: value.Donotupdatesladue, Recordtimetoint: value.Recordtimetoint,
		// 	Mstslastateid: value.Mstslastateid, Commentonrecord: value.Commentonrecord, Slastartstopindicator: value.Slastartstopindicator, Fromclientuserid: value.Fromclientuserid})
		// p("dao 4444")
		allValue = append(allValue, value)
	}
	return allValue, nil
}
func (dbc DbConn) UpdatePushTimeInHistory(historyId int64, pushtime int64) (bool, error) {
	p := fmt.Println
	rows, err := dbc.DB.Query(updatePushTimeInHistory, pushtime, historyId)
	defer rows.Close()
	p(rows)
	if err != nil {
		return false, err
	}
	return true, nil
}

//var getTrnslaentityhistoryLastPushTime = `SELECT recorddatetime, recorddatetoint FROM trnslaentityhistory WHERE clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 AND slastartstopindicator = 2 ORDER BY id DESC LIMIT 1`
var getTrnslaentityhistoryLastPushTime = "SELECT recorddatetime, recorddatetoint FROM trnslaentityhistory WHERE clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 AND slastartstopindicator = 2 AND id > (SELECT id FROM trnslaentityhistory WHERE slastartstopindicator IN (1, 3) AND  clientid = ? AND mstorgnhirarchyid = ? AND therecordid = ? AND activeflg = 1 AND deleteflg = 0 ORDER BY id DESC LIMIT 1) LIMIT 1"

func (dbc DbConn) GetTrnslaentityhistoryLastPushTime(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.TrnslaentityhistoryLastPushEntity, error) {
	// p := fmt.Println

	value := entities.TrnslaentityhistoryLastPushEntity{}
	// p("Client ID ")

	rows, err := dbc.DB.Query(getTrnslaentityhistoryLastPushTime, clientid, mstorgnhirarchyid, therecordid, clientid, mstorgnhirarchyid, therecordid)
	// p("dao 1111")
	// p(rows)
	if err != nil {
		logger.Log.Println("GetTrnslaentityhistoryLastPushTime Get Statement Prepare Error", err)
		return value, err
	}
	defer rows.Close()
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan(&value.Recorddatetime,
			&value.Recorddatetoint)
		// p("dao 3333")
		if err != nil {
			panic(err) // Error related to the iteration of rows
		}
		// p("dao 4444")
	}
	return value, nil
}

func (dbc DbConn) GetrecordFromsladue(recordId int64) ([]entities.SLASchedulerEntity, error) {
	value := entities.SLASchedulerEntity{}
	allValue := []entities.SLASchedulerEntity{}
	rows, err := dbc.DB.Query(getrecordFromsladue, recordId)
	defer rows.Close()
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetrecordFromsladue Get Statement Prepare Error", err)
		return allValue, err
	}
	for rows.Next() {
		err := rows.Scan(&value.ID, &value.ClientID, &value.Mstorgnhirarchyid, &value.RecordID)
		if err != nil {

			logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetCategoryIds Get Statement Prepare Error", err)
		}

		allValue = append(allValue, entities.SLASchedulerEntity{ID: value.ID, ClientID: value.ClientID, Mstorgnhirarchyid: value.Mstorgnhirarchyid, RecordID: value.RecordID})
	}
	return allValue, nil
}
func (dbc DbConn) UpdateExecuteFlag(ids []int64, flag int64) (bool, error) {
	// p := logger.Log.Println
	stringedIDs := fmt.Sprintf("%v", ids)
	stringedIDs = stringedIDs[1 : len(stringedIDs)-1]
	stringedIDs = strings.ReplaceAll(stringedIDs, " ", ",")
	if stringedIDs != "" {
		sqlRaw := fmt.Sprintf(updateExecuteFlag, stringedIDs)
		// fmt.Println(sqlRaw)
		rows, err := dbc.DB.Query(sqlRaw, flag)
		defer rows.Close()
		// p(rows)
		if err != nil {
			logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>UpdateExecuteFlag Get Statement Prepare Error", err)
			return false, err
		}
	}
	return true, nil
}

func (dbc DbConn) GetCategoryIds(clientid int64, mstorgnhirarchyid int64, therecordid int64) (entities.SLATabEntity, error) {
	value := entities.SLATabEntity{}
	// fmt.Println(getCategoryIds)
	rows, err := dbc.DB.Query(getCategoryIds, therecordid, therecordid, therecordid, therecordid, clientid, mstorgnhirarchyid, therecordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetMstsladue Get Statement Prepare Error", err)
		return value, err
	}
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan(&value.ClientID, &value.Mstorgnhirarchyid, &value.RecordID, &value.RecordtypeID, &value.PriorityID, &value.WorkingcatID)
		if err != nil {

			logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetCategoryIds Get Statement Prepare Error", err)
		}
	}
	return value, nil
}

/*func (dbc DbConn) UpdateRecordFullFillDetails(details entities.SLARcordFullfillUpdate) (bool, error) {
	var sql = `UPDATE recordfulldetails SET responseslameterpercentage = ?, resolutionslameterpercentage = ?, businessaging = ?, calendaraging = ?, actualeffort = ?, slaidletime = ?, respoverduetime = ?, respoverdueperc = ?, resooverduetime = ?, resooverdueperc = ? WHERE clientid = ? AND mstorgnhirarchyid = ? AND recordid = ?`
	stmt, err := dbc.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>Updaterecordfulldetails Prepare Statement  Error", err)
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(details.Responseslameterpercentage, details.Resolutionslameterpercentage, details.Businessaging, details.Calendaraging, details.Actualeffort, details.Slaidletime, details.Respoverduetime, details.Respoverdueperc, details.Resooverduetime, details.Resooverdueperc, details.ClientID, details.Mstorgnhirarchyid, details.RecordID)
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>Updaterecordfulldetails Execute Statement  Error", err)
		return false, err
	}
	return true, nil
}*/



func (dbc DbConn) UpdateRecordFullFillDetails(details entities.SLARcordFullfillUpdate) (bool, error) {
    // fmt.Println("UpdateRecordFullFillDetails >>>>>>>>>>>>>> ", details)
    // // p := logger.Log.Println
    logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>Updaterecordfulldetails >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>",details.RecordID,details.Responseslameterpercentage, details.Resolutionslameterpercentage, details.Businessaging, details.Calendaraging)

    var sql = `UPDATE recordfulldetails SET responseslameterpercentage = ?, resolutionslameterpercentage = ?, businessaging = (? - reopentimetaken) , calendaraging = (? - reopentimetakencalendar) , actualeffort = ?, slaidletime = ? , respoverduetime = ?, respoverdueperc = ?, resooverduetime = ?, resooverdueperc = ? WHERE clientid = ? AND mstorgnhirarchyid = ? AND recordid = ?`
    logger.Log.Println(sql)
    // // fmt.Println(details.Responseslameterpercentage, details.Resolutionslameterpercentage, details.Businessaging, details.Calendaraging, details.Actualeffort, details.Slaidletime, details.Respoverduetime, details.Respoverdueperc, details.Resooverduetime, details.Resooverdueperc, details.ClientID, details.Mstorgnhirarchyid)
    // rows, err := dbc.DB.Query(sql, details.Responseslameterpercentage, details.Resolutionslameterpercentage, details.Businessaging, details.Calendaraging, details.Actualeffort, details.Slaidletime, details.Respoverduetime, details.Respoverdueperc, details.Resooverduetime, details.Resooverdueperc, details.ClientID, details.Mstorgnhirarchyid, details.RecordID)
    // defer rows.Close()
    // // p(rows)
    // if err != nil {
    //  logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>recordfulldetails Update Statement Prepare Error", err)
    //  return false, err
    // }
    // // p("dao 2222")
    stmt, err := dbc.DB.Prepare(sql)
    if err != nil {
        logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>Updaterecordfulldetails Prepare Statement  Error", err)
        return false, err
    }
    defer stmt.Close()
    _, err = stmt.Exec(details.Responseslameterpercentage, details.Resolutionslameterpercentage, details.Businessaging, details.Calendaraging, details.Actualeffort, details.Slaidletime, details.Respoverduetime, details.Respoverdueperc, details.Resooverduetime, details.Resooverdueperc, details.ClientID, details.Mstorgnhirarchyid, details.RecordID)
    if err != nil {
        logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>Updaterecordfulldetails Execute Statement  Error", err)
        return false, err
    }
    return true, nil
}



var getMstsladueFirst = `SELECT startdatetimeresolution,pushtime FROM mstsladue WHERE therecordid = ? AND clientid = ? AND mstorgnhirarchyid = ? AND activeflg = 1 AND deleteflg = 0 ORDER BY id ASC LIMIT 1`

func (dbc DbConn) GetFirstStartTimeResolution(clientid int64, mstorgnhirarchyid int64, therecordid int64) (string,int64, error) {
    rows, err := dbc.DB.Query(getMstsladueFirst, therecordid, clientid, mstorgnhirarchyid)
    var Startdatetimeresolution = ""
    var pushtime int64
    defer rows.Close()
    if err != nil {
        logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetMstsladue Get Statement Prepare Error", err)
        return Startdatetimeresolution, pushtime,err
    }
    // p("dao 2222")

    for rows.Next() {
        err := rows.Scan(&Startdatetimeresolution,&pushtime)
        if err != nil {

            logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>GetMstsladue Get Statement Prepare Error", err)
        }
    }
    return Startdatetimeresolution, pushtime,nil
}





//=========================================To get last grpid=====================================================

func (dbc DbConn) GetLastSupportgrp(requestid int64,supportgroupid int64) (int64, error) {
	var value int64=0
	sql:="select b.id from mstrequesthistory c left join mstclientuser a on c.mstuserid=a.id and a.activeflag=1 and a.deleteflag=0,mstsupportgrp b where c.mainrequestid =? and c.mstgroupid not in (?) and c.mstgroupid=b.id  and c.activeflg=1 and c.deleteflg=0 and b.activeflg=1 and b.deleteflg=0  order by c.id desc limit 1;	"
	rows, err := dbc.DB.Query(sql, requestid,supportgroupid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
		return value, err
	}
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan()
		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
	}
	return value, nil
}
func (dbc DbConn) Getrequestid(therecordid int64) (int64, error) {
	var value int64 = 0
	sql := "SELECT d.id FROM mstrequest d,maprequestorecord b WHERE b.recordid=? AND d.id=b.mstrequestid AND b.activeflg=1 AND b.deleteflg=0 AND d.activeflg=1 AND d.deleteflg=0 order by b.id desc limit 1;"
	rows, err := dbc.DB.Query(sql, therecordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
		return value, err
	}
	// p("dao 2222")

	for rows.Next() {
		err := rows.Scan(&value)
		logger.Log.Println("GetMstsladue Get Statement Prepare Error", err)
	}
	return value, nil
}
//=========================================To get last grpid=====================================================

