package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var gSettings struct {
	DBDir     string `json:"db_directory"`
	Port      int    `json:"port"`
	Responses struct {
		PostStatusCode     int    `json:"post_status_code"`
		PostMessageName    string `json:"post_message_name"`
		PostMessageValue   string `json:"post_message_value"`
		PutStatusCode      int    `json:"put_status_code"`
		PutMessageName     string `json:"put_message_name"`
		PutMessageValue    string `json:"put_message_value"`
		DeleteStatusCode   int    `json:"delete_status_code"`
		DeleteMessageName  string `json:"delete_message_name"`
		DeleteMessageValue string `json:"delete_message_value"`
	} `json:"responses"`
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "s", "./db_settings.json", "Filename for Database Settings Json File.")
	flag.Parse()
	err := readSettings(fileName)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	routingRequests(gSettings.Port)
}

func readSettings(fileName string) (err error) {
	jsonDef := []byte(`{"db_directory":"./data","port":5134,"responses":{"post_status_code":405,"post_message_name":"error","post_message_value":"request forbidden","put_status_code":405,"put_message_name":"error","put_message_value":"request forbidden","delete_status_code":405,"delete_message_name":"error","delete_message_value":"request forbidden"}}`)
	rawData, err := ioutil.ReadFile(fileName)
	if err != nil {
		json.Unmarshal(jsonDef, &gSettings)
		return fmt.Errorf("file error")
	}
	err = json.Unmarshal(rawData, &gSettings)
	if err != nil {
		json.Unmarshal(jsonDef, &gSettings)
		return fmt.Errorf("format error")
	}
	return nil
}

func readDatabaseJSON(DBName string) (result interface{}, err error) {
	jsonDef := []byte(`{"null":"null"}`)
	rawData, err := ioutil.ReadFile(gSettings.DBDir + "/" + DBName + ".json")
	if err != nil {
		json.Unmarshal(jsonDef, &result)
		return result, fmt.Errorf("database notfound")
	}
	err = json.Unmarshal(rawData, &result)
	if err != nil {
		json.Unmarshal(jsonDef, &result)
		return result, fmt.Errorf("database format error")
	}
	return result, nil
}

func routingRequests(port int) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", rootRequestFunction)
	router.HandleFunc("/{request}", mainMultiRequestFunction).Methods("GET")
	router.HandleFunc("/{request}/{id}", mainSingleRequestFunction).Methods("GET")
	router.HandleFunc("/{request}", requestForbiddenPost).Methods("POST")
	router.HandleFunc("/{request}/{id}", requestForbiddenPost).Methods("POST")
	router.HandleFunc("/{request}", requestForbiddenPut).Methods("PUT")
	router.HandleFunc("/{request}/{id}", requestForbiddenPut).Methods("PUT")
	router.HandleFunc("/{request}", requestForbiddenDelete).Methods("DELETE")
	router.HandleFunc("/{request}/{id}", requestForbiddenDelete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func respondError(w http.ResponseWriter, statusCode int, errorMessage string) {
	respondJSON(w, statusCode, map[string]string{"error": errorMessage})
}

func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
}

func rootRequestFunction(w http.ResponseWriter, r *http.Request) {
	respondError(w, http.StatusForbidden, "request forbidden")
}

func mainMultiRequestFunction(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)["request"]
	json, err := readDatabaseJSON(request)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondJSON(w, http.StatusOK, json)
	}
}

func mainSingleRequestFunction(w http.ResponseWriter, r *http.Request) {
	//request := mux.Vars(r)["request"]
	//id := mux.Vars(r)["id"]
	respondJSON(w, http.StatusOK, map[string]string{"status": "めんどくさいのでまだ未実装"})
}

func requestForbiddenPost(w http.ResponseWriter, r *http.Request) {
	statusCode := gSettings.Responses.PostStatusCode
	messageName := gSettings.Responses.PostMessageName
	messageValue := gSettings.Responses.PostMessageValue
	respondJSON(w, statusCode, map[string]string{messageName: messageValue})
}

func requestForbiddenPut(w http.ResponseWriter, r *http.Request) {
	statusCode := gSettings.Responses.PutStatusCode
	messageName := gSettings.Responses.PutMessageName
	messageValue := gSettings.Responses.PutMessageValue
	respondJSON(w, statusCode, map[string]string{messageName: messageValue})
}

func requestForbiddenDelete(w http.ResponseWriter, r *http.Request) {
	statusCode := gSettings.Responses.DeleteStatusCode
	messageName := gSettings.Responses.DeleteMessageName
	messageValue := gSettings.Responses.DeleteMessageValue
	respondJSON(w, statusCode, map[string]string{messageName: messageValue})
}
