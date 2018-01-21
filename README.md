# What?

This is a small tool that runs the [Ryanteck RTK-000-00A GPIO Snowman](https://wiki.ryanteck.uk/Category:RTK-000-00A) through a couple of schemes of light flashes in random order

A pre-built binary is in the releases section, download onto your Pi and run it.

## Building

You can build your own binary from this code:

```
$ glide install
$ env GOOS=linux GOARCH=amd64 go build -tags default
```

This will build the standard binary deployable to RPIs.

There is a [Choria](https://choria.io) based management agent that can be anabled:

```
$ glide install
$ env GOOS=linux GOARCH=amd64 go build -tags choria -ldflags "-X main.middleware=demo.nats.io:4222"
```

The snowman will now connect to a management network from where you can use Choria to pause and resume it's operation.

## Contact?

R.I.Pienaar / [www.devco.net](https://www.devco.net/) / [@ripienaar](https://twitter.com/ripienaar)
