package dao

import (
	"database/sql"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"log"
	"strings"
	"time"
)

var gettimediff = "SELECT a.utcdiff as timediff,c.utcdiff as reporttimediff,coalesce(d.example,'') Timeformat,coalesce(e.example,'') Reporttimeformat from zone a,mstorgnhierarchy b,zone c,mstdatetimeformat d,mstdatetimeformat e where a.zone_id=b.timezoneid and b.reporttimezoneid=c.zone_id and b.clientid=? and b.id=? and b.timeformatid=d.id and b.reporttimeformatid=e.id"
var deletetoken = "DELETE from msttoken WHERE userid=? and deleteflg=0"
var inserttoken = "INSERT into msttoken(token,userid) values(?,?)"
var gettoken = " SELECT token from msttoken where userid=? and deleteflg=0"

func Deletetoken(tx *sql.Tx, id int64) error {
	log.Println("In side dao")
	stmt, err := tx.Prepare(deletetoken)

	if err != nil {
		logger.Log.Print("deletetoken Prepare Statement  Error", err)
		log.Print("deletetoken Prepare Statement  Error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Log.Print("deletetoken Execute Statement  Error", err)
		log.Print("deletetoken Execute Statement  Error", err)
		return err
	}
	return nil
}
func Inserttoken(tx *sql.Tx, id int64, token string) error {
	log.Println("In side dao")
	stmt, err := tx.Prepare(inserttoken)

	if err != nil {
		logger.Log.Print("inserttoken Prepare Statement  Error", err)
		log.Print("inserttoken Prepare Statement  Error", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(token, id)
	if err != nil {
		logger.Log.Print("inserttoken Execute Statement  Error", err)
		log.Print("inserttoken Execute Statement  Error", err)
		return err
	}
	return nil
}
func (mdao DbConn) Gettoken(userid int64) (error, []string) {
	var tokens []string
	stmt, err := mdao.DB.Prepare(gettoken)
	if err != nil {
		logger.Log.Print("Gettoken Statement Prepare Error", err)
		log.Print("Gettoken Statement Prepare Error", err)
		return err, tokens
	}
	defer stmt.Close()
	rows, err := stmt.Query(userid)
	if err != nil {
		logger.Log.Print("Gettoken Statement Execution Error", err)
		log.Print("Gettoken Statement Execution Error", err)
		return err, tokens
	}
	for rows.Next() {
		var token string
		rows.Scan(&token)
		tokens = append(tokens, token)
	}
	return nil, tokens
}

func (mdao DbConn) Gettimediff(ClientID int64, MstorgnhirarchyID int64) (entities.UtilityEntity, error) {
	//requestIds := []entities.UtilityEntity{}
	requestIds := entities.UtilityEntity{}
	stmt, err := mdao.DB.Prepare(gettimediff)
	if err != nil {
		logger.Log.Print("Gettimediff Statement Prepare Error", err)
		log.Print("Gettimediff Statement Prepare Error", err)
		return requestIds, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(ClientID, MstorgnhirarchyID)
	if err != nil {
		logger.Log.Print("Gettimediff Statement Execution Error", err)
		log.Print("Gettimediff Statement Execution Error", err)
		return requestIds, err
	}
	defer rows.Close()
	for rows.Next() {
		//value := entities.UtilityEntity{}
		rows.Scan(&requestIds.Timediff, &requestIds.Reporttimediff, &requestIds.Timeformat, &requestIds.Reporttimeformat)
		//requestIds = append(requestIds, value)
	}
	return requestIds, nil
}

// func Convertdate(date int64, timediff int64, format string) string {
// 	var unixTime = date + timediff
// 	t := time.Unix(unixTime, 0)
// 	//return t.Format("02-Jan-2006 15:04:05")
// 	return t.Format(format)
// }

func Convertdate(date int64, timediff int64) string {
	var unixTime = date + timediff
	t := time.Unix(unixTime, 0)
	return t.Format("02-Jan-2006 15:04:05")
	//return t.Format(format)
}

func Convertdate1(date int64, timediff int64) string {
	var unixTime = date + timediff
	t := time.Unix(unixTime, 0)
	logger.Log.Println("creationdate +++++++++++++++++++++++++++++++++++ >>>>>>>>>>>>>>>>>>>>>>>>>>>>>", t)
	logger.Log.Println("creationdate +++++++++++++++++++++++++++++++++++ >>>>>>>>>>>>>>>>>>>>>>>>>>>>>", t.String())
	var bb2 string
	if len(t.String()) > 0 {
		var bb1 = strings.Split(t.String(), " ")
		bb2 = bb1[0] + " " + bb1[1]
	}
	return bb2
	//return t.Format(format)
}
