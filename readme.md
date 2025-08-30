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

```go
package main

import (
	"fmt"
	"net/http"
)

func main() {

	port := "8080"
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("Go running on port 8080")
	http.ListenAndServe(addr, nil)

}

```

```bash
go run cmd/api
```

##### Basic http Server

cmd/api/main.go

```go
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Declare a string containing the application version number. Later in the book we'll
// generate this automatically at build time, but for now we'll just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// Declare an instance of the config struct.
	var cfg config
	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	// Initialize a new structured logger which writes log entries to the standard out
	// stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config: cfg,
		logger: logger,
	}
	// Declare a new servemux and add a /v1/healthcheck route which dispatches requests
	// to the healthcheckHandler method (which we will create in a moment).
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	// Start the HTTP server.
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

```

And in cmd/api/healthcheck.go

```go
package main

import (
	"fmt"
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
```

![image1](assets/img1.png)

```bash
# Run with flags
$ go run ./cmd/api -port=3030 -env=production
```

```bash
# -i for including headers
 curl -i localhost:4000/v1/healthcheck
```

![img2](assets/img2.png)

#### Download HttpRouter

```bash
go get github.com/julienschmidt/httprouter@v1
```

```bash
curl localhost:4000/v1/movies/123
curl -X POST localhost:4000/v1/movies
curl -i -X OPTIONS localhost:4000/v1/healthcheck
```

##### BASIC ROUTER WITH HTTPROUTER

```go
package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}

```

##### Working with Json

```go
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
// Create a fixed-format JSON response from a string. Notice how we're using a raw
// string literal (enclosed with backticks) so that we can include double-quote
// characters in the JSON without needing to escape them? We also use the %q verb to
// wrap the interpolated values in double-quotes.
js := `{"status": "available", "environment": %q, "version": %q}`
js = fmt.Sprintf(js, app.config.env, version)
// Set the "Content-Type: application/json" header on the response. If you forget to
// this, Go will default to sending a "Content-Type: text/plain; charset=utf-8"
// header instead.
w.Header().Set("Content-Type", "application/json")
// Write the JSON as the HTTP response body.
w.Write([]byte(js))
}

```

#### Encoding objects to json

```go
You can either call the json.Marshal() function, or you can declare and use a
json.Encoder type.

```

```go
package main
import (
"encoding/json" // New import
"net/http"
)
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
// Create a map which holds the information that we want to send in the response.
data := map[string]string{
"status": "available",
"environment": app.config.env,
"version": version,
}
// Pass the map to the json.Marshal() function. This returns a []byte slice
// containing the encoded JSON. If there was an error, we log it and send the client
// a generic error message.
js, err := json.Marshal(data)
if err != nil {
app.logger.Error(err.Error())
http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
return
}
// Append a newline to the JSON. This is just a small nicety to make it easier to
// view in terminal applications.
js = append(js, '\n')
// At this point we know that encoding the data worked without any problems, so we
// can safely set any necessary HTTP headers for a successful response.
w.Header().Set("Content-Type", "application/json")
// Use w.Write() to send the []byte slice containing the JSON as the response body.
w.Write(js)
}

```

##### Json Response Helper

```go
// Define a writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
// Encode the data to JSON, returning the error if there was one.
js, err := json.Marshal(data)
if err != nil {
return err
}
// Append a newline to make it easier to view in terminal applications.
js = append(js, '\n')
// At this point, we know that we won't encounter any more errors before writing the
// response, so it's safe to add any headers that we want to include. We loop
// through the header map and add each header to the http.ResponseWriter header map.
// Note that it's OK if the provided header map is nil. Go doesn't throw an error
// if you try to range over (or generally, read from) a nil map.
for key, value := range headers {
w.Header()[key] = value
}
// Add the "Content-Type: application/json" header, then write the status code and
// JSON response.
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
w.Write(js)
return nil
}

```

#### Using the encoder method to return and convert object to json (not recommended)

```go
func (app *application) exampleHandler(w http.ResponseWriter, r *http.Request) {
data := map[string]string{
"hello": "world",
}
// Set the "Content-Type: application/json" header on the response.
w.Header().Set("Content-Type", "application/json")
// Use the json.NewEncoder() function to initialize a json.Encoder instance that
// writes to the http.ResponseWriter. Then we call its Encode() method, passing in
// the data that we want to encode to JSON (which in this case is the map above). If
// the data can be successfully encoded to JSON, it will then be written to our
// http.ResponseWriter.
err := json.NewEncoder(w).Encode(data)
if err != nil {
app.logger.Error(err.Error())
http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
}
}

```

#### Using struct

```go
package data

import (
	"time"
)

// Ensure the fields are exportable i.e in Capital letter
type Movie struct {
	ID        int64     // Unique integer ID for the movie
	CreatedAt time.Time // Timestamp for when the movie is added to our database
	Title     string    // Movie title
	Year      int32     // Movie release year
	Runtime   int32     // Movie runtime (in minutes)
	Genres    []string  // Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     // The version number starts at 1 and will be incremented each
	// time the movie information is updated
}


```

```go
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Create a new instance of the Movie struct, containing the ID we extracted from
	// the URL and some dummy data. Also notice that we deliberately haven't set a
	// value for the Year field.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

```

![image4](assets/image4.png)

#### Struct tags

Probably the most common use of struct tags is to change the key names that appear in the
JSON object.

To illustrate how to do this, let’s annotate our Movies struct with struct tags so that it uses
snake_case for the keys instead.

```go
package data
...
// Annotate the Movie struct with struct tags to control how the keys appear in the
// JSON-encoded output.
type Movie struct {
ID int64 `json:"id"`
CreatedAt time.Time `json:"created_at"`
Title string `json:"title"`
Year int32 `json:"year"`
Runtime int32 `json:"runtime"`
Genres []string `json:"genres"`
Version int32 `json:"version"`
}

```

![img5](assets/img5.png)

#### Hiding struct fields in the JSON object

It’s also possible to control the visibility of individual struct fields in the JSON by using the
omitzero and - struct tag directives.

```go
type Movie struct {
ID int64 `json:"id"`
CreatedAt time.Time `json:"-"` // Use the - directive
Title string `json:"title"`
Year int32 `json:"year,omitzero"` // Add the omitzero directive
Runtime int32 `json:"runtime,omitzero"` // Add the omitzero directive
Genres []string `json:"genres,omitzero"` // Add the omitzero directive
Version int32 `json:"version"`
}
```

![img6](assets/img6.png)
Note: You can also prevent a struct field from appearing in the JSON output by simply
making it unexported. But using the json:"-" struct tag is generally a better choice:
it’s an explicit indication to both Go and any future readers of your code that you don’t
want the field included in the JSON, and it helps prevents problems if someone
changes the field to be exported in the future without realizing the consequences.

The omitempty directive is similar to omitzero but less consistent in its behavior. It
differs from omitzero in a few important ways:

```go
// For example, if you wanted to completely omit the Genres field from the JSON whenever it
// contains no values (or is nil), you could use the omitempty directive like so:

type Movie struct {
ID int64 `json:"id"`
CreatedAt time.Time `json:"-"`
Title string `json:"title"`
Year int32 `json:"year,omitzero"`
Runtime int32 `json:"runtime,omitzero"`
Genres []string `json:"genres,omitempty"` // Use the omitempty directive
Version int32 `json:"version"`
}

```

#### The string directive

A final, less-frequently-used, struct tag directive is string. You can use this on individual
struct fields to force the data to be represented as a string in the JSON output.
For example, if we wanted the value of our Runtime field to be represented as a JSON string
(instead of a number) we could use the string directive like this:

```go
type Movie struct {
ID int64 `json:"id"`
CreatedAt time.Time `json:"-"`
Title string `json:"title"`
Year int32 `json:"year,omitzero"`
Runtime int32 `json:"runtime,omitzero,string"` // Add the string directive
Genres []string `json:"genres,omitzero"`
Version int32 `json:"version"`
}

```

Result will be like:

```go
{
"id": 123,
"title": "Casablanca",
"runtime": "102", ← This is now a string
"genres": [
"drama",
"romance",
"war"
],
"version": 1
}

```

_Note_ that the string directive will only work on struct fields which have int*, uint*, float\*
or bool types. For any other type of struct field it will have no effect.

#### Json marshal indent

Help to turn actual JSON response data is all just on one line with no whitespace to readable format

```go
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
// Use the json.MarshalIndent() function so that whitespace is added to the encoded
// JSON. Here we use no line prefix ("") and tab indents ("\t") for each element.
js, err := json.MarshalIndent(data, "", "\t")
if err != nil {
return err
}
js = append(js, '\n')
for key, value := range headers {
w.Header()[key] = value
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
w.Write(js)
return nil
}
```

![img7](assets/img7.png)

json.MarshalIndent() takes 65% longer to run and uses around 30% more memory than json.Marshal(), as well as making two more heap allocations.

#### Enveloping responses

cmd/api/helpers.go

```go
package main
...
// Define an envelope type.
type envelope map[string]any
// Change the data parameter to have the type envelope instead of any.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
js, err := json.MarshalIndent(data, "", "\t")
if err != nil {
return err
}
js = append(js, '\n')
for key, value := range headers {
w.Header()[key] = value
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
w.Write(js)
return nil
```

then :

```go
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Create a new instance of the Movie struct, containing the ID we extracted from
	// the URL and some dummy data. Also notice that we deliberately haven't set a
	// value for the Year field.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

```

Also update the `healthcheck`:

```go
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

```

Call this: curl localhost:4000/v1/movies/123

![img8](assets/img8.png)

#### Customizing the Runtime field

When our Movie struct is encoded to JSON, the Runtime field (which is an int32 type) is currently formatted as a JSON number. Let’s change this so that it’s encoded as a string with the format "<runtime> mins" instead. Like so:

```go
File: internal/data/runtime.go
package data
import (
"fmt"
"strconv"
)
// Declare a custom Runtime type, which has the underlying type int32 (the same as our
// Movie struct field).
type Runtime int32
// Implement a MarshalJSON() method on the Runtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
func (r Runtime) MarshalJSON() ([]byte, error) {
// Generate a string containing the movie runtime in the required format.
jsonValue := fmt.Sprintf("%d mins", r)
// Use the strconv.Quote() function on the string to wrap it in double quotes. It
// needs to be surrounded by double quotes in order to be a valid *JSON string*.
quotedJSONValue := strconv.Quote(jsonValue)
// Convert the quoted string value to a byte slice and return it.
return []byte(quotedJSONValue), nil
}

```

#### Then change the movie struct

```go
package data
import (
"time"
)
type Movie struct {
ID int64 `json:"id"`
CreatedAt time.Time `json:"-"`
Title string `json:"title"`
Year int32 `json:"year,omitzero"`
// Use the Runtime type instead of int32. Note that the omitzero directive will
// still work on this: if the Runtime field has the underlying value 0, then it will
// be considered zero and omitted -- and the MarshalJSON() method we just made
// won't be called at all.
Runtime Runtime `json:"runtime,omitzero"`
Genres []string `json:"genres,omitzero"`
Version int32 `json:"version"`
}

```

#### SO you get

```bash
$ curl localhost:4000/v1/movies/123
{
"movie": {
	"id": 123,
	"title": "Casablanca",
	"runtime": "102 mins",
	"genres": [
	"drama",
	"romance",
	"war"
	],
	"version": 1
	}
}
```

#### - Customizing the Movie struct

