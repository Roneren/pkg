// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml_test

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ dml.ColumnMapper = (*someRecord)(nil)

type someRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr someRecord) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm.Int(&sr.SomethingID).Int64(&sr.UserID).Bool(&sr.Other).Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c {
		case "something_id":
			cm.Int(&sr.SomethingID)
		case "user_id":
			cm.Int64(&sr.UserID)
		case "other":
			cm.Bool(&sr.Other)
		default:
			return errors.NotFound.Newf("[dml_test] Column %q not found", c)
		}
	}
	return cm.Err()
}

func TestInsert_Bind(t *testing.T) {
	t.Parallel()

	objs := []someRecord{{1, 88, false}, {2, 99, true}, {3, 101, true}}
	wantArgs := []interface{}{int64(1), int64(88), false, int64(2), int64(99), true, int64(3), int64(101), true}

	t.Run("valid with multiple records", func(t *testing.T) {
		compareToSQL(t,
			dml.NewInsert("a").
				AddColumns("something_id", "user_id", "other").
				AddRecords(objs[0]).AddRecords(objs[1], objs[2]).
				AddOnDuplicateKey(
					dml.Column("something_id").Int64(99),
					dml.Column("user_id").Values(),
				),
			errors.NoKind,
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("without columns, all columns requested, with AddOnDuplicateKey", func(t *testing.T) {
		compareToSQL(t,
			dml.NewInsert("a").
				SetRecordPlaceHolderCount(3).
				AddRecords(objs[0]).AddRecords(objs[1], objs[2]).
				AddOnDuplicateKey(
					dml.Column("something_id").Int64(99),
					dml.Column("user_id").Values(),
				),
			errors.NoKind,
			"INSERT INTO `a` VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("without columns, all columns requested", func(t *testing.T) {

		customers := []*customerEntity{
			{EntityID: 11, Firstname: "Karl Gopher", StoreID: 0x7, LifetimeSales: dml.MakeNullFloat64(47.11), VoucherCodes: exampleStringSlice{"1FE9983E", "28E76FBC"}},
			{EntityID: 12, Firstname: "Fung Go Roo", StoreID: 0x7, LifetimeSales: dml.MakeNullFloat64(28.94), VoucherCodes: exampleStringSlice{"4FE7787E", "15E59FBB", "794EFDE8"}},
			{EntityID: 13, Firstname: "John Doe", StoreID: 0x6, LifetimeSales: dml.MakeNullFloat64(138.54), VoucherCodes: exampleStringSlice{""}},
		}

		compareToSQL(t,
			dml.NewInsert("customer_entity").
				SetRecordPlaceHolderCount(5). // mandatory because no columns provided!
				AddRecords(customers[0], customers[1], customers[2]),
			errors.NoKind,
			"INSERT INTO `customer_entity` VALUES (?,?,?,?,?),(?,?,?,?,?),(?,?,?,?,?)",
			"INSERT INTO `customer_entity` VALUES (11,'Karl Gopher',7,47.11,'1FE9983E|28E76FBC'),(12,'Fung Go Roo',7,28.94,'4FE7787E|15E59FBB|794EFDE8'),(13,'John Doe',6,138.54,'')",
			int64(11), "Karl Gopher", int64(7), 47.11, "1FE9983E|28E76FBC", int64(12), "Fung Go Roo", int64(7), 28.94, "4FE7787E|15E59FBB|794EFDE8", int64(13), "John Doe", int64(6), 138.54, "",
		)

	})
	t.Run("column not found", func(t *testing.T) {
		objs := []someRecord{{1, 88, false}, {2, 99, true}}
		compareToSQL(t,
			dml.NewInsert("a").AddColumns("something_it", "user_id", "other").AddRecords(objs[0]).AddRecords(objs[1]),
			errors.NotFound,
			"",
			"",
		)
	})
}

func TestInsert_Prepare(t *testing.T) {
	t.Parallel()

	t.Run("ToSQL Error", func(t *testing.T) {
		in := &dml.Insert{}
		in.AddColumns("a", "b")
		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.Empty.Match(err))
	})

	t.Run("DB Error", func(t *testing.T) {
		in := &dml.Insert{
			Into: "table",
		}
		in.DB = dbMock{
			error: errors.AlreadyClosed.Newf("Who closed myself?"),
		}
		in.AddColumns("a", "b").AddValuesUnsafe(1, true)

		stmt, err := in.Prepare(context.TODO())
		assert.Nil(t, stmt)
		assert.True(t, errors.AlreadyClosed.Match(err), "%+v", err)
	})

	t.Run("ExecArgs One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `customer_entity` (`email`,`group_id`,`created_at`) VALUES (?,?,?)"))
		prep.ExpectExec().WithArgs("a@b.c", 33, now()).WillReturnResult(sqlmock.NewResult(4, 0))
		prep.ExpectExec().WithArgs("x@y.z", 44, now().Add(time.Minute)).WillReturnResult(sqlmock.NewResult(5, 0))

		stmt, err := dml.NewInsert("customer_entity").
			AddColumns("email", "group_id", "created_at").
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			email      string
			groupID    int
			created_at time.Time
			insertID   int64
		}{
			{"a@b.c", 33, now(), 4},
			{"x@y.z", 44, now().Add(time.Minute), 5},
		}

		args := dml.MakeArgs(3)
		for i, test := range tests {
			args = args[:0]

			res, err := stmt.
				WithArguments(args.String(test.email).Int(test.groupID).Time(test.created_at)).
				ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different LastInsertIDs", i)
		}
	})

	t.Run("ExecArgs Multi Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `customer_entity` (`email`,`group_id`) VALUES (?,?),(?,?)"))
		prep.ExpectExec().WithArgs("a@b.c", 33, "d@e.f", 33).WillReturnResult(sqlmock.NewResult(6, 0))
		prep.ExpectExec().WithArgs("x@y.z", 44, "u@v.w", 44).WillReturnResult(sqlmock.NewResult(7, 0))

		stmt, err := dml.NewInsert("customer_entity").
			AddColumns("email", "group_id").
			SetRowCount(2).
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err)
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			email1   string
			groupID1 int
			email2   string
			groupID2 int
			insertID int64
		}{
			{"a@b.c", 33, "d@e.f", 33, 6},
			{"x@y.z", 44, "u@v.w", 44, 7},
		}

		args := dml.MakeArgs(4)
		for i, test := range tests {
			args = args[:0]

			res, err := stmt.
				WithArguments(args.String(test.email1).Int(test.groupID1).String(test.email2).Int(test.groupID2)).
				ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different LastInsertIDs", i)
		}
	})

	t.Run("ExecRecord One Row", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `dml_person` (`name`,`email`) VALUES (?,?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(4, 0))
		prep.ExpectExec().WithArgs("John Doe", "john@doe.go").WillReturnResult(sqlmock.NewResult(5, 0))

		stmt, err := dml.NewInsert("dml_person").
			AddColumns("name", "email").
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		tests := []struct {
			name     string
			email    string
			insertID int64
		}{
			{"Peter Gopher", "peter@gopher.go", 4},
			{"John Doe", "john@doe.go", 5},
		}

		for i, test := range tests {

			p := &dmlPerson{
				Name:  test.name,
				Email: dml.MakeNullString(test.email),
			}

			res, err := stmt.WithRecords(p).ExecContext(context.TODO())
			if err != nil {
				t.Fatalf("Index %d => %+v", i, err)
			}
			lid, err := res.LastInsertId()
			if err != nil {
				t.Fatalf("Result index %d with error: %s", i, err)
			}
			assert.Exactly(t, test.insertID, lid, "Index %d has different LastInsertIDs", i)
			assert.Exactly(t, test.insertID, p.ID, "Index %d and model p has different LastInsertIDs", i)
		}
	})

	t.Run("ExecContext", func(t *testing.T) {
		dbc, dbMock := dmltest.MockDB(t)
		defer dmltest.MockClose(t, dbc, dbMock)

		prep := dbMock.ExpectPrepare(dmltest.SQLMockQuoteMeta("INSERT INTO `dml_person` (`name`,`email`) VALUES (?,?)"))
		prep.ExpectExec().WithArgs("Peter Gopher", "peter@gopher.go").WillReturnResult(sqlmock.NewResult(4, 0))

		stmt, err := dml.NewInsert("dml_person").
			AddColumns("name", "email").
			WithDB(dbc.DB).
			Prepare(context.TODO())
		require.NoError(t, err, "failed creating a prepared statement")
		defer func() {
			require.NoError(t, stmt.Close(), "Close on a prepared statement")
		}()

		res, err := stmt.ExecContext(context.TODO(), "Peter Gopher", "peter@gopher.go")
		require.NoError(t, err, "failed to execute ExecContext")

		lid, err := res.LastInsertId()
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, int64(4), lid, "Different LastInsertIDs")
	})
}

