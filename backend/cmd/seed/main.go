package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

type userSeed struct {
	openID   string
	nickname string
	school   string
}

var users = []userSeed{
	{"seed_ly001", "河科大老学长", "河南科技大学"},
	{"seed_ly002", "洛师小书虫", "洛阳师范学院"},
	{"seed_ly003", "理工程序猿", "洛阳理工学院"},
	{"seed_ly004", "推拿小能手", "河南推拿职业学院"},
	{"seed_ly005", "镜湖边的猫", "河南科技大学"},
	{"seed_ly006", "食堂干饭第一名", "河南科技大学"},
	{"seed_ly007", "图书馆守夜人", "洛阳师范学院"},
	{"seed_ly008", "操场慢跑者", "洛阳理工学院"},
	{"seed_ly009", "自习室钉子户", "河南科技大学"},
	{"seed_ly010", "奶茶续命中", "洛阳师范学院"},
	{"seed_ly011", "期末战神", "洛阳理工学院"},
	{"seed_ly012", "校园流浪猫", "河南推拿职业学院"},
	{"seed_ly013", "早八特困生", "河南科技大学"},
	{"seed_ly014", "走廊背书人", "洛阳师范学院"},
	{"seed_ly015", "宿舍养生家", "洛阳理工学院"},
	{"seed_ly016", "社团潜水员", "河南科技大学"},
	{"seed_ly017", "咸鱼翻身中", "洛阳师范学院"},
	{"seed_ly018", "ddl战士", "洛阳理工学院"},
	{"seed_ly019", "教室后排座", "河南推拿职业学院"},
	{"seed_ly020", "路灯下读书", "河南科技大学"},
}

var postContents = []string{
	"毕业季到了，好舍不得镜湖边的樱花。大学四年转瞬即逝，学弟学妹们要珍惜校园时光啊！",
	"今天图书馆六楼靠窗位置光线正好，学习效率超高！",
	"问：有学长学姐知道计算机二级考试怎么报名吗？",
	"食堂二楼的麻辣香锅太绝了，微辣刚刚好！",
	"今晚操场有社团音乐节，大家一起来嗨！",
	"刚跑完五公里，运动真的太解压了，推荐给大家",
	"有没有一起组队参加数学建模的？私聊我",
	"新开的那家奶茶店第二杯半价，求组队！",
	"分享一个自学Python的好资源：MIT 6.0001 公开课，B站就有",
	"学校的流浪猫又生了三只小奶猫，好可爱！",
	"求推荐洛阳好吃的火锅店，要性价比高的！",
	"英语六级倒计时30天，每天打卡背单词！",
	"周末去龙门石窟玩了一圈，作为洛阳人第一次去，真的很震撼",
	"早上八点的课真的太难了，冬天根本起不来",
	"分享一下实操课心得，给学弟学妹参考",
	"校园里的牡丹开了，拍了几张照片分享给大家！",
	"论如何在期末考试前一周速成一门课（亲身经历）",
	"报名了支教活动，暑假要去山区小学支教啦",
	"想找一个能一起打羽毛球的搭子，每周三次",
	"姐妹们！学校旁边新开了个自习室，环境超好",
	"实操考试通过了！开心！",
	"洛阳的天气真的是一天四季，早上穿袄中午穿短袖",
	"求推荐好看的纪录片，最近剧荒了",
	"学校后门的烧烤摊绝了！推荐羊肉串和烤茄子",
	"关于考研择校，有没有前辈能给点建议？",
}

var schoolList = []string{"河南科技大学", "洛阳师范学院", "洛阳理工学院", "河南推拿职业学院"}

func mustExec(tx *sql.Tx, query string, args ...any) {
	_, err := tx.Exec(query, args...)
	if err != nil {
		log.Fatalf("exec error: %v\n  query: %s", err, query[:min(len(query), 80)])
	}
}

func getID(tx *sql.Tx, table, whereCol, whereVal string) int64 {
	var id int64
	err := tx.QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE %s = $1", table, whereCol), whereVal).Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func seedUsers(tx *sql.Tx) map[string]int64 {
	ids := make(map[string]int64)
	for _, u := range users {
		existing := getID(tx, "users", "open_id", u.openID)
		if existing > 0 {
			ids[u.openID] = existing
			continue
		}
		err := tx.QueryRow(
			`INSERT INTO users (open_id, nickname, school, role) VALUES ($1,$2,$3,'user') RETURNING id`,
			u.openID, u.nickname, u.school,
		).Scan(&existing)
		if err != nil {
			log.Fatalf("insert user %s: %v", u.openID, err)
		}
		ids[u.openID] = existing
	}
	return ids
}

