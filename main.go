package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	botToken    = "7951456204:AAF_39lpVE3niamHcvV-6LLb51nD758IlWs"
	adminChatID = int64(401424871)
)

type Dish struct {
	Name        string
	Price       int
	ImageURL    string
	Description string
	Time        string
	Category    string
}

var menu = []Dish{
	// Завтраки
	{
		Name:        "🍳 Яичница",
		Price:       3,
		ImageURL:    "https://sitandeat.ru/upload/resize_cache/iblock/3e8/1200_1200_2/mly7arvbo06n3xwof1h2ybkybs3b3ik7.jpg",
		Description: "С точечками кетчупа сверху :)",
		Time:        "20 мин",
		Category:    "Завтраки",
	},
	{
		Name:        "🥚 Омлет",
		Price:       4,
		ImageURL:    "https://img.delo-vcusa.ru/2019/09/omlet-s-sirom.jpg",
		Description: "Нежный омлет с молоком",
		Time:        "20 мин",
		Category:    "Завтраки",
	},
	{
		Name:        "🥣 Манная каша",
		Price:       5,
		ImageURL:    "https://static.1000.menu/img/content-v2/70/64/23581/mannaya-kasha-na-moloke-bez-komochkov_1601056376_8_max.jpg",
		Description: "Нежная каша с маслом",
		Time:        "25 мин",
		Category:    "Завтраки",
	},
	{
		Name:        "🥪 Бутербродик",
		Price:       3,
		ImageURL:    "https://cdn.botanichka.ru/wp-content/uploads/2024/01/goryachij-buterbrod-s-syrom-i-pomidorami-0.jpg",
		Description: "С тем, что найдется на нашей кухне :)",
		Time:        "10 мин",
		Category:    "Завтраки",
	},

	// Напитки
	{
		Name:        "🍵 Чай с сахаром",
		Price:       2,
		ImageURL:    "https://main-cdn.sbermegamarket.ru/big2/hlr-system/758/393/980/330/189/100023256907b1.jpg",
		Description: "Вкусный чай с сахаром",
		Time:        "5 мин",
		Category:    "Напитки",
	},
	{
		Name:        "🍵 Чай без сахара",
		Price:       1,
		ImageURL:    "https://main-cdn.sbermegamarket.ru/big2/hlr-system/758/393/980/330/189/100023256907b1.jpg",
		Description: "Чай без сахара, чтобы не было прыщиков",
		Time:        "5 мин",
		Category:    "Напитки",
	},
	{
		Name:        "💧 Водичка",
		Price:       1,
		ImageURL:    "https://stav-history.ru/wp-content/uploads/2019/03/85874599.jpg",
		Description: "Стаканчик прохладной воды",
		Time:        "2 мин",
		Category:    "Напитки",
	},

	// Хлеб
	{
		Name:        "🍞 Белый хлеб с маслом",
		Price:       1,
		ImageURL:    "https://www.m24.ru/b/d/nBkSUhL2hFghnMeyIr6BrNOp2Z318Ji-mijFnuWR9mOBdDebBizCnTY8qdJf6ReJ58vU9meMMok3Ee2nhSR6ISeO9G1N_wjJ=GkckcD-gTxuFJl0F8hqLcw.jpg",
		Description: "Свежий белый хлеб с маслом",
		Time:        "3 мин",
		Category:    "Хлеб",
	},
	{
		Name:        "🥖 Хлеб 'Тарту' с маслом",
		Price:       1,
		ImageURL:    "https://www.tablicakalorijnosti.ru/file/image/foodstuff/22492b9099f44aa99bc7421a015c0796/6c628404cd014e4abfed08b68d96fdd7",
		Description: "Ароматный хлеб 'Тарту' с маслом",
		Time:        "3 мин",
		Category:    "Хлеб",
	},
	{
		Name:        "🍞 Поджаренный белый хлеб с маслом",
		Price:       1,
		ImageURL:    "https://www.m24.ru/b/d/nBkSUhL2hFghnMeyIr6BrNOp2Z318Ji-mijFnuWR9mOBdDebBizCnTY8qdJf6ReJ58vU9meMMok3Ee2nhSR6ISeO9G1N_wjJ=GkckcD-gTxuFJl0F8hqLcw.jpg",
		Description: "Хрустящий поджаренный хлеб с маслом",
		Time:        "5 мин",
		Category:    "Хлеб",
	},
	{
		Name:        "🥖 Поджаренный 'Тарту' с маслом",
		Price:       1,
		ImageURL:    "https://www.tablicakalorijnosti.ru/file/image/foodstuff/22492b9099f44aa99bc7421a015c0796/6c628404cd014e4abfed08b68d96fdd7",
		Description: "Поджаренный хлеб 'Тарту' с маслом",
		Time:        "5 мин",
		Category:    "Хлеб",
	},
	{
		Name:        "🍞 Поджаренный тост с маслом",
		Price:       1,
		ImageURL:    "https://img.freepik.com/premium-photo/delicious-crispy-toast-with-butter-isolated-white_495423-50544.jpg",
		Description: "Золотистый тост с маслом",
		Time:        "5 мин",
		Category:    "Хлеб",
	},
}

