# Changelog

## [0.9.1](https://github.com/nuvla/api-client-go/compare/v0.9.0...v0.9.1) (2024-11-05)


### Bug Fixes

* Fixes IRS on V2 ([abf2601](https://github.com/nuvla/api-client-go/commit/abf2601049f4a89f1f93e6a574e0b22f43e699ab))

## [0.9.0](https://github.com/nuvla/api-client-go/compare/v0.8.0...v0.9.0) (2024-10-04)


### Features

* Add context to http requests ([#51](https://github.com/nuvla/api-client-go/issues/51)) ([358b88f](https://github.com/nuvla/api-client-go/commit/358b88f48468f799dcbc020bf29741fe893c29aa))

## [0.8.0](https://github.com/nuvla/api-client-go/compare/v0.7.16...v0.8.0) (2024-09-25)


### Features

* add support for sending JSON Patch ([#48](https://github.com/nuvla/api-client-go/issues/48)) ([7302d29](https://github.com/nuvla/api-client-go/commit/7302d293817c04874281d09489db84bc966dcf51))

## [0.7.16](https://github.com/nuvla/api-client-go/compare/v0.7.15...v0.7.16) (2024-09-24)


### Bug Fixes

* Remove extra logging verbosity ([b3800ea](https://github.com/nuvla/api-client-go/commit/b3800ea3ec1bda8c73c0cd980dc68389302e31c7))


### Minor Changes

* Add support for IRS in NuvlaEdge client ([becda24](https://github.com/nuvla/api-client-go/commit/becda243781af65288266ac04c8b6754f66edb07))

## [0.7.15](https://github.com/nuvla/api-client-go/compare/v0.7.14...v0.7.15) (2024-09-16)


### Bug Fixes

* Add re-authentication to client ([e8935e4](https://github.com/nuvla/api-client-go/commit/e8935e4981efa07d86362feefc14203b8bdb16ba))


### Minor Changes

* Add credential encoding to nuvlaedge client freeze and load ([1c2b7be](https://github.com/nuvla/api-client-go/commit/1c2b7be37c698f380b8187d81ecb0b86b773da99))

## [0.7.14](https://github.com/nuvla/api-client-go/compare/v0.7.13...v0.7.14) (2024-09-13)


### Bug Fixes

* Deployment parameter bug ([989fdf5](https://github.com/nuvla/api-client-go/commit/989fdf5c2843ca7f41f9acbdf43cd7ef6e9f5497))

## [0.7.13](https://github.com/nuvla/api-client-go/compare/v0.7.12...v0.7.13) (2024-09-12)


### Bug Fixes

* Nil value interface checker ([832dd47](https://github.com/nuvla/api-client-go/commit/832dd474672d61e4d5c5e99167c4fd2140cf5460))

## [0.7.12](https://github.com/nuvla/api-client-go/compare/v0.7.11...v0.7.12) (2024-09-02)


### Bug Fixes

* Deployment parameter wrong sent key ([79a88a2](https://github.com/nuvla/api-client-go/commit/79a88a27af7b86408c9e312d5a0b7d35ef8db889))
* Remove redundant LogIn and fix credential parsing in NuvlaEdge client ([2eeab86](https://github.com/nuvla/api-client-go/commit/2eeab865f03150518b55066c473286a0139b02c3))

## [0.7.11](https://github.com/nuvla/api-client-go/compare/v0.7.10...v0.7.11) (2024-07-16)


### Bug Fixes

* publish sanitise endpoint function ([539379e](https://github.com/nuvla/api-client-go/commit/539379e6dd9c72437520bd48dba3670bf34bb62e))
* Reduce unnecessary verbosity in NuvlaEdge client ([769d6f9](https://github.com/nuvla/api-client-go/commit/769d6f92551a3ae4284622b649a25159414747d3))
* unclosed response body after log in ([0bbe621](https://github.com/nuvla/api-client-go/commit/0bbe621c824b8e17297cdccd2455324184c395cb))

## [0.7.10](https://github.com/nuvla/api-client-go/compare/v0.7.9...v0.7.10) (2024-07-05)


### Minor Changes

* Add SetFailedState and general job update function to job client ([751498a](https://github.com/nuvla/api-client-go/commit/751498a4d819f5f29c329a6ed97ab4a6244efa13))

## [0.7.9](https://github.com/nuvla/api-client-go/compare/v0.7.8...v0.7.9) (2024-07-01)


### Bug Fixes

* Fix credential reading from NuvlaEdge activation ([ebb0503](https://github.com/nuvla/api-client-go/commit/ebb0503e981d53cdb0c29a846e36d53c8ae3c1f2))

## [0.7.8](https://github.com/nuvla/api-client-go/compare/v0.7.7...v0.7.8) (2024-07-01)


### Bug Fixes

* Align json tags to allow bidirectional compatibility of stored sessions with python implementation ([#38](https://github.com/nuvla/api-client-go/issues/38)) ([3939da4](https://github.com/nuvla/api-client-go/commit/3939da4870f8f9e167b141d1a3c872cf5e08ca35))

## [0.7.7](https://github.com/nuvla/api-client-go/compare/v0.7.6...v0.7.7) (2024-06-25)


### Bug Fixes

* Fix non closed http.request bodies ([#36](https://github.com/nuvla/api-client-go/issues/36)) ([7f8d49f](https://github.com/nuvla/api-client-go/commit/7f8d49f4cd3abf56651166c54d6fa2f77d519851))

## [0.7.6](https://github.com/nuvla/api-client-go/compare/v0.7.5...v0.7.6) (2024-06-12)


### Bug Fixes

* Fix nuvlaedge resource location field typing bug ([e2017fe](https://github.com/nuvla/api-client-go/commit/e2017fe822f61350554ff87fdd3b25e42eab3f28))


### Minor Changes

* Add endpoint sanitise to endpoint in session.go ([984abc1](https://github.com/nuvla/api-client-go/commit/984abc103780fa729ccda2a3a1f048bb18e9d471))

## [0.7.5](https://github.com/nuvla/api-client-go/compare/v0.7.4...v0.7.5) (2024-06-11)


### Bug Fixes

* Fix Inferred location type in NuvlaEdgeResource ([764823d](https://github.com/nuvla/api-client-go/commit/764823d7c4e4d3a293b7f7832dd5bb8ef6b3f80d))

## [0.7.4](https://github.com/nuvla/api-client-go/compare/v0.7.3...v0.7.4) (2024-06-07)


### Minor Changes

* Add Selective resource update and resource getter to nuvlaedge client ([bde055f](https://github.com/nuvla/api-client-go/commit/bde055feb6eeb934a574804c950a597a6f97258f))

## [0.7.3](https://github.com/nuvla/api-client-go/compare/v0.7.2...v0.7.3) (2024-06-04)


### Minor Changes

* Add Getter for APK key-secret on job client ([#30](https://github.com/nuvla/api-client-go/issues/30)) ([347e75b](https://github.com/nuvla/api-client-go/commit/347e75b99144fea868e3ac6ec7eb2c461d4510c7))

## [0.7.2](https://github.com/nuvla/api-client-go/compare/v0.7.1...v0.7.2) (2024-05-29)


### Bug Fixes

* Naming bug on Module Environmental Variables ([9318d6d](https://github.com/nuvla/api-client-go/commit/9318d6d8f094320778c0b32d5256a361e6bd36b1))

## [0.7.1](https://github.com/nuvla/api-client-go/compare/v0.7.0...v0.7.1) (2024-05-28)


### Minor Changes

* Add session refresh from deployment credentials to deployment client ([ccd3fca](https://github.com/nuvla/api-client-go/commit/ccd3fca2ff3cb1e8a9215d917e40bd00afb4f222))

## [0.7.0](https://github.com/nuvla/api-client-go/compare/v0.6.0...v0.7.0) (2024-05-28)


### Features

* Add deployment parameter handler to deployment client ([#27](https://github.com/nuvla/api-client-go/issues/27)) ([d89804a](https://github.com/nuvla/api-client-go/commit/d89804a77e3a7615b22de431725c6b8da1f22eb6))


### Documentation

* Update NuvlaClient operations available ([24802b2](https://github.com/nuvla/api-client-go/commit/24802b2b01b391a001255820e6c94bd57ffe4127))

## [0.6.0](https://github.com/nuvla/api-client-go/compare/v0.5.2...v0.6.0) (2024-05-17)


### Features

* add Search operation to base client ([5fdeb6a](https://github.com/nuvla/api-client-go/commit/5fdeb6ab9129bb20ae0ed25b007af3dcada9ddcc))


### Minor Changes

* Add a function to parse json or data and encode the headers accordingly ([#24](https://github.com/nuvla/api-client-go/issues/24)) ([941cda8](https://github.com/nuvla/api-client-go/commit/941cda8caeb24d5af3a08cbf601426174d3e33e5))


### Code Refactoring

* Created a custom package to group all the resource types ([#23](https://github.com/nuvla/api-client-go/issues/23)) ([d0c43ef](https://github.com/nuvla/api-client-go/commit/d0c43efec1c737e01b728d49d7b9749ca042fae9))
* Remove resources from clients and point to resources package ([#25](https://github.com/nuvla/api-client-go/issues/25)) ([4ba7dc4](https://github.com/nuvla/api-client-go/commit/4ba7dc4279f16a79bbcedf2f00607ad753058be2))


### Continuous Integration

* Add minor and patch sections for minor features or changes to bump only patch ([604f346](https://github.com/nuvla/api-client-go/commit/604f346ff7921a0dd440220e07b28d51b8618a6b))

## [0.5.2](https://github.com/nuvla/api-client-go/compare/v0.5.1...v0.5.2) (2024-04-29)


### Bug Fixes

* reduce unnecessary verbosity ([aa40ffd](https://github.com/nuvla/api-client-go/commit/aa40ffd34ea6dbebd7332fb58c77324ade601f0b))
* safe check for null pointer exception on NuvlaEdge client ([ad65b49](https://github.com/nuvla/api-client-go/commit/ad65b49676ef6c68076a438e948001439adc496c))

## [0.5.1](https://github.com/nuvla/api-client-go/compare/v0.5.0...v0.5.1) (2024-04-29)


### Bug Fixes

* allow empty NuvlaID structs when ID is empty string ([e5afa5f](https://github.com/nuvla/api-client-go/commit/e5afa5f88db2bc73ed7dc60e5630ece87e9b5ac3))
* convert NuvlaID into pointers in NuvlaEdge client ([6af22da](https://github.com/nuvla/api-client-go/commit/6af22dab57b8b8ccc20f509165b6d42fb8272672))

## [0.5.0](https://github.com/nuvla/api-client-go/compare/v0.4.1...v0.5.0) (2024-04-26)


### Features

* add client freeze capabilities to NuvlaEdgeClient ([#19](https://github.com/nuvla/api-client-go/issues/19)) ([96caabb](https://github.com/nuvla/api-client-go/commit/96caabb3cbc5da3541eaf87dba4cfcca0efbd068))


### Documentation

* update README clients implementation status ([3d5d485](https://github.com/nuvla/api-client-go/commit/3d5d4858945007d05416fe2a0fc30786195633b2))

## [0.4.1](https://github.com/nuvla/api-client-go/compare/v0.4.0...v0.4.1) (2024-04-08)


### Bug Fixes

* change deployment resource module definition ([8e2fcfd](https://github.com/nuvla/api-client-go/commit/8e2fcfd0ab42423ceef9e312ef030df50a9ba53a))

## [0.4.0](https://github.com/nuvla/api-client-go/compare/v0.3.1...v0.4.0) (2024-03-26)


### Features

* add deployment client support ([#15](https://github.com/nuvla/api-client-go/issues/15)) ([c19d358](https://github.com/nuvla/api-client-go/commit/c19d358e09b4eb4e1e8e55eb37b6211444340ff7))

## [0.3.1](https://github.com/nuvla/api-client-go/compare/v0.3.0...v0.3.1) (2024-03-26)


### Bug Fixes

* Fix deployment client lack de-sync ([72c0230](https://github.com/nuvla/api-client-go/commit/72c02300856d9056fe1aafa25d5ccc25bc31988e))


### Continuous Integration

* Fix Ci changelog-notes sections ([feb26d9](https://github.com/nuvla/api-client-go/commit/feb26d9252abec2ebc78e2a4134aa6dce6918b31))

## [0.3.0](https://github.com/nuvla/api-client-go/compare/v0.2.0...v0.3.0) (2024-03-26)


### Features

* Add job client support ([#11](https://github.com/nuvla/api-client-go/issues/11)) ([1029def](https://github.com/nuvla/api-client-go/commit/1029def4bf17bc290409b961106ce98d40ab67dd))
* Add options pattern to client creation ([1029def](https://github.com/nuvla/api-client-go/commit/1029def4bf17bc290409b961106ce98d40ab67dd))


### Bug Fixes

* Bug on edit operation in base NuvlaClient ([1029def](https://github.com/nuvla/api-client-go/commit/1029def4bf17bc290409b961106ce98d40ab67dd))

## [0.2.0](https://github.com/nuvla/api-client-go/compare/v0.1.0...v0.2.0) (2024-03-22)


### Features

* add Options pattern to clients ([#10](https://github.com/nuvla/api-client-go/issues/10)) ([ebe5fde](https://github.com/nuvla/api-client-go/commit/ebe5fdea1b1c9dd89140760df1212e366e5095db))
* add user client support ([#8](https://github.com/nuvla/api-client-go/issues/8)) ([4a86845](https://github.com/nuvla/api-client-go/commit/4a86845c0947099c7117db98c479b83fe22986c7))
* implement composition patterns for session&gt;client>clients ([ebe5fde](https://github.com/nuvla/api-client-go/commit/ebe5fdea1b1c9dd89140760df1212e366e5095db))

## [0.1.0](https://github.com/nuvla/api-client-go/compare/v0.1.0...v0.1.0) (2024-03-05)


### Features

* add compress option for requests ([386c5c6](https://github.com/nuvla/api-client-go/commit/386c5c629c2877a465fee4ccd17ebde89882fee4))


### Bug Fixes

* fix module location in go.mod to point to GH repo ([3521a23](https://github.com/nuvla/api-client-go/commit/3521a2366f7335b9f722b815572f6c3649e93f2b))
* remove unused Client interface ([e3ccc93](https://github.com/nuvla/api-client-go/commit/e3ccc93f8e74d68e0434c8e4557eb783e084892b))
* solve release issues ([#1](https://github.com/nuvla/api-client-go/issues/1)) ([eaf68b9](https://github.com/nuvla/api-client-go/commit/eaf68b95ebd70d6152e4f2b72a708b694d8c4566))

## 0.1.0 (2024-03-05)


### Features

* add compress option for requests ([386c5c6](https://github.com/nuvla/api-client-go/commit/386c5c629c2877a465fee4ccd17ebde89882fee4))


### Bug Fixes

* fix module location in go.mod to point to GH repo ([3521a23](https://github.com/nuvla/api-client-go/commit/3521a2366f7335b9f722b815572f6c3649e93f2b))
* remove unused Client interface ([e3ccc93](https://github.com/nuvla/api-client-go/commit/e3ccc93f8e74d68e0434c8e4557eb783e084892b))
* solve release issues ([#1](https://github.com/nuvla/api-client-go/issues/1)) ([eaf68b9](https://github.com/nuvla/api-client-go/commit/eaf68b95ebd70d6152e4f2b72a708b694d8c4566))

## Changelog
