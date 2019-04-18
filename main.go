package main

import (
	"database/sql"
	"fmt"
    "net/http"
    "encoding/json"
    "regexp"
    "strconv"
    "io/ioutil"
	_ "github.com/lib/pq"
)

type Test struct {
    Message string
}

type User struct {
	Id         int
	Name       string
	Email      string
	Created_at string
	Updated_at string
}

func isDigit(s string) bool{
    return regexp.MustCompile(`^([1-9][0-9]*|0)$`).Match([]byte(s))
}

func connect() *sql.DB {
	connStr := "user=dev dbname=dev password=secret sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func index(db *sql.DB) (users []User) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
    defer rows.Close()

	for rows.Next() {
		user := User{}
        err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Updated_at)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return
}

func show(db *sql.DB, id int) (user User) {
    err := db.QueryRow("SELECT * FROM users WHERE id=$1", id).Scan(&user.Id, &user.Name, &user.Email, &user.Created_at, &user.Updated_at)
    if err == sql.ErrNoRows {
        return
    } else if err != nil {
        panic(err)
    }
    return
}

func (user *User) ToJson() (s string, err error) {
    json, err := json.Marshal(&user)
    s = string(json)
    return
}

func (user *User) Create(db *sql.DB) (err error) {
    sql := "insert into users (name, email) values ($1, $2) returning id, created_at, updated_at"
    stmt, err := db.Prepare(sql)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()
    err = stmt.QueryRow(user.Name, user.Email).Scan(&user.Id, &user.Created_at, &user.Updated_at)
    if err != nil {
        fmt.Println(err)
    }
    return
}

func (user *User) Delete(db *sql.DB) (err error) {
    _, err = db.Exec("DELETE FROM users WHERE id=$1", user.Id)
    return
}

func (user *User) Update(db *sql.DB) (err error) {
    sql := "UPDATE users set name=$2, email=$3 WHERE id=$1 returning created_at, updated_at"
    stmt, err := db.Prepare(sql)
    if err != nil {
        fmt.Println(err)
    }
    defer stmt.Close()
    err = stmt.QueryRow(user.Id, user.Name, user.Email).Scan(&user.Created_at, &user.Updated_at)
    if err != nil {
        fmt.Println(err)
    }
    return
}

func userIndex(w http.ResponseWriter) {
    db := connect()
    defer db.Close()

	users := index(db)

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
    db := connect()
    defer db.Close()

	user := show(db, id)
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

    db := connect()
    defer db.Close()

    err = user.Create(db)
    fmt.Println(user)
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

    db := connect()
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

    db := connect()
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
    default: fmt.Println("no match")
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