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

// GET /users を処理する関数
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

// GET /users/:id を処理する関数
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

// POST /users/:id を処理する関数
func userCreate(w http.ResponseWriter, r *http.Request) {
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

// PUT /users/:id を処理する関数
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
        if err != nil {
            handleStatusCode(w, http.StatusInternalServerError)
        }
        handleResponse(w, http.StatusOK, response)
    }
}

// DELETE /users/:id を処理する関数
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

// "/"を処理するハンドラ。
func root(w http.ResponseWriter, r *http.Request) {
    path := string(r.URL.Path)

    // "/" 以外には404を返却
    if path != "/" {
        handleStatusCode(w, http.StatusNotFound)
        return
    }

    // Test構造体を初期化してjsonに整形し、最終的にstringに変換
    testData := Test{Message: "Hello World"}
    json, err := json.Marshal(&testData)
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
    }
    response := string(json)

    // jsonデータ(string)を200とともに返却
    handleResponse(w, http.StatusOK, response)
}

// "/users" を処理するハンドラ。
func usersHandlerRoot(w http.ResponseWriter, r *http.Request) {
    method := string(r.Method)

    // HTTP methodによる処理分岐。どこにもマッチしなければ404を返却
    switch method {
    case "GET": userIndex(w)
    case "POST": userCreate(w, r)
    default: handleStatusCode(w, http.StatusNotFound)
    }
}

// "/users/:id" を処理するハンドラ。HTTP methodによって分岐
func usersHandler(w http.ResponseWriter, r *http.Request) {
    path := string(r.URL.Path)[len("/users/"):]
    method := string(r.Method)

    // :id を特定
    id, err := strconv.Atoi(path)
    // :id に当たる部分がdigitでなければ500を返却
    if err != nil {
        handleStatusCode(w, http.StatusInternalServerError)
        return
    }

    // HTTP methodによる処理分岐。どこにもマッチしなければ404を返却
    switch method {
    case "GET": userShow(w, id)
    case "PUT": userPatch(w, r, id)
    case "DELETE": userDestroy(w, id)
    default: handleStatusCode(w, http.StatusNotFound)
    }
}

// 任意のHTTP statusをデータなしで返す
func handleStatusCode(w http.ResponseWriter, status int) {
    w.WriteHeader(status)
    fmt.Fprintf(w, "")
}

// 任意のHTTP status codeをjsonデータとともに返す
func handleResponse(w http.ResponseWriter, status int, response string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    fmt.Fprintf(w, response)
}

func main() {
/* 
    "/", "/users", "/users/:id"に対してそれぞれハンドラを登録
    それぞれのパスに対するHTTP methodによる処理分岐は各ハンドラで行う
*/
    http.HandleFunc("/", root)
    http.HandleFunc("/users", usersHandlerRoot)
    http.HandleFunc("/users/", usersHandler)

    // TCP/8081番ポートでlisten
    http.ListenAndServe(":8081", nil)
}
