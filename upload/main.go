package main

import (
	"fmt"
	//"reflect"
	"database/sql"
	"encoding/base64"
	"strings"
	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser    = "root"
	dbPass    = ""
	dbName    = "disenadorlogodb"
	dbCharset = "utf8"
)

func main() {

	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName+"?charset="+dbCharset)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT idElemento, svg FROM elementos WHERE tipo = 'ICONO'")

	if err != nil {

		panic(err.Error())
	}

	rows, err := stmt.Query()

	if err != nil {
		panic(err.Error())
	}

	var colores = []string{"#F5D327","#70C041","#51A7F9","#B36AE2"}
	var logos []map[string]string
	var contador int
	for rows.Next() {

		//logo := make(map[string]interface{})

		var id, svg string

		rows.Scan(&id, &svg)

		svgNatural, err := base64.StdEncoding.DecodeString(svg)
		if err != nil {
			
		} else{
			if !strings.HasPrefix(string(svgNatural), "<svg fill=") {

				svgRoto := "<svg" + strings.SplitN(string(svgNatural), "<svg", 2)[1]
				if !strings.HasPrefix(svgRoto, "<svg fill="){

					colorAzar := colores[rand.Intn(3)]

					svgRoto = "<svg fill=\"" + colorAzar + "\"" + strings.SplitN(string(svgNatural), "<svg", 2)[1]
					contador++

					logo := make(map[string]string)

					logo["id"] = id
					logo["svg"] = base64.StdEncoding.EncodeToString([]byte(svgRoto)) 
					logo["color"] = colorAzar

					fmt.Println(logo)

					logos = append(logos, logo)
				}
				
			}
		}

	}



	for _, logoN := range logos {
		stmt, err := db.Prepare("UPDATE elementos SET svg = ?, color = ? WHERE idElemento = ?")
		if err != nil {

			panic(err.Error())
		}

		stmt.Exec(logoN["svg"], logoN["color"], logoN["id"])

	}


	fmt.Println(contador)

	//fmt.Println(logos)

}
