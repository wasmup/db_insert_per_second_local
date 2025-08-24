
 
```sh
docker network ls
docker network rm scylla-network
docker network create scylla-network
docker run --name scylla-node1 --network scylla-network -d --restart unless-stopped scylladb/scylla:latest
# docker run --name scylla-node2 --network scylla-network -d --restart unless-stopped  scylladb/scylla:latest   --seeds=scylla-node1
# docker run --name scylla-node3 --network scylla-network -d --restart unless-stopped   scylladb/scylla:latest   --seeds=scylla-node1


docker ps
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' 1d1125344480
nc -zv 172.22.0.2 1521


docker exec -it scylla-node1 cqlsh

```

```sql
CREATE KEYSPACE ks1 WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
DESCRIBE ks1

CREATE TABLE ks1.impressions (
    impression_id VARCHAR PRIMARY KEY,
    ad_id VARCHAR,
    image_url VARCHAR,
    click_url VARCHAR,
    impression_time TIMESTAMP
);


```


Inserted 587 rows in 1.00 seconds (586.98 inserts/second)