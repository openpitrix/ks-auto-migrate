module openpitrix.io/Jobs

go 1.13

require (
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	k8s.io/api v0.18.4 // indirect
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19 // indirect
	openpitrix.io/notification v0.2.2
	openpitrix.io/openpitrix v0.4.8
)

replace github.com/gocraft/dbr => github.com/gocraft/dbr v0.0.0-20180507214907-a0fd650918f6

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190423201726-d2cfbce3f3b0

replace openpitrix.io/openpitrix => openpitrix.io/openpitrix v0.4.9-0.20200617102217-10d232395f06

replace k8s.io/client-go => k8s.io/client-go v0.17.3
