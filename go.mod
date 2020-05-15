module github.com/soundcloud/periskop

require (
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/aws/aws-sdk-go v1.30.7 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/gophercloud/gophercloud v0.10.0 // indirect
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/common v0.9.1
	github.com/prometheus/prometheus v1.8.2-0.20200213233353-b90be6f32a33
	github.com/soundcloud/periskop-go v0.0.0-20200417103225-bbd1a6d15b82
	google.golang.org/api v0.21.0 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/utils v0.0.0-20200411171748-3d5a2fe318e4 // indirect

)

go 1.13

replace github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v36.2.0+incompatible

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.0+incompatible
