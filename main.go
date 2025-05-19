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
	{
		Name:        "🥑 Тост с авокадо",
		Price:       6,
		ImageURL:    "https://example.com/avocado_toast.jpg",
		Description: "На чиабатте с кунжутом и яйцом пашот",
		Time:        "20 мин",
		Category:    "Горячие завтраки",
	},
	{
		Name:        "☕ Кофе латте",
		Price:       4,
		ImageURL:    "https://example.com/latte.jpg",
		Description: "С нежным молочным пенком и сердечком",
		Time:        "10 мин",
		Category:    "Напитки",
	},
	{
		Name:        "🍵 Чай",
		Price:       1,
		ImageURL:    "https://example.com/tea.jpg",
		Description: "Ева рекомендует",
		Time:        "10 мин",
		Category:    "Напитки",
	},
	{
		Name:        "🥐 Круассан",
		Price:       5,
		ImageURL:    "https://example.com/croissant.jpg",
		Description: "Свежая выпечка с джемом",
		Time:        "15 мин",
		Category:    "Десерты",
	},
}

var categories = []string{
	"Горячие завтраки",
	"Напитки",
	"Десерты",
}

type CartItem struct {
	Dish     Dish
	Quantity int
}

var userCarts = make(map[int64][]CartItem)
var compliments = []string{
	"Отличный выбор! 💋",
	"Ваш вкус восхитителен! 😍",
	"Добавила, солнышко! 🌸",
	"Ммм, это моё любимое! 💕",
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
	u.Timeout = 60

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
	welcomeText := `🌟 *Доброе утро, солнышко!* 🌟

Добро пожаловать в нашу летнюю кафешку! 
Сегодня на кухне шеф-повар Ева готова приготовить для тебя:

🍳 Вкуснейшие завтраки
☕ Ароматный кофе
🍓 Свежие десерты`

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

	// Создаем кнопки категорий (по 2 в ряд)
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

	// Добавляем кнопку корзины
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🛒 Корзина", "show_cart"),
	})

	msg := tgbotapi.NewMessage(chatID, "🏷 *Выберите категорию:*")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	bot.Send(msg)
}

func sendDishesByCategory(bot *tgbotapi.BotAPI, chatID int64, category string) {
	// Получаем блюда только выбранной категории
	var dishesInCategory []Dish
	for _, dish := range menu {
		if dish.Category == category {
			dishesInCategory = append(dishesInCategory, dish)
		}
	}

	if len(dishesInCategory) == 0 {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("В категории '%s' пока нет блюд 😔\nВернитесь в меню /menu", category))
		bot.Send(msg)
		return
	}

	// Отправляем заголовок категории
	headerMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🍽 *%s*\nВыберите блюдо:", category))
	headerMsg.ParseMode = "Markdown"
	bot.Send(headerMsg)

	// Отправляем каждое блюдо с фото и кнопками
	for _, dish := range dishesInCategory {
		// Создаем сообщение с фото
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(dish.ImageURL))
		photo.Caption = fmt.Sprintf(
			"*%s*\n%s\n\n⏰ %s  |  💋 %d поцелуйчиков",
			dish.Name,
			dish.Description,
			dish.Time,
			dish.Price,
		)
		photo.ParseMode = "Markdown"

		// Создаем кнопки под фото
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Добавить", fmt.Sprintf("add_%s", dish.Name)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "show_categories"),
				tgbotapi.NewInlineKeyboardButtonData("🛒 Корзина", "show_cart"),
			),
		)
		photo.ReplyMarkup = keyboard

		// Отправляем сообщение
		if _, err := bot.Send(photo); err != nil {
			log.Printf("Ошибка отправки блюда: %v", err)

			// Если не удалось отправить фото, пробуем отправить текстом
			msg := tgbotapi.NewMessage(chatID, photo.Caption)
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}

		time.Sleep(200 * time.Millisecond) // Небольшая пауза
	}
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
		sendDishesByCategory(bot, chatID, category)
	case data == "show_cart":
		showCart(bot, chatID)
	case strings.HasPrefix(data, "add_"):
		dishName := strings.TrimPrefix(data, "add_")
		addToCart(bot, chatID, dishName)
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
		// Обработка кнопки "Выполнено" у админа
		userIDStr := strings.TrimPrefix(data, "complete_")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Printf("Ошибка парсинга userID: %v", err)
			return
		}

		// Отправляем уведомление пользователю
		msg := tgbotapi.NewMessage(userID, "🎉 Ваш заказ выполнен! Приятного аппетита! 💋")
		bot.Send(msg)

		// Отвечаем админу
		answer := tgbotapi.NewCallback(callback.ID, "Заказ отмечен как выполненный")
		bot.Send(answer)

		// Удаляем сообщение с заказом у админа
		if callback.Message != nil {
			bot.Send(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID))
		}
	default:
		answer := tgbotapi.NewCallback(callback.ID, "Неизвестная команда")
		bot.Send(answer)
	}
}

func addToCart(bot *tgbotapi.BotAPI, chatID int64, dishName string) {
	for _, dish := range menu {
		if dish.Name == dishName {
			// Проверяем, есть ли уже это блюдо в корзине
			for i, item := range userCarts[chatID] {
				if item.Dish.Name == dishName {
					userCarts[chatID][i].Quantity++
					return
				}
			}
			// Если блюда нет в корзине, добавляем
			userCarts[chatID] = append(userCarts[chatID], CartItem{Dish: dish, Quantity: 1})
			return
		}
	}
}

func processOrder(bot *tgbotapi.BotAPI, chatID int64) {
	// Подготовка текста подтверждения для пользователя
	confirmationText := `🎉 *Спасибо за заказ!*

Ваш заказ принят в работу. 
Оплата: при получении 💋

Приятного аппетита и хорошего дня! 🌞

Вернуться в меню: /menu`

	// Отправляем подтверждение пользователю
	msg := tgbotapi.NewMessage(chatID, confirmationText)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	// Подготовка уведомления для админа
	var orderText strings.Builder
	orderText.WriteString(fmt.Sprintf("🔥 *Новый заказ!*\n\n"))

	// Получаем информацию о пользователе
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

	// Отправляем уведомление админу
	adminMsg := tgbotapi.NewMessage(adminChatID, orderText.String())
	adminMsg.ParseMode = "Markdown"

	// Добавляем кнопку "Выполнено" с ID пользователя
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Выполнено", fmt.Sprintf("complete_%d", chatID)),
		),
	)
	adminMsg.ReplyMarkup = keyboard

	bot.Send(adminMsg)

	// Очищаем корзину
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
