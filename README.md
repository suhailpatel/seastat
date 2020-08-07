# Seastat 🌊

Seastat is a standalone Cassandra Prometheus Exporter on top of Jolokia built for speed 🏎️

# Introduction

Seastat is a Prometheus Exporter for Cassandra written in Go. The goal was to build a standalone and opinionated Prometheus Exporter built for speed, especially if you have a lot of tables spread across lots of keyspaces.

Seastat is different to other exporters for Cassandra such as the [JMX Exporter](https://github.com/prometheus/jmx_exporter) or [cassandra-exporter](https://github.com/instaclustr/cassandra-exporter). Seastat is designed to be run standalone and updates metrics in the background (to seperate concerns between gathering metrics and serving metrics). Seastat only gathers metrics at a fixed configurable interval.

# Performance

Seastat is used for scraping metrics for more than 1,000 tables across hundreds of keyspaces every minute without sweat 😅. It is built for performance by batching queries when it makes sense and limiting the amount of data it exposes to be scalable. More metrics may be added in the future but with careful consideration to not negatively impact performance.

A very (non-scientific) test with 4000 tables across 200 keyspaces took between 10-15 seconds to scrape all stats exposed. Both the standalone Cassandra Exporter and the Prometheus JMX Exporter took over 10 minutes because they query for each MBean for each table individually which is very expensive. This test was done using Cassandra running in the Docker harness (with 4 cores and 8GB of RAM on  a completely idle cluster of 1). Your mileage may vary and you should do your own tests!

# Requirements

Seastat doesn't speak JMX directly. Instead, it uses [Jolokia](https://jolokia.org/) to translate back and forth into JMX. You will need Jolokia to be embedded as an agent into your Cassandra process. Jolokia versions 1.3+ will work just fine (the exporter has been tested with Jolokia v1.3 and v1.6).

Seastat has been designed on top of Cassandra 3.0 (specifically, 3.0.18). It may work with 3.11+ but some of the metric types may have changed between the two versions which may result in zero values.

# Metrics Exposed

Seastar exposes the metrics in categories. If you want more information about the metrics in particular, look at the [Cassandra Metrics](http://cassandra.apache.org/doc/latest/operating/metrics.html) documentation.

All metrics exported are defined in the code within a [single file](https://github.com/suhailpatel/seastat/blob/master/server/prom_metrics.go).

## Table Stat Metrics

These metrics have a labels of `keyspace` and `table` applied to them

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_table_coordinator_read_latency_seconds` | Read Latency for queries to the table which this node coordinates | Summary |
| `seastat_table_coordinator_write_latency_seconds` | Write Latency for queries to the table which this node coordinates | Summary |
| `seastat_table_coordinator_range_scan_latency_seconds` | Range Scan Latency for queries to the table which this node coordinates | Summary |
| `seastat_table_read_latency_seconds` | Read Latency for queries which this node is involved in | Summary |
| `seastat_table_write_latency_seconds` | Write Latency for queries which this node is involved in | Summary |
| `seastat_table_range_scan_latency_seconds` | Range Scan Latency for queries which this node is involved in | Summary |
| `seastat_table_cas_propose_latency_seconds` | Compare and Set Propose Latency for queries | Summary |
| `seastat_table_cas_commit_latency_seconds` | Compare and Set Commit Latency for queries | Summary |
| `seastat_table_estimated_partitions` | Number of partitions in this table (estimated) | Gauge |
| `seastat_table_pending_compactions` | Number of pending compactions on this table | Gauge |
| `seastat_table_live_disk_space_used_bytes` | Disk space used for live cells in bytes | Gauge |
| `seastat_table_total_disk_space_used_bytes` | Disk space used for all data in bytes | Gauge |
| `seastat_table_live_sstables` | Number of live SSTables | Gauge |
| `seastat_table_sstables_per_read` | Number of SSTables consulted per read query | Summary |
| `seastat_table_max_partition_size_bytes` | Max Partition Size in bytes | Gauge |
| `seastat_table_mean_partition_size_bytes` | Mean Partition Size in bytes | Gauge |
| `seastat_table_bloom_filter_false_ratio` | False positive ratio of table’s bloom filter | Gauge |
| `seastat_table_tombstones_scanned` | Number of tombstones scanned per read query | Summary |
| `seastat_table_live_cells_scanned` | Number of live cells scanned per read query | Summary |
| `seastat_table_key_cache_hit_percent` | Percent of key cache hits | Gauge
| `seastat_table_repaired_percent` | Percent of table repaired | Gauge
| `seastat_table_speculative_retries_total` | Total amount of speculative retries | Counter
| `seastat_table_speculative_failed_retries_total` | Total amount of speculative failed retries | Counter

## CQL Metrics

These CQL metrics do not have any labels

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_cql_prepared_statements` | Number of prepared statements | Gauge |
| `seastat_cql_prepared_statements_evicted_total` | Number of evicted prepared statements | Counter |
| `seastat_cql_prepared_statements_executed_total` | Number of executed prepared statements | Counter |
| `seastat_cql_regular_statements_executed_total` | Number of executed regular statements | Counter |
| `seastat_cql_prepared_statements_ratio` | Ratio of prepared statements | Gauge |

## Thread Pool Metrics

These metrics are labelled by the Thread Pool name in `name`

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_thread_pool_active_tasks` | Number of active tasks in this thread pool | Gauge |
| `seastat_thread_pool_pending_tasks` | Number of pending tasks in this thread pool | Gauge |
| `seastat_thread_pool_completed_tasks_total` | Number of completed tasks in this thread pool | Counter |
| `seastat_thread_pool_blocked_tasks_total` | Number of total blocked tasks in this thread pool | Counter |
| `seastat_thread_pool_currently_blocked_tasks` | Number of currently blocked tasks in this thread pool | Gauge |
| `seastat_thread_pool_max_pool_size` | Largest thread pool size | Gauge |

## Compaction Metrics

These Compaction metrics do not have any labels

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_compaction_bytes_compacted_total` | Total amount of bytes compacted across all compactions | Counter |
| `seastat_compaction_pending_tasks` | Number of pending compaction tasks | Gauge |
| `seastat_compaction_completed_tasks_total` | Number of completed compaction tasks | Counter |

## Client Request Metrics

These Client Request metrics are tagged by Request Type in `request_type`

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_client_request_latency_seconds` | Coordinator request latency | Summary |
| `seastat_client_request_timeout_total` | Total number of coordinated request timeouts | Counter |
| `seastat_client_request_failure_total` | Total number of coordinated request failures | Counter |
| `seastat_client_request_unavailable_total` | Total number of coordinated request unavailable | Counter |

## Connected Clients Metrics

This metric does not have any labels

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_connected_clients` | Number of connected clients | Gauge |
| `seastat_client_request_timeout_total` | Total number of coordinated request timeouts | Counter |

## Memory Metrics

These metrics are from the Java process itself and have no labels

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_memory_heap_used_bytes` | Bytes representing the used memory heap size | Gauge |
| `seastat_memory_nonheap_used_bytes` | Bytes representing the used memory non-heap size | Gauge |

## Garbage Collection Metrics

These metrics are from the Java process itself. Each metric has a single label `name` which represents the type of GC that's occurred

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_gc_total` | Total number of Garbage Collections | Counter |
| `seastat_gc_last_duration_seconds` | Duration of Last GC | Gauge |
| `seastat_gc_accumulated_duration_seconds` | Accumulated durations of GC | Counter |

## Storage Metrics

These metrics come from Cassandra's storage service which keeps track of the cluster state from the perspective of each node

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_storage_keyspaces` | Number of keyspaces reported by Cassandra | Gauge |
| `seastat_storage_tokens` | Number of tokens reported by Cassandra  | Gauge |
| `seastat_storage_node_status` | Status (`live`, `unreachable`, `joining`, `moving`, `leaving`) of each node in the cluster (tagged by node and status) | Gauge |

## Hint Metrics

These metrics come from the [Storage](https://cassandra.apache.org/doc/latest/operating/metrics.html#storage-metrics) metric which keeps track of hints, node load and storage exceptions.

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
|TotalHintsInProgress|Number of hints attemping to be sent currently.|Counter|
|TotalHints|Number of hint messages written to this node since [re]start. Includes one entry for each host to be hinted per hint.|Counter|

## Scrape Metrics

Seastat also exposes some internal metrics of how long the scrape took and the timestamp of the last scrape

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_last_scrape_timestamp` | Unix timestamp of the last scrape | Gauge |
| `seastat_last_scrape_duration_seconds` | Duration of the last scrape | Gauge |

# Usage

**Note:** Seastat is in infancy, changes to the interface will be made until it reaches 1.0.0 💪

Building Seastat is just like building any other Go application. You will need Go 1.13 or above to build Seastat.

You can use the included `make` targets

```shell
$ # To build a version for your current OS and Arch
$ make build

$ # To build a version for Linux 64-bit
$ make build-linux
```

To run Seastat

```shell
$ # To run on port 8080 (defaults to INFO logging and above)
$ ./seastat server -p 8080

$ # To run on port 8080 with debug logging
$ ./seastat server -p 8080 -v debug
```

# Things to work on

- Seastat does not support Jolokia auth
- More batching of requests can achieve more speed!
- The code has been written to be easily tested, but needs some more tests!

# Author

Suhail Patel <<me@suhailpatel.com>>
