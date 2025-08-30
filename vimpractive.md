```js
<script>console.log("hello there")</script>
```

#### Hey tool

```bash
$ go install github.com/rakyll/hey@latest
```

This is for loading test data

### Router

```bash
https://github.com/julienschmidt/httprouter
```

#### Generate skeleton directory strcuture

```bash
$ mkdir -p bin cmd/api internal migrations remote
$ touch Makefile
$ touch cmd/api/main.go

```

#### Explanation

- bin: contains compiled binary
- cmd/api: contains application specific codes - running a server, writing http requests, auth management
- internal: database interaction, data validation, sending emails
- migrations: db migrations
- remote: congig files, setup scripts
- go.mod : declare project delivery
- Makefile: automation

#### Run code manually

Add this to test your code
