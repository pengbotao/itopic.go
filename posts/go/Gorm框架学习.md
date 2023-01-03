```
{
    "url": "gorm",
    "time": "2022/12/31 21:06",
    "tag": "Golang",
    "toc": "yes",
    "public": "no"
}
```

# 一 、概述



# 二、模型定义

假设数据库中的表结构如下：

```
CREATE TABLE IF NOT EXISTS `user` (
    user_id INT(4) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    mobile VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
    nickname VARCHAR(50) NOT NULL DEFAULT '' COMMENT '昵称',
    age INT(4) NOT NULL DEFAULT 0 COMMENT '年龄',
    intro VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '简介',
    created_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (user_id)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT '用户记录表';
```

转换为模型后定义如下，

```
type User struct {
	UserId    int       `gorm:"column:id;PRIMARY_KEY"`
	Mobile    string    `gorm:"column:mobile"`
	Nickname  string    `gorm:"column:nickname"`
	Age       int       `gorm:"column:age"`
	Intro     string    `gorm:"column:intro"`
	CreatedTs time.Time `gorm:"column:created_ts"`
}

func (u User) TableName() string {
	return "user"
}
```

# 三、增删查改

## 3.1 数据插入

**1. 插入数据**

```
user := User{
	Mobile:    "13800000001",
	Nickname:  "demo",
	Age:       18,
	Intro:     "Hello World",
	CreatedTs: time.Now(),
}
res := db.Create(&user)
fmt.Println(res.Error, res.RowsAffected, user.UserId)
```

**2. 插入时SQL只使用Select中的字段插入**

```
db.Select("mobile", "nickname").Create(&user).Error
```

**3. 插入时SQL上忽略created_ts字段**

```
db.Omit("created_ts").Create(&user).Error
```

**4. 批量插入**，且插入时忽略created_ts字段。

```
users := []User{
	{Mobile: "13800000001"},
	{Mobile: "13800000002"},
	{Mobile: "13800000003"},
}
return db.Omit("created_ts").Create(&users).Error
```

批量插入时分批插入，每批次执行batchSize条：

```
db.CreateInBatches(&users, 2).Error
```

