package dao

import (
	"database/sql"
	"errors"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

//var getstatusid = "SELECT recorddiffid FROM maprecordstatetodifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND mststateid =? AND deleteflg=0 AND activeflg=1"
var getstatusid = "SELECT a.recorddiffid ,b.seqno,b.name FROM maprecordstatetodifferentiation a,mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.mststateid =? AND a.deleteflg=0 AND a.activeflg=1 AND a.recorddiffid=b.id AND a.recorddifftypeid=3"
var updatestatus = "INSERT INTO maprecordtorecorddifferentiation (clientid, mstorgnhirarchyid,recordid, recordstageid, recorddifftypeid, recorddiffid,seqno,createddate) VALUES (?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())))"
var updateactivestatus = "INSERT INTO maprecordtorecorddifferentiation (clientid, mstorgnhirarchyid,recordid, recordstageid, recorddifftypeid, recorddiffid,seqno,createddate,islatest) VALUES (?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),0)"
var lateststageID = "SELECT max(id) stageid FROM trnrecordstage WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
var updatepreviousstatus = "Update maprecordtorecorddifferentiation set islatest=0 WHERE recorddifftypeid=3 AND recordid=? AND clientid=? AND mstorgnhirarchyid=?"
var respnsemeterstatusseq = "select b.seqno as seqno FROM mstslafcrecorddiff a,mstslaindicatorterm b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddiffidstatus=? AND a.slametertypeid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.startstopindicator=b.id AND b.activeflg=1 AND b.deleteflg=0"
var recorddetails = "select a.clientid as clientid, a.mstorgnhirarchyid as mstorgnhirarchyid,a.recordid as recordid,a.recordtypeid as recordtypeid, b.workingcatid as workingcatid, c.priorityid as priorityid from (select clientid,mstorgnhirarchyid,recordid,recorddiffid as recordtypeid from maprecordtorecorddifferentiation where recordid = ? and recorddifftypeid = 2 and islatest=1) a, (select recordid, recorddiffid as workingcatid from maprecordtorecorddifferentiation where recordid = ? and isworking = 1 and islatest=1) b, (select recordid, recorddiffid as priorityid from maprecordtorecorddifferentiation where recordid = ? and recorddifftypeid = 5 and islatest=1) c where a.recordid=b.recordid and a.recordid= b.recordid"
var latesttrnhistory = "Select mstslaentityid,donotupdatesladue,recordtimetoint,mstslastateid,commentonrecord,fromclientuserid,slastartstopindicator,recorddatetime,recorddatetoint FROM trnslaentityhistory WHERE therecordid=? AND clientid=? AND mstorgnhirarchyid=? AND activeflg=1 AND deleteflg=0 order by id desc limit 1"
var recordtypediffid = "SELECT recorddiffid as ID FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=2"
var currentstateID = "select currentstateid FROM mstrequest where id in (select max(mstrequestid) From maprequestorecord where recordid=? AND clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1)"
var updatechild = "UPDATE mstparentchildmap SET isattached='N' WHERE childrecordid=? AND clientid=? AND mstorgnhirarchyid=?"
var latestsatusid = "SELECT a.recorddiffid,b.seqno,b.name FROM maprecordtorecorddifferentiation a,mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=3 AND a.recordid=? AND a.islatest=1 AND a.recorddiffid=b.id"
var activestatussql = "SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1 AND seqno=2"

func (mdao DbConn) Getrecorddiffidbystateid(ClientID int64, Mstorgnhirarchyid int64, ReordstatusID int64) (int64, int64, string, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var recorddiffid int64
	var seqno int64
	var name string
	logger.Log.Println("11111111111111111111111111111111111111111111111111111111111111111111------>", getstatusid)
	logger.Log.Println("222222222222222222222222222222222222222222222222222------>", ClientID, Mstorgnhirarchyid, ReordstatusID)
	rows, err := mdao.DB.Query(getstatusid, ClientID, Mstorgnhirarchyid, ReordstatusID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecorddiffidbystateid Get Statement Prepare Error", err)
		return recorddiffid, seqno, name, err
	}
	for rows.Next() {
		rows.Scan(&recorddiffid, &seqno, &name)

	}
	return recorddiffid, seqno, name, nil
}

func (mdao DbConn) Getemeterseqno(ClientID int64, Mstorgnhirarchyid int64, recorddiffid int64, metertypeid int64) (int64, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var seqno int64
	rows, err := mdao.DB.Query(respnsemeterstatusseq, ClientID, Mstorgnhirarchyid, recorddiffid, metertypeid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermnamesbystate Get Statement Prepare Error", err)
		return seqno, err
	}
	for rows.Next() {
		rows.Scan(&seqno)

	}
	return seqno, nil
}

func (mdao DbConn) Getrecorddetails(recordid int64) (entities.SLATabEntity, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	tz := entities.SLATabEntity{}
	stmt, error := mdao.DB.Prepare(recorddetails)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return tz, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(recordid, recordid, recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecorddetails Get Statement Prepare Error", err)
		return tz, err
	}
	for rows.Next() {
		rows.Scan(&tz.ClientID, &tz.Mstorgnhirarchyid, &tz.RecordID, &tz.RecordtypeID, &tz.WorkingcatID, &tz.PriorityID)

	}
	return tz, nil
}

func (mdao DbConn) GetLatesttrnhistory(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (entities.TrnslaentityhistoryEntity, error) {
	logger.Log.Println("In side GetLatesttrnhistory")
	tz := entities.TrnslaentityhistoryEntity{}
	rows, err := mdao.DB.Query(latesttrnhistory, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLatesttrnhistory Get Statement Prepare Error", err)
		return tz, err
	}
	for rows.Next() {
		rows.Scan(&tz.Mstslaentityid, &tz.Donotupdatesladue, &tz.Recordtimetoint, &tz.Mstslastateid, &tz.Commentonrecord, &tz.Fromclientuserid, &tz.Slastartstopindicator, &tz.Recorddatetime, &tz.Recorddatetoint)

	}
	return tz, nil
}

func (mdao DbConn) GetMaxstageID(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var stageid int64
	rows, err := mdao.DB.Query(lateststageID, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermnamesbystate Get Statement Prepare Error", err)
		return stageid, err
	}
	for rows.Next() {
		rows.Scan(&stageid)

	}
	return stageid, nil
}

func Updaterecordstatus(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, recorddiffid int64, laststageID int64) (int64, error) {
	logger.Log.Println("trnreordtracking query -->", updatestatus)
	logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, laststageID)

	stmt, err := tx.Prepare(updatestatus)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, laststageID, 3, recorddiffid, 0)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Last ID fetch Error")
	}

	return lastInsertedID, nil
}

