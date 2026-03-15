package persona

import (
	"fmt"
)

// Capability represents a specific action a persona is permitted to perform.
type Capability string

const (
	CapabilityLifecycle  Capability = "lifecycle"   // Moving tasks through todo -> review
	CapabilityAudit      Capability = "audit"       // General audit capability
	CapabilityAuditFunctional Capability = "audit_functional" // Functional gate (QA)
	CapabilityAuditSecurity   Capability = "audit_security"   // Security gate (Sec Ops)
	CapabilityBreakGlass Capability = "break_glass" // Forcing status via break-glass
	CapabilityManagement Capability = "management"  // Managing projects/milestones (Architect only)
)

var registry = map[string][]Capability{
	"architect":         {CapabilityLifecycle, CapabilityAudit, CapabilityBreakGlass, CapabilityManagement},
	"senior-engineer":   {CapabilityLifecycle},
	"junior-engineer":   {CapabilityLifecycle},
	"developer":         {CapabilityLifecycle},
	"qa-functional":     {CapabilityAudit, CapabilityAuditFunctional},
	"security-ops":      {CapabilityAudit, CapabilityAuditSecurity},
	"doc-writer":        {}, // Read-only role
}

// Linter enforces the Doctrine of Command by verifying role compliance.
type Linter struct {
}

// NewLinter creates a new Persona Linter.
func NewLinter() *Linter {
	return &Linter{}
}

// PersonaAudit verifies if the given role possesses the required capability.
func (l *Linter) PersonaAudit(role string, capability Capability) error {
	caps, ok := registry[role]
	if !ok {
		return fmt.Errorf("Request Rejected: Unknown persona '%s'. Dispatch the correct agent.", role)
	}

	for _, c := range caps {
		if c == capability {
			return nil
		}
	}

	return fmt.Errorf("Request Rejected: Outside my jurisdiction. Dispatch the correct agent. (role '%s' lacks capability '%s')", role, capability)
}

// Lint remains for backwards compatibility with command-level checks if needed,
// but Audit is preferred for engine-level enforcement.
func (l *Linter) Lint(role string, commandName string, allowedRoles []string) error {
	if len(allowedRoles) == 0 {
		return nil
	}

	permitted := false
	for _, r := range allowedRoles {
		if r == role {
			permitted = true
			break
		}
	}

	if !permitted {
		return fmt.Errorf("Request Rejected: Outside my jurisdiction. Dispatch the correct agent. (role '%s' cannot execute '%s')", role, commandName)
	}

	return nil
}
