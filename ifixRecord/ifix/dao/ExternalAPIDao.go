package dao

import (
	"errors"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
)

func (dbc DbConn) UpdateIsLatest(recordid int64) error {
	logger.Log.Println("In side UpdateMststateterm")
	sql := "UPDATE externalapitransactionlog SET islatest = 0 where recordid=? and deleteflg=0"
	stmt, err := dbc.DB.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		logger.Log.Println("UpdateMststateterm Prepare Statement  Error", err)
		return err
	}
	_, err = stmt.Exec(recordid)
	if err != nil {
		logger.Log.Println("UpdateMststateterm Execute Statement  Error", err)
		return err
	}
	return nil
}

func (mdao DbConn) InsertExternalApiLog(clintid int64, orgid int64, userid int64, recordid int64, recorddiffid int64, statusSeq int64, requestjson string, responsejson string) error {
	logger.Log.Println("In side InsertExternalApiLog")
	var sql = "insert into externalapitransactionlog( clientid, mstorgnhirarchyid, userid, recordid, recorddifftypeid, recorddiffid, statusseq, requestjson)values(?,?,?,?,?,?,?,?)"
	stmt, err := mdao.DB.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(clintid, orgid, userid, recordid, 2, recorddiffid, statusSeq, requestjson)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}
func (mdao DbConn) GetRecordDiffId(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side GetRecordDiffId")
	var recorddiffid int64
	var sql = "select recorddiffid from iFIX.maprecordtorecorddifferentiation where clientid = ? and mstorgnhirarchyid = ? and recordid = ? and recorddifftypeid = 2 and islatest = 1 and deleteflg = 0"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordDiffId Get Statement Prepare Error", err)
		return recorddiffid, err
	}
	for rows.Next() {

		err = rows.Scan(&recorddiffid)
		if err != nil {
			logger.Log.Println("GetRecordDiffId rows.next() Error", err)
		}
	}
	//logger.Log.Println("map value is ------->", v)
	return recorddiffid, nil
}

func (mdao DbConn) GetRecordId(Recordid string) (int64, error) {
	logger.Log.Println("In side GetRecordId")
	var recordid int64
	var sql = "SELECT recordid FROM recordfulldetails where ticketid = ?;"
	rows, err := mdao.DB.Query(sql, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordId Get Statement Prepare Error", err)
		return recordid, err
	}
	for rows.Next() {
		err = rows.Scan(&recordid)
		logger.Log.Println("GetRecordId rows.next() Error", err)
	}
	return recordid, nil
}

func (mdao DbConn) InsertAttachmentToRecord(ClientID int64, orgnID int64, Recordid int64, Recordstageid int64, Recordtermid int64, OriginalFileName string, NewFileName string, Createdusrid int64, CreatedgroupID int64) (string, error) {
	logger.Log.Println("In side InsertAttachmentToRecord")
	var Statusname string
	var sql = "insert into trnreordtracking (clientid,mstorgnhirarchyid,recordid,recordstageid,recordtermid,recordtrackvalue, recordtrackdescription, createdbyid,createdgrpid,createddate) values (?,?,?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now()))) ;"
	logger.Log.Println(ClientID, orgnID, Recordid, Recordstageid, Recordtermid, OriginalFileName, NewFileName, Createdusrid, CreatedgroupID)

	rows, err := mdao.DB.Query(sql, ClientID, orgnID, Recordid, Recordstageid, Recordtermid, OriginalFileName, NewFileName, Createdusrid, CreatedgroupID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("InsertAttachmentToRecord Statement Prepare Error", err)
		return Statusname, err
	}
	for rows.Next() {
		err = rows.Scan(&Statusname)
		logger.Log.Println("GetReopencount rows.next() Error", err)
	}
	return Statusname, nil
}

/*
func (mdao DbConn) GetRecordId(Recordid string) (int64, error) {
    logger.Log.Println("In side GetRecordId")
    var recordid int64
    var sql = "SELECT recordid FROM iFIX.recordfulldetails where ticketid =  ?;"
    rows, err := mdao.DB.Query(sql, Recordid)
    defer rows.Close()
    if err != nil {
        logger.Log.Println("GetRecordId Get Statement Prepare Error", err)
        return recordid, err
    }
    for rows.Next() {
        err = rows.Scan(&recordid)
        logger.Log.Println("GetRecordId rows.next() Error", err)
    }
    return recordid, nil
}*/
func (mdao DbConn) GetTermIdToRecord(clientID int64, orgnID int64) (int64, error) {
	logger.Log.Println("In side GetTermIdToRecord")
	var Recordtermid int64
	var sql = "SELECT id FROM mstrecordterms where clientid= ? AND mstorgnhirarchyid = ? AND seq = 1;"
	rows, err := mdao.DB.Query(sql, clientID, orgnID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermIdToRecord Statement Prepare Error", err)
		return Recordtermid, err
	}
	for rows.Next() {
		err = rows.Scan(&Recordtermid)
		logger.Log.Println("GetReopencount rows.next() Error", err)
	}
	return Recordtermid, nil
}

//clientid,mstorgnhirarchyid,recordid,recordstageid,recordtermid,recordtrackvalue, recordtrackdescription
func (mdao DbConn) GetStatusName(Recordid string) (string, error) {
	logger.Log.Println("In side GetStatusName")
	var Statusname string
	var sql = "SELECT status FROM recordfulldetails where ticketid = ?"
	rows, err := mdao.DB.Query(sql, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetStatusName Get Statement Prepare Error", err)
		return Statusname, err
	}
	for rows.Next() {
		err = rows.Scan(&Statusname)
		logger.Log.Println("GetReopencount rows.next() Error", err)
	}
	return Statusname, nil
}