func seedPosts(tx *sql.Tx, uid map[string]int64) []int64 {
	var postIDs []int64
	for i, content := range postContents {
		idx := i % len(users)
		u := users[idx]
		school := u.school
		topicID := (i % 5) + 1
		var id int64
		createdAt := time.Now().Add(-time.Duration(25-i) * 24 * time.Hour)
		err := tx.QueryRow(
			`INSERT INTO posts (user_id, content, images, topic_id, school, like_count, comment_count, created_at)
			 VALUES ($1,$2,'[]'::jsonb,$3,$4,0,0,$5) RETURNING id`,
			uid[u.openID], content, topicID, school, createdAt,
		).Scan(&id)
		if err != nil {
			log.Printf("insert post %d: %v", i, err)
			continue
		}
		postIDs = append(postIDs, id)
	}
	return postIDs
}

func seedComments(tx *sql.Tx, uid map[string]int64, postIDs []int64) {
	commentDefs := []struct {
		postIdx int
		userIdx int
		content string
	}{
		{0, 1, "学长毕业快乐！常回来看看"},
		{0, 5, "镜湖的樱花真的美，每年春天必打卡"},
		{1, 0, "六楼靠窗是我的固定位置！"},
		{2, 6, "全国计算机等级考试官网报名，一般在6月"},
		{2, 7, "我也想问，蹲一个答案"},
		{4, 9, "音乐节几点开始？想去！"},
		{4, 10, "已经到操场了，气氛超棒"},
		{9, 11, "在哪在哪？我要去看小奶猫！"},
		{9, 12, "求照片！猫猫太可爱了"},
		{12, 13, "龙门石窟真的值得去，建议请个讲解"},
		{15, 14, "牡丹真国色！洛阳牡丹甲天下"},
		{17, 0, "加油！支教经历会是一生的财富"},
		{23, 11, "后门烧烤真的是一绝，毕业了还怀念"},
	}
	for _, c := range commentDefs {
		if c.postIdx >= len(postIDs) {
			continue
		}
		u := users[c.userIdx]
		mustExec(tx, `INSERT INTO comments (post_id, user_id, content, created_at) VALUES ($1,$2,$3,NOW())`,
			postIDs[c.postIdx], uid[u.openID], c.content)
	}
}

