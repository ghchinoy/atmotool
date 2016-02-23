# Changelog

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