func (mdao DbConn) GetRecordstageid(Recordid int64) (int64, error) {
	logger.Log.Println("In side GetRecordstageid")
	var Recordstageid int64
	var sql = "SELECT Recordstageid FROM trnrecord where id = ?;"
	rows, err := mdao.DB.Query(sql, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordstageid Statement Prepare Error", err)
		return Recordstageid, err
	}
	for rows.Next() {
		err = rows.Scan(&Recordstageid)
		logger.Log.Println("GetReopencount rows.next() Error", err)
	}
	return Recordstageid, nil
}

/*
//clientid,mstorgnhirarchyid,recordid,recordstageid,recordtermid,recordtrackvalue, recordtrackdescription
func (mdao DbConn) InsertAttachmentToRecord(ClientID int64, orgnID int64, Recordid int64, Recordstageid int64, Recordtermid int64, OriginalFileName string, NewFileName string) (string, error) {
    logger.Log.Println("In side InsertAttachmentToRecord")
    var Statusname string
    var sql = "insert into trnreordtracking (clientid,mstorgnhirarchyid,recordid,recordstageid,recordtermid,recordtrackvalue, recordtrackdescription) values (?,?,?,?,?,?,?) ;"
    logger.Log.Println(ClientID, orgnID, Recordid, Recordstageid, Recordtermid, OriginalFileName, NewFileName)

    rows, err := mdao.DB.Query(sql, ClientID, orgnID, Recordid, Recordstageid, Recordtermid, OriginalFileName, NewFileName)
    defer rows.Close()
    if err != nil {
        logger.Log.Println("InsertAttachmentToRecord Statement Prepare Error", err)
        return Statusname, err
    }
    for rows.Next() {
        err = rows.Scan(&Statusname)
        logger.Log.Println("GetReopencount rows.next() Error", err)
    }
    return Statusname, nil
}*/
//=========================================================
func (dbc DbConn) GetExternalRecordDetailsByNo(clientID int64, orgnID int64, req *entities.RecordDetailsRequestEntityAPI) (entities.RecordDetailsEntityAPI, int64, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	v := entities.RecordDetailsEntityAPI{}

	var recordtypeid int64
	var sql = "SELECT distinct COALESCE(a.requestername,'') requestername,COALESCE(a.requesteremail,'') requesteremail,COALESCE(a.requestermobile,'') requestermobile,COALESCE(a.requesterlocation,'') requesterlocation,e.vipuser,a.source,(SELECT name FROM mstclient WHERE id=a.clientid) AS Clientname,(SELECT name FROM mstorgnhierarchy WHERE id=a.mstorgnhirarchyid) As Mstorgnhierarchyname,a.recordtitle,a.code,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1),'') priority,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1),'') status,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=7 AND n.recordid=a.id AND n.islatest=1),'') impact,a.recorddescription ,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=8 AND n.recordid=a.id AND n.islatest=1),'') urgency,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=2 AND n.recordid=a.id AND n.islatest=1),'') typename,coalesce((SELECT q.id status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=2 AND n.recordid=a.id AND n.islatest=1),0) typeid FROM trnrecord a,maprecordtorecorddifferentiation b,mstclientuser e WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.code =? AND b.recordid=a.id AND b.islatest=1 AND a.userid=e.id"
	logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql, clientID, orgnID, req.RecordNo)
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return v, recordtypeid, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.Isvipuser, &v.Source, &v.Clientname, &v.Mstorgnhirarchyname, &v.Title, &v.Recordid, &v.Priority, &v.Status, &v.Impact, &v.Description, &v.Urgency, &v.RecordType, &recordtypeid)

	}
	logger.Log.Println("Results: ", v)
	return v, recordtypeid, nil
}

func (dbc DbConn) GetExternalRecordDetailsByDate(clientID int64, orgnID int64, req *entities.RecordDetailsRequestEntityAPI) ([]entities.RecordDetailsRequestEntityAPI, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsRequestEntityAPI{}

	//var sql = "SELECT distinct COALESCE(a.requestername,'') requestername,COALESCE(a.requesteremail,'') requesteremail,COALESCE(a.requestermobile,'') requestermobile,COALESCE(a.requesterlocation,'') requesterlocation,a.source,(SELECT name FROM mstclient WHERE id=a.clientid) AS Clientname,(SELECT name FROM mstorgnhierarchy WHERE id=a.mstorgnhirarchyid) As Mstorgnhierarchyname,a.recordtitle,a.code,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1),'') priority,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1),'') status,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=7 AND n.recordid=a.id AND n.islatest=1),'') impact,a.recorddescription ,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=8 AND n.recordid=a.id AND n.islatest=1),'') urgency,coalesce((SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=2 AND n.recordid=a.id AND n.islatest=1),'') typename FROM trnrecord a,maprecordtorecorddifferentiation b,mstclientuser e WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.id  in (select x.recordid from mstrecordactivitylogs x where x.clientid=? AND x.mstorgnhirarchyid=? and x.createddate between round(UNIX_TIMESTAMP(?)-(SELECT zone.utcdiff FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=a.clientid AND mstorgnhierarchy.id=a.mstorgnhirarchyid)) and round(UNIX_TIMESTAMP(?)-(SELECT zone.utcdiff FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=a.clientid AND mstorgnhierarchy.id=a.mstorgnhirarchyid)) group by x.recordid) AND b.recordid=a.id AND b.islatest=1 AND a.userid=e.id ORDER BY a.id desc"
	var sql = "select a.ticketid,a.requestorid from recordfulldetails a where a.clientid=? and a.mstorgnhirarchyid=? and a.createddatetime between ? and ?"
	logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql, clientID, orgnID, req.Fromdate, req.Todate)
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsRequestEntityAPI{}
		rows.Scan(&v.RecordNo, &v.Userid)
		values = append(values, v)
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}

