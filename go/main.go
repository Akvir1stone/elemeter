package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sigurn/crc16"
	"go.bug.st/serial"
)

// Поиск подключенных COM портов
// func search() {
// 	ports, err := serial.GetPortsList()
// 	if err != nil {
// 		log.Fatal(err)
// 		fmt.Println("err")
// 	}
// 	if len(ports) == 0 {
// 		log.Fatal("No serial ports found!")
// 		fmt.Println("No serial ports found!")
// 	}
// 	for _, port := range ports {
// 		fmt.Printf("Found port: %v\n", port)
// 	}
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
	time.Sleep(time.Second)
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

func get_result(bytes []byte) int {
	fmt.Println(bytes)
	var res = []byte{}
	res = append(res, bytes[3])
	res = append(res, bytes[2])
	res = append(res, bytes[1])
	result := int(binary.BigEndian.Uint16(res))
	return result
}

func open_serial(com_port string) serial.Port {
	// Настройка порта
	port_settings := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		// StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(com_port, port_settings)
	if err != nil {
		log.Fatal(err)
		fmt.Println("didnt opened")
	}
	port.SetReadTimeout(time.Second * 5)
	return port
}

func open_chanel(port serial.Port, trys int, adr byte) error {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// Открыть канал связи
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00000001) // Код запроса
	req = append(req, 0b00000001) // Уровень доступа
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
	req = append(req, 0b00000001) // Пароль 6 байт
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
	if trys < 3 && len(answer) < 2 {
		trys = trys + 1
		return open_chanel(port, trys, adr)
	} else {
		if trys >= 3 {
			return errors.New("login failed several times, check connection")
		} else {
			return nil
		}
	}
}

// func reqCheck(port serial.Port, adr byte) {
// 	// Переменные для составления запроса
// 	var req []byte
// 	var crc1, crc2 uint8
// 	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

// 	// Открыть канал связи
// 	req = append(req, adr)        // Сетевой адрес
// 	req = append(req, 0b00000000) // Код запроса
// 	crc := crc16.Checksum(req, tab)
// 	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
// 	req = append(req, crc2)
// 	req = append(req, crc1)

// 	n, err := port.Write(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Sent %v bytes\n", n)

// 	answer := receive_msg(port)
// 	fmt.Println(answer)
// }

func reqPow1(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
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
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqPow2(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00001010) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqPow3(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00001011) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqVolt1(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00010001) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqVolt2(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00010010) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqVolt3(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00010011) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqCurr1(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00100001) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqCurr2(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00100010) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
	// answer := get_result(receive_msg(port))
	// fmt.Println(answer)
}

func reqCurr3(port serial.Port, adr byte) int {
	// Переменные для составления запроса
	var req []byte
	var crc1, crc2 uint8
	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

	// 4 (+2) байта мощность фазы 1
	req = append(req, adr)        // Сетевой адрес
	req = append(req, 0b00001000) // Код запроса
	req = append(req, 0b00010001) // Номер параметра
	req = append(req, 0b00100011) // BWRI запрос
	crc := crc16.Checksum(req, tab)
	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
	req = append(req, crc2)
	req = append(req, crc1)
	n, err := port.Write(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
	return get_result(receive_msg(port))
}

// func reqFreq(port serial.Port, adr byte) {
// 	// Переменные для составления запроса
// 	var req []byte
// 	var crc1, crc2 uint8
// 	tab := crc16.MakeTable(crc16.CRC16_MODBUS)

// 	// 4 (+2) байта мощность фазы 1
// 	req = append(req, adr)        // Сетевой адрес
// 	req = append(req, 0b00001000) // Код запроса
// 	req = append(req, 0b00010001) // Номер параметра
// 	req = append(req, 0b01000000) // BWRI запрос
// 	crc := crc16.Checksum(req, tab)
// 	crc1, crc2 = uint8(crc>>8), uint8(crc&0xff)
// 	req = append(req, crc2)
// 	req = append(req, crc1)
// 	n, err := port.Write(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Sent %v bytes\n", n)
// 	answer := get_result(receive_msg(port))
// 	fmt.Println(answer)
// }

// Данные счетчиков тут
// Название ком порта и сетевой адрес счетчика
func requests(device int) ([]int, error) {
	var port serial.Port
	var adr byte
	// Данные счетчиков
	// Название ком порта и сетевой адрес счетчика
	switch device {
	case 1:
		port = open_serial("COM9")
		adr = 0b01000100

	case 2:
		port = open_serial("COM9")
		adr = 0b01000100
	}
	var results = []int{}
	err := open_chanel(port, 0, adr)
	if err != nil {
		form := url.Values{}
		form.Add("device", strconv.Itoa(device))
		resp, err := http.PostForm("http://127.0.0.1:8000/no_connection", form)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(resp)
		return nil, err
	} else {

		time.Sleep(time.Second * 1)
		results = append(results, reqPow1(port, adr))

		results = append(results, reqPow2(port, adr))

		results = append(results, reqPow3(port, adr))

		results = append(results, reqVolt1(port, adr))

		results = append(results, reqVolt2(port, adr))

		results = append(results, reqVolt3(port, adr))

		results = append(results, reqCurr1(port, adr))

		results = append(results, reqCurr2(port, adr))

		results = append(results, reqCurr3(port, adr))

		port.Close()

		form := url.Values{}
		form.Add("power1", strconv.Itoa(results[0]))
		form.Add("power2", strconv.Itoa(results[1]))
		form.Add("power3", strconv.Itoa(results[2]))
		form.Add("voltage1", strconv.Itoa(results[3]))
		form.Add("voltage2", strconv.Itoa(results[4]))
		form.Add("voltage3", strconv.Itoa(results[5]))
		form.Add("current1", strconv.Itoa(results[6]))
		form.Add("current2", strconv.Itoa(results[7]))
		form.Add("current3", strconv.Itoa(results[8]))
		form.Add("device", strconv.Itoa(device))
		resp, err := http.PostForm("http://127.0.0.1:8000/rec", form)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(resp)

		return results, nil
	}
}

func routine() {
	res, err := requests(1)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)

	res, err = requests(2)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)

	// res, err = requests(3)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(res)
}

func main() {
	for {
		go routine()
		time.Sleep(time.Minute * 5)
	}
	// search()
}
