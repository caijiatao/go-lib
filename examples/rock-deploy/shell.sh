cat /root/.kube/config

kubectl create namespace sjs

# 搭建
kubectl apply -f pg-pv-pvc.yaml -n sjs

kubectl apply -f postgres.yaml -n sjs

# 执行创建数据库
# 执行sql/rock.sql

kubectl apply -f rock-backend-conf.yaml -n sjs

kubectl apply -f rock-backend.yaml -n sjs

# kubectl delete -f rock-backend.yaml -n sjs


