module github.com/giantswarm/upgrade-schedule-operator

go 1.16

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.16.0
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/cluster-api v0.3.22
	sigs.k8s.io/controller-runtime v0.10.0
)

replace (
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.24+incompatible
	github.com/coreos/etcd v3.3.13+incompatible => github.com/coreos/etcd v3.3.24+incompatible
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2
	golang/github.com/dgrijalva/jwt-go v3.2.0+incompatible => golang/github.com/dgrijalva/jwt-go v4.0.0-preview1+incompatible
)
