package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"go.bug.st/serial"
)

const POLYMATHIA_IP = "127.0.0.1:8085"
const MY_IP = "127.0.0.1:4040"

//Sensor Values
var (
	CelciusTemp = 0.0
)

var (
	//selected int      = 1
	//options  []Option = []Option{
	//	&Octopi{
	//		selectedPage: selected,
	//		lastData:     SimplePrinterData{},
	//	}, &VLC{
	//		controller: vlcctrl.VLC{},
	//	},
	//}
	VLCControl = VLC{
		controller: vlcctrl.VLC{},
		lastStatus: vlcctrl.Status{},
	}
	VLCEdit = &List{
		MyParent: nil,
		Title:    "VLC",
		Parts: []Editable{
			&StringDisplay{
				MyParent: nil,
				title:    "Now Playing",
				makeString: func() string {
					return fmt.Sprintf("-%v-", VLCControl.lastStatus.Information.Category["FileName"])
				},
			},
			&Button{
				Title:    func() string { return "Play/Pause" },
				MyParent: nil,
				OnClick: func() {
					VLCControl.controller.Pause()
				},
			},
			&Button{
				Title:    func() string { return "Stop" },
				MyParent: nil,
				OnClick: func() {
					VLCControl.controller.Stop()
				},
			},
			&FileExplorer{
				MyParent: nil,
				cwd:      "/home/rich/Music/",
				title:    "Add Song",
				selected: 0,
				Delta:    0,
				DoWithFile: func(filepath string) {
					VLCControl.controller.Add(filepath)
					err := VLCControl.controller.Play()
					check(err)
				},
				SelectDir: true,
				DoWithDir: func(path string) {
					VLCControl.controller.Add(path)
					VLCControl.controller.Play()
				},
			},
			&Button{
				Title:    func() string { return "Skip" },
				MyParent: nil,
				OnClick: func() {
					fmt.Println(VLCControl.lastStatus.Information.Category["FileName"])
					VLCControl.controller.Next()
				},
			},
			&Button{
				Title: func() string {
					//fmt.Println(VLCControl.lastStatus.Random)
					//if VLCControl.lastStatus.Random {
					//	return "Turn Off Shuffle"
					//}
					//return "Turn On Shuffle"
					return "Toggle Shuffle"
				},
				MyParent: nil,
				OnClick: func() {
					VLCControl.controller.ToggleLoop()
				},
			},

			&TimeHandler{
				VLCHandler: &VLCControl,
				MyParent:   nil,
			},
			&VolumeEdit{
				VLCHandler: &VLCControl,
				MyParent:   nil,
			},
		},
		Selected: 0,
		Delta:    0,
	}

	menu = &List{
		MyParent: nil,
		Title:    "Outer List",
		Parts: []Editable{
			&StringDisplay{
				MyParent: nil,
				title:    "Time",
				makeString: func() string {
					t := time.Now()
					return t.Format("Jan 2 15:04:05")
				},
			},
			&StringDisplay{
				MyParent: nil,
				title:    "Temp",
				makeString: func() string {
					f := CelciusTemp*1.8 + 32
					return fmt.Sprintf("%2.1f*F", f-8)
				},
			},
			VLCEdit,
		},
	}
	CurrentItem Editable = menu
)

func main() {
	VLCControl.Init()

	// Retrieve the port list
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	// Print the list of detected ports
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	// Open the first serial port detected at 9600bps N81
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports[0], mode)
	if err != nil {
		log.Fatal(err)
	}
	buff := make([]byte, 100)
	defer port.Write([]byte("Crashing...     Program exited  "))

	for {
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				fmt.Println("Disconnected")
				time.Sleep(time.Second)
				port.Close()
				port, err = serial.Open(ports[0], mode)
				check(err)
				break
			}

			cmd := string(buff[:n])
			cmd = strings.ReplaceAll(cmd, "\n", "")

			HandleSerial(port, cmd)
		}
	}
}

func HandleSerial(port serial.Port, cmd string) {
	redraw := false
	if strings.Contains(cmd, "go-up") {
		CurrentItem.Up()
		redraw = true
	}
	if strings.Contains(cmd, "go-down") {
		CurrentItem.Down()
		redraw = true
	}
	if strings.Contains(cmd, "btn-clk") {
		CurrentItem.Enter()
		fmt.Println("clk")
		redraw = true
	}
	if strings.Contains(cmd, "update") {

		tstr := strings.ReplaceAll(cmd, "update", "")
		tstr = strings.Trim(tstr, "\r ")
		temp, err := strconv.ParseFloat(tstr, 64)
		if err == nil {
			CelciusTemp = temp
		}
		redraw = true
	}
	//Resend the information
	if redraw {
		rep := [2]string{}
		rep = CurrentItem.MakeRepresentation()

		data := formatToBytes(rep)
		fmt.Println("Writing")
		port.Write(data)
	}

}

//Sanatizes
func formatToBytes(parts [2]string) []byte {
	p1 := constrainLen(parts[0], 16)
	p2 := constrainLen(parts[1], 16)
	return []byte(p1 + p2)
}

func constrainLen(s string, n int) string {
	if len(s) < n {
		s += strings.Repeat(" ", n-len(s))
	}
	return s[:n]
}
