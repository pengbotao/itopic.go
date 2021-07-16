```
{
    "url": "gitlab-in-k8s",
    "time": "2021/07/16 23:19",
    "tag": "Gitlab,Kubernetes,容器化"
}
```



```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: gitlab-pv
  labels:
    project: gitlab-pv
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  csi:
    driver: nasplugin.csi.alibabacloud.com
    volumeHandle: gitlab-pv
    volumeAttributes:
      server: "*.cn-hangzhou.nas.aliyuncs.com"
      path: "/gitlab"
  mountOptions:
  - nolock,tcp,noresvport
  - vers=3

---

kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: gitlab-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 20Gi
  selector:
    matchLabels:
      project: gitlab-pv
```





```
apiVersion: v1
kind: Service
metadata:
  name: gitlab
spec:
  selector:
    project: gitlab
  ports:
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  - name: ssh
    port: 22
    targetPort: 22
    protocol: TCP
  clusterIP: None
```





```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: gitlab
  labels:
    project: gitlab
spec:
  serviceName: gitlab
  replicas: 1
  selector:
    matchLabels:
      project: gitlab
  template:
    metadata:
      labels:
        project: gitlab
    spec:
      containers:
      - name: gitlab
        image: gitlab/gitlab-ce:9.1.3-ce.0
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          containerPort: 80
        - name: ssh
          containerPort: 22
        env:
        - name: TZ
          value: Asia/Shanghai
        - name: GITLAB_ROOT_PASSWORD
          value: password
        - name: GITLAB_ROOT_EMAIL
          value: mail@gitlab.com
        volumeMounts:
        - name: gitlab-pv
          mountPath: /var/opt/gitlab
          subPath: gitlab
        - name: gitlab-pv
          mountPath: /etc/gitlab
          subPath: etc
      volumes:
      - name: gitlab-pv
        persistentVolumeClaim:
          claimName: gitlab-pvc
```



```
$ cat /etc/gitlab/gitlab.rb
prometheus['enable'] = false
prometheus_monitoring['enable'] = false
external_url 'http://gitlab.x.com'
unicorn['worker_processes'] = 2
```





```
apiVersion: v1
kind: Service
metadata:
  name: gitlab-slb
  annotations:
    service.beta.kubernetes.io/alibaba-cloud-loadbalancer-id: slbid
    service.beta.kubernetes.io/alicloud-loadbalancer-force-override-listeners: 'true'
spec:
  selector:
    app: gitlab
  ports:
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  - name: ssh
    port: 22
    targetPort: 22
    protocol: TCP
  type: LoadBalancer
```

