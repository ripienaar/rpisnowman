# What?

This is a small tool that runs the [Ryanteck RTK-000-00A GPIO Snowman](https://wiki.ryanteck.uk/Category:RTK-000-00A) through a couple of schemes of light flashes in random order

A pre-built binary is in the releases section, download onto your Pi and run it.

## Building

You can build your own binary from this code:

```
env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags '-w -s'
```

This will build the standard binary deployable to RPIs.

## Management

A backplane is embedded using the [Choria Backplane](https://github.com/choria-io/go-backplane) system, it's setup in insecure mode so can work with any NATS server without TLS enabled like demo.nats.io.

Start the snowman as follows to enable this:

```
SNOWMAN_NAME=bob BACKPLANE_BROKERS=demo.nats.io:4222 ./rpisnowman
```

With the backplane client setup you can run:

```
$ backplane --insecure bob pause
```

To turn the light display off - later issue a `resume` action and it will start again.

## Contact?

R.I.Pienaar / [www.devco.net](https://www.devco.net/) / [@ripienaar](https://twitter.com/ripienaar)
