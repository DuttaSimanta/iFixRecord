package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
	"strings"
	"unicode"
)

var trnrecord = "INSERT INTO trnrecord (clientid, mstorgnhirarchyid, recordtitle, recorddescription,requesterinfo,userid,usergroupid,createdatetime,originaluserid,originalusergroupid,requestername,requesteremail,requestermobile,requesterlocation,source,code) VALUES (?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?,?,?,?,?,?,?,?)"
var trnrecordstage = "INSERT INTO trnrecordstage (clientid, mstorgnhirarchyid,recordid,recordtitle, recorddescription,userid,usergroupid,entrydatetime,originaluserid,originalusergroupid) VALUES (?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?,?)"
var maprecordtorecorddifferentiation = "INSERT INTO maprecordtorecorddifferentiation (clientid, mstorgnhirarchyid,recordid, recordstageid, recorddifftypeid, recorddiffid,seqno,createddate) VALUES (?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())))"
var trnreordtracking = "INSERT INTO trnreordtracking (clientid, mstorgnhirarchyid,recordid, recordstageid,recordtermid, recordtrackvalue,createdbyid,createddate,createdgrpid,recordtrackdescription) VALUES (?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?,?)"

var mstrecorddifferentiationtype = "SELECT id FROM mstrecorddifferentiationtype WHERE seqno=1 AND activeflg=1 AND deleteflg=0"
var getrecorddiffid = "SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid=? AND recorddifftypeid=? "
var mstrecordconfig = "SELECT prefix,year,month,day,configurezero FROM mstrecordconfig WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND deleteflg=0 AND activeflg=1 "

var updatetrnrecord = "UPDATE trnrecord SET code=? WHERE id=?"
var difftypemap = "SELECT id,typename FROM mstrecorddifferentiationtype WHERE clientid=? AND mstorgnhirarchyid=? AND activeflg=1 AND deleteflg=0"
var getnumber = "SELECT number FROM mstrecordautoincreament WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND activeflg=1 AND deleteflg=0"
var updatenumber = "UPDATE mstrecordautoincreament SET number = (number+1) WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND activeflg=1 AND deleteflg=0"
var recordasset = "INSERT INTO maprecordasset(clientid,mstorgnhirarchyid,recordid,recordstageid,assetid) VALUES(?,?,?,?,?)"
var recordadditional = "INSERT INTO trnreordtracking (clientid, mstorgnhirarchyid,recordid, recordstageid,recordtermid, recordtrackvalue,referenceid,referencetype,createdbyid,createddate,createdgrpid) VALUES (?,?,?,?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?)"
var updateworkinglabel = "UPDATE maprecordtorecorddifferentiation SET isworking=1 WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=? AND recordstageid=?"
var getworkingdiffid = "SELECT recorddiffid as workingdiffid FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=? AND recordstageid=?"
var updatestageid = "UPDATE trnrecord SET recordstageid=? WHERE id=?"
var gettermtype = "select (case when termtypeid=3 then 1 else 0 end) termtype FROM mstrecordterms where clientid=? AND mstorgnhirarchyid=? AND id=?"
var insertlogs = "INSERT INTO mstrecordactivitylogs(clientid,mstorgnhirarchyid,recordid,activityseqno,logValue,createdid,createddate,createdgrpid) VALUES (?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?)"
var insertlogswithgenericID = "INSERT INTO mstrecordactivitylogs(clientid,mstorgnhirarchyid,recordid,activityseqno,logValue,createdid,createddate,createdgrpid,genericid) VALUES (?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?,?)"
var updateordertracking = "UPDATE trnreordtracking SET recordtrackvalue=?,referenceid=? WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recordtermid=?"
var termsequance = "SELECT seq FROM mstrecordterms WHERE clientid=? AND mstorgnhirarchyid=? AND id=? AND deleteflg=0 AND activeflg=1"
var gettaskbycatid = "SELECT id,torecorddiffid from mstrecordtype where fromrecorddiffid=? and activeflg=1 and deleteflg=0"
var taskdetails = "SELECT title,description from msttaskpropertyvalue where recordtypeid =? and activeflg=1 and deleteflg=0"
var tickettypebycatid = "SELECT fromrecorddifftypeid,fromrecorddiffid FROM iFIX.mstrecordtype where torecorddiffid=? and fromrecorddifftypeid in (SELECT id from mstrecorddifferentiationtype where seqno=1 and activeflg=1 and deleteflg=0) and activeflg=1 and deleteflg=0;"
var difftypebyid = "SELECT recorddifftypeid from mstrecorddifferentiation where id=?"
var getparentsbycat = "SELECT parentcategoryids from mstrecorddifferentiation where id=?"
var getProcessByCategory = "SELECT mstprocessid as Processid from mstprocessrecordmap where clientid=? and mstorgnhirarchyid=? and recorddifftypeid = ? and recorddiffid =? and activeflg=1 and deleteflg=0 "
var isapproved = "SELECT count(a.id) from maprecordstatetodifferentiation b,msttransition a ,mstrecorddifferentiation c where b.clientid =? and b.mstorgnhirarchyid=? and b.activeflg=1 and b.deleteflg=0 and b.recorddiffid =c.id and (b.mststateid=a.currentstateid or b.mststateid=a.previousstateid) and a.processid=? and a.activeflg=1 and a.deleteflg=0  and c.seqno in(12,21)  and c.activeflg=1 and c.deleteflg=0"
var getdiffdetails = "SELECT recorddifftypeid,id from mstrecorddifferentiation where clientid=? and mstorgnhirarchyid=? and recorddifftypeid in (SELECT id from mstrecorddifferentiationtype where seqno=? and activeflg=1 and deleteflg=0 )  and seqno=? and  activeflg=1 and deleteflg=0"
var getseqbyid = "SELECT seqno from mstrecorddifferentiation where id =? and activeflg=1 and deleteflg=0"
var updateapproval = "UPDATE trnrecord SET isapproveworkflow=1 where id=?"

func (mdao DbConn) Updateapprovalstatus(ticketId int64) error {
	// logger.Log.Print("Updatestagingdetails ", ticketId)
	// log.Print("Updatestagingdetails", ticketId)

	stmt, err := mdao.DB.Prepare(updateapproval)
	if err != nil {
		logger.Log.Print("Updateapprovalstatus Prepare Statement  Error", err)
		log.Print("Updateapprovalstatus Prepare Statement  Error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ticketId)
	if err != nil {
		logger.Log.Print("Updateapprovalstatus Execute Statement  Error", err)
		log.Print("Updateapprovalstatus Execute Statement  Error", err)
		return err
	}
	return nil
}
func (mdao DbConn) Getseqbyid(clientid int64, orgid int64, id int64) (int64, error) {
	var seq int64
	stmt, error := mdao.DB.Prepare(getseqbyid)
	if error != nil {
		logger.Log.Println("Exception in Getseqbyid Prepare Statement..", error)
		return 0, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		logger.Log.Println("Exception in Getseqbyid Query Statement..", err)
		return 0, error
	}
	defer rows.Close()
	for rows.Next() {

		if err := rows.Scan(&seq); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}

	}
	return seq, nil
}
func (mdao DbConn) Getdiffdetailsbyseq(clientid int64, orgid int64, seq int64, typeseq int64) ([]entities.RecordSet, error) {
	taskcatids := []entities.RecordSet{}
	stmt, error := mdao.DB.Prepare(getdiffdetails)
	if error != nil {
		logger.Log.Println("Exception in Getdiffdetailsbyseq Prepare Statement..")
		return nil, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, orgid, typeseq, seq)
	if err != nil {
		logger.Log.Println("Exception in Getdiffdetailsbyseq Query Statement..")
		return nil, error
	}
	defer rows.Close()
	for rows.Next() {
		taskcatid := entities.RecordSet{}
		if err := rows.Scan(&taskcatid.ID, &taskcatid.Val); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}
		taskcatids = append(taskcatids, taskcatid)
	}
	return taskcatids, nil
}
func (dbc DbConn) Checkisapprovedprocess(tz *entities.RecordcategoryupdateEntity, processid int64) (entities.RecordcategoryupdateEntity, error) {
	//logger.Log.Println("In side GetClientsupportgroupnewCount")
	value := entities.RecordcategoryupdateEntity{}
	//logger.Log.Println("Query Checkisapprovedprocess ==========================================>", isapproved)
	//logger.Log.Println("Query Checkisapprovedprocess ==========================================>", tz.ClientID, tz.Mstorgnhirarchyid, processid)
	err := dbc.DB.QueryRow(isapproved, tz.ClientID, tz.Mstorgnhirarchyid, processid).Scan(&value.Total)
	switch err {
	case sql.ErrNoRows:
		value.Total = 0
		return value, nil
	case nil:
		return value, nil
	default:
		logger.Log.Println("Checkisapprovedprocess Get Statement Prepare Error", err)
		return value, err
	}
}
func (mdao DbConn) GetProcessByCategory(tz *entities.RecordcategoryupdateEntity) (entities.RecordData, error) {
	log.Println("In side dao")
	value := entities.RecordData{}
	err := mdao.DB.QueryRow(getProcessByCategory, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid).Scan(&value.ID)
	switch err {
	case sql.ErrNoRows:
		value.ID = 0
		return value, nil
	case nil:
		return value, nil
	default:
		logger.Log.Println("GetProcessByCategory Get Statement Prepare Error", err)
		return value, err
	}
}

