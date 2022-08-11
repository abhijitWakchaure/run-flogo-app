package config

// Constants for app config
const (
	AppName         = "run-flogo-app"
	ConfigFileName  = ".run-flogo-app"
	MaxAppsWithList = 5

	DefaultAppPatternLinux   = `^.+-linux_amd64.*$`
	DefaultAppPatternWindows = `^.+-windows_amd64.*$`
	DefaultAppPatternDarwin  = `^.+-darwin_amd64.*$`

	GithubLastestReleaseURL = "https://api.github.com/repos/abhijitWakchaure/run-flogo-app/releases/latest"
	GithubDownloadBaseURL   = "https://github.com/abhijitWakchaure/run-flogo-app/releases/download/"
	GithubBaseURL           = "https://github.com/abhijitWakchaure/run-flogo-app"
	GithubIssuesURL         = "https://github.com/abhijitWakchaure/run-flogo-app/issues"

	InstallPathLinux   = "/usr/local/bin"
	InstallPathDarwin  = "/usr/local/bin"
	InstallPathWindows = `C:\Windows\system32`
)
