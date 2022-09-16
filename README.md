# Soar ✈️
Soar is a command line interface tool for interacting with the [Pterodactyl Game Panel](https://pterodactyl.io) API. The name is inspired by Pterodactyl's Wings application because... birds.

## Installing
See the [releases page](https://github.com/PteroPackages/Soar/releases).

### Building from Source
```bash
git clone https://github.com/PteroPackages/Soar.git
cd Soar && make
# or 'go build' if on Windows without make
```

## Getting Started
After installing, run `soar config init` to generate a config file. On Linux-based systems this can be found in the user config directory (usually `$HOME/.config/.soar/config.yml`), and on Windows systems it can be found at `%APPDATA%\.soar\config.yml`. You can also specify the `--dir=` flag to generate the config in a specific directory. Next, enter your credentials for the application and client section (you can also set other options). Now you're ready to soar!

**Note:** by default Soar will check for a local config to use, if not found then it will use the global config. If you have a local config but don't want to use it, you can specify the `--global` or `-g` flag in the command to force use the global config.

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
files:chmod             changes the permissions of a file
files:compress          compresses one or more files and folders
files:contents          gets the contents of a file
files:copy              copies a file
files:create            creates an empty file
files:decompress        decompresses an archived file
files:decompress        decompresses an archived file
files:delete            deletes one or more files
files:download          downloads a file or returns the url
files:folder            creates a folder
files:info              gets the file info for a specific file
files:list              lists files on a server
files:pull              pulls a file from a remote source
files:rename            renames a file on the server
files:upload            uploads one or more files to the server
files:write             writes content to a file
servers:activity        gets the server activity logs
servers:command         sends a command to the server console
servers:get             gets account servers
servers:power           sets the server power state
servers:resources       gets server resource usage
servers:websocket       gets the server websocket data
settings:image          sets the docker image for a server
settings:reinstall      reinstalls a server
settings:rename         renames a server
startup:get             gets the startup information for a server
startup:set             updates a startup variable on a server
subusers:add            adds a subuser to the server
subusers:get            gets the server subusers
subusers:remove         removes a subuser from the server
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
