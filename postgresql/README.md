
 

```sql

CREATE TABLE impressions (
    impression_id VARCHAR(100) PRIMARY KEY,
    ad_id VARCHAR(100) NOT NULL,
    image_url VARCHAR(4000) NOT NULL,
    click_url VARCHAR(4000) NOT NULL,
    impression_time TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);

```


Inserted 1063 rows in 1.00 seconds (1062.32 inserts/second)