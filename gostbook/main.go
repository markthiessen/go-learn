package main

import (
	"labix.org/v2/mgo"
	"html/template"
	"net/http"
)

var index = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/index.html",
))

func hello(w http.ResponseWriter, req *http.Request){

	s := session.Clone()
	defer s.Close()

	coll := s.DB("gostbook").C("entries")
	query := coll.Find(nil).Sort("-timestamp")

	var entries []Entry
	if err := query.All(&entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := index.Execute(w, entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var session *mgo.Session

func main(){
	var err error
	session, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", hello)
	http.HandleFunc("/sign", sign)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}