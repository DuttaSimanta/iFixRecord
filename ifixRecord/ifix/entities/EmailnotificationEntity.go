package entities

type EmailnotificationTemplate struct {
	Clientid            int64 `json:"clientid"`
	Mstorgnhirarchyid   int64 `json:"mstorgnhirarchyid"`
	Recordid            int64 `json:"recordid"`
	Eventnotificationid int64 `json:"eventnotificationid"`
	Channeltype         int64 `json:"channeltype"`
	Priorityid          int64 `json:"priorityid"`
	Followupcount       int64 `json:"followupcount"`
	Statusid            int64 `json:"statusid"`
	TermsSeq            int64 `json:"termseq"`
}
