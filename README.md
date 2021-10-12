# kafka
docker file for kafka and sample source code

## Run Dockerfile
```shell
docker run --rm -it --net=host --name kafka \
 -e CLUSTER_ID=RuTPVymzR16F2f40eUFmgw \
 -e HOST_NAME=localhost \
 docker.pkg.github.com/hanwha-yeonji/kafka/kafka:3.0.0 
```

## Run Docker-compose
```shell
docker-compose up -d
```

## Reference
- https://github.com/wurstmeister/kafka-docker/blob/master
- https://github.com/gunnarmorling/docker-images/tree/DBZ-3903/kafka/1.7