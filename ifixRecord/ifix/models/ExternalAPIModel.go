package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/dbconfig"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	//	"bytes"
	//  "database/sql"
	"encoding/base64"
	// "encoding/json"
	//    "errors"

	FileUtils "ifixRecord/ifix/fileutils"

	"io"

	//  "log"

	"os"
	//    "time"
	// "golang.org/x/crypto/openpgp/errors"
)

func getContextPath() (string, error) {

	wd, err := os.Getwd()
	if err != nil {
		return "", errors.New("ERROR: Unable to get WD")
	}
	contextPath := strings.ReplaceAll(wd, "\\", "/") // replacing backslash by  forwardslash
	return contextPath, nil
}

func UpdateExternalRecordAttachment(tz *entities.FileAttacmentToRecordEntity) (bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }

	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
		logger.Log.Println(db)
	}
	dataAccess := dao.DbConn{DB: db}

	contextPath, contextPatherr := getContextPath()
	if contextPatherr != nil {
		logger.Log.Println(contextPatherr)
		return false, contextPatherr, ""
	}

	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.LoginID) == 0 {
		return false, err, "LoginID is missing"
	}

	if len(tz.Recordid) == 0 {
		return false, err, "RecordID is missing"
	}

	if len(tz.Logingrpname) == 0 {
		return false, err, "Group Name is missing"
	}

	statusname, err := dataAccess.GetStatusName(tz.Recordid)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Statusname is not valid"
	}
	logger.Log.Println(statusname)
	if statusname == "Closed" && statusname == "Resolved" && statusname == "Cancel" && statusname == "Rejected" {
		return false, err, "Statusname is not valid"
	}

	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Clientname is not valid"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}

	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Organization name is not valid"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}

	loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "LoginID is not valid"
	}
	if loginID == 0 {
		return false, err, "LoginID is not valid"
	}

	recordid, err := dataAccess.GetRecordId(tz.Recordid)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Ticket id is not valid"
	}
	if recordid == 0 {
		return false, err, "Ticket id is not valid"
	}

	createdgrpid, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.Logingrpname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Group id is not valid"
	}
	if createdgrpid == 0 {
		return false, err, "Group id is not valid"
	}

	for i := 0; i < len(tz.Fileattachment); i++ {
		if len(tz.Fileattachment[i].Filename) == 0 || len(tz.Fileattachment[i].Filecontent) == 0 || len(tz.Fileattachment[i].Filetype) == 0 {
			return false, err, "Please provide a valid file"
		}

		rawDecodedText, err := base64.StdEncoding.DecodeString(tz.Fileattachment[i].Filecontent)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}

		// logger.Log.Println(rawDecodedText)

		filePath := contextPath + "/ifix/resource/downloads/" + tz.Fileattachment[i].Filename

		logger.Log.Println(filePath)
		out, err := os.Create(filePath)
		if err != nil {
			// return errors.New("ERROR: Unable to create copy of file")
		}
		defer out.Close()

		// Write the content to file

		_, err = io.Copy(out, bytes.NewReader(rawDecodedText))

		// props, err := FileUtils.ReadPropertiesFile(contextPath + "/ifix/dbconfig/config.go")
		// logger.Log.Println(props["FileUploadUrl"])
		originalFileName, newFileName, err := FileUtils.FileUploadAPICall(clientID, orgnID, dbconfig.FileUploadUrl, filePath)
		if err != nil {
			logger.Log.Println("Error while downloading", "-", err)
		}
		logger.Log.Println(originalFileName, newFileName)

		recordtermid, err := dataAccess.GetTermIdToRecord(clientID, orgnID)
		logger.Log.Println(recordtermid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, ""
		}

		recordstageid, err := dataAccess.GetRecordstageid(recordid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, ""
		}

		recorddiffid, err := dataAccess.GetRecordDiffId(clientID, orgnID, recordid)
		logger.Log.Println(recorddiffid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Recorddiffid is not valid"
		}

		// attachmentSuccess, err := dataAccess.InsertAttachmentToRecord(clientID, orgnID, recordid, recordstageid, recordtermid, originalFileName, newFileName, loginID, createdgrpid)
		// logger.Log.Println(attachmentSuccess)
		// if err != nil {
		// 	logger.Log.Println(err)
		// 	return false, err, "Attachments have not been added successfully"
		// }
		if recordstageid > 0 {
			insertlog := entities.RecordcommonEntity{}
			insertlog.ClientID = clientID
			insertlog.Mstorgnhirarchyid = orgnID
			insertlog.RecordID = recordid
			insertlog.RecordstageID = recordstageid
			insertlog.Usergroupid = createdgrpid
			insertlog.Userid = tz.Userid
			insertlog.Termvalue = originalFileName
			insertlog.Termdescription = newFileName
			insertlog.Termseq = 1
			insertlog.Recorddiffid = recorddiffid
			insertlog.Recorddifftypeid = 2

			_, _, err, _ := InsertRecordTermvaluesforAttachment(&insertlog)
			if err != nil {
				logger.Log.Println("Error is --77777777777777777777777------->", err)
				return false, err, "Something Went Wrong"
			}
		}

		// id = 0
		// msg = "Attachments have been added successfully"
		// return id, status, err, msg
	}

	return true, err, "Attachments have been added successfully"

}

/*func UpdateExternalRecordAttachment(tz *entities.FileAttacmentToRecordEntity) (bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }

	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
		logger.Log.Println(db)
	}
	dataAccess := dao.DbConn{DB: db}

	contextPath, contextPatherr := getContextPath()
	if contextPatherr != nil {
		logger.Log.Println(contextPatherr)
		return false, contextPatherr, ""
	}

	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.LoginID) == 0 {
		return false, err, "LoginID is missing"
	}

	if len(tz.Recordid) == 0 {
		return false, err, "RecordID is missing"
	}

	if len(tz.Logingrpname) == 0 {
		return false, err, "Group Name is missing"
	}

	statusname, err := dataAccess.GetStatusName(tz.Recordid)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Statusname is not valid"
	}
	logger.Log.Println(statusname)
	if statusname == "Closed" && statusname == "Resolved" && statusname == "Cancel" && statusname == "Rejected" {
		return false, err, "Statusname is not valid"
	}

	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Clientname is not valid"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}

	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Organization name is not valid"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}

	loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "LoginID is not valid"
	}
	if loginID == 0 {
		return false, err, "LoginID is not valid"
	}

	recordid, err := dataAccess.GetRecordId(tz.Recordid)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Ticket id is not valid"
	}
	if recordid == 0 {
		return false, err, "Ticket id is not valid"
	}

	createdgrpid, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.Logingrpname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Group id is not valid"
	}
	if createdgrpid == 0 {
		return false, err, "Group id is not valid"
	}

	for i := 0; i < len(tz.Fileattachment); i++ {
		if len(tz.Fileattachment[i].Filename) == 0 || len(tz.Fileattachment[i].Filecontent) == 0 || len(tz.Fileattachment[i].Filetype) == 0 {
			return false, err, "Please provide a valid file"
		}

		rawDecodedText, err := base64.StdEncoding.DecodeString(tz.Fileattachment[i].Filecontent)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}

		// logger.Log.Println(rawDecodedText)

		filePath := contextPath + "/ifix/resource/downloads/" + tz.Fileattachment[i].Filename

		logger.Log.Println(filePath)
		out, err := os.Create(filePath)
		if err != nil {
			// return errors.New("ERROR: Unable to create copy of file")
		}
		defer out.Close()

		// Write the content to file

		_, err = io.Copy(out, bytes.NewReader(rawDecodedText))

		// props, err := FileUtils.ReadPropertiesFile(contextPath + "/ifix/dbconfig/config.go")
		// logger.Log.Println(props["FileUploadUrl"])
		originalFileName, newFileName, err := FileUtils.FileUploadAPICall(clientID, orgnID, dbconfig.FileUploadUrl, filePath)
		if err != nil {
			logger.Log.Println("Error while downloading", "-", err)
		}
		logger.Log.Println(originalFileName, newFileName)

		recordtermid, err := dataAccess.GetTermIdToRecord(clientID, orgnID)
logger.Log.Println(recordtermid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, ""
		}

		recordstageid, err := dataAccess.GetRecordstageid(recordid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, ""
		}

		recorddiffid, err := dataAccess.GetRecordDiffId(clientID, orgnID, recordid)
		logger.Log.Println(recorddiffid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Recorddiffid is not valid"
		}

		attachmentSuccess, err := dataAccess.InsertAttachmentToRecord(clientID, orgnID, recordid, recordstageid, recordtermid, originalFileName, newFileName, loginID, createdgrpid)
		logger.Log.Println(attachmentSuccess)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Attachments have not been added successfully"
		}
		if recordstageid > 0 {
			insertlog := entities.RecordcommonEntity{}
			insertlog.ClientID = clientID
			insertlog.Mstorgnhirarchyid = orgnID
			insertlog.RecordID = recordid
			insertlog.RecordstageID = recordstageid
			insertlog.Usergroupid = createdgrpid
			insertlog.Userid = tz.Userid
			insertlog.Termvalue = originalFileName
			insertlog.Termdescription = newFileName
			insertlog.Termseq = 1
			insertlog.Recorddiffid = recorddiffid
			insertlog.Recorddifftypeid = 2

			_, _, err, _ := InsertRecordTermvalues(&insertlog)
			if err != nil {
				logger.Log.Println("Error is --77777777777777777777777------->", err)
				return false, err, "Something Went Wrong"
			}
		}

		// id = 0
		// msg = "Attachments have been added successfully"
		// return id, status, err, msg
	}

	return true, err, "Attachments have been added successfully"

}*/

