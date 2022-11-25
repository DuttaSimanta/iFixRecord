package utility

import (
	"database/sql"
	"encoding/json"
	"ifixRecord/ifix/logger"
	"sync"
)

type NullString struct {
	sql.NullString
}

type NullInt64 struct {
	sql.NullInt64
}

type NullFloat64 struct {
	sql.NullFloat64
}

// NewDbFieldBind ...
func NewDbFieldBind() *DbFieldBind {
	return &DbFieldBind{}
}

// FieldBinding is deisgned for SQL rows.Scan() query.
type DbFieldBind struct {
	sync.RWMutex // embedded.  see http://golang.org/ref/spec#Struct_types
	FieldArr     []interface{}
	FieldPtrArr  []interface{}
	FArr         []string
	FieldCount   int64
	FArrTypes    []*sql.ColumnType
	MapFieldToID map[string]int64
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}
func (s NullInt64) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.Int64)
}
func (s NullFloat64) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.Float64)
}
func (fb *DbFieldBind) put(k string, v int64) {
	fb.Lock()
	defer fb.Unlock()
	fb.MapFieldToID[k] = v
}

// Get ...
func (fb *DbFieldBind) Get(k string) interface{} {
	fb.RLock()
	defer fb.RUnlock()
	// TODO: check map key exist and fb.FieldArr boundary.
	return fb.FieldPtrArr[fb.MapFieldToID[k]]
}

// PutFields ...
func (fb *DbFieldBind) PutFields(rows *sql.Rows) error {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	fb.FArrTypes = colTypes
	fCount := len(colTypes)
	fb.FieldPtrArr = make([]interface{}, fCount)
	fb.MapFieldToID = make(map[string]int64, fCount)

	for k, v := range colTypes {
		switch v.DatabaseTypeName() {
		case "VARCHAR":
			fb.FieldPtrArr[k] = new(NullString)
		case "TEXT":
			fb.FieldPtrArr[k] = new(NullString)
		case "TIMESTAMP":
			fb.FieldPtrArr[k] = new(NullString)
		case "INT":
			fb.FieldPtrArr[k] = new(NullInt64)
		case "FLOAT":
			fb.FieldPtrArr[k] = new(NullFloat64)
		default:
			fb.FieldPtrArr[k] = new(NullString)
			logger.Log.Println("Column Data Type:", v.DatabaseTypeName())
		}
		//fb.FieldPtrArr[k] = new(interface{})
		//	logger.Log.Println("Column Data Type:", v.DatabaseTypeName())
		fb.put(v.Name(), int64(k))
	}
	return nil
}

// GetFieldPtrArr ...
func (fb *DbFieldBind) GetFieldPtrArr() []interface{} {
	//logger.Log.Println("bf", fb.FieldPtrArr)

	return fb.FieldPtrArr
}

// GetFieldArr ...
func (fb *DbFieldBind) GetFieldArr() map[string]interface{} {
	m := make(map[string]interface{}, fb.FieldCount)

	for k, v := range fb.MapFieldToID {

		switch fb.FieldPtrArr[v].(type) {
		case NullString:
			data := fb.FieldPtrArr[v].(NullString)
			aa, _ := data.MarshalJSON()
			logger.Log.Println("------------", aa)
		case NullInt64:
			data := fb.FieldPtrArr[v].(NullInt64)
			aa, _ := data.MarshalJSON()
			logger.Log.Println("------------", aa)
		case NullFloat64:
			data := fb.FieldPtrArr[v].(NullFloat64)
			aa, _ := data.MarshalJSON()
			logger.Log.Println("------------", aa)
		}

		m[k] = fb.FieldPtrArr[v]
	}

	return m
}
