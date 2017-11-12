// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata

import (
	"encoding/json"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

// DmlgenTypes represents a single row for DB table `dmlgen_types`.
// Auto generated.
// Just another comment.
//easyjson:json
type DmlgenTypes struct {
	ID             int64           `json:"id,omitempty" `                // id int(11) NOT NULL PRI  auto_increment ""
	ColBigint1     dml.NullInt64   `json:"col_bigint_1,omitempty" `      // col_bigint_1 bigint(20) NULL    ""
	ColBigint2     int64           `json:"col_bigint_2,omitempty" `      // col_bigint_2 bigint(20) NOT NULL  DEFAULT '0'  ""
	ColBigint3     dml.NullInt64   `json:"col_bigint_3,omitempty" `      // col_bigint_3 bigint(20) unsigned NULL    ""
	ColBigint4     uint64          `json:"col_bigint_4,omitempty" `      // col_bigint_4 bigint(20) unsigned NOT NULL  DEFAULT '0'  ""
	ColBlob        dml.NullString  `json:"col_blob,omitempty" `          // col_blob blob NULL    ""
	ColDate1       dml.NullTime    `json:"col_date_1,omitempty" `        // col_date_1 date NULL    ""
	ColDate2       time.Time       `json:"col_date_2,omitempty" `        // col_date_2 date NOT NULL  DEFAULT ''0000-00-00''  ""
	ColDatetime1   dml.NullTime    `json:"col_datetime_1,omitempty" `    // col_datetime_1 datetime NULL    ""
	ColDatetime2   time.Time       `json:"col_datetime_2,omitempty" `    // col_datetime_2 datetime NOT NULL  DEFAULT ''0000-00-00 00:00:00''  ""
	ColDecimal100  dml.NullFloat64 `json:"col_decimal_10_0,omitempty" `  // col_decimal_10_0 decimal(10,0) unsigned NULL    ""
	ColDecimal124  dml.NullFloat64 `json:"col_decimal_12_4,omitempty" `  // col_decimal_12_4 decimal(12,4) NULL    ""
	Price124a      dml.Decimal     `json:"price_12_4a,omitempty" `       // price_12_4a decimal(12,4) NULL    ""
	Price124b      dml.Decimal     `json:"price_12_4b,omitempty" `       // price_12_4b decimal(12,4) NOT NULL  DEFAULT '0.0000'  ""
	ColDecimal123  float64         `json:"col_decimal_12_3,omitempty" `  // col_decimal_12_3 decimal(12,3) NOT NULL  DEFAULT '0.000'  ""
	ColDecimal206  float64         `json:"col_decimal_20_6,omitempty" `  // col_decimal_20_6 decimal(20,6) NOT NULL  DEFAULT '0.000000'  ""
	ColDecimal2412 float64         `json:"col_decimal_24_12,omitempty" ` // col_decimal_24_12 decimal(24,12) NOT NULL  DEFAULT '0.000000000000'  ""
	ColFloat       float64         `json:"col_float,omitempty" `         // col_float float NOT NULL  DEFAULT '1'  ""
	ColInt1        dml.NullInt64   `json:"col_int_1,omitempty" `         // col_int_1 int(10) NULL    ""
	ColInt2        int64           `json:"col_int_2,omitempty" `         // col_int_2 int(10) NOT NULL  DEFAULT '0'  ""
	ColInt3        dml.NullInt64   `json:"col_int_3,omitempty" `         // col_int_3 int(10) unsigned NULL    ""
	ColInt4        uint64          `json:"col_int_4,omitempty" `         // col_int_4 int(10) unsigned NOT NULL  DEFAULT '0'  ""
	ColLongtext1   dml.NullString  `json:"col_longtext_1,omitempty" `    // col_longtext_1 longtext NULL    ""
	ColLongtext2   string          `json:"col_longtext_2,omitempty" `    // col_longtext_2 longtext NOT NULL  DEFAULT ''''  ""
	ColMediumblob  dml.NullString  `json:"col_mediumblob,omitempty" `    // col_mediumblob mediumblob NULL    ""
	ColMediumtext1 dml.NullString  `json:"col_mediumtext_1,omitempty" `  // col_mediumtext_1 mediumtext NULL    ""
	ColMediumtext2 string          `json:"col_mediumtext_2,omitempty" `  // col_mediumtext_2 mediumtext NOT NULL  DEFAULT ''''  ""
	ColSmallint1   dml.NullInt64   `json:"col_smallint_1,omitempty" `    // col_smallint_1 smallint(5) NULL    ""
	ColSmallint2   int64           `json:"col_smallint_2,omitempty" `    // col_smallint_2 smallint(5) NOT NULL  DEFAULT '0'  ""
	ColSmallint3   dml.NullInt64   `json:"col_smallint_3,omitempty" `    // col_smallint_3 smallint(5) unsigned NULL    ""
	ColSmallint4   uint64          `json:"col_smallint_4,omitempty" `    // col_smallint_4 smallint(5) unsigned NOT NULL  DEFAULT '0'  ""
	HasSmallint5   bool            `json:"has_smallint_5,omitempty" `    // has_smallint_5 smallint(5) unsigned NOT NULL  DEFAULT '0'  ""
	IsSmallint5    dml.NullBool    `json:"is_smallint_5,omitempty" `     // is_smallint_5 smallint(5) NULL    ""
	ColText        dml.NullString  `json:"col_text,omitempty" `          // col_text text NULL    ""
	ColTimestamp1  time.Time       `json:"col_timestamp_1,omitempty" `   // col_timestamp_1 timestamp NOT NULL  DEFAULT 'current_timestamp()'  ""
	ColTimestamp2  dml.NullTime    `json:"col_timestamp_2,omitempty" `   // col_timestamp_2 timestamp NULL    ""
	ColTinyint1    int64           `json:"col_tinyint_1,omitempty" `     // col_tinyint_1 tinyint(1) NOT NULL  DEFAULT '0'  ""
	ColVarchar1    string          `json:"col_varchar_1,omitempty" `     // col_varchar_1 varchar(1) NOT NULL  DEFAULT ''0''  ""
	ColVarchar100  dml.NullString  `json:"col_varchar_100,omitempty" `   // col_varchar_100 varchar(100) NULL    ""
	ColVarchar16   string          `json:"col_varchar_16,omitempty" `    // col_varchar_16 varchar(16) NOT NULL  DEFAULT ''de_DE''  ""
	ColChar1       dml.NullString  `json:"col_char_1,omitempty" `        // col_char_1 char(21) NULL    ""
	ColChar2       string          `json:"col_char_2,omitempty" `        // col_char_2 char(17) NOT NULL  DEFAULT ''xchar''  ""
}

// NewDmlgenTypes creates a new pointer with pre-initialized fields. Auto
// generated.
func NewDmlgenTypes() *DmlgenTypes {
	return &DmlgenTypes{}
}

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner. Auto generated.
func (e *DmlgenTypes) AssignLastInsertID(id int64) {
	e.ID = int64(id)
}

// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *DmlgenTypes) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int64(&e.ID).NullInt64(&e.ColBigint1).Int64(&e.ColBigint2).NullInt64(&e.ColBigint3).Uint64(&e.ColBigint4).NullString(&e.ColBlob).NullTime(&e.ColDate1).Time(&e.ColDate2).NullTime(&e.ColDatetime1).Time(&e.ColDatetime2).NullFloat64(&e.ColDecimal100).NullFloat64(&e.ColDecimal124).Decimal(&e.Price124a).Decimal(&e.Price124b).Float64(&e.ColDecimal123).Float64(&e.ColDecimal206).Float64(&e.ColDecimal2412).Float64(&e.ColFloat).NullInt64(&e.ColInt1).Int64(&e.ColInt2).NullInt64(&e.ColInt3).Uint64(&e.ColInt4).NullString(&e.ColLongtext1).String(&e.ColLongtext2).NullString(&e.ColMediumblob).NullString(&e.ColMediumtext1).String(&e.ColMediumtext2).NullInt64(&e.ColSmallint1).Int64(&e.ColSmallint2).NullInt64(&e.ColSmallint3).Uint64(&e.ColSmallint4).Bool(&e.HasSmallint5).NullBool(&e.IsSmallint5).NullString(&e.ColText).Time(&e.ColTimestamp1).NullTime(&e.ColTimestamp2).Int64(&e.ColTinyint1).String(&e.ColVarchar1).NullString(&e.ColVarchar100).String(&e.ColVarchar16).NullString(&e.ColChar1).String(&e.ColChar2).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "id":
			cm.Int64(&e.ID)
		case "col_bigint_1":
			cm.NullInt64(&e.ColBigint1)
		case "col_bigint_2":
			cm.Int64(&e.ColBigint2)
		case "col_bigint_3":
			cm.NullInt64(&e.ColBigint3)
		case "col_bigint_4":
			cm.Uint64(&e.ColBigint4)
		case "col_blob":
			cm.NullString(&e.ColBlob)
		case "col_date_1":
			cm.NullTime(&e.ColDate1)
		case "col_date_2":
			cm.Time(&e.ColDate2)
		case "col_datetime_1":
			cm.NullTime(&e.ColDatetime1)
		case "col_datetime_2":
			cm.Time(&e.ColDatetime2)
		case "col_decimal_10_0":
			cm.NullFloat64(&e.ColDecimal100)
		case "col_decimal_12_4":
			cm.NullFloat64(&e.ColDecimal124)
		case "price_12_4a":
			cm.Decimal(&e.Price124a)
		case "price_12_4b":
			cm.Decimal(&e.Price124b)
		case "col_decimal_12_3":
			cm.Float64(&e.ColDecimal123)
		case "col_decimal_20_6":
			cm.Float64(&e.ColDecimal206)
		case "col_decimal_24_12":
			cm.Float64(&e.ColDecimal2412)
		case "col_float":
			cm.Float64(&e.ColFloat)
		case "col_int_1":
			cm.NullInt64(&e.ColInt1)
		case "col_int_2":
			cm.Int64(&e.ColInt2)
		case "col_int_3":
			cm.NullInt64(&e.ColInt3)
		case "col_int_4":
			cm.Uint64(&e.ColInt4)
		case "col_longtext_1":
			cm.NullString(&e.ColLongtext1)
		case "col_longtext_2":
			cm.String(&e.ColLongtext2)
		case "col_mediumblob":
			cm.NullString(&e.ColMediumblob)
		case "col_mediumtext_1":
			cm.NullString(&e.ColMediumtext1)
		case "col_mediumtext_2":
			cm.String(&e.ColMediumtext2)
		case "col_smallint_1":
			cm.NullInt64(&e.ColSmallint1)
		case "col_smallint_2":
			cm.Int64(&e.ColSmallint2)
		case "col_smallint_3":
			cm.NullInt64(&e.ColSmallint3)
		case "col_smallint_4":
			cm.Uint64(&e.ColSmallint4)
		case "has_smallint_5":
			cm.Bool(&e.HasSmallint5)
		case "is_smallint_5":
			cm.NullBool(&e.IsSmallint5)
		case "col_text":
			cm.NullString(&e.ColText)
		case "col_timestamp_1":
			cm.Time(&e.ColTimestamp1)
		case "col_timestamp_2":
			cm.NullTime(&e.ColTimestamp2)
		case "col_tinyint_1":
			cm.Int64(&e.ColTinyint1)
		case "col_varchar_1":
			cm.String(&e.ColVarchar1)
		case "col_varchar_100":
			cm.NullString(&e.ColVarchar100)
		case "col_varchar_16":
			cm.String(&e.ColVarchar16)
		case "col_char_1":
			cm.NullString(&e.ColChar1)
		case "col_char_2":
			cm.String(&e.ColChar2)
		default:
			return errors.NewNotFoundf("[testdata] DmlgenTypes Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

// DmlgenTypesCollection represents a collection type for DB table dmlgen_types
// Not thread safe. Auto generated.
// Just another comment.
//easyjson:json
type DmlgenTypesCollection struct {
	Data             []*DmlgenTypes
	BeforeMapColumns func(uint64, *DmlgenTypes) error
	AfterMapColumns  func(uint64, *DmlgenTypes) error
}

// MakeDmlgenTypesCollection creates a new initialized collection. Auto generated.
func MakeDmlgenTypesCollection() DmlgenTypesCollection {
	return DmlgenTypesCollection{
		Data: make([]*DmlgenTypes, 0, 5),
	}
}

func (cc DmlgenTypesCollection) scanColumns(cm *dml.ColumnMap, e *DmlgenTypes, idx uint64) error {
	if err := cc.BeforeMapColumns(idx, e); err != nil {
		return errors.WithStack(err)
	}
	if err := e.MapColumns(cm); err != nil {
		return errors.WithStack(err)
	}
	if err := cc.AfterMapColumns(idx, e); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// MapColumns implements dml.ColumnMapper interface. Auto generated.
func (cc DmlgenTypesCollection) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for i, e := range cc.Data {
			if err := cc.scanColumns(cm, e, uint64(i)); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		e := NewDmlgenTypes()
		if err := cc.scanColumns(cm, e, cm.Count); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			case "id":
				cm.Args = cm.Args.Int64s(cc.IDs()...)
			case "col_blob":
				cm.Args = cm.Args.Strings(cc.ColBlobs()...)
			case "col_date_2":
				cm.Args = cm.Args.Times(cc.ColDate2s()...)
			case "col_int_1":
				cm.Args = cm.Args.Int64s(cc.ColInt1s()...)
			case "col_int_2":
				cm.Args = cm.Args.Int64s(cc.ColInt2s()...)
			case "col_longtext_2":
				cm.Args = cm.Args.Strings(cc.ColLongtext2s()...)
			case "has_smallint_5":
				cm.Args = cm.Args.Bools(cc.HasSmallint5s()...)
			default:
				return errors.NewNotFoundf("[testdata] DmlgenTypesCollection Column %q not found", c)
			}
		}
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

// IDs returns a slice or appends to a slice all values.
// Auto generated.
func (cc DmlgenTypesCollection) IDs(ret ...int64) []int64 {
	if ret == nil {
		ret = make([]int64, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.ID)
	}
	return ret
}

// ColBlobs belongs to the column `col_blob`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc DmlgenTypesCollection) ColBlobs(ret ...string) []string {
	if ret == nil {
		ret = make([]string, 0, len(cc.Data))
	}

	dupCheck := make(map[string]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.ColBlob.String]; !ok {
			ret = append(ret, e.ColBlob.String)
			dupCheck[e.ColBlob.String] = struct{}{}
		}
	}
	return ret
}

