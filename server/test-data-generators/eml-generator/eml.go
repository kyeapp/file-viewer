package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhillyerd/enmime"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/gomail.v2"
)

// 3200 people x 3200 people x 100 emails = 1,024,000,000 emails
// average email is 434 words.
// names.txt
// words.txt

const NAME_FILE = "first-names.txt"
const WORD_FILE = "words.txt"
const NUM_EMAILS = 100
const NUM_PEOPLES = 320
const NUM_WORDS_PER_EMAIL = 500
const DATA_DIR = "data-eml"

var bar *progressbar.ProgressBar

func main() {
	err := os.RemoveAll(DATA_DIR)
	check(err)
	err = os.Mkdir(DATA_DIR, 0777)
	check(err)

	names := loadNamesFromFile()
	words := loadWordsFromFile()
	emailCount := NUM_PEOPLES * NUM_PEOPLES * NUM_EMAILS
	fmt.Println("people count:", NUM_PEOPLES)
	fmt.Println("word count:", len(words))
	fmt.Println("email count from A to B:", NUM_EMAILS)
	fmt.Println("total email count:", emailCount)

	bar = progressbar.Default(int64(emailCount))

	for _, sender := range names {
		for _, receiver := range names {
			for k := 0; k < NUM_EMAILS; k++ {
				subject := fmt.Sprintf("%s-%s-email%d", sender, receiver, k)
				body := generateRandomWords(words, NUM_WORDS_PER_EMAIL)
				_ = writeEml(sender, receiver, subject, body)
				bar.Add(1)
				// filename := writeEml(sender, receiver, subject, body)
				// readEml(filename)
			}
		}
	}
}

func generateRandomWords(words []string, count int) string {
	rand.Seed(time.Now().UnixNano())

	var randomWords []string
	for i := 0; i < count; i++ {
		randomIndex := rand.Intn(len(words))
		randomWord := words[randomIndex]
		randomWords = append(randomWords, randomWord)
	}

	return strings.Join(randomWords, " ")
}

func loadWordsFromFile() []string {
	// Read the names file
	data, err := ioutil.ReadFile(WORD_FILE)
	check(err)

	// Split the file content into names
	words := strings.Split(strings.TrimSpace(string(data)), "\n")

	return words
}

func loadNamesFromFile() []string {
	file, err := os.Open(NAME_FILE)
	check(err)
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		names = append(names, name)
	}

	err = scanner.Err()
	check(err)

	return names[:NUM_PEOPLES]
}

func readEml(filename string) {
	// Get the absolute path of the EML file
	filePath, err := filepath.Abs(filename)
	check(err)

	// Read the EML file
	emlData, err := ioutil.ReadFile(filePath)
	check(err)

	// Create an io.Reader from the EML data
	emlReader := bytes.NewReader(emlData)

	// Parse the EML message
	msg, err := mail.ReadMessage(emlReader)
	check(err)

	// Extract the relevant information
	from := msg.Header.Get("From")
	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")
	date, err := msg.Header.Date()
	check(err)

	// Read the body of the message
	reader := bytes.NewReader(emlData)
	env, err := enmime.ReadEnvelope(reader)
	check(err)

	body := env.Text
	body = strings.ReplaceAll(body, "=\r\n", "")
	body = strings.ReplaceAll(body, "= ", "")

	// Print the extracted information
	fmt.Println("Date:", date.Format(time.RFC1123Z))
	fmt.Println("From:", from)
	fmt.Println("To:", to)
	fmt.Println("Subject:", subject)
	fmt.Println("Body:", string(body))
}

func writeEml(sender, receiver, subject, body string) (filename string) {
	filename = filepath.Join(DATA_DIR, fmt.Sprintf("%s.eml", subject))

	// Create a new message
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s@gmail.com", sender))
	msg.SetHeader("To", fmt.Sprintf("%s@gmail.com", receiver))
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	// Set a custom date for the email
	date := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	msg.SetDateHeader("Date", date)

	// Write the EML data to a buffer
	buf := new(bytes.Buffer)
	_, err := msg.WriteTo(buf)
	check(err)

	// Write the buffer data to a file
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	check(err)

	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
