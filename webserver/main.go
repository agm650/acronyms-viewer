package webserver

import (
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
	}

	// AcroDB.
	acro2search := r.FormValue("name")
	fmt.Fprintf(w, "Searching for accronyms %s\n", acro2search)
	ctx.Warnf("Searching for %s", acro2search)
	res, err := getDefinition(acro2search)
	if err != nil {
		fmt.Fprintf(w, "Definition introuvable %s", acro2search)
		ctx.Warnf("Error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "Definition de %s : \n\t%s", acro2search, res)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Main page")
}

func StartServer(port int, dbPath string) error {
	ctx := log.WithFields(log.Fields{
		"file":     "webserver/main.go",
		"function": "handleSearch",
	})

	fmt.Printf("Starting server at port %d\n", port)
	portString := strconv.Itoa(port)

	err := loadDB(dbPath)
	if err != nil {
		ctx.Fatalf("Error loading the DB: %s", err.Error())
	}
	defer closeDB()

	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/", handleMain)

	if err := http.ListenAndServe(":"+portString, nil); err != nil {
		ctx.Fatal(err.Error())
	}
	return nil
}
