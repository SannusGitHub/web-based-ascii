DESCRIPTION:
A simple web-based ASCII art generator written in Go.
It runs a local HTTP server that lets you create ASCII art directly from your browser. Sample project features three different fonts, and a single page to send a string of text to the backend and return an altered string to the front-end

FEATURES:
* HTTP server host
* Basic index website as an interface
* Error handling
* Three ASCII font types

USAGE:
1. Clone the repository:
``
2. Run the server:
`go run .`
or
`go run main.go`
2. Open your browser and visit:
`http://localhost:8080`

DETAILS:
The algorithm is a basic (and rather inefficient) O(n^3) implementation which just goes over the text letter by letter and prints it until all eight rows are filled

This project was made when I was first learning Golang as a part of my college course covering algorithms, HTTP servers, handlers and templates. As such, the code here is heavily messy as I got a proper basic grip. Uploaded (mostly) for the sake of archival purposes.