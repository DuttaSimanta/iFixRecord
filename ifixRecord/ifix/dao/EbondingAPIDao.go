package dao

import (
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

func (mdao DbConn) GetModuleUrl(tz *entities.EbondingRecordEntity) (string, error) {
	var moduleUrl string
	var sql = "Select url from ebondingmodulemstmap where ebondingseq=? and moduleseq=? and activeflg=1 and deleteflg=0 "
	rows, err := mdao.DB.Query(sql, tz.EbondingSeq, tz.EbondingModuleSeq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetModuleUrl Get Statement Prepare Error", err)
		return moduleUrl, err
	}
	for rows.Next() {
		err = rows.Scan(&moduleUrl)
		if err != nil {
			logger.Log.Println("GetModuleUrl rows.next() Error", err)
			return moduleUrl, err
		}

	}
	return moduleUrl, nil
}
func (mdao DbConn) GetUserDetails(userID int64) (string, string, string, error) {
	var userName string
	var userEmail string
	var userContact string

	var sql = "SELECT COALESCE(loginname,'') loginname ,  COALESCE(useremail,'') requesteremail,  COALESCE(usermobileno,'') reuestercontactno FROM mstclientuser where id = ? and activeflag = 1 and deleteflag = 0;"
	rows, err := mdao.DB.Query(sql, userID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetUserDetails Get Statement Prepare Error", err)
		return userName, userEmail, userContact, err
	}
	for rows.Next() {
		err = rows.Scan(&userName, &userEmail, &userContact)
		if err != nil {
			logger.Log.Println("GetUserDetails rows.next() Error", err)
			return userName, userEmail, userContact, err
		}

	}

	return userName, userEmail, userContact, nil

}

func (mdao DbConn) GetDescription(tz *entities.EbondingRecordEntity) (string, string, int64, error) {
	var shortDesc string
	var longDesc string
	var userID int64
	var sql = "SELECT recordtitle,recorddescription,userid FROM trnrecord WHERE id = ? AND clientid = ? AND mstorgnhirarchyid = ? AND activeflg = 1 and deleteflg=0"
	rows, err := mdao.DB.Query(sql, tz.RecordID, tz.ClientID, tz.MstorgnhirarchyID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetDescription Get Statement Prepare Error", err)
		return shortDesc, longDesc, userID, err
	}
	for rows.Next() {
		err = rows.Scan(&shortDesc, &longDesc, &userID)
		if err != nil {
			logger.Log.Println("GetDescription rows.next() Error", err)
			return shortDesc, longDesc, userID, err
		}

	}
	return shortDesc, longDesc, userID, nil
}
func (mdao DbConn) GetCategories(tz *entities.EbondingRecordEntity) ([]string, error) {
	var categories []string
	var sql = "select name from ebondingdifferentiationmst where id in(select a.ebondingdiffid from ebondingdifferentiationmap" +
		" a where a.clientid=? and a.mstorgnhirarchyid=? and a.ebondingid=? and a.recorddiffid in( select id from" +
		" mstrecorddifferentiation where seqno in(3,4,5) and id in (select recorddiffid from  maprecordtorecorddifferentiation" +
		" where recordid=? and islatest=1 and recorddifftypeid  in(  SELECT id FROM mstrecorddifferentiationtype where parentid=1 " +
		" and  id in (SELECT torecorddifftypeid FROM mstrecordtype where clientid=a.clientid and mstorgnhirarchyid=a.mstorgnhirarchyid " +
		" and torecorddiffid=0 and fromrecorddiffid=?) and deleteflg=0 and activeflg=1)) and activeflg=1 and deleteflg=0) and " +
		" activeflg=1 and deleteflg=0)"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.EbondingID, tz.RecordID, tz.RecorddiffID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetCategories Get Statement Prepare Error", err)
		return categories, err
	}
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			logger.Log.Println("GetCategories rows.next() Error", err)
			return categories, err
		}
		categories = append(categories, category)

	}
	return categories, nil
}

