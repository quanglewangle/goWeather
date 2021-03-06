// Example for a FastCGI program that terminates gracefully.
package main

import (
    "net/http" 
    "net/http/fcgi"
    "html/template"
    "os" 
    "os/signal"
    "db2"
    "runtime"
    "svgo1"
    "syscall" 
    "time"
  
)

type Diagrams struct {
    WindDir     template.HTML
    WindSpeed   template.HTML
    Temperature template.HTML
    Pressure    template.HTML
}

// a simple request handler
func handler(w http.ResponseWriter, r *http.Request) {
    
    db.OpenDatabase()
        
    t, _ := template.ParseFiles("indexTemplate.template")
    t.Execute(w, Diagrams{WindSpeed:template.HTML( svgo1.WindSpeedDial()),
                          WindDir:template.HTML(svgo1.WindDial()),
                          Temperature:template.HTML(svgo1.TempDial()),
                          Pressure:template.HTML(svgo1.PressureDial()),
                          })
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
