module github.com/Octops/agones-broadcaster-http

go 1.14

require (
	agones.dev/agones v1.6.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	k8s.io/client-go v0.17.2
	github.com/Octops/agones-event-broadcaster v0.1.6-alpha
	sigs.k8s.io/controller-runtime v0.5.4 // indirect
)
