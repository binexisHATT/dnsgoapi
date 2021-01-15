# DNSGoAPI

An API written in Go for performing DNS queries

### Setup

1. Run with Docker:
    ```
    docker build . -t dnsgoapi
    docker --rm -it -p 8080:8080 dnsgoapi
    ```
    OR build the binary and run it:
    ```
    make all && ./dnsapi
    ```
2. Optional: In another terminal window, run `sudo apt install jq` to install `jq` for parsing JSON responses 

### Public DNS Servers

`dnsgoapi` current supports the following DNS servers:

- cloudflare -> `1.1.1.1:53`
- google -> `8.8.8.8:53`
- opendns -> `208.67.222.222:53`
- comodo -> `8.26.56.26:53`
- quad9 -> `9.9.9.9:53`
- verisign -> `64.6.64.6:53`

### Syntax for querying API

```
/record_type/dns_server/fully_qualified_domain_name
```
**Note**: both `record_type` and `dns_server` are case insensitive, thus `/aaaa` is equivalent to `/AAAA`

### Python Example

```python
from requests import get

url = "http://localhost:8080/"

def query_api(record_type, dns_server, fqdn):
    resp = get(f"{url}/{record_type}/{dns_server}/{fqdn}")
    if resp.status_code != 404:
        json = resp.json()
        for val in json[fqdn]:
            print(val)

if __name__ == "__main__":
    # Requesting A records
    query_api("a", "cloudflare", "google.com")
    # Requesting AAAA records
    query_api("aaaa", "cloudflare", "google.com")
    # Requesting MX records
    query_api("mx", "cloudflare", "google.com")
    # Requesting NS records
    query_api("ns", "cloudflare", "google.com")
```

### Requesting A Records

```
curl --no-progress-meter http://localhost:8080/a/cloudflare/google.com | jq
```

### Requesting AAAA Records

```
curl --no-progress-meter http://localhost:8080/aaaa/quad9/google.com | jq
```

### Requesting CNAME Records

```
curl --no-progress-meter http://localhost:8080/cname/opendns/www.github.com | jq
```

### Requesting MX Records

```
curl --no-progress-meter http://localhost:8080/mx/quad9/google.com | jq
```

### Requesting NS Records

```
curl --no-progress-meter http://localhost:8080/ns/google/google.com | jq
```

### Requesting TXT Records

```
curl --no-progress-meter http://localhost:8080/txt/verisign/google.com | jq
```

### Requesting CAA Records

```
curl --no-progress-meter http://localhost:8080/caa/cloudflare/google.com | jq
```

### Requesting SOA Records

```
curl --no-progress-meter http://localhost:8080/soa/quad9/google.com | jq
```