func UpdaterecordActivestatus(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, recorddiffid int64, laststageID int64) (int64, error) {
	logger.Log.Println("trnreordtracking query -->", updatestatus)
	logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, laststageID)

	stmt, err := tx.Prepare(updateactivestatus)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, laststageID, 3, recorddiffid, 0)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Last ID fetch Error")
	}

	return lastInsertedID, nil
}

func Updatepreviousstatus(tx *sql.Tx, RecordID int64, ClientID int64, Mstorgnhirarchyid int64) error {
	logger.Log.Println("Updatelaststatus query -->", updateworkinglabel)
	stmt, err := tx.Prepare(updatepreviousstatus)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(RecordID, ClientID, Mstorgnhirarchyid)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func Updatechildrecord(tx *sql.Tx, RecordID int64, ClientID int64, Mstorgnhirarchyid int64) error {
	logger.Log.Println("Updatelaststatus query -->", updateworkinglabel)
	stmt, err := tx.Prepare(updatechild)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(RecordID, ClientID, Mstorgnhirarchyid)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) Getrecordtypediffid(Recordid int64, ClientID int64, Mstorgnhirarchyid int64) (int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var ID int64
	rows, err := mdao.DB.Query(recordtypediffid, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return ID, nil
}

func (mdao DbConn) GetrecordlateststateID(Recordid int64, ClientID int64, Mstorgnhirarchyid int64) (int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var ID int64
	rows, err := mdao.DB.Query(currentstateID, Recordid, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return ID, nil
}

func (mdao DbConn) Getcurrentsatusid(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, int64, string, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var recorddiffid int64
	var name string
	var seqno int64
	rows, err := mdao.DB.Query(latestsatusid, ClientID, Mstorgnhirarchyid, RecordID)
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

func (mdao DbConn) GetActivestatusID(ClientID int64, Mstorgnhirarchyid int64) (int64, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var activestatusID int64
	rows, err := mdao.DB.Query(activestatussql, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getcurrentsatusid Get Statement Prepare Error", err)
		return activestatusID, err
	}
	for rows.Next() {
		rows.Scan(&activestatusID)

	}
	return activestatusID, nil
}

func InsertRecordClosure(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, RecordSeq int64, Closuredate string) error {
	sql := "INSERT INTO mstrecordautoclosure(clientid,mstorgnhirarchyid,recordid,recordseq,closuredate,closuredt) VALUES(?,?,?,?,round(UNIX_TIMESTAMP(?)),?)"
	stmt, err := tx.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, RecordSeq, Closuredate, Closuredate)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) Getprestausname(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var name string
	var sql = "SELECT b.name FROM maprecordtorecorddifferentiation a,mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.recorddifftypeid=3 AND a.islatest=0 AND a.id not in (SELECT max(id) FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=3 AND islatest=0 order by id desc) AND a.recorddiffid=b.id order by a.id desc limit 1"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getcurrentsatusid Get Statement Prepare Error", err)
		return name, err
	}
	for rows.Next() {
		rows.Scan(&name)

	}
	return name, nil
}

func (mdao DbConn) Getrecordtypeseq(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side Getrecordtypeseq")
	var typeSeq int64
	var sql = "SELECT b.seqno FROM maprecordtorecorddifferentiation a,mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.recorddifftypeid=2 AND a.islatest=1 AND a.recorddiffid=b.id"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecordtypeseq Get Statement Prepare Error", err)
		return typeSeq, err
	}
	for rows.Next() {
		rows.Scan(&typeSeq)

	}
	return typeSeq, nil
}

func (mdao DbConn) Getrecorddiffidforstask(ClientID int64, Mstorgnhirarchyid int64, ReordstatusID int64) (int64, int64, string, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var sql = "SELECT a.torecorddiffid,b.seqno,b.name FROM mstrecordtype a,mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.fromrecorddifftypeid=3 AND a.fromrecorddiffid=? AND a.torecorddifftypeid=3 AND a.activeflg=1 AND a.deleteflg=0 AND a.torecorddiffid = b.id AND b.activeflg=1 AND b.deleteflg=0"
	var recorddiffid int64
	var seqno int64
	var name string
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, ReordstatusID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecorddiffidbystateid Get Statement Prepare Error", err)
		return recorddiffid, seqno, name, err
	}
	for rows.Next() {
		rows.Scan(&recorddiffid, &seqno, &name)

	}
	return recorddiffid, seqno, name, nil
}

