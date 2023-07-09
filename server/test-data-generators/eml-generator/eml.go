package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/gomail.v2"
)

// 3200 people x 3200 people x 100 emails = 1,024,000,000 emails
// average email is 434 words.
// names.txt
// words.txt

const testfile = "test.eml"

func main() {
	writeEml()
	readEml()
	os.Remove(testfile)
}

func readEml() {
	// Get the absolute path of the EML file
	filePath, err := filepath.Abs(testfile)
	if err != nil {
		log.Fatal(err)
	}

	// Read the EML file
	emlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create an io.Reader from the EML data
	emlReader := bytes.NewReader(emlData)

	// Parse the EML message
	msg, err := mail.ReadMessage(emlReader)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the relevant information
	from := msg.Header.Get("From")
	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")
	date, err := msg.Header.Date()
	if err != nil {
		fmt.Println("Error retrieving email date:", err)
		return
	}

	// Read the body of the message
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the extracted information
	fmt.Println("Date:", date.Format(time.RFC1123Z))
	fmt.Println("From:", from)
	fmt.Println("To:", to)
	fmt.Println("Subject:", subject)
	fmt.Println("Body:", string(body))
}

func writeEml() {
	// Create a new message
	msg := gomail.NewMessage()
	msg.SetHeader("From", "sender@example.com")
	msg.SetHeader("To", "recipient@example.com")
	msg.SetHeader("Subject", "Test EML")
	msg.SetBody("text/plain", "This is the body of the email.")

	// Set a custom date for the email
	date := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
	msg.SetDateHeader("Date", date)

	// Write the EML data to a buffer
	buf := &bytes.Buffer{}
	_, err := msg.WriteTo(buf)
	if err != nil {
		fmt.Println("Error writing EML:", err)
		return
	}

	// Write the buffer data to a file
	err = ioutil.WriteFile(testfile, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	fmt.Println("EML file created successfully.")
}
