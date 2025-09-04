# Db Insert Per Second Using Local Docker Containers
Simple DB insert per second

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

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	dsn := os.ExpandEnv("user=$PSQL_USERNAME password=$PSQL_PASSWORD host=localhost port=5432 dbname=$PSQL_DB sslmode=disable")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("db open failed", "error", err)
		return
	}
	defer db.Close()

	query := `INSERT INTO impressions (impression_id, ad_id, image_url, click_url)
              VALUES ($1, $2, $3, $4)`

	maxDuration := 1 * time.Second
	startTime := time.Now()
	count := 0
	ctx := context.Background()

	for time.Since(startTime) < maxDuration {
		impressionID := uuid.New().String()
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"

		_, err := db.ExecContext(ctx, query, impressionID, adID, imageURL, clickURL)
		if err != nil {
			slog.Error("db insert failed", "error", err)
			return
		}
		count++
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d rows in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}
```

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

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
	_ "time/tzdata"

	_ "github.com/godror/godror"
	"github.com/google/uuid"
)

func main() {
	loc, err := time.LoadLocation(`UTC`)
	if err != nil {
		panic(err)
	}
	time.Local = loc

	dsn := os.ExpandEnv("$ORACLE_USERNAME/$ORACLE_PASSWORD@$ORACLE_HOST/$ORACLE_DB")
	db, err := sql.Open("godror", dsn)
	if err != nil {
		slog.Error(`db open failed`, `error`, err)
		return
	}
	defer db.Close()

	query := `INSERT INTO impressions (impression_id, ad_id, image_url, click_url )
              VALUES (:1, :2, :3, :4 )`

	maxDuration := 1 * time.Second
	startTime := time.Now()
	count := 0
	var ctx = context.Background()

	for time.Since(startTime) < maxDuration {
		impressionID := uuid.New().String()
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"

		_, err := db.ExecContext(ctx, query, impressionID, adID, imageURL, clickURL)
		if err != nil {
			slog.Error(`db insert failed`, `error`, err)
			return
		}

		count++
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d rows in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}
```

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

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "ks1"
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer session.Close()

	query := `INSERT INTO impressions (impression_id, ad_id, image_url, click_url, impression_time)
              VALUES (?, ?, ?, ?, ?)`

	maxDuration := 1 * time.Second
	startTime := time.Now()
	count := 0

	for time.Since(startTime) < maxDuration {
		impressionID := uuid.New().String()
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"
		impressionTime := time.Now().UTC()

		if err := session.Query(query, impressionID, adID, imageURL, clickURL, impressionTime).Exec(); err != nil {
			log.Fatalf("Insert failed: %v", err)
		}
		count++
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d rows in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}

```

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

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func main() {
	cluster := gocql.NewCluster("172.22.0.2")
	cluster.Keyspace = "ks1"
	cluster.Consistency = gocql.One

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer session.Close()

	query := `INSERT INTO impressions (impression_id, ad_id, image_url, click_url, impression_time)
              VALUES (?, ?, ?, ?, ?)`

	maxDuration := 1 * time.Second
	startTime := time.Now()
	count := 0

	for time.Since(startTime) < maxDuration {
		impressionID := uuid.New().String()
		adID := "ad456"
		imageURL := "https://example.com/image.jpg"
		clickURL := "https://advertiser.com/click"
		impressionTime := time.Now().UTC()

		if err := session.Query(query, impressionID, adID, imageURL, clickURL, impressionTime).Exec(); err != nil {
			log.Fatalf("Insert failed: %v", err)
		}
		count++
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Inserted %d rows in %.2f seconds (%.2f inserts/second)\n", count, elapsed.Seconds(), float64(count)/elapsed.Seconds())
}
```

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
