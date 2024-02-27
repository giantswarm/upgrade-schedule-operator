module github.com/giantswarm/upgrade-schedule-operator

go 1.16

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/giantswarm/k8smetadata v0.24.0
	github.com/go-logr/logr v0.4.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.18.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/text v0.14.0
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/cluster-api v1.0.4
	sigs.k8s.io/controller-runtime v0.10.3
)

replace (
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.27+incompatible
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2
)
