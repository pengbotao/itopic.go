```
{
    "url": "es-start",
    "time": "2018/12/01 16:46",
    "tag": "Elasticsearch",
    "toc": "yes",
    "public": "no"
}
```

# 一、概述

Elasticsearch是一个基于Lucene的搜索服务器。它提供了一个分布式多用户能力的全文搜索引擎，基于RESTful web接口。Elasticsearch是用Java语言开发的，并作为Apache许可条款下的开放源码发布，是一种流行的企业级搜索引擎。Elasticsearch用于云计算中，能够达到实时搜索，稳定，可靠，快速，安装使用方便。官方客户端在Java、.NET（C#）、PHP、Python、Apache Groovy、Ruby和许多其他语言中都是可用的。根据DB-Engines的排名显示，Elasticsearch是最受欢迎的企业搜索引擎，其次是Apache Solr，也是基于Lucene。

# 二、安装

Elasticsearch使用Java开发，需要配置JDK环境。

## 2.1 安装jdk

**1. 下载安装包**

https://www.oracle.com/java/technologies/javase/javase-jdk8-downloads.html

**2. 解压缩以及配置环境变量**

```
$ tar -zxvf jdk-8u271-linux-x64.tar.gz -C /usr/local/

$ vi /etc/profile
JAVA_HOME=/usr/local/jdk1.8.0_271/
PATH=$PATH:$JAVA_HOME/bin
CLASSPATH=$JAVA_HOME/jre/lib/ext:$JAVA_HOME/lib/tools.jar
export PATH JAVA_HOME CLASSPATH
```

**3. 刷新环境变量**

```
$ source /etc/profile
```

**4. 测试**

```
$ java -version
java version "1.8.0_271"
Java(TM) SE Runtime Environment (build 1.8.0_271-b09)
Java HotSpot(TM) 64-Bit Server VM (build 25.271-b09, mixed mode)
```

## 2.2 安装Elasticsearch

`Elasticsearch`的安装比较简单，下载对应的包直接运行即可。

**1. 下载安装包**

https://www.elastic.co/cn/downloads/past-releases#elasticsearch

**2. 解压缩**

```
$ tar -zxvf elasticsearch-6.8.13.tar.gz -C /usr/local/
```

## 2.3 启动Elasticsearch

直接使用如下命令即可启动：

```
$ /usr/local/elasticsearch-6.8.13/bin/elasticsearch
```

但可能出现以下情况：

- [1]: max file descriptors [65535] for elasticsearch process is too low, increase to at least [65536]
- [2]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
- [3]: java.lang.RuntimeException: can not run elasticsearch as root

**1. 调整文件句柄数**

```
$ vi /etc/security/limits.conf
* soft nofile 65536
* hard nofile 65536

$ ulimit -Hn
65536
```

**2. vm.max_map_count**

```
$ vi /etc/sysctl.conf
vm.max_map_count=655360
$ sysctl -p
```

**3. 切换账户**

提示不能以`root`用户启动，需切换到其他用户，放后台执行可以通过`-d`参数：

```
$ /usr/local/elasticsearch-6.8.13/bin/elasticsearch -d
```

`Elasticsearch`以HTTP的方式提供服务，后续的操作也都是以HTTP请求方式交互，先来测试服务是否启动成功：

```
$ curl 127.0.0.1:9200
{
  "name" : "Vxhzj3I",
  "cluster_name" : "elasticsearch",
  "cluster_uuid" : "cBmW_29HRTCIRMj9YGCSeA",
  "version" : {
    "number" : "6.8.13",
    "build_flavor" : "default",
    "build_type" : "tar",
    "build_hash" : "be13c69",
    "build_date" : "2020-10-16T09:09:46.555371Z",
    "build_snapshot" : false,
    "lucene_version" : "7.7.3",
    "minimum_wire_compatibility_version" : "5.6.0",
    "minimum_index_compatibility_version" : "5.0.0"
  },
  "tagline" : "You Know, for Search"
}
```

停止服务：

```
$ kill -SIGTERM pid
```

## 2.4 配置文件

上面示例配置文件在：`/usr/local/elasticsearch-6.8.13/config/elasticsearch.yml`

