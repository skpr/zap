// Package main implements entrypoint for application.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/aquasecurity/table"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	zapconfig "github.com/skpr/zap/internal/zap/autorun/config"
	zapreport "github.com/skpr/zap/internal/zap/autorun/report"
)

func main() {
	vp := viper.New()

	vp.AutomaticEnv()
	vp.SetEnvPrefix("ZAP_WRAPPER")

	var (
		endpoint = vp.GetString("endpoint")

		directory = vp.GetString("directory")
		scanType  = vp.GetString("type")
	)

	if err := run(directory, endpoint, scanType); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(directory, endpoint, scanType string) error {
	if directory == "" {
		return fmt.Errorf("directory is required")
	}

	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	if scanType == "" {
		scanType = "passive"
	}

	config, err := getConfig(endpoint, scanType)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	fmt.Println("Generating Autorun Config File")

	var (
		configPath = path.Join(directory, "zap.yaml")
		reportPath = path.Join(directory, "report.json")
	)

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error while marshalling. %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config data: %w", err)
	}

	fmt.Println("Config file generated:", configPath)

	fmt.Println("Running OWASP ZAP Task")

	cmd := exec.Command("zap.sh", "-cmd", "-autorun", configPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("zap scan failed: %w", err)
	}

	fmt.Println("Checking Results")

	report, err := zapreport.GetSummaryFromFile(reportPath)
	if err != nil {
		return fmt.Errorf("failed to get summary: %w", err)
	}

	t := table.New(os.Stdout)
	t.SetHeaders("High", "Medium", "Low", "Info")
	t.AddRow(strconv.Itoa(report.High), strconv.Itoa(report.Medium), strconv.Itoa(report.Low), strconv.Itoa(report.Info))
	t.Render()

	t = table.New(os.Stdout)
	t.SetHeaders("Severity", "Name", "Description")
	for _, detail := range report.Details {
		t.AddRow(detail.Severity, detail.Name, detail.Description)
	}
	t.Render()

	if report.High > 0 {
		return fmt.Errorf("results marked as high were found")
	}

	return nil
}

func getConfig(endpoint, scanType string) (*zapconfig.Config, error) {
	switch scanType {
	case "passive":
		return getPassiveScanConfig(endpoint)
	case "active":
		return getActiveScanConfig(endpoint)
	default:
		return nil, fmt.Errorf("scan type not supported")
	}
}

func getPassiveScanConfig(endpoint string) (*zapconfig.Config, error) {
	return &zapconfig.Config{
		Env: zapconfig.Env{
			Contexts: []zapconfig.Context{
				{
					Name: "default",
					Urls: []string{
						endpoint,
					},
					ExcludePaths: []string{}, // @todo, Implement in a future release.
					IncludePaths: []string{}, // @todo, Implement in a future release.
				},
			},
			Parameters: zapconfig.Parameters{
				FailOnError:       true,
				FailOnWarning:     false,
				ProgressToStdout:  true,
				ContinueOnFailure: false,
			},
		},
		Jobs: []zapconfig.Job{
			{
				Name: "passiveScan-config",
				Type: "passiveScan-config",
				Parameters: zapconfig.PassiveScanParameters{
					ScanOnlyInScope: true,
				},
			},
			{
				Name: "spider",
				Type: "spider",
				Parameters: zapconfig.SpiderParameters{
					MaxDuration: 60,
					MaxDepth:    1,
				},
			},
			{
				Name: "passiveScan-wait",
				Type: "passiveScan-wait",
			},
			{
				Name: "pdf",
				Type: "report",
				Parameters: zapconfig.ReportParameters{
					ReportTitle:       "Automated Vulnerability Scan",
					ReportDescription: "This is an automated report",
					Template:          "traditional-pdf",
					// @todo, This is import. We pick this report.pdf file up if/when we post to Slack.
					ReportFile: "report",
				},
			},
			{
				Name: "json",
				Type: "report",
				Parameters: zapconfig.ReportParameters{
					Template: "traditional-json",
					// @todo, This is import. We pick this report.json file up when we report on the results.
					ReportFile: "report",
				},
			},
		},
	}, nil
}

func getActiveScanConfig(endpoint string) (*zapconfig.Config, error) {
	return &zapconfig.Config{
		Env: zapconfig.Env{
			Contexts: []zapconfig.Context{
				{
					Name: "default",
					Urls: []string{
						endpoint,
					},
					ExcludePaths: []string{}, // @todo, Implement in a future release.
					IncludePaths: []string{}, // @todo, Implement in a future release.
				},
			},
			Parameters: zapconfig.Parameters{
				FailOnError:       true,
				FailOnWarning:     false,
				ProgressToStdout:  true,
				ContinueOnFailure: false,
			},
		},
		Jobs: []zapconfig.Job{
			{
				Name: "spider",
				Type: "spider",
				Parameters: zapconfig.SpiderParameters{
					MaxDuration: 60,
					MaxDepth:    1,
				},
			},
			{
				Name: "activeScan",
				Type: "activeScan",
			},
			{
				Name: "pdf",
				Type: "report",
				Parameters: zapconfig.ReportParameters{
					ReportTitle:       "Automated Vulnerability Scan",
					ReportDescription: "This is an automated report",
					Template:          "traditional-pdf",
					// @todo, This is import. We pick this report.pdf file up if/when we post to Slack.
					ReportFile: "report",
				},
			},
			{
				Name: "json",
				Type: "report",
				Parameters: zapconfig.ReportParameters{
					Template: "traditional-json",
					// @todo, This is import. We pick this report.json file up when we report on the results.
					ReportFile: "report",
				},
			},
		},
	}, nil
}
