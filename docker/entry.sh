#!/bin/bash

set -e

if [[ -z "$NODE_ID" ]]; then
    if [[ -z "$BROKER_ID" ]]; then
        NODE_ID=1
        echo "WARNING: Using default NODE_ID=1, which is valid only for non-clustered installations."
    else
        NODE_ID="$BROKER_ID"
        echo "WARNING: Using NODE_ID=$BROKER_ID, as specified via BROKER_ID variable. Please update your configuration to use the NODE_ID variable instead."
    fi
fi

# ZooKeeper mode
if [[ -z "$CLUSTER_ID" ]]; then
  CONFIG_FILE=config/server.properties
  echo "Starting in ZooKeeper mode using NODE_ID=$NODE_ID."

  if [[ -z "$ZOOKEEPER_CONNECT" ]]; then
      # Look for any environment variables set by Docker container linking. For example, if the container
      # running Zookeeper were named 'zoo' in this container, then Docker should have created several envs,
      # such as 'ZOO_PORT_2181_TCP'. If so, then use that to automatically set the 'zookeeper.connect' property.
      export ZOOKEEPER_CONNECT=$(env | grep .*PORT_2181_TCP= | sed -e 's|.*tcp://||' | uniq | paste -sd ,)
  fi
  if [[ "x$ZOOKEEPER_CONNECT" = "x" ]]; then
      export ZOOKEEPER_CONNECT=0.0.0.0:2181
  fi
  echo "Using ZOOKEEPER_CONNECT=$ZOOKEEPER_CONNECT"
# KRaft mode
else

    if [[ -z "$NODE_ROLE" ]]; then
        NODE_ROLE='combined';
    fi

    case "$NODE_ROLE" in
     'combined' ) CONFIG_FILE=config/kraft/server.properties;;
     'broker' ) CONFIG_FILE=config/kraft/broker.properties;;
     'controller' ) CONFIG_FILE=config/kraft/controller.properties;;
     *) CONFIG_FILE=config/kraft/server.properties;;
    esac

    echo "Starting in KRaft mode (EXPERIMENTAL), using CLUSTER_ID=$CLUSTER_ID, NODE_ID=$NODE_ID and NODE_ROLE=$NODE_ROLE."
fi

echo "Using configuration $CONFIG_FILE."

if [[ -n "$HEAP_OPTS" ]]; then
    sed -r -i "s/^(export KAFKA_HEAP_OPTS)=\"(.*)\"/\1=\"${HEAP_OPTS}\"/g" $KAFKA_HOME/bin/kafka-server-start.sh
    unset HEAP_OPTS
fi

export KAFKA_ZOOKEEPER_CONNECT=$ZOOKEEPER_CONNECT
export KAFKA_NODE_ID=$NODE_ID
export KAFKA_LOG_DIRS="$KAFKA_HOME/data/$KAFKA_BROKER_ID"
mkdir -p $KAFKA_LOG_DIRS
unset NODE_ID
unset ZOOKEEPER_CONNECT

if [[ -z "$ADVERTISED_HOST_NAME" ]]; then
    ADVERTISED_HOST_NAME=localhost
fi
if [[ -z "$ADVERTISED_PORT" ]]; then
    ADVERTISED_PORT=9092
fi

: ${PORT:=9092}
: ${ADVERTISED_PORT:=9092}
: ${CONTROLLER_PORT:=9093}

: ${KAFKA_ADVERTISED_PORT:=${ADVERTISED_PORT}}
: ${KAFKA_ADVERTISED_HOST_NAME:=${ADVERTISED_HOST_NAME}}

: ${KAFKA_PORT:=${PORT}}
: ${KAFKA_HOST_NAME:=${HOST_NAME}}

