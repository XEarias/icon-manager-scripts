package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

/*DB es la estructura para la configuracion de la BD  */
type DB struct {
	ip       string
	port     string
	name     string
	user     string
	password string
}

/*DBConfig son los datos de la bd del diseñador*/
var DBConfig = DB{
	ip:       "127.0.0.1",
	port:     "3306",
	name:     "disenadorlogodb_tags_prov",
	user:     "logoPro",
	password: "&rJ-fZ:1uZ24",
}

/*Tag define una tag con sus traducciones*/
type Tag struct {
	Categoria string
	ESP       string
	ENG       string
	POR       string
}

func (t *Tag) insert(tagDB *sql.DB) error {

	stmtInsert, err := tagDB.Prepare("INSERT INTO tags(categoria, ESP, ENG, POR) VALUES ( ?, LOWER(?), LOWER(?), LOWER(?))")

	if err != nil {
		return err
	}

	defer stmtInsert.Close()

	res, err := stmtInsert.Exec(t.Categoria, t.ESP, t.ENG, t.POR)

	if err != nil {
		return err
	}

	_, err = res.LastInsertId()

	if err != nil {
		return err
	}

	return nil

}

func createDB() error {

	db, err := sql.Open("mysql", DBConfig.user+":"+DBConfig.password+"@tcp("+DBConfig.ip+":"+DBConfig.port+")/")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + DBConfig.name)
	if err != nil {
		return err
	}

	_, err = db.Exec("USE " + DBConfig.name)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tags ( id  int NOT NULL PRIMARY KEY auto_increment, categoria int NOT NULL, ESP varchar(255) NOT NULL UNIQUE, ENG varchar(255) NOT NULL UNIQUE, POR varchar(255) NOT NULL UNIQUE)")
	if err != nil {
		return err
	}

	return nil

}

func main() {

	// Load a csv file.
	tagsCSV, err := os.Open("./tags.csv")

	if err != nil {
		return
	}

	fmt.Println("Checkeando la BD")
	errCreating := createDB()

	if errCreating != nil {
		fmt.Println(errCreating.Error())
		return
	}

	tagDB, err := sql.Open("mysql", DBConfig.user+":"+DBConfig.password+"@tcp("+DBConfig.ip+":"+DBConfig.port+")/"+DBConfig.name)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer tagDB.Close()

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(tagsCSV))
	for {

		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		tag := &Tag{
			Categoria: record[0],
			ESP:       record[2],
			ENG:       record[1],
			POR:       record[3],
		}

		errInsert := tag.insert(tagDB)

		if errInsert != nil {
			fmt.Println("Error al insertar tag: " + errInsert.Error())
		}
	}

	return

}
