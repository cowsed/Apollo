package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"

	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
)

const VLCPassword = "vlcpass"

const ScrollAmt = 10

type Editable interface {
	Up()
	Down()
	Enter()
	CanBeEntered() bool
	Parent() Editable
	SetParent(Editable)

	MakeRepresentation() [2]string
	MakeTitle() string
}

type Button struct {
	Title    func() string
	MyParent Editable
	OnClick  func()
}

func (b *Button) MakeTitle() string {
	return b.Title()
}
func (b *Button) Enter() {
	b.OnClick()
}
func (b *Button) Up()                           {}
func (b *Button) Down()                         {}
func (b *Button) CanBeEntered() bool            { return false }
func (b *Button) MakeRepresentation() [2]string { return [2]string{"", ""} }
func (b *Button) Parent() Editable              { return b.MyParent }
func (b *Button) SetParent(par Editable)        { b.MyParent = par }

type List struct {
	MyParent Editable
	Title    string
	Parts    []Editable
	Selected int
	Delta    int
}

func (l *List) MakeTitle() string {
	return l.Title
}

func (l *List) Up() {
	l.Delta++
	if l.Delta >= ScrollAmt {
		l.Delta %= ScrollAmt
		l.Selected++
	}
	if l.Selected >= len(l.Parts) {
		l.Selected = len(l.Parts)
	}
}
func (l *List) Down() {
	l.Delta--
	if l.Delta <= -ScrollAmt {
		l.Delta += ScrollAmt
		l.Selected--
		if l.Selected < 0 {
			l.Selected = 0
		}
	}
}
func (l *List) Enter() {
	if l.Selected == len(l.Parts) {
		if l.Parent() != nil {
			CurrentItem = l.Parent()
		}
		return
	}
	if l.Parts[l.Selected].CanBeEntered() {
		l.Parts[l.Selected].SetParent(l)
		CurrentItem = l.Parts[l.Selected]
	} else {
		l.Parts[l.Selected].Enter()
	}
}
func (l *List) CanBeEntered() bool {
	return true
}

func (l *List) Parent() Editable {
	return l.MyParent
}

func (l *List) SetParent(par Editable) {
	l.MyParent = par
}

func (l *List) MakeRepresentation() [2]string {

	var second, first string
	if l.Selected >= len(l.Parts) {
		first = l.Parts[l.Selected-1].MakeTitle()
		second = "> <-"

		return [2]string{first, second}
	}

	first = "> " + l.Parts[l.Selected].MakeTitle()
	if l.Selected+1 >= len(l.Parts) {
		second = "<-"
	} else {
		second = l.Parts[l.Selected+1].MakeTitle()
	}
	return [2]string{first, second}
}

type StringDisplay struct {
	MyParent   Editable
	title      string
	makeString func() string
}

func (s *StringDisplay) MakeRepresentation() [2]string {
	return [2]string{s.MakeTitle(), s.makeString()}
}
func (s *StringDisplay) MakeTitle() string {
	return s.title
}
func (s *StringDisplay) Up()   { s.Enter() }
func (s *StringDisplay) Down() { s.Enter() }
func (s *StringDisplay) Enter() {
	if s.MyParent != nil {
		CurrentItem = s.MyParent
	}
}
func (s *StringDisplay) CanBeEntered() bool     { return true }
func (s *StringDisplay) Parent() Editable       { return s.MyParent }
func (s *StringDisplay) SetParent(par Editable) { s.MyParent = par }

//Old stuff

type VLC struct {
	controller vlcctrl.VLC
	lastStatus vlcctrl.Status
}

func (v *VLC) Representation() [2]string {
	v.lastStatus, _ = v.controller.GetStatus()
	title := "No Title"
	if len(v.lastStatus.Information.Titles) > 0 {
		i := v.lastStatus.Information.Title
		title = fmt.Sprint(v.lastStatus.Information.Titles[i])
	}
	//v.controller
	return [2]string{title, percentageBar(float64(v.lastStatus.Volume)/255, 16)}
}
func (v *VLC) Init() {
	//Start vlc
	cmd := exec.Command("vlc", "--intf", "http", "--extraintf", "qt", "--http-password", VLCPassword)
	go cmd.Run()
	time.Sleep(time.Second)

	fmt.Println("got here")
	// Declare a local VLC instance on port 8080 with password "password"
	var err error
	v.controller, err = vlcctrl.NewVLC("127.0.0.1", 8080, VLCPassword)
	if err != nil {
		panic(err)
	}
	err = v.controller.Volume("225")
	check(err)

}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

type Octopi struct {
	lastData SimplePrinterData
}

func (o *Octopi) Init() {
	go func() {
		for {
			//Periodically request scaled down printer information
			time.Sleep(5 * time.Second)

			resp, err := http.Get("http://" + POLYMATHIA_IP + "/api/printer")
			if err != nil {
				log.Println("Octopi Error", err)
				continue
			}
			r, _ := io.ReadAll(resp.Body)
			json.Unmarshal(r, &o.lastData)

		}
	}()
}

func (o *Octopi) Representation() [2]string {
	s := [2]string{}
	s[0] = fmt.Sprintf("Printer: %v C", o.lastData.ToolActual) + "\n"
	s[1] = fmt.Sprintf("%v%% Complete", int(o.lastData.Completion*100)) + "\n"

	return s
}
func (o *Octopi) Next() {

}
func (o *Octopi) Previous() {

}
func (o *Octopi) Press() {

}

type SimplePrinterData struct {
	Filename           string
	EstimatedPrintTime float64
	Completion         float64
	State              string

	ToolActual, ToolTarget float64
	BedActual, BedTarget   float64
	Printing               bool
}
