package models

import (
	"ifixRecord/ifix/dao"
	"ifixRecord/ifix/logger"
	"ifixRecord/ifix/utility"
	"log"
	"time"
)

func CheckToken(token string, userid int64) bool {
	// if mutexutility.MutexLocked(lock) == false {
	// 	lock.Lock()
	// 	defer lock.Unlock()
	// }
	// db, err := dbconfig.ConnectMySqlDb()
	// if err != nil {
	// 	log.Println("database connection failure", err)
	// 	return false
	// }
	//defer db.Close()
	var err error
	if db == nil {
		dbcon, err := ConnectMySqlDb()
		if err != nil {
			logger.Log.Println("Error in DBConnection")
			return false
		}
		db = dbcon
	}
	dataAccess := dao.DbConn{DB: db}
	dberr, tokens := dataAccess.Gettoken(userid)
	if dberr != nil {
		log.Println("database connection failure", err)
		return false
	}
	if len(tokens) > 0 {
		if token == tokens[0] {
			claims, success, err, errtype := utility.ValidateToken(token)
			if !success {
				log.Print("\n\n token error:", err, errtype)
				logger.Log.Print("\n\n token error:", err, errtype)
				if errtype != 2 {
					return false
				} else {
					if claims.Userid != userid {
						return false
					} else {
						if claims.ExpiresAt < time.Now().Local().Unix() {
							log.Print("\n\n token expired ")
							return false
						} else {
							return false
						}
					}
				}
			} else {
				if claims.Userid != userid {
					return false
				}
				return true
			}
		} else {
			log.Print("\n\n newtoken not matched with db token:", userid)
			logger.Log.Print("\n\n newtoken not matched with db token:", userid)
			return false
		}
	} else {
		log.Print("\n\n token not found:", userid)
		logger.Log.Print("\n\n token not found:", userid)
		return false
	}
}
