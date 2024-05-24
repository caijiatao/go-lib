docker build -f ./docker-example/Dockerfile -t docker-example:latest .
docker build -t docker-example:latest .
docker run -d --name test docker-example:latest
docker exec -it test-1 /bin/bash


docker run -it docker-example:latest /bin/bash