func (mdao DbConn) GetExternalcategorynames(ClientID int64, Mstorgnhirarchyid int64, RecorddifftypeID int64, RecorddifID int64, RecordID int64) ([]entities.Categorydetails, error) {
	var sql = "SELECT distinct (SELECT typename from mstrecorddifferentiationtype WHERE id=a.torecorddifftypeid) lable,d.name FROM mstrecordtype a, mstrecorddifferentiationtype b,maprecordtorecorddifferentiation c,mstrecorddifferentiation d WHERE a.torecorddifftypeid = b.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.fromrecorddifftypeid=? AND a.fromrecorddiffid=? AND b.parentid=1 AND c.recordid=? AND c.islatest=1 AND a.torecorddifftypeid=c.recorddifftypeid and c.recorddiffid=d.id"
	v := []entities.Categorydetails{}
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecorddifftypeID, RecorddifID, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return v, err
	}
	for rows.Next() {
		values := entities.Categorydetails{}
		err = rows.Scan(&values.Label, &values.Categoryname)
		logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		v = append(v, values)
	}
	return v, nil
}

func (mdao DbConn) GetWorkflowdetails(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (entities.Workflowdetails, error) {
	var sql = "SELECT concat(c.firstname,' ',c.lastname) name,coalesce(d.name,'') as name  FROM maprequestorecord a,mstrequest b,mstclientuser c,mstsupportgrp d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.mstrequestid= b.id AND b.mstuserid=c.id AND b.mstgroupid=d.id"
	v := entities.Workflowdetails{}
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return v, err
	}
	for rows.Next() {
		err = rows.Scan(&v.Asigneename, &v.Asigneegrp)
		logger.Log.Println("GetLaststagevalue rows.next() Error", err)
	}
	return v, nil
}

func (mdao DbConn) GetClientID(Clientname string) (int64, error) {
	logger.Log.Println("In side GetClientID")
	var ClientID int64
	var sql = "SELECT id FROM mstclient WHERE name=?"
	rows, err := mdao.DB.Query(sql, Clientname)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetClientID Get Statement Prepare Error", err)
		return ClientID, err
	}
	for rows.Next() {
		err = rows.Scan(&ClientID)
		logger.Log.Println("GetReopencount rows.next() Error", err)
	}
	return ClientID, nil
}

func (mdao DbConn) GetOrgnID(Orgnname string) (int64, error) {
	logger.Log.Println("In side GetOrgnID")
	var OrgnID int64
	var sql = "SELECT id FROM mstorgnhierarchy WHERE name=?"
	rows, err := mdao.DB.Query(sql, Orgnname)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOrgnID Get Statement Prepare Error", err)
		return OrgnID, err
	}
	for rows.Next() {
		err = rows.Scan(&OrgnID)
		logger.Log.Println("GetOrgnID rows.next() Error", err)
	}
	return OrgnID, nil
}

func (mdao DbConn) GetLoginID(LogID string, clientID int64, orgnID int64) (int64, error) {
	logger.Log.Println("In side GetLoginID")
	logger.Log.Println("In side GetLoginID---------------------------------->", clientID, orgnID, LogID)
	var LoginID int64
	var sql = "SELECT id FROM mstclientuser WHERE clientid=? AND mstorgnhirarchyid=? AND loginname=? AND activeflag=1 AND deleteflag=0"
	rows, err := mdao.DB.Query(sql, clientID, orgnID, LogID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLoginID Get Statement Prepare Error", err)
		return LoginID, err
	}
	for rows.Next() {
		err = rows.Scan(&LoginID)
		logger.Log.Println("GetLoginID rows.next() Error", err)
	}
	return LoginID, nil
}

func (mdao DbConn) GetGrpID(loginID int64, ClientID int64, OrgnIDTypeID int64) (int64, error) {
	logger.Log.Println("In side GetLoginID")
	logger.Log.Println("In side GetLoginID---------------------------------->", ClientID, OrgnIDTypeID, loginID)
	var grpID int64

	var sql = "SELECT b.id FROM mstgroupmember a,mstclientsupportgroup b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.userid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.groupid=b.grpid AND b.supportgrouplevelid=1 AND b.activeflg=1 AND b.deleteflg=0"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, loginID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLoginID Get Statement Prepare Error", err)
		return grpID, err
	}
	for rows.Next() {
		err = rows.Scan(&grpID)
		logger.Log.Println("GetLoginID rows.next() Error", err)
	}
	return grpID, nil
}

func (mdao DbConn) GetRequestorInfo(loginID int64) (entities.ExternalRequestorInfo, error) {
	logger.Log.Println("In side GetRequestorInfo")
	v := entities.ExternalRequestorInfo{}
	var sql = "SELECT concat(firstname,' ',lastname) name,useremail,usermobileno,COALESCE(branch,'') FROM mstclientuser WHERE id=?"
	rows, err := mdao.DB.Query(sql, loginID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRequestorInfo Get Statement Prepare Error", err)
		return v, err
	}
	for rows.Next() {
		err = rows.Scan(&v.Requestername, &v.Requesteremail, &v.Requestermobile, &v.Requesterlocation)
		logger.Log.Println("GetRequestorInfo rows.next() Error", err)
	}
	return v, nil
}

func (mdao DbConn) GetLabelID(labelname string, ClientID int64, OrgnIDTypeID int64) (int64, error) {
	logger.Log.Println("In side GetLabelID")
	var labelID int64
	var sql = "SELECT id FROM mstrecorddifferentiationtype WHERE typename=? AND clientid=? AND mstorgnhirarchyid=?"
	rows, err := mdao.DB.Query(sql, labelname, ClientID, OrgnIDTypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLabelID Get Statement Prepare Error", err)
		return labelID, err
	}
	for rows.Next() {
		err = rows.Scan(&labelID)
		logger.Log.Println("GetLabelID rows.next() Error", err)
	}
	return labelID, nil
}

