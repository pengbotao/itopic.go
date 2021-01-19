```
{
    "url": "mongo-install",
    "time": "2018/12/10 06:37",
    "tag": "MongoDB"
}
```

通过官方下载社区版本（MongoDB Community Server），选择对应平台下载后解压即可看到可执行程序。

```
$ ls ~/data/mongodb-macos-x86_64-4.4.3/
LICENSE-Community.txt MPL-2                 README                THIRD-PARTY-NOTICES   bin

$ ll ~/data/mongodb-macos-x86_64-4.4.3/bin/
total 320808
-rwxr-xr-x@ 1 peng  staff      7683 12 22 07:42 install_compass
-rwxr-xr-x  1 peng  staff  43147504 12 22 07:39 mongo
-rwxr-xr-x  1 peng  staff  68692572 12 22 07:40 mongod
-rwxr-xr-x  1 peng  staff  52399600 12 22 07:32 mongos
```

方便访问以设置bin目录到环境变量，创建配置文件`mongod.conf`：

```
systemLog:
    destination: file
    path: "/Users/peng/data/mongodb-macos-x86_64-4.4.3/logs/mongodb.log"
    logAppend: true
storage:
    dbPath: "/Users/peng/data/mongodb-macos-x86_64-4.4.3/data/"
    journal:
        enabled: true
net:
    port: 27017
    bindIp: 127.0.0.1
processManagement:
    fork: true
    pidFilePath: "/Users/peng/data/mongodb-macos-x86_64-4.4.3/logs/mongod.pid"

```

启动服务

```
$ mongod -f ~/data/mongodb-macos-x86_64-4.4.3/mongod.conf
about to fork child process, waiting until server is ready for connections.
forked process: 49268
child process started successfully, parent exiting
```

停止服务

```
$ mongo
MongoDB shell version v4.4.3
connecting to: mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb
Implicit session: session { "id" : UUID("ef69ef5a-9a83-45f4-a1ca-b63717183535") }
MongoDB server version: 4.4.3
switched to db admin
> db
admin
> db.shutdownServer()
server should be down...
```

通过官方安装MongoDB GUI管理工具：MongoDB Compass，效果图如下：

![](../../static/uploads/mongodb-compass.png)

# 常用命令

**MongoDB与Mysql结构对比**

| MYSQL         | MongoDB           | 说明                             |
| ------------- | ----------------- | -------------------------------- |
| 库 - database | 库 - database     |                                  |
| 表 - table    | 集合 - collection | 表示一组文档                     |
| 行 - row      | 文档 - document   | 基本单元，类似关系型数据库中的行 |

## 增删查改

**查询**

```
# 显示数据库列表
> show dbs
> show databases

# 显示库中的集合
> show tables
> show collections

# 切换数据库
> use demo

> db.articles.find().pretty()
{
	"_id" : ObjectId("6006394a1e1382e154503c71"),
	"title" : "Hbase配置及数据迁移",
	"href" : "https://itopic.org/hbase.html"
}
{
	"_id" : ObjectId("6006396c1e1382e154503c72"),
	"title" : "Apache Airflow数据库迁移",
	"href" : "https://itopic.org/airflow-data-migration.html",
	"top" : 1
}
{
	"_id" : ObjectId("6006398f1e1382e154503c73"),
	"title" : "通过kubeadm部署Kubernetes集群",
	"href" : "https://itopic.org/k8s-kubeadm-install.html",
	"top" : 0
}
{
	"_id" : ObjectId("60063a931e1382e154503c74"),
	"title" : "Python入门知识点整理",
	"href" : "https://itopic.org/python.html",
	"top" : 1,
	"author" : {
		"name" : "peng",
		"age" : 18
	}
}

> db.articles.find({top:1}).pretty()
{
	"_id" : ObjectId("6006396c1e1382e154503c72"),
	"title" : "Apache Airflow数据库迁移",
	"href" : "https://itopic.org/airflow-data-migration.html",
	"top" : 1
}

> db.articles.find({'author.name':'peng'}).pretty()
{
	"_id" : ObjectId("60063a931e1382e154503c74"),
	"title" : "Python入门知识点整理",
	"href" : "https://itopic.org/python.html",
	"top" : 1,
	"author" : {
		"name" : "peng",
		"age" : 18
	}
}

> db.articles.count()
4
```

**数据插入**

```
> db.articles.insert({title: "Shell脚本入门", href: "https://itopic.org/shell-start.html"})
WriteResult({ "nInserted" : 1 })
```

**数据更新**

```
> db.articles.update({title: "Shell脚本入门"}, {title: "Shell脚本入门", href: "", top:0})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
```

**数据删除**

```
> db.articles.remove({title: "Shell脚本入门"})
WriteResult({ "nRemoved" : 1 })
```



---

- [1] [MongoDB官网](https://www.mongodb.com)
- [2] [MongoDB数据库系列](https://www.kancloud.cn/noahs/linux/1425612)