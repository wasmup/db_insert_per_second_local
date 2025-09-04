# Db Insert Per Second Using Local Docker Containers
Simple DB insert per second

---

# PostgreSQL
```sh
docker pull postgres:latest
docker run -it --rm postgres:latest psql -V
# psql (PostgreSQL) 17.6 (Debian 17.6-1.pgdg13+1)
```

# Result
Inserted 1063 rows in 1.00 seconds (1062.32 inserts/second)


```sql

CREATE TABLE impressions (
    impression_id VARCHAR(100) PRIMARY KEY,
    ad_id VARCHAR(100) NOT NULL,
    image_url VARCHAR(4000) NOT NULL,
    click_url VARCHAR(4000) NOT NULL,
    impression_time TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);

```

[Source](postgresql/main.go)


---

# Oracle
```sh
# 23.9.0.25.07
docker pull container-registry.oracle.com/database/free:latest
```

# Result
Inserted 860 rows in 1.00 seconds (858.82 inserts/second)

```sql
CREATE TABLE impressions (
    impression_id VARCHAR2(100) PRIMARY KEY,
    ad_id VARCHAR2(100) NOT NULL,
    image_url VARCHAR2(4000) NOT NULL,
    click_url VARCHAR2(4000) NOT NULL,
    impression_time TIMESTAMP WITH TIME ZONE DEFAULT (SYSTIMESTAMP AT TIME ZONE 'UTC')
);

```

[Source](oracle/main.go)


---

# Cassandra single node

# Result
Inserted 685 rows in 1.00 seconds (684.80 inserts/second)


```sh
docker pull cassandra:latest
docker run  --name cassandra1 -d --network host  cassandra:latest
docker logs -f cassandra1
docker exec -it cassandra1 cqlsh

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


[Source](cassandra/main.go)


---

# ScyllaDB single node

# Result
Inserted 587 rows in 1.00 seconds (586.98 inserts/second)
 
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
docker stop scylla-node1
docker rm scylla-node1

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

[Source](scylladb/main.go)


---

# Redis
 
[Source](redis/main.go)

# Result

Inserted 32784 impressions in 1.00 seconds (32783.21 inserts/second)

---

# Redis Pipeline
[Source](redis-pipeline/main.go)


# Result

Inserted 437000 impressions in 1.00 seconds (436891.20 inserts/second)

---
