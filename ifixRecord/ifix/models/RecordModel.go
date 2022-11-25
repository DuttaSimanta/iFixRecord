package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/dbconfig"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/mutexutility"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

func CreateRecordModel(data *entities.RecordEntity) (string, int64, error) {
	var recordset = data.RecordSets
	var lastcatid int64 = 0
	for i := 0; i < len(recordset); i++ {
		if len(recordset[i].Type) > 0 && recordset[i].ID == 1 {
			var catset = recordset[i].Type
			if len(catset) > 0 {
				lastcatid = catset[len(catset)-1].Val
			}
		}
	}
	logger.Log.Println("Mutex flag value is ---------------------------->", mutexutility.MutexLocked(lock))
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }

	// dbcon, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	logger.Log.Println("Error in DBConnection in side CreateRecordModel")
	// 	return "", 0, errors.New("Something Went Wrong")
	// }

	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return "", 0, errors.New("Something Went Wrong")
		}
		db = dbcon
	}

	dataAccess := dao.DbConn{DB: db}
	var hastask bool
	taskcatids, terr := dataAccess.Gettaskbycatid(lastcatid)
	log.Print("taskids", taskcatids)
	if terr != nil {
		return "", 0, errors.New("Something Went Wrong")
	}
	if len(taskcatids) > 0 {
		hastask = true
	} else {
		hastask = false
	}
	var parentdifftypeid int64
	logger.Log.Println(parentdifftypeid)
	var parentdiffid int64
	var parentpriodifftypeid int64
	var parentpriodiffid int64
	for i := 0; i < len(recordset); i++ {
		if recordset[i].ID == 2 {
			parentdifftypeid = recordset[i].ID
			parentdiffid = recordset[i].Val
		}
		if recordset[i].ID == 5 {
			parentpriodifftypeid = recordset[i].ID
			parentpriodiffid = recordset[i].Val
		}
	}
	seq, seqerr := dataAccess.Getseqbyid(data.ClientID, data.Mstorgnhirarchyid, parentdiffid)
	if seqerr != nil {
		return "", 0, errors.New("Something Went Wrong")
	}
	if seq == 2 && len(taskcatids) == 0 {
		return "", 0, errors.New("No task mapped with SR ticket type")
	}
	ticketID, insertedID, modelResponseError := AddRecordModelAction(data, db, parentdiffid)
	if modelResponseError != nil {
		return "", 0, errors.New("Something went wrong")
	}

	if hastask {
		/**
		Checking whether the parent record has approval process or non approval process
		*/
		isApprovedProcess := false
		for i := 0; i < len(recordset); i++ {
			if len(recordset[i].Type) > 0 {
				for j := 0; j < len(recordset[i].Type); j++ {
					if recordset[i].Type[j].ID == data.Workingcatlabelid {
						reccatentity := entities.RecordcategoryupdateEntity{}
						reccatentity.ClientID = data.ClientID
						reccatentity.Mstorgnhirarchyid = data.Mstorgnhirarchyid
						reccatentity.Recorddiffid = recordset[i].Type[j].Val
						reccatentity.Recorddifftypeid = recordset[i].Type[j].ID
						pids, err := dataAccess.GetProcessByCategory(&reccatentity)
						if err != nil {
							return "", 0, errors.New("Something Went Wrong")
						}
						count, err := dataAccess.Checkisapprovedprocess(&reccatentity, pids.ID)
						if err != nil {
							return "", 0, errors.New("Something Went Wrong")
						}
						if count.Total > 0 {
							isApprovedProcess = true
						}
						log.Print("isApprovedProcess:", isApprovedProcess)
						logger.Log.Print("isApprovedProcess:", isApprovedProcess)
						break
					}
				}
			}
		}
		if isApprovedProcess {
			err := dataAccess.Updateapprovalstatus(insertedID)
			if err != nil {
				log.Print("Error in Approval status update")
				logger.Log.Print("Error in Approval status update")
			}
		}
		for _, taskcat := range taskcatids {
			log.Print("\n\n\n------------", taskcat.ID)
			logger.Log.Print("\n\n\n------------", taskcat.ID)
			/**
			Fetching task details like title,short desc etc
			*/
			taskDetails, err := dataAccess.Gettaskdetailsbyid(taskcat.ID)
			if err != nil {
				return "", 0, errors.New("Something Went Wrong")
			}
			log.Print("details:", taskDetails)
			logger.Log.Print("details:", taskDetails)
			/**
			Fetching task ticket type id with diff type id
			*/
			ttypes, err := dataAccess.Gettickettypebycatid(taskcat.Val)
			if err != nil {
				return "", 0, errors.New("Something Went Wrong")
			}
			log.Print("ttypes", ttypes)
			if len(ttypes) > 0 {
				catentities := []entities.RecordData{}
				/**
				Fetching all the categories of the task ticket with the help of last category id
				and generating the category json
				*/
				parentcats, err := dataAccess.Getparentsbycatid(taskcat.Val)
				if err != nil {
					return "", 0, errors.New("Something Went Wrong")
				}
				log.Print("parentcats", parentcats)
				logger.Log.Print("parentcats", parentcats)

				catids := strings.Split(parentcats, "->")
				log.Print("catids", catids)
				logger.Log.Print("catids", catids)

				catentity := entities.RecordData{}
				for _, catid := range catids {
					log.Print("catid", catid)
					cat, _ := strconv.ParseInt(catid, 10, 64)
					difftype, err := dataAccess.Getdifftypebyid(cat)
					if err != nil {
						return "", 0, errors.New("Something Went Wrong")
					}
					catentity.ID = difftype
					catentity.Val = cat
					catentities = append(catentities, catentity)
				}
				difftype1, err := dataAccess.Getdifftypebyid(taskcat.Val)
				if err != nil {
					return "", 0, errors.New("Something Went Wrong")
				}
				catentity.ID = difftype1
				catentity.Val = taskcat.Val
				catentities = append(catentities, catentity)

				recEntity := entities.RecordcreaterequestEntity{}
				recEntity.Clientid = data.ClientID
				recEntity.Mstorgnhirarchyid = data.Mstorgnhirarchyid
				recEntity.Recorddifftypeid = ttypes[0].ID
				recEntity.Recorddiffid = ttypes[0].Val

				/**
				Fetching working category label
				*/
				worklabels, err := dataAccess.GetWorkingCatLabel(&recEntity)
				if err != nil {
					return "", 0, errors.New("Something Went Wrong")
				}
				if worklabels.WorkingCatLabelID > 0 {
					log.Print("worklabel", worklabels.WorkingCatLabelID)
					logger.Log.Print("worklabel", worklabels.WorkingCatLabelID)

					/**
					Fetching starting status details of ticket
					*/
					var seqno int64
					if isApprovedProcess {
						seqno = 20
					} else {
						seqno = 0
					}
					statusdet, err := dataAccess.Getdiffdetailsbyseq(data.ClientID, data.Mstorgnhirarchyid, seqno, 2)
					if err != nil {
						return "", 0, errors.New("Something Went Wrong")
					}
					if len(statusdet) > 0 {
						log.Print("statusdet:", statusdet)
						logger.Log.Print("statusdet:", statusdet)
						/**
						Fetching priority details of task
						*/
						recEntity.Recordtypeid = ttypes[0].Val
						recEntity.Recordcatid = taskcat.Val
						//typevalues, err := dataAccess.GetRecordprioritydata(&recEntity)
						//if err != nil {
						//	return "", 0, errors.New("Something Went Wrong")
						//}
						//if len(typevalues) > 0 {
						prioentity := entities.RecordSet{}
						prioentity.ID = parentpriodifftypeid
						prioentity.Val = parentpriodiffid

						/**
						Generating the create record json to create the record
						*/
						recordEntity := entities.RecordEntity{}
						recordEntity.ClientID = data.ClientID
						recordEntity.Mstorgnhirarchyid = data.Mstorgnhirarchyid
						recordEntity.Originalusergroupid = data.Originalusergroupid
						recordEntity.Originaluserid = data.Originaluserid
						recordEntity.CreatedusergroupID = data.CreatedusergroupID
						recordEntity.CreateduserID = data.CreateduserID
						if taskDetails[0].Title != "" {
							recordEntity.Recordname = taskDetails[0].Title
						} else {
							recordEntity.Recordname = data.Recordname
						}
						if taskDetails[0].Desc != "" {
							recordEntity.Recordesc = taskDetails[0].Desc
						} else {
							recordEntity.Recordesc = data.Recordesc
						}
						recordEntity.Requestername = data.Requestername
						recordEntity.Requesteremail = data.Requesteremail
						recordEntity.Requestermobile = data.Requestermobile
						recordEntity.Requesterlocation = data.Requesterlocation
						recordEntity.ParentID = insertedID
						recSets := []entities.RecordSet{}
						recSet := entities.RecordSet{}
						recSet.ID = 1
						recSet.Type = catentities
						recSets = append(recSets, recSet)
						recSets = append(recSets, ttypes[0])
						recSets = append(recSets, statusdet[0])
						recSets = append(recSets, prioentity)

						recordEntity.RecordSets = recSets
						recordEntity.Workingcatlabelid = worklabels.WorkingCatLabelID
						recordEntity.Source = data.Source
						recordEntity.Userid = data.Userid

						logger.Log.Print("\nfinal ", recordEntity)
						log.Print("\nfinal ", recordEntity)
						taskticketID, taskinsertedID, modelResponseError := AddRecordModelAction(&recordEntity, db, ttypes[0].Val)
						if modelResponseError != nil {
							//return "", 0, errors.New("Something went wrong")
							logger.Log.Print("\nError in task creation")
							log.Print("\nError in task creation")
						}

						log.Print(" task id ", taskticketID, taskinsertedID)
						logger.Log.Print(" task id ", taskticketID, taskinsertedID)

					} else {
						logger.Log.Print("\nStarting status not mapped with task")
						log.Print("\nStarting status not mapped with task")
					}
				} else {
					logger.Log.Print("\nWorking label not mapped with task")
					log.Print("\nWorking label not mapped with task")
				}

				log.Print("enttities:", catentities)
			} else {
				logger.Log.Print("\nTicket type not mapped with task")
				log.Print("\nTicket type not mapped with task")
			}
		}
	}
	return ticketID, insertedID, nil

}

