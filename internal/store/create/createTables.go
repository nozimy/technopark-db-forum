package create

import "database/sql"

func CreateTables(db *sql.DB) error {
	_ = dropAllTables(db)

	forumsQuery := `CREATE TABLE IF NOT EXISTS forums (
		id bigserial not null primary key,
		name varchar
	);`
	if _, err := db.Exec(forumsQuery); err != nil {
		return err
	}

	return nil
}

func dropAllTables(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS forums;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
