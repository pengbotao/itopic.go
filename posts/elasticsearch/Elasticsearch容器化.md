```
{
    "url": "elasticsearch-in-k8s",
    "time": "2021/05/22 19:46",
    "tag": "Elasticsearch,容器化",
    "public": "no"
}
```

# 创建pv/pvc

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: es-demo-pv
  labels:
    project: es-demo-pv
spec:
  capacity:
    storage: 2Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: es-demo-pv
    volumeAttributes:
      server: "*.cn-hangzhou.nas.aliyuncs.com"
      path: "/data"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: es-demo-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 2Gi
  selector:
    matchLabels:
      project: es-demo-pv
```

# 创建Service

```
apiVersion: v1
kind: Service
metadata:
  name: es-demo-svc
spec:
  selector:
    project: es-demo
  ports:
  - name: http
    port: 9200
    targetPort: 9200
    protocol: TCP
  - name: transport
    port: 9300
    targetPort: 9300
    protocol: TCP
  clusterIP: None
```

# 创建StatefulSet

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: es-demo
  namespace: default
  labels:
    project: es-demo
spec:
  serviceName: es-demo-svc
  replicas: 1
  selector:
    matchLabels:
      project: es-demo
  template:
    metadata:
      labels:
        project: es-demo
    spec:
      restartPolicy: Always
      containers:
      - name: elasticsearch
        image: elasticsearch:6.8.15
        imagePullPolicy: Always
        ports:
        - containerPort: 9300
          name: transport
          protocol: TCP
        - containerPort: 9200
          name: http
          protocol: TCP
        env:
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: "ES_JAVA_OPTS"
          value: "-Xms512m -Xmx512m"
        - name: discovery.type
          value: single-node
        - name: cluster.name
          value: "es-demo-${NAMESPACE}"
        - name: node.name
          value: "${POD_NAME}-${NAMESPACE}"
        - name: network.host
          value: "${POD_IP}"
        readinessProbe:
          tcpSocket:
            port: 9200
          initialDelaySeconds: 5
          timeoutSeconds: 15
          periodSeconds: 5
        livenessProbe:
          tcpSocket: 
            port: 9200
          initialDelaySeconds: 30
          timeoutSeconds: 15
          periodSeconds: 15
        volumeMounts:
        - name: es-pv
          mountPath: /usr/share/elasticsearch/data
          subPath: data
      volumes:
      - name: es-demo-pv
        persistentVolumeClaim:
          claimName: es-demo-pvc
```

# 创建kibana

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: es-demo-kibana
  namespace: default
  labels:
    project: es-demo-kibana
spec:
  replicas: 1
  selector:
    matchLabels:
      project: es-demo-kibana
  template:
    metadata:
      labels:
        project: es-demo-kibana
    spec:
      restartPolicy: Always
      containers:
      - name: kibana
        image: kibana:6.8.15
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 5601
        env:
        - name: ELASTICSEARCH_URL
          value: "http://es-demo.default.svc.cluster.local:9200"
        volumeMounts:
        - name: es-demo-pv
          mountPath: /usr/share/kibana/data
          subPath: kibana
      volumes:
      - name: es-demo-pv
        persistentVolumeClaim:
          claimName: es-demo-pvc
---

apiVersion: v1
kind: Service
metadata:
  name: es-demo-kibana
spec:
  selector:
    project: es-demo-kibana
  ports:
  - name: http
    port: 80
    targetPort: 5601
    protocol: TCP

---

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: es-demo-kibana
  namespace: default
spec:
  rules:
  - host: es.demo.com
    http:
      paths:
      - path: /
        backend:
          serviceName: es-demo-kibana
          servicePort: 80
```

# 集群模式

```
        env:
        - name: NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: "ES_JAVA_OPTS"
          value: "-Xms2g -Xmx2g"
        - name: cluster.name
          value: "es-demo-${NAMESPACE}"
        - name: node.name
          value: "${POD_NAME}-${NAMESPACE}"
        - name: network.host
          value: "${POD_IP}"
        - name: discovery.zen.ping.unicast.hosts
          value: "es-demo-svc"
        - name: discovery.zen.minimum_master_nodes
          value: "2"
```



# 索引迁移

Github: https://github.com/medcl/esm

`esm`支持跨版本进行数据迁移，对数据比较规整的情况迁移也挺快。这里有个40多G的索引跨版本升级，一直迁移不成功，机器内存升级到256G跑到一半就蹦了，当然也有可能跟数据有关系，涉及到数据迁移的不妨一试。

```
./esm -s http://192.168.0.100:9200 -d http://192.168.0.101:9200 -x orders -y orders -w 4
[05-24 17:47:04] [INF] [main.go:474,main] start data migration..
Scroll 1428110 / 1428110 [===============================================================] 100.00% 2m5s
Bulk 1428005 / 1428110 [=================================================================]  99.99% 2m5s
[05-24 17:49:09] [INF] [main.go:505,main] data migration finished.
```

