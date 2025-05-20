package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var token, channel, username string

var CmdMap = make(map[string]string)

// load les valeurs du csv dans la map (hash table)
func loadMap(){
	file, err := os.Open("dic.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	
	for i, line := range data {
		if i == 0 {
			continue // on skip le header
		}
		tempCmd := line[0]
		tempMsg := line[1]

		CmdMap[tempCmd] = tempMsg
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur chargement du .env")
	}

	token = os.Getenv("TWITCH_OAUTH_TOKEN")
	username = os.Getenv("TWITCH_USERNAME")
	channel = os.Getenv("TWITCH_CHAN_NAME")

	if token == "" || username == "" || channel == "" {
		log.Fatal("Probleme lors de la reception des variables d'environement")
	}
	loadMap()
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, "PASS %s\r\n", token)
	fmt.Fprintf(conn, "NICK %s\r\n", username)
	fmt.Fprintf(conn, "JOIN #%s\r\n", channel)

	return conn
}

func disconnect(conn net.Conn) {
	conn.Close()
}

func sendMessage(conn net.Conn, message string) {
	fmt.Fprintf(conn, "PRIVMSG #%s :%s\r\n", channel, message)
}

func processMSG(conn net.Conn, message string) {
	if response, ok := CmdMap[message]; ok {
    sendMessage(conn, response)
	}
	if message == "!exit" {
		disconnect(conn)
		log.Println("Exit du programme via intervention de l'utilisateur")
		os.Exit(0)
	}
}

func main() {
	log.Println("Debut du programme")

	conn := connect()
	defer disconnect(conn)

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
				processMSG(conn, message)
			}
		}
	}
	log.Println("Fin du programme")
}