```go
// Note that there are no struct tags on the Movie struct itself.
type Movie struct {
ID int64
CreatedAt time.Time
Title string
Year int32
Runtime int32
Genres []string
Version int32
}
// Implement a MarshalJSON() method on the Movie struct, so that it satisfies the
// json.Marshaler interface.
func (m Movie) MarshalJSON() ([]byte, error) {
// Declare a variable to hold the custom runtime string (this will be the empty
// string "" by default).
var runtime string
// If the value of the Runtime field is not zero, set the runtime variable to be a
// string in the format "<runtime> mins".
if m.Runtime != 0 {
runtime = fmt.Sprintf("%d mins", m.Runtime)
}
// Create an anonymous struct to hold the data for JSON encoding. This has exactly
// the same fields, types and tags as our Movie struct, except that the Runtime
// field here is a string, instead of an int32. Also notice that we don't include
// a CreatedAt field at all (there's no point including one, because we don't want
// it to appear in the JSON output).
aux := struct {
ID int64 `json:"id"`
Title string `json:"title"`
Year int32 `json:"year,omitzero"`
Runtime string `json:"runtime,omitzero"` // This is a string.
Genres []string `json:"genres,omitzero"`
Version int32 `json:"version"`
}{
// Set the values for the anonymous struct.
ID: m.ID,
Title: m.Title,
Year: m.Year,
Runtime: runtime, // Note that we assign the value from the runtime variable here.
Genres: m.Genres,
Version: m.Version,
}
// Encode the anonymous struct to JSON, and return it.
return json.Marshal(aux)
}

```

##### Error message

```go

```

#### ErrorHandling

```go
package main

import (
	"fmt"
	"net/http"
)

// The logError() method is a generic helper for logging an error message along
// with the current request method and URL as attributes in the log entry.
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using the any
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}
	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// The serverErrorResponse() method will be used when our application encounters an
// unexpected problem at runtime. It logs the detailed error message, then uses the
// errorResponse() helper to send a 500 Internal Server Error status code and JSON
// response (containing a generic error message) to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 Not Found status code and
// JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

```

Then change the handler

```go
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
env := envelope{
"status": "available",
"system_info": map[string]string{
"environment": app.config.env,
"version": version,
},
}
err := app.writeJSON(w, http.StatusOK, env, nil)
if err != nil {
// Use the new serverErrorResponse() helper.
app.serverErrorResponse(w, r, err)
}
}

```

#### Handling Panics

```go
package main
import (
"fmt"
"net/http"
)
func (app *application) recoverPanic(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Create a deferred function (which will always be run in the event of a panic
// as Go unwinds the stack).
defer func() {
// Use the builtin recover function to check if there has been a panic or
// not.
if err := recover(); err != nil {
// If there was a panic, set a "Connection: close" header on the
// response. This acts as a trigger to make Go's HTTP server
// automatically close the current connection after a response has been
// sent.
w.Header().Set("Connection", "close")
// The value returned by recover() has the type any, so we use
// fmt.Errorf() to normalize it into an error and call our
// serverErrorResponse() helper. In turn, this will log the error using
// our custom Logger type at the ERROR level and send the client a 500
// Internal Server Error response.
app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
}
}()
next.ServeHTTP(w, r)
})
}

```

Then in routes

```go
func (app *application) routes() http.Handler {
router := httprouter.New()
router.NotFound = http.HandlerFunc(app.notFoundResponse)
router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
// Wrap the router with the panic recovery middleware.
return app.recoverPanic(router)
}
```

Test

```bash
$ curl -i -H "Host: こんにちは" http://localhost:4000/v1/healthcheck
HTTP/1.1 400 Bad Request: malformed Host header
Content-Type: text/plain; charset=utf-8
Connection: close
400 Bad Request: malformed Host header

```

#### Decoding JSON i.e parse json to Golang's object

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
// Declare an anonymous struct to hold the information that we expect to be in the
// HTTP request body (note that the field names and types in the struct are a subset
// of the Movie struct that we created earlier). This struct will be our *target
// decode destination*.
var input struct {
Title string `json:"title"`
Year int32 `json:"year"`
Runtime int32 `json:"runtime"`
Genres []string `json:"genres"`
}
// Initialize a new json.Decoder instance which reads from the request body, and
// then use the Decode() method to decode the body contents into the input struct.
// Importantly, notice that when we call Decode() we pass a *pointer* to the input
// struct as the target decode destination. If there was an error during decoding,
// we also use our generic errorResponse() helper to send the client a 400 Bad
// Request response containing the error message.
err := json.NewDecoder(r.Body).Decode(&input)
if err != nil {
app.errorResponse(w, r, http.StatusBadRequest, err.Error())
return
}
// Dump the contents of the input struct in a HTTP response.
fmt.Fprintf(w, "%+v\n", input)
}

```

Test:

```bash
# Create a BODY variable containing the JSON data that we want to send.
$ BODY='{"title":"Moana","year":2016,"runtime":107, "genres":["animation","adventure"]}'
# Use the -d flag to send the contents of the BODY variable as the HTTP request body.
# Note that curl will default to sending a POST request when the -d flag is used.
$ curl -i -d "$BODY" localhost:4000/v1/movies
HTTP/1.1 200 OK
Date: Tue, 06 Apr 2021 17:13:46 GMT
Content-Length: 65
Content-Type: text/plain; charset=utf-8
{Title:Moana Year:2016 Runtime:107 Genres:[animation adventure]}

```

#### Handling Bad Request

```bash
# Send some XML as the request body
$ curl -d '<?xml version="1.0" encoding="UTF-8"?><note><to>Alice</to></note>' localhost:4000/v1/movies
{
"error": "invalid character '\u003c' looking for beginning of value"
}

# Send some malformed JSON (notice the trailing comma)
$ curl -d '{"title": "Moana", }' localhost:4000/v1/movies
{
"error": "invalid character '}' looking for beginning of object key string"
}
# Send a JSON array instead of an object
$ curl -d '["foo", "bar"]' localhost:4000/v1/movies
{
"error": "json: cannot unmarshal array into Go value of type struct { Title string
\"json:\\\"title\\\"\"; Year int32 \"json:\\\"year\\\"\"; Runtime int32 \"json:\\
\"runtime\\\"\"; Genres []string \"json:\\\"genres\\\"\" }"
}

```

In all these cases, we can see that our createMovieHandler is doing the right thing. When it receives an invalid request that can’t be decoded into our input struct, no further processing takes place and the client is sent a JSON response containing the error message returned by the Decode() method.

But for a public-facing API, the error messages themselves aren’t ideal.

#### Triaging the Decode error

At this point in our application build, the Decode() method could potentially return the
following five types of error:
Error types Reason

- json.SyntaxError
- io.ErrUnexpectedEOF
  There is a syntax problem with the JSON being decoded.
- json.UnmarshalTypeError A JSON value is not appropriate for the destination Go type.
- json.InvalidUnmarshalError The decode destination is not valid (usually because it is not a
  pointer). This is actually a problem with our application code,
  not the JSON itself.
- io.EOF The JSON being decoded is empty.

#### To fix this:

```go


func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
// Decode the request body into the target destination.
err := json.NewDecoder(r.Body).Decode(dst)
if err != nil {
// If there is an error during decoding, start the triage...
var syntaxError *json.SyntaxError
var unmarshalTypeError *json.UnmarshalTypeError
var invalidUnmarshalError *json.InvalidUnmarshalError
switch {
// Use the errors.As() function to check whether the error has the type
// *json.SyntaxError. If it does, then return a plain-english error message
// which includes the location of the problem.
case errors.As(err, &syntaxError):
return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
// for syntax errors in the JSON. So we check for this using errors.Is() and
// return a generic error message. There is an open issue regarding this at
// https://github.com/golang/go/issues/25956.
case errors.Is(err, io.ErrUnexpectedEOF):
return errors.New("body contains badly-formed JSON")
// Likewise, catch any \*json.UnmarshalTypeError errors. These occur when the
// JSON value is the wrong type for the target destination. If the error relates
// to a specific field, then we include that in our error message to make it
// easier for the client to debug.
case errors.As(err, &unmarshalTypeError):
if unmarshalTypeError.Field != "" {
return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
}
return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
// An io.EOF error will be returned by Decode() if the request body is empty. We
// check for this with errors.Is() and return a plain-english error message
// instead.
case errors.Is(err, io.EOF):
return errors.New("body must not be empty")
// A json.InvalidUnmarshalError error will be returned if we pass something
// that is not a non-nil pointer to Decode(). We catch this and panic,
// rather than returning an error to our handler. At the end of this chapter
// we'll talk about panicking versus returning errors, and discuss why it's an
// appropriate thing to do in this specific situation.
case errors.As(err, &invalidUnmarshalError):
panic(err)
// For anything else, return the error message as-is.
default:
return err
}
}
return nil
}
```

Then in `File: cmd/api/movies.go`

```go
package main
import (
"fmt"
"net/http"
"time"
"greenlight.alexedwards.net/internal/data"
)
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string `json:"title"`
Year int32 `json:"year"`
Runtime int32 `json:"runtime"`
Genres []string `json:"genres"`
}
// Use the new readJSON() helper to decode the request body into the input struct.
// If this returns an error we send the client the error message along with a 400
// Bad Request status code, just like before.
err := app.readJSON(w, r, &input)
if err != nil {
app.errorResponse(w, r, http.StatusBadRequest, err.Error())
return
}
fmt.Fprintf(w, "%+v\n", input)
}

```

Then retry:

```bash
# Send some XML as the request body
$ curl -d '<?xml version="1.0" encoding="UTF-8"?><note><to>Alex</to></note>' localhost:4000/v1/movies
{
"error": "body contains badly-formed JSON (at character 1)"
}
# Send some malformed JSON (notice the trailing comma)
$ curl -d '{"title": "Moana", }' localhost:4000/v1/movies
{
"error": "body contains badly-formed JSON (at character 20)"
}
# Send a JSON array instead of an object
$ curl -d '["foo", "bar"]' localhost:4000/v1/movies
{
"error": "body contains incorrect JSON type (at character 1)"
}
# Send a numeric 'title' value (instead of string)
$ curl -d '{"title": 123}' localhost:4000/v1/movies
{
"error": "body contains incorrect JSON type for \"title\""
}
# Send an empty request body
$ curl -X POST localhost:4000/v1/movies
{
"error": "body must not be empty"
}

```

##### Bad Request helper

```go
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

```

Then in movie handler

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string `json:"title"`
Year int32 `json:"year"`
Runtime int32 `json:"runtime"`
Genres []string `json:"genres"`
}
err := app.readJSON(w, r, &input)
if err != nil {
// Use the new badRequestResponse() helper.
app.badRequestResponse(w, r, err)
return
}
fmt.Fprintf(w, "%+v\n", input)
}

```

#### Restrict Unknown Fields dealing with unknown fields.

One such thing is dealing with unknown fields. For example, you can try sending a request containing the unknown field rating to our createMovieHandler, like so:

```bash
$ curl -i -d '{"title": "Moana", "rating":"PG"}' localhost:4000/v1/movies
HTTP/1.1 200 OK
Date: Tue, 06 Apr 2021 18:51:50 GMT
Content-Length: 41
Content-Type: text/plain; charset=utf-8
{Title:Moana Year:0 Runtime:0 Genres:[]}


```

##### Unknown request

```go

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1,048,576
	// bytes (1MB).
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// Decode the request body to the destination.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		// Add a new maxBytesError variable.
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message. Note that there's an open
		// issue at https://github.com/golang/go/issues/29035 regarding turning this
		// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		// Use the errors.As() function to check whether the error has the type
		// *http.MaxBytesError. If it does, then it means the request body exceeded our
		// size limit of 1MB and we return a clear error message.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

```

Then test:

```bash
$ curl -d '{"title": "Moana", "rating":"PG"}' localhost:4000/v1/movies
{
"error": "body contains unknown key \"rating\""
}
$ curl -d '{"title": "Moana"}{"title": "Top Gun"}' localhost:4000/v1/movies
{
"error": "body must only contain a single JSON value"
}
$ curl -d '{"title": "Moana"} :~()' localhost:4000/v1/movies
{
"error": "body must only contain a single JSON value"
}

```

##### The json.Unmarshaler interface

