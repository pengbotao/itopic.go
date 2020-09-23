```
{
    "url": "linux-convert",
    "time": "2015/03/26 17:50",
    "tag": "Linux"
}
```

# 安装ImageMagick
```
# yum install ImageMagick ImageMagick-devel
 
# convert -version
Version: ImageMagick 6.2.8 05/07/12 Q16 file:/usr/share/ImageMagick-6.2.8/doc/index.html
Copyright: Copyright (C) 1999-2006 ImageMagick Studio LLC
```
同时也可以安装PHP扩展 imagick: http://php.net/manual/en/book.imagick.php

# convert命令

转换图片格式
```
# convert src.jpg dst.png
```
指定处理后的图片品质
```
# convert -quality 80 src.jpg dst.png
```
缩放
```
-resize widthxheight{%} {@} {!} {<} {>} {^}
```
1、按照宽高中大的一边放大或缩小，另一边按比例处理。
```
# convert -resize 600×600 src.jpg dst.jpg
```
这里的600x600表示600px，同时也可以用百分比来表示，下面命令表示将宽缩到50%，高缩到80%
```
# convert -resize 50%x80% src.jpg dst.jpg
```
2、将图片调整到固定尺寸，可以在宽高后面加上一个感叹号!
```
# convert -resize 600×600! src.jpg dst.jpg
```
3、固定宽度或高度缩放，另一边按比例处理
```
# convert -resize 400 src.jpg dst.jpg

# convert -resize x400 src.jpg dst.jpg
```
4、使用 @ 来指定图片的像素个数。
```
# convert -resize "10000@" src.jpg dst.jpg
```
原图大小910x1365，执行后缩为81x122

5、图片大于或小于某个尺寸时缩放
```
# convert -resize "300x200>" src.jpg dst.jpg
```
300x200>表示当图片大于300x200时进行缩小

300x200<表示当图片小于300x200时进行放大 