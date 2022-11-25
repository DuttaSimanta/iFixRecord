package dao

import (
	"database/sql"
	"errors"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

func (dbc DbConn) CheckChildRecordIsExist(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side CheckChildRecordIsExist Dao")
	var value int64
	var query = "SELECT COUNT(id) total FROM mstparentchildmap WHERE deleteflg=0 AND activeflg=1 AND clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND childrecordid=? AND isattached='Y' "
	var params []interface{}

	params = append(params, req.Clientid)
	params = append(params, req.Mstorgnhirarchyid)
	params = append(params, req.RecordDiffTypeid)
	params = append(params, req.RecordDiffid)
	params = append(params, req.ChildID)

	logger.Log.Println("Main Query: ", query)
	logger.Log.Println("Main Query: ", params)
	rows, err := dbc.DB.Query(query, params...)
	if err != nil {
		logger.Log.Println("CheckChildRecordIsExist Get Statement Prepare Error", err)
		return value, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
	}
	return value, nil
}

func SaveParentRecord(tx *sql.Tx, page *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("parameters -->", page.Clientid, page.Mstorgnhirarchyid, page.RecordDiffTypeid, page.RecordDiffid, page.ParentID, page.ChildID)
	var query = "INSERT INTO mstparentchildmap(clientid, mstorgnhirarchyid, recorddifftypeid, recorddiffid, parentrecordid, childrecordid) VALUES (?,?,?,?,?,?)"
	stmt, err := tx.Prepare(query)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(page.Clientid, page.Mstorgnhirarchyid, page.RecordDiffTypeid, page.RecordDiffid, page.ParentID, page.ChildID)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	lastInsertedStageID, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Last ID fetch Error")
	}

	return lastInsertedStageID, nil
}

func (dbc DbConn) CheckParentAsChildRecord1(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side CheckParentAsChildRecord Dao")
	var value int64
	var query = "SELECT COUNT(id) total FROM mstparentchildmap WHERE deleteflg=0 AND activeflg=1 AND clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND childrecordid=? AND isattached='Y' "
	var params []interface{}

	params = append(params, req.Clientid)
	params = append(params, req.Mstorgnhirarchyid)
	params = append(params, req.RecordDiffTypeid)
	params = append(params, req.RecordDiffid)
	params = append(params, req.ParentID)

	logger.Log.Println("Main Query: ", query)
	logger.Log.Println("Main Query: ", params)
	rows, err := dbc.DB.Query(query, params...)
	if err != nil {
		logger.Log.Println("CheckParentAsChildRecord Get Statement Prepare Error", err)
		return value, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
	}
	return value, nil
}

