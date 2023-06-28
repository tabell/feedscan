package main

import (
    "github.com/gocolly/colly"
    "fmt"
    "log"
    "regexp"
)


func main() {
    included,e := regexp.Compile("^https://www\\.uptycs\\.com/blog/[[:alpha:]-]*$")
    if e != nil {
        fmt.Println("that's it")
        log.Fatalf("regex compile failed")
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

    c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
        fmt.Println(e.ChildText("body"))
    })

    c.OnXML("//h1", func(e *colly.XMLElement) {
        fmt.Println(e.Text)
    })

    //c.OnScraped(func(r *colly.Response) {
    //    fmt.Println("Finished", r.Request.URL)
    //})

    c.Visit("https://www.uptycs.com/blog/")
}