```go
File: internal/data/runtime.go
package data
import (
"errors" // New import
"fmt"
"strconv"
"strings" // New import
)
// Define an error that our UnmarshalJSON() method can return if we're unable to parse
// or convert the JSON string successfully.
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")
type Runtime int32
...
// Implement a UnmarshalJSON() method on the Runtime type so that it satisfies the
// json.Unmarshaler interface. IMPORTANT: Because UnmarshalJSON() needs to modify the
// receiver (our Runtime type), we must use a pointer receiver for this to work
// correctly. Otherwise, we will only be modifying a copy (which is then discarded when
// this method returns).
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
// We expect that the incoming JSON value will be a string in the format
// "<runtime> mins", and the first thing we need to do is remove the surrounding
// double-quotes from this string. If we can't unquote it, then we return the
// ErrInvalidRuntimeFormat error.
unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
if err != nil {
return ErrInvalidRuntimeFormat
}
// Split the string to isolate the part containing the number.
parts := strings.Split(unquotedJSONValue, " ")
// Sanity check the parts of the string to make sure it was in the expected format.
// If it isn't, we return the ErrInvalidRuntimeFormat error again.
if len(parts) != 2 || parts[1] != "mins" {
return ErrInvalidRuntimeFormat
}
// Otherwise, parse the string containing the number into an int32. Again, if this
// fails return the ErrInvalidRuntimeFormat error.
i, err := strconv.ParseInt(parts[0], 10, 32)
if err != nil {
return ErrInvalidRuntimeFormat
}
// Convert the int32 to a Runtime type and assign this to the receiver. Note that we
// use the * operator to deference the receiver (which is a pointer to a Runtime
// type) in order to set the underlying value of the pointer.
*r = Runtime(i)
return nil
}
```

so,

```bash
$ curl -d '{"title": "Moana", "runtime": "107 mins"}' localhost:4000/v1/movies
{Title:Moana Year:0 Runtime:107 Genres:[]}
```

and

```bash
$ curl -d '{"title": "Moana", "runtime": 107}' localhost:4000/v1/movies
{
"error": "invalid runtime format"
}
$ curl -d '{"title": "Moana", "runtime": "107 minutes"}' localhost:4000/v1/movies
{
"error": "invalid runtime format"
}
```

#### Validating JSON Input

e.g The movie title provided by the client is not empty and is not more than 500 bytes long.
The movie year is not empty and is between 1888 and the current year.
The movie runtime is not empty and is a positive integer.
The movie has between one and five (unique) genres.

```bash
$ mkdir internal/validator
$ touch internal/validator/validator.go
```

```go
package validator

import (
	"regexp"
	"slices"
)

// EmailRX is a regular expression for sanity checking the format of email addresses.
// Taken from https://html.spec.whatwg.org/#valid-e-mail-address
var EmailRX = regexp.MustCompile(
	`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+` + // local part
		`@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?` + // domain start
		`(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`, // domain dots
)

// Validator contains validation errors.
type Validator struct {
	Errors map[string]string
}

// New creates a new Validator instance with an empty errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (if not already present).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message only if a validation check fails.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue returns true if a specific value is in a list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]struct{})
	for _, value := range values {
		uniqueValues[value] = struct{}{}
	}
	return len(values) == len(uniqueValues)
}
```

Use it

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string `json:"title"`
Year int32 `json:"year"`
Runtime data.Runtime `json:"runtime"`
Genres []string `json:"genres"`
}
err := app.readJSON(w, r, &input)
if err != nil {
app.badRequestResponse(w, r, err)
return
}
// Initialize a new Validator instance.
v := validator.New()
// Use the Check() method to execute our validation checks. This will add the
// provided key and error message to the errors map if the check does not evaluate
// to true. For example, in the first line here we "check that the title is not
// equal to the empty string". In the second, we "check that the length of the title
// is less than or equal to 500 bytes" and so on.
v.Check(input.Title != "", "title", "must be provided")
v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")
v.Check(input.Year != 0, "year", "must be provided")
v.Check(input.Year >= 1888, "year", "must be greater than 1888")
v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")
v.Check(input.Runtime != 0, "runtime", "must be provided")
v.Check(input.Runtime > 0, "runtime", "must be a positive integer")
v.Check(input.Genres != nil, "genres", "must be provided")
v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
// Note that we're using the Unique helper in the line below to check that all
// values in the input.Genres slice are unique.
v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
// Use the Valid() method to see if any of the checks failed. If they did, then use
// the failedValidationResponse() helper to send a response to the client, passing
// in the v.Errors map.
if !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
fmt.Fprintf(w, "%+v\n", input)
}
```

With that done, we should be good to try this out. Restart the API, then go ahead and issue a
request to the POST /v1/movies endpoint containing some invalid data. Similar to this:

```bash
$ BODY='{"title":"","year":1000,"runtime":"-123 mins","genres":["sci-fi","sci-fi"]}'
$ curl -i -d "$BODY" localhost:4000/v1/movies
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json
Date: Wed, 07 Apr 2021 10:33:57 GMT
Content-Length: 180
{
"error": {
"genres": "must not contain duplicate values",
"runtime": "must be a positive integer",
"title": "must be provided",
"year": "must be greater than 1888"
}
}
```

#### Making validation rules reusable

If you’re following along, reopen the internal/data/movies.go file and add a
ValidateMovie() function containing the checks like so:

```go

package data

import (
	"time"

	"github.com/drsimplegraffiti/greenlight/internal/validator"
)

// Ensure the fields are exportable i.e in Capital letter
// type Movie struct {
// 	ID        int64     `json:"id"`
// 	CreatedAt time.Time `json:"-"` // Use the - directive
// 	Title     string    `json:"title"`
// 	Year      int32     `json:"year,omitzero"`    // Add the omitzero directive
// 	Runtime   int32     `json:"runtime,omitzero"` // Add the omitzero directive
// 	Genres    []string  `json:"genres,omitzero"`  // Add the omitzero directive
// 	Version   int32     `json:"version"`
// }

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	// Use the Runtime type instead of int32. Note that the omitzero directive will
	// still work on this: if the Runtime field has the underlying value 0, then it will
	// be considered zero and omitted -- and the MarshalJSON() method we just made
	// won't be called at all.
	Runtime Runtime  `json:"runtime,omitzero"`
	Genres  []string `json:"genres,omitzero"`
	Version int32    `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

```

Then in movie handler

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string `json:"title"`
Year int32 `json:"year"`
Runtime data.Runtime `json:"runtime"`
Genres []string `json:"genres"`
}
err := app.readJSON(w, r, &input)
if err != nil {
app.badRequestResponse(w, r, err)
return
}
// Copy the values from the input struct to a new Movie struct.
movie := &data.Movie{
Title: input.Title,
Year: input.Year,
Runtime: input.Runtime,
Genres: input.Genres,
}
// Initialize a new Validator.
v := validator.New()
// Call the ValidateMovie() function and return a response containing the errors if
// any of the checks fail.
if data.ValidateMovie(v, movie); !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
fmt.Fprintf(w, "%+v\n", input)
}
```

#### Database Connection

```bash
psql --version
psql (PostgreSQL) 16.9 (Ubuntu 16.9-0ubuntu0.24.04.1)
# Check postgres user
cat /etc/passwd | grep 'postgres'
postgres:x:105:109:PostgreSQL administrator,,,:/var/lib/postgresql:/bin/bash
# connect
sudo -u postgres psql

# check current user
SELECT current_user;

```

![image10](assets/image10.png)

##### Create Database

```bash
 CREATE DATABASE greenlight;
  \c greenlight
 CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
  CREATE EXTENSION IF NOT EXISTS citext;
