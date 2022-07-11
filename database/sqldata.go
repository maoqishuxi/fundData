package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type QueryTable struct {
	Id         int
	EquityTime string
	NetValue   float64
	GrowthRate float64
}

func NewConnPool() *sql.DB {
	db, err := sql.Open("sqlite3", "./strategy.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateTable(db *sql.DB, tableName string) {
	sqlStmt := `
	create table if not exists %s(
		id 			integer not null primary key,
		equityTime	text,
		netValue	integer,
		growthRate	integer
	);`

	sqlStmt = fmt.Sprintf(sqlStmt, tableName)
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func InsertData(db *sql.DB, tableName string, queryRows []QueryTable) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	sqlField := fmt.Sprintf("insert into %s(id, equityTime, netValue, growthRate) values(?, ?, ?, ?)", tableName)
	stmt, err := tx.Prepare(sqlField)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, query := range queryRows {
		_, err := stmt.Exec(query.Id, query.EquityTime, query.NetValue, query.GrowthRate)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func QueryData(db *sql.DB, tableName string, limits int) []QueryTable {
	sqlField := fmt.Sprintf("select * from %s order by id desc limit ?", tableName)

	var id int
	var equityTime string
	var netValue, growthRate float64
	rows, err := db.Query(sqlField, limits)
	if err != nil {
		log.Fatal(err)
	}

	result := make([]QueryTable, 0)
	for rows.Next() {
		if err = rows.Scan(&id, &equityTime, &netValue, &growthRate); err != nil {
			log.Fatal(err)
		}
		result = append(result, QueryTable{id, equityTime, netValue, growthRate})
	}

	return result
}

func RemoveData(db *sql.DB, tableName string, field int) {
	sqlField := fmt.Sprintf("delete from %s where id=?", tableName)

	stmt, err := db.Prepare(sqlField)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(field)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateData[k comparable](db *sql.DB, tableName string, update k, check k) {
	sqlField := fmt.Sprintf("update %s set ? where ?", tableName)
	stmt, err := db.Prepare(sqlField)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(update, check)
	if err != nil {
		log.Fatal(err)
	}
}
