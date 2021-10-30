package main

import (
	"html/template"
	"net/http"
)

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
    <h1>Moon Phase for {{.Date}}
    <b>{{.Phase}}
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
    err = tmpl.Execute(w, phase)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
