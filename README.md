# Soar ✈️
Soar is a command line interface tool for interacting with the [Pterodactyl Game Panel](https://pterodactyl.io) API. The name is inspired by Pterodactyl's Wings application because... birds.

## Installing
See the [releases page](https://github.com/PteroPackages/Soar/releases).

### Building from Source
```bash
git clone https://github.com/PteroPackages/Soar.git
cd Soar
go build
```

## Getting Started
After installing, run `soar config init` to generate a config file, this will be at `/etc/.soar/config.yml` on Linux-based systems or `%APPDATA%\.soar\config.yml` on Windows systems. You can also specify the `--dir=` flag to generate the config in a specific directory. Next, enter your credentials for the application and client section (you can also set other options). Now you're ready to soar! **Note:** by default Soar will use the global config for operations, to use the local config you must specify the `--local` or `-l` flag in the command.

## Commands
You can also view available commands by running `soar help <command>`.

### Application API
```
locations:get     gets panel node locations
nests:eggs:get    gets panel eggs for a nest
nests:get         gets panel nests
nodes:config      gets a node config
nodes:get         gets panel nodes
servers:delete    deletes a server
servers:get       gets panel servers
servers:reinstall reinstalls a server
servers:suspend   suspends a server
servers:unsuspend unsuspends a server
users:create      creates a user
users:delete      deletes a user
users:get         gets panel users
```

### Client API
```
account:2fa:disable
account:2fa:enable      enables two-factor on the account
account:2fa:get         gets account two-factor code
account:activity        gets the account activity logs
account:api-keys:delete deletes an api key
account:api-keys:get    gets the account api keys
account:get             gets account information
account:perms           gets system permissions
databases:get           gets server databases
files:contents          gets the contents of a file
files:copy              copies a file
files:download          downloads a file or returns the url
files:list              lists files on a server
files:rename            renames a file on the server
servers:activity        gets the server activity logs
servers:command         sends a command to the server console
servers:get             gets account servers
servers:power           sets the server power state
servers:resources       gets server resource usage
servers:websocket       gets the server websocket data
```

### Config Management
```
config          shows the config
config init     initializes a new config
```

## Contributing
1. [Fork this repo](https://github.com/PteroPackages/Soar/fork)!
2. Make a branch from `main` (`git branch -b <new feature>`)
3. Commit your changes (`git commit -am "..."`)
4. Open a PR here (`git push origin <new feature>`)

## Contributors
* [Devonte W](https://github.com/devnote-dev) - creator and maintainer

This repository is managed under the MIT license.

© 2022 PteroPackages
