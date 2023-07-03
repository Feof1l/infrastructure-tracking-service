package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr" //ключ для значения адреса HTTP-сервера в контексте

func main() {
	//mux1 := http.NewServeMux()
	//mux1.HandleFunc("/", getRoot)
	//mux1.HandleFunc("/hello", getHello)

	//mux2 := http.NewServeMux()
	//mux2.HandleFunc("/", getRoot)
	//mux2.HandleFunc("/hello", getHello)

	ctx := context.Background()
	serverOne := &http.Server{ // первый сервер
		Addr: ":3333",
		//Handler: mux1,
		BaseContext: func(l net.Listener) context.Context { // функция изменения контекста
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{ //второй сервер
		Addr: ":4444",
		//Handler: mux2,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go erorrs_er(serverOne)
	go erorrs_er(serverTwo)

	/*go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) { //завершение работы сервера
			fmt.Printf("server one closed\n")
		} else if err != nil { //иная ошибка
			fmt.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()

	go func() {
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server two closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server two: %s\n", err)
		}
		cancelCtx()
	}()*/

	<-ctx.Done() // читаем из каналом перед возвратом в main
	// продолжаем рабоут пока не завершится любая из подпрогармм сервера
	//и не будет вызвано аннулирвоание контекста
	//как только будет вызов cancelC

}
func erorrs_er(server *http.Server) {
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server two closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}

// функции обработчики
func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // создаем контекст для отслеживания серверов

	first := r.URL.Query().Get("first")
	second := r.URL.Query().Get("second")

	body, err := ioutil.ReadAll((r.Body))
	if err != nil {
		fmt.Printf("i can not read body: %s\n", err)
	}

	fmt.Printf("%s: got / request.first=%s, second=%s, body:\n%s\n", ctx.Value(keyServerAddr), first, second, body)
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-field", "myName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, fmt.Sprintf("Hello, %s!\n", myName))
}
