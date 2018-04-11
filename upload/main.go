package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

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
	name:     "disenadorlogodb",
	user:     "logoPro",
	password: "&rJ-fZ:1uZ24",
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
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS disenadorlogodb_uploads")
	if err != nil {
		return err
	}

	_, err = db.Exec("USE disenadorlogodb_uploads")
	if err != nil {
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
			/*
				var idCategory int

				stmtCategory, err := dbDisenador.Prepare("SELECT idCategoria FROM categorias WHERE nombreCategoria = ?")
				if err != nil {
					return err
				}

				defer stmtCategory.Close()

				err = stmtCategory.QueryRow(i.category).Scan(&idCategory)

				if err != nil {
					fmt.Println("error Seleccionando la categoria: " + err.Error())
					return err
				}*/

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

	i.ID = id

	return nil

}

func (i *Icon) insertTag(tagsJSON *map[string]TagJSON) ([]string, error) {

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

	request := map[string][]TagJSON{
		"etiquetas": iconTagsJSON,
	}

	requestJSON, error := json.Marshal(request)

	if error != nil {
		return nil, error
	}

	resInsert, err := client.Post("http://127.0.0.1:666/app/etiquetas", "application/json", bytes.NewBuffer(requestJSON))

	if err != nil {

		return nil, err
	}
	defer resInsert.Body.Close()

	if resInsert.StatusCode == http.StatusOK {

		resBytes, errBody := ioutil.ReadAll(resInsert.Body)

		if errBody != nil {

			return nil, errBody
		}

		r := strings.NewReader(string(resBytes))

		err := json.NewDecoder(r).Decode(&ids)

		if err != nil {
			return nil, err
		}

	} else {
		errMongo := errors.New("Error de mongo Insertando tag")
		return nil, errMongo
	}

	return ids, nil
}

func insertRelTag(Icon *Icon, tagIDs *[]string) error {

	request := struct {
		IDs    []string `json:"_ids"`
		Iconos []int    `json:"iconos"`
	}{
		IDs:    *tagIDs,
		Iconos: []int{Icon.ID},
	}

	requestJSON, error := json.Marshal(request)

	if error != nil {
		return error
	}

	resInsert, err := client.Post("http://127.0.0.1:666/app/etiquetas/iconos", "application/json", bytes.NewBuffer(requestJSON))

	if err != nil {
		return err
	}
	defer resInsert.Body.Close()

	if resInsert.StatusCode != http.StatusOK {
		errMongo := errors.New("Error de mongo Insertando relacion de tag-icono")
		return errMongo
	}

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

	//Iteramos sobre las descargas para obtener cada icono y su informacion
	for iconDescargado := range descargas {

		err = iconDescargado.insert()

		if err != nil {

			fmt.Println("Error insertando logo:" + err.Error())
			continue
		}

		idsTags, err := iconDescargado.insertTag(&TagsJSON)

		if err != nil {

			fmt.Println("Error insertando tags:" + err.Error())
			continue
		}

		errRel := insertRelTag(iconDescargado, &idsTags)

		if errRel != nil {
			fmt.Println("Error relacionando tags:" + errRel.Error())
			continue
		}

		//var tagIDs []string

		//iteramos sobre cada tag del icono en Ingles
		/*for _, tag := range iconDescargado.tags {

			//obtenemos sus traducciones
			tagMetas := findMetas(&TagsJSON, tag)
			if tagMetas == nil {
				fmt.Println(err.Error())
				continue
			}

			//insertamos el tag en la base de datos del disenador
			tagID, err := tagMetas.insert()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			tagIDs = append(tagIDs, tagID)

			/
				err = iconDescargado.insertRelTag(tagID)

				if err != nil {

					fmt.Println(err.Error())
				}/

		}*/

	}

	fmt.Println("Carga Finalizada!")

}
