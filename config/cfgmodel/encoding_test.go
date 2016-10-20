// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package cfgmodel_test

import (
	"testing"

	"encoding/json"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {

	type tt struct {
		Str   string
		Int   int
		Float float64
	}

	wantTT := &tt{
		Str:   "H3llo Gphers",
		Int:   5,
		Float: 3.14159,
	}

	var wantJSON = []byte(`{"Str":"H3llo Gphers","Int":5,"Float":3.14159}`)
	const cfgPath = "aa/bb/cc"

	b := cfgmodel.NewEncode(
		cfgPath,
		cfgmodel.WithEncoder(json.Marshal, json.Unmarshal),
		cfgmodel.WithScopeStore(),
	)
	wantPath := cfgpath.MustNewByParts(cfgPath).String() // Default Scope

	haveTT := &tt{}

	haveErr := b.Get(cfgmock.NewService(
		cfgmock.PathValue{
			wantPath: wantJSON,
		}).NewScoped(34, 4), haveTT)
	if haveErr != nil {
		t.Fatal(haveErr)
	}
	assert.Exactly(t, wantTT, haveTT)

	mw := new(cfgmock.Write)
	b.Write(mw, wantTT, scope.Store.Pack(12))
	assert.Exactly(t, wantJSON, mw.ArgValue)
	assert.Exactly(t, "stores/12/aa/bb/cc", mw.ArgPath)
}