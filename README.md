# ctcli - CraftTalk Command Line Tool

## Commands

### help
Shows help

### version
Shows version

### install
Installs a package into specified root folder. Example:
```shell
ctcli install /path/to/package.tar.gz
```

### upgrade
Upgrades a release with binaries contained in given package. Backups previous release. Example:
```shell
ctcli upgrade /path/to/package.tar.gz
```

### rollback
Install previous release
```shell
ctcli rollback
```

### list-releases
Shows full list of releases and highlights current release.
```shell
ctcli list-releases
```

### delete
Deletes a release completely. Example:
```shell
ctcli delete
```

### start
Starts current installation. Example:
```shell
ctcli start [app]
```

### stop
Stops current installation. Example:
```shell
ctcli stop [app]
```

## Root Directory File Structure

```
/backups - Backups and etc.
/config - Services' configs
/current-release - Current installation
/logs - Logs
/tmp - Temporary directory
```
