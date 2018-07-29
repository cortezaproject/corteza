package main

import (
	"bufio"
	"encoding/json"
	"github.com/SentimensRG/sigctx"
	"github.com/crusttech/crust/sam/websocket/incoming"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
)

func main() {
	var ctx = sigctx.New()

	var wsd = websocket.Dialer{}
	conn, _, err := wsd.Dial("ws://localhost:3000/websocket/?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzM5ODg4MzksInN1YiI6IjIwNDA4NTk3MDc2MzUxMzk2NSJ9.aJQE-_hH_bxbri0SDmBok_UzvYnJrvDYc16E_6meJ6w", map[string][]string{})
	must(err)

	defer conn.Close()
	go func() {

		log := log.New(os.Stderr, "<<<<", log.LstdFlags)
		// Reader
		for {
			mt, p, err := conn.ReadMessage()
			must(err)
			log.Printf("%d %s", mt, string(p))
		}
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		var chanId = ""

		// Writer
		for {
			func() {
				if chanId == "" {
					print("NO CHANNEL")
				} else {
					print(chanId)
				}
				print("> ")
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)
				if len(text) == 0 {
					return
				}

				if text[:1] == "/" {
					cmdSplit := strings.SplitN(text, " ", 2)
					switch cmdSplit[0] {
					case "/join":
						if len(cmdSplit) < 2 {
							println("/join <channelid>")
						} else {
							joinCh(conn, cmdSplit[1])
							chanId = cmdSplit[1]
							msgHistory(conn, chanId)
						}
					case "/part":
						partCh(conn, chanId)
						chanId = ""

					case "/list":
						listCh(conn)

					case "/raw":
						sendRaw(conn, cmdSplit[1])

					default:
						println("Unknown command, try:")
						println("   /join <channel-id>")
						println("   /part")
						println("   /list")
						println("   /raw <raw json to send>")
					}
				} else {
					sendMsg(conn, text, chanId)
				}
			}()
		}
	}()

	<-ctx.Done()
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func should(err error) {
	if err != nil {
		log.Printf("Warning: %+v", err)
	}
}

func sendMsg(conn *websocket.Conn, text, chanId string) {
	pb, err := json.Marshal(incoming.Payload{MessageCreate: &incoming.MessageCreate{Message: text, ChannelID: chanId}})
	should(err)
	should(conn.WriteMessage(websocket.TextMessage, pb))
}

func listCh(conn *websocket.Conn) {
	pb, err := json.Marshal(incoming.Payload{ChannelList: &incoming.ChannelList{}})
	should(err)
	should(conn.WriteMessage(websocket.TextMessage, pb))
}

func joinCh(conn *websocket.Conn, channelId string) {
	pb, err := json.Marshal(incoming.Payload{ChannelJoin: &incoming.ChannelJoin{ChannelID: channelId}})
	should(err)
	should(conn.WriteMessage(websocket.TextMessage, pb))
}

func partCh(conn *websocket.Conn, channelId string) {
	pb, err := json.Marshal(incoming.Payload{ChannelPart: &incoming.ChannelPart{ChannelID: channelId}})
	should(err)
	should(conn.WriteMessage(websocket.TextMessage, pb))
}

func msgHistory(conn *websocket.Conn, channelId string) {
	pb, err := json.Marshal(incoming.Payload{MessageHistory: &incoming.MessageHistory{ChannelID: channelId}})
	should(err)
	should(conn.WriteMessage(websocket.TextMessage, pb))
}

func sendRaw(conn *websocket.Conn, msg string) {
	should(conn.WriteMessage(websocket.TextMessage, []byte(msg)))
}
