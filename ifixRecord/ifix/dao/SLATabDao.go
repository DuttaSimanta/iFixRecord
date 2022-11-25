package dao

import (
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

var resolution = "select a.recordid,b.status,c.priority, e.resolutiontimeinsec ,g.name slaresolutionmeterstatus , coalesce((select DATE_FORMAT(duedatetimeresolution,'%d-%m-%Y %T') from mstsladue where therecordid = ? order by id desc limit 1),'') duedatetimeresolution,coalesce((SELECT IF(isresolutionviolation=1,'Yes','NO') FROM mstsladue where therecordid = ? order by id desc limit 1),'') Resolutionslaviolated from (SELECT clientid, mstorgnhirarchyid, recordid, recorddiffid as ticket_type_id FROM maprecordtorecorddifferentiation m where recorddifftypeid=2 AND recordid=? AND islatest=1) a, (SELECT n.clientid, n.mstorgnhirarchyid, n.recordid, n.recorddiffid as status_id,q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=? AND n.islatest=1) b, (SELECT o.clientid, o.mstorgnhirarchyid, o.recordid, o.recorddiffid as priority_id,p.name priority FROM maprecordtorecorddifferentiation o, mstrecorddifferentiation p where o.recorddiffid = p.id AND o.recorddifftypeid=5 AND o.recordid=? AND o.islatest=1) c , mstslafullfillmentcriteria e  , trnslaentityhistory f, mstslaclockdisplayname g where a.recordid = b.recordid and a.recordid = c.recordid and priority_id = e.mstrecorddifferentiationpriorityid and ticket_type_id = e.mstrecorddifferentiationtickettypeid and e.deleteflg = 0 and e.activeflg = 1 and e.mstrecorddifferentiationworkingcatid=0 and a.recordid = f.therecordid and a.clientid=? and a.mstorgnhirarchyid=? and a.recordid  = ? and f.id = (select max(id) from trnslaentityhistory where therecordid = ?) and f.slastartstopindicator = g.seqno and g.metertypeid=2 and g.activeflg=1 AND g.deleteflg=0 "
var response = "select distinct m.recordid,m.status,n.priority,coalesce(DATE_FORMAT(d.duedatetimeresponse,'%d-%m-%Y %T'),'') duedatetimeresponse , e.responsetimeinsec, g.name slaresponsemeterstatus,coalesce((SELECT IF(isresponseviolation=1,'Yes','NO') FROM mstsladue where therecordid = ? order by id desc limit 1),'') Responseslaviolated from (select a.recordid, a.recorddifftypeid,a.recorddiffid statusid,b.name status,a.clientid,a.mstorgnhirarchyid from maprecordtorecorddifferentiation a,mstrecorddifferentiation b where a.recordid=? and a.recorddifftypeid=3 and a.recorddiffid = b.id order by a.id desc limit 1) m, (select a.recordid, a.recorddifftypeid,a.recorddiffid priorityid,c.name priority,a.clientid,a.mstorgnhirarchyid from maprecordtorecorddifferentiation a, mstrecorddifferentiation c where a.recordid=? and a.recorddifftypeid=5 and a.recorddiffid = c.id  order by a.id desc limit 1) n, mstsladue d , mstslafullfillmentcriteria e,mstslaclockdisplayname g where m.recordid = n.recordid and m.clientid=? and m.mstorgnhirarchyid=? and m.recordid =d.therecordid and n.priorityid = e.mstrecorddifferentiationpriorityid and e.mstrecorddifferentiationtickettypeid = (select max(recorddiffid) from maprecordtorecorddifferentiation where recordid = ? and recorddifftypeid = 2) AND d.isresponsecomplete = g.seqno and g.metertypeid=1 and g.activeflg=1 AND g.deleteflg=0 and e.activeflg=1 AND e.deleteflg=0 order by d.id desc limit 1"

//var holiday = "select from_unixtime(dateofholiday,'%Y-%m-%d') as holidate,starttime,endtime from mstclientholiday where (dateofholiday+endtimeinteger) between (select createdatetime from trnrecord where clientid=? AND mstorgnhirarchyid=? AND id=?) AND (select duedatetimeresolutionint from mstsladue where clientid=? AND mstorgnhirarchyid=? AND therecordid=?)"
var holiday = "select dateofholiday,starttime,endtime from mstclientholiday where (dateofholiday+endtimeinteger) between (select createdatetime from trnrecord where clientid=? AND mstorgnhirarchyid=? AND id=?) AND (select duedatetimeresolutionint from mstsladue where clientid=? AND mstorgnhirarchyid=? AND therecordid=?)"

func (mdao DbConn) GetResolutiondetails(page *entities.SLATabEntity,resultTime int64) (entities.SLAResolutionmeterEntity, error) {
	logger.Log.Println("GetResolutiondetails query -->", resolution)
	logger.Log.Println("GetResolutiondetails parameters -->", page.RecordID, page.RecordID, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.RecordID)
	value := entities.SLAResolutionmeterEntity{}
	rows, err := mdao.DB.Query(resolution, page.RecordID, page.RecordID, page.RecordID, page.RecordID, page.RecordID, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetResolutiondetails Get Statement Prepare Error", err)
		return value, err
	}

	logger.Log.Println("Result time value is ---------------------       111111111111111111                   -->", resultTime)
	for rows.Next() {
		err = rows.Scan(&value.RecordID, &value.Status, &value.Priority, &value.Resolutiontime, &value.Resolutionclockstatus, &value.Resolutionduetime, &value.Resolutionslaviolated)
		logger.Log.Println("GetResolutiondetails error -->", err)
		if resultTime == 0 || resultTime < 0 {
			if value.Resolutionclockstatus != "Stopped" {
				value.Resolutionclockstatus = "Paused"
			}

		}
	}
	return value, nil
}

func (mdao DbConn) GetResponsedetails(page *entities.SLATabEntity,resultTime int64) (entities.SLAResponsemeterEntity, error) {
	logger.Log.Println("GetResponsedetails query -->", response)
	logger.Log.Println("GetResponsedetails parameters -->", page.RecordID, page.RecordID, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.RecordID,resultTime)

	value := entities.SLAResponsemeterEntity{}
	rows, err := mdao.DB.Query(response, page.RecordID, page.RecordID, page.RecordID, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetResolutiondetails Get Statement Prepare Error", err)
		return value, err
	}
logger.Log.Println("Result time value is ---------------------       111111111111111111                   -->", resultTime)
	for rows.Next() {
		err = rows.Scan(&value.RecordID, &value.Status, &value.Priority, &value.Responseduetime, &value.Responsetime, &value.Responseclockstatus, &value.Responseslaviolated)
		logger.Log.Println("GetResponsedetails error -->", err)
		if resultTime == 0 || resultTime < 0 {
			if value.Responseclockstatus != "Stopped" {
				value.Responseclockstatus = "Paused"
			}

		}
	}
	return value, nil
}

func (mdao DbConn) GetHolidaydetails(page *entities.SLATabEntity) ([]entities.SLAHolidayEntity, error) {
	logger.Log.Println("GetHolidaydetails query -->", holiday)
	logger.Log.Println("GetHolidaydetails parameters -->", page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)

	values := []entities.SLAHolidayEntity{}
	rows, err := mdao.DB.Query(holiday, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetHolidaydetails Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.SLAHolidayEntity{}
		err = rows.Scan(&value.Holiday, &value.Starttime, &value.Endtime)
		logger.Log.Println("GetResponsedetails error -->", err)
		values = append(values, value)
	}
	return values, nil
}



func (mdao DbConn) GetCurrentSupportGRP(RecordID int64) (int64, error) {
	logger.Log.Println("In side GetCurrentSupportGRP")
	var sql = "SELECT mstgroupid FROM mstrequesthistory where mainrequestid in (SELECT mstrequestid FROM maprequestorecord where recordid=?) order by id desc limit 1"
	var grpID int64
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetCurrentSupportGRP Get Statement Prepare Error", err)
		return grpID, err
	}
	for rows.Next() {
		rows.Scan(&grpID)

	}
	return grpID, nil
}
