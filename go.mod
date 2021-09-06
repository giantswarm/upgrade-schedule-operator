module github.com/giantswarm/upgrade-schedule-operator

go 1.16

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489 // indirect
	go.uber.org/tools v0.0.0-20190618225709-2cfd321de3ee // indirect
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/cluster-api v0.3.22
	sigs.k8s.io/controller-runtime v0.10.0
)
