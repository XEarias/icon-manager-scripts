package main

import "flag"
import "fmt"
import "./find"

func main() {

	action := flag.String("action", "", "Valores aceptados:\"find\" ó \"upload\"")

	flag.Parse()
	
	fmt.Println(*action)

}
