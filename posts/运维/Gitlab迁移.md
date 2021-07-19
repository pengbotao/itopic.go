```
{
    "url": "gitlab-migration",
    "time": "2019/06/01 11:52",
    "tag": "Gitlab,运维"
}
```

目标：Gitlab服务器迁移，需要在新服务器安装及数据同步。

Step1：原Gitlab数据备份

```
$ gitlab-rake gitlab:backup:create
Dumping database ...
Dumping PostgreSQL database gitlabhq_production ... [DONE]
done
Dumping repositories ...
......
done
Dumping uploads ...
done
Dumping builds ...
done
Dumping artifacts ...
done
Dumping pages ...
done
Dumping lfs objects ...
done
Dumping container registry images ...
[DISABLED]
Creating backup archive: 1622519542_2021_06_01_gitlab_backup.tar ... done
Uploading backup archive to remote storage  ... skipped
Deleting tmp directories ... done
done
done
done
done
done
done
done
Deleting old backups ... skipping

$ ls -lh /var/opt/gitlab/backups/
总用量 639M
-rw------- 1 git git 269M 6月   1 11:52 1622519542_2021_06_01_gitlab_backup.tar
```

Step2: 安装新Gitlab

```
$ rpm -ivh gitlab-ce-9.1.3-ce.0.el6.x86_64.rpm
Preparing...                ########################################### [100%]
   1:gitlab-ce              ########################################### [100%]


       *.                  *.
      ***                 ***
     *****               *****
    .******             *******
    ********            ********
   ,,,,,,,,,***********,,,,,,,,,
  ,,,,,,,,,,,*********,,,,,,,,,,,
  .,,,,,,,,,,,*******,,,,,,,,,,,,
      ,,,,,,,,,*****,,,,,,,,,.
         ,,,,,,,****,,,,,,
            .,,,***,,,,
                ,*,.

     _______ __  __          __
    / ____(_) /_/ /   ____ _/ /_
   / / __/ / __/ /   / __ `/ __ \
  / /_/ / / /_/ /___/ /_/ / /_/ /
  \____/_/\__/_____/\__,_/_.___/


gitlab: Thank you for installing GitLab!
gitlab: To configure and start GitLab, RUN THE FOLLOWING COMMAND:

sudo gitlab-ctl reconfigure

gitlab: GitLab should be reachable at http://demo
gitlab: Otherwise configure GitLab for your system by editing /etc/gitlab/gitlab.rb file
gitlab: And running reconfigure again.
gitlab:
gitlab: For a comprehensive list of configuration options please see the Omnibus GitLab readme
gitlab: https://gitlab.com/gitlab-org/omnibus-gitlab/blob/master/README.md
gitlab:
It looks like GitLab has not been configured yet; skipping the upgrade script.
```

安装之后调整配置，配置文件：`/etc/gitlab/gitlab.rb`

```
gitlab_rails['backup_path'] = "/var/opt/gitlab/backups"

external_url 'http://gitlab.test.com'

git_data_dirs({ "default" => { "path" => "/data/git-data", 'gitaly_address' => 'unix:/var/opt/gitlab/gitaly/gitaly.socket' } })
```

然后执行：

```
$ gitlab-ctl reconfigure
```

然后会自动创建用户以及git-data目录。查看状态，可以看到已经都启动了

```
$ gitlab-ctl status
run: gitaly: (pid 13300) 3618777s; run: log: (pid 8252) 3619408s
run: gitlab-monitor: (pid 13310) 3618777s; run: log: (pid 8616) 3619388s
run: gitlab-workhorse: (pid 13322) 3618776s; run: log: (pid 8286) 3619406s
run: logrotate: (pid 16000) 764s; run: log: (pid 8347) 3619402s
run: nginx: (pid 13347) 3618776s; run: log: (pid 8320) 3619404s
run: node-exporter: (pid 13360) 3618775s; run: log: (pid 8486) 3619394s
run: postgres-exporter: (pid 13372) 3618775s; run: log: (pid 8583) 3619390s
run: postgresql: (pid 13383) 3618774s; run: log: (pid 7979) 3619438s
run: prometheus: (pid 13398) 3618774s; run: log: (pid 8453) 3619396s
run: redis: (pid 13414) 3618774s; run: log: (pid 7907) 3619444s
run: redis-exporter: (pid 13423) 3618773s; run: log: (pid 8530) 3619392s
run: sidekiq: (pid 13434) 3618773s; run: log: (pid 8237) 3619410s
run: unicorn: (pid 13437) 3618772s; run: log: (pid 8196) 3619412s
```

Step3: 数据迁移。同步备份文件到backups目录，按下面操作输入2个yes后即可。

```
$ gitlab-ctl stop unicorn 
$ gitlab-ctl stop sidekiq
$ gitlab-rake gitlab:backup:restore BACKUP=1622519542_2021_06_01
$ gitlab-ctl restart
```

