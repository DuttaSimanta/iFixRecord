package dao

import (
	"database/sql"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"strings"
)

func (dbc DbConn) GetCurrentStageId(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side GetCurrentStageId Dao")
	var stageid int64
	var sql1 = "SELECT MAX(trnrecordstage.id) currentstageid	FROM trnrecordstage	WHERE trnrecordstage.activeflg=1 AND trnrecordstage.deleteflg=0 AND  trnrecordstage.recordid=?"
	//logger.Log.Println("GetCurrentStageId Query: ", sql1)
	//logger.Log.Println("GetCurrentStageId Params: ", req.Recordid)
	err := dbc.DB.QueryRow(sql1, req.Recordid).Scan(&stageid)
	switch err {
	case sql.ErrNoRows:
		stageid = 0
		return stageid, nil
	case nil:
		return stageid, nil
	default:
		logger.Log.Println("GetCurrentStageId Get Statement Prepare Error", err)
		return stageid, err
	}
}

func (dbc DbConn) GetProcessidbyworkingcatid(req *entities.RecordDetailsRequestEntity) (entities.WorkFlowEntity, error) {
	logger.Log.Println("In side GetProcessidbyworkingcatid Dao")
	workflowid := entities.WorkFlowEntity{}
	//var sql1 = "SELECT mstprocessrecordmap.mstprocessid,maprecordtorecorddifferentiation.recorddifftypeid,maprecordtorecorddifferentiation.recorddiffid FROM maprecordtorecorddifferentiation JOIN mstprocessrecordmap ON maprecordtorecorddifferentiation.clientid = mstprocessrecordmap.clientid AND maprecordtorecorddifferentiation.mstorgnhirarchyid = mstprocessrecordmap.mstorgnhirarchyid AND maprecordtorecorddifferentiation.recorddiffid=mstprocessrecordmap.recorddiffid AND maprecordtorecorddifferentiation.recorddifftypeid=mstprocessrecordmap.recorddifftypeid WHERE maprecordtorecorddifferentiation.activeflg = 1 AND maprecordtorecorddifferentiation.deleteflg = 0 AND mstprocessrecordmap.activeflg = 1 AND mstprocessrecordmap.deleteflg = 0 AND maprecordtorecorddifferentiation.isworking = 1 AND maprecordtorecorddifferentiation.clientid = ? AND maprecordtorecorddifferentiation.mstorgnhirarchyid = ? AND maprecordtorecorddifferentiation.recordid = ? AND maprecordtorecorddifferentiation.recordstageid = ?"
	var sql1 = "SELECT mstprocessrecordmap.mstprocessid,maprecordtorecorddifferentiation.recorddifftypeid,maprecordtorecorddifferentiation.recorddiffid FROM maprecordtorecorddifferentiation JOIN mstprocessrecordmap ON maprecordtorecorddifferentiation.clientid = mstprocessrecordmap.clientid AND maprecordtorecorddifferentiation.mstorgnhirarchyid = mstprocessrecordmap.mstorgnhirarchyid AND maprecordtorecorddifferentiation.recorddiffid=mstprocessrecordmap.recorddiffid AND maprecordtorecorddifferentiation.recorddifftypeid=mstprocessrecordmap.recorddifftypeid WHERE maprecordtorecorddifferentiation.activeflg = 1 AND maprecordtorecorddifferentiation.deleteflg = 0 AND mstprocessrecordmap.activeflg = 1 AND mstprocessrecordmap.deleteflg = 0 AND maprecordtorecorddifferentiation.isworking = 1 AND maprecordtorecorddifferentiation.islatest=1 AND maprecordtorecorddifferentiation.clientid = ? AND maprecordtorecorddifferentiation.mstorgnhirarchyid = ? AND maprecordtorecorddifferentiation.recordid = ? " //AND maprecordtorecorddifferentiation.recordstageid = ?
	//logger.Log.Println("GetProcessidbyworkingcatid Query: ", sql1)
	//logger.Log.Println("GetProcessidbyworkingcatid Params: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid)
	err := dbc.DB.QueryRow(sql1, req.Clientid, req.Mstorgnhirarchyid, req.Recordid).Scan(&workflowid.WorkFlowID, &workflowid.CatTypeID, &workflowid.CatID) // req.RecordStageID
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

func (dbc DbConn) GetProcessidbyworkingcatidNew(req *entities.RecordDetailsRequestEntity, OrgnID int64) (entities.WorkFlowEntity, error) {
	logger.Log.Println("In side GetProcessidbyworkingcatid Dao")
	workflowid := entities.WorkFlowEntity{}
	var sql1 = "SELECT mstprocessrecordmap.mstprocessid,maprecordtorecorddifferentiation.recorddifftypeid,maprecordtorecorddifferentiation.recorddiffid FROM maprecordtorecorddifferentiation JOIN mstprocessrecordmap ON maprecordtorecorddifferentiation.clientid = mstprocessrecordmap.clientid AND maprecordtorecorddifferentiation.mstorgnhirarchyid = mstprocessrecordmap.mstorgnhirarchyid AND maprecordtorecorddifferentiation.recorddiffid=mstprocessrecordmap.recorddiffid AND maprecordtorecorddifferentiation.recorddifftypeid=mstprocessrecordmap.recorddifftypeid WHERE maprecordtorecorddifferentiation.activeflg = 1 AND maprecordtorecorddifferentiation.deleteflg = 0 AND mstprocessrecordmap.activeflg = 1 AND mstprocessrecordmap.deleteflg = 0 AND maprecordtorecorddifferentiation.isworking = 1 AND maprecordtorecorddifferentiation.clientid = ? AND maprecordtorecorddifferentiation.mstorgnhirarchyid = ? AND maprecordtorecorddifferentiation.recordid = ? AND maprecordtorecorddifferentiation.recordstageid = ?"
	//	logger.Log.Println("GetProcessidbyworkingcatid Query: ", sql1)
	//	logger.Log.Println("GetProcessidbyworkingcatid Params: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordStageID)
	err := dbc.DB.QueryRow(sql1, req.Clientid, OrgnID, req.Recordid, req.RecordStageID).Scan(&workflowid.WorkFlowID, &workflowid.CatTypeID, &workflowid.CatID)
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

func (dbc DbConn) GetResolRespBreachCount(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side GetResolRespBreachCount Dao")
	var count int64
	var query = "SELECT COUNT(trnreordtracking.id) counter FROM trnreordtracking JOIN mstrecordterms ON trnreordtracking.recordtermid=mstrecordterms.id AND trnreordtracking.clientid=mstrecordterms.clientid AND trnreordtracking.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 AND trnreordtracking.clientid=? AND trnreordtracking.mstorgnhirarchyid=? AND trnreordtracking.recordid=? AND mstrecordterms.seq IN (?)"
	//logger.Log.Println("GetResolRespBreachCount Query: ", query)
	//logger.Log.Println("GetProcessidbyworkingcatid Params: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.TermsSeq)
	err := dbc.DB.QueryRow(query, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.TermsSeq).Scan(&count)
	switch err {
	case sql.ErrNoRows:
		return count, nil
	case nil:
		return count, nil
	default:
		logger.Log.Println("GetResolRespBreachCount Get Statement Prepare Error", err)
		return count, err
	}
}

func (dbc DbConn) GetRecordDetails(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsEntity{}

	//var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=recordtype.recorddiffid) typeseq FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddiffid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddiffid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND trnrecord.id =?"
	//var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=recordtype.recorddiffid) typeseq FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddiffid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddiffid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND trnrecord.id =?"
	//var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, COALESCE(c.recorddiffid, 0) statusid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = c.recorddiffid),'') status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,COALESCE(d.recorddifftypeid, 0) typedifftypeid,COALESCE(d.recorddiffid, 0) typeid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = d.recorddiffid),'') typename,COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=d.recorddiffid) typeseq FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddifftypeid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddifftypeid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN maprecordtorecorddifferentiation c ON trnrecord.id = c.recordid AND c.islatest = 1 AND c.recorddifftypeid = 3 JOIN maprecordtorecorddifferentiation d ON trnrecord.id = d.recordid AND d.islatest = 1 AND d.recorddifftypeid = 2 LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND trnrecord.id =?"
	
	var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, COALESCE(c.recorddiffid, 0) statusid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = c.recorddiffid),'') status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,COALESCE(d.recorddifftypeid, 0) typedifftypeid,COALESCE(d.recorddiffid, 0) typeid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = d.recorddiffid),'') typename,COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=d.recorddiffid) typeseq,(SELECT IF(deleteflag = 0, 'N', 'Y') FROM mstclientuser WHERE id = trnrecord.userid) isuserdelete,(SELECT IF(deleteflag = 0, 'N', 'Y') FROM mstclientuser WHERE id = trnrecord.originaluserid) isoriginaluseriddelete FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddifftypeid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddifftypeid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN maprecordtorecorddifferentiation c ON trnrecord.id = c.recordid AND c.islatest = 1 AND c.recorddifftypeid = 3 JOIN maprecordtorecorddifferentiation d ON trnrecord.id = d.recordid AND d.islatest = 1 AND d.recorddifftypeid = 2 LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND trnrecord.id = ?"
	//logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.Recordid)
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsEntity{}
		rows.Scan(&v.OriginalUserID, &v.OrgRequestorName, &v.OrgRequestorEmail, &v.OrgRequestorMobile, &v.OrgRequestorLocation, &v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.GroupLevelID, &v.GroupLevel, &v.PriorityTypeID, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.Status, &v.ImpactID, &v.Impact, &v.RequestorInfo, &v.Description, &v.UrgencyID, &v.Urgency, &v.RecordTypeDiffTypeID, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.AssignedGroupLevelID, &v.AssignedGroupLevel, &v.RecordStageID, &v.Duedate, &v.TypeSeqNo,&v.Isuserdelete, &v.Isoriginaluseriddelete)
		values = append(values, v)
	}
	//logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) RemoveChildRecord(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side RemoveChildRecord Dao")
	values := []entities.RecordDetailsEntity{}
	var query = "UPDATE mstparentchildmap SET deleteflg=1, activeflg=0,isattached='N' WHERE  clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND parentrecordid=? "
	var params []interface{}

	params = append(params, req.Clientid)
	params = append(params, req.Mstorgnhirarchyid)
	params = append(params, req.RecordDiffTypeid)
	params = append(params, req.RecordDiffid)
	params = append(params, req.ParentID)
	if len(req.ChildIDS) > 0 {
		var subquery []string
		for _, v := range req.ChildIDS {
			subquery = append(subquery, " ? ")
			params = append(params, v)
		}
		query = query + " AND childrecordid IN ( " + strings.Join(subquery, ",") + " ) "
	}

	//logger.Log.Println("Main Query: ", query)
	//logger.Log.Println("Main Query: ", params)
	rows, err := dbc.DB.Query(query, params...)
	if err != nil {
		logger.Log.Println("RemoveChildRecord Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	return values, nil
}
func (dbc DbConn) CheckDuplicateChildRecord(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side CheckDuplicateChildRecord Dao")
	var value int64
	var query = "SELECT COUNT(id) total FROM mstparentchildmap WHERE deleteflg=0 AND activeflg=1 AND isattached='Y' AND clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? "
	var params []interface{}

	params = append(params, req.Clientid)
	params = append(params, req.Mstorgnhirarchyid)
	params = append(params, req.RecordDiffTypeid)
	params = append(params, req.RecordDiffid)
	//params = append(params, req.ParentID)
	if len(req.ChildIDS) > 0 {
		var subquery []string
		for _, v := range req.ChildIDS {
			subquery = append(subquery, " ? ")
			params = append(params, v)
		}
		query = query + " AND childrecordid IN ( " + strings.Join(subquery, ",") + " ) "
	}

	//	logger.Log.Println("Main Query: ", query)
	//	logger.Log.Println("Main Query: ", params)
	rows, err := dbc.DB.Query(query, params...)
	if err != nil {
		logger.Log.Println("CheckDuplicateChildRecord Get Statement Prepare Error", err)
		return value, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
	}
	return value, nil
}

func (dbc DbConn) CheckDuplicateParentRecord(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side CheckDuplicateChildRecord Dao")
	var value int64
	var query = "SELECT COUNT(id) total FROM mstparentchildmap WHERE deleteflg=0 AND activeflg=1 AND clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND parentrecordid=? AND isattached='Y'"
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
		logger.Log.Println("CheckDuplicateChildRecord Get Statement Prepare Error", err)
		return value, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
	}
	return value, nil
}

