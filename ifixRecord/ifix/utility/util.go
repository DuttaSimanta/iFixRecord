package utility

import (
	"ifixRecord/ifix/dbconfig"
	"ifixRecord/ifix/entities"
	"log"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//var insertlogs = "INSERT INTO mstrecordactivitylogs(clientid,mstorgnhirarchyid,recordid,activityseqno,logValue,createdid,createddate,createdgrpid) VALUES (?,?,?,?,?,?,round(UNIX_TIMESTAMP(now())),?)"

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func ItemExists(arrayType interface{}, item interface{}) (bool, int) {
	arr := reflect.ValueOf(arrayType)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true, i
		}
	}
	return false, -1
}

// func InsertActivityLogs(tx *sql.Tx, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, ActivityID int64, Logvalue string, CreatedID int64, CreatedgrpID int64) error {
// 	logger.Log.Println("InsertActivityLogs query -->", insertlogs)
// 	logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID)

// 	stmt, err := tx.Prepare(insertlogs)

// 	if err != nil {
// 		logger.Log.Println(err)
// 		return err
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, ActivityID, Logvalue, CreatedID, CreatedgrpID)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return err
// 	}

// 	return nil
// }

func CreateToken(userid int64) (string, error) {
	claims := &entities.Claims{
		Userid: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(dbconfig.TOKEN_EXPIRE_TIME)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(dbconfig.JWT_KEY)
	return tokenString, err
}
func ValidateToken(token string) (*entities.Claims, bool, error, int) {
	var claims = &entities.Claims{}
	_, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(dbconfig.JWT_KEY), nil
		})
	if err != nil {
		log.Print("\n\n err:", err)
		if err == jwt.ErrSignatureInvalid {
			return claims, false, err, 1
		}
		return claims, false, err, 2
	}
	//if !tkn.Valid {
	//	return claims, false, nil, 3
	//}
	return claims, true, nil, 0
}

// func CheckToken(token string, userid int64) bool {
// 	db, err := dbconfig.ConnectMySqlDb()
// 	if err != nil {
// 		log.Println("database connection failure", err)
// 		return false
// 	}
// 	defer db.Close()
// 	dataAccess := dao.DbConn{DB: db}
// 	dberr, tokens := dataAccess.Gettoken(userid)
// 	if dberr != nil {
// 		log.Println("database connection failure", err)
// 		return false
// 	}
// 	if len(tokens) > 0 {
// 		if token == tokens[0] {
// 			claims, success, err, errtype := ValidateToken(token)
// 			if !success {
// 				log.Print("\n\n token error:", err, errtype)
// 				logger.Log.Print("\n\n token error:", err, errtype)
// 				if errtype != 2 {
// 					return false
// 				} else {
// 					if claims.Userid != userid {
// 						return false
// 					} else {
// 						if claims.ExpiresAt < time.Now().Local().Unix() {
// 							log.Print("\n\n token expired ")
// 							return false
// 						} else {
// 							return false
// 						}
// 					}
// 				}
// 			} else {
// 				if claims.Userid != userid {
// 					return false
// 				}
// 				return true
// 			}
// 		} else {
// 			log.Print("\n\n newtoken not matched with db token:", userid)
// 			logger.Log.Print("\n\n newtoken not matched with db token:", userid)
// 			return false
// 		}
// 	} else {
// 		log.Print("\n\n token not found:", userid)
// 		logger.Log.Print("\n\n token not found:", userid)
// 		return false
// 	}
// }
