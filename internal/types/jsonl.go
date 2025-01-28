package types

type OT114K_LINE struct {
	Conversations []OT114K_DATA `json:"conversations"`
}

type OT114K_DATA struct {
	From  string `json:"from"`
	Value string `json:"value"`
}

type OT_PROCESSED_01_LINE struct {
	PrefixID  int    `json:"prefix_id"`
	Question  string `json:"question"`
	Reasoning string `json:"reasoning"`
	Response  string `json:"response"`
}
