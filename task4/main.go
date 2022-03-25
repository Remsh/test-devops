package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//gloabl var
var iMap = make(map[string]int)
var port = 9000
var currentName string

type Person struct {
	Name string `json:"name"`
}

func main() {
	fmt.Println("Version: v0.2.\nPlease browse 127.0.0.1:9000/welcome")

	create_table() //create db test04 first
	setupRoutes()

}

func setupRoutes() {
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/", rootHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	address := ":" + strconv.Itoa(port)
	http.ListenAndServe(address, nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/welcome" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		if len(currentName) > 0 {
			w.Write([]byte(fmt.Sprintf("Hello, %s!\n", currentName)))

		} else {
			w.Write([]byte("Hello, Wiredcraft!"))
		}
		iMap[currentName] += 1
	case "POST":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", reqBody)

		var person Person

		if err := json.Unmarshal(reqBody, &person); err != nil {
			log.Fatalf("error %v", err)
			w.Write([]byte("Prase Json Fail\n"))
		} else {

			name := person.Name
			currentName = name
			iMap[currentName] += 1 //Todo, load iMap from db first in real app

			if iMap[currentName] == 1 {
				save2Db(currentName)
			}

			w.Write([]byte("Received a POST request\n"))
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}

}

func save2Db(currentName string) {
	db, err := sql.Open("mysql", "root:example@tcp(192.168.1.109:3306)/test04?tls=skip-verify&autocommit=true")
	if err != nil {
		panic(err.Error()) // Just for example purpose.
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO persons (name) VALUES( ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in real app
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(currentName)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in real app
	}
	return

}

func create_table() {
	connStr := "root:example@tcp(192.168.1.109:3306)/test04?tls=skip-verify&autocommit=true"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}

	const query = `CREATE TABLE IF NOT EXISTS persons (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL)`

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()
}

/*
func drop_table() {
	connStr := "root:example@tcp(192.168.1.109:3306)/test01?tls=skip-verify&autocommit=true"
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS pages")
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()
}
*/
