package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type PageVariables struct {
	AsciiOutput string
}

func main() {
	http.HandleFunc("/", home)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("data"))))
	fmt.Println("currently running on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

var asciiCharactersInFile = []rune{
	' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.',
	'/', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>',
	'?', '@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '[', '\\', ']', '^',
	'_', '`', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
	'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '{', '|', '}', '~',
}

var fontFile = readFile("data/standard.txt")

func generateAscii(desiredMessage string, fontType string) string {
	fontFile = readFile("data/" + fontType + ".txt")

	// here we get the argument and we do a lil' bit of formatting
	stringArgument := strings.Replace(desiredMessage, `\n`, "\n", -1) // we replace any escaped versions of a line character
	parts := delimiter(stringArgument, "\n")                          // ...and we delimiter the string into two halves, since newlines are required by the project

	assembledAsciiArt := ""
	// we take both parts and for loop them over, then toss 'em to the textToAscii function for printing out
	// remember to account for newline splits.
	for _, part := range parts {
		assembledAsciiArt += textToAscii(part)
	}

	return assembledAsciiArt
}

func textToAscii(userRequestedText string) string {
	asciiRow := ""
	// if X is less than length of string (userRequestedText) then use function "readAsciiCharacterFromFile(_, row, _)",
	// else if X is longer than string append a new line onto it and add +1 to row,
	// then repeat until row is 8.

	// this is jank code in my opinion.
	// no, I will not be fixing it. if it works it works.
	// i've spent countless hours up at night trying to make it work.

	// oh my god i have butchered this code even more
	filteredText := ""
	for _, currentCharacter := range userRequestedText {
		for _, currentASCIICharacter := range asciiCharactersInFile {
			if currentCharacter == currentASCIICharacter {
				filteredText += string(currentCharacter)
			}
		}
	}

	// cycle through rows 1 to 8, since every ascii art character is 8 rows consistently.
	counter := 0
	for row := 1; row <= 8; row++ {
		for _, currentCharacter := range filteredText { // we get the current character in the user requested text
			for index, currentASCIICharacter := range asciiCharactersInFile { // we get the current ascii character in the table above
				if currentCharacter == currentASCIICharacter { // we check if current character matches any of the characters in the table
					// this will run if they're valid...

					// if the counter (which is the amount of characters printed) is not the length,
					if counter < len(filteredText) {
						// we print out the ascii character (index), with row # (row), from font file). then add 1 to counter.
						counter++
						asciiRow += readAsciiCharacterFromFile(index, row, fontFile)
					}

					// if the counter is more or equal (signifies that its the last character the user has requested)...
					if counter >= len(filteredText) {
						// we print a new line, and restart the counter
						counter = 0
						asciiRow += "\n"
					}
				}
			}
		}
	}

	return asciiRow
}

func readAsciiCharacterFromFile(characterToGet int, row int, dataToRead string) string {
	scanner := bufio.NewScanner(strings.NewReader(dataToRead))

	line := 0
	asciiCharacterToUse := ((characterToGet * 8) + characterToGet) + 1

	for scanner.Scan() {
		line++

		if line == asciiCharacterToUse+row {
			return scanner.Text()
		}
	}

	return ""
}

func readFile(fileName string) string {
	if fileName == "" {
		log.Fatalf("Provide a valid file name.")
	}

	fileContents, errorResult := os.ReadFile(fileName)
	if errorResult != nil {
		log.Fatalf("Unable to read file: %s.", fileName)
	}

	return string(fileContents)
}

func delimiter(stringToDelimiter string, delimiter string) []string {
	var result []string
	parts := strings.Split(stringToDelimiter, delimiter)

	for i, part := range parts {
		result = append(result, part)

		if i < len(parts)-1 {
			result = append(result, delimiter)
		}
	}

	return result
}

func handleForm(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		http.Error(writer, "Error parsing form data", http.StatusBadRequest)
		return
	}
}

func handleError(writer http.ResponseWriter, request *http.Request, status int) {
	writer.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(writer, "You have encountered the super secret and rare 404 page! Not Found!")
	}

	if status == http.StatusBadRequest {
		fmt.Fprint(writer, "You have encountered the even secreter and rarer 400 page! Bad Request!")
	}

	if status == http.StatusInternalServerError {
		fmt.Fprint(writer, "You have encountered the secretest and rarest 500 page! Internal Server Error!")
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		tmpl, err := template.ParseFiles("data/index.html")
		if err != nil {
			http.Error(writer, "Error loading HTML template", http.StatusInternalServerError)
			return
		}

		if request.URL.Path != "/" {
			handleError(writer, request, http.StatusNotFound)
			return
		}

		pv := PageVariables{
			AsciiOutput: "",
		}

		tmpl.Execute(writer, pv)
		return
	}

	handleForm(writer, request)

	typeOfFont := request.Form.Get("fontType")
	asciiTextRequest := request.Form.Get("userTextField")
	AsciiOutput := generateAscii(asciiTextRequest, typeOfFont)

	tmpl, err := template.ParseFiles("data/index.html")
	if err != nil {
		http.Error(writer, "Error loading HTML template", http.StatusInternalServerError)
		return
	}

	pv := PageVariables{
		AsciiOutput: AsciiOutput,
	}

	tmpl.Execute(writer, pv)
}
