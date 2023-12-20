# <img src="https://static-00.iconduck.com/assets.00/prometheus-icon-511x512-1vmxbcxr.png" width="50"/>  Auto Prometheus Exporter
[![Release](https://github.com/MGSousa/prom-exporter/actions/workflows/release.yml/badge.svg)](https://github.com/MGSousa/prom-exporter/actions/workflows/release.yml)

Prometheus exporter that fetches and auto-resolves JSON data from selected services by converting metrics/stats in Prometheus format, being ready to be scrapped.

Compatible with Elastic Beats plugins (Filebeat, Auditbeat, Packetbeat, Metricbeat, etc.) and any other tool that has any HTTP endpoint with output in JSON.

## Install
### Linux / Darwin 
  - specify version
```sh
VERSION=1.0.0
wget -nv https://github.com/MGSousa/prom-exporter/releases/download/v$VERSION/prom-exporter_${VERSION}_$(uname | awk '{print tolower($0)}')_amd64 -O prom-exporter && chmod +x prom-exporter
```

  - or fetch the latest release
```sh
curl -fsL "https://api.github.com/repos/MGSousa/prom-exporter/releases/latest" |\
    jq -r ".assets[] | select(.name|contains(\"$(uname | awk '{print tolower($0)}')\")) | .url" |\
    wget --header="Accept: application/octet-stream" -O prom-exporter -nv -i - && chmod +x prom-exporter
```
### Windows
```sh
curl -o prom-exporter.exe https://github.com/MGSousa/prom-exporter/releases/download/v1.0.0/prom-exporter_1.0.0_windows_amd64.exe
```

## Usage
 - About args
```sh
Usage of ./prom-exporter:
  -debug
    	Enable debug mode.
  -listen-address string
    	Address to listen on to be scraped. (default ":19100")
  -service-metrics-path string
    	Service path to scrape metrics from. (default "metrics")
  -service-name string
    	Remote service name to reference.
  -service-port string
    	HTTP Port of the remote service. (default "80")
  -service-protocol string
    	HTTP Schema of the remote service (http or https). (default "http")
  -service-uri string
    	Endpoint address of the remote service.
  -service-version-scrape
    	Enable whether the service will be internally scraped for fetching remote build version or not.
  -telemetry-path string
    	Base path under which to expose metrics. (default "/metrics")
```

 - Run it via docker compose / swarm
### filebeat.yaml
```yaml
fields_under_root: true
filebeat.autodiscover.providers:
   - type: docker
     templates:
       - condition.and:
          - not.contains:
              docker.container.labels.name: "filebeat"
         config:
           - type: container
             paths:
               - "/var/lib/docker/containers/${data.docker.container.id}/*.log"
processors:
  - add_docker_metadata:
      host: "unix:///var/run/docker.sock"
output.console:
  pretty: true
logging:
  level: "info"
  metrics.enabled: false
http:
  enabled: true
  host: 0.0.0.0
```
### compose.yaml
```yaml
# configuration example

networks:
  filebeat:
    external: true

services:
  filebeat:
    hostname: &filebeat_endpoint "CUSTOM_HOSTNAME"
    container_name: filebeat
    image: docker.elastic.co/beats/filebeat:8.11.1
    command: -environment container
    user: root
    volumes:
      - ${PWD}/filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - /var/lib/docker:/var/lib/docker:ro
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - filebeat
  exporter:
    container_name: filebeat_exporter
    image: filebeat-exporter:latest
    command: |
      -listen-address ":9201"
      -service-name "filebeat"
      -service-port "5066"
      -service-metrics-path "stats"
    environment:
      SERVICE_ENDPOINT: *filebeat_endpoint
    networks:
      - filebeat
    ports:
      - '9201:9201'
```
 - Run it via executable
```sh
./prom-exporter -service-name "SERVICE_NAME" -service-uri "REMOTE_HOST" -service-port REMOTE_PORT -service-metrics-path "REMOTE_PATH"
```

## TODO
[] Set desired metrics with custom ValueType (Gauge, Histogram, Counter, etc.) via mapping