//AddRecordModelAction for implements all fields validation and business logic
func AddRecordModelAction(data *entities.RecordEntity, dbcon *sql.DB, RecordtypeID int64) (string, int64, error) {
	logger.Log.Println("Inside AddRecordModelAction Models")
	var workflowflag bool
	var workingdiffid int64
	var recordtypeid int64
	var recordpriorityid int64
	var errormsg string
	var fstlevelcatID int64
	var ticketID string
	var recordset = data.RecordSets
	for i := 0; i < len(recordset); i++ {
		if len(recordset[i].Type) > 0 {
			for j := 0; j < len(recordset[i].Type); j++ {
				if recordset[i].Type[j].ID == data.Workingcatlabelid {
					workingdiffid = recordset[i].Type[j].Val
					//Encode the data
					postBody, _ := json.Marshal(map[string]int64{
						"clientid":          data.ClientID,
						"mstorgnhirarchyid": data.Mstorgnhirarchyid,
						"recorddifftypeid":  recordset[i].Type[j].ID,
						"recorddiffid":      recordset[i].Type[j].Val,
					})

					responseBody := bytes.NewBuffer(postBody)
					resp, err := http.Post(dbconfig.MASTER_URL+"/checkworkflow", "application/json", responseBody)
					if err != nil {
						logger.Log.Println("An Error Occured --->", err)
						return "", 0, errors.New("Something went wrong")
					}
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						logger.Log.Println("response body ------> ", err)
						return "", 0, errors.New("Something went wrong")
					}
					sb := string(body)
					wfres := entities.WorkflowResponse{}
					json.Unmarshal([]byte(sb), &wfres)
					fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>", wfres)
					logger.Log.Println("strB value is -11-->", wfres.Success)
					workflowflag = wfres.Success
					errormsg = wfres.Message
				}
			}
		}
	}

	//SLA Checking done here
	for j := 0; j < len(recordset); j++ {
		if recordset[j].ID == 2 {
			recordtypeid = recordset[j].Val
		}
		if recordset[j].ID == 5 {
			recordpriorityid = recordset[j].Val
		}

	}
	logger.Log.Println("recordtypeid value is  --->", recordtypeid)
	logger.Log.Println("recordpriorityid value is  --->", recordpriorityid)
	logger.Log.Println("workflowflag value is  --->", workflowflag)
	//SLA checking end here

	//Workflow checking strt here
	if workflowflag == true {

		//lock.Lock()
		//defer lock.Unlock()
		// dbcon, err := dbconfig.ConnectMySqlDb()
		// logger.Log.Println("111111111111111111111111111-------------------->", workflowflag)
		// if err != nil {
		// 	//dbcon.Close()
		// 	logger.Log.Println("Error in DBConnection in side AddRecordModelAction")
		// 	return "", 0, errors.New("Connection Error")
		// }

		var cerror error
		ticketID, cerror = CreateRecordCode(data, RecordtypeID, dbcon) // generate ticket number here
		if cerror != nil {
			logger.Log.Println(cerror)
			return "", 0, errors.New("Ticket number generation failures in AddRecordModelAction")
		}

		logger.Log.Println("ticketID1111111111====================>", ticketID)

		tx, err := dbcon.Begin()
		if err != nil {
			logger.Log.Println("Transaction creation error in AddRecordModelAction", err)
			return "", 0, err
		}
		a := makeTimestamp()
		var Transactionnumber = strconv.Itoa(int(a))
		insertedID, err := dao.InsertTrnRecord(tx, data, Transactionnumber, ticketID)
		if err != nil {
			logger.Log.Println(err)
			tx.Rollback()
			return "", 0, errors.New("Insertion failure in trnoder table")
		}
		logger.Log.Println("Last id value is :", insertedID)

		if insertedID > 0 {
			lastInsertedStageID, err := dao.InsertTrnRecordStage(tx, data, insertedID, Transactionnumber)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("Insertion failure in trnoderstage table")
			}
			//here we manage record priority,type,category,status etc.
			for i := 0; i < len(recordset); i++ {
				if len(recordset[i].Type) > 0 {
					fstlevelcatID, err = TypeArrayDataHandle(tx, data, insertedID, lastInsertedStageID, recordset[i].ID, recordset[i].Type, Transactionnumber)
					if err != nil {
						logger.Log.Println(err)
						tx.Rollback()
						return "", 0, errors.New("Insertion failure in type array data in AddRecordModelAction")
					}
				} else {
					err = ValueArrayDataHandle(tx, data, insertedID, lastInsertedStageID, recordset[i].ID, recordset[i].Val, Transactionnumber)
					if err != nil {
						logger.Log.Println(err)
						tx.Rollback()
						return "", 0, errors.New("Insertion failure in value array data in AddRecordModelAction")
					}
				}

			}
			var recordfields = data.Recordfields // here we manage record attachment,extra fields here
			err = RecordfieldsModels(tx, data, insertedID, lastInsertedStageID, data.CreateduserID, recordfields, Transactionnumber, dbcon)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("Insertion failure in records fields in AddRecordModelAction")
			}

			// here we manage assetdtls
			err = TypeAssetDataHandle(tx, data, insertedID, lastInsertedStageID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("Insertion failure in records fields in AddRecordModelAction")
			}

			//here we manage additional fields
			var additionalfields = data.Additionalfields
			err = RecordAdditionalFields(tx, data, additionalfields, insertedID, lastInsertedStageID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("Insertion failure in records fields in AddRecordModelAction")
			}

			//here we manage extra record ids
			err = AdditionalRecordids(tx, data, insertedID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("Insertion failure in records fields in AddRecordModelAction")
			}

			// here we update working category label
			err = dao.Updateisworking(tx, data.ClientID, data.Mstorgnhirarchyid, insertedID, lastInsertedStageID, data.Workingcatlabelid)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("Updateisworking updation failure.")
			}

			sgtable := entities.StagetableEntity{}
			dataAccess := dao.DbConn{DB: dbcon}
			var recordtype int64
			hashmap, err := dataAccess.Getdiffdtls(data.ClientID, data.Mstorgnhirarchyid)
			if err != nil {
				return "", 0, errors.New("Somethig went wrong.")
			}
			logger.Log.Println(sgtable, recordtype)
			Username, err := dataAccess.GetUsername(data.CreateduserID)
			if err != nil {
				return "", 0, errors.New("Somethig went wrong.")
			}
			logger.Log.Println(Username)
			for i := 0; i < len(recordset); i++ {
				if len(recordset[i].Type) == 0 {
					if recordset[i].ID == 2 {
						sgtable.Tickettypeid = recordset[i].Val
						sgtable.Tickettype = hashmap[recordset[i].Val]
						recordtype = recordset[i].Val
					}

					if recordset[i].ID == 3 {
						sgtable.Statusid = recordset[i].Val
						sgtable.Status = hashmap[recordset[i].Val]
					}

					if recordset[i].ID == 5 {
						sgtable.Priorityid = recordset[i].Val
						sgtable.Priority = hashmap[recordset[i].Val]
						impactID, urgencyID, err := dao.GetImpactUrgencydtls(tx, data.ClientID, data.Mstorgnhirarchyid, recordtype, recordset[i].Val)
						if err != nil {
							logger.Log.Println(err)
							return "", 0, errors.New("Somethig went wrong.")
						}
						sgtable.Impact = hashmap[impactID]
						sgtable.Impactid = impactID
						sgtable.Urgency = hashmap[urgencyID]
						sgtable.Urgencyid = urgencyID
					}

				}
			}

			// err = dao.UpdateRecordID(tx, insertedID, Transactionnumber)
			// if err != nil {
			// 	logger.Log.Println(err)
			// 	return "", 0, errors.New("Ticket random number updation failures....")
			// }

			//here insert logdata into tables

			err = dao.Updatestageid(tx, insertedID, lastInsertedStageID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("stageid updation failure.")
			}

			// New Addition 02.08.2021

			sgtable.ClientID = data.ClientID
			sgtable.OrgnID = data.Mstorgnhirarchyid
			sgtable.RecordID = insertedID
			sgtable.TicketID = ticketID
			sgtable.Source = data.Source
			sgtable.RequestorID = data.CreateduserID
			sgtable.Requestorname = data.Requestername
			sgtable.Requestorlocation = data.Requesterlocation
			sgtable.Requestorphone = data.Requestermobile
			sgtable.Requestoremail = data.Requesteremail
			sgtable.Shortdescription = data.Recordname

			//&value.Orgcreatorphone, &value.Orgcreatorlocation
			originalInfo, err := dataAccess.GetOriginalInfo(data.Originaluserid)
			if err != nil {
				logger.Log.Println(err)
				return "", 0, errors.New("Somethig went wrong.")
			}
			sgtable.Orgcreatorloginid = originalInfo.Orgcreatorloginid
			sgtable.Orgcreatorname = originalInfo.Orgcreatorname
			sgtable.Orgcreatoremail = originalInfo.Orgcreatoremail
			sgtable.Orgcreatorphone = originalInfo.Orgcreatorphone
			sgtable.Orgcreatorlocation = originalInfo.Orgcreatorlocation
			sgtable.Orgcreatorid = data.Originaluserid
			sgtable.LastuserID = data.CreateduserID
			sgtable.Lastusername = Username

			creatorInfo, err := dataAccess.GetCreatorInfo(data.CreateduserID)
			if err != nil {
				logger.Log.Println(err)
				return "", 0, errors.New("Somethig went wrong.")
			}

			catname, err := dataAccess.GetFstlevelCatName(fstlevelcatID)
			if err != nil {
				logger.Log.Println(err)
				return "", 0, errors.New("Somethig went wrong.")
			}
			sgtable.Fstlevelcategorynm = catname
			sgtable.Requestorloginid = creatorInfo.Requestorloginid
			sgtable.Vipticket = creatorInfo.Vipticket
			err = dao.InsertStageTbl(tx, sgtable)
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				return "", 0, errors.New("Activity log data insertion failure")
			}
			// New Addition 02.08.2021

			//here database commit done here.
			err = tx.Commit()
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", 0, errors.New("Transaction commit failure in AddRecordModelAction")
			} else {

				// Unique record number generation logic implementation 25.12.2021

				// ticket id generation done here.
				/*for {
					ticketID, err = CreateRecordNumber(data, insertedID, Transactionnumber, recordtype, dbcon) // generate ticket number here
					if err != nil {
						logger.Log.Println(err)
						return "", 0, errors.New("Ticket number generation failures in AddRecordModelAction")
					}

					//record number duplicate checking............
					count, err := dataAccess.Checkrecordnumberduplicacy(ticketID)
					if err != nil {
						logger.Log.Println(err)
						return "", 0, errors.New("Error in Checkrecordnumberduplicacy..")
					}
					if count == 0 {
						err = dataAccess.UpdateRecordID(insertedID, ticketID)
						if err != nil {
							logger.Log.Println(err)
							if strings.Contains(fmt.Sprint(err), "code_UNIQUE") {
								// do something
							}
							return "", 0, errors.New("Ticket random number updation failures....")
						}
						err = dataAccess.UpdateRecordIDINStage(insertedID, ticketID)
						if err != nil {
							logger.Log.Println(err)
							return "", 0, errors.New("Ticket random number updation failures....")
						}
						break
					}
				}*/

				util, err := dataAccess.Gettimediff(data.ClientID, data.Mstorgnhirarchyid)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					dataAccess.UpdateTrnrecorddeleteflg(insertedID)
					dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
					return "", 0, errors.New("Something Went Wrong")
				}

				creationdate, err := dataAccess.GetRecordcreationdate(insertedID, util.Timediff)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					dataAccess.UpdateTrnrecorddeleteflg(insertedID)
					dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
					return "", 0, errors.New("Something Went Wrong")
				}
				//	logger.Log.Println("creationdate +++++++++++++++++++++++++++++++++++ >>>>>>>>>>>>>>>>>>>>>>>>>>>>>", creationdate)
				if len(creationdate) > 0 {
					err := dataAccess.UpdateStageCreatedate(insertedID, creationdate)
					if err != nil {
						logger.Log.Println("error is ----->", err)
						dataAccess.UpdateTrnrecorddeleteflg(insertedID)
						dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
						return "", 0, errors.New("Something Went Wrong")
					}
				}

				//ticketID, err = CreateRecordNumber(data, insertedID, Transactionnumber, recordtype, dbcon) // generate ticket number here
				// ticketID, err = CreateRecordNumberNew(data, insertedID, Transactionnumber, recordtype, dbcon) // generate ticket number here
				// if err != nil {
				// 	logger.Log.Println(err)
				// 	dataAccess.UpdateTrnrecorddeleteflg(insertedID)
				// 	return "", 0, errors.New("Ticket number generation failures in AddRecordModelAction")
				// }
				// .......................................

				// ...........................................
				//here we get record create date
				err = dataAccess.InsertRecordActivityLogs(data.ClientID, data.Mstorgnhirarchyid, insertedID, 1, "Ticketid is: "+ticketID, data.CreateduserID, data.CreatedusergroupID)
				if err != nil {
					logger.Log.Println("error is ----->", err)
					dataAccess.UpdateTrnrecorddeleteflg(insertedID)
					dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
					return "", 0, errors.New("Activity log data insertion failure")
				}
				//-----------------------------------------------------------------

				// new addition ---------------------------------------------------------
				diffid, err := dataAccess.Getrecordtypediffid(data.ParentID, data.ClientID, data.Mstorgnhirarchyid)
				if err != nil {
					dataAccess.UpdateTrnrecorddeleteflg(insertedID)
					dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
					return "", 0, errors.New("Something Went Wrong")
				}
				childdiffid, err := dataAccess.Getrecordtypediffid(insertedID, data.ClientID, data.Mstorgnhirarchyid)
				if err != nil {
					dataAccess.UpdateTrnrecorddeleteflg(insertedID)
					dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
					return "", 0, errors.New("Something Went Wrong")
				}
				childseq, err := dataAccess.Getseqbyid(data.ClientID, data.Mstorgnhirarchyid, childdiffid)
				if err != nil {
					dataAccess.UpdateTrnrecorddeleteflg(insertedID)
					dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
					return "", 0, errors.New("Something Went Wrong")
				}

				if childseq == 3 || childseq == 5 {
					var ids []int64
					rdrEntity := entities.RecordDetailsRequestEntity{}
					rdrEntity.Clientid = data.ClientID
					rdrEntity.Mstorgnhirarchyid = data.Mstorgnhirarchyid
					rdrEntity.RecordDiffTypeid = 2
					rdrEntity.RecordDiffid = diffid
					rdrEntity.ParentID = data.ParentID
					ids = append(ids, insertedID)
					rdrEntity.ChildIDS = ids
					_, parErr := dataAccess.SaveChildRecord(&rdrEntity)
					if parErr != nil {
						//return "", 0, errors.New("Something went wrong")
						logger.Log.Print("\nError in task child mapping ", ticketID)

					}
				}

				isApproveworkflow, err := dataAccess.GetIsapprovalworkflow(data.ParentID)
				if err != nil {
					dataAccess.UpdateTrnrecorddeleteflg(insertedID)
					dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
					return "", 0, errors.New("Something Went Wrong")
				}

				isApproveflag, err := dataAccess.GetIsapprovalFlag(data.ParentID)
				if err != nil {
					return "", 0, errors.New("Something Went Wrong")
				}
				logger.Log.Println("childseq       --->", childseq)
				logger.Log.Println("isApproveworkflow       --->", isApproveworkflow)
				logger.Log.Println("isApproveflag       --->", isApproveflag)
				// new addition ---------------------------------------------------------

				if workingdiffid > 0 && childseq != 3 {
					if isApproveflag == 1 && childseq == 5 {
						inactivestateID, err := dataAccess.GetPreviousStateID(data.ClientID, data.Mstorgnhirarchyid, 20)
						if err != nil {
							logger.Log.Println(err)
							dataAccess.UpdateTrnrecorddeleteflg(insertedID)
							dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
							return "", 0, errors.New("Something Went Wrong")
						}

						openstateID, err := dataAccess.GetPreviousStateID(data.ClientID, data.Mstorgnhirarchyid, 29)
						if err != nil {
							logger.Log.Println(err)
							return "", 0, errors.New("Something Went Wrong")
						}

						postBody, _ := json.Marshal(map[string]int64{
							"clientid":             data.ClientID,
							"mstorgnhirarchyid":    data.Mstorgnhirarchyid,
							"recorddifftypeid":     data.Workingcatlabelid,
							"recorddiffid":         workingdiffid,
							"previousstateid":      inactivestateID,
							"manualstateselection": 0,
							"transactionid":        insertedID,
							"mstgroupid":           data.CreatedusergroupID,
							"mstuserid":            data.CreateduserID,
							"userid":               data.CreateduserID,
							"createdgroupid":       data.CreatedusergroupID,
							"currentstateid":       openstateID,
						})
						b, err := json.MarshalIndent(postBody, "", "  ")
						if err != nil {
							fmt.Println(err)
						}
						logger.Log.Println("postBody value is --->", string(b))
						responseBody := bytes.NewBuffer(postBody)
						logger.Log.Println("postBody       --->", responseBody)
						resp, err := http.Post(dbconfig.MASTER_URL+"/moveWorkflow", "application/json", responseBody)
						if err != nil {
							logger.Log.Println("An Error Occured --->", err)
						}
						defer resp.Body.Close()
						//Read the response body
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							logger.Log.Println("response body ------> ", err)
						}
						sb := string(body)
						logger.Log.Println("sb body value is --33333333333333->", sb)
					} else {
						//Encode the data
						postBody, _ := json.Marshal(map[string]int64{
							"clientid":             data.ClientID,
							"mstorgnhirarchyid":    data.Mstorgnhirarchyid,
							"recorddifftypeid":     data.Workingcatlabelid,
							"recorddiffid":         workingdiffid,
							"previousstateid":      -1,
							"manualstateselection": 0,
							"transactionid":        insertedID,
							"mstgroupid":           data.CreatedusergroupID,
							"mstuserid":            data.CreateduserID,
							"userid":               data.CreateduserID,
							"createdgroupid":       data.CreatedusergroupID,
						})
						b, err := json.MarshalIndent(postBody, "", "  ")
						if err != nil {
							fmt.Println(err)
						}
						logger.Log.Println("postBody value is --->", string(b))
						responseBody := bytes.NewBuffer(postBody)
						logger.Log.Println("postBody       --->", responseBody)
						resp, err := http.Post(dbconfig.MASTER_URL+"/moveWorkflow", "application/json", responseBody)
						if err != nil {
							logger.Log.Println("An Error Occured --->", err)
						}
						defer resp.Body.Close()
						//Read the response body
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							logger.Log.Println("response body ----11111--> ", err)
						}
						sb := string(body)
						logger.Log.Println("sb body value is --->", sb)
					}

				} else if workingdiffid > 0 && childseq == 3 {
					if isApproveworkflow == 1 {
						if isApproveflag == 1 {
							logger.Log.Println("In side if flaggggggggggggggggggggggggggggggggggggg       --->")
							prestateID, err := dataAccess.GetPreviousStateID(data.ClientID, data.Mstorgnhirarchyid, 20)
							if err != nil {
								return "", 0, errors.New("Something Went Wrong")
							}

							postBody, _ := json.Marshal(map[string]int64{
								"clientid":               data.ClientID,
								"mstorgnhirarchyid":      data.Mstorgnhirarchyid,
								"recorddifftypeid":       data.Workingcatlabelid,
								"recorddiffid":           workingdiffid,
								"previousstateid":        prestateID,
								"manualstateselection":   0,
								"transactionid":          insertedID,
								"mstgroupid":             data.CreatedusergroupID,
								"mstuserid":              data.CreateduserID,
								"userid":                 data.CreateduserID,
								"creatticketIDedgroupid": data.CreatedusergroupID,
							})
							b, err := json.MarshalIndent(postBody, "", "  ")
							if err != nil {
								fmt.Println(err)
							}
							logger.Log.Println("postBody value is --->", string(b))
							responseBody := bytes.NewBuffer(postBody)
							logger.Log.Println("postBody       --->", responseBody)
							resp, err := http.Post(dbconfig.MASTER_URL+"/moveWorkflow", "application/json", responseBody)
							if err != nil {
								logger.Log.Println("An Error Occured --->", err)
							}
							defer resp.Body.Close()
							//Read the response body
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								logger.Log.Println("response body ------> ", err)
							}
							sb := string(body)
							logger.Log.Println("sb body value is --33333333333333->", sb)
						} else {
							logger.Log.Println("In side Else flaggggggggggggggggggggggggggggggggggggg       --->")
							postBody, _ := json.Marshal(map[string]int64{
								"clientid":             data.ClientID,
								"mstorgnhirarchyid":    data.Mstorgnhirarchyid,
								"recorddifftypeid":     data.Workingcatlabelid,
								"recorddiffid":         workingdiffid,
								"previousstateid":      -1,
								"manualstateselection": 0,
								"transactionid":        insertedID,
								"mstgroupid":           data.CreatedusergroupID,
								"mstuserid":            data.CreateduserID,
								"userid":               data.CreateduserID,
								"createdgroupid":       data.CreatedusergroupID,
							})
							b, err := json.MarshalIndent(postBody, "", "  ")
							if err != nil {
								fmt.Println(err)
							}
							logger.Log.Println("postBody value is --->", string(b))
							responseBody := bytes.NewBuffer(postBody)
							logger.Log.Println("postBody       --->", responseBody)
							resp, err := http.Post(dbconfig.MASTER_URL+"/moveWorkflow", "application/json", responseBody)
							if err != nil {
								logger.Log.Println("An Error Occured --->", err)
							}
							defer resp.Body.Close()
							//Read the response body
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								logger.Log.Println("response body ------> ", err)
							}
							sb := string(body)
							logger.Log.Println("sb body value is -222222222-->", sb)
						}
					} else {
						prestateID, err := dataAccess.GetPreviousStateID(data.ClientID, data.Mstorgnhirarchyid, 20)
						if err != nil {
							return "", 0, errors.New("Something Went Wrong")
						}

						postBody, _ := json.Marshal(map[string]int64{
							"clientid":             data.ClientID,
							"mstorgnhirarchyid":    data.Mstorgnhirarchyid,
							"recorddifftypeid":     data.Workingcatlabelid,
							"recorddiffid":         workingdiffid,
							"previousstateid":      prestateID,
							"manualstateselection": 0,
							"transactionid":        insertedID,
							"mstgroupid":           data.CreatedusergroupID,
							"mstuserid":            data.CreateduserID,
							"userid":               data.CreateduserID,
							"createdgroupid":       data.CreatedusergroupID,
						})
						b, err := json.MarshalIndent(postBody, "", "  ")
						if err != nil {
							fmt.Println(err)
						}
						logger.Log.Println("postBody value is --->", string(b))
						responseBody := bytes.NewBuffer(postBody)
						logger.Log.Println("postBody       --->", responseBody)
						resp, err := http.Post(dbconfig.MASTER_URL+"/moveWorkflow", "application/json", responseBody)
						if err != nil {
							logger.Log.Println("An Error Occured --->", err)
						}
						defer resp.Body.Close()
						//Read the response body
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							logger.Log.Println("response body ------> ", err)
						}
						sb := string(body)
						logger.Log.Println("sb body value is --33333333333333->", sb)
					}
				}

			}

			// ebonding code start here 10.06.2022
			workingdifftypeID, workingdiffID, err := dataAccess.GetWorkingCategoryDetails(insertedID)
			if err != nil {
				dataAccess.UpdateTrnrecorddeleteflg(insertedID)
				dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
				return "", 0, errors.New("Something Went Wrong")
			}

			ebondingSeq, err := dataAccess.GetEbondingSeq(data.ClientID, data.Mstorgnhirarchyid, workingdifftypeID, workingdiffID)
			if err != nil {
				dataAccess.UpdateTrnrecorddeleteflg(insertedID)
				dataAccess.UpdateRecordfulldetailsddeleteflg(insertedID)
				return "", 0, errors.New("Something Went Wrong")
			}
			logger.Log.Println("ebonding seq number is =====> ", ebondingSeq)
			if ebondingSeq > 0 {
				ebonding := entities.EbondingRecordEntity{}
				ebonding.ClientID = data.ClientID
				ebonding.MstorgnhirarchyID = data.Mstorgnhirarchyid
				ebonding.RecorddifftypeID = 2
				ebonding.RecorddiffID = recordtypeid
				ebonding.RecordID = insertedID
				ebonding.RecordStagedID = lastInsertedStageID
				ebonding.RecordCode = ticketID
				ebonding.EbondingSeq = ebondingSeq
				ebonding.EbondingModuleSeq = 2
				go Ebonding(&ebonding)

			}

			//ebonding code end here 10.06.2022

			return ticketID, insertedID, nil
		}
		return "", 0, errors.New("Record creation failure")
	} else {
		return "", 0, errors.New(errormsg)
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

var recordtype int64

//ValueArrayDataHandle handle priority,ticket type,ticket seq etc..
func ValueArrayDataHandle(tx *sql.Tx, data *entities.RecordEntity, insertedID int64, lastInsertedStageID int64, id int64, val int64, Transactionnumber string) error {
	err := dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, insertedID, lastInsertedStageID, id, val, 0, Transactionnumber)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("Insertion failure in value array data")
	}

	if id == 2 {
		recordtype = val
	}

	logger.Log.Println("GetImpactUrgencydtls parameters --------------------------222222222222222222222222222222222222222------------------------------->", id, val)
	if id == 5 {
		impactID, urgencyID, err := dao.GetImpactUrgencydtls(tx, data.ClientID, data.Mstorgnhirarchyid, recordtype, val)
		if err != nil {
			logger.Log.Println(err)
			return errors.New("Insertion failure in value array data")
		}
		if impactID > 0 {
			err := dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, insertedID, lastInsertedStageID, 7, impactID, 0, Transactionnumber)
			if err != nil {
				logger.Log.Println(err)
				return errors.New("Insertion failure in value array data")
			}
		}

		if urgencyID > 0 {
			err := dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, insertedID, lastInsertedStageID, 8, urgencyID, 0, Transactionnumber)
			if err != nil {
				logger.Log.Println(err)
				return errors.New("Insertion failure in value array data")
			}
		}
	}
	return nil
}

