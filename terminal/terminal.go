package main

import (
	//"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	//"os"
	"os/exec"
	//"os/signal"
	"strings"
	"time"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct{
	BorderColor lipgloss.Color
	InputField lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#2563EB")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct{
	response Response
	width int
	height int
	contentField textinput.Model
	styles *Styles
}

func New(response Response) *model {
	styles := DefaultStyles()
	contentField := textinput.New()
	contentField.Placeholder = "Enter today's entry here!"
	contentField.Focus()
	return &model{
		response: response,
		contentField: contentField,
		styles: styles,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

type Response struct {
	prompt string
	output string
}

func NewResponse(prompt string) Response {
	return Response{prompt: prompt}
}
/*
func newTrueResponse(prompt string) {
	response = NewResponse(prompt)
	model := NewLongAnswerField()
	response.input = model
	return response
}
*/
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	content := &m.response
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			content.output = m.contentField.Value()
			m.contentField.SetValue("")
			return m, tea.Quit
		case "ctrl+c":
			content.output = "99876-BAD"
			return m, tea.Quit
		}
	}
	m.contentField, cmd = m.contentField.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center, 
			m.response.prompt,
			m.styles.InputField.Render(m.contentField.View()),
		),
	)
}

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

func timeString(currentTime time.Time) string {
	weekday := currentTime.Weekday()
	month := currentTime.Month()
	day := currentTime.Day()
	year := currentTime.Year()
	retString := "Welcome to your daily terminal journal!\nToday is " + string(weekday) + " " +  string(month) + " " + day + ", " + year + "\n"
	return retString
}

func main() {
	//Interupt to prevent people from escaping hahaha

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
	
/*
	if complete.IsCompleted {
		fmt.Printf("%sYou have already completed your journal for the day! Goodbye!%s\n", Green, Reset)
		return
	}
*/
	currentTime := time.Now()
	prompt := timeString(currentTime)
	//New Model pointer
	response := NewResponse(prompt)
	m := New(response)

	//Creates new Bubble Tea Program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil{
		fmt.Printf("Tea has failed :(\n)")	
	}

	var content string
	content = m.response.output
	
	for len(strings.Split(content, " ")) < 5 {
		var prompt2 string
		if content == "99876-BAD"{
			prompt2 = Red + "Nice try with CTRL+C, finish your entry!" + Reset + "\nPlease write a journal entry - it must be at least 50 words\n"
		} else {
			prompt2 = Red + "NOT ENOUGH WORDS" + Reset + "\nPlease write a journal entry - it must be at least 50 words\n"
		}
		
		response2 := NewResponse(prompt2)
		m2 := New(response2)

		p2 := tea.NewProgram(m2, tea.WithAltScreen())
		if _, err := p2.Run(); err != nil{
			fmt.Printf("Tea has failed :(\n)")	
		}
		content = m2.response.output
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
	if res.StatusCode == 201{
		fmt.Println("Sucessful Journal Upload!")
	} else {
		fmt.Println("Failed to save to Journal Database!")
	}

	//fmt.Println(jsonString)
}
