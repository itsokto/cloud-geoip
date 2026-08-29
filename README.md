# cloud-geoip

IP prefix lists for cloud/CDN providers, generated from live BGP announcement data.

## Supported providers

| Provider | Source |
|----------|--------|
| Akamai | `AS-AKAMAI` |
| Alibaba Cloud | `AS37963`, `AS45102`, `AS24429`, `AS134963`, `AS203513` |
| Tencent | `AS45090`, `AS132203`, `AS133478`, `AS137876` |
| Cognosphere / HoYoverse | `AS203923` |
| UCloud | `AS135377`, `AS139327` |

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
      // or https://cdn.jsdelivr.net/gh/itsokto/cloud-geoip@release/srs/akamai.srs
      "url": "https://raw.githubusercontent.com/itsokto/cloud-geoip/release/srs/akamai.srs"
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

All entries are in a single file. The tag after the colon is the target name (e.g. `akamai`, `alibaba`, `tencent`, `cognosphere`, `ucloud`).

## Local usage

```bash
go run . -output out
```

This generates `out/plain/*.txt`, `out/srs/*.srs`, and `out/cloud-geoip.dat`.

## Flags

| Flag | Default | Description                                        |
|------|---------|----------------------------------------------------|
| `-output` | `output` | Output directory                                   |
| `-no-v4` | `false` | Skip IPv4                                          |
| `-no-v6` | `false` | Skip IPv6                                          |
| `-no-aggregate` | `false` | Keep prefixes as published instead of merging them |
| `-v` | `false` | Verbose logging                                    |
