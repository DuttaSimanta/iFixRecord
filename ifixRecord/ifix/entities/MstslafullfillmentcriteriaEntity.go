package entities

import (
	"encoding/json"
	"io"
)

type MstslafullfillmentcriteriaEntity struct {
	Id                                   int64  `json:"id"`
	Clientid                             int64  `json:"clientid"`
	Mstorgnhirarchyid                    int64  `json:"mstorgnhirarchyid"`
	Slaid                                int64  `json:"slaid"`
	Mstrecorddifferentiationtickettypeid int64  `json:"mstrecorddifferentiationtickettypeid"`
	Mstrecorddifferentiationpriorityid   int64  `json:"mstrecorddifferentiationpriorityid"`
	Mstrecorddifferentiationworkingcatid int64  `json:"mstrecorddifferentiationworkingcatid"`
	Responsetimeinhr                     int64  `json:"responsetimeinhr"`
	Responsetimeinmin                    int64  `json:"responsetimeinmin"`
	Responsetimeinsec                    int64  `json:"responsetimeinsec"`
	Resolutiontimeinhr                   int64  `json:"resolutiontimeinhr"`
	Resolutiontimeinmin                  int64  `json:"resolutiontimeinmin"`
	Resolutiontimeinsec                  int64  `json:"resolutiontimeinsec"`
	Supportgroupspecific                 int64  `json:"supportgroupspecific"`
	Activeflg                            int64  `json:"activeflg"`
	Offset                               int64  `json:"offset"`
	Limit                                int64  `json:"limit"`
	Clientname                           string `json:"clientname"`
	Mstorgnhirarchyname                  string `json:"mstorgnhirarchyname"`
	Tickettypename                       string `json:"tickettypename"`
	Priorityname                         string `json:"priorityname"`
	Workingcatname                       string `json:"workingcatname"`
	Slaname                              string `json:"slaname"`
}

type MstslafullfillmentcriteriaEntities struct {
	Total  int64                              `json:"total"`
	Values []MstslafullfillmentcriteriaEntity `json:"values"`
}

type MstslafullfillmentcriteriaResponse struct {
	Success bool                               `json:"success"`
	Message string                             `json:"message"`
	Details MstslafullfillmentcriteriaEntities `json:"details"`
}

type MstslafullfillmentcriteriaResponseInt struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

func (w *MstslafullfillmentcriteriaEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
