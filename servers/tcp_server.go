package main

import (
	"bytes"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Addr         string
	IdleTimeout  time.Duration // через какое время бездействия закрывать коннект
	MaxReadBytes int64         // макс объем данных
	conns        map[*conn]struct{}
	listener     net.Listener
	inShutdown   bool
	mu           sync.Mutex
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	log.Printf("starting server on %v\n", addr)
	// Устанавливаем прослушивание порта
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()
	srv.listener = listener
	for {
		if srv.inShutdown {
			break
		}
		// Принимаем входящее соединение
		newConn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		log.Printf("accepted connection from %v", newConn.RemoteAddr())
		conn := &conn{
			Conn:          newConn,
			IdleTimeout:   srv.IdleTimeout,
			MaxReadBuffer: srv.MaxReadBytes,
		}
		//кол-во соединений
		srv.trackConn(conn)
		// дедлайн и на чтение, и на запись
		conn.SetDeadline(time.Now().Add(conn.IdleTimeout))
		// запускаем функцию process(conn)   как горутину
		go srv.process(conn)
	}
	return nil
}
func (srv *Server) process(conn *conn) error {
	// функция запускается как горутина
	defer func() {
		log.Printf("closing connection from %v", conn.RemoteAddr())
		conn.Close()
		srv.deleteConn(conn)
	}()
	deadline := time.After(conn.IdleTimeout)
	for {
		select {
		case <-deadline:
			return nil
		default:
			// считываем полученные в запросе данные
			input := make([]byte, (1024 * 4))
			n, err := conn.Read(input, srv.IdleTimeout)
			if n == 0 || err != nil {
				log.Println("Read error:", err)
				time.Sleep(time.Second * 3)
				break
			}

			// выводим на консоль сервера диагностическую информацию
			log.Println(string(input))
			// отправляем данные клиенту
			//time.Sleep(time.Second * 6)
			input = bytes.ToUpper(input)
			conn.Write([]byte(input), srv.IdleTimeout)
			deadline = time.After(conn.IdleTimeout)
		}

	}
	return nil
}

// перестаем принимать новые соединения.
// Опрашиваем оставшиеся соединения и ждем пока они прекратят свою работу.
// Как только все соединения закрыты, можно останавливать наш сервер.
func (srv *Server) Shutdown() {

	srv.inShutdown = true
	log.Println("shutting down...")
	srv.listener.Close()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		//новые соединения не принимаются, текущие пока продолжают свою работу.
		//Дальше, мы опрашиваем счетчик текущих соединений каждые 500ms.
		//Как только счетчик доходит до 0 мы останавливаем сервер.
		select {
		case <-ticker.C:
			log.Printf("waiting on %v connections", len(srv.conns))
		}
		if len(srv.conns) == 0 {
			return
		}
	}
}

// отслеживаем количество соединений
func (srv *Server) trackConn(con *conn) {
	// отложенное разблокированиме
	defer srv.mu.Unlock()
	//Для блокирования доступа к общему разделяемому ресурсу
	srv.mu.Lock()
	if srv.conns == nil {
		srv.conns = make(map[*conn]struct{})
	}
	srv.conns[con] = struct{}{}
}

// удаляем соединение после закрытия(мертвые)
func (srv *Server) deleteConn(conn *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	delete(srv.conns, conn)
}
