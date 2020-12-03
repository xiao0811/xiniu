package model

// Member 客户结构体
// 审核的状态都是  0待审核  1审核通过  2审核驳回 3删除
type Member struct {
	ID                uint       `gorm:"primarykey" json:"id"`
	UUID              string     `json:"uuid" gorm:"type:char(26);unique"`           // 用户编号
	Name              string     `gorm:"type:varchar(50)" json:"name"`               // 门店名称
	City              string     `json:"city" gorm:"type:varchar(10)"`               // 所在城市
	FirstCategory     string     `json:"first_category" gorm:"type:varchar(20)"`     // 一级类目
	SecondaryCategory string     `json:"secondary_category" gorm:"type:varchar(20)"` // 二级类目
	BusinessScope     string     `json:"business_scope" gorm:"type:varchar(20)"`     // 主营范围
	Stores            uint8      `json:"stores"`                                     // 门店数量
	Accounts          uint8      `json:"accounts"`                                   // 账户数量
	Bosses            uint8      `json:"bosses"`                                     // 老板人数
	Brands            uint8      `json:"brands"`                                     // 品牌数量
	OperationsGroup   int        `json:"operations_group"`                           // 运营组
	OperationsStaff   int        `json:"operations_staff"`                           // 运营人员
	BusinessGroup     int        `json:"business_group"`                             // 业务组
	BusinessPeople    int        `json:"business_people"`                            // 业务人员
	ReviewAccount     string     `json:"review_account" gorm:"type:varchar(30)"`     // 点评账号
	CommentPassword   string     `json:"comment_password" gorm:"type:varchar(30)"`   // 点评密码
	Email             string     `json:"email" gorm:"type:varchar(30)"`              // 客户邮箱
	Phone             string     `json:"phone" gorm:"type:varchar(11)"`              // 客户手机号码
	OtherTags         string     `json:"other_tags"`                                 // 其他标签
	ReviewTime        MyTime     `json:"review_time"`                                // 审核时间
	Auditors          uint       `json:"auditors"`                                   // 审核人员
	Type              int8       `json:"type"`                                       // 备注信息
	Contracts         []Contract `json:"contracts"`
	FirstCreate       MyTime     `json:"first_creat"`         // 第一次创建合约时间
	ExpireTime        MyTime     `json:"expire_time"`         // 到期时间
	NumberOfContracts uint8      `json:"number_of_contracts"` // 合约数
	Refund            MyTime     `json:"refund"`              // 退款
	Status            int8       `json:"status"`
	Remarks           string     `json:"remarks"`
	Reason            string     `json:"reason"` // 审核备注
	CreatedAt         MyTime     `json:"created_at"`
	UpdatedAt         MyTime     `json:"updated_at"`
}

// ALTER TABLE members ADD COLUMN `first_create` datetime(3);
// ALTER TABLE members ADD COLUMN `expire_time` datetime(3);
// ALTER TABLE members ADD COLUMN `number_of_contracts` tinyint unsigned;
// ALTER TABLE members ADD COLUMN `refund` datetime(3);
