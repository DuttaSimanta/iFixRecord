package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/utility"
	"log"
	"strings"
	"time"
)

func (dbc DbConn) GetAssetAttributesByAseetID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go GetAssetAttributesByAseetID")
	var query = "SELECT mstrecorddifferentiation.id,mstrecorddifferentiation.name,mstrecorddifferentiation.seqno,mstrecorddifferentiation.recorddifftypeid parent FROM mstrecorddifferentiation JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id AND mstrecorddifferentiation.clientid=mstrecorddifferentiationtype.clientid AND mstrecorddifferentiation.mstorgnhirarchyid=mstrecorddifferentiationtype.mstorgnhirarchyid JOIN trnasset ON trnasset.mstdifferentiationtypeid=mstrecorddifferentiationtype.id WHERE mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiation.clientid=? AND mstrecorddifferentiation.mstorgnhirarchyid=? AND trnasset.id =? ORDER BY mstrecorddifferentiation.seqno "
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["assetid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["assetid"])
	if err != nil {
		logger.Log.Println("GetAssetAttributesByAseetID Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	value := map[string]interface{}{"id": 0, "name": "Asset ID", "parent": 0, "seqno": 0}
	values = append(values, value)
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

func (dbc DbConn) GetOnlyAssetIDByID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side side RecordassetDao.go GetOnlyAssetIDByID")
	var query = "SELECT trnasset.id, trnasset.assetid,trnasset.mstdifferentiationtypeid FROM trnasset WHERE trnasset.activeflg = 1 AND trnasset.deleteflg = 0 AND trnasset.clientid = ? AND trnasset.mstorgnhirarchyid = ? AND trnasset.id=?"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["assetid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["assetid"])

	if err != nil {
		logger.Log.Println("GetOnlyAssetIDByID Get Statement Prepare Error", err)
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
	logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) GetAssetAttributesByTypeID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go GetAssetAttributesByTypeID")
	var query = "SELECT mstrecorddifferentiation.id,mstrecorddifferentiation.name,mstrecorddifferentiation.seqno,mstrecorddifferentiation.recorddifftypeid parent FROM mstrecorddifferentiation JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id AND mstrecorddifferentiation.clientid=mstrecorddifferentiationtype.clientid AND mstrecorddifferentiation.mstorgnhirarchyid=mstrecorddifferentiationtype.mstorgnhirarchyid WHERE mstrecorddifferentiationtype.activeflg=1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiation.clientid=? AND mstrecorddifferentiation.mstorgnhirarchyid=? AND mstrecorddifferentiation.recorddifftypeid=? ORDER BY mstrecorddifferentiation.seqno "
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAssetAttributesByTypeID Get Statement Prepare Error", err)
		return values, err
	}
	value := map[string]interface{}{"id": 0, "name": "Asset ID", "parent": 0, "seqno": 0}
	values = append(values, value)
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

func (dbc DbConn) CountAssetByRecordID(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side RecordassetDao.go CountAssetByRecordID")
	var query = "SELECT COUNT(DISTINCT trnasset.id) total FROM trnasset,mapassetdifferentiation,maprecordasset WHERE trnasset.activeflg = 1 AND trnasset.deleteflg = 0 AND mapassetdifferentiation.activeflg = 1 AND mapassetdifferentiation.deleteflg = 0 AND maprecordasset.activeflg =1 AND maprecordasset.deleteflg=0 AND maprecordasset.clientid=trnasset.clientid AND maprecordasset.mstorgnhirarchyid = trnasset.mstorgnhirarchyid AND maprecordasset.assetid=trnasset.id AND trnasset.id = mapassetdifferentiation.trnassetid AND trnasset.clientid = mapassetdifferentiation.clientid AND trnasset.mstorgnhirarchyid = mapassetdifferentiation.mstorgnhirarchyid AND trnasset.clientid = ? AND trnasset.mstorgnhirarchyid = ? AND mapassetdifferentiation.mstdifferentiationtypeid = ? AND maprecordasset.recordid=?"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"], req["recordid"])
	var value int64

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"], req["recordid"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("CountAssetByRecordID Statement Prepare Error", err)
		return value, err
	}

	for rows.Next() {
		rows.Scan(&value)
	}
	logger.Log.Println("Results: ", value)
	return value, nil
}
func (dbc DbConn) GetAllRecordAssetByID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side side RecordassetDao.go GetRecordAssetByID")
	var query = "SELECT trnasset.id, trnasset.assetid FROM trnasset,mapassetdifferentiation,maprecordasset WHERE trnasset.activeflg = 1 AND trnasset.deleteflg = 0 AND mapassetdifferentiation.activeflg = 1 AND mapassetdifferentiation.deleteflg = 0 AND maprecordasset.activeflg =1 AND maprecordasset.deleteflg=0 AND maprecordasset.clientid=trnasset.clientid AND maprecordasset.mstorgnhirarchyid = trnasset.mstorgnhirarchyid AND maprecordasset.assetid=trnasset.id AND trnasset.id = mapassetdifferentiation.trnassetid AND trnasset.clientid = mapassetdifferentiation.clientid AND trnasset.mstorgnhirarchyid = mapassetdifferentiation.mstorgnhirarchyid AND trnasset.clientid = ? AND trnasset.mstorgnhirarchyid = ? AND mapassetdifferentiation.mstdifferentiationtypeid = ? AND maprecordasset.recordid=? GROUP BY trnasset.id"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"], req["recordid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"], req["recordid"])

	if err != nil {
		logger.Log.Println("GetRecordAssetByID Get Statement Prepare Error", err)
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
	logger.Log.Println("Results: ", values)
	return values, nil
}
func (dbc DbConn) GetRecordAssetByID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side side RecordassetDao.go GetRecordAssetByID")
	var query = "SELECT trnasset.id, trnasset.assetid FROM trnasset,mapassetdifferentiation,maprecordasset WHERE trnasset.activeflg = 1 AND trnasset.deleteflg = 0 AND mapassetdifferentiation.activeflg = 1 AND mapassetdifferentiation.deleteflg = 0 AND maprecordasset.activeflg =1 AND maprecordasset.deleteflg=0 AND maprecordasset.clientid=trnasset.clientid AND maprecordasset.mstorgnhirarchyid = trnasset.mstorgnhirarchyid AND maprecordasset.assetid=trnasset.id AND trnasset.id = mapassetdifferentiation.trnassetid AND trnasset.clientid = mapassetdifferentiation.clientid AND trnasset.mstorgnhirarchyid = mapassetdifferentiation.mstorgnhirarchyid AND trnasset.clientid = ? AND trnasset.mstorgnhirarchyid = ? AND mapassetdifferentiation.mstdifferentiationtypeid = ? AND maprecordasset.recordid=? GROUP BY trnasset.id LIMIT ?,?"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"], req["recordid"], req["offset"], req["limit"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"], req["recordid"], req["offset"], req["limit"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordAssetByID Get Statement Prepare Error", err)
		return values, err
	}

	for rows.Next() {
		fb := utility.NewDbFieldBind()
		err = fb.PutFields(rows)
		if err != nil {
			return values, err
		}
		rows.Scan(fb.GetFieldPtrArr()...)
		values = append(values, fb.GetFieldArr())
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) GetAssetTypeIDByRecordID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side side RecordassetDao.go GetAssetTypeIDByRecordID")
	var query = "SELECT trnasset.mstdifferentiationtypeid id, mstrecorddifferentiationtype.typename name FROM trnasset JOIN mstrecorddifferentiationtype ON trnasset.mstdifferentiationtypeid=mstrecorddifferentiationtype.id JOIN maprecordasset ON maprecordasset.clientid = trnasset.clientid AND maprecordasset.mstorgnhirarchyid = trnasset.mstorgnhirarchyid AND maprecordasset.assetid = trnasset.id JOIN mapassetdifferentiation ON trnasset.id = mapassetdifferentiation.trnassetid AND trnasset.clientid = mapassetdifferentiation.clientid AND trnasset.mstorgnhirarchyid = mapassetdifferentiation.mstorgnhirarchyid WHERE trnasset.activeflg = 1 AND trnasset.deleteflg = 0 AND mapassetdifferentiation.activeflg = 1 AND mapassetdifferentiation.deleteflg = 0 AND maprecordasset.activeflg = 1 AND maprecordasset.deleteflg = 0 AND trnasset.clientid = ? AND trnasset.mstorgnhirarchyid = ? AND maprecordasset.recordid = ? GROUP BY trnasset.mstdifferentiationtypeid ORDER BY trnasset.mstdifferentiationtypeid ASC , trnasset.id ASC"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["recordid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordid"])

	if err != nil {
		logger.Log.Println("GetAssetTypeIDByRecordID Get Statement Prepare Error", err)
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

/*func (dbc DbConn) GetAssetAttributesByAssetID(req map[string]interface{}) ([]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go GetAssetAttributesByAssetID")
	var query = "SELECT  mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename,mstrecorddifferentiation.id attrid,mstrecorddifferentiation.name attrname,IF(mapassetdifferentiation.value IS NULL,'',mapassetdifferentiation.value) value FROM mstrecorddifferentiationtype JOIN mstrecorddifferentiation ON mstrecorddifferentiationtype.id=mstrecorddifferentiation.recorddifftypeid AND mstrecorddifferentiationtype.clientid=mstrecorddifferentiation.clientid AND mstrecorddifferentiationtype.mstorgnhirarchyid=mstrecorddifferentiation.mstorgnhirarchyid  LEFT JOIN trnasset ON mstrecorddifferentiationtype.id = trnasset.mstdifferentiationtypeid AND trnasset.id=? AND trnasset.activeflg=1 AND trnasset.deleteflg=0 LEFT JOIN mapassetdifferentiation ON trnasset.id=mapassetdifferentiation.trnassetid AND mapassetdifferentiation.mstdifferentiationid=mstrecorddifferentiation.id AND mapassetdifferentiation.mstdifferentiationtypeid=mstrecorddifferentiationtype.id AND mapassetdifferentiation.activeflg=1 AND mapassetdifferentiation.deleteflg=0 WHERE mstrecorddifferentiationtype.activeflg =1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0  AND mstrecorddifferentiation.clientid=? AND mstrecorddifferentiationtype.mstorgnhirarchyid=? AND mstrecorddifferentiationtype.id=? ORDER BY mstrecorddifferentiation.seqno ASC"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["id"], req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"])
	var values []interface{}

	rows, err := dbc.DB.Query(query, req["id"], req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAssetAttributesByAssetID Get Statement Prepare Error", err)
		return values, err
	}

	for rows.Next() {
		fb := utility.NewDbFieldBind()
		err = fb.PutFields(rows)
		if err != nil {
			return values, err
		}
		rows.Scan(fb.GetFieldPtrArr()...)
		values = append(values, fb.GetFieldArr())
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}*/

func (dbc DbConn) GetAssetAttributesByAssetID(req map[string]interface{}) ([]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go GetAssetAttributesByAssetID")
	var query = "SELECT  mstrecorddifferentiationtype.id typeid,mstrecorddifferentiationtype.typename,mstrecorddifferentiation.id attrid,mstrecorddifferentiation.name attrname,IF(mapassetdifferentiation.value IS NULL,'',mapassetdifferentiation.value) value FROM mstrecorddifferentiationtype JOIN mstrecorddifferentiation ON mstrecorddifferentiationtype.id=mstrecorddifferentiation.recorddifftypeid AND mstrecorddifferentiationtype.clientid=mstrecorddifferentiation.clientid AND mstrecorddifferentiationtype.mstorgnhirarchyid=mstrecorddifferentiation.mstorgnhirarchyid  LEFT JOIN trnasset ON mstrecorddifferentiationtype.id = trnasset.mstdifferentiationtypeid AND trnasset.id=? AND trnasset.activeflg=1 AND trnasset.deleteflg=0 LEFT JOIN mapassetdifferentiation ON trnasset.id=mapassetdifferentiation.trnassetid AND mapassetdifferentiation.mstdifferentiationid=mstrecorddifferentiation.id AND mapassetdifferentiation.mstdifferentiationtypeid=mstrecorddifferentiationtype.id AND mapassetdifferentiation.activeflg=1 AND mapassetdifferentiation.deleteflg=0 WHERE mstrecorddifferentiationtype.activeflg =1 AND mstrecorddifferentiationtype.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0  AND mstrecorddifferentiation.clientid=? AND mstrecorddifferentiationtype.mstorgnhirarchyid=? AND mstrecorddifferentiationtype.id=? ORDER BY mstrecorddifferentiation.seqno ASC"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["id"], req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"])
	var values []interface{}

	rows, err := dbc.DB.Query(query, req["id"], req["clientid"], req["mstorgnhirarchyid"], req["mstdifferentiationtypeid"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAssetAttributesByAssetID Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		fb := utility.NewDbFieldBind()
		err = fb.PutFields(rows)
		if err != nil {
			return values, err
		}
		rows.Scan(fb.GetFieldPtrArr()...)
		values = append(values, fb.GetFieldArr())
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) GetRecordIDByID(req int64) (string, error) {
	var query = "SELECT trnrecord.code FROM trnrecord WHERE trnrecord.id=?"
	var recordID = ""
	rows, err := dbc.DB.Query(query, req)
	if err != nil {
		logger.Log.Println("getRecordIDByID Get Statement Prepare Error", err)
		return recordID, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&recordID)
	}
	return recordID, nil
}

func (dbc DbConn) GetAssetIDByID(req int64) (string, error) {
	var query = "SELECT trnasset.assetid FROM trnasset WHERE id=?"
	var recordID = ""
	rows, err := dbc.DB.Query(query, req)
	if err != nil {
		logger.Log.Println("GetAssetIDByID Get Statement Prepare Error", err)
		return recordID, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&recordID)
	}
	return recordID, nil
}

func (dbc DbConn) GetRecordIDByIDInterface(req map[string]interface{}) (string, error) {
	var query = "SELECT trnrecord.code FROM trnrecord WHERE trnrecord.id=?"
	var recordID = ""
	rows, err := dbc.DB.Query(query, req["recordid"])
	if err != nil {
		logger.Log.Println("getRecordIDByID Get Statement Prepare Error", err)
		return recordID, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&recordID)
	}
	return recordID, nil
}

