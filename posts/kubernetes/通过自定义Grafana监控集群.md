```
{
    "url": "grafana-k8s",
    "time": "2021/04/25 21:18",
    "tag": "Kubernetes,容器化"
}
```

阿里云Ack容器集成了Prometheus监控，但监控数据的展示需要进入到阿里云后台，涉及到多账号配置等。我们可以将监控数据配置到自己的Grafana上，只需要添加Dashboard配置对应的数据源即可。

首先，需要先安装Grafana：

1、创建PV

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: grafana-pv
  labels:
    project: grafana
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: grafana-pv
    volumeAttributes:
      server: "*.cn-hangzhou.nas.aliyuncs.com"
      path: "/grafana"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: grafana-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 5Gi
  selector:
    matchLabels:
      project: grafana

```

2、创建Service

```
apiVersion: v1
kind: Service
metadata:
  name: grafana
spec:
  selector:
    app: grafana
  ports:
  - port: 80
    targetPort: 3000
    protocol: TCP

---

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: grafana
spec:
  rules:
  - host: grafana.demo.com
    http:
      paths:
      - path: /
        backend:
          serviceName: grafana
          servicePort: 80

```

3、创建Deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
  labels:
    app: grafana
spec:
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
      - name: grafana
        image: grafana/grafana:7.4.5
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 3000
          name: grafana
        env:
        - name: GF_SECURITY_ADMIN_USER
          value: admin
        - name: GF_SECURITY_ADMIN_PASSWORD
          value: admin
        readinessProbe:
          failureThreshold: 10
          httpGet:
            path: /api/health
            port: 3000
            scheme: HTTP
          initialDelaySeconds: 60
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 30
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /api/health
            port: 3000
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources:
          requests:
            cpu: 100m
            memory: 256Mi
          limits:
            cpu: 1000m
            memory: 1000Mi
        volumeMounts:
        - mountPath: /var/lib/grafana
          subPath: grafana
          name: storage
      securityContext:
        fsGroup: 472
        runAsUser: 472
      volumes:
      - name: storage
        persistentVolumeClaim:
          claimName: grafana-pvc
```

Grafana创建后接下来要做的是将数据源配置到自建的Grafana中。

1、集群通过Prometheus存储监控数据，需要先找到Prometheus的API接口地址：

```
容器服务 - Kubernetes -> 点击对应集群 -> Prometheus监控 -> 跳转到应用实施监控服务 ARMS -> 设置 -> Agent设置
```

2、在自建的Grafana中配置数据源

```
Configuration -> Data Sources
```

3、在自建的Grafana中创建Folder，用来展示该集群的所有监控图表

```
Dashboards -> Manage -> New Folder
```

4、Dashboard导入

```
1、在阿里云的Grafana界面打开想要导入的图表，左上有个 Share Dashboard的，进入Export Tab页面
2、拷贝JSON并修改数据源为前面配置的数据源名称
3、回到自建的Grafana，进入前一步创建的Folder，点击Import完成导入。
```

正常应该就可以看到监控数据了，根据需要导入想要展示的图表即可。