package types

type Server struct {
    Advanced        bool        `json:"advanced"`
    AntiCheatOn     bool        `json:"anti_cheat_on"`
    BonusFrequency  uint16      `json:"bonus_frequency"`
    Country         string      `json:"country"`
    CurrentMap      string      `json:"current_map"`
    GameStyle       string      `json:"game_style"`
    IP              string      `json:"ip"`
    Info            string      `json:"info"`
    MaxPlayers      uint8       `json:"max_players"`
    Name            string      `json:"name"`
    NumBots         uint8       `json:"num_bots"`
    OS              string      `json:"os"`
    Players         []string    `json:"players"`
    Port            uint16      `json:"port"`
    Private         bool        `json:"private"`
    Realistic       bool        `json:"realistic"`
    Respawn         uint32      `json:"respawn"`
    Survival        bool        `json:"survival"`
    Version         string      `json:"version"`
    WM              bool        `json:"wm"`
    UpdatedAt       int64       `json:"-"`
}