func (dbc DbConn) GetAssetIDByIDInterface(req map[string]interface{}) (string, error) {
	var query = "SELECT trnasset.assetid FROM trnasset WHERE id=?"
	var recordID = ""
	rows, err := dbc.DB.Query(query, req["assetid"])
	if err != nil {
		logger.Log.Println("GetAssetIDByID Get Statement Prepare Error", err)
		return recordID, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&recordID)
	}
	return recordID, nil
}
func (dbc DbConn) InsertTrnAsset(tx *sql.Tx, trnAsset *entities.TrnAsset) (int64, error) {

	var InsertTrnAssetQuery string = "INSERT INTO trnasset(`clientid`,`mstorgnhirarchyid`,`mstdifferentiationtypeid`,`assetid`,`additionalattr`,`deleteflg`,`activeflg`) VALUES(?,?,?,?,?,?,?)"
	stmtInsertTrnAsset, stmtErr := tx.Prepare(InsertTrnAssetQuery)
	if stmtErr != nil {
		tx.Rollback()
		logger.Log.Println(stmtErr)
		return 0, errors.New("ERROR: Unable to Prepare Statement")
	}
	defer stmtInsertTrnAsset.Close()
	InsertTrnAssetResultSet, insertErr := stmtInsertTrnAsset.Exec(trnAsset.ClientId, trnAsset.MstOrgnHirarchyId, trnAsset.MstDifftypeid, trnAsset.AssetId, trnAsset.AdditionalAttr, trnAsset.DeleteFlag, trnAsset.ActiveFlag)
	if insertErr != nil {
		tx.Rollback()
		logger.Log.Println(insertErr)
		return 0, errors.New("ERROR: InsertTrnAssetResultSet  Error")
	}
	lastInsertedTrnAssetId, lastIdError := InsertTrnAssetResultSet.LastInsertId()
	if lastIdError != nil {
		tx.Rollback()
		logger.Log.Println(lastIdError)
		return 0, errors.New("ERROR: last Inserted Id Error in TrnAsset")
	}
	return lastInsertedTrnAssetId, nil
}
func (dbc DbConn) GetAssetIdToInsert(tx *sql.Tx, clientID int64, mstOrgnHirarchyId int64) (string, error) {
	//var code string
	var lastID int64
	var code string
	var getLastAssetIdQuery string = "select code,uid from uidgen where clientid=? and mstorgnhirarchyid=? and difftypeid=6"
	err := dbc.DB.QueryRow(getLastAssetIdQuery, clientID, mstOrgnHirarchyId).Scan(&code, &lastID)
	if err != nil {
		logger.Log.Println(err)
		return "", errors.New("ERROR: Unable to Scan Asset Id Details")
	}
	//log.Println("lastid %0.7d",lastID)
	lastID++
	assetID := fmt.Sprintf("%0.7d", lastID)
	assetID = code + assetID
	var updateAssetLastIdQuery string = "update uidgen set uid=? where clientid=? and mstorgnhirarchyid=? and difftypeid=6"
	stmtAssetLastId, stmtErr := tx.Prepare(updateAssetLastIdQuery)
	if stmtErr != nil {
		logger.Log.Println(stmtErr)
		return "", errors.New("ERROR: Unable to Prepare Statement")
	}
	defer stmtAssetLastId.Close()
	//var scanLastID int64=0
	res, err := stmtAssetLastId.Exec(lastID, clientID, mstOrgnHirarchyId)
	if err != nil {
		logger.Log.Println(err)
		return "", errors.New("ERROR: ResultSet Fetching Error")
	}
	count, err := res.RowsAffected()
	if err != nil {
		logger.Log.Println(err)
		return "", errors.New("ERROR: ResultSet Fetching Error")
	}
	if count != 1 {
		logger.Log.Println("Row not Updated Properly", count)
		return "", errors.New("ERROR: Row not Updated Properly")
	}
	logger.Log.Println("RowsAffeCted==> ", count)
	return assetID, nil
}
func (dbc DbConn) InsertMapAssetDiff(tx *sql.Tx, mapAssetDiff *entities.MapAssetDifferentiation) error {

	var InsertMapAssetDiffQuery string = "INSERT INTO mapassetdifferentiation(`clientid`,`mstorgnhirarchyid`,`mstdifferentiationtypeid`,`mstdifferentiationid`,`trnassetid`,`value`,`deleteflg`,`activeflg`,`audittransactionid`) VALUES(?,?,?,?,?,?,?,?,?)"
	stmtInsertMapAssetDiff, stmtErr := tx.Prepare(InsertMapAssetDiffQuery)
	if stmtErr != nil {
		tx.Rollback()
		logger.Log.Println(stmtErr)
		return errors.New("ERROR: Unable to Prepare Statement")
	}
	defer stmtInsertMapAssetDiff.Close()
	InsertMapAssetDiffResultSet, insertErr := stmtInsertMapAssetDiff.Exec(mapAssetDiff.Clientid, mapAssetDiff.Mstorgnhirarchyid, mapAssetDiff.Mstdifferentiationtypeid, mapAssetDiff.Mstdifferentiationid, mapAssetDiff.Trnassetid, mapAssetDiff.Value, mapAssetDiff.Deleteflg, mapAssetDiff.Activeflg, mapAssetDiff.AuditTransactionId)
	if insertErr != nil {
		tx.Rollback()
		logger.Log.Println(insertErr)
		return errors.New("ERROR: InsertTrnAssetResultSet  Error")
	}
	lastInsertedMapAssetDIffId, lastIdError := InsertMapAssetDiffResultSet.LastInsertId()
	if lastIdError != nil {
		tx.Rollback()
		logger.Log.Println(lastIdError)
		return errors.New("ERROR: Last Inserted MapAssetDiff id  fetch  Error")
	}
	if lastInsertedMapAssetDIffId == 0 {
		tx.Rollback()
		logger.Log.Println("lastInsertedMapAssetDIffId is zero")
		return errors.New("ERROR: Last Inserted MapAssetDiff id  fetch  Error")
	}
	return nil
}
func (dbc DbConn) InsertRecordAsset(tx *sql.Tx, req *entities.InsertRecordAssetEntity, assetIDS string) ([]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go InsertRecordAsset")
	var values []interface{}
	var ticketID string
	getTicketIDQuery := "select code from trnrecord where clientid=? and mstorgnhirarchyid=? and id=?"
	getTicketIDErr := dbc.DB.QueryRow(getTicketIDQuery, req.Clientid, req.Mstorgnhirarchyid, req.Recordid).Scan(&ticketID)
	if getTicketIDErr != nil {
		logger.Log.Println("getTicketIDErr in UpdateRecordAsset==> ", getTicketIDErr)
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
	}
	logger.Log.Println("AssetID===>", req.Assetid)
	var count int64
	checkRecordAssetMapQuery := "select count(id) as count from maprecordasset where clientid=? and mstorgnhirarchyid=? and recordid=? and assetid=? and activeflg=1 and deleteflg=0 "
	checkRecordAssetMapErr := dbc.DB.QueryRow(checkRecordAssetMapQuery, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.Assetid).Scan(&count)
	if checkRecordAssetMapErr != nil {
		logger.Log.Println("checkRecordAssetMapErr in UpdateRecordAsset==> ", checkRecordAssetMapErr)
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
		return values, errors.New("SQL Execution Error")
	}
	if count <= 0 {
		var addAssetWithRecordquery = "INSERT INTO maprecordasset (clientid, mstorgnhirarchyid, recordid, recordstageid, assetid, trnassetid, deleteflg, activeflg) VALUES(?, ?, ?, ?, ?, ?, 0, 1)"
		addAssetWithRecordstmnt, addAssetWithRecordPreparestmntErr := tx.Prepare(addAssetWithRecordquery)
		if addAssetWithRecordPreparestmntErr != nil {
			logger.Log.Println("addAssetWithRecordPreparestmntErr in UpdateRecordAsset==> ", addAssetWithRecordPreparestmntErr)
			tx.Rollback()
			return values, addAssetWithRecordPreparestmntErr

		}
		defer addAssetWithRecordstmnt.Close()
		_, addAssetWithRecordstmnterr := addAssetWithRecordstmnt.Exec(req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordStageid, req.Assetid, req.Assetid)
		if addAssetWithRecordstmnterr != nil {
			logger.Log.Println(addAssetWithRecordstmnterr)
			tx.Rollback()
			return values, errors.New("SQL Execution Error")
		}
		recordCode := ticketID
		var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"
		err := InsertActivityLogs(tx, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, 7, logData, req.Userid, req.GroupID)
		if err != nil {
			log.Println("error is ----->", err)
			tx.Rollback()
			return values, err
		}
		if req.ParentRecordID != 0 {
			var parentRecordStagedID int64
			var parentTicketID string
			getTicketIDQuery := "select code,recordstageid from trnrecord where clientid=? and mstorgnhirarchyid=? and id=?"
			getTicketIDErr := dbc.DB.QueryRow(getTicketIDQuery, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID).Scan(&parentTicketID, &parentRecordStagedID)
			if getTicketIDErr != nil {
				logger.Log.Println("getTicketIDErr in UpdateRecordAsset==> ", getTicketIDErr)
				//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
			}

			checkRecordAssetMapQuery := "select count(id) as count from maprecordasset where clientid=? and mstorgnhirarchyid=? and recordid=? and assetid=? and activeflg=1 and deleteflg=0 "
			checkRecordAssetMapErr := dbc.DB.QueryRow(checkRecordAssetMapQuery, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, req.Assetid).Scan(&count)
			if checkRecordAssetMapErr != nil {
				logger.Log.Println("checkRecordAssetMapErr in UpdateRecordAsset==> ", checkRecordAssetMapErr)
				//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
				return values, errors.New("SQL Execution Error")
			}
			var addAssetWithRecordquery = "INSERT INTO maprecordasset (clientid, mstorgnhirarchyid, recordid, recordstageid, assetid, trnassetid, deleteflg, activeflg) VALUES(?, ?, ?, ?, ?, ?, 0, 1)"
			addAssetWithRecordstmnt, addAssetWithRecordPreparestmntErr := tx.Prepare(addAssetWithRecordquery)
			if addAssetWithRecordPreparestmntErr != nil {
				logger.Log.Println("addAssetWithRecordPreparestmntErr in UpdateRecordAsset==> ", addAssetWithRecordPreparestmntErr)
				tx.Rollback()
				return values, addAssetWithRecordPreparestmntErr

			}
			defer addAssetWithRecordstmnt.Close()
			_, addAssetWithRecordstmnterr := addAssetWithRecordstmnt.Exec(req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, parentRecordStagedID, req.Assetid, req.Assetid)
			if addAssetWithRecordstmnterr != nil {
				logger.Log.Println(addAssetWithRecordstmnterr)
				tx.Rollback()
				return values, errors.New("SQL Execution Error")
			}
			recordCode := parentTicketID
			var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"
			err := InsertActivityLogs(tx, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, 7, logData, req.Userid, req.GroupID)
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				return values, err
			}
		}

	}

	return values, nil
}

