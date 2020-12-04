package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiao0811/xiniu/config"
	"github.com/xiao0811/xiniu/handle"
	"github.com/xiao0811/xiniu/model"
	"gorm.io/gorm"
)

// CreateContract 创建新的合约
func CreateContract(c *gin.Context) {
	var r struct {
		MemberID                 uint    `json:"member_id"`                   // 门店名称
		CooperationTime          string  `json:"cooperation_time"`            // 合作时间
		ExpireTime               string  `json:"expire_time"`                 // 到期时间
		IsStartService           bool    `json:"is_start_service"`            // 是否开始服务
		DelayTime                string  `json:"delay_time"`                  // 延期后到期时间
		ContractAmount           uint    `json:"contract_amount"`             // 签约金额
		Arrives                  bool    `json:"arrives"`                     // 是否到账
		Arrears                  uint    `json:"arrears"`                     // 有无欠款
		CurrentStoreCollections  uint    `json:"current_store_collections"`   // 目前门店收藏量
		CurrentNumber            uint    `json:"current_number"`              // 目前评价数
		CurrentStar              float32 `json:"current_star"`                // 目前星级
		CurrentLeaderboard       string  `json:"current_leaderboard"`         // 目前排行榜
		StoreCollections         uint    `json:"store_collections"`           // 门店收藏量
		InformationFlow          uint    `json:"information_flow"`            // 信息流
		BigVReview               uint    `json:"big_v_review"`                // 大V评论
		GroupBuyingVolume        uint    `json:"group_buying_volume"`         // 团购卖量
		Like                     uint    `json:"like"`                        // 点赞
		FollowPeers              string  `json:"follow_peers"`                // 关注同行
		CurrentStatusOfPromotion int8    `json:"current_status_of_promotion"` // 推广通现状
		Upgrade                  bool    `json:"upgrade"`                     // 是否提升金牌店铺
		IncludeDetailsPage       bool    `json:"include_details_page"`        // 是否包含详情页
		Remarks                  string  `json:"remarks"`                     // 备注
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "请求数据不正确", c)
		return
	}
	db := config.MysqlConn
	var _c model.Contract
	if err := db.Where("status = 1 AND DelayTime > ?", time.Now()).First(&_c).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "存在未完成的合约", c)
		return
	}
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	// 没传延期时间则设置为到期时间
	if r.DelayTime == "" {
		r.DelayTime = r.ExpireTime
	}
	_CooperationTime, _ := time.ParseInLocation("2006-01-02", r.CooperationTime, time.Local)
	_ExpireTime, _ := time.ParseInLocation("2006-01-02", r.ExpireTime, time.Local)
	_DelayTime, _ := time.ParseInLocation("2006-01-02", r.DelayTime, time.Local)
	var operations model.User
	var business model.User
	var member model.Member
	if r.MemberID != 0 {
		db.Where("id = ?", r.MemberID).First(&member)
	}
	db.Where("id = ?", member.OperationsStaff).First(&operations)
	db.Where("id = ?", member.BusinessPeople).First(&business)
	if member.FirstCreate.String() == "0001-01-01 00:00:00 +0000 UTC" {
		member.FirstCreate = model.MyTime{Time: _CooperationTime}
	}
	// 查看用户sort
	var _sort model.Contract
	var sort uint8
	if err := db.Where("id = ?", r.MemberID).Order("created_at desc").First(&_sort).Error; err != nil {
		sort = 0
	}
	sort = _sort.Sort + 1
	con := model.Contract{
		UUID:                     "XINIU-ORD-" + time.Now().Format("200601021504") + strconv.Itoa(handle.RandInt(1000, 9999)),
		MemberID:                 r.MemberID,
		CooperationTime:          model.MyTime{Time: _CooperationTime},
		ExpireTime:               model.MyTime{Time: _ExpireTime},
		IsStartService:           r.IsStartService,
		DelayTime:                model.MyTime{Time: _DelayTime},
		ContractAmount:           r.ContractAmount,
		Arrives:                  r.Arrives,
		Arrears:                  r.Arrears,
		CurrentStoreCollections:  r.CurrentStoreCollections,
		CurrentNumber:            r.CurrentNumber,
		CurrentStar:              float32(r.CurrentStar),
		CurrentLeaderboard:       r.CurrentLeaderboard,
		StoreCollections:         r.StoreCollections,
		InformationFlow:          r.InformationFlow,
		BigVReview:               r.BigVReview,
		GroupBuyingVolume:        r.GroupBuyingVolume,
		Like:                     r.Like,
		FollowPeers:              r.FollowPeers,
		CurrentStatusOfPromotion: r.CurrentStatusOfPromotion,
		Upgrade:                  r.Upgrade,
		IncludeDetailsPage:       r.IncludeDetailsPage,
		OperationsStaff:          operations.RealName,
		BusinessPeople:           business.RealName,
		Remarks:                  r.Remarks,
		Status:                   0,
		Sort:                     sort,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Create(&con).Error; err != nil {
			// handle.ReturnError(http.StatusBadRequest, "门店创建失败", c)
			return err
		}
		member.ExpireTime = model.MyTime{Time: _DelayTime}
		member.NumberOfContracts++
		member.Refund = model.MyTime{}
		if err := db.Save(&member).Error; err != nil {
			return err
		}

		// 创建用户记录
		l := model.UserLog{
			Operator: token.UserID,
			Action:   "Create Contact",
			Contract: con.ID,
			Remarks:  token.FullName + "创建合约: " + con.UUID,
		}
		if err := db.Create(&l).Error; err != nil {
			// handle.ReturnError(http.StatusBadRequest, "门店创建失败", c)
			return err
		}

		return nil
	})

	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "门店创建失败", c)
		return
	}

	handle.ReturnSuccess("ok", con, c)
}

