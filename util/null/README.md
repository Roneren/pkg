## null-extended

`import "github.com/corestoreio/csfw/util/null"`

null-extended is a library with reasonable options for dealing with nullable SQL
and JSON values

There are two packages: `null` and its subpackage `zero`.

Types in `null` will only be considered null on null input, and will JSON encode
to `null`. If you need zero and null be considered separate values, use these.

Types in `zero` are treated like zero values in Go: blank string input will
produce a null `zero.String`, and null Strings will JSON encode to `""`. Zero
values of these types will be considered null to SQL. If you need zero and null
treated the same, use these.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this
library in place of `sql.NullXXX`. All types also implement:
`encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`,
`json.Unmarshaler` and `sql.Scanner`.

### null package

`import "github.com/corestoreio/csfw/util/null"`

#### null.String Nullable string.

Marshals to JSON null if SQL source data is null. Zero (blank) input will not
produce a null String. Can unmarshal from `sql.NullString` JSON input or string
input.

#### null.Bool Nullable bool.

Marshals to JSON null if SQL source data is null. False input will not produce a
null Bool. Can unmarshal from `sql.NullBool` JSON input.

#### null.Time

Marshals to JSON null if SQL source data is null. Uses `time.Time`'s marshaler.
Can unmarshal from `pq.NullTime` and similar JSON input.

#### null.Float32 Nullable float32.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Float32. Can unmarshal from `null.NullFloat32` JSON input.

#### null.Float64 Nullable float64.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Float64. Can unmarshal from `sql.NullFloat64` JSON input.

#### null.Int Nullable int.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Int. Can unmarshal from `null.NullInt` JSON input.

#### null.Int8 Nullable int8.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Int8. Can unmarshal from `null.NullInt8` JSON input.

#### null.Int16 Nullable int16.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Int16. Can unmarshal from `null.NullInt16` JSON input.

#### null.Int32 Nullable int32.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Int32. Can unmarshal from `null.NullInt32` JSON input.

#### null.Int64 Nullable int64.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Int64. Can unmarshal from `sql.NullInt64` JSON input.

#### null.Uint Nullable uint.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Uint. Can unmarshal from `null.NullUint` JSON input.

#### null.Uint8 Nullable uint8.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Uint8. Can unmarshal from `null.NullUint8` JSON input.

#### null.Uint16 Nullable uint16.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Uint16. Can unmarshal from `null.NullUint16` JSON input.

#### null.Uint32 Nullable int32.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Uint32. Can unmarshal from `null.NullUint32` JSON input.

#### null.Int64 Nullable uint64.

Marshals to JSON null if SQL source data is null. Zero input will not produce a
null Uint64. Can unmarshal from `null.NullUint64` JSON input.

### zero package

`import "github.com/corestoreio/csfw/util/null/zero"`

#### zero.String Nullable string.

Will marshal to a blank string if null. Blank string input produces a null
String. Null values and zero values are considered equivalent. Can unmarshal
from `sql.NullString` JSON input.

#### zero.Bool Nullable bool.

Will marshal to false if null. `false` produces a null Float. Null values and
zero values are considered equivalent. Can unmarshal from `sql.NullBool` JSON
input.

#### zero.Time

Will marshal to the zero time if null. Uses `time.Time`'s marshaler. Can
unmarshal from `pq.NullTime` and similar JSON input.

#### zero.Float32 Nullable float32.

Will marshal to 0 if null. 0.0 produces a null Float32. Null values and zero
values are considered equivalent. Can unmarshal from `zero.NullFloat32` JSON
input.

#### zero.Float64 Nullable float64.

Will marshal to 0 if null. 0.0 produces a null Float64. Null values and zero
values are considered equivalent. Can unmarshal from `sql.NullFloat64` JSON
input.

#### zero.Int Nullable int.

Will marshal to 0 if null. 0 produces a null Int. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullInt` JSON input.

#### zero.Int8 Nullable int8.

Will marshal to 0 if null. 0 produces a null Int8. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullInt8` JSON input.

#### zero.Int16 Nullable int16.

Will marshal to 0 if null. 0 produces a null Int16. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullInt16` JSON input.

#### zero.Int32 Nullable int32.

Will marshal to 0 if null. 0 produces a null Int32. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullInt32` JSON input.

#### zero.Int64 Nullable int64.

Will marshal to 0 if null. 0 produces a null Int64. Null values and zero values
are considered equivalent. Can unmarshal from `sql.NullInt64` JSON input.

#### zero.Uint Nullable uint.

Will marshal to 0 if null. 0 produces a null Uint. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullUint` JSON input.

#### zero.Uint8 Nullable uint8.

Will marshal to 0 if null. 0 produces a null Uint8. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullUint8` JSON input.

#### zero.Uint16 Nullable uint16.

Will marshal to 0 if null. 0 produces a null Uint16. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullUint16` JSON input.

#### zero.Uint32 Nullable uint32.

Will marshal to 0 if null. 0 produces a null Uint32. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullUint32` JSON input.

#### zero.Uint64 Nullable uint64.

Will marshal to 0 if null. 0 produces a null Uint64. Null values and zero values
are considered equivalent. Can unmarshal from `zero.NullUint64` JSON input.

### Bugs `json`'s `",omitempty"` struct tag does not work correctly right now.
### It will never omit a null or empty String. This might be [fixed
### eventually](https://github.com/golang/go/issues/4357).

### License
BSD
