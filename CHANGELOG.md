# Changelog

### Deprecation Notice 201701

Atmotool is being deprecated and replaced with [rwctl](https://github.com/ghchinoy/rwctl). Atmotool will remain in maintenance until dependent tools are updated (ex. [yeoman theme generator](https://www.npmjs.com/package/generator-akana-theme))

### 1.7.6
* API details, basic info

### 1.7.5
* policies list and list policies; refactor in process

### 1.7.4 20161107
* Create API methods now have better result output
* Delete users method added

### 1.7.3 20161106
* Create API with OAI spec file, proxy pattern
* Updated control login to return userinfo; updated various methods to use new return info
* API listings (list listversions) have initial dynamic table display

### 1.7.2 20161104
* removing some debug statements from output when debug flag was not set

### 1.7.1 20161007
* fix for reset and rebuild theme

### 1.7.0 20161001
* Add API variations and related CMS additions and refactoring
* explicit method on curl debug output
* uplift addCsrfHeader function to control package

### 1.6.1 20160928
* `version` as a non-flag cmd added

### 1.6.0 20160926
* initial refactor of command structure (ex. `apis list`), while keeping old one (ex. `list apis`)
* added `apis list`, `apis metrics <apiID>`, `apis logs <apiID>`
* added `cms list`
* homebrew installation for OS X (ref [tap](https://github.com/ghchinoy/homebrew-akana))

### 1.5.1 20160918
* formatting for `list apis` and `list policies`

### 1.5.0 20160822

* added a CMS list function

### 1.4.6 20160818

* Rebuilding styles outputs which theme is being rebuilt
* Rebuilding & resetting functions now look to config struct (config file) first for theme, can still be overridden by cli param

### 1.4.5 20160610

* minor change to listUsers - outputs response to stdout rather than log

### 1.4.4 20160322

* added `list topapis` which makes a call to `/api/businesses/tenantbusiness.enterpriseapi/metrics`

### 1.4.3 20160316

* added listUsers which uses `/api/search` to list users, outputs json.

### 1.4.2 20160314

* debug statements in loginToCM and uploadFile available via `--debug` flag

### 1.4.1 20160226

* `--debug` flag added for more verbose output; relegated some existing output to behind debug fence

### 1.4.0 20160222

* CSRF now properly implemented

### 1.3.4 20151230

* added a debugging function curlThis() which outputs the equivalent curl command
* added an unused method called listUsers() - there's no equivalent CM API that lists all users, and will remove this in the future

### 1.3.3 20151116

* Changed file upload methods to use logged in http.Client
* Added query param Type with proper values to 'list policies' method

### 1.3.2 20151116

* Changed how login is handled (returns http.Client)

### 1.3.1 20150918

* Login now occurs prior to rebuilding styles
* added a Go cross-compile script for ease of building releases

### 1.3.0 20150810
* added `reset` action to reset CM to blank (i18n, favicon, landing/index.htm, and custom.less)

### 1.2.0
* renamed to `atmotool` throughout
* removed `--dir` flag from `zip`
* added license

### 1.1.1

* fixes issue #5, adds `.conf` to exclude list for zipping, so as not to include `local.conf`
* fixes issue #6, remove trailing slash from dir when zipping

### 1.1.0

* `download` action, downloads a zip of a `--path`

### 1.0.2

* `config` flag made optional

### 1.0.1

* uploads files

### 1.0.0

* uploads less