func seedLikes(tx *sql.Tx, uid map[string]int64, postIDs []int64) {
	pairs := []struct{ u, p int }{
		{0, 1}, {0, 4}, {1, 0}, {1, 4}, {5, 0}, {5, 1}, {6, 2}, {6, 4},
		{7, 0}, {7, 9}, {9, 4}, {9, 12}, {10, 9}, {11, 9}, {11, 15},
		{12, 12}, {13, 15}, {14, 0}, {15, 15}, {16, 23},
	}
	for _, p := range pairs {
		if p.p >= len(postIDs) {
			continue
		}
		u := users[p.u]
		mustExec(tx, `INSERT INTO likes (user_id, post_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			uid[u.openID], postIDs[p.p])
	}
	for _, pid := range postIDs {
		var cnt int
		tx.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id=$1", pid).Scan(&cnt)
		mustExec(tx, "UPDATE posts SET like_count=$1 WHERE id=$2", cnt, pid)
	}
}

func seedCircles(tx *sql.Tx, uid map[string]int64) []int64 {
	circleDefs := []struct {
		name        string
		desc        string
		creatorIdx  int
		memberCount int
	}{
		{"河科大小窝", "河南科技大学校园生活交流群", 0, 8},
		{"洛师文艺角", "洛阳师范学院文艺青年聚集地", 5, 4},
		{"理工编程社", "洛阳理工学院编程爱好者的乐园", 6, 4},
		{"推拿养生堂", "河南推拿职业学院交流圈", 9, 2},
		{"洛阳吃货联盟", "洛阳高校美食探索小分队", 1, 10},
		{"考研互助组", "考研路上互相鼓励一起上岸", 11, 4},
		{"摄影爱好者", "用镜头记录校园美好瞬间", 12, 3},
		{"运动健身圈", "跑步打球撸铁，一起动起来", 13, 5},
	}
	var circleIDs []int64
	for _, c := range circleDefs {
		u := users[c.creatorIdx]
		var id int64
		err := tx.QueryRow(
			`INSERT INTO circles (name, description, creator_id, member_count) VALUES ($1,$2,$3,0) RETURNING id`,
			c.name, c.desc, uid[u.openID],
		).Scan(&id)
		if err != nil {
			log.Printf("insert circle %s: %v", c.name, err)
			continue
		}
		circleIDs = append(circleIDs, id)
		mustExec(tx, `INSERT INTO circle_members (circle_id, user_id, role) VALUES ($1,$2,'owner')`, id, uid[u.openID])
	}
	return circleIDs
}

func seedCircleMembers(tx *sql.Tx, uid map[string]int64, circleIDs []int64) {
	memberMap := [][]int{
		{1, 3, 4, 7, 11, 13, 18},          // circle 0
		{3, 8, 12},                          // circle 1
		{2, 13, 16},                         // circle 2
		{17},                                // circle 3
		{0, 4, 5, 6, 8, 9, 10, 12, 13, 15}, // circle 4
		{7, 14, 18},                         // circle 5
		{12, 15},                            // circle 6
		{6, 16, 7, 13},                      // circle 7
	}
	for ci, members := range memberMap {
		if ci >= len(circleIDs) {
			continue
		}
		for _, mi := range members {
			u := users[mi]
			mustExec(tx, `INSERT INTO circle_members (circle_id, user_id, role) VALUES ($1,$2,'member') ON CONFLICT DO NOTHING`,
				circleIDs[ci], uid[u.openID])
		}
	}
}

func seedCirclePosts(tx *sql.Tx, uid map[string]int64, circleIDs []int64) {
	cpDefs := []struct {
		circleIdx int
		userIdx   int
		content   string
	}{
		{0, 0, "欢迎新同学加入河科大小窝！"},
		{0, 1, "图书馆新开了24小时自习室，方便考研党通宵学习了"},
		{0, 3, "学校的樱花大道要开始维护了，大家最近抓紧时间拍照"},
		{1, 5, "征集校园诗歌作品，优秀的会发表在洛师校报上"},
		{1, 8, "推荐一本好书：《你当像鸟飞往你的山》"},
		{2, 6, "周末组织了一个小型hackathon，有兴趣的报名！"},
		{2, 2, "分享一个Vue3的实战项目，适合新手练手"},
		{3, 9, "这周学习了颈椎推拿手法，感觉收获很大"},
		{4, 1, "洛阳必吃榜：小街锅贴、老雒阳面馆、管记水席"},
		{4, 4, "学校门口的煎饼果子加鸡蛋火腿，才6块钱！"},
		{5, 11, "今天开始每天打卡学习12小时，有没有一起的"},
		{5, 7, "求推荐英语一复习资料"},
		{6, 12, "发一组校园春景"},
		{7, 13, "今晚操场夜跑，有一起的吗？"},
		{7, 3, "打卡健身第30天，瘦了5斤！"},
	}
	for _, cp := range cpDefs {
		if cp.circleIdx >= len(circleIDs) {
			continue
		}
		u := users[cp.userIdx]
		mustExec(tx, `INSERT INTO circle_posts (circle_id, user_id, content) VALUES ($1,$2,$3)`,
			circleIDs[cp.circleIdx], uid[u.openID], cp.content)
	}
	for _, cid := range circleIDs {
		var cnt int
		tx.QueryRow("SELECT COUNT(*) FROM circle_posts WHERE circle_id=$1", cid).Scan(&cnt)
		mustExec(tx, "UPDATE circles SET post_count=$1 WHERE id=$2", cnt, cid)
		var mcnt int
		tx.QueryRow("SELECT COUNT(*) FROM circle_members WHERE circle_id=$1", cid).Scan(&mcnt)
		mustExec(tx, "UPDATE circles SET member_count=$1 WHERE id=$2", mcnt, cid)
	}
}

func seedProducts(tx *sql.Tx, uid map[string]int64) {
	type product struct {
		userIdx int
		title   string
		desc    string
		price   float64
		op      float64
		cat     string
		cond    string
	}
	items := []product{
		{11, "高等数学第七版（同济）", "九成新，只翻了几页", 25, 48, "教材", "九成新"},
		{12, "英语四级真题卷", "2024版全新，带解析", 15, 39.8, "教材", "九成新"},
		{13, "吉他入门款", "买了没弹几次，适合新手送教程书", 200, 450, "乐器", "轻微使用"},
		{14, "宿舍用小风扇", "去年夏天用的，风力大噪音小", 30, 79, "生活", "轻微使用"},
		{15, "考研数学复习全书", "2025版几乎全新送配套网课", 50, 128, "教材", "九五成新"},
		{16, "宿舍台灯", "LED护眼台灯三档调节", 40, 99, "生活", "轻微使用"},
		{18, "自行车出售", "毕业季出骑了两年车况良好", 150, 500, "出行", "轻微使用"},
		{0, "程序员面试宝典", "Java方向内容丰富很有帮助", 35, 79, "教材", "八五成新"},
		{1, "羽毛球拍一对", "尤尼克斯入门款送一筒球", 120, 299, "体育", "轻微使用"},
		{5, "素描工具包", "学画画买的含铅笔素描本橡皮等", 45, 120, "兴趣", "九成新"},
	}
	for _, p := range items {
		u := users[p.userIdx]
		mustExec(tx, `INSERT INTO products (user_id,title,description,price,original_price,category,condition)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			uid[u.openID], p.title, p.desc, p.price, p.op, p.cat, p.cond)
	}
}

