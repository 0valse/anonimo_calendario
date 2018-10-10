package main

import (
    "text/template"
    "net/http"
    "path"
    "log"
    )


var (
    // компилируем шаблоны, если не удалось, то выходим
    post_template = template.Must(template.ParseFiles(path.Join("templates", "index.html")))
)


type Data struct {
    Description string
    Message string
    Title string
}



func postHandler(w http.ResponseWriter, r *http.Request) {
    // обработчик запросов
    
    d := Data{
        Description: "dfdfs",
        Message: "lkjlklj",
        Title: "54756756765",
    }
    
    if err := post_template.ExecuteTemplate(w, "index.html", d); err != nil {
        log.Println(err.Error())
        http.Error(w, http.StatusText(500), 500)
    }
}


func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    
    http.HandleFunc("/", postHandler)
    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
    
}

