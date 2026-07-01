<h1 align="center">
  <br>Nali<br>
</h1>

<h4 align="center">An offline tool for querying IP geographic information and CDN provider.</h4>

<p align="center">
  <a href="https://github.com/zu1k/nali/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/zu1k/nali/go.yml?branch=master&style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/zu1k/nali">
    <img src="https://goreportcard.com/badge/github.com/zu1k/nali?style=flat-square">
  </a>
  <a href="https://github.com/zu1k/nali/releases">
    <img src="https://img.shields.io/github/release/zu1k/nali/all.svg?style=flat-square">
  </a>
  <a href="https://github.com/zu1k/nali/releases">
    <img src="https://img.shields.io/github/downloads/zu1k/nali/total?style=flat-square">
  </a>
</p>

#### [中文文档](https://github.com/zu1k/nali/blob/master/README.md)

## Features

- Multi database support
  - Chunzhen qqwry IPv4 database
  - ZX IPv6 database
  - Geoip2 city database
  - IPIP free database
  - ip2region database
  - DB-IP database
  - IP2Location DB3 LITE database
- CDN provider query
- Pipeline support
- Interactive query
- Both IPv4 and IPv6 supported
- Multilingual support
- Offline query
- Full platform support (Linux / macOS / Windows / FreeBSD, including multiple MIPS architectures)
- Colorized output
- JSON format output
- Shell completion support
- Built-in LRU cache with output buffering for high performance

## Install

### Install from source

Nali requires Go >= 1.23. You can build it from source:

```sh
$ go install github.com/zu1k/nali@latest
```

### Install pre-built binary

