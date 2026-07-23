//go:build pro

package pro

import (
	_ "embed"
	"encoding/json"
)

//go:embed web/flow_builder.html
var flowBuilderHTML string

//go:embed web/omni_inbox.html
var omniInboxHTML string

//go:embed web/agency.html
var agencyHTML string

//go:embed web/templates_extra.json
var templatesExtraJSON string

func init() {
	var extra []FlowTemplate
	if err := json.Unmarshal([]byte(templatesExtraJSON), &extra); err == nil {
		templates = append(templates, extra...)
	}
}