func (mdao DbConn) GetImpact(tz *entities.EbondingRecordEntity) (string, error) {
	var impact string
	var sql = "select name from ebondingdifferentiationmst where id in(select a.ebondingdiffid from ebondingdifferentiationmap a where a.clientid=? and a.mstorgnhirarchyid=? and a.ebondingid=? and a.recorddiffid in (select recorddiffid from  maprecordtorecorddifferentiation where recordid=? and islatest=1 and recorddifftypeid = 7 and activeflg=1 and deleteflg=0) and  a.deleteflg=0 and a.activeflg=1)"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.EbondingID, tz.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetImpact Get Statement Prepare Error", err)
		return impact, err
	}
	for rows.Next() {

		err = rows.Scan(&impact)
		if err != nil {
			logger.Log.Println("GetImpact rows.next() Error", err)
			return impact, err
		}

	}
	return impact, nil
}
func (mdao DbConn) GetUrgency(tz *entities.EbondingRecordEntity) (string, error) {
	var urgency string
	var sql = "select name from ebondingdifferentiationmst where id in(select a.ebondingdiffid from ebondingdifferentiationmap a where a.clientid=? and a.mstorgnhirarchyid=? and a.ebondingid=? and a.recorddiffid in (select recorddiffid from  maprecordtorecorddifferentiation where recordid=? and islatest=1 and recorddifftypeid = 8 and activeflg=1 and deleteflg=0) and  a.deleteflg=0 and a.activeflg=1)"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.EbondingID, tz.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetUrgency Get Statement Prepare Error", err)
		return urgency, err
	}
	for rows.Next() {

		err = rows.Scan(&urgency)
		if err != nil {
			logger.Log.Println("GetUrgency rows.next() Error", err)
			return urgency, err
		}

	}
	return urgency, nil
}

func (mdao DbConn) GetUpdatedTicketStatusFromLog(tz *entities.EbondingRecordEntity) (int64, error) {
	var currentStatusSeq int64

	var sql = "select statusseq  from externalapitransactionlog where clientid =? and mstorgnhirarchyid=? and recordid=? and islatest=1 and activeflg=1 and deleteflg=0"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetStateID Get Statement Prepare Error", err)
		return currentStatusSeq, err
	}
	for rows.Next() {

		err = rows.Scan(&currentStatusSeq)
		if err != nil {
			logger.Log.Println("GetUpdatedTicketStatusFromLog rows.next() Error", err)
			return currentStatusSeq, err
		}

	}

	return currentStatusSeq, nil
}
func (mdao DbConn) GetStateID(tz *entities.EbondingRecordEntity) (string, error) {
	var stateID string
	var sql = "select name from ebondingdifferentiationmst where id in (select ebondingdiffid from ebondingdifferentiationmap where clientid=? and mstorgnhirarchyid=? and recorddiffid=? and ebondingid=? and recorddifftypeid=3 and activeflg=1 and deleteflg=0) and activeflg=1 and deleteflg=0"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.StatusID, tz.EbondingID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetStateID Get Statement Prepare Error", err)
		return stateID, err
	}
	for rows.Next() {

		err = rows.Scan(&stateID)
		if err != nil {
			logger.Log.Println("GetStateID rows.next() Error", err)
			return stateID, err
		}

	}
	return stateID, nil
}
func (mdao DbConn) InsertTransactionLog(transactionEntity *entities.EbondingTransactionLog) (string, error) {
	logger.Log.Println("In side InsertTransactionLog")
	var success string
	var sql = "insert into ebondingtransactionlogsmst (ebondingid,ifixrecordid,requestjson,responsejson,responsecode) values (?,?,?,?,?);"

	rows, err := mdao.DB.Query(sql, transactionEntity.Ebondingid, transactionEntity.RecordID, transactionEntity.Requestjson, transactionEntity.Responsejson, transactionEntity.Responsecode)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("InsertTransactionLog Statement Prepare Error", err)
		return success, err
	}
	for rows.Next() {
		err = rows.Scan(&success)
		if err != nil {
			logger.Log.Println("InsertTransactionLog rows.next() Error", err)
			return success, err
		}
		logger.Log.Println("InsertTransactionLog rows.next() Error", err)
	}
	return success, nil
}

