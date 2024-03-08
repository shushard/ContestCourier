package telegram

import (
	"fmt"
	"io"

	"tester/internal/metrics"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type alertInfo struct {
	service string
	message string
	file    fileInfo
}

type fileInfo struct {
	name       string
	readCloser io.ReadCloser
	bytes      []byte
	path       string
}

func (b *Bot) alert(service, message string, file fileInfo) {
	alert := alertInfo{
		service: service,
		message: message,
		file:    file,
	}

	select {
	case <-b.shutdownCtx.Done():
		return
	case b.alertsChan <- alert:
	}
}

func (b *Bot) processAlerts() {
	defer b.wg.Done()

	for alert := range b.alertsChan {
		for _, chatID := range b.config.Subscribers {
			b.sendAlertToChat(alert, chatID)
		}
	}
}

func (b *Bot) sendAlertToChat(alert alertInfo, chatID int64) {
	var (
		text = fmt.Sprintf("[Alerting] Alert from %s service\n%s", alert.service, alert.message)
		msg  tgbotapi.Chattable
	)
	if alert.file.name == "" {
		alert.file.name = "file"
	}
	switch {
	case alert.file.readCloser != nil:
		defer func() {
			err := alert.file.readCloser.Close()
			if err != nil {
				b.logger.Error("can't close alert file data",
					"alert", alert,
					"error", err,
				)
			}
		}()
		fr := tgbotapi.FileReader{
			Name:   alert.file.name,
			Reader: alert.file.readCloser,
			Size:   -1,
		}
		doc := tgbotapi.NewDocumentUpload(chatID, fr)
		doc.Caption = text

		msg = doc
	case alert.file.bytes != nil:
		fb := tgbotapi.FileBytes{
			Name:  alert.file.name,
			Bytes: alert.file.bytes,
		}
		doc := tgbotapi.NewDocumentUpload(chatID, fb)
		doc.Caption = text

		msg = doc
	case alert.file.path != "":
		doc := tgbotapi.NewDocumentUpload(chatID, alert.file.path)
		doc.Caption = text

		msg = doc
	default:
		msg = tgbotapi.NewMessage(chatID, text)
	}

	_, err := b.bot.Send(msg)
	if err == nil {
		b.logger.Debug("upload file to chat", "chatID", chatID, "alert", alert)
		return
	}

	b.logger.Error("can't send alert",
		"chatID", chatID,
		"alert", alert,
		"error", err,
	)
	metrics.IncError(ServiceName,
		fmt.Errorf("can't send alert=%v to chatID=%d: %w",
			alert,
			chatID,
			err,
		),
	)
}
