# Shared metadata contract and AudioMeta v4 API

The unified SDK ships under the approved v4 module path. The v3 module retains its v3.0.5 source.

The SDK already depends on four independent adapter modules. Put the canonical identifiers, validation, typed values, representation registry and editing contract in a standard-library-only `github.com/gcottom/audiometa-core/v2` module. An SDK-owned contract would create an import cycle and duplicating definitions in adapters would create unequal validation and mappings. The core has no dependency on format libraries. Each adapter still parses its own container, manages offsets, and handles streaming and output errors.

Keep the `Tag` interface and original `OpenTag` shape in v4. Add `OpenMetadata` to expose canonical metadata through each adapter's equivalent opening function. The new API uses a tagged Go value with typed constructors and getters. One type-erased field dispatch is necessary for enumeration and heterogeneous transfer, and generic methods are unavailable in Go. Convenience methods use the same editor. Field values and custom values returned to callers are owned copies.

The registry identifies semantics and storage precedence. Adapters retain the original reader and record only requested edits. Save renders from the original source and replays edits on repeated saves. Unknown metadata and alternate encodings are retained where valid. Clear operations and a field's aliases have documented scope. Reject layouts for which the adapter cannot preserve opaque data and maintain playback integrity. A stream destination can be partially written on I/O failure.

The namespace is `audiometa/1/<field>` in ID3 descriptions and iTunes freeform names, under the reverse-DNS mean `io.github.gcottom.audiometa`. Vorbis-style comments use `AUDIOMETA_1_<FIELD>`. Custom values have a separate prefix and case-insensitive ASCII keys. Existing native and interoperable keys take precedence on read and write. Opaque native metadata is not implicitly transferred across formats.

Local development uses a multi-module Go workspace. Publish core first, then adapters, then the SDK. Keep release refs and proprietary metadata behavior independent of local testing.