var categories = []string{
	"Завтраки",
	"Напитки",
	"Хлеб",
}

type CartItem struct {
	Dish     Dish
	Quantity int
}

var userCarts = make(map[int64][]CartItem)
var compliments = []string{
	"Отличный выбор! 💋",
	"Добавил, зайка! 🌸",
	"Ммм, вкуснятина! 😋",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		sendWelcome(bot, msg.Chat.ID)
	case "/menu":
		sendCategoryMenu(bot, msg.Chat.ID)
	case "/cart":
		showCart(bot, msg.Chat.ID)
	default:
		sendDefaultResponse(bot, msg.Chat.ID)
	}
}

func sendWelcome(bot *tgbotapi.BotAPI, chatID int64) {
	welcomeText := `🌟 *Доброе утро, Анечка!* 🌟

Добро пожаловать в Ресторан Гюсто! 
Сегодня на кухне шеф-повар Влад готов приготовить для тебя:

🍳 Вкуснейшие завтраки
☕ Ароматный чай
🍞 Хлебушек`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍽️ Открыть меню", "show_categories"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, welcomeText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func sendCategoryMenu(bot *tgbotapi.BotAPI, chatID int64) {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопки категорий (по 2 в ряд)
	for i := 0; i < len(categories); i += 2 {
		var row []tgbotapi.InlineKeyboardButton
		if i < len(categories) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(categories[i], "category_"+categories[i]))
		}
		if i+1 < len(categories) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(categories[i+1], "category_"+categories[i+1]))
		}
		rows = append(rows, row)
	}

	// Кнопка корзины
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🛒 Корзина", "show_cart"),
	})

	msg := tgbotapi.NewMessage(chatID, "🏷 *Выберите категорию:*")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	bot.Send(msg)
}

func sendDishesMenu(bot *tgbotapi.BotAPI, chatID int64, category string) {
	var dishesInCategory []int // индексы блюд в menu
	for i, dish := range menu {
		if dish.Category == category {
			dishesInCategory = append(dishesInCategory, i)
		}
	}

	if len(dishesInCategory) == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("В категории '%s' пока нет блюд 😔", category))
		bot.Send(msg)
		return
	}

	// Создаем кнопки для блюд (по 2 в ряд)
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(dishesInCategory); i += 2 {
		var row []tgbotapi.InlineKeyboardButton
		idx := dishesInCategory[i]
		dish := menu[idx]
		btnText := fmt.Sprintf("%s - 💋%d", dish.Name, dish.Price)
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(btnText, fmt.Sprintf("dish_%d", idx)))
		if i+1 < len(dishesInCategory) {
			idx2 := dishesInCategory[i+1]
			dish2 := menu[idx2]
			btnText2 := fmt.Sprintf("%s - 💋%d", dish2.Name, dish2.Price)
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(btnText2, fmt.Sprintf("dish_%d", idx2)))
		}
		rows = append(rows, row)
	}

	// Кнопки навигации
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "show_categories"),
		tgbotapi.NewInlineKeyboardButtonData("🛒 Корзина", "show_cart"),
	})

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🍽 *%s*\nВыберите блюдо:", category))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	bot.Send(msg)
}

