package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type TaskInfo struct {
	TaskId int
	State  int
}

type ModUniqueTask struct {
	MyTaskInfo map[int]*TaskInfo

	player *Player
	path   string
}

func (self *ModUniqueTask) IsTaskFinish(taskId int) bool {
	if taskId == 10001 || taskId == 10002 {
		return true
	}

	task, ok := self.MyTaskInfo[taskId]
	if !ok {
		return false
	}
	return task.State == TASK_STATE_FINISH
}

func (self *ModUniqueTask) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (self *ModUniqueTask) LoadData(player *Player) {

	self.player = player
	self.path = self.player.localPath + "/uniquetask.json"

	configFile, err := ioutil.ReadFile(self.path)
	if err != nil {
		fmt.Println("error")
		return
	}
	err = json.Unmarshal(configFile, &self)
	if err != nil {
		self.InitData()
		return
	}

	if self.MyTaskInfo == nil {
		self.MyTaskInfo = make(map[int]*TaskInfo)
	}
	return
}

func (self *ModUniqueTask) InitData() {

}