func (dbc DbConn) CheckChildAsParentRecord(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side CheckChildAsParentRecord Dao")
	var value int64
	var query = "SELECT COUNT(id) total FROM mstparentchildmap WHERE deleteflg=0 AND activeflg=1 AND isattached='Y' AND clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? "
	var params []interface{}

	params = append(params, req.Clientid)
	params = append(params, req.Mstorgnhirarchyid)
	params = append(params, req.RecordDiffTypeid)
	params = append(params, req.RecordDiffid)
	//params = append(params, req.ParentID)
	if len(req.ChildIDS) > 0 {
		var subquery []string
		for _, v := range req.ChildIDS {
			subquery = append(subquery, " ? ")
			params = append(params, v)
		}
		query = query + " AND parentrecordid IN ( " + strings.Join(subquery, ",") + " ) "
	}

	//	logger.Log.Println("Main Query: ", query)
	//	logger.Log.Println("Main Query: ", params)
	rows, err := dbc.DB.Query(query, params...)
	if err != nil {
		logger.Log.Println("CheckChildAsParentRecord Get Statement Prepare Error", err)
		return value, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&value)
	}
	return value, nil
}

func (dbc DbConn) SaveChildRecord(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side SaveChildRecord Dao")
	values := []entities.RecordDetailsEntity{}
	if len(req.ChildIDS) > 0 {
		var query = "INSERT INTO mstparentchildmap(clientid, mstorgnhirarchyid, recorddifftypeid, recorddiffid, parentrecordid, childrecordid,deleteflg, activeflg) VALUES "
		var subquery []string
		var params []interface{}
		for _, v := range req.ChildIDS {
			subquery = append(subquery, " (?, ?, ?, ?, ?, ?, 0,1) ")
			params = append(params, req.Clientid)
			params = append(params, req.Mstorgnhirarchyid)
			params = append(params, req.RecordDiffTypeid)
			params = append(params, req.RecordDiffid)
			params = append(params, req.ParentID)
			params = append(params, v)
		}
		query = query + strings.Join(subquery, ",")
		//		logger.Log.Println("Main Query: ", query)
		//		logger.Log.Println("Main Query: ", params)
		rows, err := dbc.DB.Query(query, params...)
		if err != nil {
			logger.Log.Println("SaveChildRecord Get Statement Prepare Error", err)
			return values, err
		}
		defer rows.Close()
	}
	return values, nil
}

