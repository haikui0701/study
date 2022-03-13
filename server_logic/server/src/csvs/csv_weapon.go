package csvs

import "server/utils"

type ConfigWeapon struct {
	WeaponId int `json:"WeaponId"`
	Type     int `json:"Type"`
	Star     int `json:"Star"`
}

type ConfigWeaponLevel struct {
	WeaponStar    int `json:"WeaponStar"`
	Level         int `json:"Level"`
	NeedExp       int `json:"NeedExp"`
	NeedStarLevel int `json:"NeedStarLevel"`
}

type ConfigWeaponStar struct {
	WeaponStar int `json:"WeaponStar"`
	StarLevel  int `json:"StarLevel"`
	Level      int `json:"Level"`
}

var (
	ConfigWeaponMap        map[int]*ConfigWeapon
	ConfigWeaponLevelSlice []*ConfigWeaponLevel
	ConfigWeaponStarSlice  []*ConfigWeaponStar
)

func init() {
	ConfigWeaponMap = make(map[int]*ConfigWeapon)
	utils.GetCsvUtilMgr().LoadCsv("Weapon", &ConfigWeaponMap)

	utils.GetCsvUtilMgr().LoadCsv("WeaponLevel", &ConfigWeaponLevelSlice)
	utils.GetCsvUtilMgr().LoadCsv("WeaponStar", &ConfigWeaponStarSlice)
	return
}

func GetWeaponConfig(weaponId int) *ConfigWeapon {
	return ConfigWeaponMap[weaponId]
}
