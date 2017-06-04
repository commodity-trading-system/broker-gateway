# broker-gateway
This is broker gateway of a Distributed Commodities OTC Electronic Trading System.
It contains three parts.
- Receiver
- Executor
- Querier

Receiver takes the responsibility of receiving consignations from trader gateway and 
send it into the `redis queue` after making some necessary validation.

Executor pops the consignation from `redis queue` and matching it. At the same time, it
will publish the latest depth and status onto `Etcd`

Querier is a simple http server. Trader gateway can fetch information of futrues, consignations and orders.


## Architecture
![Architecture](http://ojiqea97q.bkt.clouddn.com/docker/Screen%20Shot%202017-06-04%20at%202.34.15%20PM.png)


## How to install
First, clone the repository.

```
git clone https://github.com/commodity-trading-system/broker-gateway.git
```

#### Environment
You should prepare `Redis` `Etcd` `Mysql` by your self.
Here is a reference
```
# install redis
docker run -p 6739:6736 --name cts-redis  redis

# install etcd
docker run -itd --hostname etcd -p 2379:2379 --name etcd index.tenxcloud.com/coreos/etcd:2.3.1 /usr/local/bin/etcd -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 -advertise-client-urls http://127.0.0.1:2379,http://127.0.0.1:4001
```

#### Receiver

Add `.env` file and change it according to your environment
```
cp cmd/receiver/.env.example /cmp/receiver/.env
```

Build the dockerfile and run the container
```
docker build -f ./ReceiverDockerfile -t receiver --rm  .
docker run -p 5001:80 receiver
```

#### Executor

Add `.env` file and change it according to your environment
```
cp cmd/executor/.env.example /cmp/executor/.env
```

Build the dockerfile and run the container
```
docker build -f ./ExecutorDockerfile -t executor --rm  .
docker run executor
```

#### Querier

Add `.env` file and change it according to your environment
```
cp cmd/executor/.env.example /cmp/executor/.env
```

Build the dockerfile and run the container
```
docker build -f ./QuerierDockerfile -t querier --rm  .
docker run -p 5002:80 querier
```