func (mdao DbConn) GetMstRecordtermId(tz *entities.EbondingRecordEntity) (int64, error) {
	logger.Log.Println("In side GetMstRecordtermId")
	var Recordtermid int64
	var sql = "SELECT id FROM mstrecordterms where clientid= ? AND mstorgnhirarchyid = ? AND seq = 85;"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetMstRecordtermId Statement Prepare Error", err)
		return Recordtermid, err
	}
	for rows.Next() {
		err = rows.Scan(&Recordtermid)
		if err != nil {
			logger.Log.Println("GetReopencount rows.next() Error", err)
			return Recordtermid, err
		}

	}
	return Recordtermid, nil
}

func (mdao DbConn) InsertExternalTicketId(tz *entities.EbondingRecordEntity, Recordtermid int64, ExternalTicket string) (string, error) {
	logger.Log.Println("In side InsertExternalTicketId")
	var success string
	var sql = "insert into trnreordtracking (clientid,mstorgnhirarchyid,recordid,recordstageid,recordtermid,recordtrackvalue) values (?,?,?,?,?,?) ;"
	logger.Log.Println(tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID, tz.RecordStagedID, Recordtermid, ExternalTicket)

	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID, tz.RecordStagedID, Recordtermid, ExternalTicket)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("InsertExternalTicketId Statement Prepare Error", err)
		return success, err
	}
	for rows.Next() {
		err = rows.Scan(&success)
		if err != nil {
			logger.Log.Println("InsertExternalTicketId rows.next() Error", err)

			return success, err
		}
	}
	return success, nil
}

func (mdao DbConn) UpdateExternalTicketId(tz *entities.EbondingRecordEntity, Recordtermid int64, ExternalTicket string) (string, error) {
	logger.Log.Println("In side UpdateExternalTicketId")
	var success string
	var sql = "update trnreordtracking set recordtrackvalue=? where clientid=? and mstorgnhirarchyid=? and recordid=? and recordtermid=?"
	logger.Log.Println(ExternalTicket, tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID, Recordtermid)

	rows, err := mdao.DB.Query(sql, ExternalTicket, tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID, Recordtermid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("UpdateExternalTicketId Statement Prepare Error", err)
		return success, err
	}
	for rows.Next() {
		err = rows.Scan(&success)
		if err != nil {
			logger.Log.Println("UpdateExternalTicketId rows.next() Error", err)

			return success, err
		}
	}
	return success, nil
}
func (mdao DbConn) GetStateWiseComment(tz *entities.EbondingRecordEntity, termseq int64) (string, error) {
	logger.Log.Println("In side GetStateWiseComment")
	var stateWiseComment string
	sql := "SELECT a.recordtrackvalue FROM trnreordtracking a,mstrecordterms b where a.clientid=? and a.mstorgnhirarchyid=? and a.recordid=? and a.recordtermid=b.id and b.seq=? and a.deleteflg=0 and a.activeflg=1 ;"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID, termseq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("stateWiseComment Get Statement Prepare Error", err)
		return stateWiseComment, err
	}

	for rows.Next() {

		err = rows.Scan(&stateWiseComment)
		if err != nil {
			logger.Log.Println("stateWiseComment rows.next() Error", err)

			return stateWiseComment, err
		}
	}
	return stateWiseComment, nil
}
func (mdao DbConn) GetExternalTicketID(tz *entities.EbondingRecordEntity) (string, error) {
	logger.Log.Println("In side GetExternalTicketID")
	var externalTicketID string
	sql := "SELECT a.recordtrackvalue FROM trnreordtracking a,mstrecordterms b where a.clientid=? and a.mstorgnhirarchyid=? and a.recordid=? and a.recordtermid=b.id and b.seq=85 and a.deleteflg=0 and a.activeflg=1 ;"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetExternalTicketID Get Statement Prepare Error", err)
		return externalTicketID, err
	}

	for rows.Next() {

		err = rows.Scan(&externalTicketID)
		if err != nil {
			logger.Log.Println("GetExternalTicketID rows.next() Error", err)

			return externalTicketID, err
		}
	}
	return externalTicketID, nil
}
func (mdao DbConn) GetExternalWorkNote(tz *entities.EbondingRecordEntity) (string, error) {
	logger.Log.Println("In side GetExternalWorkNote")
	var externalWorkNote string
	sql := "SELECT a.recordtrackvalue FROM iFIX.trnreordtracking a,mstrecordterms b where a.clientid=? and a.mstorgnhirarchyid=? and a.recordid=? and a.recordtermid=b.id and b.seq=11 and a.deleteflg=0 and a.activeflg=1 ;"
	rows, err := mdao.DB.Query(sql, tz.ClientID, tz.MstorgnhirarchyID, tz.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetExternalWorkNote Get Statement Prepare Error", err)
		return externalWorkNote, err
	}

	for rows.Next() {

		err = rows.Scan(&externalWorkNote)
		if err != nil {
			logger.Log.Println("GetExternalWorkNote rows.next() Error", err)

			return externalWorkNote, err
		}
	}
	return externalWorkNote, nil
}