if [[ -z "$CLUSTER_ID" ]]; then
    : ${KAFKA_LISTENERS:=PLAINTEXT://$KAFKA_HOST_NAME:$KAFKA_PORT}
else
    case "$NODE_ROLE" in
     'combined' ) : ${KAFKA_LISTENERS:=PLAINTEXT://$KAFKA_HOST_NAME:$KAFKA_PORT,CONTROLLER://$KAFKA_HOST_NAME:$CONTROLLER_PORT};;
     'broker' ) : ${KAFKA_LISTENERS:=PLAINTEXT://$KAFKA_HOST_NAME:$KAFKA_PORT};;
     'controller' ) : ${KAFKA_LISTENERS:=PLAINTEXT://$KAFKA_HOST_NAME:$CONTROLLER_PORT};;
     *) : ${KAFKA_LISTENERS:=PLAINTEXT://$KAFKA_HOST_NAME:$KAFKA_PORT,CONTROLLER://$KAFKA_HOST_NAME:$CONTROLLER_PORT};;
    esac
fi

: ${KAFKA_ADVERTISED_LISTENERS:=PLAINTEXT://$KAFKA_ADVERTISED_HOST_NAME:$KAFKA_ADVERTISED_PORT}

export KAFKA_LISTENERS KAFKA_ADVERTISED_LISTENERS
unset HOST_NAME ADVERTISED_HOST_NAME KAFKA_HOST_NAME KAFKA_ADVERTISED_HOST_NAME PORT ADVERTISED_PORT KAFKA_PORT KAFKA_ADVERTISED_PORT CONTROLLER_PORT NODE_ROLE

echo "Using KAFKA_LISTENERS=$KAFKA_LISTENERS and KAFKA_ADVERTISED_LISTENERS=$KAFKA_ADVERTISED_LISTENERS"

# Copy config files if not provided in volume
cp -r $KAFKA_HOME/config.orig/* $KAFKA_HOME/config

# Configure the log files ...
if [[ -z "$LOG_LEVEL" ]]; then
    LOG_LEVEL="INFO"
fi
sed -i -r -e "s|=INFO, stdout|=$LOG_LEVEL, stdout|g" $KAFKA_HOME/config/log4j.properties
sed -i -r -e "s|^(log4j.appender.stdout.threshold)=.*|\1=${LOG_LEVEL}|g" $KAFKA_HOME/config/log4j.properties
export KAFKA_LOG4J_OPTS="-Dlog4j.configuration=file:$KAFKA_HOME/config/log4j.properties"
unset LOG_LEVEL

# Add missing EOF at the end of the config file
echo "" >> $KAFKA_HOME/$CONFIG_FILE

# Process all environment variables that start with 'KAFKA_' (but not 'KAFKA_HOME' or 'KAFKA_VERSION'):
for VAR in `env`
do
  env_var=`echo "$VAR" | sed "s/=.*//"`
  if [[ $env_var =~ ^KAFKA_ && $env_var != "KAFKA_VERSION" && $env_var != "KAFKA_HOME"  && $env_var != "KAFKA_LOG4J_OPTS" && $env_var != "KAFKA_JMX_OPTS" ]]; then
    prop_name=`echo "$VAR" | sed -r "s/^KAFKA_(.*)=.*/\1/g" | tr '[:upper:]' '[:lower:]' | tr _ .`
    if egrep -q "(^|^#)$prop_name=" $KAFKA_HOME/$CONFIG_FILE; then
        #note that no config names or values may contain an '@' char
        sed -r -i "s%(^|^#)($prop_name)=(.*)%\2=${!env_var}%g" $KAFKA_HOME/$CONFIG_FILE
    else
        #echo "Adding property $prop_name=${!env_var}"
        echo "$prop_name=${!env_var}" >> $KAFKA_HOME/$CONFIG_FILE
    fi
  fi
done

if [[ ! -z "$CLUSTER_ID" && ! -f "$KAFKA_LOG_DIRS/meta.properties" ]]; then
  echo "No meta.properties found in $KAFKA_LOG_DIRS; going to format the directory"
  $KAFKA_HOME/bin/kafka-storage.sh format -t $CLUSTER_ID -c $KAFKA_HOME/$CONFIG_FILE
fi

exec $KAFKA_HOME/bin/kafka-server-start.sh $KAFKA_HOME/$CONFIG_FILE
;;