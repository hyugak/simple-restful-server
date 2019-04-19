package main

import (
    . "./db"
    . "./model"
	"fmt"
    "net/http"
    "encoding/json"
    "strconv"
    "io/ioutil"
	_ "github.com/lib/pq"
)

func userIndex(w http.ResponseWriter) {
    db := Connect()
    defer db.Close()

	users := Index(db)

    json, err := json.Marshal(&users)
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "")
    }
    response := string(json)

    handleResponse(w, http.StatusOK, response)
}

func userShow(w http.ResponseWriter, id int) {
    db := Connect()
    defer db.Close()

	user := Show(db, id)
    if user.Id == 0 {
        handleStatusCode(w, http.StatusNotFound)
    } else {
        response, err := user.ToJson()
        if err != nil {
            handleStatusCode(w, http.StatusInternalServerError)
        } else {
            handleResponse(w, http.StatusOK, response)
        }
    }
}

func userPost(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    user := User{}
    err = json.Unmarshal(body, &user)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
    }

    db := Connect()
    defer db.Close()

    err = user.Create(db)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
    } else {
        response, err := user.ToJson()
        if err != nil { handleStatusCode(w, http.StatusInternalServerError) }
        handleResponse(w, http.StatusCreated, response)
    }
}

func userPatch(w http.ResponseWriter, r *http.Request, id int) {
    body, err := ioutil.ReadAll(r.Body)
    user := User{}
    err = json.Unmarshal(body, &user)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
    }
    user.Id = id

    db := Connect()
    defer db.Close()

    err = user.Update(db)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
    } else {
        response, err := user.ToJson()
        if err != nil { handleStatusCode(w, http.StatusInternalServerError) }
        handleResponse(w, http.StatusOK, response)
    }
}

func userDestroy(w http.ResponseWriter, id int) {
    user := User{}
    user.Id = id

    db := Connect()
    defer db.Close()

    err := user.Delete(db)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
    } else {
        handleStatusCode(w, http.StatusNoContent)
    }
}

func root(w http.ResponseWriter, r *http.Request) {
    path := string(r.URL.Path)
    if path != "/" {
        handleStatusCode(w, http.StatusNotFound)
        return
    }
    testData := Test{Message: "Hello World"}
    json, err := json.Marshal(&testData)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
    }
    response := string(json)

    handleResponse(w, http.StatusOK, response)
}

func usersHandlerRoot(w http.ResponseWriter, r *http.Request) {
    method := string(r.Method)
    switch method {
    case "GET": userIndex(w)
    case "POST": userPost(w, r)
    default: handleStatusCode(w, http.StatusNotFound)
    }
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    path := string(r.URL.Path)[len("/users/"):]
    method := string(r.Method)
    id, err := strconv.Atoi(path)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
        return
    }

    switch method {
    case "GET": userShow(w, id)
    case "PUT": userPatch(w, r, id)
    case "DELETE": userDestroy(w, id)
    default: handleStatusCode(w, http.StatusNotFound)
    }
}

func handleStatusCode(w http.ResponseWriter, status int) {
    w.WriteHeader(status)
    fmt.Fprintf(w, "")
}

func handleResponse(w http.ResponseWriter, status int, response string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    fmt.Fprintf(w, response)
}

func main() {
    http.HandleFunc("/", root)
    http.HandleFunc("/users", usersHandlerRoot)
    http.HandleFunc("/users/", usersHandler)

    http.ListenAndServe(":8081", nil)
}
