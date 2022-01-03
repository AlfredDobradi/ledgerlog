# CHANGELOG

<!--- next entry here -->

## 0.4.0
2022-01-03

### Features

- Implementing websocket handler (e40b4126790e6c3cab65958d988607032243d51c)
- Use connection pools (019f923a03b5ce6f3696d675bfed054aab0048ef)

### Fixes

- Fix post loop (97e5c25b67c48b8b736643deeb01b2b5800912ce)

## 0.3.0
2021-12-31

### Features

- Index page (80b7d09c35f4757b5a2529399bd4a70c524dbde0)

### Fixes

- Carriage return after signal (d1493f4b7f97bd78f2067129615179d2dbf379c2)
- Remove logging from GetPosts (4141eba57563faf1b53ef8dc32ef844c40fc7103)

## 0.2.0
2021-12-31

### Features

- Binary names (71523253eacd8dc22402125bc2992c2f8058b793)

## 0.1.0
2021-12-31

### Features

- gitignore (27b904a95775631d9349e3e72deed5a23717f423)
- Authorization and registration, Badger debug (56fb63f9165ffa0789c19e2f711e5a62f6f0f622)
- Client commands, send and show posts (6391b68741a18cabfe28f20084aa145f402af345)
- Configuration (521012f9ca598aab55a326015cedc6470c9810de)
- added postgres driver, barely (56472ed583f73b552307bd6a8e94e823d9fa4914)
- Sort pgx/crdbpgx, register user (bc2e10ca0d9c6569e45d0a842cd02ec41b5ababf)
- GetPublicKey method for crdb (182a3997d2b8cc0210ae5ae96763db5daa4293a9)
- Implement AddPost in crdb (47d3dd8ef5f7dd0ea64980f86cacb96418e6a0f1)
- Get user model when handling request (650e01cd69505ad8d74eff5bb29e3cee808cb03c)
- implement GetPosts (1bfc98709c0040077e53f41c57d71d97db996b0f)
- Add Magefile (bf0d8f0f80e6070e34c712710eb287ba3109b91d)
- Magefile error messages (be26cc99e346399322b0f75e92c667b67daead65)
- Client errors (3c28c65a1a8da40492c47f59d218d73d4b594581)
- Added CHANGELOG.md (9e3ea31b1f01cb07c2bfaca780bcfc316b7d9916)

### Fixes

- Move web service code (f1fbb7a660eea213611895dec9bbcb3b23437f3e)
- CLI commands (6f1f63870ff1f64711d1275754ade7a3c34b3438)
- Route constants, instance url in client (d0a0b8bfa2abb5369a3d588884f4bd880aae74af)
- DB connections, config usage (28946c41b3045c9a97a7fcd377f98882c5a116d1)
- Focus on crdb driver (f0d6ee9ea9812c548b86a4283f46ef68cfaebf43)
- Refactoring a bit of everything (b2431018b5659cdfa0835317756a58b0a17aa817)
- Use separate connections (c22dc5d07406f88168ea6ebc6470de1b0266d56e)
- split server and client (fbd545381d551b889b5a65854fe34e7c6027384a)
- Remove unused (51c0ef6fb10612aac836ac4e908fae79057126ee)
- New server has no error path (ccb20a38c551eb5e9b29df15a0e4066cda652c69)