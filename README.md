# meeko #

A Meeko CLI management utility.

## Overview ##

meeko is a command line utility for managing Meeko agents.

Under the hook, meeko is implemented as a short-lived Meeko agent that uses
the Meeko RPC service to communicate with the agent supervisor component.

### Commands ###

```
  env        show agent variable values
  info       show agent info
  install    install a new agent
  list       list installed agents
  remove     uninstall an agent
  restart    restart a running agent
  set        set agent variable
  start      start an agent
  status     show agent status
  stop       stop a running agent
  unset      unset agent variable
  upgrade    upgrade an installed agent
  watch      stream agent logs
```

## Installation ##

meeko must be installed from sources right now. It uses [Godep](https://github.com/tools/godep) for vendoring.

```
$ cd meeko
$ godep go build
$ ./meeko -h
```

## Configuration ##

meeko expects a file called `.meekorc` to be present in the user's home directory.
It is a YAML file that must contain the following keys:

```yaml
endpoint_address: Meeko RPC service WebSocket endpoint address (URL)
access_token:     Meeko RPC service WebSocket endpoint token
management_token: Meeko management token
```

## License ##

MIT, see the `LICENSE` file.