func (mdao DbConn) GetrecorddiffidforstaskNew(ClientID int64, Mstorgnhirarchyid int64, ReordstatusID int64, ParentrcordtypeID int64, ChildrecordtypeID int64) (int64, int64, string, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	//var sql = "select d.id,d.seqno,d.name from mstrecordtype a ,mstrecordtype b,mstrecordtype c,mstrecorddifferentiation d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.fromrecorddifftypeid = 2 AND a.fromrecorddiffid = ? AND a.torecorddifftypeid=3 AND a.torecorddiffid=? AND a.torecorddiffid = b.fromrecorddiffid AND c.fromrecorddifftypeid = 2 AND c.fromrecorddiffid = ? AND c.torecorddiffid = b.torecorddiffid AND c.torecorddiffid = d.id AND a.activeflg=1 AND a.deleteflg=0 AND b.activeflg=1 AND b.deleteflg=0 AND c.activeflg=1 AND c.deleteflg=0 AND d.activeflg=1 AND d.deleteflg=0"
	var sql = "select d.id,d.seqno,d.name from mstrecordtype a ,mstrecordtype b,mstrecordtype c,mstrecorddifferentiation d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.fromrecorddifftypeid = 2 AND a.fromrecorddiffid = ? AND a.torecorddifftypeid=3 AND a.torecorddiffid=? AND a.torecorddiffid = b.fromrecorddiffid AND c.fromrecorddifftypeid = 2 AND c.fromrecorddiffid = ? AND c.torecorddiffid = b.torecorddiffid AND c.torecorddiffid = d.id AND d.seqno NOT IN (0) AND a.activeflg=1 AND a.deleteflg=0 AND b.activeflg=1 AND b.deleteflg=0 AND c.activeflg=1 AND c.deleteflg=0 AND d.activeflg=1 AND d.deleteflg=0"
	var recorddiffid int64
	var seqno int64
	var name string
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, ParentrcordtypeID, ReordstatusID, ChildrecordtypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecorddiffidbystateid Get Statement Prepare Error", err)
		return recorddiffid, seqno, name, err
	}
	for rows.Next() {
		rows.Scan(&recorddiffid, &seqno, &name)

	}
	return recorddiffid, seqno, name, nil
}

