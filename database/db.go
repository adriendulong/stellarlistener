package database

import (
	"database/sql"

	// Postgres
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var d *sql.DB

// GetDB return an instance of the database
func GetDB() (db *sql.DB, err error) {
	if d == nil {
		log.Info("CREATE DATABASE CONN")
		d, err = sql.Open("postgres", "user=adriendulong dbname=stellar sslmode=disable")
		if err != nil {
			log.Error(err)
			return d, err
		}

		if err := d.Ping(); err != nil {
			log.Error(err)
			return d, err
		}
	}
	return d, nil
}

// CreateDatabase create a database
func CreateDatabase() {
	database, err := GetDB()
	if err != nil {
		log.Error(err)
	}

	_, err = database.Exec(`CREATE TABLE operations (
		id 					text NOT NULL PRIMARY KEY,
		type				text,
		type_i				integer,
		starting_balance	text,	
		funder				text,
		asset_type			text,
		asset_code			text,
		asset_issuer		text,
		from_account		text,
		to_account			text,
		amount				text,
		sellingasset_code	text,
		buying_asset_code	text,
		price				text,
		price_r				text
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

//DeleteDatabase drop tables
func DeleteDatabase() {
	database, err := GetDB()
	if err != nil {
		log.Error(err)
	}

	_, err = database.Exec(`DROP TABLE operations`)
	if err != nil {
		log.Fatal(err)
	}
}