func (mdao DbConn) Getdifftypebyid(id int64) (int64, error) {
	var difftype int64
	stmt, error := mdao.DB.Prepare(difftypebyid)
	if error != nil {
		logger.Log.Println("Exception in Getdifftypebyid Prepare Statement..")
		return 0, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		logger.Log.Println("Exception in Getdifftypebyid Query Statement..")
		return 0, error
	}
	defer rows.Close()
	for rows.Next() {
		//taskcatid :=entities.RecordData{}
		if err := rows.Scan(&difftype); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}
	}
	return difftype, nil
}
func (mdao DbConn) Getparentsbycatid(id int64) (string, error) {
	log.Print("parentcat", id)
	var parentcat string
	stmt, error := mdao.DB.Prepare(getparentsbycat)
	if error != nil {
		logger.Log.Println("Exception in Getparentsbycatid Prepare Statement..")
		return "", error
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		logger.Log.Println("Exception in Getparentsbycatid Query Statement..")
		return "", error
	}
	defer rows.Close()
	for rows.Next() {
		//taskcatid :=entities.RecordData{}
		if err := rows.Scan(&parentcat); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}
	}
	return parentcat, nil
}
func (mdao DbConn) Gettickettypebycatid(id int64) ([]entities.RecordSet, error) {
	log.Print("id:", id)
	taskcatids := []entities.RecordSet{}
	stmt, error := mdao.DB.Prepare(tickettypebycatid)
	if error != nil {
		logger.Log.Println("Exception in Gettickettypebycatid Prepare Statement..")
		return nil, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		logger.Log.Println("Exception in Gettickettypebycatid Query Statement..")
		return nil, error
	}
	defer rows.Close()
	for rows.Next() {
		taskcatid := entities.RecordSet{}
		if err := rows.Scan(&taskcatid.ID, &taskcatid.Val); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}
		taskcatids = append(taskcatids, taskcatid)
	}
	return taskcatids, nil
}
func (mdao DbConn) Gettaskdetailsbyid(id int64) ([]entities.TaskdetailsEntity, error) {
	taskcatids := []entities.TaskdetailsEntity{}
	stmt, error := mdao.DB.Prepare(taskdetails)
	if error != nil {
		logger.Log.Println("Exception in Gettaskdetailsbyid Prepare Statement..")
		return nil, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		logger.Log.Println("Exception in Gettaskdetailsbyid Query Statement..")
		return nil, error
	}
	defer rows.Close()
	for rows.Next() {
		taskcatid := entities.TaskdetailsEntity{}
		if err := rows.Scan(&taskcatid.Title, &taskcatid.Desc); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}
		taskcatids = append(taskcatids, taskcatid)
	}
	return taskcatids, nil
}
func (mdao DbConn) Gettaskbycatid(catid int64) ([]entities.RecordData, error) {
	taskcatids := []entities.RecordData{}
	stmt, error := mdao.DB.Prepare(gettaskbycatid)
	if error != nil {
		logger.Log.Println("DB Exception ......", error)
		return nil, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(catid)
	if err != nil {
		logger.Log.Println("Exception in Gettaskbycatid Query Statement..", err)
		return nil, error
	}
	defer rows.Close()
	for rows.Next() {
		taskcatid := entities.RecordData{}
		if err := rows.Scan(&taskcatid.ID, &taskcatid.Val); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}
		taskcatids = append(taskcatids, taskcatid)

	}
	return taskcatids, nil
}

//InsertTrnRecord data insertd in trnorder table
func InsertTrnRecord(tx *sql.Tx, rec *entities.RecordEntity, lognumber string, ticketID string) (int64, error) {
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("trnrecord query -->", trnrecord)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, rec.Recordname, rec.Recordesc, rec.CreateduserID, rec.CreatedusergroupID, rec.Originaluserid, rec.Originalusergroupid, rec.Requestername, rec.Requesteremail, rec.Requestermobile, rec.Requesterlocation, rec.Source)
	logger.Log.Println("Record description value is   *************************************88 ************************** -->", rec.Recordesc)
	
	longdesc := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, rec.Recordesc)
	stmt, err := tx.Prepare(trnrecord)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, rec.Recordname,strings.TrimSpace(longdesc), rec.Requesterinfo, rec.CreateduserID, rec.CreatedusergroupID, rec.Originaluserid, rec.Originalusergroupid, rec.Requestername, rec.Requesteremail, rec.Requestermobile, rec.Requesterlocation, rec.Source, ticketID)
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

//InsertTrnRecordStage data inserted in trnorderstage table
func InsertTrnRecordStage(tx *sql.Tx, rec *entities.RecordEntity, lastInsertedID int64, lognumber string) (int64, error) {
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("trnrecordstage query -->", trnrecordstage)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, rec.Recordname, rec.Recordesc, lastInsertedID, rec.CreateduserID, rec.CreatedusergroupID, rec.Originaluserid, rec.Originalusergroupid)
	
	longdesc := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, rec.Recordesc)
	stmt, err := tx.Prepare(trnrecordstage)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, lastInsertedID, rec.Recordname, strings.TrimSpace(longdesc), rec.CreateduserID, rec.CreatedusergroupID, rec.Originaluserid, rec.Originalusergroupid)
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

//InsertTrnRecordMapDifferrtiation data inserted in maprecordtorecorddifferentiation table
func InsertTrnRecordMapDifferrtiation(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, lastInsertedID int64, lastInsertedStageID int64, difftypeID int64, diffid int64, seqno int64, lognumber string) error {
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("maprecordtorecorddifferentiation query -->", maprecordtorecorddifferentiation)
	// logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, lastInsertedID, lastInsertedStageID, difftypeID, diffid, seqno)

	stmt, err := tx.Prepare(maprecordtorecorddifferentiation)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, lastInsertedID, lastInsertedStageID, difftypeID, diffid, seqno)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

//InsertRecordTermvalues data inserted in trnreordtracking table
func InsertRecordTermvalues(tx *sql.Tx, rec *entities.RecordEntity, lastInsertedID int64, lastInsertedStageID int64, createdbyID int64, termID int64, val string, val1 string, lognumber string) error {
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("trnreordtracking query -->", trnreordtracking)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, lastInsertedID, lastInsertedStageID, termID, val, createdbyID, rec.CreatedusergroupID)

	stmt, err := tx.Prepare(trnreordtracking)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, lastInsertedID, lastInsertedStageID, termID, val, createdbyID, rec.CreatedusergroupID, val1)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

/** Without transaction****/
//GetRecorddifftypeID type id from mstrecorddifferentiationtype table
func (mdao DbConn) GetRecorddifftypeIDnormal(rec *entities.RecordEntity, lastInsertedID int64, lognumber string) (int64, error) {
	var recorddifftypeID int64
	stmt, error := mdao.DB.Prepare(mstrecorddifferentiationtype)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return recorddifftypeID, error
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		logger.Log.Println("Exception in GetTickeType Query Statement..")
		return recorddifftypeID, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&recorddifftypeID); err != nil {
			logger.Log.Println("Error in fetching data")
		}

	}
	fmt.Println("tickettypeID value is :", recorddifftypeID)
	return recorddifftypeID, nil
}

