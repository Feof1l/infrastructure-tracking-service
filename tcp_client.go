package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Server_client struct {
	Addr         string
	IdleTimeout  time.Duration // через какое время бездействия закрывать коннект
	MaxReadBytes int64         // макс объем данных
}

func main() {

	srv := Server_client{
		Addr:         ":1111",
		IdleTimeout:  300 * time.Millisecond,
		MaxReadBytes: 1000,
	}
	// Подключаемся к сокету
	newConn, err1 := net.Dial("tcp", srv.Addr)
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	conn := &conn{
		Conn:          newConn,
		IdleTimeout:   srv.IdleTimeout,
		MaxReadBuffer: srv.MaxReadBytes,
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
		conn.SetWriteDeadline(time.Now().Add(srv.IdleTimeout))
		// отправляем сообщение серверу
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			fmt.Println(err)
			return
		}

		// Прослушиваем ответ
		conn.SetDeadline(time.Now().Add(srv.IdleTimeout))
		//conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		//клиент может ожидать данные на чтение от сервера в течении 900 мс
		//По истечении этого времени операция чтения генерирует ошибку и соответственно происходит выход из цикла,
		//где мы пытаемся прочитать данные от сервера.
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				break
			}
			fmt.Print(string(buff[0:n]))
			conn.SetDeadline(time.Now().Add(time.Millisecond * 700))
			//conn.SetReadDeadline(time.Now().Add(time.Second * 3))
			//после прочтения первых 1024 байт таймаут сбрасывается до 300 миллисекунд.
			//То есть если в течение последующих 300 милисекунд сервер не пришлет никаких данных, то происходит выход из цикла и соответственно чтение данных заканчивается.

		}
		fmt.Println()
	}
}
