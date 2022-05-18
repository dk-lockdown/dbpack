# ChangeLog

# v0.1.0

### Bug Fixes

* should use db connection rather than tx connection exec sql request ([#8](https://github.com/cectc/dbpack/pull/8)) ([7e2b42d](https://github.com/cectc/dbpack/commit/52d78cab0bc414d92a5c59230f2827c8332c2bde))
* add terminationDrainDuration config ([#4](https://github.com/cectc/dbpack/issues/4)) [6604ce8](https://github.com/cectc/dbpack/commit/5c607d48d1149218cff3988dcb00d83da571a561))
* when receive ComQuit request, should return connection ([#51](https://github.com/cectc/dbpack/pull/51)) [627adc2](https://github.com/cectc/dbpack/commit/e8f07086ccf76a7112f00512e3ed3f6e94aff410))

### Features

* distributed transaction support etcd watch ([#11](https://github.com/cectc/dbpack/pull/11)) ([ce10990](https://github.com/cectc/dbpack/commit/e9910501e32d23741f99f5fe9ece1077ba1b348c))
* support leader election, only leader can commit and rollback ([#19](https://github.com/cectc/dbpack/pull/19)) ([b89c672](https://github.com/cectc/dbpack/commit/d7ab60b6ed5547f1bc9a6c426e1fb9ee21d6f4f3))
* support tcc branch commit & rollback ([#12](https://github.com/cectc/dbpack/issues/12)) ([c0bfdf9](https://github.com/cectc/dbpack/commit/feab7aefe819bf3217363994c67515b887f8adb9))
* add prometheus metric ([#25](https://github.com/cectc/dbpack/issues/25)) ([627adc2](https://github.com/cectc/dbpack/commit/627adc2ced9da499e6b658f718b23417e7df9903))
* support GlobalLock hint ([#14](https://github.com/cectc/dbpack/issues/14)) ([8369f8f](https://github.com/cectc/dbpack/commit/5c7c96797539943ed75495d1cfa92f6094ff548e))

### Changes

* update branch session process logic ([#17](https://github.com/cectc/dbpack/pull/17)) ([c6a6626](https://github.com/cectc/dbpack/commit/06d624511c65a379e73dae91c2be4fb3785b9bf0))
