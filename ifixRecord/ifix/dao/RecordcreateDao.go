package dao

import (
	"database/sql"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

var categorysql = "SELECT mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno typeseq,mstrecorddifferentiation.id,mstrecorddifferentiation.name,mstrecorddifferentiation.seqno FROM mstrecordtype,mstrecorddifferentiationtype,mstrecorddifferentiation WHERE mstrecordtype.clientid=? AND mstrecordtype.mstorgnhirarchyid=? AND mstrecordtype.activeflg=1 AND mstrecordtype.deleteflg=0 AND mstrecordtype.fromrecorddifftypeid=? AND mstrecordtype.fromrecorddiffid=? AND mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiationtype.id=mstrecordtype.torecorddifftypeid AND mstrecorddifferentiationtype.parentid=1 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiation.id=mstrecordtype.torecorddiffid ORDER BY mstrecorddifferentiationtype.seqno ASC,mstrecorddifferentiation.name ASC"

var tickettypesql = "SELECT  mstrecorddifferentiation.id, mstrecorddifferentiation.name, mstrecorddifferentiation.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstrecorddifferentiation, mstrecorddifferentiationtype WHERE mstrecorddifferentiation.activeflg = 1 AND mstrecorddifferentiation.deleteflg = 0 AND mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id AND mstrecorddifferentiation.clientid = ? AND mstrecorddifferentiation.mstorgnhirarchyid = ? AND mstrecorddifferentiation.recorddifftypeid = ? "

var getCatchildsql = "SELECT mstrecorddifferentiation.id,mstrecorddifferentiation.name,mstrecorddifferentiation.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstrecordtype,mstrecorddifferentiationtype,mstrecorddifferentiation WHERE mstrecordtype.clientid=? AND mstrecordtype.mstorgnhirarchyid=? AND mstrecordtype.activeflg=1 AND mstrecordtype.deleteflg=0 AND mstrecordtype.fromrecorddifftypeid=? AND mstrecordtype.fromrecorddiffid=? AND mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiationtype.id=mstrecordtype.torecorddifftypeid AND mstrecorddifferentiationtype.parentid=1 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiation.id=mstrecordtype.torecorddiffid AND mstrecorddifferentiation.parentid=?  ORDER BY mstrecorddifferentiationtype.seqno ASC,mstrecorddifferentiation.name ASC"

//var getdirectionsql = "SELECT mstbusinessdirection.direction FROM mstbusinessdirection WHERE mstbusinessdirection.activeflg=1 AND mstbusinessdirection.deleteflg=0 AND mstbusinessdirection.clientid=? AND mstbusinessdirection.mstorgnhirarchyid=? AND mstbusinessdirection.mstrecorddifferentiationtypeid=? AND mstbusinessdirection.mstrecorddifferentiationid=? LIMIT 1"

var getdirectionsql = "SELECT mstbusinessdirection.direction FROM mstbusinessdirection WHERE mstbusinessdirection.activeflg=1 AND mstbusinessdirection.deleteflg=0 AND mstbusinessdirection.clientid=? AND mstbusinessdirection.mstorgnhirarchyid=? AND mstbusinessdirection.mstrecorddifferentiationtypeid=? AND mstbusinessdirection.mstrecorddifferentiationid=? and mstbusinessdirection.baseconfig=?"

var getimpacturgencysql = "SELECT  mstrecorddifferentiation.id, mstrecorddifferentiation.name, mstrecorddifferentiation.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstrecordtype, mstrecorddifferentiationtype, mstrecorddifferentiation WHERE mstrecordtype.clientid = ? AND mstrecordtype.mstorgnhirarchyid = ? AND mstrecordtype.activeflg = 1 AND mstrecordtype.deleteflg = 0 AND mstrecordtype.fromrecorddifftypeid = ? AND mstrecordtype.fromrecorddiffid = ? AND mstrecorddifferentiationtype.activeflg = 1 AND mstrecorddifferentiationtype.deleteflg = 0 AND mstrecorddifferentiationtype.id = mstrecordtype.torecorddifftypeid AND mstrecorddifferentiationtype.seqno = ? AND mstrecorddifferentiation.activeflg = 1 AND mstrecorddifferentiation.deleteflg = 0 AND mstrecorddifferentiation.id = mstrecordtype.torecorddiffid ORDER BY mstrecorddifferentiationtype.seqno ASC , mstrecorddifferentiation.name ASC"

var getprioritysql = "SELECT  prioritydiff.id, prioritydiff.name, prioritydiff.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstbusinessmatrix, mstrecorddifferentiation prioritydiff, mstrecorddifferentiationtype WHERE mstbusinessmatrix.activeflg = 1 AND prioritydiff.activeflg = 1 AND mstrecorddifferentiationtype.activeflg=1 AND mstbusinessmatrix.deleteflg = 0 AND prioritydiff.deleteflg = 0 AND mstrecorddifferentiationtype.deleteflg=0 AND mstbusinessmatrix.mstrecorddifferentiationpriorityid = prioritydiff.id AND prioritydiff.recorddifftypeid=mstrecorddifferentiationtype.id AND mstbusinessmatrix.clientid = ? AND mstbusinessmatrix.mstorgnhirarchyid = ? AND mstbusinessmatrix.mstrecorddifferentiationtickettypeid = ? AND mstbusinessmatrix.mstrecorddifferentiationimpactid = ? AND mstbusinessmatrix.mstrecorddifferentiationurgencyid = ? AND mstbusinessmatrix.mstrecorddifferentiationcatid = ? "

var gettermlistsql = "SELECT  mstrecordterms.id, mstrecordterms.termname, mstrecordterms.termtypeid, msttermtype.termtypename, mstrecordterms.termvalue FROM mststateterm JOIN mstrecordterms ON mststateterm.recordtermid = mstrecordterms.id AND mststateterm.clientid = mstrecordterms.clientid AND mststateterm.mstorgnhirarchyid = mstrecordterms.mstorgnhirarchyid JOIN msttermtype ON mstrecordterms.termtypeid = msttermtype.id WHERE mststateterm.activeflg = 1 AND mststateterm.deleteflg = 0 AND mstrecordterms.activeflg = 1 AND mstrecordterms.deleteflg = 0 AND msttermtype.id = 3 AND mstrecordterms.clientid = ? AND mstrecordterms.mstorgnhirarchyid = ? AND mststateterm.recorddifftypeid = ? AND mststateterm.recorddiffid = ? ORDER BY mstrecordterms.termname ASC"

var getfirststatussql = "SELECT mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename,mstrecorddifferentiationtype.seqno typeseq,mstrecorddifferentiation.id,mstrecorddifferentiation.name,mstrecorddifferentiation.seqno FROM mstrecordtype JOIN mstrecorddifferentiation ON mstrecordtype.torecorddiffid=mstrecorddifferentiation.id AND mstrecordtype.torecorddifftypeid=mstrecorddifferentiation.recorddifftypeid AND mstrecordtype.clientid=mstrecorddifferentiation.clientid AND mstrecordtype.mstorgnhirarchyid=mstrecorddifferentiation.mstorgnhirarchyid JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id WHERE mstrecordtype.activeflg=1 AND mstrecordtype.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.seqno=2 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecordtype.clientid=? AND mstrecordtype.mstorgnhirarchyid=? AND mstrecordtype.fromrecorddifftypeid=? AND mstrecordtype.fromrecorddiffid=? ORDER BY mstrecorddifferentiation.seqno ASC LIMIT 1"

var isassetsql = "SELECT IF(COUNT(mstrecordtype.id)>0,1,0) counter FROM mstrecordtype JOIN mstrecorddifferentiationtype ON mstrecordtype.torecorddifftypeid=mstrecorddifferentiationtype.id WHERE mstrecordtype.activeflg=1 AND mstrecordtype.deleteflg=0 AND mstrecordtype.torecorddiffid=0 AND mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiationtype.seqno=5 AND mstrecordtype.clientid=? AND mstrecordtype.mstorgnhirarchyid=? AND mstrecordtype.fromrecorddifftypeid=? AND mstrecordtype.fromrecorddiffid=?"

var getworkinglabelsql = "SELECT mstworkdifferentiation.mainrecorddifftypeid  FROM mstworkdifferentiation WHERE mstworkdifferentiation.activeflg=1 AND mstworkdifferentiation.deleteflg=0 AND mstworkdifferentiation.clientid=? AND mstworkdifferentiation.mstorgnhirarchyid=? AND mstworkdifferentiation.forrecorddifftypeid=? AND mstworkdifferentiation.forrecorddiffid=?"

func (dbc DbConn) GetAllRecordtypes(req *entities.RecordcreaterequestEntity) ([]entities.RecordtypedetailsEntity, error) {
	logger.Log.Println("In side GetAllRecordTypes")
	values := []entities.RecordtypedetailsEntity{}
	rows, err := dbc.DB.Query(tickettypesql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid)

	if err != nil {
		logger.Log.Println("GetAllRecordTypes Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RecordtypedetailsEntity{}
		rows.Scan(&value.ID, &value.Name, &value.Seqno, &value.Typeid, &value.Typename, &value.Typeseq)
		values = append(values, value)
	}
	return values, nil
}

func (dbc DbConn) GetAllRecordtermslist(req *entities.RecordcreaterequestEntity) ([]entities.RecordtermlistEntity, error) {
	logger.Log.Println("In side GetAllRecordtermslist")
	values := []entities.RecordtermlistEntity{}
	logger.Log.Println("Query String: ", gettermlistsql)
	logger.Log.Println("Params : ", req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid)
	rows, err := dbc.DB.Query(gettermlistsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAllRecordtermslist Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordtermlistEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Termtypeid, &value.Termtypename, &value.Termvalue)
		values = append(values, value)
	}
	return values, nil
}

func (dbc DbConn) GetAllRecordcategories(req *entities.RecordcreaterequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetAllRecordcategories")
	values := []entities.RecordcatchildEntity{}
	//	logger.Log.Println("query is ---1111111111111111111111-->", categorysql)
	//	logger.Log.Println("query is ---1111111111111111111111-->", req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid)
	rows, err := dbc.DB.Query(categorysql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAllRecordcategories Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordcatchildEntity{}
		err := rows.Scan(&value.Typeid, &value.Typename, &value.Typeseq, &value.ID, &value.Name, &value.Seqno)
		logger.Log.Println("err is ---1111111111111111111111-->", err)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

// func (dbc DbConn) Getdirection(req *entities.RecordcreaterequestEntity) (int64, error) {
// 	logger.Log.Println("In side Getdirection")
// 	var values int64
// 	//rows, err := dbc.DB.Query(getdirectionsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid,)
// 	rows, err := dbc.DB.Query(getdirectionsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid, req.BaseConfig)
// 	defer rows.Close()
// 	if err != nil {
// 		logger.Log.Println("Getdirection Get Statement Prepare Error", err)
// 		return values, err
// 	}
// 	for rows.Next() {
// 		rows.Scan(&values)
// 	}
// 	//	logger.Log.Println("values is ----->", values)
// 	return values, nil
// }

func (dbc DbConn) Getdirection(req *entities.RecordcreaterequestEntity) (int64, error) {

	logger.Log.Println("In side Getdirection")

	req.BaseConfig = 1

	logger.Log.Println("req ----->", req)

	var values int64

	//rows, err := dbc.DB.Query(getdirectionsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid,)

	rows, err := dbc.DB.Query(getdirectionsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid, req.BaseConfig)

	defer rows.Close()

	if err != nil {

		logger.Log.Println("Getdirection Get Statement Prepare Error", err)

		return values, err

	}

	for rows.Next() {

		rows.Scan(&values)

	}

	logger.Log.Println("getdirectionsql ----->", getdirectionsql)
	return values, nil
}
func (dbc DbConn) GetRecordimpact(req *entities.RecordcreaterequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetRecordimpact")
	values := []entities.RecordcatchildEntity{}
	rows, err := dbc.DB.Query(getimpacturgencysql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid, 6)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordimpact Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordcatchildEntity{}
		rows.Scan(&value.ID, &value.Name, &value.Seqno, &value.Typeid, &value.Typename, &value.Typeseq)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) GetRecordurgency(req *entities.RecordcreaterequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetRecordurgency")
	values := []entities.RecordcatchildEntity{}
	rows, err := dbc.DB.Query(getimpacturgencysql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid, 7)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordurgency Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordcatchildEntity{}
		rows.Scan(&value.ID, &value.Name, &value.Seqno, &value.Typeid, &value.Typename, &value.Typeseq)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) GetRecordprioritydata(req *entities.RecordcreaterequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetRecordprioritydata")
	values := []entities.RecordcatchildEntity{}
	var query string
	var params []interface{}
	if req.Recordimpactid > 0 && req.Recordurgencyid > 0 {
		query = "SELECT  prioritydiff.id, prioritydiff.name, prioritydiff.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstbusinessmatrix, mstrecorddifferentiation prioritydiff, mstrecorddifferentiationtype WHERE mstbusinessmatrix.activeflg = 1 AND prioritydiff.activeflg = 1 AND mstrecorddifferentiationtype.activeflg=1 AND mstbusinessmatrix.deleteflg = 0 AND prioritydiff.deleteflg = 0 AND mstrecorddifferentiationtype.deleteflg=0 AND mstbusinessmatrix.mstrecorddifferentiationpriorityid = prioritydiff.id AND prioritydiff.recorddifftypeid=mstrecorddifferentiationtype.id AND mstbusinessmatrix.clientid = ? AND mstbusinessmatrix.mstorgnhirarchyid = ? AND mstbusinessmatrix.mstrecorddifferentiationtickettypeid = ? AND mstbusinessmatrix.mstrecorddifferentiationimpactid = ? AND mstbusinessmatrix.mstrecorddifferentiationurgencyid = ? "
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.Recordtypeid)
		params = append(params, req.Recordimpactid)
		params = append(params, req.Recordurgencyid)
	} else {
		query = "SELECT  prioritydiff.id, prioritydiff.name, prioritydiff.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstbusinessmatrix, mstrecorddifferentiation prioritydiff, mstrecorddifferentiationtype WHERE mstbusinessmatrix.activeflg = 1 AND prioritydiff.activeflg = 1 AND mstrecorddifferentiationtype.activeflg=1 AND mstbusinessmatrix.deleteflg = 0 AND prioritydiff.deleteflg = 0 AND mstrecorddifferentiationtype.deleteflg=0 AND mstbusinessmatrix.mstrecorddifferentiationpriorityid = prioritydiff.id AND prioritydiff.recorddifftypeid=mstrecorddifferentiationtype.id AND mstbusinessmatrix.clientid = ? AND mstbusinessmatrix.mstorgnhirarchyid = ? AND mstbusinessmatrix.mstrecorddifferentiationtickettypeid = ?  AND mstbusinessmatrix.mstrecorddifferentiationcatid = ? "
		params = append(params, req.Clientid)
		params = append(params, req.Mstorgnhirarchyid)
		params = append(params, req.Recordtypeid)
		params = append(params, req.Recordcatid)
	}
	rows, err := dbc.DB.Query(query, params...)
	if err != nil {
		logger.Log.Println("GetRecordprioritydata Get Statement Prepare Error", err)
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

func (dbc DbConn) GetRecordstatusdata(req *entities.RecordcreaterequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetRecordstatusdata", getfirststatussql)
	values := []entities.RecordcatchildEntity{}
	rows, err := dbc.DB.Query(getfirststatussql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordstatusdata Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordcatchildEntity{}
		rows.Scan(&value.Typeid, &value.Typename, &value.Typeseq, &value.ID, &value.Name, &value.Seqno)
		values = append(values, value)
	}
	//	logger.Log.Println("values is ----->", values)
	return values, nil
}

func (dbc DbConn) GetRecordcatchild(req *entities.RecordcreaterequestEntity) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetRecordcatchild")
	values := []entities.RecordcatchildEntity{}
	rows, err := dbc.DB.Query(getCatchildsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid, req.Recorddiffparentid)

	if err != nil {
		logger.Log.Println("GetAllRecordcategories Get Statement Prepare Error", err)
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

func (dbc DbConn) GetEstimateEffort(req *entities.RecordcreaterequestEntity) ([]string, []string, []string, error) {
	logger.Log.Println("In side GetEstimateEffort")
	var estimatedEfforts []string
	var slaCompliances []string
	var changeTypes []string
	var query = "SELECT mapcategorywithestimatetime.estimatedtime,mapcategorywithestimatetime.efficiency,mapcategorywithestimatetime.changetype FROM mapcategorywithestimatetime WHERE mapcategorywithestimatetime.activeflg=1 AND mapcategorywithestimatetime.deleteflg=0 AND mapcategorywithestimatetime.clientid=? AND mapcategorywithestimatetime.mstorgnhirarchyid=? AND mapcategorywithestimatetime.recorddiffid=?"
	rows, err := dbc.DB.Query(query, req.Clientid, req.Mstorgnhirarchyid, req.Recordcatid)
	if err != nil {
		logger.Log.Println("GetEstimateEffort Get Statement Prepare Error", err)
		return estimatedEfforts, slaCompliances, slaCompliances, err
	}
	defer rows.Close()
	for rows.Next() {
		var estimatedEffort string
		var slaCompliance string
		var changeType string
		scanErr := rows.Scan(&estimatedEffort, &slaCompliance, &changeType)
		if scanErr != nil {
			logger.Log.Println("GetEstimateEffort Scan Error", err)
			return estimatedEfforts, slaCompliances, changeTypes, err
		}
		estimatedEfforts = append(estimatedEfforts, estimatedEffort)
		slaCompliances = append(slaCompliances, slaCompliance)
		changeTypes = append(changeTypes, changeType)
	}
	//	logger.Log.Println("values is ----->", values)
	return estimatedEfforts, slaCompliances, changeTypes, nil
}

func (dbc DbConn) CheckAssetCount(req *entities.RecordcreaterequestEntity) (entities.RecordcreateEntity, error) {
	logger.Log.Println("In side GetAssetCount")
	value := entities.RecordcreateEntity{}
	err := dbc.DB.QueryRow(isassetsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid).Scan(&value.AssetAttached)
	switch err {
	case sql.ErrNoRows:
		value.AssetAttached = 0
		return value, nil
	case nil:
		return value, nil
	default:
		logger.Log.Println("CheckAssetCount Get Statement Prepare Error", err)
		return value, err
	}
}

func (dbc DbConn) CheckAssetCountForDetails(req *entities.RecordDetailsRequestEntity) (entities.RecordcreateEntity, error) {
	logger.Log.Println("In side GetAssetCount")
	value := entities.RecordcreateEntity{}
	err := dbc.DB.QueryRow(isassetsql, req.Clientid, req.Mstorgnhirarchyid, req.RecordDiffTypeid, req.RecordDiffid).Scan(&value.AssetAttached)
	switch err {
	case sql.ErrNoRows:
		value.AssetAttached = 0
		return value, nil
	case nil:
		return value, nil
	default:
		logger.Log.Println("CheckAssetCount Get Statement Prepare Error", err)
		return value, err
	}
}

func (dbc DbConn) GetWorkingCatLabel(req *entities.RecordcreaterequestEntity) (entities.RecordcreateEntity, error) {
	logger.Log.Println("In side GetWorkingCatLabel")
	value := entities.RecordcreateEntity{}
	err := dbc.DB.QueryRow(getworkinglabelsql, req.Clientid, req.Mstorgnhirarchyid, req.Recorddifftypeid, req.Recorddiffid).Scan(&value.WorkingCatLabelID)
	switch err {
	case sql.ErrNoRows:
		value.WorkingCatLabelID = 0
		return value, nil
	case nil:
		return value, nil
	default:
		logger.Log.Println("GetWorkingCatLabel Get Statement Prepare Error", err)
		return value, err
	}
}
