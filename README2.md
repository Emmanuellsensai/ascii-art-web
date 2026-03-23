# ascii-art-web — Complete Beginner's Guide

> A web-based ASCII art generator built in Go. Type text in a browser, pick a font, get big ASCII art.

---

## Table of Contents

1. [What This Project Does](#what-this-project-does)
2. [Project Structure — Every File Explained](#project-structure)
3. [How a Web Server Works (The Big Picture)](#how-a-web-server-works)
4. [main.go — Line by Line](#maingo--line-by-line)
5. [The HTML Template — index.html Explained](#the-html-template)
6. [ascii/render.go — The Art Engine](#asciirendergo--the-art-engine)
7. [HTTP Status Codes — What They Mean](#http-status-codes)
8. [How to Run the Project](#how-to-run-the-project)
9. [Simplified CSS — Shorter Version](#simplified-css)
10. [Can We Use ANSI Colors?](#can-we-use-ansi-colors)
11. [Learning Resources](#learning-resources)

---

## What This Project Does

Imagine you type "Hello" into a website and the server converts it into this:

```
 _    _          _  _          
| |  | |        | || |         
| |__| |   ___  | || |  ___    
|  __  |  / _ \ | || | / _ \   
| |  | | |  __/ | || || (_) |  
|_|  |_|  \___| |_||_| \___/   
```

That is what this project does. Here is how it works, step by step:

1. The user opens `http://localhost:8080` in a browser
2. They type text and pick a font style
3. They click Generate
4. The browser sends the text + font choice to the Go server
5. The Go server reads a font file, converts each letter into art rows, and builds the result
6. The server sends the result back and the browser displays it

---

## Project Structure

```
ascii-art-web/
├── main.go               ← The web server. Routes, handlers, startup.
├── go.mod                ← Tells Go the module name and version.
├── ascii/
│   └── render.go         ← Reads font files and builds the ASCII art.
├── templates/
│   └── index.html        ← The webpage the user sees.
└── banners/
    ├── standard.txt      ← Font file: standard style
    ├── shadow.txt        ← Font file: shadow style
    └── thinkertoy.txt    ← Font file: thinkertoy style
```

### What each file is responsible for

| File | Job |
|------|-----|
| `main.go` | Starts the server, listens for requests, calls the render logic |
| `ascii/render.go` | Opens `.txt` font files and draws the art character by character |
| `templates/index.html` | The HTML page — has the form, and shows the result |
| `banners/*.txt` | Each file is a font. Every printable character is drawn using 8 rows of text |
| `go.mod` | Required by Go. Declares the module name so imports work |

---

## How a Web Server Works

Before reading the code, you need to understand what a web server actually does.

### The request-response cycle

```
Browser                           Go Server
  |                                   |
  |--- GET /  ----------------------->|   "Give me the home page"
  |<-- 200 OK + HTML ---------------  |   Server sends the page back
  |                                   |
  |--- POST /ascii-art (text+font) -->|   "Here is the form data"
  |<-- 200 OK + HTML with result ---  |   Server sends page with art
```

Every time you click something or type a URL, your browser sends an **HTTP request**. The server reads it, does work, and sends back an **HTTP response**.

There are two types of requests used here:

- **GET**: "Please give me a page." No data sent. Used when loading the home page.
- **POST**: "Here is data, process it." Form data is included. Used when clicking Generate.

📖 Read more: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods  
📺 Watch: https://www.youtube.com/watch?v=iYM2zFP3Zn0 (HTTP Crash Course)

### What is a handler?

A handler is just a Go function that runs when a specific URL is visited. You tell Go: "When someone visits `/`, run `homeHandler`." That is called **routing**.

```go
http.HandleFunc("/", homeHandler)
http.HandleFunc("/ascii-art", asciiArtHandler)
```

📖 Read more: https://gowebexamples.com/hello-world/

---

## main.go — Line by Line

Here is the full file with every part explained:

```go
package main
```
Every Go file starts with a package name. `main` is special — it is the entry point, the file Go runs first.

```go
import (
    "ascii-art-web/ascii"   // our own render.go package
    "html/template"          // Go's built-in HTML template engine
    "net/http"               // Go's built-in HTTP server
)
```
`import` brings in packages — bundles of code you want to use.  
- `net/http` is everything you need for an HTTP server. No Express, no Flask — it is built into Go.  
- `html/template` lets you write HTML with `{{placeholders}}` that Go fills in.  
- `ascii-art-web/ascii` is our own package inside the `ascii/` folder.

📖 Read more: https://pkg.go.dev/net/http  
📖 Read more: https://pkg.go.dev/html/template

---

### The PageData struct

```go
type PageData struct {
    Result string
    Error  string
}
```

A **struct** is a container that groups related values together under one name.  
`PageData` is what we send to the HTML template. The template can then read `.Result` and `.Error`.

Think of it like a small parcel you hand to the HTML page:  
- `Result` = the finished ASCII art string  
- `Error` = an error message to show the user (if something went wrong)

If both are empty, the page shows nothing below the form — just a blank starting state.

---

### Loading the template

```go
var tmpl = template.Must(template.ParseFiles("templates/index.html"))
```

This runs **once when the program starts**, before any request arrives.

- `template.ParseFiles(...)` reads `index.html` from disk and parses it
- `template.Must(...)` is a wrapper — if the file is missing or has a syntax error, it crashes immediately with a clear message instead of silently breaking later

Why load it once at startup? Because reading a file from disk on every single request would be slow. Load it once, reuse it forever.

📖 Read more: https://gowebexamples.com/templates/

---

### The main() function

```go
func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/ascii-art", asciiArtHandler)
    http.ListenAndServe(":8080", nil)
}
```

`main()` is where execution begins.

**`http.HandleFunc`** registers a route. It says: "When someone visits this URL path, call this function."

| URL path | Handler function | What happens |
|----------|-----------------|--------------|
| `/` | `homeHandler` | Shows the empty form |
| `/ascii-art` | `asciiArtHandler` | Processes form data, returns art |

**`http.ListenAndServe(":8080", nil)`** starts the server on port 8080. `:8080` means "listen on all network interfaces, port 8080". The `nil` means "use the default router" (the one we just configured above).

This line **blocks forever** — the program stays running and keeps waiting for requests. That is intentional.

📖 Read more: https://gowebexamples.com/hello-world/  
📺 Watch: https://www.youtube.com/watch?v=5BIylxkudaE (Go Web Servers)

---

### homeHandler

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
```

Every handler receives exactly two arguments:
- `w http.ResponseWriter` — this is what you write your response INTO. Think of it as the envelope you fill and send back.
- `r *http.Request` — this contains everything about the incoming request: the URL, method, form data, headers.

```go
    if r.URL.Path != "/" {
        http.Error(w, "404 - Page Not Found", http.StatusNotFound)
        return
    }
```

Go's default router is not strict. If someone visits `/banana`, it will still match the `/` handler. So we manually check that the path is exactly `/`. If it is not, we send a 404.

`http.Error(w, message, statusCode)` is a shortcut that writes an error response in one line.

```go
    if r.Method != http.MethodGet {
        http.Error(w, "400 - Bad Request", http.StatusBadRequest)
        return
    }
```

Only GET requests should load the home page. If someone sends a POST (or DELETE, or anything else), we reject it with a 400 error.

`http.MethodGet` is just the string `"GET"` — using the constant makes the code cleaner and avoids typos.

```go
    err := tmpl.Execute(w, PageData{})
    if err != nil {
        http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
    }
```

`tmpl.Execute(w, data)` fills in the HTML template with the data you provide, and writes the result into `w` (the response).

We pass `PageData{}` with nothing filled in — both `Result` and `Error` are empty strings. This gives the user the blank starting page.

If Execute fails (extremely rare, but possible if the template breaks at runtime), we send a 500.

---

### asciiArtHandler

```go
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "400 - Bad Request", http.StatusBadRequest)
        return
    }
```

This route only accepts POST. If someone visits `/ascii-art` directly in the browser (which is a GET), they get a 400.

```go
    text := r.FormValue("text")
    banner := r.FormValue("banner")
```

`r.FormValue("name")` reads a form field by its `name=""` attribute from the HTML.

The HTML form has:
```html
<textarea name="text">...</textarea>
<input type="radio" name="banner" value="standard">
```

So `r.FormValue("text")` gives us what the user typed, and `r.FormValue("banner")` gives us `"standard"`, `"shadow"`, or `"thinkertoy"`.

⚠️ **The `name` attribute in HTML must match exactly what you pass to `r.FormValue`.** This is a common mistake. If your HTML says `name="input"` but your Go says `r.FormValue("text")`, you will get an empty string.

📖 Read more: https://gowebexamples.com/forms/

```go
    if text == "" || banner == "" {
        http.Error(w, "400 - Bad Request: text and banner are required", http.StatusBadRequest)
        return
    }
```

Both fields must exist. If either is empty, we return a 400.

```go
    allowed := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
    if !allowed[banner] {
        http.Error(w, "400 - Bad Request: invalid banner", http.StatusBadRequest)
        return
    }
```

A `map[string]bool` is like a lookup table. We use it to check if the banner name is one of the three valid values. This prevents someone sending `banner=../../etc/passwd` or any other unexpected value.

```go
    bannerLines, err := ascii.ReadBanner("banners/" + banner + ".txt")
    if err != nil {
        http.Error(w, "404 - Banner Not Found", http.StatusNotFound)
        return
    }
```

We build the file path by combining `"banners/"` + banner name + `".txt"`. For example: `"banners/standard.txt"`.

`ascii.ReadBanner` reads the file and returns its lines. If the file does not exist, `err` is not `nil` and we send a 404.

```go
    asciiMap := ascii.BuildAsciiMap(bannerLines)
    result := ascii.PrintAscii(text, asciiMap)
```

These two calls do all the art work:
1. `BuildAsciiMap` organises the font file lines into a lookup table: character → its 8 art rows
2. `PrintAscii` loops through the input text, looks up each character, and builds the final string

```go
    err = tmpl.Execute(w, PageData{Result: result})
    if err != nil {
        http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
    }
```

Send the result to the template. This time `Result` is filled in, so `{{if .Result}}` in the HTML will be true and the art box will appear.

---

## The HTML Template

Go's `html/template` package lets you embed Go variables into HTML using double curly braces `{{ }}`.

### How template data flows

```
Go code                             HTML template
----------                          -------------
tmpl.Execute(w, PageData{           {{if .Result}}
    Result: "  _   \n | | \n...",       <pre>{{.Result}}</pre>
})                           →      {{end}}
```

The dot `.` refers to the PageData struct you passed in.  
`.Result` is the `Result` field. `.Error` is the `Error` field.

### Key template directives

```html
{{if .Result}}
    <!-- only shows this block if Result is not empty -->
    <pre>{{.Result}}</pre>
{{end}}

{{if .Error}}
    <p>{{.Error}}</p>
{{end}}
```

### Why `<pre>` is critical

```html
<pre>{{.Result}}</pre>
```

The `<pre>` tag (preformatted text) tells the browser: **preserve every space and newline exactly as written**. Without it, HTML collapses all spaces into one and ignores newlines. ASCII art would be completely destroyed. Always use `<pre>` for ASCII art output.

### How the form works

```html
<form action="/ascii-art" method="POST">
    <textarea name="text" ...></textarea>
    <input type="radio" name="banner" value="standard" checked>
    <input type="radio" name="banner" value="shadow">
    <input type="radio" name="banner" value="thinkertoy">
    <button type="submit">Generate</button>
</form>
```

| Attribute | What it does | Must match in Go |
|-----------|--------------|-----------------|
| `action="/ascii-art"` | Where the form sends data | The route in `http.HandleFunc` |
| `method="POST"` | How it sends it | `r.Method != http.MethodPost` check |
| `name="text"` | Field identifier | `r.FormValue("text")` |
| `name="banner"` | Field identifier | `r.FormValue("banner")` |
| `value="standard"` | What gets sent for that radio | The value in `allowed` map |

📖 Read more: https://gowebexamples.com/forms/  
📺 Watch: https://www.youtube.com/watch?v=fNcJuPIZ2WE (HTML Forms explained)

---

## ascii/render.go — The Art Engine

### ReadBanner

```go
func ReadBanner(file string) ([]string, error) {
    data, err := os.ReadFile(file)
    if err != nil {
        return nil, err
    }
    lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
    return lines, nil
}
```

1. `os.ReadFile(file)` reads the entire `.txt` font file as raw bytes
2. `strings.ReplaceAll(..., "\r\n", "\n")` normalises Windows line endings to Unix
3. `strings.Split(..., "\n")` splits it into a slice of lines

The result is every line of the font file as a separate string in a slice.

### BuildAsciiMap

```go
func BuildAsciiMap(lines []string) map[rune][]string {
    asciiMap := make(map[rune][]string)
    char := 32
    for i := 1; i < len(lines); i += 9 {
        asciiMap[rune(char)] = lines[i : i+8]
        char++
    }
    return asciiMap
}
```

The font files have a very specific structure. Every character takes exactly **9 lines**: 8 art rows + 1 blank separator. The first character in every font file is space (ASCII code 32).

So:
- Space (ASCII 32) is at lines 1–8
- `!` (ASCII 33) is at lines 10–17
- `"` (ASCII 34) is at lines 19–26
- ...and so on

The formula: `position = (ASCII value - 32) * 9 + 1`

`BuildAsciiMap` loops through the file in jumps of 9 and creates a map where:  
- **Key**: the character as a `rune` (Go's type for a Unicode character)  
- **Value**: a slice of 8 strings — the 8 art rows for that character

### PrintAscii

```go
func PrintAscii(text string, asciiMap map[rune][]string) string {
    var result strings.Builder

    for i, line := range strings.Split(text, "\\n") {
        if line == "" {
            if i != 0 {
                result.WriteString("\n")
            }
            continue
        }
        for row := 0; row < 8; row++ {
            for _, ch := range line {
                if art, ok := asciiMap[ch]; ok {
                    result.WriteString(art[row])
                }
            }
            result.WriteString("\n")
        }
    }
    return result.String()
}
```

The user types `\n` to mean "new line in the art". This splits the input on that literal string `\n`.

For each segment of text (between `\n`s):
- Loop through row 0 to row 7 (all 8 art rows)
- For each row, loop through every character in the line
- Look up that character in `asciiMap` and write its art row
- After all characters in a row are written, add a real newline

This builds the art row by row — because all characters on one line share the same 8 rows.

📖 Read more about strings.Builder: https://pkg.go.dev/strings#Builder  
📖 Read more about maps: https://go.dev/blog/maps

---

## HTTP Status Codes

Status codes are numbers the server sends back to tell the browser what happened.

| Code | Name | When we use it |
|------|------|----------------|
| 200 | OK | Request succeeded. Go sends this automatically. |
| 400 | Bad Request | The user/client sent invalid data (missing fields, bad banner name) |
| 404 | Not Found | The URL doesn't exist, or the banner file is missing |
| 500 | Internal Server Error | Something unexpected broke on the server side |

```go
// Sending a 404 in Go
http.Error(w, "404 - Page Not Found", http.StatusNotFound)

// http.StatusNotFound    = 404
// http.StatusBadRequest  = 400
// http.StatusInternalServerError = 500
```

Using named constants like `http.StatusNotFound` instead of raw numbers `404` makes the code easier to read and harder to mis-type.

📖 Read more: https://gowebexamples.com/http-status-codes/  
📺 Watch: https://www.youtube.com/watch?v=LtNSd_4txVc (HTTP status codes explained)

---

## How to Run the Project

```bash
# 1. Clone the project
git clone <your-repo-url>
cd ascii-art-web

# 2. Run the server
go run .

# 3. Open in browser
# http://localhost:8080
```

If you change a `.go` file, stop the server (Ctrl+C) and run `go run .` again.

### Common errors and fixes

| Error | Cause | Fix |
|-------|-------|-----|
| `no required module provides package` | Missing `go.mod` or wrong module name | Run `go mod init ascii-art-web` at the project root only |
| `open banners/standard.txt: no such file` | Running from wrong folder | Make sure you `cd` into `ascii-art-web` first |
| `template: pattern matches no files` | Wrong path to template | Run `go run .` from the project root, not from a subfolder |
| `address already in use` | Port 8080 is taken | Stop any other process using 8080, or change the port |

---

## Simplified CSS

Yes — the CSS can be shortened. Here is a simpler version that achieves the same look with about half the lines:

```css
<style>
  * { box-sizing: border-box; }

  body {
    font-family: Arial, sans-serif;
    max-width: 800px;
    margin: 40px auto;
    padding: 0 20px;
    background: #f5f5f5;
  }

  h1 { color: #333; }

  form, .result-box {
    background: white;
    padding: 24px;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,.1);
    margin-top: 32px;
  }

  label { display: block; margin-bottom: 6px; font-weight: bold; color: #555; }

  textarea {
    width: 100%; height: 80px;
    padding: 10px; font-size: 16px;
    border: 1px solid #ccc; border-radius: 4px;
    resize: vertical;
  }

  .banners { display: flex; gap: 20px; margin: 16px 0; }
  .banners label { font-weight: normal; cursor: pointer; }

  button {
    background: #333; color: white;
    border: none; padding: 12px 28px;
    font-size: 16px; border-radius: 4px;
    cursor: pointer; margin-top: 8px;
  }
  button:hover { background: #555; }

  pre { overflow-x: auto; font-size: 14px; line-height: 1.2; }
</style>
```

**What changed:**
- `* { box-sizing: border-box; }` is a universal reset that removes the need for `box-sizing` on individual elements
- `form, .result-box` are grouped — they share the same card style, so no need to write it twice
- The `pre` block dropped explicit `background: white` and `color: black` because those are already the browser defaults
- `.banners label` lost the flex/align properties — radio buttons + text align fine without it in all modern browsers

---

## Can We Use ANSI Colors?

**Short answer: No — not in the browser.**

ANSI color codes like `\033[31m` (red) or `\033[32m` (green) are escape sequences designed for **terminals**. They work in `bash`, `zsh`, and any terminal emulator because the terminal reads them and sets the color.

A browser does **not** interpret ANSI codes. If you put `\033[31mHello\033[0m` inside a `<pre>` tag, the browser displays it literally as garbage characters: `←[31mHello←[0m`.

### The right tool for each environment

| Environment | How to add color |
|-------------|-----------------|
| Terminal output (`fmt.Println(...)`) | ANSI codes: `"\033[32m green \033[0m"` |
| Browser (`<pre>` in HTML) | CSS: `color: green` or `<span style="color:green">` |

### What you CAN do for colored art in the browser

If you want each character of the ASCII art to appear in different colors in the browser, the approach is:
1. Change `PrintAscii` to return HTML instead of plain text (wrap each character's rows in `<span style="color:...">`)
2. Use `template.HTML` type in Go so the template does not escape your `<span>` tags
3. Change `<pre>{{.Result}}</pre>` to `<pre>{{.Result}}</pre>` but pass `template.HTML(result)`

This is a valid enhancement but changes the architecture significantly. For the current project, plain text in `<pre>` is clean, correct, and simple — keep it.

---

## Learning Resources

### Go fundamentals
| Topic | Link |
|-------|------|
| Official Go tour (interactive) | https://go.dev/tour/ |
| Go by Example | https://gobyexample.com/ |
| Effective Go | https://go.dev/doc/effective_go |

### Web server in Go
| Topic | Link |
|-------|------|
| Go web examples (all topics) | https://gowebexamples.com/ |
| Hello World server | https://gowebexamples.com/hello-world/ |
| Templates | https://gowebexamples.com/templates/ |
| Forms | https://gowebexamples.com/forms/ |
| Status codes | https://gowebexamples.com/http-status-codes/ |
| Official Go wiki article | https://go.dev/doc/articles/wiki/ |

### HTTP concepts
| Topic | Link |
|-------|------|
| HTTP methods (GET, POST) | https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods |
| HTTP status codes | https://developer.mozilla.org/en-US/docs/Web/HTTP/Status |
| How the web works | https://developer.mozilla.org/en-US/docs/Learn/Getting_started_with_the_web/How_the_Web_works |

### HTML + CSS
| Topic | Link |
|-------|------|
| HTML forms | https://developer.mozilla.org/en-US/docs/Learn/Forms |
| The `<pre>` element | https://developer.mozilla.org/en-US/docs/Web/HTML/Element/pre |
| CSS box model | https://developer.mozilla.org/en-US/docs/Learn/CSS/Building_blocks/The_box_model |

### Video courses
| Topic | Link |
|-------|------|
| Go full course (freeCodeCamp) | https://www.youtube.com/watch?v=un6ZyFkqFKo |
| HTTP explained simply | https://www.youtube.com/watch?v=iYM2zFP3Zn0 |
| HTML + CSS crash course | https://www.youtube.com/watch?v=916GWv2Qs08 |

---

## Authors

- Emmanuel ([@emmanuellsensai](https://github.com/emmanuellsensai))