func (mdao DbConn) GetrecorddiffidforstaskNewFromStask(ClientID int64, Mstorgnhirarchyid int64, ReordstatusID int64, ParentrcordtypeID int64, ChildrecordtypeID int64) (int64, int64, string, error) {
	var sql = "select d.id,d.seqno,d.name from mstrecordtype a ,mstrecordtype b,mstrecordtype c,mstrecorddifferentiation d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.fromrecorddifftypeid = 2 AND a.fromrecorddiffid = ? AND a.torecorddifftypeid=3 AND a.torecorddiffid=? AND a.torecorddiffid = b.fromrecorddiffid AND c.fromrecorddifftypeid = 2 AND c.fromrecorddiffid = ? AND c.torecorddiffid = b.torecorddiffid AND c.torecorddiffid = d.id AND a.activeflg=1 AND a.deleteflg=0 AND b.activeflg=1 AND b.deleteflg=0 AND c.activeflg=1 AND c.deleteflg=0 AND d.activeflg=1 AND d.deleteflg=0"
	logger.Log.Println("In side GetrecorddiffidforstaskNewFromStask >>>>>>>>>>>>>>>>>>.", sql)
	logger.Log.Println("In side GetrecorddiffidforstaskNewFromStask >>>>>>>>>>>>>>>>>>.", ClientID, Mstorgnhirarchyid, ChildrecordtypeID, ReordstatusID, ParentrcordtypeID)
	var recorddiffid int64
	var seqno int64
	var name string
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, ChildrecordtypeID, ReordstatusID, ParentrcordtypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecorddiffidbystateid Get Statement Prepare Error", err)
		return recorddiffid, seqno, name, err
	}
	for rows.Next() {
		err := rows.Scan(&recorddiffid, &seqno, &name)
		logger.Log.Println("Rows error is  ----+++++++++++++++++++++++++++++++++++++++++-->", err)
	}
	return recorddiffid, seqno, name, nil
}

func (mdao DbConn) Getparentrecordids(RecordID int64, ClientID int64, Mstorgnhirarchyid int64, Recorddifftypeid int64, Recorddiffid int64) ([]int64, error) {
	logger.Log.Println("Parameter is ------>", RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)

	var parentids []int64
	var getparentids = "select parentrecordid from mstparentchildmap where childrecordid=? and clientid=? and mstorgnhirarchyid=? and deleteflg=0 and activeflg=1 and isattached='Y' and parentrecordid !=0"
	rows, err := mdao.DB.Query(getparentids, RecordID, ClientID, Mstorgnhirarchyid)
	logger.Log.Println("Rows error is  ----+++++++++++++++++++++++++++++++++++++++++-->", err)

	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return parentids, err
	}
	for rows.Next() {
		var ID int64
		err = rows.Scan(&ID)
		parentids = append(parentids, ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return parentids, nil
}

func (mdao DbConn) GetSLAdataexist(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var ID int64
	var sql = "SELECT id FROM mstsladue WHERE clientid=? AND mstorgnhirarchyid=? AND therecordid=? "
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return ID, nil
}

func (mdao DbConn) GetStatusPriority(ClientID int64, Mstorgnhirarchyid int64, Recordid int64, TypeID int64, ChildRecordID int64) (int64, int64, error) {
	//logger.Log.Println("In side GetStatusPriority---------------------->")
	logger.Log.Println("In side GetStatusPriority---------------------->", ClientID, Mstorgnhirarchyid, Recordid, 2, TypeID, ChildRecordID)
	var statusID int64
	var priorityID int64
	var sql = "SELECT c.differentiationid,c.priority FROM mstparentchildmap a,maprecordtorecorddifferentiation b,mstrecorddifferentiationpriority c WHERE  a.clientid=? AND a.mstorgnhirarchyid=? AND a.parentrecordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.childrecordid = b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid = c.differentiationid AND c.typedifferentiationtypeid=? AND c.typedifferentiationid=? AND c.activeflg=1 AND c.deleteflg=0 AND a.childrecordid not in(?) AND a.isattached='Y' ORDER BY c.priority ASC LIMIT 1"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Recordid, 2, TypeID, ChildRecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetStatusPriority Get Statement Prepare Error", err)
		return statusID, priorityID, err
	}
	for rows.Next() {
		err = rows.Scan(&statusID, &priorityID)
		logger.Log.Println("GetStatusPriority rows.next() Error", err)
	}
	return statusID, priorityID, nil
}

