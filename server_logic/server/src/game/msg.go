package game

type MsgLogin struct {
	MsgId    int    `json:"msgid" `
	Account  string `json:"account" `
	Password string `json:"password" `
	UserId   int64  `json:"userid" `
}
