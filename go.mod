module github.com/Octops/agones-broadcaster-http

go 1.14

require (
	agones.dev/agones v1.5.0
	github.com/Octops/agones-event-broadcaster v0.0.0-20200524110828-b366e7885189
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	k8s.io/client-go v11.0.1-0.20191029005444-8e4128053008+incompatible
)

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.3.0