```

In this project we’re going to use the citext extension. This adds a case-insensitive character string type to PostgreSQL, which we will use later in the book to store user email addresses.
Go ahead and run the following commands to create a new greenlight user with a specific
password and add the citext extension to our database.

To ensure our password and username works well

```bash
➜  gox psql --host=localhost --dbname=greenlight --username=greenlight
Password for user greenlight:
psql (16.9 (Ubuntu 16.9-0ubuntu0.24.04.1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, compression: off)
Type "help" for help.

greenlight=>
```

#### PG driver

```bash
go get github.com/lib/pq@v1
```

Connection string `postgres://greenlight:pa55word@localhost/greenlight`

```go
package main

import (
	"context"      // New import
	"database/sql" // New import
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

const version = "1.0.0"

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:pa55word@localhost/greenlight", "PostgreSQL DSN")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()
	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Info("database connection pool established")
	app := &application{
		config: cfg,
		logger: logger,
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	// Because the err variable is now already declared in the code above, we need
	// to use the = operator here, instead of the := operator.
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error. If we get this error, or any other, we close the connection pool and
	// return the error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}

```

Note: If you receive the error message pq: SSL is not enabled on the server you should set your DSN to:
postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable.

If you’re following along, create a new GREENLIGHT_DB_DSN environment variable by adding
the following line to either your $HOME/.profile or $HOME/.bashrc files:
File: $HOME/.profile

```bash
export GREENLIGHT_DB_DSN='postgres://greenlight:pa55word@localhost/greenlight'
source $HOME/.profile
export GREENLIGHT_DB_DSN='postgres://xyborg:Bassguitar1@localhost:5432/greenlight'
source $HOME/.profile

```

Confirm the settings:

```bash
echo $GREENLIGHT_DB_DSN
# You should get: postgres://greenlight:pa55word@localhost/greenlight
```

#### set using os.Getenv

```bash
// Use the value of the GREENLIGHT_DB_DSN environment variable as the default value
// for our db-dsn command-line flag.
flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
```

You can also try specifying the -help flag when running the application. This should output the descriptive text and default values for our three command-line flags, including the DSN value pulled through from the environment variable. Similar to this:

```bash
$ go run ./cmd/api -help
Usage of /tmp/go-build417842398/b001/exe/api:
-db-dsn string
PostgreSQL DSN (default "postgres://greenlight:pa55word@localhost/greenlight")
-env string
Environment (development|staging|production) (default "development")
-port int
API server port (default 4000)

```

##### Additional Information

Using the DSN with psql
A nice side effect of storing the DSN in an environment variable is that you can use it to
easily connect to the greenlight database as the greenlight user, rather than specifying all
the connection options manually when running psql. Like so:

```bash

$ psql $GREENLIGHT_DB_DSN
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=>
```

#### Add db config

```go
package main

import (
	"context"      // New import
	"database/sql" // New import
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

const version = "1.0.0"

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	// flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:pa55word@localhost/greenlight", "PostgreSQL DSN")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice that the default values we're using are the ones we discussed above?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()
	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Info("database connection pool established")
	app := &application{
		config: cfg,
		logger: logger,
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	// Because the err variable is now already declared in the code above, we need
	// to use the = operator here, instead of the := operator.
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

```

Run the app:

```bash
$ go run ./cmd/api #or
$ go run ./cmd/api -db-max-open-conns=50 -db-max-idle-conns=50 -db-max-idle-time=2h30m
```

#### Golang migrate

```bash
brew install golang-migrate #mac


 cd /tmp
 curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
 mv migrate ~/go/bin/

```

Then check migrate

```bash
 migrate -version
 4.16.2
```

#### Migration

```bash
 migrate create -seq -ext=.sql -dir=./migrations create_movies_table
 #you get
 /mnt/c/Users/CAPTAIN ENJOY/Downloads/gox/migrations/000001_create_movies_table.up.sql
/mnt/c/Users/CAPTAIN ENJOY/Downloads/gox/migrations/000001_create_movies_table.down.sql

```

In this command:
The -seq flag indicates that we want to use sequential numbering like 0001, 0002, ...
for the migration files (instead of a Unix timestamp, which is the default).
The -ext flag indicates that we want to give the migration files the extension .sql.
The -dir flag indicates that we want to store the migration files in the ./migrations
directory (which will be created automatically if it doesn’t already exist).
The name create_movies_table is a descriptive label that we give the migration files to
signify their contents.
If you look in your migrations directory, you should now see a pair of new ‘up’ and ‘down’
migration files like so:
./migrations/
├── 000001_create_movies_table.down.sql
└── 000001_create_movies_table.up.sql

At the moment these two new files are completely empty. Let’s edit the ‘up’ migration file to contain the necessary CREATE TABLE statement for our movies table, like so:
In the up `File: migrations/000001_create_movies_table.up.sql`

```sql
CREATE TABLE IF NOT EXISTS movies (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
title text NOT NULL,
year integer NOT NULL,
runtime integer NOT NULL,
genres text[] NOT NULL,
version integer NOT NULL DEFAULT 1
);

```

```bash
migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints
```

#### Executing the migrations

```bash
 migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
 # you should get:
1/u create_movies_table (26.298722ms)
2/u add_movies_check_constraints (36.597614ms)

```

#### To confirm

```bash
$ psql $GREENLIGHT_DB_DSN
Password for user greenlight:
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=> \dt
List of relations
Schema | Name | Type | Owner
--------+-------------------+-------+------------
public | movies | table | greenlight
public | schema_migrations | table | greenlight
(2 rows)

```

![image12](assets/image12.png)

#### Check the movies table

```sql
greenlight=# \d movies
                                        Table "public.movies"
   Column   |            Type             | Collation | Nullable |              Default
------------+-----------------------------+-----------+----------+------------------------------------
 id         | bigint                      |           | not null | nextval('movies_id_seq'::regclass)
 created_at | timestamp(0) with time zone |           | not null | now()
 title      | text                        |           | not null |
 year       | integer                     |           | not null |
 runtime    | integer                     |           | not null |
 genres     | text[]                      |           | not null |
 version    | integer                     |           | not null | 1
Indexes:
    "movies_pkey" PRIMARY KEY, btree (id)
Check constraints:
    "genres_length_check" CHECK (array_length(genres, 1) >= 1 AND array_length(genres, 1) <= 5)
    "movies_runtime_check" CHECK (runtime >= 0)
    "movies_year_check" CHECK (year >= 1888 AND year::double precision <= date_part('year'::text, now()))

greenlight=#
```

#### Migrating to a specific version

```bash
see which
# migration version your database is currently on
$ migrate -path=./migrations -database=$EXAMPLE_DSN version
```

You can also migrate up or down to a specific version by using the goto command:

```bash
$ migrate -path=./migrations -database=$EXAMPLE_DSN goto 1
$ migrate -path=./migrations -database =$EXAMPLE_DSN down 1
$ migrate -path=./migrations -database=$EXAMPLE_DSN force 1
```

#### Add db access

```go
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Note that the movie variable contains a *pointer* to a Movie struct.
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated movie struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

```

Call it:

```bash
$ BODY='{"title":"Moana","year":2016,"runtime":"107 mins", "genres":["animation","adventure"]}'
$ curl -i -d "$BODY" localhost:4000/v1/movies
HTTP/1.1 201 Created
Content-Type: application/json
Location: /v1/movies/1
Date: Wed, 07 Apr 2021 19:21:41 GMT
Content-Length: 156
{
"movie": {
"id": 1,
"title": "Moana",
"year": 2016,
"runtime": "107 mins",
"genres": [
"animation",
"adventure"
],
"version": 1
}
}
```

#### Check in DB

```bash
$ psql $GREENLIGHT_DB_DSN
Password for user greenlight:
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=> SELECT * FROM movies;
id | created_at | title | year | runtime | genres | version
----+------------------------+--------------------+------+---------+-----------------------+---------
1 | 2021-04-07 21:21:41+02 | Moana | 2016 | 107 | {animation,adventure} | 1
2 | 2021-04-07 21:28:28+02 | Black Panther | 2018 | 134 | {action,adventure} | 1
3 | 2021-04-07 21:28:36+02 | Deadpool | 2016 | 108 | {action,comedy} | 1
4 | 2021-04-07 21:28:44+02 | The Breakfast Club | 1986 | 96 | {drama} | 1
(4 rows)
```

#### Get single movie

```go

func (m MovieModel) Insert(movie *Movie) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

```

```bash
➜  gox curl -i localhost:4000/v1/movies/42
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Fri, 29 Aug 2025 09:54:16 GMT
Content-Length: 58

{
        "error": "the requested resource could not be found"
}
➜  gox curl -i localhost:4000/v1/movies/4
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 29 Aug 2025 09:54:19 GMT
Content-Length: 148

{
        "movie": {
                "id": 4,
                "title": "The Breakfast Club",
                "year": 1986,
                "runtime": "96 mins",
                "genres": [
                        "drama"
                ],
                "version": 1
        }
}

```

#### Update

```bash
$ BODY='{"title":"Black Panther","year":2018,"runtime":"134 mins","genres":["sci-fi","action","adventure"]}'
$ curl -X PUT -d "$BODY" localhost:4000/v1/movies/2
{
"movie": {
"id": 2,
"title": "Black Panther",
"year": 2018,
"runtime": "134 mins",
"genres": [
"sci-fi",
"action",
"adventure"
],
"version": 2
}
}

```

```go
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Declare an input struct to hold the expected data from the c
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the movie
	// record.
	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres
	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated movie record to our new Update() method.
	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the updated movie record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

```

#### Delete

```bash
curl -X DELETE localhost:4000/v1/movies/3
```

```go

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the movie from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
```

```go

func (m MovieModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM movies
		WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as
	// the value for the placeholder parameter. The Exec() method returns a sql.Result
	// object.
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
```

#### Partial Update

```go

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Retrieve the movie record as normal.
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// If the input.Title value is nil then we know that no corresponding "title" key/
	// value pair was provided in the JSON request body. So we move on and leave the
	// movie record unchanged. Otherwise, we update the movie record with the new title
	// value. Importantly, because input.Title is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our movie record.
	if input.Title != nil {
		movie.Title = *input.Title
	}
	// We also do the same for the other fields in the input struct.
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres // Note that we don't need to dereference a slice.
	}
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

```

#### Test it

```bash
$ curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/movies/4
{
"movie": {
"id": 4,
"title": "The Breakfast Club",
"year": 1985,
"runtime": "96 mins",
"genres": [
"drama"
],
"version": 2
}
}
```

#### Race condition

```bash
 xargs -I % -P8 curl -X PATCH -d '{"runtime": "97 mins"}' "localhost:4000/v1/movies/4" < <(printf '%s\n' {1..8})
```

You get :

```bash
 xargs -I % -P8 curl -X PATCH -d '{"runtime": "97 mins"}' "localhost:4000/v1/movies/4" < <(printf '%s\n' {1..8})
{
        "movie": {
{
{
                "id": 4,
        "error": "unable to update the record due to an edit conflict, please try again"
        "error": "unable to update the record due to an edit conflict, please try again"
                "title": "The Breakfast Club",
}
}
                "year": 1985,
                "runtime": "97 mins",
                "genres": [
                        "drama"
                ],
                "version": 5
        }
}
{
        "error": "unable to update the record due to an edit conflict, please try again"
}
{
        "error": "unable to update the record due to an edit conflict, please try again"
}
```

```go

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Retrieve the movie record as normal.
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// If the input.Title value is nil then we know that no corresponding "title" key/
	// value pair was provided in the JSON request body. So we move on and leave the
	// movie record unchanged. Otherwise, we update the movie record with the new title
	// value. Importantly, because input.Title is a now a pointer to a string, we need
	// to dereference the pointer using the * operator to get the underlying value
	// before assigning it to our movie record.
	if input.Title != nil {
		movie.Title = *input.Title
	}
	// We also do the same for the other fields in the input struct.
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres // Note that we don't need to dereference a slice.
	}
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

```

##### Simulating a long running request

```go
func (m MovieModel) Get(id int64) (*Movie, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts
	// auto-incrementing at 1 by default, so we know that no movies will have ID values
	// less than that. To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	// Update the query to return pg_sleep(8) as the first value.
	query := `
		SELECT pg_sleep(8), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1`
	// query := `
	// 	SELECT id, created_at, title, year, runtime, genres, version
	// 	FROM movies
	// 	WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var movie Movie
	// Execute the query using the QueryRow() method, passing in the provided id valu
	// Execute the query using the QueryRow() method, passing in the provided id value
	// as a placeholder parameter, and scan the response data into the fields of the
	// Movie struct. Importantly, notice that we need to convert the scan target for the
	// genres column using the pq.Array() adapter function again.
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &movie, nil
}

```

Then test:

```bash
$ curl -w '\nTime: %{time_total}s \n' localhost:4000/v1/movies/1
{
"movie": {
"id": 1,
"title": "Moana",
"year": 2015,
"runtime": "107 mins",
"genres": [
"animation",
"adventure"
],
"version": 1
}
}
Time: 8.013534s
```

Adding a query timeout
Now that we’ve got some code that mimics a long-running query, let’s enforce a timeout so
that the SQL query is automatically canceled if it doesn’t complete within 3 seconds.
To do this we need to:

1. Use the context.WithTimeout() function to create a context.Context instance with a 3-
   second timeout deadline.
2. Execute the SQL query using the QueryRowContext() method, passing the
   context.Context instance as a parameter.

```go
func (m MovieModel) Get(id int64) (*Movie, error) {
if id < 1 {
return nil, ErrRecordNotFound
}
query := `
SELECT pg_sleep(8), id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1`
var movie Movie
// Use the context.WithTimeout() function to create a context.Context which carries a
// 3-second timeout deadline. Note that we're using the empty context.Background()
// as the 'parent' context.
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// Importantly, use defer to make sure that we cancel the context before the Get()
// method returns.
defer cancel()
// Use the QueryRowContext() method to execute the query, passing in the context
// with the deadline as the first argument.
err := m.DB.QueryRowContext(ctx, query, id).Scan(
&[]byte{},
&movie.ID,
&movie.CreatedAt,
&movie.Title,
&movie.Year,
&movie.Runtime,
pq.Array(&movie.Genres),
&movie.Version,
)
if err != nil {
switch {
case errors.Is(err, sql.ErrNoRows):
return nil, ErrRecordNotFound
default:
return nil, err
}
}
return &movie, nil
}
```

```bash
➜  gox curl -w '\nTime: %{time_total}s \n' localhost:4000/v1/movies/1
{
        "error": "the server encountered a problem and could not process your request"
}

Time: 3.014782s
```

And you will see this in the terminal running the app

```bash
➜  gox go run ./cmd/api
time=2025-08-29T17:27:12.442+01:00 level=INFO msg="database connection pool established"
time=2025-08-29T17:27:12.442+01:00 level=INFO msg="starting server" addr=:4000 env=development
time=2025-08-29T17:27:19.255+01:00 level=ERROR msg="pq: canceling statement due to user request" method=GET uri=/v1/movies/
```

#### Adding pagination and sorting

```go
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
// To keep things consistent with our other handlers, we'll define an input struct
// to hold the expected values from the request query string.
var input struct {
Title string
Genres []string
Page int
PageSize int
Sort string
}
// Initialize a new Validator instance.
v := validator.New()
// Call r.URL.Query() to get the url.Values map containing the query string data.
qs := r.URL.Query()
// Use our helpers to extract the title and genres query string values, falling back
// to defaults of an empty string and an empty slice respectively if they are not
// provided by the client.
input.Title = app.readString(qs, "title", "")
input.Genres = app.readCSV(qs, "genres", []string{})
// Get the page and page_size query string values as integers. Notice that we set
// the default page value to 1 and default page_size to 20, and that we pass the
// validator instance as the final argument here.
input.Page = app.readInt(qs, "page", 1, v)
input.PageSize = app.readInt(qs, "page_size", 20, v)
// Extract the sort query string value, falling back to "id" if it is not provided
// by the client (which will imply a ascending sort on movie ID).
input.Sort = app.readString(qs, "sort", "id")
// Check the Validator instance for any errors and use the failedValidationResponse()
// helper to send the client a response if necessary.
if !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
// Dump the contents of the input struct in a HTTP response.
fmt.Fprintf(w, "%+v\n", input)
}
```

#### Test it

```bash
curl "localhost:4000/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=year"
{Title:godfather Genres:[crime drama] Page:1 PageSize:5 Sort:year}
```

##### Creating a filter struct

```bash
touch internal/data/filters.go
```

```go
type Filters struct {
Page int
PageSize int
Sort string
}
```

Then

```go
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
// Embed the new Filters struct.
var input struct {
Title string
Genres []string
data.Filters
}
v := validator.New()
qs := r.URL.Query()
input.Title = app.readString(qs, "title", "")
input.Genres = app.readCSV(qs, "genres", []string{})
// Read the page and page_size query string values into the embedded struct.
input.Filters.Page = app.readInt(qs, "page", 1, v)
input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
// Read the sort query string value into the embedded struct.
input.Filters.Sort = app.readString(qs, "sort", "id")
if !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
fmt.Fprintf(w, "%+v\n", input)
}
```

##### Test it

```bash
`curl "localhost:4000/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=year"
{Title:godfather Genres:[crime drama] Filters:{Page:1 PageSize:5 Sort:year}}``
```

