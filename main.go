package main

import (
	"bufio"
	"fmt"
	"log"
	"net/textproto"
	"os"
	"strings"

	"twitch_chat_bot/bot"
	"twitch_chat_bot/config"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur chargement du .env")
	}

	config.Token = os.Getenv("TWITCH_OAUTH_TOKEN")
	config.Username = os.Getenv("TWITCH_USERNAME")
	config.Channel = os.Getenv("TWITCH_CHAN_NAME")

	if config.Token == "" || config.Username == "" || config.Channel == "" {
		log.Fatal("Probleme lors de la reception des variables d'environement")
	}
	bot.LoadMap()
}


func main() {
	log.Println("Debut du programme")

	conn := bot.Connect()
	defer bot.Disconnect(conn)

	tp := textproto.NewReader(bufio.NewReader(conn))
	for { // boucle infinie
		status, err := tp.ReadLine()
		
		if err != nil {
			panic(err)
		}
		fmt.Println(status)

		if strings.HasPrefix(status, "PING") {
			fmt.Fprintf(conn, "PONG :tmi.twitch.tv\r\n")
		}
		
		if strings.Contains(status, "PRIVMSG") {
			messageParts := strings.Split(status, " :")
			if len(messageParts) > 1 {
				message := messageParts[1]
				bot.ProcessMSG(conn, message)
			}
		}
	}
	log.Println("Fin du programme")
}