func (mdao DbConn) CheckParentOrNot(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side CheckParentOrNot")
	var count int64
	var query = "SELECT count(c.id) FROM trnrecord a,maprecordtorecorddifferentiation b,mstparentchildmap c Where a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.code=? AND a.id=b.recordid AND b.recorddifftypeid=? AND b.recorddiffid=? AND a.id = c.parentrecordid AND c.activeflg=1 AND c.deleteflg=0"
	rows, err := mdao.DB.Query(query, req.Clientid, req.Mstorgnhirarchyid, req.RecordNo, req.RecordDiffTypeid, req.RecordDiffid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("CheckParentOrNot Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("CheckParentOrNot rows.next() Error", err)
	}
	return count, nil
}

func (dbc DbConn) GetProcessidbyworkingcatidforsearch(req *entities.ChildRecordSearchEntity) (entities.WorkFlowEntity, error) {
	logger.Log.Println("In side GetProcessidbyworkingcatid Dao")
	workflowid := entities.WorkFlowEntity{}
	var sql1 = "SELECT mstprocessrecordmap.mstprocessid,maprecordtorecorddifferentiation.recorddifftypeid,maprecordtorecorddifferentiation.recorddiffid FROM maprecordtorecorddifferentiation JOIN mstprocessrecordmap ON maprecordtorecorddifferentiation.clientid = mstprocessrecordmap.clientid AND maprecordtorecorddifferentiation.mstorgnhirarchyid = mstprocessrecordmap.mstorgnhirarchyid AND maprecordtorecorddifferentiation.recorddiffid=mstprocessrecordmap.recorddiffid AND maprecordtorecorddifferentiation.recorddifftypeid=mstprocessrecordmap.recorddifftypeid WHERE maprecordtorecorddifferentiation.activeflg = 1 AND maprecordtorecorddifferentiation.deleteflg = 0 AND mstprocessrecordmap.activeflg = 1 AND mstprocessrecordmap.deleteflg = 0 AND maprecordtorecorddifferentiation.isworking = 1 AND maprecordtorecorddifferentiation.clientid = ? AND maprecordtorecorddifferentiation.mstorgnhirarchyid = ? AND maprecordtorecorddifferentiation.recordid = ? AND maprecordtorecorddifferentiation.recordstageid = ?"
	logger.Log.Println("GetProcessidbyworkingcatid Query: ", sql1)
	logger.Log.Println("GetProcessidbyworkingcatid Params: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordStageID)
	err := dbc.DB.QueryRow(sql1, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordStageID).Scan(&workflowid.WorkFlowID, &workflowid.CatTypeID, &workflowid.CatID)
	switch err {
	case sql.ErrNoRows:
		return workflowid, nil
	case nil:
		return workflowid, nil
	default:
		logger.Log.Println("GetProcessidbyworkingcatid Get Statement Prepare Error", err)
		return workflowid, err
	}
}

func (dbc DbConn) GetParentRecordDetailsByNo(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsEntity{}

	var sql = "SELECT   COALESCE(trnrecord.requestername,'') requestername,COALESCE(trnrecord.requesteremail,'') requesteremail,COALESCE(trnrecord.requestermobile,'') requestermobile,COALESCE(trnrecord.requesterlocation,'') requesterlocation,  trnrecord.id,  trnrecord.source,  trnrecord.clientid,  trnrecord.mstorgnhirarchyid,  trnrecord.recordtitle,  trnrecord.code,  trnrecord.usergroupid,  mstclientsupportgroup.supportgroupname,  supportgrouplevel.id supportgrouplevelid,  supportgrouplevel.name levelname,  prioritytab.recorddiffid priorityid,  prioritytab.name priority,  mstclientuser.id creatorid,  mstclientuser.loginname createdby,  COALESCE(mstclientuser.vipuser, 'N') vipuser,  COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime+zone.utcdiff),'%d-%m-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=trnrecord.clientid AND mstorgnhierarchy.id=trnrecord.mstorgnhirarchyid),  '') createdatetime  ,  statustab.recorddiffid statusid,statustab.seqno statusseqno,  statustab.name status,  COALESCE(impacttab.recorddiffid, 0) impactid,  COALESCE(impacttab.name, '') impact,  trnrecord.requesterinfo,  trnrecord.recorddescription,  COALESCE(urgencytab.recorddiffid, 0) urgencyid,  COALESCE(urgencytab.name, '') urgency,  typetab.recorddiffid typeid,  typetab.name typename,  COALESCE(mstrequest.mstgroupid, 0) assignedgrpid,  COALESCE(mstrequest.mstuserid, 0) assigneduid,  COALESCE(reqclient.loginname, '') assigneduser,  COALESCE(clisupgroup.supportgroupname, '') assignedgrpname,  COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid,  COALESCE(supgrouplvl.name, '') assignedsupgrlvl,  trnrecord.recordstageid,  COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%m-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = trnrecord.id),  '') duedate FROM  trnrecord  JOIN  mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1  AND mstclientsupportgroup.deleteflg = 0  AND trnrecord.usergroupid = mstclientsupportgroup.id  JOIN  supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id  JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 4 AND maprecordtorecorddifferentiation.islatest = 1) prioritytab ON prioritytab.recordid = trnrecord.id  AND trnrecord.clientid = prioritytab.clientid  AND trnrecord.mstorgnhirarchyid = prioritytab.mstorgnhirarchyid  JOIN  mstclientuser ON mstclientuser.activeflag = 1  AND mstclientuser.deleteflag = 0  AND mstclientuser.clientid = trnrecord.clientid  AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid  AND trnrecord.userid = mstclientuser.id  JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid,mstrecorddifferentiation.seqno  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 2 AND maprecordtorecorddifferentiation.islatest=1) statustab ON statustab.recordid = trnrecord.id  AND trnrecord.clientid = statustab.clientid  AND trnrecord.mstorgnhirarchyid = statustab.mstorgnhirarchyid  LEFT JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 6 AND maprecordtorecorddifferentiation.islatest = 1) impacttab ON impacttab.recordid = trnrecord.id  AND trnrecord.clientid = impacttab.clientid  AND trnrecord.mstorgnhirarchyid = impacttab.mstorgnhirarchyid  LEFT JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 7 AND maprecordtorecorddifferentiation.islatest = 1) urgencytab ON urgencytab.recordid = trnrecord.id  AND trnrecord.clientid = urgencytab.clientid  AND trnrecord.mstorgnhirarchyid = urgencytab.mstorgnhirarchyid  JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 1 AND maprecordtorecorddifferentiation.islatest = 1) typetab ON typetab.recordid = trnrecord.id  AND trnrecord.clientid = typetab.clientid  AND trnrecord.mstorgnhirarchyid = typetab.mstorgnhirarchyid  LEFT JOIN  maprequestorecord ON maprequestorecord.activeflg = 1  AND maprequestorecord.deleteflg = 0  AND maprequestorecord.clientid = trnrecord.clientid  AND maprequestorecord.mstorgnhirarchyid = trnrecord.id  LEFT JOIN  mstrequest ON mstrequest.activeflg = 1  AND mstrequest.deleteflg = 0  AND maprequestorecord.clientid = mstrequest.clientid  AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid  AND maprequestorecord.mstrequestid = mstrequest.id  LEFT JOIN  mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1  AND clisupgroup.deleteflg = 0  AND mstrequest.mstgroupid = clisupgroup.id  LEFT JOIN  supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id  LEFT JOIN  mstclientuser reqclient ON reqclient.activeflag = 1  AND reqclient.deleteflag = 0  AND reqclient.clientid = trnrecord.clientid  AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid  AND mstrequest.userid = reqclient.id WHERE  trnrecord.activeflg = 1  AND trnrecord.deleteflg = 0  AND trnrecord.clientid = ?  AND trnrecord.mstorgnhirarchyid = ? AND  typetab.recorddiffid = ?  AND trnrecord.code = ? AND trnrecord.id in (SELECT distinct recordid FROM maprecordtorecorddifferentiation WHERE recorddiffid in (SELECT id FROM mstrecorddifferentiation WHERE recorddifftypeid=3 AND seqno!=1))"

	logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffid, req.RecordNo)
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsEntity{}
		rows.Scan(&v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.GroupLevelID, &v.GroupLevel, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.StatusSeqNo, &v.Status, &v.ImpactID, &v.Impact, &v.RequestorInfo, &v.Description, &v.UrgencyID, &v.Urgency, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.AssignedGroupLevelID, &v.AssignedGroupLevel, &v.RecordStageID, &v.Duedate)
		values = append(values, v)
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}

func (mdao DbConn) GetIdAgainstRecordNo(Clientid int64, Mstorgnhirarchyid int64, RecordNo string) (int64, int64, error) {
	logger.Log.Println("In side CheckParentOrNot")
	var id int64
	var createdatetime int64

	var query = "SELECT id,createdatetime FROM trnrecord WHERE clientid=? AND mstorgnhirarchyid=? AND code=?"
	rows, err := mdao.DB.Query(query, Clientid, Mstorgnhirarchyid, RecordNo)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("CheckParentOrNot Get Statement Prepare Error", err)
		return id, createdatetime, err
	}
	for rows.Next() {
		err = rows.Scan(&id, &createdatetime)
		logger.Log.Println("CheckParentOrNot rows.next() Error", err)
	}
	return id, createdatetime, nil
}

