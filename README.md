# SOA Software Community Manager Customizer

A tool to customize SOA Software's Community Manager.


## Installation

Clone git repo to `$GOPATH/src/bitbucket.org/ghchinoy` then `cd atmotool`

    go get github.com/docopt/docopt-go
    go install

Requires a config file (typically called environment`.config`, ex. `local.config` or `eap.config`) that contains CM url, username, and password, in JSON format:

```
{
	"url": "http://local.cm.demo:9900/atmosphere",
	"email": "administrator@cm.demo",
	"password": "password"
}
```

## Capabilities

### Build zipfiles

Builds zipfiles, suitable for uploading to Community Manager

    atmotool zip --prefix <prefix> --config <config> [--dir <dir>]

* prefix: Prefix for zip to be created, PREFIX_resourcesThemeDefault.zip will be generated
* config: Config file, see above
* dir: base directory of the CM customization files, defaults to the current directory

Example usage creating two zipfiles

    atmotool zip --prefix test --config ./testdata/local.conf --dir ./testdata/testfiles

Outputs would be, in the current working directory:

* `test_resourcesThemeDefault.zip` for use to upload to CM's CMS at /resources/theme/default
* `test_contentHomeLanding.zip` for use to upload to CM's CMS at /content/home/landing

### Upload Less file

    atmotool upload less <file> --config <config>

Will upload a `.less` file to Community Manager. Automatically names less file `custom.less` when uploading.

### Upload to CM CMS

Uploads to CM's CMS, allowing user to specify file name and path.

    atmotool upload file --path <path> --config <config> <files>...

Note, that if the filename ends in `.zip`, zip expansion will occur at the target path.

Example usage of uploading customization zipfile to `/content/home/landing`

    atmotool upload file --path /content/home/landing --config local.conf prospect_contentHomeLanding.zip

### Upload customizations to CM

Looks for and uploads `custom.less`, `PREFIX_resourcesThemDefault.zip`, `PREFIX_contentHomeLanding.zip` to Community Manager located at `ATMO_BASE_URL`

    atmotool upload all --config <config> [--dir <dir>]

* config: config file, as above
* dir: base directory for the CM customizations to upload, defaults to current directory




## Development

### Notes

Using `docopt.org` for command line argument processing

