# ctcli - CraftTalk Command Line Tool

[![Actions Status](https://github.com/crafttalk/ctcli/workflows/build-and-test/badge.svg)](https://github.com/crafttalk/ctcli/actions)
[![codecov](https://codecov.io/gh/crafttalk/ctcli/branch/master/graph/badge.svg)](https://codecov.io/gh/crafttalk/ctcli)

## Commands

### help
Shows help

### version
Shows version

### init 
Initializes specified root directory as a ctcli dir.
```shell
ctcli --root /path/to/dir init
```

### install
Installs a package into specified root folder. Example:
```shell
ctcli install /path/to/package.tar.gz
```

First it creates tmp dir and extracts an archive. Then it creates oci image config and rootfs folder. 
Then it copies unpacked images to current-release folder. Makes symlinks for logs folder and configs.
Creates a zero backup.

### upgrade
Upgrades a release with binaries contained in given package. Backups previous release. Example:
```shell
ctcli upgrade /path/to/package.tar.gz
```

Creates a tmp dir and extracts the archive. Then it creates oci image configs and rootfs folders. 
Next it checks diff between apps in package and apps in current release and prompts for upgrade.
If user decided to upgrade, creates a backup, deletes apps from current release and copies files 
from tmp folder to current release.

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

### logs
Shows logs (stdout/stderr) of specified service. Example:
```shell
ctcli logs [--tail <n>] [-f] <app>
```

Options:
* `--tail 100` - Show latest 100 bytes
* `-f` - Follow log file 

### status
Show current status. Example:
```shell
ctcli status
```

### exec
Execute a command inside a container. Example:
```shell
ctcli exec busybox cat /etc/hosts
```

### backup
Makes an archive containing `current-release/`, `data/` and `config/` folders and puts it into `releases/` folder  

* `--ignore-data` flag disables `data/` folder backup

```shell
ctcli backup [--ignore-data]
```

### release-info
Shows current release app list and each app version
```shell
ctcli release-info
```

### rollback
Install previous release
```shell
ctcli rollback /path/to/archive.tar.gz
```

## To Be Done

### list-releases
Shows full list of releases and highlights current release.
```shell
ctcli list-releases
```

## Root Directory File Structure

```
/releases           - Previous releases and etc.
/config             - Services' configs
/current-release    - Current installation
/logs               - Logs
/data               - app persistent data
/runc-root          - runc containers root dir 
/tmp                - Temporary directory
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
            ./skopeo                - contains oci image that umoci will then extract
```

## Backups

Backup is an archive and has following structure:
```
/config     - app configs
/apps       - apps
/data
```
