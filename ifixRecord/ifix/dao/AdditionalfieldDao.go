package dao

import (
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"strconv"
)

func (dbc DbConn) GetAdditionalFields(req *entities.AdditionalfieldRequestEntity) ([]entities.AdditionalFieldEntity, error) {
	logger.Log.Println("In side GetAdditionalFields")
	values := []entities.AdditionalFieldEntity{}
	diffLength := len(req.Mstdifferentiationset)
	var sql1 = ""
	for _, v := range req.Mstdifferentiationset {
		if sql1 != "" {
			sql1 = sql1 + " + "
		}
		sql1 = sql1 + "SUM(IF(mstrecordfielddiff.recorddifftypeid = " + strconv.Itoa(int(v.Mstdifferentiationtypeid)) + " AND mstrecordfielddiff.recorddiffid = " + strconv.Itoa(int(v.Mstdifferentiationid)) + " ,1,0)) "
	}
	var sql = "SELECT mstrecordfielddiff.mstrecordfieldid FROM mstrecordfielddiff WHERE mstrecordfielddiff.activeflg = 1 AND mstrecordfielddiff.deleteflg = 0 AND mstrecordfielddiff.clientid = " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfielddiff.mstorgnhirarchyid = " + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " GROUP BY mstrecordfielddiff.mstrecordfieldid HAVING (" + sql1 + ") = " + strconv.Itoa(diffLength)
	var getadditionalsubquery = "SELECT mstrecordfielddiff.mstrecordfieldid FROM mstrecordfielddiff WHERE mstrecordfielddiff.activeflg = 1	AND mstrecordfielddiff.deleteflg = 0 AND mstrecordfielddiff.clientid = " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfielddiff.mstorgnhirarchyid = " + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " AND mstrecordfielddiff.mstrecordfieldid IN (" + sql + ") GROUP BY mstrecordfielddiff.mstrecordfieldid HAVING COUNT(mstrecordfielddiff.id) = " + strconv.Itoa(diffLength)

	//var getadditionafiledsql = "SELECT mstrecordfield.id FieldID,mstrecordterms.id TermsID,mstrecordterms.termname TermsName,mstrecordterms.termvalue TermsValue,msttermtype.id TermsTypeID,msttermtype.termtypename TermsTypeName,mstrecordfield.ismandatory FROM mstrecordfield JOIN mstrecordterms ON mstrecordfield.recordtermid=mstrecordterms.id AND mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE mstrecordfield.activeflg=1 AND mstrecordfield.deleteflg=0 AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD' AND mstrecordfield.clientid= " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfield.mstorgnhirarchyid=" + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " AND mstrecordfield.id IN (" + getadditionalsubquery + ")"
	var getadditionafiledsql = "SELECT mstrecordfield.id FieldID,mstrecordterms.id TermsID,mstrecordterms.seq TermSeqNo, mstrecordterms.termname TermsName,mstrecordterms.termvalue TermsValue,msttermtype.id TermsTypeID,msttermtype.termtypename TermsTypeName,mststateterm.iscompulsory FROM mststateterm,mstrecordfield JOIN mstrecordterms ON mstrecordfield.recordtermid=mstrecordterms.id AND mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE mstrecordfield.activeflg=1 AND mstrecordfield.deleteflg=0 AND mstrecordterms.activeflg=1 AND mststateterm.deleteflg = 0 AND mststateterm.activeflg = 1 AND mstrecordterms.deleteflg=0 AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD' AND mststateterm.recordtermid=mstrecordterms.id AND mstrecordfield.clientid= " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfield.mstorgnhirarchyid=" + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " AND mstrecordfield.id IN (" + getadditionalsubquery + ")"
	logger.Log.Println("Sub Query: ", getadditionalsubquery)
	logger.Log.Println("Main Query: ", getadditionafiledsql)
	rows, err := dbc.DB.Query(getadditionafiledsql)
	if err != nil {
		logger.Log.Println("GetAllRecordTypes Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.AdditionalFieldEntity{}
		rows.Scan(&value.FieldID, &value.TermsID, &value.TermSeqNo, &value.TermsName, &value.TermsValue, &value.TermsTypeID, &value.TermsTypeName, &value.IsMandatory)
		values = append(values, value)
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}

/*func (dbc DbConn) GetAdditionalFields(req *entities.AdditionalfieldRequestEntity) ([]entities.AdditionalFieldEntity, error) {
	logger.Log.Println("In side GetAdditionalFields")
	values := []entities.AdditionalFieldEntity{}
	diffLength := len(req.Mstdifferentiationset)
	var sql1 = ""
	for _, v := range req.Mstdifferentiationset {
		if sql1 != "" {
			sql1 = sql1 + " + "
		}
		sql1 = sql1 + "SUM(IF(mstrecordfielddiff.recorddifftypeid = " + strconv.Itoa(int(v.Mstdifferentiationtypeid)) + " AND mstrecordfielddiff.recorddiffid = " + strconv.Itoa(int(v.Mstdifferentiationid)) + " ,1,0)) "
	}
	var sql = "SELECT mstrecordfielddiff.mstrecordfieldid FROM mstrecordfielddiff WHERE mstrecordfielddiff.activeflg = 1 AND mstrecordfielddiff.deleteflg = 0 AND mstrecordfielddiff.clientid = " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfielddiff.mstorgnhirarchyid = " + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " GROUP BY mstrecordfielddiff.mstrecordfieldid HAVING (" + sql1 + ") = " + strconv.Itoa(diffLength)
	var getadditionalsubquery = "SELECT mstrecordfielddiff.mstrecordfieldid FROM mstrecordfielddiff WHERE mstrecordfielddiff.activeflg = 1	AND mstrecordfielddiff.deleteflg = 0 AND mstrecordfielddiff.clientid = " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfielddiff.mstorgnhirarchyid = " + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " AND mstrecordfielddiff.mstrecordfieldid IN (" + sql + ") GROUP BY mstrecordfielddiff.mstrecordfieldid HAVING COUNT(mstrecordfielddiff.id) = " + strconv.Itoa(diffLength)

	//var getadditionafiledsql = "SELECT mstrecordfield.id FieldID,mstrecordterms.id TermsID,mstrecordterms.termname TermsName,mstrecordterms.termvalue TermsValue,msttermtype.id TermsTypeID,msttermtype.termtypename TermsTypeName,mstrecordfield.ismandatory FROM mstrecordfield JOIN mstrecordterms ON mstrecordfield.recordtermid=mstrecordterms.id AND mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE mstrecordfield.activeflg=1 AND mstrecordfield.deleteflg=0 AND mstrecordterms.activeflg=1 AND mstrecordterms.deleteflg=0 AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD' AND mstrecordfield.clientid= " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfield.mstorgnhirarchyid=" + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " AND mstrecordfield.id IN (" + getadditionalsubquery + ")"
	var getadditionafiledsql = "SELECT mstrecordfield.id FieldID,mstrecordterms.id TermsID,mstrecordterms.termname TermsName,mstrecordterms.termvalue TermsValue,msttermtype.id TermsTypeID,msttermtype.termtypename TermsTypeName,mststateterm.iscompulsory FROM mststateterm,mstrecordfield JOIN mstrecordterms ON mstrecordfield.recordtermid=mstrecordterms.id AND mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE mstrecordfield.activeflg=1 AND mstrecordfield.deleteflg=0 AND mstrecordterms.activeflg=1 AND mststateterm.deleteflg = 0 AND mststateterm.activeflg = 1 AND mstrecordterms.deleteflg=0 AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD' AND mststateterm.recordtermid=mstrecordterms.id AND mstrecordfield.clientid= " + strconv.Itoa(int(req.Clientid)) + " AND mstrecordfield.mstorgnhirarchyid=" + strconv.Itoa(int(req.Mstorgnhirarchyid)) + " AND mstrecordfield.id IN (" + getadditionalsubquery + ")"
	logger.Log.Println("Sub Query: ", getadditionalsubquery)
	logger.Log.Println("Main Query: ", getadditionafiledsql)
	rows, err := dbc.DB.Query(getadditionafiledsql)
	if err != nil {
		logger.Log.Println("GetAllRecordTypes Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.AdditionalFieldEntity{}
		rows.Scan(&value.FieldID, &value.TermsID, &value.TermsName, &value.TermsValue, &value.TermsTypeID, &value.TermsTypeName, &value.IsMandatory)
		values = append(values, value)
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}*/

func (dbc DbConn) GetAdditionalFieldsByDiffId(req *entities.AdditionalfieldRequestEntity) ([]entities.AdditionalFieldEntity, error) {
	logger.Log.Println("In side GetAdditionalFieldsByDiffId")
	values := []entities.AdditionalFieldEntity{}
	var sql = "SELECT distinct mstrecordfield.id FieldID,mstrecordterms.id TermsID,mstrecordterms.termname TermsName,mstrecordterms.termvalue TermsValue,msttermtype.id TermsTypeID,msttermtype.termtypename TermsTypeName,mststateterm.iscompulsory FROM mststateterm,mstrecordfield JOIN mstrecordterms ON mstrecordfield.recordtermid=mstrecordterms.id AND mstrecordfield.clientid=mstrecordterms.clientid AND mstrecordfield.mstorgnhirarchyid=mstrecordterms.mstorgnhirarchyid JOIN msttermtype ON mstrecordterms.termtypeid=msttermtype.id WHERE mstrecordfield.activeflg=1 AND mstrecordfield.deleteflg=0 AND mstrecordterms.activeflg=1 AND mststateterm.deleteflg = 0 AND mststateterm.activeflg = 1 AND mstrecordterms.deleteflg=0 AND mstrecordfield.mstrecordfieldtype='ADDITIONAL FIELD' AND mststateterm.recordtermid=mstrecordterms.id AND mstrecordfield.clientid= '" + strconv.Itoa(int(req.Clientid)) + "' AND mstrecordfield.mstorgnhirarchyid='" + strconv.Itoa(int(req.Mstorgnhirarchyid)) + "' AND mstrecordfield.id IN (SELECT bb.mstrecordfieldid FROM (SELECT mstrecordfieldid,GROUP_CONCAT(CONCAT(recorddifftypeid,'_',recorddiffid)) aa FROM iFIX.mstrecordfielddiff WHERE deleteflg=0 AND activeflg=1 AND clientid='" + strconv.Itoa(int(req.Clientid)) + "' AND mstorgnhirarchyid='" + strconv.Itoa(int(req.Mstorgnhirarchyid)) + "' GROUP BY mstrecordfieldid) bb WHERE "

	var sql1 = ""
	for _, v := range req.Mstdifferentiationset {
		if sql1 != "" {
			sql1 = sql1 + " AND "
		}
		sql1 = sql1 + " bb.aa LIKE '%" + strconv.Itoa(int(v.Mstdifferentiationtypeid)) + "_" + strconv.Itoa(int(v.Mstdifferentiationid)) + "%' "
	}
	sql = sql + sql1 + ")"
	logger.Log.Println("Main Query: ", sql)
	rows, err := dbc.DB.Query(sql)
	if err != nil {
		logger.Log.Println("GetAllRecordTypes Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.AdditionalFieldEntity{}
		rows.Scan(&value.FieldID, &value.TermsID, &value.TermsName, &value.TermsValue, &value.TermsTypeID, &value.TermsTypeName, &value.IsMandatory)
		values = append(values, value)
	}
	logger.Log.Println("Results: ", values)
	return values, nil
}
