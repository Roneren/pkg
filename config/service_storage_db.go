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

package config

import (
	"fmt"
	"time"

	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cast"
)

// DBStorage connects the MySQL DB with the config.Service type.
type DBStorage struct {
	// All is a SQL statement for the all keys query
	All *csdb.ResurrectStmt
	// Read is a SQL statement for selecting a value from a path/key
	Read *csdb.ResurrectStmt
	// Write statement inserts or updates a value
	Write *csdb.ResurrectStmt
}

// NewDBStorage creates a new pointer with resurrecting prepared SQL statements.
// Default logger for the three underlying ResurrectStmt type is the PkgLog.
//
// All has an idle time of 15s. Read an idle time of 10s. Write an idle time of 30s.
func NewDBStorage(p csdb.Preparer) (*DBStorage, error) {
	// todo: instead of logging the error we may write it into an
	// error channel and the gopher who calls NewDBStorage is responsible
	// for continuously reading from the error channel. or we accept an error channel
	// as argument here and then writing to it ...

	dbs := &DBStorage{
		All: csdb.NewResurrectStmt(p, fmt.Sprintf(
			"SELECT scope,scope_id,path FROM `%s` ORDER BY scope,scope_id,path",
			TableCollection.Name(TableIndexCoreConfigData),
		)),
		Read: csdb.NewResurrectStmt(p, fmt.Sprintf(
			"SELECT `value` FROM `%s` WHERE `scope`=? AND `scope_id`=? AND `path`=?",
			TableCollection.Name(TableIndexCoreConfigData),
		)),

		Write: csdb.NewResurrectStmt(p, fmt.Sprintf(
			"INSERT INTO `%s` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=?",
			TableCollection.Name(TableIndexCoreConfigData),
		)),
	}
	dbs.All.Idle = time.Second * 15
	dbs.All.Log = PkgLog
	dbs.Read.Idle = time.Second * 10
	dbs.Read.Log = PkgLog
	dbs.Write.Idle = time.Second * 30
	dbs.Write.Log = PkgLog
	// in the future we may add errors ... just to have for now the func signature
	return dbs, nil
}

// MustNewDBStorage same as NewDBStorage but panics on error
func MustNewDBStorage(p csdb.Preparer) *DBStorage {
	s, err := NewDBStorage(p)
	if err != nil {
		panic(err)
	}
	return s
}

// Start starts the internal idle time checker for the resurrecting SQL statements.
func (dbs *DBStorage) Start() *DBStorage {
	dbs.All.StartIdleChecker()
	dbs.Read.StartIdleChecker()
	dbs.Write.StartIdleChecker()
	return dbs
}

// Stop stops the internal goroutines for idle time checking. Returns the
// first occurring sql.Stmt.Close() error.
func (dbs *DBStorage) Stop() error {
	if err := dbs.All.StopIdleChecker(); err != nil {
		return err
	}
	if err := dbs.Read.StopIdleChecker(); err != nil {
		return err
	}
	if err := dbs.Write.StopIdleChecker(); err != nil {
		return err
	}
	return nil
}

// Set sets a value with its key. Database errors get logged as Info message.
// Enabled debug level logs the insert ID or rows affected.
func (dbs *DBStorage) Set(key path.Path, value interface{}) error {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.Write.StartStmtUse()
	defer dbs.Write.StopStmtUse()

	valStr, err := cast.ToStringE(value)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.Set.ToString", "err", err, "SQL", dbs.Write.SQL, "value", value)
		}
		return err
	}

	stmt, err := dbs.Write.Stmt()
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.Set.Write.getStmt", "err", err, "SQL", dbs.Write.SQL)
		}
		return err
	}

	pl, err := key.Level(-1)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.Set.ReverseFQPath", "err", err, "key", key)
		}
		return err
	}

	result, err := stmt.Exec(key.Scope, key.ID, pl, valStr, valStr)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.Set.Write.Exec", "err", err, "SQL", dbs.Write.SQL, "key", key, "value", value)
		}
		return err
	}
	if PkgLog.IsDebug() {
		li, err1 := result.LastInsertId()
		ra, err2 := result.RowsAffected()
		PkgLog.Debug("config.DBStorage.Set.Write.Result", "lastInsertID", li, "lastInsertIDErr", err1, "rowsAffected", ra, "rowsAffectedErr", err2, "SQL", dbs.Write.SQL, "key", key, "value", value)
	}
	return nil
}

// Get returns a value from the database by its key. It is guaranteed that the
// type in the empty interface is a string. It returns nil on error but errors
// get logged as info message
func (dbs *DBStorage) Get(key path.Path) (interface{}, error) {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.Read.StartStmtUse()
	defer dbs.Read.StopStmtUse()

	stmt, err := dbs.Read.Stmt()
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.Get.Read.getStmt", "err", err, "SQL", dbs.Read.SQL)
		}
		return nil, err
	}

	pl, err := key.Level(-1)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.Get.path.Level", "err", err, "key", key)
		}
		return nil, err
	}

	var data dbr.NullString
	err = stmt.QueryRow(key.Scope, key.ID, pl).Scan(&data)
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.Get.QueryRow", "err", err, "key", key, "path_level", pl)
		}
		return nil, err
	}
	if data.Valid {
		return data.String, nil
	}
	return nil, ErrKeyNotFound
}

// AllKeys returns all available keys. Database errors get logged as info message.
func (dbs *DBStorage) AllKeys() (path.PathSlice, error) {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.All.StartStmtUse()
	defer dbs.All.StopStmtUse()

	stmt, err := dbs.All.Stmt()
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.AllKeys.All.getStmt", "err", err, "SQL", dbs.All.SQL)
		}
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("config.DBStorage.AllKeys.All.Query", "err", err, "SQL", dbs.All.SQL)
		}
		return nil, err
	}
	defer rows.Close()

	const maxCap = 750 // Just a guess the 750
	var ret = make([]path.Path, 0, maxCap)
	var sqlScope dbr.NullString
	var sqlScopeID dbr.NullInt64
	var sqlPath dbr.NullString
	i := 0
	for rows.Next() {
		if err := rows.Scan(&sqlScope, &sqlScopeID, &sqlPath); err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("config.DBStorage.AllKeys.All.Rows.Scan", "err", err, "SQL", dbs.All.SQL)
			}
			return nil, err
		}
		if sqlPath.Valid {
			p, err := path.NewByParts(sqlPath.String)
			if err != nil {
				return ret, err
			}
			p = p.Bind(scope.FromString(sqlScope.String), sqlScopeID.Int64)
			if i < maxCap {
				ret[i] = p
			} else {
				ret = append(ret, p)
			}
		}
		sqlScope.String = ""
		sqlScope.Valid = false
		sqlScopeID.Int64 = 0
		sqlScopeID.Valid = false
		sqlPath.String = ""
		sqlPath.Valid = false
		i++
	}
	return ret, nil
}
