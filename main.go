package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func createBotClient(token string) *bot.Bot {
	opts := []bot.Option{}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	return b
}

func routeHandler(fn func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Send any text message to the bot after the bot has been started

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":3000", "--addr [host]:port")
	flag.Parse()

	log.SetLevel(log.DebugLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ev, err := LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	robot := createBotClient(ev.TelegramBotToken)
	_ = robot

	robot.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: fmt.Sprintf("%s/receive-from-telegram", ev.PublicWebhookURL),
	})

	go func() {
		log.Infof("starting telegram webhook bot")
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok hi"))
	})

	mux.HandleFunc("/push/text", routeHandler(func(w http.ResponseWriter, r *http.Request) error {
		if _, err := robot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: ev.TelegramChatID,
			Text:   r.URL.Query().Get("message"),
		}); err != nil {
			return err
		}
		return nil
	}))

	mux.HandleFunc("/receive-from-telegram", routeHandler(func(w http.ResponseWriter, r *http.Request) error {
		log.Debugf("[route] %s", r.RequestURI)
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		update := &models.Update{}
		if err := json.Unmarshal(b, update); err != nil {
			return err
		}

		if _, err := robot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.ChannelPost.Chat.ID,
			Text:   "hello, welcome to kloudlite alerts",
		}); err != nil {
			return err
		}
		w.Write([]byte("ok"))
		return nil
	}))

	log.Infof("starting http server on (%s)", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
