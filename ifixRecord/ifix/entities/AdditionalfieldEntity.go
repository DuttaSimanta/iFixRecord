package entities

import (
	"encoding/json"
	"io"
)

type AdditionalfieldRequestEntity struct {
	Clientid              int64                       `json:"clientid"`
	Mstorgnhirarchyid     int64                       `json:"mstorgnhirarchyid"`
	RecordTypeDiffTypeID  int64                       `json:"recordtypedifftypeid"`
	RecordTypeDiffID      int64                       `json:"recordtypediffid"`
	RecordCatSet          []AdditionalFieldDiffEntity `json:"recordcatset"`
	Mstdifferentiationset []AdditionalFieldDiffEntity `json:"mstdifferentiationset"`
}

type AdditionalFieldDiffEntity struct {
	Mstdifferentiationtypeid int64 `json:"mstdifferentiationtypeid"`
	Mstdifferentiationid     int64 `json:"mstdifferentiationid"`
}

type AdditionalFieldEntity struct {
	FieldID       int64  `json:"fieldid"`
	TermsID       int64  `json:"termsid"`
	TermsName     string `json:"termsname"`
	TermsValue    string `json:"termsvalue"`
	TermsTypeID   int64  `json:"termstypeid"`
	TermsTypeName string `json:"termstypename"`
	TermSeqNo     int64  `json:"seqno"`
	IsMandatory   int64  `json:"ismandatory"`
}

type AdditionalFieldResponeData struct {
	Status   bool                    `json:"success"`
	Message  string                  `json:"message"`
	Response []AdditionalFieldEntity `json:"response"`
}

func (w *AdditionalfieldRequestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
