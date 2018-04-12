package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

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
	name:     "testdisenadorlogodb",
	user:     "logoPro",
	password: "&rJ-fZ:1uZ24",
}

/*TCPclient es la estructura para la configuracion de las conexiones TCP a mongo */
type TCPclient struct {
	ip   string
	port string
}

/*TCPconfig son los datos de la del cliente de mongodb*/
var TCPconfig = TCPclient{
	ip:   "127.0.0.1",
	port: "666",
}

var client = http.Client{}

var colors = []string{"#000000", "#ff0000", "#020100"}

/*Icon representa un icono de NounProject*/
type Icon struct {
	ID       int
	svg      string
	tags     []string
	category string
	color    string
	nounID   int
}

/*IconJSON representa un JSON individual de Icono*/
type IconJSON struct {
	URL      string   `json:"url"`
	Tags     []string `json:"tags"` //tags en ingles
	Category string   `json:"category"`
}

/*TagJSON representa un JSON individiual de Tag */
type TagJSON struct {
	Categoria string `json:"categoria,omitempty"`
	ESP       string `json:"ESP"`
	ENG       string `json:"ENG"`
	POR       string `json:"POR"`
}

func createDB() error {

	db, err := sql.Open("mysql", DBConfig.user+":"+DBConfig.password+"@tcp("+DBConfig.ip+":"+DBConfig.port+")/")
	if err != nil {
		fmt.Println("1")
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS disenadorlogodb_uploads")
	if err != nil {
		fmt.Println("2")
		return err
	}

	_, err = db.Exec("USE disenadorlogodb_uploads")
	if err != nil {
		fmt.Println("3")
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS icons_uploads ( id  int NOT NULL PRIMARY KEY auto_increment, nounId int NOT NULL UNIQUE, disenadorId int NOT NULL UNIQUE)")
	if err != nil {
		return err
	}

	return nil

}

/*LeerImagen abre realiza una peticion Get por una imagen, la transforma en string y devuelve 1 objeto icono */
func leerImagen(descargas *chan *Icon, wg *sync.WaitGroup, url string, tags []string, category string, nounID int) {

	defer wg.Done()

	resImg, err := client.Get(url)

	if err != nil {

		fmt.Println(err.Error())

		return
	}
	defer resImg.Body.Close()

	if resImg.StatusCode == http.StatusOK {

		imgBytes, errBody := ioutil.ReadAll(resImg.Body)

		if errBody != nil {
			fmt.Println(errBody.Error())
			return
		}

		svgParsed, color := svgParser(string(imgBytes))

		svgBase64 := base64.StdEncoding.EncodeToString([]byte(svgParsed))

		iconoDescargado := Icon{svg: svgBase64, tags: tags, category: category, color: color, nounID: nounID}

		*descargas <- &iconoDescargado

	}

}

func svgParser(svg string) (svgParsed string, color string) {

	color = colors[rand.Intn(len(colors))]

	svgParsed = `<svg fill="` + color + `" ` + strings.SplitN(svg, "<svg", 2)[1]

	return

}

func findMetas(source *map[string]TagJSON, tagEng string) *TagJSON {

	tags := *source

	tag, exist := tags[tagEng]

	if !exist {
		return nil
	}

	return &tag

}

func (i *Icon) insert() error {

	dbUploads, err := sql.Open("mysql", DBConfig.user+":"+DBConfig.password+"@tcp("+DBConfig.ip+":"+DBConfig.port+")/disenadorlogodb_uploads")

	if err != nil {
		return err
	}
	defer dbUploads.Close()

	stmt, err := dbUploads.Prepare("SELECT disenadorId FROM icons_uploads WHERE nounId = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	var id int

	err = stmt.QueryRow(i.nounID).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {

			dbDisenador, err := sql.Open("mysql", DBConfig.user+":"+DBConfig.password+"@tcp("+DBConfig.ip+":"+DBConfig.port+")/"+DBConfig.name)

			if err != nil {
				return err
			}
			defer dbDisenador.Close()

			stmtInsert, err := dbDisenador.Prepare("INSERT INTO elementos(nombre, svg, color, tipo, comprado, categorias_idCategoria) VALUES ('Icono', ?, ?, 'ICONO', 0, ?)")

			if err != nil {
				return err
			}

			defer stmtInsert.Close()

			res, err := stmtInsert.Exec(i.svg, i.color, i.category)

			if err != nil {
				fmt.Println("error insertando icono: " + err.Error())
				return err
			}

			lastID, err := res.LastInsertId()

			if err != nil {
				return err
			}

			i.ID = int(lastID)

			stmtInsertUpload, err := dbUploads.Prepare("INSERT INTO icons_uploads(nounId,disenadorId) VALUES (?,?)")

			if err != nil {
				fmt.Println("error insertando icon_upload: " + err.Error())
				return err
			}

			defer stmtInsertUpload.Close()

			_, err = stmtInsertUpload.Exec(i.nounID, i.ID)

			if err != nil {
				return err
			}

			return nil

		}

		return err

	}

	fmt.Println("Icono Insertado!")
	i.ID = id

	return nil

}

func (i *Icon) insertTag(tagsJSON *map[string]TagJSON, TCPconn *net.Conn) ([]string, error) {

	var iconTagsJSON []TagJSON

	for _, tag := range i.tags {
		tagMetas := findMetas(tagsJSON, tag)
		if tagMetas == nil {
			continue
		}
		tagMetas.Categoria = ""
		iconTagsJSON = append(iconTagsJSON, *tagMetas)
	}

	var ids []string

	request := struct {
		ID     int64  `json:"id"`
		Action string `json:"action"`
		Data   struct {
			Etiquetas []TagJSON `json:"etiquetas"`
		} `json:"data"`
	}{
		ID:     time.Now().Unix(),
		Action: "guardar",
	}

	request.Data.Etiquetas = iconTagsJSON

	requestJSON, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(requestJSON))
	fmt.Fprintf(*TCPconn, string(requestJSON))

	// listen for reply
	resMongo, err := bufio.NewReader(*TCPconn).ReadString('\n')

	if err != nil {
		return nil, err
	}

	var resMongoJSON struct {
		Status float64  `json:"status"`
		Action string   `json:"action"`
		ID     float64  `json:"id"`
		Data   []string `json:"data"`
		Error  string   `json:"error"`
	}

	errMongo := json.Unmarshal([]byte(resMongo), &resMongoJSON)

	if errMongo != nil {

		return nil, errMongo

	}

	if int64(resMongoJSON.ID) != request.ID {
		errIDMongoJSON := errors.New("ID de la respuesta y peticion no coincide")
		return nil, errIDMongoJSON
	}

	if int(resMongoJSON.Status) != 200 {

		errStatusMongoJSON := errors.New(resMongoJSON.Error)
		return nil, errStatusMongoJSON
	}

	ids = resMongoJSON.Data
	fmt.Println("Tags Insertadas!")
	return ids, nil
}

func insertRelTag(Icon *Icon, tagIDs *[]string, TCPconn *net.Conn) error {

	request := struct {
		ID   int64
		Data struct {
			IDs    []string `json:"_ids"`
			Iconos []int    `json:"iconos"`
		} `json:"data"`
	}{
		ID: time.Now().Unix(),
	}

	request.Data.IDs = *tagIDs
	request.Data.Iconos = []int{Icon.ID}

	requestJSON, error := json.Marshal(request)

	if error != nil {
		return error
	}

	fmt.Fprintf(*TCPconn, string(requestJSON))

	// listen for reply
	resMongo, err := bufio.NewReader(*TCPconn).ReadString('\n')

	if err != nil {
		return err
	}

	var resMongoJSON struct {
		Status float64  `json:"status"`
		Action string   `json:"action"`
		ID     float64  `json:"id"`
		Data   []string `json:"data"`
		Error  string   `json:"error"`
	}

	errMongo := json.Unmarshal([]byte(resMongo), &resMongoJSON)

	if errMongo != nil {

		return errMongo

	}

	if int64(resMongoJSON.ID) != request.ID {
		errIDMongoJSON := errors.New("ID de la respuesta y peticion no coincide")
		return errIDMongoJSON
	}

	if int(resMongoJSON.Status) != 200 {

		errStatusMongoJSON := errors.New(resMongoJSON.Error)
		return errStatusMongoJSON
	}
	fmt.Println("Relacion Creada!")
	return nil

}

func main() {
	//creamos un grupo para sincronizar todas las rutinas
	var wg sync.WaitGroup

	//creamos un map contenedor de los registros JSON de iconos
	IconsJSON := make(map[string]IconJSON)

	//Abrimos el archivo producido de la busqueda de NODEJS
	iconsFile, err := os.Open("../find/icons.json")
	if err != nil {
		panic(err.Error())
	}

	//decodificamos el json
	jsonParser := json.NewDecoder(iconsFile)

	err = jsonParser.Decode(&IconsJSON)

	if err != nil {

		panic(err.Error())

	}

	//contamos la cantidad de elementos, creamos un canal de comunicacion entre todas las rutinas y sincronizamos
	cantidad := len(IconsJSON)
	descargas := make(chan *Icon, cantidad)
	wg.Add(cantidad)

	//Iniciamos las descargas asincronicamente
	fmt.Println("Iniciando descargas!")

	for nounID, icon := range IconsJSON {
		id, err := strconv.Atoi(nounID)

		if err != nil {
			continue
		}

		go leerImagen(&descargas, &wg, icon.URL, icon.Tags, icon.Category, id)
	}

	fmt.Println("Descargas en proceso...")

	//esperamos a que todas las descargas se completen
	wg.Wait()

	//cerramos el canal
	close(descargas)

	fmt.Println("Proceso de descarga terminó!")

	//creamos un contenedor para el JSON de tags
	TagsJSON := make(map[string]TagJSON)

	//Abrimos el archivo JSOn de tagas
	tagsFile, err := os.Open("../find/tags.json")
	if err != nil {
		panic(err.Error())
	}

	//Decodificamos el archivo
	jsonParser = json.NewDecoder(tagsFile)
	err = jsonParser.Decode(&TagsJSON)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Verificando DB intermediaria")
	//Revisamos si la BD intermedia de registros de descargas esta creada, si no lo está, la creamos
	err = createDB()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Cargando logos y etiquetas a DB final")

	//Conectandose al servidor para cargar Logos
	connTCP, err := net.Dial("tcp", TCPconfig.ip+":"+TCPconfig.port)

	if err != nil {
		panic(err.Error())
	}

	//Iteramos sobre las descargas para obtener cada icono y su informacion
	for iconDescargado := range descargas {

		err = iconDescargado.insert()

		if err != nil {

			fmt.Println("Error insertando logo:" + err.Error())
			continue
		}

		idsTags, err := iconDescargado.insertTag(&TagsJSON, &connTCP)

		if err != nil {

			fmt.Println("Error insertando tags:" + err.Error())
			continue
		}

		errRel := insertRelTag(iconDescargado, &idsTags, &connTCP)

		if errRel != nil {
			fmt.Println("Error relacionando tags:" + errRel.Error())
			continue
		}

	}

	fmt.Println("Carga Finalizada!")

}
