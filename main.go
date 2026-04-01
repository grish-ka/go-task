package main

import (
	"fmt"
	"os"
	"sort"
	"time"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	tasks  []Task
	cursor int
}

type Task struct {
	Name  string
	State string
	Ram   float64
	Cpu   float64
	Gpu   float64
}

var blueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

type timerFinishedMsg struct {
	index    int
	newState string
}

func waitAndTransition(index int, sec int, newState string) tea.Cmd {
	return func() tea.Msg {
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
				return m, waitAndTransition(m.cursor, 2, "Running")
			} else if (m.tasks[m.cursor].State != "Idle" && m.tasks[m.cursor].State != "Stopped" || m.tasks[m.cursor].State != "Stopping") {
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
				return m, waitAndTransition(m.cursor, 2, "Running")
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
	return m, nil
}

func (m model) View() string {
    // 1. Calculate totals first
    var totalCpu, totalRam, totalGpu float64
    for _, t := range m.tasks {
        totalCpu += t.Cpu
        totalRam += t.Ram
        totalGpu += t.Gpu
    }

    // 2. Sort tasks by RAM
    sort.Slice(m.tasks, func(i, j int) bool {
        return m.tasks[i].Ram > m.tasks[j].Ram
    })

    // 3. Define settings
    width := 20
    maxRam := 33.7

    // --- CPU Bar Logic ---
    cpuRatio := totalCpu / 100.0
    if cpuRatio > 1.0 { cpuRatio = 1.0 }
    cpuCnt := int(cpuRatio * float64(width))
    cpuBar := strings.Repeat("-", width)
    if cpuCnt > 0 {
        cpuBar = strings.Repeat("#", cpuCnt-1) + ">" + strings.Repeat("-", width-cpuCnt)
    }

    // --- GPU Bar Logic ---
    gpuRatio := totalGpu / 100.0 // Fixed: use totalGpu
    if gpuRatio > 1.0 { gpuRatio = 1.0 }
    gpuCnt := int(gpuRatio * float64(width))
    gpuBar := strings.Repeat("-", width)
    if gpuCnt > 0 {
        gpuBar = strings.Repeat("#", gpuCnt-1) + ">" + strings.Repeat("-", width-gpuCnt)
    }

    // --- RAM Bar Logic ---
    ramRatio := totalRam / maxRam
    if ramRatio > 1.0 { ramRatio = 1.0 }
    ramCnt := int(ramRatio * float64(width))
    ramPer := ramRatio * 100.0 // Fixed: calculate actual percentage
    ramBar := strings.Repeat("-", width)
    if ramCnt > 0 {
        ramBar = strings.Repeat("#", ramCnt-1) + ">" + strings.Repeat("-", width-ramCnt)
    }

    // Start building our final string 's'
    s := "Go-Task Manager 🚀\n\n"
    s += fmt.Sprintf("%-17s | %s %5.1f%%\n", "CPU", blueStyle.Render(cpuBar), totalCpu)
    s += fmt.Sprintf("%-17s | %s %5.1f%%\n", fmt.Sprintf("RAM %.1f/%0.1fGB", totalRam, maxRam), blueStyle.Render(ramBar), ramPer)
    s += fmt.Sprintf("%-17s | %s %5.1f%%\n", "GPU", blueStyle.Render(gpuBar), totalGpu)
    s += "\n"


	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		icon := ""
		switch task.State {
		case "Running":
			icon = ""
		case "Stopped":
			icon = ""
		case "Idle":
			icon = "󰒲"
		default:
			icon = "󰔚"
		}

		// Row updated to include GPU!
		s += fmt.Sprintf("%s [%s %-10s] %-20s | RAM: %4.2fGB | CPU: %5.1f%% | GPU: %5.1f%%\n",
			cursor, icon, task.State, task.Name, task.Ram, task.Cpu, task.Gpu)
	}

	// This is where we'll put the Total row next...
	s += "--------------------------------------------------------------------------------\n"

	s += "\nDISCLAIMER: Idled/Stopped/Stopping/Starting or Idling Tasks do not take the resources that they say they do,\nit was the last seen amount\n"

	s += "\npress j/k to move\n<enter> to stop/start\ni to idle/run q to quit\n"
	return s
}

func main() {
	initialModel := model{
		tasks: []Task{
			{Name: "Web Server", State: "Running", Ram: 1.2, Cpu: 0.5, Gpu: 0.0},
			{Name: "Docker Engine", State: "Idle", Ram: 2.4, Cpu: 0.1, Gpu: 0.0},
			{Name: "Discord Bot", State: "Stopped", Ram: 0.0, Cpu: 0.0, Gpu: 0.0},
			{Name: "AI Model", State: "Running", Ram: 8.5, Cpu: 15.2, Gpu: 85.5}, // Added a heavy GPU task
		},
		cursor: 0,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}