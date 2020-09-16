```
{
    "url": "geonames",
    "time": "2019/02/19 06:52",
    "tag": "小知识"
}
```

# 一、概述

GeoNames是一个免费的全球地理数据库。GeoNames的目标是把各种来源的免费数据进行集成并制作成一个数据库或一系列的Web服务。

# 二、接口文档

- API说明：https://www.geonames.org/export/
- API列表：https://www.geonames.org/export/ws-overview.html

Terms and Conditions

- **free** : GeoNames data is free, the data is available without costs.
- **cc-by licence** (creative commons attributions license). You should give credit to GeoNames when using data or web services with a link or another reference to GeoNames.
- commercial usage is allowed
- **'as is'** : The data is provided "as is" without warranty or any representation of accuracy, timeliness or completeness.
- **20'000** credits daily limit per application (identified by the parameter 'username'), the hourly limit is 1000 credits. A credit is a web service request hit for most services. An exception is thrown when the limit is exceeded.
- **Service Level Agreement** is available for our premium web services.


# 三、 根据经纬度获取时区(Timezone)

## 3.1 接口调用

http://api.geonames.org/timezoneJSON?lat=13.724712798375322&lng=100.63311079999994&username=xxx


返回示例：
```
{
	• sunrise: "2020-02-20 06:43",
	• lng: -74.00552399999998,
	• countryCode: "US",
	• gmtOffset: -5,
	• rawOffset: -5,
	• sunset: "2020-02-20 17:36",
	• timezoneId: "America/New_York",
	• dstOffset: -4,
	• countryName: "United States",
	• time: "2020-02-19 20:28",
	• lat: 40.713425
}
```

字段|说明
---|---
countryCode| ISO countrycode
countryName|name (language can be set with param lang)
timezoneId|name of the timezone (according to olson), this information is sufficient to work with the timezone and defines DST rules, consult the documentation of your development environment. Many programming environments include functions based on the olson timezoneId (example java TimeZone)
time|the local current time
sunset|sunset local time (date)
sunrise|sunrise local time (date)
rawOffset|the amount of time in hours to add to UTC to get standard time in this time zone. Because this value is not affected by daylight saving time, it is called raw offset.
gmtOffset|offset to GMT at 1. January (deprecated)
dstOffset|offset to GMT at 1. July (deprecated)


## 3.2 夏令时说明

夏令时，表示为了节约能源，人为规定时间的意思。也叫夏时制，夏时令（Daylight Saving Time：DST），又称“日光节约时制”和“夏令时间”，在这一制度实行期间所采用的统一时间称为“夏令时间”。一般在天亮早的夏季人为将时间调快一小时，可以使人早起早睡，减少照明量，以充分利用光照资源，从而节约照明用电。各个采纳夏时制的国家具体规定不同。目前全世界有近110个国家每年要实行夏令时。

1986年4月，中国中央有关部门发出“在全国范围内实行夏时制的通知”，具体作法是：每年从四月中旬第一个星期日的凌晨2时整（北京时间），将时钟拨快一小时，即将表针由2时拨至3时，夏令时开始；到九月中旬第一个星期日的凌晨2时整（北京夏令时），再将时钟拨回一小时，即将表针由2时拨至1时，夏令时结束。

从1986年到1991年的六个年度，除1986年因是实行夏时制的第一年，从5月4日开始到9月14日结束外，其它年份均按规定的时段施行。在夏令时开始和结束前几天，新闻媒体均刊登有关部门的通告。1992年起，夏令时暂停实行。

## 3.3 GMT-格林尼治标准时

GMT 的全名是格林威治标准时间或格林威治平时 （Greenwich Mean Time），这个时间系统的概念在 1884 年确立，由英国伦敦的格林威治皇家天文台计算并维护，并在往后的几十年往欧陆其他国家扩散。在 1924 年开始，格林威治天文台每小时就会向全世界播报时间。

在刚开始的几十年，GMT 的测量方法非常简单：观测者随时监控太阳在天空的位置，并且把每天太阳爬升到仰角最高的时候记录下来，这个时间点称呼为“过中天”。一般人对于一天 24 小时的理解，大致上就相等于两次太阳过中天的时间间隔。不过由于地球是以椭圆轨道绕着太阳，在轨道上的行进速率不一，导致一年之中会有“比较长的一天”与“比较短的一天”，所以格林威治的观测者必须要至少连续观测一年，然后求取 365 个长度不一的“天”，再把他们全部平均后，得到固定的一天长度，之后再细分成时、分、秒等单位。这个就是 GMT。

GMT 12:00 就是指的是英国伦敦郊区的皇家格林尼治天文台当地的中午12:00，而GMT+8 12:00，则是指的东八区的北京当地时间的12:00。

## 3.4 UTC-协调世界时

自从 1967 年国际度量衡大会把秒的定义改成铯原子进行固定震荡次数的时间后，时间的测量就可以与星球的自转脱节了。只利用原子钟计算时间与日期的系统，称作国际原子时 （International Atomic Time），这是一种只有“天”的系统，时分秒都以“天”的小数点零头来表示。以国际原子时为计算基准，把时间格式与 UT1 对齐，让一般人都方便使用的时间系统，就叫做协调世界时 （Universal Time Coordinated），也就是 UTC。这也就是 UTC 为什么与 GMT 几乎一样的关係。由于 UTC 直接与国际度量衡标准相联繫，所以目前所有的国际通讯系统，像是卫星、航空、GPS 等等，全部都协议採用 UTC 时间。
协调世界时，又称世界统一时间、世界标准时间、国际协调时间。协调世界时，即以我为基准，向我看齐的意思。（英语：Coordinated Universal Time，法语：Temps Universel Coordonné，简称UTC）是最主要的世界时间标准，由于英文（CUT）和法文（TUC）的缩写不同，作为妥协，简称UTC。
