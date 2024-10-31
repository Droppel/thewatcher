package archipelago

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"watcher/discordbot"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type Message map[string]interface{}
type Item struct {
	Flags    int `json:"flags"`
	Item     int `json:"item"`
	Location int `json:"location"`
	Player   int `json:"player"`
}

type Version struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Build int    `json:"build"`
	Class string `json:"class"`
}

type RoomInfo struct {
	Version               Version
	Tags                  []string
	Password              bool
	Permissions           map[string]int
	Hint_cost             int
	Location_check_points int
	Games                 []string
	Datapackage_cecksums  map[string]string
	Seed_name             string
	Time                  float64
}

type Game struct {
	Item_name_to_id     map[string]int `json:"item_name_to_id"`
	Location_name_to_id map[string]int `json:"location_name_to_id"`
}

type IdMap struct {
	Item_id_to_name     map[int]string `json:"item_id_to_name"`
	Location_id_to_name map[int]string `json:"location_id_to_name"`
}

type ConnectPacket struct {
	Cmd            string   `json:"cmd"`
	Password       string   `json:"password"`
	Game           string   `json:"game"`
	Name           string   `json:"name"`
	Uuid           string   `json:"uuid"`
	Version        Version  `json:"version"`
	Items_handling int      `json:"items_handling"`
	Tags           []string `json:"tags"`
	Slot_data      bool     `json:"slot_data"`
}

type Player struct {
	Slot int    `json:"slot"`
	Name string `json:"name"`
}

var roomInfo RoomInfo
var IdMaps IdMap
var players map[int]Player
var discordMessageCh chan discordbot.DiscordAction

func Connect(discMessageCh chan discordbot.DiscordAction) {
	discordMessageCh = discMessageCh

	var messagesCh chan []byte
	_, sender, messagesCh, done := start_websocket()

	for {
		select {
		case m := <-messagesCh:
			// Parse the message
			messages := []Message{}
			err := json.Unmarshal(m, &messages)
			if err != nil {
				log.Println("json:", err)
				return
			}
			for _, message := range messages {
				cmd := message["cmd"]
				messageBytes, err := json.Marshal(message)
				if err != nil {
					log.Println("json:", err)
					return
				}
				switch cmd {
				case "RoomInfo":
					HandleRoomInfo(messageBytes)

					connectPacket := ConnectPacket{
						Cmd:            "Connect",
						Password:       "",
						Game:           "",
						Name:           os.Getenv("ARCHIPELAGO_NAME"),
						Uuid:           uuid.New().String(),
						Version:        roomInfo.Version,
						Items_handling: 0b111,
						Tags:           []string{"TextOnly"},
						Slot_data:      false,
					}
					connectPacketBytes, err := json.Marshal([]ConnectPacket{connectPacket})
					if err != nil {
						log.Println("json:", err)
						return
					}

					sender <- connectPacketBytes
				case "Connected":
					players = make(map[int]Player, 0)
					playerdata := message["players"].([]interface{})
					for _, player := range playerdata {
						var playerData Player
						err = mapstructure.Decode(player, &playerData)
						if err != nil {
							log.Println("json:", err)
							return
						}
						players[playerData.Slot] = playerData
					}
					games := strings.Join(roomInfo.Games, `","`)
					sender <- []byte(fmt.Sprintf(`[{"cmd":"GetDataPackage","games":["%s"]}]`, games))
				case "DataPackage":
					data := message["data"].(map[string]interface{})["games"].(map[string]interface{})
					IdMaps = IdMap{
						Item_id_to_name:     make(map[int]string),
						Location_id_to_name: make(map[int]string),
					}
					for _, game := range data {
						var gameData Game
						err = mapstructure.Decode(game, &gameData)
						if err != nil {
							log.Println("json:", err)
							return
						}
						for name, id := range gameData.Item_name_to_id {
							IdMaps.Item_id_to_name[id] = name
						}
						for name, id := range gameData.Location_name_to_id {
							IdMaps.Location_id_to_name[id] = name
						}
					}
				case "PrintJSON":
					HandlePrintJson(messageBytes)
				default:
					log.Println("unknown message:", cmd)
				}
			}
		case <-done:
			_, sender, messagesCh, done = start_websocket()
		}
	}
}

func HandleRoomInfo(m []byte) {
	err := json.Unmarshal(m, &roomInfo)
	if err != nil {
		log.Println("json:", err)
		return
	}
}

func HandlePrintJson(m []byte) {
	var result map[string]interface{}
	err := json.Unmarshal(m, &result)
	if err != nil {
		log.Println("json:", err)
		return
	}
	switch result["type"] {
	case "ItemSend":
		HandleItemSend(result)
	case "Join":
		HandleJoin(result)
	case "Part":
		HandlePart(result)
	}
}

func HandleItemSend(result map[string]interface{}) {
	var item Item
	err := mapstructure.Decode(result["item"], &item)
	if err != nil {
		log.Println("decode:", err)
		return
	}

	if item.Flags == 0 || item.Flags == 4 { // Filler or Trap
		return
	}

	isProgression := item.Flags == 1

	recvPlayer := result["receiving"].(float64)
	slotName := players[int(recvPlayer)].Name
	discMsg := fmt.Sprintf("%s received %s", slotName, IdMaps.Item_id_to_name[item.Item])
	discordMessageCh <- discordbot.DiscordAction{
		Type: "message",
		Message: discordbot.DiscordMessage{
			Slot:    int(recvPlayer),
			Message: discMsg,
			Silent:  !isProgression, // Send silent messages for useful items
		},
	}

	discordMessageCh <- discordbot.DiscordAction{
		Type: "channel_topic",
		ChannelTopicEdit: discordbot.DiscordChannelTopicEdit{
			Slot:  int(recvPlayer),
			Topic: discordbot.UNKNOWN_STATUS,
		},
	}
}

func HandleJoin(result map[string]interface{}) {
	slot := result["slot"].(float64)
	slotName := players[int(slot)].Name
	text := result["data"].([]interface{})[0].(map[string]interface{})["text"].(string)

	if isTextOnly(text) {
		return
	}

	discMsg := fmt.Sprintf("%s joined", slotName)
	discordMessageCh <- discordbot.DiscordAction{
		Type: "message",
		Message: discordbot.DiscordMessage{
			Slot:    int(slot),
			Message: discMsg,
			Silent:  true,
		},
	}
}

func HandlePart(result map[string]interface{}) {
	slot := result["slot"].(float64)
	slotName := players[int(slot)].Name
	text := result["data"].([]interface{})[0].(map[string]interface{})["text"].(string)

	if isTextOnly(text) {
		return
	}

	discMsg := fmt.Sprintf("%s left", slotName)
	discordMessageCh <- discordbot.DiscordAction{
		Type: "message",
		Message: discordbot.DiscordMessage{
			Slot:    int(slot),
			Message: discMsg,
			Silent:  true,
		},
	}
}

func isTextOnly(msg string) bool {
	return strings.Contains(msg, "TextOnly") || strings.Contains(msg, "IgnoreGame") || strings.Contains(msg, "Tracker")
}