/*func (dbc DbConn) UpdateRecordAsset(tx *sql.Tx, req *entities.UpdateRecordAssetEntity, timediff int64) ([]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go UpdateRecordAsset")
	var values []interface{}
	var ticketID string
	var assetHeaderName string
	var assetHistory string
	var assetOldValue string
	var updatedBy string
	getTicketIDQuery := "select code from trnrecord where clientid=? and mstorgnhirarchyid=? and id=?"
	getTicketIDErr := dbc.DB.QueryRow(getTicketIDQuery, req.Clientid, req.Mstorgnhirarchyid, req.Recordid).Scan(&ticketID)
	if getTicketIDErr != nil {
		logger.Log.Println("getTicketIDErr in UpdateRecordAsset==> ", getTicketIDErr)
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
	}
	getAssetHeaderQuery := "select name from mstrecorddifferentiation where clientid=? and mstorgnhirarchyid=? and id=?"
	getAssetHeaderErr := dbc.DB.QueryRow(getAssetHeaderQuery, req.Clientid, req.Mstorgnhirarchyid, req.AssetHeaderId).Scan(&assetHeaderName)
	if getAssetHeaderErr != nil {
		logger.Log.Println("getAssetHeaderErr in UpdateRecordAsset==> ", getAssetHeaderErr)
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
	}
	getOldAssetValueQuery := "select value from mapassetdifferentiation where clientid=? and mstorgnhirarchyid=? and mstdifferentiationid=? and trnassetid =?"
	getOldAssetValueErr := dbc.DB.QueryRow(getOldAssetValueQuery, req.Clientid, req.Mstorgnhirarchyid, req.AssetHeaderId, req.Assetid).Scan(&assetOldValue)
	if getOldAssetValueErr != nil {
		logger.Log.Println("getOldAssetValueErr in UpdateRecordAsset==> ", getOldAssetValueErr)
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
	}
	getUserNameQuery := "select name from mstclientuser where clientid=? and mstorgnhirarchyid=? and id=?"
	getUserNameErr := dbc.DB.QueryRow(getUserNameQuery, req.Clientid, req.Mstorgnhirarchyid, req.Userid).Scan(&updatedBy)
	if getUserNameErr != nil {
		logger.Log.Println("getUserNameErr in UpdateRecordAsset==> ", getUserNameErr)
		return values, getUserNameErr
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
	}

	updateAssetQuery := "update mapassetdifferentiation set value=? where clientid=? and mstorgnhirarchyid=? and mstdifferentiationid=? and trnassetid=?"
	updateAssetstmt, updateAssetPrepareStmnterr := tx.Prepare(updateAssetQuery)
	if updateAssetPrepareStmnterr != nil {
		logger.Log.Println("prepareStmnterr in UpdateRecordAsset==> ", updateAssetPrepareStmnterr)
		return values, updateAssetPrepareStmnterr

	}
	//stmt, err := tx.Prepare(insertlogs)
	defer updateAssetstmt.Close()
	_, updateAssetstmterr := updateAssetstmt.Exec(req.UpdatedValue, req.Clientid, req.Mstorgnhirarchyid, req.AssetHeaderId, req.Assetid)
	if updateAssetstmterr != nil {
		logger.Log.Println(updateAssetstmterr)
		return values, errors.New("SQL Execution Error")
	}
	var assetIDS string
	getAssetHistoryQuery := "select coalesce(history,'') as history,assetid from trnasset where clientid=? and mstorgnhirarchyid=? and  id=?"
	getAssetHistoryErr := dbc.DB.QueryRow(getAssetHistoryQuery, req.Clientid, req.Mstorgnhirarchyid, req.Assetid).Scan(&assetHistory, &assetIDS)
	if getAssetHistoryErr != nil {
		logger.Log.Println("getAssetHistoryErr in UpdateRecordAsset==> ", getAssetHistoryErr)
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
		return values, getAssetHistoryErr
	}
	//var util = entities.UtilityEntity{}

	now := time.Now()
	time := Convertdate(int64(now.Unix()), timediff)
	//now := time.Now()
	//stringTime := now.Local().Format("03-Oct-2021 19:41:14") //string(now.Unix())
	var logvalue = "" + "Updated Date/Time : " + time + " " + "Column Name : " + assetHeaderName + " has been changed from " + assetOldValue + " to " + req.UpdatedValue + "\n" + "Ref. ID : " + ticketID + "\n" + "Updated By : " + updatedBy

	assetHistory = logvalue + "::" + assetHistory
	updateAssetHistoryQuery := "update trnasset set history=? where clientid=? and mstorgnhirarchyid=? and id=?	"

	updateAssetHistorystmnt, updateAssetHistoryPreparestmntErr := tx.Prepare(updateAssetHistoryQuery)
	if updateAssetHistoryPreparestmntErr != nil {
		logger.Log.Println("updateAssetHistoryPreparestmntErr in UpdateRecordAsset==> ", updateAssetHistoryPreparestmntErr)
		tx.Rollback()
		return values, updateAssetHistoryPreparestmntErr

	}
	defer updateAssetHistorystmnt.Close()
	_, updateAssetHistorystmnterr := updateAssetHistorystmnt.Exec(assetHistory, req.Clientid, req.Mstorgnhirarchyid, req.Assetid)
	if updateAssetHistorystmnterr != nil {
		logger.Log.Println(updateAssetHistorystmnterr)
		tx.Rollback()
		return values, errors.New("SQL Execution Error")
	}
	var count int64
	checkRecordAssetMapQuery := "select count(id) as count from maprecordasset where clientid=? and mstorgnhirarchyid=? and recordid=? and assetid=? and activeflg=1 and deleteflg=0 "
	checkRecordAssetMapErr := dbc.DB.QueryRow(checkRecordAssetMapQuery, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.Assetid).Scan(&count)
	if checkRecordAssetMapErr != nil {
		logger.Log.Println("checkRecordAssetMapErr in UpdateRecordAsset==> ", checkRecordAssetMapErr)
		//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
		return values, errors.New("SQL Execution Error")
	}
	if count <= 0 {
		var addAssetWithRecordquery = "INSERT INTO maprecordasset (clientid, mstorgnhirarchyid, recordid, recordstageid, assetid, trnassetid, deleteflg, activeflg) VALUES(?, ?, ?, ?, ?, ?, 0, 1)"
		addAssetWithRecordstmnt, addAssetWithRecordPreparestmntErr := tx.Prepare(addAssetWithRecordquery)
		if addAssetWithRecordPreparestmntErr != nil {
			logger.Log.Println("addAssetWithRecordPreparestmntErr in UpdateRecordAsset==> ", addAssetWithRecordPreparestmntErr)
			tx.Rollback()
			return values, updateAssetHistoryPreparestmntErr

		}
		defer addAssetWithRecordstmnt.Close()
		_, addAssetWithRecordstmnterr := addAssetWithRecordstmnt.Exec(req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordStageid, req.Assetid, req.Assetid)
		if addAssetWithRecordstmnterr != nil {
			logger.Log.Println(addAssetWithRecordstmnterr)
			tx.Rollback()
			return values, errors.New("SQL Execution Error")
		}
		recordCode := ticketID
		var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"
		err := InsertActivityLogs(tx, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, 7, logData, req.Userid, req.GroupID)
		if err != nil {
			log.Println("error is ----->", err)
			tx.Rollback()
			return values, err
		}
		if req.ParentRecordID != 0 {
			var parentRecordStagedID int64
			var parentTicketID string
			getTicketIDQuery := "select code,recordstageid from trnrecord where clientid=? and mstorgnhirarchyid=? and id=?"
			getTicketIDErr := dbc.DB.QueryRow(getTicketIDQuery, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID).Scan(&parentTicketID, &parentRecordStagedID)
			if getTicketIDErr != nil {
				logger.Log.Println("getTicketIDErr in UpdateRecordAsset==> ", getTicketIDErr)
				//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
			}

			checkRecordAssetMapQuery := "select count(id) as count from maprecordasset where clientid=? and mstorgnhirarchyid=? and recordid=? and assetid=? and activeflg=1 and deleteflg=0 "
			checkRecordAssetMapErr := dbc.DB.QueryRow(checkRecordAssetMapQuery, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, req.Assetid).Scan(&count)
			if checkRecordAssetMapErr != nil {
				logger.Log.Println("checkRecordAssetMapErr in UpdateRecordAsset==> ", checkRecordAssetMapErr)
				//	return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")
				return values, errors.New("SQL Execution Error")
			}
			var addAssetWithRecordquery = "INSERT INTO maprecordasset (clientid, mstorgnhirarchyid, recordid, recordstageid, assetid, trnassetid, deleteflg, activeflg) VALUES(?, ?, ?, ?, ?, ?, 0, 1)"
			addAssetWithRecordstmnt, addAssetWithRecordPreparestmntErr := tx.Prepare(addAssetWithRecordquery)
			if addAssetWithRecordPreparestmntErr != nil {
				logger.Log.Println("addAssetWithRecordPreparestmntErr in UpdateRecordAsset==> ", addAssetWithRecordPreparestmntErr)
				tx.Rollback()
				return values, updateAssetHistoryPreparestmntErr

			}
			defer addAssetWithRecordstmnt.Close()
			_, addAssetWithRecordstmnterr := addAssetWithRecordstmnt.Exec(req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, parentRecordStagedID, req.Assetid, req.Assetid)
			if addAssetWithRecordstmnterr != nil {
				logger.Log.Println(addAssetWithRecordstmnterr)
				tx.Rollback()
				return values, errors.New("SQL Execution Error")
			}
			recordCode := parentTicketID
			var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"
			err := InsertActivityLogs(tx, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, 7, logData, req.Userid, req.GroupID)
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				return values, err
			}
		}

	}

	return values, nil
}*/

