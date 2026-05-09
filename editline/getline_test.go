package editline_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/chalk-ai/bubbline"
)

func TestBubblineBasic(t *testing.T) {
	m := bubbline.New()

	// Ensure the prompt is rendered by checking the view.
	view := m.View().Content
	if !strings.Contains(view, ">") {
		t.Errorf("expected view to contain default prompt '>', got: %q", view)
	}
	// Helper to send a sequence of keystrokes.
	typeString := func(s string) {
		for _, r := range s {
			msg := tea.KeyPressMsg{Text: string(r)}
			// Bubbletea uses Code to represent actual runes for typing.
			if r >= 'a' && r <= 'z' {
				msg.Code = r
			} else if r == ' ' {
				msg.Code = ' '
				msg.Text = " "
			}
			newM, _ := m.Update(msg)
			m = newM.(*bubbline.Editor)
		}
	}

	// Type "hello"
	typeString("hello")
	if m.Value() != "hello" {
		t.Errorf("expected value to be 'hello', got: %q", m.Value())
	}

	// Type " world"
	typeString(" world")
	if m.Value() != "hello world" {
		t.Errorf("expected value to be 'hello world', got: %q", m.Value())
	}

	// Press Enter to complete input
	newM, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(*bubbline.Editor)

	if m.Value() != "hello world" {
		t.Errorf("expected value to be 'hello world' after enter, got: %q", m.Value())
	}

	// Normally enter yields tea.Quit or similar, depending on the model's logic.
	// Bubbline's Update method should return a tea.Quit command upon Enter.
	_ = cmd
}

func TestBubblineHistory(t *testing.T) {
	m := bubbline.New()
	m.AddHistoryEntry("say hello to the world")
	m.AddHistoryEntry("peter parker was not spiderman")

	history := m.GetHistory()
	if len(history) != 2 {
		t.Errorf("expected 2 history entries, got: %d", len(history))
	}
	if history[0] != "say hello to the world" {
		t.Errorf("expected first history entry to match")
	}
}
