package entities

type DynamicqueryEntity struct {
	ID         int64  `json:"id"`
	Query      string `json:"query"`
	QueryParam string `json:"queryparams"`
	JoinQuery  string `json:"joinquery"`
	QueryType  int64  `json:"querytype"`
}

type QueryResultEntity struct {
	ID                  int64                 `json:"id"`
	Recordtitle         string                `json:"recordtitle"`
	Code                string                `json:"code"`
	Supportgroupname    string                `json:"supportgroupname"`
	Levelname           string                `json:"levelname"`
	Priority            string                `json:"priority"`
	Createdby           string                `json:"createdby"`
	Createdatetime      string                `json:"createdatetime"`
	Status              string                `json:"status"`
	Impact              string                `json:"impact"`
	Requesterinfo       string                `json:"requesterinfo"`
	Recorddescription   string                `json:"recorddescription"`
	Urgency             string                `json:"urgency"`
	Recordtype          string                `json:"recordtype"`
	Assignee            string                `json:"assignee"`
	AssignedGroupLevel  string                `json:"assignedgrouplevel"`
	AssignedGroup       string                `json:"assignedgroup"`
	StageID             int64                 `json:"stageid"`
	PriorityChangeCount int64                 `json:"prioritychangecount"`
	FollowUpCount       int64                 `json:"followupcount"`
	OutBoundCount       int64                 `json:"outboundcount"`
	ReOpenCount         int64                 `json:"reopencount"`
	ResolutionViolated  string                `json:"resolutionviolated"`
	Aging               int                   `json:"aging"`
	Categories          map[int64]interface{} `json:"category"`
}

type CategoryNameOnly struct {
	TypeName string `json:"typename"`
	Name     string `json:"name"`
}

type QueryResultResponeData struct {
	Status   bool               `json:"success"`
	Message  string             `json:"message"`
	Response QueryDetailsEntity `json:"details"`
}

type QueryDetailsEntity struct {
	Total  int64               `json:"total"`
	Values []QueryResultEntity `json:"values"`
}

type QueryResponeData struct {
	Status   bool          `json:"success"`
	Message  string        `json:"message"`
	Response []interface{} `json:"details"`
}

type QueryCountResultResponeData struct {
	Status   bool                   `json:"success"`
	Message  string                 `json:"message"`
	Response map[string]interface{} `json:"details"`
}
