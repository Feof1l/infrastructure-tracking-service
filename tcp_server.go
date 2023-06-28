package main

import (
	"log"
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
			log.Println("Read error:", err)
			break
		}

		// выводим на консоль сервера диагностическую информацию
		log.Println(string(input))
		// отправляем данные клиенту
		//time.Sleep(time.Second * 6)
		conn.Write([]byte(input))
	}
}
func shutdown(done chan bool, ln net.Listener) chan bool {
	done <- true // We can advance past this because we gave it buffer of 1
	ln.Close()   // Now it the Accept will have an error above
	return done
}
func main() {
	stop_chan := make(chan bool, 1)
	log.Println("Start server...")

	// Устанавливаем прослушивание порта
	ln, err1 := net.Listen("tcp", ":1111")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	//log.Println("Start server...")

	// Запускаем цикл обработки соединений
	for {
		// Принимаем входящее соединение
		//stop_chan = shutdown(stop_chan, ln)
		conn, err2 := ln.Accept()
		if err2 != nil {
			select {
			case <-stop_chan:
				log.Println("Start is stoped...")
			default:
				log.Fatal("Accept failed:", err2.Error())
			}
			return
		}
		/*if err2 != nil {
			log.Println(err2.Error())
			conn.Close()
			continue
		}*/
		// запускаем функцию process(conn)   как горутину
		go process(conn)
	}
}

/*
select {
		//case <-stopChan:
		//	log.Println("Stopping server")
		//	return
		case <-stop:
			go func() {
				defer signal.Stop(stop)
				sig := <-stop
				action()
			}()
		default:
			// Принимаем входящее соединение
			conn, err2 := ln.Accept()
			if err2 != nil {
				log.Println(err2.Error())
				conn.Close()
				continue
			}
			// запускаем функцию process(conn)   как горутину
			go process(conn)
		}
	}*/
