package logger

type EnvType string

const (
	EnvTypeUnspecified EnvType = EnvType("Unspecified")
	EnvTypeLocal               = "Local"
	EnvTypeTesting             = "Testing"
	EnvTypeProd                = "Prod"
)
