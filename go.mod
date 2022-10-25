module github.com/SecurityBrewery/catalyst

go 1.19

replace github.com/xeipuuv/gojsonschema => github.com/warjiang/gojsonschema v1.2.1-0.20210329105853-aa9f9a8cfec7

require (
	github.com/alecthomas/kong v0.7.0
	github.com/alecthomas/kong-yaml v0.1.1
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220816024939-bc8df83d7b9d
	github.com/arangodb/go-driver v1.4.0
	github.com/aws/aws-sdk-go v1.44.123
	github.com/blevesearch/bleve/v2 v2.3.5
	github.com/coreos/go-oidc/v3 v3.4.0
	github.com/docker/docker v17.12.0-ce-rc1.0.20201201034508-7d75c1d40d88+incompatible
	github.com/go-chi/chi/v5 v5.0.7
	github.com/gobwas/ws v1.1.0
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/icza/dyno v0.0.0-20220812133438-f0b6f8a18845
	github.com/imdario/mergo v0.3.13
	github.com/jonas-plum/maut v0.0.0-20221001191853-1856efe6da0b
	github.com/mingrammer/commonregex v1.0.1
	github.com/stretchr/testify v1.8.1
	github.com/tidwall/gjson v1.14.3
	github.com/tidwall/sjson v1.2.5
	github.com/tus/tusd v1.10.0
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/exp v0.0.0-20221025133541-111beb427cde
	golang.org/x/oauth2 v0.1.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/RoaringBitmap/roaring v0.9.4 // indirect
	github.com/arangodb/go-velocypack v0.0.0-20200318135517-5af53c29c67e // indirect
	github.com/bits-and-blooms/bitset v1.2.1 // indirect
	github.com/blevesearch/bleve_index_api v1.0.4 // indirect
	github.com/blevesearch/geo v0.1.15 // indirect
	github.com/blevesearch/go-porterstemmer v1.0.3 // indirect
	github.com/blevesearch/gtreap v0.1.1 // indirect
	github.com/blevesearch/mmap-go v1.0.4 // indirect
	github.com/blevesearch/scorch_segment_api/v2 v2.1.3 // indirect
	github.com/blevesearch/segment v0.9.0 // indirect
	github.com/blevesearch/snowballstem v0.9.0 // indirect
	github.com/blevesearch/upsidedown_store_api v1.0.1 // indirect
	github.com/blevesearch/vellum v1.0.9 // indirect
	github.com/blevesearch/zapx/v11 v11.3.6 // indirect
	github.com/blevesearch/zapx/v12 v12.3.6 // indirect
	github.com/blevesearch/zapx/v13 v13.3.6 // indirect
	github.com/blevesearch/zapx/v14 v14.3.6 // indirect
	github.com/blevesearch/zapx/v15 v15.3.6 // indirect
	github.com/bmizerany/pat v0.0.0-20210406213842-e4b6760bdd6f // indirect
	github.com/containerd/containerd v1.6.8 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/geo v0.0.0-20210211234256-740aa86cb551 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mschoch/smat v0.2.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	golang.org/x/crypto v0.0.0-20220926161630-eccd6366d1be // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220920201722-2b89144ce006 // indirect
	google.golang.org/grpc v1.50.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
