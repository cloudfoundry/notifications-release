module github.com/cloudfoundry/notifications-release/src/acceptance/v81

go 1.24.4

require (
	bitbucket.org/chrj/smtpd v0.0.0-20170817182725-9ddcdbda0f7a
	github.com/cloudfoundry/notifications-release/src/notifications/v81 v81.0.0
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/onsi/ginkgo/v2 v2.23.4
	github.com/onsi/gomega v1.37.0
	github.com/pivotal-cf-experimental/warrant v0.0.0-20211122194707-17385443920f
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cloudfoundry/notifications-release/src/notifications/v81 => ../notifications
