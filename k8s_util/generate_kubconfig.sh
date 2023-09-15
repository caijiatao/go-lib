server=https://192.168.13.213:6443

ca=$(kubectl get secret $(kubectl get serviceaccount airec-server-client -o jsonpath='{.secrets[0].name}') -o jsonpath='{.data.ca\.crt}')
token=$(kubectl get secret $(kubectl get serviceaccount airec-server-client -o jsonpath='{.secrets[0].name}') -o jsonpath='{.data.token}' | base64 --decode)

echo "
apiVersion: v1
kind: Config
clusters:
- name: default-cluster
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: default-context
  context:
    cluster: default-cluster
    namespace: default
    user: default-user
current-context: default-context
users:
- name: default-user
  user:
    token: ${token}
" > airec-server-client.kubeconfig


redis-cli --cluster create 10.244.1.232:6379 10.244.1.233:6379 10.244.1.234:6379
