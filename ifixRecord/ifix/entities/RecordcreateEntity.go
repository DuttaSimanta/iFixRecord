package entities

import (
	"encoding/json"
	"io"
)

type RecordcreaterequestEntity struct {
	Clientid           int64 `json:"clientid"`
	Mstorgnhirarchyid  int64 `json:"mstorgnhirarchyid"`
	Recorddifftypeid   int64 `json:"recorddifftypeid"`
	Recorddiffid       int64 `json:"recorddiffid"`
	Recorddiffparentid int64 `json:"recorddiffparentid"`
	Recordtypeid       int64 `json:"recordtypeid"`
	Recordimpactid     int64 `json:"recordimpactid"`
	Recordurgencyid    int64 `json:"recordurgencyid"`
	Recordcatid        int64 `json:"recordcatid"`
	BaseConfig         int64 `json:"baseconfig"`
}

type RecordcreateEntity struct {
	Recordcategory          []RecordcategorydetailsEntity `json:"recordcategory"`
	Recordstatus            []RecordcatchildEntity        `json:"recordstatus"`
	Recordcatpos            int64                         `json:"recordcatpos"`
	Recordurgency           []RecordcatchildEntity        `json:"recordurgency"`
	Recordimpact            []RecordcatchildEntity        `json:"recordimpact"`
	Recordterms             []RecordtermlistEntity        `json:"recordterms"`
	Businessmatrixdirection int64                         `json:"configtype"`
	AssetAttached           int64                         `json:"isassetattached"`
	WorkingCatLabelID       int64                         `json:"workingcatlabelid"`
	AdditionalFields        []AdditionalFieldEntity       `json:"additionalfields"`
}

type RecordtermlistEntity struct {
	ID           int64  `json:"id"`
	Termname     string `json:"termname"`
	Termtypeid   int64  `json:"termtypeid"`
	Termtypename string `json:"termtypename"`
	Termvalue    string `json:"termvalue"`
}

type RecordtypeEntity struct {
	Recordtype []RecordtypedetailsEntity `json:"recordtype"`
}

type RecordcategorydetailsEntity struct {
	ID         int64                         `json:"id"`
	Title      string                        `json:"title"`
	Sequanceno int64                         `json:"sequanceno"`
	Child      []RecordcategorydetailsEntity `json:"child"`
	IsDisabled bool                          `json:"isDisabled"`
}

type RecordtypedetailsEntity struct {
	Typeid   int64  `json:"typeid"`
	Typename string `json:"typetitle"`
	Typeseq  int64  `json:"typeseqno"`
	ID       int64  `json:"id"`
	Name     string `json:"typename"`
	Seqno    int64  `json:"seqno"`
}

type RecordcatchildEntity struct {
	Typeid   int64  `json:"typeid"`
	Typename string `json:"typetitle"`
	Typeseq  int64  `json:"typeseqno"`
	ID       int64  `json:"id"`
	Name     string `json:"title"`
	Seqno    int64  `json:"seqno"`
}

type RecordcatchildNEstimatedEfforEntity struct {
	ChildCategories  []RecordcatchildEntity `json:"priority"`
	EstimatedEfforts []string               `json:"estimatedefforts"`
	SlaCompliance    []string               `json:"efficiency"`
	ChangeType       []string               `json:"changetype"`
}

//RecordcreateResponeData is final response structure of Record details
type RecordcreateResponeData struct {
	Status   bool               `json:"success"`
	Message  string             `json:"message"`
	Response RecordcreateEntity `json:"response"`
}

//RecordtypeResponeData is final response structure of Record details
type RecordtypeResponeData struct {
	Status   bool                      `json:"success"`
	Message  string                    `json:"message"`
	Response []RecordtypedetailsEntity `json:"response"`
}

//RecordcatchildResponeData is final response structure of Record details
type RecordcatchildResponeData struct {
	Status   bool                   `json:"success"`
	Message  string                 `json:"message"`
	Response []RecordcatchildEntity `json:"response"`
}

type RecordPriotiryResponeData struct {
	Status   bool                                `json:"success"`
	Message  string                              `json:"message"`
	Response RecordcatchildNEstimatedEfforEntity `json:"response"`
}

func (w *RecordcreaterequestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
