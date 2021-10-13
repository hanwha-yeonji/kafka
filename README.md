# kafka
docker file for kafka and sample source code

## Run Dockerfile
```shell
docker run --rm -it --name kafka \
 -e CLUSTER_ID=RuTPVymzR16F2f40eUFmgw \
 -e BROKER_ID=1 \
 -p 19092:9092 \
 -p 19093:9093 \
 docker.pkg.github.com/hanwha-yeonji/kafka/kafka:3.0.0 
```

## Run Docker-compose
```shell
docker-compose up -d
```

## Reference
- https://github.com/wurstmeister/kafka-docker/blob/master
- https://github.com/gunnarmorling/docker-images/tree/DBZ-3903/kafka/1.7