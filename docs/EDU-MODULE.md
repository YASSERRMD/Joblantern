# Education & Visa Mill Detection

Fake universities and "study abroad" visa mills exploit students
through the same playbook as recruitment scams: promises of foreign
opportunity, urgent decision, and large upfront payments.

## Rule pack

See [rules/millpack.go](../internal/edu/rules/millpack.go). High-
severity codes:

- `unaccredited` — not in any recognised national registry.
- `life-experience-degree` — degrees for life experience.
- `dubious-accreditor` — accreditor matches mill pattern.
- `program-not-in-catalogues` — claimed program absent from public catalogues.

## Visa pathway plausibility

[visa.Check](../internal/edu/visa/pathway.go) confirms the institution
can plausibly issue the form named in the offer (I-20, CAS, COE, LOA,
Zulassungsbescheid).

## Illegal agent commissions

Origin-country regulators frequently prohibit student-paid agent
commissions. See [commission.Flag](../internal/edu/commission/illegal.go).

## Bundled study-and-work

The "graduate-and-work" pattern is treated as a higher risk because
it ties admission to a specific named employer at the time of offer
— a signature of visa-mill scams. See
[combined.Merge](../internal/edu/combined/student_to_job.go).

## Worked example

```text
Submission:
  institution: "Pacific Western University"
  program:     "Doctorate of Management"
Result:
  band: red
  flags: unaccredited, life-experience-degree, dubious-accreditor,
         young-domain-unaccredited, program-not-in-catalogues
  citation: chea.org registry lookup returned no result
```
