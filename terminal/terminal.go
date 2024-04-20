package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Progress Bar constants
const (
	padding  = 4
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2563EB")).Render

type tickMsg time.Time

type progressModel struct {
	progress progress.Model
}

func (mp progressModel) Init() tea.Cmd {
	return tickCmd()
}

func (mp progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		mp.progress.Width = msg.Width - padding*2 - 4
		if mp.progress.Width > maxWidth {
			mp.progress.Width = maxWidth
		}
		return mp, nil

	case tickMsg:
		if mp.progress.Percent() == 1.0 {
			return mp, tea.Quit
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := mp.progress.IncrPercent(1)
		return mp, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := mp.progress.Update(msg)
		mp.progress = progressModel.(progress.Model)
		return mp, cmd

	default:
		return mp, nil
	}
}

func (mp progressModel) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + mp.progress.View() + "\n\n" +
		pad + "Loading Your Journal..."
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Used for BubbleTea Wizard
type Styles struct {
	BorderColor lipgloss.Color
	BackColor   lipgloss.Color
	InputField  lipgloss.Style
}

// Used for BubbleTea Wizard
func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#22D3EE")
	s.BackColor = lipgloss.Color(Red)
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.DoubleBorder()).Padding(1).Width(80)
	return s
}

// Used for BubbleTea Wizard
type model struct {
	response     Response
	width        int
	height       int
	contentField textinput.Model
	styles       *Styles
}

// Used for BubbleTea Wizard
func New(response Response) *model {
	styles := DefaultStyles()
	contentField := textinput.New()
	contentField.Placeholder = "Enter today's entry here!"
	contentField.Focus()
	return &model{
		response:     response,
		contentField: contentField,
		styles:       styles,
	}
}

// Used for BubbleTea Wizard
func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

// Used for BubbleTea Wizard
type Response struct {
	prompt string
	output string
	ctrlC bool
}

// Used for BubbleTea Wizard
func NewResponse(prompt string) Response {
	return Response{prompt: prompt}
}

// One of the three major Bubble Tea functions, runs the update loop for the wizard
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
			//content.output = "99876-BAD"
			content.ctrlC = true
			return m, tea.Quit
		}
	}
	m.contentField, cmd = m.contentField.Update(msg)
	return m, cmd
}

// This decides how the View port is structured for Bubble Tea
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

// Color Constants
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// Entry Struct, same in Server
type Entry struct {
	Content   string    `json:"content"`
	UserEmail string    `json:"userEmail"`
	Time      time.Time `json:"time"`
}

// Completed struct, keeps track if the user has already completed their journal entry for the day
type Completed struct {
	IsCompleted bool `json"isCompleted"`
}

func (e *Entry) setEntry(content string, email string, time time.Time) {
	e.Content = content
	e.UserEmail = email
	e.Time = time
}

func timeString(currentTime time.Time) string {
	weekday := currentTime.Weekday()
	month := currentTime.Month()
	day := currentTime.Day()
	year := currentTime.Year()
	retString := fmt.Sprintf("Welcome to your daily terminal journal!\nToday is %s %s %d, %d \n", weekday, month, day, year)
	return retString
}

func main() {
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

	checkRes, err := http.Get("http://127.0.0.1:1234/timecheck/" + email)
	var complete Completed
	json.NewDecoder(checkRes.Body).Decode(&complete)

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#22D3EE")).
		Padding(2, 2, 2, 2).
		Width(45).
		Border(lipgloss.ThickBorder(), true, true).
		BorderForeground(lipgloss.Color("#2563EB")).
		BorderStyle(lipgloss.RoundedBorder())
	var output string

	if complete.IsCompleted {
		output = fmt.Sprintf("You have already completed your journal for the day! Goodbye!\n")
		fmt.Println(style.Render(output))
		return
	}

	mp := progressModel{
		progress: progress.New(progress.WithGradient("#22D3EE", "#2563EB")),
	}

	if _, err := tea.NewProgram(mp, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Progress Bar goofed", err)
		os.Exit(1)
	}

	currentTime := time.Now()
	prompt := timeString(currentTime)
	//New Model pointer
	response := NewResponse(prompt)
	m := New(response)

	//Creates new Bubble Tea Program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Tea has failed :(\n)")
	}

	var content string
	content = m.response.output

	for len(strings.Split(content, " ")) < 5 {
		var prompt2 string
		if m.response.ctrlC {
			prompt2 = "\033[31m" + "Nice try with CTRL+C, finish your entry!" + "\033[0m" + "\nPlease write a journal entry - it must be at least 50 words\n"
		} else {
			prompt2 = "\033[31m" + "NOT ENOUGH WORDS" + "\033[0m" + "\nPlease write a journal entry - it must be at least 50 words\n"
		}

		response2 := NewResponse(prompt2)
		m2 := New(response2)
		m2.contentField.SetValue(content)
		p2 := tea.NewProgram(m2, tea.WithAltScreen())
		if _, err := p2.Run(); err != nil {
			fmt.Printf("Tea has failed :(\n)")
		}
		content = m2.response.output
		m = m2
	}

	var entry Entry
	entry.setEntry(content, email, currentTime)
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		panic(err)
	}
	jsonString := string(jsonBytes)
	res, err := http.Post("http://127.0.0.1:1234/newentry", "application/json", strings.NewReader(jsonString))
	if err != nil {
		panic(err)
	}

	if res.StatusCode == 201 {
		output = fmt.Sprintf("Sucessful Journal Upload!\n")
	} else {
		output = fmt.Sprintf("Failed to save to Journal Database!\n")
		style.Background(lipgloss.Color("#FF0000"))
	}
	fmt.Println(style.Render(output))
}