func (dbc DbConn) UpdateRecordAsset(tx *sql.Tx, req *entities.UpdateRecordAssetEntity, timediff int64) ([]interface{}, error) {

	logger.Log.Println("In side RecordassetDao.go UpdateRecordAsset")

	var values []interface{}

	var ticketID string

	var assetHeaderName string

	var assetHistory string

	var assetOldValue string

	var updatedBy string

	getTicketIDQuery := "select code from trnrecord where clientid=? and mstorgnhirarchyid=? and id=?"

	getTicketIDErr := dbc.DB.QueryRow(getTicketIDQuery, req.Clientid, req.Mstorgnhirarchyid, req.Recordid).Scan(&ticketID)

	if getTicketIDErr != nil {

		logger.Log.Println("getTicketIDErr in UpdateRecordAsset==> ", getTicketIDErr)

		//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

	}

	getAssetHeaderQuery := "select name from mstrecorddifferentiation where clientid=? and mstorgnhirarchyid=? and id=?"

	getAssetHeaderErr := dbc.DB.QueryRow(getAssetHeaderQuery, req.Clientid, req.Mstorgnhirarchyid, req.AssetHeaderId).Scan(&assetHeaderName)

	if getAssetHeaderErr != nil {

		logger.Log.Println("getAssetHeaderErr in UpdateRecordAsset==> ", getAssetHeaderErr)

		//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

	}

	getOldAssetValueQuery := "select value from mapassetdifferentiation where clientid=? and mstorgnhirarchyid=? and mstdifferentiationid=? and trnassetid =?"

	getOldAssetValueErr := dbc.DB.QueryRow(getOldAssetValueQuery, req.Clientid, req.Mstorgnhirarchyid, req.AssetHeaderId, req.Assetid).Scan(&assetOldValue)

	if getOldAssetValueErr != nil {

		logger.Log.Println("getOldAssetValueErr in UpdateRecordAsset==> ", getOldAssetValueErr)

		//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

	}

	getUserNameQuery := "select name from mstclientuser where id=?"

	getUserNameErr := dbc.DB.QueryRow(getUserNameQuery, req.Userid).Scan(&updatedBy)

	if getUserNameErr != nil {

		logger.Log.Println("getUserNameErr in UpdateRecordAsset==> ", getUserNameErr)

		return values, getUserNameErr

		//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

	}

	updateAssetQuery := "update mapassetdifferentiation set value=? where clientid=? and mstorgnhirarchyid=? and mstdifferentiationid=? and trnassetid=?"

	updateAssetstmt, updateAssetPrepareStmnterr := tx.Prepare(updateAssetQuery)

	if updateAssetPrepareStmnterr != nil {

		logger.Log.Println("prepareStmnterr in UpdateRecordAsset==> ", updateAssetPrepareStmnterr)

		return values, updateAssetPrepareStmnterr

	}

	//stmt, err := tx.Prepare(insertlogs)

	defer updateAssetstmt.Close()

	_, updateAssetstmterr := updateAssetstmt.Exec(req.UpdatedValue, req.Clientid, req.Mstorgnhirarchyid, req.AssetHeaderId, req.Assetid)

	if updateAssetstmterr != nil {

		logger.Log.Println(updateAssetstmterr)

		return values, errors.New("SQL Execution Error")

	}

	var assetIDS string

	getAssetHistoryQuery := "select coalesce(history,'') as history,assetid from trnasset where clientid=? and mstorgnhirarchyid=? and  id=?"

	getAssetHistoryErr := dbc.DB.QueryRow(getAssetHistoryQuery, req.Clientid, req.Mstorgnhirarchyid, req.Assetid).Scan(&assetHistory, &assetIDS)

	if getAssetHistoryErr != nil {

		logger.Log.Println("getAssetHistoryErr in UpdateRecordAsset==> ", getAssetHistoryErr)

		//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

		return values, getAssetHistoryErr

	}

	//var util = entities.UtilityEntity{}

	now := time.Now()

	time := Convertdate(int64(now.Unix()), timediff)

	//now := time.Now()

	//stringTime := now.Local().Format("03-Oct-2021 19:41:14") //string(now.Unix())

	var logvalue = "" + "Updated Date/Time : " + time + " " + "Column Name : " + assetHeaderName + " has been changed from " + assetOldValue + " to " + req.UpdatedValue + "\n" + "Ref. ID : " + ticketID + "\n" + "Updated By : " + updatedBy

	assetHistory = logvalue + "::" + assetHistory

	updateAssetHistoryQuery := "update trnasset set history=? where clientid=? and mstorgnhirarchyid=? and id=?    "

	updateAssetHistorystmnt, updateAssetHistoryPreparestmntErr := tx.Prepare(updateAssetHistoryQuery)

	if updateAssetHistoryPreparestmntErr != nil {

		logger.Log.Println("updateAssetHistoryPreparestmntErr in UpdateRecordAsset==> ", updateAssetHistoryPreparestmntErr)

		tx.Rollback()

		return values, updateAssetHistoryPreparestmntErr

	}

	defer updateAssetHistorystmnt.Close()

	_, updateAssetHistorystmnterr := updateAssetHistorystmnt.Exec(assetHistory, req.Clientid, req.Mstorgnhirarchyid, req.Assetid)

	if updateAssetHistorystmnterr != nil {

		logger.Log.Println(updateAssetHistorystmnterr)

		tx.Rollback()

		return values, errors.New("SQL Execution Error")

	}

	var count int64

	checkRecordAssetMapQuery := "select count(id) as count from maprecordasset where clientid=? and mstorgnhirarchyid=? and recordid=? and assetid=? and activeflg=1 and deleteflg=0 "

	checkRecordAssetMapErr := dbc.DB.QueryRow(checkRecordAssetMapQuery, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.Assetid).Scan(&count)

	if checkRecordAssetMapErr != nil {

		logger.Log.Println("checkRecordAssetMapErr in UpdateRecordAsset==> ", checkRecordAssetMapErr)

		//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

		return values, errors.New("SQL Execution Error")

	}

	if count <= 0 {

		var addAssetWithRecordquery = "INSERT INTO maprecordasset (clientid, mstorgnhirarchyid, recordid, recordstageid, assetid, trnassetid, deleteflg, activeflg) VALUES(?, ?, ?, ?, ?, ?, 0, 1)"

		addAssetWithRecordstmnt, addAssetWithRecordPreparestmntErr := tx.Prepare(addAssetWithRecordquery)

		if addAssetWithRecordPreparestmntErr != nil {

			logger.Log.Println("addAssetWithRecordPreparestmntErr in UpdateRecordAsset==> ", addAssetWithRecordPreparestmntErr)

			tx.Rollback()

			return values, updateAssetHistoryPreparestmntErr

		}

		defer addAssetWithRecordstmnt.Close()

		_, addAssetWithRecordstmnterr := addAssetWithRecordstmnt.Exec(req.Clientid, req.Mstorgnhirarchyid, req.Recordid, req.RecordStageid, req.Assetid, req.Assetid)

		if addAssetWithRecordstmnterr != nil {

			logger.Log.Println(addAssetWithRecordstmnterr)

			tx.Rollback()

			return values, errors.New("SQL Execution Error")

		}

		recordCode := ticketID

		var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"

		err := InsertActivityLogs(tx, req.Clientid, req.Mstorgnhirarchyid, req.Recordid, 7, logData, req.Userid, req.GroupID)

		if err != nil {

			log.Println("error is ----->", err)

			tx.Rollback()

			return values, err

		}

		if req.ParentRecordID != 0 {

			var parentRecordStagedID int64

			var parentTicketID string

			getTicketIDQuery := "select code,recordstageid from trnrecord where clientid=? and mstorgnhirarchyid=? and id=?"

			getTicketIDErr := dbc.DB.QueryRow(getTicketIDQuery, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID).Scan(&parentTicketID, &parentRecordStagedID)

			if getTicketIDErr != nil {

				logger.Log.Println("getTicketIDErr in UpdateRecordAsset==> ", getTicketIDErr)

				//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

			}

			checkRecordAssetMapQuery := "select count(id) as count from maprecordasset where clientid=? and mstorgnhirarchyid=? and recordid=? and assetid=? and activeflg=1 and deleteflg=0 "

			checkRecordAssetMapErr := dbc.DB.QueryRow(checkRecordAssetMapQuery, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, req.Assetid).Scan(&count)

			if checkRecordAssetMapErr != nil {

				logger.Log.Println("checkRecordAssetMapErr in UpdateRecordAsset==> ", checkRecordAssetMapErr)

				//           return userTOEmails, userCCEmails, errors.New("ERROR: Unable to fetch EmailAddress of Creator")

				return values, errors.New("SQL Execution Error")

			}

			var addAssetWithRecordquery = "INSERT INTO maprecordasset (clientid, mstorgnhirarchyid, recordid, recordstageid, assetid, trnassetid, deleteflg, activeflg) VALUES(?, ?, ?, ?, ?, ?, 0, 1)"

			addAssetWithRecordstmnt, addAssetWithRecordPreparestmntErr := tx.Prepare(addAssetWithRecordquery)

			if addAssetWithRecordPreparestmntErr != nil {

				logger.Log.Println("addAssetWithRecordPreparestmntErr in UpdateRecordAsset==> ", addAssetWithRecordPreparestmntErr)

				tx.Rollback()

				return values, updateAssetHistoryPreparestmntErr

			}

			defer addAssetWithRecordstmnt.Close()

			_, addAssetWithRecordstmnterr := addAssetWithRecordstmnt.Exec(req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, parentRecordStagedID, req.Assetid, req.Assetid)

			if addAssetWithRecordstmnterr != nil {

				logger.Log.Println(addAssetWithRecordstmnterr)

				tx.Rollback()

				return values, errors.New("SQL Execution Error")

			}

			recordCode := parentTicketID

			var logData = "(Ticket ID:" + recordCode + ", CI ID:" + assetIDS + ")"

			err := InsertActivityLogs(tx, req.Clientid, req.Mstorgnhirarchyid, req.ParentRecordID, 7, logData, req.Userid, req.GroupID)

			if err != nil {

				log.Println("error is ----->", err)

				tx.Rollback()

				return values, err

			}

		}

	}

	return values, nil

}