// ContractList 合约列表
func ContractList(c *gin.Context) {
	var r struct {
		UUID     string `json:"uuid"`
		MemberID uint   `json:"member_id"`
		Status   int    `json:"status"`
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "用户名密码输入不正确", c)
		return
	}
	db := config.MysqlConn
	var contracts []model.Contract
	var count int64
	var pages int
	sql := db.Preload("Member")
	if r.UUID != "" {
		sql = sql.Where("uuid = ?", r.UUID)
	}
	if r.MemberID != 0 {
		sql = sql.Where("member_id = ?", r.MemberID)
	}
	if r.Status != -1 {
		sql = sql.Where("status = ?", r.Status)
	}
	sql.Offset((r.Page - 1) * 10).Find(&contracts).Count(&count)
	if count == 0 {
		handle.ReturnSuccess("ok", nil, c)
		return
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	if int(count)%r.Limit != 0 {
		pages = int(count)/r.Limit + 1
	} else {
		pages = int(count) / r.Limit
	}
	currPage := r.Page/r.Limit + 1
	handle.ReturnSuccess("ok", gin.H{"contracts": contracts, "pages": pages, "currPage": currPage}, c)
}

// UpdateContract 更新合约
func UpdateContract(c *gin.Context) {
	var r model.Contract
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, err.Error(), c)
		return
	}
	var co model.Contract
	db := config.MysqlConn
	if err := db.Where("id = ?", r.ID).First(&co).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约不存在", c)
		return
	}
	if err := db.Updates(&r).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约更新失败", c)
		return
	}
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Update Contact",
		Contract: r.ID,
		Remarks:  token.FullName + "更新合约: " + r.UUID,
	}
	db.Create(&l)

	handle.ReturnSuccess("ok", r, c)
}

