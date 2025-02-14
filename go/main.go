package main

import (
	"fmt"
	"log"

	"time"

	"github.com/sigurn/crc16"
	"go.bug.st/serial"
)

// Поиск подключенных COM портов
//
// ports, err := serial.GetPortsList()
// if err != nil {
// 	log.Fatal(err)
// 	fmt.Println("err")
// }
// if len(ports) == 0 {
// 	log.Fatal("No serial ports found!")
// 	fmt.Println("No serial ports found!")
// }
// for _, port := range ports {
// 	fmt.Printf("Found port: %v\n", port)
// }

func checkcrc(arr []byte) bool {
	if len(arr) < 4 {
		return false
	}
	crc2 := arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	crc1 := arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)
	crc := crc16.Checksum(arr, tab)
	tcrc1, tcrc2 := uint8(crc>>8), uint8(crc&0xff)
	if crc2 == tcrc1 && crc1 == tcrc2 {
		return true
	} else {
		return false
	}
}

func receive_msg(port serial.Port) []byte {
	// var buffer []byte = nil
	var buffer = make([]byte, 100)
	// var result []byte
	n, err := port.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	// n += 1
	if checkcrc(buffer[:n]) {
		return buffer[:n]
	} else {
		return nil
	}

	// return buffer
}

func open_serial() serial.Port {
	// Настройка порта
	port_settings := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		// StopBits: serial.OneStopBit,
	}

	port, err := serial.Open("COM4", port_settings)
	if err != nil {
		log.Fatal(err)
		fmt.Println("didint opened")
	}
	port.SetReadTimeout(time.Second * 5)
	return port
}

func open_chanel(port serial.Port) {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)
	crc := crc16.Checksum(req, tab)

	// Открыть канал связи
	req = append(req, 0b01000100) // Сетевой адрес
	req = append(req, 0b00000001) // Код запроса
	req = append(req, 0b00000001) // Уровень доступа
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	crc = crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)

	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	// answer := receive_msg(port)
	// fmt.Println(answer)
}

func reqCheck(port serial.Port) {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)
	crc := crc16.Checksum(req, tab)

	// Открыть канал связи
	req = append(req, 0b01000100) // Сетевой адрес
	req = append(req, 0b00000000) // Код запроса
	crc = crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)

	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	answer := receive_msg(port)
	fmt.Println(answer)
}

func reqPow1(port serial.Port) {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)
	// crc := crc16.Checksum(req, tab)

	// 4 (+2) байта мощность фазы 1
	req = append(req, 0b01000100) // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00001001) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	answer := receive_msg(port)
	fmt.Println(answer)
}

func reqTime(port serial.Port) {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)
	crc := crc16.Checksum(req, tab)

	// 4 (+2) байта мощность фазы 1
	req = append(req, 0b01000100) // Сетевой адрес
	req = append(req, 0b00000100) // Код запроса
	req = append(req, 0b00000001) // Номер параметра
	crc = crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	answer := receive_msg(port)
	fmt.Println(answer)
}

func reqSerNum(port serial.Port) {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)
	crc := crc16.Checksum(req, tab)

	// 4 (+2) байта мощность фазы 1
	req = append(req, 0b01000100) // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010100) // Номер параметра
	req = append(req, 0b00010001) // BWRI запрос
	// tab = crc16.MakeTable(crc16.CRC16_MODBUS)
	crc = crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	answer := receive_msg(port)
	fmt.Println(answer)
}

func requests() {
	port := open_serial()
	time.Sleep(time.Second * 1)
	open_chanel(port)
	port.ResetInputBuffer()
	// time.Sleep(time.Second * 1)
	// reqCheck(port)
	time.Sleep(time.Second * 1)
	reqSerNum(port)
	// port.Close()
}

func main() {
	requests()
}
