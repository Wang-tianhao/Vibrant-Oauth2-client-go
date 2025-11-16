# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-16

### Added
- Initial stable release of Vibrant OAuth2 Client for Go
- OAuth2 client credentials grant type support
- Automatic token caching mechanism
- Automatic token refresh with 60-second expiration buffer
- Thread-safe concurrent access with mutex synchronization
- Simple one-function API: `GetToken()`
- Environment variable configuration for credentials
  - `VIBRANT_CLIENT_ID`
  - `VIBRANT_CLIENT_SECRET`
- Vibrant API endpoint: `https://api.vibrant-wellness.com/v1/oauth2/token`
- Comprehensive documentation and examples
- No external dependencies (standard library only)

### Features
- `NewClient()` - Creates a new OAuth2 client from environment variables
- `GetToken()` - Returns a valid access token (cached or freshly fetched)
- `ClearCache()` - Clears cached token to force refresh
- `IsExpired()` - Checks if cached token has expired

### Documentation
- Complete README with installation and usage instructions
- Example application demonstrating all features
- Inline code documentation for all public APIs

[1.0.0]: https://github.com/Wang-tianhao/Vibrant-Oauth2-client-go/releases/tag/v1.0.0
