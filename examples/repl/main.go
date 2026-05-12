package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Balaji01-4D/bubbline/computil"
	"github.com/Balaji01-4D/bubbline/editline"
	"github.com/alecthomas/chroma/v2/quick"
)

var pgKeywords = []string{
	"SELECT", "FROM", "WHERE", "INSERT", "UPDATE", "DELETE",
	"CREATE", "DROP", "ALTER", "TABLE", "INDEX",
	"JOIN", "LEFT", "RIGHT", "INNER", "OUTER", "ON",
	"GROUP", "BY", "ORDER", "HAVING", "LIMIT", "OFFSET",
	"VALUES", "INTO", "DISTINCT", "AND", "OR", "NOT", "NULL",
	"TRUE", "FALSE", "AS",
}

func postgresAutocomplete(v [][]rune, line, col int) (string, editline.Completions) {
	word, wstart, wend := computil.FindWord(v, line, col)
	if word == "" {
		return "", nil
	}

	upperWord := strings.ToUpper(word)
	var matches []string
	for _, kw := range pgKeywords {
		if strings.HasPrefix(kw, upperWord) {
			matches = append(matches, kw)
		}
	}

	if len(matches) == 0 {
		return "", nil
	}

	return "", editline.SimpleWordsCompletion(matches, "Keywords", col, wstart, wend)
}

var (
	userInputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#908CAA"))
	appOutputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0DEF4"))
	statusBarStle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#908CAA")).
			Background(lipgloss.Color("#2A273F")).
			Padding(0, 1)
)

type Model struct {
	editline *editline.Model
	width    int
}

func NewModel() Model {
	el := editline.New(0, 0)
	el.SetHelpDisabled(true)
	el.SetHighlighter(func(s string) string {
		var buf bytes.Buffer

		// Parameters: Output Buffer, Source String, Language, Formatter, Theme
		err := quick.Highlight(&buf, s, "postgresql", "terminal256", "monokai")
		if err != nil {
			return s
		}

		return buf.String()
	})
	el.AutoComplete = postgresAutocomplete
	return Model{
		editline: el,
	}
}

func (m Model) Init() tea.Cmd {
	return m.editline.Focus()
}

func (m Model) Update(imsg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := imsg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.editline.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		// Adding explicit quit handling just in case editline swallows it
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	case editline.InputCompleteMsg:
		// 1. Retrieve the text
		// Assuming bubbline has a .Value() method. If the text is contained
		// inside the msg itself (e.g., string(msg)), you can swap this out.
		input := strings.TrimSpace(m.editline.Value())

		if input == "" {
			// Reset the input and wait for the next interaction
			m.editline.Reset()
			return m, nil
		}

		// 2. Process logic
		output := processInput(input)

		// 3. Format interaction for the terminal history
		formattedInteraction := fmt.Sprintf(
			"%s\n\n%s\n",
			userInputStyle.Render("> "+input),
			appOutputStyle.Render("✦ "+output),
		)

		// 4. Reset the editline component entirely to clear the old input
		m.editline.Reset()
		// 5. Fire tea.Printf to print the history above the active prompt safely
		// We explicitly use "%s" to prevent format injection bugs!
		return m, tea.Printf("%s", formattedInteraction)
	}

	// Update the editline component and capture its new state
	var nextCmd tea.Cmd
	m.editline, nextCmd = m.editline.Update(imsg)

	return m, nextCmd
}

func (m Model) View() tea.View {
	statusBar := statusBarStle.Width(m.width).Render(
		"workspace: ~/projects/pgxcli    branch: feat/47    quota: 6% used",
	)

	// Keeping your exact view logic for bubbline
	str := fmt.Sprintf("\n%s\n\n%s", m.editline.View(), statusBar)
	return tea.NewView(str)
}

func processInput(input string) string {
	if strings.ToLower(input) == "exit" {
		os.Exit(0)
	}
	return "Simulated AI response for: " + input
}

func main() {
	p := tea.NewProgram(NewModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
