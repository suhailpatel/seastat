# Seastat 🌊

Seastat is a standalone Cassandra Prometheus Exporter on top of Jolokia built for speed 🏎️

## Introduction

Seastat is a Prometheus Exporter for Cassandra written in Go. The goal was to build a standalone and opinionated Prometheus Exporter built for speed, especially if you have a lot of tables spread across lots of keyspaces. 

Seastat is different to other exporters for Cassandra such as the [JMX Exporter](https://github.com/prometheus/jmx_exporter) or [cassandra-exporter](https://github.com/instaclustr/cassandra-exporter). Seastat is designed to be run standalone and updates metrics in the background (to seperate concerns between gathering metrics and serving metrics). Seastat only gathers metrics at a fixed configurable interval.

## Performance

Seastat is used for scraping metrics for more than 1,000 tables across hundreds of keyspaces every minute without sweat 😅. It is built for performance by batching queries when it makes sense and limiting the amount of data it exposes to be scalable. More metrics may be added in the future but with careful consideration to not negatively impact performance. 

A very (non-scientific) test with 4000 tables across 200 keyspaces took approximately 7.5 seconds to scrape all stats exposed. Both the standalone Cassandra Exporter and the Prometheus JMX Exporter took over 10 minutes because they query for each MBean for each table individually which is very expensive. This test was done using Cassandra running in the Docker harness (with 4 cores and 8GB of RAM). Your mileage may vary and you should do your own tests!

## Requirements

Seastat doesn't speak JMX directly. Instead, it uses [Jolokia](https://jolokia.org/) to translate back and forth into JMX. You will need Jolokia to be embedded as an agent into your Cassandra process. Jolokia versions 1.3+ will work just fine (the exporter has been tested with Jolokia v1.3 and v1.6).

Seastat has been designed on top of Cassandra 3.0 (specifically, 3.0.18). It may work with 3.11+ but some of the metric types may have changed between the two versions which may result in zero values.

## Metrics Exposed

TODO

## Usage

TODO

## Things to work on

- Seastat does not support Jolokia auth
- More batching of requests can achieve more speed!

## Author

Suhail Patel <<me@suhailpatel.com>>
