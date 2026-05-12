package history

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/Balaji01-4D/bubbline/editline"
)

var (
	t0 = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	t1 = time.Date(2024, 6, 1, 12, 5, 0, 0, time.UTC)
)

func TestLoadHistory(t *testing.T) {
	testCases := []struct {
		input  string
		exp    []editline.HistoryEntry
		expErr string
	}{
		// empty file
		{"", nil, ""},
		// malformed JSON — skipped entirely
		{"notjson\n", nil, ""},
		// single valid entry
		{
			`{"text":"ls -la","timestamp":"2024-06-01T12:00:00Z"}`,
			[]editline.HistoryEntry{{Text: "ls -la", Timestamp: t0}},
			"",
		},
		// two valid entries with trailing newline
		{
			"{\"text\":\"ls -la\",\"timestamp\":\"2024-06-01T12:00:00Z\"}\n{\"text\":\"git status\",\"timestamp\":\"2024-06-01T12:05:00Z\"}\n",
			[]editline.HistoryEntry{
				{Text: "ls -la", Timestamp: t0},
				{Text: "git status", Timestamp: t1},
			},
			"",
		},
		// blank lines are skipped
		{
			"\n{\"text\":\"ls -la\",\"timestamp\":\"2024-06-01T12:00:00Z\"}\n\n",
			[]editline.HistoryEntry{{Text: "ls -la", Timestamp: t0}},
			"",
		},
		// malformed entries are skipped, valid ones kept
		{
			"{\"text\":\"ls -la\",\"timestamp\":\"2024-06-01T12:00:00Z\"}\nnotjson\n{\"text\":\"git status\",\"timestamp\":\"2024-06-01T12:05:00Z\"}",
			[]editline.HistoryEntry{
				{Text: "ls -la", Timestamp: t0},
				{Text: "git status", Timestamp: t1},
			},
			"",
		},
	}

	for _, tc := range testCases {
		buf := bytes.NewBufferString(tc.input)
		h, err := loadHistoryFromFile(buf)
		if tc.expErr != "" {
			if err == nil {
				t.Errorf("%q: expected error, got no error", tc.input)
			} else if err.Error() != tc.expErr {
				t.Errorf("%q: expected error:\n%s\ngot:\n%v", tc.input, tc.expErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q: expected no error, got: %v", tc.input, err)
		}
		if !reflect.DeepEqual(tc.exp, h) {
			t.Errorf("%q: expected:\n%+v\ngot:\n%+v", tc.input, tc.exp, h)
		}
	}
}

func TestSaveHistory(t *testing.T) {
	testCases := []struct {
		input  []editline.HistoryEntry
		exp    string
		expErr string
	}{
		// nil slice — nothing written
		{nil, "", ""},
		// single entry
		{
			[]editline.HistoryEntry{{Text: "ls -la", Timestamp: t0}},
			"{\"text\":\"ls -la\",\"timestamp\":\"2024-06-01T12:00:00Z\"}\n",
			"",
		},
		// two entries
		{
			[]editline.HistoryEntry{
				{Text: "ls -la", Timestamp: t0},
				{Text: "git status", Timestamp: t1},
			},
			"{\"text\":\"ls -la\",\"timestamp\":\"2024-06-01T12:00:00Z\"}\n{\"text\":\"git status\",\"timestamp\":\"2024-06-01T12:05:00Z\"}\n",
			"",
		},
		// entry with spaces — JSON handles it natively, no octal encoding
		{
			[]editline.HistoryEntry{{Text: "SELECT * FROM users", Timestamp: t0}},
			"{\"text\":\"SELECT * FROM users\",\"timestamp\":\"2024-06-01T12:00:00Z\"}\n",
			"",
		},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		err := saveHistoryToFile(&buf, tc.input)
		if tc.expErr != "" {
			if err == nil {
				t.Errorf("%v: expected error, got no error", tc.input)
			} else if err.Error() != tc.expErr {
				t.Errorf("%v: expected error:\n%s\ngot:\n%v", tc.input, tc.expErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%v: expected no error, got: %v", tc.input, err)
		}
		if result := buf.String(); result != tc.exp {
			t.Errorf("%v: expected:\n%q\ngot:\n%q", tc.input, tc.exp, result)
		}
	}
}

func TestAppendHistory(t *testing.T) {
	var buf bytes.Buffer
	initial := []editline.HistoryEntry{{Text: "initial", Timestamp: t0}}
	_ = saveHistoryToFile(&buf, initial)

	newEntries := []editline.HistoryEntry{{Text: "new", Timestamp: t1}}
	err := saveHistoryToFile(&buf, newEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\"text\":\"initial\",\"timestamp\":\"2024-06-01T12:00:00Z\"}\n{\"text\":\"new\",\"timestamp\":\"2024-06-01T12:05:00Z\"}\n"
	if buf.String() != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, buf.String())
	}
}
