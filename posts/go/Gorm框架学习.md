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

## 3.2 数据更新

**1. 更新单列**

插入表名或者结构体

```
# UPDATE `user` SET `mobile`='13800000002' WHERE id = 1
db.Table("user").Where("id = ?", 1).Update("mobile", "13800000002")
db.Model(&User{}).Where("id = ?", 1).Update("mobile", "13800000002")

db.Model(&user).Update("mobile", "13800000002")
```

**2. 更新多列**

`Updates` 方法支持 `struct` 和 `map[string]interface{}` 参数。当使用 `struct` 更新时，默认情况下，GORM 只会更新非零值的字段

```
db.Model(&user).Updates(User{Age: 1})
```

如果上面年龄传0则不会做任何更新，因为0为整型的零值。此时可以使用map来更新。

```
# UPDATE `user` SET `age`=0 WHERE `user_id` = 1
db.Model(&user).Updates(map[string]interface{}{"age": 0})
```

同时也支持`Select()`、`Omit()`、`Table("user")`、`Table(&User{})`的方式，如：

```
# UPDATE `user` SET `age`=99 WHERE user_id in (1,2,3) AND user_id < 3
db.Model(&User{}).Select("age").
	Where("user_id in ?", []int{1, 2, 3}).
	Where("user_id < 3").
	Updates(User{Nickname: "Test", Age: 99})
```

**3. 存在更新，不存在插入**

如果主键为零值则执行插入，否则执行更新操作。

```
db.Save(&user)
```

## 3.3 删除数据

**1. 按主键删除**

```
# DELETE FROM `user` WHERE `user`.`user_id` = 1
db.Delete(&User{}, 1)

# DELETE FROM `user` WHERE `user`.`user_id` IN (1,2,3)
db.Delete(&User{}, []int{1, 2, 3})
```

**2. 传入结构体删除**

```
# DELETE FROM `user` WHERE `user`.`user_id` = 1
db.Delete(&user)
```