//TypeArrayDataHandle handle category details data.
func TypeArrayDataHandle(tx *sql.Tx, data *entities.RecordEntity, insertedID int64, lastInsertedStageID int64, id int64, typedata []entities.RecordData, Transactionnumber string) (int64, error) {
	var fstlevelcat int64

	for i := 0; i < len(typedata); i++ {
		fmt.Println("individual val value :", typedata[i].ID)
		if i == 0 {
			fstlevelcat = typedata[i].Val
		}
		err := dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, insertedID, lastInsertedStageID, typedata[i].ID, typedata[i].Val, 0, Transactionnumber)
		if err != nil {
			logger.Log.Println(err)
			return fstlevelcat, errors.New("Insertion failure in value array data")
		}
	}
	return fstlevelcat, nil
}

//RecordAdditionalFields handle additional details data.
func RecordAdditionalFields(tx *sql.Tx, data *entities.RecordEntity, additionalfields []entities.RecordAdditional, insertedID int64, lastInsertedStageID int64) error {
	for i := 0; i < len(additionalfields); i++ {
		var a = &additionalfields[i]
		//if len(a.Val) > 0 {
		err := dao.InsertRecordAdditional(tx, data, a, insertedID, lastInsertedStageID)
		if err != nil {
			logger.Log.Println(err)
			return errors.New("Insertion failure in RecordAdditionalFields data")
		}
		//}

	}
	return nil
}

