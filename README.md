# DNSGoAPI

An API written in Go for performing DNS queries

### Setup

1. Run the dockerfile with:
    `docker build Dockerfile`
2. Optional: `sudo apt install jq` for parsing the returned JSON responses 

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
curl --no-progress-meter http://localhost:8080/ns/comodo/google.com | jq
```
