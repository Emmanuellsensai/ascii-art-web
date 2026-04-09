package main

import (
	"ascii-art-web/ascii"
	"html/template"
	"log"
	"net/http"
)

// PageData holds what we send to the HTML template
type PageData struct {
	Result string
	Error  string
	Text   string
	Banner string
}

// templates are loaded once at startup
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// homeHandler handles GET / — shows the main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// reject any path other than "/"
	if r.URL.Path != "/" {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	// only allow GET on this route
	if r.Method != http.MethodGet {
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}

	// render the page with no result yet
	err := tmpl.Execute(w, PageData{})
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
	}
}

// asciiArtHandler handles POST /ascii-art — processes the form
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// only allow POST on this route
	if r.Method != http.MethodPost {
		http.Error(w, "400 - Bad Request", http.StatusBadRequest)
		return
	}

	// read form values — these match the name="" in the HTML form
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// validate — both fields must be present
	if text == "" || banner == "" {
		http.Error(w, "400 - Bad Request: text and banner are required", http.StatusBadRequest)
		return
	}

	// validate banner is one of the three allowed values
	allowed := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
	if !allowed[banner] {
		http.Error(w, "400 - Bad Request: invalid banner", http.StatusBadRequest)
		return
	}

	// load the banner file
	bannerLines, err := ascii.ReadBanner("banners/" + banner + ".txt")
	if err != nil {
		http.Error(w, "404 - Banner Not Found", http.StatusNotFound)
		return
	}

	// build the lookup table and render the art
	asciiMap := ascii.BuildAsciiMap(bannerLines)
	result := ascii.PrintAscii(text, asciiMap)

	// send the result back to the page
	tmpl.Execute(w, PageData{
		Result: result,
		Text:   text,
		Banner: banner,
	})

}

func main() {
	mux := http.NewServeMux()
	// register our two routes
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("/ascii-art", asciiArtHandler)

	// start the server
	log.Println("Server running at http://localhost:8080\nCheck your browser on this port")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