func (mdao DbConn) GetLabelIDAgainstTickettype(labelname string, ClientID int64, OrgnIDTypeID int64, TypedifftypeID int64, TypediffID int64) (int64, error) {
	logger.Log.Println("In side GetLabelID")
	var labelID int64
	var sql = "SELECT torecorddifftypeid FROM mstrecordtype WHERE clientid=? AND mstorgnhirarchyid=? AND fromrecorddifftypeid=? AND fromrecorddiffid=? AND torecorddifftypeid in (SELECT id FROM mstrecorddifferentiationtype WHERE clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1 AND typename=?) AND torecorddiffid=0 AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, TypedifftypeID, TypediffID, ClientID, OrgnIDTypeID, labelname)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLabelID Get Statement Prepare Error", err)
		return labelID, err
	}
	for rows.Next() {
		err = rows.Scan(&labelID)
		logger.Log.Println("GetLabelID rows.next() Error", err)
	}
	return labelID, nil
}

func (mdao DbConn) GetLabelvalueID(labelname string) (int64, error) {
	logger.Log.Println("In side GetLabelvalueID")
	var labelvalID int64
	var sql = "SELECT id FROM mstrecorddifferentiation WHERE name=?"
	rows, err := mdao.DB.Query(sql, labelname)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLabelvalueID Get Statement Prepare Error", err)
		return labelvalID, err
	}
	for rows.Next() {
		err = rows.Scan(&labelvalID)
		logger.Log.Println("GetLabelvalueID rows.next() Error", err)
	}
	return labelvalID, nil
}

func (mdao DbConn) GetDifferentionID(ClientID int64, OrgnIDTypeID int64, TypeID int64, Name string, parentID int64) (int64, error) {
	logger.Log.Println("In side GetLabelvalueID")
	logger.Log.Println("In side GetLabelvalueID----->", Name, TypeID, ClientID, OrgnIDTypeID, parentID)
	var labelvalID int64
	var sql = "SELECT id FROM mstrecorddifferentiation WHERE name=? AND recorddifftypeid=? AND clientid=? AND mstorgnhirarchyid=? AND parentid=?"
	rows, err := mdao.DB.Query(sql, Name, TypeID, ClientID, OrgnIDTypeID, parentID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLabelvalueID Get Statement Prepare Error", err)
		return labelvalID, err
	}
	for rows.Next() {
		err = rows.Scan(&labelvalID)
		logger.Log.Println("GetLabelvalueID rows.next() Error", err)
	}
	return labelvalID, nil
}

func (mdao DbConn) GetTermID(ClientID int64, OrgnIDTypeID int64, Name string) (int64, int64, error) {
	logger.Log.Println("In side GetTermID")
	logger.Log.Println("In side GetTermID----->", Name)
	var termID int64
	var seq int64
	var sql = "SELECT id,seq FROM mstrecordterms WHERE termname=? AND clientid=? AND mstorgnhirarchyid=?"
	rows, err := mdao.DB.Query(sql, Name, ClientID, OrgnIDTypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermID Get Statement Prepare Error", err)
		return termID, seq, err
	}
	for rows.Next() {
		err = rows.Scan(&termID, &seq)
		logger.Log.Println("GetTermID rows.next() Error", err)
	}
	return termID, seq, nil
}

func (mdao DbConn) GetAdditional(ClientID int64, OrgnIDTypeID int64, termID int64) (int64, error) {
	logger.Log.Println("In side GetAdditional")
	logger.Log.Println("In side GetAdditional----->", termID)
	var additionID int64
	var sql = "SELECT id FROM mstrecordfield WHERE recordtermid=? AND clientid=? AND mstorgnhirarchyid=?"
	rows, err := mdao.DB.Query(sql, termID, ClientID, OrgnIDTypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAdditional Get Statement Prepare Error", err)
		return additionID, err
	}
	for rows.Next() {
		err = rows.Scan(&additionID)
		logger.Log.Println("GetAdditional rows.next() Error", err)
	}
	return additionID, nil
}

func (mdao DbConn) GetRequestIdbyRecordId(recordID int64) (error, []entities.TransactionEntity) {
	var getRequestIdByStageId = "SELECT coalesce(mstrequestid,0) as Requestid from maprequestorecord where recordid=? and deleteflg=0 and activeflg=1 "
	requestIds := []entities.TransactionEntity{}
	stmt, err := mdao.DB.Prepare(getRequestIdByStageId)
	if err != nil {
		log.Print("GetRequestIdbyRecordId isWorkflowFirstStep Statement Prepare Error", err)
		return err, requestIds
	}
	defer stmt.Close()
	rows, err := stmt.Query(recordID)
	if err != nil {
		log.Print("GetRequestIdbyRecordId isWorkflowFirstStep Save Statement Execution Error", err)
		return err, requestIds
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.TransactionEntity{}
		rows.Scan(&value.Requestid)
		requestIds = append(requestIds, value)
	}
	return nil, requestIds
}

