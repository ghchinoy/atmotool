# Akana Community Manager CLI

A command-line tool to manage [Akana](https://www.akana.com/)'s [Community Manager API Portal](https://www.akana.com/products/api-portal).

For downloads, please see the [releases](https://github.com/ghchinoy/atmotool/releases).


```
Usage:
  atmotool zip --prefix <prefix> <dir>
  atmotool upload less <file> [--config <config>] [--debug]
  atmotool upload file --path <path> <files>... [--config <config>] [--debug]
  atmotool download --path <path> <filename> [--config <config>] [--debug]
  atmotool list apis [--config <config>] [--debug]
  atmotool list topapis [--config <config>] [--debug]
  atmotool list apps [--config <config>] [--debug]
  atmotool list users [--config <config>] [--debug]
  atmotool list policies [--config <config>] [--debug]
  atmotool list cms [<path>] [--config <config>] [--debug]
  atmotool rebuild [<theme>] [--config <config>] [--debug]
  atmotool reset [<theme>] [--config <config>] [--debug]
  atmotool -h | --help
  atmotool --version
```

## Capabilities

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


## Installation

Clone git repo to `$GOPATH/src/bitbucket.org/ghchinoy`, change into the 
directory (`cd atmotool`), get the prerequisite, and issue the go 
install command.

    go get github.com/docopt/docopt-go
    go install

Requires a config file (typically called _environment_`.config`, ex. `local.config` or `eap.config`) that contains CM url, username, and password, in JSON format:

```
{
    "url": "http://local.cm.demo:9900",
    "email": "administrator@cm.demo",
    "password": "password"
}
```

Note, no CM context (ex. `/atmosphere` or `/enterpriseapi`).



## Development Notes


Using `docopt.org` for command line argument processing

