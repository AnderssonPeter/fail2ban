# Fail2ban plugin for traefik

[![Build Status](https://travis-ci.com/tomMoulard/fail2ban.svg?branch=main)](https://travis-ci.com/tomMoulard/fail2ban)

This plugin is a small implementation of a fail2ban instance as a middleware
plugin for Traefik.

## Configuration

### Whitelist
You can whitelist some IP using this:
```yml
testData:
  whitelist:
    files:
      - "tests/test-ipfile.txt"
    ip:
      - "::1"
      - "127.0.0.1"
```

Where you can use some IP in an array of files or directly in the
configuration.

### Blacklist
Like whitelist, you can blacklist some IP using this:
```yml
testData:
  blacklist:
    files:
      - "tests/test-ipfile.txt"
    ip:
      - "::1"
      - "127.0.0.1"
```

Where you can use some IP in an array of files or directly in the
configuration.

### LogLevel
In order to help you and us when building and using the plugin, we added some
logs on stdout.
You can choose the level of logging with this:

```yml
testData:
  logLevel: "INFO"
```

<details>

There is 3 level of logging :

#### `NONE`
The plugin will not output *any* logs.

```
INFO[0000] Configuration loaded from file: ./traefik.yml
```

#### `INFO`
Configuration informations will be displayed.

```
INFO[0000] Configuration loaded from file: ./traefik.yml
INFO: Fail2Ban: restricted.go:51: Whitelisted: '127.0.0.2/32'
INFO: Fail2Ban: restricted.go:51: Blacklisted: '127.0.0.3/32'
INFO: Fail2Ban: restricted.go:51: Bantime: 3h0m0s
INFO: Fail2Ban: restricted.go:51: Findtime: 3h0m0s
INFO: Fail2Ban: restricted.go:51: FailToBan Rules : '{Xbantime:3h0m0s Xfindtime:3h0m0s Xurlregexp:[localhost:5000/whoami] Xmaxretry:4 Xenabled:true}'
INFO: Fail2Ban: restricted.go:52: Plugin: FailToBan is up and running
INFO: Fail2Ban: restricted.go:51: Whitelisted: '127.0.0.2/32'
INFO: Fail2Ban: restricted.go:51: Blacklisted: '127.0.0.3/32'
INFO: Fail2Ban: restricted.go:51: Bantime: 3h0m0s
INFO: Fail2Ban: restricted.go:51: Findtime: 3h0m0s
INFO: Fail2Ban: restricted.go:51: FailToBan Rules : '{Xbantime:3h0m0s Xfindtime:3h0m0s Xurlregexp:[localhost:5000/whoami] Xmaxretry:4 Xenabled:true}'
INFO: Fail2Ban: restricted.go:52: Plugin: FailToBan is up and running
```

#### `DEBUG`
Every event will be logged.

Warning, all IPs will be prompted in clear text with this option.

```
INFO[0000] Configuration loaded from file: ./traefik.yml
INFO: Fail2Ban: restricted.go:51: Whitelisted: '127.0.0.2/32'
INFO: Fail2Ban: restricted.go:51: Blacklisted: '127.0.0.3/32'
INFO: Fail2Ban: restricted.go:51: Bantime: 3s
INFO: Fail2Ban: restricted.go:51: Findtime: 3h0m0s
INFO: Fail2Ban: restricted.go:51: FailToBan Rules : '{Xbantime:3s Xfindtime:3h0m0s Xurlregexp:[localhost:5000/whoami] Xmaxretry:4 Xenabled:true}'
INFO: Fail2Ban: restricted.go:52: Plugin: FailToBan is up and running
DEBUG: Fail2Ban: restricted.go:51: New request: &{GET /whoami HTTP/1.1 1 1
DEBUG: Fail2Ban: restricted.go:51: welcome ::1
DEBUG: Fail2Ban: restricted.go:51: New request: &{GET /whoami HTTP/1.1 1 1
DEBUG: Fail2Ban: restricted.go:51: welcome back ::1 for the 2 time
DEBUG: Fail2Ban: restricted.go:51: New request: &{GET /whoami HTTP/1.1 1 1
DEBUG: Fail2Ban: restricted.go:51: welcome back ::1 for the 3 time
DEBUG: Fail2Ban: restricted.go:51: New request: &{GET /whoami HTTP/1.1 1 1
DEBUG: Fail2Ban: restricted.go:52: ::1 is now banned temporarily
DEBUG: Fail2Ban: restricted.go:51: New request: &{GET /whoami HTTP/1.1 1 1
DEBUG: Fail2Ban: restricted.go:51: ::1 is still banned since 2021-04-23T21:40:55+02:00, 5 request
DEBUG: Fail2Ban: restricted.go:51: New request: &{GET /whoami HTTP/1.1 1 1
DEBUG: Fail2Ban: restricted.go:52: ::1 is no longer banned
```

</details>

## Fail2ban
We plan to use all [default fail2ban configuration]() but at this time only a
few features are implemented:
```yml
testData:
  logLevel: "INFO"
  rules:
    urlregexps:
    - regexp: "/no"
      mode: block
    - regexp: "/yes"
      mode: allow
    bantime: "3h"
    findtime: "10m"
    maxretry: 4
    enabled: true
```

Where:
 - `findtime`: is the time slot used to count requests (if there is too many
requests with the same ip in this slot of time, the ip goes into ban). You can
use 'smart' strings: "4h", "2m", "1s", ...
 - `bantime`: correspond to the amount of time the IP is in Ban mode.
 - `maxretry`: number of request before Ban mode.
 - `enabled`: allow to enable or disable the plugin (must be set to `true` to
enable the plugin).
 - `urlregexp`: a regexp list to block / allow requests with regexps on the url
 - `logLevel`: is used to show the correct level of logs (`DEBUG`, `INFO`,
`NONE`)

#### URL Regexp
Urlregexp are used to defined witch part of your website will be either
allowed, blocked or filtered :
- allow : all requests where the url match the regexp will be forwarded to the
backend without any check
- block : all requests where the url match the regexp will be stopped

##### No definitions

```yml
testData:
  rules:
    bantime: "3h"
    findtime: "10m"
    maxretry: 4
    enabled: true
```

By default, fail2ban will be applied.

##### Multiple definition

```yml
testData:
  rules:
    urlregexps:
    - regexp: "/whoami"
      mode: allow
    - regexp: "/do-not-access"
      mode: block
    bantime: "3h"
    findtime: "10m"
    maxretry: 4
    enabled: true
```

In the case where you define multiple regexp on the same url, the order of
process will be :
1. Block
2. Allow

In this example, all requests to `/do-not-access` will be denied and all
requests to `/whoami` will be allowed without any fail2ban interaction.

#### Schema
First request, IP is added to the Pool, and the `findtime` timer is started:
```
A |------------->
  ↑
```

Second request, `findtime` is not yet finished thus the request is fine:
```
A |--x---------->
     ↑
```

Third request, `maxretry` is now full, this request is fine but the next wont.
```
A |--x--x------->
        ↑
```

Fourth request, too bad, now it's jail time, next request will go through after
`bantime`:
```
A |--x--x--x---->
           ↓
B          |------------->
```

Fifth request, the IP is in Ban mode, nothing happen:
```
A |--x--x--x---->
B          |--x---------->
              ↑
```

Last request, the `bantime` is now over, another `findtime` is started:
```
A |--x--x--x---->            |------------->
                             ↑
B          |--x---------->
```

## Dev `traefik.yml` configuration file for traefik

```yml
pilot:
  token: [REDACTED]

experimental:
  devPlugin:
    goPath: /home/${USER}/go
    moduleName: github.com/tomMoulard/fail2ban

entryPoints:
  http:
    address: ":8000"
    forwardedHeaders:
      insecure: true

api:
  dashboard: true
  insecure: true

providers:
  file:
    filename: rules-fail2ban.yaml
```

## How to dev
```bash
$ docker run -d --network host containous/whoami -port 5000
# traefik --configfile traefik.yml
```

# Authors
| Tom Moulard | Clément David | Martin Huvelle | Alexandre Bossut-Lasry |
|-------------|---------------|----------------|------------------------|
|[![](img/gopher-tom_moulard.png)](https://tom.moulard.org)|[![](img/gopher-clement_david.png)](https://github.com/cledavid)|[![](img/gopher-martin_huvelle.png)](https://github.com/nitra-mfs)|[![](img/gopher-alexandre_bossut-lasry.png)](https://www.linkedin.com/in/alexandre-bossut-lasry/)|
