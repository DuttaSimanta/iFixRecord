package entities

import (
	"encoding/json"
	"io"
)

// type SLATabEntity struct {
// 	ClientID          int64 `json:"clientid"`
// 	Mstorgnhirarchyid int64 `json:"mstorgnhirarchyid"`
// 	RecordID          int64 `json:"recordid"`
// 	RecordtypeID      int64 `json:"recordtypeid"`
// 	WorkingcatID      int64 `json:"workingcatid"`
// 	PriorityID        int64 `json:"priorityid"`
// }

type SLATabEntity struct {
	ClientID          int64 `json:"clientid"`
	Mstorgnhirarchyid int64 `json:"mstorgnhirarchyid"`
	RecordID          int64 `json:"recordid"`
	RecordtypeID      int64 `json:"recordtypeid"`
	WorkingcatID      int64 `json:"workingcatid"`
	PriorityID        int64 `json:"priorityid"`
	SupportgroupId    int64 `json:"supportgroupid"`
}

func (p *SLATabEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type SLAResponsemeterEntity struct {
	RecordID            int64  `json:"recordid"`
	Priority            string `json:"priority"`
	Status              string `json:"status"`
	Responsetime        int64  `json:"responsetime"`
	Responseduetime     string `json:"responseduetime"`
	Responseslaviolated string `json:"responseslaviolated"`
	Responseclockstatus string `json:"responseclockstatus"`
}

type SLAHolidayEntity struct {
	Holiday   int64  `json:"holiday"`
	Starttime string `json:"starttime"`
	Endtime   string `json:"endtime"`
}

type SLAResolutionmeterEntity struct {
	RecordID              int64  `json:"recordid"`
	Priority              string `json:"priority"`
	Status                string `json:"status"`
	Resolutiontime        int64  `json:"resolutiontime"`
	Resolutionduetime     string `json:"resolutionduetime"`
	Resolutionslaviolated string `json:"resolutionslaviolated"`
	Resolutionclockstatus string `json:"resolutionclockstatus"`
}

type SLATabresponsesEntity struct {
	Responsedetails  SLAResponsemeterEntity   `json:"responsedetails"`
	Resolutionetails SLAResolutionmeterEntity `json:"resolutiondetails"`
	Holidaydetails   []SLAHolidayEntity       `json:"holidaydetails"`
}

type SLAMeterEntity struct {
	Remainresolutiontime int64   `json:"remainresolutiontime"`
	Resolutionpercent    float64 `json:"resolutionpercent"`
	Remainresponsetime   int64   `json:"remainresponsetime"`
	Responsepercent      float64 `json:"responsepercent"`
	RecordID             int64   `json:"recordid"`
}

//RecordcommonAllResponse is defined for response of API
type SLATabAllResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Details SLATabresponsesEntity `json:"details"`
}

type SLAMeterAllResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Details SLAMeterEntity `json:"details"`
}



type SLASchedulerEntity struct {
    ID                int64 `json:"id"`
    ClientID          int64 `json:"clientid"`
    Mstorgnhirarchyid int64 `json:"mstorgnhirarchyid"`
    RecordID          int64 `json:"recordid"`
}

type SLARcordFullfillUpdate struct {
    ClientID                     int64
    Mstorgnhirarchyid            int64
    RecordID                     int64
    Responseslameterpercentage   float64
    Resolutionslameterpercentage float64
    Businessaging                int64
    Calendaraging                int64
    Actualeffort                 int64
    Slaidletime                  int64
    Respoverduetime              int64
    Respoverdueperc              float64
    Resooverduetime              int64
    Resooverdueperc              float64
}
