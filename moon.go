package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/fcgi"
	"os"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gorilla/mux"
)

type MoonDatePhase struct {
	Date         time.Time `json:"date"`
	Phase        string    `json:"phase"`
	Illumination float64   `json:"Illumination"`
	Age          float64   `json:"Age"`
	Distance     float64   `json:"Distance"`
	Diameter     float64   `json:"Diameter"`
	SunDistance  float64   `json:"SunDistance"`
	SunDiameter  float64   `json:"SunDiameter"`
	NextNewMoon  time.Time `json:"NextNewMoon"`
	NextFullMoon time.Time `json:"NextFullMoon"`
	ZodiacSign   string    `json:"ZodiacSign"`
}

type MoonPhaseResponse struct {
	Method    string `json:"method"`
	URI       string `json:"uri"`
	Vars      string `json:"vars"`
	PhaseDate MoonDatePhase
}

var router *mux.Router

const defaultTimeZone = "America/Chicago"

func requestLocation(vars map[string]string) (*time.Location, error) {
	tz := vars["tz"]
	if tz == "" {
		tz = vars["TZ"]
	}
	if tz == "" {
		tz = vars["timezone"]
	}
	if tz == "" {
		tz = defaultTimeZone
	}
	return time.LoadLocation(tz)
}

func hasTimeComponent(str string) bool {
	lower := strings.ToLower(str)
	return strings.Contains(str, ":") ||
		strings.Contains(str, "T") ||
		strings.Contains(lower, "am") ||
		strings.Contains(lower, "pm")
}

func parseDate(str string, loc *time.Location) (time.Time, error) {
	if len(str) == 0 {
		return time.Now().In(loc), nil
	}
	date, err := dateparse.ParseIn(str, loc)
	if err != nil {
		return time.Time{}, err
	}
	if !hasTimeComponent(str) {
		date = time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, loc)
	}
	return date, nil
}

func moonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var results MoonPhaseResponse
	results.Method = r.Method
	results.URI = r.URL.String()
	varStr := ""
	for key, value := range vars {
		varStr += key + "=" + value + ", "
	}
	queryParams := r.URL.Query()
	for key, value := range queryParams {
		vars[key] = strings.Join(value, ",")
		//varStr += key + "=" + strings.Join(value, ",") + ", "
	}
	results.Vars = varStr

	loc, err := requestLocation(vars)
	if err != nil {
		http.Error(w, "invalid timezone", http.StatusBadRequest)
		return
	}

	date, err := parseDate(vars["date"], loc)
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}
	m := New(date)
	phase := MoonDatePhase{Date: date, Phase: m.PhaseName(),
		Illumination: m.Illumination(), Age: m.Age(),
		Distance: m.Distance(), Diameter: m.Diameter(),
		SunDistance: m.SunDistance(), SunDiameter: m.SunDiameter(),
		NextNewMoon:  time.Unix(int64(m.NextNewMoon()), 0).In(loc),
		NextFullMoon: time.Unix(int64(m.NextFullMoon()), 0).In(loc),
		ZodiacSign:   m.ZodiacSign()}

	if _, ok := vars["html"]; ok {
		moonHtmlHandler(w, r, phase)
		return
	}
	// htmlParam := r.URL.Query().Get("html")
	// if htmlParam != "" {
	//     moonHtmlHandler(w, r, phase)
	//     return
	// }

	results.PhaseDate = phase

	b, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
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
