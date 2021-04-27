# ctcli - CraftTalk Command Line Tool

## Commands

### help
Shows help

### version
Shows version

### install
Installs a package into specified root folder. Example:
```shell
ctcli install --root /myfolder/ /path/to/package.tar.gz
```

### delete
Deletes a completely. Example:
```shell
ctcli delete --root /myfolder/
```

### start
Starts current installation. Example:
```shell
ctcli start --root /myfolder/
```

### stop
Stops current installation. Example:
```shell
ctcli stop --root /myfolder/
```

## Root Directory File Structure

```
/backups - Backups and etc.
/config - Services' configs
/current-release - Current installation
/logs - Logs
/tmp - Temporary directory
```