| 配置项                             | 说明                                              | 示例                 |
| ---------------------------------- | ------------------------------------------------- | -------------------- |
| cluster.name                       | 集群名称，默认是elasticsearch，同名称才能组成集群 |                      |
| node.name                          | 节点名称                                          |                      |
| node.master                        | 是否有资格为master，默认为true                    |                      |
| node.data                          | 是否存储数据，默认为true                          |                      |
| path.data                          | 数据目录，多个目录逗号分隔                        |                      |
| path.logs                          | 日志目录                                          |                      |
| bootstrap.memory_lock              | 锁定物理内存                                      |                      |
| network.host                       | 绑定主机                                          |                      |
| http.port                          | 绑定端口                                          |                      |
| http.enabled                       | 是否启用HTTP服务                                  | false                |
| discovery.zen.ping.unicast.hosts   | 单播主机列表                                      | `["host1", "host2"]` |
| discovery.zen.minimum_master_nodes | 要选举Master至少需要多少个节点同意，默认为1       |                      |

比如，调整集群和节点名称后：

```
$ curl 127.0.0.1:9200
{
  "name" : "es-demo-node1",
  "cluster_name" : "es-demo",
  "cluster_uuid" : "cBmW_29HRTCIRMj9YGCSeA",
  "version" : {
    "number" : "6.8.13",
    "build_flavor" : "default",
    "build_type" : "tar",
    "build_hash" : "be13c69",
    "build_date" : "2020-10-16T09:09:46.555371Z",
    "build_snapshot" : false,
    "lucene_version" : "7.7.3",
    "minimum_wire_compatibility_version" : "5.6.0",
    "minimum_index_compatibility_version" : "5.0.0"
  },
  "tagline" : "You Know, for Search"
}
```

# 三、基本概念

## 3.1 索引（Index）

索引可以理解成索引库，是相同结构的文档数据的集合。

列出当前节点下所有的`Index`。这里只创建了一个`index`，名称为`area`。存储的是省市区的联动数据。

```
$ curl localhost:9200/_cat/indices?v
health status index uuid                   pri rep docs.count docs.deleted store.size pri.store.size
yellow open   area  U9zBRy56SoqB5COncvQR2w   5   1       1440            1    275.8kb        275.8kb
```

### 3.1.1 创建索引

创建索引的方式如下，创建名为`demo`的索引：

```
$ curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/demo'
{
    "acknowledged": true,
    "shards_acknowledged": true,
    "index": "demo"
}
```

**创建`Index`**，同时指定`Type`

```
$ cat << EOF | curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/demo2?pretty' -d @- 
{
    "mappings": {
        "_doc": {
            "properties": {
                "title": {
                    "type" : "keyword"
                },
                "content": {
                    "type": "text"
                }
            }
        }
    }
}
EOF

{
  "acknowledged" : true,
  "shards_acknowledged" : true,
  "index" : "demo2"
}
```

### 3.1.2 删除索引

**删除`Index`：**

```
$ curl -X DELETE 'http://localhost:9200/demo'
{
    "acknowledged": true
}
```

## 3.2 类型（Type）

介于索引和文档之间，相当于对文档进行分组。一个索引库下有多个分组，每个分组下存储多个文档，就可以按Mysql的库-表-记录来理解了。

列出`area`索引下所有的`Type`，也可以去掉`area`。`/_mapping`查看所有，`/area/_mapping`查看area索引下的Type。

```
$ curl localhost:9200/area/_mapping?pretty
{
  "area" : {
    "mappings" : {
      "demo" : {
        "properties" : {
          "code" : {
            "type" : "text",
            "fields" : {
              "keyword" : {
                "type" : "keyword",
                "ignore_above" : 256
              }
            }
          },
          "name" : {
            "type" : "text",
            "fields" : {
              "keyword" : {
                "type" : "keyword",
                "ignore_above" : 256
              }
            }
          },
          "parent" : {
            "type" : "text",
            "fields" : {
              "keyword" : {
                "type" : "keyword",
                "ignore_above" : 256
              }
            }
          },
          "pinyin" : {
            "type" : "text",
            "fields" : {
              "keyword" : {
                "type" : "keyword",
                "ignore_above" : 256
              }
            }
          }
        }
      }
    }
  }
}
```

`Elasticsearch`6.x版本每个`Index`只允许包含一个`Type`，7.x版本将会彻底移除`Type`。在测试版本中提交不同`Type`的数据时会报错：

```
{
    "type": "illegal_argument_exception",
    "reason": "Rejecting mapping update to [area] as the final mapping would have more than 1 type: [demo, demo2]"
}
```

## 3.3 文档（Document）