Pre-built binaries are available here: [Release](https://github.com/zu1k/nali/releases)

Download the binary compatible with your platform, unpack and copy to the directory in path.

### Arch Linux

We have published 3 packages in AUR:

- `nali-go`: release version, compile when installing
- `nali-go-bin`: release version, pre-compiled binary
- `nali-go-git`: the latest master branch version, compile when installing

### Docker

```sh
$ docker pull ghcr.io/zu1k/nali:latest
$ docker run --rm nali 1.2.3.4
```

Or build the image from source:

```sh
$ docker build -t nali .
$ docker run --rm nali 1.2.3.4
```

## Usage

### Query a simple IP address

```sh
$ nali 1.2.3.4
1.2.3.4 [澳大利亚 APNIC Debogon-prefix网络]
```

#### or use `pipe`

```sh
$ echo IP 6.6.6.6 | nali
IP 6.6.6.6 [美国 亚利桑那州华楚卡堡市美国国防部网络中心]
```

### Query multiple IP addresses

```sh
$ nali 1.2.3.4 4.3.2.1 123.23.3.0
1.2.3.4 [澳大利亚 APNIC Debogon-prefix网络]
4.3.2.1 [美国 新泽西州纽瓦克市Level3Communications]
123.23.3.0 [越南 越南邮电集团公司]
```

### Interactive query

use `exit` or `quit` to quit

```sh
$ nali
123.23.23.23
123.23.23.23 [越南 越南邮电集团公司]
1.0.0.1
1.0.0.1 [美国 APNIC&CloudFlare公共DNS服务器]
8.8.8.8
8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]
quit
```

### JSON output

Use `-j` or `--json` for JSON formatted output, making it easy to integrate with other tools:

```sh
$ nali -j 1.2.3.4
{"ip":"1.2.3.4","text":"澳大利亚 APNIC Debogon-prefix网络","source":"qqwry"}
```

Pipeline mode is also supported:

```sh
$ echo 1.2.3.4 6.6.6.6 | nali -j
{"ip":"1.2.3.4","text":"澳大利亚 APNIC Debogon-prefix网络","source":"qqwry"}
{"ip":"6.6.6.6","text":"美国 亚利桑那州华楚卡堡市美国国防部网络中心","source":"qqwry"}
```

### Use with `dig`

```sh
$ dig nali.zu1k.com +short | nali
104.28.2.115 [美国 CloudFlare公司CDN节点]
104.28.3.115 [美国 CloudFlare公司CDN节点]
172.67.135.48 [美国 CloudFlare节点]
```

### Use with `nslookup`

```sh
$ nslookup nali.zu1k.com 8.8.8.8 | nali
Server:         8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]
Address:        8.8.8.8 [美国 加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器]#53

Non-authoritative answer:
Name:   nali.zu1k.com
Address: 104.28.3.115 [美国 CloudFlare公司CDN节点]
Name:   nali.zu1k.com
Address: 104.28.2.115 [美国 CloudFlare公司CDN节点]
Name:   nali.zu1k.com
Address: 172.67.135.48 [美国 CloudFlare节点]
```

### Use with any other program

Because nali can read the contents of the `stdin` pipeline, it can be used with any program.

```sh
bash abc.sh | nali
```

Nali will insert IP information after IP address, and CDN provider information after CDN domains.

### IPv6 support

Use like IPv4. Nali also automatically converts NAT64 addresses (`64:ff9b::/96`) to IPv4 for lookup.

```sh
$ nslookup google.com | nali
Server:         127.0.0.53 [局域网 IP]
Address:        127.0.0.53 [局域网 IP]#53

Non-authoritative answer:
Name:   google.com
Address: 216.58.211.110 [美国 Google全球边缘网络]
Name:   google.com
Address: 2a00:1450:400e:809::200e [荷兰Amsterdam Google Inc. 服务器网段]
```

### Query CDN provider

Since CDN services usually use CNAME domain resolution, it is recommended to use with `nslookup` or `dig`. Can also be used standalone when you already know the CNAME.

```sh
$ nslookup www.gov.cn | nali
Server:         127.0.0.53 [局域网 IP]
Address:        127.0.0.53 [局域网 IP]#53

Non-authoritative answer:
www.gov.cn      canonical name = www.gov.cn.bsgslb.cn [白山云 CDN].
www.gov.cn.bsgslb.cn [白山云 CDN]       canonical name = zgovweb.v.bsgslb.cn [白山云 CDN].
Name:   zgovweb.v.bsgslb.cn [白山云 CDN]
Address: 103.104.170.25 [新加坡 ]
Name:   zgovweb.v.bsgslb.cn [白山云 CDN]
Address: 2001:428:6402:21b::5 [美国Louisiana州Monroe Qwest Communications Company, LLC (CenturyLink)]
Name:   zgovweb.v.bsgslb.cn [白山云 CDN]
Address: 2001:428:6402:21b::6 [美国Louisiana州Monroe Qwest Communications Company, LLC (CenturyLink)]
```

## User Interface

### Help

```sh
$ nali --help
An offline tool for querying IP geographic information.

Find document on: https://github.com/zu1k/nali

Usage:
  nali [flags]
  nali [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  info        get the necessary information of nali
  update      update qqwry, zxipv6wry, ip2region ip database and cdn, update nali to latest version if -v

Flags:
      --gbk       Use GBK decoder
  -h, --help      help for nali
  -j, --json      Output in JSON format
  -v, --version   version for nali

Use "nali [command] --help" for more information about a command.
```

### Show current configuration

```sh
$ nali info
Nali Version:      v1.0.0
Config Dir Path:   /Users/user/Library/Application Support/nali
DB Data Dir Path:  /Users/user/Library/Application Support/nali
Selected IPv4 DB:  qqwry
Selected IPv6 DB:  zxipv6wry
Selected CDN DB:   cdn
```

### Database configuration

After nali runs for the first time, a configuration file `config.yaml` will be generated in the config directory (use `nali info` to check the exact path). The configuration file defines the database information.

A database is defined as follows:

```yaml
databases:
  - name: geoip
    name-alias:
    - geolite
    - geolite2
    format: mmdb
    file: GeoLite2-City.mmdb
    languages:
    - ALL
    types:
    - IPv4
    - IPv6
```

### Update database

Update all databases if available:

```sh
$ nali update
```

Update specified databases:

```sh
$ nali update --db qqwry,cdn
```

### Update Nali itself

Use the `-v` flag to also update nali to the latest version:

```sh
$ nali update -v
```

Or update both database and nali:

```sh
$ nali update --db qqwry,cdn -v
```

### Specify database

Users can specify which database to use via environment variables `NALI_DB_IP4`, `NALI_DB_IP6`, `NALI_DB_CDN`, or through the `selected` field in the `config.yaml` configuration file.

Supported databases:

- Geoip2 `['geoip', 'geoip2']`
- Chunzhen `['chunzhen', 'qqwry']`
- IPIP `['ipip']`
- Ip2Region `['ip2region', 'i2r']`
- DBIP `['dbip', 'db-ip']`
- IP2Location `['ip2location']`
- CDN `['cdn']`

#### Windows

##### Use geoip db

```sh
set NALI_DB_IP4=geoip

or use powershell

$env:NALI_DB_IP4="geoip"
```

##### Use ipip db

```sh
set NALI_DB_IP6=ipip

or use powershell

$env:NALI_DB_IP6="ipip"
```

#### Linux

##### Use geoip db

```sh
export NALI_DB_IP4=geoip
```

##### Use ipip db

```sh
export NALI_DB_IP6=ipip
```

##### Use custom CDN db

```sh
export NALI_DB_CDN=cdn
```

### Multilingual support

Specify the language to be used by modifying the environment variable `NALI_LANG`. When using a non-Chinese language, only the GeoIP2 database is supported.

The values that can be set for this parameter can be found in the list of supported languages for GeoIP2:

```sh
# NALI_LANG=en nali 1.1.1.1
1.1.1.1 [Australia]
```

### Change working directory

Set the environment variable `NALI_HOME` to specify the working directory where both the configuration file and database files are stored. You can also use absolute paths in the configuration file to specify other database paths.

To set them separately, use `NALI_CONFIG_HOME` for the configuration file directory and `NALI_DB_HOME` for the database file directory.

If no environment variable is specified, the XDG specification will be used, with the configuration file directory in `$XDG_CONFIG_HOME/nali` and the database file directory in `$XDG_DATA_HOME/nali`.

```sh
set NALI_HOME=D:\nali

or

export NALI_HOME=/var/nali
```

### Shell completion

nali supports generating autocompletion scripts for bash, zsh, fish, and powershell:

```sh
# bash
$ source <(nali completion bash)

# zsh
$ source <(nali completion zsh)

# fish
$ nali completion fish | source

# powershell
$ nali completion powershell | Out-String | Invoke-Expression
```

To make it permanent, add the corresponding completion script to your shell configuration file.

## Thanks

- [纯真QQIP离线数据库](http://www.cz88.net)
- [qqwry纯真数据库解析](https://github.com/yinheli/qqwry)
- [ZX公网ipv6数据库](https://ip.zxinc.org/ipquery/)
- [Geoip2 city数据库](https://www.maxmind.com/en/geoip2-precision-city-service)
- [geoip2-golang解析器](https://github.com/oschwald/geoip2-golang)
- [CDN provider数据库](https://github.com/SukkaLab/cdn)
- [IPIP数据库](https://www.ipip.net/product/ip.html)
- [IPIP数据库解析](https://github.com/ipipdotnet/ipdb-go)
- [ip2region数据库](https://github.com/lionsoul2014/ip2region)
- [IP2Location DB3 LITE](https://lite.ip2location.com/database/db3-ip-country-region-city)
- [Cobra CLI库](https://github.com/spf13/cobra)

Thanks to JetBrains for the Open Source License 

<a href="https://www.jetbrains.com/?from=nali">
  <img src="assets/GoLand.svg">
</a>

## Author

**Nali** © [zu1k](https://github.com/zu1k), Released under the [MIT](./LICENSE) License.<br>
