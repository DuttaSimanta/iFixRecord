package dao

import (
	"database/sql"
	"encoding/json"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/utility"
	"log"
	"strings"
	"time"
)

var getmstdashdtlqsql = "SELECT id,query,COALESCE(queryparam,''),joinquery,querytype FROM mstdashboarddtls WHERE activeflg=1 AND deleteflg=0 AND clientid=? AND mstorgnhirarchyid=? AND mstrecorddifferentiationid=? AND mapfunctionalityid=? AND querytype=? "

func (dbc DbConn) GetQueryNParams(req map[string]interface{}) ([]entities.DynamicqueryEntity, error) {
	logger.Log.Println("In side GetQueryNParams")

	var recorddiffid int64
	if int64(req["iscatalog"].(float64)) == 1 {
		recorddiffid = 0
		//req["recorddiffid"] = 0
	} else if int64(req["iscatalog"].(float64)) == 0 {
		recorddiffid = int64(req["recorddiffid"].(float64))
	}
	//logger.Log.Println("Statement :=====>", getmstdashdtlqsql)
	//logger.Log.Println("Params :=====>", req["clientid"], req["mstorgnhirarchyid"], recorddiffid, req["menuid"], req["querytype"])
	values := []entities.DynamicqueryEntity{}

	rows, err := dbc.DB.Query(getmstdashdtlqsql, req["clientid"], req["mstorgnhirarchyid"], recorddiffid, req["menuid"], req["querytype"])

	if err != nil {
		logger.Log.Println("GetQueryNParams Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.DynamicqueryEntity{}
		rows.Scan(&value.ID, &value.Query, &value.QueryParam, &value.JoinQuery, &value.QueryType)
		values = append(values, value)
	}
	//logger.Log.Println("Results: ", values)
	return values, nil
}

func (dbc DbConn) GetQueryResuls(req map[string]interface{}, qry entities.DynamicqueryEntity) ([]interface{}, error) {
	logger.Log.Println("In side GetQueryResuls")

	var values = []interface{}{}
	pint := []interface{}{}
	if qry.QueryParam != "" {
		params := strings.Split(qry.QueryParam, ",")

		for _, v := range params {
			pint = append(pint, req[v])
		}

	}
	//logger.Log.Println("Query:===>", qry.Query)
	//logger.Log.Println("Params:===>", pint)
	rows, err := dbc.DB.Query(qry.Query, pint...)
	if err != nil {
		logger.Log.Println("GetQueryResuls Get Statement Prepare Error", err)
		return values, err
	}
	var fArr []string
	//var ftArr []*sql.ColumnType
	fb := utility.NewFieldBinding()

	if fArr, err = rows.Columns(); err != nil {
		return values, err
	}
	/*if ftArr, err = rows.ColumnTypes(); err != nil {
		return values, err
	}*/
	/*for i, v := range ftArr {
		tmp, _ := v.ScanType().FieldByName(fArr[i])
		if tmp.Type==reflect.String{

		}
		log.Println(tmp)
	}*/

	fb.PutFields(fArr)

	defer rows.Close()

	for rows.Next() {
		rows.Scan(fb.GetFieldPtrArr()...)
		values = append(values, fb.GetFieldArr())

	}
	return values, nil
}

func (dbc DbConn) GetCategoryNameByRecordID(req map[string]interface{}) (map[int64]interface{}, error) {
	logger.Log.Println("In side DynamicqueryDao.go")
	values := make(map[int64]interface{})
	var query = "SELECT mstrecorddifferentiation.name AS recorddiffname, mstrecorddifferentiationtype.typename AS recorddifftypename,mstrecorddifferentiationtype.id recorddifftypeid FROM  maprecordtorecorddifferentiation  JOIN mstrecorddifferentiation  ON  maprecordtorecorddifferentiation.recorddiffid  = mstrecorddifferentiation.id  AND  maprecordtorecorddifferentiation.recorddifftypeid  = mstrecorddifferentiation.recorddifftypeid  AND  maprecordtorecorddifferentiation.clientid  = mstrecorddifferentiation.clientid  AND  maprecordtorecorddifferentiation.mstorgnhirarchyid  = mstrecorddifferentiation.mstorgnhirarchyid  JOIN mstrecorddifferentiationtype  ON  mstrecorddifferentiation.recorddifftypeid  = mstrecorddifferentiationtype.id  WHERE  maprecordtorecorddifferentiation.activeflg  = 1 AND  maprecordtorecorddifferentiation.deleteflg  = 0 AND  mstrecorddifferentiation.activeflg  = 1 AND  mstrecorddifferentiation.deleteflg  = 0 AND  mstrecorddifferentiationtype.parentid = 1 AND maprecordtorecorddifferentiation.clientid = ? AND maprecordtorecorddifferentiation.mstorgnhirarchyid = ? AND  maprecordtorecorddifferentiation.recordid = ? AND maprecordtorecorddifferentiation.recordstageid = ?"
	//logger.Log.Println("Query:===>", query)
	//logger.Log.Println("Params:===>", req["clientid"], req["mstorgnhirarchyid"], req["recordid"], req["recordstageid"])
	rows, err := dbc.DB.Query(query, req["clientid"], req["mstorgnhirarchyid"], req["recordid"], req["recordstageid"])
	if err != nil {
		logger.Log.Println("GetCategoryNameByRecordID Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	var typename string
	var name string
	var typeid int64
	for rows.Next() {
		rows.Scan(&name, &typename, &typeid)
		values[typeid] = name
	}
	return values, nil
}
func ConvertDateByDiff(date int64, timediff int64) string {
	var unixTime = date + timediff
	t := time.Unix(unixTime, 0)
	return t.Format("02-Jan-2006 15:04:05")
}
func (dbc DbConn) GetQueryResulsTemp(req map[string]interface{}, qry entities.DynamicqueryEntity, timediff int64) ([]entities.QueryResultEntity, error) {
	logger.Log.Println("In side GetQueryResuls")

	values := []entities.QueryResultEntity{}
	pint := []interface{}{}
	if qry.QueryParam != "" {
		params := strings.Split(qry.QueryParam, ",")

		for _, v := range params {
			pint = append(pint, req[v])
		}

	}
	//logger.Log.Println("Query:===>", qry.Query)
	//logger.Log.Println("Params:===>", pint)
	rows, err := dbc.DB.Query(qry.Query, pint...)
	if err != nil {
		logger.Log.Println("GetQueryResuls Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.QueryResultEntity{}
		var datewithdiff int64
		rows.Scan(&value.ID, &value.Recordtitle, &value.Code, &value.Recorddescription, &value.Requesterinfo, &datewithdiff, &value.Createdby, &value.Supportgroupname, &value.Levelname, &value.Recordtype, &value.Status, &value.Impact, &value.Urgency, &value.Priority, &value.AssignedGroup, &value.AssignedGroupLevel, &value.Assignee, &value.StageID, &value.PriorityChangeCount, &value.FollowUpCount, &value.OutBoundCount, &value.ReOpenCount, &value.ResolutionViolated, &value.Aging)
		value.Createdatetime = ConvertDateByDiff(datewithdiff, timediff)
		values = append(values, value)

	}
	return values, nil
}

func (dbc DbConn) GetQueryResulsCountTemp(req map[string]interface{}, qry entities.DynamicqueryEntity) (entities.QueryResultEntity, error) {
	logger.Log.Println("In side GetQueryResuls")

	values := entities.QueryResultEntity{}
	pint := []interface{}{}
	if qry.QueryParam != "" {
		params := strings.Split(qry.QueryParam, ",")

		for _, v := range params {
			pint = append(pint, req[v])
		}

	}
	//logger.Log.Println("Query:===>", qry.Query)
	//logger.Log.Println("Params:===>", pint)
	rows, err := dbc.DB.Query(qry.Query, pint...)
	if err != nil {
		logger.Log.Println("GetQueryResuls Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&values.ID)
	}
	return values, nil
}

func (dbc DbConn) GetCategory(req map[string]interface{}) ([]entities.DynamicqueryEntity, error) {
	logger.Log.Println("In side GetQueryNParams")
	//logger.Log.Println("Statement :=====>", getmstdashdtlqsql)
	//logger.Log.Println("Params :=====>", req)
	values := []entities.DynamicqueryEntity{}
	rows, err := dbc.DB.Query(getmstdashdtlqsql, req["clientid"], req["mstorgnhirarchyid"], req["recorddiffid"], req["menuid"], req["querytype"])
	if err != nil {
		logger.Log.Println("GetQueryNParams Get Statement Prepare Error", err)
		return values, err
	}
	defer rows.Close()
	for rows.Next() {
		value := entities.DynamicqueryEntity{}
		rows.Scan(&value.ID, &value.Query, &value.QueryParam)
		values = append(values, value)
	}
	//logger.Log.Println("Results: ", values)
	return values, nil
}

func (mdao DbConn) Gettimezonediff(Clientid interface{}, Mstorgnhirarchyid interface{}) (int64, int64, error) {
	var timediff int64
	var reportTimediff int64
	var gettimediff = "select a.utcdiff as timediff,c.utcdiff as reporttimediff from zone a,mstorgnhierarchy b,zone c where a.zone_id=b.timezoneid and b.reporttimezoneid=c.zone_id and b.clientid=? and b.id=? "
	stmt, err := mdao.DB.Prepare(gettimediff)
	if err != nil {
		logger.Log.Print("Gettimediff Statement Prepare Error", err)
		log.Print("Gettimediff Statement Prepare Error", err)
		return timediff, reportTimediff, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(Clientid, Mstorgnhirarchyid)
	if err != nil {
		logger.Log.Print("Gettimediff Statement Execution Error", err)
		log.Print("Gettimediff Statement Execution Error", err)
		return timediff, reportTimediff, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&timediff, &reportTimediff)
	}
	return timediff, reportTimediff, nil
}

func (dbc DbConn) RecordGridResult(req map[string]interface{}, qry entities.DynamicqueryEntity, getOrgnID []int64) ([]map[string]interface{}, error) {
	logger.Log.Println("In side RecordGridResult")
	values := []map[string]interface{}{}
	pint := []interface{}{}
	var queryStr = "SELECT "

	if req["headers"] != nil {
		for i, v := range req["headers"].([]interface{}) {
			if i > 0 {
				queryStr = queryStr + " , "
			}
			queryStr = queryStr + " recordfulldetails." + v.(string) + " "
		}
		if len(req["headers"].([]interface{})) == 0 {
			queryStr = queryStr + " recordfulldetails.* "
		}
	} else {
		queryStr = queryStr + " recordfulldetails.* "
	}
	queryStr = queryStr + " FROM recordfulldetails "
	if qry.JoinQuery != "" {
		queryStr = queryStr + " " + qry.JoinQuery + " "
	}
	var ids string = ""
	// for i, state := range getOrgnID {
	// 	if i > 0 {
	// 		ids += ","
	// 	}
	// 	ids += strconv.Itoa(int(state))
	// }
	ids = req["searchmstorgnhirarchyid"].(string)

	if int64(req["iscatalog"].(float64)) == 1 {
		//queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (?) AND recordfulldetails.tickettypeid=? "
		queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (" + ids + ") AND recordfulldetails.tickettypeid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid in (" + ids + ") AND recorddifftypeid=2 AND seqno in(1,2) AND deleteflg=0 AND activeflg=1) AND recordfulldetails.statusid in (SELECT id FROM mstrecorddifferentiation WHERE clientid = ? AND mstorgnhirarchyid IN (" + ids + ") AND recorddifftypeid=3 "
		pint = append(pint, req["clientid"])
		//pint = append(pint, req["mstorgnhirarchyid"])
		//ids := strings.Join(getOrgnID, "','")

		//pint = append(pint, ids)
		//pint = append(pint, req["recorddiffid"])
		pint = append(pint, req["clientid"])
		//pint = append(pint, ids)
		//pint = append(pint, req["recorddiffidseq"])
		pint = append(pint, req["clientid"])

		//req["recorddiffid"] = 0
	} else if int64(req["iscatalog"].(float64)) == 0 {
		//queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (?) AND recordfulldetails.tickettypeid=? "
		queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (" + ids + ") AND recordfulldetails.tickettypeid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid in (" + ids + ") AND recorddifftypeid=2 AND seqno=? AND deleteflg=0 AND activeflg=1) AND recordfulldetails.statusid in (SELECT id FROM mstrecorddifferentiation WHERE clientid = ? AND mstorgnhirarchyid IN (" + ids + ") AND recorddifftypeid=3 "
		pint = append(pint, req["clientid"])
		//pint = append(pint, req["mstorgnhirarchyid"])
		//ids := strings.Join(getOrgnID, "','")

		//pint = append(pint, ids)
		//pint = append(pint, req["recorddiffid"])
		pint = append(pint, req["clientid"])
		//pint = append(pint, ids)
		pint = append(pint, req["recorddiffidseq"])
		pint = append(pint, req["clientid"])
	}

	/*//queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (?) AND recordfulldetails.tickettypeid=? "
	queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (" + ids + ") AND recordfulldetails.tickettypeid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid in (" + ids + ") AND recorddifftypeid=2 AND seqno=? AND deleteflg=0 AND activeflg=1) AND recordfulldetails.statusid in (SELECT id FROM mstrecorddifferentiation WHERE clientid = ? AND mstorgnhirarchyid IN (" + ids + ") AND recorddifftypeid=3 "
	pint = append(pint, req["clientid"])
	//pint = append(pint, req["mstorgnhirarchyid"])
	//ids := strings.Join(getOrgnID, "','")

	//pint = append(pint, ids)
	//pint = append(pint, req["recorddiffid"])
	pint = append(pint, req["clientid"])
	//pint = append(pint, ids)
	pint = append(pint, req["recorddiffidseq"])
	pint = append(pint, req["clientid"])*/

	var paramCounter int64
	var supportgrpids string

	if qry.QueryParam != "" {
		params := strings.Split(qry.QueryParam, ",")
		for _, v := range params {
			if v == "clientid" || v == "mstorgnhirarchyid" || v == "recorddiffid" || v == "limit" || v == "offset" || v == "where" || v == "order" {
				continue
			} else if v == "supportgrpid" {
				supportgrpids = req["supportgrpid"].(string)
				qry.Query = strings.Replace(qry.Query, "recordfulldetails.assignedgroupid=?", " recordfulldetails.assignedgroupid IN ("+supportgrpids+") ", -1)
				//logger.Log.Print("Query =============================================>", qry.Query)
				qry.Query = strings.Replace(qry.Query, "mstrequesthistory.mstgroupid=?", " mstrequesthistory.mstgroupid IN ("+supportgrpids+") ", -1)
				logger.Log.Print("Query =======================3333333333333333333333333333333======================>", qry.Query)
				qry.Query = strings.Replace(qry.Query, "recordfulldetails.resogroupid=?", " recordfulldetails.resogroupid IN ("+supportgrpids+") ", -1)
				//continue
			} else {
				paramCounter = paramCounter + 1
				pint = append(pint, req[v])
			}

		}
		// if paramCounter > 0 {
		// 	queryStr = queryStr + " AND " + qry.Query
		// }

	}
	queryStr = queryStr + " AND " + qry.Query

	if req["where"] != nil {
		whereMap := req["where"].([]interface{})
		qstr, prm := GenerateCustomWhereQuery(whereMap)
		queryStr = queryStr + " " + qstr
		pint = append(pint, prm...)
	}
	queryStr = queryStr + " AND recordfulldetails.statusid = maprecordstatetodifferentiation.recorddiffid GROUP BY recordfulldetails.id  "
	if req["order"] != nil {
		orderMap := req["order"].([]interface{})
		qstr := GenerateCustomOrderQuery(orderMap)
		if len(qstr) > 0 {
			queryStr = queryStr + " " + qstr
		} else {
			queryStr = queryStr + " " + "ORDER BY recordfulldetails.lastupdateddatetime desc"
		}
	}

	if req["limit"] != nil && req["offset"] != nil {
		queryStr = queryStr + " LIMIT ?,?"
		pint = append(pint, req["offset"])
		pint = append(pint, req["limit"])
	}
logger.Log.Println("Query:===>", queryStr)
logger.Log.Println("Params:===>", pint)
	rows, err := dbc.DB.Query(queryStr, pint...)
	if err != nil {
		logger.Log.Println("GetQueryResuls Get Statement Prepare Error", err)
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
		//logger.Log.Println("values", fb.GetFieldArr())
		// values = append(values, fb.GetFieldArr())
		res := fb.GetFieldArr()
		record := res["recordid"].(*utility.NullInt64)
		recordid, _ := record.Value()
		tickettypeid := res["tickettypeid"].(*utility.NullInt64)
		diffid, _ := tickettypeid.Value()
		client := res["clientid"].(*utility.NullInt64)
		clientid, _ := client.Value()
		org := res["mstorgnhirarchyid"].(*utility.NullInt64)
		orgid, _ := org.Value()
		categories, err := GetCategoryNames(dbc, clientid.(int64), orgid.(int64), 2, diffid.(int64), recordid.(int64))
		if err != nil {
			return values, err
		}
		res["categories"] = categories
		values = append(values, res)
	}
	return values, nil
}

func (dbc DbConn) RecordGridResultOnly(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side RecordGridResultOnly")
	values := []map[string]interface{}{}
	pint := []interface{}{}

	var queryStr = "SELECT "

	if req["headers"] != nil {
		for i, v := range req["headers"].([]interface{}) {
			if i > 0 {
				queryStr = queryStr + " , "
			}
			queryStr = queryStr + " recordfulldetails." + v.(string) + " "
		}
		if len(req["headers"].([]interface{})) == 0 {
			queryStr = queryStr + " recordfulldetails.* "
		}
	} else {
		queryStr = queryStr + " recordfulldetails.* "
	}
	queryStr = queryStr + " FROM recordfulldetails  WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 "

	if req["where"] != nil {
		whereMap := req["where"].([]interface{})
		qstr, prm := GenerateCustomWhereQuery(whereMap)
		queryStr = queryStr + " " + qstr
		pint = append(pint, prm...)
	}

	if req["order"] != nil {
		orderMap := req["order"].([]interface{})
		qstr := GenerateCustomOrderQuery(orderMap)
		queryStr = queryStr + " " + qstr
	}
	if req["limit"] != nil && req["offset"] != nil {
		queryStr = queryStr + " LIMIT ?,?"
		pint = append(pint, req["offset"])
		pint = append(pint, req["limit"])
	}
	//logger.Log.Println("Query:===>", queryStr)
	//logger.Log.Println("Params:===>", pint)
	rows, err := dbc.DB.Query(queryStr, pint...)
	if err != nil {
		logger.Log.Println("GetQueryResuls Get Statement Prepare Error", err)
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
		res := fb.GetFieldArr()
		record := res["recordid"].(*utility.NullInt64)
		recordid, _ := record.Value()
		tickettypeid := res["tickettypeid"].(*utility.NullInt64)
		diffid, _ := tickettypeid.Value()
		client := res["clientid"].(*utility.NullInt64)
		clientid, _ := client.Value()
		org := res["mstorgnhirarchyid"].(*utility.NullInt64)
		orgid, _ := org.Value()
		categories, err := GetCategoryNames(dbc, clientid.(int64), orgid.(int64), 2, diffid.(int64), recordid.(int64))
		if err != nil {
			return values, err
		}
		res["categories"] = categories
		values = append(values, res)
	}
	return values, nil
}
func GetCategoryNames(mdao DbConn, ClientID int64, Mstorgnhirarchyid int64, RecorddifftypeID int64, RecorddifID int64, RecordID int64) ([]map[string]interface{}, error) {
	var sql = "SELECT distinct (SELECT typename from mstrecorddifferentiationtype WHERE id=a.torecorddifftypeid) lable,d.name,a.torecorddifftypeid as lebelid FROM mstrecordtype a, mstrecorddifferentiationtype b,maprecordtorecorddifferentiation c,mstrecorddifferentiation d WHERE a.torecorddifftypeid = b.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.fromrecorddifftypeid=? AND a.fromrecorddiffid=? AND b.parentid=1 AND c.recordid=? AND c.islatest=1 AND a.torecorddifftypeid=c.recorddifftypeid and c.recorddiffid=d.id"
	values := []map[string]interface{}{}
	logger.Log.Println(sql)
	logger.Log.Println("PARAMETER", ClientID, Mstorgnhirarchyid, RecorddifftypeID, RecorddifID, RecordID)
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecorddifftypeID, RecorddifID, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLaststagevalue Get Statement Prepare Error", err)
		return values, err
	}
	// logger.Log.Println(rows)
	for rows.Next() {
		fb := utility.NewDbFieldBind()
		err = fb.PutFields(rows)
		if err != nil {
			return values, err
		}
		rows.Scan(fb.GetFieldPtrArr()...)
		res := fb.GetFieldArr()
		values = append(values, res)

	}
	return values, nil
}
func GenerateCustomWhereOrQuery(req []interface{}) (string, []interface{}) {
	var queryStr string
	pint := []interface{}{}
	queryStr = ""
	for i, v := range req {
		c := v.(map[string]interface{})
		x := c["op"].(string)
		if i > 0 {
			queryStr = queryStr + " OR "
		}
		switch x {
		case ">", "<", ">=", "<=", "!=", "=":
			pint = append(pint, c["val"].(string))
			queryStr = queryStr + " recordfulldetails." + c["field"].(string) + " " + c["op"].(string) + " " + "? "
		case "like":
			queryStr = queryStr + " recordfulldetails." + c["field"].(string) + " LIKE '%" + c["val"].(string) + "%' "
		case "in":
			rawval := c["val"].(string)
			var splitval = strings.Split(rawval, ",")
			var joinedval = " ('" + strings.Join(splitval, "','") + "') "
			queryStr = queryStr + " recordfulldetails." + c["field"].(string) + " IN " + joinedval
		case "notin":
			rawval := c["val"].(string)
			var splitval = strings.Split(rawval, ",")
			var joinedval = " ('" + strings.Join(splitval, "','") + "') "
			queryStr = queryStr + " recordfulldetails." + c["field"].(string) + " NOT IN " + joinedval
		case "between":
			rawval := c["val"].(string)
			var splitval = strings.Split(rawval, ",")
			if len(splitval) == 2 {
				queryStr = queryStr + " recordfulldetails." + c["field"].(string) + " BETWEEN  ? AND ? "
				pint = append(pint, splitval[0])
				pint = append(pint, splitval[1])
			}
		}

	}

	return queryStr, pint
}

func GenerateCustomWhereQuery(req []interface{}) (string, []interface{}) {
	var queryStr string
	pint := []interface{}{}
	queryStr = ""
	for _, v := range req {
		c := v.(map[string]interface{})
		x := c["op"].(string)
		switch x {
		case ">", "<", ">=", "<=", "!=", "=":
			pint = append(pint, c["val"].(string))
			queryStr = queryStr + " AND recordfulldetails." + c["field"].(string) + " " + c["op"].(string) + " " + "? "
		case "like":
			queryStr = queryStr + " AND recordfulldetails." + c["field"].(string) + " LIKE '%" + c["val"].(string) + "%' "
		case "in":
			rawval := c["val"].(string)
			var splitval = strings.Split(rawval, ",")
			var joinedval = " ('" + strings.Join(splitval, "','") + "') "
			queryStr = queryStr + " AND recordfulldetails." + c["field"].(string) + " IN " + joinedval
		case "notin":
			rawval := c["val"].(string)
			var splitval = strings.Split(rawval, ",")
			var joinedval = " ('" + strings.Join(splitval, "','") + "') "
			queryStr = queryStr + " AND recordfulldetails." + c["field"].(string) + " NOT IN " + joinedval
		case "between":
			rawval := c["val"].(string)
			var splitval = strings.Split(rawval, ",")
			if len(splitval) == 2 {
				queryStr = queryStr + " AND recordfulldetails." + c["field"].(string) + " BETWEEN  ? AND ? "
				pint = append(pint, splitval[0])
				pint = append(pint, splitval[1])
			}
		case "or":
			rawval := c["val"].([]interface{})
			resQuery, resParams := GenerateCustomWhereOrQuery(rawval)
			pint = append(pint, resParams...)
			queryStr = queryStr + " AND ( " + resQuery + " ) "
		}

	}

	return queryStr, pint
}

func GenerateCustomOrderQuery(req []interface{}) string {
	var queryStr string
	for i, v := range req {
		if i == 0 {
			queryStr = queryStr + " ORDER BY "
		} else {
			queryStr = queryStr + " , "
		}
		c := v.(map[string]interface{})
		queryStr = queryStr + " recordfulldetails." + c["field"].(string) + " " + c["dir"].(string)
	}
	logger.Log.Println(" WHERE Query:===>", queryStr)
	return queryStr
}

func (dbc DbConn) RecordGridCount(req map[string]interface{}, qry entities.DynamicqueryEntity, getOrgnID []int64) (int64, error) {
	logger.Log.Println("In side RecordGridCount")
	var values int64
	pint := []interface{}{}
	var queryStr = "SELECT COUNT(DISTINCT recordfulldetails.recordid) total FROM recordfulldetails "
	if qry.JoinQuery != "" {
		queryStr = queryStr + " " + qry.JoinQuery + " "
	}
	var ids string = ""
	// for i, state := range getOrgnID {
	// 	if i > 0 {
	// 		ids += ","
	// 	}
	// 	ids += strconv.Itoa(int(state))
	// }

	ids = req["searchmstorgnhirarchyid"].(string)

	if int64(req["iscatalog"].(float64)) == 1 {
		//queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (?) AND recordfulldetails.tickettypeid=? "
		queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (" + ids + ") AND recordfulldetails.tickettypeid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid in (" + ids + ") AND recorddifftypeid=2 AND seqno in (1,2) AND deleteflg=0 AND activeflg=1) AND recordfulldetails.statusid in (SELECT id FROM mstrecorddifferentiation WHERE clientid = ? AND mstorgnhirarchyid IN (" + ids + ") AND recorddifftypeid=3 "

		//ids := strings.Join(getOrgnID, "','")

		pint = append(pint, req["clientid"])
		//pint = append(pint, req["mstorgnhirarchyid"])
		//pint = append(pint, ids)
		//pint = append(pint, req["recorddiffid"])
		pint = append(pint, req["clientid"])
		//pint = append(pint, ids)
		//pint = append(pint, req["recorddiffidseq"])
		pint = append(pint, req["clientid"])
		//req["recorddiffid"] = 0
	} else if int64(req["iscatalog"].(float64)) == 0 {
		//queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (?) AND recordfulldetails.tickettypeid=? "
		queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (" + ids + ") AND recordfulldetails.tickettypeid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid in (" + ids + ") AND recorddifftypeid=2 AND seqno=? AND deleteflg=0 AND activeflg=1) AND recordfulldetails.statusid in (SELECT id FROM mstrecorddifferentiation WHERE clientid = ? AND mstorgnhirarchyid IN (" + ids + ") AND recorddifftypeid=3 "

		//ids := strings.Join(getOrgnID, "','")

		pint = append(pint, req["clientid"])
		//pint = append(pint, req["mstorgnhirarchyid"])
		//pint = append(pint, ids)
		//pint = append(pint, req["recorddiffid"])
		pint = append(pint, req["clientid"])
		//pint = append(pint, ids)
		pint = append(pint, req["recorddiffidseq"])
		pint = append(pint, req["clientid"])
	}

	/*//queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (?) AND recordfulldetails.tickettypeid=? "
	queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 AND recordfulldetails.clientid=? AND recordfulldetails.mstorgnhirarchyid in (" + ids + ") AND recordfulldetails.tickettypeid in (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid in (" + ids + ") AND recorddifftypeid=2 AND seqno=? AND deleteflg=0 AND activeflg=1) AND recordfulldetails.statusid in (SELECT id FROM mstrecorddifferentiation WHERE clientid = ? AND mstorgnhirarchyid IN (" + ids + ") AND recorddifftypeid=3 "

	//ids := strings.Join(getOrgnID, "','")

	pint = append(pint, req["clientid"])
	//pint = append(pint, req["mstorgnhirarchyid"])
	//pint = append(pint, ids)
	//pint = append(pint, req["recorddiffid"])
	pint = append(pint, req["clientid"])
	//pint = append(pint, ids)
	pint = append(pint, req["recorddiffidseq"])
	pint = append(pint, req["clientid"])*/
	var paramCounter int64
	var supportgrpids string
	//logger.Log.Println("Query parameter :=111111111111111111111111111111111111111111==>", qry.QueryParam)
	//logger.Log.Println("Query parameter :=111111111111111111111111111111111111111111==>", qry.Query)
	if qry.QueryParam != "" {
		params := strings.Split(qry.QueryParam, ",")

		for _, v := range params {
			if v == "clientid" || v == "mstorgnhirarchyid" || v == "recorddiffid" || v == "limit" || v == "offset" || v == "where" || v == "order" {
				continue
			} else if v == "supportgrpid" {
				supportgrpids = req["supportgrpid"].(string)
				//logger.Log.Println("queryStr :=111111111111111111111111111111111111111111==>", queryStr)
				qry.Query = strings.Replace(qry.Query, "recordfulldetails.assignedgroupid=?", " recordfulldetails.assignedgroupid IN ("+supportgrpids+") ", -1)
				qry.Query = strings.Replace(qry.Query, "mstrequesthistory.mstgroupid=?", " mstrequesthistory.mstgroupid IN ("+supportgrpids+") ", -1)
				qry.Query = strings.Replace(qry.Query, "recordfulldetails.resogroupid=?", " recordfulldetails.resogroupid IN ("+supportgrpids+") ", -1)
				//continue
			} else {
				paramCounter = paramCounter + 1
				pint = append(pint, req[v])
			}

		}
		// if paramCounter > 0 {
		// 	queryStr = queryStr + " AND " + qry.Query
		// }

	}
	queryStr = queryStr + " AND " + qry.Query

	if req["where"] != nil {
		whereMap := req["where"].([]interface{})
		qstr, prm := GenerateCustomWhereQuery(whereMap)
		queryStr = queryStr + " " + qstr
		pint = append(pint, prm...)
	}

	queryStr = queryStr + " AND recordfulldetails.statusid = maprecordstatetodifferentiation.recorddiffid "
	logger.Log.Println("Query:===>", queryStr)
	logger.Log.Println("Params:===>", pint)

	err := dbc.DB.QueryRow(queryStr, pint...).Scan(&values)
	switch err {
	case sql.ErrNoRows:
		values = 0
		return values, nil
	case nil:
		return values, nil
	default:
		logger.Log.Println("RecordGridCount Get Statement Prepare Error", err)
		return values, err

	}
}

func (dbc DbConn) RecordGridCountOnly(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side RecordGridCountOnly")
	var values int64
	pint := []interface{}{}
	var queryStr = "SELECT COUNT(DISTINCT recordfulldetails.recordid) total FROM recordfulldetails "
	queryStr = queryStr + " WHERE recordfulldetails.activeflg=1 AND recordfulldetails.deleteflg=0 "

	if req["where"] != nil {
		whereMap := req["where"].([]interface{})
		qstr, prm := GenerateCustomWhereQuery(whereMap)
		queryStr = queryStr + " " + qstr
		pint = append(pint, prm...)
	}

	//logger.Log.Println("Query:===>", queryStr)
	//	logger.Log.Println("Params:===>", pint)

	err := dbc.DB.QueryRow(queryStr, pint...).Scan(&values)
	switch err {
	case sql.ErrNoRows:
		values = 0
		return values, nil
	case nil:
		return values, nil
	default:
		logger.Log.Println("RecordGridCountOnly Get Statement Prepare Error", err)
		return values, err

	}
}

func (dbc DbConn) RecordFilterNameCheck(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side RecordFilterNameCheck")
	var values int64
	pint := []interface{}{}
	var queryStr = "SELECT COUNT(DISTINCT recordsearchfilter.id) total FROM recordsearchfilter "
	queryStr = queryStr + " WHERE recordsearchfilter.activeflg=1 AND recordsearchfilter.deleteflg=0 AND recordsearchfilter.userid=? AND recordsearchfilter.name = ? "

	pint = append(pint, req["userid"])
	pint = append(pint, req["name"])

	//logger.Log.Println("Query:===>", queryStr)
	//logger.Log.Println("Params:===>", pint)

	err := dbc.DB.QueryRow(queryStr, pint...).Scan(&values)
	switch err {
	case sql.ErrNoRows:
		values = 0
		return values, nil
	case nil:
		return values, nil
	default:
		logger.Log.Println("RecordGridCountOnly Get Statement Prepare Error", err)
		return values, err

	}
}

func (dbc DbConn) RecordFilterAdd(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side RecordFilterAdd")
	var values int64
	pint := []interface{}{}
	var queryStr = "INSERT INTO recordsearchfilter(name, userid,filter,savedfilters,activeflg,deleteflg,createdate) VALUES (?,?,?,?,1,0,NOW())"

	pint = append(pint, req["name"])
	pint = append(pint, req["userid"])
	jsonData, err := json.Marshal(req["filter"])
	if err != nil {
		logger.Log.Println("RecordFilterAdd Get Statement Prepare Error", err)
		return values, err
	}
	pint = append(pint, string(jsonData))
	pint = append(pint, req["savedfilters"])

	//logger.Log.Println("Query:===>", queryStr)
	//logger.Log.Println("Params:===>", pint)

	res, err := dbc.DB.Exec(queryStr, pint...)
	if err != nil {
		logger.Log.Println("RecordFilterAdd Get Statement Prepare Error", err)
		return values, err
	}
	values, err = res.LastInsertId()
	if err != nil {
		logger.Log.Println("RecordFilterAdd Get Statement Prepare Error", err)
		return values, err
	}
	return values, err
}

func (dbc DbConn) RecordFilterCount(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side RecordFilterCount")
	var values int64
	pint := []interface{}{}
	var queryStr = "SELECT COUNT(DISTINCT recordsearchfilter.id) total FROM recordsearchfilter "
	queryStr = queryStr + " WHERE recordsearchfilter.activeflg=1 AND recordsearchfilter.deleteflg=0 AND recordsearchfilter.userid=? "

	pint = append(pint, req["userid"])

	//logger.Log.Println("Query:===>", queryStr)
	//logger.Log.Println("Params:===>", pint)

	err := dbc.DB.QueryRow(queryStr, pint...).Scan(&values)
	switch err {
	case sql.ErrNoRows:
		values = 0
		return values, nil
	case nil:
		return values, nil
	default:
		logger.Log.Println("RecordFilterCount Get Statement Prepare Error", err)
		return values, err

	}
}

func (dbc DbConn) RecordFilterList(req map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Println("In side RecordFilterList")
	values := []map[string]interface{}{}
	pint := []interface{}{}
	var queryStr = "SELECT recordsearchfilter.* FROM recordsearchfilter "

	queryStr = queryStr + " WHERE recordsearchfilter.activeflg=1 AND recordsearchfilter.deleteflg=0 AND recordsearchfilter.userid=? "

	pint = append(pint, req["userid"])
	queryStr = queryStr + " ORDER BY createdate DESC "

	if req["limit"] != nil && req["offset"] != nil {
		queryStr = queryStr + " LIMIT ?,?"
		pint = append(pint, req["offset"])
		pint = append(pint, req["limit"])
	}
	//logger.Log.Println("Query:===>", queryStr)
	//logger.Log.Println("Params:===>", pint)
	rows, err := dbc.DB.Query(queryStr, pint...)
	if err != nil {
		logger.Log.Println("RecordFilterList Get Statement Prepare Error", err)
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

func (dbc DbConn) RecordFilterDelete(req map[string]interface{}) (bool, error) {
	logger.Log.Println("In side RecordFilterDelete")
	pint := []interface{}{}
	var queryStr = "UPDATE recordsearchfilter SET activeflg=0,deleteflg=1,createdate=NOW() WHERE id=?"

	pint = append(pint, req["id"])

	//logger.Log.Println("Query:===>", queryStr)
	//logger.Log.Println("Params:===>", pint)

	_, err := dbc.DB.Exec(queryStr, pint...)
	if err != nil {
		logger.Log.Println("RecordFilterDelete Get Statement Prepare Error", err)
		return false, err
	}

	return true, err
}

func (dbc DbConn) RecordFilterUpdate(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side RecordFilterUpdate")
	var values int64
	pint := []interface{}{}
	var queryStr = "UPDATE recordsearchfilter SET name=?,filter=?,savedfilters=? WHERE id=?"

	pint = append(pint, req["name"])
	// pint = append(pint, req["userid"])
	jsonData, err := json.Marshal(req["filter"])
	if err != nil {
		logger.Log.Println("RecordFilterUpdate Get Statement Prepare Error", err)
		return values, err
	}
	pint = append(pint, string(jsonData))
	pint = append(pint, req["savedfilters"])
	pint = append(pint, req["id"])

	logger.Log.Println("Query:===>", queryStr)
	logger.Log.Println("Params:===>", pint)

	res, err := dbc.DB.Exec(queryStr, pint...)
	if err != nil {
		logger.Log.Println("RecordFilterUpdate Get Statement Prepare Error", err)
		return values, err
	}
	values, err = res.LastInsertId()
	if err != nil {
		logger.Log.Println("RecordFilterUpdate Get Statement Prepare Error", err)
		return values, err
	}
	return values, err
}



/*func (dbc DbConn) RecordFilterUpdate(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side RecordFilterUpdate")
	var values int64
	pint := []interface{}{}
	var queryStr = "UPDATE recordsearchfilter SET name=?,filter=?,savedfilters=? WHERE id=?"

	pint = append(pint, req["name"])
	// pint = append(pint, req["userid"])
	jsonData, err := json.Marshal(req["filter"])
	if err != nil {
		logger.Log.Println("RecordFilterUpdate Get Statement Prepare Error", err)
		return values, err
	}
	pint = append(pint, string(jsonData))
	pint = append(pint, req["savedfilters"])
	pint = append(pint, req["id"])

	logger.Log.Println("Query:===>", queryStr)
	logger.Log.Println("Params:===>", pint)

	res, err := dbc.DB.Exec(queryStr, pint...)
	if err != nil {
		logger.Log.Println("RecordFilterUpdate Get Statement Prepare Error", err)
		return values, err
	}
	values, err = res.LastInsertId()
	if err != nil {
		logger.Log.Println("RecordFilterUpdate Get Statement Prepare Error", err)
		return values, err
	}
	return values, err
}*/
