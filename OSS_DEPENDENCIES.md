# Open Source Software Dependencies Changelog

This document records all top-level dependencies included in this Terraform provider project, including license information, inclusion rationale, security assessments, and approval details.

## Dependencies

### github.com/hashicorp/terraform-plugin-docs v0.20.1

- **License:** MPL-2.0 (Mozilla Public License 2.0)
- **Date of Decision to Include:** 2025-11-04
- **Why it's needed:** Required for generating Terraform provider documentation from code annotations and examples. Used in the `go generate` command to automatically create documentation for the Terraform registry.
- **Security Concerns:** No known security vulnerabilities found. HashiCorp maintains this package as part of their official Terraform provider tooling.
- **Approver:** @madsquist

---

### github.com/hashicorp/terraform-plugin-framework v1.13.0

- **License:** MPL-2.0 (Mozilla Public License 2.0)
- **Date of Decision to Include:** 2025-11-04
- **Why it's needed:** Core framework dependency for building Terraform providers. Provides the foundation for defining resources, data sources, and provider schemas. This is the primary framework used throughout the provider codebase for implementing all Terraform resources and data sources.
- **Security Concerns:** No known security vulnerabilities found. This is the official HashiCorp framework for building Terraform providers and is actively maintained.
- **Approver:** @madsquist

---

### github.com/hashicorp/terraform-plugin-go v0.25.0

- **License:** MPL-2.0 (Mozilla Public License 2.0)
- **Date of Decision to Include:** 2025-11-04
- **Why it's needed:** Low-level Terraform plugin protocol implementation. Required for communication between Terraform Core and the provider plugin. Provides the gRPC interface and protocol buffer definitions necessary for the provider to function.
- **Security Concerns:** No known security vulnerabilities found. This is an official HashiCorp package that implements the Terraform plugin protocol.
- **Approver:** @madsquist

---

### github.com/hashicorp/terraform-plugin-log v0.9.0

- **License:** MPL-2.0 (Mozilla Public License 2.0)
- **Date of Decision to Include:** 2025-11-04
- **Why it's needed:** Provides structured logging functionality specifically designed for Terraform providers. Used throughout the codebase (e.g., `tflog` package) for logging provider operations, API calls, and debugging information. Essential for troubleshooting and monitoring provider behavior.
- **Security Concerns:** No known security vulnerabilities found. This is an official HashiCorp logging library for Terraform providers.
- **Approver:** @madsquist

---

### github.com/hashicorp/terraform-plugin-testing v1.11.0

- **License:** MPL-2.0 (Mozilla Public License 2.0)
- **Date of Decision to Include:** 2025-11-04
- **Why it's needed:** Testing utilities and helpers for writing Terraform provider acceptance tests. Provides the `helper/resource` package used in test files to run Terraform configurations and validate provider behavior. Critical for ensuring the provider functions correctly.
- **Security Concerns:** No known security vulnerabilities found. This is an official HashiCorp testing library for Terraform providers.
- **Approver:** @madsquist

---

### github.com/google/uuid v1.6.0

- **License:** Apache-2.0 (Apache License 2.0)
- **Date of Decision to Include:** 2025-11-04
- **Why it's needed:** Provides UUID generation and parsing functionality. Used in the provider codebase (e.g., `utils.go` and various resource files) for generating unique identifiers and validating UUID format in Terraform configurations. Essential for resource identification and validation.
- **Security Concerns:** No known security vulnerabilities found. This is a well-maintained Google package with widespread adoption in the Go ecosystem.
- **Approver:** @madsquist

---

### github.com/hashicorp/terraform-plugin-framework-validators v0.16.0

- **License:** MPL-2.0 (Mozilla Public License 2.0)
- **Date of Decision to Include:** 2025-11-04
- **Why it's needed:** Provides validation functions for Terraform provider schemas. Used extensively throughout the provider codebase (e.g., `stringvalidator`, `listvalidator`, `int64validator`) to validate user input in Terraform configurations. Ensures data integrity and provides helpful error messages to users.
- **Security Concerns:** No known security vulnerabilities found. This is an official HashiCorp validation library for Terraform providers.
- **Approver:** @madsquist

---

## Notes

- All HashiCorp dependencies are licensed under MPL-2.0, which is compatible with most open-source licenses.
- The Google UUID package is licensed under Apache-2.0, which is permissive and compatible with MPL-2.0.
- All dependencies have been reviewed for security concerns as of the initial setup date (2025-11-04).
- This document should be updated when new dependencies are added or existing ones are updated.

  
---

**Document Maintained By:** Maximilian Beller  
**Last Updated:** 2025-11-04
