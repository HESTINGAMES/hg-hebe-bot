package api

import (
	"encoding/json"
	"fmt"

	"github.com/hestingames/hg-hebe-bot/internal/apiclient"
)

var (
	ApiBaseUrl string
)

func InitializeCsgoApi(BaseUrl string) {
	ApiBaseUrl = BaseUrl
	apiclient.InitializeApiClient()

	// Development api is using a self signed certificate
	apiclient.DisableCertificateCheck()
}

func GetMatchakingQueueStatus() ([]MatmakingQueueStatus, error) {
	var queueStatus []MatmakingQueueStatus
	bytes, err := apiclient.DoRequest("GET", fmt.Sprintf("%squery/queues", ApiBaseUrl))
	if err != nil {
		return queueStatus, err
	}

	err = json.Unmarshal(bytes, &queueStatus)
	return queueStatus, err
}

func GetPlayingNow() (int, error) {
	playingNow := 0
	var csgoServers CsgoServersResponse
	bytes, err := apiclient.DoRequest("GET", fmt.Sprintf("%squery/servers", ApiBaseUrl))
	if err != nil {
		return playingNow, err
	}

	err = json.Unmarshal(bytes, &csgoServers)
	if err != nil {
		return playingNow, err
	}

	for i := range csgoServers.Servers {
		playingNow += len(csgoServers.Servers[i].PlayersId)
	}

	return playingNow, err
}

func ParseServerStatus(matchmakingStatus []MatmakingQueueStatus) CsgoServerStatus {
	var serverStatus CsgoServerStatus

	for _, queue := range matchmakingStatus {
		serverStatus.PlayingNow += uint32(queue.CurrentPlaying)
		serverStatus.SearchingNow += uint32(queue.CurrentSearching)

		switch queue.GameType {
		case int(Sigma):
			serverStatus.SigmaSearching += uint32(queue.CurrentSearching)
			serverStatus.SigmaPlaying += uint32(queue.CurrentPlaying)
		case int(Delta):
			serverStatus.DeltaSearching += uint32(queue.CurrentSearching)
			serverStatus.DeltaPlaying += uint32(queue.CurrentPlaying)
		case int(DustII):
			serverStatus.DustIISearching += uint32(queue.CurrentSearching)
			serverStatus.DustIIPlaying += uint32(queue.CurrentPlaying)
		case int(Hostages):
			serverStatus.HostagesSearching += uint32(queue.CurrentSearching)
			serverStatus.HostagesPlaying += uint32(queue.CurrentPlaying)
		}
	}

	return serverStatus
}
