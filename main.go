package main

import (
	"ascii-art-web-export-file/Ascii"
	"ascii-art-web-export-file/Export"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"ascii-art-web-export-file/Error"
)

type PageData struct {
	Result string
}

var err error

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error.RenderErrorPage(w, http.StatusNotFound, "Page Not Found")
		return
	}

	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles("template/web.html")
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		}

	} else {
		Error.RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

}

func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			Error.RenderErrorPage(w, http.StatusBadRequest, "Bad Request")
			return
		}

		text := r.FormValue("Text")
		banner := r.FormValue("Font")

		result, err := Ascii.AsciiArt(text, banner)
		if err != nil {
			Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		data := PageData{
			Result: result,
		}
		renderMainPageWithData(w, data)

	} else {
		Error.RenderErrorPage(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}
}

func renderMainPageWithData(w http.ResponseWriter, data PageData) {
	tmpl, err := template.ParseFiles("template/web.html")
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

func handleExport(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		Error.RenderErrorPage(w, http.StatusBadRequest, "Bad Request")
		return
	}

	text := r.FormValue("Text")
	banner := r.FormValue("Font")
	docType := r.FormValue("ExportedFile")

	result, err := Ascii.AsciiArt(text, banner)
	if err != nil {
		Error.RenderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := PageData{
		Result: result,
	}
	export.Export(w, r, docType, data.Result)
}

func main() {
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/ascii-art", handleAsciiArt)
	http.HandleFunc("/export", handleExport)
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/style/", http.StripPrefix("/style/", fs))
	fmt.Println("starting server at port 8125\n")
	err = http.ListenAndServe(":8125", nil)
	if err != nil {
		log.Fatal(err)
	}
}
