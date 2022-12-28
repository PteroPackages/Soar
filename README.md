# Soar ✈️
Soar is a command line interface tool for interacting with the [Pterodactyl Game Panel](https://pterodactyl.io) API. The name is inspired by Pterodactyl's Wings application because... birds. This tool covers most of the application and client API, including support for creating resources from your terminal!

## Installing
See the [releases page](https://github.com/PteroPackages/Soar/releases) for available downloads.

### Building from Source
```bash
git clone https://github.com/PteroPackages/Soar.git
cd Soar && make
```

## Getting Started
After installing, run `soar config init` to generate a config file. On Linux-based systems this can be found in the user config directory (usually `$HOME/.config/.soar/config.yml`), and on Windows systems it can be found at `%APPDATA%\.soar\config.yml`. You can also specify the `--dir=` flag to generate the config in a specific directory. Next, enter your credentials for the application and client section (you can also set other options). Now you're ready to soar!

**Note:** by default Soar will check for a local config to use, if not found then it will use the global config. If you have a local config but don't want to use it, you can specify the `--global` or `-g` flag in the command to force use the global config.

## Usage
Soar has a convinient naming convention for its commands:

```
soar <api> <resource>:<action> [args]
      ---   --------   ------   ----
       |        |        |      additional arguments or flags
       |        |        |
       |        |      the action to execute (get, create, delete, etc)
       |        |
       |    the resource to target (users, servers, nodes, etc)
       |
      the API to target (application or client)
```

This naming convention is designed to be compact and readable, so you don't need to memorize every command or search the help command to figure out what it does (you can still do this if you want to, though). Some resource commands are flattened for convinience like the `soar client files:list` command which lists the files of a specified server, and is much quicker to type than `soar client servers:files:list`.

## Supported Resources

### Application
* users
* * [X] get
* * [X] create
* * [ ] update
* * [X] delete
* servers
* * [X] get
* * [ ] create
* * [ ] update build
* * [ ] update details
* * [ ] update startup
* * [X] suspend/unsuspend
* * [X] reinstall
* * [X] delete
* nodes
* * [X] get
* * [X] get configuration
* * [ ] create
* * [ ] update
* * [X] delete
* locations
* * [X] get
* * [X] create
* * [ ] update
* * [X] delete
* nests
* * [X] get
* * eggs
* * [X] get

### Client

* account
* * [X] get
* * [X] get permissions
* * [X] get activities
* * [ ] update email
* * [ ] update password
* * 2FA
* * * [X] get
* * * [X] enable
* * * [X] disable
* * API Keys
* * * [X] get
* * * [ ] create
* * * [ ] delete
* * SSH Keys
* * * [ ] get
* * * [ ] create
* * * [ ] delete
* servers
* * [X] get
* * [X] get activities
* * [X] get resource usage
* * [X] get server websocket auth
* * [X] send server command
* * [X] send server power state
* * databases
* * [X] get
* * files
* * * [X] get
* * * [X] download
* * * [X] rename
* * * [X] copy
* * * [X] write
* * * [X] create
* * * [X] compress
* * * [X] decompress
* * * [X] delete
* * * [X] create folder
* * * [X] change file permissions
* * * [X] pull remote file
* * * [X] upload files
* * subusers
* * * [X] get
* * * [X] add
* * * [X] remove
* * startup
* * * [X] get
* * * [X] set variable
* * settings
* * * [X] set docker image
* * * [X] rename
* * * [X] reinstall

This list of commands is subject to change as the API develops, if there is a command or feature that isn't here feel free to open an issue or PR requesting it!

## Contributing
1. [Fork this repo](https://github.com/PteroPackages/Soar/fork)!
2. Make a branch from `main` (`git branch -b <new feature>`)
3. Commit your changes (`git commit -am "..."`)
4. Open a PR here (`git push origin <new feature>`)

## Contributors
* [Devonte W](https://github.com/devnote-dev) - creator and maintainer

This repository is managed under the MIT license.

© 2022-present PteroPackages
