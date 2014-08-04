package pushgo

import (
	"fmt"
	"net/http"
    "encoding/json"
    "log"
    "time"
)

func init() {
    log.Printf("Started %d", time.Now().Local().Unix());
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        response, _ := json.Marshal(true)
	    fmt.Fprint(w, string(response))
    } )
}

