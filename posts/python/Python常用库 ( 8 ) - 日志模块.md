```
{
    "url": "python-logging",
    "time": "2016/04/25 23:32",
    "tag": "Python"
}
```

# 一、简介

# 二、logging模块

```
import logging

# 输出到文件
fh = logging.FileHandler('log.txt', encoding='utf-8')

# 输出到屏幕
sh = logging.StreamHandler()

logging.basicConfig(level=logging.DEBUG, format="%(asctime)s - %(levelname)s - %(message)s", datefmt="%Y-%m-%d %H:%M:%S", handlers=[fh, sh])

logging.debug("debug log")
logging.info("info log")
logging.warning("warning log")
logging.error("error log")
logging.critical("critical log")
```

```
import logging

logger = logging.getLogger('loggername')
logger.setLevel(logging.DEBUG)
ft = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')

fh = logging.FileHandler('log.txt', encoding='utf-8')
fh.setLevel(logging.DEBUG)

# 控制输出到屏幕的错误等级
sh = logging.StreamHandler()
sh.setLevel(logging.INFO)

fh.setFormatter(ft)
sh.setFormatter(ft)

logger.addHandler(fh)
logger.addHandler(sh)

logger.debug("debug log")
logger.info("info log")
logger.warning("warning log")
logger.error("error log")
logger.critical("critical log")
```