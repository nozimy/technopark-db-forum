package create

import (
	"database/sql"
	"io/ioutil"
)

func CreateTables(db *sql.DB) error {
	//_ = dropAllTables(db)

	//forumsQuery := `CREATE TABLE IF NOT EXISTS forums (
	//	id bigserial not null primary key,
	//	name varchar
	//);`
	//if _, err := db.Exec(forumsQuery); err != nil {
	//	return err
	//}

	file, err := ioutil.ReadFile("./assets/db/postgres/base.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(file))
	if err != nil {
		return err
	}

	//requests := strings.Split(string(file), ";")
	//for _, request := range requests {
	//	_, err := db.Exec(request)
	//	if err != nil {
	//		return err
	//	}
	//}

	//if _, err := db.Exec(string(file)); err != nil {
	//	return err
	//}

	return nil
}

func dropAllTables(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS forums;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