func (mdao DbConn) Gethopcount(id int64) (error, []int64) {
	var getgrpidsbymainreqstid = "select coalesce(mstgroupid,0) from mstrequesthistory where mainrequestid=? and activeflg =1 and deleteflg=0 and mstgroupid IN (select grpid from mstclientsupportgroup where supportgrouplevelid>1 and activeflg =1 and deleteflg=0)"
	var S []int64
	stmt, err := mdao.DB.Prepare(getgrpidsbymainreqstid)
	if err != nil {
		logger.Log.Print("Gethopcount Statement Prepare Error", err)
		log.Print("Gethopcount Statement Prepare Error", err)
		return err, S
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		logger.Log.Print("Gethopcount Statement Execution Error", err)
		log.Print("Gethopcount Statement Execution Error", err)
		return err, S
	}
	defer rows.Close()
	for rows.Next() {
		var value int64
		rows.Scan(&value)
		S = append(S, value)
	}
	return nil, S
}

func (mdao DbConn) GetRecordTypeDetails(ClientID int64, OrgnIDTypeID int64, lastlabelID int64, lastlabelcatID int64) (int64, int64, error) {
	logger.Log.Println("In side GetAdditional")
	logger.Log.Println("In side GetAdditional----->", lastlabelID, lastlabelcatID)
	var fromrecorddifftypeid int64
	var fromrecorddiffid int64
	var sql = "SELECT fromrecorddifftypeid,fromrecorddiffid FROM mstrecordtype where clientid=? AND mstorgnhirarchyid=? AND torecorddifftypeid=? AND torecorddiffid=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, lastlabelID, lastlabelcatID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAdditional Get Statement Prepare Error", err)
		return fromrecorddifftypeid, fromrecorddiffid, err
	}
	for rows.Next() {
		err = rows.Scan(&fromrecorddifftypeid, &fromrecorddiffid)
		logger.Log.Println("GetAdditional rows.next() Error", err)
	}
	return fromrecorddifftypeid, fromrecorddiffid, nil
}

func (mdao DbConn) GetRecordTypeDetailsAgainstTickettype(ClientID int64, OrgnIDTypeID int64, Tickettype string) (int64, int64, error) {
	logger.Log.Println("In side GetRecordTypeDetailsAgainstTickettype")
	logger.Log.Println("In side GetRecordTypeDetailsAgainstTickettype----->", ClientID, OrgnIDTypeID, Tickettype)
	var recorddifftypeid int64
	var recorddiffid int64
	var sql = "SELECT recorddifftypeid,id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=2 AND name=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, Tickettype)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordTypeDetailsAgainstTickettype Get Statement Prepare Error", err)
		return recorddifftypeid, recorddiffid, err
	}
	for rows.Next() {
		err = rows.Scan(&recorddifftypeid, &recorddiffid)
		logger.Log.Println("GetRecordTypeDetailsAgainstTickettype rows.next() Error", err)
	}
	return recorddifftypeid, recorddiffid, nil
}

func (mdao DbConn) GetDifferentionIDbySeq(ClientID int64, OrgnIDTypeID int64, TypeID int64, Sequance int64) (int64, error) {
	logger.Log.Println("In side GetLabelvalueID")
	logger.Log.Println("In side GetLabelvalueID----->", Sequance, TypeID, ClientID, OrgnIDTypeID)
	var labelvalID int64
	var sql = "SELECT id FROM mstrecorddifferentiation WHERE seqno=? AND recorddifftypeid=? AND clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, Sequance, TypeID, ClientID, OrgnIDTypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLabelvalueID Get Statement Prepare Error", err)
		return labelvalID, err
	}
	for rows.Next() {
		err = rows.Scan(&labelvalID)
		logger.Log.Println("GetLabelvalueID rows.next() Error", err)
	}
	return labelvalID, nil
}

func (mdao DbConn) GetPriorityID(ClientID int64, OrgnIDTypeID int64, lastlabelcatID int64) (int64, error) {
	logger.Log.Println("In side GetLabelvalueID")
	logger.Log.Println("In side GetLabelvalueID----->", lastlabelcatID)
	var priorityID int64
	var sql = "SELECT mstrecorddifferentiationpriorityid FROM mstbusinessmatrix WHERE clientid=? AND mstorgnhirarchyid=? AND mstrecorddifferentiationcatid=? AND activeflg=1 AND deleteflg=0"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, lastlabelcatID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLabelvalueID Get Statement Prepare Error", err)
		return priorityID, err
	}
	for rows.Next() {
		err = rows.Scan(&priorityID)
		logger.Log.Println("GetLabelvalueID rows.next() Error", err)
	}
	return priorityID, nil
}

func (mdao DbConn) GetLastupdatednames(ClientID int64, OrgnIDTypeID int64, recordID int64) (string, error) {
	logger.Log.Println("In side GetLastupdatednames")
	logger.Log.Println("In side GetLastupdatednames----->", recordID)
	var name string
	var sql = "SELECT concat(b.firstname,' ',b.lastname) name FROM mstrecordactivitylogs a,mstclientuser b where a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.createdid=b.id order by a.id desc limit 1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, recordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLabelvalueID Get Statement Prepare Error", err)
		return name, err
	}
	for rows.Next() {
		err = rows.Scan(&name)
		logger.Log.Println("GetLabelvalueID rows.next() Error", err)
	}
	return name, nil
}

func (mdao DbConn) GetRecordTypeID(ClientID int64, OrgnIDTypeID int64, recordID int64) (int64, error) {
	logger.Log.Println("In side GetRecordTypeID")
	logger.Log.Println("In side GetRecordTypeID----->", recordID)
	var typeID int64
	var sql = "SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=2 AND islatest=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, recordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordTypeID Get Statement Prepare Error", err)
		return typeID, err
	}
	for rows.Next() {
		err = rows.Scan(&typeID)
		logger.Log.Println("GetRecordTypeID rows.next() Error", err)
	}
	return typeID, nil
}

