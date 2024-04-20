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
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
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
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#22D3EE")).
		Padding(4, 4, 4, 4).
		Width(100).
		Border(lipgloss.ThickBorder(), true, true).
		BorderForeground(lipgloss.Color("#2563EB")).
		BorderStyle(lipgloss.RoundedBorder())
	output := fmt.Sprintf("Welcome to your daily terminal journal!\nToday is %s %s %d, %d\n", weekday, month, day, year)
	fmt.Println(style.Render(output))
}

func main() {
	//Interupt to prevent people from escaping hahaha
	signal.Ignore(os.Interrupt)

	currentTime := time.Now()
	//Message for start of the program

	startupMessage(currentTime)

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

	checkRes, err := http.Get("http://127.0.0.1:8080/timecheck/" + email)
	var complete Completed
	json.NewDecoder(checkRes.Body).Decode(&complete)
	if complete.IsCompleted {
		fmt.Printf("%sYou have already completed your journal for the day! Goodbye!%s\n", Green, Reset)
		return
	}
	fmt.Println("Please Write a journal entry - it must be at least 50 words")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var content string
	content = scanner.Text()

	for len(strings.Split(content, " ")) < 5 {
		fmt.Printf("%sNOT ENOUGH WORDS%s\nPlease write a journal entry - it must be at least 50 words\n", Red, Reset)
		scanner.Scan()
		content = scanner.Text()
	}

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