/*func UpdateExternalRecordAttachment(tz *entities.FileAttacmentToRecordEntity) (bool, error, string) {
    if mutexutility.MutexLocked(lock) == false {
        lock.Lock()
        defer lock.Unlock()
    }

    db, err :=ConnectMySqlDb()
    if err != nil {
        logger.Log.Println("database connection failure", err)
        return false, err, "Something Went Wrong"
    }

    dataAccess := dao.DbConn{DB: db}

    contextPath, contextPatherr := getContextPath()
    if contextPatherr != nil {
        logger.Log.Println(contextPatherr)
        return false, contextPatherr, ""
    }

    if len(tz.Clientname) == 0 {
        return false, err, "Client name is missing"
    }
    if len(tz.Mstorgnhirarchyname) == 0 {
        return false, err, "Organization name is missing"
    }
    if len(tz.LoginID) == 0 {
        return false, err, "LoginID is missing"
    }

    if len(tz.Recordid) == 0 {
        return false, err, "RecordID is missing"
    }

    if len(tz.Logingrpname) == 0 {
        return false, err, "Group Name is missing"
    }

    statusname, err := dataAccess.GetStatusName(tz.Recordid)
    if err != nil {
        logger.Log.Println(err)
        return false, err, "Statusname is not valid"
    }

    if statusname == "Closed" && statusname == "Resolved" && statusname == "Cancel" && statusname == "Rejected" {
        return false, err, "Statusname is not valid"
    }

    clientID, err := dataAccess.GetClientID(tz.Clientname)
    if err != nil {
        logger.Log.Println(err)
        return false, err, "Clientname is not valid"
    }

    if clientID == 0 {
        return false, err, "Clientname is not valid"
    }

    orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
    if err != nil {
        logger.Log.Println(err)
        return false, err, "Organization name is not valid"
    }
    if orgnID == 0 {
        return false, err, "Organization name is not valid"
    }

    loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
    if err != nil {
        logger.Log.Println(err)
        return false, err, "LoginID is not valid"
    }
    if loginID == 0 {
        return false, err, "LoginID is not valid"
    }

    recordid, err := dataAccess.GetRecordId(tz.Recordid)
    if err != nil {
        logger.Log.Println(err)
        return false, err, "Ticket id is not valid"
    }
    if recordid == 0 {
        return false, err, "Ticket id is not valid"
    }

    createdgrpid, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.Logingrpname)
    if err != nil {
        logger.Log.Println(err)
        return false, err, "Group id is not valid"
    }
    if createdgrpid == 0 {
        return false, err, "Group id is not valid"
    }

    for i := 0; i < len(tz.Fileattachment); i++ {
        if len(tz.Fileattachment[i].Filename) == 0 || len(tz.Fileattachment[i].Filecontent) == 0 || len(tz.Fileattachment[i].Filetype) == 0 {
            return false, err, "Please provide a valid file"
        }

        rawDecodedText, err := base64.StdEncoding.DecodeString(tz.Fileattachment[i].Filecontent)
        if err != nil {
            logger.Log.Println(err)
            return false, err, "Something Went Wrong"
        }

        // logger.Log.Println(rawDecodedText)

        filePath := contextPath + "/ifix/resource/downloads/" + tz.Fileattachment[i].Filename

        logger.Log.Println(filePath)
        out, err := os.Create(filePath)
        if err != nil {
            // return errors.New("ERROR: Unable to create copy of file")
        }
        defer out.Close()

        // Write the content to file

        _, err = io.Copy(out, bytes.NewReader(rawDecodedText))

        // props, err := FileUtils.ReadPropertiesFile(contextPath + "/ifix/dbconfig/config.go")
        // logger.Log.Println(props["FileUploadUrl"])
        originalFileName, newFileName, err := FileUtils.FileUploadAPICall(clientID, orgnID, dbconfig.FileUploadUrl, filePath)
        if err != nil {
            logger.Log.Println("Error while downloading", "-", err)
        }
        logger.Log.Println(originalFileName, newFileName)

        recordtermid, err := dataAccess.GetTermIdToRecord(clientID, orgnID)
        if err != nil {
            logger.Log.Println(err)
            return false, err, ""
        }

        recordstageid, err := dataAccess.GetRecordstageid(recordid)
        if err != nil {
            logger.Log.Println(err)
            return false, err, ""
        }

        // currentTime := time.Now().UTC()

        attachmentSuccess, err := dataAccess.InsertAttachmentToRecord(clientID, orgnID, recordid, recordstageid, recordtermid, originalFileName, newFileName, loginID, createdgrpid)

        if err != nil {
            logger.Log.Println(err)
            return false, err, "Attachments have not been added successfully"
        }
        logger.Log.Println(attachmentSuccess)
    }

    return true, err, "Attachments have been added successfully"

}*/

//===================================================
func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetExternalRecordDetailsByNo(tz *entities.RecordDetailsRequestEntityAPI) (entities.RecordDetailsEntityAPI, bool, error, string) {
	logger.Log.Println("<===============START=============>")

	logger.Log.Println("In side GetRecordDetailsByNo", tz)
	t := entities.RecordDetailsEntityAPI{}
	var err error
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return t, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	var parentchildcontrol []int64
	data, success, _, msg := GetExternalRecordDetailsByNoRecursive(tz, db, parentchildcontrol)
	logger.Log.Println("<===============END=============>")

	return data, success, err, msg

}

