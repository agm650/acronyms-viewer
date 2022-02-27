package webserver

import (
	"database/sql"
	"os"

	"github.com/apex/log"

	_ "modernc.org/sqlite"
)

/*
Copyright © 2022 DANDOY Luc

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

var AcroDB *sql.DB = nil

func loadDB(path string) error {
	ctx := log.WithFields(log.Fields{
		"file":     "webserver/database.go",
		"function": "loadDB",
	})

	_, err := os.Stat(path)
	if err != nil {
		// Database doesn't exist, create it
		ctx.Warnf("Database doesnt exist. Creating it")

		AcroDB, err = sql.Open("sqlite", path)
		if err != nil {
			return err
		}

		// Creating an empty database
		sqlStmt := `
		pragma synchronous = OFF;
		pragma journal_mode = MEMORY;
		pragma temp_store = MEMORY;
		pragma locking_mode = EXCLUSIVE;
		pragma auto_vacuum = 1;

		CREATE TABLE acronym (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			acronym TEXT NOT NULL,
			description TEXT NOT NULL,
			contributor TEXT
		);
		CREATE VIRTUAL TABLE acronym_fts USING fts5(
			acronym,
			description,
			contributor UNINDEXED,
			content='acronym',
			content_rowid='id'
		);
		CREATE TRIGGER acronym_ai AFTER INSERT ON acronym
		BEGIN
				INSERT INTO acronym_fts (rowid, acronym, description)
				VALUES (new.id, new.acronym, new.description);
		END;

		CREATE TRIGGER acronym_ad AFTER DELETE ON acronym
		BEGIN
				INSERT INTO acronym_fts (acronym_fts, rowid, acronym, description)
				VALUES ('delete', old.id, old.acronym, old.description);
		END;

		CREATE TRIGGER acronym_au AFTER UPDATE ON acronym
		BEGIN
				INSERT INTO acronym_fts (acronym_fts, rowid, acronym, description)
				VALUES ('delete', old.id, old.acronym, old.description);
				INSERT INTO acronym_fts (rowid, acronym, description)
				VALUES (new.id, new.acronym, new.description);
		END;
		`
		_, err = AcroDB.Exec(sqlStmt)
		if err != nil {
			return err
		}

	} else {
		ctx.Warnf("Database alsready exist. Using it")
		AcroDB, err = sql.Open("sqlite", path)
		if err != nil {
			return err
		}
	}

	return err
}

func closeDB() error {
	return AcroDB.Close()
}

func getDefinition(acro string) ([]string, error) {
	ctx := log.WithFields(log.Fields{
		"file":     "webserver/database.go",
		"function": "getDefinition",
	})
	defToReturn := make([]string, 0, 20)

	if AcroDB == nil {
		ctx.Fatal("DB not opened")
	}
	stmt, err := AcroDB.Prepare("SELECT description FROM acronym_fts WHERE acronym MATCH ? ORDER BY rank")
	if err != nil {
		return defToReturn, err
	}
	defer stmt.Close()

	result, err := stmt.Query(acro)
	if err != nil {
		ctx.Warnf("%d definitions found for %s", len(defToReturn), acro)
		return defToReturn, err
	}

	for result.Next() {
		var definition string
		err = result.Scan(&definition)
		if err != nil {
			continue
		}
		defToReturn = append(defToReturn, definition)
	}

	ctx.Warnf("%d definitions found for %s", len(defToReturn), acro)
	return defToReturn, nil
}

func addDefinition(acro string, definition string, contrib string) (int64, error) {
	stmt, err := AcroDB.Prepare("INSTER INTO acronym(acronym, description, contributor) VALUES(?,?,?);")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(acro)
	if err != nil {
		return -1, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return lastID, nil
}
