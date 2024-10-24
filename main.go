package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// Styles
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	border = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	cardStyle = lipgloss.NewStyle().
			BorderStyle(border).
			BorderForeground(highlight).
			Padding(1).
			Width(30)

	titleStyle = lipgloss.NewStyle().
			Foreground(special).
			Bold(true).
			MarginLeft(1)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF")).
			Bold(true)
)

type model struct {
	cpuUsage    float64
	memoryUsage float64
	diskUsage   float64
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		if msg.(tea.KeyMsg).Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case tickMsg:
		// Update CPU usage
		cpuPercent, err := cpu.Percent(0, false)
		if err == nil {
			m.cpuUsage = cpuPercent[0]
		}

		// Update Memory usage
		memInfo, err := mem.VirtualMemory()
		if err == nil {
			m.memoryUsage = memInfo.UsedPercent
		}

		// Update Disk usage
		diskInfo, err := disk.Usage("/")
		if err == nil {
			m.diskUsage = diskInfo.UsedPercent
		}

		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	// Create cards for each metric
	cpuCard := cardStyle.Render(
		titleStyle.Render("CPU Usage") + "\n" +
			valueStyle.Render(fmt.Sprintf("%.2f%%", m.cpuUsage)),
	)

	memoryCard := cardStyle.Render(
		titleStyle.Render("Memory Usage") + "\n" +
			valueStyle.Render(fmt.Sprintf("%.2f%%", m.memoryUsage)),
	)

	diskCard := cardStyle.Render(
		titleStyle.Render("Disk Usage") + "\n" +
			valueStyle.Render(fmt.Sprintf("%.2f%%", m.diskUsage)),
	)

	// Layout all cards horizontally
	row := lipgloss.JoinHorizontal(
		lipgloss.Center,
		cpuCard,
		memoryCard,
		diskCard,
	)

	// Add a title and instructions
	header := lipgloss.NewStyle().
		Foreground(special).
		Bold(true).
		Padding(1).
		Render("System Monitor")

	footer := lipgloss.NewStyle().
		Foreground(subtle).
		Render("Press Ctrl+C to quit")

	// Join everything vertically
	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		row,
		footer,
	)
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