func (mdao DbConn) GetStatusPriorityFromSR(ClientID int64, Mstorgnhirarchyid int64, Recordid int64, TypeID int64) (int64, int64, error) {
	logger.Log.Println("In side GetStatusPriority")
	var statusID int64
	var priorityID int64
	var sql = "SELECT c.differentiationid,c.priority FROM mstparentchildmap a,maprecordtorecorddifferentiation b,mstrecorddifferentiationpriority c WHERE  a.clientid=? AND a.mstorgnhirarchyid=? AND a.parentrecordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.childrecordid = b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid = c.differentiationid AND c.typedifferentiationtypeid=? AND c.typedifferentiationid=? AND c.activeflg=1 AND c.deleteflg=0 ORDER BY c.priority ASC LIMIT 1"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Recordid, 2, TypeID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetStatusPriority Get Statement Prepare Error", err)
		return statusID, priorityID, err
	}
	for rows.Next() {
		err = rows.Scan(&statusID, &priorityID)
		logger.Log.Println("GetStatusPriority rows.next() Error", err)
	}
	return statusID, priorityID, nil
}

func (mdao DbConn) GetUpdateStatusPriority(ClientID int64, Mstorgnhirarchyid int64, TypeID int64, StatusID int64) (int64, error) {
	logger.Log.Println("In side GetUpdateStatusPriority")
	var priorityID int64
	var sql = "SELECT priority FROM mstrecorddifferentiationpriority WHERE clientid=? AND mstorgnhirarchyid=? AND typedifferentiationtypeid=? AND typedifferentiationid=? AND differentiationtypeid=? AND differentiationid=? AND activeflg=1 AND deleteflg=0"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, 2, TypeID, 3, StatusID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetStatusPriority Get Statement Prepare Error", err)
		return priorityID, err
	}
	for rows.Next() {
		err = rows.Scan(&priorityID)
		logger.Log.Println("GetStatusPriority rows.next() Error", err)
	}
	return priorityID, nil
}

func UpdateTaskflag(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, RecorddiffID int64) error {
	var sql = "UPDATE maprecordtorecorddifferentiation SET istask=1  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=3 AND recorddiffid=? AND islatest=1"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, OrgnID, RecordID, RecorddiffID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetcurrentTaskflag(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var sql = "SELECT istask FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=3 AND islatest=1"
	var taskID int64
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getcurrentsatusid Get Statement Prepare Error", err)
		return taskID, err
	}
	for rows.Next() {
		rows.Scan(&taskID)

	}
	return taskID, nil
}

func UpdateApproveflag(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE trnrecord SET isapprove=1  WHERE clientid=? AND mstorgnhirarchyid=? AND id=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateApproveflagForOne(ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE trnrecord SET isapprove=1  WHERE clientid=? AND mstorgnhirarchyid=? AND id=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateApproveflagzero(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE trnrecord SET isapprove=0  WHERE clientid=? AND mstorgnhirarchyid=? AND id=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetIsapprovalworkflow(RecordID int64) (int64, error) {
	logger.Log.Println("In side GetIsapprovalworkflow")
	var sql = "SELECT isapproveworkflow FROM trnrecord WHERE id=?"
	var taskID int64
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetIsapprovalworkflow Get Statement Prepare Error", err)
		return taskID, err
	}
	for rows.Next() {
		rows.Scan(&taskID)

	}
	return taskID, nil
}

func (mdao DbConn) GetIsapprovalFlag(RecordID int64) (int64, error) {
	logger.Log.Println("In side GetIsapprovalworkflow")
	var sql = "SELECT isapprove FROM trnrecord WHERE id=?"
	var taskID int64
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetIsapprovalworkflow Get Statement Prepare Error", err)
		return taskID, err
	}
	for rows.Next() {
		rows.Scan(&taskID)

	}
	return taskID, nil
}

func (mdao DbConn) Getrecorddifferation(ClientID int64, Mstorgnhirarchyid int64, Statusseq int64) (int64, int64, string, error) {
	logger.Log.Println("In side Getrecorddifferation")
	var sql = "SELECT id,seqno,name FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=3 AND seqno=? AND activeflg=1 AND deleteflg=0"
	var recorddiffid int64
	var seqno int64
	var name string
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Statusseq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecorddifferation Get Statement Prepare Error", err)
		return recorddiffid, seqno, name, err
	}
	for rows.Next() {
		rows.Scan(&recorddiffid, &seqno, &name)

	}
	return recorddiffid, seqno, name, nil
}

