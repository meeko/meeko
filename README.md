# ciderapp

`ciderapp` is a command line utility for managing local Cider applications.

Cider RPC is being used as the transport for the management calls, so `ciderapp`
is actually a Cider application itself. It just connects to a Cider RPC endpoint,
executes a single management call and exits.

```text
APPLICATION:
  ciderapp - Cider applications management utility

USAGE:
  ciderapp [-debug] [-endpoint ENDPOINT] SUBCMD

VERSION:
  0.0.1

OPTIONS:
  -debug=false: print debug output
  -endpoint="": Cider ZeroMQ 3.x RPC endpoint
  -h=false: print help and exit

DESCRIPTION:
  ciderapp is a command line utility for managing local Cider instance,
  or rather the Cider applications running on it.

  This tool expects the local Cider instance's management token to be saved
  in .cider_token file places in the current user's home directory.

SUBCOMMANDS:
  env            show app variable values
  info           show app info
  install        install a new app
  list           list installed apps
  remove         uninstall an app
  restart        restart a running app
  set            set app variable
  start          start an app
  status         show app status
  stop           stop a running app
  unset          unset app variable
  upgrade        upgrade an existing app
  
```

# License

MIT, see the `LICENSE` file.
