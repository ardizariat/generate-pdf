# Security policy

## Supported versions

Only the latest release receives fixes.

## Reporting a vulnerability

Please do not open a public issue for security problems. Use GitHub's private
vulnerability reporting on the repository's **Security** tab, or email the
maintainer at the address on the GitHub profile.

Include the affected version, a description of the issue, and reproduction steps.
You can expect an initial response within seven days.

## Scope

`pdfgen` generates local PDF files from fictional data. It does not make network
requests and does not handle real personal data. Relevant reports include path
traversal through output paths, unsafe file overwrites, and dependency
vulnerabilities.
