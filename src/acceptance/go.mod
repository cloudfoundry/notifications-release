module github.com/cloudfoundry/notifications-release/src/acceptance/v81

go 1.25.0

require (
	bitbucket.org/chrj/smtpd v0.0.0-20170817182725-9ddcdbda0f7a
	github.com/cloudfoundry/notifications-release/src/notifications/v81 v81.0.0
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/onsi/ginkgo/v2 v2.25.1
	github.com/onsi/gomega v1.38.2
	github.com/pivotal-cf-experimental/warrant v0.0.0-20211122194707-17385443920f
)

require (
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
)

replace github.com/cloudfoundry/notifications-release/src/notifications/v81 => ../notifications
