package main

import (
    "fmt"
    "log"
    "os"
    "bufio"
    "net/http"
    //"regexp"
    //"strings"
    //"github.com/gocolly/colly"
)

type Blog struct {
    baseUrl string
    extension string
    found bool
}

func checkForRSS(blog *Blog) error {
    var extensions [3] string
    extensions[0] = "/rss"
    extensions[1] = "/blog"
    extensions[2] = "?format=rss"

    for _,extension := range(extensions) {
        testUrl := blog.baseUrl + extension
        //fmt.Println("Checking " + testUrl)
        resp, err := http.Get(testUrl)
        if err != nil {
            log.Fatalf("http get error: %v", err)
        } else if resp.StatusCode == 200 {
            blog.extension = extension
            blog.found = true
//            fmt.Printf("Got a hit at %s\n", testUrl)
            return nil
        }
    }
    return nil
}

func readBlogs(filename string) []*Blog {
    file,err := os.Open("blogs.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    blogs := make([]*Blog, 0)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        blogs = append(blogs, &Blog{scanner.Text(), "",  false})
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return blogs
}

func main() {
        // TODO: things will break if there's no trailing slash
    blogs := readBlogs("blogs.txt")
    for _,b := range(blogs) {
        //fmt.Println("Checking for RSS on", b.baseUrl)
        checkForRSS(b)
        if b.found == true {
            fmt.Println(b.baseUrl + b.extension)
        }

    }
}
