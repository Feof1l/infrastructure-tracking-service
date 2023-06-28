package main

import (
	"time"
)

func main() {
	srv := Server{
		Addr:         ":1111",
		IdleTimeout:  5 * time.Second,
		MaxReadBytes: 1000,
	}
	go srv.ListenAndServe()
	time.Sleep(4 * time.Second)
	//srv.Shutdown()
	select {}
}
