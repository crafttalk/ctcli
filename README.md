# ctcli - CraftTalk Command Line Tool

## Commands

### help
Shows help

### version
Shows version

### install
Installs a package into specified root folder. Example:
```shell
domain install /path/to/package.tar.gz
```

### upgrade
Upgrades a release with binaries contained in given package. Backups previous release. Example:
```shell
domain upgrade /path/to/package.tar.gz
```

### rollback
Install previous release
```shell
domain rollback
```

### list-releases
Shows full list of releases and highlights current release.
```shell
domain list-releases
```

### delete
Deletes a release completely. Example:
```shell
domain delete
```

### start
Starts current installation. Example:
```shell
domain start [app]
```

### stop
Stops current installation. Example:
```shell
domain stop [app]
```

## Root Directory File Structure

```
/backups - Backups and etc.
/config - Services' configs
/current-release - Current installation
/logs - Logs
/tmp - Temporary directory
```
