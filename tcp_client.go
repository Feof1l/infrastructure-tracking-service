package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	// Подключаемся к сокету
	conn, err1 := net.Dial("tcp", "127.0.0.1:1111")
	if err1 != nil {
		fmt.Println(err1.Error())
	}

	// отложенное закрытие соединения, которое срабатывает при выходе из функции
	defer conn.Close()
	for {
		var source string
		fmt.Print("Введите слово: ")
		_, err := fmt.Scanln(&source)
		if err != nil {
			fmt.Println("Некорректный ввод", err)

		}
		// отправляем сообщение серверу
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			fmt.Println(err)
			return
		}
		// Прослушиваем ответ
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 800))
		//conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		//клиент может ожидать данные на чтение от сервера в течении 800 мс
		//По истечении этого времени операция чтения генерирует ошибку и соответственно происходит выход из цикла,
		//где мы пытаемся прочитать данные от сервера.
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				break
			}
			fmt.Print(string(buff[0:n]))
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 300))
			//conn.SetReadDeadline(time.Now().Add(time.Second * 3))
			//после прочтения первых 1024 байт таймаут сбрасывается до 300 миллисекунд.
			//То есть если в течение последующих 300 милисекунд сервер не пришлет никаких данных, то происходит выход из цикла и соответственно чтение данных заканчивается.

		}
		fmt.Println()
	}
}
