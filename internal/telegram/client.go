package telegram

import "io"

type AlertClient struct {
	bot         *Bot
	serviceName string
}

func (b *Bot) NewAlertClient(serviceName string) *AlertClient {
	return &AlertClient{
		bot:         b,
		serviceName: serviceName,
	}
}

func (c *AlertClient) Alert(message string) {
	c.bot.alert(c.serviceName, message, fileInfo{})
}

func (c *AlertClient) AlertWithFilePath(message, filePath string) {
	c.bot.alert(c.serviceName, message, fileInfo{path: filePath})
}

func (c *AlertClient) AlertWithReadCloser(message, fileName string, fileData io.ReadCloser) {
	c.bot.alert(c.serviceName, message, fileInfo{
		name:       fileName,
		readCloser: fileData,
	})
}

func (c *AlertClient) AlertWithBytes(message, fileName string, fileBytes []byte) {
	c.bot.alert(c.serviceName, message, fileInfo{
		name:  fileName,
		bytes: fileBytes,
	})
}
