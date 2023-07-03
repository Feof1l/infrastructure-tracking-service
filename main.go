package main

import (
	"fmt"
	"log"
	"sort"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type FileConfig struct {
	FileID string
}

func main() {
	//var test1 checker = task_tcp{}
	//var test2 checker = task_http{}
	//_ = test1
	//_ = test2
	res := []data{}
	tasks := read_json_file("conf.json")
	tcp, http := load_params(tasks)
	_ = tcp
	_ = http
	//prov_http_task(http)
	//prov_tcp_task(tcp)

	/////////////////////////////////////
	telegramBotToken := "5821021624:AAFNG2VEx4h9l23ios0mb892PNgYMETGCqA"
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go Run(bot)
	//go message_func(test1, test2, "127.0.0.1", "4444")
	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, _ := bot.GetUpdatesChan(u)
	//var file_new FileConfig
	//file_new.FileID = "sdsdsfddfsfs"
	//file,error:=bot.GetFile(tgbotapi.FileConfig(file_new))
	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем

	for update := range updates {
		//for !flag {
		// универсальный ответ на любое сообщение
		reply := ""
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я - телеграм-бот, слежу за состоянием инфраструктуры"
		case "all_checks":
			var structStr string
			sort.Slice(res, func(i, j int) (less bool) {
				return res[i].id < res[j].id
			})
			//сортировка по id задач

			for _, value := range res {
				structStr = fmt.Sprint(value)
				reply = reply + structStr

			}

		case "tcp_connect":
			res_cur := []data{}
			_, _, res, res_cur = prov_tcp_task(tcp, res_cur)
			var structStr string
			for _, value := range res_cur {
				if value.type_prov == "tcp" {
					structStr = fmt.Sprint(value)
					reply = reply + structStr
				}

			}

		case "http_rec":

			res_cur := []data{}
			_, _, res, res_cur = prov_http_task(http, res_cur)
			var structStr string
			for _, value := range res_cur {
				if value.type_prov == "http" {
					structStr = fmt.Sprint(value)
					reply = reply + structStr
				}

			}
		}

		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		// отправляем
		bot.Send(msg)
	}

	//}
	//////////////////////

	/*ch := make(chan struct{})
	go func(c chan struct{}) {
		Run()
		close(c)
	}(ch)
	<-ch
	*/

	////////////////////////////////////

}
