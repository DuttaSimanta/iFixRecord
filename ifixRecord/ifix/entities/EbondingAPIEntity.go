package entities

import (
	"encoding/json"
	"io"
)

type EbondingApiRespone struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details int64  `json:"details"`
}

func (w *EbondingRecordEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

//

// type FileAttachmentEntity struct {
// 	Clientid          int64  `json:"clientid"`
// 	Mstorgnhirarchyid string `json:"mstorgnhirarchyid"`
// 	UploadedFileName  string `json:"uploadedfilename"`
// 	OriginalFileName  string `json:"originalfilename"`
// }

// type FileuploadEntity struct {
// 	Id                 int64  `json:"id"`
// 	Clientid           int64  `json:"clientid"`
// 	Mstorgnhirarchyid  int64  `json:"mstorgnhirarchyid"`
// 	Credentialtype     string `json:"credentialtype"`
// 	Credentialaccount  string `json:"credentialaccount"`
// 	Credentialpassword string `json:"credentialpassword"`
// 	Credentialkey      string `json:"credentialkey"`
// 	Activeflg          int64  `json:"activeflg"`
// 	Originalfile       string `json:"originalfile"`
// 	Filename           string `json:"filename"`
// 	Path               string `json:"path"`
// }

// type FileuploadResponse struct {
// 	Success bool             `json:"success"`
// 	Message string           `json:"message"`
// 	Details FileuploadEntity `json:"details"`
// }

type EbondingRecordEntity struct {
	ClientID          int64  `json:"clientid"`
	MstorgnhirarchyID int64  `json:"mstorgnhirarchyid"`
	RecorddifftypeID  int64  `json:"recorddifftypeid"`
	RecorddiffID      int64  `json:"recorddiffid"`
	RecordID          int64  `json:"recordid"`
	RecordStagedID    int64  `json:"recordstageid"`
	RecordCode        string `json:"code"`
	EbondingID        int64  `json:"ebondingid"`
	EbondingModule    string `json:"ebondingmodule"`
	UploadedFileName  string `json:"uploadedfilename"`
	OriginalFileName  string `json:"originalfilename"`
	Worknote          string `json:"worknote"`
	EbondingSeq       int64  `json:"ebondingseq"`
	EbondingModuleSeq int64  `json:"ebondingmoduleseq"`
	StatusID          int64  `json:"statusid"`
	StatusSeq         int64  `json:"statusseq"`
	
}

type EbondingOtonCreateRecord struct {
	First_Name       string `json:"First_Name"`
	Last_Name        string `json:"Last_Name"`
	ShortDescription string `json:"Description"`
	LongDescription  string `json:"Detailed_Decription"`
	Impact           string `json:"Impact"`
	Urgency          string `json:"Urgency"`
	Status           string `json:"Status"`
	Reported_Source  string `json:"Reported Source"`
	Service_Type     string `json:"Service_Type"`
	Cat3             string `json:"Categorization Tier 1"`
	Cat4             string `json:"Categorization Tier 2"`
	Cat5             string `json:"Categorization Tier 3"`
	IfixTicketID     string `json:"External_System_ID"`
}

type EbondingOtonCreateRecordEntity struct {
	Values EbondingOtonCreateRecord `json:"values"`
}
type EbondingTransactionLog struct {
	Ebondingid   int64  `json:"ebondingid"`
	RecordID     int64  `json:"ifixrecordid"`
	Requestjson  string `json:"requestjson"`
	Responsejson string `json:"responsejson"`
	Responsecode int64  `json:"responsecode"`
}

type EbondingTicketID struct {
	IncidentNumber string `json:"Incident Number"`
}

type EbondingSelf struct {
	Self []EbondingSelfLink `json:"self"`
}

type EbondingSelfLink struct {
	Href string `json:"href"`
}
type EbondingOtonCreateTicketResponseEntities struct {
	Values EbondingTicketID `json:"values"`
	Links  EbondingSelf     `json:"_links"`
}
type OtonWorkLogEntity struct {
	IncidentNumber      string `json:"Incident Number"`
	WorkLogType         string `json:"Work Log Type"`
	Z1DAction           string `json:"z1D Action"`
	ViewAccess          string `json:"View Access"`
	SecureWorkLog       string `json:"Secure Work Log"`
	DetailedDescription string `json:"Detailed Description"`
}

type OtonWorkLogEntities struct {
	Values OtonWorkLogEntity `json:"values"`
}

type EbondingTDLCreateRecord struct {
	ShortDescription string                             `json:"short_description"`
	LongDescription  string                             `json:"description"`
	Impact           string                             `json:"impact"`
	Urgency          string                             `json:"urgency"`
	Email            string                             `json:"caller_email_id"`
	LoginID          string                             `json:"caller_id"`
	RequestorMobile  string                             `json:"contact_number"`
	TicketID         string                             `json:"aai_id"`
	Notes            string                             `json:"notes"`
	Fileattachment   []EbondingTDLFileAttachmentDetails `json:"attachment"`
}

type EbondingTDLCreateRecordResponse struct {
	Message     string `json:"message"`
	TDLTicketID string `json:"incident_number"`
}
type EbondingTDLCreateRecordResponseEntity struct {
	Result EbondingTDLCreateRecordResponse `json:"result"`
}

type EbondingTDLUpdateRecordResponse struct {
	Message     string `json:"message"`
	TDLTicketID string `json:"sid"`
}
type EbondingTDLUpdateRecordResponseEntity struct {
	Result EbondingTDLCreateRecordResponse `json:"result"`
}

type EbondingTDLFileAttachmentDetails struct {
	Filename    string `json:"name"`
	Filetype    string `json:"type"`
	Filecontent string `json:"content"`
}

func (w *EbondingTDLCreateRecord) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

type EbondingTDLUpdateRecord struct {
	ServicenowID   string                             `json:"servicenow_id"`
	TicketID       string                             `json:"aai_id"`
	StateID        string                             `json:"state"`
	Impact         string                             `json:"impact"`
	Urgency        string                             `json:"urgency"`
	Notes          string                             `json:"notes"`
	Fileattachment []EbondingTDLFileAttachmentDetails `json:"attachment"`
}

func (w *EbondingTDLUpdateRecord) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

type FileAttachmentEntity struct {
	Clientid          int64  `json:"clientid"`
	Mstorgnhirarchyid string `json:"mstorgnhirarchyid"`
	UploadedFileName  string `json:"uploadedfilename"`
	OriginalFileName  string `json:"originalfilename"`
}

/*type FileuploadEntity struct {
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

type FileuploadResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Details FileuploadEntity `json:"details"`
}*/



type OtonAttachmentEntity struct {
	IncidentNumber      string `json:"Incident Number"`
	WorkLogType         string `json:"Work Log Type"`
	Z1DAction           string `json:"z1D Action"`
	ViewAccess          string `json:"View Access"`
	SecureWorkLog       string `json:"Secure Work Log"`
	DetailedDescription string `json:"Detailed Description"`
	Z2AFWorkLog01       string `json:"z2AF Work Log01"`
}

type OtonAttachmentEntitis struct {
	Values OtonAttachmentEntity `json:"values"`
}
