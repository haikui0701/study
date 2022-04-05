package game

import (
	"fmt"
	"regexp"
	"server/csvs"
	"time"
)

var manageBanWord *ManageBanWord

type ManageBanWord struct {
	BanWordBase  []string //配置生成
	BanWordExtra []string //更新
	Test         map[int]int
	MsgChannel   chan int
}

func GetManageBanWord() *ManageBanWord {
	if manageBanWord == nil {
		manageBanWord = new(ManageBanWord)
		manageBanWord.BanWordExtra = []string{"外挂", "工具", "原神"}
		manageBanWord.Test = make(map[int]int)
		manageBanWord.MsgChannel = make(chan int)
	}
	return manageBanWord
}

func (self *ManageBanWord) IsBanWord(txt string) bool {
	for _, v := range self.BanWordBase {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println("发现违禁词:", v)
		}
		if match {
			return match
		}
	}
	for _, v := range self.BanWordExtra {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println("发现违禁词:", v)
		}
		if match {
			return match
		}
	}
	return false
}

func (self *ManageBanWord) Run() {
	GetServer().AddGo()
	self.BanWordBase = csvs.GetBanWordBase()
	//基础词库的更新
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			if time.Now().Unix()%10 == 0 {
				//fmt.Println("违禁词库更新完成")
				GetServer().UpdateBanWord(self.BanWordBase)
			} else {
				//manageBanWord.Test[1] = 1
			}
		case _,ok:=<-self.MsgChannel:
			if !ok{
				GetServer().GoDone()
				return
			}
		}
	}
}

func (self *ManageBanWord) Close() {
	close(self.MsgChannel)
}
