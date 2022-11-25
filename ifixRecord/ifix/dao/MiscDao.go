package dao

import (
	"database/sql"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/utility"
)

func (dbc DbConn) GetFileTermByrecordTypeID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side GetFileTermByrecordTypeID")
	var values []map[string]interface{}
	var query = "SELECT  mstrecordterms.id, mstrecordterms.termname, mstrecordterms.termtypeid, msttermtype.termtypename, mstrecordterms.termvalue FROM mststateterm JOIN mstrecordterms ON mststateterm.recordtermid = mstrecordterms.id AND mststateterm.clientid = mstrecordterms.clientid AND mststateterm.mstorgnhirarchyid = mstrecordterms.mstorgnhirarchyid JOIN msttermtype ON mstrecordterms.termtypeid = msttermtype.id WHERE mststateterm.activeflg = 1 AND mststateterm.deleteflg = 0 AND mstrecordterms.activeflg = 1 AND mstrecordterms.deleteflg = 0 AND msttermtype.id = 3 AND mstrecordterms.clientid = ? AND mstrecordterms.mstorgnhirarchyid = ? AND mststateterm.recorddifftypeid = ? AND mststateterm.recorddiffid = ? ORDER BY mstrecordterms.termname ASC"
	logger.Log.Println("Query String: ", query)
	logger.Log.Println("Params : ", req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"])
	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"])
	if err != nil {
		logger.Log.Println("GetFileTermByrecordTypeID Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		fb := utility.NewDbFieldBind()
		err = fb.PutFields(rows)
		if err != nil {
			return values, err
		}
		rows.Scan(fb.GetFieldPtrArr()...)
		values = append(values, fb.GetFieldArr())
	}
	return values, nil
}

func (dbc DbConn) GetFirstStatusByrecordTypeID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side GetFirstStatusByrecordTypeID")
	var values []map[string]interface{}
	var query = "SELECT mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename,mstrecorddifferentiationtype.seqno typeseq,mstrecorddifferentiation.id,mstrecorddifferentiation.name,mstrecorddifferentiation.seqno FROM mstrecordtype JOIN mstrecorddifferentiation ON mstrecordtype.torecorddiffid=mstrecorddifferentiation.id AND mstrecordtype.torecorddifftypeid=mstrecorddifferentiation.recorddifftypeid AND mstrecordtype.clientid=mstrecorddifferentiation.clientid AND mstrecordtype.mstorgnhirarchyid=mstrecorddifferentiation.mstorgnhirarchyid JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id WHERE mstrecordtype.activeflg=1 AND mstrecordtype.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.seqno=2 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecordtype.clientid=? AND mstrecordtype.mstorgnhirarchyid=? AND mstrecordtype.fromrecorddifftypeid=? AND mstrecordtype.fromrecorddiffid=? ORDER BY mstrecorddifferentiation.seqno ASC LIMIT 1"
	logger.Log.Println("Query String: ", query)
	logger.Log.Println("Params : ", req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"])
	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"])
	if err != nil {
		logger.Log.Println("GetFirstStatusByrecordTypeID Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		fb := utility.NewDbFieldBind()
		err = fb.PutFields(rows)
		if err != nil {
			return values, err
		}
		rows.Scan(fb.GetFieldPtrArr()...)
		values = append(values, fb.GetFieldArr())
	}
	return values, nil
}

func (dbc DbConn) GetRecorddifferentiationbyparent(req map[string]interface{}) ([]entities.RecorddifferentionSingle, error) {
	logger.Log.Println("In side GelAllRecorddifferentiation")
	var getdiffvaluebyparentid = "SELECT id,recorddifftypeid,parentid,name from mstrecorddifferentiation where id=? and activeflg=1 and deleteflg=0 "
	logger.Log.Println(getdiffvaluebyparentid)
	values := []entities.RecorddifferentionSingle{}
	rows, err := dbc.DB.Query(getdiffvaluebyparentid, req["recordcatdiffid"])
	if err != nil {
		logger.Log.Println("GetRecorddifferentiationbyparent Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.RecorddifferentionSingle{}
		rows.Scan(&value.Id, &value.Recorddifftypeid, &value.Parentid, &value.Name)
		values = append(values, value)
	}
	return values, nil
}

func (dbc DbConn) GetPriorityByrecordTypeNCatID(req map[string]interface{}) ([]entities.RecordcatchildEntity, error) {
	logger.Log.Println("In side GetPriorityByrecordTypeNCatID")
	values := []entities.RecordcatchildEntity{}
	var query = "SELECT  prioritydiff.id, prioritydiff.name, prioritydiff.seqno,mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename typename,mstrecorddifferentiationtype.seqno FROM mstbusinessmatrix, mstrecorddifferentiation prioritydiff, mstrecorddifferentiationtype WHERE mstbusinessmatrix.activeflg = 1 AND prioritydiff.activeflg = 1 AND mstrecorddifferentiationtype.activeflg=1 AND mstbusinessmatrix.deleteflg = 0 AND prioritydiff.deleteflg = 0 AND mstrecorddifferentiationtype.deleteflg=0 AND mstbusinessmatrix.mstrecorddifferentiationpriorityid = prioritydiff.id AND prioritydiff.recorddifftypeid=mstrecorddifferentiationtype.id AND mstbusinessmatrix.clientid = ? AND mstbusinessmatrix.mstorgnhirarchyid = ? AND mstbusinessmatrix.mstrecorddifferentiationtickettypeid = ?  AND mstbusinessmatrix.mstrecorddifferentiationcatid = ? "

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordtypediffid"], req["recordcatdiffid"])

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
	return values, nil
}

func (dbc DbConn) GetWorkCatLabel(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side GetWorkCatLabel")
	var query = "SELECT mstworkdifferentiation.mainrecorddifftypeid  FROM mstworkdifferentiation WHERE mstworkdifferentiation.activeflg=1 AND mstworkdifferentiation.deleteflg=0 AND mstworkdifferentiation.clientid=? AND mstworkdifferentiation.mstorgnhirarchyid=? AND mstworkdifferentiation.forrecorddifftypeid=? AND mstworkdifferentiation.forrecorddiffid=?"
	var value int64
	logger.Log.Println("Query String: ", query)
	logger.Log.Println("Params : ", req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"])
	err := dbc.DB.QueryRow(query, req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"]).Scan(&value)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return value, nil
	default:
		logger.Log.Println("GetWorkCatLabel Get Statement Prepare Error", err)
		return 0, err
	}
}

func (dbc DbConn) GetAssetCountByRecordType(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side GetAssetCountByRecordType")
	var query = "SELECT IF(COUNT(mstrecordtype.id)>0,1,0) counter FROM mstrecordtype JOIN mstrecorddifferentiationtype ON mstrecordtype.torecorddifftypeid=mstrecorddifferentiationtype.id WHERE mstrecordtype.activeflg=1 AND mstrecordtype.deleteflg=0 AND mstrecordtype.torecorddiffid=0 AND mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiationtype.seqno=5 AND mstrecordtype.clientid=? AND mstrecordtype.mstorgnhirarchyid=? AND mstrecordtype.fromrecorddifftypeid=? AND mstrecordtype.fromrecorddiffid=?"
	var value int64
	logger.Log.Println("Query String: ", query)
	logger.Log.Println("Params : ", req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"])
	err := dbc.DB.QueryRow(query, req["clientid"], req["mstorgnhirarchyid"], req["recordtypedifftypeid"], req["recordtypediffid"]).Scan(&value)
	switch err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return value, nil
	default:
		logger.Log.Println("GetAssetCountByRecordType Get Statement Prepare Error", err)
		return 0, err
	}
}