func (mdao DbConn) GetrecordnxtstateID(ClientID int64, OrgnIDTypeID int64, Statusname string) (int64, error) {
	logger.Log.Println("In side GetrecordnxtstateID")
	logger.Log.Println("In side GetrecordnxtstateID----->", Statusname)
	var stateID int64
	var sql = "SELECT mststateid FROM maprecordstatetodifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=3 AND recorddiffid in(SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=3 AND name=? AND deleteflg=0 AND activeflg=1 ) AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, ClientID, OrgnIDTypeID, Statusname)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetrecordnxtstateID Get Statement Prepare Error", err)
		return stateID, err
	}
	for rows.Next() {
		err = rows.Scan(&stateID)
		logger.Log.Println("GetrecordnxtstateID rows.next() Error", err)
	}
	return stateID, nil
}

func (mdao DbConn) GetStageID(ClientID int64, OrgnIDTypeID int64, recordID int64) (int64, error) {
	logger.Log.Println("In side getStageID")
	logger.Log.Println("In side getStageID----->", recordID)
	var stageID int64
	var sql = "SELECT recordstageid FROM trnrecord WHERE clientid=? AND mstorgnhirarchyid=? AND id=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, recordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("getStageID Get Statement Prepare Error", err)
		return stageID, err
	}
	for rows.Next() {
		err = rows.Scan(&stageID)
		logger.Log.Println("getStageID rows.next() Error", err)
	}
	return stageID, nil
}

func (mdao DbConn) GetGrpIDByName(ClientID int64, OrgnIDTypeID int64, Grpname string) (int64, error) {
	logger.Log.Println("In side GetGrpIDByName")
	logger.Log.Println("In side GetGrpIDByName----->", ClientID, OrgnIDTypeID, Grpname)
	var grpID int64
	//var sql = "SELECT id FROM mstsupportgrp WHERE clientid=? AND mstorgnhirarchyid=? AND name=? AND deleteflg=0 AND activeflg=1"
	var sql = "SELECT a.id FROM mstsupportgrp a,mstclientsupportgroup b WHERE b.clientid=? AND b.mstorgnhirarchyid=? AND a.name=? AND a.id=b.grpid AND a.deleteflg=0 AND a.activeflg=1 AND b.deleteflg=0 AND b.activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, Grpname)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetGrpIDByName Get Statement Prepare Error", err)
		return grpID, err
	}
	for rows.Next() {
		err = rows.Scan(&grpID)
		logger.Log.Println("GetGrpIDByName rows.next() Error", err)
	}
	return grpID, nil
}

func (mdao DbConn) GetWokinglabel(ClientID int64, OrgnIDTypeID int64, RecordID int64) (int64, int64, error) {
	logger.Log.Println("In side GetWokinglabel")
	logger.Log.Println("In side GetWokinglabel----->", RecordID)
	var workingtypeID int64
	var workingcatID int64
	var sql = "SELECT recorddifftypeid,recorddiffid FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid =? AND isworking=1 order by id desc limit 1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetWokinglabel Get Statement Prepare Error", err)
		return workingtypeID, workingcatID, err
	}
	for rows.Next() {
		err = rows.Scan(&workingtypeID, &workingcatID)
		logger.Log.Println("GetWokinglabel rows.next() Error", err)
	}
	return workingtypeID, workingcatID, nil
}

func (mdao DbConn) GetMstGrpID(ClientID int64, OrgnIDTypeID int64, workingtypeID int64, workingcatID int64, currentStateID int64, nxtStateID int64) (int64, error) {
	logger.Log.Println("In side GetGrpIDByName")
	logger.Log.Println("In side GetGrpIDByName =================================== *****************************************          ----->", ClientID, OrgnIDTypeID, workingtypeID, workingcatID, currentStateID, nxtStateID)
	var grpID int64
	var sql = "SELECT c.mstgroupid FROM mstprocessrecordmap a,msttransition b,maprecorddifferentiongroup c WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=? AND a.recorddiffid=?  AND a.deleteflg=0 AND a.activeflg=1 AND a.mstprocessid = b.processid AND b.currentstateid=? AND b.previousstateid=? AND b.id = c.transitionid AND c.mstgroupid != 0 AND b.deleteflg=0 AND b.activeflg=1 AND c.deleteflg=0 AND c.activeflg=1 limit 1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, workingtypeID, workingcatID, nxtStateID, currentStateID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetGrpIDByName Get Statement Prepare Error", err)
		return grpID, err
	}
	for rows.Next() {
		err = rows.Scan(&grpID)
		logger.Log.Println("GetGrpIDByName rows.next() Error", err)
	}
	return grpID, nil
}

func (mdao DbConn) GetRecordCurrentGrpID(ClientID int64, OrgnIDTypeID int64, recordID int64) (int64, int64, error) {
	logger.Log.Println("In side getStageID")
	logger.Log.Println("In side getStageID----->", recordID)
	var grpID int64
	var userID int64
	var sql = "select b.mstgroupid,b.mstuserid from maprequestorecord a,mstrequest b where a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.mstrequestid = b.id AND b.activeflg=1 AND b.deleteflg=0"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, recordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("getStageID Get Statement Prepare Error", err)
		return grpID, userID, err
	}
	for rows.Next() {
		err = rows.Scan(&grpID, &userID)
		logger.Log.Println("getStageID rows.next() Error", err)
	}
	return grpID, userID, nil
}

func (mdao DbConn) ValidateAssigneUser(ClientID int64, OrgnIDTypeID int64, changedgrpID int64, changeloginID int64) (int64, error) {
	logger.Log.Println("In side ValidateAssigneUser")
	logger.Log.Println("In side ValidateAssigneUser----->", changedgrpID)
	var ID int64
	var sql = "SELECT id FROM mstgroupmember WHERE clientid=? AND mstorgnhirarchyid=? AND groupid=? AND userid=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, changedgrpID, changeloginID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("ValidateAssigneUser Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("ValidateAssigneUser rows.next() Error", err)
	}
	return ID, nil
}