// ColDate2s belongs to the column `col_date_2`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc DmlgenTypesCollection) ColDate2s(ret ...time.Time) []time.Time {
	if ret == nil {
		ret = make([]time.Time, 0, len(cc.Data))
	}

	dupCheck := make(map[time.Time]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.ColDate2]; !ok {
			ret = append(ret, e.ColDate2)
			dupCheck[e.ColDate2] = struct{}{}
		}
	}
	return ret
}

// ColInt1s belongs to the column `col_int_1`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc DmlgenTypesCollection) ColInt1s(ret ...int64) []int64 {
	if ret == nil {
		ret = make([]int64, 0, len(cc.Data))
	}

	dupCheck := make(map[int64]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.ColInt1.Int64]; !ok {
			ret = append(ret, e.ColInt1.Int64)
			dupCheck[e.ColInt1.Int64] = struct{}{}
		}
	}
	return ret
}

// ColInt2s belongs to the column `col_int_2`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc DmlgenTypesCollection) ColInt2s(ret ...int64) []int64 {
	if ret == nil {
		ret = make([]int64, 0, len(cc.Data))
	}

	dupCheck := make(map[int64]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.ColInt2]; !ok {
			ret = append(ret, e.ColInt2)
			dupCheck[e.ColInt2] = struct{}{}
		}
	}
	return ret
}

