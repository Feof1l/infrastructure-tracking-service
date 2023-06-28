package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	t.Run("incorrect addr", func(t *testing.T) {
		//t.Parallel()
		//t.Log("incorrect addr")
		srv := Server{
			Addr:         ":1111",
			IdleTimeout:  5 * time.Second,
			MaxReadBytes: 1000,
		}
		go srv.ListenAndServe()

		time.Sleep(1 * time.Second) // hack to wait for server to start
		conn, err := net.Dial("tcp", "2222")
		if err != nil {
			t.Fatal(err)
		}
		for {
			source := "dsjdksjkdsjksd"
			//fmt.Scanln(&source)
			if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			for {
				buff := make([]byte, 1024)
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Print(string(buff[0:n]))
			}
			fmt.Println()
		}
	})
	t.Run("slow", func(t *testing.T) {
		//t.Parallel()
		//t.Log("slow")
		srv := Server{
			Addr:         ":3333",
			IdleTimeout:  5 * time.Second,
			MaxReadBytes: 1000,
		}
		go srv.ListenAndServe()

		time.Sleep(1 * time.Second) // hack to wait for server to start
		conn, err := net.Dial("tcp", srv.Addr)
		if err != nil {
			t.Fatal(err)
		}
		for {
			source := "dsjdksjkdsjksd"
			time.Sleep(6 * time.Second)
			//fmt.Scanln(&source)
			if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			for {
				buff := make([]byte, 1024)
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Print(string(buff[0:n]))
			}
		}
	})
	t.Run("empty content", func(t *testing.T) {
		//t.Parallel()
		//t.Log("empty content")
		srv := Server{
			Addr:         ":4444",
			IdleTimeout:  5 * time.Second,
			MaxReadBytes: 1000,
		}
		go srv.ListenAndServe()

		time.Sleep(1 * time.Second) // hack to wait for server to start
		conn, err := net.Dial("tcp", srv.Addr)
		if err != nil {
			t.Fatal(err)
		}
		for {
			source := ""
			time.Sleep(6 * time.Second)
			//fmt.Scanln(&source)
			if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			for {
				buff := make([]byte, 1024)
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Print(string(buff[0:n]))
			}
		}
	})
	t.Run("stop", func(t *testing.T) {
		//t.Parallel()
		//t.Log("stop")
		srv := Server{
			Addr:         ":5555",
			IdleTimeout:  5 * time.Second,
			MaxReadBytes: 1000,
		}
		go srv.ListenAndServe()

		time.Sleep(1 * time.Second) // hack to wait for server to start
		conn, err := net.Dial("tcp", srv.Addr)
		if err != nil {
			t.Fatal(err)
		}
		for {
			source := "dsjdksjkdsjksd"
			time.Sleep(6 * time.Second)
			//fmt.Scanln(&source)
			srv.Shutdown()
			if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			for {
				buff := make([]byte, 1024)
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Print(string(buff[0:n]))
			}
		}
	})
	t.Run("bigdata with litle MaxReadBytes", func(t *testing.T) {
		//t.Parallel()
		//t.Log("bigdata")
		srv := Server{
			Addr:         ":5555",
			IdleTimeout:  5 * time.Second,
			MaxReadBytes: 10,
		}
		go srv.ListenAndServe()

		time.Sleep(1 * time.Second) // hack to wait for server to start
		conn, err := net.Dial("tcp", srv.Addr)
		if err != nil {
			t.Fatal(err)
		}
		for {
			source := "dsjdksjksdflkjskdfjksdfjksldfjklsjdfksdfjsakfjkdlafjkogjokfdgjokdfgjdofgljfkdgjodfgjdsjksd"
			//time.Sleep(6 * time.Second)
			//fmt.Scanln(&source)
			//srv.Shutdown()
			if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			for {
				buff := make([]byte, 1024)
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Print(string(buff[0:n]))
			}
		}
	})
	t.Run("bigdata with medium MaxReadBytes", func(t *testing.T) {
		//t.Parallel()
		//t.Log("stop")
		srv := Server{
			Addr:         ":6666",
			IdleTimeout:  5 * time.Second,
			MaxReadBytes: 100,
		}
		go srv.ListenAndServe()

		time.Sleep(1 * time.Second) // hack to wait for server to start
		conn, err := net.Dial("tcp", srv.Addr)
		if err != nil {
			t.Fatal(err)
		}
		for {
			source := ""
			for i := 0; i < 1000; i++ {
				source = source + (string("s"))
			}
			time.Sleep(6 * time.Second)
			//fmt.Scanln(&source)
			//srv.Shutdown()
			if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			for {
				buff := make([]byte, 1024)
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				fmt.Print(string(buff[0:n]))
			}
		}
	})
	t.Run("slow sleeping", func(t *testing.T) {
		//t.Parallel()
		//t.Log("slow")
		srv := Server{
			Addr:         ":7777",
			IdleTimeout:  5 * time.Second,
			MaxReadBytes: 1000,
		}
		go srv.ListenAndServe()

		time.Sleep(1 * time.Second) // hack to wait for server to start
		conn, err := net.Dial("tcp", srv.Addr)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Second)
		for {
			source := "dsjdksjkdsjksd"
			time.Sleep(10 * time.Second)
			//fmt.Scanln(&source)
			if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(10 * time.Second)
			for {
				buff := make([]byte, 1024)
				n, err := conn.Read(buff)
				if err != nil {
					break
				}
				time.Sleep(10 * time.Second)
				fmt.Print(string(buff[0:n]))
			}
		}
	})
}
