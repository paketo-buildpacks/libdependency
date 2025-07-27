module github.com/paketo-buildpacks/libdependency

go 1.24.5

// This is required because of a breaking change in a newer version
replace github.com/ekzhu/minhash-lsh => github.com/ekzhu/minhash-lsh v0.0.0-20171225071031-5c06ee8586a1

// https://golang.testcontainers.org/quickstart/#2-install-testcontainers-for-go
replace github.com/docker/docker => github.com/docker/docker v20.10.3-0.20221013203545-33ab36d6b304+incompatible // 22.06 branch

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/Masterminds/semver/v3 v3.4.0
	github.com/anchore/packageurl-go v0.1.1-0.20250220190351-d62adb6e1115
	github.com/go-enry/go-license-detector/v4 v4.3.0
	github.com/onsi/gomega v1.37.0
	github.com/paketo-buildpacks/occam v0.13.3
	github.com/paketo-buildpacks/packit/v2 v2.22.0
	github.com/sclevine/spec v1.4.0
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/ForestEckhardt/freezer v0.0.12 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Microsoft/hcsshim v0.11.7 // indirect
	github.com/ProtonMail/go-crypto v1.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/containerd/containerd v1.7.27 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/cyphar/filepath-securejoin v0.4.1 // indirect
	github.com/dgryski/go-minhash v0.0.0-20190315135803-ad340ca03076 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v28.2.2+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/ekzhu/minhash-lsh v0.0.0-20190924033628-faac2c6342f8 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.2 // indirect
	github.com/go-git/go-git/v5 v5.16.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hhatto/gorst v0.0.0-20181029133204-ca9f730cac5b // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jdkato/prose v1.2.1 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/pjbgf/sha1cd v0.3.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/shogo82148/go-shuffle v1.0.1 // indirect
	github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af // indirect
	github.com/skeema/knownhosts v1.3.1 // indirect
	github.com/testcontainers/testcontainers-go v0.17.0 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
	gopkg.in/neurosnap/sentences.v1 v1.0.7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