func seedActivities(tx *sql.Tx, uid map[string]int64) {
	type activity struct {
		userIdx int
		title   string
		desc    string
		atype   string
		loc     string
		max     int
		start   string
		end     string
	}
	items := []activity{
		{0, "河科大校园音乐节", "音乐节乐队表演弹唱街舞", "娱乐", "河科大操场", 200, "2026-06-01 18:00:00+08", "2026-06-01 22:00:00+08"},
		{5, "洛师读书分享会", "每月读书分享活动本期主题东野圭吾", "学习", "洛师图书馆", 30, "2026-06-03 14:00:00+08", "2026-06-03 17:00:00+08"},
		{6, "编程入门讲座", "零基础Python入门面向全校", "学习", "理工教学楼A201", 50, "2026-06-05 15:00:00+08", "2026-06-05 17:00:00+08"},
		{9, "中医推拿体验日", "免费体验推拿按摩讲解日常养生", "公益", "推拿学院实训楼", 40, "2026-06-08 09:00:00+08", "2026-06-08 12:00:00+08"},
		{13, "校园篮球友谊赛", "洛阳理工vs河科大来给选手加油", "体育", "洛阳理工体育馆", 100, "2026-06-10 16:00:00+08", "2026-06-10 18:00:00+08"},
		{11, "考研经验分享会", "邀请已上岸学长学姐分享经验", "学习", "河科大图书馆报告厅", 80, "2026-05-28 19:00:00+08", "2026-05-28 21:00:00+08"},
		{12, "校园摄影大赛", "主题校园四季投稿评选", "娱乐", "线上投稿", 50, "2026-05-20 00:00:00+08", "2026-06-20 23:59:59+08"},
		{1, "洛阳美食探店活动", "小街天府吃小吃费用AA", "娱乐", "小街天府", 15, "2026-06-02 11:00:00+08", "2026-06-02 14:00:00+08"},
	}
	for _, a := range items {
		u := users[a.userIdx]
		mustExec(tx, `INSERT INTO activities (user_id,title,description,activity_type,location,max_participants,start_time,end_time,status)
			VALUES ($1,$2,$3,$4,$5,$6,$7::timestamptz,$8::timestamptz,0)`,
			uid[u.openID], a.title, a.desc, a.atype, a.loc, a.max, a.start, a.end)
	}
}

