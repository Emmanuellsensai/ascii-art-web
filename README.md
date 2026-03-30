# ASCII Art Web

A web application built with **Go (Golang)** that converts user input text into **ASCII art** using different banner styles.

---

## 📌 Project Purpose

This project is designed to help you understand:

* How to build a **web server in Go**
* How to handle **HTTP requests (GET & POST)**
* How to use **HTML templates**
* How to process user input and generate dynamic output
* How ASCII art works internally

👉 In simple terms:

```
User types text → Server processes it → ASCII art is generated → Displayed in browser
```

---

## 🚀 How the Application Works

```
Browser → HTTP Request → Go Server → ASCII Logic → HTML Template → Browser
```

---

## 📁 Project Structure

```
ascii-art-web/
│
├── ascii/              # ASCII generation logic
├── banners/            # ASCII font files
├── templates/
│   └── index.html      # Frontend UI
│
├── main.go             # Backend logic
├── go.mod
└── README.md
```

---

# 🧠 `main.go` — FULL Line-by-Line Explanation

---

## 📦 Package Declaration

```go
package main
```

* Defines the main package
* Required for executable Go programs

---

## 📥 Imports

```go
import (
    "ascii-art-web/ascii"
    "html/template"
    "log"
    "net/http"
)
```

* `ascii-art-web/ascii` → your custom ASCII logic
* `html/template` → renders HTML pages
* `log` → prints messages to terminal
* `net/http` → handles web server and requests

---

## 📊 Struct Definition

```go
type PageData struct {
    Result string
    Error  string
    Text   string
    Banner string
}
```

This struct is used to send data from Go → HTML.

| Field  | Purpose               |
| ------ | --------------------- |
| Result | ASCII output          |
| Error  | Error message         |
| Text   | Keeps user input      |
| Banner | Keeps selected option |

---

## 📄 Load Template

```go
var tmpl = template.Must(template.ParseFiles("templates/index.html"))
```

* Loads the HTML file once when the app starts
* `ParseFiles` reads the file
* `Must` crashes the app if loading fails

---

## 🚀 Main Function

```go
func main() {
```

Program entry point.

---

### Register Routes

```go
http.HandleFunc("/", homeHandler)
```

* Handles homepage `/`

---

```go
http.HandleFunc("/ascii-art", asciiArtHandler)
```

* Handles form submission

---

### Start Server

```go
log.Println("Server running at http://localhost:8080\nCheck your browser on this port")
```

* Prints server status in terminal

---

```go
log.Fatal(http.ListenAndServe(":8080", nil))
```

* Starts server on port 8080
* Stops program if error occurs

---

# 🏠 Home Handler

```go
func homeHandler(w http.ResponseWriter, r *http.Request)
```

Handles requests to `/`.

---

```go
if r.URL.Path != "/"
```

* Ensures only exact `/` is allowed
* Prevents invalid URLs

---

```go
if r.Method != http.MethodGet
```

* Only allows GET requests

---

```go
err := tmpl.Execute(w, PageData{})
```

* Sends empty page (no result yet)

---

```go
if err != nil
```

* Handles template errors

---

# 🎯 ASCII Handler

```go
func asciiArtHandler(w http.ResponseWriter, r *http.Request)
```

Handles form submission.

---

```go
if r.Method != http.MethodPost
```

* Only allows POST requests

---

```go
text := r.FormValue("text")
banner := r.FormValue("banner")
```

* Retrieves user input from form

---

```go
if text == "" || banner == ""
```

* Ensures user entered data

---

```go
allowed := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
```

* Defines valid banner styles

---

```go
if !allowed[banner]
```

* Prevents invalid banner input

---

```go
bannerLines, err := ascii.ReadBanner("banners/" + banner + ".txt")
```

* Loads ASCII font file

---

```go
asciiMap := ascii.BuildAsciiMap(bannerLines)
```

* Converts file into character map

---

```go
result := ascii.PrintAscii(text, asciiMap)
```

* Generates ASCII output

---

```go
tmpl.Execute(w, PageData{
    Result: result,
    Text:   text,
    Banner: banner,
})
```

* Sends data back to HTML
* Preserves user input (important fix you implemented)

---

```go
if err != nil
```

⚠️ Minor issue:

* This `err` refers to the earlier variable
* You should capture error from `Execute` again for accuracy

---

# 🌐 `index.html` — FULL Line-by-Line Explanation

---

## 📄 Document Setup

```html
<!DOCTYPE html>
<html lang="en">
```

* Defines HTML5 document
* Sets language

---

## 🧠 Head Section

```html
<meta charset="UTF-8">
```

* Supports all characters

```html
<title>ASCII Art Generator</title>
```

* Page title

---


```html
<body style="text-align: center;">
```

---

## 🏷️ Title

```html
<h1>ASCII Art Generator</h1>
```

* Displays page heading

---

## 📝 Form

```html
<form action="/ascii-art" method="POST">
```

* Sends data to server

---

## ✏️ Text Input

```html
<textarea ...>{{.Text}}</textarea>
```

* User enters text
* `{{.Text}}` preserves input after submission

---

## 🎛️ Banner Selection

```html
<input type="radio" name="banner" value="standard"
{{if eq .Banner "standard"}}checked{{end}}>
```

* Lets user select banner
* Keeps selected option after submit

---

## 📦 Flexbox Layout

```html
<div class="banner" style="display: flex; align-items: center; justify-content: center;">
```

* Centers radio buttons horizontally

---

## 🔘 Button

```html
<button type="submit">Generate</button>
```

* Submits form

---

## ⚠️ Minor HTML Issue

```html
</div>
```

* Extra closing div (not opened properly)

---

## 📤 Result Section

```html
{{if .Result}}
```

* Only shows result if available

---

```html
<pre>{{.Result}}</pre>
```

* Displays ASCII art
* `<pre>` preserves spacing (very important)

---

## ❌ Error Section

```html
{{if .Error}}
```

* Displays errors

---


```html
<div class="result" style="text-align: justify;">
```

---

# 🔤 ASCII Logic Summary

* Each character = 8 lines
* Stored in banner files
* Converted into map
* Printed line-by-line

---

## ▶️ How to Run

```bash
git clone https://github.com/Emmanuellsensai/ascii-art-web.git
cd ascii-art-web
go run main.go
```

Open:

```
http://localhost:8080
```

---

## 📚 Resources

* Go Docs → https://golang.org/doc/
* net/http → https://pkg.go.dev/net/http
* html/template → https://pkg.go.dev/html/template
* ASCII Art → https://en.wikipedia.org/wiki/ASCII_art

---

## 👨‍💻 Author

**Emmanuel Usang**
https://github.com/Emmanuellsensai

---

## 📄 License

Open-source and for educational purposes.
