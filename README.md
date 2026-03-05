# cloud-geoip

IP prefix lists for cloud/CDN providers, generated from IRR data using [bgpq4](https://github.com/bgp/bgpq4).

## Supported providers

- Akamai
- Alibaba Cloud
- Cognosphere / HoYoverse

## Downloads

Prebuilt files are published to the [`release`](../../tree/release) branch and updated daily.

| Format | Path |
|--------|------|
| Plain text | `plain/<name>.txt` |
| [sing-box](https://sing-box.sagernet.org/) SRS | `srs/<name>.srs` |
| [Xray-core](https://github.com/XTLS/Xray-core) dat | `cloud-geoip.dat` |

### sing-box usage

```json
{
  "rule_set": [
    {
      "tag": "akamai",
      "type": "remote",
      "format": "binary",
      "url": "https://raw.githubusercontent.com/<owner>/cloud-geoip/release/srs/akamai.srs"
    }
  ]
}
```

### Xray-core usage

Place `cloud-geoip.dat` in the Xray asset directory (or set `XRAY_LOCATION_ASSET`), then reference it with the `ext:` prefix:

```json
{
  "routing": {
    "rules": [
      {
        "type": "field",
        "outboundTag": "direct",
        "ip": ["ext:cloud-geoip.dat:akamai"]
      }
    ]
  }
}
```

All entries are in a single file. The tag after the colon is the target name (e.g. `akamai`, `alibaba`, `cognosphere`).

## Local usage

```bash
# requires: bgpq4
go run . -output out
```

This generates `out/plain/*.txt`, `out/srs/*.srs`, and `out/cloud-geoip.dat`.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-output` | `output` | Output directory |
| `-S` | | IRR sources (passed to bgpq4 `-S`) |
| `-h` | | IRR server (passed to bgpq4 `-h`) |
| `-no-v4` | `false` | Skip IPv4 |
| `-no-v6` | `false` | Skip IPv6 |
| `-no-aggregate` | `false` | Disable bgpq4 prefix aggregation (`-A`) |
