```shell
redis-cli --cluster create redis-cluster-access-0.exploit-ghrec-cas.svc.cluster.local:6379 redis-cluster-access-1.exploit-ghrec-cas.svc.cluster.local:6379 redis-cluster-access-2.exploit-ghrec-cas.svc.cluster.local:6379 --cluster-replicas 0

redis-cli --cluster create redis-cluster-headless-0:6379 redis-cluster-headless-1:6379 redis-cluster-headless-2:6379 --cluster-replicas 0


redis-cli --cluster create 	10.244.5.168:6379 10.244.3.223:6379 10.244.3.224:6379
```