#### Validate Query String parameters

```go
File: internal/data/filters.go
package data
import (
"greenlight.alexedwards.net/internal/validator" // New import
)
// Add a SortSafelist field to hold the supported sort values.
type Filters struct {
Page int
PageSize int
Sort string
SortSafelist []string
}
func ValidateFilters(v *validator.Validator, f Filters) {
// Check that the page and page_size parameters contain sensible values.
v.Check(f.Page > 0, "page", "must be greater than zero")
v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
// Check that the sort parameter matches a value in the safelist.
v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
```

Then

```go
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string
Genres []string
data.Filters
}
v := validator.New()
qs := r.URL.Query()
input.Title = app.readString(qs, "title", "")
input.Genres = app.readCSV(qs, "genres", []string{})
input.Filters.Page = app.readInt(qs, "page", 1, v)
input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
input.Filters.Sort = app.readString(qs, "sort", "id")
// Add the supported sort values for this endpoint to the sort safelist.
input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
// Execute the validation checks on the Filters struct and send a response
// containing the errors if necessary.
if data.ValidateFilters(v, input.Filters); !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
fmt.Fprintf(w, "%+v\n", input)
}
```

If you restart the API and try making a request with some invalid page, page_size and sort
parameters, you should now receive an error response containing the relevant validation
failure messages. Similar to this:

```bash
$ curl "localhost:4000/v1/movies?page=-1&page_size=-1&sort=foo"
{
"error": {
"page": "must be greater than zero",
"page_size": "must be greater than zero",
"sort": "invalid sort value"
}
}
```

#### Getting multiple records from the database

```go
// Create a new GetAll() method which returns a slice of movies. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
	// Construct the SQL query to retrieve all movie records.
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		ORDER BY id`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	movies := []*Movie{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var movie Movie
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}
		// Add the Movie struct to the slice.
		movies = append(movies, &movie)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the slice of movies.
	return movies, nil
}
```

Then the listMoviesHandler

```go
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string
Genres []string
data.Filters
}
v := validator.New()
qs := r.URL.Query()
input.Title = app.readString(qs, "title", "")
input.Genres = app.readCSV(qs, "genres", []string{})
input.Filters.Page = app.readInt(qs, "page", 1, v)
input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
input.Filters.Sort = app.readString(qs, "sort", "id")
input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
if data.ValidateFilters(v, input.Filters); !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
// Call the GetAll() method to retrieve the movies, passing in the various filter
// parameters.
movies, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
if err != nil {
app.serverErrorResponse(w, r, err)
return
}
// Send a JSON response containing the movie data.
err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
if err != nil {
app.serverErrorResponse(w, r, err)
}
}
```

#### Test it

```bash
$ curl localhost:4000/v1/movies
{
"movies": [
	{
	"id": 1,
	"title": "Moana",
	"year": 2016,
	"runtime": "107 mins",
	"genres": [
	"animation",
	"adventure"
	],
	"version": 1
	},
	{
	"id": 2,
	"title": "Black Panther",
	"year": 2018,
	"runtime": "134 mins",
	"genres": [
	"sci-fi",
	"action",
	"adventure"
	],
	"version": 2
	},
	{
	"id": 4,
	"title": "The Breakfast Club",
	"year": 1985,
	"runtime": "97 mins",
	"genres": [
	"comedy"
	],
	"version": 5
	}
	]
}
```

##### Dynamic filtering in the SQL query

```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
	// Update the SQL query to include the filter conditions.
	query := `
			SELECT id, created_at, title, year, runtime, genres, version
			FROM movies
			WHERE (LOWER(title) = LOWER($1) OR $1 = '')
			AND (genres @> $2 OR $2 = '{}')
			ORDER BY id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Pass the title and genres as the placeholder parameter values.
	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

```

```bash
$ curl "localhost:4000/v1/movies?title=black+panther"
{
"movies": [
{
"id": 2,
"title": "Black Panther",
"year": 2018,
"runtime": "134 mins",
"genres": [
"sci-fi",
"action",
"adventure"
],
"version": 2
}
]
```

```bash
curl "localhost:4000/v1/movies?title=black+panther"
curl "localhost:4000/v1/movies?genres=adventure"
```

#### Full-Text Search

```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
	// Update the SQL query to include the filter conditions.
	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY id`
	// query := `
	// 		SELECT id, created_at, title, year, runtime, genres, version
	// 		FROM movies
	// 		WHERE (LOWER(title) = LOWER($1) OR $1 = '')
	// 		AND (genres @> $2 OR $2 = '{}')
	// 		ORDER BY id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Pass the title and genres as the placeholder parameter values.
	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

```

#### Test it

```bash
$ curl "localhost:4000/v1/movies?title=panther"
{
"movies": [
{
"id": 2,
"title": "Black Panther",
"year": 2018,
"runtime": "134 mins",
"genres": [
"sci-fi",
"action",
"adventure"
],
"version": 2
}
]
}
$ curl "localhost:4000/v1/movies?title=the+club"
{
"movies": [
{
"id": 4,
"title": "The Breakfast Club",
"year": 1985,
"runtime": "97 mins",
"genres": [
"comedy"
],
"version": 5
}
]
}
```

##### Add indexes

```bash
$ migrate create -seq -ext .sql -dir ./migrations add_movies_indexes
```

File: migrations/000003_add_movies_indexes.up.sql
CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);
File: migrations/000003_add_movies_indexes.down.sql
DROP INDEX IF EXISTS movies_title_idx;
DROP INDEX IF EXISTS movies_genres_idx;

Then run:

```bash
$ migrate -path ./migrations -database $GREENLIGHT_DB_DSN up
3/u add_movies_indexes (38.638777ms)
```

#### Sorting and sorting direction

```go
package data
import (
"strings" // New import
"greenlight.alexedwards.net/internal/validator"
)
type Filters struct {
Page int
PageSize int
Sort string
SortSafelist []string
}
// Check that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
for _, safeValue := range f.SortSafelist {
if f.Sort == safeValue {
return strings.TrimPrefix(f.Sort, "-")
```

#### Then update the GetAll() method

```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
// Add an ORDER BY clause and interpolate the sort column and direction. Importantly
// notice that we also include a secondary sort on the movie ID to ensure a
// consistent ordering.
query := fmt.Sprintf(`
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())
// Nothing else below needs to change.
...
}
```

##### Try it

```bash
$ curl "localhost:4000/v1/movies?sort=-title"
{
"movies": [
{
"id": 4,
"title": "The Breakfast Club",
"year": 1985,
"runtime": "97 mins",
"genres": [
"comedy"
],
"version": 5
},
{
"id": 1,
"title": "Moana",
"year": 2016,
"runtime": "107 mins",
"genres": [
"animation",
"adventure"
],
"version": 1
},
{
"id": 2,
"title": "Black Panther",
"year": 2018,
"runtime": "134 mins",
"genres": [
"sci-fi",
"action",
"adventure"
],
"version": 2
}
]
}
```

#### Sort with runtime

```bash
curl "localhost:4000/v1/movies?sort=-runtime"
```

#### Paging list

```go
func (f Filters) limit() int {
return f.PageSize
}
func (f Filters) offset() int {
return (f.Page - 1) * f.PageSize
}
```

Then update the GetAll() method

```go
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
// Update the SQL query to include the LIMIT and OFFSET clauses with placeholder
// parameter values.
query := fmt.Sprintf(`
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
// As our SQL query now has quite a few placeholder parameters, let's collect the
// values for the placeholders in a slice. Notice here how we call the limit() and
// offset() methods on the Filters struct to get the appropriate values for the
// LIMIT and OFFSET clauses.
args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
// And then pass the args slice to QueryContext() as a variadic parameter.
rows, err := m.DB.QueryContext(ctx, query, args...)
if err != nil {
return nil, err
}
```

Test:

```bash
curl "localhost:4000/v1/movies?page_size=2"
curl "localhost:4000/v1/movies?page_size=2&page=2"
```

#### Returning Pagination Metadata

```go
File: internal/data/filters.go
package data
import (
"math" // New import
"strings"
"greenlight.alexedwards.net/internal/validator"
)
...
// Define a new Metadata struct for holding the pagination metadata.
type Metadata struct {
CurrentPage int `json:"current_page,omitzero"`
PageSize int `json:"page_size,omitzero"`
FirstPage int `json:"first_page,omitzero"`
LastPage int `json:"last_page,omitzero"`
TotalRecords int `json:"total_records,omitzero"`
}
// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that when the last page value is calculated we are dividing two int values, and
// when dividing integer types in Go the result will also be an integer type, with
// the modulus (or remainder) dropped. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be (12+5-1)/5 = 3.2, which is then
// truncated to 3 by Go.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
if totalRecords == 0 {
// Note that we return an empty Metadata struct if there are no records.
return Metadata{}
}
return Metadata{
CurrentPage: page,
PageSize: pageSize,
FirstPage: 1,
LastPage: (totalRecords + pageSize - 1) / pageSize,
TotalRecords: totalRecords,
}
}
```

##### Update the GetAll() method

```go
// Update the function signature to return a Metadata struct.
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
// Update the SQL query to include the window function which counts the total
// (filtered) records.
query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}
rows, err := m.DB.QueryContext(ctx, query, args...)
if err != nil {
return nil, Metadata{}, err // Update this to return an empty Metadata struct.
}
defer rows.Close()
// Declare a totalRecords variable.
totalRecords := 0
movies := []*Movie{}
for rows.Next() {
var movie Movie
err := rows.Scan(
&totalRecords, // Scan the count from the window function into totalRecords.
&movie.ID,
&movie.CreatedAt,
&movie.Title,
&movie.Year,
&movie.Runtime,
pq.Array(&movie.Genres),
&movie.Version,
)
if err != nil {
return nil, Metadata{}, err // Update this to return an empty Metadata struct.
}
movies = append(movies, &movie)
}
if err = rows.Err(); err != nil {
return nil, Metadata{}, err // Update this to return an empty Metadata struct.
}
// Generate a Metadata struct, passing in the total record count and pagination
// parameters from the client.
metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
// Include the metadata struct when returning.
return movies, metadata, nil
}``
```

#### Update the listMoviesHandler() method

```go
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
var input struct {
Title string
Genres []string
data.Filters
}
v := validator.New()
qs := r.URL.Query()
input.Title = app.readString(qs, "title", "")
input.Genres = app.readCSV(qs, "genres", []string{})
input.Filters.Page = app.readInt(qs, "page", 1, v)
input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
input.Filters.Sort = app.readString(qs, "sort", "id")
input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
if data.ValidateFilters(v, input.Filters); !v.Valid() {
app.failedValidationResponse(w, r, v.Errors)
return
}
// Accept the metadata struct as a return value.
movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
if err != nil {
app.serverErrorResponse(w, r, err)
return
}
// Include the metadata in the response envelope.
err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
if err != nil {
app.serverErrorResponse(w, r, err)
}
}
```

#### Test it

```bash
curl "localhost:4000/v1/movies?page=1&page_size=2"
{
"metadata": {
"current_page": 1,
"page_size": 2,
"first_page": 1,
"last_page": 2,
"total_records": 3
},
"movies": [
{
"id": 1,
"title": "Moana",
"year": 2015,
"runtime": "107 mins",
"genres": [
"animation",
"adventure"
],
"version": 1
},
{
"id": 2,
"title": "Black Panther",
"year": 2018,
"runtime": "134 mins",
"genres": [
"sci-fi",
"action",
"adventure"
],
"version": 2
}
]
}
```

#### Rate Limiting

Global Rate Limiting Middleware

```bash
go get golang.org/x/time/rate@latest
```

```go
func (app *application) rateLimit(next http.Handler) http.Handler {
// Initialize a new rate limiter which allows an average of 2 requests per second,
// with a maximum of 4 requests in a single ‘burst’.
limiter := rate.NewLimiter(2, 4)
// The function we are returning is a closure, which 'closes over' the limiter
// variable.
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Call limiter.Allow() to see if the request is permitted, and if it's not,
// then we call the rateLimitExceededResponse() helper to return a 429 Too Many
// Requests response (we will create this helper in a minute).
if !limiter.Allow() {
app.rateLimitExceededResponse(w, r)
return
}
next.ServeHTTP(w, r)
})
}
```

```go
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
message := "rate limit exceeded"
app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
```

#### Then in the routes

```go
File: cmd/api/routes.go
package main
...
func (app *application) routes() http.Handler {
router := httprouter.New()
router.NotFound = http.HandlerFunc(app.notFoundResponse)
router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
// Wrap the router with the rateLimit() middleware.
return app.recoverPanic(app.rateLimit(router))
}
```

#### Then test

```bash for i in {1..6}; do curl -s http://localhost:4000/v1/healthcheck | jq -c .; done
{"status":"available","system_info":{"environment":"development","version":"1.0.0"}}
{"status":"available","system_info":{"environment":"development","version":"1.0.0"}}
{"status":"available","system_info":{"environment":"development","version":"1.0.0"}}
{"status":"available","system_info":{"environment":"development","version":"1.0.0"}}
{"error":"rate limit exceeded"}
{"error":"rate limit exceeded"}
```

#### Ip rate limiting

```bash
go get github.com/tomasen/realip@latest
```

```go
func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client) // map[ip]client
	)

	// Launch a background goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			// Loop through all clients. If they haven't been seen within the last three minutes, delete them.
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)

		mu.Lock()
		c, found := clients[ip]
		if !found {
			// Create and add a new client struct if it doesn't already exist.
			c = &client{limiter: rate.NewLimiter(2, 4)}
			clients[ip] = c
		}

		// Update the last seen time for the client.
		c.lastSeen = time.Now()

		if !c.limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

