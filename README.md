# ciderapp

`ciderapp` is a command line utility for managing local Cider applications.

The idea is simple - it uses the same Cider client as any other Cider application,
but it only executed a single management command, waits for the result and exits.

```text
$ ./ciderapp 
APPLICATION:
  ciderapp - Cider applications management utility

USAGE:
  ciderapp SUBCMD [options] [arguments]

VERSION:
  0.0.1

OPTIONS:
  -h=false: print help and exit
  -verbose=false: print more verbose output

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

MIT, see the LICENSE file.
