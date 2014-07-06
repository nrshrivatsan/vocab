
package main

import (
    "fmt"
    "log"
    "github.com/PuerkitoBio/goquery"
    "strings"
    "net/http"
    "html/template"
    "encoding/json"  
)

const baseUrl string = "http://en.wikipedia.org"
const searchURLPrefix string = "http://en.wikipedia.org/wiki/"
const selector string = "#mw-content-text p a"
func main() {
       
    http.HandleFunc("/", viewHandler)    
    http.HandleFunc("/search", search)    
    http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
    http.ListenAndServe(":8080", nil)
}

func search(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.URL.Query())
    if r.URL.Query()["q"] != nil && r.URL.Query()["q"][0] != ""{
            queryString := r.URL.Query()["q"][0]
            
            result :=    scrape(searchURLPrefix+queryString,selector) 
            b,_ := json.Marshal(result)
            w.Header().Set("Content-Type", "application/json")
            w.Write(b)
        }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {    
    t, _ := template.ParseFiles("index.html")
    t.Execute(w,nil)
}

func scrape(url,selector string) (map[string]string) {
    var doc *goquery.Document
    var e error

    if doc, e = goquery.NewDocument(url); e != nil {
        log.Fatal(e)
    }
    urlMap := make(map[string]string)

    doc.Find(selector).Each(func(i int, s *goquery.Selection) {
        if s.Text() != "" {
            link,_ := s.Attr("href")                 
            if !strings.HasPrefix(link,"//") && strings.HasPrefix(link,"/wiki"){
            
            // link="http:"+link   
                urlMap[s.Text()] = baseUrl+link
            }
            
        }       
    })

    for k,v := range urlMap {
        fmt.Println(k,"->",v)
    }
    return urlMap
}
