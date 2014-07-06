
package main

import (
    "fmt"
    "log"
    "./goquery"
    "strings"
    "net/http"
    "html/template"
    "encoding/json"  
)

const baseUrl string = "http://en.wikipedia.org"
const searchURLPrefix string = "http://en.wikipedia.org/wiki/"
const selector string = "#mw-content-text p a"
const imageSelector string = "#mw-content-text table.infobox tbody tr td a img"
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

func scrape(url,selector string) (map[string]interface{}) {
    var doc *goquery.Document
    var e error
    responseMap := make(map[string]interface{})
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

    responseMap["links"] = urlMap;

    
    doc.Find(imageSelector).Each(func(i int, s *goquery.Selection) {
        imageLink,_ := s.Attr("src") 

        // fmt.Println(s.Text())
        if imageLink != "" && responseMap["imageURL"] == nil {            
            responseMap["imageURL"] = imageLink                        
        }
    });

    // for k,v := range urlMap {
    //     fmt.Println(k,"->",v)
    // }
    return responseMap
}