func seedLostItems(tx *sql.Tx, uid map[string]int64) {
	type item struct {
		userIdx  int
		title    string
		desc     string
		cat      string
		itemType string
		loc      string
		contact  string
	}
	items := []item{
		{7, "寻：校园卡（河科大）", "图书馆丢失", "lost", "证件", "河科大图书馆", "QQ:99999999"},
		{8, "捡到：黑色钱包", "食堂捡到黑色钱包有现金和学生证", "found", "钱包", "洛师一食堂", "电话:13800004444"},
		{9, "寻：蓝色保温杯", "操场看台丢失杯子上有贴纸", "lost", "水杯", "洛阳理工操场", "微信:test005"},
		{0, "捡到：白色充电宝", "教学楼A301捡到已交给楼管", "found", "电子", "河科大教学楼A301", "直接找楼管"},
		{11, "寻：黑色双肩包", "有笔记本电脑和书急", "lost", "包", "洛师图书馆", "电话:13800005555"},
		{12, "捡到：红色雨伞", "二号食堂门口捡到放失物招领处", "found", "雨伞", "二号食堂", "直接去失物招领处"},
		{13, "寻：蓝牙耳机", "AirPods二代充电仓有贴纸", "lost", "电子", "理工田径场", "微信:test006"},
	}
	for _, it := range items {
		u := users[it.userIdx]
		mustExec(tx, `INSERT INTO lost_items (user_id,title,description,category,item_type,location,contact)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			uid[u.openID], it.title, it.desc, it.cat, it.itemType, it.loc, it.contact)
	}
}

func seedWhispers(tx *sql.Tx, uid map[string]int64) {
	type whisper struct {
		userIdx int
		content string
	}
	items := []whisper{
		{7, "最近真的好累啊，感觉压力好大，但又不知道跟谁说"},
		{11, "今天图书馆看到一个小姐姐，好心动但不敢要微信"},
		{0, "还有一个月毕业了，真的不想离开学校"},
		{5, "室友每天晚上打游戏到两点，我该怎么办"},
		{13, "暗恋一个人三年了，现在毕业了可能再也见不到了"},
		{1, "大二了还没有方向，看着周围同学都那么目标明确，好焦虑"},
		{9, "表白被拒了，但至少我说出口了，不后悔"},
		{6, "今天体测1000米跑了4分05秒，比上学期进步了！"},
		{12, "真的很感谢我的室友，每次难过的时候都是她陪着我"},
		{15, "考过了！一个学期没白费，努力真的有用"},
		{14, "好想吃妈妈做的红烧肉啊，想家了"},
		{16, "有没有人和我一样，越到考试越不想复习"},
	}
	codenames := []string{"路过的一只猫", "图书馆常驻", "期末战神", "早八特困生", "操场慢跑者",
		"学术垃圾话", "路灯下读书", "摸鱼专业户", "宿舍守夜人", "ddl战士"}
	for i, w := range items {
		u := users[w.userIdx]
		cn := codenames[i%len(codenames)]
		createdAt := time.Now().Add(-time.Duration(12-i) * 24 * time.Hour)
		mustExec(tx, `INSERT INTO whispers (user_id, content, is_anonymous, codename, whisper_type, created_at)
			VALUES ($1,$2,true,$3,0,$4)`,
			uid[u.openID], w.content, cn, createdAt)
	}
}

func seedFollows(tx *sql.Tx, uid map[string]int64) {
	pairs := [][2]int{
		{0, 1}, {0, 3}, {0, 5}, {0, 11}, {1, 0}, {1, 3}, {1, 6},
		{3, 0}, {3, 1}, {3, 4}, {5, 0}, {5, 8}, {5, 12},
		{6, 1}, {6, 2}, {6, 13}, {7, 0}, {7, 11}, {7, 14},
		{8, 5}, {8, 3}, {9, 0}, {9, 3}, {10, 1}, {10, 6},
		{11, 0}, {11, 1}, {11, 7}, {12, 5}, {12, 8},
		{13, 6}, {13, 1}, {14, 0}, {14, 11}, {15, 5},
		{16, 6}, {16, 13}, {18, 0}, {18, 11},
	}
	for _, p := range pairs {
		f := users[p[0]]
		g := users[p[1]]
		mustExec(tx, `INSERT INTO follows (follower_id, following_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			uid[f.openID], uid[g.openID])
	}
}

func seedFriends(tx *sql.Tx, uid map[string]int64) {
	pairs := [][2]int{
		{0, 1}, {0, 3}, {0, 11}, {1, 0}, {1, 3}, {1, 6},
		{3, 0}, {3, 1}, {5, 8}, {5, 12}, {6, 1}, {6, 13},
		{7, 11}, {7, 14}, {9, 0}, {9, 3}, {12, 5}, {12, 8},
		{13, 6}, {13, 1}, {14, 0}, {14, 11},
	}
	for _, p := range pairs {
		f := users[p[0]]
		g := users[p[1]]
		mustExec(tx, `INSERT INTO friends (user_id, friend_id, status) VALUES ($1,$2,'accepted') ON CONFLICT DO NOTHING`,
			uid[f.openID], uid[g.openID])
		mustExec(tx, `INSERT INTO friends (user_id, friend_id, status) VALUES ($1,$2,'accepted') ON CONFLICT DO NOTHING`,
			uid[g.openID], uid[f.openID])
	}
}

func seedMessages(tx *sql.Tx, uid map[string]int64) {
	type msg struct {
		from int
		to   int
		c    string
	}
	convs := [][]msg{
		{{0, 1, "你好！看到你也在河科大，交个朋友吧"}, {1, 0, "哈哈好啊，学长好！"}, {0, 1, "你哪个专业的？"}, {1, 0, "计算机专业的，学长呢？"}, {0, 1, "我也是计院的！缘分啊"}},
		{{5, 8, "你的奶茶店组队贴我看到了，一起去啊"}, {8, 5, "好呀好呀，下午三点？"}, {5, 8, "可以，我下课了去找你"}},
	}
	for _, conv := range convs {
		for _, m := range conv {
			f := users[m.from]
			t := users[m.to]
			mustExec(tx, `INSERT INTO messages (sender_id, receiver_id, content, created_at)
				VALUES ($1,$2,$3,NOW() - interval '1 hour' * random() * 48)`,
				uid[f.openID], uid[t.openID], m.c)
		}
	}
}

func seedNotifications(tx *sql.Tx, uid map[string]int64) {
	notifs := []struct {
		to      int
		from    int
		ntype   string
		content string
	}{
		{0, 1, "follow", "小明关注了你"},
		{0, 3, "friend", "镜湖边的猫请求加你为好友"},
		{1, 0, "like", "河科大老学长赞了你的帖子"},
		{5, 9, "follow", "奶茶续命中关注了你"},
		{6, 1, "follow", "小明关注了你"},
	}
	for _, n := range notifs {
		to := users[n.to]
		from := users[n.from]
		mustExec(tx, `INSERT INTO notifications (user_id, from_user_id, type, content)
			VALUES ($1,$2,$3,$4)`,
			uid[to.openID], uid[from.openID], n.ntype, n.content)
	}
}

func seedBanners(tx *sql.Tx) {
	banners := []string{"河科大音乐节", "考研经验分享会", "校园摄影大赛", "二手市集开业啦"}
	for i, b := range banners {
		mustExec(tx, `INSERT INTO banners (title, image_url, sort_order, is_active)
			VALUES ($1,'',$2,true) ON CONFLICT DO NOTHING`, b, i+1)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/university_social?sslmode=disable")
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("begin: %v", err)
	}
	defer tx.Rollback()

	log.Println("Seeding users...")
	uid := seedUsers(tx)

	log.Println("Seeding posts...")
	postIDs := seedPosts(tx, uid)
	log.Printf("  %d posts", len(postIDs))

	log.Println("Seeding comments...")
	seedComments(tx, uid, postIDs)

	log.Println("Seeding likes...")
	seedLikes(tx, uid, postIDs)

	log.Println("Seeding circles...")
	circleIDs := seedCircles(tx, uid)
	log.Printf("  %d circles", len(circleIDs))

	log.Println("Seeding circle members...")
	seedCircleMembers(tx, uid, circleIDs)

	log.Println("Seeding circle posts...")
	seedCirclePosts(tx, uid, circleIDs)

	log.Println("Seeding products...")
	seedProducts(tx, uid)

	log.Println("Seeding activities...")
	seedActivities(tx, uid)

	log.Println("Seeding lost items...")
	seedLostItems(tx, uid)

	log.Println("Seeding whispers...")
	seedWhispers(tx, uid)

	log.Println("Seeding follows...")
	seedFollows(tx, uid)

	log.Println("Seeding friends...")
	seedFriends(tx, uid)

	log.Println("Seeding messages...")
	seedMessages(tx, uid)

	log.Println("Seeding notifications...")
	seedNotifications(tx, uid)

	log.Println("Seeding banners...")
	seedBanners(tx)

	if err := tx.Commit(); err != nil {
		log.Fatalf("commit: %v", err)
	}
	log.Println("Seed complete!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