func (dbc DbConn) AddAssetWithRecord(req *entities.RecordAssetRequestEntity) ([]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go AddAssetWithRecord")
	var query = "INSERT INTO maprecordasset (clientid, mstorgnhirarchyid, recordid, recordstageid, assetid, trnassetid, deleteflg, activeflg) VALUES "
	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, v := range req.AssetID {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, 0, 1)")

		valueArgs = append(valueArgs, req.Clientid)
		valueArgs = append(valueArgs, req.Mstorgnhirarchyid)
		valueArgs = append(valueArgs, req.Recordid)
		valueArgs = append(valueArgs, req.RecordStageid)
		valueArgs = append(valueArgs, v)
		valueArgs = append(valueArgs, v)
	}
	query = query + strings.Join(valueStrings, ",")
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", valueArgs)
	var values []interface{}

	_, err := dbc.DB.Exec(query, valueArgs...)
	if err != nil {
		logger.Log.Println("AddAssetWithRecord Get Statement Prepare Error", err)
		return values, err
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) DeleteAssetFromRecord(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go DeleteAssetFromRecord")
	var query = "UPDATE maprecordasset SET deleteflg=1,activeflg=0 WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND assetid=?"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["recordid"], req["assetid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordid"], req["assetid"])
	if err != nil {
		logger.Log.Println("DeleteAssetFromRecord Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	return values, nil
}
func (dbc DbConn) GetAssetHistroyByAssetID(page *entities.FetchAssetHistoryRequest) (map[string]interface{}, error) {
	var values = make(map[string]interface{})
	var history string
	var query = "select coalesce(history,'') as history from trnasset where clientid=? and mstorgnhirarchyid=? and  id=?"
	err := dbc.DB.QueryRow(query, page.Clientid, page.Mstorgnhirarchyid, page.Assetid).Scan(&history)
	if err != nil {
		logger.Log.Println("GetAssetHistroyByAssetID Scan Error", err)
		return values, err
	}
	values["history"] = history

	return values, nil
}
func (dbc DbConn) InsertActivityLogForRecordAsset(req map[string]interface{}, activitySeq int64) ([]map[string]interface{}, error) {
	logger.Log.Println("In side RecordassetDao.go InsertActivityLogForRecordAsset")
	var query = "INSERT INTO mstrecordactivitylogs(clientid,mstorgnhirarchyid,recordid,activityseqno,logValue,createdid,createddate,createdgrpid) VALUES (?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?)"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["recordid"], req["assetid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordid"], activitySeq, req["logval"], req["userid"], req["groupid"])
	if err != nil {
		logger.Log.Println("InsertActivityLogForRecordAsset Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	return values, nil
}

func (dbc DbConn) GetAssetTypeByRecordID(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side side RecordassetDao.go GetAssetTypeByRecordID")
	var query = "SELECT mstrecorddifferentiationtype.id, mstrecorddifferentiationtype.typename name, mstrecorddifferentiationtype.seqno, difftype.id parent FROM mstrecorddifferentiationtype JOIN mstrecorddifferentiationtype difftype ON mstrecorddifferentiationtype.parentid = difftype.id AND difftype.seqno = 5 WHERE mstrecorddifferentiationtype.activeflg = 1 AND mstrecorddifferentiationtype.deleteflg = 0 AND mstrecorddifferentiationtype.clientid = ? AND mstrecorddifferentiationtype.mstorgnhirarchyid = ? AND mstrecorddifferentiationtype.id IN (SELECT trnasset.mstdifferentiationtypeid FROM trnasset JOIN maprecordasset ON trnasset.id=maprecordasset.assetid WHERE maprecordasset.activeflg=1 AND maprecordasset.deleteflg=0 AND trnasset.activeflg=1 AND trnasset.deleteflg=0 AND maprecordasset.recordid=?) ORDER BY mstrecorddifferentiationtype.seqno ASC"
	logger.Log.Println("Statement :=====>", query)
	logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], req["recordid"])
	var values []map[string]interface{}

	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordid"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAssetTypeByRecordID Get Statement Prepare Error", err)
		return values, err
	}

	for rows.Next() {
		fb := utility.NewDbFieldBind()
		err = fb.PutFields(rows)
		if err != nil {
			return values, err
		}
		rows.Scan(fb.GetFieldPtrArr()...)
		values = append(values, fb.GetFieldArr())
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) GetAssetfieldSpecificDataBYRecordID(page *entities.AssetAttrNameValRequestEntity) ([]entities.AssetAttrNameValEntity, error) {
	logger.Log.Println("In side side RecordassetDao.go GetAssetfieldSpecificDataBYRecordID")
	values := []entities.AssetAttrNameValEntity{}
	if len(page.AssetFieldsNames) > 0 {
		var subquery = ""
		for i := 0; i < len(page.AssetFieldsNames); i++ {
			if subquery != "" {
				subquery = subquery + ","
			}
			subquery = subquery + "'" + page.AssetFieldsNames[i] + "'"
		}
		var query = "SELECT trnasset.id, mstrecorddifferentiation.id attrid, mstrecorddifferentiation.name name, mapassetdifferentiation.value FROM trnasset JOIN maprecordasset ON trnasset.id=maprecordasset.assetid JOIN mapassetdifferentiation ON maprecordasset.assetid=mapassetdifferentiation.trnassetid JOIN mstrecorddifferentiation ON mapassetdifferentiation.mstdifferentiationid=mstrecorddifferentiation.id JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id JOIN mstrecorddifferentiationtype difftype ON mstrecorddifferentiationtype.parentid = difftype.id AND difftype.seqno = 5 WHERE trnasset.activeflg=1 AND trnasset.deleteflg=0 AND maprecordasset.activeflg=1 AND maprecordasset.deleteflg=0 AND mapassetdifferentiation.activeflg=1 AND mapassetdifferentiation.deleteflg=0 AND mstrecorddifferentiationtype.activeflg = 1 AND mstrecorddifferentiationtype.deleteflg = 0 AND trnasset.clientid = ? AND trnasset.mstorgnhirarchyid = ? AND maprecordasset.recordid=? AND LOWER(mstrecorddifferentiation.name) IN (" + subquery + ") ORDER BY trnasset.id ASC, mstrecorddifferentiation.id ASC"
		logger.Log.Println("Statement :=====>", query)
		logger.Log.Println("Params :=====>", page.Clientid, page.Mstorgnhirarchyid, page.Recordid)
		rows, err := dbc.DB.Query(query, page.Clientid, page.Mstorgnhirarchyid, page.Recordid)
		if err != nil {
			logger.Log.Println("GetAssetfieldSpecificDataBYRecordID Get Statement Prepare Error", err)
			return values, err
		}
		defer rows.Close()
		for rows.Next() {
			value := entities.AssetAttrNameValEntity{}
			rows.Scan(&value.Id, &value.AttrID, &value.Name, &value.Value)
			values = append(values, value)
		}
		logger.Log.Println("Results :=====>", values)
	}
	return values, nil
}
