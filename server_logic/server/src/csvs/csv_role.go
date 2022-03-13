package csvs

import "server/utils"

type ConfigRole struct {
	RoleId          int   `json:"RoleId"`
	Star            int   `json:"Star"`
	Stuff           int   `json:"Stuff"`
	StuffNum        int64 `json:"StuffNum"`
	StuffItem       int   `json:"StuffItem"`
	StuffItemNum    int64 `json:"StuffItemNum"`
	MaxStuffItem    int   `json:"MaxStuffItem"`
	MaxStuffItemNum int64 `json:"MaxStuffItemNum"`
	Type            int   `json:"Type"`
}

var (
	ConfigRoleMap map[int]*ConfigRole
)

func init() {
	ConfigRoleMap = make(map[int]*ConfigRole)
	utils.GetCsvUtilMgr().LoadCsv("Role", &ConfigRoleMap)
	return
}

func GetRoleConfig(roleId int) *ConfigRole {
	return ConfigRoleMap[roleId]
}
