package entities

import (
	"encoding/json"
	"io"
)

type RecordstatusEntity struct {
	ClientID          int64 `json:"clientid"`
	Mstorgnhirarchyid int64 `json:"mstorgnhirarchyid"`
	RecordID          int64 `json:"recordid"`
	ReordstatusID     int64 `json:"reordstatusid"`
	UserID            int64 `json:"userid"`
	Usergroupid       int64 `json:"usergroupid"`
	Changestatus      int64 `json:"changestatus"`
	Issrrequestor     int64 `json:"issrrequestor"`
}

//FromJSON is used for convert data into JSON format
func (p *RecordstatusEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

//RecordResponeData is final response structure of Record details
type RecordstatusResponeData struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

type ParentchildEntity struct {
	Parentid       int64   `json:"parentid"`
	Childids       []int64 `json:"childids"`
	Userid         int64   `json:"userid"`
	Createdgroupid int64   `json:"createdgroupid"`
	Isupdate       bool    `json:"isupdate"`
	Transactionid  int64   `json:"transactionid"`
	Usergroupid    int64   `json:"usergroupid"`
	IsAttaching    int64   `json:"isattaching"`
}
