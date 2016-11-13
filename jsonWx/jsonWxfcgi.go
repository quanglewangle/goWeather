// Example for a FastCGI program that terminates gracefully.
package main

import (
    "net/http" 
    "net/http/fcgi"
  
 //   "html/template"
    "os" 
    "os/signal"
    "github.com/quanglewangle/goWeather/wx"
    "github.com/quanglewangle/goWeather/db2"
    "runtime"
  //  "github.com/quanglewangle/goWeather/svgo1"
 //   "github.com/quanglewangle/goWeather/zambra"
    "syscall" 
    "time"
    "encoding/json"
  
)
type Wx struct {
  TemperatureC    float64
}

// a simple request handler
func handler(w http.ResponseWriter, r *http.Request) {
    
    db.OpenDatabase()
    
    var wxRecord wx.CurrentWeather
    wxRecord = db.GetCurrentSummary()

    js, err := json.Marshal(wxRecord)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)

}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())     // use all CPU cores
    n := runtime.NumGoroutine() + 1          // initial number of Goroutines

    // install signal handler
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGTERM)

    // Spawn request handler
    go func() {
        err := fcgi.Serve(nil, http.HandlerFunc(handler))
        if err != nil {
            panic(err)
        }
    }()

   
    // catch signal
    _ = <-c

    // give pending requests in fcgi.Serve() some time to enter the request handler
    time.Sleep(time.Millisecond * 100)

    // wait at most 3 seconds for request handlers to finish
    for i := 0; i < 30; i++ {
        if runtime.NumGoroutine() <= n {
            return
        }
        time.Sleep(time.Millisecond * 100)
    }
}
