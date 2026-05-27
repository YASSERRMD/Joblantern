# Rental & Housing Scams

Migrant workers commonly hit paired scams: a fake recruitment offer
plus a fake housing arrangement run by the same network. Phase 39
extends Joblantern's evidence model so housing can be scored against
the same agent infrastructure.

## What's covered

- Listings posted on open aggregators (where the ToS permits crawling).
- Listings shared as a URL or as a screenshot/photo of an ad.
- Listings paired with a job verdict to detect [cross-links](../internal/rental/crosslink/cross.go).

## Rule pack highlights

| Code                     | Severity | Meaning |
|--------------------------|----------|---------|
| `deposit-wire-only`      | 5        | Wire-only deposits (Western Union, gift cards). |
| `no-viewing`             | 4        | No in-person or video viewing offered. |
| `reverse-image-hit`      | 4        | Listing photos appear on unrelated listings. |
| `rent-way-below-market`  | 3        | Rent < 50% of local median. |
| `urgency-pressure`       | 2        | Repeated "decide now" / "others waiting" phrasing. |

## Perceptual image hashing

[imagehash/phash](../internal/rental/imagehash/phash.go) computes a
64-bit pHash. Listings whose hash differs by ≤ 6 bits are treated as
the same photo for cross-listing detection.

## Combined verdict

When the same submission carries both a job offer and a rental
listing, [combined.Merge](../internal/rental/combined/verdict.go)
returns one band: the worst of the two, bumped one notch up if any
shared contacts are found.

## See also

- [THREAT-MODEL](THREAT-MODEL.md)
- [REGULATOR-INTEGRATION](REGULATOR-INTEGRATION.md)
