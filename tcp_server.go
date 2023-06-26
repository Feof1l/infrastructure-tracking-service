package main

import (
	"fmt"
	"net"
)

// функция запускается как горутина
func process(conn net.Conn) {
	defer conn.Close()
	for {
		// считываем полученные в запросе данные
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Read error:", err)
			break
		}
		// выводим на консоль сервера диагностическую информацию
		fmt.Println(string(input))
		// отправляем данные клиенту
		//time.Sleep(time.Second * 6)
		conn.Write([]byte(input))
	}
}

func main() {

	fmt.Println("Start server...")

	// Устанавливаем прослушивание порта
	ln, err1 := net.Listen("tcp", ":1111")
	if err1 != nil {
		fmt.Println(err1.Error())
	}

	// Запускаем цикл обработки соединений
	for {
		// Принимаем входящее соединение
		conn, err2 := ln.Accept()
		if err2 != nil {
			fmt.Println(err2.Error())
			conn.Close()
			continue
		}
		// запускаем функцию process(conn)   как горутину
		go process(conn)
	}
}
