package main

import (
	"fmt"
	// "go.bug.st/serial"
)

func main() {
	// Поиск подключенных COM портов
	//
	// ports, err := serial.GetPortsList()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if len(ports) == 0 {
	// 	log.Fatal("No serial ports found!")
	// }
	// for _, port := range ports {
	// 	fmt.Printf("Found port: %v\n", port)
	// }

	// Настройка порта
	// port_settings := &serial.Mode{
	// 	BaudRate: 57600,
	// 	Parity:   serial.EvenParity,
	// 	DataBits: 7,
	// 	StopBits: serial.OneStopBit,
	// }

	req := [6]byte{0b00000001, 0b00001000, 0b00001011, 0b00001001}
	fmt.Println(req)
	// port, err := serial.Open("/dev/ttyUSB0", port_settings)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
