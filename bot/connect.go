package bot

import(
	"net"
	"fmt"

	"twitch_chat_bot/config"
)


func Connect() net.Conn {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, "PASS %s\r\n", config.Token)
	fmt.Fprintf(conn, "NICK %s\r\n", config.Username)
	fmt.Fprintf(conn, "JOIN #%s\r\n", config.Channel)

	return conn
}

func Disconnect(conn net.Conn) {
	conn.Close()
}
