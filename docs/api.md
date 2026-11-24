# Unified metadata in AudioMeta v3

Use `OpenMetadata` for a format-independent editor. Its source must implement `io.ReadSeeker`. Reopen the saved output to inspect the actual on-disk result. The SDK and its adapters leave caller-provided handles open.

```go
tag, err := audiometa.OpenMetadata(input)
if err != nil {
    return err
}
if err := tag.SetTitle("Example"); err != nil {
    return err
}
if err := tag.SetField(audiometa.Artist,
    audiometa.TextValue("First artist"),
    audiometa.TextValue("Second artist"),
); err != nil {
    return err
}
if err := tag.SetISRC("USABC2600001"); err != nil {
    return err
}
if err := tag.SetCustom("COLLECTION", "Favorites", "Archive"); err != nil {
    return err
}
return tag.Save(output)
```

`GetField` reads the first value and reports presence. `Get` returns the full ordered list. `SetField` replaces a list, `AppendField` appends to a repeatable field, `RemoveValue` removes one matching occurrence, and `RemoveField` removes the field. `Populated` returns a snapshot of the interpreted fields. The custom equivalents are `GetCustom`, `SetCustom`, `AppendCustom`, `RemoveCustomValue`, `RemoveCustom`, `Custom` and `ClearCustom`. Returned maps, slices and structured values are caller-owned.

Set an empty text value to store an explicit empty, or remove the field to make it absent. The same distinction applies to zero integers and false booleans. Track/disc numbers and totals are independent. Partial dates preserve the supplied precision. Use the separate recording, release and original release identifiers when known. Artist lists are ordered and independent of literal artist credits.

`OpenTag` retains the v3 `Tag` contract and existing concrete adapter types. The new convenience setters return validation errors. Use `OpenMetadata` when you need multi-value reads, precise removals, transfer, or the preservation guarantees documented here.

## Transfer

Iterate over `source.Populated()`. Check `audiometa.Lookup(field)` and set mutable values on the destination with `destination.SetField(field, values...)`. Copy the `Custom()` snapshot with `SetCustom`. Each save operates on that destination's audio. Preserve unknown metadata in its original container. Unknown payloads do not have a portable interpretation.

## Preservation and errors

- Reads resolve standardized storage and conventional aliases before the fallback, with documented exceptions when a native representation cannot represent an exact value. Some codecs write a native compatibility projection with a corresponding lossless extension. An unknown or invalid entry remains in its container unless an explicit clear can remove it safely.
- A set replaces the written slot and can retain competing entries as opaque or lower-precedence aliases. Removing a field removes its recognized representations to prevent older entries from reappearing. When multiple fields share one physical slot, such as track number and total, clearing one retains the other when representable.
- Custom fields are namespaced. Calling `SetCustom("TITLE", ...)` leaves a native title alone. Arbitrary pre-existing frame, atom and comment names outside the registry are preserved in their original representation. They have no portable edit operation.
- `ClearAll(true)` clears interpreted mutable canonical values and namespaced custom values and keeps opaque descriptive metadata. `ClearAll(false)` requests removal of opaque descriptive metadata too. Technical data, uncertain application blocks and playback-critical values stay intact. Unsafe clear operations fail explicitly.
- The editor reports validation failures via `ErrInvalidValue`, unknown fields via `ErrInvalidField`, and edits to immutable fields via `ErrReadOnly`. A supported semantic value can still fail at Save with `ErrUnsupported` for a specific container layout, representation limit or conflicting protected entry.
- Use a separate output. A stream save may leave partial bytes on error. The SDK makes no atomic-replacement promise. Keep the original input unchanged for repeated saves. The same editor and source require external synchronization.

## Fallback storage

| Format | Canonical field | Custom key |
| --- | --- | --- |
| ID3v2 | TXXX (TXX in v2.2), description `audiometa/1/<field>` | Description `audiometa/1/custom/<KEY>` |
| MP4 | iTunes freeform item with mean `io.github.gcottom.audiometa` and name `audiometa/1/<field>` | Same mean, name `audiometa/1/custom/<KEY>` |
| FLAC, Vorbis, Opus | Comment field `AUDIOMETA_1_<FIELD>` | Comment field `AUDIOMETA_1_CUSTOM_<KEY>` |

Custom keys are ASCII, case-insensitive, and contain one to 128 letters, digits, periods, underscores, or hyphens, starting with a letter or digit. Repeated values retain order. Fallbacks are lossless when read through AudioMeta. External tools may show the opaque data. They need audio metadata support for the canonical field to interpret it.
