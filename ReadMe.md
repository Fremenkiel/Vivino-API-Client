# Introduction

This is a simple GO client for the [Vivino API](https://www.vivino.com/api/).

It bypassed the AWS WAF challange by generating a token before each request.

To prevent a user agent from being blocked, the client uses a random user agent and rotates it every time a request is made or if the request is rejected.

# Usage

```shell
go run main.go
```

# License

MIT
