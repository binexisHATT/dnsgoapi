from requests import get

url = "http://localhost:8080/"

def query_api(record_type, dns_server, fqdn):
    resp = get(f"{url}/{record_type}/{dns_server}/{fqdn}")
    if resp.status_code != 404:
        json = resp.json()
        for val in json[fqdn]:
            print(val)

if __name__ == "__main__":
    query_api("a", "cloudflare", "google.com")
    query_api("aaaa", "cloudflare", "google.com")
    query_api("mx", "cloudflare", "google.com")
    query_api("ns", "cloudflare", "google.com")
