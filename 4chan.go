package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"time"

	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/cheggaaa/pb"
)

// Download4Chan Download all the images of a 4chan thread
// cmdArgs - command line arguments
func Download4Chan(cmdArgs []string) error {
	if len(cmdArgs) < 2 {
		return errors.New("must inform thread url")
	}

	threadURL := cmdArgs[1]

	threadTitle := getThreadTitle(threadURL)
	if threadTitle == "" {
		return errors.New("cannot get title of thread")
	}

	dir, err := createThreadImagesDir(threadTitle)
	if err != nil {
		return err
	}

	images, err := getImagesURL(threadURL)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	successChan := make(chan bool)
	remainingItems := len(images)
	progressBar := pb.StartNew(remainingItems)

	for _, img := range images {
		go downloadImage(dir, img, successChan, errChan)
	}

	for {
		select {
		case _ = <-successChan:
			remainingItems--
			progressBar.Increment()
		case err := <-errChan:
			return err
		}

		if remainingItems == 0 {
			progressBar.FinishPrint(dir + " saved")
			break
		}
	}

	return nil
}

func downloadImage(dir string, img string, successChan chan bool, errChan chan error) {
	imgURL, err := url.Parse(img)
	if err != nil {
		errChan <- err
		return
	}
	// create file
	frags := strings.Split(imgURL.Path, "/")
	imgName := dir + string(os.PathSeparator) + frags[len(frags)-1]
	file, err := os.Create(imgName)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	// download image data
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := check.Get("http:" + img) // add a filter to check redirect
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		errChan <- err
		return
	}

	successChan <- true
}

func createThreadImagesDir(title string) (dir string, err error) {
	dir = title + strconv.FormatInt(time.Now().UnixNano(), 10)
	err = os.Mkdir(dir, 0777)
	return
}

func getThreadTitle(thread string) (title string) {
	threadURL, err := url.Parse(thread)
	if err != nil {
		return thread
	}
	frags := strings.Split(threadURL.Path, "/")
	return frags[len(frags)-1]
}

func getImagesURL(thread string) (images []string, err error) {
	// get thread url
	resp, err := http.Get(thread)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// parse html
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	nodeIsFileThumb := func(n *html.Node) (isFileThumb bool, href string) {
		if n.Type == html.ElementNode && n.Data == "a" {
			isFileThumb = false
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "fileThumb" {
					isFileThumb = true
				}
				if a.Key == "href" {
					href = a.Val
				}
			}
		}
		return
	}

	var parse func(*html.Node)
	parse = func(n *html.Node) {
		fileThumb, href := nodeIsFileThumb(n)
		if fileThumb {
			images = append(images, href)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}
	}
	parse(doc)
	return
}

func main() {
	err := Download4Chan(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
