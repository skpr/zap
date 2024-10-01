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
		domain    = vp.GetString("domain")
		directory = vp.GetString("directory")
	)

	if err := run(directory, domain); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(directory, domain string) error {
	if directory == "" {
		return fmt.Errorf("directory is required")
	}

	if domain == "" {
		return fmt.Errorf("domain is required")
	}

	fmt.Println("Writing autorun file to: sdfsdfdsf")

	content := zapconfig.Config{
		Env: zapconfig.Env{
			Contexts: []zapconfig.Context{
				{
					Name: "default",
					Urls: []string{
						domain,
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
				Parameters: map[string]bool{
					"scanOnlyInScope": true,
				},
			},
			{
				Name: "spider",
				Type: "spider",
				Parameters: map[string]int{
					"maxDuration": 60,
					"maxDepth":    1,
				},
			},
			{
				Name: "passiveScan-wait",
				Type: "passiveScan-wait",
			},
			{
				Name: "pdf",
				Type: "report",
				Parameters: map[string]string{
					"reportTitle":       "Automated Vulnerability Scan",
					"reportDescription": "This is an automated report",
					"template":          "traditional-pdf",
					// @todo, This is import. We pick this report.pdf file up if/when we post to Slack.
					"reportFile": "report",
				},
			},
			{
				Name: "json",
				Type: "report",
				Parameters: map[string]string{
					"template": "traditional-json",
					// @todo, This is import. We pick this report.json file up when we report on the results.
					"reportFile": "report",
				},
			},
		},
	}

	var (
		configPath = path.Join(directory, "zap.yaml")
		reportPath = path.Join(directory, "report.json")
	)

	data, err := yaml.Marshal(&content)
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

	return nil
}
