// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// @todo hide password and other sensitive fields in JSON struct tags

const tplCopyPkg = `// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package {{ .Package }}
`

const tplHeader = `
// Auto generated via tableToStruct

import (
	"sort"
	"time"
    {{ if .HasTypeCodeValueTables }}
	"github.com/corestoreio/csfw/eav"{{end}}
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/storage/money"
)

var _ = (*time.Time)(nil)
var _ = (*money.Currency)(nil)

// TableIndex... is the index to a table. These constants are guaranteed
// to stay the same for all Magento versions. Please access a table via this
// constant instead of the raw table name. TableIndex iotas must start with 0.
const (
    {{ range $k,$v := .Tables }}TableIndex{{$v.Name}} {{ if eq $k 0 }}csdb.Index = iota{{ end }} // Table: {{$v.NameRaw}}
{{ end }}	TableIndexZZZ  // the maximum index, which is not available.
)

func init(){
    TableCollection = csdb.NewTableManager(
    {{ range $k,$v := .Tables }} csdb.AddTableByName(TableIndex{{.Name}}, "{{.NameRaw}}"),
    {{ end }} )
    // Don't forget to call TableCollection.ReInit(...) in your code to load the column definitions.
}`

const tplTable = `
// {{.Struct}} and {{.Slice}}, a type for DB table {{ .NameRaw }}
type (
    {{.Slice}} []*{{.Struct}}
    {{.Struct}} struct {
        {{ range .GoColumns }}{{.GoName}} {{.GoType}} {{ $.Tick }}db:"{{.Field.String}}" json:",omitempty"{{ $.Tick }} {{.Comment}}
        {{ end }} }
)

var _ sort.Interface = (*{{.Slice}})(nil)

// {{ typePrefix "SQLSelect" }} fills this slice with data from the database
func (s *{{.Slice}}) {{ typePrefix "SQLSelect" }}(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return csdb.LoadSlice(dbrSess, TableCollection, TableIndex{{.NameRaw | prepareVar}}, &(*s), cbs...)
}

// {{ typePrefix "SQLInsert" }} inserts all records into the database @todo
func (s *{{.Slice}}) {{ typePrefix "SQLInsert" }}(dbrSess dbr.SessionRunner, cbs ...csdb.DbrInsertCb) (int, error) {
	return 0, nil
}

// {{ typePrefix "SQLUpdate" }} updates all record in the database @todo
func (s *{{.Slice}}) {{ typePrefix "SQLUpdate" }}(dbrSess dbr.SessionRunner, cbs ...csdb.DbrUpdateCb) (int, error) {
	return 0, nil
}

// {{ typePrefix "SQLDelete" }} deletes all record from the database @todo
func (s *{{.Slice}}) {{ typePrefix "SQLDelete" }}(dbrSess dbr.SessionRunner, cbs ...csdb.DbrDeleteCb) (int, error) {
	return 0, nil
}

{{if (.FindByPk) ne ""}}
// {{ typePrefix .FindByPk }} searches the primary keys and returns a *{{.Struct}} if found or an error
func (s {{.Slice}}) {{ typePrefix .FindByPk }}(
{{range $k,$v := .Columns.PrimaryKeys}} {{ $v.Name }} {{$v.GetGoPrimitive false}},
{{end}}	) (*{{.Struct}}, error) {
	for _, u := range s {
		if u != nil {{ range $c := .Columns.PrimaryKeys }} && u.{{ $c.Name | camelize }}{{dbrType $c}} == {{$c.Name}} {{ end }} {
			return u, nil
		}
	}
	return nil, csdb.NewError("ID not found in {{.Slice}}")
}
{{end}}

{{ range $k,$c := .Columns.UniqueKeys }}
// {{ findBy $c.Name | typePrefix }} searches through this unique key and returns
// a *{{$.Struct}} if found or an error
func (s {{$.Slice}}) {{ findBy $c.Name | typePrefix }} ( {{ $c.Name }} {{$c.GetGoPrimitive false}} ) (*{{$.Struct}}, error) {
	for _, u := range s {
		if u != nil && u.{{ $c.Name | camelize }}{{ dbrType $c }} == {{$c.Name}} {
			return u, nil
		}
	}
	return nil, csdb.NewError("ID not found in {{$.Slice}}")
}
{{ end }}

// {{ typePrefix "Len" }} returns the length and  will satisfy the sort.Interface
func (s {{.Slice}}) {{ typePrefix "Len" }}() int { return len(s) }

// {{ typePrefix "Less" }} will satisfy the sort.Interface and compares via
// the primary key
func (s {{.Slice}}) {{ typePrefix "Less" }}(i, j int) bool {
	return {{ range $c := .Columns.PrimaryKeys }} s[i].{{ $c.Name | camelize }}{{dbrType $c}} < s[j].{{ $c.Name | camelize }}{{dbrType $c}} && {{ end }} 1 == 1
}

// {{ typePrefix "Swap" }} will satisfy the sort.Interface
func (s {{.Slice}}) {{ typePrefix "Swap" }}(i, j int) { s[i], s[j] = s[j], s[i] }

// {{ typePrefix "Sort" }} will sort {{.Slice}}
func (s {{.Slice}}) {{ typePrefix "Sort" }}() { sort.Sort(s) }

// {{ typePrefix "FilterThis" }} filters the current slice by predicate f without memory allocation
func (s {{.Slice}}) {{ typePrefix "FilterThis" }} (f func(*{{.Struct}}) bool) {{.Slice}} {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}

// {{ typePrefix "Filter" }} returns a new slice filtered by predicate f
func (s {{.Slice}}) {{ typePrefix "Filter" }} (f func(*{{.Struct}}) bool) {{.Slice}} {
	sl := make({{.Slice}}, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}

// {{ typePrefix "FilterNot" }} will return a new {{.Slice}} that do not match
// by calling the function f
func (s {{.Slice}}) {{ typePrefix "FilterNot" }}(f func(*{{.Struct}}) bool) {{.Slice}} {
	sl := make({{.Slice}}, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}

// {{ typePrefix "Map" }} will run function f on all items in {{.Slice}}
func (s {{.Slice}}) Map(f func(*{{.Struct}}) ) {{.Slice}} {
	for i := range s {
		f(s[i])
	}
	return s
}

// {{ typePrefix "Cut" }} will remove items i through j-1
func (s *{{.Slice}}) {{ typePrefix "Cut" }}(i, j int) {
	z := *s // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this should avoid the memory leak
	}
	z = z[:len(z)-j+i]
	*s = z
}

// {{ typePrefix "Delete" }} will remove item i
func (s *{{.Slice}}) {{ typePrefix "Delete" }}(i int) {
	z := *s // copy the slice header
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	*s = z
}

// {{ typePrefix "Insert" }} will place a new item at position i
func (s *{{.Slice}}) {{ typePrefix "Insert" }}(n *{{.Struct}}, i int) {
	z := *s // copy the slice header
	z = append(z, &{{.Struct}}{})
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}

// {{ typePrefix "Append" }} will add a new item at the end of {{.Slice}}
func (s *{{.Slice}}) {{ typePrefix "Append" }}(n ...*{{.Struct}}) {
	*s = append(*s, n...)
}

// {{ typePrefix "Prepend" }} will add a new item at the beginning of {{.Slice}}
func (s *{{.Slice}}) {{ typePrefix "Prepend" }}(n *{{.Struct}}) {
	s.Insert(n, 0)
}

type Extract{{.NameRaw | camelize}} struct {
{{ range $k,$c := .Columns }} {{$c.Name | camelize }} func() []{{$c.GetGoPrimitive false}}
{{end}} }

// {{ typePrefix "Extract" }} generates slices from all fields
func (s {{$.Slice}}) {{ typePrefix "Extract" }}() Extract{{.NameRaw | camelize}} {
	return Extract{{.NameRaw | camelize}} {
		{{ range $k,$c := .Columns }} {{$c.Name | camelize }} : func() []{{$c.GetGoPrimitive false}} {
			ext := make([]{{$c.GetGoPrimitive false}}, 0, len(s))
			for _, v := range s {
				ext = append(ext, v.{{ $c.Name | camelize }}{{dbrType $c}})
			}
			return ext
		},
		{{end}} }
}

`

const tplEAValues = `
{{range $typeCode,$valueTables := .TypeCodeValueTables}}
// Get{{ $typeCode | prepareVar }}ValueStructure returns for an EAV index the table structure.
// Important also if you have custom value tables
func Get{{ $typeCode | prepareVar }}ValueStructure(i eav.ValueIndex) (*csdb.Table, error) {
	switch i {
	{{range $vt,$v := $valueTables }}case eav.EntityType{{ $v | prepareVar }}:
		return TableCollection.Structure(TableIndex{{ $vt | prepareVar }})
    {{end}}	}
	return nil, eav.ErrEntityTypeValueNotFound
}
{{end}}
`

//// IDs returns an Int64Slice with all primary key ids
//func (s {{.Slice}}) IDs() utils.Int64Slice {
//id := make(utils.Int64Slice, len(s))
//for i, r := range s {
//id[i] = r.{{ .primary key }}
//}
//return id
//}
