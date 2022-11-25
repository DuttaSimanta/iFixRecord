package dao

import (
	"database/sql"
	"errors"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

var catlevel = "SELECT distinct a.torecorddifftypeid as levelID FROM mstrecordtype a, mstrecorddifferentiationtype b WHERE a.torecorddifftypeid = b.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.fromrecorddifftypeid=? AND a.fromrecorddiffid=? AND b.parentid=1"
var updatelevel = "UPDATE maprecordtorecorddifferentiation SET islatest=0 WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=?"
var workingcatID = "SELECT recorddiffid as workingcatagoryID FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid =? AND isworking=1 order by id desc limit 1"
var lateststageval = "SELECT recordtitle,recorddescription FROM trnrecordstage WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND deleteflg=0 AND activeflg=1 ORDER BY id DESC LIMIT 1"
var categoryname = "SELECT distinct d.name FROM mstrecordtype a, mstrecorddifferentiationtype b,maprecordtorecorddifferentiation c,mstrecorddifferentiation d WHERE a.torecorddifftypeid = b.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.fromrecorddifftypeid=? AND a.fromrecorddiffid=? AND b.parentid=1 AND c.recordid=? AND c.islatest=1 AND a.torecorddifftypeid=c.recorddifftypeid and c.recorddiffid=d.id"
var diffname = "SELECT group_concat(name) as name FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1 AND id in(?)"
var oldvalue = "SELECT recordtrackvalue FROM trnreordtracking WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recordtermid=?"

func (mdao DbConn) GetcatlevelidagainstrecordID(ClientID int64, Mstorgnhirarchyid int64, Recorddifftypeid int64, Recorddiffid int64) ([]int64, error) {
	logger.Log.Println("In side GetcatlevelidagainstrecordID")
	logger.Log.Println("Query is ------>", catlevel)
	logger.Log.Println("Parameter is ------>", ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)

	var levelIDs []int64
	rows, err := mdao.DB.Query(catlevel, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetcatlevelidagainstrecordID Get Statement Prepare Error", err)
		return levelIDs, err
	}
	for rows.Next() {
		var levelID int64
		err = rows.Scan(&levelID)
		levelIDs = append(levelIDs, levelID)
		logger.Log.Println("GetcatlevelidagainstrecordID rows.next() Error", err)
	}
	return levelIDs, nil
}

func Updatepreviousrecord(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, levelID int64) error {
	logger.Log.Println("Updatepreviouscatlevel query -->", updatelevel)
	logger.Log.Println("Updatepreviouscatlevel parameters -->", ClientID, Mstorgnhirarchyid, RecordID, levelID)
	stmt, err := tx.Prepare(updatelevel)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, levelID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetLastWorkingcatID(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side GetLastWorkingcatID")
	logger.Log.Println("Query is ------>", workingcatID)
	logger.Log.Println("Parameter is ------>", ClientID, Mstorgnhirarchyid, RecordID)

	var workingcatagoryID int64
	rows, err := mdao.DB.Query(workingcatID, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLastWorkingcatID Get Statement Prepare Error", err)
		return workingcatagoryID, err
	}
	for rows.Next() {
		err = rows.Scan(&workingcatagoryID)
		logger.Log.Println("GetLastWorkingcatID rows.next() Error", err)
	}
	return workingcatagoryID, nil
}

func (mdao DbConn) GetLaststagevalue(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, string, error) {
	logger.Log.Println("In side GetLaststagevalue")
	logger.Log.Println("Query is ------>", lateststageval)
	logger.Log.Println("Parameter is ------>", ClientID, Mstorgnhirarchyid, RecordID)

	var recordtitle string
	var recorddescription string
	rows, err := mdao.DB.Query(lateststageval, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return recordtitle, recorddescription, err
	}
	for rows.Next() {
		err = rows.Scan(&recordtitle, &recorddescription)
		logger.Log.Println("GetLaststagevalue rows.next() Error", err)
	}
	return recordtitle, recorddescription, nil
}

func UpdateRecordStage(tx *sql.Tx, rec *entities.RecordcategoryupdateEntity, recordtitle string, recorddescription string) (int64, error) {
	logger.Log.Println("trnrecordstage query -->", trnrecordstage)
	logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, rec.RecordID, recordtitle, recorddescription, rec.UserID, rec.UsergroupID, rec.UserID, rec.UsergroupID)

	stmt, err := tx.Prepare(trnrecordstage)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, rec.RecordID, recordtitle, recorddescription, rec.UserID, rec.UsergroupID, rec.UserID, rec.UsergroupID)
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

func (mdao DbConn) Getcategorynames(ClientID int64, Mstorgnhirarchyid int64, RecorddifftypeID int64, RecorddifID int64, RecordID int64) (string, error) {
	logger.Log.Println("In side Getcategorynames")
	logger.Log.Println("Query is ------>", categoryname)
	logger.Log.Println("Parameter is ------>", ClientID, Mstorgnhirarchyid, RecorddifftypeID, RecorddifID, RecordID)

	var name string
	var catname string
	rows, err := mdao.DB.Query(categoryname, ClientID, Mstorgnhirarchyid, RecorddifftypeID, RecorddifID, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return catname, err
	}
	for rows.Next() {
		err = rows.Scan(&name)
		logger.Log.Println("GetLaststagevalue rows.next() Error", err)
		catname = catname + "," + name
	}
	return catname, nil
}

func Getdiffname(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, ids string) (string, error) {
	logger.Log.Println("In side Getcategorynames")
	logger.Log.Println("Query is ------>", diffname)
	logger.Log.Println("Parameter is ------>", ClientID, Mstorgnhirarchyid, ids[:len(ids)-1])

	var name string
	rows, err := tx.Query(diffname, ClientID, Mstorgnhirarchyid, ids[:len(ids)-1])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return name, err
	}
	for rows.Next() {
		err = rows.Scan(&name)
		logger.Log.Println("Getdiffname rows.next() Error", err)

	}
	return name, nil
}

