package entities

import (
	"encoding/json"
	"io"
	"time"
)

type MstsladueEntity struct {
	Id                       int64
	Clientid                 int64
	Mstorgnhirarchyid        int64
	Mstslaentityid           int64
	Therecordid              int64
	Latestone                int64
	Startdatetimeresponse    string
	Startdatetimeresolution  string
	Duedatetimeresponse      string
	Duedatetimeresponseint   int64
	Duedatetimeresolution    string
	DuedatetimeresolutionInt int64
	Duedatetimetominute      int64
	Resoltiondone            int64
	Resolutiondatetime       string
	Lastupdatedattime        string
	Trnslaentityhistoryid    int64
	Remainingtime            int64
	Completepercent          float64
	Responseremainingtime    int64
	Responsepercentage       float64
	Isresponsecomplete       int
	Isresolutioncomplete     int
	ResponseCompleteTime     string
	Activeflg                int64
	Audittransactionid       int64
	Offset                   int64
	Limit                    int64
	PushTime                 int64
}

type ZoneEntity struct {
	Id      int64
	UTCdiff int64
}

type SLAEntity struct {
	Clientid          int64  `json:"clientid"`
	Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
	RecordID          int64  `json:"recordid"`
	Resolutiondone    int64  `json:"resolutiondone"`
	Fromclientuserid  int64  `json:"fromclientuserid"`
	Slastateid        int64  `json:"slastateid"`
	Createdate        string `json:"createdate"`
	Recordtypeid      int64  `json:"recordtypeid"`
	Workingdiffid     int64  `json:"workingdiffid"`
	Wecordpriorityid  int64  `json:"recordpriorityid"`
}

type SLAResponseEntity struct {
	SLAResponseStartTime    time.Time
	SLAResolutioneStartTime time.Time
	ResponseTimeDate        time.Time
	SLAResponseEndTime      time.Time
	ResolutionTimeDate      time.Time
	SLAResolutioneEndTime   time.Time
}

type SLAResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Details SLAEntity `json:"details"`
}

type MstsladueEntities struct {
	Total  int64               `json:"total"`
	Values []SLAResponseEntity `json:"values"`
}

type MstsladueResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Details MstsladueEntities `json:"details"`
}

type MstsladueResponseInt struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

func (w *MstsladueEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

func (w *SLAEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
