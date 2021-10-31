package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type PhaseTemplateData struct {
    Date         string `json:"date"`
    Phase        string `json:"phase"`
    Illumination string `json:"Illumination"`
    Age          string `json:"Age"`
    Distance     string `json:"Distance"`
    Diameter     string `json:"Diameter"`
    SunDistance  string `json:"SunDistance"`
    SunDiameter  string `json:"SunDiameter"`
    NextNewMoon  string `json:"NextNewMoon"`
    NextFullMoon string `json:"NextFullMoon"`
    ZodiacSign   string `json:"ZodiacSign"`
}

const (
    dateLayout = "January 2, 2006"
    timeLayout = "January 2, 2006 at 12:04"
)

func PrepareData(input MoonDatePhase) PhaseTemplateData {
    output := PhaseTemplateData{Phase: input.Phase, ZodiacSign: input.ZodiacSign}
    output.Date = input.Date.Format(dateLayout)
    output.Illumination = fmt.Sprintf("%d%% of the surface", int(input.Illumination * 100))
    output.Age = fmt.Sprintf("%.1f days", input.Age)
    output.Distance = fmt.Sprintf("%d miles (%d km)", 
        int(input.Distance *  0.6214), int(input.Distance))
    output.Diameter = fmt.Sprintf("%.2f", input.Diameter)
    output.SunDistance = fmt.Sprintf("%d miles (%d km)", 
        int(input.SunDistance *  0.6214), int(input.SunDistance))
    output.SunDiameter = fmt.Sprintf("%.2f", input.SunDiameter)
    output.NextNewMoon = input.NextNewMoon.Format(timeLayout)
    output.NextFullMoon = input.NextFullMoon.Format((timeLayout))

    return output
}

func get_minimal_template() string {
    minimal_template := `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Moon Phase for {{.Date}}</title>
    <meta name="description" content="A simple moon phase website.">
    <meta name="author" content="Joel York">
    <meta property="og:title" content="Moon Phase">
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://www.sitepoint.com/a-basic-html5-template/">
    <meta property="og:description" content="moon phase">
</head>
<body>
    <h1>Moon Phase for {{.Date}}</h1>
    <b>{{.Phase}}</b>
    <table>
        <tr><td>Illumination:</td><td>{{.Illumination}}</td></tr>
        <tr><td>Next New Moon:</td><td>{{.NextNewMoon}}</td></tr>
        <tr><td>Next Full Moon:</td><td>{{.NextFullMoon}}</td></tr>
        <tr><td>Age:</td><td>{{.Age}}</td></tr>
        <tr><td>Distance:</td><td>{{.Distance}}</td></tr>
        <tr><td>Diameter:</td><td>{{.Diameter}}</td></tr>
        <tr><td>Sun Distance:</td><td>{{.SunDistance}}</td></tr>
        <tr><td>Sun Diameter:</td><td>{{.SunDiameter}}</td></tr>
        <tr><td>ZodiacSign:</td><td>{{.ZodiacSign}}</td></tr>
    </table>
</body>
</html>
`
    return minimal_template
}

func moonHtmlHandler(w http.ResponseWriter, r *http.Request, phase MoonDatePhase) {
    tmpl := template.New("minimal")
    tmpl, err := tmpl.Parse(get_minimal_template())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmplData := PrepareData(phase)
    err = tmpl.Execute(w, tmplData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
