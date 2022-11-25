package dao

import (
	"database/sql"
	"errors"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

var stageinsert = "INSERT INTO trnrecordstage (clientid, mstorgnhirarchyid,recordid,recordtitle, recorddescription,userid,usergroupid,entrydatetime,originaluserid,originalusergroupid) VALUES (?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?,?)"
var updatepriorityflag = "UPDATE maprecordtorecorddifferentiation SET islatest=0 WHERE id =?"
var latestID = "SELECT max(id) as ID,recorddiffid as Recorddiffid FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=? AND islatest=1"
var trnhistory = "Select mstslaentityid,donotupdatesladue,recordtimetoint,mstslastateid,commentonrecord,fromclientuserid,recorddatetime,recorddatetoint,clientid,mstorgnhirarchyid,therecordid,slastartstopindicator FROM trnslaentityhistory WHERE therecordid=? AND clientid=? AND mstorgnhirarchyid=? AND activeflg=1 AND deleteflg=0 order by id desc limit 1"
var getpreviouspriortynm = "SELECT a.name FROM mstrecorddifferentiation a,maprecordtorecorddifferentiation b WHERE b.id=? AND b.recorddiffid=a.id"
var getpriortynm = "SELECT name FROM mstrecorddifferentiation WHERE id=? "
var createdtquery = "SELECT FROM_UNIXTIME(createdatetime) as createdt  FROM trnrecord WHERE clientid=? AND mstorgnhirarchyid=? AND id=?"
var reopendatequery = "SELECT FROM_UNIXTIME(createddate) as createdt FROM maprecordtorecorddifferentiation WHERE recordid=? AND recorddifftypeid =3 AND recorddiffid in (SELECT id FROM iFIX.mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=3 AND seqno=10)"

//InsertTrnRecordStage data inserted in trnorderstage table
func InsertRecordStage(tx *sql.Tx, rec *entities.RecordpriorityEntity) (int64, error) {
	logger.Log.Println("InsertRecordStage query -->", stageinsert)
	logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, rec.Recordname, rec.Recordesc, rec.RecordID, rec.Userid, rec.Usergroupid, rec.Originaluserid, rec.Originalusergroupid)

	stmt, err := tx.Prepare(stageinsert)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, rec.RecordID, rec.Recordname, rec.Recordesc, rec.Userid, rec.Usergroupid, rec.Originaluserid, rec.Originalusergroupid)
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

func GetlatestDiffID(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, DifftypeID int64) (int64, int64, error) {
	logger.Log.Println("GetlatestpriorityID query -->", latestID)
	logger.Log.Println("GetlatestpriorityID parameters -->", ClientID, Mstorgnhirarchyid, RecordID)

	var ID int64
	var Recorddiffid int64
	stmt, error := tx.Prepare(latestID)
	if error != nil {
		logger.Log.Println("Exception in GetlatestpriorityID Prepare Statement..")
		return ID, Recorddiffid, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(ClientID, Mstorgnhirarchyid, RecordID, DifftypeID)
	if err != nil {
		logger.Log.Println("Exception in GetlatestpriorityID Query Statement..")
		return ID, Recorddiffid, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ID, &Recorddiffid); err != nil {
			logger.Log.Println("Error in fetching data")
		}

	}
	logger.Log.Println("GetlatestpriorityID  value is :", ID)
	return ID, Recorddiffid, nil
}

func Updateoldpriorityflag(tx *sql.Tx, ID int64) error {
	logger.Log.Println("Updateoldpriorityflag query -->", updatepriorityflag)
	logger.Log.Println("Updateoldpriorityflag parameters -->", ID)
	stmt, err := tx.Prepare(updatepriorityflag)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateoldpriorityflagNew(tx *sql.Tx, RecordID int64, DifftypeID int64) error {
	logger.Log.Println("Updateoldpriorityflag query -->", updatepriorityflag)
	logger.Log.Println("Updateoldpriorityflag parameters -->", RecordID)
	var sql = "UPDATE maprecordtorecorddifferentiation SET islatest=0 WHERE recordid =? AND recorddifftypeid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(RecordID, DifftypeID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func InsertRecordMapDifferrtiation(tx *sql.Tx, rec *entities.RecordpriorityEntity, lastInsertedStageID int64, seqno int64) (int64, error) {
	logger.Log.Println("maprecordtorecorddifferentiation query -->", maprecordtorecorddifferentiation)
	logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, lastInsertedStageID, seqno)

	stmt, err := tx.Prepare(maprecordtorecorddifferentiation)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, rec.RecordID, lastInsertedStageID, rec.Recorddifftypeid, rec.Recorddiffid, seqno)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Last ID fetch Error")
	}

	return id, nil
}

func (mdao DbConn) GetLatesttrnhistoryAll(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (entities.TrnslaentityhistoryEntity, error) {
	logger.Log.Println("In side GetLatesttrnhistory")
	tz := entities.TrnslaentityhistoryEntity{}
	rows, err := mdao.DB.Query(trnhistory, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLatesttrnhistory Get Statement Prepare Error", err)
		return tz, err
	}
	for rows.Next() {
		rows.Scan(&tz.Mstslaentityid, &tz.Donotupdatesladue, &tz.Recordtimetoint, &tz.Mstslastateid, &tz.Commentonrecord, &tz.Fromclientuserid, &tz.Recorddatetime, &tz.Recorddatetoint, &tz.Clientid, &tz.Mstorgnhirarchyid, &tz.Therecordid, &tz.Slastartstopindicator)

	}
	return tz, nil
}

func Getpriorityname(tx *sql.Tx, ID int64) (string, error) {
	logger.Log.Println("Getpriorityname query -->", getpriortynm)
	logger.Log.Println("Getpriorityname parameters -->", ID)

	var name string
	stmt, error := tx.Prepare(getpriortynm)
	if error != nil {
		logger.Log.Println("Exception in Getpriorityname Prepare Statement..")
		return name, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(ID)
	if err != nil {
		logger.Log.Println("Exception in Getpriorityname Query Statement..")
		return name, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			logger.Log.Println("Error in fetching data")
		}

	}
	logger.Log.Println("Getpriorityname  value is :", name)
	return name, nil
}

func GetPreviouspriorityname(tx *sql.Tx, ID int64) (string, error) {
	logger.Log.Println("Getpriorityname query -->", getpreviouspriortynm)
	logger.Log.Println("Getpriorityname parameters -->", ID)

	var name string
	stmt, error := tx.Prepare(getpreviouspriortynm)
	if error != nil {
		logger.Log.Println("Exception in Getpriorityname Prepare Statement..")
		return name, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(ID)
	if err != nil {
		logger.Log.Println("Exception in Getpriorityname Query Statement..")
		return name, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			logger.Log.Println("Error in fetching data")
		}

	}
	logger.Log.Println("Getpriorityname  value is :", name)
	return name, nil
}

func (mdao DbConn) Getrecordcreatedate(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, error) {
	logger.Log.Println("In side Getrecordcreatedate")
	var createdt string
	rows, err := mdao.DB.Query(createdtquery, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecordcreatedate Get Statement Prepare Error", err)
		return createdt, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&createdt)

	}
	return createdt, nil
}

func (mdao DbConn) Getrecordreopencreatedate(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, error) {
	logger.Log.Println("In side Getrecordreopencreatedate")
	var createdt string
	rows, err := mdao.DB.Query(reopendatequery, RecordID, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getrecordreopencreatedate Get Statement Prepare Error", err)
		return createdt, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&createdt)

	}
	return createdt, nil
}