// ColLongtext2s belongs to the column `col_longtext_2`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc DmlgenTypesCollection) ColLongtext2s(ret ...string) []string {
	if ret == nil {
		ret = make([]string, 0, len(cc.Data))
	}

	dupCheck := make(map[string]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.ColLongtext2]; !ok {
			ret = append(ret, e.ColLongtext2)
			dupCheck[e.ColLongtext2] = struct{}{}
		}
	}
	return ret
}

// HasSmallint5s belongs to the column `has_smallint_5`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc DmlgenTypesCollection) HasSmallint5s(ret ...bool) []bool {
	if ret == nil {
		ret = make([]bool, 0, len(cc.Data))
	}

	dupCheck := make(map[bool]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.HasSmallint5]; !ok {
			ret = append(ret, e.HasSmallint5)
			dupCheck[e.HasSmallint5] = struct{}{}
		}
	}
	return ret
}

// UnmarshalJSON implements interface json.Unmarshaler.
func (cc *DmlgenTypesCollection) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, cc.Data)
}

// MarshalJSON implements interface json.Marshaler.
func (cc *DmlgenTypesCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(cc.Data)
}

// TODO add MarshalText and UnmarshalText.
// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (cc *DmlgenTypesCollection) UnmarshalBinary(data []byte) error {
	return cc.Unmarshal(data) // Implemented via github.com/gogo/protobuf
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (cc *DmlgenTypesCollection) MarshalBinary() (data []byte, err error) {
	return cc.Marshal() // Implemented via github.com/gogo/protobuf
}

// GobDecode kept for Go 1 compatibility reasons.
// deprecated in Go 2, use UnmarshalBinary
func (cc *DmlgenTypesCollection) GobDecode(data []byte) error {
	return cc.Unmarshal(data) // Implemented via github.com/gogo/protobuf
}

// GobEncode kept for Go 1 compatibility reasons.
// deprecated in Go 2, use MarshalBinary
func (cc *DmlgenTypesCollection) GobEncode() ([]byte, error) {
	return cc.Marshal() // Implemented via github.com/gogo/protobuf
}
