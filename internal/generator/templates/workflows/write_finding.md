---
description: The Security Ops agent's subroutine for writing a detailed, formal vulnerability finding.
---

### Doctrine: The Scribe of Weakness

When you find a vulnerability, you become an auditor. Your purpose is to create a formal finding that allows the engineer to understand the weakness, its potential impact, and how to mitigate it.

1.  **The Law of the CVE:** Your finding MUST identify the specific type of vulnerability (e.g., "SQL Injection," "Cross-Site Scripting"). Reference a CWE or CVE if possible.
2.  **The Law of the Attack Vector:** Your finding MUST describe the specific location of the vulnerability in the code (file and line number) and the method used to exploit it.
3.  **The Law of Mitigation:** Your finding MUST include a clear recommendation for how to fix the vulnerability.

### Sequence: The Finding Protocol

1.  **Log the Vulnerability Finding**
    *   `castra note add --role security-ops --project "%%project_id%%" --task "%%task_id%%" --content "%%vulnerability_finding%%" --tags "finding,security,vulnerability"`

### Variables

*   `%%project_id%%`: **[Input]** The ID of the project.
*   `%%task_id%%`: **[Input]** The ID of the compromised task.
*   `%%vulnerability_finding%%`: **[Input]** A structured markdown block containing: 1. Vulnerability Type (CWE), 2. Location/Attack Vector, 3. Recommended Mitigation.

