package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type Entry struct {
	Content   string    `json:"content"`
	UserEmail string    `json:"userEmail"`
	Time      time.Time `json:"time"`
}

func (e *Entry) setEntry(content string, email string, time time.Time) {
	e.Content = content
	e.UserEmail = email
	e.Time = time
}

func main() {

	signal.Ignore(os.Interrupt)

	fmt.Println("write a journal entry - it must be at least 50 words")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var content string
	content = scanner.Text()

	for len(strings.Split(content, " ")) < 5 {
		fmt.Println("write a journal entry - it must be at least 50 words")
		scanner.Scan()
		content = scanner.Text()
	}

	//get github username
	cmd := exec.Command("git", "config", "user.email")
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	email := outBuffer.String()

	currentTime := time.Now()
	fmt.Println("====================")

	var entry Entry
	entry.setEntry(content, email, currentTime)
	jsonString, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonString))
}
