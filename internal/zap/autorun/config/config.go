// Package config implements config which is given to OWASP Zap.
package config

// Config provided to OWASP Zap.
type Config struct {
	Env  Env   `yaml:"env"`
	Jobs []Job `yaml:"jobs"`
}

// Env config provided to OWASP Zap.
type Env struct {
	Contexts   []Context  `yaml:"contexts"`
	Parameters Parameters `yaml:"parameters"`
}

// Context config provided to OWASP Zap.
type Context struct {
	Name         string   `yaml:"name"`
	Urls         []string `yaml:"urls"`
	ExcludePaths []string `yaml:"excludePaths,omitempty"`
	IncludePaths []string `yaml:"includePaths,omitempty"`
}

// Parameters for Env config provided to OWASP Zap.
type Parameters struct {
	FailOnError       bool `yaml:"failOnError"`
	FailOnWarning     bool `yaml:"failOnWarning"`
	ProgressToStdout  bool `yaml:"progressToStdout"`
	ContinueOnFailure bool `yaml:"continueOnFailure"`
}

// Job config provided to OWASP Zap.
type Job struct {
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Parameters interface{} `yaml:"parameters"`
}

// SpiderParameters to Job config.
type SpiderParameters struct {
	MaxDuration int32 `yaml:"maxDuration,omitempty"`
	MaxDepth    int32 `yaml:"maxDepth,omitempty"`
}

// PassiveScanParameters to Job config.
type PassiveScanParameters struct {
	ScanOnlyInScope bool `yaml:"scanOnlyInScope,omitempty"`
}

// ReportParameters to Job config.
type ReportParameters struct {
	ReportTitle       string `yaml:"reportTitle,omitempty"`
	ReportDescription string `yaml:"reportDescription,omitempty"`
	Template          string `yaml:"template,omitempty"`
	ReportFile        string `yaml:"reportFile,omitempty"`
}
