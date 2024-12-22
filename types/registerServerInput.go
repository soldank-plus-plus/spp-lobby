package types

type RegisterServerInput struct {
    Advanced        *bool       `json:"advanced" binding:"required"`
    AntiCheatOn     *bool       `json:"anti_cheat_on" binding:"required"`
    BonusFrequency  *uint16     `json:"bonus_frequency" binding:"required"`
    Country         *string     `json:"country" binding:"required"`
    CurrentMap      string      `json:"current_map" binding:"required"`
    GameStyle       string      `json:"game_style" binding:"required"`
    Info            *string     `json:"info" binding:"required"`
    MaxPlayers      uint8       `json:"max_players" binding:"required"`
    Name            string      `json:"name" binding:"required"`
    NumBots         *uint8      `json:"num_bots" binding:"required"`
    OS              string      `json:"os" binding:"required"`
    Players         []string    `json:"players" binding:"required"`
    Port            uint16      `json:"port" binding:"required"`
    Private         *bool       `json:"private" binding:"required"`
    Realistic       *bool       `json:"realistic" binding:"required"`
    Respawn         *uint32     `json:"respawn" binding:"required"`
    Survival        *bool       `json:"survival" binding:"required"`
    Version         string      `json:"version" binding:"required"`
    WM              *bool       `json:"wm" binding:"required"`
}

const MaxCountryLength = 2
const MaxMapSize = 16
const MaxGameStyleSize = 3
const MaxInfoSize = 255
const MaxNameSize = 30
const MaxOSSize = 10
const MaxPlayerNameSize = 16
const MaxVersionSize = 10

func ConvertRegisterServerInputToServer(registerServerInput RegisterServerInput) Server {
    return Server{
        Advanced: *registerServerInput.Advanced,
        AntiCheatOn: *registerServerInput.AntiCheatOn,
        BonusFrequency: *registerServerInput.BonusFrequency,
        Country: *registerServerInput.Country,
        CurrentMap: registerServerInput.CurrentMap,
        GameStyle: registerServerInput.GameStyle,
        Info: *registerServerInput.Info,
        MaxPlayers: registerServerInput.MaxPlayers,
        Name: registerServerInput.Name,
        NumBots: *registerServerInput.NumBots,
        OS: registerServerInput.OS,
        Players: registerServerInput.Players,
        Port: registerServerInput.Port,
        Private: *registerServerInput.Private,
        Realistic: *registerServerInput.Realistic,
        Respawn: *registerServerInput.Respawn,
        Survival: *registerServerInput.Survival,
        Version: registerServerInput.Version,
        WM: *registerServerInput.WM,
    }
}

func ValidateRegisterServerInput(registerServerInput RegisterServerInput) bool {
    if len(*registerServerInput.Country) > MaxCountryLength {
        return false
    }

    if len(registerServerInput.CurrentMap) > MaxMapSize {
        return false
    }

    if len(registerServerInput.GameStyle) > MaxGameStyleSize {
        return false
    }

    if len(*registerServerInput.Info) > MaxInfoSize {
        return false
    }

    if len(registerServerInput.Name) > MaxNameSize {
        return false
    }

    if len(registerServerInput.OS) > MaxOSSize {
        return false
    }

    if len(registerServerInput.Players) > int(registerServerInput.MaxPlayers) {
        return false
    }

    for _, playerName := range registerServerInput.Players {
        if len(playerName) > MaxPlayerNameSize {
            return false
        }
    }

    if len(registerServerInput.Version) > MaxVersionSize {
        return false
    }

    return true
}

