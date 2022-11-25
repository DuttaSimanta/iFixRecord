package entities


type TransactionEntity struct{
	Id                  int64    `json:"id"`
	Clientid            int64    `json:"clientid"`
	Mstorgnhirarchyid   int64    `json:"mstorgnhirarchyid"`
	Recordtitle         string   `json:"recordtitle"`
	Recorddescription   string   `json:"recorddescription"`
	Requesterinfo       string   `json:"requesterinfo"`
	Recordstageid       int64    `json:"recordstageid"`
	Requestid           int64    `json:"requestid"`
	Daterequested       int64    `json:"daterequested"`
	Userid              int64    `json:"userid"`
	Groupid             int64    `json:"groupid"`
	Currentstateid      int64    `json:"currentstateid"`
	Deleteflg           int64    `json:"deleteflg"`
	Activeflg  			int64    `json:"activeflg"`
	Audittransactionid  int64    `json:"audittransactionid"`
}
type TransactionRespEntity struct{
	Id                  int64    `json:"id"`
	Clientid            int64    `json:"clientid"`
	Mstorgnhirarchyid   int64    `json:"mstorgnhirarchyid"`
	Recorddifftypeid   	int64    `json:"recorddifftypeid"`
	Recorddiffid   		int64    `json:"recorddiffid"`
	Seqno   			int64    `json:"seqno"`
	Statename           string   `json:"statename"`
	Supportgroupname    string   `json:"supportgroupname"`
	Lastgroupname       string   `json:"lastgroupname"`
	Status              string   `json:"status"`
	Username            string    `json:"username"`
	Lastusername        string    `json:"lastusername"`
	Transitionid        int64    `json:"transitionid"`
	Userid              int64    `json:"userid"`
	Groupid             int64    `json:"groupid"`
	Currentstateid      int64    `json:"currentstateid"`
	Grplevel      int64    `json:"grplevel"`
}
type TransactionEntityResponse struct {
	Success  	bool `json:"success"`
	Message 	string `json:"message"`
	Details 	[]TransactionRespEntity `json:"details"`
}
