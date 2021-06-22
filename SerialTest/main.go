package main

import (
	"fmt"
	"log"
	"strings"

	"go.bug.st/serial"
)

func main() {

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

	// Send the string to the serial port
	n, err := port.Write([]byte("--------------------------------"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	// Read and print the response

	buff := make([]byte, 100)
	var counter uint8 = 0
	for {
		for {
			// Reads up to 100 bytes
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				fmt.Println("\nEOF")
				break
			}

			cmd := string(buff[:n])
			cmd = strings.ReplaceAll(cmd, "\n", "")
			redraw := false
			if strings.Contains(cmd, "go-up") {
				counter++
				redraw = true
			}
			if strings.Contains(cmd, "go-down") {
				counter--
				redraw = true
			}
			if strings.Contains(cmd, "btn-clk") {
				counter = 0
				redraw = true
			}
			if redraw {
				w := fmt.Sprintf("%d", counter)
				if len(w) < 32 {
					w += strings.Repeat(" ", 32-len(w))
				}
				w = w[0:32]
				port.Write([]byte(w))
				redraw = false
			}

			fmt.Printf("%sEOF\n", cmd)

			// If we receive a newline stop reading
			if strings.Contains(string(buff[:n]), "\n") {
				break
			}
			//w := fmt.Sprint("%02d")
			//w = fmt.Sprintln("'%-4s'", "john")

		}
	}

}
