package entities

import (
	"encoding/json"
	"io"
)

type RecordDetailsRequestEntity struct {
	Clientid          int64   `json:"clientid"`
	Mstorgnhirarchyid int64   `json:"mstorgnhirarchyid"`
	Recordid          int64   `json:"recordid"`
	RecordDiffid      int64   `json:"recorddiffid"`
	RecordDiffTypeid  int64   `json:"recorddifftypeid"`
	RecordStageID     int64   `json:"recordstageid"`
	RecordNo          string  `json:"recordno"`
	ParentID          int64   `json:"parentid"`
	ChildID           int64   `json:"childid"`
	ChildIDS          []int64 `json:"childids"`
	Userid            int64   `json:"userid"`
	GroupID           int64   `json:"groupid"`
	TermsSeq          int64   `json:"termsseq"`
	BaseConfig        int64   `json:"baseconfig"`
}

type WorkFlowEntity struct {
	WorkFlowID int64 `json:"workflowid"`
	CatID      int64 `json:"catid"`
	CatTypeID  int64 `json:"cattypeid"`
}
type RecordDetailsEntity struct {
	Clientid             int64          `json:"clientid"`
	Mstorgnhirarchyid    int64          `json:"mstorgnhirarchyid"`
	Recordid             int64          `json:"recordid"`
	Title                string         `json:"title"`
	Description          string         `json:"description"`
	RecordTypeDiffTypeID int64          `json:"typedifftypeid"`
	RecordTypeID         int64          `json:"recordtypeid"`
	RecordType           string         `json:"recordtype"`
	GroupLevelID         int64          `json:"grouplevelid"`
	GroupLevel           string         `json:"grouplevel"`
	GroupID              int64          `json:"groupid"`
	Group                string         `json:"group"`
	PriorityID           int64          `json:"priorityid"`
	PriorityTypeID       int64          `json:"prioritytypeid"`
	Priority             string         `json:"priority"`
	StatusID             int64          `json:"statusid"`
	Status               string         `json:"status"`
	StatusSeqNo          int64          `json:"statusseqno"`
	ImpactID             int64          `json:"impactid"`
	Impact               string         `json:"impact"`
	UrgencyID            int64          `json:"urgencyid"`
	Urgency              string         `json:"urgency"`
	SourceType           string         `json:"source"`
	AssigneeID           int64          `json:"assigneeid"`
	Assignee             string         `json:"assignee"`
	RequestorInfo        string         `json:"requestorinfo"`
	AssignedGroupLevelID int64          `json:"assignedgrouplevelid"`
	AssignedGroupLevel   string         `json:"assignedgrouplevel"`
	AssignedGroupID      int64          `json:"assignedgroupid"`
	AssignedGroup        string         `json:"assignedgroup"`
	CreatedBy            string         `json:"createdby"`
	CreatorID            int64          `json:"creatorid"`
	ID                   int64          `json:"id"`
	Code                 string         `json:"code"`
	CreatedDateTime      string         `json:"createddatetime"`
	RecordStageID        int64          `json:"recordstageid"`
	WorkFlowDetails      WorkFlowEntity `json:"workflowdetails"`
	Vipuser              string         `json:"isvip"`
	Duedate              string         `json:"duedate"`
	RequestorName        string         `json:"requestername"`
	RequestorEmail       string         `json:"requesteremail"`
	RequestorMobile      string         `json:"requestermobile"`
	RequestorLocation    string         `json:"requesterlocation"`
	OrgRequestorName     string         `json:"orgrequestername"`
	OrgRequestorEmail    string         `json:"orgrequesteremail"`
	OrgRequestorMobile   string         `json:"orgrequestermobile"`
	OrgRequestorLocation string         `json:"orgrequesterlocation"`
	IsReslBreach         bool           `json:"resobreachcomment"`
	IsRespBreach         bool           `json:"respbreachcomment"`
	OriginalUserID       int64          `json:"originaluserid"`
	TypeSeqNo            int64          `json:"typeseqno"`
	Haspermission        bool           `json:"haspermission"`
	Isparent             string         `json:"isparent"`
	Ischild              string         `json:"ischild"`
       Isuserdelete           string         `json:"isuserdelete"`
	Isoriginaluseriddelete string         `json:"isoriginaluseriddelete"`
}