func AdditionalRecordids(tx *sql.Tx, data *entities.RecordEntity, insertedID int64) error {
	for i := 0; i < len(data.RecordIds); i++ {
		err := dao.InsertExtraRecordID(tx, data.ClientID, data.Mstorgnhirarchyid, insertedID, data.RecordIds[i])
		if err != nil {
			logger.Log.Println(err)
			return errors.New("Insertion failure in RecordAdditionalFields data")
		}
		//}

	}
	return nil
}

//TypeAssetDataHandle handle asset details data.
func TypeAssetDataHandle(tx *sql.Tx, data *entities.RecordEntity, insertedID int64, lastInsertedStageID int64) error {
	for i := 0; i < len(data.AssetIds); i++ {
		err := dao.InsertRecordAsset(tx, data, data.AssetIds[i], insertedID, lastInsertedStageID)
		if err != nil {
			logger.Log.Println(err)
			return errors.New("Insertion failure in TypeAssetDataHandle data")
		}
	}
	return nil
}

//RecordfieldsModels handle attachment & extra fields details data.
func RecordfieldsModels(tx *sql.Tx, data *entities.RecordEntity, insertedID int64, lastInsertedStageID int64, createdbyID int64, recordfields []entities.RecordField, Transactionnumber string, dbcon *sql.DB) error {
	for i := 0; i < len(recordfields); i++ {
		var fields = recordfields[i].Val
		for k := 0; k < len(fields); k++ {
			err := dao.InsertRecordTermvalues(tx, data, insertedID, lastInsertedStageID, createdbyID, recordfields[i].TermID, fields[k].OriginalName, fields[k].FileName, Transactionnumber)
			if err != nil {
				logger.Log.Println(err)
				return errors.New("Insertion failure in RecordfieldsModels data")
			}

			// here we checking term value is attachment or not

			// db, err := dbconfig.ConnectMySqlDb()
			// if err != nil {
			// 	logger.Log.Println("database connection failure", err)
			// 	return errors.New("Insertion failure in RecordfieldsModels data")
			// }
			// lock.Lock()
			// defer lock.Unlock()
			//defer db.Close()
			dataAccess := dao.DbConn{DB: dbcon}
			termsequance, err := dataAccess.GetRecordtermSequance(data.ClientID, data.Mstorgnhirarchyid, recordfields[i].TermID)
			logger.Log.Println("Term sequance value is >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", termsequance)
			if err != nil {
				logger.Log.Println("error is ----->", err)
				return errors.New("Insertion failure in RecordfieldsModels data")
			}
			if termsequance == 1 {
				termnm, err := dataAccess.Gettermnamebyid(recordfields[i].TermID, data.ClientID, data.Mstorgnhirarchyid)
				if err != nil {
					return errors.New("Insertion failure in RecordfieldsModels data")
				}

				err1 := dao.InsertActivityLogsfromterms(tx, data.ClientID, data.Mstorgnhirarchyid, insertedID, 100, termnm+" :: "+fields[k].OriginalName, data.CreateduserID, data.CreatedusergroupID, recordfields[i].TermID)
				if err1 != nil {
					logger.Log.Println("error is ----->", err1)
					return errors.New("Insertion failure in RecordfieldsModels data")
				}
			}

		} // end for loop
	}
	return nil
}

