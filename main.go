package main

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Dish struct {
	Name  string
	Price int
}

var menu = []Dish{
	{"Яйца пашот", 3},
	{"Круассан", 5},
	{"Кофе латте", 4},
	{"Тост с авокадо", 6},
}

var userCarts = make(map[int64]map[string]int) // userID -> map[dishName]count

const adminChatID int64 = 401424871 // тут твой ID для получения заказов

func main() {
	bot, err := tgbotapi.NewBotAPI("7987892591:AAFa2TQKrUypxuBTh5Soi_4OKEiMohsv0Kc")
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // обработка команд /start и т.п.
			if update.Message.Text == "/start" {
				sendMenu(bot, update.Message.Chat.ID)
			}
		} else if update.CallbackQuery != nil {
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

func sendMenu(bot *tgbotapi.BotAPI, chatID int64) {
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, dish := range menu {
		btnText := dish.Name + " - " + strconv.Itoa(dish.Price) + " поцелуев"
		callbackData := "add_" + dish.Name
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(btnText, callbackData),
		))
	}
	// Добавим кнопку оформить заказ
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Оформить заказ", "checkout"),
	))

	msg := tgbotapi.NewMessage(chatID, "Выберите блюда для корзины:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	bot.Send(msg)
}

func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	data := callback.Data

	if strings.HasPrefix(data, "add_") {
		dishName := strings.TrimPrefix(data, "add_")

		if userCarts[userID] == nil {
			userCarts[userID] = make(map[string]int)
		}
		userCarts[userID][dishName]++

		answer := "Добавлено в корзину: " + dishName
		callbackConfig := tgbotapi.NewCallback(callback.ID, answer)
		bot.Send(callbackConfig)

	} else if data == "checkout" {
		cart := userCarts[userID]
		if len(cart) == 0 {
			callbackConfig := tgbotapi.NewCallback(callback.ID, "Корзина пуста")
			bot.Send(callbackConfig)
			return
		}

		orderText := "Новый заказ от @" + callback.From.UserName + ":\n"
		total := 0
		for dish, count := range cart {
			price := 0
			for _, d := range menu {
				if d.Name == dish {
					price = d.Price
					break
				}
			}
			orderText += dish + " x" + strconv.Itoa(count) + " = " + strconv.Itoa(price*count) + " поцелуев\n"
			total += price * count
		}
		orderText += "\nИтого: " + strconv.Itoa(total) + " поцелуев"

		msg := tgbotapi.NewMessage(adminChatID, orderText)
		bot.Send(msg)

		delete(userCarts, userID)

		callbackConfig := tgbotapi.NewCallback(callback.ID, "Заказ отправлен!")
		bot.Send(callbackConfig)
	}
}
