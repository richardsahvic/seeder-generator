package logic

type RequestData struct {
	TMMID       string `json:"tmmID"`
	ChannelCode string `json:"channelCode"`
}

type Request struct {
	PgAccountID string        `json:"accountID"`
	Data        []RequestData `json:"data"`
}