func (mdao DbConn) Getadditionaloldvalue(Termid int64, ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, error) {
	logger.Log.Println("In side Getadditionaloldvalue")
	logger.Log.Println("Query is ------>", oldvalue)
	logger.Log.Println("Parameter is ------>", ClientID, Mstorgnhirarchyid, RecordID, Termid)

	var value string
	rows, err := mdao.DB.Query(oldvalue, ClientID, Mstorgnhirarchyid, RecordID, Termid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getadditionaloldvalue Get Statement Prepare Error", err)
		return value, err
	}
	for rows.Next() {
		err = rows.Scan(&value)
		logger.Log.Println("Getadditionaloldvalue rows.next() Error", err)

	}
	return value, nil
}

func (mdao DbConn) GetparentrecordInfo(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (string, string, string, string, string, int64, int64, int64, int64, error) {
	var query = "SELECT requestername,requesteremail,requestermobile,requesterlocation,source,userid,usergroupid,originaluserid,originalusergroupid FROM trnrecord WHERE clientid=? AND mstorgnhirarchyid=? AND id=? AND deleteflg=0 AND activeflg=1"
	var requestername string
	var requesteremail string
	var requestermobile string
	var requesterlocation string
	var source string
	var userid int64
	var usergroupid int64
	var originaluserid int64
	var originalusergroupid int64
	rows, err := mdao.DB.Query(query, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetparentrecordInfo Get Statement Prepare Error", err)
		return requestername, requesteremail, requestermobile, requesterlocation, source, userid, usergroupid, originaluserid, originalusergroupid, err
	}
	for rows.Next() {
		err = rows.Scan(&requestername, &requesteremail, &requestermobile, &requesterlocation, &source, &userid, &usergroupid, &originaluserid, &originalusergroupid)
		logger.Log.Println("GetparentrecordInfo rows.next() Error", err)

	}
	return requestername, requesteremail, requestermobile, requesterlocation, source, userid, usergroupid, originaluserid, originalusergroupid, nil
}

func (mdao DbConn) Checklatestlastcatcount(RecordID int64, ClientID int64, Mstorgnhirarchyid int64, Lastcatlevelid int64, lastcatvalueid int64) (int64, error) {
	logger.Log.Println(" Checklatestlastcatcount Parameter is ------>", ClientID, Mstorgnhirarchyid, RecordID, Lastcatlevelid, lastcatvalueid)
	var count int64
	var query = "SELECT count(*) countval FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=? AND recorddiffid=? AND islatest=1"
	rows, err := mdao.DB.Query(query, ClientID, Mstorgnhirarchyid, RecordID, Lastcatlevelid, lastcatvalueid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Checklatestlastcatcount Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("Checklatestlastcatcount rows.next() Error", err)
	}
	return count, nil
}



func (mdao DbConn) GetLatestTwologsIds(RecordId int64) ([]int64, error) {
	logger.Log.Println("In side GetLatestTwologsIds")
	logger.Log.Println("Query is ------>", catlevel)
	logger.Log.Println("Parameter is ------>", RecordId)

	var IDs []int64
	var query = "SELECT * FROM mstrecordactivitylogs where recordid=? AND id > (SELECT max(id) FROM mstrecordactivitylogs where recordid=? AND activityseqno =6)  order by id desc limit 2"
	rows, err := mdao.DB.Query(query, RecordId, RecordId)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLatestTwologsIds Get Statement Prepare Error", err)
		return IDs, err
	}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		IDs = append(IDs, id)
		logger.Log.Println("GetLatestTwologsIds rows.next() Error", err)
	}
	return IDs, nil
}

func (mdao DbConn) UpdateDeleteLogs(ID int64) error {
	var sql = "UPDATE mstrecordactivitylogs SET deleteflg=1 WHERE id=?"
	stmt, err := mdao.DB.Prepare(sql)
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
