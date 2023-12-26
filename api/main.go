package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var UsersDB DBProvider[User]
var ScheduleDB DBProvider[Lesson]

var BadAuthε      = errors.New("Privileges insufficient")
var NothingFoundε = errors.New("Query yielded no results")
var BadRequestε   = errors.New("Malformed request")

func main() {

    client := P2(mongo.Connect(context.Background(),
        options.Client().ApplyURI(os.Getenv("MONGODB_URI"))))
    defer L1(client.Disconnect(context.Background()))
    dbHandle := client.Database("main")

    UsersDB    = P2(NewMongoDB[User](dbHandle, "users"))
    ScheduleDB = P2(NewMongoDB[Lesson](dbHandle,"schedule"))

    // Behold. API
    http.HandleFunc("/schedule/", schedule)
    http.ListenAndServe(":8080", nil)
}

type privs struct {
    Admin   bool
    Self    bool
    Teacher bool
}

func get_priv(c context.Context, RequesterID uint64, RequestedID uint64) (p privs, ε ε) {
    list := P2(UsersDB.Find(c, User{GitHubID: RequesterID}))
    if len(list) > 1 {
        log.Fatal(fmt.Sprintf("More than one user has the ID: %v",RequesterID))
    }
    return privs{
        Admin: list[0].Roles["admin"],
        Self: RequesterID == RequestedID,
        Teacher: list[0].Roles["teacher"],
    }, nil
}

func respond(w *http.ResponseWriter, body any, ε ε) {
    if errors.Is(ε,        BadAuthε                 ) {
        (*w).WriteHeader(   http.StatusUnauthorized ); L1(ε)
    } else if errors.Is(ε, BadRequestε              ) {
        (*w).WriteHeader(   http.StatusBadRequest   ); L1(ε)
    } else if errors.Is(ε, NothingFoundε            ) {
        (*w).WriteHeader(   http.StatusNoContent    ); L1(ε)
    } else if ε ==         nil                        {
        (*w).WriteHeader(   http.StatusOK           )
    } else { P1(ε) }
    switch body.(type) {
        case string:              fmt.Fprint(*w, body.(string)                 ); break
        case nil:                 fmt.Fprint(*w, "[]"                          ); break
        default:                  fmt.Fprint(*w, string(P2(json.Marshal(body))))
    }
}


func schedule(w http.ResponseWriter, r *http.Request) {
    defer func(){ if r:=recover(); r!=nil {respond(&w, nil, BadRequestε)} }()
    string_paths := strings.Split(strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/schedule"), "/"), "/")
    PA(len(string_paths) <= 3)
    var day, lesson uint64
    group := string_paths[0]
    PA(ValidGroup.MatchString(string_paths[0]))
    if len(string_paths) > 1 {
        day    = P2(strconv.ParseUint(string_paths[1], 10, 64))
    if len(string_paths) > 2 {
        lesson = P2(strconv.ParseUint(string_paths[2], 10, 64))
    }}
    
    var ε ε
    var resp string
    if r.Method == "GET" {
        // Days, lessons, etc. 1-indexed to differentiate the zero-value
        if           day==0 { var l [][]Lesson; l, ε = getWeek  (r.Context(), group);          resp = string(P2(json.Marshal(l)))
        } else if lesson==0 { var l []Lesson; l, ε = getDay   (r.Context(), group, day);       resp = string(P2(json.Marshal(l)))
        } else              { var l Lesson; l, ε = getLesson(r.Context(), group, day, lesson); resp = string(P2(json.Marshal(l)))
        }
    } else if r.Method == "POST" {
        if lesson==0 {respond(&w, nil, BadRequestε); return} 
        ε = modifyLesson(r.Context(), group, day, lesson, r.Body)
    }

    respond(&w, resp, ε)
}

func getWeek(c context.Context, group string) (r [][]Lesson, ε ε) {
    for i := uint64(1); i<=14; i++ {
        t, tε := getDay(c, group, i)
        r = append(r,t); ε = errors.Join(ε, tε)
    }
    return
}

func getDay(c context.Context, group string, day uint64) (r []Lesson, ε ε) {
    for i := uint64(1); i<=8; i++ {
        t, tε := getLesson(c, group, day, i)
        r = append(r,t); ε = errors.Join(ε, tε)
    }
    return
}

func getLesson(c context.Context, group string, day uint64, period uint64) (r Lesson, ε ε) {
    defer ə(&ε)
    list := P2(ScheduleDB.Find(c, Lesson{Group: group, Day: day, Period: period}))
    if len(list) > 1 {
        panic(fmt.Sprintf("Matched more than one lesson for the combination: %v, %v, %v (Zero values with omitempty?)", group, day, period))
    }
    r = list[0]
    return
}

type LessonRequest struct {
    Change Lesson
    RequesterID uint64
}

func modifyLesson(c context.Context, group string, day uint64, period uint64, b io.Reader) (ε ε) {
    defer ə(&ε)
    var body LessonRequest
    P1(json.NewDecoder(b).Decode(&body))
    priv := P2(get_priv(c, body.RequesterID, uint64(0)))
    if priv.Admin {
        ScheduleDB.Edit(c, Lesson{Group: group, Day: day, Period: period}, body.Change)
    } else if priv.Teacher {
        ScheduleDB.Edit(c, Lesson{Group: group, Day: day, Period: period}, Lesson{Note: body.Change.Note})
    } else {
        return BadAuthε
    }
    
    return
}
