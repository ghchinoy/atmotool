# Akana API Platform Portal CLI

A command-line tool to manage [Akana](https://www.akana.com/)'s [Community Manager API Portal](https://www.akana.com/products/api-portal).

** * Deprecation notice * **
01/2017

Atmotool is being deprecated and replaced with [rwctl](https://github.com/ghchinoy/rwctl). Atmotool will remain in maintenance until dependent tools are updated (ex. [yeoman theme generator](https://www.npmjs.com/package/generator-akana-theme))

## Installation

### Download

For downloads, please see the [releases](https://github.com/ghchinoy/atmotool/releases) - Linux, Windows and OS X binaries are available.


### Via Homebrew (for OS X)

```
brew update
brew install ghchinoy/akana/atmotool
```

### From Source

Clone git repo to `$GOPATH/src/github.com/ghchinoy`, change into the 
directory (`cd atmotool`), get the prerequisite, and issue the go 
install command.

    go get github.com/docopt/docopt-go
    go install

## Usage

Requires a config file (typically called _environment_`.config`, ex. `local.config` or `eap.config`) that contains CM url, username, and password, in JSON format:

```
{
    "url": "http://local.cm.demo:9900",
    "email": "administrator@cm.demo",
    "password": "password"
}
```

Note, no CM context (ex. `/atmosphere` or `/enterpriseapi`) is needed in the `url`.


```
Usage:
  atmotool zip --prefix <prefix> <dir>
  atmotool upload less <file> [--config <config>] [--debug]
  atmotool upload file --path <path> <files>... [--config <config>] [--debug]
  atmotool download --path <path> <filename> [--config <config>] [--debug]
  atmotool apis list [--config <config>] [--debug]
  atmotool apis metrics <apiId> [--config <config>] [--debug]
  atmotool apis logs <apiId> [--config <config>] [--debug]
  atmotool list apps [--config <config>] [--debug]
  atmotool list users [--config <config>] [--debug]
  atmotool list policies [--config <config>] [--debug]
  atmotool cms list [<path>] [--config <config>] [--debug]
  atmotool rebuild [<theme>] [--config <config>] [--debug]
  atmotool reset [<theme>] [--config <config>] [--debug]
  atmotool -h | --help
  atmotool --version
```

## Capabilities

This section contains a partial description of capabililites.

### Upload Less file

    atmotool upload less <file> [--config <config>]

Will upload a `.less` file to Community Manager using the specified config file. Automatically names less file `custom.less` when uploading.

The config file is optional, will default to looking for `local.conf` in the current directory.

### Upload to CM CMS

Uploads to CM's CMS, allowing user to specify file name and path.

    atmotool upload file --path <path> <files>... [--config <config>]

Note, that if the filename ends in `.zip`, zip expansion will occur at the target path.

The config file is optional, will default to looking for `local.conf` in the current directory.


Example usage of uploading customization zipfile to `/content/home/landing`

    atmotool upload file --path /content/home/landing prospect_contentHomeLanding.zip

### Download a cms path as zip

Downloads a zipfile for the indicated CMS path

    atmotool download --path <path> <filename> [--config <config>]

* path: CM CMS path
* filename: output filename, does not need `.zip`

Example usage

    atmotool download --path /content/home/landing contentHomeLanding

Output would be a zip file `contentHomeLanding.zip` which will contain a zip of the contents of the CMS directory `/content/home/landing`


### Rebuild Styles

Rebuilds the CM styles that already exist for a particular theme; no uploading, see [Upload Less File](#uploadless)

    atmotool rebuild [<theme>] [--config <config>]

* theme: defaults to `default`


### Reset CM

Deletes a pre-selected list of items in CM to "reset" the UI to default out-of-the-box state.

    atmotool reset [<theme>] [--config <config>]

Items deleted (note the `default` theme is modifiable via flag)

* resources/theme/default/i18n
* resources/theme/default/style/images/favico.ico
* resources/theme/default/less/custom.les
* content/home/landing/index.htm

### Build zipfiles

Builds zipfiles, suitable for uploading to Community Manager

    atmotool zip --prefix <prefix> [--dir <dir>] --config <config>

* prefix: Prefix for zip to be created, PREFIX_resourcesThemeDefault.zip will be generated
* config: Config file, see above
* dir: base directory of the CM customization files, defaults to the current directory

Example usage creating two zipfiles

    atmotool zip --prefix test --config ./testdata/local.conf --dir ./testdata/testfiles

Outputs would be, in the current working directory:

* `test_resourcesThemeDefault.zip` for use to upload to CM's CMS at /resources/theme/default
* `test_contentHomeLanding.zip` for use to upload to CM's CMS at /content/home/landing


### Upload customizations to CM

Looks for and uploads `custom.less`, `PREFIX_resourcesThemeDefault.zip`, `PREFIX_contentHomeLanding.zip` to Community Manager located at `ATMO_BASE_URL`

    atmotool upload all --config <config> [--dir <dir>]

* config: config file, as above
* dir: base directory for the CM customizations to upload, defaults to current directory



## Development Notes


Atmotool uses `docopt.org` for command line argument processing.

