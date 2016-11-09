#!/bin/bash
if [[ $# -eq 2  && -f $1 ]]
then

        file="$1"
        path=${file%%.*} ##upload/product/20141112/9590565ac926283fc3381d6e235f8892
        ext=${file#*.} ##png
        #echo $file
        #echo $path
        #echo $ext

        if [ $2 -eq 1 ]
        then
                convert -quality 80 -resize 120x120! ${file} ${path}_120x120.${ext} 2>&1 >/dev/null
                convert -quality 80 -resize 240x240! ${file} ${path}_240x240.${ext} 2>&1 >/dev/null
                convert -quality 80 -resize 480x480! ${file} ${path}_480x480.${ext} 2>&1 >/dev/null
                convert -quality 80 -resize 960x960! ${file} ${path}_960x960.${ext} 2>&1 >/dev/null
        elif [ $2 -eq 2 ]
        then
                convert -quality 80 -resize 120x ${file} ${path}_120.${ext} 2>&1 >/dev/null
                convert -quality 80 -resize 240x ${file} ${path}_240.${ext} 2>&1 >/dev/null
                convert -quality 80 -resize 480x ${file} ${path}_480.${ext} 2>&1 >/dev/null
                convert -quality 80 -resize 960x ${file} ${path}_960.${ext} 2>&1 >/dev/null
        elif [ $2 -eq 3 ]
        then

                convert -quality 80 -resize 640x ${file} ${path}_640.${ext} 2>&1 >/dev/null
        else
                echo "Parameter error"
                exit 1
        fi
else
        echo "system error"
        exit 2
fi
