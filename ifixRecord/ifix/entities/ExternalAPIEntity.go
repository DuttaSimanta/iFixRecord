package entities

import (
	"encoding/json"
	"io"
)
type FileAttacmentToRecordEntity struct {
    Clientname          string                        `json:"clientname"`
    Mstorgnhirarchyname string                        `json:"mstorgnhirarchyname"`
    Recordid            string                        `json:"recordid"`
    LoginID             string                        `json:"loginid"`
    Userid              int64                         `json:"userid"`
    Logingrpname        string                        `json:"logingrpname"`
    Fileattachment      []FileAttachmentDetailsEntity `json:"fileattachment"`
}

type FileAttachmentDetailsEntity struct {
    Filename    string `json:"filename"`
    Filetype    string `json:"filetype"`
    Filecontent string `json:"content"`
}

func (w *FileAttacmentToRecordEntity) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(w)
}

type FileuploadResponse struct {
    Success bool             `json:"success"`
    Message string           `json:"message"`
    Details FileuploadEntity `json:"details"`
}

type FileuploadEntity struct {
    Id                 int64  `json:"id"`
    Clientid           int64  `json:"clientid"`
    Mstorgnhirarchyid  int64  `json:"mstorgnhirarchyid"`
    Credentialtype     string `json:"credentialtype"`
    Credentialaccount  string `json:"credentialaccount"`
    Credentialpassword string `json:"credentialpassword"`
    Credentialkey      string `json:"credentialkey"`
    Activeflg          int64  `json:"activeflg"`
    Originalfile       string `json:"originalfile"`
    Filename           string `json:"filename"`
    Path               string `json:"path"`
}

//==========================================
type RecordDetailsEntityAPI struct {
	Clientname                      string                   `json:"clientname"`
	Mstorgnhirarchyname             string                   `json:"mstorgnhirarchyname"`
	Recordid                        string                   `json:"recordid"`
	Title                           string                   `json:"title"`
	Description                     string                   `json:"description"`
	RecordType                      string                   `json:"recordtype"`
	Priority                        string                   `json:"priority"`
	Status                          string                   `json:"status"`
	Impact                          string                   `json:"impact"`
	Urgency                         string                   `json:"urgency"`
	RequestorName                   string                   `json:"requestername"`
	Requestorloginid                string                   `json:"requesterloginid"`
	RequestorEmail                  string                   `json:"requesteremail"`
	RequestorMobile                 string                   `json:"requestermobile"`
	RequestorLocation               string                   `json:"requesterlocation"`
	Assignee                        string                   `json:"assignee"`
	AssignedGroup                   string                   `json:"assignedgroup"`
	Prioritycount                   int64                    `json:"prioritycount"`
	Followupcount                   int64                    `json:"followupcount"`
	Reopencount                     int64                    `json:"reopencount"`
	Outboundcount                   int64                    `json:"outboundcount"`
	Aging                           int64                    `json:"aging"`
	Source                          string                   `json:"source"`
	Visiblecommentdaycount          int64                    `json:"visiblecommentdaycount"`
	Hopcount                        int64                    `json:"hopcount"`
	Latestupdatedby                 string                   `json:"latestupdatedby"`
	Createdatetime                  string                   `json:"createdatetime"`
	Termsdetails                    map[string]interface{}   `json:"termsdetails"`
	Categories                      []Categorydetails        `json:"categories"`
	OriginalCreatedByLoginID        string                   `json:"originalcreatedbyloginid"`
	OriginalCreatedByFullName       string                   `json:"originalcreatedbyfullname"`
	OriginalCreatedByPrimaryContact string                   `json:"originalcreatedbyprimarycontact"`
	Isvipuser                       string                   `json:"isvipuser"`
	ResolvedByGroup                 string                   `json:"resolvedbygroup"`
	ResolvedByUser                  string                   `json:"resolvedbyuser"`
	ResponseDueDate                 string                   `json:"responseduedate"`
	ResolutionDueDate               string                   `json:"resolutionduedate"`
	LastModifiedDateTime            string                   `json:"lastmodifieddatetime"`
	ResponseClockStatus             string                   `json:"responseclockstatus"`
	ResponseSLABreachedStatus       string                   `json:"responseslabreachedstatus"`
	ResponseSLAOverdue              string                   `json:"responseslaoverdue"`
	ResolutionClockStatus           string                   `json:"resolutionclockstatus"`
	ResolutionSLABreachedStatus     string                   `json:"resolutionslabreachedstatus"`
	ResolutionSLAOverdue            string                   `json:"resolutionslaoverdue"`
	LastModifiedDate                string                   `json:"lastmodifieddate"`
	Linkedtickets                   []string                 `json:"linkedtickets"`
	Responseslacompliance           string                   `json:"responseslacompliance"`
	Resolutionslacompliance         string                   `json:"resolutionslacompliance"`
	TotalEffort                     string                   `json:"totaleffort"`
	IsParent                        string                   `json:"isparent"`
	ChildCount                      int64                    `json:"childcount"`
	Childdetails                    []RecordDetailsEntityAPI `json:"childdetails"`
	IsChild                         string                   `json:"ischild"`
	ParentCount                     int64                    `json:"parentcount"`
	Parentdetails                   []RecordDetailsEntityAPI `json:"parentdetails"`
}