func showDishDetails(bot *tgbotapi.BotAPI, chatID int64, dishIdx int) {
	if dishIdx < 0 || dishIdx >= len(menu) {
		msg := tgbotapi.NewMessage(chatID, "Ошибка: блюдо не найдено.")
		bot.Send(msg)
		return
	}
	selectedDish := menu[dishIdx]

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(selectedDish.ImageURL))
	photo.Caption = fmt.Sprintf(
		"*%s*\n%s\n\n⏰ %s  |  💋 %d поцелуйчиков",
		selectedDish.Name,
		selectedDish.Description,
		selectedDish.Time,
		selectedDish.Price,
	)
	photo.ParseMode = "Markdown"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить в корзину", fmt.Sprintf("add_%d", dishIdx)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к меню", "category_"+selectedDish.Category),
			tgbotapi.NewInlineKeyboardButtonData("🛒 Корзина", "show_cart"),
		),
	)
	photo.ReplyMarkup = keyboard

	bot.Send(photo)
}

func showCart(bot *tgbotapi.BotAPI, chatID int64) {
	cart := userCarts[chatID]
	if len(cart) == 0 {
		msg := tgbotapi.NewMessage(chatID, "🛒 Ваша корзина пока пуста...")
		bot.Send(msg)
		return
	}

	var total int
	var cartText strings.Builder
	cartText.WriteString("🌟 *Ваша корзина*\n\n")

	for i, item := range cart {
		itemTotal := item.Dish.Price * item.Quantity
		cartText.WriteString(
			fmt.Sprintf("%d. *%s*\n   x%d  |  💋 %d\n\n",
				i+1,
				item.Dish.Name,
				item.Quantity,
				itemTotal,
			),
		)
		total += itemTotal
	}

	// Расчет времени приготовления
	maxTime := 0
	for _, item := range cart {
		if t, err := strconv.Atoi(strings.TrimSuffix(item.Dish.Time, " мин")); err == nil {
			if t > maxTime {
				maxTime = t
			}
		}
	}
	totalTime := maxTime + 10

	cartText.WriteString(fmt.Sprintf("💋 *Итого:* %d поцелуйчиков\n⏰ *Время приготовления:* %d мин", total, totalTime))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Оформить заказ", "checkout"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Очистить корзину", "clear_cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить ещё", "show_categories"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, cartText.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	switch {
	case data == "show_categories":
		bot.Send(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID))
		sendCategoryMenu(bot, chatID)
	case strings.HasPrefix(data, "category_"):
		category := strings.TrimPrefix(data, "category_")
		bot.Send(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID))
		sendDishesMenu(bot, chatID, category)
	case strings.HasPrefix(data, "dish_"):
		dishIdxStr := strings.TrimPrefix(data, "dish_")
		dishIdx, err := strconv.Atoi(dishIdxStr)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка: блюдо не найдено.")
			bot.Send(msg)
			return
		}
		bot.Send(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID))
		showDishDetails(bot, chatID, dishIdx)
	case data == "show_cart":
		showCart(bot, chatID)
	case strings.HasPrefix(data, "add_"):
		dishIdxStr := strings.TrimPrefix(data, "add_")
		dishIdx, err := strconv.Atoi(dishIdxStr)
		if err != nil {
			return
		}
		addToCart(bot, chatID, dishIdx)
		answer := tgbotapi.NewCallback(callback.ID, compliments[rand.Intn(len(compliments))])
		bot.Send(answer)
	case data == "clear_cart":
		delete(userCarts, chatID)
		answer := tgbotapi.NewCallback(callback.ID, "Корзина очищена 🌸")
		bot.Send(answer)
		showCart(bot, chatID)
	case data == "checkout":
		processOrder(bot, chatID)
		answer := tgbotapi.NewCallback(callback.ID, "Заказ оформлен! 💌")
		bot.Send(answer)
	case strings.HasPrefix(data, "complete_"):
		userIDStr := strings.TrimPrefix(data, "complete_")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Printf("Ошибка парсинга userID: %v", err)
			return
		}

		msg := tgbotapi.NewMessage(userID, "🎉 Ваш заказ выполнен! Приятного аппетита! 💋")
		bot.Send(msg)

		answer := tgbotapi.NewCallback(callback.ID, "Заказ отмечен как выполненный")
		bot.Send(answer)

		if callback.Message != nil {
			bot.Send(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID))
		}
	default:
		answer := tgbotapi.NewCallback(callback.ID, "Неизвестная команда")
		bot.Send(answer)
	}
}

