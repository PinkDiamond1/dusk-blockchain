# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Make BlockGasLimit configurable [#1418]
- Adapt to the new and improved reward strategy [#1320]
- Include BlockGenerator reward into the header [#1368]
- Search txs up to 10.000 blocks [#1401]
- Mempool updates its state at start-up automatically [#1258]
- Mempool discards any transaction with repeated nullifier [#1388]
- Detect any occurrences of missed fallback procedure [#1413] 

### Changed
- Extraction from mempool consider Gas expenditure estimation [#1421]
- Add GasLimit to the Block's header [#1416]
- Reduce wire messages needed for propagating an Agreement vote
- Improve fallback procedure to revert multiple ephemeral blocks [#1343]
- Disable provisioners records storing in api.db  [#1361]
- Optimize VST calls of a Provisioner per a consensus iteration [#1371]
- Call VerifyStateTransition in Reduction 2 only if it is a committee member [#1357]
- Drop topics.Block messages with invalid header.hash value [#1425]
- Request candidate block from arbitrary active nodes [#1359]
- Allow verifyFn to report errors correctly #1436
- Clone message.Agreement on notifying the consumer [#1433]


### Security
- Ensure candidate block hash is equal to the BlockHash of the msg.header [#1364]
- Check quorum per each reduction step on agreement msg verification [#1373]
- Extend reduction 1 verification procedure with VST call [#1358]

### Fixed
- Default configuration values loading [#1419] 

[#1421]: (https://github.com/dusk-network/dusk-blockchain/issues/1421)
[#1419]: (https://github.com/dusk-network/dusk-blockchain/issues/1419)
[#1418]: (https://github.com/dusk-network/dusk-blockchain/issues/1418)
[#1416]: (https://github.com/dusk-network/dusk-blockchain/issues/1416)
[#1401]: (https://github.com/dusk-network/dusk-blockchain/issues/1401)
[#1388]: (https://github.com/dusk-network/dusk-blockchain/issues/1388)
[#1373]: (https://github.com/dusk-network/dusk-blockchain/issues/1373)
[#1371]: (https://github.com/dusk-network/dusk-blockchain/issues/1371)
[#1368]: (https://github.com/dusk-network/dusk-blockchain/issues/1368)
[#1364]: (https://github.com/dusk-network/dusk-blockchain/issues/1364)
[#1361]: (https://github.com/dusk-network/dusk-blockchain/issues/1361)
[#1358]: (https://github.com/dusk-network/dusk-blockchain/issues/1358)
[#1357]: (https://github.com/dusk-network/dusk-blockchain/issues/1357)
[#1343]: (https://github.com/dusk-network/dusk-blockchain/issues/1343)
[#1320]: (https://github.com/dusk-network/dusk-blockchain/issues/1320)
[#1258]: (https://github.com/dusk-network/dusk-blockchain/issues/1258)


<!-- Releases -->

[Unreleased]: https://github.com/dusk-network/dusk-blockchain/compare/daybreak-20220321...HEAD
