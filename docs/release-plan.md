# Candidate release plan

The five upstream default branches and their release tags remain release baselines. Local feature branches hold the AI-assisted changes. The shared module is local pending publication approval.

| Order | Go module | Proposed version | Depends on |
| --- | --- | --- | --- |
| 1 | github.com/gcottom/audiometa-core | v0.1.0 | Standard library |
| 2 | github.com/gcottom/mp3meta | v0.1.0 | Core v0.1.0 |
| 2 | github.com/gcottom/mp4meta | v0.1.0 | Core v0.1.0 |
| 2 | github.com/gcottom/oggmeta | v0.1.0 | Core v0.1.0 |
| 2 | github.com/gcottom/flacmeta | v0.1.0 | Core v0.1.0 |
| 3 | github.com/gcottom/audiometa/v3 | v3.1.0 | Core and four adapters v0.1.0 |

After approval, publish core source and its first tag, then the tested adapter versions, then AudioMeta. Release dependency hashes must resolve against those actual module artifacts. Repeat the external consumer test with a fresh module cache and the real public Go proxy before publishing AudioMeta. Keep existing release tags intact. The SDK still lives under `v3/`.

The development workspace uses local module replacements for proposed versions. The consumer verification script instead packages the exact trees into an isolated file proxy and resolves the candidate versions with a clean cache. It tests module packaging and minimum declared Go version separately from public distribution. The candidate artifacts must never be pushed from this offline test.

Existing `OpenTag` callers require no migration. Use `OpenMetadata` to adopt typed fields, custom fields, and precise artwork or metadata preservation, and handle the new convenience setters’ returned errors.