func (mdao DbConn) GetGrpname(GrpID int64) (string, error) {
	logger.Log.Println("In side Getrecorddiffidbystateid")
	var name string
	var sql = "SELECT name FROM mstsupportgrp WHERE id=?"
	rows, err := mdao.DB.Query(sql, GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getcurrentsatusid Get Statement Prepare Error", err)
		return name, err
	}
	for rows.Next() {
		rows.Scan(&name)

	}
	return name, nil
}

func (mdao DbConn) GetPreviousstatusdate(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, error) {
	logger.Log.Println("In side GetPreviousstatusdate")
	var name string
	var sql = "SELECT FROM_UNIXTIME(createddate) FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=3 AND recordid=? AND islatest=1"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetPreviousstatusdate Get Statement Prepare Error", err)
		return name, err
	}
	for rows.Next() {
		rows.Scan(&name)

	}
	return name, nil
}

func (mdao DbConn) GetReopendate(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, recorddiffseq int64) (string, error) {
	logger.Log.Println("In side GetPreviousstatusdate")
	var name string
	var sql = "SELECT FROM_UNIXTIME(a.createddate) FROM maprecordtorecorddifferentiation a,mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=3 AND a.recordid=? AND b.clientid=? AND b.mstorgnhirarchyid=? AND b.recorddifftypeid=3 AND b.seqno=? AND b.activeflg=1 AND b.deleteflg=0 AND a.recorddiffid=b.id order by a.id desc limit 1"
	logger.Log.Println("Reopen query is --------------->", sql)
	logger.Log.Println("Reopen query is --------------->", ClientID, Mstorgnhirarchyid, RecordID, ClientID, Mstorgnhirarchyid, recorddiffseq)
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID, ClientID, Mstorgnhirarchyid, recorddiffseq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetPreviousstatusdate Get Statement Prepare Error", err)
		return name, err
	}
	for rows.Next() {
		rows.Scan(&name)

	}
	return name, nil
}

func (mdao DbConn) GetFirstResponseValue(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, error) {
	logger.Log.Println("In side GetPreviousstatusdate")
	var date string
	var sql = "SELECT COALESCE(firstresponsedatetime,'NA') FROM recordfulldetails WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetPreviousstatusdate Get Statement Prepare Error", err)
		return date, err
	}
	for rows.Next() {
		rows.Scan(&date)

	}
	return date, nil
}

func (mdao DbConn) GetFirstResolutionValue(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, error) {
	logger.Log.Println("In side GetPreviousstatusdate")
	var date string
	var sql = "SELECT COALESCE(firstresodatetime,'NA') FROM recordfulldetails WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetPreviousstatusdate Get Statement Prepare Error", err)
		return date, err
	}
	for rows.Next() {
		rows.Scan(&date)

	}
	return date, nil
}

func (mdao DbConn) FetchRecordtypeID(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, TypeID int64) (int64, error) {
	logger.Log.Println("In side GetUpdateStatusPriority")
	var typeID int64
	var sql = "SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recordid=? AND deleteflg=0 AND activeflg=1 AND islatest=1"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, TypeID, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchRecordtypeID Get Statement Prepare Error", err)
		return typeID, err
	}
	for rows.Next() {
		err = rows.Scan(&typeID)
		logger.Log.Println("FetchRecordtypeID rows.next() Error", err)
	}
	return typeID, nil
}

func (mdao DbConn) FetchDifferentiationIDBySeq(ClientID int64, Mstorgnhirarchyid int64, TypeID int64, Sequance int64) (int64, error) {
	logger.Log.Println("In side GetUpdateStatusPriority")
	var typeID int64
	var sql = "SELECT id FROM mstrecorddifferentiation where clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND seqno=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, TypeID, Sequance)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchRecordtypeID Get Statement Prepare Error", err)
		return typeID, err
	}
	for rows.Next() {
		err = rows.Scan(&typeID)
		logger.Log.Println("FetchRecordtypeID rows.next() Error", err)
	}
	return typeID, nil
}

