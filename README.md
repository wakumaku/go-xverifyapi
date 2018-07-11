# xverify API Client [![Build Status](https://travis-ci.org/wakumaku/go-xverifyapi.svg?branch=master)](https://travis-ci.org/wakumaku/go-xverifyapi) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/9b66f7d42dcb413bbf96f8f4d1471020)](https://www.codacy.com/app/wakumaku/go-xverifyapi?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=wakumaku/go-xverifyapi&amp;utm_campaign=Badge_Grade) [![Code Coverage](https://scrutinizer-ci.com/g/wakumaku/go-xverifyapi/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/wakumaku/go-xverifyapi/?branch=master) [![GoDoc](https://godoc.org/github.com/wakumaku/go-xverifyapi?status.svg)](https://godoc.org/github.com/wakumaku/go-xverifyapi)
### Source: http://docs.xverify.com/

```
go get github.com/wakumaku/go-xverify
```

```
client = xverifyapi.New(apiKey, domain, nil)
verified, err := client.IsEmailVerified("email@domain.tld")
```

Makefile:
* `make test` Runs tests
