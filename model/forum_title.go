package model

type ForumTitle struct {
	ID               uint           `gorm:"primarykey" json:"id" binding:"required"`
	Title            string         `json:"title"`                             // 标题
	Content          string         `json:"content" gorm:"type:text"`          // 内容
	Type             uint8          `json:"type"`                              // 分类
	Recommended      bool           `json:"recommended"`                       // 是否推荐
	Pageviews        uint           `json:"pageviews"`                         // 浏览量
	TotalIntegral    uint8          `json:"total_integral"`                    // 总积分
	GivenIntegral    uint8          `json:"given_integral"`                    // 已给积分
	OriginalPoster   string         `json:"original_poster"`                   // 楼主
	OriginalPosterID uint           `json:"original_poster_id"`                // 楼主ID
	Like             uint           `json:"like"`                              // 点赞
	Contract         uint           `json:"contract"`                          // 关联合约
	Images           string         `json:"images"`                            // 图片
	Label            uint           `json:"label"`                             // 标签
	LabelText        string         `json:"label_text"`                        // 标签文本
	IsCarousel       bool           `json:"is_carousel"`                       // 是否轮播
	Comment          []ForumComment `json:"comment" gorm:"foreignKey:TitleID"` // 评论
	Likes            []ForumLike    `json:"likes" gorm:"foreignKey:TitleID"`   // 点赞
	CreatedAt        MyTime         `json:"created_at"`
	UpdatedAt        MyTime         `json:"updated_at"`
}