type RecordCatDetailsEntity struct {
	Recordcategory     []RecordcatEntity      `json:"recordcategory"`
	RecordCreateStatus []RecordcatchildEntity `json:"recordstatus"`
	Recordcatpos       int64                  `json:"recordcatpos"`
	Recordurgency      []RecordcatchildEntity `json:"recordurgency"`
	Recordimpact       []RecordcatchildEntity `json:"recordimpact"`
	RecordFields       []RecordFieldEntity    `json:"recordfields"`
	WorkFlowDetails    WorkFlowEntity         `json:"workflowdetails"`
	//Recordterms             []RecordtermlistEntity `json:"recordterms"`
	Businessmatrixdirection int64 `json:"configtype"`
	AssetAttached           int64 `json:"isassetattached"`
	//WorkingCatLabelID       int64                  `json:"workingcatlabelid"`
	EstimatedEfforts []string `json:"estimatedefforts"`
	SlaCompliances   []string `json:"slacomplainces"`
	ChangeTypes      []string `json:"changetypes"`
	Haspermission    bool     `json:"haspermission"`
}

type RecordFieldEntity struct {
	FieldID       int64  `json:"fieldid"`
	TermsID       int64  `json:"termsid"`
	TermsName     string `json:"termsname"`
	TermsValue    string `json:"termsvalue"`
	TermsTypeID   int64  `json:"termstypeid"`
	TermsTypeName string `json:"termstypename"`
	Value         string `json:"value"`
	CatSeq        string `json:"catSeq"`
	TermSeqNo     int64  `json:"seqno"`
}

type RecordcatEntity struct {
	ID         int64             `json:"id"`
	Title      string            `json:"title"`
	Sequanceno int64             `json:"sequanceno"`
	IsSelected int64             `json:"selected"`
	Child      []RecordcatEntity `json:"child"`
	IsDisabled bool              `json:"isDisabled"`
}

type RawDiffEntity struct {
	Typeid   int64  `json:"typeid"`
	Typename string `json:"typetitle"`
	Typeseq  int64  `json:"typeseqno"`
	ID       int64  `json:"id"`
	Name     string `json:"title"`
	Seqno    int64  `json:"seqno"`
	Selected int64  `json:"selected"`
	ParentID int64  `json:"parentid"`
}

type RecordDetailsResponeData struct {
	Status  bool                  `json:"success"`
	Message string                `json:"message"`
	Details []RecordDetailsEntity `json:"details"`
}

type RecordCatDetailsResponeData struct {
	Status  bool                   `json:"success"`
	Message string                 `json:"message"`
	Details RecordCatDetailsEntity `json:"details"`
}

func (w *RecordDetailsRequestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

type RecordDetailsParentResponeData struct {
	Status  bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

type ChildRecordSearchEntity struct {
	Clientid          int64  `json:"clientid"`
	Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
	Recordid          int64  `json:"recordid"`
	RecordDiffid      int64  `json:"recorddiffid"`
	RecordDiffTypeid  int64  `json:"recorddifftypeid"`
	RecordNo          string `json:"recordno"`
	GroupID           int64  `json:"groupid"`
	Requestername     string `json:"requestername"`
	RequesterID       string `json:"requesterid"`
	Requesterlocation string `json:"requesterlocation"`
	ShortDescription  string `json:"shortdescription"`
	Priority          int64  `json:"priority"`
	RecordStageID     int64  `json:"recordstageid"`
	Fromdate          string `json:"fromdate"`
	Todate            string `json:"todate"`
	CategorylabelID   int64  `json:"categorylabelid"`
	CategoryID        int64  `json:"categoryid"`
}

func (w *ChildRecordSearchEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

type RecordAccessEntity struct {
	Clientid          int64  `json:"clientid"`
	Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
	RecordNo          string `json:"recordno"`
	Haspermission     bool   `json:"haspermission"`
}

type RecordAccessDetailsResponeData struct {
	Status  bool               `json:"success"`
	Message string             `json:"message"`
	Details RecordAccessEntity `json:"details"`
}