//CreateRecordNumber is used for generate custom record number
func CreateRecordNumber(data *entities.RecordEntity, insertedID int64, Transactionnumber string, recordtype int64, dbcon *sql.DB) (string, error) {
	var ticketID string
	dataAccess := dao.DbConn{DB: dbcon}

	log.Println("CreateRecordNumber >>>>>>>>>data.ClientID >>>>>>>>>>", data.ClientID)
	log.Println("CreateRecordNumber >>>>>>>>>recordtype >>>>>>>>>>", recordtype)
	for {
		ticketID = ""
		Isclient, err := dataAccess.IsClientSpecificOrNot(data.ClientID, 2, recordtype)
		if err != nil {
			logger.Log.Println(err)
			return ticketID, errors.New("error in fetch recorddiffID from table")
		}

		log.Println("Is client specific or not >>>>>>>>>>>>>>>>>>>>>>>>>", Isclient)
		var number int
		var prefixarray []string
		var typediffID int64
		log.Println("Is client specific or not >>>>>>>>>>>>>>>>>>>>>>>>>", number, prefixarray, typediffID)
		if Isclient == 0 {
			number, err = dataAccess.GetNumber(data.ClientID, data.Mstorgnhirarchyid, 2, recordtype)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in fetch recorddiffID from table")
			}

			prefixarray, err = dataAccess.GetPrefixnormal(data, 2, recordtype, Transactionnumber)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in fetch prefixarray from table")
			}
		} else {
			typediffID, err = dataAccess.GetParentDiffID(data.ClientID, data.Mstorgnhirarchyid, 2, recordtype)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in fetch prefixarray from table")
			}
			log.Println("typediffID >>>>>>>>>>>>>>>>>>>>>>>>>", typediffID)

			number, err = dataAccess.GetNumberbyclientID(data.ClientID, 2, typediffID)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in fetch recorddiffID from table")
			}
			log.Println("number >>>>>>>>>>>>>>>>>>>>>>>>>", number)

			prefixarray, err = dataAccess.GetPrefixnormalbyclientID(data.ClientID, 2, typediffID)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in fetch prefixarray from table")
			}

			log.Println("prefixarray >>>>>>>>>>>>>>>>>>>>>>>>>", prefixarray)
		}

		currentTime := time.Now()
		var aa = currentTime.Format("20060102")
		if len(prefixarray[0]) > 0 {
			ticketID = ticketID + prefixarray[0]
		}

		if len(prefixarray[2]) > 0 {
			if prefixarray[2] != "NA" {
				var mm = string(aa[4:6])
				ticketID = ticketID + mm
			}
		}

		if len(prefixarray[3]) > 0 {
			if prefixarray[3] != "NA" {
				var day = string(aa[6:8])
				ticketID = ticketID + day
			}
		}

		if len(prefixarray[1]) > 0 {
			if prefixarray[1] != "NA" {
				if len(prefixarray[1]) == 2 {
					var yy = string(aa[2:4])
					ticketID = ticketID + yy
				} else {
					var yy = string(aa[0:4])
					ticketID = ticketID + yy
				}

			}
		}

		if len(prefixarray[4]) > 0 {
			var s = prefixarray[4]
			var numberlength = len(strconv.Itoa(number))
			var safeSubstring = string(s[0 : len(s)-numberlength])
			ticketID = ticketID + safeSubstring + strconv.Itoa(number)
		}

		logger.Log.Println("Final ticket id ------------->", ticketID)

		/*err = dao.UpdateRecordID(tx, insertedID, ticketID, Transactionnumber)
		if err != nil {
			logger.Log.Println(err)
			return ticketID, errors.New("error in update trnrecord table with ticket id")
		}*/

		if Isclient == 0 {
			err = dataAccess.Updatenumber(data.ClientID, data.Mstorgnhirarchyid, 2, recordtype)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in update trnrecord table with ticket id")
			}
		} else {
			err = dataAccess.UpdatenumberbyclientID(data.ClientID, 2, typediffID)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in update trnrecord table with ticket id")
			}
		}

		logger.Log.Println("Count val >>>>>>>>>111111111111111111>>>>>>>>>>")

		//record number duplicate checking............

		count, err := dataAccess.Checkrecordnumberduplicacy(ticketID)
		if err != nil {
			logger.Log.Println(err)
			return "", errors.New("Error in Checkrecordnumberduplicacy..")
		}
		logger.Log.Println("Count val >>>>>>>>>>>>>>>>>>>", count)
		if count == 0 {
			err = dataAccess.UpdateRecordID(insertedID, ticketID)
			if err != nil {
				logger.Log.Println(err)
				if strings.Contains(fmt.Sprint(err), "1062") {
					logger.Log.Println("In side code_UNIQUE if loop >>>>>>>>>>>>>>>>>>>")
					//CreateRecordNumber(data, insertedID, Transactionnumber, recordtype, dbcon) // generate ticket number here
					continue
				} else {
					return "", errors.New("Ticket random number updation failures....")
				}

			}
			err = dataAccess.UpdateRecordIDINStage(insertedID, ticketID)
			if err != nil {
				logger.Log.Println(err)
				return "", errors.New("Ticket random number updation failures....")
			}
			logger.Log.Println("*********************************************************************************************")
			break

		} else {
			//CreateRecordNumber(data, insertedID, Transactionnumber, recordtype, dbcon) // generate ticket number herecontinue
			continue
		}

	}
	logger.Log.Println("22222222222222222222222222222222222222******************************************************")
	return ticketID, nil
}

