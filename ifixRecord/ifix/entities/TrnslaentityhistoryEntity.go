package entities

import (
	"encoding/json"
	"io"
)

type TrnslaentityhistoryEntity struct {
	Id                    int64
	Clientid              int64
	Mstorgnhirarchyid     int64
	Mstslaentityid        int64
	Therecordid           int64
	Recorddatetime        string
	Recorddatetoint       int64
	Donotupdatesladue     int64
	Recordtimetoint       int64
	Mstslastateid         int64
	Commentonrecord       string
	Slastartstopindicator int64
	Fromclientuserid      int64
	PushTime              int64
	Activeflg             int64
	Audittransactionid    int64
	Offset                int64
	Limit                 int64
}

type TrnslaentityhistoryLastPushEntity struct {
	Recorddatetime  string
	Recorddatetoint int64
}

type TrnslaentityhistoryEntities struct {
	Total  int64                       `json:"total"`
	Values []TrnslaentityhistoryEntity `json:"values"`
}

type TrnslaentityhistoryResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Details TrnslaentityhistoryEntities `json:"details"`
}

type TrnslaentityhistoryResponseInt struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

func (w *TrnslaentityhistoryEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
