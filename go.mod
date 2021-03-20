module github.com/openpitrix/ks-auto-migrate

go 1.13

require (
	github.com/aws/aws-sdk-go v1.25.21
	github.com/golang/protobuf v1.3.2
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	openpitrix.io/notification v0.2.2
	openpitrix.io/openpitrix v0.4.8
)

replace github.com/gocraft/dbr => github.com/gocraft/dbr v0.0.0-20180507214907-a0fd650918f6

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190423201726-d2cfbce3f3b0

replace openpitrix.io/openpitrix => openpitrix.io/openpitrix v0.4.9-0.20200617102217-10d232395f06

replace k8s.io/client-go => k8s.io/client-go v0.17.3
