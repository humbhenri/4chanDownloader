package main

import (
	"os"
	"path"
	"testing"
)

func TestGetJson(t *testing.T) {
	url := "http://boards.4chan.org/g/thread/50843193/what-is-the-biggest-tech-failure-of-the-last-5"
	data, err := GetJson(url)
	if err != nil {
		t.Errorf("get json shoud not return an error in this context: %s\n",
			err.Error())
	}
	if len(data.Posts) == 0 {
		t.Errorf("should be more than one post but was none")
	}
	expectedTitle := "what-is-the-biggest-tech-failure-of-the-last-5"
	if data.Title != expectedTitle {
		t.Errorf("title should be %s but was %s\n", expectedTitle, data.Title)
	}
	for _, post := range data.Posts {
		if post.Board != "g" {
			t.Errorf("board should be %s but was %s\n", "g", post.Board)
		}
	}
}

func TestDownloadImage(t *testing.T) {
	post := JSON4ChanPost{
		Tim:      1445029221779,
		Filename: "serveimage",
		Ext:      ".jpg",
		Board:    "g",
	}
	os.Chdir(os.TempDir())
	err := DownloadPostImage(post)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	filename := "serveimage.jpg"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("no such file %s\n", filename)
	}
}

func TestDownloadThread(t *testing.T) {
	url := "http://boards.4chan.org/g/thread/50843193/what-is-the-biggest-tech-failure-of-the-last-5"
	os.Chdir(os.TempDir())
	err := DownloadThread(url)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	dir := path.Join(os.TempDir(), "what-is-the-biggest-tech-failure-of-the-last-5")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("no such file %s\n", dir)
	}
}
