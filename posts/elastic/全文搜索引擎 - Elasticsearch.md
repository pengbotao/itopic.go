```
{
    "url": "elasticsearch",
    "time": "2018/12/01 16:46",
    "tag": "Elasticsearch,ELK",
    "toc": "yes"
}
```

# 一、概述

Elasticsearch是一个基于Lucene的搜索服务器。它提供了一个分布式多用户能力的全文搜索引擎，基于RESTful web接口。Elasticsearch是用Java语言开发的，并作为Apache许可条款下的开放源码发布，是一种流行的企业级搜索引擎。Elasticsearch用于云计算中，能够达到实时搜索，稳定，可靠，快速，安装使用方便。官方客户端在Java、.NET（C#）、PHP、Python、Apache Groovy、Ruby和许多其他语言中都是可用的。根据DB-Engines的排名显示，Elasticsearch是最受欢迎的企业搜索引擎，其次是Apache Solr，也是基于Lucene。

> 当前文档使用Elasticsearch版本：6.8.12

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
$ tar -zxvf elasticsearch-6.8.12.tar.gz -C /usr/local/
```

## 2.3 启动Elasticsearch

直接使用如下命令即可启动：

```
$ /usr/local/elasticsearch-6.8.12/bin/elasticsearch
```

但可能出现以下情况：

- [1]: max file descriptors [65535] for elasticsearch process is too low, increase to at least [65536]
- [2]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
- [3]: max number of threads [1024] for user [es] is too low, increase to at least [4096]
- [4]: java.lang.RuntimeException: can not run elasticsearch as root

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

**3. 线程数**

```
$ vi /etc/security/limits.d/90-nproc.conf
$ cat /etc/security/limits.d/90-nproc.conf
# Default limit for number of user's processes to prevent
# accidental fork bombs.
# See rhbz #432903 for reasoning.

