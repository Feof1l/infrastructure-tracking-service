package main

import (
	"log"
	"time"
)

func main() {
	now_date := time.Now()
	log.Println("Main started")

	log.Println("Current date and time is: ", now_date.Format("2006-01-02 15:04:05"))

	/*ch := make(chan struct{})
	go func(c chan struct{}) {
		Run()
		close(c)
	}(ch)
	<-ch
	*/
	Run()
	time.Sleep(1 * time.Second)
	var test1 checker = tcp_prov{}
	var test2 checker = http_prov{}
	test1.check("127.0.0.1", "3333")
	test2.check("127.0.0.1", "3333")

	test1.check("127.0.0.1", "4444")
	test2.check("127.0.0.1", "4444")

	log.Println("Main out!")
}
