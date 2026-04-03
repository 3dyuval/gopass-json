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

## Prerequisites

- [gopass](https://github.com/gopasspw/gopass) — the vault backend
- [jq](https://jqlang.github.io/jq/) — for field filtering

## Install

```bash
go install github.com/3dyuval/gopass-json/cmd/gopass-json@latest
```

## Usage

| Command | Description |
|---|---|
| `gopass-json get <entry>` | All fields as a JSON object |
| `gopass-json get <entry> <jq-filter>` | Field value via any jq expression — `get cloud/infra .host` |
| `gopass-json list [pattern]` | JSON array of all entry paths — pipeable to `jq '.[]'` |
| `gopass-json find <query>` | JSON array of matching paths — `find cloud \| jq '.[0]'` |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

Contributions are welcome. To propose a change:

1. Fork the repo
2. Create a branch (`git checkout -b feature/my-change`)
3. Commit your changes (`git commit -m 'add my change'`)
4. Push the branch (`git push origin feature/my-change`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Acknowledgments

- [gopass-jsonapi](https://github.com/gopasspw/gopass-jsonapi) — this project is built on top of the gopass JSON API and uses its Go library to access the vault directly
- [gopass](https://github.com/gopasspw/gopass) — the underlying password manager
- [Best-README-Template](https://github.com/othneildrew/Best-README-Template)

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
