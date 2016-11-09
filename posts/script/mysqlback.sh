#!/bin/bash
#by hz

suffix=`date "+%Y%m%d"`.sql
dir="/data/mysql_backup"
pro="/usr/local/server/mysql/bin/mysqldump"
myDBuser="root"
myDBpassword=''
db_list="/root/db.txt"

if [[ ! -d $dir ]]
then
        mkdir $dir
fi

function do_bk()
{
cd $dir
cat ${db_list}|grep -v "#"|while read db
do
        $pro -u${myDBuser} -p${myDBpassword} ${db}|gzip >${db}-${suffix}.gz
done
}

function do_cb()
{
        find ${dir}  -mtime +180 -exec rm -rf {} \;
}
##backup
do_bk
do_cb
