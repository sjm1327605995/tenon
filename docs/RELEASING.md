# Releasing Tenon

Tenon remains pre-1.0. A version tag must not imply capabilities that are still listed as pending in `ROADMAP.md` or `docs/PLATFORM_SUPPORT.md`.

## Verification commands

Run the consolidated development check:

```sh
make release-check
```

The script checks that `backend/g3n/go.mod` has no development `replace` directive (evaluated with `GOWORK=off`), formatting, root build/vet/tests, interaction bounds, independent quick-start and versioned-plugin consumers (including invalid payload rejection and close/remount execution), race when CGO is available, G3N when available, and real BiDi, multiwindow, interactive-catalog, gradient, tilted 3D-card and projected-image captures when a display is available. `REQUIRE_WINDOW=1` also runs native G3N plugin, lifecycle, UI interaction and projected-image smokes when CGO/G3N are enabled. Every `SKIP` is explicitly not 1.0 release evidence.

For a release candidate, require every applicable gate:

```sh
REQUIRE_RACE=1 REQUIRE_G3N=1 REQUIRE_WINDOW=1 \
  TENON_RELEASE_WINDOW=1 make release-check
```

On headless Linux, install Xvfb and Gio's documented native dependencies. G3N additionally requires its documented compiler, OpenGL, OpenAL and Vorbis setup.

## Stable release checklist

- `ROADMAP.md` phases intended for the release are implemented and their acceptance checks pass; no required check is skipped.
- Linux, Windows and macOS CI passes, including configured real-window jobs. Inspect the screenshot artifacts.
- Root race tests and native G3N build/tests pass in suitable environments.
- A clean external consumer module builds the quick start using published module versions, not a repository-local `replace`.
- Root and `backend/g3n` modules use the same release version. Remove development-only local replace directives from published metadata.
- `CHANGELOG.md` contains the public additions, fixes, performance changes and migration notes.
- README, `docs/PLATFORM_SUPPORT.md` and `docs/PERFORMANCE.md` match the delivered code and actual verification.
- Pre-1.0 tags explicitly mention remaining API and platform limitations.

## Public API changes

`docs/public-api.txt` is the checked-in signature baseline for the root public packages, including exported struct fields but excluding implementation-only fields. `make api-check` generates a fresh snapshot and fails on additions, removals or signature changes; it is also part of `make release-check`.

When a public API change is intentional, run `make api-update`, review `docs/public-api.txt`, and include migration notes in `CHANGELOG.md`. Do not update the baseline merely to suppress an unintended difference.

## Versioning

Keep the root module and optional backend tags synchronized. Backend consumers should be able to select a single Tenon version without replacing it with a checkout path. Development pseudo-versions must use the canonical 12-character revision hash and have matching `go.sum` entries; short hashes are not valid for external consumers. During development, an ignored `go.work` maps the current pseudo-version to the local checkout; `backend/g3n/go.mod` itself contains no `replace` directive. Release metadata keeps synchronized versions instead of development-only replacements.

After tagging, verify both module paths from an external module:

```sh
go get github.com/sjm1327605995/tenon@<version>
go get github.com/sjm1327605995/tenon/backend/g3n@<version>
go build ./...
```

Do not mark a release complete until these consumer builds and the required native checks pass. The published G3N pseudo-version consumer currently demonstrates why synchronized versions are required: its older root requirement lacks APIs used by the backend. The workspace updates that requirement, but this fix must be committed and verified from the resulting published version.
