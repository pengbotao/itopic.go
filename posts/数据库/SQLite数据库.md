```
{
    "url": "sqlite",
    "time": "2017/03/03 22:22",
    "tag": "数据库,SQLite"
}
```

# 一、简介

`SQLite`是一个零配置的关系型数据库，支持大部分`SQL`语句。就像Linux系统下会自带Python一样，通常也会安装`Sqlite3`，可以通过`sqlite3`命令来确认是否已经安装：

```
$ sqlite3
SQLite version 3.24.0 2018-06-04 14:10:15
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
sqlite>
```

虽然`SQLite`支持的功能大部分`Mysql`都有，但对比`SQLite`就会发现，`SQLite`小巧、零配置、移植方便、不需要额外启动服务端进程、功能也相当完善，较擅长在一些独立项目上提供本地存储，本纯文本方式方便，比`Mysql`清爽。

安装上可直接从[官网](https://www.sqlite.org/download.html)上下载，相关文档可从[SQLite TuTorial](http://www.sqlitetutorial.net/)上查看。操作工具可以直接使用`命令行`或者[SQLite Studio](https://sqlitestudio.pl/index.rvt?act=download)或者`Navicat`。

# 二、SQLite



