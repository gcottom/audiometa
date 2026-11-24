# AudioMeta v4 release sequence

| Order | Go module | Release | Dependencies |
| --- | --- | --- | --- |
| 1 | `github.com/gcottom/audiometa-core/v2` | v2.0.0 | Standard library |
| 2 | `github.com/gcottom/mp3meta/v2` | v2.0.0 | Core v2.0.0 |
| 2 | `github.com/gcottom/mp4meta/v2` | v2.0.0 | Core v2.0.0 |
| 2 | `github.com/gcottom/oggmeta/v2` | v2.0.0 | Core v2.0.0 |
| 2 | `github.com/gcottom/flacmeta/v2` | v2.0.0 | Core v2.0.0 |
| 3 | `github.com/gcottom/audiometa/v4` | v4.0.0 | Core and four adapters v2.0.0 |

Core and the adapters use root modules with the Go `/v2` major-version suffix. The SDK lives under `v4/`. The existing v3 SDK remains under `v3/`. Go can resolve each major version independently.

Publish Core and verify it through the public Go proxy, then publish the four adapters and verify their resolved module versions and checksums. Test the SDK from its own module with workspace resolution disabled before pushing the v4 tag. Verify the pushed SDK tag from a fresh external module using the public Go proxy and checksum database before publishing the SDK's GitHub release.

The local consumer test packages all six repositories into a temporary file proxy. The public mode uses a clean module cache and public resolution. See the verification report for commands and tested toolchains.

To adopt v4, replace the v3 import with `github.com/gcottom/audiometa/v4`. Applications that also import format adapters directly should update those imports to `/v2` where they want to exchange concrete tag types with the v4 SDK. The `OpenTag` usage pattern is available in v4. `OpenMetadata` adds typed fields, custom values, and structured metadata.