func (mdao DbConn) GetDifferentiationdstatustls(ClientID int64, OrgnIDTypeID int64, Statusname string) (int64, int64, error) {
	logger.Log.Println("In side GetDifferentiationdstatustls")
	var id int64
	var seq int64
	var sql = "SELECT id,seqno FROM mstrecorddifferentiation WHERE  clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=3 AND name=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, Statusname)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetDifferentiationdstatustls Get Statement Prepare Error", err)
		return id, seq, err
	}
	for rows.Next() {
		err = rows.Scan(&id, &seq)
		logger.Log.Println("GetDifferentiationdstatustls rows.next() Error", err)
	}
	return id, seq, nil
}

func (mdao DbConn) GetTermnamesbystatusID(ClientID int64, OrgnIDTypeID int64, StatusID int64, TypeID int64) ([]string, error) {
	logger.Log.Println("In side GetDifferentiationdstatustls")
	var values []string
	//var sql = "SELECT b.termname FROM mststateterm a,mstrecordterms b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.recorddifftypeid=3 AND a.recorddiffid=? AND a.recordtermid=b.id AND a.iscompulsory=1"
	var sql = "SELECT b.termname FROM mststateterm a,mstrecordterms b,mststateterm c WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.recorddifftypeid=3 AND a.recorddiffid=? AND c.recorddifftypeid=2 AND c.recorddiffid=? AND c.activeflg=1 AND c.deleteflg=0 AND a.recordtermid=b.id AND c.recordtermid=b.id AND a.iscompulsory=1"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnIDTypeID, StatusID, TypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetDifferentiationdstatustls Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		var val string
		err = rows.Scan(&val)
		values = append(values, val)
		logger.Log.Println("GetDifferentiationdstatustls rows.next() Error", err)
	}
	return values, nil
}

func (mdao DbConn) GetTermsdetails(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (map[string]interface{}, error) {
	var sql = "SELECT mstrecordterms.termname,COALESCE(trnreordtracking.recordtrackvalue,'NA') as Termvalue FROM trnreordtracking,mstrecordterms WHERE trnreordtracking.id IN (SELECT MAX(id) AS id FROM trnreordtracking WHERE clientid = ? AND mstorgnhirarchyid=? AND recordid=? GROUP BY recordtermid) AND trnreordtracking.recordtermid = mstrecordterms.id AND mstrecordterms.seq NOT IN(29,30) ORDER BY mstrecordterms.id"
	v := map[string]interface{}{}
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return v, err
	}
	for rows.Next() {
		var termname string
		var termvalue string
		err = rows.Scan(&termname, &termvalue)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
		//logger.Log.Println("Termname is -------->", termname)
		//logger.Log.Println("Termvalue is ------->", termvalue)
		v[termname] = termvalue

	}
	//logger.Log.Println("map value is ------->", v)
	return v, nil
}
func (mdao DbConn) GetTermsdetailss(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (map[string]interface{}, error) {
	var sql = "select  coalesce(a.termname,''),coalesce(b.recordtrackvalue,'') as termvalue from  mststateterm c,recordfulldetails d,mstrecordterms a left join trnreordtracking b on b.id=(select max(e.id ) from trnreordtracking e where e.clientid=a.clientid and e.mstorgnhirarchyid=a.mstorgnhirarchyid and e.recordid=? and e.activeflg=1 and e.deleteflg=0 and e.recordtermid=a.id) and b.clientid=a.clientid and b.mstorgnhirarchyid=a.mstorgnhirarchyid and b.activeflg=1 and b.deleteflg=0 where d.recordid=? and d.clientid=? and d.mstorgnhirarchyid=? and d.activeflg=1 and d.deleteflg=0 and d.clientid=c.clientid and d.mstorgnhirarchyid=c.mstorgnhirarchyid and d.tickettypeid=c.recorddiffid and c.activeflg=1 and c.deleteflg=0 and c.recordtermid=a.id and c.clientid=a.clientid and c.mstorgnhirarchyid=a.mstorgnhirarchyid"
	v := map[string]interface{}{}
	rows, err := mdao.DB.Query(sql, RecordID, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return v, err
	}
	for rows.Next() {
		var termname string
		var termvalue string
		err = rows.Scan(&termname, &termvalue)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
		//logger.Log.Println("Termname is -------->", termname)
		//logger.Log.Println("Termvalue is ------->", termvalue)
		v[termname] = termvalue

	}
	//logger.Log.Println("map value is ------->", v)
	return v, nil
}
func (mdao DbConn) GetFromRecordfulldetails(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, values *entities.RecordDetailsEntityAPI) error {
	var sql = "SELECT coalesce(assigneduser,''),coalesce(assignedgroup,''),createddatetime, coalesce(orgcreatorname,'') as orgcreatorname,coalesce(requestorloginid,'') as requestorloginid,coalesce(orgcreatorphone,'') as orgcreatorphone,coalesce(requestorloginid,'') as  requestorloginid,coalesce(resogroup,'') as resogroup,coalesce(resolveduser,'') as resolveduser,coalesce(respsladuedatetime,'') as respsladuedatetime,coalesce(resosladuedatetime,'') as resosladuedatetime,coalesce(respclockstatus,'') as respclockstatus,coalesce(respslabreachstatus,'') as respslabreachstatus,coalesce(respoverduetime,'') as respoverduetime,coalesce(resoclockstatus,'') as resoclockstatus,coalesce(resolslabreachstatus,'') as resolslabreachstatus,coalesce(resooverduetime,'') as resooverduetime FROM recordfulldetails where recordid=? and clientid=? and mstorgnhirarchyid=? and activeflg=1 and deleteflg=0"
	rows, err := mdao.DB.Query(sql, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return err
	}
	for rows.Next() {

		err = rows.Scan(&values.Assignee,&values.AssignedGroup,&values.Createdatetime, &values.OriginalCreatedByFullName, &values.OriginalCreatedByLoginID, &values.OriginalCreatedByPrimaryContact, &values.Requestorloginid, &values.ResolvedByGroup, &values.ResolvedByUser, &values.ResponseDueDate, &values.ResolutionDueDate, &values.ResponseClockStatus, &values.ResponseSLABreachedStatus, &values.ResponseSLAOverdue, &values.ResolutionClockStatus, &values.ResolutionSLABreachedStatus, &values.ResolutionSLAOverdue)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
	}
	//logger.Log.Println("map value is ------->", v)
	return nil
}
func (mdao DbConn) Getchild(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) ([]int64, error) {
	var values []int64
	var sql = "SELECT distinct childrecordid FROM mstparentchildmap where parentrecordid=? and clientid=? and mstorgnhirarchyid=? and activeflg=1 and deleteflg=0;"
	rows, err := mdao.DB.Query(sql, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		var v int64
		err = rows.Scan(&v)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
		values = append(values, v)
	}
	//logger.Log.Println("map value is ------->", v)
	return values, nil
}
func (mdao DbConn) GetRecordNo(RecordID int64) (string, error) {
	var values string
	var sql = "SELECT  coalesce(ticketid,'') FROM recordfulldetails where recordid=? and activeflg=1 and deleteflg=0;"
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		// var v string
		err = rows.Scan(&values)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
		// values = append(values, v)
	}
	//logger.Log.Println("map value is ------->", v)
	return values, nil
}

