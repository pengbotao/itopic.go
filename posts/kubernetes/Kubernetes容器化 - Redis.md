```
{
    "url": "redis-in-k8s",
    "time": "2021/06/05 13:14",
    "tag": "Redis,Kubernetes,容器化"
}
```

这里配置为单节点配置。

第一步，创建ConfigMap

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-demo-config
  namespace: default
data:
  redis.conf: |
    daemonize no
    port 6379
    tcp-backlog 511
    bind 0.0.0.0
    protected-mode no
    timeout 0
    tcp-keepalive 300
    loglevel notice
    maxmemory 4G
    maxmemory-policy volatile-lru
    maxmemory-samples 5
    maxmemory-eviction-tenacity 10
    requirepass yourpwd
    databases 16
    save 900 1
    save 300 10
    save 60 10000
    stop-writes-on-bgsave-error no
    rdbcompression yes
    rdbchecksum yes
    dbfilename dump.rdb
    dir /data
    appendonly no
    appendfilename "appendonly.aof"
    appendfsync everysec
```

第二步，创建Service

```
apiVersion: v1
kind: Service
metadata:
  name: redis-demo
  namespace: default
  labels:
    project: redis-demo
spec:
  selector:
    project: redis-demo
  ports:
  - name: redis
    port: 6379
    protocol: TCP
  clusterIP: None
```

第三步，创建持久化存储

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: redis-demo-pv
  labels:
    project: redis-demo-pv
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: redis-demo-pv
    volumeAttributes:
      server: "nasid.cn-hangzhou.nas.aliyuncs.com"
      path: "/redis-demo"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: redis-demo-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  selector:
    matchLabels:
      project: redis-demo-pv

```

第四步，编写StatefulSet

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis-demo
  namespace: default
  labels:
    project: redis-demo
spec:
  serviceName: redis-demo
  replicas: 1
  selector:
    matchLabels:
      project: redis-demo
  template:
    metadata:
      labels:
        project: redis-demo
    spec:
      containers:
      - name: redis
        image: redis:6.2.4
        imagePullPolicy: IfNotPresent
        command: ["/bin/sh", "-c", "/usr/local/bin/redis-server /usr/local/etc/redis/redis.conf"]
        ports:
        - containerPort: 6379
          name: redis
          protocol: TCP
        readinessProbe:
          tcpSocket:
            port: 6379
          initialDelaySeconds: 5
          timeoutSeconds: 15
          periodSeconds: 5
        livenessProbe:
          tcpSocket: 
            port: 6379
          initialDelaySeconds: 30
          timeoutSeconds: 15
          periodSeconds: 15
        volumeMounts:
        - name: redis-demo-config
          mountPath: /usr/local/etc/redis/redis.conf
          subPath: redis.conf
        - name: redis-demo-pv
          mountPath: /data
      volumes:
      - name: redis-demo-config
        configMap:
          name: redis-demo-config
      - name: redis-demo-pv
        persistentVolumeClaim:
          claimName: redis-demo-pvc
```

