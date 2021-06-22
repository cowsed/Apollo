package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
)

type FileExplorer struct {
	MyParent   Editable
	title      string
	cwd        string
	selected   int
	files      []fs.FileInfo
	Delta      int
	DoWithFile func(filepath string)
	SelectDir  bool
	DoWithDir  func(path string)
}

func (f *FileExplorer) Parent() Editable       { return f.MyParent }
func (f *FileExplorer) SetParent(par Editable) { f.MyParent = par }

func (f *FileExplorer) MakeTitle() string {
	return f.title
}
func (f *FileExplorer) MakeRepresentation() [2]string {

	var first, second string
	if f.files == nil {
		f.files, _ = ioutil.ReadDir(f.cwd)
	}
	//Top is selected
	if f.selected == 0 {
		parts := strings.Split(f.cwd, "/")
		first = ">" + parts[len(parts)-2] + "/"
		if f.files[0].IsDir() {
			second = "&"
		}
		second += f.files[0].Name()
		return [2]string{first, second}
	}
	//Last Back button is selected
	if f.selected > len(f.files) {
		if f.files[len(f.files)-1].IsDir() {
			first = "&"
		}
		first += f.files[len(f.files)-1].Name()
		second = "> <-"
		return [2]string{first, second}
	}
	//Last Selected, show back button
	if f.selected == len(f.files) {
		first = "> "
		if f.files[len(f.files)-1].IsDir() {
			first += "&"
		}
		first += f.files[len(f.files)-1].Name()
		second = "<-"
		return [2]string{first, second}
	}
	//Base Case - normal
	first = "> "
	if f.files[f.selected-1].IsDir() {
		first += "&"
	}
	first += f.files[f.selected-1].Name()

	if f.files[f.selected].IsDir() {
		second = "&"
	}
	second += f.files[f.selected].Name()

	return [2]string{first, second}
}

func (f *FileExplorer) CanBeEntered() bool { return true }
func (f *FileExplorer) Up() {
	f.Delta++
	if f.Delta > ScrollAmt {
		f.Delta -= ScrollAmt
		f.selected++
	}
	if f.selected > len(f.files)+2 {
		f.selected = len(f.files) + 2
	}
}
func (f *FileExplorer) Down() {
	f.Delta--
	if f.Delta < ScrollAmt {
		f.Delta += ScrollAmt
		f.selected--
	}
	if f.selected < 0 {
		f.selected = 0
	}

}
func (f *FileExplorer) Enter() {
	//On the cwd
	if f.selected == 0 {
		if f.SelectDir {
			f.DoWithDir(f.cwd)
			CurrentItem = f.MyParent
		}
		return
	}
	//Go Up file tree
	if f.selected > len(f.files) {
		if f.cwd != "/" {
			f.cwd = f.cwd[0:strings.LastIndex(f.cwd, "/")]
			f.cwd = f.cwd[0:strings.LastIndex(f.cwd, "/")] + "/"
		}
		fmt.Println(f.cwd)
		f.files = nil
		return
	}
	if f.files[f.selected-1].IsDir() {
		f.cwd += f.files[f.selected-1].Name() + "/"
		f.files = nil
		return
	}
	//Selected a file
	f.DoWithFile(f.cwd + f.files[f.selected-1].Name())
	CurrentItem = f.MyParent
}

type TimeHandler struct {
	VLCHandler *VLC
	MyParent   Editable
}

func (v *TimeHandler) CanBeEntered() bool { return true }

func (v *TimeHandler) Enter() {
	//Return to parent
	if v.MyParent != nil {
		CurrentItem = v.MyParent
	}
}
func (v *TimeHandler) Up() {
	v.VLCHandler.controller.Seek("+10")
}
func (v *TimeHandler) Down() {
	v.VLCHandler.controller.Seek("-10")
}
func (v *TimeHandler) Parent() Editable {
	return v.MyParent
}
func (v *TimeHandler) SetParent(par Editable) {
	v.MyParent = par
}
func (v *TimeHandler) MakeTitle() string { return "Completion" }

func (v *TimeHandler) MakeRepresentation() [2]string {
	v.VLCHandler.lastStatus, _ = v.VLCHandler.controller.GetStatus()
	return [2]string{"Completion", percentageBar(v.VLCHandler.lastStatus.Position, 16)}
}

type VolumeEdit struct {
	VLCHandler *VLC
	MyParent   Editable
}

func (v *VolumeEdit) CanBeEntered() bool { return true }

func (v *VolumeEdit) Enter() {
	//Return to parent
	if v.MyParent != nil {
		CurrentItem = v.MyParent
	}
}
func (v *VolumeEdit) Up() {
	v.VLCHandler.controller.Volume("+10")
}
func (v *VolumeEdit) Down() {
	v.VLCHandler.controller.Volume("-10")
}
func (v *VolumeEdit) Parent() Editable {
	return v.MyParent
}
func (v *VolumeEdit) SetParent(par Editable) {
	v.MyParent = par
}
func (v *VolumeEdit) MakeTitle() string { return "Volume" }

func (v *VolumeEdit) MakeRepresentation() [2]string {
	v.VLCHandler.lastStatus, _ = v.VLCHandler.controller.GetStatus()
	return [2]string{"Volume", percentageBar(float64(v.VLCHandler.lastStatus.Volume)/255, 16)}
}