func (dbc DbConn) GetChildRecordsBYParentID(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsEntity{}

	//var sql = "SELECT COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, mstclientsupportgroup.supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(recordimpact.recorddiffid, 0) impactid, COALESCE(recordimpact.recorddiffname, '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(recordurgency.recorddiffid, 0) urgencyid, COALESCE(recordurgency.recorddiffname, '') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE(clisupgroup.supportgroupname, '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.id JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN recordimpact ON recordimpact.recordid = trnrecord.id AND trnrecord.clientid = recordimpact.clientid AND trnrecord.mstorgnhirarchyid = recordimpact.mstorgnhirarchyid LEFT JOIN recordurgency ON recordurgency.recordid = trnrecord.id AND trnrecord.clientid = recordurgency.clientid AND trnrecord.mstorgnhirarchyid = recordurgency.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.id LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid = ? AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached = 'Y'"
	//var sql = "SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(recordimpact.recorddiffid, 0) impactid, COALESCE(recordimpact.recorddiffname, '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(recordurgency.recorddiffid, 0) urgencyid, COALESCE(recordurgency.recorddiffname, '') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN recordimpact ON recordimpact.recordid = trnrecord.id AND trnrecord.clientid = recordimpact.clientid AND trnrecord.mstorgnhirarchyid = recordimpact.mstorgnhirarchyid LEFT JOIN recordurgency ON recordurgency.recordid = trnrecord.id AND trnrecord.clientid = recordurgency.clientid AND trnrecord.mstorgnhirarchyid = recordurgency.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid = ? AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached = 'Y'"
	//var sql = "SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status,COALESCE(a.recorddiffid, 0) impactid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid  AND a.islatest = 1 AND a.recorddiffid = 7 LEFT JOIN  maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest = 1  AND b.recorddiffid = 8 JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid = ? AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached = 'Y'"
	//var sql = "SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status,COALESCE(a.recorddiffid, 0) impactid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid  AND a.islatest = 1 AND a.recorddiffid = 7 LEFT JOIN  maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest = 1  AND b.recorddiffid = 8 JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid = ? AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached IN ('Y','N')"
	//var sql = "SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status,COALESCE(a.recorddiffid, 0) impactid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename,COALESCE(recordfulldetails.assignedgroupid, 0) assignedgrpid,COALESCE(recordfulldetails.assigneduserid, 0) assigneduid,COALESCE(recordfulldetails.assigneduser, '') assigneduser,COALESCE(recordfulldetails.assignedgroup,'') assignedgrpname,trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid  AND a.islatest = 1 AND a.recorddiffid = 7 LEFT JOIN  maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest = 1  AND b.recorddiffid = 8 JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN recordfulldetails ON recordfulldetails.recordid=trnrecord.id AND trnrecord.clientid = recordfulldetails.clientid AND trnrecord.mstorgnhirarchyid = recordfulldetails.mstorgnhirarchyid JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid = ? AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached IN ('Y')"
	var sql ="SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status,COALESCE(a.recorddiffid, 0) impactid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename,COALESCE(recordfulldetails.assignedgroupid, 0) assignedgrpid,COALESCE(recordfulldetails.assigneduserid, 0) assigneduid,COALESCE(recordfulldetails.assigneduser, '') assigneduser,COALESCE(recordfulldetails.assignedgroup,'') assignedgrpname,trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid  AND a.islatest = 1 AND a.recorddiffid = 7 LEFT JOIN  maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest = 1  AND b.recorddiffid = 8 JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN recordfulldetails ON recordfulldetails.recordid=trnrecord.id AND trnrecord.clientid = recordfulldetails.clientid AND trnrecord.mstorgnhirarchyid = recordfulldetails.mstorgnhirarchyid JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid = ? AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached IN ('Y')"
	//logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffid, req.ParentID)
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsEntity{}
		err := rows.Scan(&v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.Status, &v.ImpactID, &v.Impact, &v.RequestorInfo, &v.Description, &v.UrgencyID, &v.Urgency, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.RecordStageID, &v.Duedate) //, &v.AssignedGroupLevelID, &v.AssignedGroupLevel
		//&v.GroupLevelID, &v.GroupLevel
		if err != nil {
			logger.Log.Println("GetRecordDetails Scan Error", err)
			return values, err
		}
		values = append(values, v)
	}
	//	logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) GetStaskChildRecordsBYParentID(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsEntity{}

	//var sql = "SELECT COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, mstclientsupportgroup.supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(recordimpact.recorddiffid, 0) impactid, COALESCE(recordimpact.recorddiffname, '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(recordurgency.recorddiffid, 0) urgencyid, COALESCE(recordurgency.recorddiffname, '') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE(clisupgroup.supportgroupname, '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.id JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN recordimpact ON recordimpact.recordid = trnrecord.id AND trnrecord.clientid = recordimpact.clientid AND trnrecord.mstorgnhirarchyid = recordimpact.mstorgnhirarchyid LEFT JOIN recordurgency ON recordurgency.recordid = trnrecord.id AND trnrecord.clientid = recordurgency.clientid AND trnrecord.mstorgnhirarchyid = recordurgency.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.id LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid = ? AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached = 'Y'"
	//var sql = "SELECT COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(recordimpact.recorddiffid, 0) impactid, COALESCE(recordimpact.recorddiffname, '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(recordurgency.recorddiffid, 0) urgencyid, COALESCE(recordurgency.recorddiffname, '') urgency, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN recordimpact ON recordimpact.recordid = trnrecord.id AND trnrecord.clientid = recordimpact.clientid AND trnrecord.mstorgnhirarchyid = recordimpact.mstorgnhirarchyid LEFT JOIN recordurgency ON recordurgency.recordid = trnrecord.id AND trnrecord.clientid = recordurgency.clientid AND trnrecord.mstorgnhirarchyid = recordurgency.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=2 AND seqno in (2,3) AND deleteflg=0 AND activeflg=1) AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached = 'Y'"
	//var sql = "SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status,COALESCE(a.recorddiffid, 0) impactid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid), '') impact,trnrecord.requesterinfo,trnrecord.recorddescription,COALESCE(b.recorddiffid, 0) urgencyid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') urgency,recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest = 1 AND a.recorddiffid = 7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid  AND b.islatest = 1 AND b.recorddiffid = 8 JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=2 AND seqno in (2,3,4,5) AND deleteflg=0 AND activeflg=1) AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached  IN ('Y','N')" //AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid
	//var sql = "SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status,COALESCE(a.recorddiffid, 0) impactid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid), '') impact,trnrecord.requesterinfo,trnrecord.recorddescription,COALESCE(b.recorddiffid, 0) urgencyid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') urgency,recordtype.recorddiffid typeid, recordtype.recorddiffname typename,COALESCE(recordfulldetails.assignedgroupid, 0) assignedgrpid,COALESCE(recordfulldetails.assigneduserid, 0) assigneduid,COALESCE(recordfulldetails.assigneduser, '') assigneduser,COALESCE(recordfulldetails.assignedgroup,'') assignedgrpname,trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest = 1 AND a.recorddiffid = 7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid  AND b.islatest = 1 AND b.recorddiffid = 8 JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN recordfulldetails ON recordfulldetails.recordid=trnrecord.id AND trnrecord.clientid = recordfulldetails.clientid AND trnrecord.mstorgnhirarchyid = recordfulldetails.mstorgnhirarchyid JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=2 AND seqno in (2,3,4,5) AND deleteflg=0 AND activeflg=1) AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached  IN ('Y','N')"
	
	var sql = "SELECT distinct COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, recordstatus.recorddiffid statusid, recordstatus.recorddiffname status,COALESCE(a.recorddiffid, 0) impactid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid), '') impact,trnrecord.requesterinfo,trnrecord.recorddescription,COALESCE(b.recorddiffid, 0) urgencyid,COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id = a.recorddiffid),'') urgency,recordtype.recorddiffid typeid, recordtype.recorddiffname typename,COALESCE(recordfulldetails.assignedgroupid, 0) assignedgrpid,COALESCE(recordfulldetails.assigneduserid, 0) assigneduid,COALESCE(recordfulldetails.assigneduser, '') assigneduser,COALESCE(recordfulldetails.assignedgroup,'') assignedgrpname,trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%b-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid=trnrecord.id AND trnrecord.clientid=recordpriority.clientid AND trnrecord.mstorgnhirarchyid=recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid= trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest = 1 AND a.recorddiffid = 7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid  AND b.islatest = 1 AND b.recorddiffid = 8 JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND maprequestorecord.recordid=trnrecord.id LEFT JOIN recordfulldetails ON recordfulldetails.recordid=trnrecord.id AND trnrecord.clientid = recordfulldetails.clientid AND trnrecord.mstorgnhirarchyid = recordfulldetails.mstorgnhirarchyid JOIN mstparentchildmap ON mstparentchildmap.childrecordid = trnrecord.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND mstparentchildmap.activeflg = 1 AND mstparentchildmap.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND recordtype.recorddiffid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=2 AND seqno in (2,3,4,5) AND deleteflg=0 AND activeflg=1) AND mstparentchildmap.parentrecordid = ? AND mstparentchildmap.isattached  IN ('Y','N')"
	//logger.Log.Println("Main Query: ", sql)
	//, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.Clientid, req.Mstorgnhirarchyid, req.ParentID)
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsEntity{}
		err := rows.Scan(&v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.Status, &v.ImpactID, &v.Impact, &v.RequestorInfo, &v.Description, &v.UrgencyID, &v.Urgency, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.RecordStageID, &v.Duedate) //&v.AssignedGroupLevelID, &v.AssignedGroupLevel
		//, &v.GroupLevelID, &v.GroupLevel
		if err != nil {
			logger.Log.Println("GetRecordDetails Scan Error", err)
			return values, err
		}
		values = append(values, v)
	}
	//	logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) GetRecordDetailsByNoForlinkrecord(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsEntity{}

	//var sql = "SELECT 	trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, mstclientsupportgroup.supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(recordimpact.recorddiffid, 0) impactid, COALESCE(recordimpact.recorddiffname, '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(recordurgency.recorddiffid, 0) urgencyid, COALESCE(recordurgency.recorddiffname, '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE(clisupgroup.supportgroupname, '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.id JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN recordimpact ON recordimpact.recordid = trnrecord.id AND trnrecord.clientid = recordimpact.clientid AND trnrecord.mstorgnhirarchyid = recordimpact.mstorgnhirarchyid LEFT JOIN recordurgency ON recordurgency.recordid = trnrecord.id AND trnrecord.clientid = recordurgency.clientid AND trnrecord.mstorgnhirarchyid = recordurgency.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.id LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND originalclient.clientid = trnrecord.clientid AND originalclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.originaluserid = originalclient.id WHERE  trnrecord.activeflg = 1  AND trnrecord.deleteflg = 0  AND trnrecord.clientid = ?  AND trnrecord.mstorgnhirarchyid = ? AND trnrecord.code = ?"

	//var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=recordtype.recorddiffid) typeseq FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddiffid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddiffid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ?  AND trnrecord.code =?" //AND trnrecord.mstorgnhirarchyid = ?

	var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=recordtype.recorddiffid) typeseq FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddiffid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddiffid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.mstorgnhirarchyid = ? AND trnrecord.code =?"

	//logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordNo) //req.Mstorgnhirarchyid,
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsEntity{}
		rows.Scan(&v.OriginalUserID, &v.OrgRequestorName, &v.OrgRequestorEmail, &v.OrgRequestorMobile, &v.OrgRequestorLocation, &v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.GroupLevelID, &v.GroupLevel, &v.PriorityTypeID, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.Status, &v.ImpactID, &v.Impact, &v.RequestorInfo, &v.Description, &v.UrgencyID, &v.Urgency, &v.RecordTypeDiffTypeID, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.AssignedGroupLevelID, &v.AssignedGroupLevel, &v.RecordStageID, &v.Duedate, &v.TypeSeqNo)
		//rows.Scan(&v.OriginalUserID, &v.OrgRequestorName, &v.OrgRequestorEmail, &v.OrgRequestorMobile, &v.OrgRequestorLocation, &v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.GroupLevelID, &v.GroupLevel, &v.PriorityTypeID, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.Status, &v.RequestorInfo, &v.Description, &v.RecordTypeDiffTypeID, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.AssignedGroupLevelID, &v.AssignedGroupLevel, &v.RecordStageID, &v.Duedate)
		values = append(values, v)
	}
	return values, nil
}

