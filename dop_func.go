package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	alarm_time = final_time_time()
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
func final_time() string { // формирование строки врмеени, по которому будет срабатывать будтльник
	time_c := time.Now()
	// берем текущее время и дату, вытаскиваем оттуда только дату
	// и конкатенируем с тем врменем, по которому должен сработать будильник
	time_str := time_c.String()
	time_str = strings.TrimSpace(time_str)
	slice := strings.Split(time_str, "")
	slice = slice[:10]
	var time_fin = []string{"20:45:00"}
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

//type http_prov struct{}

//type tcp_prov struct{}

func Run(bot *tgbotapi.BotAPI) {
	i := 0

	log.Println("Alarm will ring in ", time.Until(alarm_time))
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
		alarm_time = alarm_time.AddDate(0, 0, 1)
		log.Println(alarm_time)
		alarm := "Звенит будильник.Пожалуйста, проверьте статус серверов"
		bot.Send(tgbotapi.NewMessage(653924346, alarm))
		log.Println(time.Now().Format("2006/01/02 15:04:05.999999999"), " ALARM \n")
		//final_time()
		log.Println("Next alarm will ring in ", time.Until(alarm_time))
		break
	}

}
func read_json_file(content []byte) map[string]map[string]string {
	var list_conf map[string]map[string]string
	// данные с файла считываемв мапу мап - ключ - id задачи
	// значение - мапа строк с параметрами задачи
	err := json.Unmarshal(content, &list_conf)
	if err != nil {
		log.Fatal("Ошибка во время выполнения Unmarshal(): ", err)
	}
	//log.Println(list_conf["2"]["sdsd"])
	return list_conf
}

// структуры для храненения параметров для разных проверок
type task_tcp struct {
	id             string
	ip_name        string
	ip_port        string
	answer_timeout string
}
type task_http struct {
	id             string
	ip_name        string
	ip_port        string
	url            string
	answer_timeout string
}
type task_new_type struct {
	id             string
	ip_name        string
	ip_port        string
	url            string
	answer_timeout string
}

func load_params(list_conf map[string]map[string]string) ([]task_tcp, []task_http) {

	tasks_tcp := []task_tcp{} // создаем слайсы для хранения стурктур задач
	tasks_http := []task_http{}
	i := 0
	//c := append(slice, time_fin...)
	for Id, map_params_task := range list_conf {
		//проходимся по мапе
		// проверяем тип, если он tcp, записываем параметры из мапы в соот структуру
		//также сохраняем ids
		i++
		if map_params_task["type"] == "tcp_connect" {
			tasks_tcp = append(tasks_tcp, task_tcp{id: Id, ip_name: map_params_task["ip_name"],
				ip_port:        map_params_task["ip_port"],
				answer_timeout: map_params_task["answer_timeout"]})
		} else if map_params_task["type"] == "http_rec" {
			tasks_http = append(tasks_http, task_http{id: Id, ip_name: map_params_task["ip_name"],
				ip_port:        map_params_task["ip_port"],
				url:            map_params_task["url"],
				answer_timeout: map_params_task["answer_timeout"]})
		}

	}
	//log.Println(tasks_tcp)
	//log.Println(tasks_http)
	return tasks_tcp, tasks_http
}
