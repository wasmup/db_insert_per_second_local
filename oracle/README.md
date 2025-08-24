
 

```sql
CREATE TABLE impressions (
    impression_id VARCHAR2(100) PRIMARY KEY,
    ad_id VARCHAR2(100) NOT NULL,
    image_url VARCHAR2(4000) NOT NULL,
    click_url VARCHAR2(4000) NOT NULL,
    impression_time TIMESTAMP WITH TIME ZONE DEFAULT (SYSTIMESTAMP AT TIME ZONE 'UTC')
);

```



Inserted 860 rows in 1.00 seconds (858.82 inserts/second)