# Deprecation policy

When a stable surface (see `docs/STABILITY.md`) is removed or changed:

1. **Announce.** Mark deprecated in source comments + the doc page.
   Open a tracking issue.
2. **Warn at runtime.** When the deprecated path is exercised, emit
   `slog.Warn` with the message `"deprecated: <surface> will be removed in vX.0"`.
3. **Hold for one minor release.** A surface deprecated in `v1.2`
   may not be removed before `v1.3`. Removal lands at `v2.0`.
4. **Document the migration.** The CHANGELOG entry for the
   removing release must include a "Migration from `<surface>`"
   subsection with a working example.

Exceptions are limited to security vulnerabilities where the only
defensible fix is breaking the API. Those land per
`SECURITY.md`'s disclosure process.

## Deprecation log

(none yet)
