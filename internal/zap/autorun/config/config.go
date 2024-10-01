package config

type Config struct {
	Env  Env   `yaml:"env"`
	Jobs []Job `yaml:"jobs"`
}

type Env struct {
	Contexts   []Context  `yaml:"contexts"`
	Parameters Parameters `yaml:"parameters"`
}

type Context struct {
	Name         string   `yaml:"name"`
	Urls         []string `yaml:"urls"`
	ExcludePaths []string `yaml:"excludePaths,omitempty"`
	IncludePaths []string `yaml:"includePaths,omitempty"`
}

type Parameters struct {
	FailOnError       bool `yaml:"failOnError"`
	FailOnWarning     bool `yaml:"failOnWarning"`
	ProgressToStdout  bool `yaml:"progressToStdout"`
	ContinueOnFailure bool `yaml:"continueOnFailure"`
}

type Job struct {
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Parameters interface{} `yaml:"parameters"`
}

type SpiderParameters struct {
	MaxDuration int32 `yaml:"maxDuration,omitempty"`
	MaxDepth    int32 `yaml:"maxDepth,omitempty"`
}

type PassiveScanParameters struct {
	ScanOnlyInScope bool `yaml:"scanOnlyInScope,omitempty"`
}

type ReportParameters struct {
	ReportTitle       string `yaml:"reportTitle,omitempty"`
	ReportDescription string `yaml:"reportDescription,omitempty"`
	Template          string `yaml:"template,omitempty"`
	ReportFile        string `yaml:"reportFile,omitempty"`
}
