package application

// Metadata holds build metadatas.
type Metadata struct {
	Branch     string `json:"branch"`
	Compiler   string `json:"compiler"`
	CompiledAt string `json:"compiledAt"`
	Sha        string `json:"sha"`
}

// ToMap returns a map from given metadata.
func (m Metadata) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Branch":     m.Branch,
		"Compiler":   m.Compiler,
		"CompiledAt": m.CompiledAt,
		"Sha":        m.Sha,
	}
}
