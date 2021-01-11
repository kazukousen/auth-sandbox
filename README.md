
```shell script
make build.docker name=wasm-ext-opa
make run name=wasm-ext-opa
```

```shell script
BOB_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4iLCJzdWIiOiJZbTlpIiwibmJmIjoxNTE0ODUxMTM5LCJleHAiOjE2NDEwODE1Mzl9.WCxNAveAVAdRCmkpIObOTaSd0AJRECY2Ch2Qdic3kU8"
curl -iS http://localhost:18000 -H "authorization: Bearer ${BOB_TOKEN}"
```
