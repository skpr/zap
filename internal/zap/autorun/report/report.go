// Package report implements reporting which given by OWASP Zap.
package report

import (
	"encoding/json"
	"fmt"
	"os"
)

// Summary for reporting the scans result.
type Summary struct {
	High    int
	Medium  int
	Low     int
	Info    int
	Details []Detail
}

// Detail from the report.
type Detail struct {
	Name        string
	Description string
	Severity    string
}

// Report provided by the scan.
type Report struct {
	Site []Site `json:"site"`
}

// Site provided in the report.
type Site struct {
	Alerts []Alert `json:"alerts"`
}

// Alert provided in the report.
type Alert struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	RiskCode    string `json:"riskcode"` // WHY IS THIS A STRING!!!!
}

// GetSummaryFromFile from a report file.
func GetSummaryFromFile(filePath string) (Summary, error) {
	var summary Summary

	file, err := os.Open(filePath)
	if err != nil {
		return summary, fmt.Errorf("failed to load file: %w", err)
	}

	var report Report

	if err = json.NewDecoder(file).Decode(&report); err != nil {
		return summary, fmt.Errorf("failed to decode json: %w", err)
	}

	for _, site := range report.Site {
		for _, alert := range site.Alerts {
			detail := Detail{
				Name:        alert.Name,
				Description: alert.Description,
			}

			switch alert.RiskCode {
			case "0":
				detail.Severity = "Info"
				summary.Info++
			case "1":
				detail.Severity = "Low"
				summary.Low++
			case "2":
				detail.Severity = "Medium"
				summary.Medium++
			case "3":
				detail.Severity = "High"
				summary.High++
			}

			summary.Details = append(summary.Details, detail)
		}
	}

	return summary, nil
}
