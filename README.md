# ascii-art-web

A web-based ASCII art generator built in Go. Type any text in the browser, choose a font style, and the server renders it as large ASCII art characters.

---

## Description

This project takes the core ascii-art rendering logic and wraps it in a web server. The user interacts through a webpage — no terminal required. The server handles form submissions, renders the ASCII art using the chosen font, and sends the result back to the browser.

---

## Authors

- Emmanuel ([@emmanuellsensai](https://github.com/emmanuellsensai))

---

## Usage

**1. Clone the project**
```bash
git clone <your-repo-url>
cd ascii-art-web
```

**2. Make sure the banner files are in the banners/ folder**
```
banners/
├── standard.txt
├── shadow.txt
└── thinkertoy.txt
```

**3. Run the server**
```bash
go run .
```

**4. Open your browser and go to**
```
http://localhost:8080
```

**5. Use the page**
- Type your text in the input box
- Pick a banner style (standard, shadow, or thinkertoy)
- Click Generate
- The ASCII art appears below the form

**To use a new line in your art, type `\n` in the text box:**
```
Hello\nWorld
```

---

## Project Structure

```
ascii-art-web/
├── main.go               ← HTTP server, routes, handlers
├── go.mod                ← module declaration
├── ascii/
│   └── render.go         ← core rendering logic (unchanged from ascii-art)
├── templates/
│   └── index.html        ← the webpage the user sees
└── banners/
    ├── standard.txt      ← standard font
    ├── shadow.txt        ← shadow font
    └── thinkertoy.txt    ← thinkertoy font
```

---

## Implementation Details

### How the Server Works

The server uses Go's built-in `net/http` package — no external frameworks needed.

Two routes are registered:

| Route | Method | What it does |
|---|---|---|
| `/` | GET | Serves the main HTML page |
| `/ascii-art` | POST | Receives form data, renders art, returns result |

```go
http.HandleFunc("/", homeHandler)
http.HandleFunc("/ascii-art", asciiArtHandler)
http.ListenAndServe(":8080", nil)
```

**Learn more:** https://gowebexamples.com/hello-world/

---

### How the Template Works

Go's `html/template` package lets us write HTML files with placeholders that Go fills in at runtime.

```go
// load the template once at startup
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// send data to the template
tmpl.Execute(w, PageData{Result: result})
```

In the HTML file, `{{.Result}}` is replaced by the actual ASCII art string:

```html
<pre>{{.Result}}</pre>
```

The `<pre>` tag preserves all spaces and newlines — essential for ASCII art to display correctly.

**Learn more:** https://gowebexamples.com/templates/

---

### How the Form Works

The HTML form sends data to the server when the user clicks Generate:

```html
<form action="/ascii-art" method="POST">
    <textarea name="text"></textarea>
    <input type="radio" name="banner" value="standard">
    <button type="submit">Generate</button>
</form>
```

- `action="/ascii-art"` — where the form sends the data
- `method="POST"` — how it sends it
- `name="text"` — Go reads this with `r.FormValue("text")`
- `name="banner"` — Go reads this with `r.FormValue("banner")`

**Learn more:** https://gowebexamples.com/forms/

---

### How Form Data is Read in Go

```go
text   := r.FormValue("text")
banner := r.FormValue("banner")
```

The `name` attribute in the HTML must match exactly what you pass to `r.FormValue`.

---

### How the ASCII Art is Rendered

The core rendering logic is unchanged from the original ascii-art project:

**Step 1 — ReadBanner:** reads the `.txt` font file and splits it into lines

**Step 2 — BuildAsciiMap:** organises the lines into a lookup table where each character maps to its 8 art rows

**Step 3 — PrintAscii:** loops through the input text, looks up each character, and builds the final art string row by row

```go
bannerLines, err := ascii.ReadBanner("banners/" + banner + ".txt")
asciiMap := ascii.BuildAsciiMap(bannerLines)
result := ascii.PrintAscii(text, asciiMap)
```

**The formula that finds any character in the font file:**
```
position = (ASCII value of character - 32) * 9 + 1 + current row
```

Space is ASCII 32 and is the first character in the file. Every character after it is 9 lines further (8 art rows + 1 blank separator).

---

### HTTP Status Codes

The server returns the correct status code for every situation:

| Situation | Code | How |
|---|---|---|
| Page loads fine | 200 | Automatic |
| Art rendered fine | 200 | Automatic |
| Unknown URL | 404 | `http.Error(w, "...", http.StatusNotFound)` |
| Banner file missing | 404 | `http.Error(w, "...", http.StatusNotFound)` |
| Empty or missing form fields | 400 | `http.Error(w, "...", http.StatusBadRequest)` |
| Invalid banner name | 400 | `http.Error(w, "...", http.StatusBadRequest)` |
| Template fails to render | 500 | `http.Error(w, "...", http.StatusInternalServerError)` |

**Learn more:** https://gowebexamples.com/http-status-codes/

---

### Allowed Packages

Only Go standard library packages are used:

| Package | Used for |
|---|---|
| `net/http` | HTTP server, handlers, status codes |
| `html/template` | Rendering HTML templates with Go data |
| `os` | Reading banner files from disk |
| `strings` | Splitting and building strings |

---

## Learning Resources

| Topic | Resource |
|---|---|
| Everything in one place | https://gowebexamples.com/ |
| HTTP server basics | https://gowebexamples.com/hello-world/ |
| Templates | https://gowebexamples.com/templates/ |
| Forms | https://gowebexamples.com/forms/ |
| Status codes | https://gowebexamples.com/http-status-codes/ |
| Official Go web tutorial | https://go.dev/doc/articles/wiki/ |