package main

import "github.com/Maksym-Perehinets/discord_ai_assistant/bot"

func main() {
	b := bot.StartBot()
	exitHandler := b.Run()
	b.SetUpSlashCommands("1350120716497846374")
	exitHandler()
}
