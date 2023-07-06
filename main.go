package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	res := []data{}
	var tasks map[string]map[string]string
	var tcp []task_tcp
	var http_t []task_http

	/////////////////////////////////////
	telegramBotToken := "5821021624:AAFNG2VEx4h9l23ios0mb892PNgYMETGCqA"
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	go Run(bot)
	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, _ := bot.GetUpdatesChan(u)
	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем

	for update := range updates {
		// универсальный ответ на любое сообщение
		reply := ""
		if update.Message == nil {
			continue
		}
		//загрузка файла
		if update.Message.Document != nil {
			fileID := update.Message.Document.FileID
			log.Printf("File %s downloaded", update.Message.Document.FileName)
			mg := tgbotapi.NewMessage(update.Message.Chat.ID, "File downloaded")
			bot.Send(mg)
			fileURL, err := bot.GetFileDirectURL(fileID)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			response, err := http.Get(fileURL)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			defer response.Body.Close()
			fileBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			// читаем фйл и загружаем оттуда необходимые параметры
			tasks = read_json_file(fileBytes)
			tcp, http_t = load_params(tasks)
			//удаляем файл
			err = os.Remove(update.Message.Document.FileName)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			log.Printf("File %s was removed", update.Message.Document.FileName)
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
			_, _, res, res_cur = prov_http_task(http_t, res_cur)
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

	////////////////////////////////////

}
