# Configuration

## Example

```
interval: 15m
log:
  debug: false
  type: text
  file:
    enable: false
    path: default.log
elasticsearch:
  host: http://localhost:9200
  ssl_certificate_verification: true
  auth:
    enable: false
    username: elastic
    password: secret
  indices:
  - name: '*'
    rollover_pattern: '[0-9]{4}.[0-9]{2}.[0-9]{2}'
kibana:
  host: http://localhost:5601
  ssl_certificate_verification: true
  auth:
    enable: false
    username: elastic
    password: secret
xpack:
  enable: false
  spaces:
  - name: global
    pattern: '*'
    timestamp: timestamp
opendistro:
  enable: false
  tenants:
  - name: global
    pattern: '*'
    timestamp: timestamp

```
## Config Explanation

### interval

Type: string

Description: Interval Duration between current and the next job execution for scheduled job or cron job. Example: 30s, 5m, 1d.

---

### log.debug (optional)

Type: boolean

Description: Enable debug log.

### log.type

Type: string

Description: Choose the log format type plain 'text' or 'json'.

Possible values:

- text
- json

### log.file.enable (optional)

Type: boolean

Description: Enable to store log to file.

### log.file.path (optional)

Type: string

Description: Full path of the location where the logs will store.

---

### elasticsearch.host

Type: string

Description: The host of elasticsearch (full url with protocol and port if necessary).

### elasticsearch.ssl_certificate_verification (optional)

Type: boolean

Description: enable or disable ssl certificate verification.

### elasticsearch.auth.enable (optional)

Type: boolean

Description: enable or disable authentication when connecting to the elasticsearch.

### elasticsearch.auth.username (optional)

Type: string

Description: Username to be used for authentication.

### elasticsearch.auth.password (optional)

Type: string

Description: Password to be used for authentication.

### elasticsearch.indices[].name

Type: string

Description: Provide the pattern (regex) for getting the list of indexes from elasticsearch. For example:
- `'*'` pattern will get all index names.
- `'log.core'` pattern will get all index names starting from 'log.core'.

### elasticsearch.indices[].rollover_pattern (optional / recommended)

Type: string

Description: Provide the pattern (regex) for filter the rollover pattern from the index name, so the result for the index pattern will be clean instead of creating multiple index patterns for each rollover index. For example:
- `'[0-9]{4}.[0-9]{2}.[0-9]{2}'` pattern will match and truncate the rollover index from index name, for example: `log.nginx-2021.07.01` will be `log.nginx`.

---

### kibana.host

Type: string

Description: The host of kibana (full url with protocol and port if necessary).

### kibana.ssl_certificate_verification (optional)

Type: boolean

Description: enable or disable ssl certificate verification.

### kibana.auth.enable (optional)

Type: boolean

Description: enable or disable authentication when connecting to the kibana.

### kibana.auth.username (optional)

Type: string

Description: Username to be used for authentication.

### kibana.auth.password (optional)

Type: string

Description: Password to be used for authentication.

---

### xpack.enable

Type: boolean

Description: Enable sync kibana using default provider or xpack.

### xpack.spaces[].name

Type: string

Description: Space name for sync the index, if you are not using any kibana space, please write it with `global` to sync with the global space.

### xpack.spaces[].pattern

Type: string

Description: Second filter for filtering the index name, sometimes you need it when syncing to multiple spaces. For example:
- `'*'` pattern will store all index names matched.
- `'log.core'` pattern will store all index names starting from 'log.core'.

### xpack.spaces[].timestamp

Type: string

Description: Index pattern needs timestamp to use for display the data in kibana, you need to provide it, by default is using `@timestamp`. For example:
- `'log.timestamp'` will use 'log.timestamp' field for the index pattern.

---

### opendistro.enable

Type: boolean

Description: Enable sync kibana using opendistro tenants.

### opendistro.tenants[].name

Type: string

Description: Tenant name for sync the index, if you are not using any opendistro tenant, please write it with `global` to sync with the global tenant.

### opendistro.tenants[].pattern

Type: string

Description: Second filter for filtering the index name, sometimes you need it when syncing to multiple tenants. For example:
- `'*'` pattern will store all index names matched.
- `'log.core'` pattern will store all index names starting from 'log.core'.

### opendistro.tenants[].timestamp

Type: string

Description: Index pattern needs timestamp to use for display the data in kibana, you need to provide it, by default is using `@timestamp`. For example:
- `'log.timestamp'` will use 'log.timestamp' field for the index pattern.