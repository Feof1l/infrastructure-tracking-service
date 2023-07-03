package main

import (
	"time"
)

func main() {
	srv := Server{
		Addr:         ":1111",
		IdleTimeout:  13 * time.Second,
		MaxReadBytes: 5,
	}
	go srv.ListenAndServe()
	time.Sleep(4 * time.Second)
	//srv.Shutdown()
	select {}
}