func addToCart(bot *tgbotapi.BotAPI, chatID int64, dishIdx int) {
	if dishIdx < 0 || dishIdx >= len(menu) {
		return
	}
	dish := menu[dishIdx]
	for i, item := range userCarts[chatID] {
		if item.Dish.Name == dish.Name {
			userCarts[chatID][i].Quantity++
			return
		}
	}
	userCarts[chatID] = append(userCarts[chatID], CartItem{Dish: dish, Quantity: 1})
}

func processOrder(bot *tgbotapi.BotAPI, chatID int64) {
	confirmationText := `🎉 *Спасибо за заказ!*

Ваш заказ принят в работу. 
Оплата: при получении 💋

Приятного аппетита и хорошего дня! 🌞`

	// Отправляем фото подтверждения
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL("https://www.cdn.memify.ru/media/chuUwsZJgwFASPiDQrBXFg/20240927/5454042673154484697.jpg"))
	photo.Caption = confirmationText
	photo.ParseMode = "Markdown"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍽 Вернуться в меню", "show_categories"),
		),
	)
	photo.ReplyMarkup = keyboard

	bot.Send(photo)

	// Отправляем уведомление админу
	var orderText strings.Builder
	orderText.WriteString(fmt.Sprintf("🔥 *Новый заказ!*\n\n"))

	chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	})

	if err != nil {
		log.Printf("Ошибка получения информации о чате: %v", err)
		orderText.WriteString(fmt.Sprintf("От: пользователь (ID: %d)\n\n", chatID))
	} else {
		orderText.WriteString(fmt.Sprintf("От: @%s (ID: %d)\n\n", chat.UserName, chatID))
	}

	total := 0
	for _, item := range userCarts[chatID] {
		itemTotal := item.Dish.Price * item.Quantity
		orderText.WriteString(fmt.Sprintf("• %s x%d = 💋%d\n", item.Dish.Name, item.Quantity, itemTotal))
		total += itemTotal
	}

	orderText.WriteString(fmt.Sprintf("\n💋 *Итого:* %d поцелуйчиков", total))

	adminMsg := tgbotapi.NewMessage(adminChatID, orderText.String())
	adminMsg.ParseMode = "Markdown"

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Выполнено", fmt.Sprintf("complete_%d", chatID)),
		),
	)
	adminMsg.ReplyMarkup = keyboard

	bot.Send(adminMsg)

	delete(userCarts, chatID)
}

func sendDefaultResponse(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, `Я не понимаю эту команду 😔

Вот что я умею:
/start - начать работу
/menu - показать меню
/cart - показать корзину`)
	bot.Send(msg)
}
