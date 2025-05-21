package bot

import (
	"net"
	"fmt"
	"os"
	"log"
	"encoding/csv"

	"twitch_chat_bot/config"
)

var CmdMap = make(map[string]string)

// load les valeurs du csv dans la map (hash table)
func LoadMap(){
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

func SendMessage(conn net.Conn, message string) {
	fmt.Fprintf(conn, "PRIVMSG #%s :%s\r\n", config.Channel, message)
}

func ProcessMSG(conn net.Conn, message string) {
	if response, ok := CmdMap[message]; ok {
    SendMessage(conn, response)
	}
	if message == "!exit" {
		Disconnect(conn)
		log.Println("Exit du programme via intervention de l'utilisateur")
		os.Exit(0)
	}
}
