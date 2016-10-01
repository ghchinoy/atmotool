package version

const (
	version     = "1.6.1"
	versionName = "cirrus"
)

// Version returns the version
func Version() string {
	return version + " " + versionName
}
