package entities

import (
	"encoding/json"
	"io"
)
type BulkRecordcommonEntity struct {
	RecordID    []int64 `json:"recordid"`
	Termvalue   string  `json:"termvalue"`
	Userid      int64   `"json:userid"`
	Usergroupid int64   `json:"usergroupid"`
	Termseq     int64   `json:"termseq"`
}
func (p *BulkRecordcommonEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

//RecordcommonEntity contains all required data fields
type RecordcommonEntity struct {
	ClientID           int64   `json:"clientid"`
	Mstorgnhirarchyid  int64   `json:"mstorgnhirarchyid"`
	RecordID           int64   `json:"recordid"`
	RecordstageID      int64   `json:"recordstageid"`
	TermID             int64   `json:"termid"`
	Termvalue          string  `json:"termvalue"`
	ForuserID          int64   `json:"foruserid"`
	Recorddifftypeid   int64   `json:"recorddifftypeid"`
	Recorddiffid       int64   `json:"recorddiffid"`
	Termdescription    string  `json:"termdescription"`
	Userid             int64   `"json:userid"`
	Usergroupid        int64   `json:"usergroupid"`
	Termseq            int64   `json:"termseq"`
	ID                 int64   `json:"id"`
	Sequance           []int64 `json:"sequance"`
	Recordno           string  `json:"recordno"`
	GrpID              int64   `json:"grpid"`
	Recordstatustypeid int64   `json:"recordstatustypeid"`
	Recordstatusid     int64   `json:"recordstatusid"`
}

type RecordmultiplecommonEntity struct {
	ClientID          int64                   `json:"clientid"`
	Mstorgnhirarchyid int64                   `json:"mstorgnhirarchyid"`
	RecordID          int64                   `json:"recordid"`
	RecordstageID     int64                   `json:"recordstageid"`
	Details           []RecordTermnamesEntity `json:"details"`
	ForuserID         int64                   `json:"foruserid"`
	Userid            int64                   `"json:userid"`
	Recorddifftypeid  int64                   `json:"recorddifftypeid"`
	Recorddiffid      int64                   `json:"recorddiffid"`
	Usergroupid       int64                   `json:"usergroupid"`
}

type RecordTermvaluesEntity struct {
	TermID    int64  `json:"termid"`
	Termvalue string `json:"termvalue"`
	ForuserID int64  `json:"foruserid"`
}

type RecordcommonstateEntity struct {
	ClientID                   int64 `json:"clientid"`
	Mstorgnhirarchyid          int64 `json:"mstorgnhirarchyid"`
	Recordtickettypedifftypeid int64 `json:"recordtickettypedifftypeid"`
	Recordtickettypediffid     int64 `json:"recordtickettypediffid"`
	Recordstatusdifftypeid     int64 `json:"recordstatusdifftypeid"`
	Recordstatusdiffid         int64 `json:"recordstatusdiffid"`
	Userid                     int64 `"json:userid"`
}

type RecordTermnamesEntity struct {
	ID              int64  `json:"id"`
	Termname        string `json:"tername"`
	Recordtermvalue string `json:"recordtermvalue"`
	Iscompulsory    int64  `json:"iscompulsory"`
	Termtypename    string `json:"termtypename"`
	Termtypeid      int64  `json:"termtypeid"`
	Insertedvalue   string `json:"insertedvalue"`
	Seq             int64  `json:"seq"`
	Termdescription string `json:"termdescription"`
}

type RecordcommonresponseEntity struct {
	ID             int64  `json:"id"`
	TermID         int64  `json:"termid"`
	Termvalue      string `json:"termvalue"`
	ForuserID      string `json:"foruserid"`
	Termname       string `json:"termname"`
	Recorddiffname string `json:"recorddiffname"`
	Username       string `json:"username"`
	Createddate    int64  `json:"createddate"`
}

type RecordScheduleTabTermnamesEntity struct {
	ID              int64  `json:"id"`
	Termname        string `json:"tername"`
	Recordtermvalue string `json:"recordtermvalue"`
	Iscompulsory    int64  `json:"iscompulsory"`
	Termtypename    string `json:"termtypename"`
	Termtypeid      int64  `json:"termtypeid"`
	Readpermission  int64  `json:"readpermission"`
	Writepermission int64  `json:"writepermission"`
	FieldID         int64  `json:"fieldid"`
	Val             string `json:"val"`
	Seq             int    `json:"seq"`
}

type RecordPlanTabTermnamesEntity struct {
	ID              int64  `json:"id"`
	Termname        string `json:"tername"`
	Recordtermvalue string `json:"recordtermvalue"`
	Iscompulsory    int64  `json:"iscompulsory"`
	Termtypename    string `json:"termtypename"`
	Termtypeid      int64  `json:"termtypeid"`
	Readpermission  int64  `json:"readpermission"`
	Writepermission int64  `json:"writepermission"`
	FieldID         int64  `json:"fieldid"`
	Val             string `json:"val"`
	Seq             int    `json:"seq"`
}

type RecordTabTermsEntity struct {
	ScheduleTab []RecordScheduleTabTermnamesEntity `json:"scheduletab"`
	PlanTab     []RecordPlanTabTermnamesEntity     `json:"plantab"`
}

type ChildRecordEntity struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

//FromJSON is used for convert data into JSON format
func (p *RecordcommonEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *RecordcommonstateEntity) FromstateJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *RecordmultiplecommonEntity) FrommultipleJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

