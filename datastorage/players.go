package datastorage

type Player struct {
	Slot int    `json:"slot"`
	Name string `json:"name"`
}

var SlotNumbersToAPSlots map[int]Player