//GetRecorddiffID type id from maprecordtorecorddifferentiation table
func (mdao DbConn) GetRecorddiffIDnormal(lastInsertedID int64, recorddifftypeid int64, lognumber string) (int64, error) {
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("recorddiffid query -->", getrecorddiffid)
	// logger.Log.Println("GetRecorddiffIDnormal-------------------------------- -->", lastInsertedID, recorddifftypeid)

	var recorddiffid int64
	stmt, error := mdao.DB.Prepare(getrecorddiffid)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return recorddiffid, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(lastInsertedID, recorddifftypeid)
	if err != nil {
		logger.Log.Println(err)
		return recorddiffid, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&recorddiffid); err != nil {
			logger.Log.Println(err)
		}

	}
	logger.Log.Println("recorddiffid value is :", recorddiffid)
	return recorddiffid, nil
}

//GetPrefix type id from maprecordtorecorddifferentiation table
func (mdao DbConn) GetPrefixnormal(rec *entities.RecordEntity, recorddifftypeID int64, recorddiffID int64, lognumber string) ([]string, error) {
	var prefixarr []string
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("mstrecordconfig query -->", mstrecordconfig)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, recorddifftypeID, recorddiffID)
	stmt, error := mdao.DB.Prepare(mstrecordconfig)
	if error != nil {
		logger.Log.Println("Exception in GetPrefix Prepare Statement..")
		return prefixarr, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(rec.ClientID, rec.Mstorgnhirarchyid, recorddifftypeID, recorddiffID)
	if err != nil {
		logger.Log.Println("Exception in GetPrefix Query Statement..")
		return prefixarr, error
	}
	defer rows.Close()
	if rows.Next() {
		var prefix string
		var year string
		var month string
		var day string
		var configurezero string

		if err := rows.Scan(&prefix, &year, &month, &day, &configurezero); err != nil {
			logger.Log.Println("Error in fetching data")
		}
		prefixarr = append(prefixarr, prefix)
		prefixarr = append(prefixarr, year)
		prefixarr = append(prefixarr, month)
		prefixarr = append(prefixarr, day)
		prefixarr = append(prefixarr, configurezero)
	}
	logger.Log.Println("prefixarr value is :", prefixarr)
	return prefixarr, nil
}

//GetRecorddifftypeID type id from mstrecorddifferentiationtype table
func GetRecorddifftypeID(tx *sql.Tx, rec *entities.RecordEntity, lastInsertedID int64, lognumber string) (int64, error) {
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("mstrecorddifferentiationtype query -->", mstrecorddifferentiationtype)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid)

	var recorddifftypeID int64
	stmt, error := tx.Prepare(mstrecorddifferentiationtype)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return recorddifftypeID, error
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		logger.Log.Println("Exception in GetTickeType Query Statement..")
		return recorddifftypeID, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&recorddifftypeID); err != nil {
			logger.Log.Println("Error in fetching data")
		}

	}
	fmt.Println("tickettypeID value is :", recorddifftypeID)
	return recorddifftypeID, nil
}

//GetRecorddiffID type id from maprecordtorecorddifferentiation table
func GetRecorddiffID(tx *sql.Tx, lastInsertedID int64, recorddifftypeid int64, lognumber string) (int64, error) {
	// logger.Log.Println("recorddiffid query -->", getrecorddiffid)
	// logger.Log.Println("parameters -->", lastInsertedID, recorddifftypeid)

	var recorddiffID int64
	stmt, error := tx.Prepare(getrecorddiffid)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return recorddiffID, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(lastInsertedID, recorddifftypeid)
	if err != nil {
		logger.Log.Println("Exception in GetTickeType Query Statement..")
		return recorddiffID, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&recorddiffID); err != nil {
			logger.Log.Println("Error in fetching data")
		}

	}
	fmt.Println("recorddiffID value is :", recorddiffID)
	return recorddiffID, nil
}

