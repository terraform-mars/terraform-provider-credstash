# Release Process

## Prerequisites

### 1. Generate GPG Key

```bash
gpg --full-generate-key
```

Choose:
- RSA and RSA
- 4096 bits
- Key does not expire (or set expiration)
- Real name and email

### 2. Export GPG Key

```bash
# Get your key ID
gpg --list-secret-keys --keyid-format=long

# Export private key (armor format)
gpg --armor --export-secret-keys YOUR_KEY_ID
```

### 3. Add to GitHub Repository Secrets

Go to: `https://github.com/terraform-mars/terraform-provider-credstash/settings/secrets/actions`

Add two secrets:
- `GPG_PRIVATE_KEY`: Paste the full output from the export command (including BEGIN/END lines)
- `PASSPHRASE`: Your GPG key passphrase

### 4. Publish Public Key to Ubuntu Keyserver

```bash
gpg --keyserver keyserver.ubuntu.com --send-keys YOUR_KEY_ID
```

### 5. Add GPG Public Key to Terraform Registry

1. Go to: https://registry.terraform.io/settings/gpg-keys
2. Add your public GPG key
3. Associate it with your namespace

## Creating a Release

### Delete and Recreate the Current Tag

Since the existing 0.5.1 release has no assets:

```bash
# Delete local tag
git tag -d 0.5.1

# Delete remote tag
git push origin :refs/tags/0.5.1

# Delete the release on GitHub
# Go to https://github.com/terraform-mars/terraform-provider-credstash/releases
# and delete the 0.5.1 release

# Create new tag
git tag v0.5.1

# Push tag (will trigger GitHub Actions)
git push origin v0.5.1
```

### Future Releases

```bash
git tag v0.5.2
git push origin v0.5.2
```

The GitHub Action will automatically:
- Build binaries for all platforms (including darwin_arm64)
- Create checksums
- Sign with GPG
- Create GitHub release with all assets
