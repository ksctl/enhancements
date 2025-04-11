package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Plan represents a pricing plan with details
type Plan struct {
	price      string
	hourlyRate string
	ram        string
	cpu        string
	storage    string
	transfer   string
}

// Model represents the application state
type Model struct {
	plans        []Plan
	currentPlan  int
	windowWidth  int
	windowHeight int
	help         help.Model
	keys         keyMap
	quitting     bool
	selectedPlan int // -1 means no selection yet
}

type keyMap struct {
	left     key.Binding
	right    key.Binding
	selected key.Binding
	quit     key.Binding
	help     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.left, k.right, k.selected},
		{k.help, k.quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
		left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "previous plan"),
		),
		right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next plan"),
		),
		selected: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "select plan"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc", "quit"),
		),
	}
}

func newModel() Model {
	// Create pricing plans similar to the provided image
	plans := []Plan{
		{
			price:      "$8",
			hourlyRate: "$0.012/hour",
			ram:        "1 GB",
			cpu:        "1 Intel CPU",
			storage:    "35 GB NVMe SSDs",
			transfer:   "1000 GB transfer",
		},
		{
			price:      "$16",
			hourlyRate: "$0.024/hour",
			ram:        "2 GB",
			cpu:        "1 Intel CPU",
			storage:    "70 GB NVMe SSDs",
			transfer:   "2 TB transfer",
		},
		{
			price:      "$24",
			hourlyRate: "$0.036/hour",
			ram:        "2 GB",
			cpu:        "2 Intel CPUs",
			storage:    "90 GB NVMe SSDs",
			transfer:   "3 TB transfer",
		},
		{
			price:      "$32",
			hourlyRate: "$0.048/hour",
			ram:        "4 GB",
			cpu:        "2 Intel CPUs",
			storage:    "120 GB NVMe SSDs",
			transfer:   "4 TB transfer",
		},
		{
			price:      "$48",
			hourlyRate: "$0.071/hour",
			ram:        "8 GB",
			cpu:        "2 Intel CPUs",
			storage:    "160 GB NVMe SSDs",
			transfer:   "5 TB transfer",
		},
		{
			price:      "$64",
			hourlyRate: "$0.095/hour",
			ram:        "8 GB",
			cpu:        "4 Intel CPUs",
			storage:    "240 GB NVMe SSDs",
			transfer:   "6 TB transfer",
		},
	}

	return Model{
		plans:        plans,
		currentPlan:  0,
		help:         help.New(),
		keys:         newKeyMap(),
		selectedPlan: -1,
	}
}

func (m Model) Init() tea.Cmd {
	// Initialize with a window size command
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keys.help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keys.left):
			// Move to previous plan
			if m.currentPlan > 0 {
				m.currentPlan--
			}
			return m, nil

		case key.Matches(msg, m.keys.right):
			// Move to next plan
			if m.currentPlan < len(m.plans)-1 {
				m.currentPlan++
			}
			return m, nil

		case key.Matches(msg, m.keys.selected):
			// Select current plan
			m.selectedPlan = m.currentPlan
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		if m.selectedPlan >= 0 && m.selectedPlan < len(m.plans) {
			plan := m.plans[m.selectedPlan]
			return fmt.Sprintf("Selected plan: %s/mo (%s)\nThank you for your selection!\n", plan.price, plan.hourlyRate)
		}
		return "Exited without selecting a plan.\n"
	}

	// Calculate dimensions based on window size
	maxWidth := m.windowWidth
	if maxWidth == 0 {
		maxWidth = 100 // Default width
	}

	// Card styling
	cardWidth := 24
	// cardGap := 2
	visibleCards := 3 // Force exactly 3 visible cards

	// Add some top padding
	var builder strings.Builder
	builder.WriteString("\n\n")

	// Determine which cards to show
	var startIdx int
	if m.currentPlan == 0 {
		startIdx = 0
	} else if m.currentPlan == len(m.plans)-1 {
		startIdx = max(0, len(m.plans)-visibleCards)
	} else {
		startIdx = max(0, m.currentPlan-1) // Center the selected card
	}
	endIdx := min(len(m.plans), startIdx+visibleCards)

	// Create cards
	cards := make([]string, endIdx-startIdx)

	// Build each card
	for i := startIdx; i < endIdx; i++ {
		plan := m.plans[i]
		isActive := i == m.currentPlan

		// Define colors
		var borderColor, textColor, separatorColor lipgloss.Color
		var borderStyle lipgloss.Border

		if isActive {
			borderColor = lipgloss.Color("#0066FF") // Blue for active
			textColor = lipgloss.Color("#EEEEC7")   // Light color for text
			separatorColor = lipgloss.Color("#0066FF")
			borderStyle = lipgloss.RoundedBorder()
		} else {
			borderColor = lipgloss.Color("#444444") // Dark gray for inactive
			textColor = lipgloss.Color("#888888")   // Muted text for inactive
			separatorColor = lipgloss.Color("#444444")
			borderStyle = lipgloss.RoundedBorder()
		}

		// Price section
		priceStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Width(cardWidth - 4). // Adjusted for better border spacing
			Align(lipgloss.Center).
			PaddingTop(1)

		priceStr := fmt.Sprintf("%s/mo\n%s", plan.price, plan.hourlyRate)
		priceSection := priceStyle.Render(priceStr)

		// Separator
		separator := lipgloss.NewStyle().
			Foreground(separatorColor).
			Width(cardWidth - 4). // Adjusted for better border spacing
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1).
			Render(strings.Repeat("─", cardWidth-6)) // Adjusted width

		// Specs section
		specsStyle := lipgloss.NewStyle().
			Foreground(textColor).
			Width(cardWidth - 4). // Adjusted for better border spacing
			Align(lipgloss.Left).
			PaddingLeft(1)

		specsStr := fmt.Sprintf("%s / %s\n%s\n%s",
			plan.ram, plan.cpu, plan.storage, plan.transfer)
		specsSection := specsStyle.Render(specsStr)

		// Card container style
		cardStyle := lipgloss.NewStyle().
			Border(borderStyle).
			BorderForeground(borderColor).
			Width(cardWidth).
			Padding(0, 1) // Add horizontal padding inside the border

		// Put the card together
		cardContent := lipgloss.JoinVertical(lipgloss.Left,
			priceSection,
			separator,
			specsSection,
		)

		cards[i-startIdx] = cardStyle.Render(cardContent)
	}

	var rowContent strings.Builder

	rowContent.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cards...))

	centeredStyle := lipgloss.NewStyle().Width(maxWidth).Align(lipgloss.Center)
	builder.WriteString(centeredStyle.Render(rowContent.String()))
	builder.WriteString("\n\n")

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Center).
		Width(maxWidth)

	instructions := "← → to navigate • enter to select plan"
	builder.WriteString(instructionStyle.Render(instructions))

	return builder.String()
}

func main() {
	fmt.Println("Welcome to the Pricing Plan Selection!")
	model := newModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program and capture the final model state
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	// Cast the final model to our Model type
	m, ok := finalModel.(Model)
	if !ok {
		fmt.Println("Could not get final model state")
		return
	}

	// Print the selection message after the program exits
	if m.selectedPlan >= 0 && m.selectedPlan < len(m.plans) {
		plan := m.plans[m.selectedPlan]
		fmt.Printf("\nSelected plan: %s/mo (%s)\nThank you for your selection!\n",
			plan.price, plan.hourlyRate)
	} else if m.quitting {
		fmt.Println("\nExited without selecting a plan.")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