func (mdao DbConn) FetchDifferentiationDetailsByID(ID int64) (int64, int64, string, error) {
	logger.Log.Println("In side Getrecorddifferation")
	var sql = "SELECT id,seqno,name FROM mstrecorddifferentiation WHERE id=? AND activeflg=1 AND deleteflg=0"
	var recorddiffid int64
	var seqno int64
	var name string
	rows, err := mdao.DB.Query(sql, ID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecorddifferation Get Statement Prepare Error", err)
		return recorddiffid, seqno, name, err
	}
	for rows.Next() {
		rows.Scan(&recorddiffid, &seqno, &name)

	}
	return recorddiffid, seqno, name, nil
}

func (mdao DbConn) FetchCurrentGrpID(ID int64) (int64, error) {
	logger.Log.Println("In side FetchCurrentGrpID")
	var sql = "SELECT a.mstgroupid FROM mstrequest a,maprequestorecord b WHERE b.recordid=? AND a.id=b.mstrequestid AND b.activeflg=1 AND b.deleteflg=0 AND a.activeflg=1 AND a.deleteflg=0 order by b.id desc limit 1"
	var grpID int64
	rows, err := mdao.DB.Query(sql, ID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchCurrentGrpID Get Statement Prepare Error", err)
		return grpID, err
	}
	for rows.Next() {
		rows.Scan(&grpID)

	}
	return grpID, nil
}

func (mdao DbConn) FetchResponseResolutionCompleteValue(ClientID int64, MstorgnhirarchyID int64, RecordID int64) (int64, int64, error) {
	logger.Log.Println("In side FetchResponseResolutionCompleteValue")
	var sql = "SELECT isresponsecomplete,isresolutioncomplete FROM mstsladue WHERE clientid=? AND mstorgnhirarchyid=? AND therecordid=?"
	var isresponsecomplete int64
	var isresolutioncomplete int64
	rows, err := mdao.DB.Query(sql, ClientID, MstorgnhirarchyID, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchResponseResolutionCompleteValue Get Statement Prepare Error", err)
		return isresponsecomplete, isresolutioncomplete, err
	}
	for rows.Next() {
		rows.Scan(&isresponsecomplete, &isresolutioncomplete)

	}
	return isresponsecomplete, isresolutioncomplete, nil
}

func (mdao DbConn) FetchSLADueRow(ClientID int64, MstorgnhirarchyID int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side FetchSLADueRow")
	var sql = "SELECT id FROM mstsladue WHERE clientid=? AND mstorgnhirarchyid=? AND therecordid=?"
	var id int64
	rows, err := mdao.DB.Query(sql, ClientID, MstorgnhirarchyID, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchSLADueRow Get Statement Prepare Error", err)
		return id, err
	}
	for rows.Next() {
		rows.Scan(&id)

	}
	return id, nil
}

func (mdao DbConn) FetchAutoCloseRecordCount(RecordID int64) (int64, error) {
	logger.Log.Println("In side FetchAutoCloseRecordCount")
	var sql = "SELECT count(id) count FROM mstrecordautoclosure WHERE recordid=? AND islatest=1"
	var count int64
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchSLADueRow Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		rows.Scan(&count)

	}
	return count, nil
}

func (mdao DbConn) UpdateIslatestFlag(RecordID int64) error {
	var sql = "UPDATE mstrecordautoclosure SET islatest=0 WHERE recordid=?"
	logger.Log.Println("UpdateIslatestFlag query -->", sql)
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) DeleteFromClosureTable(RecordID int64) error {
	var sql = "DELETE FROM mstrecordautoclosure WHERE recordid=?"
	logger.Log.Println("UpdateIslatestFlag query -->", sql)
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) FetchResolvedStateID(ClientID int64, OrgnID int64, ResolveStatusSeq int64) (int64, error) {
	logger.Log.Println("In side FetchResolvedStateID")
	var sql = "SELECT a.mststateid FROM maprecordstatetodifferentiation a,mstrecorddifferentiation b WHERE b.clientid=? AND b.mstorgnhirarchyid=? AND b.seqno=? AND b.recorddifftypeid=3 AND  a.recorddiffid =b.id AND a.deleteflg=0 AND a.activeflg=1"
	var resolvestateID int64
	rows, err := mdao.DB.Query(sql, ClientID, OrgnID, ResolveStatusSeq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchResolvedStateID Get Statement Prepare Error", err)
		return resolvestateID, err
	}
	for rows.Next() {
		rows.Scan(&resolvestateID)

	}
	return resolvestateID, nil
}