```

Then test it

```bash
for i in {1..6}; do curl -s http://localhost:4000/v1/healthcheck | jq -c .; done
```

#### Add rate limiter config in main.go

```go

package main

import (
	"context"      // New import
	"database/sql" // New import
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	"github.com/drsimplegraffiti/greenlight/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}
type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	// flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:pa55word@localhost/greenlight", "PostgreSQL DSN")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice that the default values we're using are the ones we discussed above?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()
	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Info("database connection pool established")
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	// Because the err variable is now already declared in the code above, we need
	// to use the = operator here, instead of the := operator.
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
```

In `middleware.go`:

```go
func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client) // map[ip]client
	)

	// Launch a background goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			// Loop through all clients. If they haven't been seen within the last three minutes, delete them.
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only carry out the check if rate limiting is enabled.
		if app.config.limiter.enabled {
			ip := realip.FromRequest(r)
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					// Use the requests-per-second and burst values from the config
					// struct.
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}
```

Test :

```bash
go run ./cmd/api/ -limiter-burst=2
go run ./cmd/api/ -limiter-enabled=false # Disable rate limiting
```

#### Graceful Shutdown

```go

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"        // New import
	"os/signal" // New import
	"syscall"   // New import
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	// Start a background goroutine.
	go func() {
		// Create a quit channel which carries os.Signal values.
		quit := make(chan os.Signal, 1)
		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel. Any other signals will not be caught by
		// signal.Notify() and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// Read the signal from the quit channel. This code will block until a signal is
		// received.
		s := <-quit
		// Log a message to say that the signal has been caught. Notice that we also
		// call the String() method on the signal to get the signal name and include it
		// in the log entry attributes.
		app.logger.Info("caught signal", "signal", s.String())
		// Exit the application with a 0 (success) status code.
		os.Exit(0)
	}()
	// Start the server as normal.
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	return srv.ListenAndServe()
}
```

#### main.go

```go

package main

import (
	"context"      // New import
	"database/sql" // New import
	"flag"
	"log/slog"
	"os"
	"time"

	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	"github.com/drsimplegraffiti/greenlight/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}
type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
```

#### Excecuting Shutdown

```go

package main

import (
	"context" // New import
	"errors"  // New import
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	// Create a shutdownError channel. We will use this to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)
	go func() {
		// Intercept the signals, as before.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		// Update the log entry to say "shutting down server" instead of "caught signal".
		app.logger.Info("shutting down server", "signal", s.String())
		// Create a context with a 30-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 30-second context deadline is
		// hit). We relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}
	// At this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message.
	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
}
```

#### Setup User

```bash
migrate create -seq -ext=.sql -dir=./migrations create_citext_extension
migrate create -seq -ext=.sql -dir=./migrations create_users_table
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN force 4
```

File: migrations/000004_create_users_table.up.sql

```sql

CREATE TABLE IF NOT EXISTS users (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
name text NOT NULL,
email citext UNIQUE NOT NULL,
password_hash bytea NOT NULL,
activated bool NOT NULL,
version integer NOT NULL DEFAULT 1
);
```

File: migrations/000004_create_users_table.down.sql

```sql
DROP TABLE IF EXISTS users;
```

The email column has the type citext (case-insensitive text). This type stores text data
exactly as it is inputted — without changing the case in any way — but comparisons
against the data are always case-insensitive… including lookups on associated indexes. 2. We’ve also got a UNIQUE constraint on the email column. Combined with the citext
type, this means that no two rows in the database can have the same email value — even
if they have different cases. This essentially enforces a database-level business rule that
no two users should exist with the same email address. 3. The password_hash column has the type bytea (binary string). In this column we’ll store
a one-way hash of the user’s password generated using bcrypt — not the plaintext
password itself.

```go
package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/drsimplegraffiti/greenlight/internal/data"
	"github.com/drsimplegraffiti/greenlight/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// convert email to lowercase
	user.Email = strings.ToLower(user.Email)

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

```

##### User request

```bash
$ BODY='{"name": "Alice Smith", "email": "alice@example.com", "password": "pa55word"}'
$ curl -i -d "$BODY" localhost:4000/v1/users
HTTP/1.1 201 Created
Content-Type: application/json
Date: Mon, 15 Mar 2021 14:42:58 GMT
Content-Length: 152
{
"user": {
"id": 1,
"created_at": "2021-03-15T15:42:58+01:00",
"name": "Alice Smith",
"email": "alice@example.com",
"activated": false
}
}
```

#### Check database

```bash
$ psql $GREENLIGHT_DB_DSN
Password for user greenlight:
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=> SELECT * FROM users;
id | created_at | name | email | password_hash | activated | version
----+------------------------+-------------+-------------------+-------------------------------------+-----------+---------
1 | 2021-04-11 14:29:45+02 | Alice Smith | alice@example.com | \x24326124313224526157784d67356d... | f | 1
(1 row)
```

##### Email template

```user_welcome.tmpl
File: internal/mailer/templates/user_welcome.tmpl
{{define "subject"}}Welcome to Greenlight!{{end}}
{{define "plainBody"}}
Hi,
Thanks for signing up for a Greenlight account. We're excited to have you on board!
For future reference, your user ID number is {{.ID}}.
Thanks,
The Greenlight Team
{{end}}
{{define "htmlBody"}}
<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width" />
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
</head>
<body>
<p>Hi,</p>
<p>Thanks for signing up for a Greenlight account. We're excited to have you on board!</p>
<p>For future reference, your user ID number is {{.ID}}.</p>
<p>Thanks,</p>
<p>The Greenlight Team</p>
</body>
</html>
{{end}}
```

##### Test

```bash
$ BODY='{"name": "Bob Jones", "email": "bob@example.com", "password": "pa55word"}'
$ curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
{
"user": {
"id": 3,
"created_at": "2021-04-11T20:26:22+02:00",
"name": "Bob Jones",
"email": "bob@example.com",
"activated": false
}
}
Time: 2.331957
```

#### server

```go

package main

import (
	"context" // New import
	"errors"  // New import
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	// Create a shutdownError channel. We will use this to receive any errors returned
	// by the graceful Shutdown() function.
	shutdownError := make(chan error)
	go func() {
		// Intercept the signals, as before.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		// Update the log entry to say "shutting down server" instead of "caught signal".
		app.logger.Info("shutting down server", "signal", s.String())
		// Create a context with a 30-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 30-second context deadline is
		// hit). We relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	// Otherwise, we wait to receive the return value from Shutdown() on the
	// shutdownError channel. If return value is an error, we know that there was a
	// problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}
	// At this point we know that the graceful shutdown completed successfully and we
	// log a "stopped server" message.
	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
}
```

#### User activation

```bash
migrate create -seq -ext .sql -dir ./migrations create_tokens_table
```

File: migrations/000005_create_tokens_table.up.sql

```sql

CREATE TABLE IF NOT EXISTS tokens (
hash bytea PRIMARY KEY,
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
expiry timestamp(0) with time zone NOT NULL,
scope text NOT NULL
);
```

File: migrations/000005_create_tokens_table.down.sql

```sql

DROP TABLE IF EXISTS tokens;
```

#### Add token

```go

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// convert email to lowercase
	user.Email = strings.ToLower(user.Email)

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		// As there are now multiple pieces of data that we want to pass to our email
		// templates, we create a map to act as a 'holding structure' for the data. This
		// contains the plaintext version of the activation token for the user, along
		// with their ID.
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		// Send the welcome email, passing in the map above as dynamic data.
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})
	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
```

```bash
$ psql $GREENLIGHT_DB_DSN
Password for user greenlight:
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=> SELECT * FROM tokens;
hash | user_id | expiry | scope
--------------------------------------------------------------------+---------+------------------------+------------
\x09bcb40206b25fe511bfef4d56cbe8c4a141869fc29612fa984b371ef086f5f5 | 7 | 2021-04-18 20:25:41+02 | activation
```

```go

package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"time"

	"github.com/drsimplegraffiti/greenlight/internal/validator"
)

// Define constants for the token scope. For now we just define the scope "activation"
// but we'll add additional scopes later in the book.
const (
	ScopeActivation = "activation"
)

// Define a Token struct to hold the data for an individual token. This includes the
// plaintext and hashed versions of the token, associated user ID, expiry time and
// scope.
type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) *Token {
	// Create a Token instance. In this, we set the Plaintext field to be a random
	// token generated by rand.Text(), and also set values for the user ID, expiry, and
	// scope of the token. Notice that we add the provided ttl (time-to-live) duration
	// parameter to the current time to get the expiry time?
	token := &Token{
		Plaintext: rand.Text(),
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}
	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to
	// work with we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

// Define the TokenModel type.
type TokenModel struct {
	DB *sql.DB
}

// The New() method is a shortcut which creates a new Token struct and then inserts the
// data in the tokens table.
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := generateToken(userID, ttl, scope)
	err := m.Insert(token)
	return token, err
}

