package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type Message map[string]interface{}

type Version struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Build int    `json:"build"`
	Class string `json:"class"`
}

type RoomInfo struct {
	Version               Version
	GenratorVersion       Version
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

var roomInfo RoomInfo

func main() {
	_, sender, messagesCh := start_websocket()

	for m := range messagesCh {
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
					Game:           "Minecraft",
					Name:           "Minecraft",
					Uuid:           uuid.New().String(),
					Version:        roomInfo.Version,
					Items_handling: 0,
					Tags:           []string{"TextOnly"},
					Slot_data:      false,
				}
				connectPacketBytes, err := json.Marshal([]ConnectPacket{connectPacket})
				if err != nil {
					log.Println("json:", err)
					return
				}

				sender <- connectPacketBytes
			case "PrintJSON":
				HandlePrintJson(messageBytes)
			default:
				log.Println("unknown message:", cmd)
			}
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
	log.Println(result["data"])
}