//GetPrefix type id from maprecordtorecorddifferentiation table
func GetPrefix(tx *sql.Tx, rec *entities.RecordEntity, recorddifftypeID int64, recorddiffID int64, lognumber string) ([]string, error) {
	var prefix []string
	// logger.Log.Println("Transaction log number -->", lognumber)
	// logger.Log.Println("mstrecordconfig query -->", mstrecordconfig)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, recorddifftypeID, recorddiffID)
	stmt, error := tx.Prepare(mstrecordconfig)
	if error != nil {
		logger.Log.Println("Exception in GetPrefix Prepare Statement..")
		return prefix, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(rec.ClientID, rec.Mstorgnhirarchyid, recorddifftypeID, recorddiffID)
	if err != nil {
		logger.Log.Println("Exception in GetPrefix Query Statement..")
		return prefix, error
	}
	defer rows.Close()
	if rows.Next() {
		var prefix1 string
		var prefix2 string
		var prefix3 string

		if err := rows.Scan(&prefix1, &prefix2, &prefix3); err != nil {
			logger.Log.Println("Error in fetching data")
		}
		prefix = append(prefix, prefix1)
		prefix = append(prefix, prefix2)
		prefix = append(prefix, prefix3)

	}
	fmt.Println("recorddiffID value is :", recorddiffID)
	return prefix, nil
}

//UpdateRecordID update record number table
func UpdateRecordID(tx *sql.Tx, lastInsertedID int64, recordno string) error {
	// logger.Log.Println("updatetrnrecord query -->", updatetrnrecord)
	// logger.Log.Println("parameters -->", recordno, lastInsertedID)

	stmt, err := tx.Prepare(updatetrnrecord)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(recordno, lastInsertedID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

//GetClientWiseDifferationType to get differationtype map client wise
func GetClientWiseDifferationType(dbcon *sql.DB, rec *entities.RecordEntity, lognumber string) (map[string]int64, error) {
	//logger.Log.Println("Transaction log number -->", lognumber)
	m := make(map[string]int64)
	tx, err := dbcon.Begin()
	if err != nil {
		dbcon.Close()
		logger.Log.Println("Transaction creation error in GetClientWiseDifferationType", err)
		return m, err
	}

	stmt, error := tx.Prepare(difftypemap)
	if error != nil {
		logger.Log.Println("Exception in GetClientWiseDifferationType Prepare Statement..")
		return m, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(rec.ClientID, rec.Mstorgnhirarchyid)
	if err != nil {
		logger.Log.Println("Exception in GetClientWiseDifferationType Query Statement..")
		return m, error
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var typename string
		if err := rows.Scan(&id, &typename); err != nil {
			logger.Log.Println("Error in fetching data")
		}
		m[typename] = id
	}
	//fmt.Println("Map values --->", m)

	return m, nil
}

//GetDiffID method is used for getting recorddiffid against recordid & recorddifftypeid
func GetDiffID(dbcon *sql.DB, lastInsertedID int64, recorddifftypeid int64, lognumber string) (int64, error) {
	//logger.Log.Println("Transaction log number -->", lognumber)
	tx, err := dbcon.Begin()
	if err != nil {
		logger.Log.Println("Transaction creation error in GetDiffID", err)
		return 0, err
	}

	recorddiffid, err := GetRecorddiffID(tx, lastInsertedID, recorddifftypeid, lognumber)
	if err != nil {
		logger.Log.Println("Error of fetching data in GetDiffID", err)
		return 0, err
	}
	return recorddiffid, nil
}

//GetRecorddiffID type id from maprecordtorecorddifferentiation table
func (mdao DbConn) GetNumber(clientid int64, mtorgnhirarchyid int64, difftypeID int64, diffid int64) (int, error) {
	// logger.Log.Println("recorddiffid query -->", getnumber)
	// logger.Log.Println("GetRecorddiffIDnormal-------------------------------- -->", clientid, mtorgnhirarchyid, difftypeID, diffid)

	var number int
	stmt, error := mdao.DB.Prepare(getnumber)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return number, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, mtorgnhirarchyid, difftypeID, diffid)
	if err != nil {
		logger.Log.Println(err)
		return number, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&number); err != nil {
			logger.Log.Println(err)
		}

	}
	logger.Log.Println("number value is :", number)
	return number, nil
}

func (mdao DbConn) Updatenumber(clientid int64, mtorgnhirarchyid int64, difftypeID int64, diffid int64) error {
	stmt, err := mdao.DB.Prepare(updatenumber)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientid, mtorgnhirarchyid, difftypeID, diffid)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

// InsertRecordAsset is used for insert asset data
func InsertRecordAsset(tx *sql.Tx, rec *entities.RecordEntity, assetID int64, insertedID int64, lastInsertedStageID int64) error {
	// logger.Log.Println("InsertRecordAsset query -->", recordasset)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, assetID, insertedID, lastInsertedStageID)

	stmt, err := tx.Prepare(recordasset)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, insertedID, lastInsertedStageID, assetID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

// InsertRecordAdditional is used for insert additional data
func InsertRecordAdditional(tx *sql.Tx, rec *entities.RecordEntity, additionalfields *entities.RecordAdditional, insertedID int64, lastInsertedStageID int64) error {
	// logger.Log.Println("InsertRecordAsset query -->", recordasset)
	// logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, insertedID, lastInsertedStageID, additionalfields.Termsid, additionalfields.Val, additionalfields.ID, "Additional")

	stmt, err := tx.Prepare(recordadditional)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, insertedID, lastInsertedStageID, additionalfields.Termsid, additionalfields.Val, additionalfields.ID, "Additional", rec.CreateduserID, rec.CreatedusergroupID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func InsertExtraRecordID(tx *sql.Tx, clientid int64, mtorgnhirarchyid int64, insertedID int64, extrarecordid int64) error {
	var sql = "INSERT INTO maprecordtolinkrecords(clientid,mstorgnhirarchyid,recordid,linkrecordid) VALUES (?,?,?,?)"
	stmt, err := tx.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientid, mtorgnhirarchyid, insertedID, extrarecordid)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func Updateisworking(tx *sql.Tx, clientid int64, mtorgnhirarchyid int64, insertedID int64, lastInsertedStageID int64, workinglabel int64) error {
	stmt, err := tx.Prepare(updateworkinglabel)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientid, mtorgnhirarchyid, insertedID, workinglabel, lastInsertedStageID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

//GetRecorddiffID type id from maprecordtorecorddifferentiation table
func (mdao DbConn) GetWorkingdiffid(clientid int64, mtorgnhirarchyid int64, insertedID int64, lastInsertedStageID int64, Workingcatlabelid int64) (int64, error) {
	// logger.Log.Println("recorddiffid query -->", getworkingdiffid)
	// logger.Log.Println("GetRecorddiffIDnormal-------------------------------- -->", clientid, mtorgnhirarchyid, insertedID, lastInsertedStageID, Workingcatlabelid)

	var workingdiffid int64
	stmt, error := mdao.DB.Prepare(getworkingdiffid)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return workingdiffid, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, mtorgnhirarchyid, insertedID, lastInsertedStageID, Workingcatlabelid)
	if err != nil {
		logger.Log.Println(err)
		return workingdiffid, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&workingdiffid); err != nil {
			logger.Log.Println(err)
		}

	}
	logger.Log.Println("workingdiffid value is :", workingdiffid)
	return workingdiffid, nil
}

//GetRecorddiffID type id from maprecordtorecorddifferentiation table
func GetWorkingdiffid(tx *sql.Tx, clientid int64, mtorgnhirarchyid int64, insertedID int64, lastInsertedStageID int64, Workingcatlabelid int64) (int64, error) {
	// logger.Log.Println("recorddiffid query -->", getworkingdiffid)
	// logger.Log.Println("GetRecorddiffIDnormal-------------------------------- -->", clientid, mtorgnhirarchyid, insertedID, lastInsertedStageID, Workingcatlabelid)

	var workingdiffid int64
	stmt, error := tx.Prepare(getworkingdiffid)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return workingdiffid, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, mtorgnhirarchyid, insertedID, Workingcatlabelid, lastInsertedStageID)
	if err != nil {
		logger.Log.Println(err)
		return workingdiffid, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&workingdiffid); err != nil {
			logger.Log.Println(err)
		}

	}
	logger.Log.Println("workingdiffid value is :", workingdiffid)
	return workingdiffid, nil
}

func Updatestageid(tx *sql.Tx, insertedID int64, lastInsertedStageID int64) error {
	// logger.Log.Println("Updatestageid query -->", updateworkinglabel)
	// logger.Log.Println("Updatestageid parameters -->", insertedID, lastInsertedStageID)
	stmt, err := tx.Prepare(updatestageid)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(lastInsertedStageID, insertedID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

//GetRecordtermtype type id from mstrecorddifferentiationtype table
func (mdao DbConn) GetRecordtermtype(ClientID int64, Mstorgnhirarchyid int64, Termid int64) (int64, error) {
	// logger.Log.Println("mstrecorddifferentiationtype query -->", mstrecorddifferentiationtype)
	// logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, Termid)

	var termtype int64
	stmt, error := mdao.DB.Prepare(gettermtype)
	if error != nil {
		logger.Log.Println("Exception in GetRecordtermtype Prepare Statement..")
		return termtype, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(ClientID, Mstorgnhirarchyid, Termid)
	if err != nil {
		logger.Log.Println("Exception in GetRecordtermtype Query Statement..")
		return termtype, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&termtype); err != nil {
			logger.Log.Println("Error in fetching data")
		}

	}
	fmt.Println("tickettypeID value is :", termtype)
	return termtype, nil
}

func (mdao DbConn) GetRecordtermSequance(ClientID int64, Mstorgnhirarchyid int64, Termid int64) (int64, error) {
	// logger.Log.Println("mstrecorddifferentiationtype query -->", mstrecorddifferentiationtype)
	// logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, Termid)

	var termseq int64
	stmt, error := mdao.DB.Prepare(termsequance)
	if error != nil {
		logger.Log.Println("Exception in GetRecordtermtype Prepare Statement..")
		return termseq, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(ClientID, Mstorgnhirarchyid, Termid)
	if err != nil {
		logger.Log.Println("Exception in GetRecordtermtype Query Statement..", err)
		return termseq, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&termseq); err != nil {
			logger.Log.Println("Error in fetching data", err)
		}

	}
	fmt.Println("termsequance value is :", termseq)
	return termseq, nil
}

//InsertActivityLogs used for insert data into log table
func (mdao DbConn) InsertActivityLogs(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, ActivityID int64, Logvalue string, CreatedID int64, CreatedgrpID int64, TermID int64) error {
	// logger.Log.Println("InsertActivityLogs query -->", insertlogs)
	// logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID, TermID)

	stmt, err := mdao.DB.Prepare(insertlogswithgenericID)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID, TermID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) InsertRecordActivityLogs(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, ActivityID int64, Logvalue string, CreatedID int64, CreatedgrpID int64) error {
	// logger.Log.Println("InsertActivityLogs query -->", insertlogs)
	// logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID)

	stmt, err := mdao.DB.Prepare(insertlogs)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func InsertActivityLogs(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, ActivityID int64, Logvalue string, CreatedID int64, CreatedgrpID int64) error {
	 logger.Log.Println("InsertActivityLogs query -->", insertlogs)
         logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID)

	stmt, err := tx.Prepare(insertlogs)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func InsertActivityLogsfromterms(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, ActivityID int64, Logvalue string, CreatedID int64, CreatedgrpID int64, TermID int64) error {
	// logger.Log.Println("InsertActivityLogs query -->", insertlogs)
	// logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID, TermID)

	stmt, err := tx.Prepare(insertlogswithgenericID)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID, TermID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateRecordAdditional(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, ReferenceID int64, TermsID int64, Termsvalue string, RecordID int64) error {
	// logger.Log.Println("Updateoldpriorityflag query -->", updateordertracking)
	// logger.Log.Println("Updateoldpriorityflag parameters -->", Termsvalue, ReferenceID, ClientID, Mstorgnhirarchyid, RecordID, TermsID)
	stmt, err := tx.Prepare(updateordertracking)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Termsvalue, ReferenceID, ClientID, Mstorgnhirarchyid, RecordID, TermsID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func GetImpactUrgencydtls(tx *sql.Tx, ClientID int64, OrgnID int64, ReocrdType int64, PriorityID int64) (int64, int64, error) {
	//logger.Log.Println("GetImpactUrgencydtls parameters --------------------------222222222222222222222222222222222222222------------------------------->", ClientID, OrgnID, ReocrdType, PriorityID)
	var impactID int64
	var urgencyID int64
	var sql = "SELECT mstrecorddifferentiationimpactid,mstrecorddifferentiationurgencyid FROM mstbusinessmatrix WHERE clientid=? AND mstorgnhirarchyid=? AND mstrecorddifferentiationtickettypeid=? AND mstrecorddifferentiationcatid=0 AND mstrecorddifferentiationpriorityid=? AND activeflg=1 AND deleteflg=0"
	stmt, error := tx.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in GetImpactUrgencydtls Prepare Statement..")
		return impactID, urgencyID, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(ClientID, OrgnID, ReocrdType, PriorityID)
	if err != nil {
		logger.Log.Println("Exception in GetImpactUrgencydtls Query Statement..")
		return impactID, urgencyID, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&impactID, &urgencyID); err != nil {
			logger.Log.Println("GetImpactUrgencydtls Error in fetching data")
		}

	}
	return impactID, urgencyID, nil
}

// func (mdao DbConn) Getdiffdtls(ClientID int64, Mstorgnhirarchyid int64) (map[int64]string, error) {
// 	var t = make(map[int64]string)
// 	var ID int64
// 	var Name string
// 	var sql = "SELECT id,name FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1"
// 	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid)
// 	defer rows.Close()
// 	if err != nil {
// 		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
// 		return t, err
// 	}
// 	for rows.Next() {
// 		err = rows.Scan(&ID, &Name)
// 		t[ID] = Name
// 	}
// 	//logger.Log.Println("Hashmap value is---------------->", t)
// 	return t, nil
// }

func (mdao DbConn) Getdiffdtls(ClientID int64, Mstorgnhirarchyid int64) (map[int64]string, error) {
	var t = make(map[int64]string)
	var ID int64
	var Name string
	var sql = "SELECT id,name FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1"
	stmt, error := mdao.DB.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in GetImpactUrgencydtls Prepare Statement..", error)
		return t, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(ClientID, Mstorgnhirarchyid)
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ID, &Name)
		t[ID] = Name
	}
	//logger.Log.Println("Hashmap value is---------------->", t)
	return t, nil
}

func (mdao DbConn) GetOriginalInfo(OrginalID int64) (entities.StagetableEntity, error) {
	value := entities.StagetableEntity{}
	var sql = "SELECT loginname,name,useremail,usermobileno,branch FROM mstclientuser WHERE id=?"
	rows, err := mdao.DB.Query(sql, OrginalID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOriginalInfo Get Statement Prepare Error", err)
		return value, err
	}
	for rows.Next() {
		err = rows.Scan(&value.Orgcreatorloginid, &value.Orgcreatorname, &value.Orgcreatoremail, &value.Orgcreatorphone, &value.Orgcreatorlocation)

	}
	return value, nil
}

func (mdao DbConn) GetCreatorInfo(CreatorID int64) (entities.StagetableEntity, error) {
	value := entities.StagetableEntity{}
	var sql = "SELECT loginname,vipuser FROM mstclientuser WHERE id=?"
	rows, err := mdao.DB.Query(sql, CreatorID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOriginalInfo Get Statement Prepare Error", err)
		return value, err
	}
	for rows.Next() {
		err = rows.Scan(&value.Requestorloginid, &value.Vipticket)

	}
	return value, nil
}

func InsertStageTbl(tx *sql.Tx, sgtable entities.StagetableEntity) error {
	var sql = "Insert into recordfulldetails (clientid,mstorgnhirarchyid,recordid,ticketid,source,requestorid,requestorname,requestorlocation,requestorphone,requestoremail,shortdescription,tickettypeid,tickettype,statusid,status,priorityid,priority,impact,impactid,urgency,urgencyid,createddatetime,requestorloginid,orgcreatorlocation,orgcreatorphone,orgcreatorid,orgcreatorname,orgcreatorloginid,orgcreatoremail,vipticket,lastuserid,lastuser,lastupdateddatetime,levelonecatename) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,?,?,?,?,?,?,?,?,?,now(),?)"
	stmt, err := tx.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(sgtable.ClientID, sgtable.OrgnID, sgtable.RecordID, sgtable.TicketID, sgtable.Source, sgtable.RequestorID, sgtable.Requestorname, sgtable.Requestorlocation, sgtable.Requestorphone, sgtable.Requestoremail, sgtable.Shortdescription, sgtable.Tickettypeid, sgtable.Tickettype, sgtable.Statusid, sgtable.Status, sgtable.Priorityid, sgtable.Priority, sgtable.Impact, sgtable.Impactid, sgtable.Urgency, sgtable.Urgencyid, sgtable.Requestorloginid, sgtable.Orgcreatorlocation, sgtable.Orgcreatorphone, sgtable.Orgcreatorid, sgtable.Orgcreatorname, sgtable.Orgcreatorloginid, sgtable.Orgcreatoremail, sgtable.Vipticket, sgtable.LastuserID, sgtable.Lastusername, sgtable.Fstlevelcategorynm)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func Updatecategorychangecout(tx *sql.Tx, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET categorychangecount=(categorychangecount+1) WHERE recordid=?"
	stmt, err := tx.Prepare(sql)
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

func UpdateStagePriority(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, PriorityID int64, Priority string) error {
	var sql = "UPDATE recordfulldetails SET prioritycount=(prioritycount+1),priorityid=?,priority=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(PriorityID, Priority, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateStageImpact(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, ImpactID int64, Impact string) error {
	var sql = "UPDATE recordfulldetails SET impactid=?,impact=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ImpactID, Impact, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateStageUrgency(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, UrgencyID int64, Urgency string) error {
	var sql = "UPDATE recordfulldetails SET urgencyid=?,urgency=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(UrgencyID, Urgency, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateStageStatus(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, StatusID int64, Status string) error {
	var sql = "UPDATE recordfulldetails SET statusid=?,status=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(StatusID, Status, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateFollowupcountFromCommon(ClientID int64, OrgnID int64, RecordID int64, Fcount int64) error {
	var sql = "UPDATE recordfulldetails SET followupcount=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Fcount, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateOutboundcount(ClientID int64, OrgnID int64, RecordID int64, Ocount int64) error {
	var sql = "UPDATE recordfulldetails SET outboundcount=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Ocount, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateAging(ClientID int64, OrgnID int64, RecordID int64, Aging int64) error {
	var sql = "UPDATE recordfulldetails SET calendaraging=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Aging, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateWorknotedaycount(ClientID int64, OrgnID int64, RecordID int64, Daycount int64) error {
	var sql = "UPDATE recordfulldetails SET worknotenotupdated=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Daycount, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateSLAFields(ClientID int64, OrgnID int64, RecordID int64, Responseduedate string, Responseclockstatus string, Resolutionduedate string, ResolutionClockstatus string) error {
	var sql = "UPDATE recordfulldetails SET respsladuedatetime=?,respclockstatus=?,resosladuedatetime=?,resoclockstatus=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>2222222222222222222222222222222222>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", ClientID, OrgnID, RecordID, Responseduedate, Responseclockstatus, Resolutionduedate, ResolutionClockstatus)
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>", err)
		return errors.New("SQL Prepare Error")
	}


	var aa4 string
	var bb4 string

	if len(Responseduedate) > 0 {
	var aa = Responseduedate
	var aa1 = strings.Split(aa, " ")
	var aa2 = aa1[0]
	var aa3 = strings.Split(aa2, "-")
	aa4 = aa3[2] + "-" + aa3[1] + "-" + aa3[0] + " " + aa1[1]
	}

	if len(Resolutionduedate) > 0 {
	var bb = Resolutionduedate
	var bb1 = strings.Split(bb, " ")
	var bb2 = bb1[0]
	var bb3 = strings.Split(bb2, "-")
	bb4 = bb3[2] + "-" + bb3[1] + "-" + bb3[0] + " " + bb1[1]
	}
	defer stmt.Close()
	_, err = stmt.Exec(aa4, Responseclockstatus, bb4, ResolutionClockstatus, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateSLAFields(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, Responseduedate string, Responseclockstatus string, Resolutionduedate string, ResolutionClockstatus string) error {
	var sql = "UPDATE recordfulldetails SET respsladuedatetime=?,respclockstatus=?,resosladuedatetime=?,resoclockstatus=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	//logger.Log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>2222222222222222222222222222222222>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", ClientID, OrgnID, RecordID, Responseduedate, Responseclockstatus, Resolutionduedate, ResolutionClockstatus)
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>", err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	var aa4 string
	var bb4 string
	if Responseduedate != "" {
		var aa = Responseduedate
		var aa1 = strings.Split(aa, " ")
		var aa2 = aa1[0]
		var aa3 = strings.Split(aa2, "-")
		aa4 = aa3[2] + "-" + aa3[1] + "-" + aa3[0] + " " + aa1[1]
	}
	if Resolutionduedate != "" {
		var bb = Resolutionduedate
		var bb1 = strings.Split(bb, " ")
		var bb2 = bb1[0]
		var bb3 = strings.Split(bb2, "-")
		bb4 = bb3[2] + "-" + bb3[1] + "-" + bb3[0] + " " + bb1[1]
	}
	_, err = stmt.Exec(aa4, Responseclockstatus, bb4, ResolutionClockstatus, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println("555555555555555555555555555555555555555555555555-------------------------------------------->", err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateFirstResponse(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, responsetiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET firstresponsedatetime=now(),latestresponsedatetime=now(),responsetime=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsetiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateLatestResponse(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, responsetiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET latestresponsedatetime=now(),responsetime=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsetiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateFirstResolution(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, resolutiontiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET firstresodatetime=now(),latestresodatetime=now(),resolutiontime=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutiontiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateLatestResolution(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, resolutiontiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET latestresodatetime=now(),resolutiontime=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutiontiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateClosureCode(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, Closurecode string) error {
	var sql = "UPDATE recordfulldetails SET csatscore=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Closurecode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateClosureComment(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, Closurecomment string) error {
	var sql = "UPDATE recordfulldetails SET csatcomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Closurecomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateReopenCount(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET reopencount=(reopencount+1),reopendatetime=now(),reopenedflag='Y'  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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

func (mdao DbConn) IsClientSpecificOrNot(clientid int64, RecorddifftypeID int64, RecorddiffID int64) (int64, error) {
	var flag int64
	var sql = "SELECT isclient FROM mstrecordconfig WHERE clientid=? AND recorddifftypeid=? AND recorddiffid=? AND activeflg=1 AND deleteflg=0"
	stmt, error := mdao.DB.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in IsClientSpecificOrNot Prepare Statement..")
		return flag, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, RecorddifftypeID, RecorddiffID)
	if err != nil {
		logger.Log.Println(err)
		return flag, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&flag); err != nil {
			logger.Log.Println(err)
		}

	}
	return flag, nil
}

func (mdao DbConn) GetNumberbyclientID(clientid int64, difftypeID int64, diffID int64) (int, error) {
	var number int
	var sql = "SELECT number FROM mstrecordautoincreament WHERE clientid=? AND recorddifftypeid=? AND recorddiffid=? AND activeflg=1 AND deleteflg=0"
	stmt, error := mdao.DB.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return number, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, difftypeID, diffID)
	if err != nil {
		logger.Log.Println(err)
		return number, error
	}
	defer rows.Close()
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&number); err != nil {
			logger.Log.Println(err)
		}

	}
	logger.Log.Println("number value is :", number)
	return number, nil
}

func (mdao DbConn) GetPrefixnormalbyclientID(clientid int64, difftypeID int64, diffID int64) ([]string, error) {
	var prefixarr []string
	var sql = "SELECT prefix,year,month,day,configurezero FROM mstrecordconfig WHERE clientid=? AND recorddifftypeid=? AND recorddiffid=? AND deleteflg=0 AND activeflg=1 "
	stmt, error := mdao.DB.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in GetPrefix Prepare Statement..")
		return prefixarr, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, difftypeID, diffID)
	if err != nil {
		logger.Log.Println("Exception in GetPrefix Query Statement..")
		return prefixarr, error
	}
	defer rows.Close()
	if rows.Next() {
		var prefix string
		var year string
		var month string
		var day string
		var configurezero string

		if err := rows.Scan(&prefix, &year, &month, &day, &configurezero); err != nil {
			logger.Log.Println("Error in fetching data")
		}
		prefixarr = append(prefixarr, prefix)
		prefixarr = append(prefixarr, year)
		prefixarr = append(prefixarr, month)
		prefixarr = append(prefixarr, day)
		prefixarr = append(prefixarr, configurezero)
	}
	logger.Log.Println("prefixarr value is :", prefixarr)
	return prefixarr, nil
}

func (mdao DbConn) GetParentDiffID(clientid int64, OrgnID int64, difftypeID int64, diffID int64) (int64, error) {
	var number int64
	var sql = "SELECT fromrecorddiffid FROM mstrecorddifferentiationmap WHERE toclientid=? AND tomstorgnhirarchyid=? AND torecorddifftypeid=? AND torecorddiffid=? AND activeflg=1 AND deleteflg=0"
	stmt, error := mdao.DB.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..")
		return number, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(clientid, OrgnID, difftypeID, diffID)
	if err != nil {
		logger.Log.Println(err)
		return number, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&number); err != nil {
			logger.Log.Println(err)
		}

	}
	logger.Log.Println("number value is :", number)
	return number, nil
}

func (mdao DbConn) UpdatenumberbyclientID(clientid int64, difftypeID int64, diffid int64) error {
	var sql = "UPDATE mstrecordautoincreament SET number = (number+1) WHERE clientid=? AND recorddifftypeid=? AND recorddiffid=? AND activeflg=1 AND deleteflg=0"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(clientid, difftypeID, diffid)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetApproveFlag(Recordid int64, ClientID int64, Mstorgnhirarchyid int64) (int64, error) {
	logger.Log.Println("In side GetApproveFlag")
	var ID int64
	var sql = "SELECT isapprove FROM trnrecord WHERE clientid=? AND mstorgnhirarchyid=? AND id=?"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetApproveFlag Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("GetApproveFlag rows.next() Error", err)
	}
	return ID, nil
}

func (mdao DbConn) GetPreviousStateID(ClientID int64, Mstorgnhirarchyid int64, TypeSeq int64) (int64, error) {
	logger.Log.Println("In side GetApproveFlag")
	var ID int64
	var sql = "SELECT mststateid FROM maprecordstatetodifferentiation WHERE clientid=? and mstorgnhirarchyid=? and recorddifftypeid=3 and recorddiffid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? and mstorgnhirarchyid=? and recorddifftypeid=3 and seqno=?)"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, ClientID, Mstorgnhirarchyid, TypeSeq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetApproveFlag Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("GetApproveFlag rows.next() Error", err)
	}
	return ID, nil
}

func UpdateStageResolver(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, UserID int64, UserName string, GrpID int64, Grpname string) error {
	var sql = "UPDATE recordfulldetails SET resogroupid=?,resogroup=?,resolveduserid=?,resolveduser=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(GrpID, Grpname, UserID, UserName, ClientID, Mstorgnhirarchyid, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateUserreplydatetime(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, replytimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET userreplieddatetime=now(),userreplytimetaken=?,assigneduserid=0 ,assigneduser=NULL  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(replytimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdatePendinguserAction(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET pendinguserdatetime=now(),pendingusercount=(pendingusercount+1)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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

func UpdateFollowupcount(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET pendingvendorcount=(pendingvendorcount+1)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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

func UpdateUserInfo(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, UserID int64, Username string) error {
	var sql = "UPDATE recordfulldetails SET lastupdateddatetime=now(),lastuserid=?,lastuser=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserID, Username, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateUserInfoWithoutTrn(ClientID int64, OrgnID int64, RecordID int64, UserID int64, Username string) error {
	var sql = "UPDATE recordfulldetails SET lastupdateddatetime=now(),lastuserid=?,lastuser=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserID, Username, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetUsername(UserID int64) (string, error) {
	logger.Log.Println("In side GetUsername")
	var name string
	var sql = "SELECT name FROM mstclientuser WHERE id=?"
	rows, err := mdao.DB.Query(sql, UserID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetUsername Get Statement Prepare Error", err)
		return name, err
	}
	for rows.Next() {
		err = rows.Scan(&name)
		logger.Log.Println("GetUsername rows.next() Error", err)
	}
	return name, nil
}

func UpdateFollowuptimetaken(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, followuptimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET followuprespdatetime=now(),followuptimetaken=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(followuptimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateCloseddate(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET closedatetime=now()  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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

func (mdao DbConn) GetFstlevelCatName(catagoryID int64) (string, error) {
	var fstcategoryname string
	var sql = "SELECT name FROM mstrecorddifferentiation WHERE id=?"
	rows, err := mdao.DB.Query(sql, catagoryID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOriginalInfo Get Statement Prepare Error", err)
		return fstcategoryname, err
	}
	for rows.Next() {
		err = rows.Scan(&fstcategoryname)

	}
	return fstcategoryname, nil
}

func (mdao DbConn) Checkrecordnumberduplicacy(recordnumber string) (int64, error) {
	var count int64
	var sql = "SELECT count(id) count FROM trnrecord WHERE code=? "
	stmt, error := mdao.DB.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in Checkrecordnumberduplicacy....")
		return count, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(recordnumber)
	if err != nil {
		logger.Log.Println(err)
		return count, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			logger.Log.Println(err)
		}

	}
	return count, nil
}

func (mdao DbConn) UpdateRecordID(lastInsertedID int64, recordno string) error {
	logger.Log.Println("updatetrnrecord query -->", updatetrnrecord)
	logger.Log.Println("parameters -->", recordno, lastInsertedID)

	stmt, err := mdao.DB.Prepare(updatetrnrecord)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(recordno, lastInsertedID)
	if err != nil {
		logger.Log.Println(err)
		return err
	}

	return nil
}

func (mdao DbConn) UpdateRecordIDINStage(lastInsertedID int64, recordno string) error {
	//logger.Log.Println("UpdateRecordIDINStage query -->", updatetrnrecord)
	//logger.Log.Println("UpdateRecordIDINStage -->", recordno, lastInsertedID)
	var sql = "UPDATE recordfulldetails SET ticketid=? WHERE recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(recordno, lastInsertedID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateTrnrecorddeleteflg(ID int64) error {
	var sql = "UPDATE trnrecord SET deleteflg=1 WHERE id=?"
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

func (mdao DbConn) UpdateALLSLAFields(ClientID int64, OrgnID int64, RecordID int64, Responseduedate string, Responseclockstatus string, Resolutionduedate string, ResolutionClockstatus string) error {
	var sql = "UPDATE recordfulldetails SET respsladuedatetime=?,respclockstatus=?,resosladuedatetime=?,resoclockstatus=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(">>>>>>>>>>>>>>>>", err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	var aa4 string
	var bb4 string
	if Responseduedate != "" {
		var aa = Responseduedate
		var aa1 = strings.Split(aa, " ")
		var aa2 = aa1[0]
		var aa3 = strings.Split(aa2, "-")
		aa4 = aa3[2] + "-" + aa3[1] + "-" + aa3[0] + " " + aa1[1]
	}
	if Resolutionduedate != "" {
		var bb = Resolutionduedate
		var bb1 = strings.Split(bb, " ")
		var bb2 = bb1[0]
		var bb3 = strings.Split(bb2, "-")
		bb4 = bb3[2] + "-" + bb3[1] + "-" + bb3[0] + " " + bb1[1]
	}
	_, err = stmt.Exec(aa4, Responseclockstatus, bb4, ResolutionClockstatus, ClientID, OrgnID, RecordID)
	if err != nil {
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetRecordcreationdate(RecordID int64, Timediff int64) (string, error) {
	var createdate string
	var sql = "SELECT createdatetime FROM trnrecord WHERE id=?"
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordcreationdate Get Statement Prepare Error", err)
		return createdate, err
	}
	for rows.Next() {
		var datedata int64
		err = rows.Scan(&datedata)
		logger.Log.Println(err)
		createdate = Convertdate1(datedata, Timediff)

	}
	return createdate, nil
}

func (mdao DbConn) UpdateStageCreatedate(RecordID int64, Createdate string) error {
	var sql = "UPDATE recordfulldetails SET createddatetime=? WHERE recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Createdate, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

// func (mdao DbConn) GetParentRecordNo(difftypeID int64, diffID int64) (int64, string, error) {
// 	var Recordno string
// 	var ID int64
// 	var sql = "SELECT id,code FROM mstrecordcode WHERE recorddifftypeid=? AND recorddiffid=? AND isuse=0 AND activeflg=1 AND deleteflg=0 limit 1 FOR UPDATE" //FOR UPDATE
// 	stmt, error := mdao.DB.Prepare(sql)
// 	if error != nil {
// 		logger.Log.Println("Exception in GetParentRecordNo Prepare Statement..")
// 		return ID, Recordno, error
// 	}
// 	defer stmt.Close()
// 	rows, err := stmt.Query(difftypeID, diffID)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return ID, Recordno, error
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		if err := rows.Scan(&ID, &Recordno); err != nil {
// 			logger.Log.Println(err)
// 		}

// 	}
// 	logger.Log.Println("Recordno value is :", Recordno)
// 	return ID, Recordno, nil
// }

// func (mdao DbConn) UpdateRecordnoTB(ID int64) error {
// 	var sql = "UPDATE mstrecordcode SET isuse=1 WHERE id=?"
// 	stmt, err := mdao.DB.Prepare(sql)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return errors.New("SQL Prepare Error")
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(ID)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return errors.New("SQL Execution Error")
// 	}

// 	return nil
// }

func (mdao DbConn) GetParentRecordNo(difftypeID int64, diffID int64, tx *sql.Tx) (int64, string, error) {
	var Recordno string
	var ID int64
	var sql = "SELECT id,code FROM mstrecordcode WHERE recorddifftypeid=? AND recorddiffid=? AND isuse=0 AND activeflg=1 AND deleteflg=0 limit 1 FOR UPDATE" //FOR UPDATE
	stmt, error := tx.Prepare(sql)
	if error != nil {
		logger.Log.Println("Exception in GetParentRecordNo Prepare Statement..")
		return ID, Recordno, error
	}
	defer stmt.Close()
	rows, err := stmt.Query(difftypeID, diffID)
	if err != nil {
		logger.Log.Println(err)
		return ID, Recordno, error
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ID, &Recordno); err != nil {
			logger.Log.Println(err)
			//mdao.UpdateRecordnoTB(ID)
		}

	}
	logger.Log.Println("Recordno value is :", Recordno)
	return ID, Recordno, nil
}

func (mdao DbConn) UpdateRecordnoTB(ID int64, tx *sql.Tx) error {
	var sql = "UPDATE mstrecordcode SET isuse=1 WHERE id=?"
	stmt, err := tx.Prepare(sql)
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

func UpdatePendingvendorname(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, vendorname string) error {
	var sql = "UPDATE recordfulldetails SET vendorname=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(vendorname, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdatePendingvendorticketID(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, vendorticketID string) error {
	var sql = "UPDATE recordfulldetails SET vendorticketid=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(vendorticketID, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateResolutioncode(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, resolutioncode string) error {
	var sql = "UPDATE recordfulldetails SET resolutioncode=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutioncode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateResolutioncomment(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, resolutioncomment string) error {
	var sql = "UPDATE recordfulldetails SET resolutioncomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutioncomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateRecordfulldetailsddeleteflg(ID int64) error {
	var sql = "UPDATE recordfulldetails SET deleteflg=1 WHERE recordid=?"
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

// Changes On 18.05.2022 Start From Here ..........................

func (mdao DbConn) UpdateStageStatus(ClientID int64, OrgnID int64, RecordID int64, StatusID int64, Status string) error {
	var sql = "UPDATE recordfulldetails SET statusid=?,status=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(StatusID, Status, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateFirstResponse(ClientID int64, OrgnID int64, RecordID int64, responsetiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET firstresponsedatetime=now(),latestresponsedatetime=now(),responsetime=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsetiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateLatestResponse(ClientID int64, OrgnID int64, RecordID int64, responsetiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET latestresponsedatetime=now(),responsetime=(responsetime+?)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsetiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

/*func (mdao DbConn) UpdateFirstResolution(ClientID int64, OrgnID int64, RecordID int64, resolutiontiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET firstresodatetime=now(),latestresodatetime=now(),resolutiontime=(resolutiontime+?)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutiontiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}*/

/*func (mdao DbConn) UpdateLatestResolution(ClientID int64, OrgnID int64, RecordID int64, resolutiontiimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET latestresodatetime=now(),resolutiontime=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutiontiimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}*/

func (mdao DbConn) UpdateFirstResolution(ClientID int64, OrgnID int64, RecordID int64, resolutiontiimetaken int64, resolutionTimeTakenWithoutIdleTime int64) error {
		var sql = "UPDATE recordfulldetails SET firstresodatetime=now(),latestresodatetime=now(),resolutiontime=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
			logger.Log.Println("UpdateFirstResolution===================================================================>", resolutiontiimetaken, resolutionTimeTakenWithoutIdleTime, ClientID, OrgnID, RecordID)

				stmt, err := mdao.DB.Prepare(sql)
					if err != nil {
								logger.Log.Println(err)
										return errors.New("SQL Prepare Error")
											}
												defer stmt.Close()
													_, err = stmt.Exec(resolutiontiimetaken, ClientID, OrgnID, RecordID)
														if err != nil {
																	logger.Log.Println(err)
																			return errors.New("SQL Execution Error")
																				}

																					return nil
																				}
func (mdao DbConn) GetIdleTime(ClientID int64, OrgnID int64, RecordID int64) (int64, error) {
	logger.Log.Println("GetIdleTime ========================>")
	var idleTime int64
	var sql = "select sum(pushtime) from mstsladue where clientid=? and  mstorgnhirarchyid=? and  therecordid=? and activeflg=1 and deleteflg=0"
	rows, err := mdao.DB.Query(sql, ClientID, OrgnID, RecordID)

	if err != nil {
		logger.Log.Println("GetIdleTime Get Statement Prepare Error==============>", err)
		return idleTime, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&idleTime)
		if err != nil {
			logger.Log.Println("GetIdleTime rows.next() Error===========>", err)
			return idleTime, err
		}

	}
		logger.Log.Println("GetIdleTime===================================================================>", idleTime, ClientID, OrgnID, RecordID)

	return idleTime, nil
}

func (mdao DbConn) UpdateLatestResolution(ClientID int64, OrgnID int64, RecordID int64, resolutiontiimetaken int64, resolutionTimeTakenWithoutIdleTime int64) error {
	
	var sql = "UPDATE recordfulldetails SET latestresodatetime=now(),resolutiontime=(resolutiontime+?) WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
		//logger.Log.Println("UpdateLatestResolution===================================================================>", resolutiontiimetaken, resolutionTimeTakenWithoutIdleTime, ClientID, OrgnID, RecordID)

	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutiontiimetaken,  ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}
func (mdao DbConn) UpdateUserreplydatetime(ClientID int64, OrgnID int64, RecordID int64, replytimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET userreplieddatetime=now(),userreplytimetaken=?,assigneduserid=0 ,assigneduser=NULL  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(replytimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateFollowuptimetaken(ClientID int64, OrgnID int64, RecordID int64, followuptimetaken int64) error {
	var sql = "UPDATE recordfulldetails SET followuprespdatetime=now(),followuptimetaken=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(followuptimetaken, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateCloseddate(ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET closedatetime=now()  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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

func (mdao DbConn) UpdateStageResolver(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, UserID int64, UserName string, GrpID int64, Grpname string) error {
	var sql = "UPDATE recordfulldetails SET resogroupid=?,resogroup=?,resolveduserid=?,resolveduser=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(GrpID, Grpname, UserID, UserName, ClientID, Mstorgnhirarchyid, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateReopenCount(ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET reopencount=(reopencount+1),reopendatetime=now(),reopenedflag='Y'  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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


func (mdao DbConn) UpdateReopenCount1(ClientID int64, OrgnID int64, RecordID int64, reopentimetaken int64,reopentimetakencalendar int64) error {
	var sql = "UPDATE recordfulldetails SET reopencount=(reopencount+1),reopendatetime=now(),reopenedflag='Y',reopentimetaken=(reopentimetaken+?),reopentimetakencalendar=(reopentimetakencalendar+?)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(reopentimetaken,reopentimetakencalendar, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdatePendinguserAction(ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET pendinguserdatetime=now(),pendingusercount=(pendingusercount+1)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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

func (mdao DbConn) UpdateFollowupcount(ClientID int64, OrgnID int64, RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET pendingvendorcount=(pendingvendorcount+1)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
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

func (mdao DbConn) UpdateUserInfo(ClientID int64, OrgnID int64, RecordID int64, UserID int64, Username string) error {
	var sql = "UPDATE recordfulldetails SET lastupdateddatetime=now(),lastuserid=?,lastuser=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(UserID, Username, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateClosureCode(ClientID int64, OrgnID int64, RecordID int64, Closurecode string) error {
	var sql = "UPDATE recordfulldetails SET csatscore=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Closurecode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateClosureComment(ClientID int64, OrgnID int64, RecordID int64, Closurecomment string) error {
	var sql = "UPDATE recordfulldetails SET csatcomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(Closurecomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdatePendingvendorname(ClientID int64, OrgnID int64, RecordID int64, vendorname string) error {
	var sql = "UPDATE recordfulldetails SET vendorname=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(vendorname, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdatePendingvendorticketID(ClientID int64, OrgnID int64, RecordID int64, vendorticketID string) error {
	var sql = "UPDATE recordfulldetails SET vendorticketid=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(vendorticketID, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateResolutioncode(ClientID int64, OrgnID int64, RecordID int64, resolutioncode string) error {
	var sql = "UPDATE recordfulldetails SET resolutioncode=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutioncode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateResolutioncomment(ClientID int64, OrgnID int64, RecordID int64, resolutioncomment string) error {
	var sql = "UPDATE recordfulldetails SET resolutioncomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutioncomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateStagePriority(ClientID int64, OrgnID int64, RecordID int64, PriorityID int64, Priority string) error {
	var sql = "UPDATE recordfulldetails SET prioritycount=(prioritycount+1),priorityid=?,priority=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(PriorityID, Priority, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateStageImpact(ClientID int64, OrgnID int64, RecordID int64, ImpactID int64, Impact string) error {
	var sql = "UPDATE recordfulldetails SET impactid=?,impact=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ImpactID, Impact, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateStageUrgency(ClientID int64, OrgnID int64, RecordID int64, UrgencyID int64, Urgency string) error {
	var sql = "UPDATE recordfulldetails SET urgencyid=?,urgency=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(UrgencyID, Urgency, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) Updatecategorychangecout(RecordID int64) error {
	var sql = "UPDATE recordfulldetails SET categorychangecount=(categorychangecount+1) WHERE recordid=?"
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



func (mdao DbConn) GetWorkingCategoryDetails(RecordID int64) (int64, int64, error) {
	logger.Log.Println("In side GetWorkingCategoryDetails")
	var workingdifftypeID int64
	var workingdiffID int64
	var sql = "SELECT recorddifftypeid,recorddiffid FROM maprecordtorecorddifferentiation where recordid=? AND isworking=1 AND islatest=1 AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetWorkingCategoryDetails Get Statement Prepare Error", err)
		return workingdifftypeID, workingdiffID, err
	}
	for rows.Next() {
		err = rows.Scan(&workingdifftypeID, &workingdiffID)
		logger.Log.Println("GetWorkingCategoryDetails rows.next() Error", err)
	}
	return workingdifftypeID, workingdiffID, nil
}

func (mdao DbConn) GetEbondingSeq(ClientID int64, Mstorgnhirarchyid int64, WorkingdifftypeID int64, WorkingdiffID int64) (int64, error) {
	logger.Log.Println("In side GetEbondingSeq")
	var ebondingseq int64
	var sql = "SELECT seqno FROM ebondingmst WHERE id in (SELECT distinct ebondingid FROM ebondingdifferentiationmap WHERE clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND activeflg=1 AND deleteflg=0) AND activeflg=1 AND deleteflg=0"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, WorkingdifftypeID, WorkingdiffID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetEbondingSeq Get Statement Prepare Error", err)
		return ebondingseq, err
	}
	for rows.Next() {
		err = rows.Scan(&ebondingseq)
		logger.Log.Println("GetEbondingSeq rows.next() Error", err)
	}
	return ebondingseq, nil
}


func (mdao DbConn) GetRecordCodeANDTypeID(RecordID int64) (string, int64, error) {
	logger.Log.Println("In side GetRecordCodeANDTypeID")
	var typeID int64
	var code string
	var sql = "SELECT a.code,b.recorddiffid FROM trnrecord a,maprecordtorecorddifferentiation b WHERE a.id=? AND a.id=b.recordid AND b.recorddifftypeid=2 AND b.islatest=1"
	rows, err := mdao.DB.Query(sql, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordCodeANDTypeID Get Statement Prepare Error", err)
		return code, typeID, err
	}
	for rows.Next() {
		err = rows.Scan(&code, &typeID)
		logger.Log.Println("GetRecordCodeANDTypeID rows.next() Error", err)
	}
	return code, typeID, nil
}



func (mdao DbConn) Updateresotimeexcludeidletime(ClientID int64, OrgnID int64, RecordID int64, idleTime int64) error {
	var sql = "UPDATE recordfulldetails SET resotimeexcludeidletime=(resolutiontime-?)  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	logger.Log.Println("UpdateFirstResolution===================================================================>", ClientID, OrgnID, RecordID, idleTime)

	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(idleTime, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}



func (mdao DbConn) UpdateStageResolverForSR(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, UserID int64, UserName string, GrpID int64, Grpname string) error {
	var sql = "UPDATE recordfulldetails SET resogroupid=?,resogroup=?,resolveduserid=?,resolveduser=?,firstresodatetime=now(),latestresodatetime=now()  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(GrpID, Grpname, UserID, UserName, ClientID, Mstorgnhirarchyid, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

