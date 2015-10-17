package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"io/ioutil"
	"strings"

	"github.com/cheggaaa/pb"

	"encoding/json"
)

type JSON4ChanPost struct {
	Tim      int
	Filename string
	Ext      string
	Board    string
}

type JSON4ChanThread struct {
	Title string
	Posts []*JSON4ChanPost
}

func GetJson(url string) (*JSON4ChanThread, error) {
	urlArr := strings.Split(url, "/")
	jsonlink := fmt.Sprintf("http://a.4cdn.org/%s/thread/%s.json",
		urlArr[3], urlArr[5])
	res, err := http.Get(jsonlink)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var data JSON4ChanThread
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	for _, p := range data.Posts {
		p.Board = urlArr[3]
	}
	data.Title = urlArr[len(urlArr)-1]
	return &data, nil
}

func DownloadPostImage(post JSON4ChanPost) error {
	link := fmt.Sprintf("http://images.4chan.org/%s/src/%d%s",
		post.Board, post.Tim, post.Ext)
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	file, err := os.Create(post.Filename + post.Ext)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file.Name(), contents, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func DownloadThread(url string) error {
	thread, err := GetJson(url)
	if err != nil {
		return err
	}
	if _, err := os.Stat(thread.Title); os.IsNotExist(err) {
		err = os.Mkdir(thread.Title, 0777)
		if err != nil {
			return err
		}
	}
	os.Chdir(thread.Title)
	count := 0
	for _, post := range thread.Posts {
		if post.Tim != 0 {
			count = count + 1
		}
	}
	bar := pb.StartNew(count)
	for _, post := range thread.Posts {
		if post.Tim != 0 {
			err = DownloadPostImage(*post)
			if err != nil {
				return err
			}
			bar.Increment()
		}
	}
	bar.Finish()
	return nil
}

// Download4Chan Download all the images of a 4chan thread
// cmdArgs - command line arguments
func Download4Chan(cmdArgs []string) error {
	if len(cmdArgs) < 2 {
		return errors.New("must inform thread url")
	}

	threadURL := cmdArgs[1]
	return DownloadThread(threadURL)
}

func main() {
	err := Download4Chan(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
