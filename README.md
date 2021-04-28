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
/releases - Previous releases and etc.
/config - Services' configs
/current-release - Current installation
/logs - Logs
/tmp - Temporary directory
```

## Package structure
```
meta.json   - contains PackageMeta object and has commitSha, imageTag and etc
package/    - contains package content
    ./runc.amd64    - runc binary
    ./umoci.amd64   - umoci binary
    ./apps          - apps in container
        ./{app}     - app in package
            ./package-config.json   - contains AppPackageConfig object
            ./version.json          - contains info about app image version 
```