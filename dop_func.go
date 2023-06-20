package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	DateTime = "2006-01-02 15:04:05"
	//final_time = "2023-06-19 15:42:00"
)

func final_time_time() time.Time {
	t, _ := time.Parse(DateTime, final_time())
	t = t.Add(time.Hour * -3)
	return t
}
func unslice(old string) string {
	new := make([]byte, len(old))
	copy(new, old)
	return string(old)
}
func final_time() string {
	time_c := time.Now()
	time_str := time_c.String()
	time_str = strings.TrimSpace(time_str)
	slice := strings.Split(time_str, "")
	slice = slice[:10]
	var time_fin = []string{"17:35:00"}
	c := append(slice, time_fin...)
	str2 := strings.Join(c, " ")
	str2 = strings.TrimSpace(str2)
	dlina := 0
	s := ""
	str2 = strings.Trim(str2, "\n\r")
	r := []rune(str2)
	for _, i := range str2 {
		_ = i
		dlina++
	}
	for i := 0; i < 19; i++ {
		s = s + strings.TrimSpace(string((r[i])))
	}
	s = s + string(r[19])
	for i := 20; i < dlina; i++ {
		s = s + (string((r[i])))
	}
	return s
}

type checker interface {
	check(host, port string) bool
}

type http_prov struct{}

type tcp_prov struct{}

func (c http_prov) check(host, port string) bool {
	log.Println("http_prov satrted the work")
	url := ("http://" + host + ":" + port)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	_ = body
	log.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("HTTP Status is in the 2xx range\n")
		return true
	} else {
		log.Println("Something has broken\n")
		return false
	}

}

func (a tcp_prov) check(host, port string) bool {
	log.Println("tcp_prov satrted the work")
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", host+":"+port, timeout)
	if err != nil {
		log.Printf("%s %s %s", host, "not responding on port ", port, err.Error())
		log.Println()
		return false
	} else {
		log.Printf("%s %s %s\n", host, "responding on port:", port)
		log.Println()
		return true
	}
}
func message_func(test1, test2 checker, host, port string) (string, string) {
	msg1 := "Сообщений о проверке пока нет"
	msg2 := "Сообщений о проверке пока нет"
	if test1.check(host, port) {
		msg1 = "Проверка прошла успешно, tcp connect работает, на указанном порту открыт сервер\n"
	} else {
		msg1 = "Ошибка, tcp connect не рабоатет, на указанном порту не удается обнаружить сервер\n"
	}
	if test2.check(host, port) {
		msg2 = "Проверка прошла успешно, http запрос успешно дошел и был обработан, сервер обрабатывает сообщения в штатном режиме\n"
	} else {
		msg2 = "Ошибка, функция обработчик не работает или что-то пошло не так\n"
	}
	return msg1, msg2

}
func Run(bot *tgbotapi.BotAPI) {
	i := 0
	log.Println("Alarm will ring in ", time.Until(final_time_time()))
	tg := time.Until(final_time_time())
	ticker := time.AfterFunc(tg, func() {
		//i = i.Add(1 * time.Second)
		i++
	})

	for {
		select {
		case <-ticker.C:
			log.Println("nsmei")
		case <-time.After((tg)):
			if i == 1 {
				ticker.Reset(1 * time.Second)
				continue
			}
			goto BRK
		}
	BRK:
		ticker.Stop()
		alarm := "Звенит будильник.Пожалуйста, проверьте статус серверов"
		bot.Send(tgbotapi.NewMessage(653924346, alarm))
		log.Println(time.Now().Format("2006/01/02 15:04:05.999999999"), " ALARM \n")
		break
	}

}
