package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

const OCTOPI_IP = "http://10.0.0.168/"
const API_KEY = "AB0FB46BDA7A4382BDF935274D372CC3"

//go:embed notes.html
var notesTemplateSrc string

var notesTemplater *template.Template = template.New("notes template")

var notes []Note = []Note{{Title: "Example Note", Body: "Here is an example note."}, {Title: "Note 2", Body: "Second Note."}}

var SendingPrinterData SimpleData

func main() {
	//Retrieve data

	go QueryOctopi(&SendingPrinterData)

	var err error
	notesTemplater, err = notesTemplater.Parse(notesTemplateSrc)
	if err != nil {
		panic(err)
	}

	//Api
	http.HandleFunc("/api/printer", printHandler)

	//Note editing
	http.HandleFunc("/notes", noteHandler)

	//Css stuff
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	http.ListenAndServe("localhost:8080", nil)
}

func noteHandler(w http.ResponseWriter, r *http.Request) {
	//Update Note
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	if r.Form.Get("NoteTitle") != "" {
		title := r.Form.Get("NoteTitle")
		fmt.Println(title)
		body := r.Form.Get("NoteBody")
		for i := range notes {
			if notes[i].Title == title {
				fmt.Println("Update")
				notes[i].Body = body
			}
		}
		fmt.Println(notes)
		//Redirect to base notes page
		http.Redirect(w, r, "/notes", http.StatusSeeOther)

	}

	err = notesTemplater.Execute(w, notes)
	if err != nil {
		log.Println(err)
	}

}

func printHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.MarshalIndent(SendingPrinterData, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, string(data))

}

func QueryOctopi(data *SimpleData) {
	var CurrentJob JobData
	var CurrentServer ServerData
	for {
		GetOctopiJobInformation(&CurrentJob)
		GetOctopiServerInformation(&CurrentServer)

		data.Filename = CurrentJob.Job.File.Name
		data.EstimatedPrintTime = CurrentJob.Job.EstimatedPrintTime
		data.Completion = CurrentJob.Progress.Completion
		data.State = CurrentJob.State

		data.ToolActual = CurrentServer.Temperature.Tool0.Actual
		data.ToolTarget = CurrentServer.Temperature.Tool0.Target

		data.BedActual = CurrentServer.Temperature.Bed.Actual
		data.BedTarget = CurrentServer.Temperature.Bed.Target

		data.Printing = CurrentServer.State.Flags.Printing

		time.Sleep(time.Second)
	}
}

func GetOctopiJobInformation(job *JobData) {
	const FOLDER = "api/job"
	//Request
	url := OCTOPI_IP + FOLDER + "?apikey=" + API_KEY
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	//Read and save
	b := resp.Body
	r, err := io.ReadAll(b)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(r, &job)
	if err != nil {
		log.Println(err)
		return
	}

}

func GetOctopiServerInformation(server *ServerData) {
	const FOLDER = "api/printer"
	//Request
	url := OCTOPI_IP + FOLDER + "?apikey=" + API_KEY
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	//Read and save
	b := resp.Body
	r, err := io.ReadAll(b)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(r, &server)
	if err != nil {
		log.Println(err)
		return
	}

}