// ContractReview 合约审核
func ContractReview(c *gin.Context) {
	var r struct {
		ID     int    `json:"id" binding:"required"`
		Status int8   `json:"status" binding:"required"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	var co model.Contract
	if err := db.Where("id = ?", r.ID).First(&co).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约不存在", c)
		return
	}
	var _c model.Contract
	if err := db.Where("status = 1 AND DelayTime > ?", time.Now()).First(&_c).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "存在未完成的合约", c)
		return
	}
	co.Status = r.Status
	if err := db.Save(&co).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "审核失败", c)
		return
	}
	// 创建用户记录
	_token, _ := c.Get("token")
	token, _ := _token.(*handle.JWTClaims)
	var msg string
	if r.Status == 1 {
		msg = token.FullName + "审核合约通过: " + strconv.Itoa(r.ID)
	} else if r.Status == 2 {
		msg = token.FullName + "审核用户拒绝: " + r.Remark
	}
	l := model.UserLog{
		Operator: token.UserID,
		Action:   "Review Contract",
		Member:   co.ID,
		Remarks:  msg,
	}
	db.Create(&l)
	handle.ReturnSuccess("ok", co, c)
}

// GetContractDetails 合约详情
func GetContractDetails(c *gin.Context) {
	var r struct {
		ID int `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	var co model.Contract
	db.Where("id = ?", r.ID).First(&co)
	handle.ReturnSuccess("ok", co, c)
}

// ContractExtension 合约延期
func ContractExtension(c *gin.Context) {
	var r struct {
		ID        uint   `json:"id" binding:"required"`
		Extension string `json:"extension" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	var co model.Contract
	db := config.MysqlConn
	if err := db.Where("id = ?", r.ID).First(&co).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "合约不存在", c)
		return
	}
	dt, _ := time.ParseInLocation("2006-01-02", r.Extension, time.Local)
	co.ExpireTime = model.MyTime{Time: dt}
	if err := db.Save(&co).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "延期修改失败", c)
		return
	}

	handle.ReturnSuccess("ok", co, c)
}

// Task 七日任务
type Task struct {
	ID   int    `json:"id" binding:"required"`
	Task string `json:"task"`
}

// GetContractTask 获取合约七日任务任务
func GetContractTask(c *gin.Context) {
	var r Task
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	fmt.Println(r.ID)
	db := config.MysqlConn
	var co model.Contract
	if err := db.Where("id = ?", r.ID).First(&co).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约ID不正确", c)
		return
	}
	handle.ReturnSuccess("ok", co, c)
}

// UpdateContractTask 获取合约七日任务任务
func UpdateContractTask(c *gin.Context) {
	var r Task
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	var co model.Contract
	if err := db.Where("id = ?", r.ID).First(&co).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约ID不正确", c)
		return
	}
	co.Task = r.Task
	if err := db.Save(&co).Error; err != nil {
		handle.ReturnError(http.StatusInternalServerError, "七日任务修改失败", c)
		return
	}
	handle.ReturnSuccess("ok", co, c)
}

// GetContractByStatus 获取不同状态的合约 - 客户形式显示
func GetContractByStatus(c *gin.Context) {
	var r struct {
		Type        string `json:"type" binding:"required"`
		Date        string `json:"date"`
		Marshalling string `json:"marshalling"`
		Page        int    `json:"page"`
		Limit       int    `json:"limit"`
		Name        string `json:"name"`
		City        string `json:"city"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	var contracts []model.Contract
	sql := db.Preload("Member")
	var date string
	if r.Date != "" {
		date = r.Date
	} else {
		date = time.Now().Format("2006-01")
	}
	dt, _ := time.ParseInLocation("2006-01", date, time.Local)
	start := dt
	end := dt.AddDate(0, 1, 0)

	switch r.Type {
	case "newly": // 新签客户
		sql.Where("cooperation_time >= ? AND cooperation_time < ?", start, end).Where("sort = 0")
	case "inserve": // 服务中客户
		sql.Where("expire_time >= ?", end)
	case "beexpire": // 即将断约
		sql.Where("expire_time >= ? AND expire_time < ?", start, end)
	case "renewal": // 续约客户
		sql.Where("cooperation_time >= ? AND cooperation_time < ?", start, end).Where("sort > 0")
	case "break": // 断约客户
		if r.Date == time.Now().Format("2006-01") || r.Date == "" {
			end = time.Now()
		}
		sql.Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Where("expire_time >= ? AND expire_time < ?", start, end)
		}).Where("cooperation_time >= ? AND cooperation_time < ?", start, end)
	case "return": // 退款客户
		sql.Where("refund >= ? AND refund < ?", start, end)
	case "recycle": // 回收站
	default:

	}

	if r.City != "" {
		sql.Preload("Member", func(db *gorm.DB) *gorm.DB {
			return db.Where("city like ?", "%"+r.City+"%")
		})
	}

	if r.Marshalling != "" {
		sql.Where("operations_staff like ?", "%"+r.Marshalling+"%").
			Or("business_people like ?", "%"+r.Marshalling+"%")
	}

	if r.Name != "" {
		sql.Where("name like ?", "%"+r.Name+"%")
	}
	var count int64
	var pages int
	var page int
	_count := sql
	_count.Find(&contracts).Order("id desc").Count(&count)
	if r.Page == 0 {
		page = 1
	} else {
		page = r.Page
	}
	sql.Limit(10).Offset((page - 1) * 10).Find(&contracts)

	if int(count)%10 != 0 {
		pages = int(count)/10 + 1
	} else {
		pages = int(count) / 10
	}

	handle.ReturnSuccess("ok", gin.H{"contracts": contracts, "pages": pages, "currPage": page}, c)
}

