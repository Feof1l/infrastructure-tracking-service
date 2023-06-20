package main

import (
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	var test1 checker = tcp_prov{}
	var test2 checker = http_prov{}
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
	updates := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	//flag := Run()
	//go func(msg_http, msg_tcp string) {
	for update := range updates {
		//for !flag {
		// универсальный ответ на любое сообщение
		reply := "Не знаю, что сказать"
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
		case "tcp_connect":
			reply, _ = message_func(test1, test2, "127.0.0.1", "4444")
		case "http_rec":
			_, reply = message_func(test1, test2, "127.0.0.1", "4444")
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

}
