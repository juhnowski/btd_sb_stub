package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ident struct {
	ID         string `json:"ID"`
	FirebaseID string `json:"FirebaseId"`
	SberbankID string `json:"SberbankId"`
	BeruID     string `json:"BeruId"`
	DocDocID   string `json:"DocDocId"`
}

type allIdents []ident

var idents = allIdents{
	{
		ID:         "1",
		FirebaseID: "3xYkQsV3hpOVtWNpQ76GpbmRLxH3",
		SberbankID: "SberbankId-123",
		BeruID:     "BeruId-123",
		DocDocID:   "DocDocId-123",
	},
}

var identCounter = 1

func createIdent(w http.ResponseWriter, r *http.Request) {
	var newIdent ident
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "enter data with the: FirebaseID, SberbankId, BeruId, DocDocId")
	}

	json.Unmarshal(reqBody, &newIdent)
	identCounter = identCounter + 1

	newIdent.ID = strconv.Itoa(identCounter)

	idents = append(idents, newIdent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newIdent)
}

func aliveLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "alive")
}

func getAllIdents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(idents)
}

func deleteIdent(w http.ResponseWriter, r *http.Request) {
	identID := mux.Vars(r)["id"]

	for i, singleIdent := range idents {
		if singleIdent.ID == identID {
			idents = append(idents[:i], idents[i+1:]...)
			fmt.Fprintf(w, "The ident with ID %v has been deleted successfully", identID)
		}
	}
}

func updateIdent(w http.ResponseWriter, r *http.Request) {
	identID := mux.Vars(r)["id"]
	var updatedIdent ident

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the idents field: FirebaseID, SberbankId, BeruId, DocDocId to update")
	}
	json.Unmarshal(reqBody, &updatedIdent)

	for i, singleIdent := range idents {
		if singleIdent.ID == identID {
			singleIdent.FirebaseID = updatedIdent.FirebaseID
			singleIdent.SberbankID = updatedIdent.SberbankID
			singleIdent.BeruID = updatedIdent.BeruID
			singleIdent.DocDocID = updatedIdent.DocDocID
			idents = append(idents[:i], singleIdent)
			json.NewEncoder(w).Encode(singleIdent)
		}
	}
}

func getOneIdent(w http.ResponseWriter, r *http.Request) {
	identID := mux.Vars(r)["id"]

	for _, singleIdent := range idents {
		if singleIdent.ID == identID {
			json.NewEncoder(w).Encode(singleIdent)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", aliveLink)
	router.HandleFunc("/ident", createIdent).Methods("POST")
	router.HandleFunc("/idents", getAllIdents).Methods("GET")
	router.HandleFunc("/idents/{id}", getOneIdent).Methods("GET")
	router.HandleFunc("/idents/{id}", updateIdent).Methods("PATCH")
	router.HandleFunc("/idents/{id}", deleteIdent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", router))
}
