FROM azul/zulu-openjdk-alpine:15.0.2-15.29.15

ARG KAFKA_VERSION=3.0.0
ARG SCALA_VERSION=2.13

ENV KAFKA_VERSION=$KAFKA_VERSION \
    SCALA_VERSION=$SCALA_VERSION \
    KAFKA_HOME=/kafka \
    KAFKA_URL_PATH=kafka/$KAFKA_VERSION/ \
    KAFKA_FILE_NAME=kafka_$SCALA_VERSION-$KAFKA_VERSION.tgz
ENV KAFKA_URL_PATH=kafka/$KAFKA_VERSION/$KAFKA_FILE_NAME

ENV PATH=$PATH:$KAFKA_HOME/bin


# create home directory
RUN mkdir $KAFKA_HOME \
 && mkdir $KAFKA_HOME/logs \
 && mkdir $KAFKA_HOME/data

# install utils
RUN apk add --no-cache bash curl jq docker

# download kafka
RUN curl -fSL -o /tmp/$KAFKA_FILE_NAME $(curl --stderr /dev/null https://www.apache.org/dyn/closer.cgi\?as_json\=1 | sed -rn 's/.*"preferred":.*"(.*)"/\1/p')$KAFKA_URL_PATH || curl -fSL -o /tmp/kafka.tgz https://archive.apache.org/dist/$KAFKA_URL_PATH

# install kafka
RUN tar -xvzf /tmp/$KAFKA_FILE_NAME -C $KAFKA_HOME --strip-components 1 \
 && rm -f /tmp/$KAFKA_FILE_NAME

COPY ./log4j.properties $KAFKA_HOME/config/log4j.properties
RUN mkdir $KAFKA_HOME/config.orig \
 && mv $KAFKA_HOME/config/* $KAFKA_HOME/config.orig

# remove unnecessary files
RUN rm -rf $KAFKA_HOME/libs/*-{sources,javadoc,scaladoc}.jar* \
 && rm -rf $KAFKA_HOME/site-docs

# expose the ports and set up volumes
EXPOSE 9092
VOLUME ["/kafka/data","/kafka/logs","/kafka/config"]

COPY ./entry.sh /usr/bin
WORKDIR $KAFKA_HOME

CMD ["entry.sh"]