func (mdao DbConn) Getparent(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) ([]int64, error) {
	var values []int64
	var sql = "SELECT distinct parentrecordid   FROM mstparentchildmap where childrecordid=? and clientid=? and mstorgnhirarchyid=? and activeflg=1 and deleteflg=0;"
	rows, err := mdao.DB.Query(sql, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		var v int64
		err = rows.Scan(&v)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
		values = append(values, v)
	}
	//logger.Log.Println("map value is ------->", v)
	return values, nil
}
func (mdao DbConn) Getlastmodified(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	var lastmodified int64
	var sql = "select createddate from mstrecordactivitylogs where recordid=? and clientid=? and mstorgnhirarchyid=? and activeflg=1 and deleteflg=0 order by id desc limit 1"
	rows, err := mdao.DB.Query(sql, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return lastmodified, err
	}
	for rows.Next() {

		err = rows.Scan(&lastmodified)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
	}
	//logger.Log.Println("map value is ------->", v)
	return lastmodified, nil
}
func (mdao DbConn) Getlikrecord(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, values *entities.RecordDetailsEntityAPI) error {
	var sql = "select distinct coalesce(a.ticketid,'') as rocordno from recordfulldetails a,maprecordtolinkrecords b where b.recordid=? and b.clientid=? and b.mstorgnhirarchyid=? and b.activeflg=1 and b.deleteflg=0 and b.linkrecordid=a.recordid and a.clientid=b.clientid and a.mstorgnhirarchyid=b.mstorgnhirarchyid and a.activeflg=1 and a.deleteflg=0"
	rows, err := mdao.DB.Query(sql, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return err
	}
	for rows.Next() {
		var value string
		err = rows.Scan(&values)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
		values.Linkedtickets = append(values.Linkedtickets, value)
	}
	//logger.Log.Println("map value is ------->", v)
	return nil
}
func (mdao DbConn) GetSlacompliance(ClientID int64, Mstorgnhirarchyid int64, recordtype int64, values *entities.RecordDetailsEntityAPI) error {
	var sql = "select coalesce(responsecompliance,'') as responsecompliance,coalesce(resolutioncompliance,'') as resolutioncompliance from mstslafullfillmentcriteria where mstrecorddifferentiationtickettypeid=? and clientid=? and mstorgnhirarchyid=? and deleteflg=0"
	rows, err := mdao.DB.Query(sql, recordtype, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return err
	}
	for rows.Next() {

		err = rows.Scan(&values.Responseslacompliance, &values.Resolutionslacompliance)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
	}
	//logger.Log.Println("map value is ------->", v)
	return nil
}
func (mdao DbConn) GetTotaleffort(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) ([]string, error) {
	var totaleffort []string
	var sql = "SELECT a.recordtrackvalue FROM trnreordtracking a,mstrecordterms b where recordtermid=b.id   and a.clientid=b.clientid and a.mstorgnhirarchyid=b.mstorgnhirarchyid and b.seq=31 and a.clientid=? and a.mstorgnhirarchyid=? and a.recordid=? and a.deleteflg=0;"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return totaleffort, err
	}
	for rows.Next() {

		err = rows.Scan(&totaleffort)
		if err != nil {
			logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		}
	}
	//logger.Log.Println("map value is ------->", v)
	return totaleffort, nil
}

/*func (mdao DbConn) GetRecordDiffId(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side GetRecordDiffId")
	var recorddiffid int64
	var sql = "select recorddiffid from maprecordtorecorddifferentiation where clientid = ? and mstorgnhirarchyid = ? and recordid = ? and recorddifftypeid = 2 and islatest = 1 and deleteflg = 0"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordDiffId Get Statement Prepare Error", err)
		return recorddiffid, err
	}
	for rows.Next() {

		err = rows.Scan(&recorddiffid)
		if err != nil {
			logger.Log.Println("GetRecordDiffId rows.next() Error", err)
		}
	}
	//logger.Log.Println("map value is ------->", v)
	return recorddiffid, nil
}*/
