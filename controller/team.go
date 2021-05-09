package controller

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
)

// CreateTeam 创建队伍
func CreateTeam(c *gin.Context) {
	var r struct {
		Name  string `json:"name"`  // 组名
		Users string `json:"users"` // 用户
		Month uint8  `json:"month"` // 月份
		Key   string `json:"key"`   // 标识
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	team := model.Team{
		Name:  r.Name,
		Users: r.Users,
		Month: r.Month,
		Key:   r.Key,
	}
	if err := db.Create(&team).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "创建队伍成功", c)
		return
	}

	handle.ReturnSuccess("ok", team, c)
}

// UpdateTeam 更新队伍
func UpdateTeam(c *gin.Context) {
	var team model.Team

	if err := c.ShouldBind(&team); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}

	db := config.GetMysql()
	if err := db.Save(&team).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "组队更新失败", c)
		return
	}

	handle.ReturnSuccess("ok", team, c)
}

// DeleteTeam 批量删除分组
func DeleteTeam(c *gin.Context) {
	var r struct {
		Teams []uint `json:"teams" binding:"required"`
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()

	if err := db.Delete(&model.Team{}, r.Teams).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "删除失败", c)
	}

	handle.ReturnSuccess("ok", r, c)
}

// GetTeamRank 获取队伍积分排名
func GetTeamRank(c *gin.Context) {
	_ranks := GetIntegral()
	rankes := make(map[uint]uint8)
	for _, r := range _ranks {
		rankes[r.ID] = r.Integral
	}

	var teams []model.Team
	db := config.GetMysql()
	db.Find(&teams)
	var teamRanks TeamRanks
	for _, t := range teams {
		var usersID []uint
		json.Unmarshal([]byte(t.Users), &usersID)
		var intergarl uint8

		for _, u := range usersID {
			intergarl += rankes[u]
		}

		teamRanks = append(teamRanks, TeamRank{
			ID:       t.ID,
			Name:     t.Name,
			Integral: intergarl,
			Key:      t.Key,
		})
	}

	sort.Sort(teamRanks)

	handle.ReturnSuccess("ok", teamRanks, c)
}

// GetTeams 获取队伍
func GetTeams(c *gin.Context) {
	var r struct {
		Name  string `json:"name"`  // 组名
		Month uint8  `json:"month"` // 月份
		Key   string `json:"key"`   // 标识
	}

	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.GetMysql()
	var teams []model.Team
	if r.Name != "" {
		db = db.Where("name LIKE ?", "%"+r.Name+"%")
	}

	if r.Month != 0 {
		db = db.Where("month = ?", r.Month)
	}

	if r.Key != "" {
		db = db.Where("key LIKE ?", "%"+r.Key+"%")
	}

	db.Order("created_at DESC").Find(&teams)
	handle.ReturnSuccess("ok", teams, c)
}

type TeamRank struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Integral uint8  `json:"integral"`
	Key      string `json:"key"`
}

type TeamRanks []TeamRank

// Len()
func (s TeamRanks) Len() int {
	return len(s)
}

// Less():成绩将有低到高排序
func (s TeamRanks) Less(i, j int) bool {
	return s[i].Integral > s[j].Integral
}

// Swap()
func (s TeamRanks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