func (mdao DbConn) FetchCurrentGrpIDForReopen(RecordID int64, ResolveStateID int64) (int64, error) {
	logger.Log.Println("In side FetchResolvedStateID")
	var sql = "SELECT a.mstgroupid FROM mstrequesthistory a,maprequestorecord b WHERE b.recordid=? AND a.mainrequestid=b.mstrequestid AND currentstateid != ? AND b.activeflg=1 AND b.deleteflg=0 AND a.activeflg=1 AND a.deleteflg=0 order by a.id desc limit 1"
	var GrpID int64
	rows, err := mdao.DB.Query(sql, RecordID, ResolveStateID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchResolvedStateID Get Statement Prepare Error", err)
		return GrpID, err
	}
	for rows.Next() {
		rows.Scan(&GrpID)

	}
	return GrpID, nil
}

func (mdao DbConn) GetSupportgrpdayofweekcount(ClientID int64, OrgnID int64, GrpID int64) (int64, error) {
	logger.Log.Println("In side GetSupportgrpdayofweekcount")
	var sql = "SELECT COUNT(id) AS daycount FROM mstclientsupportgroupdayofweek WHERE clientid=? AND mstorgnhirarchyid=? AND mstclientsupportgroupid=? AND deleteflg=0 AND activeflg=1"
	var daycount int64
	rows, err := mdao.DB.Query(sql, ClientID, OrgnID, GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetSupportgrpdayofweekcount Get Statement Prepare Error", err)
		return daycount, err
	}
	for rows.Next() {
		rows.Scan(&daycount)

	}
	return daycount, nil
}

func (mdao DbConn) GetOrganizationdayofweekcount(ClientID int64, OrgnID int64) (int64, error) {
	logger.Log.Println("In side GetOrganizationdayofweekcount")
	var sql = "SELECT COUNT(id) FROM mstclientdayofweek WHERE clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1"
	var daycount int64
	rows, err := mdao.DB.Query(sql, ClientID, OrgnID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOrganizationdayofweekcount Get Statement Prepare Error", err)
		return daycount, err
	}
	for rows.Next() {
		rows.Scan(&daycount)

	}
	return daycount, nil
}



func (mdao DbConn) GetLatestActivitylogSeq(RecordID int64) (int64, error) {
	logger.Log.Println("In side GetLatestActivitylogSeq")
	var sql = "SELECT activityseqno  FROM mstrecordactivitylogs where recordid=? AND activeflg=1 AND deleteflg=0 order by id desc limit 1"
	var seqID int64
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLatestActivitylogSeq Get Statement Prepare Error", err)
		return seqID, err
	}
	for rows.Next() {
		rows.Scan(&seqID)

	}
	return seqID, nil
}



func (mdao DbConn) InsertErrorByRecordID(ClientID int64, OrgnID int64, Actionname string, RecordID int64, Errordesc string) error {
	sql := "INSERT INTO msterrortracking(clientid,mstorgnhirarchyid,recordid,actionname,description) VALUES(?,?,?,?,?)"
	stmt, err := mdao.DB.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, OrgnID, RecordID, Actionname, Errordesc)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) InsertRecordClosure(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, RecordSeq int64, Closuredate string) error {
	sql := "INSERT INTO mstrecordautoclosure(clientid,mstorgnhirarchyid,recordid,recordseq,closuredate,closuredt) VALUES(?,?,?,?,round(UNIX_TIMESTAMP(?)),?)"
	stmt, err := mdao.DB.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, RecordSeq, Closuredate, Closuredate)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}



func (mdao DbConn) InsertRecordClosureForSR(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, RecordSeq int64) error {
	sql := "INSERT INTO mstrecordautoclosure(clientid,mstorgnhirarchyid,recordid,recordseq,closuredate,closuredt) VALUES(?,?,?,?,round(UNIX_TIMESTAMP(TIMESTAMPADD(MINUTE,4650,NOW()))),TIMESTAMPADD(MINUTE,4650,NOW()))"
	stmt, err := mdao.DB.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, RecordSeq)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}
