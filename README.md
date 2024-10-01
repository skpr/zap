Zap
===

A toolkit for integrating with OWASP Zap.

## Github Actions Integration

```yaml
name: Dynamic Application Security Testing

# Add your triggers here.

jobs:
  dast:
    runs-on: ubuntu-latest
    container:
      image: ghcr.io/skpr/zap:main
      options: --user root
      env:
        ZAP_WRAPPER_ENDPOINT: https://example.com
        # ZAP_WRAPPER_TYPE: active
    steps:
      - name: Run Scan
        run: |
          zap-wrapper
      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: report
          path: /tmp/report.pdf
```