func (mdao DbConn) GetAttachment(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) ([]entities.FileAttachmentEntity, error) {
	logger.Log.Println("In side GetAttachment")
	values := []entities.FileAttachmentEntity{}
	recordtypediffid := "SELECT a.clientid,a.mstorgnhirarchyid,a.recordtrackvalue,a.recordtrackdescription FROM iFIX.trnreordtracking a,mstrecordterms b where a.clientid=? and a.mstorgnhirarchyid=? and a.recordid=? and a.recordtermid=b.id and b.seq=1 and a.deleteflg=0 and a.activeflg=1 ;"
	rows, err := mdao.DB.Query(recordtypediffid, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAttachment Get Statement Prepare Error", err)
		return values, err
	}

	for rows.Next() {
		value := entities.FileAttachmentEntity{}

		err = rows.Scan(&value.Clientid, &value.Mstorgnhirarchyid, &value.OriginalFileName, &value.UploadedFileName)
		logger.Log.Println("GetAttachment rows.next() Error", err)
		values = append(values, value)
	}
	return values, nil
}
func (mdao DbConn) GetEbondingModeleName(tz *entities.EbondingRecordEntity) (string, error) {
	var modulename string
	var sql = "SELECT modulename  FROM ebondingmodulemst where seqno=? and deleteflg=0;"
	rows, err := mdao.DB.Query(sql, tz.EbondingModuleSeq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetCategories Get Statement Prepare Error", err)
		return modulename, err
	}
	for rows.Next() {

		err = rows.Scan(&modulename)
		if err != nil {
			logger.Log.Println("GetCategories rows.next() Error", err)
			return modulename, err
		}

	}
	return modulename, nil
}
func (mdao DbConn) GetebondingId(tz *entities.EbondingRecordEntity) (int64, error) {
	var ebondingSeqName int64
	var sql = "SELECT id  FROM ebondingmst where seqno=? and deleteflg=0;"
	rows, err := mdao.DB.Query(sql, tz.EbondingSeq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetCategories Get Statement Prepare Error", err)
		return ebondingSeqName, err
	}
	for rows.Next() {

		err = rows.Scan(&ebondingSeqName)
		if err != nil {
			logger.Log.Println("GetCategories rows.next() Error", err)
			return ebondingSeqName, err
		}

	}
	return ebondingSeqName, nil
}
