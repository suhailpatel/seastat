# Seastat 🌊

Seastat is a standalone Cassandra Prometheus Exporter on top of Jolokia built for speed 🏎️

## Introduction

Seastat is a Prometheus Exporter for Cassandra written in Go. The goal was to build a standalone and opinionated Prometheus Exporter built for speed, especially if you have a lot of tables spread across lots of keyspaces.

Seastat is different to other exporters for Cassandra such as the [JMX Exporter](https://github.com/prometheus/jmx_exporter) or [cassandra-exporter](https://github.com/instaclustr/cassandra-exporter). Seastat is designed to be run standalone and updates metrics in the background (to seperate concerns between gathering metrics and serving metrics). Seastat only gathers metrics at a fixed configurable interval.

## Performance

Seastat is used for scraping metrics for more than 1,000 tables across hundreds of keyspaces every minute without sweat 😅. It is built for performance by batching queries when it makes sense and limiting the amount of data it exposes to be scalable. More metrics may be added in the future but with careful consideration to not negatively impact performance.

A very (non-scientific) test with 4000 tables across 200 keyspaces took between 7-10 seconds to scrape all stats exposed. Both the standalone Cassandra Exporter and the Prometheus JMX Exporter took over 10 minutes because they query for each MBean for each table individually which is very expensive. This test was done using Cassandra running in the Docker harness (with 4 cores and 8GB of RAM). Your mileage may vary and you should do your own tests!

## Requirements

Seastat doesn't speak JMX directly. Instead, it uses [Jolokia](https://jolokia.org/) to translate back and forth into JMX. You will need Jolokia to be embedded as an agent into your Cassandra process. Jolokia versions 1.3+ will work just fine (the exporter has been tested with Jolokia v1.3 and v1.6).

Seastat has been designed on top of Cassandra 3.0 (specifically, 3.0.18). It may work with 3.11+ but some of the metric types may have changed between the two versions which may result in zero values.

## Metrics Exposed

Seastar exposes the metrics in categories. If you want more information about the metrics in particular, look at the [Cassandra Metrics](http://cassandra.apache.org/doc/latest/operating/metrics.html) documentation

### Table Stat Metrics

These metrics have a labels of `keyspace` and `table` applied to them

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_table_coordinator_read_latency_seconds` | Read Latency for queries to the table which this node coordinates | Summary |
| `seastat_table_coordinator_write_latency_seconds` | Write Latency for queries to the table which this node coordinates | Summary |
| `seastat_table_coordinator_range_scan_latency_seconds` | Range Scan Latency for queries to the table which this node coordinates | Summary |
| `seastat_table_read_latency_seconds` | Read Latency for queries which this node is involved in | Summary |
| `seastat_table_write_latency_seconds` | Write Latency for queries which this node is involved in | Summary |
| `seastat_table_range_scan_latency_seconds` | Range Scan Latency for queries which this node is involved in | Summary |
| `seastat_table_estimated_partitions` | Number of partitions in this table (estimated) | Gauge |
| `seastat_table_pending_compactions` | Number of pending compactions on this table | Gauge |
| `seastat_table_max_partition_size_bytes` | Max Partition Size in bytes | Gauge |
| `seastat_table_mean_partition_size_bytes` | Mean Partition Size in bytes | Gauge |
| `seastat_table_bloom_filter_false_ratio` | False positive ratio of table’s bloom filter | Gauge |
| `seastat_table_key_cache_hit_percent` | Percent of key cache hits | Gauge
| `seastat_table_repaired_percent` | Percent of table repaired | Gauge

### CQL Metrics

These CQL metrics do not have any labels

| Name          | Description   | Type |
| ------------- | ------------- | ---- |
| `seastat_cql_prepared_statements` | Number of prepared statements | Gauge |
| `seastat_cql_prepared_statements_evicted_total` | Number of evicted prepared statements | Counter |
| `seastat_cql_prepared_statements_executed_total` | Number of executed prepared statements | Counter |
| `seastat_cql_regular_statements_executed_total` | Number of executed regular statements | Counter |
| `seastat_cql_prepared_statements_ratio` | Ratio of prepared statements | Gauge |

## Usage

TODO

## Things to work on

- Seastat does not support Jolokia auth
- More batching of requests can achieve more speed!

## Author

Suhail Patel <<me@suhailpatel.com>>
