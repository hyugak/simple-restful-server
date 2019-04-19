package model

import (
    "fmt"
    "encoding/json"
    "database/sql"
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

func Index(db *sql.DB) (users []User) {
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

func Show(db *sql.DB, id int) (user User) {
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
