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

The fantastic ORM library for Golang aims to be developer friendly.

## 1.1 安装

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

## 1.2 连接

```
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```

## 1.3 调试模式

```
db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger:logger.Default.LogMode(logger.Info),
})
```

或者

```
db.Debug().First(&user)
```

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

标签定义，更多可[参考官网](https://gorm.io/zh_CN/docs/models.html#%E5%AD%97%E6%AE%B5%E6%A0%87%E7%AD%BE)

| 标签名     | 说明             |
| :--------- | :--------------- |
| column     | 指定 db 列名     |
| primaryKey | 将列定义为主键   |
| unique     | 将列定义为唯一键 |
| default    | 定义列的默认值   |

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

**3. 按字段删除**

```
# DELETE FROM `user` WHERE age < 3
db.Where("age < ?", 3).Delete(&User{})
```

## 3.4 查询数据

GORM 提供了 `First`、`Take`、`Last` 方法，以便从数据库中检索单个对象。当查询数据库时它添加了 `LIMIT 1` 条件，且没有找到记录时，它会返回 `ErrRecordNotFound` 错误

**1. 查询单条记录**

```
// 获取一条记录，没有指定排序字段
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// 获取第一条记录（主键升序）
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// 获取最后一条记录（主键降序）
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

result := db.First(&user)
result.RowsAffected // 返回找到的记录数
result.Error        // returns error or nil

// 检查 ErrRecordNotFound 错误
errors.Is(result.Error, gorm.ErrRecordNotFound)

//查询1条，但是不抛错
db.Limit(1).Find(&user)
```

**2. 查询多条记录**

```
users := []User{}
db.Find(&users)
```

也可以设置一些筛选条件，通过链式操作

```
# SELECT `user_id`,`nickname` FROM `user` WHERE id > 0 ORDER BY user_id desc LIMIT 10 OFFSET 0
db.Where("user_id > ?", 0).
	Select("user_id", "nickname").
	Order("user_id desc").
	Limit(10).Offset(0).
	Find(&user)
```

**3. 查询Count**

```
var total int64 = 0
# SELECT count(1) FROM `user`
db.Model(User{}).Count(&total)
```

**4. 分组Group**

```
type result struct {
	Age   int
	Count int
}
var res result
db.Model(User{}).Select("age, count(*) as count").Group("age").Scan(&res)
```

## 3.5 原生SQL

原生查询 SQL 和 `Scan`

```
type Result struct {
  ID   int
  Name string
  Age  int
}

var result Result
db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&result)

db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&result)

var age int
db.Raw("SELECT SUM(age) FROM users WHERE role = ?", "admin").Scan(&age)

var users []User
db.Raw("UPDATE users SET name = ? WHERE age = ? RETURNING id, name", "jinzhu", 20).Scan(&users)
```

`Exec` 原生 SQL

```
db.Exec("DROP TABLE users")
db.Exec("UPDATE orders SET shipped_at = ? WHERE id IN ?", time.Now(), []int64{1, 2, 3})

// Exec with SQL Expression
db.Exec("UPDATE users SET money = ? WHERE name = ?", gorm.Expr("money * ? + ?", 10000, 1), "jinzhu")
```

