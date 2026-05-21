package model

import "time"

type Whisper struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	Images       []string  `json:"images"`
	IsAnonymous  bool      `json:"is_anonymous"`
	Codename     string    `json:"codename,omitempty"`
	WhisperType  int       `json:"whisper_type"`
	TargetUserID *int64    `json:"target_user_id,omitempty"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	IsLiked      bool      `json:"is_liked,omitempty"`
}

type CreateWhisperRequest struct {
	Content      string   `json:"content" binding:"required"`
	Images       []string `json:"images"`
	IsAnonymous  bool     `json:"is_anonymous"`
	WhisperType  int      `json:"whisper_type"`
	TargetUserID *int64   `json:"target_user_id"`
}

type WhisperComment struct {
	ID        int64     `json:"id"`
	WhisperID int64     `json:"whisper_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	Codename  string    `json:"codename,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type WhisperCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

var CodenamePool = []string{
	"路过的一只猫", "今晚不熬夜", "图书馆常驻", "食堂干饭人",
	"早八特困生", "学术垃圾话", "奶茶续命中", "ddl战士",
	"明天一定早睡", "课代表本尊", "宿舍守夜人", "校园流浪猫",
	"操场听歌人", "期末战神", "摸鱼专业户", "走廊等风来",
	"教室后排座", "实验室幽灵", "快递站常客", "校门口卖瓜",
	"阳台晾袜子", "食堂加鸡腿", "图书馆守门", "开水房歌王",
	"草坪晒太阳", "路灯下读书", "操场慢跑者", "天台吹风人",
	"黑板擦守护", "自习室钉子", "雨伞交换生", "电动车漂移",
	"食堂品鉴师", "选修课休眠", "打印机杀手", "路演围观者",
	"校园外卖侠", "咸鱼翻身中", "摆烂艺术家", "教学楼的猫",
	"镜湖喂鱼人", "篮球场饮水", "宿舍养生家", "教材传教士",
	"课桌涂鸦师", "空调遥控器", "走廊背书人", "周末留校生",
	"社团潜水员", "后排伏案者", "深夜树洞回", "晨跑永远鸽",
}
