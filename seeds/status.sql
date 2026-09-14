SELECT pid state, query_start, query_start
FROM pg_stat_activity
WHERE
    query LIKE '%pg_sleep'
    AND state != 'idle';