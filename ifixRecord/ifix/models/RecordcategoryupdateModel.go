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

	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Recordcategoryupdate(data *entities.RecordcategoryupdateEntity) (int64, bool, error, string) {
	var errormsg string
	var workflowflag bool
	var Priority bool
	var recordset = data.RecordSets
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	for i := 0; i < len(recordset); i++ {
		if len(recordset[i].Type) > 0 {
			for j := 0; j < len(recordset[i].Type); j++ {
				if recordset[i].Type[j].ID == data.Workingcatlabelid {
					//workingdiffid = recordset[i].Type[j].Val
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
						return 0, false, err, "Something Went Wrong"
					}
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						logger.Log.Println("response body ------> ", err)
						return 0, false, err, "Something Went Wrong"
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

	if workflowflag == true {
		// dbcon, err := dbconfig.ConnectMySqlDb()
		// if err != nil {
		// 	logger.Log.Println("Error in DBConnection in side Recordcategoryupdate")
		// 	return 0, false, err, "Something Went Wrong"
		// }
		//defer dbcon.Close()
		if db == nil {
			dbcon, err := ConnectMySqlDb()
			if err != nil {
				logger.Log.Println("Error in DBConnection")
				return 0, false, err, "Something Went Wrong"
			}
			db = dbcon
		}
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Println("Transaction creation error in Recordcategoryupdate", err)
			//dbcon.Close()
			return 0, false, err, "Something Went Wrong"
		}
		dataAccess := dao.DbConn{DB: db}
		hashmap, err := dataAccess.Getdiffdtls(data.ClientID, data.Mstorgnhirarchyid)
		if err != nil {
			return 0, false, err, "Something Went Wrong"
		}
		lastworkingcatID, err := dataAccess.GetLastWorkingcatID(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
		logger.Log.Println("Last working category id is --->", lastworkingcatID)
		if err != nil {
			logger.Log.Println("Error is --->", err)
			//dbcon.Close()
			return 0, false, err, "Something Went Wrong"
		}

		recordtitle, recorddescription, err := dataAccess.GetLaststagevalue(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
		logger.Log.Println("Record title is --->", recordtitle)
		logger.Log.Println("Record description is --->", recorddescription)
		if err != nil {
			logger.Log.Println("Error is --->", err)
			//dbcon.Close()
			return 0, false, err, "Something Went Wrong"
		}

		lastcategoryname, err := dataAccess.Getcategorynames(data.ClientID, data.Mstorgnhirarchyid, data.Recorddifftypeid, data.Recorddiffid, data.RecordID)
		logger.Log.Println("Last category name is --->", lastcategoryname)
		if err != nil {
			logger.Log.Println("Error is --->", err)
			////dbcon.Close()
			return 0, false, err, "Something Went Wrong"
		}
		//category update here
		if len(data.RecordSets) > 0 {
			arr, err := dataAccess.GetcatlevelidagainstrecordID(data.ClientID, data.Mstorgnhirarchyid, data.Recorddifftypeid, data.Recorddiffid)
			if err != nil {
				logger.Log.Println("Error is --->", err)
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}
			//update islatest flag for old cat level
			for k := 0; k < len(arr); k++ {
				err = dao.Updatepreviousrecord(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, arr[k])
				if err != nil {
					logger.Log.Println("Error is --->", err)
					//dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}
			//insert new record into stage table
			stageID, err := dao.UpdateRecordStage(tx, data, recordtitle, recorddescription)
			if err != nil {
				logger.Log.Println("Error is --->", err)
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}
			var currentcatnames string
			for i := 0; i < len(data.RecordSets); i++ {
				if len(data.RecordSets[i].Type) > 0 {
					for m := 0; m < len(data.RecordSets[i].Type); m++ {
						err := dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, stageID, data.RecordSets[i].Type[m].ID, data.RecordSets[i].Type[m].Val, 0, "")
						if err != nil {
							logger.Log.Println("Error is --->", err)
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						nms, err := dao.Getpriorityname(tx, data.RecordSets[i].Type[m].Val)
						logger.Log.Println("current ")
						if err != nil {
							logger.Log.Println(err)
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}
						currentcatnames = currentcatnames + nms + ","
					}

				}

				//For Priority update
				if data.RecordSets[i].ID == 5 {
					pid, recorddiffid, err := dao.GetlatestDiffID(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 5)
					if err != nil {
						tx.Rollback()
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}

					if recorddiffid != data.RecordSets[i].Val {
						prepirotynm, err := dao.GetPreviouspriorityname(tx, pid)
						if err != nil {
							logger.Log.Println(err)
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}
						err = dao.Updatepreviousrecord(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, data.RecordSets[i].ID)
						if err != nil {
							logger.Log.Println("Error is --->", err)
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}
						err = dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, stageID, data.RecordSets[i].ID, data.RecordSets[i].Val, 0, "")
						if err != nil {
							logger.Log.Println("Error is --->", err)
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						curpirotynm, err := dao.Getpriorityname(tx, data.RecordSets[i].Val)
						if err != nil {
							logger.Log.Println(err)
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						//activity log entry here

						err = dao.InsertActivityLogs(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 10, "From "+prepirotynm+" To "+curpirotynm, data.UserID, data.UsergroupID)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						//Update Stage TBL For Priority

						//err = dao.UpdateStagePriority(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, data.RecordSets[i].Val, curpirotynm)
						err = dataAccess.UpdateStagePriority(data.ClientID, data.Mstorgnhirarchyid, data.RecordID, data.RecordSets[i].Val, curpirotynm)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						//Update Stage TBL For Priority

						// Newly addition in 18.07.2021
						//typeID, _, err := dao.GetlatestDiffID(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 2)
						typeID, err := dataAccess.GetRecordTypeID(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
						if err != nil {
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						impactID, urgencyID, err := dao.GetImpactUrgencydtls(tx, data.ClientID, data.Mstorgnhirarchyid, typeID, data.RecordSets[i].Val)
						if err != nil {
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						// Impact ID
						preimpactID, _, err := dao.GetlatestDiffID(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 7)
						if err != nil {
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						if preimpactID > 0 {
							err = dao.Updateoldpriorityflag(tx, preimpactID)
							if err != nil {
								tx.Rollback()
								//dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							err := dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, stageID, 7, impactID, 0, "")
							if err != nil {
								tx.Rollback()
								//dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							//Update Stage TBL For Impact

							//err = dao.UpdateStageImpact(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, impactID, hashmap[impactID])
							err = dataAccess.UpdateStageImpact(data.ClientID, data.Mstorgnhirarchyid, data.RecordID, impactID, hashmap[impactID])
							if err != nil {
								log.Println("error is ----->", err)
								tx.Rollback()
								//dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							//Update Stage TBL For Impact

						}

						//Urgency ID
						preurgencyID, _, err := dao.GetlatestDiffID(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 8)
						if err != nil {
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}

						if preurgencyID > 0 {
							err = dao.Updateoldpriorityflag(tx, preurgencyID)
							if err != nil {
								tx.Rollback()
								//dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							err := dao.InsertTrnRecordMapDifferrtiation(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, stageID, 8, urgencyID, 0, "")
							if err != nil {
								tx.Rollback()
								//dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							//Update Stage TBL For Impact

							//err = dao.UpdateStageUrgency(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, urgencyID, hashmap[urgencyID])
							err = dataAccess.UpdateStageUrgency(data.ClientID, data.Mstorgnhirarchyid, data.RecordID, urgencyID, hashmap[urgencyID])
							if err != nil {
								log.Println("error is ----->", err)
								tx.Rollback()
								//dbcon.Close()
								return 0, false, err, "Something Went Wrong"
							}

							//Update Stage TBL For Impact
						}

						Priority = true
					} // not equal

				}

			}

			//update working level flag
			err = dao.Updateisworking(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, stageID, data.Workingcatlabelid)
			if err != nil {
				logger.Log.Println("Error is --->", err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			// New Addition 04.08.2020
			//err = dao.Updatecategorychangecout(tx, data.RecordID)
			err = dataAccess.Updatecategorychangecout(data.RecordID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}
			Username, err := dataAccess.GetUsername(data.UserID)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			// err = dataAccess.UpdateUserInfo(data.ClientID, data.Mstorgnhirarchyid, data.RecordID, data.UserID, Username)
			// if err != nil {
			// 	logger.Log.Println(err)
			// 	tx.Rollback()
			// 	//dbcon.Close()
			// 	return 0, false, err, "Something Went Wrong"
			// }
			// New Addition 04.08.2020

			logger.Log.Println("current --------------------------->", currentcatnames)
			err = dao.InsertActivityLogs(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 6, "<b>From</b> "+lastcategoryname[1:]+"<br/><b> To </b>"+currentcatnames, data.UserID, data.UsergroupID)
			if err != nil {
				log.Println("error is ----->", err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			// For Additional Fields

			// for i := 0; i < len(data.Additionalfields); i++ {
			// 	var a = data.Additionalfields[i]
			// 	if len(a.Val) > 0 {
			// 		oldvalue, err := dataAccess.Getadditionaloldvalue(a.Termsid, data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
			// 		if err != nil {
			// 			logger.Log.Println(err)
			// 			tx.Rollback()
			// 			//dbcon.Close()
			// 			return 0, false, err, "Something Went Wrong"
			// 		}

			// 		err = dao.UpdateRecordAdditional(tx, data.ClientID, data.Mstorgnhirarchyid, a.ID, a.Termsid, a.Val, data.RecordID)
			// 		if err != nil {
			// 			logger.Log.Println(err)
			// 			tx.Rollback()
			// 			//dbcon.Close()
			// 			return 0, false, err, "Something Went Wrong"
			// 		}

			// 		fildsname, err := dao.Gettermnamebyid(tx, a.Termsid, data.ClientID, data.Mstorgnhirarchyid)
			// 		if err != nil {
			// 			logger.Log.Println(err)
			// 			tx.Rollback()
			// 			//dbcon.Close()
			// 			return 0, false, err, "Something Went Wrong"
			// 		}

			// 		err = dao.InsertActivityLogs(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID, 5, fildsname+" has been changed from "+oldvalue+" to "+a.Val, data.UserID, data.UsergroupID)
			// 		if err != nil {
			// 			log.Println("error is ----->", err)
			// 			tx.Rollback()
			// 			//dbcon.Close()
			// 			return 0, false, err, "Something Went Wrong"
			// 		}
			// 	}

			// }

			// New code for create auto task due to category change 06.11.2021
			recordtypeSeq, err := dataAccess.Getrecordtypeseq(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
			if err != nil {
				logger.Log.Println(err)
				return 0, false, err, "Something Went Wrong"
			}
			var lastcatid int64 = 0
			var lastcattypeid int64 = 0
			var recordset = data.RecordSets
			for i := 0; i < len(recordset); i++ {
				if len(recordset[i].Type) > 0 && recordset[i].ID == 1 {
					var catset = recordset[i].Type
					if len(catset) > 0 {
						lastcatid = catset[len(catset)-1].Val
						lastcattypeid = catset[len(catset)-1].ID
					}
				}
			}
			checklatestlastcat, err := dataAccess.Checklatestlastcatcount(data.RecordID, data.ClientID, data.Mstorgnhirarchyid, lastcattypeid, lastcatid)
			if err != nil {
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}
			logger.Log.Println("checklatestlastcat value in category upadate is ------------------------------>", checklatestlastcat)
			if recordtypeSeq == 2 || recordtypeSeq == 4 {
				if checklatestlastcat == 0 {

					diffid, err := dataAccess.Getrecordtypediffid(data.RecordID, data.ClientID, data.Mstorgnhirarchyid)
					if err != nil {
						logger.Log.Println("database connection failure", err)
						tx.Rollback()
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}
					childids, err := dataAccess.Getchildrecordids(data.RecordID, data.ClientID, data.Mstorgnhirarchyid, 2, diffid)
					if err != nil {
						logger.Log.Println("database connection failure", err)
						tx.Rollback()
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}
					currentstatusID, currentseq, _, err := dataAccess.Getrecorddifferation(data.ClientID, data.Mstorgnhirarchyid, 3)
					logger.Log.Println("currentstatusID  is -----------------category change         --------------------------------------------------------------->", currentstatusID)
					logger.Log.Println("currentseq  is ------------------------------- category change     ------------------------------------------------->", currentseq)
					if err != nil {
						logger.Log.Println(err)
						tx.Rollback()
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}
					logger.Log.Println(childids)
					//logger.Log.Println(currentseq)
					for i := 0; i < len(childids); i++ {
						logger.Log.Println("Child ID IS ==========================================>",childids[i])
						_, err = ChildSRrecordstatusupdation(tx, data.ClientID, data.Mstorgnhirarchyid, childids[i], currentstatusID, data.UserID, data.UsergroupID, db, data.RecordID, currentseq, recordtypeSeq)
						if err != nil {
							logger.Log.Println(err)
							tx.Rollback()
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}
					} // End For loop here.....
					if recordtypeSeq == 2 {
						err = dao.UpdateApproveflagzero(tx, data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
						if err != nil {
							log.Println("error is ----->", err)
							tx.Rollback()
							return 0, false, err, "Something Went Wrong"
						}
					}

					err = tx.Commit()
					if err != nil {
						logger.Log.Println("Error is --->", err)
						tx.Rollback()
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}

					var workflowflag bool
					var errormsg string
					if len(childids) > 0 {
						reqbd := &entities.ParentchildEntity{}
						reqbd.Parentid = data.RecordID
						reqbd.Childids = childids
						reqbd.Userid = data.UserID
						reqbd.Isupdate = true
						reqbd.Createdgroupid = data.UsergroupID
						postBody, _ := json.Marshal(reqbd)

						logger.Log.Println("Record status request body -->", reqbd)
						responseBody := bytes.NewBuffer(postBody)
						resp, err := http.Post(dbconfig.MASTER_URL+"/updatetaskstatus", "application/json", responseBody)
						if err != nil {
							logger.Log.Println("An Error Occured --->", err)
							return 0, false, err, "Something Went Wrong"
						}
						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							logger.Log.Println("response body ------> ", err)
							return 0, false, err, "Something Went Wrong"
						}
						sb := string(body)
						wfres := entities.WorkflowResponse{}
						json.Unmarshal([]byte(sb), &wfres)
						workflowflag = wfres.Success
						errormsg = wfres.Message
						logger.Log.Println("Record status response message -->", workflowflag)
						logger.Log.Println("Record status response error message -->", errormsg)
					}

					// New Addition 23.12.2021 -----------------------------------------------------------

					for i := 0; i < len(childids); i++ {
						err := dataAccess.Updatechildrecordflag(data.ClientID, data.Mstorgnhirarchyid, data.RecordID, childids[i])
						if err != nil {
							logger.Log.Println("database connection failure", err)

						}
					}
					// New Addition 23.12.2021 -----------------------------------------------------------

				} // count logic end here....
			} else {
				err = tx.Commit()
				if err != nil {
					logger.Log.Println("Error is --->", err)
					tx.Rollback()
					//dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
			}

			err = dataAccess.UpdateUserInfo(data.ClientID, data.Mstorgnhirarchyid, data.RecordID, data.UserID, Username)
			if err != nil {
				logger.Log.Println(err)
				tx.Rollback()
				//dbcon.Close()
				return 0, false, err, "Something Went Wrong"
			}

			if recordtypeSeq == 2 || recordtypeSeq == 4 { //|| recordtypeSeq == 1
				if checklatestlastcat == 0 {
					// SR Workflow change logic
					lastworkingcatID, err := dataAccess.GetLastWorkingcatID(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
					logger.Log.Println("Last working category id is --->", lastworkingcatID)
					if err != nil {
						logger.Log.Println("Error is --->", err)
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}
					postBody, _ := json.Marshal(map[string]int64{
						"clientid":             data.ClientID,
						"mstorgnhirarchyid":    data.Mstorgnhirarchyid,
						"recorddifftypeid":     data.Workingcatlabelid,
						"recorddiffid":         lastworkingcatID,
						"previousstateid":      -1,
						"manualstateselection": 0,
						"transactionid":        data.RecordID,
						"mstgroupid":           data.UsergroupID,
						"mstuserid":            data.UserID,
						"userid":               data.UserID,
						"createdgroupid":       data.UsergroupID,
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

					taskcatids, terr := dataAccess.Gettaskbycatid(lastcatid)
					log.Print("taskids", taskcatids)
					if terr != nil {
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}

					//--------------------------------------- 20.12.2021 ---------------------------------------------------------------
					_, recorddiffseq, _, err := dataAccess.Getcurrentsatusid(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
					if err != nil {
						logger.Log.Println(err)
						return 0, false, err, "Something Went Wrong"
					}
					if recorddiffseq == 15 || recorddiffseq == 17 || recorddiffseq == 26 { //|| isApproveworkflow == 1

						err = dataAccess.UpdateApproveflagForOne(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
						if err != nil {
							log.Println("error is ----->", err)
							return 0, false, err, "Something Went Wrong"
						}
					}

					//------------------------------------------------------------------------------------------------------------------
					if len(taskcatids) > 0 {
						err := CreateNewTaskRecord(db, data, recordtitle, recorddescription, taskcatids)
						if err != nil {
							//dbcon.Close()
							return 0, false, err, "Something Went Wrong"
						}
					}
				}
			}

			if recordtypeSeq == 1 {
				lastestworkingcatID, err := dataAccess.GetLastWorkingcatID(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
				logger.Log.Println("Last working category id is --->", lastestworkingcatID)
				if err != nil {
					logger.Log.Println("Error is --->", err)
					//dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
				if lastworkingcatID != lastestworkingcatID {
					postBody, _ := json.Marshal(map[string]int64{
						"clientid":             data.ClientID,
						"mstorgnhirarchyid":    data.Mstorgnhirarchyid,
						"recorddifftypeid":     data.Workingcatlabelid,
						"recorddiffid":         lastestworkingcatID,
						"previousstateid":      -1,
						"manualstateselection": 0,
						"transactionid":        data.RecordID,
						"mstgroupid":           data.UsergroupID,
						"mstuserid":            data.UserID,
						"userid":               data.UserID,
						"createdgroupid":       data.UsergroupID,
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

			}
			// New code for create auto task due to category change 06.11.2021

			// New logic implemented in 09.05.2022
			if recordtypeSeq == 2 {
				ids, err := dataAccess.GetLatestTwologsIds(data.RecordID)
				if err != nil {
					//dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}
				for i := 0; i < len(ids); i++ {
					err := dataAccess.UpdateDeleteLogs(ids[i])
					if err != nil {
						return 0, false, err, "Something Went Wrong"

					}
				}

			}

			//New logic implemented in 09.05.2022

			// For SLA calling
			if recordtypeSeq == 1 || recordtypeSeq == 3 {
				historyrecord, err := dataAccess.GetLatesttrnhistoryAll(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
				currentTime := time.Now()
				zonediff, _, _, _ := Getutcdiff(data.ClientID, data.Mstorgnhirarchyid)
				datetime := AddSubSecondsToDate(currentTime, zonediff.UTCdiff)

				createdt, err := dataAccess.Getrecordcreatedate(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
				if err != nil {
					//dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}

				reopendt, err := dataAccess.Getrecordreopencreatedate(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
				if err != nil {
					//dbcon.Close()
					return 0, false, err, "Something Went Wrong"
				}

				if len(createdt) > 0 && Priority {
					logger.Log.Println("createdt  1111111111111111---------------------------------------------------->", createdt)
					datetime1 := AddSubSecondsToDate(TimeParse(createdt, ""), zonediff.UTCdiff)
					// value assign
					histrn := entities.TrnslaentityhistoryEntity{}
					histrn.Clientid = data.ClientID
					histrn.Mstorgnhirarchyid = data.Mstorgnhirarchyid
					histrn.Therecordid = data.RecordID
					histrn.Mstslaentityid = historyrecord.Mstslaentityid
					histrn.Donotupdatesladue = historyrecord.Donotupdatesladue
					histrn.Mstslastateid = historyrecord.Mstslastateid
					histrn.Commentonrecord = historyrecord.Commentonrecord
					histrn.Fromclientuserid = historyrecord.Fromclientuserid
					histrn.Recorddatetime = datetime.Format("2006-01-02 15:04:05")
					histrn.Recorddatetoint = TimeParse(datetime.Format("2006-01-02 15:04:05"), "").Unix()
					histrn.Slastartstopindicator = historyrecord.Slastartstopindicator
					trnid, err := dataAccess.InsertTrnslaentityhistory(&histrn)
					if err != nil {
						logger.Log.Println(err)
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}
					logger.Log.Println("history table id---------------------------------------------------->", trnid)
					res, err := dataAccess.Getrecorddetails(data.RecordID)
					if err != nil {
						//dbcon.Close()
						return 0, false, err, "Something Went Wrong"
					}
					grpID, err := dataAccess.FetchCurrentGrpID(data.RecordID)
					if err != nil {
						return 0, false, err, "Something Went Wrong"
					}

					returnValue, _, _, _ := SLACriteriaRespResl(data.ClientID, data.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID)
					if returnValue.Supportgroupspecific == 1 {
						count, err := dataAccess.GetSupportgrpdayofweekcount(data.ClientID, data.Mstorgnhirarchyid, grpID)
						if err != nil {
							return 0, false, err, "Something Went Wrong"
						}
						if count < 7 {
							return 0, false, err, "Day Of Week Not Properly Configured.Please Check."
						}
					} else {
						count, err := dataAccess.GetOrganizationdayofweekcount(data.ClientID, data.Mstorgnhirarchyid)
						if err != nil {
							return 0, false, err, "Something Went Wrong"
						}
						if count < 7 {
							return 0, false, err, "Day Of Week Not Properly Configured.Please Check."
						}
					}
					if len(reopendt) > 0 {
						datetime2 := AddSubSecondsToDate(TimeParse(reopendt, ""), zonediff.UTCdiff)
						SLADueTimeCalculation(data.RecordID, 0, 1, 3, datetime2.Format("2006-01-02 15:04:05"), data.ClientID, data.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "P", grpID)
					} else {
						SLADueTimeCalculation(data.RecordID, 0, 1, 3, datetime1.Format("2006-01-02 15:04:05"), data.ClientID, data.Mstorgnhirarchyid, res.RecordtypeID, res.WorkingcatID, res.PriorityID, "P", grpID)
					}

					t := entities.SLATabEntity{}
					t.ClientID = data.ClientID
					t.Mstorgnhirarchyid = data.Mstorgnhirarchyid
					t.RecordID = data.RecordID
					sladtls, _, err, _ := GetSLATabvalues(&t)
					if err != nil {
						logger.Log.Println(err)
					}
					err = dataAccess.UpdateSLAFields(data.ClientID, data.Mstorgnhirarchyid, data.RecordID, sladtls.Responsedetails.Responseduetime, sladtls.Responsedetails.Responseclockstatus, sladtls.Resolutionetails.Resolutionduetime, sladtls.Resolutionetails.Resolutionclockstatus)
					if err != nil {
						logger.Log.Println(err)
					}

					//dbcon.Close()
					return stageID, true, err, "Category update successfully done."
				} else {
					//dbcon.Close()
					return stageID, true, err, "Category update successfully done."
				}
			} else {
				//dbcon.Close()
				return stageID, true, err, "Category update successfully done."
			}
			//End SLA calling

		} else {
			//dbcon.Close()
			return 0, false, err, "Something Went Wrong"
		}

	} else {
		return 0, false, errors.New(errormsg), "Workflow not configured for this category."
	}

}

func CreateNewTaskRecord(dbcon *sql.DB, data *entities.RecordcategoryupdateEntity, recordtitle string, recorddescription string, taskcatids []entities.RecordData) error {
	var recordset = data.RecordSets
	// var lastcatid int64 = 0
	// for i := 0; i < len(recordset); i++ {
	// 	if len(recordset[i].Type) > 0 && recordset[i].ID == 1 {
	// 		var catset = recordset[i].Type
	// 		if len(catset) > 0 {
	// 			lastcatid = catset[len(catset)-1].Val
	// 		}
	// 	}
	// }
	dataAccess := dao.DbConn{DB: dbcon}

	// taskcatids, terr := dataAccess.Gettaskbycatid(lastcatid)
	// log.Print("taskids", taskcatids)
	// if terr != nil {
	// 	errors.New("Something Went Wrong")
	// }
	// var hastask bool
	// if len(taskcatids) > 0 {
	// 	hastask = true
	// } else {
	// 	hastask = false
	// }
	var parentpriodifftypeid int64
	var parentpriodiffid int64
	for i := 0; i < len(recordset); i++ {
		if recordset[i].ID == 5 {
			parentpriodifftypeid = recordset[i].ID
			parentpriodiffid = recordset[i].Val
		}
	}
	seq, seqerr := dataAccess.Getseqbyid(data.ClientID, data.Mstorgnhirarchyid, data.Recorddiffid)
	if seqerr != nil {
		return errors.New("Something Went Wrong")
	}
	if seq == 2 && len(taskcatids) == 0 {
		return errors.New("No task mapped with SR ticket type")
	}

	requestername, requesteremail, requestermobile, requesterlocation, source, userid, usergroupid, originaluserid, originalusergroupid, err := dataAccess.GetparentrecordInfo(data.ClientID, data.Mstorgnhirarchyid, data.RecordID)
	if err != nil {
		return errors.New("Something Went Wrong")
	}
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
						return errors.New("Something Went Wrong")
					}
					count, err := dataAccess.Checkisapprovedprocess(&reccatentity, pids.ID)
					if err != nil {
						return errors.New("Something Went Wrong")
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
		err := dataAccess.Updateapprovalstatus(data.RecordID)
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
			return errors.New("Something Went Wrong")
		}
		log.Print("details:", taskDetails)
		logger.Log.Print("details:", taskDetails)
		/**
		Fetching task ticket type id with diff type id
		*/
		ttypes, err := dataAccess.Gettickettypebycatid(taskcat.Val)
		if err != nil {
			return errors.New("Something Went Wrong")
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
				return errors.New("Something Went Wrong")
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
					return errors.New("Something Went Wrong")
				}
				catentity.ID = difftype
				catentity.Val = cat
				catentities = append(catentities, catentity)
			}
			difftype1, err := dataAccess.Getdifftypebyid(taskcat.Val)
			if err != nil {
				return errors.New("Something Went Wrong")
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
				return errors.New("Something Went Wrong")
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
					seqno = 1
				}
				statusdet, err := dataAccess.Getdiffdetailsbyseq(data.ClientID, data.Mstorgnhirarchyid, seqno, 2)
				if err != nil {
					return errors.New("Something Went Wrong")
				}
				if len(statusdet) > 0 {
					log.Print("statusdet:", statusdet)
					logger.Log.Print("statusdet:", statusdet)
					/**
					Fetching priority details of task
					*/
					recEntity.Recordtypeid = ttypes[0].Val
					recEntity.Recordcatid = taskcat.Val
					prioentity := entities.RecordSet{}
					prioentity.ID = parentpriodifftypeid
					prioentity.Val = parentpriodiffid

					/**
					Generating the create record json to create the record
					*/
					recordEntity := entities.RecordEntity{}
					recordEntity.ClientID = data.ClientID
					recordEntity.Mstorgnhirarchyid = data.Mstorgnhirarchyid
					recordEntity.Originalusergroupid = originalusergroupid
					recordEntity.Originaluserid = originaluserid
					recordEntity.CreatedusergroupID = usergroupid
					recordEntity.CreateduserID = userid
					if taskDetails[0].Title != "" {
						recordEntity.Recordname = taskDetails[0].Title
					} else {
						recordEntity.Recordname = recordtitle
					}
					if taskDetails[0].Desc != "" {
						recordEntity.Recordesc = taskDetails[0].Desc
					} else {
						recordEntity.Recordesc = recorddescription
					}
					recordEntity.Requestername = requestername
					recordEntity.Requesteremail = requesteremail
					recordEntity.Requestermobile = requestermobile
					recordEntity.Requesterlocation = requesterlocation
					recordEntity.ParentID = data.RecordID
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
					recordEntity.Source = source
					recordEntity.Userid = data.UserID

					logger.Log.Print("\nfinal ", recordEntity)
					log.Print("\nfinal ", recordEntity)
					taskticketID, taskinsertedID, modelResponseError := AddRecordModelAction(&recordEntity, dbcon, ttypes[0].Val)
					if modelResponseError != nil {
						logger.Log.Print("\nError in task creation")
						log.Print("\nError in task creation")
						return errors.New("Something went wrong")

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
	return nil
}