*          soft    nproc     102400
root       soft    nproc     unlimited
```

**4. 切换账户**

提示不能以`root`用户启动，需切换到其他用户，放后台执行可以通过`-d`参数：

```
$ /usr/local/elasticsearch-6.8.12/bin/elasticsearch -d
```

`Elasticsearch`以HTTP的方式提供服务，后续的操作也都是以HTTP请求方式交互，先来测试服务是否启动成功：

```
$ curl 127.0.0.1:9200
{
  "name" : "Vxhzj3I",
  "cluster_name" : "elasticsearch",
  "cluster_uuid" : "cBmW_29HRTCIRMj9YGCSeA",
  "version" : {
    "number" : "6.8.12",
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

上面示例配置文件在：`/usr/local/elasticsearch-6.8.12/config/elasticsearch.yml`

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
    "number" : "6.8.12",
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

**3.1.1 查看索引**

```
$ curl localhost:9200/_cat/indices?v
health status index uuid                   pri rep docs.count docs.deleted store.size pri.store.size
yellow open   area  8A9sotMOQN-UuirmEmOF9Q   5   1       3447            0      757kb          757kb

$ curl localhost:9200/_cat/
=^.^=
/_cat/allocation
/_cat/shards
/_cat/shards/{index}
/_cat/master
/_cat/nodes
/_cat/tasks
/_cat/indices
/_cat/indices/{index}
...
```

**3.1.2 创建索引**

创建名为`demo`的索引：

```
$ curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/demo'
{
    "acknowledged": true,
    "shards_acknowledged": true,
    "index": "demo"
}
```

**创建`Index`**，同时指定`Type=_doc`，定义`mapping`：包含2个字段：`title`、`content`。

```
$ cat << EOF | curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/demo?pretty' -d @- 
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
  "index" : "demo"
}
```

**3.1.3 删除索引**

**删除`Index`：**

```
$ curl -X DELETE 'http://localhost:9200/demo'
{
    "acknowledged": true
}
```

## 3.2 类型（Type）

类型介于索引和文档之间，相当于对文档进行分组。一个索引库下有多个分组，每个分组下存储多个文档，就可以按Mysql的`库-表-记录`来理解了。类型包含类型名称和`mapping`，`mapping`用来定义字段的属性。

`Elasticsearch`6.x版本每个`Index`只允许包含一个`Type`（推荐名称是`_doc`），7.x版本将会彻底移除`Type`。在测试版本中提交不同`Type`的数据时会报错：

```
{
    "type": "illegal_argument_exception",
    "reason": "Rejecting mapping update to [demo] as the final mapping would have more than 1 type: [_doc, test]"
}
```

**3.2.1 读取mapping**

列出`demo`索引下的`mapping`，由于目前只支持1对1，如果早期版本有多个`Type`，可以在索引后指定`Type`名称。

```
$ curl localhost:9200/demo/_mapping?pretty
{
  "demo" : {
    "mappings" : {
      "_doc" : {
        "properties" : {
          "content" : {
            "type" : "text"
          },
          "title" : {
            "type" : "keyword"
          }
        }
      }
    }
  }
}
```

**3.2.2 更新mapping**

`mapping`在创建索引时或者第一次提交数据时确定，创建后存在的字段将无法更改，但可以向已有的`mapping`里新增字段。

```
cat << EOF | curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/demo/_doc/_mapping?pretty' -d @- 
{
    "properties": {
        "author": {
            "type": "text"
        }
    }
}
EOF

{
  "acknowledged" : true
}
```

**3.2.3 索引迁移**

前面说到已经提交的字段无法修改，但可以通过起个别名来迁移，分为5个步骤：

**1. 给已有的索引起个别名**

```
$ cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/_aliases?pretty' -d @-
{
    "actions": [
        { "add" : { "index" : "demo", "alias" : "test" } }
    ]
}
EOF

{
  "acknowledged" : true
}
```

**2. 创建新的索引**

这里指定使用ik分词，需要安装后面的`ik`分词插件。因为这是新建的，所以类型等字段就随自定义了。

```
$ cat << EOF | curl -X PUT -H "Content-type: application/json" 'http://localhost:9200/beta?pretty' -d @- 
{
    "mappings": {
        "_doc": {
            "properties": {
                "title": {
                    "type" : "keyword",
                    "analyzer": "ik_max_word",
                    "search_analyzer": "ik_smart"
                },
                "content": {
                    "type": "text",
                    "analyzer": "ik_max_word",
                    "search_analyzer": "ik_smart"
                }
            }
        }
    }
}
EOF
```

**3. 重建新索引**

```
$ cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/_reindex?pretty' -d @-
{
    "source": {
        "index": "demo"
    },
    "dest": {
        "index": "beta"
    }
}
EOF

{
  "took" : 44,
  "timed_out" : false,
  "total" : 1,
  "updated" : 0,
  "created" : 1,
  "deleted" : 0,
  "batches" : 1,
  "version_conflicts" : 0,
  "noops" : 0,
  "retries" : {
    "bulk" : 0,
    "search" : 0
  },
  "throttled_millis" : 0,
  "requests_per_second" : -1.0,
  "throttled_until_millis" : 0,
  "failures" : [ ]
}
```

**4. 将别名指向到新的索引**

```
$ cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/_aliases?pretty' -d @-
{
    "actions": [
        { "add" : { "index" : "beta", "alias" : "test" } }
    ]
}
EOF

{
  "acknowledged" : true
}
```

**5. 删除旧索引**

一个别名可以关联多个索引，经过上一步后则别名关联了2个索引，删除旧的索引别名。

```
$ cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/_aliases?pretty' -d @-
{
    "actions": [
        {"remove": {"index": "demo", "alias": "test"}}
    ]
}
EOF

{
    "acknowledged": true
}
```

如果在定义索引之前没有设置别名，通过该方式代码端还是需要有一些调整，但如果有定义别名并使用别名，则只需做ES服务端的迁移即可。

## 3.3 文档（Document）

类似`Mysql`的行数据，一行数据就是一个文档，数据格式通常使用`JSON`格式，同一个索引+类型下可以存储多个文档数据，可以通过接口操作文档数据。

**3.3.1 添加数据**

给`demo`索引下`Type`为`_doc`添加数据：

```
$ cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/demo/_doc/10001?pretty' -d @-
{
    "title": "11月将发生这些大事！第一条就与你有关",
    "content": "全国人口普查开始正式登记;养老机构将建立老年人健康档案"
}
EOF

{
  "_index" : "demo",
  "_type" : "_doc",
  "_id" : "10001",
  "_version" : 1,
  "result" : "created",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 0,
  "_primary_term" : 1
}
```

上面的地址为：`demo/_doc/10001`，10001为指定记录的ID，也可以不指定ID，ES会自己生成。

```
$ cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/demo/_doc/?pretty' -d @-
{
    "title": "11月将发生这些大事！第一条就与你有关",
    "content": "全国人口普查开始正式登记;养老机构将建立老年人健康档案"
}
EOF

{
  "_index" : "demo",
  "_type" : "_doc",
  "_id" : "yc_BeHUB5PnW6VuaGLjA",
  "_version" : 1,
  "result" : "created",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 1,
  "_primary_term" : 1
}
```

**3.3.2 获取数据**

获取`Index=demo`, `Type=_doc`, `_id=10001`的记录。

```
$ curl 'http://localhost:9200/demo/_doc/10001?pretty'
{
  "_index" : "demo",
  "_type" : "_doc",
  "_id" : "10001",
  "_version" : 2,
  "_seq_no" : 1,
  "_primary_term" : 1,
  "found" : true,
  "_source" : {
    "title" : "11月将发生这些大事！第一条就与你有关",
    "content" : "全国人口普查开始正式登记;养老机构将建立老年人健康档案"
  }
}
```

**3.3.3 更新数据**

将前面数据指定ID再次提交即可，更新之后版本号变为了2.

```
$ cat << EOF | curl -X POST -H "Content-type: application/json" 'http://localhost:9200/demo/_doc/10001?pretty' -d @-
{
    "title": "11月将发生这些大事！第一条就与你有关",
    "content": "全国人口普查开始正式登记;养老机构将建立老年人健康档案"
}
EOF

{
  "_index" : "demo",
  "_type" : "_doc",
  "_id" : "10001",
  "_version" : 2,
  "result" : "updated",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 1,
  "_primary_term" : 1
}
```

**3.3.4 删除数据**

```
$ curl -X DELETE 'http://localhost:9200/demo/_doc/yc_BeHUB5PnW6VuaGLjA?pretty'
{
  "_index" : "demo",
  "_type" : "_doc",
  "_id" : "yc_BeHUB5PnW6VuaGLjA",
  "_version" : 2,
  "result" : "deleted",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "_seq_no" : 1,
  "_primary_term" : 1
}
```

删除之后再次执行查询操作，可以看到`found=false`

```
$ curl 'http://localhost:9200/demo/_doc/yc_BeHUB5PnW6VuaGLjA?pretty'

{
  "_index" : "demo",
  "_type" : "_doc",
  "_id" : "yc_BeHUB5PnW6VuaGLjA",
  "found" : false
}
```

# 四、中文分词

创建或更新文档时会对文档进行分词，查询时会对查询语句进行分词，分词的效果会最终影响到查询的结果，比如中文字符在默认规则下会按照单个汉字来分：`苹` + `果` + `笔` + `记` + `本`

```
cat << EOF | curl -H "Content-type: application/json" 'http://localhost:9200/area/_analyze?pretty' -d @-
{
    "text": "苹果笔记本",
    "analyzer": "standard"
}
EOF

{
  "tokens" : [
    {
      "token" : "苹",
      "start_offset" : 0,
      "end_offset" : 1,
      "type" : "<IDEOGRAPHIC>",
      "position" : 0
    },
    ...
}
```

所以，需要对中文分词的支持，通常可用：ik分词、smartcn[<sup>[3]</sup>](#refer)。

## 4.1 安装ik中文分词

Github上选择对应版本即可：`https://github.com/medcl/elasticsearch-analysis-ik/releases`

```
$ ./bin/elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v6.8.12/elasticsearch-analysis-ik-6.8.12.zip
-> Downloading https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v6.8.12/elasticsearch-analysis-ik-6.8.12.zip
[=================================================] 100%
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@     WARNING: plugin requires additional permissions     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
* java.net.SocketPermission * connect,resolve
See http://docs.oracle.com/javase/8/docs/technotes/guides/security/permissions.html
for descriptions of what these permissions allow and the associated risks.

Continue with installation? [y/N]y
-> Installed analysis-ik
```

安装会需要重启。

## 4.2 ik分析器

有两种分词方式：

`ik_max_word`: 会将文本做最细粒度的拆分，比如会将“中华人民共和国国歌”拆分为“中华人民共和国,中华人民,中华,华人,人民共和国,人民,人,民,共和国,共和,和,国国,国歌”，会穷尽各种可能的组合，适合 Term Query；

`ik_smart`: 会做最粗粒度的拆分，比如会将“中华人民共和国国歌”拆分为“中华人民共和国,国歌”，适合 Phrase 查询。

来看看ik的两种示例：用`ik_smart`方式，苹果笔记本分成了 `苹果` + `笔记本`，用`ik_max_word`方式分成了： `苹果` + `笔记本` + `笔记` + `本`

```
$ cat << EOF | curl -H "Content-type: application/json" 'http://localhost:9200/area/_analyze?pretty' -d @-
{
    "text": "苹果笔记本",
    "analyzer": "ik_smart"
}
EOF

{
  "tokens" : [
    {
      "token" : "苹果",
      "start_offset" : 0,
      "end_offset" : 2,
      "type" : "CN_WORD",
      "position" : 0
    },
    {
      "token" : "笔记本",
      "start_offset" : 2,
      "end_offset" : 5,
      "type" : "CN_WORD",
      "position" : 1
    }
  ]
}

$ cat << EOF | curl -H "Content-type: application/json" 'http://localhost:9200/area/_analyze?pretty' -d @-
{
    "text": "苹果笔记本",
    "analyzer": "ik_max_word"

}
EOF

{
  "tokens" : [
    {
      "token" : "苹果",
      "start_offset" : 0,
      "end_offset" : 2,
      "type" : "CN_WORD",
      "position" : 0
    },
    {
      "token" : "笔记本",
      "start_offset" : 2,
      "end_offset" : 5,
      "type" : "CN_WORD",
      "position" : 1
    },
    {
      "token" : "笔记",
      "start_offset" : 2,
      "end_offset" : 4,
      "type" : "CN_WORD",
      "position" : 2
    },
    {
      "token" : "本",
      "start_offset" : 4,
      "end_offset" : 5,
      "type" : "CN_CHAR",
      "position" : 3
    }
  ]
}
```

# 五、索引定义

## 5.1 索引类型

### 5.1.1 Keyword

- keyword: 通常用来存储结构化数据，比如ID、Email、状态码、标签等
- constant_keyword: 始终包含相同值的字段
- wildcard: 通配符

Keyword字段常用来做排序、聚合、Term级别查询，避免将keyword用于全文搜索，全文搜索可以使用text类型。

### 5.1.2 Text

会进行全文索引的字段，会对字段进行分词，然后索引，可以模糊查询。

# 六、查询语句

## 6.1 查询列表

```
{
  "query": {
    "match_all": {}
  }
}
```

## 6.2 精确查询 - Term

Term查询用于查询确定的值。相当于`Where Name = peng`。

```
{
  "query": {
    "term": {"name": "Peng"}
  }
}
```

相当于`Where Name in ('peng', 'lion')`

```
{
  "query": {
    "terms": {"name": ["Peng", "Lion"]}
  }
}
```

但需要注意的是被查询的字段会否会分词，默认`text`会分词、忽略大小写，这就可能造成无法返回结果。文本类型有`text`和`keyword`类型，其中`keyword`取代了不需要分词的`string`。比如：

```
"name" : {
  "type" : "text",
  "fields" : {
    "keyword" : {
      "type" : "keyword",
      "ignore_above" : 256
    }
  }
},
```

上面例子如果要精确匹配可以使用下面的`keyword`：`name.keyword`。

## 6.3 全文查询 - Match

会进行分词后查询，默认为或的操作。

```
{
  "query": {
    "match": {"name": "Bobby Peng"}
  }
}
```

可调整为`and`

```
{
  "query": {
     "match": {
        "name": {
          "query": "Bobby Peng",
          "operator" : "and"
       }
      }
    }
}
```

match还有2个变种：`match_phrase` 和 `multi_match`，match_phrase用于同时满足所有词，同时也可以增加slop调整少一个也行。

```
{
  "query": {
     "match_phrase": {
        "name": {
          "query": "Bobby Peng",
          "slop" : 1
       }
      }
    }
}
```

multi_match用于多个字段匹配，有一个字段满足即可

```
{
  "query": {
     "multi_match": {
          "query": "Bobby Peng",
          "fields" : ["name", "title"]
      }
    }
}
```

## 6.4 范围查询

```
{
  "query": {
      "range": {
          "age": {
              "gte": 18,
              "lt": 30 
          }
      }
  }
}
```



## 6.5 复合查询

bool 复合查询用于组合叶子查询语句或复合查询语句。如：must, should, must_not, or filter。

- must 必须匹配。
- should 至少匹配一个文档。
- filter 必须匹配，忽略相关性评分。
- must_not 必须不匹配，忽略相关性评分。

```
GET /_search
{
  "query": { 
    "bool": { 
      "must": [
        { "match": { "title":   "Search"}},
        { "match": { "content": "Elasticsearch" }}
      ],
      "filter": [ 
        { "term":  { "status": "published"}},
        { "range": { "publish_date": { "gte": "2015-01-01" }}}
      ]
    }
  },
  "from": 0,
  "size": 10,
  "sort": {
    "timestamp": "desc"
  }
}
```

---

<div id="refer"></div>

- [1] [Elasticsearch Reference 6.8](https://www.elastic.co/guide/en/elasticsearch/reference/6.8/index.html)
- [2] [elasticsearch-analysis-ik](https://github.com/medcl/elasticsearch-analysis-ik)
- [3] [Smart Chinese Analysis Plugin](https://www.elastic.co/guide/en/elasticsearch/plugins/current/analysis-smartcn.html)
- [4] [全文搜索引擎 Elasticsearch 入门教程](http://www.ruanyifeng.com/blog/2017/08/elasticsearch.html)
- [5] [Elasticsearch: analyzer](https://www.cnblogs.com/sanduzxcvbnm/p/12084607.html)
- [6] [ES Mapping、字段类型Field type详解](https://blog.csdn.net/ZYC88888/article/details/83059040)
- [7] [Elasticsearch DSL 查询详解](https://blog.csdn.net/lamp_yang_3533/article/details/97618687)