func (dbc DbConn) GetRecordDetailsByNo(req *entities.RecordDetailsRequestEntity) ([]entities.RecordDetailsEntity, error) {
	logger.Log.Println("In side GetRecordDetails Dao")
	values := []entities.RecordDetailsEntity{}

	//var sql = "SELECT 	trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, mstclientsupportgroup.supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(recordimpact.recorddiffid, 0) impactid, COALESCE(recordimpact.recorddiffname, '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(recordurgency.recorddiffid, 0) urgencyid, COALESCE(recordurgency.recorddiffname, '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE(clisupgroup.supportgroupname, '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate FROM trnrecord JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.id JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid LEFT JOIN recordimpact ON recordimpact.recordid = trnrecord.id AND trnrecord.clientid = recordimpact.clientid AND trnrecord.mstorgnhirarchyid = recordimpact.mstorgnhirarchyid LEFT JOIN recordurgency ON recordurgency.recordid = trnrecord.id AND trnrecord.clientid = recordurgency.clientid AND trnrecord.mstorgnhirarchyid = recordurgency.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.id LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND reqclient.clientid = trnrecord.clientid AND reqclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND originalclient.clientid = trnrecord.clientid AND originalclient.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.originaluserid = originalclient.id WHERE  trnrecord.activeflg = 1  AND trnrecord.deleteflg = 0  AND trnrecord.clientid = ?  AND trnrecord.mstorgnhirarchyid = ? AND trnrecord.code = ?"

	//var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=recordtype.recorddiffid) typeseq FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddiffid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddiffid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND mstclientuser.mstorgnhirarchyid = trnrecord.mstorgnhirarchyid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ?  AND trnrecord.code =?" //AND trnrecord.mstorgnhirarchyid = ?

//	var sql = "SELECT 	distinct trnrecord.originaluserid,COALESCE(originalclient.name, '') orgrequestername, COALESCE(originalclient.useremail, '') orgrequesteremail, COALESCE(originalclient.usermobileno, '') orgrequestermobile, COALESCE(originalclient.city, '') orgrequesterlocation, COALESCE(trnrecord.requestername, '') requestername, COALESCE(trnrecord.requesteremail, '') requesteremail, COALESCE(trnrecord.requestermobile, '') requestermobile, COALESCE(trnrecord.requesterlocation, '') requesterlocation, trnrecord.id, trnrecord.source source, trnrecord.clientid, trnrecord.mstorgnhirarchyid, trnrecord.recordtitle, trnrecord.code, trnrecord.usergroupid, (SELECT name FROM mstsupportgrp WHERE id = mstclientsupportgroup.grpid) supportgroupname, supportgrouplevel.id supportgrouplevelid, supportgrouplevel.name levelname, recordpriority.recorddifftypeid prioritytypeid, recordpriority.recorddiffid priorityid, recordpriority.recorddiffname priority, mstclientuser.id creatorid, mstclientuser.loginname createdby, COALESCE(mstclientuser.vipuser, 'N') vipuser, COALESCE((SELECT DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff), '%d-%b-%Y %H:%i:%s') FROM zone JOIN mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid WHERE mstorgnhierarchy.clientid = trnrecord.clientid AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid), '') createdatetime, 	recordstatus.recorddiffid statusid, recordstatus.recorddiffname status, COALESCE(a.recorddiffid, 0) impactid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') impact, trnrecord.requesterinfo, trnrecord.recorddescription, COALESCE(b.recorddiffid, 0) urgencyid, COALESCE((SELECT name FROM mstrecorddifferentiation WHERE id=a.recorddiffid), '') urgency,recordtype.recorddifftypeid typedifftypeid, recordtype.recorddiffid typeid, recordtype.recorddiffname typename, COALESCE(mstrequest.mstgroupid, 0) assignedgrpid, COALESCE(mstrequest.mstuserid, 0) assigneduid, COALESCE(reqclient.loginname, '') assigneduser, COALESCE((SELECT name FROM mstsupportgrp WHERE id = clisupgroup.grpid), '') assignedgrpname, COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid, COALESCE(supgrouplvl.name, '') assignedsupgrlvl, trnrecord.recordstageid, COALESCE((SELECT DATE_FORMAT(MAX(mstsladue.duedatetimeresolution), '%d-%m-%Y %H:%i:%s') FROM mstsladue WHERE mstsladue.therecordid = trnrecord.id), '') duedate,(SELECT seqno FROM mstrecorddifferentiation WHERE id=recordtype.recorddiffid) typeseq FROM trnrecord  LEFT JOIN maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid AND a.islatest =1 AND a.recorddiffid =7 LEFT JOIN maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid AND b.islatest =1 AND b.recorddiffid =8 JOIN mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1 AND mstclientsupportgroup.deleteflg = 0 AND trnrecord.usergroupid = mstclientsupportgroup.grpid JOIN supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id JOIN recordpriority ON recordpriority.recordid = trnrecord.id AND trnrecord.clientid = recordpriority.clientid AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid JOIN mstclientuser ON mstclientuser.activeflag = 1 AND mstclientuser.deleteflag = 0 AND mstclientuser.clientid = trnrecord.clientid AND trnrecord.userid = mstclientuser.id JOIN recordstatus ON recordstatus.recordid = trnrecord.id AND trnrecord.clientid = recordstatus.clientid AND trnrecord.mstorgnhirarchyid = recordstatus.mstorgnhirarchyid JOIN recordtype ON recordtype.recordid = trnrecord.id AND trnrecord.clientid = recordtype.clientid AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid LEFT JOIN maprequestorecord ON maprequestorecord.activeflg = 1 AND maprequestorecord.deleteflg = 0 AND maprequestorecord.clientid = trnrecord.clientid AND maprequestorecord.recordid = trnrecord.id LEFT JOIN mstrequest ON mstrequest.activeflg = 1 AND mstrequest.deleteflg = 0 AND maprequestorecord.clientid = mstrequest.clientid AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid AND maprequestorecord.mstrequestid = mstrequest.id LEFT JOIN mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1 AND clisupgroup.deleteflg = 0 AND mstrequest.mstgroupid = clisupgroup.grpid LEFT JOIN supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id LEFT JOIN mstclientuser reqclient ON reqclient.activeflag = 1 AND reqclient.deleteflag = 0 AND mstrequest.userid = reqclient.id LEFT JOIN mstclientuser originalclient ON originalclient.activeflag = 1 AND originalclient.deleteflag = 0 AND trnrecord.originaluserid = originalclient.id WHERE trnrecord.activeflg = 1 AND trnrecord.deleteflg = 0 AND trnrecord.clientid = ? AND trnrecord.code =?"
var sql = `SELECT DISTINCT
    trnrecord.originaluserid,
    COALESCE(originalclient.name, '') orgrequestername,
    COALESCE(originalclient.useremail, '') orgrequesteremail,
    COALESCE(originalclient.usermobileno, '') orgrequestermobile,
    COALESCE(originalclient.city, '') orgrequesterlocation,
    COALESCE(trnrecord.requestername, '') requestername,
    COALESCE(trnrecord.requesteremail, '') requesteremail,
    COALESCE(trnrecord.requestermobile, '') requestermobile,
    COALESCE(trnrecord.requesterlocation, '') requesterlocation,
    trnrecord.id,
    trnrecord.source source,
    trnrecord.clientid,
    trnrecord.mstorgnhirarchyid,
    trnrecord.recordtitle,
    trnrecord.code,
    trnrecord.usergroupid,
    (SELECT 
            name
        FROM
            mstsupportgrp
        WHERE
            id = mstclientsupportgroup.grpid) supportgroupname,
    supportgrouplevel.id supportgrouplevelid,
    supportgrouplevel.name levelname,
    recordpriority.recorddifftypeid prioritytypeid,
    recordpriority.recorddiffid priorityid,
    recordpriority.recorddiffname priority,
    mstclientuser.id creatorid,
    mstclientuser.loginname createdby,
    COALESCE(mstclientuser.vipuser, 'N') vipuser,
    COALESCE((SELECT 
                    DATE_FORMAT(FROM_UNIXTIME(trnrecord.createdatetime + zone.utcdiff),
                                '%d-%b-%Y %H:%i:%s')
                FROM
                    zone
                        JOIN
                    mstorgnhierarchy ON zone.zone_id = mstorgnhierarchy.timezoneid
                WHERE
                    mstorgnhierarchy.clientid = trnrecord.clientid
                        AND mstorgnhierarchy.id = trnrecord.mstorgnhirarchyid),
            '') createdatetime,
    COALESCE(c.recorddiffid, 0) statusid,
    COALESCE((SELECT 
                    name
                FROM
                    mstrecorddifferentiation
                WHERE
                    id = c.recorddiffid),
            '') status,
    COALESCE(a.recorddiffid, 0) impactid,
    COALESCE((SELECT 
                    name
                FROM
                    mstrecorddifferentiation
                WHERE
                    id = a.recorddiffid),
            '') impact,
    trnrecord.requesterinfo,
    trnrecord.recorddescription,
    COALESCE(b.recorddiffid, 0) urgencyid,
    COALESCE((SELECT 
                    name
                FROM
                    mstrecorddifferentiation
                WHERE
                    id = a.recorddiffid),
            '') urgency,
    recordtype.recorddifftypeid typedifftypeid,
    recordtype.recorddiffid typeid,
    recordtype.recorddiffname typename,
    COALESCE(mstrequest.mstgroupid, 0) assignedgrpid,
    COALESCE(mstrequest.mstuserid, 0) assigneduid,
    COALESCE(reqclient.loginname, '') assigneduser,
    COALESCE((SELECT 
                    name
                FROM
                    mstsupportgrp
                WHERE
                    id = clisupgroup.grpid),
            '') assignedgrpname,
    COALESCE(clisupgroup.supportgrouplevelid, 0) assignedgrplvlid,
    COALESCE(supgrouplvl.name, '') assignedsupgrlvl,
    trnrecord.recordstageid,
    COALESCE((SELECT 
                    DATE_FORMAT(MAX(mstsladue.duedatetimeresolution),
                                '%d-%m-%Y %H:%i:%s')
                FROM
                    mstsladue
                WHERE
                    mstsladue.therecordid = trnrecord.id),
            '') duedate,
    (SELECT 
            seqno
        FROM
            mstrecorddifferentiation
        WHERE
            id = recordtype.recorddiffid) typeseq,
      (SELECT 
				IF(deleteflag = 0, 'N', 'Y')
			FROM
				mstclientuser
			WHERE
				id = trnrecord.userid) isuserdelete,
		(SELECT 
				IF(deleteflag = 0, 'N', 'Y')
			FROM
				mstclientuser
			WHERE
				id = trnrecord.originaluserid) isoriginaluseriddelete
FROM
    trnrecord
        LEFT JOIN
    maprecordtorecorddifferentiation a ON trnrecord.id = a.recordid
        AND a.islatest = 1
        AND a.recorddifftypeid = 7
        LEFT JOIN
    maprecordtorecorddifferentiation b ON trnrecord.id = b.recordid
        AND b.islatest = 1
        AND b.recorddifftypeid = 8
        JOIN
    mstclientsupportgroup ON mstclientsupportgroup.activeflg = 1
        AND mstclientsupportgroup.deleteflg = 0
        AND trnrecord.usergroupid = mstclientsupportgroup.grpid
        JOIN
    supportgrouplevel ON mstclientsupportgroup.supportgrouplevelid = supportgrouplevel.id
        JOIN
    recordpriority ON recordpriority.recordid = trnrecord.id
        AND trnrecord.clientid = recordpriority.clientid
        AND trnrecord.mstorgnhirarchyid = recordpriority.mstorgnhirarchyid
        JOIN
    mstclientuser ON mstclientuser.activeflag = 1
        AND mstclientuser.clientid = trnrecord.clientid
        AND trnrecord.userid = mstclientuser.id
        JOIN
    maprecordtorecorddifferentiation c ON trnrecord.id = c.recordid
        AND c.islatest = 1
        AND c.recorddifftypeid = 3
        JOIN
    recordtype ON recordtype.recordid = trnrecord.id
        AND trnrecord.clientid = recordtype.clientid
        AND trnrecord.mstorgnhirarchyid = recordtype.mstorgnhirarchyid
        LEFT JOIN
    maprequestorecord ON maprequestorecord.activeflg = 1
        AND maprequestorecord.deleteflg = 0
        AND maprequestorecord.clientid = trnrecord.clientid
        AND maprequestorecord.recordid = trnrecord.id
        LEFT JOIN
    mstrequest ON mstrequest.activeflg = 1
        AND mstrequest.deleteflg = 0
        AND maprequestorecord.clientid = mstrequest.clientid
        AND mstrequest.mstorgnhirarchyid = maprequestorecord.mstorgnhirarchyid
        AND maprequestorecord.mstrequestid = mstrequest.id
        LEFT JOIN
    mstclientsupportgroup clisupgroup ON clisupgroup.activeflg = 1
        AND clisupgroup.deleteflg = 0
        AND mstrequest.mstgroupid = clisupgroup.grpid
        LEFT JOIN
    supportgrouplevel supgrouplvl ON clisupgroup.supportgrouplevelid = supgrouplvl.id
        LEFT JOIN
    mstclientuser reqclient ON reqclient.activeflag = 1
        AND reqclient.deleteflag = 0
        AND mstrequest.userid = reqclient.id
        LEFT JOIN
    mstclientuser originalclient ON originalclient.activeflag = 1
        AND trnrecord.originaluserid = originalclient.id
WHERE
    trnrecord.activeflg = 1
        AND trnrecord.deleteflg = 0
        AND trnrecord.clientid = ?
        AND trnrecord.code = ?;`
	//logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql, req.Clientid, req.RecordNo) //req.Mstorgnhirarchyid,
	if err != nil {
		logger.Log.Println("GetRecordDetails Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		v := entities.RecordDetailsEntity{}
		rows.Scan(&v.OriginalUserID, &v.OrgRequestorName, &v.OrgRequestorEmail, &v.OrgRequestorMobile, &v.OrgRequestorLocation, &v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.GroupLevelID, &v.GroupLevel, &v.PriorityTypeID, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.Status, &v.ImpactID, &v.Impact, &v.RequestorInfo, &v.Description, &v.UrgencyID, &v.Urgency, &v.RecordTypeDiffTypeID, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.AssignedGroupLevelID, &v.AssignedGroupLevel, &v.RecordStageID, &v.Duedate, &v.TypeSeqNo,&v.Isuserdelete, &v.Isoriginaluseriddelete)
		//rows.Scan(&v.OriginalUserID, &v.OrgRequestorName, &v.OrgRequestorEmail, &v.OrgRequestorMobile, &v.OrgRequestorLocation, &v.RequestorName, &v.RequestorEmail, &v.RequestorMobile, &v.RequestorLocation, &v.ID, &v.SourceType, &v.Clientid, &v.Mstorgnhirarchyid, &v.Title, &v.Code, &v.GroupID, &v.Group, &v.GroupLevelID, &v.GroupLevel, &v.PriorityTypeID, &v.PriorityID, &v.Priority, &v.CreatorID, &v.CreatedBy, &v.Vipuser, &v.CreatedDateTime, &v.StatusID, &v.Status, &v.RequestorInfo, &v.Description, &v.RecordTypeDiffTypeID, &v.RecordTypeID, &v.RecordType, &v.AssignedGroupID, &v.AssigneeID, &v.Assignee, &v.AssignedGroup, &v.AssignedGroupLevelID, &v.AssignedGroupLevel, &v.RecordStageID, &v.Duedate)
		values = append(values, v)
	}
	return values, nil
}

func (dbc DbConn) GetAllTypeWiseCategories(req *entities.RecordDetailsRequestEntity) ([]entities.RawDiffEntity, error) {
	logger.Log.Println("In side GetAllRecordcategories")
	var sql = "SELECT mstrecorddifferentiationtype.id typeid, mstrecorddifferentiationtype.typename typename, mstrecorddifferentiationtype.seqno typeseq, mstrecorddifferentiation.id, mstrecorddifferentiation.name, mstrecorddifferentiation.seqno, IF(maprecordtorecorddifferentiation.id IS NULL, 0, 1) selected,mstrecorddifferentiation.parentid FROM mstrecordtype JOIN mstrecorddifferentiation ON mstrecorddifferentiation.id = mstrecordtype.torecorddiffid JOIN mstrecorddifferentiationtype ON mstrecorddifferentiationtype.id = mstrecordtype.torecorddifftypeid AND mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id LEFT JOIN maprecordtorecorddifferentiation ON maprecordtorecorddifferentiation.activeflg = 1 AND maprecordtorecorddifferentiation.deleteflg = 0 AND maprecordtorecorddifferentiation.clientid = ? AND maprecordtorecorddifferentiation.mstorgnhirarchyid = ? AND maprecordtorecorddifferentiation.recorddiffid = mstrecorddifferentiation.id AND maprecordtorecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id AND maprecordtorecorddifferentiation.recordid = ? AND maprecordtorecorddifferentiation.islatest=1 WHERE mstrecordtype.clientid = ? AND mstrecordtype.mstorgnhirarchyid = ? AND mstrecordtype.activeflg = 1 AND mstrecordtype.deleteflg = 0 AND mstrecordtype.fromrecorddifftypeid = ? AND mstrecordtype.fromrecorddiffid = ? AND mstrecorddifferentiationtype.activeflg = 1 AND mstrecorddifferentiationtype.deleteflg = 0 AND mstrecorddifferentiationtype.parentid = 1 AND mstrecorddifferentiation.activeflg = 1 AND mstrecorddifferentiation.deleteflg = 0 ORDER BY mstrecorddifferentiationtype.seqno ASC , selected DESC , name ASC"
	//logger.Log.Println("Main Query: ", sql)
	//logger.Log.Println("Param: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid)
	values := []entities.RawDiffEntity{}
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid)
	if err != nil {
		logger.Log.Println("GetAllTypeWiseCategories Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RawDiffEntity{}
		rows.Scan(&value.Typeid, &value.Typename, &value.Typeseq, &value.ID, &value.Name, &value.Seqno, &value.Selected, &value.ParentID)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}
func (dbc DbConn) GetAllRecordFields(req *entities.RecordDetailsRequestEntity) ([]entities.RecordFieldEntity, error) {
	logger.Log.Println("In side GetAllRecordField")
	//var sql = "SELECT mstrecordfield.id,mstrecordterms.id termid,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename,trnreordtracking.recordtrackvalue FROM mstrecordfield JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND trnreordtracking.referencetype='Additional'  JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0"
	//var sql = "SELECT mstrecordfield.id,mstrecordterms.id termid,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename,trnreordtracking.recordtrackvalue FROM mstrecordfield JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD'  JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 "

	var sql = "SELECT distinct(mstrecordfield.id),mstrecordterms.id termid,mstrecordterms.seq TermSeqNo,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename, " +
		" trnreordtracking.recordtrackvalue,mstrecorddifferentiation.seqno FROM mstrecordfield " +
		" JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND " +
		" mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND " +
		" mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD' JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND " +
		" mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN" +
		" msttermtype ON mstrecordterms.termtypeid=msttermtype.id join mstrecordfielddiff on mstrecordfield.id=mstrecordfielddiff.mstrecordfieldid" +
		" join mstrecorddifferentiation on mstrecordfielddiff.recorddiffid=mstrecorddifferentiation.id" +
		" WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND " +
		" trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND  mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 and " +
		" mstrecordfielddiff.recorddifftypeid !=? AND mstrecordfielddiff.deleteflg =0 AND mstrecorddifferentiation.deleteflg=0"

	//logger.Log.Println("Main Query: ", sql)
	//logger.Log.Println("Param: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordDiffTypeid)
	values := []entities.RecordFieldEntity{}
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordDiffTypeid)
	if err != nil {
		logger.Log.Println("GetAllRecordField Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RecordFieldEntity{}
		rows.Scan(&value.FieldID, &value.TermsID, &value.TermSeqNo, &value.TermsName, &value.TermsValue, &value.TermsTypeID, &value.TermsTypeName, &value.Value, &value.CatSeq)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

/*func (dbc DbConn) GetAllRecordFields(req *entities.RecordDetailsRequestEntity) ([]entities.RecordFieldEntity, error) {
	logger.Log.Println("In side GetAllRecordField")
	//var sql = "SELECT mstrecordfield.id,mstrecordterms.id termid,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename,trnreordtracking.recordtrackvalue FROM mstrecordfield JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND trnreordtracking.referencetype='Additional'  JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0"
	//var sql = "SELECT mstrecordfield.id,mstrecordterms.id termid,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename,trnreordtracking.recordtrackvalue FROM mstrecordfield JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD'  JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 "

	var sql = "SELECT distinct(mstrecordfield.id),mstrecordterms.id termid,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename, " +
		" trnreordtracking.recordtrackvalue,mstrecorddifferentiation.seqno FROM mstrecordfield " +
		" JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND " +
		" mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND " +
		" mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD' JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND " +
		" mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN" +
		" msttermtype ON mstrecordterms.termtypeid=msttermtype.id join mstrecordfielddiff on mstrecordfield.id=mstrecordfielddiff.mstrecordfieldid" +
		" join mstrecorddifferentiation on mstrecordfielddiff.recorddiffid=mstrecorddifferentiation.id" +
		" WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND " +
		" trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND  mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 and " +
		" mstrecordfielddiff.recorddifftypeid !=? AND mstrecordfielddiff.deleteflg =0 AND mstrecorddifferentiation.deleteflg=0"

	logger.Log.Println("Main Query: ", sql)
	logger.Log.Println("Param: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordDiffTypeid)
	values := []entities.RecordFieldEntity{}
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordDiffTypeid)
	if err != nil {
		logger.Log.Println("GetAllRecordField Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RecordFieldEntity{}
		rows.Scan(&value.FieldID, &value.TermsID, &value.TermsName, &value.TermsValue, &value.TermsTypeID, &value.TermsTypeName, &value.Value, &value.CatSeq)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}*/

/*func (dbc DbConn) GetAllRecordFields(req *entities.RecordDetailsRequestEntity) ([]entities.RecordFieldEntity, error) {
	logger.Log.Println("In side GetAllRecordField")
	//var sql = "SELECT mstrecordfield.id,mstrecordterms.id termid,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename,trnreordtracking.recordtrackvalue FROM mstrecordfield JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND trnreordtracking.referencetype='Additional'  JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0"
	var sql = "SELECT mstrecordfield.id,mstrecordterms.id termid,mstrecordterms.termname,mstrecordterms.termvalue,msttermtype.id termtypeid,msttermtype.termtypename,trnreordtracking.recordtrackvalue FROM mstrecordfield JOIN trnreordtracking ON mstrecordfield.clientid=trnreordtracking.clientid AND mstrecordfield.mstorgnhirarchyid=trnreordtracking.mstorgnhirarchyid AND mstrecordfield.id=trnreordtracking.referenceid AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD'  JOIN mstrecordterms ON mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid AND mstrecordterms.id=mstrecordfield.recordtermid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE trnreordtracking.activeflg=1 AND trnreordtracking.deleteflg=0 AND  trnreordtracking.clientid=? AND trnreordtracking.mstorgnhirarchyid=? AND  trnreordtracking.recordid=? AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 "
	logger.Log.Println("Main Query: ", sql)
	logger.Log.Println("Param: ", req.Clientid, req.Mstorgnhirarchyid, req.Recordid)
	values := []entities.RecordFieldEntity{}
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.Recordid)
	if err != nil {
		logger.Log.Println("GetAllRecordField Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RecordFieldEntity{}
		rows.Scan(&value.FieldID, &value.TermsID, &value.TermsName, &value.TermsValue, &value.TermsTypeID, &value.TermsTypeName, &value.Value)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}*/

/*func (dbc DbConn) GetRecordDirection(req *entities.RecordDetailsRequestEntity) (int64, error) {
	logger.Log.Println("In side GetRecordDirection dao model")
	//var sql = "SELECT mstbusinessdirection.direction FROM mstbusinessdirection WHERE mstbusinessdirection.activeflg=1 AND mstbusinessdirection.deleteflg=0 AND mstbusinessdirection.clientid=? AND mstbusinessdirection.mstorgnhirarchyid=? AND mstbusinessdirection.mstrecorddifferentiationtypeid=? AND mstbusinessdirection.mstrecorddifferentiationid=? LIMIT 1"
	//var values int64
	//rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid)

	var sql = "SELECT mstbusinessdirection.direction FROM mstbusinessdirection WHERE mstbusinessdirection.activeflg=1 AND mstbusinessdirection.deleteflg=0 AND mstbusinessdirection.clientid=? AND mstbusinessdirection.mstorgnhirarchyid=? AND mstbusinessdirection.mstrecorddifferentiationtypeid=? AND mstbusinessdirection.mstrecorddifferentiationid=?  and mstbusinessdirection.baseconfig=?"
	var values int64
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid, req.BaseConfig)
	if err != nil {
		logger.Log.Println("Getdirection Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&values)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}*/

func (dbc DbConn) GetRecordDirection(req *entities.RecordDetailsRequestEntity) (int64, error) {

	logger.Log.Println("In side GetRecordDirection dao model")

	//var sql = "SELECT mstbusinessdirection.direction FROM mstbusinessdirection WHERE mstbusinessdirection.activeflg=1 AND mstbusinessdirection.deleteflg=0 AND mstbusinessdirection.clientid=? AND mstbusinessdirection.mstorgnhirarchyid=? AND mstbusinessdirection.mstrecorddifferentiationtypeid=? AND mstbusinessdirection.mstrecorddifferentiationid=? LIMIT 1"

	//var values int64

	//rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid)

	req.BaseConfig = 1

	var sql = "SELECT mstbusinessdirection.direction FROM mstbusinessdirection WHERE mstbusinessdirection.activeflg=1 AND mstbusinessdirection.deleteflg=0 AND mstbusinessdirection.clientid=? AND mstbusinessdirection.mstorgnhirarchyid=? AND mstbusinessdirection.mstrecorddifferentiationtypeid=? AND mstbusinessdirection.mstrecorddifferentiationid=?  and mstbusinessdirection.baseconfig=?"

	var values int64

	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid, req.BaseConfig)

	if err != nil {

		logger.Log.Println("Getdirection Get Statement Prepare Error", err)

		return values, err

	}

	defer rows.Close()

	for rows.Next() {

		rows.Scan(&values)

	}

	//            logger.Log.Println("values is ----->", values)

	return values, nil

}

func (dbc DbConn) GetOrgnIDbyrecordID(RecordID int64) (int64, error) {
	logger.Log.Println("In side GetOrgnIDbyrecordID dao model")
	var sql = "SELECT mstorgnhirarchyid FROM trnrecord WHERE id=?"
	var values int64
	rows, err := dbc.DB.Query(sql, RecordID)
	if err != nil {
		logger.Log.Println("Getdirection Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&values)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) GetOrgnIDbyrecordCode(Recordcode string) (int64, error) {
	logger.Log.Println("In side GetOrgnIDbyrecordID dao model")
	var sql = "SELECT mstorgnhirarchyid FROM trnrecord WHERE code=?"
	var values int64
	rows, err := dbc.DB.Query(sql, Recordcode)
	if err != nil {
		logger.Log.Println("Getdirection Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&values)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) CheckOrgnIDbyuserID(Clientid int64, orgnID int64, Userid int64) ([]int64, error) {
	logger.Log.Println("In side CheckOrgnIDbyuserID dao model")
	var sql = "SELECT id FROM mstgroupmember WHERE clientid=? AND mstorgnhirarchyid=? AND userid=? AND deleteflg=0 AND activeflg=1"
	var values []int64
	rows, err := dbc.DB.Query(sql, Clientid, orgnID, Userid)
	if err != nil {
		logger.Log.Println("Getdirection Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		var value int64
		rows.Scan(&value)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) GetRecordWiseImpact(req *entities.RecordDetailsRequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetRecordimpact")
	var sql = "SELECT  mstrecorddifferentiation.id, mstrecorddifferentiation.name, mstrecorddifferentiation.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstrecordtype, mstrecorddifferentiationtype, mstrecorddifferentiation WHERE mstrecordtype.clientid = ? AND mstrecordtype.mstorgnhirarchyid = ? AND mstrecordtype.activeflg = 1 AND mstrecordtype.deleteflg = 0 AND mstrecordtype.fromrecorddifftypeid = ? AND mstrecordtype.fromrecorddiffid = ? AND mstrecorddifferentiationtype.activeflg = 1 AND mstrecorddifferentiationtype.deleteflg = 0 AND mstrecorddifferentiationtype.id = mstrecordtype.torecorddifftypeid AND mstrecorddifferentiationtype.seqno = ? AND mstrecorddifferentiation.activeflg = 1 AND mstrecorddifferentiation.deleteflg = 0 AND mstrecorddifferentiation.id = mstrecordtype.torecorddiffid ORDER BY mstrecorddifferentiationtype.seqno ASC , mstrecorddifferentiation.seqno ASC"
	values := []entities.RecordcatchildEntity{}
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid, 6)
	if err != nil {
		logger.Log.Println("GetRecordimpact Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RecordcatchildEntity{}
		rows.Scan(&value.ID, &value.Name, &value.Seqno, &value.Typeid, &value.Typename, &value.Typeseq)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) GetRecordWiseUrgency(req *entities.RecordDetailsRequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetRecordurgency")
	var sql = "SELECT  mstrecorddifferentiation.id, mstrecorddifferentiation.name, mstrecorddifferentiation.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstrecordtype, mstrecorddifferentiationtype, mstrecorddifferentiation WHERE mstrecordtype.clientid = ? AND mstrecordtype.mstorgnhirarchyid = ? AND mstrecordtype.activeflg = 1 AND mstrecordtype.deleteflg = 0 AND mstrecordtype.fromrecorddifftypeid = ? AND mstrecordtype.fromrecorddiffid = ? AND mstrecorddifferentiationtype.activeflg = 1 AND mstrecorddifferentiationtype.deleteflg = 0 AND mstrecorddifferentiationtype.id = mstrecordtype.torecorddifftypeid AND mstrecorddifferentiationtype.seqno = ? AND mstrecorddifferentiation.activeflg = 1 AND mstrecorddifferentiation.deleteflg = 0 AND mstrecorddifferentiation.id = mstrecordtype.torecorddiffid ORDER BY mstrecorddifferentiationtype.seqno ASC , mstrecorddifferentiation.seqno ASC"
	values := []entities.RecordcatchildEntity{}
	rows, err := dbc.DB.Query(sql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid, 7)
	if err != nil {
		logger.Log.Println("GetRecordurgency Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RecordcatchildEntity{}
		rows.Scan(&value.ID, &value.Name, &value.Seqno, &value.Typeid, &value.Typename, &value.Typeseq)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) CheckParentAsChildRecord(req *entities.RecordDetailsRequestEntity) (int64, error) {
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

func (dbc DbConn) GetRecordCatByLastID(req *entities.ChildRecordSearchEntity) ([]entities.RawDiffEntity, error) {
	logger.Log.Println("In side GetRecordCatByLastID")
	var sql = "SELECT mstrecorddifferentiationtype.id typeid, mstrecorddifferentiationtype.typename typename, mstrecorddifferentiationtype.seqno typeseq, mstrecorddifferentiation.id,mstrecorddifferentiation.name,mstrecorddifferentiation.seqno,IF(mstrecorddifferentiation.id=mrd.id,1,0) selected,mstrecorddifferentiation.parentid FROM mstrecordtype JOIN mstrecorddifferentiation ON mstrecorddifferentiation.id = mstrecordtype.torecorddiffid JOIN mstrecorddifferentiation mrd ON mstrecorddifferentiation.parentid=mrd.parentid AND mrd.id=? AND mrd.activeflg=1 AND mrd.deleteflg=0 JOIN mstrecorddifferentiationtype ON mstrecorddifferentiationtype.id = mstrecordtype.torecorddifftypeid  AND mstrecorddifferentiation.recorddifftypeid = mstrecorddifferentiationtype.id WHERE mstrecordtype.clientid = ? AND mstrecordtype.mstorgnhirarchyid = ? AND mstrecordtype.activeflg = 1 AND mstrecordtype.deleteflg = 0 AND mstrecordtype.fromrecorddifftypeid = ? AND mstrecordtype.fromrecorddiffid = ? AND mstrecorddifferentiationtype.activeflg = 1 AND mstrecorddifferentiationtype.deleteflg = 0 AND mstrecorddifferentiationtype.parentid = 1 AND mstrecorddifferentiation.activeflg = 1 AND mstrecorddifferentiation.deleteflg = 0 ORDER BY mstrecorddifferentiationtype.seqno ASC ,selected DESC,  mstrecorddifferentiation.seqno ASC"
	//logger.Log.Println("Main Query: ", sql)
	//logger.Log.Println("Param: ", req.CategoryID, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid)
	values := []entities.RawDiffEntity{}
	rows, err := dbc.DB.Query(sql, req.CategoryID, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid)
	if err != nil {
		logger.Log.Println("GetRecordCatByLastID Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RawDiffEntity{}
		rows.Scan(&value.Typeid, &value.Typename, &value.Typeseq, &value.ID, &value.Name, &value.Seqno, &value.Selected, &value.ParentID)
		values = append(values, value)
	}
	return values, nil
}