//GetRecordResponse generate Record response details here
func GetRecordResponse(dbcon *sql.DB, data *entities.RecordEntity, insertedID int64, Transactionnumber string) (entities.RecordRespone, error) {
	var arr []entities.ResponseDetail
	var reponsedtls = entities.RecordRespone{}
	for i := 0; i < len(data.ResponsedifftypeID); i++ {
		recorddiffID, err := dao.GetDiffID(dbcon, insertedID, data.ResponsedifftypeID[i], Transactionnumber)
		if err != nil {
			logger.Log.Println(err)
			return reponsedtls, errors.New("error in update trnrecord table with ticket id")
		}
		resdetails := entities.ResponseDetail{}
		resdetails.DifftypeID = data.ResponsedifftypeID[i]
		resdetails.DiffID = recorddiffID
		arr = append(arr, resdetails)
	}
	reponsedtls.ID = insertedID
	reponsedtls.ClientID = data.ClientID
	reponsedtls.Mstorgnhirarchyid = data.Mstorgnhirarchyid
	reponsedtls.ResponseDetails = arr
	return reponsedtls, nil
}

// func getWorkingDiffvalue(ClientID int64, Mstorgnhirarchyid int64, insertedID int64, lastInsertedStageID int64, Workingcatlabelid int64) (int64, error) {
// 	db, err := dbconfig.ConnectMySqlDb()
// 	if err != nil {
// 		log.Println("database connection failure", err)
// 		return 0, err
// 	}
// 	defer db.Close()
// 	dataAccess := dao.DbConn{DB: db}
// 	workingdiffid, err := dataAccess.GetWorkingdiffid(ClientID, Mstorgnhirarchyid, insertedID, lastInsertedStageID, Workingcatlabelid)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return 0, errors.New("error in fetch recorddiffid from table")
// 	}
// 	return workingdiffid, nil
// }