类似Mysql的行数据，一行数据就是一个文档，数据格式通常使用JSON格式，同一个索引+类型下可以存储多个文档数据，可以通过接口操作文档数据。

### 3.3.1 添加数据

给`area`索引下`Type`为`demo`添加数据：

```
$ cat << EOF | curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/area/demo/10001?pretty' -d @-
{
    "code": "10001",
    "name": "测试城市",
    "parent": "0",
    "pinyin": "ceshichengshi"
}
EOF

{
  "_index" : "area",
  "_type" : "demo",
  "_id" : "10001",
  "_version" : 1,
  "result" : "created",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 295,
  "_primary_term" : 4
}
```

上面的地址为：`area/demo/10001`，10001为指定记录的ID，也可以不指定ID，ES会自己生成，请求方法需要调整为`POST`。

```
cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/area/demo?pretty' -d @-
{
    "code": "10002",
    "name": "测试城市2",
    "parent": "0",
    "pinyin": "ceshichengshi2"
}
EOF

{
  "_index" : "area",
  "_type" : "demo",
  "_id" : "Eja7bnUBQ99FRyN-H2AJ",
  "_version" : 1,
  "result" : "created",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 286,
  "_primary_term" : 4
}
```

### 3.3.2 获取数据

获取`Index=area`, `Type=demo`, `_id=10001`的记录。

```
$ curl 'http://localhost:9200/area/demo/10001?pretty'
{
  "_index" : "area",
  "_type" : "demo",
  "_id" : "10001",
  "_version" : 1,
  "_seq_no" : 295,
  "_primary_term" : 4,
  "found" : true,
  "_source" : {
    "code" : "10001",
    "name" : "测试城市",
    "parent" : "0",
    "pinyin" : "ceshichengshi"
  }
}
```

### 3.3.3 更新数据

将`4.3`的数据指定ID再次提交即可，更新之后版本号变为了2.

```
$ cat << EOF | curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/area/demo/10001?pretty' -d @-
{
     "code": "10001",
     "name": "测试城市-update",
     "parent": "0",
     "pinyin": "ceshichengshi"
}
EOF

{
  "_index" : "area",
  "_type" : "demo",
  "_id" : "10001",
  "_version" : 2,
  "result" : "updated",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 296,
  "_primary_term" : 4
}
```

### 3.3.4 删除数据

```
$ curl -X DELETE 'http://localhost:9200/area/demo/10001?pretty'
{
  "_index" : "area",
  "_type" : "demo",
  "_id" : "10001",
  "_version" : 3,
  "result" : "deleted",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 294,
  "_primary_term" : 4
}
```

删除之后再次执行查询操作，可以看到`found=false`

```
$ curl 'http://localhost:9200/area/demo/10001?pretty'
{
  "_index" : "area",
  "_type" : "demo",
  "_id" : "10001",
  "found" : false
}
```

### 3.3.5 搜索数据

```
$ cat << EOF | curl -H "Content-type: application/json" 'http://localhost:9200/area/_search?pretty' -d @-
{
    "query": {
        "match": {
            "name": "房山"
        }
    },
    "from": 0,
    "size": 2
}
EOF

{
  "took" : 8,
  "timed_out" : false,
  "_shards" : {
    "total" : 5,
    "successful" : 5,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : 93,
    "max_score" : 8.325781,
    "hits" : [
      {
        "_index" : "area",
        "_type" : "demo",
        "_id" : "110111",
        "_score" : 8.325781,
        "_source" : {
          "code" : "110111",
          "name" : "房山区",
          "pinyin" : "fangshanqu",
          "parent" : "1101"
        }
      },
      {
        "_index" : "area",
        "_type" : "demo",
        "_id" : "230108",
        "_score" : 5.666169,
        "_source" : {
          "code" : "230108",
          "name" : "平房区",
          "pinyin" : "pingfangqu",
          "parent" : "2301"
        }
      }
    ]
  }
}
```

# 四、查询语句



# 五、中文分词 - ik

```
$ ./bin/elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v6.8.12/elasticsearch-analysis-ik-6.8.12.zip
```

# 六、倒排索引



# 七、集群管理

## 7.1 分布式

分片

副本



## 7.2 服务发现





- [1] [全文搜索引擎 Elasticsearch 入门教程](http://www.ruanyifeng.com/blog/2017/08/elasticsearch.html)
- [2] [Elasticsearch Reference 6.8](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/index.html)

