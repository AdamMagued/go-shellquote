package shellquote

import (
	"os/exec"
	"testing"
)

func TestSimpleJoin(t *testing.T) {
	for _, elem := range simpleJoinTest {
		output := Join(elem.input...)
		if output != elem.output {
			t.Errorf("Input %q, got %q, expected %q", elem.input, output, elem.output)
		}
	}
}

var simpleJoinTest = []struct {
	input  []string
	output string
}{
	{[]string{"test"}, "test"},
	{[]string{"hello goodbye"}, "'hello goodbye'"},
	{[]string{"hello", "goodbye"}, "hello goodbye"},
	{[]string{"don't you know the dewey decimal system?"}, "'don'\\''t you know the dewey decimal system?'"},
	{[]string{"don't", "you", "know", "the", "dewey", "decimal", "system?"}, "don\\'t you know the dewey decimal system\\?"},
	{[]string{"~user", "u~ser", " ~user", "!~user"}, "\\~user u~ser ' ~user' \\!~user"},
	{[]string{"#comment", "c#omment", " #comment", "!#comment", "#", "##", "###"}, "\\#comment c#omment ' #comment' \\!#comment \\# \\## \\###"},
	{[]string{"cd", "#metoo"}, "cd \\#metoo"},
	{[]string{"foo*", "M{ovies,usic}", "ab[cd]", "%3"}, "foo\\* M\\{ovies,usic} ab\\[cd] %3"},
	{[]string{"one", "", "three"}, "one '' three"},
	{[]string{"some(parentheses)"}, "some\\(parentheses\\)"},
	{[]string{"$some_ot~her_)spe!cial_*_characters"}, "\\$some_ot~her_\\)spe\\!cial_\\*_characters"},
	{[]string{"' "}, "\\'' '"},
	{[]string{"~#", "#~"}, "\\~# \\#~"},
}

func TestJoinLeadingCommentHash(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"cd", "#metoo"},
			expected: "cd \\#metoo",
		},
		{
			input:    []string{"#"},
			expected: "\\#",
		},
		{
			input:    []string{"##"},
			expected: "\\##",
		},
		{
			input:    []string{"#tag", "plain#tag"},
			expected: "\\#tag plain#tag",
		},
		{
			input:    []string{"#with spaces"},
			expected: "'#with spaces'",
		},
	}

	for _, tc := range testCases {
		got := Join(tc.input...)
		if got != tc.expected {
			t.Errorf("Join(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestJoinLeadingCommentInSh(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not found")
	}
	args := []string{"printf", "%s\n", "cd", "#metoo", "##double", "mid#dle", "#"}
	cmdStr := Join(args...)
	out, err := exec.Command(sh, "-c", cmdStr).Output()
	if err != nil {
		t.Fatalf("failed to run %q: %v", cmdStr, err)
	}
	expected := "cd\n#metoo\n##double\nmid#dle\n#\n"
	if string(out) != expected {
		t.Errorf("got %q, expected %q", string(out), expected)
	}
}