func getWorkingDiffvalue(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, insertedID int64, lastInsertedStageID int64, Workingcatlabelid int64) (int64, error) {
	workingdiffid, err := dao.GetWorkingdiffid(tx, ClientID, Mstorgnhirarchyid, insertedID, lastInsertedStageID, Workingcatlabelid)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("error in fetch recorddiffid from table")
	}
	return workingdiffid, nil
}

// SLA Integration codebase

// func CreateRecordNumberNew(data *entities.RecordEntity, insertedID int64, Transactionnumber string, recordtype int64, dbcon *sql.DB) (string, error) {
// 	var ticketID string
// 	var typediffID int64
// 	var ID int64
// 	dataAccess := dao.DbConn{DB: dbcon}
// 	for {
// 		Isclient, err := dataAccess.IsClientSpecificOrNot(data.ClientID, 2, recordtype)
// 		if err != nil {
// 			logger.Log.Println(err)
// 			return ticketID, errors.New("error in fetch recorddiffID from table")
// 		}
// 		if Isclient == 0 {
// 			ID, ticketID, err = dataAccess.GetParentRecordNo(2, recordtype)
// 			if err != nil {
// 				logger.Log.Println(err)
// 				return ticketID, errors.New("error in fetch from table")
// 			}
// 			//logger.Log.Println(ID)
// 		} else {
// 			typediffID, err = dataAccess.GetParentDiffID(data.ClientID, data.Mstorgnhirarchyid, 2, recordtype)
// 			if err != nil {
// 				logger.Log.Println(err)
// 				return ticketID, errors.New("error in fetch from table")
// 			}
// 			log.Println("typediffID >>>>>>>>>>>>>>>>>>>>>>>>>", typediffID)
// 			ID, ticketID, err = dataAccess.GetParentRecordNo(2, typediffID)
// 			if err != nil {
// 				logger.Log.Println(err)
// 				return ticketID, errors.New("error in fetch from table")
// 			}
// 		}

