# go-network-route-monitor
This exporter reads a yaml file to get a list of hosts and checks if host is accessible on the port or not and exports metrics for prometheus at :2112/metrics

`network_route_up{endpoint="google.com:80",name="Google"} 1`
`network_route_up{endpoint="172.21.3.40:8080",name="tunnel3"} 1`

## hosts.yaml
```
hosts:
  - name: Google
	address: google.com
    port: 80
  - name: tunnel3
    address: 172.21.3.40
    port: 8080
```

## Usage
gonetcheck -config path/to/hosts.yaml

