package entities

import (
	"encoding/json"
	"io"
)

//RecordEntity contains all required data fields
type RecordEntity struct {
	ID                  int64              `json:"id"`
	ClientID            int64              `json:"clientid"`
	Mstorgnhirarchyid   int64              `json:"mstorgnhirarchyid"`
	Requesterinfo       string             `json:"requesterinfo"`
	RecordTypeSeq       int64              `json:"recordtypeseq"`
	RecordTypeID        int64              `json:"recordtypeid"`
	Recordname          string             `json:"recordname"`
	Recordesc           string             `json:"recordesc"`
	RecordPriorityID    string             `json:"recordpriorityid"`
	Recordattachpath    string             `json:"recordattachpath"`
	Recordcategorydtls  string             `json:"recordcategorydtls"`
	Recordpsdetails     string             `json:"recordpsdetails"`
	Recordsourcetype    string             `json:"recordsourcetype"`
	RecordcreatedID     int64              `json:"recordcreatedid"`
	RecordoriginalID    int64              `json:"recordoriginalid"`
	Recordrequestinfo   string             `json:"recordrequestinfo"`
	RecordSets          []RecordSet        `json:"recordsets"`
	Recordfields        []RecordField      `json:"recordfields"`
	ResponsedifftypeID  []int64            `json:"responsedifftypeid"`
	Userid              int64              `"json:userid"`
	Usergroupid         int64              `json:"usergroupid"`
	Originaluserid      int64              `json:"originaluserid"`
	Originalusergroupid int64              `json:"originalusergroupid"`
	AssetIds            []int64            `json:"assetIds"`
	Workingcatlabelid   int64              `json:"workingcatlabelid"`
	Additionalfields    []RecordAdditional `json:"additionalfields"`
	Requestername       string             `json:"requestername"`
	Requesteremail      string             `json:"requesteremail"`
	Requestermobile     string             `json:"requestermobile"`
	Requesterlocation   string             `json:"requesterlocation"`
	Source              string             `json:"source"`
	CreateduserID       int64              `"json:createduserid"`
	CreatedusergroupID  int64              `json:"createdusergroupid"`
	Lastlevelcatid      int64              `json:"lastlevelcatid"`
	ParentID            int64              `json:"parentid"`
	RecordIds           []int64            `json:"recordids"`
}

type TaskdetailsEntity struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type RecordcategoryupdateEntity struct {
	ClientID          int64              `json:"clientid"`
	Mstorgnhirarchyid int64              `json:"mstorgnhirarchyid"`
	Recorddifftypeid  int64              `json:"recorddifftypeid"`
	Recorddiffid      int64              `json:"recorddiffid"`
	RecordID          int64              `json:"recordid"`
	UserID            int64              `"json:userid"`
	UsergroupID       int64              `json:"usergroupid"`
	Total             int64              `json:"total"`
	RecordSets        []RecordSet        `json:"recordsets"`
	Recordfields      []RecordField      `json:"recordfields"`
	Workingcatlabelid int64              `json:"workingcatlabelid"`
	Additionalfields  []RecordAdditional `json:"additionalfields"`
}

//RecordAllResponse is defined for response of API
type RecordAllResponse struct {
	Status   bool           `json:"success"`
	Message  string         `json:"message"`
	Response []RecordEntity `json:"response"`
}

//FromJSON is used for convert data into JSON format
func (p *RecordEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

//FromJSON is used for convert data into JSON format
func (p *RecordcategoryupdateEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

//RecordData is a details value of Recordsets entity
type RecordData struct {
	ID  int64 `json:"id"`
	Val int64 `json:"val"`
}

//RecordSet is a details value of Recordsets entity
type RecordSet struct {
	ID   int64        `json:"id"`
	Type []RecordData `json:"type"`
	Val  int64        `json:"val"`
}

//RecordAdditional is a details value of RecordAdditional entity
type RecordAdditional struct {
	ID      int64  `json:"id"`
	Termsid int64  `json:"termsid"`
	Val     string `json:"val"`
}

//RecordField is a details value of Recordsets entity
type RecordField struct {
	// DifftypeID int64        `json:"difftypeid"`
	// DiffID     int64        `json:"diffid"`
	// Terms      []RecordTerm `json:"terms"`
	TermID int64        `json:"termid"`
	Val    []RecordTerm `json:"val"`
}

//RecordTerm is a details value of Recordsets entity
type RecordTerm struct {
	OriginalName string `json:"originalName"`
	FileName     string `json:"fileName"`
}

//RecordRespone is response of Record details
type RecordRespone struct {
	ID                int64            `json:"id"`
	ClientID          int64            `json:"clientid"`
	Mstorgnhirarchyid int64            `json:"mstorgnhirarchyid"`
	ResponseDetails   []ResponseDetail `json:responsedetails`
}

//ResponseDetail is response of Record details
type ResponseDetail struct {
	DifftypeID int64 `json:"difftypeid"`
	DiffID     int64 `json:"diffid"`
}

//RecordResponeData is final response structure of Record details
type RecordResponeData struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Response string `json:"response"`
	ID       int64  `json:"id"`
}

type WorkflowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type RecordpriorityEntity struct {
	ClientID            int64  `json:"clientid"`
	Mstorgnhirarchyid   int64  `json:"mstorgnhirarchyid"`
	Recorddifftypeid    int64  `json:"recorddifftypeid"`
	Recorddiffid        int64  `json:"recorddiffid"`
	RecordID            int64  `json:"recordid"`
	Userid              int64  `"json:userid"`
	Usergroupid         int64  `json:"usergroupid"`
	Originaluserid      int64  `json:"originaluserid"`
	Originalusergroupid int64  `json:"originalusergroupid"`
	Recordname          string `json:"recordname"`
	Recordesc           string `json:"recordesc"`
}

//RecordResponeData is final response structure of Record details
type RecordPriorityResponeEntity struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	StageID int64  `json:"stageid"`
}

//FromJSON is used for convert data into JSON format
func (p *RecordpriorityEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type StagetableEntity struct {
	ClientID          int64
	OrgnID            int64
	RecordID          int64
	TicketID          string
	Source            string
	Requestorloginid  string
	RequestorID       int64
	Requestorname     string
	Requestorlocation string
	Requestorphone    string
	Requestoremail    string

	Orgcreatorlocation string
	Orgcreatorphone    string
	Orgcreatorid       int64
	Orgcreatorname     string
	Orgcreatorloginid  string
	Orgcreatoremail    string

	Tickettype   string
	Tickettypeid int64
	Priorityid   int64
	Priority     string
	Statusid     int64
	Status       string
	Vipticket    string
	Urgencyid    int64
	Urgency      string
	Impactid     int64
	Impact       string

	Shortdescription             string
	assigneduserloginid          int64
	createddatetime              int64
	createddate                  int64
	lastupdateddatetime          int64
	lastupdateddate              int64
	reopencount                  int64
	reassigncount                int64
	prioritycount                int64
	followupcount                int64
	outboundcount                int64
	isparent                     string
	childcount                   int64
	parentticketid               int64
	responseslameterpercentage   int64
	resolutionslameterpercentage int64
	worknotenotupdated           int64
	LastuserID                   int64
	Lastusername                 string
	Fstlevelcategorynm           string
}
