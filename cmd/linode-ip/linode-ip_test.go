package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFetchNoMatches(t *testing.T) {

	linodes := []Linode{
		Linode{Label: "foo", Ipv4: []string{"bizz", "buzz "}},
		Linode{Label: "bar", Ipv4: []string{"fizz", "fuzz "}},
		Linode{Label: "bum", Ipv4: []string{"big", "bug "}},
	}

	ipv4 := Fetch("bam", linodes...)
	if ipv4 != "" {
		t.Errorf("Expected no matches but found %s", ipv4)
	}

}

func TestFetchSingleMatch(t *testing.T) {

	linodes := []Linode{
		Linode{Label: "foo", Ipv4: []string{"bizz", "buzz "}},
		Linode{Label: "bar", Ipv4: []string{"fizz", "fuzz "}},
		Linode{Label: "bum", Ipv4: []string{"big", "bug "}},
	}

	ipv4 := Fetch("bar", linodes...)
	if ipv4 != "fizz" {
		t.Errorf("Expected fizz but found %s", ipv4)
	}

}

func TestMultipleMatcheSelectOne(t *testing.T) {

	linodes := []Linode{
		Linode{Label: "foo", Ipv4: []string{"bizz", "buzz "}},
		Linode{Label: "bar", Ipv4: []string{"fizz", "fuzz "}},
		Linode{Label: "bum", Ipv4: []string{"big", "bug "}},
	}

	// Mock User Input
	content := []byte("1\n")
	tmpfile, err := ioutil.TempFile("", "example")

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	os.Stdin = tmpfile

	ipv4 := Fetch("b", linodes...)

	if ipv4 != "big" {
		t.Errorf("Expected big but found %s", ipv4)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

}

func MockStdIn(input string, test func() string) string {
	content := []byte(input)
	tmpfile, err := ioutil.TempFile("", "example")

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin
	os.Stdin = tmpfile

	test_result := test()

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return test_result
}

func TestMultipleMatchSelectExit(t *testing.T) {

	linodes := []Linode{
		Linode{Label: "foo", Ipv4: []string{"bizz", "buzz "}},
		Linode{Label: "bar", Ipv4: []string{"fizz", "fuzz "}},
		Linode{Label: "bum", Ipv4: []string{"big", "bug "}},
	}

	ipv4 := MockStdIn("e\n", func() string { return Fetch("b", linodes...) })

	if ipv4 != "" {
		t.Errorf("Expected no ipv4 but found %s", ipv4)
	}
}

func TestMultipleMatchUpdateMatcher(t *testing.T) {

	linodes := []Linode{
		Linode{Label: "foo", Ipv4: []string{"bizz", "buzz "}},
		Linode{Label: "bar", Ipv4: []string{"fizz", "fuzz "}},
		Linode{Label: "bum", Ipv4: []string{"big", "bug "}},
	}

	ipv4 := MockStdIn("u\nfoo\n", func() string { return Fetch("b", linodes...) })

	if ipv4 != "bizz" {
		t.Errorf("Expected ipv4 to be bizz but found %s", ipv4)
	}
}
