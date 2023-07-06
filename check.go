package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	result      = []string{}
	result_data = []data{}
)

type checker interface {
	check(result []string) bool
}

func (task task_http) check(res_cur []data) (bool, data, []data, []data) {
	// res_cur для хранения результатов текущей проверки,
	// result_data - для хранения всех проверок
	var result data

	timeout, err1 := strconv.Atoi(task.answer_timeout)
	if err1 != nil {
		log.Println(err1.Error())
	}
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	time_c := time.Now()
	time_format_c := time_c.Format(time.DateTime)
	log.Println("http_prov satrted the work")
	result.id = task.id
	result.ip_name = task.ip_name
	result.ip_port = task.ip_port
	result.type_prov = "http"
	result.time = time_format_c
	result.answer_timeout = task.answer_timeout
	resp, err := client.Get(task.url)
	if err != nil {
		result.check = "Error"
		result.message = "impossible get url "
		result.status_code = 404
		if !contain(result_data, result) {
			result_data = append(result_data, result)
			res_cur = append(res_cur, result)
		}
		return false, result, result_data, res_cur
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		result.check = "Error"
		result.message = "impossible read the Body"
		result.status_code = 404
		if !contain(result_data, result) {
			result_data = append(result_data, result)
			res_cur = append(res_cur, result)
		}
		return false, result, result_data, res_cur
	}
	_ = body
	log.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("HTTP Status is in the 2xx range\n")
		result.check = "OK"
		result.status_code = resp.StatusCode
		if !contain(result_data, result) {
			result_data = append(result_data, result)
			res_cur = append(res_cur, result)
		}
		return true, result, result_data, res_cur
	} else {
		log.Println("Something has broken\n")
		result.check = "Error"
		result.message = "Something has broken"
		result.status_code = resp.StatusCode
		if !contain(result_data, result) {
			result_data = append(result_data, result)
			res_cur = append(res_cur, result)
		}
		return false, result, result_data, res_cur
	}

}

func (task task_tcp) check(res_cur []data) (bool, data, []data, []data) {
	// res_cur для хранения результатов текущей проверки,
	// result_data - для хранения всех проверок
	var result data
	timeout, err1 := strconv.Atoi(task.answer_timeout)
	if err1 != nil {
		log.Println(err1.Error())
	}
	time_c := time.Now()
	time_format_c := time_c.Format(time.DateTime)

	result.id = task.id
	result.ip_name = task.ip_name
	result.ip_port = task.ip_port
	result.type_prov = "tcp"
	result.time = time_format_c
	result.answer_timeout = task.answer_timeout
	log.Println("tcp_prov satrted the work")
	// преобразование таймаута к типу time.Duration в миллисекундах
	timeout_format := time.Duration.Milliseconds(time.Duration(timeout))
	timeout_time := time.Duration(timeout_format)
	_, err := net.DialTimeout("tcp", task.ip_name+":"+task.ip_port, timeout_time)
	if err != nil {
		log.Printf("%s %s %s", task.ip_name, "not responding on port ", task.ip_port, err.Error())
		log.Println()
		result.check = "Error"
		result.message = "not responding on current port "
		result.status_code = 404
		if !contain(result_data, result) {
			result_data = append(result_data, result)
			res_cur = append(res_cur, result)
		}

		return false, result, result_data, res_cur
	} else {
		log.Printf("%s %s %s\n", task.ip_name, "responding on port:", task.ip_port)
		log.Println()
		result.check = "OK "
		result.status_code = 200
		if !contain(result_data, result) {
			result_data = append(result_data, result)
			res_cur = append(res_cur, result)
		}

		return true, result, result_data, res_cur
	}

}

type data struct {
	id             string
	ip_name        string
	ip_port        string
	type_prov      string
	check          string
	answer_timeout string
	message        string
	time           string
	status_code    int
}

func (res data) String() string {
	return fmt.Sprintf("{ id: %s,type: %s,result: %s,message: %s,time: %s, status code: %d }",
		res.id, res.type_prov, res.check, res.message, res.time, res.status_code)
}
func prov_tcp_task(params_tcp []task_tcp, res []data) (bool, data, []data, []data) {
	var prov bool
	var data_t data
	res_cur := []data{}
	for _, value := range params_tcp {
		//log.Println(value.check())
		prov, data_t, res, res_cur = value.check(res_cur)
		log.Println(prov, data_t)
	}
	log.Println(res)
	return prov, data_t, res, res_cur
}
func prov_http_task(params_http []task_http, res []data) (bool, data, []data, []data) {
	var prov bool
	var data_t data
	res_cur := []data{}
	for _, value := range params_http {
		//log.Println(value)
		//log.Println(value.check())
		prov, data_t, res, res_cur = value.check(res_cur)
		log.Println(prov, data_t)

	}
	return prov, data_t, res, res_cur
}

// функция для проверки, находится ли элемент в слайсе
// используется в функции check для того, чтобы не добавлять одни и те же результаты проверки в слайс
// при выполнении проверки
func contain(slice []data, x data) bool {
	for _, value := range slice {
		if value == x {
			return true
		}
	}
	return false
}
