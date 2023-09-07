# loki-range
Run loki_range queries with no limit

# usage
query 'sum(count_over_time({app="app-name"} | json | msg="messaging success")) by (msg)' --range "10m"

# rnage
- today
- yesterday
- {x}d
- {x}h
- {x}m

# start, end
RFC3339 time

