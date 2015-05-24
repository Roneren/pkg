// Copyright 2015 CoreStore Authors
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

package config_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/stretchr/testify/assert"
)

func TestScopeBits(t *testing.T) {
	const (
		scope1 config.ScopeID = iota + 1
		scope2
		scope3
		scope4
		scope5
	)

	tests := []struct {
		have    []config.ScopeID
		want    config.ScopeID
		notWant config.ScopeID
		human   []string
	}{
		{[]config.ScopeID{scope1, scope2}, scope2, scope3, []string{"ScopeDefault", "ScopeWebsite"}},
		{[]config.ScopeID{scope3, scope4}, scope3, scope2, []string{"ScopeGroup", "ScopeStore"}},
		{[]config.ScopeID{scope4, scope5}, scope4, scope2, []string{"ScopeStore", "ScopeID(5)"}},
	}

	for _, test := range tests {
		var b config.ScopePerm
		b.Set(test.have...)
		if b.Has(test.want) == false {
			t.Errorf("%d should contain %d", b, test.want)
		}
		if b.Has(test.notWant) {
			t.Errorf("%d should not contain %d", b, test.notWant)
		}
		assert.EqualValues(t, test.human, b.Human())
	}
}

func TestScopeApplyDefaults(t *testing.T) {
	pkgCfg := config.NewConfiguration(
		&config.Section{
			ID: "contact",
			Groups: config.GroupSlice{
				&config.Group{
					ID: "contact",
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `contact/contact/enabled`,
							ID:      "enabled",
							Default: true,
						},
					},
				},
				&config.Group{
					ID: "email",
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `contact/email/recipient_email`,
							ID:      "recipient_email",
							Default: `hello@example.com`,
						},
						&config.Field{
							// Path: `contact/email/sender_email_identity`,
							ID:      "sender_email_identity",
							Default: 2.7182818284590452353602874713527,
						},
						&config.Field{
							// Path: `contact/email/email_template`,
							ID:      "email_template",
							Default: 4711,
						},
					},
				},
			},
		},
	)
	s := config.NewManager()
	s.ApplyDefaults(pkgCfg)
	cer, err := pkgCfg.FindFieldByPath("contact", "email", "recipient_email")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Exactly(t, cer.Default.(string), s.GetString(config.Path("contact/email/recipient_email")))
}
