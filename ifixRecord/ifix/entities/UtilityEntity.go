package entities

import "github.com/dgrijalva/jwt-go"

type UtilityEntity struct {
	Id                int64  `json:"id"`
	Clientid          int64  `json:"clientid"`
	Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
	Date              int64  `json:"date"`
	Timediff          int64  `json:"timediff"`
	Reporttimediff    int64  `json:"reporttimediff"`
	Reporttimeformat  string `json:"reporttimeformat"`
	Timeformat        string `json:"timeformat"`
}

type Claims struct {
	Userid int64 `json:"userid"`
	jwt.StandardClaims
}