func TestInsert_WithLogger(t *testing.T) {
	uniID := new(int32)
	rConn := createRealSession(t)
	defer dmltest.Close(t, rConn)

	var uniqueIDFunc = func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 4))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	require.NoError(t, rConn.Options(dml.WithLogger(lg, uniqueIDFunc)))

	t.Run("ConnPool", func(t *testing.T) {
		d := rConn.InsertInto("dml_people").Replace().AddColumns("email", "name").AddValuesUnsafe("a@b.c", "John")

		t.Run("Exec", func(t *testing.T) {
			defer func() {
				buf.Reset()
				d.IsInterpolate = false
			}()
			_, err := d.Interpolate().Exec(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ04\" insert_id: \"UNIQ08\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ08*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()
			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ04\" insert_id: \"UNIQ08\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ08*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := rConn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				_, err := tx.InsertInto("dml_people").Replace().AddColumns("email", "name").AddValuesUnsafe("a@b.c", "John").Interpolate().Exec(context.TODO())
				return err
			}))
			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ04\" tx_id: \"UNIQ12\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" tx_id: \"UNIQ12\" insert_id: \"UNIQ16\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ16*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\nDEBUG Commit conn_pool_id: \"UNIQ04\" tx_id: \"UNIQ12\" duration: 0\n",
				buf.String())
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.TODO())
		require.NoError(t, err)

		d := conn.InsertInto("dml_people").Replace().AddColumns("email", "name").AddValuesUnsafe("a@b.zeh", "J0hn")

		t.Run("Exec", func(t *testing.T) {
			defer func() {
				buf.Reset()
				d.IsInterpolate = false
			}()

			_, err := d.Interpolate().Exec(context.TODO())
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ24*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\n",
				buf.String())
		})

		t.Run("Prepare", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ24*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\n",
				buf.String())
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			defer buf.Reset()

			stmt, err := d.Prepare(context.TODO())
			require.NoError(t, err)
			defer stmt.Close()

			_, err = stmt.Exec(context.TODO(), "mail@e.de", "Hans")
			require.NoError(t, err)

			assert.Exactly(t, "DEBUG Prepare conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ24*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" insert_id: \"UNIQ24\" table: \"dml_people\" duration: 0 arg_len: 2\n",
				buf.String())
		})

		t.Run("Tx Commit", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.NoError(t, tx.Wrap(func() error {
				_, err := tx.InsertInto("dml_people").Replace().AddColumns("email", "name").AddValuesUnsafe("a@b.c", "John").Interpolate().Exec(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ28\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ28\" insert_id: \"UNIQ32\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ32*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\nDEBUG Commit conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ28\" duration: 0\n",
				buf.String())
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			defer buf.Reset()
			tx, err := conn.BeginTx(context.TODO(), nil)
			require.NoError(t, err)
			require.Error(t, tx.Wrap(func() error {
				_, err := tx.InsertInto("dml_people").Replace().AddColumns("email", "name").
					AddValues(dml.MakeArgs(1).String("only one arg provided")).Interpolate().Exec(context.TODO())
				return err
			}))

			assert.Exactly(t, "DEBUG BeginTx conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ36\"\nDEBUG Exec conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ36\" insert_id: \"UNIQ40\" table: \"dml_people\" duration: 0 sql: \"REPLACE /*ID:UNIQ40*/ INTO `dml_people` (`email`,`name`) VALUES (?,?)\"\nDEBUG Rollback conn_pool_id: \"UNIQ04\" conn_id: \"UNIQ20\" tx_id: \"UNIQ36\" duration: 0\n",
				buf.String())
		})
	})
}
