# mk #

A Meeko PaaS CLI management utility.

## Overview ##

mk is a command line utility for managing Meeko agents.

Under the hook, mk is implemented as a short-lived Meeko agent that uses
the Meeko RPC service to communicate with the agent supervisor component.

## Installation ##

mk must be installed from sources right now. It uses [Godep](https://github.com/tools/godep) for vendoring.

```
$ cd mk
$ godep go build
$ ./mk -h
```

## License ##

MIT, see the `LICENSE` file.
