module github.com/cloudfoundry/notifications-release/src/acceptance/v81

go 1.25.6

require (
	bitbucket.org/chrj/smtpd v0.0.0-20170817182725-9ddcdbda0f7a
	github.com/cloudfoundry/notifications-release/src/notifications/v81 v81.0.0
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/onsi/ginkgo/v2 v2.28.1
	github.com/onsi/gomega v1.39.0
	github.com/pivotal-cf-experimental/warrant v0.0.0-20211122194707-17385443920f
)

require (
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20260115054156-294ebfa9ad83 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
)

replace github.com/cloudfoundry/notifications-release/src/notifications/v81 => ../notifications
