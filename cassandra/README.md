
 
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


Inserted 685 rows in 1.00 seconds (684.80 inserts/second)