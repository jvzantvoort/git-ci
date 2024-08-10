package main

import (
	"strings"
	"testing"
)

type ticketMatchingTest struct {
	inputstr  string
	outputstr string
	err       error
}

type buildMessageStringTest struct {
	num        int
	branchname string
	message    string
	outputstr  string
	err        error
	errstr     string
}

var ticketMatchingTests = []ticketMatchingTest{
	ticketMatchingTest{"feature/ABCD-1234-is-lala", "ABCD-1234", nil},
	ticketMatchingTest{"bugfix/ABCD-1234-is-lala", "ABCD-1234", nil},
	ticketMatchingTest{"release/ABCD-1234-is-lala", "ABCD-1234", nil},
}

var buildMessageStringTests = []buildMessageStringTest{
	buildMessageStringTest{1, "feature/ABCD-1234-is-lala", "foo to the bar", "ABCD-1234 foo to the bar", nil, ""},
	buildMessageStringTest{2, "feature/garbage", "foo to the bar", "foo to the bar", nil, "cannot find ticket in string"},
	buildMessageStringTest{3, "", "foo to the bar", "foo to the bar", nil, "branch is empty"},
}

// ErrorContains checks if the error message in out contains the text in
// want.
//
// This is safe when out is nil. Use an empty string for want if you want to
// test that err is nil.
func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

func TestPat1(t *testing.T) {
	m := &Message{}
	for _, test := range ticketMatchingTests {
		gotstr, err := m.ExtractTicket(test.inputstr)
		if gotstr != test.outputstr {
			t.Errorf("got %q, wanted %q", gotstr, test.outputstr)
		}
		if err != test.err {
			t.Errorf("got %q, wanted %q", err, test.err)
		}

	}
}

func TestBM(t *testing.T) {
	for _, test := range buildMessageStringTests {
		var err error
		m := &Message{Message: test.message}
		err = m.SetBranch(test.branchname)
		gotstr := m.String()
		if gotstr != test.outputstr {
			t.Errorf("got %q, wanted %q -> %d", gotstr, test.outputstr, test.num)
		}
		if !ErrorContains(err, test.errstr) {
			t.Errorf("Unexpected error %v (%d)", err, test.num)
		}

	}

}
