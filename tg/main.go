package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "os"
    "bytes"
    "strings"
    "log"
    "reflect"
    "time"
    "strconv"
    "os/signal"
    //"errors"
)

/* Telegram HTTP responses */

type tgUser struct {
    ID          int     `json:"id"`
    First_Name  string  `json:"first_name"`
    Last_Name   string  `json:"last_name"`
    Username    string  `json:"username"`
}

type tgChat struct {
    ID      int `json:"id"`
}

type tgUpdate struct {
    Update_ID    uint64             `json:"update_id"`
    Message      struct {
        Message_ID          int     `json:"message_id"`
        From                tgUser  `json:"from"`
        Date                int     `json:"date"`
        Chat                tgChat  `json:"chat"`
        Text                string  `json:"text"`
        Caption             string  `json:"caption"`
    } `json:"message"`
    Inline_Query struct {
        ID      string  `json:"id"`
        From    tgUser  `json:"from"`
        Query   string  `json:"query"`
        Offset  string  `json:"offset"`
    } `json:"inline_query"`
}

var HandlingID uint64 = 0

func main() {
    
    HandlingID = L2(strconv.ParseUint(strings.TrimSpace(string(L2(os.ReadFile("state")))),10,64))
    // I'm actually appaled that AtExit does not exist in stdlib but you can
    // still make it using your own way.
    // Go feels more like a compile target than a language meant to be used by
    // humans. GPT is probably as efficient at writing go as the best go
    // programmer.
    go func() {
        sigchan := make(chan os.Signal)
        signal.Notify(sigchan, os.Interrupt)
        <-sigchan
        log.Println("Exiting")
        
        fh := F2(os.OpenFile("state",os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644))
        F1(fh.Truncate(0))
        F2(fh.Seek(0, 0))
        F2(fmt.Fprintf(fh, strconv.FormatUint(HandlingID,10)))
        
        os.Exit(0)
    }()

    tgapi := API{
        Client: http.Client{},
        BaseURL: func() string {
            token := strings.TrimSpace(os.Getenv("TELEGRAM_TOKEN"))
            if token=="" { log.Fatal("TELEGRAM_TOKEN environment variable unset or empty, can't proceed.") }
            return "https://api.telegram.org/bot" + token
        }(),
        Headers: map[string][]string {
            "Accept": {"application/json"},
            "Content-Type": {"application/json"},
        },
    }

    for true {
        tgapi.Call([]string{"getUpdates"}, map[string]string{}, "GET", "", updatesHandler)
        time.Sleep(1000 * time.Millisecond)
    }
}

func updatesHandler(x *API, r *http.Response) {
    defer r.Body.Close()
    
    u := struct {
        Ok      bool        `json:"ok"`
        Result  []tgUpdate  `json:"result"`
    }{}
    F1(json.NewDecoder(r.Body).Decode(&u))

    if int(r.StatusCode / 100) != 2 || !u.Ok { log.Fatal(r.Status) }

    for _, u := range u.Result {
        if u.Update_ID <= HandlingID { continue }; HandlingID = u.Update_ID
        if !reflect.ValueOf(u.Message).IsZero() {
            u := u.Message
            L1(commandHandler(x, Or(u.Text,u.Caption), u.From.ID, u.Date, u.Chat.ID))
        }
        if !reflect.ValueOf(u.Inline_Query).IsZero() {
            u := u.Inline_Query
            L1(commandHandler(x, u.Query, u.From.ID, int(time.Now().Unix()), 0))
        }
    }
}

func commandHandler(x *API, text string, userid int, date int, chatid int) error {
    // XXX: Dummy code
    fmt.Println(userid)
    fmt.Println(date)
    fmt.Println(text)
    fmt.Println(chatid)
    x.Call([]string{"sendMessage"}, map[string]string{
        "chat_id": strconv.Itoa(chatid),
        "text": "Testing message",
    }, "POST", "",
        func (x *API, r *http.Response) {
            buf := bytes.NewBuffer(nil)
            buf.ReadFrom(r.Body)
            fmt.Println(string(buf.Bytes()))
            r.Body.Close()
        })

    return nil
}
