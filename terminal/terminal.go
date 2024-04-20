package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type Entry struct {
	Content   string
	UserEmail string
	Time      time.Time
}

func main() {

	signal.Ignore(os.Interrupt)

	fmt.Println("write a journal entry - it must be at least 50 words")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var entry string
	entry = scanner.Text()

	for len(strings.Split(entry, " ")) < 50 {
		fmt.Println("write a journal entry - it must be at least 50 words")
		scanner.Scan()
		entry = scanner.Text()
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
	fmt.Println(entry)
	fmt.Println(email)
	fmt.Println(currentTime)
}
