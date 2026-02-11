# Release Process

This document describes how to create a new release of the HCP Terraform CLI.

## Prerequisites

- Push access to the GitHub repository
- All changes merged to `main` branch
- All CI checks passing
- CHANGELOG.md updated with release notes

## Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (x.0.0): Incompatible API changes
- **MINOR** version (0.x.0): New functionality in a backward compatible manner
- **PATCH** version (0.0.x): Backward compatible bug fixes

Examples:
- `v0.1.0` - Initial release
- `v0.2.0` - New features added
- `v0.2.1` - Bug fixes
- `v1.0.0` - First stable release

## Release Steps

### 1. Update Version Information

Update the version in `CHANGELOG.md`:

```markdown
## [0.2.0] - 2025-02-11

### Added
- New feature descriptions

### Changed
- Changed feature descriptions

### Fixed
- Bug fix descriptions
```

### 2. Commit Changes

```bash
git add CHANGELOG.md
git commit -m "chore: prepare for v0.2.0 release"
git push origin main
```

### 3. Create and Push Tag

```bash
# Create an annotated tag
git tag -a v0.2.0 -m "Release v0.2.0"

# Push the tag to GitHub
git push origin v0.2.0
```

### 4. Automated Release Process

Once the tag is pushed, GitHub Actions will automatically:

1. Run all tests
2. Build binaries for multiple platforms:
   - Linux (amd64, arm64, arm)
   - macOS (amd64, arm64)
   - Windows (amd64)
   - FreeBSD (amd64)
3. Create release archives (tar.gz for Unix, zip for Windows)
4. Generate SHA256 checksums
5. Create a GitHub release with:
   - Release notes from CHANGELOG.md
   - Compiled binaries attached
   - Checksum file

### 5. Verify Release

After the GitHub Action completes:

1. Go to https://github.com/thrashr888/hcptf-cli/releases
2. Verify the release appears with all binaries
3. Download and test a binary:

```bash
# Example: Download macOS binary
curl -LO https://github.com/thrashr888/hcptf-cli/releases/download/v0.2.0/hcptf_0.2.0_darwin_amd64.tar.gz

# Extract
tar -xzf hcptf_0.2.0_darwin_amd64.tar.gz

# Test
./hcptf version
```

### 6. Announce Release

Consider announcing the release:
- GitHub Discussions
- HashiCorp Community Forum
- Twitter/Social Media
- Internal channels

## Release Assets

Each release includes:

### Binaries

Platform-specific binaries in archives:
- `hcptf_VERSION_linux_amd64.tar.gz`
- `hcptf_VERSION_linux_arm64.tar.gz`
- `hcptf_VERSION_darwin_amd64.tar.gz` (macOS Intel)
- `hcptf_VERSION_darwin_arm64.tar.gz` (macOS Apple Silicon)
- `hcptf_VERSION_windows_amd64.zip`
- `hcptf_VERSION_freebsd_amd64.tar.gz`

### Documentation

Each archive includes:
- README.md
- LICENSE
- CHANGELOG.md
- docs/ directory with all guides

### Checksums

- `hcptf_VERSION_SHA256SUMS` - SHA256 checksums for all binaries

## Verifying Downloads

Users can verify downloads using the checksum file:

```bash
# Download checksum file
curl -LO https://github.com/thrashr888/hcptf-cli/releases/download/v0.2.0/hcptf_0.2.0_SHA256SUMS

# Verify specific binary
shasum -a 256 -c hcptf_0.2.0_SHA256SUMS 2>&1 | grep hcptf_0.2.0_darwin_amd64.tar.gz
```

## Hotfix Releases

For critical bug fixes that need immediate release:

1. Create a hotfix branch from the release tag:
   ```bash
   git checkout -b hotfix/v0.2.1 v0.2.0
   ```

2. Make the fix and commit:
   ```bash
   git commit -m "fix: critical bug description"
   ```

3. Update CHANGELOG.md

4. Merge to main and follow normal release process

## Pre-release Versions

For beta or release candidate versions:

```bash
# Create pre-release tag
git tag -a v0.3.0-beta.1 -m "Release v0.3.0-beta.1"
git push origin v0.3.0-beta.1
```

GoReleaser will automatically mark these as pre-releases on GitHub.

## Rolling Back a Release

If a release has critical issues:

1. Delete the tag locally and remotely:
   ```bash
   git tag -d v0.2.0
   git push origin :refs/tags/v0.2.0
   ```

2. Delete the GitHub release via the web interface

3. Fix the issue

4. Create a new patch release (e.g., v0.2.1)

## Local Testing

To test the release process locally without publishing:

```bash
# Install GoReleaser
brew install goreleaser

# Test release locally
goreleaser release --snapshot --clean

# Check output in dist/ directory
ls -la dist/
```

## Troubleshooting

### Release workflow fails

1. Check GitHub Actions logs
2. Verify all tests pass: `go test ./...`
3. Verify gofmt: `gofmt -l .`
4. Verify GoReleaser config: `goreleaser check`

### Missing binaries

Check `.goreleaser.yml` for platform configuration. Some platforms may be excluded if unsupported.

### Version not injected correctly

Verify ldflags in `.goreleaser.yml` match the variable names in `main.go`:
- `main.Version`
- `main.VersionPrerelease`

## Related Files

- `.goreleaser.yml` - GoReleaser configuration
- `.github/workflows/release.yml` - GitHub Actions release workflow
- `main.go` - Version variables
- `CHANGELOG.md` - Release notes

## References

- [GoReleaser Documentation](https://goreleaser.com)
- [Semantic Versioning](https://semver.org)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
