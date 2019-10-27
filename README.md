# Seasat 🌊

A standalone Cassandra Prometheus Exporter on top of Jolokia built for speed 🏎️

## Introduction

Seastat is a Prometheus Exporter for Cassandra written in Go. The goal was to build a standalone and opinionated Prometheus Exporter built for speed, especially if you have a lot of tables. 

Seastat is used for scraping metrics for more than 1,000 tables across hundreds of keyspaces every minute without sweat 😅

Seastat is different to other exporters for Cassandra such as the [JMX Exporter](https://github.com/prometheus/jmx_exporter) or [cassandra-exporter](https://github.com/instaclustr/cassandra-exporter). 

Seastat is designed to be run standalone and updates metrics in the background (to seperate concerns between gathering metrics and serving metrics). Seastat only gathers metrics at a fixed configurable interval.

## Requirements

Seastat doesn't speak JMX directly. Instead, it uses [Jolokia](https://jolokia.org/) to translate back and forth into JMX. You will need Jolokia to be embedded into your Cassandra process

**Note:** Seastat does not support Jolokia auth. I'm hoping to add this soon!

## Metrics Exposed

Seastat exposes a few of the JMX metrics exposed by Cassandra. It focuses on 

## Usage

TODO

## Author

Suhail Patel <<me@suhailpatel.com>>
