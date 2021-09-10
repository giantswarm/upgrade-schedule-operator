module github.com/giantswarm/upgrade-schedule-operator

go 1.16

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/giantswarm/apiextensions/v3 v3.33.0
	github.com/go-logr/logr v0.4.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/cluster-api v0.3.22
	sigs.k8s.io/controller-runtime v0.10.0
)

replace (
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2
)
