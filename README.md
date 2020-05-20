<h1 align="center">
  Vanity Algos
</h1>
<p align="center">
  Go-based Algorand vanity address generator
</p>

<h4 align="center">


[![Last Commit](https://img.shields.io/github/last-commit/DJRHails/vanity-algos.svg?style=plasticr)](https://github.com/DJRHails/vanity-algos/commits/master)
[![Stars](https://img.shields.io/github/stars/DJRHails/vanity-algos.svg?style=plasticr)](https://github.com/DJRHails/vanity-algos/stargazers)
[![GitHub issues](https://img.shields.io/github/issues-raw/DJRHails/vanity-algos?style=flat)](https://github.com/DJRHails/vanity-algos/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/DJRHails/vanity-algos)](https://github.com/DJRHails/vanity-algos/pulls)
</h4>

## What's a vanity address?
A vanity address is an address which contains some personalisation, making it look less random.

### Examples
**HAILS**  `HAILSN476AQKWQ46TSIG4XGOCBLGG2CICATHB6YCAWN4SQNGH66QDKHGPQ`
**ENCODE** `ENCODE2SRL3QOTLN67FLSURTQ5R7NO3WTSAU6ZAKTOIVBPQFA3G326TGXY`

## Usage
```bash
go get github.com/DJRHails/vanity-algos

# Generate 3 new addresses with prefix FOO
vanity-algos gen -n 3 FOO

# Generate 2 new addresses matching regex
vanity-algos gen -n 2 '^[0-9]{5}'
```