// Insert() adds the data for a specific token to the tokens table.
func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser() deletes all tokens for a specific user and scope.
func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}
```

#### activate user

```bash
$ curl -X PUT -d '{"token": "invalid"}' localhost:4000/v1/users/activated
{
"error": {
"token": "must be 26 bytes long"
}
}
$ curl -X PUT -d '{"token": "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}' localhost:4000/v1/users/activated
{
"error": {
"token": "invalid or expired activation token"
}
}
```

```bas
➜  nvim curl -X PUT -d '{"token": "AMYIMQE44LHIO7ZYYMQXYZBJJV"}' localhost:4000/v1/users/activated
{
        "user": {
                "id": 11,
                "created_at": "2025-08-30T10:44:50+01:00",
                "name": "Faith Smith",
                "email": "faith@example.com",
                "activated": true
        }
}
```

##### Authenticate token

```bash
$ BODY='{"email": "alice@example.com", "password": "pa55word"}'
$ curl -i -d "$BODY" localhost:4000/v1/tokens/authentication
HTTP/1.1 201 Created
Content-Type: application/json
Date: Fri, 16 Apr 2021 09:03:36 GMT
Content-Length: 125
{
"authentication_token": {
"token": "IEYZQUBEMPPAKPOAWTPV6YJ6RM",
"expiry": "2021-04-17T11:03:36.767078518+02:00"
}
}
```

```bash
Before we continue, let’s quickly have a look at the tokens table in our PostgreSQL database
to check that the authentication token has been created.
$ psql $GREENLIGHT_DB_DSN
Password for user greenlight:
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=> SELECT * FROM tokens WHERE scope = 'authentication';
\x4390d2ff4af7346dd4238ffccb8a5b18e8c3af9aa8cf57852895ad0f8ee2c50d | 1 | 2021-04-17 11:03:37+02 | authentication
```

##### auth bearer token

```bash
➜  nvim curl -H "Authorization: Bearer ENQMSZ2CXPANZK3NPBNYVQELIG" localhost:4000/v1/healthcheck
{
        "status": "available",
        "system_info": {
                "environment": "development",
                "version": "1.0.0"
        }
}
```

#### Add middleware authentication

```go

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requireActivatedUser(app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requireActivatedUser(app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requireActivatedUser(app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requireActivatedUser(app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requireActivatedUser(app.deleteMovieHandler))

	// Add the route for the POST /v1/users endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	// Wrap the router with the rateLimit() middleware.
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
```

Test:

```bash
$ curl -i localhost:4000/v1/movies/1
HTTP/1.1 401 Unauthorized
Content-Type: application/json
Vary: Authorization
Www-Authenticate: Bearer
Date: Fri, 16 Apr 2021 15:59:33 GMT
Content-Length: 66
{
"error": "you must be authenticated to access this resource"
}
```

To use:

```bash
$ BODY='{"email": "alice@example.com", "password": "pa55word"}'
$ curl -d "$BODY" localhost:4000/v1/tokens/authentication
{
"authentication_token": {
"token": "2O4YHHWDHVVWWDNKN2UZR722BU",
"expiry": "2021-04-17T18:03:09.598843181+02:00"
}
}
$ curl -i -H "Authorization: Bearer 2O4YHHWDHVVWWDNKN2UZR722BU" localhost:4000/v1/movies/1
HTTP/1.1 403 Forbidden
Content-Type: application/json
Vary: Authorization
Date: Fri, 16 Apr 2021 16:03:45 GMT
Content-Length: 76
{
"error": "your user account must be activated to access this resource"
}
```

```bash
➜  nvim curl -i -H "Authorization: Bearer SL4OYNOCP4GKKKSH2NTOJXQZHS" localhost:4000/v1/movies/1
HTTP/1.1 200 OK
Content-Type: application/json
Vary: Authorization
Date: Sat, 30 Aug 2025 10:59:32 GMT
Content-Length: 155

{
        "movie": {
                "id": 1,
                "title": "The Matrix",
                "year": 1999,
                "runtime": "136 mins",
                "genres": [
                        "Action",
                        "Sci-Fi"
                ],
                "version": 1
        }
}
```

##### Permission

```bash
migrate create -seq -ext .sql -dir ./migrations add_permissions
/home/alex/Projects/greenlight/migrations/000006_add_permissions.up.
```

```bash
File: migrations/000006_add_permissions.up.sql
CREATE TABLE IF NOT EXISTS permissions (
id bigserial PRIMARY KEY,
code text NOT NULL
);
CREATE TABLE IF NOT EXISTS users_permissions (
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
PRIMARY KEY (user_id, permission_id)
);
-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
('movies:read'),
('movies:write');
```

Drop

```bash
File: migrations/000006_add_permissions.down.sql
DROP TABLE IF EXISTS users_permissions;
DROP TABLE IF EXISTS permissions;
```

Migrate:

```bash
migrate -path ./migrations -database $GREENLIGHT_DB_DSN up
```

#### Give permissions

```bash
$ psql $GREENLIGHT_DB_DSN
Password for user greenlight:
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=>
And execute the following statements:
-- Set the activated field for alice@example.com to true.
UPDATE users SET activated = true WHERE email = 'alice@example.com';
-- Give all users the 'movies:read' permission
INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'movies:read') FROM users;
-- Give faith@example.com the 'movies:write' permission
INSERT INTO users_permissions
VALUES (
(SELECT id FROM users WHERE email = 'faith@example.com'),
(SELECT id FROM permissions WHERE code = 'movies:write')
);
-- List all activated users and their permissions.
SELECT email, array_agg(permissions.code) as permissions
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.activated = true
GROUP BY email;
```

![image14](assets/image14.png)

```bash
To begin with, let’s try making some requests as alice@example.com to our
GET /v1/movies/1 and DELETE /v1/movies/1 endpoints. The first request should work
correctly, but the second should fail because the user doesn’t have the necessary
movies:write permission.
$ BODY='{"email": "alice@example.com", "password": "pa55word"}'
$ curl -d "$BODY" localhost:4000/v1/tokens/authentication
{
"authentication_token": {
"token": "OPFXEPOYZWMGNWXWKMYIMEGATU",
"expiry": "2021-04-17T20:49:39.963768416+02:00"
}
}
$ curl -H "Authorization: Bearer OPFXEPOYZWMGNWXWKMYIMEGATU" localhost:4000/v1/movies/1
{
"movie": {
"id": 1,
"title": "Moana",
"year": 2016,
"runtime": "107 mins",
"genres": [
"animation",
"adventure"
],
"version": 1
}
}
$ curl -X DELETE -H "Authorization: Bearer OPFXEPOYZWMGNWXWKMYIMEGATU" localhost:4000/v1/movies/1
{
"error": "your user account doesn't have the necessary permissions to access this resource"
}
```

```bash
$ BODY='{"name": "Grace Smith", "email": "grace@example.com", "password": "pa55word"}'
$ curl -d "$BODY" localhost:4000/v1/users
{
"user": {
"id": 8,
"created_at": "2021-04-16T21:32:56+02:00",
"name": "Grace Smith",
"email": "grace@example.com",
"activated": false
}
}
```

```bash
$ psql $GREENLIGHT_DB_DSN
Password for user greenlight:
psql (15.4 (Ubuntu 15.4-1.pgdg22.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.
greenlight=> SELECT email, code FROM users
greenlight-> INNER JOIN users_permissions ON users.id = users_permissions.user_id
greenlight-> INNER JOIN permissions ON users_permissions.permission_id = permissions.id
greenlight-> WHERE users.email = 'grace@example.com';
email | code
-------------------+-------------
grace@example.com | movies:read
(1 row)
```

```sql
SELECT email, code FROM users
INNER JOIN users_permissions ON users.id = users_permissions.user_id
INNER JOIN permissions ON users_permissions.permission_id = permissions.id
WHERE users.email = 'grace@example.com';
```

#### Cors Demonstration

OK, let’s try this out. Go ahead and start up the new application:

```bash

$ go run ./cmd/examples/cors/simple
2021/04/17 17:23:14 starting server on :9000
```

And then open a second terminal window and start our regular API application at the same
time:

```bash

```

At this point, you should now have the API running on the origin http://localhost:4000
and the webpage with the JavaScript running on the origin http://localhost:9000.
Because the ports are different, these are two different origins.
So, when you visit http://localhost:9000 in your web browser, the fetch() action to
http://localhost:4000/v1/healthcheck should be forbidden by the same-origin policy.
Specifically, our API should receive and process the request, but your web browser should
block the response from being read by the JavaScript code.
Let’s take a look. If you open your web browser and visit http://localhost:9000, you
should see the Simple CORS header followed by an error message similar to this:

![image17](assets/image17.png)

```go
func (app *application) enableCORS(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Access-Control-Allow-Origin", "*")
next.ServeHTTP(w, r)
})
}
```

And in the route

```go
//...
return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
```

```bash
$ go run ./cmd/api -cors-trusted-origins="https://www.example.com https://staging.example.com"

```

Update cors

```go
func (app *application) enableCORS(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Add the "Vary: Origin" header.
w.Header().Add("Vary", "Origin")
// Get the value of the request's Origin header.
origin := r.Header.Get("Origin")
// Only run this if there's an Origin request header present.
if origin != "" {
// Loop through the list of trusted origins, checking to see if the request
// origin exactly matches one of them. If there are no trusted origins, then
// the loop won't be iterated.
for i := range app.config.cors.trustedOrigins {
if origin == app.config.cors.trustedOrigins[i] {
// If there is a match, then set a "Access-Control-Allow-Origin"
// response header with the request origin as the value and break
// out of the loop.
w.Header().Set("Access-Control-Allow-Origin", origin)
break
}
}
}
// Call the next handler in the chain.
next.ServeHTTP(w, r)
})
}
```

```bash
go run ./cmd/api -cors-trusted-origins="http://localhost:9000 http://localhost:9001"
```

#### cors and preflight

```go

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		// Add the "Vary: Access-Control-Request-Method" header.
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")
		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat
					// it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response headers, as discussed
						// previously.
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						// Write the headers along with a 200 OK status and return from
						// the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
```

then test :

```bash
➜  gox go run ./cmd/examples/cors/simple
2025/08/30 13:09:31 starting server on :900
➜  gox go run ./cmd/api -cors-trusted-origins="http://localhost:9000"
time=2025-08-30T13:09:12.107+01:00 level=INFO msg="database connection pool established"
time=2025-08-30T13:09:12.107+01:00 level=INFO msg="starting server" addr=:4000 env=development
```

![image19](assets/image19.png)

#### Metrics

use the `expvar` package

```go
// Register a new GET /debug/vars endpoint pointing to the expvar handler.
router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
```

```bash
Go ahead and restart the API, passing in a couple of command line flags for demonstration
purposes. Like so:
$ go run ./cmd/api -limiter-enabled=false -port=4000
time=2023-09-10T10:59:13.722+02:00 level=INFO msg="database connection pool established"
time=2023-09-10T10:59:13.722+02:00 level=INFO msg="starting server" addr=:4000 env=development
And if you visit http://localhost:4000/debug/vars in your web browser, you should see a
JSON response containing information about your running application.
```

![metric](assets/metric.png)

#### Create custom metric

```go
package main
import (
"context"
"database/sql"
"expvar" // New import
"flag"
"log/slog"
"os"
"strings"
"sync"
"time"
"greenlight.alexedwards.net/internal/data"
"greenlight.alexedwards.net/internal/mailer"
_ "github.com/lib/pq"
)
// Remember, our version number is just a constant string (for now).
const version = "1.0.0"
...
func main() {
...
// Publish a new "version" variable in the expvar handler containing our application
// version number (currently the constant "1.0.0").
expvar.NewString("version").Set(version)
app := &application{
config: cfg,
logger: logger,
models: data.NewModels(db),
mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
}
err = app.serve()
if err != nil {
logger.Error(err.Error())
os.Exit(1)
}
}
```

##### Send batch request using `hey`

```bash

$ BODY='{"email": "alice@example.com", "password": "pa55word"}'
$ hey -d "$BODY" -m "POST" http://localhost:4000/v1/tokens/authentication
```

```bash

