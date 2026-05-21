CREATE TABLE IF NOT EXISTS sensitive_words (
    id          BIGSERIAL PRIMARY KEY,
    word        VARCHAR(128) NOT NULL UNIQUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO sensitive_words (word) VALUES
('赌博'), ('毒品'), ('色情'), ('暴力'), ('恐怖'),
('枪支'), ('走私'), ('诈骗'), ('传销'), ('邪教'),
('代考'), ('替课'), ('作弊'), ('卖答案'), ('代写论文'),
('约炮'), ('找小姐'), ('裸聊'), ('一夜情'), ('招嫖'),
('办证'), ('假证'), ('发票'), ('套现'), ('洗钱')
ON CONFLICT (word) DO NOTHING;
