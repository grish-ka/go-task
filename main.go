package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tasks  []Task // A list of our background processes
	cursor int    // Which task our cursor is currently pointing at
}

type Task struct {
	Name  string
	State string
	Ram   float64
	Cpu   float64
	Gpu   float64
}

type timerFinishedMsg struct {
    index    int
    newState string
}
func waitAndTransition(index int, sec int, newState string) tea.Cmd {
	return func() tea.Msg {
		// Convert our int to a time.Duration
		time.Sleep(time.Duration(sec) * time.Second) 
		return timerFinishedMsg{index: index, newState: newState}
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timerFinishedMsg:
		m.tasks[msg.index].State = msg.newState
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "i":
			if m.tasks[m.cursor].State == "Idle" {
				m.tasks[m.cursor].State = "Starting"
				return m, waitAndTransition(m.cursor, 2, " Running")
			} else if (m.tasks[m.cursor].State != "Idle" && m.tasks[m.cursor].State != "Stopped" || m.tasks[m.cursor].State != "Stopping"){
				m.tasks[m.cursor].State = "Idling" 
				return m, waitAndTransition(m.cursor, 2, "Idle")
			} else {
				m.tasks[m.cursor].State = "Starting"
				return m, waitAndTransition(m.cursor, 2, "Idle")
			}
		case "enter":
			if m.tasks[m.cursor].State != "Stopped" {
				m.tasks[m.cursor].State = "Stopping" 
				return m, waitAndTransition(m.cursor, 2, "Stopped")
			} else {
				m.tasks[m.cursor].State = "Starting"
				return m, waitAndTransition(m.cursor, 2, " Running")
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		}
	}
	
	// FIX 1: We must return the model even if no valid key was pressed!
	return m, nil
}

func (m model) View() string {
	s := "Go-Task Manager\n\n"

	for i, task := range m.tasks {
		// Is the cursor pointing at this task?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s [%-9s] %-18s | RAM: %4.2fGB | CPU: %5.1f%%\n",
			cursor, task.State, task.Name, task.Ram, task.Cpu)
	}

	s += "\npress j/k to move\n<enter> to stop/start\ni to idle/run q to quit\n"
	return s
}

// FIX 2: Added the missing main() function which uses the "os" import
func main() {
	// 1. Create initial dummy data
	initialModel := model{
		tasks: []Task{
			{Name: "Web Server", State: " Running", Ram: 1.2, Cpu: 0.5},
			{Name: "Docker Engine", State: "Idle", Ram: 2.4, Cpu: 0.1},
			{Name: "Discord Bot", State: "Stopped", Ram: 0.0, Cpu: 0.0},
		},
		cursor: 0,
	}

	// 2. Start the Bubble Tea program
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1) // This is where "os" is used!
	}
}