$ BODY='{"email": "alice@example.com", "password": "pa55word"}'
$ hey -d "$BODY" -m "POST" http://localhost:4000/v1/tokens/authentication
Summary:
Total: 8.0979 secs
Slowest: 2.4612 secs
Fastest: 1.6169 secs
Average: 1.9936 secs
Requests/sec: 24.6977
Total data: 24975 bytes
Size/request: 124 bytes
Response time histogram:
1.617 [1] |■
1.701 [6] |■■■■■
1.786 [10] |■■■■■■■■■
1.870 [26] |■■■■■■■■■■■■■■■■■■■■■■■
1.955 [36] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
2.039 [46] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
2.123 [36] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
2.208 [21] |■■■■■■■■■■■■■■■■■■
2.292 [12] |■■■■■■■■■■
2.377 [4] |■■■
2.461 [2] |■■
Latency distribution:
10% in 1.8143 secs
25% in 1.8871 secs
50% in 1.9867 secs
75% in 2.1000 secs
90% in 2.2017 secs
95% in 2.2642 secs
99% in 2.3799 secs
Details (average, fastest, slowest):
DNS+dialup: 0.0009 secs, 1.6169 secs, 2.4612 secs
DNS-lookup: 0.0005 secs, 0.0000 secs, 0.0030 secs
req write: 0.0002 secs, 0.0000 secs, 0.0051 secs
resp wait: 1.9924 secs, 1.6168 secs, 2.4583 secs
resp read: 0.0000 secs, 0.0000 secs, 0.0001 secs
Status code distribution:
[201] 200 responses
```

```bash
If you visit the GET /debug/vars endpoint while the hey tool is running, you should see that
your application metrics now look quite different:
{
"cmdline": ["/home/xyborg/.cache/go-build/e7/e7d892484457ded0f23b056922c2cc10cf3f395c70481455e33d805c93509284-d/api"],
"database": {"MaxOpenConnections":25,"OpenConnections":4,"InUse":0,"Idle":4,"WaitCount":0,"WaitDuration":0,"MaxIdleClosed":0,"MaxIdleTimeClosed":0,"MaxLifetimeClosed":0},
"goroutines": 8,
"memstats": {"Alloc":1945864,"TotalAlloc":2657176,"Sys":12737552,"Lookups":0,"Mallocs":14946,"Frees":5444,"HeapAlloc":1945864,"HeapSys":7536640,"HeapIdle":3244032,"HeapInuse":4292608,"HeapReleased":2154496,"HeapObjects":9502,"StackInuse":851968,"StackSys":851968,"MSpanInuse":136640,"MSpanSys":146880,"MCacheInuse":9664,"MCacheSys":15704,"BuckHashSys":3838,"GCSys":2136048,"OtherSys":2046474,"NextGC":4194304,"LastGC":1756556623895552537,"PauseTotalNs":577402,"PauseNs":[577402,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"PauseEnd":[1756556623895552537,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"NumGC":1,"NumForcedGC":0,"GCCPUFraction":0.00003070069452852875,"EnableGC":true,"DebugGC":false,"BySize":[{"Size":0,"Mallocs":0,"Frees":0},{"Size":8,"Mallocs":165,"Frees":19},{"Size":16,"Mallocs":3438,"Frees":1259},{"Size":24,"Mallocs":822,"Frees":233},{"Size":32,"Mallocs":1065,"Frees":512},{"Size":48,"Mallocs":2386,"Frees":596},{"Size":64,"Mallocs":908,"Frees":218},{"Size":80,"Mallocs":845,"Frees":258},{"Size":96,"Mallocs":590,"Frees":130},{"Size":112,"Mallocs":689,"Frees":393},{"Size":128,"Mallocs":383,"Frees":86},{"Size":144,"Mallocs":283,"Frees":83},{"Size":160,"Mallocs":57,"Frees":7},{"Size":176,"Mallocs":12,"Frees":0},{"Size":192,"Mallocs":24,"Frees":9},{"Size":208,"Mallocs":118,"Frees":12},{"Size":224,"Mallocs":412,"Frees":122},{"Size":240,"Mallocs":7,"Frees":3},{"Size":256,"Mallocs":56,"Frees":10},{"Size":288,"Mallocs":228,"Frees":80},{"Size":320,"Mallocs":226,"Frees":80},{"Size":352,"Mallocs":613,"Frees":215},{"Size":384,"Mallocs":5,"Frees":1},{"Size":416,"Mallocs":11,"Frees":2},{"Size":448,"Mallocs":126,"Frees":1},{"Size":480,"Mallocs":8,"Frees":1},{"Size":512,"Mallocs":25,"Frees":9},{"Size":576,"Mallocs":8,"Frees":4},{"Size":640,"Mallocs":5,"Frees":1},{"Size":704,"Mallocs":6,"Frees":3},{"Size":768,"Mallocs":15,"Frees":5},{"Size":896,"Mallocs":32,"Frees":3},{"Size":1024,"Mallocs":6,"Frees":1},{"Size":1152,"Mallocs":62,"Frees":5},{"Size":1280,"Mallocs":5,"Frees":1},{"Size":1408,"Mallocs":13,"Frees":3},{"Size":1536,"Mallocs":18,"Frees":2},{"Size":1792,"Mallocs":9,"Frees":1},{"Size":2048,"Mallocs":34,"Frees":0},{"Size":2304,"Mallocs":2,"Frees":0},{"Size":2688,"Mallocs":1,"Frees":0},{"Size":3072,"Mallocs":1,"Frees":1},{"Size":3200,"Mallocs":0,"Frees":0},{"Size":3456,"Mallocs":1,"Frees":0},{"Size":4096,"Mallocs":109,"Frees":0},{"Size":4864,"Mallocs":4,"Frees":0},{"Size":5376,"Mallocs":1,"Frees":0},{"Size":6144,"Mallocs":1,"Frees":1},{"Size":6528,"Mallocs":9,"Frees":5},{"Size":6784,"Mallocs":0,"Frees":0},{"Size":6912,"Mallocs":0,"Frees":0},{"Size":8192,"Mallocs":19,"Frees":1},{"Size":9472,"Mallocs":9,"Frees":0},{"Size":9728,"Mallocs":0,"Frees":0},{"Size":10240,"Mallocs":0,"Frees":0},{"Size":10880,"Mallocs":0,"Frees":0},{"Size":12288,"Mallocs":1,"Frees":0},{"Size":13568,"Mallocs":0,"Frees":0},{"Size":14336,"Mallocs":0,"Frees":0},{"Size":16384,"Mallocs":0,"Frees":0},{"Size":18432,"Mallocs":2,"Frees":1}]},
"timestamp": 1756556729,
"version": "1.0.0"
}
```

##### Request-level Metrics

The total number of requests received.
The total number of responses sent.
The total (cumulative) time taken to process all requests in microseconds.
In `middleware.go` :

```go
func (app *application) metrics(next http.Handler) http.Handler {
// Initialize the new expvar variables when the middleware chain is first built.
var (
totalRequestsReceived = expvar.NewInt("total_requests_received")
totalResponsesSent = expvar.NewInt("total_responses_sent")
totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
)
// The following code will be run for every request...
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Record the time that we started to process the request.
start := time.Now()
// Use the Add() method to increment the number of requests received by 1.
totalRequestsReceived.Add(1)
// Call the next handler in the chain.
next.ServeHTTP(w, r)
// On the way back up the middleware chain, increment the number of responses
// sent by 1.
totalResponsesSent.Add(1)
// Calculate the number of microseconds since we began to process the request,
// then increment the total processing time by this amount.
duration := time.Since(start).Microseconds()
totalProcessingTimeMicroseconds.Add(duration)
})
}
```

#### Use the metric

```go
// Use the new metrics() middleware at the start of the chain.
return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
```

Run the app

```bash
go run ./cmd/api -limiter-enabled=false
```

then go to the debug/var endpoint:

![metrix](assets/metrix.png)

#### Using makefile

```bash
brew install make
choco install make
sudo apt install make
```

```make


## help: print this help message

help:
@echo 'Usage:'
@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# Create the new confirm target.

confirm:
@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application

run/api:
go run ./cmd/api

## db/psql: connect to the database using psql

db/psql:
psql ${GREENLIGHT_DB_DSN}

migration:
@echo 'Creating migration files for ${name}...'
migrate create -seq -ext=.sql -dir=./migrations ${name}
up:
@echo 'Running up migrations...'
migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

## db/migrations/new name=$1: create a new database migration

db/migrations/new:
@echo 'Creating migration files for ${name}...'
migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations

db/migrations/up: confirm
@echo 'Running up migrations...'
migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up
```

And if you now execute the help target, you should get a response which lists all the
available targets and the corresponding help text. Similar to this:
$ make help
Usage:
help print this help message
run/api run the cmd/api application
db/psql connect to the database using psql
db/migrations/new name=$1 create a new database migration
db/migrations/up apply all up database migrations
I should also point out that positioning the help rule as the first thing in the Makefile is a
deliberate move. If you run make without specifying a target then it will default to executing
the first rule in the file.
So this means that if you try to run make without a target you’ll now be presented with the
help information, like so:
$ make
Usage:
help print this help message
run/api run the cmd/api application
db/psql connect to the database using psql
db/migrations/new name=$1 create a new database migration
db/migrations/up apply all up database migrations

##### Quality

```bash
go get -tool honnef.co/go/tools/cmd/staticcheck@latest
go tool staticcheck --version
go tool staticcheck ./...
```

#### BUILD

```bash
$ ./bin/api -port=4040 -db-dsn=postgres://greenlight:pa55word@localhost/greenlight
```

#### Cross platform compilation

```bash
go tool dist list

```

To see a list of all the operating system/architecture combinations that Go supports, you can
run the go tool dist list command like so:
$ go tool dist list
aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/amd64

And you can specify the operating system and architecture that you want to create the
binary for by setting GOOS and GOARCH environment variables when running go build . For
example:
$ GOOS=linux GOARCH=amd64 go build {args}
In the next section of the book, we’re going to walk through how to deploy an executable
binary on an Ubuntu Linux server hosted by Digital Ocean. For this we’ll need a binary which
is designed to run on a machine with a linux/amd64 OS and architecture combination.
So let’s update our make build/api rule so that it creates two binaries — one for use on your
local machine, and another for deploying to the Ubuntu Linux server.

## build/api: build the cmd/api application

```bash

.PHONY: build/api
build/api:
@echo 'Building cmd/api...'
go build -ldflags='-s' -o=./bin/api ./cmd/api
GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
```

You should see that two binaries are now created — with the cross-compiled binary located
under the ./bin/linux_amd64 directory, like so:
$ make build/api
Building cmd/api...
go build -ldflags='-s' -o=./bin/api ./cmd/api
GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
$ tree ./bin
./bin
├── api
└── linux_amd64
└── api

##### Version number

```bash
$ make build/api
Building cmd/api...
go build -ldflags="-s" -o="./bin/api" ./cmd/api
GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o="./bin/linux_amd64/api" ./cmd/api
$ ./bin/api -version
Version: 1.0.0
```

#### Go version

```bash
➜  gox go version -m ./bin/api
./bin/api: go1.24.1
        path    github.com/drsimplegraffiti/greenlight/cmd/api
        mod     github.com/drsimplegraffiti/greenlight  (devel)
        dep     github.com/julienschmidt/httprouter     v1.3.0
        dep     github.com/lib/pq       v1.10.9
        dep     github.com/tomasen/realip       v0.0.0-20180522021738-f0c99a92ddce
        dep     github.com/wneessen/go-mail     v0.6.2
        dep     golang.org/x/crypto     v0.41.0
        dep     golang.org/x/text       v0.28.0
        dep     golang.org/x/time       v0.12.0
        build   -buildmode=exe
        build   -compiler=gc
        build   -ldflags=-s
        build   CGO_ENABLED=1
        build   CGO_CFLAGS=
        build   CGO_CPPFLAGS=
        build   CGO_CXXFLAGS=
        build   CGO_LDFLAGS=
        build   GOARCH=amd64
        build   GOOS=linux
        build   GOAMD64=v
```

##### Jwt

```bash
go get github.com/pascaldekloe/jwt@v1
```
