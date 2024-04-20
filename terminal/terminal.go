package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"strings"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

var st string = Red + "Hello" + Reset

type Entry struct {
	Content   string    `json:"content"`
	UserEmail string    `json:"userEmail"`
	Time      time.Time `json:"time"`
}

type Completed struct {
	IsCompleted bool `json"isCompleted"`
}

func (e *Entry) setEntry(content string, email string, time time.Time) {
	e.Content = content
	e.UserEmail = email
	e.Time = time
}

func startupMessage(currentTime time.Time) {
	weekday := currentTime.Weekday()
	month := currentTime.Month()
	day := currentTime.Day()
	year := currentTime.Year()
	fmt.Printf("%sWelcome to your daily terminal journal!\nToday is %s %s %d, %d%s\n", Green, weekday, month, day, year, Reset)
	fmt.Println("Please Write a journal entry - it must be at least 50 words")
}

func handler(signal os.Signal) {
	if signal == syscall.SIGTERM{
		fmt.Printf("%sYOU CANNOT LEAVE UNTIL YOU WRITE YOUR ENTRY!%s", Red, Reset)
	}
}

func main() {
	//Interupt to prevent people from escaping hahaha
	signal.Ignore(os.Interrupt)
	

	currentTime := time.Now()
	//Message for start of the program
	startupMessage(currentTime)
	
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var content string
	content = scanner.Text()

	for len(strings.Split(content, " ")) < 5{
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
	email = email[:len(email)-1]

	fmt.Println("====================")

	var entry Entry
	entry.setEntry(content, email, currentTime)
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	jsonString := string(jsonBytes)
	res, err := http.Post("http://127.0.0.1:8080/newentry", "application/json", strings.NewReader(jsonString))
	if err != nil {
		panic(err)
	}
	fmt.Println(res.StatusCode)

	fmt.Println(jsonString)
}
