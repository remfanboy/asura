package main

import (
	_ "asura/src/commands" // Initialize all commands and put them into an array
	"asura/src/database"
	"asura/src/handler"
	"asura/src/telemetry"
	"asura/src/utils/rinha"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"time"
)

func onReady(session disgord.Session, evt *disgord.Ready) {
	handler.Client.UpdateStatusString("Use j!comandos para ver meus comandos")
	telemetry.Info(fmt.Sprintf("%s Started", evt.User.Username), map[string]string{
		"eventType": "ready",
	})
	go telemetry.MetricUpdate(handler.Client)

}
func main() {
	rand.Seed(time.Now().UnixNano())
	rinha.DailyGalo = rinha.GetRand()
	// If it's not in production so it's good to read a ".env" file
	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}
	// Initialize datalog services for telemetry of the application
	telemetry.Init()
	database.Init()
	fmt.Println("Starting bot...")
	client := disgord.New(disgord.Config{
		RejectEvents: []string{"PRESENCE_UPDATE", "GuildMessageTyping"},
		BotToken:     os.Getenv("TOKEN"),
	})
	defer client.StayConnectedUntilInterrupted()
	handler.Client = client
	client.Gateway().MessageCreate(handler.OnMessage)
	client.Gateway().MessageUpdate(handler.OnMessageUpdate)
	client.Gateway().MessageReactionAdd(handler.OnReactionAdd)
	client.Gateway().MessageReactionRemove(handler.OnReactionRemove)
	client.Gateway().Ready(onReady)
	
}
