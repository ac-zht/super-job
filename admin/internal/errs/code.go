package errs

const (
	TaskInternalServerError = 10000
	TaskInvalidInput        = 10001

	ExecutorInternalServerError = 11000
	ExecutorInvalidInput        = 11001
	ExecutorRequiredNotInput    = 11002

	SettingInternalServerError = 12000
	SettingInvalidInput        = 12001

	InstallInternalServerError = 13000
	InstallOccurred            = 13001
)
