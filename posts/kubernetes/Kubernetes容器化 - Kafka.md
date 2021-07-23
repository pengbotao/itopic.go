```
{
    "url": "kafka-in-k8s",
    "time": "2021/06/10 19:46",
    "tag": "Kafka,Kubernetes,容器化"
}
```

除最新的2.8不需要ZooKeeper，这之前的版本都需要用到ZooKeeper，可以参考前一篇[Kubernetes容器化 - Zookeeper](https://itopic.org/zookeeper-in-k8s.html)安装ZooKeeper。ZooKeeper与Kafka版本之间貌似有一些对应关系，未找到原出处。在容器中安装Kafka首先要看有没有合适的镜像，主要看到的有2个镜像：

- `bitnami/kafka`
- [wurstmeister/kafka](https://github.com/wurstmeister/kafka-docker)

都不是官方镜像，两个都做过尝试，最终选择的是第二个镜像。

- 项目根目录在`/opt/kafka`
- 数据存储在`/kafka`，可以挂载容器外存储

**第一步，创建pv**，这里和前面ZooKeeper的操作一样，找个合适的存储创建pv即可。

**第二步，创建Service**

```
apiVersion: v1
kind: Service
metadata:
  name: kafka
  namespace: default
  labels:
    project: kafka
spec:
  selector:
    project: kafka
  ports:
  - name: kafka
    port: 9092
    targetPort: 9092
    protocol: TCP
  clusterIP: None
```

**第三步，创建StatefulSet**，这里测试的版本不新，但这个镜像最近有更新，理论上和其他版本安装差不多，镜像使用说明：

>1、通过StatefulSet来部署集群，通过BROKER_ID_COMMAND环境变量可以实现broker.id的动态设置
>
>2、类似KAFKA_LOG_RETENTION_HOURS和KAFKA_DELETE_TOPIC_ENABLE，配置文件中的变量可以以类似方式注入到配置文件
>
>3、目前集群内部访问正常，外部访问Broker失败，应该是Pod内部IP无法直接与外部通信。可以排查KAFKA_ADVERTISED_LISTENERS的设置。

```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: kafka
  namespace: default
  labels:
    project: kafka
spec:
  serviceName: kafka
  replicas: 3
  selector:
    matchLabels:
      project: kafka
  template:
    metadata:
      labels:
        project: kafka
    spec:
      restartPolicy: Always
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 80
            podAffinityTerm:
              topologyKey: kubernetes.io/hostname
              labelSelector:
                matchExpressions:
                - key: project
                  operator: In
                  values: 
                  - kafka
      containers:
      - name: kafka
        image: "wurstmeister/kafka:2.11-0.11.0.3"
        imagePullPolicy: "IfNotPresent"
        ports:
        - name: kafka
          containerPort: 9092
        env:
        - name: BROKER_ID_COMMAND
          value: "hostname | awk -F '-' '{print $NF}'"
        - name: KAFKA_PORT
          value: "9092"
        - name: KAFKA_LISTENERS
          value: "PLAINTEXT://:9092"
        - name: KAFKA_ADVERTISED_PORT
          value: "9092"
        - name: KAFKA_ADVERTISED_LISTENERS
          value: "PLAINTEXT://:9092"
        - name: KAFKA_ZOOKEEPER_CONNECT
          value: "zookeeper.default.svc.cluster.local:2181"
        - name: KAFKA_LOG_RETENTION_HOURS
          value: "168"
        - name: KAFKA_DELETE_TOPIC_ENABLE
          value: "false"
        - name: KAFKA_DEFAULT_REPLICATION_FACTOR
          value: "2"
        volumeMounts:
        - name: kafka-pvc
          mountPath: /kafka
  volumeClaimTemplates:
  - metadata:
      name: kafka-pvc
    spec:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 50Gi
        selector:
          matchLabels:
            project: kafka-pv

```



---

- [1] [K8S 搭建 Kafka:2.13-2.6.0 和 Zookeeper:3.6.2 集群](https://www.debugger.wiki/article/html/1604647080423925)