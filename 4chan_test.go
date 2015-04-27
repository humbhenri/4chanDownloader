package main

import "testing"

func TestThreadUrlShouldBePresent(t *testing.T) {
	osArgs := []string{}
	err := Download4Chan(osArgs)
	noURLError := "must inform thread url"
	if err.Error() != noURLError {
		t.Error("Must inform " + noURLError)
	}
}

func TestGetTitleFromEmptyUrl(t *testing.T) {
	title := getThreadTitle("")
	if title != "" {
		t.Error("empty url")
	}
}

func TestGetTitleFromThreadUrl(t *testing.T) {
	title := getThreadTitle("abcd")
	if title != "abcd" {
		t.Errorf("expected %s but was %s", "abcd", title)
	}
}

func TestGetTitleFromThreadUrl2(t *testing.T) {
	title := getThreadTitle("xyz/abcd")
	if title != "abcd" {
		t.Errorf("expected %s but was %s", "abcd", title)
	}
}

func TestGetTitleFromThreadUrl3(t *testing.T) {
	title := getThreadTitle("foo/xyz/abcd")
	if title != "abcd" {
		t.Errorf("expected %s but was %s", "abcd", title)
	}
}