//RecordcommonAllResponse is defined for response of API
type RecordcommonAllResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Details []RecordcommonresponseEntity `json:"details"`
}

type RecordcommonResponseInt struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

//RecordTermnames is defined for response of API
type RecordTermnamesAllResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Details []RecordTermnamesEntity `json:"details"`
}

type RecordTabTermnamesAllResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Details RecordTabTermsEntity `json:"details"`
}

type RecordcountAllResponse struct {
	Success                  bool   `json:"success"`
	Message                  string `json:"message"`
	Prioritycount            int64  `json:"prioritycount"`
	Followupcount            int64  `json:"followupcount"`
	Reopencount              int64  `json:"reopencount"`
	Pendingvendoractioncount int64  `json:"pendingvendoractioncount"`
	Outboundcount            int64  `json:"outboundcount"`
	Aging                    int64  `json:"aging"`
}

type RecordcountEntity struct {
	Prioritycount            int64 `json:"prioritycount"`
	Followupcount            int64 `json:"followupcount"`
	Reopencount              int64 `json:"reopencount"`
	Pendingvendoractioncount int64 `json:"pendingvendoractioncount"`
	Outboundcount            int64 `json:"outboundcount"`
	Aging                    int64 `json:"aging"`
}

type RecentrecordEntity struct {
	ID             int64  `json:"id"`
	Code           string `json:"code"`
	Title          string `json:"title"`
	Status         string `json:"status"`
	Seq            int64  `json:"seq"`
	Createdate     int64  `json:"createdate"`
	Showcreatedate string `json:"showcreatedate"`
	OrgnID         int64  `json:"orgnid"`
}

type RecentrecordAllResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Details []RecentrecordEntity `json:"details"`
}

type RecordlogsEntity struct {
	ID               int64  `json:"id"`
	Logvalue         string `json:"logvalue"`
	Name             string `json:"name"`
	Activitydesc     string `json:"activitydesc"`
	Createdate       int64  `json:"createdate"`
	Supportgroupname string `json:"supportgroupname"`
}

type RecordlogsAllResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Details []RecordlogsEntity `json:"details"`
}

type FrequentRecordEntity struct {
	Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
	LastlevelID      int64  `json:"lastlevelid"`
	ParentcatID      string `json:"parentcatid"`
	Parentcatname    string `json:"parentcatname"`
	Count            int64  `json:"count"`
	Recorddifftypeid int64  `json:"recorddifftypeid"`
	Recorddiffid     int64  `json:"recorddiffid"`
	Seq              int64  `json:"seq"`
}

type FrequentrecordsAllResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Details []FrequentRecordEntity `json:"details"`
}

type ParentticketEntity struct {
	ID           int64  `json:"id"`
	Recordnumber string `json:"recordnumber"`
}

type ParentticketEntityAllResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Details []ParentticketEntity `json:"details"`
}

type Recordactivitymst struct {
	ID          int64                   `json:"id"`
	Description string                  `json:"description"`
	Seq         int64                   `json:"seq"`
	LogType     string                  `json:"logtype"`
	Details     []RecordTermnamesEntity `json:"details"`
}
type ActivitynamesAllResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Details []Recordactivitymst `json:"details"`
}

type NewActivitylogsEntity struct {
	ID               int64  `json:"id"`
	RecordID         int64  `json:"recordid"`
	Logvalue         string `json:"logvalue"`
	Createddate      int64  `json:"createddate"`
	Name             string `json:"name"`
	Activitydesc     string `json:"activitydesc"`
	Supportgroupname string `json:"supportgroupname"`
	Termname         string `json:"termname"`
	Status           string `json:"status"`
	Termdescription  string `json:"termdescription"`
	Showcreatedate   string `json:"showcreatedate"`
	Code             string `json:"code"`
}

type NewActivityAllResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Details []NewActivitylogsEntity `json:"details"`
}

//InsertTermvalue,InsertMultipleTermvalue

type Activitylogsearchcriteria struct {
	RecordID          int64           `json:"recordid"`
	Recordcode        string          `json:"recordcode"`
	ClientID          int64           `json:"clientid"`
	Mstorgnhirarchyid int64           `json:"mstorgnhirarchyid"`
	Userid            int64           `"json:userid"`
	Usergroupid       int64           `json:"usergroupid"`
	Searchfilter      []Searchfilters `json:"searchfilter"`
	Searchbyseq       []int64         `json:"searchbyseq"`
}

