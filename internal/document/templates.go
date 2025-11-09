package document

import "law_for_it/internal/model"

var defaultTemplates = []model.Template{
	{
		Name:        "gdpr_basic",
		DisplayName: "GDPR Compliance Policy",
		Description: "A concise GDPR compliance policy covering data collection, storage and user rights.",
		Tags:        []string{"gdpr", "privacy", "compliance"},
		Fields: []model.TemplateField{
			{Key: "company_name", Label: "Company Name", Type: "text", Required: true},
			{Key: "contact_email", Label: "Contact Email", Type: "email", Required: true},
			{Key: "data_protection_officer", Label: "Data Protection Officer", Type: "text", Required: false},
			{Key: "processing_purposes", Label: "Processing Purposes", Type: "textarea", Required: true},
			{Key: "retention_period", Label: "Retention Period", Type: "text", Required: true},
			{Key: "legal_basis", Label: "Legal Basis", Type: "textarea", Required: true},
			{Key: "user_rights", Label: "User Rights Summary", Type: "textarea", Required: true},
		},
		Content: `# {{.company_name}} GDPR Compliance Policy\n\n## Contact Information\n- Contact Email: {{.contact_email}}\n{{- if .data_protection_officer}}\n- Data Protection Officer: {{.data_protection_officer}}\n{{- end}}\n\n## Processing Purposes\n{{.processing_purposes}}\n\n## Legal Basis\n{{.legal_basis}}\n\n## Data Retention\nWe retain personal data for {{.retention_period}}.\n\n## User Rights\n{{.user_rights}}\n`,
	},
	{
		Name:        "cookie_standard",
		DisplayName: "Cookie Policy",
		Description: "Standard cookie policy describing categories, purposes and consent.",
		Tags:        []string{"cookie", "tracking", "website"},
		Fields: []model.TemplateField{
			{Key: "website_name", Label: "Website Name", Type: "text", Required: true},
			{Key: "cookie_categories", Label: "Cookie Categories", Type: "textarea", Required: true},
			{Key: "consent_mechanism", Label: "Consent Mechanism", Type: "textarea", Required: true},
			{Key: "third_parties", Label: "Third-Party Services", Type: "textarea", Required: false},
			{Key: "update_frequency", Label: "Policy Update Frequency", Type: "text", Required: true},
		},
		Content: `# {{.website_name}} Cookie Policy\n\n## Cookie Categories\n{{.cookie_categories}}\n\n## Consent\n{{.consent_mechanism}}\n\n## Third-Party Services\n{{if .third_parties}}{{.third_parties}}{{else}}No third-party tracking cookies are in use.{{end}}\n\n## Updates\nThis policy is reviewed {{.update_frequency}}.\n`,
	},
}