func (dbc DbConn) GetRecordDetailsByOthers(req *entities.ChildRecordSearchEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsEntity{}
	var sql string
	var params []interface{}
	//var sql1 = "SELECT   COALESCE(trnrecord.requestername,'') requestername,COALESCE(trnrecord.requesteremail,'') requesteremail,COALESCE(trnrecord.requestermobile,'') requestermobile,COALESCE(trnrecord.requesterlocation,'') requesterlocation,  trnrecord.id,  trnrecord.source,  trnrecord.clientid,  trnrecord.mstorgnhirarchyid,  trnrecord.recordtitle,  trnrecord.code,  trnrecord.usergroupid,  mstclientsupportgroup.supportgroupname,  supportgrouplevel.id supportgrouplevelid,  supportgrouplevel.name levelname,  prioritytab.recorddiffid priorityid,  prioritytab.name priority,  mstclientuser.id creatorid,  mstclientuser.loginname createdby,  COALESCE(mstclientuser.vipuser, 'N') vipuser,  COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=trnrecord.clientid AND mstorgnhierarchy.id=trnrecord.mstorgnhirarchyid),  '') createdatetime  ,  statustab.recorddiffid statusid,statustab.seqno statusseqno,  statustab.name status,  COALESCE(impacttab.recorddiffid, 0) impactid,  COALESCE(impacttab.name, '') impact,  trnrecord.requesterinfo,  trnrecord.recorddescription,  COALESCE(urgencytab.recorddiffid, 0) urgencyid,  COALESCE(urgencytab.name, '') urgency,  typetab.recorddiffid typeid,  typetab.name typename,  COALESCE(mstrequest.mstgroupid, 0) assignedgrpid,  COALESCE(mstrequest.mstuserid, 0) assigneduid,  COALESCE(reqclient.loginname, '') assigneduser,  COALESCE(clisupgroup.supportgroupname, '') assignedgrpname,  COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid,  COALESCE(supgrouplvl.name, '') assignedsupgrlvl,  trnrecord.recordstageid,  COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = trnrecord.id),  '') duedate FROM  trnrecord  JOIN  mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1  AND mstclientsupportgroup.deleteflg = 0  AND trnrecord.usergroupid = mstclientsupportgroup.id  JOIN  supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id  JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 4 AND maprecordtorecorddifferentiation.islatest = 1) prioritytab ON prioritytab.recordid = trnrecord.id  AND trnrecord.clientid = prioritytab.clientid  AND trnrecord.mstorgnhirarchyid = prioritytab.mstorgnhirarchyid  JOIN  mstclientuser ON mstclientuser.activeflag = 1  AND mstclientuser.deleteflag = 0  AND mstclientuser.clientid = trnrecord.clientid  AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid  AND trnrecord.userid = mstclientuser.id  JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid,mstrecorddifferentiation.seqno  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 2 AND maprecordtorecorddifferentiation.islatest=1) statustab ON statustab.recordid = trnrecord.id  AND trnrecord.clientid = statustab.clientid  AND trnrecord.mstorgnhirarchyid = statustab.mstorgnhirarchyid  LEFT JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 6 AND maprecordtorecorddifferentiation.islatest = 1) impacttab ON impacttab.recordid = trnrecord.id  AND trnrecord.clientid = impacttab.clientid  AND trnrecord.mstorgnhirarchyid = impacttab.mstorgnhirarchyid  LEFT JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 7 AND maprecordtorecorddifferentiation.islatest = 1) urgencytab ON urgencytab.recordid = trnrecord.id  AND trnrecord.clientid = urgencytab.clientid  AND trnrecord.mstorgnhirarchyid = urgencytab.mstorgnhirarchyid  JOIN  (SELECT   maprecordtorecorddifferentiation.recordid,  mstrecorddifferentiation.name,  maprecordtorecorddifferentiation.clientid,  maprecordtorecorddifferentiation.mstorgnhirarchyid,  maprecordtorecorddifferentiation.recorddifftypeid,  maprecordtorecorddifferentiation.recorddiffid,  maprecordtorecorddifferentiation.recordstageid  FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation ON maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id  JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id  WHERE  mstrecorddifferentiation.activeflg = 1  AND mstrecorddifferentiation.deleteflg = 0  AND mstrecorddifferentiationtype.seqno = 1 AND maprecordtorecorddifferentiation.islatest = 1) typetab ON typetab.recordid = trnrecord.id  AND trnrecord.clientid = typetab.clientid  AND trnrecord.mstorgnhirarchyid = typetab.mstorgnhirarchyid  LEFT JOIN  maprequestorecord ON maprequestorecord.activeflg = 1  AND maprequestorecord.deleteflg = 0  AND maprequestorecord.clientid = trnrecord.clientid  AND maprequestorecord.mstorgnhirarchyid = trnrecord.id  LEFT JOIN  mstrequest ON mstrequest.activeflg = 1  AND mstrequest.deleteflg = 0  AND maprequestorecord.clientid = mstrequest.clientid  AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid  AND maprequestorecord.mstrequestid = mstrequest.id  LEFT JOIN  mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1  AND clisupgroup.deleteflg = 0  AND mstrequest.mstgroupid = clisupgroup.id  LEFT JOIN  supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id  LEFT JOIN  mstclientuser reqclient ON reqclient.activeflag = 1  AND reqclient.deleteflag = 0  AND reqclient.clientid = trnrecord.clientid  AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid  AND mstrequest.userid = reqclient.id"
	//var sql2 = " WHERE  trnrecord.activeflg = 1  AND trnrecord.deleteflg = 0  AND trnrecord.clientid = ?  AND trnrecord.mstorgnhirarchyid = ? AND  typetab.recorddiffid = ?"

	if len(req.RecordNo) > 0 {
		//sql = sql1 + sql2 + " AND trnrecord.code = ?"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,e.name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND a.code =? and e.id=a.userid AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id"

		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.RecordNo)
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)

	}

	if len(req.RequesterID) > 0 {
		//sql = sql1 + " ,mstclientuser e " + sql2 + " AND e.loginname like ? AND (trnrecord.userid=e.id OR trnrecord.originaluserid=e.id) AND trnrecord.id!=2"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,(Select name from mstclientuser where id=a.userid) name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND (a.userid=e.id OR a.originaluserid=e.id) and e.loginname  like ? AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id order by a.id desc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, "%"+req.RequesterID+"%")
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
	}

	if len(req.Requestername) > 0 {
		//sql = sql1 + " ,mstclientuser e " + sql2 + " AND e.name  like ? AND (trnrecord.userid=e.id OR trnrecord.originaluserid=e.id) AND trnrecord.id!=2"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,(Select name from mstclientuser where id=a.userid) name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND (a.userid=e.id OR a.originaluserid=e.id) and e.name  like ? AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id order by a.id desc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, "%"+req.Requestername+"%")
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
	}

	if len(req.Requesterlocation) > 0 {
		//sql = sql1 + " ,mstclientuser e " + sql2 + " AND e.loginname like ? AND (trnrecord.userid=e.id OR trnrecord.originaluserid=e.id) AND trnrecord.id!=2"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,(Select name from mstclientuser where id=a.userid) name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND (a.userid=e.id OR a.originaluserid=e.id) and e.branch like ? AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id order by a.id desc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, "%"+req.Requesterlocation+"%")
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
	}

	if len(req.ShortDescription) > 0 {
		//sql = sql1 + sql2 + " AND trnrecord.id in (SELECT distinct id FROM trnrecord where recordtitle like ?) AND trnrecord.id!=2"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,e.name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND a.recordtitle like ? and e.id=a.userid AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id order by a.id desc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, "%"+req.ShortDescription+"%")
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
	}

	if req.Priority > 0 {
		//sql = sql1 + sql2 + " AND trnrecord.id in (SELECT distinct recordid FROM maprecordtorecorddifferentiation where recorddifftypeid=5 AND recorddiffid=?) AND trnrecord.id!=2"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,e.name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND a.id in (SELECT recordid FROM maprecordtorecorddifferentiation WHERE recorddifftypeid = 5 AND recorddiffid = ? AND clientid = ? AND mstorgnhirarchyid = ? AND islatest=1) and e.id=a.userid AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id order by a.id desc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.Priority)
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
	}

	if len(req.Fromdate) > 0 && len(req.Todate) > 0 {
		//sql = sql1 + sql2 + " AND trnrecord.id in (SELECT distinct id FROM trnrecord where createdatetime BETWEEN round(UNIX_TIMESTAMP(?)) AND round(UNIX_TIMESTAMP(?))) AND trnrecord.id!=2"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,e.name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND a.createdatetime BETWEEN round(UNIX_TIMESTAMP(?)) AND round(UNIX_TIMESTAMP(?)) and e.id=a.userid AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id order by a.createdatetime asc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.Fromdate)
		params = append(params, req.Todate)
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
	}

	if req.CategorylabelID > 0 && req.CategoryID > 0 {
		//	sql = sql1 + sql2 + " AND trnrecord.id in (SELECT distinct recordid FROM maprecordtorecorddifferentiation where recorddifftypeid=? AND recorddiffid=?) AND trnrecord.id!=2"
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,e.name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstclientuser e WHERE b.recordid=a.id AND b.islatest=1 AND a.id in (SELECT recordid FROM maprecordtorecorddifferentiation WHERE recorddifftypeid = ? AND recorddiffid = ? AND clientid = ? AND mstorgnhirarchyid = ? AND islatest=1) and e.id=a.userid AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,811) AND b.recorddiffid=d.id order by a.id desc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.CategorylabelID)
		params = append(params, req.CategoryID)
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)

	}

	if req.GroupID > 0 {
		sql = "SELECT distinct a.id,a.recordstageid,a.code,a.recordtitle,COALESCE((SELECT   DATE_FORMAT(FROM_UNIXTIME(a.createdatetime+zone.utcdiff),'%d-%b-%Y %H:%i:%s')  FROM zone JOIN mstorgnhierarchy ON zone.zone_id=mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid=? AND mstorgnhierarchy.id=?),  '') createdatetime,(Select name from mstclientuser where id=a.userid) name,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=3 AND n.recordid=a.id AND n.islatest=1) status,(SELECT q.name status FROM maprecordtorecorddifferentiation n, mstrecorddifferentiation q where n.recorddiffid = q.id AND n.recorddifftypeid=5 AND n.recordid=a.id AND n.islatest=1) priority,(select distinct mstgroupid from mstrequest m,maprequestorecord n where n.recordid=a.id AND m.id=n.mstrequestid) grpid,COALESCE((SELECT   DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),'%d-%b-%Y %H:%i:%s')  FROM  mstsladue  WHERE  mstsladue.therecordid = a.id),  '') duedate,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=a.id AND recorddifftypeid=2) recordtype,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE parentrecordid=a.id ORDER BY id desc limit 1),'NO') isparent,coalesce((SELECT IF(isattached='Y','Yes','NO') FROM mstparentchildmap WHERE childrecordid=a.id ORDER BY id desc limit 1),'NO') ischild FROM trnrecord a,maprecordtorecorddifferentiation b ,mstrecorddifferentiation d,mstrequest e,maprequestorecord f WHERE e.mstgroupid=? AND e.id=f.mstrequestid AND f.recordid=a.id AND b.recordid=a.id AND b.islatest=1 AND d.clientid=? AND d.mstorgnhirarchyid=? AND d.recorddifftypeid=3 AND d.seqno NOT IN (3,8,11) AND b.recorddiffid=d.id order by a.id desc"
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.GroupID)
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
	}

	logger.Log.Println("Main Query: ", sql)
	logger.Log.Println("Query parameter : ", params)
	rows, err := dbc.DB.Query(sql, params...)
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsEntity{}
		rows.Scan(&v.ID, &v.RecordStageID, &v.Code, &v.Title, &v.CreatedDateTime, &v.CreatedBy, &v.Status, &v.Priority, &v.GroupID, &v.Duedate, &v.RecordTypeID, &v.Isparent, &v.Ischild)
		if req.RecordDiffid == v.RecordTypeID {
			values = append(values, v)
		}

	}

	logger.Log.Println("Results: ", values)
	return values, nil
}
