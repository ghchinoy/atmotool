# Akana Community Manager Customizer

A tool to customize Akana's Community Manager.


## Installation

Clone git repo to `$GOPATH/src/bitbucket.org/apihussain` then `cd atmotool`

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

## Capabilities - _Working_

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


### Build zipfiles

Builds zipfiles, suitable for uploading to Community Manager

    atmotool zip --prefix <prefix> [--dir <dir>]

* prefix: Prefix for zip to be created, PREFIX_resourcesThemeDefault.zip will be generated
* dir: base directory of the CM customization files, defaults to the current directory

Example usage creating a single zipfile

    atmotool zip --prefix TEST --dir ./resources/theme/default

Output would be, in the current working directory:

    TEST_-resources-theme-default.zip


## Capabilities - _Planned_


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




## Development

### 1.0.2

* `config` flag made optional

### 1.0.1

* uploads files

### 1.0.0

* uploads less

### Notes

Using `docopt.org` for command line argument processing