// ContractRefund 合约退款
func ContractRefund(c *gin.Context) {
	var r struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	var contract model.Contract
	if err := db.Where("id = ?", r.ID).First(&contract).Error; err == nil {
		handle.ReturnError(http.StatusBadRequest, "合约ID不存在", c)
		return
	}
	contract.Refund = model.MyTime{Time: time.Now()}
	if err := db.Save(&contract).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "合约ID不存在", c)
		return
	}
	handle.ReturnSuccess("ok", contract, c)
}

// ChangeManagement 更换管理人员
func ChangeManagement(c *gin.Context) {
	var r struct {
		ID              uint `json:"id" binding:"required"`
		OperationsStaff int  `json:"operations_staff"` // 运营人员
		BusinessPeople  int  `json:"business_people"`  // 业务人员
	}
	if err := c.ShouldBind(&r); err != nil {
		handle.ReturnError(http.StatusBadRequest, "输入数据格式不正确", c)
		return
	}
	db := config.MysqlConn
	var m model.Member
	var u model.Contract
	if err := db.Where("id = ?", r.ID).First(&m).Error; err != nil {
		handle.ReturnError(http.StatusBadRequest, "客户ID不存在", c)
		return
	}
	if r.OperationsStaff != 0 {
		m.OperationsStaff = r.OperationsStaff
		var operationsStaff model.User
		db.Where("id = ?", r.OperationsStaff).First(&operationsStaff)
		u.OperationsStaff = operationsStaff.RealName
	}
	if r.BusinessPeople != 0 {
		m.BusinessPeople = r.BusinessPeople
		var businessPeople model.User
		db.Where("id = ?", r.BusinessPeople).First(&businessPeople)
		u.OperationsStaff = businessPeople.RealName
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&m).Error; err != nil {
			return err
		}

		if err := db.Model(&model.Contract{}).Where("member_id = ?", r.ID).Updates(u).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		handle.ReturnError(http.StatusBadRequest, "更新失败", c)
		return
	}
	handle.ReturnSuccess("ok", "更新成功", c)
}

// ALTER TABLE contracts ADD COLUMN `operations_staff` VARCHAR(10);
// ALTER TABLE members ADD COLUMN `refund` datetime(3);