func GetExternalRecordDetailsByNoRecursive(tz *entities.RecordDetailsRequestEntityAPI, db *sql.DB, parentchildcontrol []int64) (entities.RecordDetailsEntityAPI, bool, error, string) {
	logger.Log.Println("In side GetRecordDetailsByNo")
	t := entities.RecordDetailsEntityAPI{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return t, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	dataAccess := dao.DbConn{DB: db}

	if len(tz.Clientname) == 0 {
		return t, false, errors.New("Client name is missing"), "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return t, false, errors.New("Organization name is missing"), "Organization name is missing"
	}

	if len(tz.RecordNo) == 0 {
		return t, false, errors.New("Record number is missing"), "Record number is missing"
	}

	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return t, false, err, "Something Went Wrong"
	}
	// fmt.Println("client:", clientID)
	if clientID == 0 {
		return t, false, err, "Clientname is not valid"
	}
	tz.Userid = 0
	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return t, false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return t, false, err, "Organization name is not valid"
	}

	values, recordtypeid, err := dataAccess.GetExternalRecordDetailsByNo(clientID, orgnID, tz)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	// fmt.Println("values:", values)
	id, _, err := dataAccess.GetIdAgainstRecordNo(clientID, orgnID, tz.RecordNo)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	pcount, err := dataAccess.GetPrioritycount(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	rcount, err := dataAccess.GetReopencount(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	fcount, err := dataAccess.GetFollowupcount(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	obcount, err := dataAccess.GetOutboundcount(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	aging, err := dataAccess.GetAging(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	categories, err := dataAccess.GetExternalcategorynames(clientID, orgnID, 2, recordtypeid, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	lastupdated, err := dataAccess.GetLastupdatednames(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	c := entities.RecordcommonEntity{}
	c.ClientID = clientID
	c.Mstorgnhirarchyid = orgnID
	c.RecordID = id
	visiblecomment, _, err, _ := Customervisiblecomment(&c)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	workflows, err := dataAccess.GetWorkflowdetails(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	hopcount, _, err, _ := Gethopcount(id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	values.Prioritycount = pcount
	values.Reopencount = rcount
	values.Followupcount = fcount
	values.Outboundcount = obcount
	values.Aging = aging
	values.Categories = categories
	values.Visiblecommentdaycount = visiblecomment[0].Daycount
	values.Assignee = workflows.Asigneename
	values.AssignedGroup = workflows.Asigneegrp
	values.Hopcount = hopcount
	values.Latestupdatedby = lastupdated
	// New Addition in 02.02.2022
	// termvalues, err := dataAccess.GetTermsdetails(clientID, orgnID, id)
	// if err != nil {
	// 	return t, false, err, "Something Went Wrong"
	// }
	termvalues, err := dataAccess.GetTermsdetailss(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	values.Termsdetails = termvalues
	err1 := dataAccess.GetFromRecordfulldetails(clientID, orgnID, id, &values)
	if err1 != nil {
		return t, false, err, "Something Went Wrong"
	}
	// logger.Log.Println(values)
	lastmodified, err := dataAccess.Getlastmodified(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	totaleffort, err := dataAccess.GetTotaleffort(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	var totalhour int64
	var totalminute int64
	for i := 0; i < len(totaleffort); i++ {
		if totaleffort[i] != "" && len(totaleffort[i]) == 5 {
			hourtime := strings.Split(totaleffort[i], ":")
			hour, _ := strconv.ParseInt(hourtime[0], 10, 64)
			totalhour = totalhour + hour
			minute, _ := strconv.ParseInt(hourtime[1], 10, 64)
			totalminute = totalminute + minute
		}
	}
	// logger.Log.Println("totaleffort bofore:", totalhour, totalminute)
	totalhour = totalhour + (totalminute / 60)
	totalminute = totalminute % 60
	// logger.Log.Println("totaleffort now:", totalhour, totalminute)
	s := strconv.Itoa(int(totalhour))
	m := strconv.Itoa(int(totalminute))

	effort := s + ":" + m
	util, err1 := dataAccess.Gettimediff(clientID, orgnID)
	if err1 != nil {
		return t, false, err1, "Something Went Wrong"
	}
	values.TotalEffort = effort
	values.LastModifiedDate = dao.Convertdate(lastmodified, util.Timediff)
	err1 = dataAccess.Getlikrecord(clientID, orgnID, id, &values)
	if err1 != nil {
		return t, false, err, "Something Went Wrong"
	}
	err1 = dataAccess.GetSlacompliance(clientID, orgnID, recordtypeid, &values)
	if err1 != nil {
		return t, false, err, "Something Went Wrong"
	}
	child, err := dataAccess.Getchild(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}

	parentchildcontrol = append(parentchildcontrol, id)
	if len(child) > 0 {

		logger.Log.Println("child record is:", child)
		for i := 0; i < len(child); i++ {
			values.IsParent = "yes"
			values.ChildCount = int64(len(child))
			temp := entities.RecordDetailsRequestEntityAPI{}
			temp.Clientname = tz.Clientname
			temp.Mstorgnhirarchyname = tz.Mstorgnhirarchyname
			RecordNo, err := dataAccess.GetRecordNo(child[i])
			temp.RecordNo = RecordNo
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
			if contains(parentchildcontrol, child[i]) {
				continue
			}
			// logger.Log.Println("Recordno is:", RecordNo)
			// logger.Log.Println("Next call value:", temp)
			data := entities.RecordDetailsEntityAPI{}
			data, success, _, msg := GetExternalRecordDetailsByNoRecursive(&temp, db, parentchildcontrol)
			if success == false {
				logger.Log.Println(msg, err)
				return t, false, err, "Something Went Wrong"
			}
			values.Childdetails = append(values.Childdetails, data)

		}
	} else {
		values.IsParent = "no"
	}
	parent, err := dataAccess.Getparent(clientID, orgnID, id)
	if err != nil {
		return t, false, err, "Something Went Wrong"
	}
	// logger.Log.Println("parent record is:", parent)
	if len(parent) > 0 {

		for i := 0; i < len(parent); i++ {
			values.IsChild = "yes"
			values.ParentCount = int64(len(parent))
			temp := entities.RecordDetailsRequestEntityAPI{}
			temp.Clientname = tz.Clientname
			temp.Mstorgnhirarchyname = tz.Mstorgnhirarchyname
			RecordNo, err := dataAccess.GetRecordNo(parent[i])
			temp.RecordNo = RecordNo
			if err != nil {
				return t, false, err, "Something Went Wrong"
			}
			if contains(parentchildcontrol, parent[i]) {
				continue
			}
			// logger.Log.Println("Recordno is:", RecordNo)
			// logger.Log.Println("Next call value:", temp)
			data := entities.RecordDetailsEntityAPI{}
			data, success, _, msg := GetExternalRecordDetailsByNoRecursive(&temp, db, parentchildcontrol)
			if success == false {
				logger.Log.Println(msg, err)
				return t, false, err, "Something Went Wrong"
			}
			values.Parentdetails = append(values.Parentdetails, data)

		}
	} else {
		values.IsChild = "no"
	}

	// New Addition in 02.02.2022
	//logger.Log.Println("=========+++++++++++ ************** =======>       ", termvalues)
	//logger.Log.Println("=========+++++++++++ ************** =======>       ", values.Termdetails)
	// t = append(t, values)
	// logger.Log.Println("=========+++++++++++=======>       ", values)
	// logger.Log.Panic("Not Error")
	values.Recordid = tz.RecordNo
	return values, true, err, ""
}

func GetExternalRecordDetailsByDate(tz *entities.RecordDetailsRequestEntityAPI) ([]entities.RecordDetailsEntityAPI, bool, error, string) {
	logger.Log.Println("<===============START=============>")

	logger.Log.Println("In side GetRecordDetailsByNo")
	t := []entities.RecordDetailsRequestEntityAPI{}
	values := []entities.RecordDetailsEntityAPI{}
	var err error
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return values, false, err, "Something Went Wrong"
	// }
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection ")
			return values, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	//defer db.Close()
	dataAccess := dao.DbConn{DB: db}
	if len(tz.Clientname) == 0 {
		return values, false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return values, false, err, "Organization name is missing"
	}
	logger.Log.Println("In side GetRecordDetailsByNo----------->", len(tz.Fromdate))
	if len(tz.Fromdate) == 0 {
		return values, false, err, "From date is missing"
	}

	if len(tz.Todate) == 0 {
		return values, false, err, "To date is missing"
	}
	logger.Log.Println("fromdate is:", tz.Fromdate)
	if len(tz.Fromdate) != 0 {
		_, err = time.Parse("2006-01-02 15:04:05", tz.Fromdate)
		if err != nil {
			logger.Log.Println(err)
			return values, false, err, "From date format is not valid"
		}
	}

	if len(tz.Todate) != 0 {
		_, err = time.Parse("2006-01-02 15:04:05", tz.Todate)
		if err != nil {
			return values, false, err, "To date format is not valid"
		}
	}

	formdateInt := TimeParse(tz.Fromdate, "").Unix()
	todateInt := TimeParse(tz.Todate, "").Unix()

	if formdateInt > todateInt {
		return values, false, err, "From date is greater then To date."
	}

	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return values, false, err, "Something Went Wrong"
	}

	if clientID == 0 {
		return values, false, err, "Clientname is not valid"
	}
	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return values, false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return values, false, err, "Organization name is not valid"
	}
	util, err1 := dataAccess.Gettimediff(clientID, orgnID)
	if err1 != nil {
		return values, false, err1, "Something Went Wrong"
	}
	logger.Log.Println("timediffis", util.Timediff)
	fromdate, err := Getutctime(clientID, orgnID, tz.Fromdate, util.Timediff)
	if err != nil {
		return values, false, err1, "Something Went Wrong"
	}
	todate, err := Getutctime(clientID, orgnID, tz.Todate, util.Timediff)
	if err != nil {
		return values, false, err1, "Something Went Wrong"
	}
	tz.Fromdate = fromdate
	tz.Todate = todate
	// fmt.Println(tz, fromdate, todate)
	t, err = dataAccess.GetExternalRecordDetailsByDate(clientID, orgnID, tz)
	if err != nil {
		return values, false, err, "Something Went Wrong"
	}

	for i := 0; i < len(t); i++ {
		// tz.RecordNo = t[i].Recordid
		t[i].Clientname = tz.Clientname
		t[i].Mstorgnhirarchyname = tz.Mstorgnhirarchyname
		var parentchildcontrol []int64
		data, success, _, _ := GetExternalRecordDetailsByNoRecursive(&(t[i]), db, parentchildcontrol)
		if success == false || err != nil {
			return values, false, err, "Something Went Wrong"

		}
		values = append(values, data)
		// id, createdatetime, err := dataAccess.GetIdAgainstRecordNo(clientID, orgnID, t[i].Recordid)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }

		// pcount, err := dataAccess.GetPrioritycount(clientID, orgnID, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		// rcount, err := dataAccess.GetReopencount(clientID, orgnID, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		// fcount, err := dataAccess.GetFollowupcount(clientID, orgnID, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }

		// obcount, err := dataAccess.GetOutboundcount(clientID, orgnID, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }

		// aging, err := dataAccess.GetAging(clientID, orgnID, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		// categories, err := dataAccess.GetExternalcategorynames(clientID, orgnID, 2, 4, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		// lastupdated, err := dataAccess.GetLastupdatednames(clientID, orgnID, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		// c := entities.RecordcommonEntity{}
		// c.ClientID = clientID
		// c.Mstorgnhirarchyid = orgnID
		// c.RecordID = id
		// visiblecomment, _, err, _ := Customervisiblecomment(&c)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		// workflows, err := dataAccess.GetWorkflowdetails(clientID, orgnID, id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }
		// hopcount, _, err, _ := Gethopcount(id)
		// if err != nil {
		// 	return t, false, err, "Something Went Wrong"
		// }

		// t[i].Prioritycount = pcount
		// t[i].Reopencount = rcount
		// t[i].Followupcount = fcount
		// t[i].Outboundcount = obcount
		// t[i].Aging = aging
		// t[i].Categories = categories
		// t[i].Visiblecommentdaycount = visiblecomment[0].Daycount
		// t[i].Assignee = workflows.Asigneename
		// t[i].AssignedGroup = workflows.Asigneegrp
		// t[i].Hopcount = hopcount
		// t[i].Latestupdatedby = lastupdated
		// t[i].Createdatetime = dao.Convertdate(createdatetime, util.Timediff)
	}
	logger.Log.Println("<===============END=============>")

	return values, true, err, ""
}

func ExternalRecordCreate(tz *entities.ExternalCreateRecord) (bool, error, string) {
	v := entities.RecordEntity{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.RequestorID) == 0 {
		return false, err, "RequestorID is missing"
	}
	if len(tz.ShortDescription) == 0 {
		return false, err, "Short description value is missing."
	}
	if len(tz.LongDescription) == 0 {
		return false, err, "Long description value is missing."
	}
	if len(tz.TickettypeID) == 0 {
		return false, err, "Ticket type is missing."
	}
	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	// New addition

	if len(tz.OriginalID) == 0 {
		return false, err, "OriginalID is missing"
	}

	if len(tz.Originalgrpname) == 0 {
		return false, err, "OriginalgrpID is missing"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}
	v.ClientID = clientID
	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}
	v.Mstorgnhirarchyid = orgnID
	loginID, err := dataAccess.GetLoginID(tz.RequestorID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if loginID == 0 {
		return false, err, "RequestorID is not valid"
	}

	originalloginID, err := dataAccess.GetLoginID(tz.OriginalID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if originalloginID == 0 {
		return false, err, "OriginalID is not valid"
	}

	v.Originaluserid = originalloginID
	v.CreateduserID = loginID
	grpID, err := dataAccess.GetGrpID(loginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	logger.Log.Println("")
	if grpID == 0 {
		return false, err, "The User does not have the permission to raise a request"
	}

	originalgrpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.Originalgrpname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	logger.Log.Println("")
	if originalgrpID == 0 {
		return false, err, "Originalgrpname is not valid"
	}

	//typedifftypeID, typediffID, err := dataAccess.GetRecordTypeDetails(clientID, orgnID, lastlabelID, lastlabelcatID)
	typedifftypeID, typediffID, err := dataAccess.GetRecordTypeDetailsAgainstTickettype(clientID, orgnID, tz.TickettypeID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if typedifftypeID == 0 {
		return false, err, "Record type not found."
	}
	if typediffID == 0 {
		return false, err, "Record type not found."
	}

	v.CreatedusergroupID = grpID
	v.Originalusergroupid = originalgrpID

	requestorinfo, err := dataAccess.GetRequestorInfo(loginID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if len(requestorinfo.Requestername) == 0 {
		return false, err, "Something Went Wrong"
	}
	v.Requesteremail = requestorinfo.Requesteremail
	v.Requestername = requestorinfo.Requestername
	v.Requesterlocation = requestorinfo.Requesterlocation
	v.Requestermobile = requestorinfo.Requestermobile
	v.Source = "Kibana"
	v.Recordname = tz.ShortDescription
	v.Recordesc = tz.LongDescription
	rsets := []entities.RecordSet{}
	rset := entities.RecordSet{}
	values := []entities.RecordData{}
	var lastlabelID int64
	var lastlabelcatID int64
	var parentID int64
	var additionalFlag int64
	for i := 0; i < len(tz.ExternalRecordSets); i++ {
		if len(tz.ExternalRecordSets[i].Type) > 0 {
			length := (len(tz.ExternalRecordSets[0].Type) - 1)
			if len(tz.ExternalRecordSets[0].Type) < 5 {
				return false, err, "Please provide all 5 level categories."
			}
			for k := 0; k < len(tz.ExternalRecordSets[i].Type); k++ {
				v := entities.RecordData{}
				//labelID, err := dataAccess.GetLabelID(tz.ExternalRecordSets[i].Type[k].Labelname, clientID, orgnID)
				labelID, err := dataAccess.GetLabelIDAgainstTickettype(tz.ExternalRecordSets[i].Type[k].Labelname, clientID, orgnID, typedifftypeID, typediffID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if labelID == 0 {
					return false, err, "Category is missing"
				}
				if tz.ExternalRecordSets[i].Type[k].Labelvalue == "Datacenter Services" {
					additionalFlag = 1
				}
				labelvalID, err := dataAccess.GetDifferentionID(clientID, orgnID, labelID, tz.ExternalRecordSets[i].Type[k].Labelvalue, parentID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if labelvalID == 0 {
					return false, err, "Categoryyyyyy is missing"
				}
				v.ID = labelID
				v.Val = labelvalID
				values = append(values, v)
				parentID = labelvalID
				if length == k {
					lastlabelID = labelID
					lastlabelcatID = labelvalID
				}

			}
		}
	}
	rset.ID = 1
	rset.Type = values
	rsets = append(rsets, rset)

	logger.Log.Println("Last level catagory label id ::::::::", lastlabelID)
	logger.Log.Println("Last level catagory value id ::::::::", lastlabelcatID)

	// Ticket Type start here
	rset1 := entities.RecordSet{}

	rset1.ID = typedifftypeID
	rset1.Val = typediffID
	rsets = append(rsets, rset1)
	// Ticket Type end here

	//status start here
	rset2 := entities.RecordSet{}
	statusID, err := dataAccess.GetDifferentionIDbySeq(clientID, orgnID, 3, 1)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if statusID == 0 {
		return false, err, "Record status not found."
	}
	rset2.ID = 3
	rset2.Val = statusID
	rsets = append(rsets, rset2)
	// status end here

	//priority start here
	rset3 := entities.RecordSet{}
	priorityID, err := dataAccess.GetPriorityID(clientID, orgnID, lastlabelcatID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if priorityID == 0 {
		return false, err, "Please check last level category name."
	}
	rset3.ID = 5
	rset3.Val = priorityID
	rsets = append(rsets, rset3)
	//priority end here
	v.RecordSets = rsets
	additional := []entities.RecordAdditional{}
	if len(tz.ExternalAdditionalfields) < 3 && additionalFlag == 1 {
		return false, err, "All additional fields are required."
	}
	if additionalFlag == 1 {
		var flag bool
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			if tz.ExternalAdditionalfields[a].Termname == "Asset ID" {
				//return false, err, "Asset ID is missing."
				flag = true
				break
			}
		}
		if flag == false {
			return false, err, "Asset ID is missing."
		}
	}

	if additionalFlag == 1 {
		var flag bool
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			if tz.ExternalAdditionalfields[a].Termname != "Asset IP" {
				//return false, err, "Asset IP is missing."
				flag = true
				break
			}
		}
		if flag == false {
			return false, err, "Asset IP is missing."
		}
	}

	if additionalFlag == 1 {
		var flag bool
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			if tz.ExternalAdditionalfields[a].Termname != "Asset Type" {
				//return false, err, "Asset Type is missing."
				flag = true
				break
			}

		}
		if flag == false {
			return false, err, "Asset IP is missing."
		}
	}

	for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
		ad := entities.RecordAdditional{}
		if len(tz.ExternalAdditionalfields[a].Termname) == 0 {
			return false, err, "Termname is missing"
		}
		gettermid, _, err := dataAccess.GetTermID(clientID, orgnID, tz.ExternalAdditionalfields[a].Termname)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		if gettermid == 0 {
			return false, err, "This  Termname  is not valid"
		}

		additionalID, err := dataAccess.GetAdditional(clientID, orgnID, gettermid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		ad.ID = additionalID
		ad.Termsid = gettermid
		ad.Val = tz.ExternalAdditionalfields[a].Value
		additional = append(additional, ad)
	}
	v.Additionalfields = additional
	v.Workingcatlabelid = 34
	logger.Log.Println("Ticket Details ------------>", v)

	ticketID, _, err := AddRecordModelAction(&v, db, typediffID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	logger.Log.Println("ticketID ------------>", ticketID)
	return true, err, ticketID
}

func ExternalRecordCreateForSN(tz *entities.ExternalCreateRecord) (bool, error, string) {
	v := entities.RecordEntity{}
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}

	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.RequestorID) == 0 {
		return false, err, "RequestorID is missing"
	}
	if len(tz.ShortDescription) == 0 {
		return false, err, "Short description value is missing."
	}
	if len(tz.LongDescription) == 0 {
		return false, err, "Long description value is missing."
	}
	if len(tz.TickettypeID) == 0 {
		return false, err, "Ticket type is missing."
	}
	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	// New addition

	if len(tz.OriginalID) == 0 {
		return false, err, "OriginalID is missing"
	}

	if len(tz.Originalgrpname) == 0 {
		return false, err, "OriginalgrpID is missing"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}
	v.ClientID = clientID
	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}
	v.Mstorgnhirarchyid = orgnID
	loginID, err := dataAccess.GetLoginID(tz.RequestorID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if loginID == 0 {
		return false, err, "RequestorID is not valid"
	}

	originalloginID, err := dataAccess.GetLoginID(tz.OriginalID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if originalloginID == 0 {
		return false, err, "OriginalID is not valid"
	}

	v.Originaluserid = originalloginID
	v.CreateduserID = loginID
	grpID, err := dataAccess.GetGrpID(loginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	logger.Log.Println("")
	if grpID == 0 {
		return false, err, "The User does not have the permission to raise a request"
	}

	originalgrpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.Originalgrpname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	logger.Log.Println("")
	if originalgrpID == 0 {
		return false, err, "Originalgrpname is not valid"
	}

	//typedifftypeID, typediffID, err := dataAccess.GetRecordTypeDetails(clientID, orgnID, lastlabelID, lastlabelcatID)
	typedifftypeID, typediffID, err := dataAccess.GetRecordTypeDetailsAgainstTickettype(clientID, orgnID, tz.TickettypeID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if typedifftypeID == 0 {
		return false, err, "Record type not found."
	}
	if typediffID == 0 {
		return false, err, "Record type not found."
	}

	v.CreatedusergroupID = grpID
	v.Originalusergroupid = originalgrpID

	requestorinfo, err := dataAccess.GetRequestorInfo(loginID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if len(requestorinfo.Requestername) == 0 {
		return false, err, "Something Went Wrong"
	}
	v.Requesteremail = requestorinfo.Requesteremail
	v.Requestername = requestorinfo.Requestername
	v.Requesterlocation = requestorinfo.Requesterlocation
	v.Requestermobile = requestorinfo.Requestermobile
	v.Source = "SN"
	v.Recordname = tz.ShortDescription
	v.Recordesc = tz.LongDescription
	rsets := []entities.RecordSet{}
	rset := entities.RecordSet{}
	values := []entities.RecordData{}
	var lastlabelID int64
	var lastlabelcatID int64
	var parentID int64
	var additionalFlag int64
	for i := 0; i < len(tz.ExternalRecordSets); i++ {
		if len(tz.ExternalRecordSets[i].Type) > 0 {
			length := (len(tz.ExternalRecordSets[0].Type) - 1)
			if len(tz.ExternalRecordSets[0].Type) < 5 {
				return false, err, "Please provide all 5 level categories."
			}
			for k := 0; k < len(tz.ExternalRecordSets[i].Type); k++ {
				v := entities.RecordData{}
				//labelID, err := dataAccess.GetLabelID(tz.ExternalRecordSets[i].Type[k].Labelname, clientID, orgnID)
				labelID, err := dataAccess.GetLabelIDAgainstTickettype(tz.ExternalRecordSets[i].Type[k].Labelname, clientID, orgnID, typedifftypeID, typediffID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if labelID == 0 {
					return false, err, "Category is missing"
				}
				if tz.ExternalRecordSets[i].Type[k].Labelvalue == "Datacenter Services" {
					additionalFlag = 1
				}
				labelvalID, err := dataAccess.GetDifferentionID(clientID, orgnID, labelID, tz.ExternalRecordSets[i].Type[k].Labelvalue, parentID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if labelvalID == 0 {
					return false, err, "Categoryyyyyy is missing"
				}
				v.ID = labelID
				v.Val = labelvalID
				values = append(values, v)
				parentID = labelvalID
				if length == k {
					lastlabelID = labelID
					lastlabelcatID = labelvalID
				}

			}
		}
	}
	rset.ID = 1
	rset.Type = values
	rsets = append(rsets, rset)

	logger.Log.Println("Last level catagory label id ::::::::", lastlabelID)
	logger.Log.Println("Last level catagory value id ::::::::", lastlabelcatID)

	// Ticket Type start here
	rset1 := entities.RecordSet{}

	rset1.ID = typedifftypeID
	rset1.Val = typediffID
	rsets = append(rsets, rset1)
	// Ticket Type end here

	//status start here
	rset2 := entities.RecordSet{}
	statusID, err := dataAccess.GetDifferentionIDbySeq(clientID, orgnID, 3, 0)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if statusID == 0 {
		return false, err, "Record status not found."
	}
	rset2.ID = 3
	rset2.Val = statusID
	rsets = append(rsets, rset2)
	// status end here

	//priority start here
	rset3 := entities.RecordSet{}
	priorityID, err := dataAccess.GetPriorityID(clientID, orgnID, lastlabelcatID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if priorityID == 0 {
		return false, err, "Please check last level category name."
	}
	rset3.ID = 5
	rset3.Val = priorityID
	rsets = append(rsets, rset3)
	//priority end here
	v.RecordSets = rsets
	additional := []entities.RecordAdditional{}
	if len(tz.ExternalAdditionalfields) < 3 && additionalFlag == 1 {
		return false, err, "All additional fields are required."
	}
	if additionalFlag == 1 {
		var flag bool
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			if tz.ExternalAdditionalfields[a].Termname == "Asset ID" {
				//return false, err, "Asset ID is missing."
				flag = true
				break
			}
		}
		if flag == false {
			return false, err, "Asset ID is missing."
		}
	}

	if additionalFlag == 1 {
		var flag bool
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			if tz.ExternalAdditionalfields[a].Termname != "Asset IP" {
				//return false, err, "Asset IP is missing."
				flag = true
				break
			}
		}
		if flag == false {
			return false, err, "Asset IP is missing."
		}
	}

	if additionalFlag == 1 {
		var flag bool
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			if tz.ExternalAdditionalfields[a].Termname != "Asset Type" {
				//return false, err, "Asset Type is missing."
				flag = true
				break
			}

		}
		if flag == false {
			return false, err, "Asset IP is missing."
		}
	}

	for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
		ad := entities.RecordAdditional{}
		if len(tz.ExternalAdditionalfields[a].Termname) == 0 {
			return false, err, "Termname is missing"
		}
		gettermid, _, err := dataAccess.GetTermID(clientID, orgnID, tz.ExternalAdditionalfields[a].Termname)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		if gettermid == 0 {
			return false, err, "This  Termname  is not valid"
		}

		additionalID, err := dataAccess.GetAdditional(clientID, orgnID, gettermid)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		ad.ID = additionalID
		ad.Termsid = gettermid
		ad.Val = tz.ExternalAdditionalfields[a].Value
		additional = append(additional, ad)
	}
	v.Additionalfields = additional
	v.Workingcatlabelid = 34
	logger.Log.Println("Ticket Details ------------>", v)

	ticketID, _, err := AddRecordModelAction(&v, db, typediffID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	logger.Log.Println("ticketID ------------>", ticketID)
	return true, err, ticketID
}

func Gethopcount(recordID int64) (int64, bool, error, string) {
	logger.Log.Print("Gethopcount:", recordID)
	var count int64 = 0
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()

	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	log.Println("database connection failure", err)
	// 	return 0, false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection ")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	err, requestIds := dataAccess.GetRequestIdbyRecordId(recordID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if len(requestIds) > 0 {
		err, S := dataAccess.Gethopcount(requestIds[0].Requestid)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		} else {
			for i := 0; i < len(S)-1; i++ {
				if S[i] != S[i+1] {
					count = count + 1
				}
			}
			return count, true, nil, ""
		}
	} else {
		return 0, false, nil, "Ticket is not mapped with process"
	}
}

// func ExternalRecordStatusUpdate(tz *entities.ExternalCreateRecord) (bool, error, string) {
// 	db, err := dbconfig.ConnectMySqlDb()
// 	if err != nil {
// 		logger.Log.Println("database connection failure", err)
// 		return false, err, "Something Went Wrong"
// 	}
// 	defer db.Close()
// 	dataAccess := dao.DbConn{DB: db}
// 	if len(tz.Clientname) == 0 {
// 		return false, err, "Client name is missing"
// 	}
// 	if len(tz.Mstorgnhirarchyname) == 0 {
// 		return false, err, "Organization name is missing"
// 	}
// 	if len(tz.LoginID) == 0 {
// 		return false, err, "LoginID is missing"
// 	}

// 	if len(tz.Recordid) == 0 {
// 		return false, err, "RecordID is missing"
// 	}
// 	if len(tz.Statusname) == 0 {
// 		return false, err, "Statusname is missing"
// 	}

// 	if len(tz.AssigneGroupname) == 0 {
// 		return false, err, "Assignee Group name is missing"
// 	}

// 	clientID, err := dataAccess.GetClientID(tz.Clientname)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}

// 	if clientID == 0 {
// 		return false, err, "Clientname is not valid"
// 	}

// 	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}
// 	if orgnID == 0 {
// 		return false, err, "Organization name is not valid"
// 	}

// 	loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}
// 	if loginID == 0 {
// 		return false, err, "LoginID is not valid"
// 	}

// 	id, err := dataAccess.GetIdAgainstRecordNo(clientID, orgnID, tz.Recordid)
// 	if err != nil {
// 		return false, err, "Something Went Wrong"
// 	}
// 	if id == 0 {
// 		return false, err, "RecordID is not valid"
// 	}
// 	// Get TicketType
// 	typeID, err := dataAccess.GetRecordTypeID(clientID, orgnID, id)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}
// 	if typeID == 0 {
// 		return false, err, "Something Went Wrong"
// 	}
// 	//Get Ticket Type
// 	currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}

// 	nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}

// 	workingtypeID, workingcatID, err := dataAccess.GetWokinglabel(clientID, orgnID, id)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}

// 	grpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.AssigneGroupname)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}
// 	if grpID == 0 {
// 		return false, err, "LoginID is not valid"
// 	}

// 	validateassigne, err := dataAccess.ValidateAssigneUser(clientID, orgnID, grpID, loginID)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}

// 	if validateassigne == 0 {
// 		return false, err, "Current Assignee Group And Current Assignee are not mapped."
// 	}

// 	currentgrp, err := dataAccess.GetRecordCurrentGrpID(clientID, orgnID, id)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return false, err, "Something Went Wrong"
// 	}
// 	if currentgrp == 0 {
// 		logger.Log.Println("Not Found Current Support Group..........")
// 		return false, err, "Something Went Wrong"
// 	}
// 	//|| currentgrp != grpID && nxtStateID != 12
// 	//if (currentgrp != grpID && nxtStateID != 11) || (currentgrp != grpID && nxtStateID != 12) {
// 	if currentgrp != grpID && nxtStateID != 11 {
// 		return false, err, "Current Assignee Group And Change Assignee Group are not same."
// 	}

// 	reqbd := &entities.RequestBody{}
// 	reqbd.ClientID = clientID
// 	reqbd.MstorgnhirarchyID = orgnID
// 	reqbd.RecorddifftypeID = workingtypeID
// 	reqbd.RecorddiffID = workingcatID
// 	reqbd.PreviousstateID = currentStateID
// 	reqbd.CurrentstateID = nxtStateID
// 	reqbd.TransactionID = id
// 	reqbd.CreatedgroupID = grpID
// 	reqbd.MstgroupID = grpID
// 	reqbd.MstuserID = loginID
// 	//Record status request body --> &{2 2 2 4 2 4 1065 1 1 4 false 0}

// 	postBody, _ := json.Marshal(reqbd)
// 	logger.Log.Println("Record status request body -->", reqbd)
// 	responseBody := bytes.NewBuffer(postBody)
// 	resp, err := http.Post("http://localhost:8082/api/moveWorkflow", "application/json", responseBody)
// 	if err != nil {
// 		logger.Log.Println("Error is ---111111111111111111------>", err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		logger.Log.Println("Error is --22222222222222222222------->", err)
// 	}
// 	sb := string(body)
// 	wfres := entities.WorkflowResponse{}
// 	json.Unmarshal([]byte(sb), &wfres)
// 	var workflowflag = wfres.Success
// 	//var errormsg = wfres.Message
// 	logger.Log.Println("Error is --333333333333333333333333333333333------->", workflowflag)
// 	if workflowflag == true {

// 		bn := entities.RecordstatusEntity{}
// 		bn.ClientID = clientID
// 		bn.Mstorgnhirarchyid = orgnID
// 		bn.RecordID = id
// 		bn.ReordstatusID = nxtStateID
// 		bn.UserID = loginID
// 		bn.Usergroupid = grpID
// 		logger.Log.Println("Error is --99999999999999999999999999------->", bn)
// 		_, _, err, _ = Updaterecordstatus(&bn)
// 		if err != nil {
// 			logger.Log.Println("Error is --4444444444444444444444444444------->", err)
// 			return false, err, "Something Went Wrong"
// 		}

// 		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
// 			gettermid, _, err := dataAccess.GetTermID(clientID, orgnID, tz.ExternalAdditionalfields[a].Termname)
// 			if err != nil {
// 				logger.Log.Println("Error is --555555555555555555555555555------->", err)
// 				return false, err, "Something Went Wrong"
// 			}
// 			if gettermid == 0 {
// 				return false, err, "This  Termname is not valid"
// 			}
// 			stageID, err := dataAccess.GetMaxstageID(clientID, orgnID, id)
// 			if err != nil {
// 				logger.Log.Println("Error is --666666666666666666666666------->", err)
// 				return false, err, "Something Went Wrong"
// 			}
// 			if gettermid > 0 && stageID > 0 {
// 				page := entities.RecordcommonEntity{}
// 				page.ClientID = clientID
// 				page.Mstorgnhirarchyid = orgnID
// 				page.TermID = gettermid
// 				page.Termvalue = tz.ExternalAdditionalfields[a].Value
// 				page.RecordID = id
// 				page.RecordstageID = stageID
// 				page.Usergroupid = grpID
// 				page.Userid = loginID
// 				_, _, err, _ := InsertRecordTermvalues(&page)
// 				if err != nil {
// 					logger.Log.Println("Error is --77777777777777777777777------->", err)
// 					return false, err, "Something Went Wrong"
// 				}
// 			}

// 		}
// 	} else {
// 		return false, err, "Workflow does not support this changes11111111111111111."
// 	}
// 	return true, err, "Record status update successfully"
// }
func InsertExternalApiLog(clintid int64, orgid int64, userid int64, recordid int64, recorddiffid int64, statusSeq int64, requestjson string, responsejson string, db *sql.DB) error {
	dataAccess := dao.DbConn{DB: db}
	err := dataAccess.UpdateIsLatest(recordid)
	if err != nil {
		logger.Log.Println(err)
		return err
	}
	err = dataAccess.InsertExternalApiLog(clintid, orgid, userid, recordid, recorddiffid, statusSeq, requestjson, responsejson)
	if err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}
func ExternalRecordStatusUpdate(tz *entities.ExternalCreateRecord) (bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.LoginID) == 0 {
		return false, err, "LoginID is missing"
	}

	if len(tz.Recordid) == 0 {
		return false, err, "RecordID is missing"
	}
	if len(tz.Statusname) == 0 {
		return false, err, "Statusname is missing"
	}

	if len(tz.AssigneGroupname) == 0 {
		return false, err, "Assignee Group name is missing"
	}

	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}

	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}

	statusID, statusSeq, err := dataAccess.GetDifferentiationdstatustls(clientID, orgnID, tz.Statusname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	logger.Log.Println(statusID)

	loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if loginID == 0 {
		return false, err, "LoginID is not valid"
	}

	grpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.AssigneGroupname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if grpID == 0 {
		return false, err, "AssigneGroupname is not valid"
	}

	id, _, err := dataAccess.GetIdAgainstRecordNo(clientID, orgnID, tz.Recordid)
	if err != nil {
		return false, err, "Something Went Wrong"
	}
	if id == 0 {
		return false, err, "RecordID is not valid"
	}

	// Get TicketType
	typeID, err := dataAccess.GetRecordTypeID(clientID, orgnID, id)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if typeID == 0 {
		return false, err, "Something Went Wrong"
	}
	currentgrp, currentuID, err := dataAccess.GetRecordCurrentGrpID(clientID, orgnID, id)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if currentgrp == 0 {
		logger.Log.Println("Not Found Current Support Group..........")
		return false, err, "Something Went Wrong"
	}
	workingtypeID, workingcatID, err := dataAccess.GetWokinglabel(clientID, orgnID, id)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	// For New or REOPEN
	if statusSeq == 1 || statusSeq == 10 {
		return false, err, "NEW or REOPEN status change NOT allowed"
	}
	var additional []string
	if len(tz.ExternalAdditionalfields) > 0 {
		for i := 0; i < len(tz.ExternalAdditionalfields); i++ {
			additional = append(additional, tz.ExternalAdditionalfields[i].Termname)
		}
	}
	diffID, diffseq, diffname, err := dataAccess.Getcurrentsatusid(clientID, orgnID, id)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	//Other status checking
	if statusSeq != 9 {
		if loginID != currentuID {
			return false, err, "Current Assignee User and Changed Assignee User both are not same."
		}
		if grpID != currentgrp {
			return false, err, "Current Assignee Group and Changed Assignee Group both are not same."
		}

		logger.Log.Println("values is --------------------->", diffID, diffseq)
		if statusSeq == 3 {
			logger.Log.Println("################################################################################################################################", statusSeq, diffseq)

			if diffseq == 2 {
				termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
					//For log STart===============================
					resp := entities.ExternalRecordEntityResponse{}
					byte, _ := json.Marshal(tz)
					resp.Message = ""
					resp.Success = true
					//respopnse, _ := json.Marshal(resp)
					err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					//For log END==============================
					for i := 0; i < len(termnames); i++ {
						_, found := Find(additional, termnames[i])
						if !found {
							return false, err, "Termname is missing."
						}
					}
					currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, tz, 0)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					return true, err, _msg

				} else {
					return false, err, "Additional fields are missing."
				}
			} else if diffseq == 4 { //for direct pending vendor to resolved
				var ebond entities.EbondingRecordEntity
				ebond.ClientID = clientID
				ebond.MstorgnhirarchyID = orgnID
				ebond.RecordID = id
				updateTicketID, externalTicketIDErr := dataAccess.GetExternalTicketID(&ebond)
				if externalTicketIDErr != nil {
					logger.Log.Println("externalTicketIDErr =====>", externalTicketIDErr)
					return false, err, "Something Went Wrong"
				}
				if len(updateTicketID) > 0 {
					//status change from pending vendor to active
					currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, "Active")
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					_, err, _ = Updatestatuswithoutterms(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, db)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					// pending vendor to active completed
					//Now Active to resolved started
					termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
						for i := 0; i < len(termnames); i++ {
							_, found := Find(additional, termnames[i])
							if !found {
								return false, err, "Termname is missing."
							}
						}

						//For log STart===============================
						resp := entities.ExternalRecordEntityResponse{}
						byte, _ := json.Marshal(tz)
						resp.Message = ""
						resp.Success = true
						//respopnse, _ := json.Marshal(resp)
						err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}

						//For log END==============================
						currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}

						nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}
						_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, tz, 0)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}
						// //For log===============================
						// resp := entities.ExternalRecordEntityResponse{}
						// byte, _ := json.Marshal(tz)
						// resp.Message = ""
						// resp.Success = true
						// respopnse, _ := json.Marshal(resp)
						// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
						// if err != nil {
						// 	logger.Log.Println(err)
						// 	return false, err, "Something Went Wrong"
						// }
						// //For log==============================

						return true, err, _msg

					} else {
						return false, err, "Additional fields are missing."
					}

				} else {
					return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname

				}
			} else {
				return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname
			}
		} else if statusSeq == 5 {
			if diffseq == 2 {
				termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
					// for i := 0; i < len(tz.ExternalAdditionalfields); i++ {
					// 	_, found := Find(termnames, tz.ExternalAdditionalfields[i].Termname)
					// 	if !found {
					// 		return false, err, "Termname is missing."
					// 	}
					// }
					for i := 0; i < len(termnames); i++ {
						_, found := Find(additional, termnames[i])
						if !found {
							return false, err, "Termname is missing."
						}
					}

					//For log STart===============================
					resp := entities.ExternalRecordEntityResponse{}
					byte, _ := json.Marshal(tz)
					resp.Message = ""
					resp.Success = true
					//respopnse, _ := json.Marshal(resp)
					err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					//For log END==============================
					currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, tz, 0)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					// //For log===============================
					// resp := entities.ExternalRecordEntityResponse{}
					// byte, _ := json.Marshal(tz)
					// resp.Message = ""
					// resp.Success = true
					// respopnse, _ := json.Marshal(resp)
					// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
					// if err != nil {
					// 	logger.Log.Println(err)
					// 	return false, err, "Something Went Wrong"
					// }
					//For log==============================

					return true, err, _msg

				} else {
					return false, err, "Additional fields are missing."
				}
			} else {
				return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname
			}
		} else if statusSeq == 4 {
			if diffseq == 2 {
				termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
					for i := 0; i < len(termnames); i++ {
						_, found := Find(additional, termnames[i])
						if !found {
							return false, err, "Termname is missing."
						}
					}
					//For log STart===============================
					resp := entities.ExternalRecordEntityResponse{}
					byte, _ := json.Marshal(tz)
					resp.Message = ""
					resp.Success = true
					//respopnse, _ := json.Marshal(resp)
					err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					//For log END==============================
					currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, tz, 0)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					//For log===============================
					// resp := entities.ExternalRecordEntityResponse{}
					// byte, _ := json.Marshal(tz)
					// resp.Message = ""
					// resp.Success = true
					// respopnse, _ := json.Marshal(resp)
					// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
					// if err != nil {
					// 	logger.Log.Println(err)
					// 	return false, err, "Something Went Wrong"
					// }
					//For log==============================

					return true, err, _msg

				} else {
					return false, err, "Additional fields are missing."
				}
			} else {
				return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname
			}
		} else if statusSeq == 2 {
			if diffseq == 9 || diffseq == 4 {
				// termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID)
				// if err != nil {
				// 	logger.Log.Println(err)
				// 	return false, err, "Something Went Wrong"
				// }
				//For log STart===============================
				resp := entities.ExternalRecordEntityResponse{}
				byte, _ := json.Marshal(tz)
				resp.Message = ""
				resp.Success = true
				//respopnse, _ := json.Marshal(resp)
				err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}

				//For log END==============================
				currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}

				nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				_, err, _msg := Updatestatuswithoutterms(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, db)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				// //For log===============================
				// resp := entities.ExternalRecordEntityResponse{}
				// byte, _ := json.Marshal(tz)
				// resp.Message = ""
				// resp.Success = true
				// respopnse, _ := json.Marshal(resp)
				// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
				// if err != nil {
				// 	logger.Log.Println(err)
				// 	return false, err, "Something Went Wrong"
				// }
				//For log==============================

				return true, err, _msg

			} else {
				return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname
			}
		} else if statusSeq == 10 {
			if diffseq == 3 {
				termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
					for i := 0; i < len(termnames); i++ {
						_, found := Find(additional, termnames[i])
						if !found {
							return false, err, "Termname is missing."
						}
					}

					//For log STart===============================
					resp := entities.ExternalRecordEntityResponse{}
					byte, _ := json.Marshal(tz)
					resp.Message = ""
					resp.Success = true
					//respopnse, _ := json.Marshal(resp)
					err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					//For log END==============================
					currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, tz, 0)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					// //For log===============================
					// resp := entities.ExternalRecordEntityResponse{}
					// byte, _ := json.Marshal(tz)
					// resp.Message = ""
					// resp.Success = true
					// respopnse, _ := json.Marshal(resp)
					// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
					// if err != nil {
					// 	logger.Log.Println(err)
					// 	return false, err, "Something Went Wrong"
					// }
					//For log==============================

					return true, err, _msg

				} else {
					return false, err, "Additional fields are missing."
				}
			} else {
				return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname
			}
		} else if statusSeq == 8 {
			if diffseq == 3 {
				termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
				logger.Log.Println("Closed term names ----------------------------->", termnames)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
					for i := 0; i < len(termnames); i++ {
						_, found := Find(additional, termnames[i])
						if !found {
							return false, err, "Termname is missing."
						}
					}
					//For log STart===============================
					resp := entities.ExternalRecordEntityResponse{}
					byte, _ := json.Marshal(tz)
					resp.Message = ""
					resp.Success = true
					//respopnse, _ := json.Marshal(resp)
					err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					//For log END==============================
					currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, tz, 0)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					// //For log===============================
					// resp := entities.ExternalRecordEntityResponse{}
					// byte, _ := json.Marshal(tz)
					// resp.Message = ""
					// resp.Success = true
					// respopnse, _ := json.Marshal(resp)
					// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
					// if err != nil {
					// 	logger.Log.Println(err)
					// 	return false, err, "Something Went Wrong"
					// }
					//For log==============================

					return true, err, _msg

				} else {
					return false, err, "Additional fields are missing."
				}
			} else {
				return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname
			}
		} else if statusSeq == 11 {
			if diffseq == 2 {
				logingrpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.LoginGrpname)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if logingrpID == 0 {
					return false, err, "LoginGroupname is not valid"
				}

				termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
				logger.Log.Println("Closed term names ----------------------------->", termnames)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
					for i := 0; i < len(termnames); i++ {
						_, found := Find(additional, termnames[i])
						if !found {
							return false, err, "Termname is missing."
						}
					}
					//For log STart===============================
					resp := entities.ExternalRecordEntityResponse{}
					byte, _ := json.Marshal(tz)
					resp.Message = ""
					resp.Success = true
					//respopnse, _ := json.Marshal(resp)
					err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					//For log END==============================

					currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}

					nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, logingrpID, loginID, tz, 1)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					//For log===============================
					// resp := entities.ExternalRecordEntityResponse{}
					// byte, _ := json.Marshal(tz)
					// resp.Message = ""
					// resp.Success = true
					// respopnse, _ := json.Marshal(resp)
					// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
					// if err != nil {
					// 	logger.Log.Println(err)
					// 	return false, err, "Something Went Wrong"
					// }
					//For log==============================

					return true, err, _msg

				} else {
					return false, err, "Additional fields are missing."
				}
			} else {
				var ebond entities.EbondingRecordEntity
				ebond.ClientID = clientID
				ebond.MstorgnhirarchyID = orgnID
				ebond.RecordID = id
				updateTicketID, externalTicketIDErr := dataAccess.GetExternalTicketID(&ebond)
				if externalTicketIDErr != nil {
					logger.Log.Println("externalTicketIDErr =====>", externalTicketIDErr)
					return false, err, "Something Went Wrong"
				}
				if len(updateTicketID) > 0 {
					logingrpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.LoginGrpname)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					if logingrpID == 0 {
						return false, err, "LoginGroupname is not valid"
					}

					termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
					logger.Log.Println("Closed term names ----------------------------->", termnames)
					if err != nil {
						logger.Log.Println(err)
						return false, err, "Something Went Wrong"
					}
					if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
						for i := 0; i < len(termnames); i++ {
							_, found := Find(additional, termnames[i])
							if !found {
								return false, err, "Termname is missing."
							}
						}
						//For log STart===============================
						resp := entities.ExternalRecordEntityResponse{}
						byte, _ := json.Marshal(tz)
						resp.Message = ""
						resp.Success = true
						//respopnse, _ := json.Marshal(resp)
						err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}

						//For log END==============================
						currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}

						nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}
						_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, logingrpID, loginID, tz, 1)
						if err != nil {
							logger.Log.Println(err)
							return false, err, "Something Went Wrong"
						}
						//For log===============================
						// resp := entities.ExternalRecordEntityResponse{}
						// byte, _ := json.Marshal(tz)
						// resp.Message = ""
						// resp.Success = true
						// respopnse, _ := json.Marshal(resp)
						// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
						// if err != nil {
						// 	logger.Log.Println(err)
						// 	return false, err, "Something Went Wrong"
						// }
						//For log==============================

						return true, err, _msg

					} else {
						return false, err, "Additional fields are missing."
					}

				} else {
					return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname

				}
			}
		} else {
			return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname

		}
	} else {
		if diffseq == 5 {
			termnames, err := dataAccess.GetTermnamesbystatusID(clientID, orgnID, statusID, typeID)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			if len(tz.ExternalAdditionalfields) > 0 && len(termnames) > 0 {
				for i := 0; i < len(termnames); i++ {
					_, found := Find(additional, termnames[i])
					if !found {
						return false, err, "Termname is missing."
					}
				}
				//For log STart===============================
				resp := entities.ExternalRecordEntityResponse{}
				byte, _ := json.Marshal(tz)
				resp.Message = ""
				resp.Success = true
				//respopnse, _ := json.Marshal(resp)
				err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), "", db)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}

				//For log END==============================
				currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}

				nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, tz.Statusname)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				_, err, _msg := Updatestatus(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, tz, 0)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
				//For log===============================
				// resp := entities.ExternalRecordEntityResponse{}
				// byte, _ := json.Marshal(tz)
				// resp.Message = ""
				// resp.Success = true
				// respopnse, _ := json.Marshal(resp)
				// err = InsertExternalApiLog(clientID, orgnID, loginID, id, typeID, statusSeq, string(byte), string(respopnse), db)
				// if err != nil {
				// 	logger.Log.Println(err)
				// 	return false, err, "Something Went Wrong"
				// }
				//For log==============================

				return true, err, _msg

			} else {
				return false, err, "Additional fields are missing."
			}
		} else {
			return false, err, "The Requested status change is not allowed for " + diffname + " to " + tz.Statusname
		}
	}

	//return true, err, "Record status update successfully"

}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func Updatestatus(clientID int64, orgnID int64, workingtypeID int64, workingcatID int64, currentStateID int64, nxtStateID int64, id int64, grpID int64, loginID int64, tz *entities.ExternalCreateRecord, manualVal int64) (bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	reqbd := &entities.RequestBody{}
	reqbd.ClientID = clientID
	reqbd.MstorgnhirarchyID = orgnID
	reqbd.RecorddifftypeID = workingtypeID
	reqbd.RecorddiffID = workingcatID
	reqbd.PreviousstateID = currentStateID
	reqbd.CurrentstateID = nxtStateID
	reqbd.TransactionID = id
	reqbd.CreatedgroupID = grpID
	reqbd.MstgroupID = grpID
	reqbd.MstuserID = loginID
	reqbd.UserID = loginID
	reqbd.Manualstateselection = manualVal
	//Record status request body --> &{2 2 2 4 2 4 1065 1 1 4 false 0}

	postBody, _ := json.Marshal(reqbd)
	logger.Log.Println("Record status request body -->", reqbd)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(dbconfig.MASTER_URL+"/moveWorkflow", "application/json", responseBody)
	if err != nil {
		logger.Log.Println("Error is ---111111111111111111------>", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Println("Error is --22222222222222222222------->", err)
	}
	sb := string(body)
	wfres := entities.WorkflowResponse{}
	json.Unmarshal([]byte(sb), &wfres)
	var workflowflag = wfres.Success
	//var errormsg = wfres.Message
	logger.Log.Println("Error is --333333333333333333333333333333333------->", workflowflag)
	if workflowflag == true {

		bn := entities.RecordstatusEntity{}
		bn.ClientID = clientID
		bn.Mstorgnhirarchyid = orgnID
		bn.RecordID = id
		bn.ReordstatusID = nxtStateID
		bn.UserID = loginID
		bn.Usergroupid = grpID
		logger.Log.Println("Error is --99999999999999999999999999------->", bn)
		//	_, _, err, _ = Updaterecordstatus(&bn, db)
		//		if err != nil {
		//			logger.Log.Println("Error is --4444444444444444444444444444------->", err)
		//			return false, err, "Something Went Wrong"
		//		}

		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			gettermid, _, err := dataAccess.GetTermID(clientID, orgnID, tz.ExternalAdditionalfields[a].Termname)
			if err != nil {
				logger.Log.Println("Error is --555555555555555555555555555------->", err)
				return false, err, "Something Went Wrong"
			}
			if gettermid == 0 {
				return false, err, "This  Termname is not valid"
			}
			stageID, err := dataAccess.GetMaxstageID(clientID, orgnID, id)
			if err != nil {
				logger.Log.Println("Error is --666666666666666666666666------->", err)
				return false, err, "Something Went Wrong"
			}
			if gettermid > 0 && stageID > 0 {
				page := entities.RecordcommonEntity{}
				page.ClientID = clientID
				page.Mstorgnhirarchyid = orgnID
				page.TermID = gettermid
				page.Termvalue = tz.ExternalAdditionalfields[a].Value
				page.RecordID = id
				page.RecordstageID = stageID
				page.Usergroupid = grpID
				page.Userid = loginID
				_, _, err, _ := InsertRecordTermvalues(&page)
				if err != nil {
					logger.Log.Println("Error is --77777777777777777777777------->", err)
					return false, err, "Something Went Wrong"
				}
			}

		}
	} else {
		return false, err, "Workflow does not support this changes"
	}
	return true, err, "Record status updated successfully"

}

