package entities

import (
	"encoding/json"
	"io"
)

type TrnAsset struct {
	Id                 int64
	ClientId           int64
	MstOrgnHirarchyId  int64
	MstDifftypeid      int64
	AssetId            string
	AdditionalAttr     string
	DeleteFlag         int64
	ActiveFlag         int64
	AuditTransactionId int64
}
type MapAssetDifferentiation struct {
	Id                       int64  `json:"id"`
	Clientid                 int64  `json:"clientid"`
	Mstorgnhirarchyid        int64  `json:"mstorgnhirarchyid"`
	Mstdifferentiationtypeid int64  `json:"mstdifferentiationtypeid"`
	Mstdifferentiationid     int64  `json:"mstdifferentiationid"`
	Trnassetid               int64  `json:"trnassetid"`
	Value                    string `json:"value"`
	Deleteflg                int64  `json:"deleteflg"`
	Activeflg                int64  `json:"activeflg"`
	AuditTransactionId       int64  `json:"audittransactionid"`
}
type InsertRecordAssetEntity struct {
	Clientid           int64                  `json:"clientid"`
	Mstorgnhirarchyid  int64                  `json:"mstorgnhirarchyid"`
	Recordid           int64                  `json:"recordid"`
	RecordStageid      int64                  `json:"recordstageid"`
	ParentRecordID     int64                  `json:"parentrecordid"`
	TicketTypeSequence int64                  `json:"tickettypeseq"`
	Assetid            int64                  `json:"assetid"`
	AssetHeaderId      int64                  `json:"assetheaderid"`
	Userid             int64                  `json:"userid"`
	GroupID            int64                  `json:"groupid"`
	MstDiffTypeID      int64                  `json:"mstdifftypeid"`
	AssetDetails       map[string]interface{} `json:"asset"`
}
type FetchAssetHistoryRequest struct {
	Clientid          int64 `json:"clientid"`
	Mstorgnhirarchyid int64 `json:"mstorgnhirarchyid"`
	Assetid           int64 `json:"assetid"`
}
type UpdateRecordAssetEntity struct {
	Clientid           int64  `json:"clientid"`
	Mstorgnhirarchyid  int64  `json:"mstorgnhirarchyid"`
	Recordid           int64  `json:"recordid"`
	RecordStageid      int64  `json:"recordstageid"`
	ParentRecordID     int64  `json:"parentrecordid"`
	TicketTypeSequence int64  `json:"tickettypeseq"`
	Assetid            int64  `json:"assetid"`
	AssetHeaderId      int64  `json:"assetheaderid"`
	UpdatedValue       string `json:"updatedvalue"`
	Userid             int64  `json:"userid"`
	GroupID            int64  `json:"groupid"`
}

type RecordAssetRequestEntity struct {
	Clientid          int64   `json:"clientid"`
	Mstorgnhirarchyid int64   `json:"mstorgnhirarchyid"`
	Recordid          int64   `json:"recordid"`
	RecordStageid     int64   `json:"recordstageid"`
	AssetID           []int64 `json:"assetid"`
	Userid            int64   `json:"userid"`
	GroupID           int64   `json:"groupid"`
}

type RecordAssetEntity struct {
	Clientid          int64 `json:"clientid"`
	Mstorgnhirarchyid int64 `json:"mstorgnhirarchyid"`
	Recordid          int64 `json:"recordid"`
	RecordStageid     int64 `json:"recordstageid"`
	AssetID           int64 `json:"assetid"`
	TrnAssetID        int64 `json:"trnassetid"`
}

func (w *RecordAssetRequestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}

func (w *UpdateRecordAssetEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
func (w *InsertRecordAssetEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
func (w *FetchAssetHistoryRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