// 		logger.Log.Println("Final ticket id ------------->", ticketID)

// 		if len(ticketID) > 0 {
// 			err = dataAccess.UpdateRecordID(insertedID, ticketID)
// 			if err != nil {
// 				logger.Log.Println(err)
// 				//return "", errors.New("Ticket random number updation failures....")
// 				continue
// 			}
// 			err = dataAccess.UpdateRecordIDINStage(insertedID, ticketID)
// 			if err != nil {
// 				logger.Log.Println(err)
// 				return "", errors.New("Ticket random number updation failures....")
// 			}
// 			err = dataAccess.UpdateRecordnoTB(ID)
// 			if err != nil {
// 				logger.Log.Println(err)
// 				return "", errors.New("Record no table updation failures....")
// 			}
// 		}
// 		break
// 	}
// 	return ticketID, nil
// }

func CreateRecordCode(data *entities.RecordEntity, recordtype int64, dbcon *sql.DB) (string, error) {
	var ticketID string
	var typediffID int64
	var ID int64
	dataAccess := dao.DbConn{DB: dbcon}
	for {
		Isclient, err := dataAccess.IsClientSpecificOrNot(data.ClientID, 2, recordtype)
		if err != nil {
			logger.Log.Println(err)
			return ticketID, errors.New("error in fetch recorddiffID from table")
		}
		if Isclient == 0 {
			tx, _ := dbcon.Begin()
			ID, ticketID, err = dataAccess.GetParentRecordNo(2, recordtype, tx)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return ticketID, errors.New("error in fetch from table")
			}
			err = dataAccess.UpdateRecordnoTB(ID, tx)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", errors.New("Record no table updation failures....")
			}
			err := tx.Commit()
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", errors.New("Commit  failures....!!!!")
			}
			//logger.Log.Println(ID)
		} else {

			typediffID, err = dataAccess.GetParentDiffID(data.ClientID, data.Mstorgnhirarchyid, 2, recordtype)
			if err != nil {
				logger.Log.Println(err)
				return ticketID, errors.New("error in fetch from table")
			}
			//log.Println("typediffID >>>>>>>>>>>>>>>>>>>>>>>>>", typediffID)
			tx, _ := dbcon.Begin()
			ID, ticketID, err = dataAccess.GetParentRecordNo(2, typediffID, tx)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return ticketID, errors.New("error in fetch from table")
			}
			err = dataAccess.UpdateRecordnoTB(ID, tx)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", errors.New("Record no table updation failures....")
			}

			err := tx.Commit()
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				return "", errors.New("Commit Failed....!!!")
			}
		}

		logger.Log.Println("Final ticket id ------------->", ticketID)

		if len(ticketID) == 0 {
			logger.Log.Println("Recursive calling processing ---------------------   --------------------------------   ------------->")
			time.Sleep(25 * time.Millisecond)
			continue

		}
		break
	}
	return ticketID, nil
}
