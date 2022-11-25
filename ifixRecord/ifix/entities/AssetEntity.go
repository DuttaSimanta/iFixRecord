package entities

import (
	"encoding/json"
	"io"
)

type AssetAttrNameValEntity struct {
	Id     int64  `json:"id"`
	AttrID int64  `json:"attrid"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

type AssetAttrNameValRequestEntity struct {
	Clientid          int64    `json:"clientid"`
	Mstorgnhirarchyid int64    `json:"mstorgnhirarchyid"`
	Recordid          int64    `json:"recordid"`
	AssetFieldsNames  []string `json:"assetfieldsnames"`
}

func (w *AssetAttrNameValRequestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
