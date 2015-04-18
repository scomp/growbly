/*
export RETHINKDB_URL="localhost:28015"
export RETHINKDB_AUTHKEY="ca005692-5b32-4516-9dcc-8aee93203fa3"
*/

package main

import (
    "github.com/gocraft/web"
    "net/http"
    "strings"
    "fmt"
    "log"
    "os"
    "time"
    r "github.com/dancannon/gorethink"
)

var session *r.Session
var url, authKey string
type user struct {
        UserName string
        UserPassword string
        UserAccess string
        Email string
    };
type Context struct {
    user *User // Assumes you've defined a User type as well
};
func init() {
    // Needed for wercker. By default url is "localhost:28015"
    url = os.Getenv("RETHINKDB_URL")
    if url == "" {
        url = "localhost:28015"
    }
    // Needed for running tests for RethinkDB with a non-empty authkey
    authKey = os.Getenv("RETHINKDB_AUTHKEY")
}

func main() {
    session, err := r.Connect(r.ConnectOpts{
        Address: url,
        AuthKey: authKey,
        Database: "growbly_dev",
    })
    if err != nil {
        log.Fatalf("Error connecting to DB: %s", err)
    }

    res, err := r.Expr("Hello World").Run(session)
    if err != nil {
        log.Fatalln(err.Error())
    }


    u := user{"test","test1234","god","jemduff@gmail.com"}
    start := time.Now()
    _,err =  r.Table("users").Insert(u).RunWrite(session);
    log.Printf("Insert took %s", time.Since(start));
    if err != nil {
              log.Fatalln(err.Error())
    }
    fmt.Println("did insert");
    var response string
    err = res.One(&response)
    if err != nil {
        log.Fatalln(err.Error())
    }
    ress,_ := r.Table("users").Run(session);
    var rows []interface{}
    _ = ress.All(&rows)

    fmt.Println(rows)
    fmt.Println(response)
            router := web.New(Context{}).                   // Create your router
            Middleware(web.LoggerMiddleware).           // Use some included middleware
            Middleware(web.ShowErrorsMiddleware).       // ...
            Middleware((*Context).IsLoggedIn).       // Your own middleware!
            Get("/", (*Context).SayHello) 
    http.ListenAndServe("localhost:3000", router)
}