func Updatestatuswithoutterms(clientID int64, orgnID int64, workingtypeID int64, workingcatID int64, currentStateID int64, nxtStateID int64, id int64, grpID int64, loginID int64, db *sql.DB) (bool, error, string) {

	reqbd := &entities.RequestBody{}
	reqbd.ClientID = clientID
	reqbd.MstorgnhirarchyID = orgnID
	reqbd.RecorddifftypeID = workingtypeID
	reqbd.RecorddiffID = workingcatID
	reqbd.PreviousstateID = currentStateID
	reqbd.CurrentstateID = nxtStateID
	reqbd.TransactionID = id
	reqbd.CreatedgroupID = grpID
	reqbd.MstgroupID = grpID
	reqbd.MstuserID = loginID
	reqbd.UserID = loginID

	//Record status request body --> &{2 2 2 4 2 4 1065 1 1 4 false 0}

	postBody, _ := json.Marshal(reqbd)
	logger.Log.Println("Record status request body -->", reqbd)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(dbconfig.MASTER_URL+"/moveWorkflow", "application/json", responseBody)
	if err != nil {
		logger.Log.Println("Error is ---111111111111111111------>", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Println("Error is --22222222222222222222------->", err)
	}
	sb := string(body)
	wfres := entities.WorkflowResponse{}
	json.Unmarshal([]byte(sb), &wfres)
	var workflowflag = wfres.Success
	//var errormsg = wfres.Message
	logger.Log.Println("Error is --333333333333333333333333333333333------->", workflowflag)
	if workflowflag == true {

		bn := entities.RecordstatusEntity{}
		bn.ClientID = clientID
		bn.Mstorgnhirarchyid = orgnID
		bn.RecordID = id
		bn.ReordstatusID = nxtStateID
		bn.UserID = loginID
		bn.Usergroupid = grpID
		//	logger.Log.Println("Error is --99999999999999999999999999------->", bn)
		//	_, _, err, _ = Updaterecordstatus(&bn, db)
		//	if err != nil {
		//		logger.Log.Println("Error is --4444444444444444444444444444------->", err)
		//		return false, err, "Something Went Wrong"
		//	}

	} else {
		return false, err, "Workflow does not support this changes"
	}
	return true, err, "Record status updated successfully"

}

func ExternalRecordGrpUpdate(tz *entities.ExternalCreateRecord) (bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.LoginID) == 0 {
		return false, err, "LoginID is missing"
	}

	if len(tz.Recordid) == 0 {
		return false, err, "RecordID is missing"
	}
	if len(tz.LoginGrpname) == 0 {
		return false, err, "Login Group name is missing"
	}
	if len(tz.AssigneGroupname) == 0 {
		return false, err, "Work Flow doesnt support this change"
	}
	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}

	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}
	loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if loginID == 0 {
		return false, err, "LoginID is not valid"
	}

	//grpID, err := dataAccess.GetGrpID(loginID, clientID, orgnID)
	grpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.LoginGrpname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if grpID == 0 {
		return false, err, "LoginID is not valid"
	}
	id, _, err := dataAccess.GetIdAgainstRecordNo(clientID, orgnID, tz.Recordid)
	if err != nil {
		return false, err, "RecordID is not valid"
	}
	if id == 0 {
		return false, err, "RecordID is not valid"
	}

	_, currentstatusseq, _, err := dataAccess.Getcurrentsatusid(clientID, orgnID, id)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if currentstatusseq != 2 {
		return false, err, "Group changed is not possible in this status."
	}

	changedgrpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.AssigneGroupname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if changedgrpID == 0 {
		return false, err, "Group name is not valid."
	}

	if grpID == changedgrpID {
		return false, err, "Current Group And Changed Group both are same."
	}

	reqbd := &entities.RequestBody1{}
	reqbd.TransactionID = id
	reqbd.CreatedgroupID = grpID
	reqbd.MstgroupID = changedgrpID
	reqbd.MstuserID = 0
	reqbd.Samegroup = false
	reqbd.UserID = loginID
	postBody, _ := json.Marshal(reqbd)
	logger.Log.Println("Record status request body -->", reqbd)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(dbconfig.MASTER_URL+"/changerecordgroup", "application/json", responseBody)
	if err != nil {
		logger.Log.Println("Error is --------->", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Println("Error is --------->", err)
	}
	sb := string(body)
	wfres := entities.WorkflowResponse{}
	json.Unmarshal([]byte(sb), &wfres)
	var workflowflag = wfres.Success
	if workflowflag == false {
		return false, err, "Something Went Wrong"
	} else {
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			gettermid, _, err := dataAccess.GetTermID(clientID, orgnID, tz.ExternalAdditionalfields[a].Termname)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			if gettermid == 0 {
				return false, err, "This  Termname is not valid"
			}

			//
			stageID, err := dataAccess.GetMaxstageID(clientID, orgnID, id)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			if gettermid > 0 && stageID > 0 {
				page := entities.RecordcommonEntity{}
				page.ClientID = clientID
				page.Mstorgnhirarchyid = orgnID
				page.TermID = gettermid
				page.Termvalue = tz.ExternalAdditionalfields[a].Value
				page.RecordID = id
				page.RecordstageID = stageID
				page.Usergroupid = grpID
				page.Userid = loginID
				_, _, err, _ := InsertRecordTermvalues(&page)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
			}

		}
	}

	return true, err, "Support Group changed successfully"
}

