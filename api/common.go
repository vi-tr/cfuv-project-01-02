package main

import (
    "context"
    "encoding/hex"
    "net/http"
    "net/url"
    "os"
    "regexp"
    "strconv"
    "strings"
    sha "crypto/sha512"
)

type DBProvider[Schema any] interface {
    Find  (c context.Context, query any) (result []Schema, ε ε)
    Add   (c context.Context, addition Schema) (ε ε)
    Edit  (c context.Context, query any, change Schema) (ε ε)
    Remove(c context.Context, query any) (ε ε)
}

func Hash(orig uint64, salt string) string {
    h := sha.Sum512_256([]byte(strconv.FormatUint(orig,10)+salt+strings.TrimSpace(os.Getenv("HASH_SALT"))))
    return hex.EncodeToString(h[:])
}

type User struct {
    GitHubID        uint64          `json:"github_id,omitempty"`
    TelegramIDs     []uint64        `json:"telegram_id,omitempty"`
    Roles           map[string]bool `json:"roles,omitempty"`
    Name            string          `json:"name,omitempty"`
    Group           string          `json:"group,omitempty"`
}

var ValidGroup *regexp.Regexp = regexp.MustCompile(`[а-яА-Я]+-[б]-[оз]-\d+(\(\d+\))?`)
var ValidKind  *regexp.Regexp = regexp.MustCompile(`(ПР|ЛЕ|СМ)`)

type Lesson struct {
    Group   string `json:"group,omitempty"`
    Day     uint64 `json:"day,omitempty"`
    Period  uint64 `json:"period,omitempty"`
    Title   string `json:"title,omitempty"`
    Speaker string `json:"speaker,omitempty"`
    Kind    string `json:"kind,omitempty"`
    Note    string `json:"note,omitempty"`
}


type API struct {
    Client  http.Client
    BaseURL string
    Params  map[string]string
    Headers map[string][]string
}

func (x *API) Call(endpoint []string,
                   extraparams map[string]string,
                   method string, body string,
                   callback ...func(x *API, r *http.Response)) (ε ε) {
    defer ə(&ε)

    q := make(url.Values)
    for k, v := range x.Params    { q.Add(k, v) }
    for k, v := range extraparams { q.Add(k, v) }

    cleanURL := P2(url.JoinPath(x.BaseURL, endpoint...))
    req := P2(http.NewRequest(method, cleanURL + "?" + q.Encode(), strings.NewReader(body)))
    for k, strings := range x.Headers { 
        for _, v := range strings { req.Header.Add(k, v) }
    }
    resp := P2(x.Client.Do(req))

    if len(callback) > 0 { callback[0](x, resp) }
    return
}
