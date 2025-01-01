# PostgreSQL Troubleshooting & Query Guide

A lightweight diagnostic guide to identify table bloat, low HOT updates, stale autovacuum runs, unused index overhead, and high-frequency queries.

---

## 1. Enable Query Monitoring

Run this once per database to track statement executions, latency, and runtime statistics.

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

```

---

## 2. Table Health Dashboard Query

Lists all tables larger than 10MB with health status flags, data vs. index sizes, HOT update percentages, dead row bloat, and a list of unused index names.

```sql
WITH unused_indexes AS (
    SELECT
        s.relname AS table_name,
        COUNT(*) AS unused_count,
        pg_size_pretty(SUM(pg_relation_size(s.indexrelid))) AS wasted_size,
        string_agg(s.indexrelname, ', ' ORDER BY pg_relation_size(s.indexrelid) DESC) AS unused_list
    FROM pg_stat_user_indexes s
    JOIN pg_index i ON s.indexrelid = i.indexrelid
    WHERE s.idx_scan < 50               -- Threshold: Less than 50 scans total
      AND i.indisunique IS FALSE        -- Exclude Primary Keys & Unique constraints
    GROUP BY s.relname
)
SELECT
    t.relname AS table_name,
    CASE
        WHEN ui.unused_count > 0 THEN 'UNUSED INDEXES'
        WHEN t.n_tup_upd > 10000 AND (t.n_tup_hot_upd::numeric / t.n_tup_upd) < 0.70 THEN 'LOW HOT RATE'
        WHEN (t.n_live_tup + t.n_dead_tup) > 0 AND (t.n_dead_tup::numeric / (t.n_live_tup + t.n_dead_tup)) > 0.10 THEN 'HIGH BLOAT'
        WHEN t.last_autovacuum < NOW() - INTERVAL '14 days' OR t.last_autovacuum IS NULL THEN 'STALE VACUUM'
        ELSE 'HEALTHY'
    END AS status_flag,
    pg_size_pretty(pg_relation_size(t.relid)) AS data_size,
    pg_size_pretty(pg_total_relation_size(t.relid) - pg_relation_size(t.relid)) AS index_size,
    CASE
        WHEN t.n_tup_upd > 0 THEN ROUND((t.n_tup_hot_upd::numeric / t.n_tup_upd) * 100, 0) || '%'
        ELSE 'N/A'
    END AS hot_update_pct,
    CASE
        WHEN (t.n_live_tup + t.n_dead_tup) > 0
        THEN ROUND((t.n_dead_tup::numeric / (t.n_live_tup + t.n_dead_tup)) * 100, 1) || '%'
        ELSE '0%'
    END AS dead_row_pct,
    COALESCE(ui.wasted_size, '0 MB') AS wasted_index_space,
    COALESCE(ui.unused_count, 0) AS unused_index_cnt,
    COALESCE(ui.unused_list, 'None') AS unused_index_names,
    COALESCE(to_char(t.last_autovacuum, 'YYYY-MM-DD'), 'Never') AS last_autovacuum
FROM pg_stat_user_tables t
LEFT JOIN unused_indexes ui ON t.relname = ui.table_name
WHERE pg_total_relation_size(t.relid) > 10 * 1024 * 1024 -- >10 MB
ORDER BY pg_total_relation_size(t.relid) DESC;

```

---

## 3. Unused & Low-Scan Indexes Query

Lists all non-unique indexes with under 50 scans, sorted by disk space wasted.

```sql
SELECT
    s.schemaname,
    s.relname AS table_name,
    s.indexrelname AS index_name,
    s.idx_scan AS total_scans,
    s.idx_tup_read AS tuples_read,
    s.idx_tup_fetch AS tuples_fetched,
    pg_size_pretty(pg_relation_size(s.indexrelid)) AS index_size
FROM pg_stat_user_indexes s
JOIN pg_index i ON s.indexrelid = i.indexrelid
WHERE s.idx_scan < 50
  AND i.indisunique IS FALSE
ORDER BY pg_relation_size(s.indexrelid) DESC;

```

---

## 4. Most Frequent & Slow Queries

Finds the queries running most frequently along with total rows affected, total execution time, and average latency.

```sql
SELECT
    calls,
    rows AS total_rows_affected,
    ROUND(rows::numeric / calls, 2) AS avg_rows_per_call,
    ROUND(total_exec_time::numeric, 2) AS total_time_ms,
    ROUND(mean_exec_time::numeric, 2) AS avg_time_ms,
    query
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat%'
ORDER BY calls DESC
LIMIT 20;

```

---

## 5. Inspect Index Definitions (DDL)

Shows the exact `CREATE INDEX` statements for a specific table to check if updated columns exist inside index keys or partial `WHERE` clauses.

```sql
SELECT
    indexname,
    indexdef AS index_ddl
FROM pg_indexes
WHERE tablename = 'YOUR_TABLE_NAME'
  AND schemaname = 'public';
```