func ExternalRecordUserUpdate(tz *entities.ExternalCreateRecord) (bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.LoginID) == 0 {
		return false, err, "LoginID is missing"
	}

	if len(tz.Recordid) == 0 {
		return false, err, "RecordID is missing"
	}
	if len(tz.AssigneGroupname) == 0 {
		return false, err, "Group name is missing"
	}
	if len(tz.LoginGrpname) == 0 {
		return false, err, "Login Group name is missing"
	}
	if len(tz.AssigneUserID) == 0 {
		return false, err, "Change LoginID is missing"
	}

	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}

	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}
	loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if loginID == 0 {
		return false, err, "LoginID is not valid"
	}

	//grpID, err := dataAccess.GetGrpID(loginID, clientID, orgnID)
	grpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.LoginGrpname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if grpID == 0 {
		return false, err, "Loginsupportgroup is not valid"
	}
	id, _, err := dataAccess.GetIdAgainstRecordNo(clientID, orgnID, tz.Recordid)
	if err != nil {
		return false, err, "RecordID is not valid"
	}
	if id == 0 {
		return false, err, "RecordID is not valid"
	}

	changedgrpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.AssigneGroupname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if changedgrpID == 0 {
		return false, err, "Group name is not valid."
	}
	currentgrp, _, err := dataAccess.GetRecordCurrentGrpID(clientID, orgnID, id)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if currentgrp == 0 {
		logger.Log.Println("Not Found Current Support Group..........")
		return false, err, "Something Went Wrong"
	}
	if currentgrp != changedgrpID {
		return false, err, "Current Group And Changed Group both are not same."
	}

	changeloginID, err := dataAccess.GetLoginID(tz.AssigneUserID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if changeloginID == 0 {
		return false, err, "Change LoginID is not valid"
	}

	validateassigne, err := dataAccess.ValidateAssigneUser(clientID, orgnID, changedgrpID, changeloginID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	if validateassigne == 0 {
		return false, err, "Current Assignee Group And Current Assignee are not mapped."
	}

	reqbd := &entities.RequestBody1{}
	reqbd.TransactionID = id
	reqbd.CreatedgroupID = grpID
	reqbd.MstgroupID = changedgrpID
	reqbd.MstuserID = changeloginID
	reqbd.Samegroup = true
	reqbd.UserID = loginID
	postBody, _ := json.Marshal(reqbd)
	logger.Log.Println("Record status request body -->", reqbd)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(dbconfig.MASTER_URL+"/changerecordgroup", "application/json", responseBody)
	if err != nil {
		logger.Log.Println("Error is --------->", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Println("Error is --------->", err)
	}
	sb := string(body)
	wfres := entities.WorkflowResponse{}
	json.Unmarshal([]byte(sb), &wfres)
	var workflowflag = wfres.Success
	if workflowflag == false {
		return false, err, "Something Went Wrong"
	} else {
		for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
			gettermid, _, err := dataAccess.GetTermID(clientID, orgnID, tz.ExternalAdditionalfields[a].Termname)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			if gettermid == 0 {
				return false, err, "This  Termname is not valid"
			}
			stageID, err := dataAccess.GetMaxstageID(clientID, orgnID, id)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			if gettermid > 0 && stageID > 0 {
				page := entities.RecordcommonEntity{}
				page.ClientID = clientID
				page.Mstorgnhirarchyid = orgnID
				page.TermID = gettermid
				page.Termvalue = tz.ExternalAdditionalfields[a].Value
				page.RecordID = id
				page.RecordstageID = stageID
				page.Usergroupid = grpID
				page.Userid = loginID
				_, _, err, _ := InsertRecordTermvalues(&page)
				if err != nil {
					logger.Log.Println(err)
					return false, err, "Something Went Wrong"
				}
			}

		} // for additional field end here

		// new changes done here
		// Get TicketType

		diffID, diffseq, diffname, err := dataAccess.Getcurrentsatusid(clientID, orgnID, id)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		logger.Log.Println("In side user update ----------->", diffID, diffname)
		if diffseq == 1 || diffseq == 9 || diffseq == 10 {
			typeID, err := dataAccess.GetRecordTypeID(clientID, orgnID, id)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			if typeID == 0 {
				return false, err, "Something Went Wrong"
			}
			workingtypeID, workingcatID, err := dataAccess.GetWokinglabel(clientID, orgnID, id)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}

			currentStateID, err := dataAccess.GetrecordlateststateID(id, clientID, orgnID)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}

			nxtStateID, err := dataAccess.GetrecordnxtstateID(clientID, orgnID, "Active")
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			_, err, _msg := Updatestatuswithoutterms(clientID, orgnID, workingtypeID, workingcatID, currentStateID, nxtStateID, id, grpID, loginID, db)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
			return true, err, _msg
		}

		// ---------------------
	}

	return true, err, "Assigne User changed successfully"
}

