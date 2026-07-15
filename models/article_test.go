package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestExtractTagsResultSalaryRoundTrip guards the JSON contract between the
// prompts' Output Format and the Go struct.
func TestExtractTagsResultSalaryRoundTrip(t *testing.T) {
	raw := `{
		"collaboration_types": [0, 1],
		"work_locations": ["臺北市"],
		"salary": "月薪 6 萬起",
		"salary_detail": {"type": "monthly", "min": 60000, "currency": "TWD", "negotiable": false}
	}`

	var r ExtractTagsResult
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r.Salary == nil || *r.Salary != "月薪 6 萬起" {
		t.Errorf("salary = %v, want 月薪 6 萬起", r.Salary)
	}
	if r.SalaryDetail == nil || r.SalaryDetail.Type != "monthly" || r.SalaryDetail.Min == nil || *r.SalaryDetail.Min != 60000 {
		t.Errorf("salary_detail not parsed: %+v", r.SalaryDetail)
	}

	// 薪資面議 case: negotiable only, no amounts.
	raw = `{"salary": "薪資面議", "salary_detail": {"negotiable": true}}`
	r = ExtractTagsResult{}
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatalf("unmarshal negotiable: %v", err)
	}
	if r.SalaryDetail == nil || !r.SalaryDetail.Negotiable {
		t.Errorf("negotiable not parsed: %+v", r.SalaryDetail)
	}
}

// TestExtractTagsPromptsMentionSalary makes sure both platform prompts carry
// the salary instructions and output keys.
func TestExtractTagsPromptsMentionSalary(t *testing.T) {
	for _, pt := range []PlatformType{PlatformTypeApen, PlatformTypeNurse, PlatformTypePhar} {
		p := GetExtractTagsSystemPrompt(pt)
		for _, must := range []string{`"salary"`, `"salary_detail"`, "薪資面議"} {
			if !strings.Contains(p, must) {
				t.Errorf("%s prompt missing %q", pt, must)
			}
		}
	}
}
