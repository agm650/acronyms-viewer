package webserver

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apex/log"
)

/*
Copyright © 2022 DANDOY Luc

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

var staticContent *embed.FS

func handleHello(w http.ResponseWriter, r *http.Request) {
	ctx := log.WithFields(log.Fields{
		"file":     "webserver/main.go",
		"function": "handleHello",
	})
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	ctx.Warn("Printting Hello")

	fmt.Fprintf(w, "Hello!")
}

func handleAddition(w http.ResponseWriter, r *http.Request) {
	ctx := log.WithFields(log.Fields{
		"file":     "webserver/main.go",
		"function": "handleAddition",
	})
	if r.URL.Path != "/add" {
		ctx.Warn("404 not found.")
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		ctx.Warnf("Method %s is not supported.", r.Method)
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error Parsing Form.", http.StatusInternalServerError)
		ctx.Warnf("Error: %s", err.Error())
		return
	}

	name := r.FormValue("name")
	definition := r.FormValue("definition")
	contributor := r.FormValue("contrib")

	ctx.Warnf("New accronym %s (user: %s):\t %s", name, contributor, definition)
	id, err := addDefinition(name, definition, contributor)
	if err != nil {
		http.Error(w, "Error Adding definition in DB.", http.StatusInternalServerError)
		ctx.Warnf("Error: %s", err.Error())
	}

	ctx.Warnf("Definition Ajouté pour %s (id: %d)", name, id)
	fmt.Fprintf(w, "Definition Ajouté pour %s (id: %d)", name, id)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	ctx := log.WithFields(log.Fields{
		"file":     "webserver/main.go",
		"function": "handleSearch",
	})
	if r.URL.Path != "/search" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error Parsing Form.", http.StatusInternalServerError)
		return
	}

	// AcroDB.
	acro2search := r.FormValue("name")
	// fmt.Fprintf(w, "Searching for accronyms %s\n", acro2search)
	ctx.Warnf("Searching for %s", acro2search)
	res, err := getDefinition(acro2search)
	if err != nil {
		fmt.Fprintf(w, "Definition introuvable %s", acro2search)
		ctx.Warnf("Error: %s", err.Error())
		return
	}

	// fmt.Fprintf(w, "Definition de %s : \n\t%s", acro2search, res)
	// Build a json response
	type jsonReply struct {
		Name        string   `json:"name"`
		Definitions []string `json:"definitions"`
	}

	var reply jsonReply
	reply.Name = acro2search
	for i := 0; i < len(res); i++ {
		reply.Definitions = append(reply.Definitions, res[i])
	}
	jsonData, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, "Error Generating JSON.", http.StatusInternalServerError)
		ctx.Warnf("Error: %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonData))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		data, _ := staticContent.ReadFile("embed/index.html")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, string(data))
		return
	}

	if r.URL.Path == "/acronyms.css" {
		data, _ := staticContent.ReadFile("embed/acronyms.css")
		w.Header().Set("Content-Type", "text/css")
		fmt.Fprint(w, string(data))
		return
	}

	if r.URL.Path == "/acronyms.js" {
		data, _ := staticContent.ReadFile("embed/acronyms.js")
		w.Header().Set("Content-Type", "text/javascript")
		fmt.Fprint(w, string(data))
		return
	}
}

func StartServer(port int, webContent *embed.FS, dbPath string) error {
	ctx := log.WithFields(log.Fields{
		"file":     "webserver/main.go",
		"function": "handleSearch",
	})
	staticContent = webContent
	fmt.Printf("Starting server at port %d\n", port)
	portString := strconv.Itoa(port)

	err := loadDB(dbPath)
	if err != nil {
		ctx.Fatalf("Error loading the DB: %s", err.Error())
	}
	defer closeDB()

	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/add", handleAddition)
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/acronyms.css", handleMain)
	http.HandleFunc("/acronyms.js", handleMain)
	http.HandleFunc("/index.html", handleMain)

	if err := http.ListenAndServe(":"+portString, nil); err != nil {
		ctx.Fatal(err.Error())
	}
	return nil
}