func ExternalRecordInternalCommentUpdate(tz *entities.ExternalCreateRecord) (bool, error, string) {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("database connection failure", err)
	// 	return false, err, "Something Went Wrong"
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	if len(tz.Clientname) == 0 {
		return false, err, "Client name is missing"
	}
	if len(tz.Mstorgnhirarchyname) == 0 {
		return false, err, "Organization name is missing"
	}
	if len(tz.LoginID) == 0 {
		return false, err, "LoginID is missing"
	}

	if len(tz.Recordid) == 0 {
		return false, err, "RecordID is missing"
	}
	if len(tz.LoginGrpname) == 0 {
		return false, err, "Login Group name is missing"
	}

	if len(tz.ExternalAdditionalfields) == 0 {
		return false, err, "Internal comment is missing"
	}

	clientID, err := dataAccess.GetClientID(tz.Clientname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}

	if clientID == 0 {
		return false, err, "Clientname is not valid"
	}

	orgnID, err := dataAccess.GetOrgnID(tz.Mstorgnhirarchyname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if orgnID == 0 {
		return false, err, "Organization name is not valid"
	}
	loginID, err := dataAccess.GetLoginID(tz.LoginID, clientID, orgnID)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if loginID == 0 {
		return false, err, "LoginID is not valid"
	}

	//grpID, err := dataAccess.GetGrpID(loginID, clientID, orgnID)
	grpID, err := dataAccess.GetGrpIDByName(clientID, orgnID, tz.LoginGrpname)
	if err != nil {
		logger.Log.Println(err)
		return false, err, "Something Went Wrong"
	}
	if grpID == 0 {
		return false, err, "LoginID is not valid"
	}

	id, _, err := dataAccess.GetIdAgainstRecordNo(clientID, orgnID, tz.Recordid)
	if err != nil {
		return false, err, "RecordID is not valid"
	}
	if id == 0 {
		return false, err, "RecordID is not valid"
	}

	for a := 0; a < len(tz.ExternalAdditionalfields); a++ {
		logger.Log.Println("Termname is --------------------------------->", tz.ExternalAdditionalfields[a].Termname)
		gettermid, seq, err := dataAccess.GetTermID(clientID, orgnID, tz.ExternalAdditionalfields[a].Termname)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		/*              if seq != 22 || seq != 11 {//seq 11 is added to update 'Work Notes'

			return false, err, "This  Termname is not valid"
		}*/
		if seq != 22 {
			if seq != 11 {

				return false, err, "This  Termname is not valid"
			}
		}
		if gettermid == 0 {
			return false, err, "This  Termname is not valid"
		}
		stageID, err := dataAccess.GetMaxstageID(clientID, orgnID, id)
		if err != nil {
			logger.Log.Println(err)
			return false, err, "Something Went Wrong"
		}
		logger.Log.Println("Termname is --------------gettermid------------------->", gettermid)
		logger.Log.Println("Termname is --------------stageID------------------->", stageID)
		if gettermid > 0 && stageID > 0 {
			page := entities.RecordcommonEntity{}
			page.ClientID = clientID
			page.Mstorgnhirarchyid = orgnID
			page.TermID = gettermid
			page.Termvalue = tz.ExternalAdditionalfields[a].Value
			page.RecordID = id
			page.RecordstageID = stageID
			page.Usergroupid = grpID
			page.Userid = loginID
			_, _, err, _ := InsertRecordTermvalues(&page)
			if err != nil {
				logger.Log.Println(err)
				return false, err, "Something Went Wrong"
			}
		}

	}

	return true, err, "Comment added successfully"
}
func Getutctime(clientid int64, mstorgnhierarchyid int64, datetime string, Timediff int64) (string, error) {

	// t, err1 := dataAccess.Gettimediff(clientID, orgnID)
	// if err1 != nil {
	// 	return t, false, err1, "Something Went Wrong"
	// }
	layout := "2006-01-02 15:04:05"
	parsetime, err := time.Parse(layout, datetime)
	if err != nil {
		logger.Log.Println("parsetime error:", err)
		return "", err

	}
	logger.Log.Println("Parsetime is", parsetime)
	unixtime := parsetime.Unix()
	logger.Log.Println("unixtime:", unixtime, Timediff)
	// time := dao.Convertdate(int64(parsetime.Unix()), Timediff)
	// logger.Log.Println("Time before:" + datetime + "   Time Now:" + time)
	unixTime := int64(parsetime.Unix()) - Timediff
	logger.Log.Println("unixtime:", unixTime)

	t := time.Unix(unixTime, 0)
	// return t.Format("02-Jan-2006 15:04:05"), nil
	return t.Format(layout), nil

}
func InsertRecordTermvaluesforAttachment(tz *entities.RecordcommonEntity) (int64, bool, error, string) {
	logger.Log.Println("In side Mstslastatemodel")
	// if mutexutility.MutexLocked(lock) == false {
	//  lock.Lock()
	//  defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	//  log.Println("database connection failure", err)
	//  return 0, false, err, "Something Went Wrong"
	// }
	//defer db.Close()

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return 0, false, err, "Something Went Wrong"
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	//childids, err := dataAccess.Getchildrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	childids, err := dataAccess.GetchildrecordidsForCommon(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	//logger.Log.Println("Child ticket size is ----------->", len(childids))
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	hashmap, err := dataAccess.Termsequance(tz.ClientID, tz.Mstorgnhirarchyid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	if len(hashmap) == 0 {
		return 0, false, err, "Something Went Wrong"
	}
	if tz.TermID == 0 {
		id := hashmap[tz.Termseq]
		tz.TermID = id
	}
	//logger.Log.Println("Termid value is ==============================>", tz.TermID)
	//logger.Log.Println("Term seuance value is ==============================>", tz.Termseq)

	releationID, err := dataAccess.Checktermreleation(tz.TermID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
	logger.Log.Println("Child ticket size is ----------->", releationID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	termnm, err := dataAccess.Gettermnamebyid(tz.TermID, tz.ClientID, tz.Mstorgnhirarchyid)
	//logger.Log.Println("Term name is  -----2222222222222222222222222222222------>", termnm)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	grplevel, err := dataAccess.GetgrplevelID(tz.ClientID, tz.Mstorgnhirarchyid, tz.Usergroupid)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	//logger.Log.Println("Term id  is in InsertRecordTermvalues ----------->", tz.TermID)

	id, err := dataAccess.InsertRecordTermvalues(tz)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}
	err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	if tz.Termseq == 29 {
		fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if fcount > 0 {
			//Update Stage TBL For Impact
			err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, 17, "The value is :: "+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

		}

	}

	Username, err := dataAccess.GetUsername(tz.Userid)
	if err != nil {
		logger.Log.Println(err)
		return 0, false, err, "Something Went Wrong"
	}

	err = dataAccess.UpdateUserInfoWithoutTrn(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
	if err != nil {
		logger.Log.Println(err)
		//return 0, false, err, "Something Went Wrong"
	}
	//logger.Log.Println("childids values is  ----------->", childids)
	//logger.Log.Println("releationID values is  ----------->", releationID)
	//logger.Log.Println("grplevel values is  ----------->", grplevel)

	if len(childids) > 0 && releationID > 0 && grplevel > 1 {
		for k := 0; k < len(childids); k++ {
			//  logger.Log.Println("222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222222  ----------->")
			_, err := dataAccess.InsertRecordTermvaluesForchilds(tz, childids[k])
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

			err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

			if tz.Termseq == 11 {
				go CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
			}
			if tz.Termseq == 1 {
				go FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
			}
			err = dataAccess.UpdateUserInfoWithoutTrn(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
			if err != nil {
				logger.Log.Println(err)
				//return 0, false, err, "Something Went Wrong"
			}

			if tz.Termseq == 29 {
				fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
				if fcount > 0 {
					err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 17, "The value is :: "+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}

				}

			}

		}
	} // child for loop

	// For Stask ticket type ......
	typeseq, err := dataAccess.Getrecordtypeseq(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
	if err != nil {
		return 0, false, err, "Something Went Wrong"
	}

	// For Stask ticket type ......
	if len(childids) > 0 && releationID > 0 && grplevel == 0 { // For SR ticket type...
		if typeseq == 2 || typeseq == 4 {
			for k := 0; k < len(childids); k++ {
				//  logger.Log.Println("333333333333333333333333333333333333333333333333333333333333333333333333333333333333  ----------->", typeseq)
				_, err := dataAccess.InsertRecordTermvaluesForchilds(tz, childids[k])
				if err != nil {
					return 0, false, err, "Something Went Wrong"
				}

				err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
				if err != nil {
					return 0, false, err, "Something Went Wrong"
				}

				if tz.Termseq == 29 {
					fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
					if fcount > 0 {
						err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, childids[k], 17, "The value is :: "+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
						if err != nil {
							return 0, false, err, "Something Went Wrong"
						}

					}

				}

				if tz.Termseq == 11 {
					go CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
				}
				if tz.Termseq == 1 {
					go FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, childids[k])
				}
				err = dataAccess.UpdateUserInfoWithoutTrn(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, tz.Userid, Username)
				if err != nil {
					logger.Log.Println(err)
					//return 0, false, err, "Something Went Wrong"
				}

			}
		}

	} // child for loop

	if typeseq == 3 || typeseq == 5 {
		_, Seq, _, err := dataAccess.Getcurrentsatusid(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		//3,8,11,14
		if Seq != 3 && Seq != 8 && Seq != 11 && Seq != 14 {
			parentid, err := dataAccess.Getparentrecordids(tz.RecordID, tz.ClientID, tz.Mstorgnhirarchyid, tz.Recorddifftypeid, tz.Recorddiffid)
			//logger.Log.Println("Child ticket size is ----------->", len(parentid))
			if err != nil {
				return 0, false, err, "Something Went Wrong"
			}

			if len(parentid) > 0 && releationID > 0 {
				for k := 0; k < len(parentid); k++ {
					_, err := dataAccess.InsertRecordTermvaluesForchilds(tz, parentid[k])
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}

					err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], 100, termnm+" :: "+tz.Termvalue, tz.Userid, tz.Usergroupid, tz.TermID)
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}

					if tz.Termseq == 29 {
						fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k])
						if fcount > 0 {
							err = dataAccess.InsertActivityLogs(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k], 17, "The value is ::"+strconv.FormatInt(fcount, 10), tz.Userid, tz.Usergroupid, 0)
							if err != nil {
								return 0, false, err, "Something Went Wrong"
							}

						}

					}

					if tz.Termseq == 11 {
						go CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k])
					}
					if tz.Termseq == 1 {
						go FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, parentid[k])
					}

				}
			} // parent for loop
		} // Seq if condition end here...
	} // typeseq if condition
	// For Stask ticket type ............

	if id > 0 {
		if tz.Termseq == 11 {
			CustomerVisibleWorkNotesEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		}

		if tz.Termseq == 1 {
			FileAttachmentEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
		}
		if tz.Termseq == 29 {
			fcount, _ := dataAccess.GetFollowupcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if fcount > 0 {
				//Update Stage TBL For Impact

				err = dataAccess.UpdateFollowupcountFromCommon(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, fcount)
				if err != nil {
					logger.Log.Println("error is ----->", err)
				}

				//Update Stage TBL For Impact
				go FollowupCountEmail(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, fcount)
			}

		}

		if tz.Termseq == 30 {
			obcount, err := dataAccess.GetOutboundcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID)
			if err != nil {
				logger.Log.Println("error is ----->", err)
			}
			err = dataAccess.UpdateOutboundcount(tz.ClientID, tz.Mstorgnhirarchyid, tz.RecordID, obcount)
			if err != nil {
				logger.Log.Println("error is ----->", err)
			}
		}
		return id, true, err, ""
	} else {
		return 0, false, err, "Something Went Wrong"
	}

}
