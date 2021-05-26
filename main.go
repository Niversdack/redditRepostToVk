package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/otiai10/gosseract"
	"github.com/turnage/graw/reddit"
)

func main() {
	bot, err := reddit.NewScript(".agent", 0)
	if err != nil {
		fmt.Println("Failed to create bot handle: ", err)
		return
	}
	dir:=os.Getenv("GOPATH")
	path := path.Join(dir, "/src/redditRepostToVk")
	harvest, err := bot.Listing("/r/memes/", "")
	if err != nil {
		fmt.Println("Failed to fetch /r/golang: ", err)
		return
	}
	var wg sync.WaitGroup
	for _, post := range harvest.Posts[:5] {
		fmt.Println(post.URL)
		go SaveImg(path,post.URL,post.Title,&wg)
		wg.Add(1)
	}
	wg.Wait()
}
func SaveImg(path string,url string,title string, wg *sync.WaitGroup)  {
	defer wg.Done()
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()
	//open a file for writing
	file, err := os.Create(path+"/tmp/"+strings.TrimSpace(title)+".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(path+"/tmp/"+strings.TrimSpace(title)+".jpg")
	text, _ := client.Text()
	fmt.Println(text)

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")

}