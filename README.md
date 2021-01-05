# MCPT Central Client

Simple client application for generating a wireguard config and interfacing with mcpt/central, outputs the config to stdout.

You must have a template for the config, an example is shown below

```
[Interface]
Address = {{ .IP }}
PrivateKey = {{ .PrivateKey }}

[Peer]
PublicKey = e5Vc7tbl0saT/sJVuG6KE/3uBs6xv0H0BDSFAnXnAAA=
AllowedIPs = {{ .AllowedIPs }}
Endpoint = 192.168.47.142:443
```