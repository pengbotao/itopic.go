```
{
    "url": "jumpserver-in-k8s",
    "time": "2021/08/10 23:14",
    "tag": "JumpServer,Kubernetes,容器化",
    "public": "no"
}
```

和SSHD的功能类似，JumpServer的作用则是在本地可以实现对远程服务的访问，同时可以方便的进行授权、审计等。

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jumpserver
  labels:
    app: jumpserver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jumpserver
  template:
    metadata:
      labels:
        app: jumpserver
    spec:
      containers:
      - name: jumpserver
        image: jumpserver/jms_all:v2.11.4
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
          name: jumpserver
        - containerPort: 2222
          name: ssh
        env:
        - name: SECRET_KEY
          # cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 50
          value: iJ2YaSaNLCALNOVSqsw7sgwk3cX5gK6nCcey57UZiujF20I32n
        - name: BOOTSTRAP_TOKEN
          # cat /dev/urandom | tr -dc A-Za-z0-9 | head -c 16
          value: S7sLMAH9J0mTaqSZ
        - name: DB_HOST
          value: mysql-js.default.svc.cluster.local
        - name: DB_PORT
          value: "3306"
        - name: DB_USER
          value: jumpserver
        - name: DB_PASSWORD
          value: "password"
        - name: DB_NAME
          value: jumpserver
        - name: REDIS_HOST
          value: redis-js.default.svc.cluster.local
        - name: REDIS_PORT
          value: "6379"
        - name: REDIS_PASSWORD
          value: ""
        volumeMounts:
        - name: jumpserver-pv
          mountPath: /opt/jumpserver/data/media
      volumes:
      - name: jumpserver-pv
        persistentVolumeClaim:
          claimName: jumpserver-pvc
```



```
apiVersion: v1
kind: Service
metadata:
  name: jumpserver
spec:
  selector:
    app: jumpserver
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
```



```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: jumpserver-pv
  labels:
    app: jumpserver
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: jumpserver-pv
    volumeAttributes:
      server: "nasid.cn-hangzhou.nas.aliyuncs.com"
      path: "/jumpserver"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: jumpserver-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 5Gi
  selector:
    matchLabels:
      app: jumpserver

```

