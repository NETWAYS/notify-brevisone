package main

type AuthTokenAnswer struct {
	Token    string `json:"jwt"`
	ExpireAt uint   `json:"expireAt"`
}

type Recipient struct {
	To     string `json:"to"`
	Target string `json:"target"`
}

type Message struct {
	Recipients []Recipient `json:"recipients"`
	Text       string      `json:"text"`
	Provider   string      `json:"provider"`
	Type       string      `json:"type"`
}
