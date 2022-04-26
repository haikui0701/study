package game

type MsgHead struct {
	MsgId    int `json:"msgid" `
}

type MsgLogin struct {
	MsgId    int    `json:"msgid" `
	Account  string `json:"account" `
	Password string `json:"password" `
	UserId   int64  `json:"userid" `
}

type MsgPool struct {
	MsgId    int `json:"msgid" `
	PoolType int `json:"pooltype" `
}