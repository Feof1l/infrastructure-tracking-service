package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	DateTime = "2006-01-02 15:04:05"
	//final_time = "2023-06-19 15:42:00"
)

var (
	result = []string{}
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
	check(result []string) bool
}

//type http_prov struct{}

//type tcp_prov struct{}

func (task task_http) check() (bool, []string) {

	timeout, err1 := strconv.Atoi(task.answer_timeout)
	if err1 != nil {
		log.Println(err1.Error())
	}
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Nanosecond,
	}

	time_c := time.Now()
	time_format_c := time_c.Format(time.DateTime)
	msg := "Номер задачи: " + task.id + " тип проверки: http " + "текущая дата: " + time_format_c + " результат проверки: "

	log.Println("http_prov satrted the work")

	resp, err := client.Get(task.url)
	if err != nil {
		log.Println(err)
		msg = msg + "Error " + "комментарии: " + "impossible get url " + task.url + "\n"
		result = append(result, msg)
		return false, result
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		msg = msg + "Error " + "комментарии: " + "impossible read the Body" + "/n"
		result = append(result, msg)
		return false, result
	}
	_ = body
	log.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("HTTP Status is in the 2xx range\n")
		msg = msg + "OK " + "комментарии: " + "StatusCode: " + strconv.Itoa(resp.StatusCode) + "\n"
		result = append(result, msg)
		return true, result
	} else {
		log.Println("Something has broken\n")
		msg = msg + "Error\n"
		result = append(result, msg)
		return false, result
	}

}

func (task task_tcp) check() (bool, []string) {
	timeout, err1 := strconv.Atoi(task.answer_timeout)
	if err1 != nil {
		log.Println(err1.Error())
	}
	time_c := time.Now()
	time_format_c := time_c.Format(time.DateTime)
	msg := "Номер задачи: " + task.id + " тип проверки: tcp " + "текущая дата: " + time_format_c + " результат проверки: "
	log.Println("tcp_prov satrted the work")
	timeout_time := time.Duration(timeout)
	_, err := net.DialTimeout("tcp", task.ip_name+":"+task.ip_port, timeout_time)
	if err != nil {
		log.Printf("%s %s %s", task.ip_name, "not responding on port ", task.ip_port, err.Error())
		log.Println()
		msg = msg + "Error " + "комментарии:" + " not responding on current port " + task.ip_port + "\n"
		result = append(result, msg)
		return false, result
	} else {
		log.Printf("%s %s %s\n", task.ip_name, "responding on port:", task.ip_port)
		log.Println()
		msg = msg + "OK\n"
		result = append(result, msg)
		return true, result
	}

}
func (task task_new_type) check() (bool, []string) {

	timeout, err1 := strconv.Atoi(task.answer_timeout)
	if err1 != nil {
		log.Println(err1.Error())
	}
	_=timeout

	time_c := time.Now()
	time_format_c := time_c.Format(time.DateTime)
	msg := "Номер задачи: " + task.id + " тип проверки: http " + "текущая дата: " + time_format_c + " результат проверки: "

	log.Println("http_prov satrted the work")

	resp, err := http.Get(task.url)
	if err != nil {
		log.Println(err)
		msg = msg + "Error " + "комментарии: " + "impossible get url " + task.url + "\n"
		result = append(result, msg)
		return false, result
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		msg = msg + "Error " + "комментарии: " + "impossible read the Body" + "/n"
		result = append(result, msg)
		return false, result
	}
	_ = body
	log.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("HTTP Status is in the 2xx range\n")
		msg = msg + "OK " + "комментарии: " + "StatusCode: " + strconv.Itoa(resp.StatusCode) + "\n"
		result = append(result, msg)
		return true, result
	} else {
		log.Println("Something has broken\n")
		msg = msg + "Error\n"
		result = append(result, msg)
		return false, result
	}

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
func read_json_file(file_name string) map[string]map[string]string {
	content, err := ioutil.ReadFile(file_name) // чтение файла
	if err != nil {
		log.Fatal("ошибка открытия файла: ", err)
	}

	var list_conf map[string]map[string]string
	// данные с файла считываемв мапу мап - ключ - id задачи
	// значение - мапа строк с параметрами задачи
	err = json.Unmarshal(content, &list_conf)
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
type task_new_type struct{
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
		//также сохраняем id
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

/*
func convert_params(mas []task_tcp){

		for _, value := range mas{


		}
		i, err := strconv.Atoi(s)
	    if err != nil {
	        // ... handle error
	        panic(err)
	    }
	}
*/
func prov_tcp_task(params_tcp []task_tcp) {
	for _, value := range params_tcp {
		log.Println(value.check())
	}
}
func prov_http_task(params_http []task_http) {
	for _, value := range params_http {
		log.Println(value.check())
	}
}
