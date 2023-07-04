package main

import (
    "github.com/gocolly/colly"
    "fmt"
    "log"
    "regexp"
    "strings"
    "os"
    "bufio"
    "net/http"
)

func ScanForKeywords(needle []string, haystack string) bool {
    for _,n := range needle {
        if strings.Contains(strings.ToLower(haystack), strings.ToLower(n)) {
            return true
        }
    }
    return false
}

func crawlUrl(url string) {
    included,e := regexp.Compile("^" + url + "[[:alpha:]-]*$")
    if e != nil {
        fmt.Println("that's it")
        log.Fatalf("regex compile failed")
    } else {
        fmt.Printf("regex compiled: %v\n", included)
    }

    c := colly.NewCollector(
        colly.URLFilters(included),
        colly.MaxDepth(2),
    //    colly.DisallowedURLFilters(excluded),
    )

    //c.OnRequest(func(r *colly.Request) {
    //    fmt.Println("Visiting", r.URL)
    //})

    c.OnError(func(_ *colly.Response, err error) {
        log.Println("Something went wrong:", err)
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Visited", r.Request.URL)
    })

    //
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        e.Request.Visit(e.Attr("href"))
    })

    c.OnXML("//h1", func(e *colly.XMLElement) {
        fmt.Println("scanning: " + e.Text)
        needles := []string{"linux"}
        if ScanForKeywords(needles, e.Text) {
            fmt.Println("--- Found match: " + e.Text)
        }
    })

    //c.OnScraped(func(r *colly.Response) {
    //    fmt.Println("Finished", r.Request.URL)
    //})

    fmt.Println("Starting with " + url)
    c.Visit(url)
}

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
    //        fmt.Println("Okay!")
            return nil
     //   } else {
      //      fmt.Printf("%+v\n", resp.Status)
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
    //    crawlUrl("https://www.malwarebytes.com/blog/")
    blogs := readBlogs("blogs.txt")
    for _,b := range(blogs) {
//        fmt.Println("TODO: Checking for RSS on", b.baseUrl)
        checkForRSS(b)
        if b.found == true {
            fmt.Println(b.baseUrl + b.extension)
        }

    }
}
