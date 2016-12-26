package monitoring

import (
	"encoding/json"
	"net/http"
)

func channelStatus(w http.ResponseWriter, r *http.Request) {
	allChannelsResponseData := make(map[string]Channel)
	for k, v := range monitoringStruct.Channels {
		var oneChannel Channel
		oneChannel.NumberOfMessages = len(v)
		allChannelsResponseData[k] = oneChannel
	}

	if err := json.NewEncoder(w).Encode(allChannelsResponseData); err != nil {
		panic(err)
	}
}
