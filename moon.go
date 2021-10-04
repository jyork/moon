package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/fcgi"
	"os"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gorilla/mux"
	"gopkg.in/ini.v1"
)

type MoonPhaseResponse struct {
    Method       string    `json:"method"`
    URI          string    `json:"uri"`
    Date         time.Time `json:"date"`
    Phase        string    `json:"phase"`
    Illumination float64   `json:"Illumination"`
    Age          float64   `json:"Age"`
    Distance     float64   `json:"Distance"`
    Diameter     float64   `json:"Diameter"`
    SunDistance  float64   `json:"SunDistance"`
    SunDiameter  float64   `json:"SunDiameter"`
    NewMoon      float64   `json:"NewMoon"`
    ZodiacSign   string    `json:"ZodiacSign"`
}

var router *mux.Router

func parseDate(str string) time.Time {
    if len(str) == 0 {
        return time.Now()
    }
    date, err := dateparse.ParseAny(str)
    if err != nil {
        return time.Now()
    }
    return date
}

func moonHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    var results MoonPhaseResponse
    results.Method = r.Method
    results.URI = r.URL.String()

    date := parseDate(vars["date"])
    
    results.Date = date
    m := New(date)
    results.Phase = m.PhaseName()
    results.Illumination = m.Illumination()
    results.Age = m.Age()
    results.Distance = m.Distance()
    results.Diameter = m.Diameter()
    results.SunDistance = m.SunDistance()
    results.SunDiameter = m.SunDiameter()
    results.NewMoon = m.NewMoon()
    results.ZodiacSign = m.ZodiacSign()
  
    b, err := json.MarshalIndent(results, "", "    ")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-type", "application/json")
    io.WriteString(w, string(b))
  
}

func handler(w http.ResponseWriter, r *http.Request) {
    _ = fcgi.ProcessEnv(r)
    log.Print("handler processing "+r.URL.String()+"; request=\n", r)
    router.ServeHTTP(w, r)
    //  r.URL =
}

func main() {
    file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()
    log.SetOutput(file)

    cfg, err := ini.Load("config/moon.ini")
    if err != nil {
        log.Fatalf("Fail to read file: %v", err)
        os.Exit(1)
    }

    log.Printf("test key is %s", cfg.Section("").Key("test").String())
    /*
     *
     *  Everything that is done here should be setup code etc. which is retained between calls
     *
     */
    router = mux.NewRouter()
    router.HandleFunc("/moon.fcgi/moon/{date}", moonHandler)
    router.HandleFunc("/moon.fcgi", moonHandler)
    http.HandleFunc("/", handler)
    //http.HandleFunc("/pokedex/", PokedexHandler)
    // This is what actually concurrently handles requests
    if err := fcgi.Serve(nil, nil); err != nil {
        panic(err)
    }
}
