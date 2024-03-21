SERVICES := chat chat_api_server notification
DB := mysql
ETCD := etcd

.PHONY: all $(SERVICES) $(DB) $(ETCD)

all: $(DB) $(ETCD) $(SERVICES)

$(DB):
    @echo "Starting $(DB) service..."
    # 部署一个DB，可以用docker启动
    docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7.22

$(ETCD):
    @echo "Starting $(ETCD) service..."
    # 部署一个ETCD，可以用docker启动
    docker run -d -p 2379:2379 --name etcd quay.io/coreos/etcd:v3.3.9

$(SERVICES): $(DB) $(ETCD)
    @echo "Starting $@ service..."
    cd $@ && go run main.go &


.PHONY: clean

clean:
    @echo "Stopping all services..."