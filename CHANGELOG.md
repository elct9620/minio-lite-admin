# Changelog

## [0.1.1](https://github.com/elct9620/minio-lite-admin/compare/v0.1.0...v0.1.1) (2025-06-28)


### Bug Fixes

* remove incorrect package-name from release-please workflow ([d7d498f](https://github.com/elct9620/minio-lite-admin/commit/d7d498f3326dbf2470e2bba6e0d19e07c084de48))

## 0.1.0 (2025-06-28)


### Features

* add Chi HTTP router with middleware ([a9b8aab](https://github.com/elct9620/minio-lite-admin/commit/a9b8aab0b150c94d53334772606eac48065d50e1))
* add comprehensive tests for access keys creation feature ([2c84afb](https://github.com/elct9620/minio-lite-admin/commit/2c84afb94a2fcf477a6f6e725eab78dd6b30841b))
* add disk usage monitoring with optimized API integration ([7724106](https://github.com/elct9620/minio-lite-admin/commit/7724106294b6a1e396065fababc500f8ac89d227))
* add Docker Compose development setup with watch mode ([7b0438b](https://github.com/elct9620/minio-lite-admin/commit/7b0438b9b2004be3ecaca3240f54a2ef3caeb79d))
* add GitHub Actions workflow for Docker builds and GHCR publishing ([850e319](https://github.com/elct9620/minio-lite-admin/commit/850e3198b92af38cef4161b3060923002dfa20e0))
* add graceful shutdown to HTTP server ([7cd2234](https://github.com/elct9620/minio-lite-admin/commit/7cd2234cb331fda823935bb10a796f6dbf115b85))
* add production Dockerfile with multi-stage build ([0e3a780](https://github.com/elct9620/minio-lite-admin/commit/0e3a7800a10f3686b20ac3c15e2e164e786c50b8))
* add release-please automation for semantic versioning ([b9070a5](https://github.com/elct9620/minio-lite-admin/commit/b9070a54008525ca6e682290efef7136cfadf22e))
* add sidebar navigation to dashboard with focused menu items ([9edc169](https://github.com/elct9620/minio-lite-admin/commit/9edc1693726d0527474da144489cce03e5becffb))
* add SPA routing support for Vue Router client-side navigation ([b762476](https://github.com/elct9620/minio-lite-admin/commit/b7624769937bd932405846f1546cec678bdcb5fc))
* add Vue 3 + TypeScript frontend with Vite ([6629e17](https://github.com/elct9620/minio-lite-admin/commit/6629e17e0f303d2d140437757c12bf50e5037b37))
* configure release-please with issue permissions ([9b07cb4](https://github.com/elct9620/minio-lite-admin/commit/9b07cb4833000242c432a21ff9f0793d716737f8))
* implement access key and secret display after creation and rotation ([856cd5c](https://github.com/elct9620/minio-lite-admin/commit/856cd5c1a16eb102e4c40313161aff520213a51a))
* implement access keys API with comprehensive MinIO integration ([1b6fe29](https://github.com/elct9620/minio-lite-admin/commit/1b6fe29bc30005232e2c4edab24721f179f7073b))
* implement comprehensive edit access keys feature ([18660a2](https://github.com/elct9620/minio-lite-admin/commit/18660a27b86593f80b15a9fb8dee6d5a33cf8f28))
* implement conditional asset embedding with build tags ([dca52a4](https://github.com/elct9620/minio-lite-admin/commit/dca52a451864f9605877e4b3338d99e59db974c5))
* implement create access keys feature with comprehensive service account support ([6a70ff5](https://github.com/elct9620/minio-lite-admin/commit/6a70ff51fb250c203801c26ed621bf889d06baa8))
* implement delete access keys feature with comprehensive confirmation workflow ([df81f13](https://github.com/elct9620/minio-lite-admin/commit/df81f13c5974509d50edf07fbbda75e336f45471))
* implement disk usage dashboard with real-time data display ([ace2443](https://github.com/elct9620/minio-lite-admin/commit/ace2443950a7490fc4092803b8fbd377278ac1e9))
* implement Vue Router 4 with navigation for MinIO admin pages ([1e6d578](https://github.com/elct9620/minio-lite-admin/commit/1e6d5785095eff0b1488b331e48475f7d08f169e))
* integrate Go backend with Vite using olivere/vite ([e0bfbfa](https://github.com/elct9620/minio-lite-admin/commit/e0bfbfa182fdb64cc09f6b0e77adb1baf35530aa))
* integrate MinIO Admin SDK with service layer architecture ([5e500d4](https://github.com/elct9620/minio-lite-admin/commit/5e500d42cc5168274fb302f0b521410ef7c9a776))
* integrate TailwindCSS 4 with modern dashboard UI ([72bb67a](https://github.com/elct9620/minio-lite-admin/commit/72bb67a453189ab1e09312fb994d341678768ddc))
* integrate zerolog for structured logging ([c3781e1](https://github.com/elct9620/minio-lite-admin/commit/c3781e18b4ddc25ef64fd34a03f71cf0b5c07b39))
* migrate all SVG icons to Heroicons ([dc303ae](https://github.com/elct9620/minio-lite-admin/commit/dc303ae1f25a10d1ac6e9b3e8d4f481ba05f10cd))


### Bug Fixes

* copy entire source directory in Dockerfile to resolve build failure ([3af7bd6](https://github.com/elct9620/minio-lite-admin/commit/3af7bd6ff5efbedcce657b0b350696e90a6db343))
* correct permissions in release-please workflow ([3e92993](https://github.com/elct9620/minio-lite-admin/commit/3e92993e03823e539ddcbbc56024391cd63e1fa6))
* improve state isolation between create and edit access key modals ([0920ff4](https://github.com/elct9620/minio-lite-admin/commit/0920ff4278bb95d9b5a5aeb0e95cc79153b3b3bf))
* resolve golangci-lint errors to pass CI pipeline ([aab73e5](https://github.com/elct9620/minio-lite-admin/commit/aab73e5d97e5f581461de4ead87e3d31337a8f10))
* resolve infinite loading state when no access keys exist ([6d19c79](https://github.com/elct9620/minio-lite-admin/commit/6d19c799053355460ef202e1762669a9d96025b7))
* update GitHub repository URLs to use correct owner elct9620 ([f6e5848](https://github.com/elct9620/minio-lite-admin/commit/f6e584883ff24bacc96fbb6833e2abeb89bb8c0d))
* use Unix timestamps for expiration dates instead of RFC3339 strings ([28c1950](https://github.com/elct9620/minio-lite-admin/commit/28c1950c8462630094dc7127c410f7a099e75be9))


### Miscellaneous Chores

* initialize release versioning ([7fdbf99](https://github.com/elct9620/minio-lite-admin/commit/7fdbf99e96e28108f7c42e27e78400a43014dc58))
