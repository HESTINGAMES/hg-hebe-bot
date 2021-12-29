package api

type MatmakingQueueStatus struct {
	GameType         int     `json:"gameType"`
	CurrentPlaying   int     `json:"currentPlaying"`
	CurrentSearching int     `json:"currentSearching"`
	ServerList       []int64 `json:"serverList"`
}

type CsgoServer struct {
	ServerId          int64   `json:"serverId"`
	MapName           string  `json:"map"`
	GameType          int     `json:"gameType"`
	PlayersId         []int64 `json:"playersId"`
	MatchmakingServer bool    `json:"matchmakingServer"`
}

type CsgoServersResponse struct {
	Servers []CsgoServer `json:"servers"`
}

type CsgoServerStatus struct {
	SearchingNow      uint32
	PlayingNow        uint32
	SigmaSearching    uint32
	SigmaPlaying      uint32
	DeltaSearching    uint32
	DeltaPlaying      uint32
	DustIISearching   uint32
	DustIIPlaying     uint32
	HostagesSearching uint32
	HostagesPlaying   uint32
}

type GameType int

const (
	DustII   GameType = 519
	Hostages GameType = 135200775
	Delta    GameType = 269530119
	Sigma    GameType = 1635794951
)
