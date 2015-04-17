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

// Package store implements the handling of websites, groups and stores
package store

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils"
)

const (
	DefaultStoreId int64 = 0
)

type (
	// StoreBucket contains two maps for faster retrieving of the store index and the store collection
	// Only used in generated code. Implements interface StoreGetter.
	StoreBucket struct {
		// store collection
		s  TableStoreSlice
		im *indexMap
	}
	// StoreGetter methods to retrieve a store pointer
	StoreGetter interface {
		Get(int64, ...string) (*TableStore, error)
		Collection() TableStoreSlice
	}
)

var (
	ErrStoreNotFound = errors.New("Store not found")
)
var _ StoreGetter = (*StoreBucket)(nil)

// NewStoreBucket returns a new pointer to a StoreBucket.
func NewStoreBucket(s TableStoreSlice) *StoreBucket {
	return &StoreBucket{
		im: (&indexMap{}).populateStore(s),
		s:  s,
	}
}

// Get uses the  store id as 1st arg to return a TableStore struct or if a store code is provided
// as 2nd then ignores the store id
func (s *StoreBucket) Get(sID int64, sc ...string) (*TableStore, error) {
	if len(sc) == 1 {
		if i, ok := s.im.code[sc[0]]; ok {
			return s.s[i], nil
		}
		return nil, ErrStoreNotFound
	}
	if i, ok := s.im.id[sID]; ok {
		return s.s[i], nil
	}
	return nil, ErrStoreNotFound
}

// Collection returns the TableStoreSlice
func (s *StoreBucket) Collection() TableStoreSlice { return s.s }

// Load uses a dbr session to load all data from the core_store table into the current slice.
// The variadic 2nd argument can be a call back function to manipulate the select.
// Additional columns or joins cannot be added. This method receiver should only be used in development.
// @see app/code/Magento/Store/Model/Resource/Store/Collection.php::Load() for sort order
func (s *TableStoreSlice) Load(dbrSess dbr.SessionRunner, cbs ...csdb.DbrSelectCb) (int, error) {
	return loadSlice(dbrSess, TableIndexStore, &(*s), append(cbs, func(sb *dbr.SelectBuilder) *dbr.SelectBuilder {
		sb.OrderBy("CASE WHEN main_table.store_id = 0 THEN 0 ELSE 1 END ASC")
		sb.OrderBy("main_table.sort_order ASC")
		return sb.OrderBy("main_table.name ASC")
	})...)
}

// Len returns the length
func (s TableStoreSlice) Len() int { return len(s) }

// ByGroupID returns a new slice with all stores belonging to a group id
func (s TableStoreSlice) FilterByGroupID(id int64) TableStoreSlice {
	return s.Filter(func(store *TableStore) bool {
		return store.GroupID == id
	})
}

// FilterByWebsiteID returns a new slice with all stores belonging to a website id
func (s TableStoreSlice) FilterByWebsiteID(id int64) TableStoreSlice {
	return s.Filter(func(store *TableStore) bool {
		return store.WebsiteID == id
	})
}

// Filter returns a new slice filtered by predicate f
func (s TableStoreSlice) Filter(f func(*TableStore) bool) TableStoreSlice {
	var tss TableStoreSlice
	for _, v := range s {
		if v != nil && f(v) {
			tss = append(tss, v)
		}
	}
	return tss
}

// Codes returns a StringSlice with all store codes
func (s TableStoreSlice) Codes() utils.StringSlice {
	c := make(utils.StringSlice, len(s))
	for i, store := range s {
		c[i] = store.Code.String
	}
	return c
}

// IDs returns an Int64Slice with all store ids
func (s TableStoreSlice) IDs() utils.Int64Slice {
	id := make(utils.Int64Slice, len(s))
	for i, store := range s {
		id[i] = store.StoreID
	}
	return id
}

func (s TableStore) IsDefault() bool {
	return s.StoreID == DefaultStoreId
}

/*
	@todo implement Magento\Store\Model\Store
*/
