package internal

// appVersion is set at build time by goreleaser (-X), so it must be a var
var appVersion = "v0.0.0"

func GetAppVersion() string {
	return appVersion
}