func (p *Activitylogsearchcriteria) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type Searchfilters struct {
	ID  int64 `json:"id"`
	Seq int64 `json:"seq"`
}

type Searchresults struct {
	Total  int                     `json:"total"`
	Values []NewActivitylogsEntity `json:"values"`
}

type ActivitySearchAllResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Details []NewActivitylogsEntity `json:"details"`
}

type Pendingstatustermvalue struct {
	Termname  string `json:"termname"`
	Termvalue string `json:"termvalue"`
	Seq       int64  `json:"seq"`
	ID        int64  `json:"id"`
}

type PendingstatusAllResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Details []Pendingstatustermvalue `json:"details"`
}

type RecordAttachmentfiles struct {
	ID                      int64  `json:"id"`
	RecordtermID            int64  `json:"recordtermid"`
	Originalname            string `json:"originalname"`
	Uploadname              string `json:"uploadname"`
	Createdate              int64  `json:"createdate"`
	Showcreatedate          string `json:"showcreatedate"`
	Createdbyid             int64  `json:"createdbyid"`
	Createdgrpid            int64  `json:"createdgrpid"`
	Name                    string `json:"name"`
	Supportgrouplevelid     int64  `json:"supportgrouplevelid"`
	RecordID                int64  `json:"recordid"`
	ClientID                int64  `json:"clientid"`
	Mstorgnhirarchyid       int64  `json:"mstorgnhirarchyid"`
	Userid                  int64  `"json:userid"`
	Usergroupid             int64  `json:"usergroupid"`
	RecorduserID            int64  `json:"recorduserid"`
	RecordusergrpID         int64  `json:"recordusergrpid"`
	RecordoriginaluserID    int64  `json:"recordoriginaluserid"`
	RecordoriginalusergrpID int64  `json:"recordoriginalusergrpid"`
}

func (p *RecordAttachmentfiles) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type RecordattachmentAllResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Details []RecordAttachmentfiles `json:"details"`
}

type DocumentupdateAllResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Customervisiblecomment struct {
	ClientID          int64  `json:"clientid"`
	Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
	RecordID          int64  `json:"recordid"`
	Createddate       string `json:"createddate"`
	Daycount          int64  `json:"daycount"`
	Recordtrackvalue  string `json:"recordtrackvalue"`
	Recordcreatedate  string `json:"recordcreatedate"`
}

//Customervisiblecomment is defined for response of API
type CustomervisiblecommentAllResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Details []Customervisiblecomment `json:"details"`
}

type Recordtermseqvalue struct {
	Recordtermvalue string `json:"recordtermvalue"`
	Createddate     int64  `json:"createddate"`
	Showcreatedate  string `json:"showcreatedate"`
	Name            string `json:"name"`
}

type RecordtermseqvalueAllResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Details []Recordtermseqvalue `json:"details"`
}

type ParentchildSearchAllResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Details []NewActivitylogsEntity `json:"details"`
}

type LinkRecordEntity struct {
	ClientID          int64  `json:"clientid"`
	Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
	RecordID          int64  `json:"recordid"`
	LinkrecordID      int64  `json:"linkrecordid"`
	Userid            int64  `"json:userid"`
	Usergroupid       int64  `json:"usergroupid"`
	Linkrecordno      string `json:"linkrecordno"`
}

type LinkRecordDetailsEntity struct {
	Recordtitle  string `json:"recordtitle"`
	Recordcode   string `json:"recordcode"`
	Recordtype   string `json:"recordtype"`
	LinkrecordID int64  `json:"linkrecordid"`
}

func (p *LinkRecordEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

type LinkRecordDetailsAllResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Details []LinkRecordDetailsEntity `json:"details"`
}

type LinkRecordSaveAllResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

type ParentRecordInfoEntity struct {
	Recordtitle      string `json:"recordtitle"`
	Recordcode       string `json:"recordcode"`
	Recordstatus     string `json:"recordstatus"`
	ID               int64  `json:"id"`
	PlannedStartDate string `json:"plannedstartdate"`
	PlannedEndDate   string `json:"plannedenddate"`
}

type ParentRecordInfoAllResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Details ParentRecordInfoEntity `json:"details"`
}
//below entity is only for updating sladata of ticket which doesn't caculate sla
type RecordNoEntity struct {
    RecordNo []string `json:"recordno"`
}
type RecordInfoEntity struct {
    ClientID          int64  `json:"clientid"`
    Mstorgnhirarchyid int64  `json:"mstorgnhirarchyid"`
    RecordID          int64  `json:"recordid"`
    Datetime          string `json:"datetime"`
}
func (p *RecordNoEntity) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(p)
}
