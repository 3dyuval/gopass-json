<div id="readme-top"></div>

[![Stars][stars-shield]][stars-url]
[![Forks][forks-shield]][forks-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

# gopass-json

Treat your gopass vault as JSON routes. Query any secret field with `jq`.

```bash
gopass-json get cloud/infra .host
gopass-json get cloud/infra '.["api-token"]'
gopass-json list
gopass-json find cloud
```

## Install

```bash
go install github.com/3dyuval/gopass-json/cmd/gopass-json@latest
```

## Usage

| Command | Output |
|---|---|
| `gopass-json get <entry>` | All fields as JSON |
| `gopass-json get <entry> <jq-filter>` | Single field value |
| `gopass-json list [pattern]` | JSON array of matching entry paths |
| `gopass-json find <query>` | JSON array of entries matching query |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS -->
[stars-shield]: https://img.shields.io/github/stars/3dyuval/gopass-json.svg?style=for-the-badge
[stars-url]: https://github.com/3dyuval/gopass-json/stargazers
[forks-shield]: https://img.shields.io/github/forks/3dyuval/gopass-json.svg?style=for-the-badge
[forks-url]: https://github.com/3dyuval/gopass-json/network/members
[issues-shield]: https://img.shields.io/github/issues/3dyuval/gopass-json.svg?style=for-the-badge
[issues-url]: https://github.com/3dyuval/gopass-json/issues
[license-shield]: https://img.shields.io/github/license/3dyuval/gopass-json.svg?style=for-the-badge
[license-url]: https://github.com/3dyuval/gopass-json/blob/master/LICENSE
