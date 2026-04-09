# ASCII Art Web

A web application built with Go that converts text into ASCII art using selectable banner styles. Users enter text in a browser, choose a font, and receive stylised ASCII art output rendered in real time.

## Features

- Three built-in banner styles: **Standard**, **Shadow**, and **Thinkertoy**
- Preserves user input and selection after form submission
- Server-side input validation and banner whitelisting
- XSS protection via Go's `html/template` engine
- Clean separation of concerns: routing, rendering logic, and templates

## Prerequisites

- [Go](https://go.dev/dl/) 1.22 or later

## Getting Started

```bash
git clone https://github.com/Emmanuellsensai/ascii-art-web.git
cd ascii-art-web
go run main.go
```

Open your browser and navigate to `http://localhost:8080`.

## Project Structure

```
ascii-art-web/
├── main.go              # HTTP server, routing, and request handling
├── go.mod               # Go module definition
├── ascii/
│   └── render.go        # Banner file parser, character map builder, art renderer
├── banners/
│   ├── standard.txt     # Standard font data
│   ├── shadow.txt       # Shadow font data
│   └── thinkertoy.txt   # Thinkertoy font data
├── templates/
│   └── index.html       # HTML template (UI)
└── README.md
```

## How It Works

1. The user visits `/` — the server renders an empty form.
2. The user enters text, selects a banner style, and submits.
3. The server receives a `POST` request at `/ascii-art`.
4. The selected banner file is read and parsed into a character map (each printable ASCII character mapped to 8 lines of art).
5. The input text is converted row-by-row into ASCII art using the map.
6. The result is injected into the HTML template and returned to the browser.

## Architecture

```
Browser ──GET /──> homeHandler ──> Render empty template
Browser ──POST /ascii-art──> asciiArtHandler ──> ReadBanner ──> BuildAsciiMap ──> PrintAscii ──> Render template with result
```

## Banner Format

Each banner file defines characters 32 (space) through 126 (~) in ASCII order. Every character occupies exactly **8 lines of art** separated by **1 blank line** (9 lines total per character). The file begins with a blank line.

## API Endpoints

| Method | Path         | Description                          |
| ------ | ------------ | ------------------------------------ |
| GET    | `/`          | Serves the main page                 |
| POST   | `/ascii-art` | Accepts text and banner, returns art |

## Validation and Security

- HTTP method enforcement on both endpoints
- Empty input rejection
- Banner value whitelisted against `standard`, `shadow`, `thinkertoy` to prevent path traversal
- `html/template` auto-escapes output to prevent XSS

## Technologies

- **Language:** Go 1.22
- **Standard Library:** `net/http`, `html/template`, `os`, `strings`, `log`
- **No external dependencies**

## Contributing

1. Fork this repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -m "Add my feature"`
4. Push to the branch: `git push origin feature/my-feature`
5. Open a Pull Request

## Author

**Emmanuel Usang** — [GitHub](https://github.com/Emmanuellsensai)

**Emmanuel Oyelade** — [GitHub](https://github.com/oyehimar-bot)

## License

This project is open-source and available for educational purposes.
