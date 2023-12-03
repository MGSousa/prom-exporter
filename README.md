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
```sh
./prom-exporter -service-uri "REMOTE_HOST:PORT"-service-name "SERVICE_NAME" -service-metrics-path "REMOTE_PATH"