type RecordDetailsRequestEntityAPI struct {
	Clientname          string `json:"clientname"`
	Mstorgnhirarchyname string `json:"mstorgnhirarchyname"`
	RecordNo            string `json:"recordno"`
	Fromdate            string `json:"fromdate"`
	Todate              string `json:"todate"`
	Userid              int64  `json:"userid"`
}

func (w *RecordDetailsRequestEntityAPI) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

type Categorydetails struct {
	Label        string `json:"label"`
	Categoryname string `json:"Categoryname"`
}

type Workflowdetails struct {
	Asigneename string
	Asigneegrp  string
}

type RecordDetailsEntityResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Details RecordDetailsEntityAPI `json:"details"`
}
type TokenRecordDetailsEntityResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Details []RecordDetailsEntityAPI `json:"details"`
}

type ExternalRecordEntityResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type ExternalCreateRecord struct {
	Clientname               string                     `json:"clientname"`
	Mstorgnhirarchyname      string                     `json:"mstorgnhirarchyname"`
	ShortDescription         string                     `json:"shortDescription"`
	LongDescription          string                     `json:"longDescription"`
	LoginID                  string                     `json:"loginid"`
	LoginGrpname             string                     `json:"logingrpname"`
	ExternalRecordSets       []ExternalRecordSet        `json:"recordsets"`
	ExternalAdditionalfields []ExternalRecordAdditional `json:"additionalfields"`
	Recordid                 string                     `json:"recordid"`
	Statusname               string                     `json:"statusname"`
	AssigneGroupname         string                     `json:"assigneegrpname"`
	AssigneUserID            string                     `json:"assigneeuserid"`

	// For Create ticket
	RequestorID     string `json:"requestorid"`
	OriginalID      string `json:"originalid"`
	Originalgrpname string `json:"originalgrpname"`
	Userid          int64  `json:"userid"`
	TickettypeID    string `json:"tickettypeid"`
}

type ExternalRecordSet struct {
	Typename string               `json:"typename"`
	Type     []ExternalRecordData `json:"type"`
	Value    string               `json:"value"`
}

type ExternalRequestorInfo struct {
	Requestername     string `json:"requestername"`
	Requesteremail    string `json:"requesteremail"`
	Requestermobile   string `json:"requestermobile"`
	Requesterlocation string `json:"requesterlocation"`
}

func (w *ExternalCreateRecord) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

type ExternalRecordData struct {
	Labelname  string `json:"labelname"`
	Labelvalue string `json:"labelvalue"`
}

type ExternalRecordAdditional struct {
	Termname string `json:"termname"`
	Value    string `json:"value"`
}

type RequestBody struct {
	ClientID             int64 `json:"clientid"`
	MstorgnhirarchyID    int64 `json:"mstorgnhirarchyid"`
	RecorddifftypeID     int64 `json:"recorddifftypeid"`
	RecorddiffID         int64 `json:"recorddiffid"`
	PreviousstateID      int64 `json:"previousstateid"`
	CurrentstateID       int64 `json:"currentstateid"`
	TransactionID        int64 `json:"transactionid"`
	CreatedgroupID       int64 `json:"createdgroupid"`
	MstgroupID           int64 `json:"mstgroupid"`
	MstuserID            int64 `json:"mstuserid"`
	Manualstateselection int64 `json:"manualstateselection"`
	//Samegroup         bool  `json:"samegroup"`
        UserID            int64 `json:"userid"`
}

type RequestBody1 struct {
	ClientID          int64 `json:"clientid"`
	MstorgnhirarchyID int64 `json:"mstorgnhirarchyid"`
	RecorddifftypeID  int64 `json:"recorddifftypeid"`
	RecorddiffID      int64 `json:"recorddiffid"`
	PreviousstateID   int64 `json:"previousstateid"`
	CurrentstateID    int64 `json:"currentstateid"`
	TransactionID     int64 `json:"transactionid"`
	CreatedgroupID    int64 `json:"createdgroupid"`
	MstgroupID        int64 `json:"mstgroupid"`
	MstuserID         int64 `json:"mstuserid"`
	Samegroup         bool  `json:"samegroup"`
	UserID            int64 `json:"userid"`
}
