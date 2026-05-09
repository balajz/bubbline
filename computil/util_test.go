package computil

import (
	"strings"
	"testing"
)

func TestConvertToCompletionQuery(t *testing.T) {
	testCases := []struct {
		input    string
		line     int
		col      int
		expected string
		pos      int
	}{
		{
			input:    "select * from user",
			line:     0,
			col:      18,
			expected: "select * from user🛇",
			pos:      18,
		},
		{
			input:    "line1\nline2",
			line:     1,
			col:      2,
			expected: "line1␤line2🛇",
			pos:      8,
		},
	}

	for _, tc := range testCases {
		rows := strings.Split(tc.input, "\n")
		var r [][]rune
		for _, row := range rows {
			r = append(r, []rune(row))
		}

		s, pos := Flatten(r, tc.line, tc.col)
		s = strings.ReplaceAll(s, "\n", "␤") + "🛇"

		if s != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, s)
		}
		if pos != tc.pos {
			t.Errorf("expected pos %d, got %d", tc.pos, pos)
		}
	}
}

func TestFindLongestCommonPrefix(t *testing.T) {
	td := []struct {
		a, b string
		ci   bool
		exp  string
	}{
		{``, ``, false, ``},
		{`a`, ``, false, ``},
		{``, `b`, false, ``},
		{`aab`, `ab`, false, `a`},
		{`aab`, `ab`, true, `a`},
		{`aab`, `aa`, false, `aa`},
		{`aab`, `aa`, true, `aa`},
		{`aab`, `Aba`, false, `a`},
		{`aab`, `Aba`, true, ``},
		{"\xc3\xb8", "\xc3\x98", true, ""},
		{"\xc3\xb8", "\xc3\x98", false, "ø"},
	}

	for i, tc := range td {
		p := FindLongestCommonPrefix(tc.a, tc.b, tc.ci)
		if p != tc.exp {
			t.Fatalf("%d: expected %q, got %q", i, tc.exp, p)
		}
	}
}

func TestFindWord(t *testing.T) {
	text := [][]rune{
		[]rune("there's no place"),
		[]rune("like home"),
	}

	word, s, e := FindWord(text, 1, 6)
	if word != "home" || s != 5 || e != 9 {
		t.Fatal("bad")
	}
	word, s, e = FindWord(text, 1, 2)
	if word != "like" || s != 0 || e != 4 {
		t.Fatal("bad")
	}
}
