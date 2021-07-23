```
{
    "url": "mongo-in-k8s",
    "time": "2021/05/30 21:10",
    "tag": "MongoDB,Kubernetes,容器化"
}
```

尝试在本地运行`Mongo`:

```
$ docker pull mongo:3.4.24

$ docker run --name mongo -d mongo:3.4.24
$ docker ps
CONTAINER ID   IMAGE          COMMAND                  CREATED         STATUS         PORTS       NAMES
50c4b61834a0   mongo:3.4.24   "docker-entrypoint.s…"   3 seconds ago   Up 2 seconds   27017/tcp   mongo
```

系统为：`Ubuntu 16.04.6 LTS`，通过命令行也可以连接：

```
$ mongo --host 127.0.0.1 --port 27017
MongoDB shell version v3.4.24
connecting to: mongodb://127.0.0.1:27017/
MongoDB server version: 3.4.24
Welcome to the MongoDB shell.
```

同`Mysql`类似，也可以通过环境变量来设置`root`账号：

`MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD`

> These variables, used in conjunction, create a new user and set that user's password. This user is created in the `admin` authentication database and given the role of `root`), which is a "superuser" role).

尝试运行之后就可以在容器里部署了，这里部署为单节点Mongo服务。

**第一步，**创建pv/pvc

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mongo-test-pv
  labels:
    project: mongo-test-pv
spec:
  capacity:
    storage: 50Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: mongo-test-pv
    volumeAttributes:
      server: "nasid.cn-hangzhou.nas.aliyuncs.com"
      path: "/mongo-test"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: mongo-test-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi
  selector:
    matchLabels:
      project: mongo-test-pv
```

**第二步**，创建Service

```
apiVersion: v1
kind: Service
metadata:
  name: mongo-test
  namespace: default
  labels:
    project: mongo-test
spec:
  selector:
    project: mongo-test
  ports:
  - name: mongo
    port: 27017
    protocol: TCP
  clusterIP: None
```

**第三步**，创建ConfigMap，如果开启授权的话可以在运行安装成功之后再开启。

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: mongo-test-config
  namespace: default
data:
  mongodb.conf: |
    storage:
        dbPath: "/data/db"
        journal:
          enabled: true
    systemLog:
        destination: file
        logAppend: true
        path: "/data/mongodb.log"
    net:
        port: 27017
        bindIp: 0.0.0.0
    processManagement:
        fork: false
        pidFilePath: "/data/mongod.pid"
    #security:
    #    authorization: enabled
    operationProfiling:
        slowOpThresholdMs: 3000
        mode: slowOp

```

**第四步**，创建StatefulSet，`sync-host-time`从宿主机同步了时间。

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mongo-test
  namespace: default
  labels:
    project: mongo-test
spec:
  serviceName: mongo-test
  replicas: 1
  selector:
    matchLabels:
      project: mongo-test
  template:
    metadata:
      labels:
        project: mongo-test
    spec:
      restartPolicy: Always
      containers:
      - name: mongo
        image: mongo:3.4.24
        imagePullPolicy: IfNotPresent
        command: ["/bin/sh", "-c", "/usr/bin/mongod -f /data/conf/mongodb.conf"]
        ports:
        - containerPort: 27017
        volumeMounts:
        - name: mongo-test-pv
          mountPath: /data/db/
        - name: mongo-test-config
          mountPath: /data/conf/
        - name: sync-host-time
          mountPath: /etc/localtime
          readOnly: true
      volumes:
      - name: mongo-test-config
        configMap:
          name: mongo-test-config
      - name: mongo-test-pv
        persistentVolumeClaim:
          claimName: mongo-test-pvc
      - name: sync-host-time
        hostPath:
          path: /etc/localtime
```

运行上面编写的Yaml文件即可，正常一个单节点的Mongo就应该创建好了。
