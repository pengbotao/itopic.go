```
{
    "url": "python-excel",
    "time": "2017/03/18 23:32",
    "tag": "Python"
}
```

常用的操作`Excel`的库有：`openpyxl`、`xlsxwriter`、`xlrd/xlwt`、`xlutils`，从`python-excel.org`上可以看到一些区别，主流应该还是`openpyxl`：


Packages|说明
---|---
`openpyxl`|The `recommended package` for reading and writing Excel 2010 files (ie: .xlsx)
`xlsxwriter`|XlsxWriter is a Python module for writing files in the Excel 2007+ XLSX file format.
`xlrd`|This library currently has no active maintainers. You are advised to use `OpenPyXL` instead. 
`xlwt`|This package is for writing data and formatting information to `older Excel files` (ie: .xls)
`xlutils`|对`xlrd/xlwt`的封装，NB: In general, these use cases are now covered by openpyxl!

**Excel `xls`和`xlsx`的区别在于：**

> 1、核心结构上：xls 是一个特有的二进制格式，其核心结构是复合文档类型的结构，而 xlsx 的核心结构是 XML 类型的结构，采用的是基于 XML 的压缩方式，使其占用的空间更小。xlsx 中最后一个 x 的意义就在于此。

> 2、版本上：xls是2003版本下的文件 ，不管有没有宏程序的话都是xls文件 ，从2007开始做了区分，XLSM文件和XLSX文件都是excel2007及其以后的文件，但前者是含有宏启用，Excel中默认情况下不自动启用宏，默认是XLSX。VBA中，如果不想保存代码，可以保存为xlsx，即可自动删除其中VBA代码，反之则保存为XLSM文件。

# 一、OpenPyXL

- 项目地址：`https://bitbucket.org/openpyxl/openpyxl`
- 文档地址：`https://openpyxl.readthedocs.io/en/stable`

## 1.1 安装及示例
```
$ pip install openpyxl
```

To be able to include images (jpeg, png, bmp,…) into an openpyxl file, you will also need the `pillow` library that can be installed with:

```
$ pip install pillow
```

## 1.2 创建工作表

```
from openpyxl import Workbook

wb = Workbook()

# 默认工作表
ws1 = wb.active
ws1.title = "工作表-1"
# A1 即表示A1单元格
ws1["A1"] = "A1 Hello World"

ws2 = wb.create_sheet("工作表-2")
ws2["A1"] = "A2 Hello World"

# 创建工作表3并排列在工作表的第一位
ws3 = wb.create_sheet("工作表-3", 0)
ws3["A1"] = "A3 Hello World"
wb.save("sample.xlsx")
```

## 1.3 写入工作表

```
from openpyxl import Workbook

wb = Workbook()

# 默认工作表
ws1 = wb.active
ws1.title = "工作表-1"
ws2 = wb.create_sheet("工作表-2")
ws3 = wb.create_sheet("工作表-3", 0)


row = ["ID", "Name", "Birth"]
ws1.append(row)

for r in [
    [1, "Judy", "2020-01-01"],
    [2, "Jason", "2022-01-01"]
]:
    ws1.append(r)

wb.save("sample.xlsx")
```

## 1.4 读取工作表

```
from openpyxl import load_workbook

wb = load_workbook("sample.xlsx")

sheets = wb.get_sheet_names()
print(sheets) # ['工作表-3', '工作表-1', '工作表-2']

sheet1 = wb['工作表-1']
print(sheet1.title) # 工作表-1

# 获取工作表-1的A1单元格的数据
print(sheet1["A2"].value)

# 通过cell(row, column)方式获取A1的数据
print(sheet1.cell(2, 1).value)
```