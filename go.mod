module github.com/isnlan/coral

go 1.14

require (
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/armon/go-metrics v0.3.2 // indirect
	github.com/assembla/cony v0.3.2
	github.com/benweissmann/memongo v0.1.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/cloudflare/cfssl v1.5.0
	//github.com/cloudflare/cfssl v0.0.0
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/dgraph-io/ristretto v0.0.3
	github.com/docker/docker v1.4.2-0.20191101170500-ac7306503d23
	github.com/dsnet/compress v0.0.1
	github.com/elazarl/goproxy v0.0.0-20200809112317-0581fc3aee2d // indirect
	github.com/fsouza/go-dockerclient v1.6.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/go-git/go-git/v5 v5.3.0
	github.com/go-playground/form v3.1.4+incompatible
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-redis/redis/v8 v8.1.1
	github.com/go-redsync/redsync/v4 v4.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/go-immutable-radix v1.1.0 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hyperledger/fabric v2.1.1+incompatible
	github.com/hyperledger/fabric-amcl v0.0.0-20200424173818-327c9e2cf77a // indirect
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20200511190512-bcfeb58dd83a
	github.com/hyperledger/fabric-protos-go v0.0.0-20191121202242-f5500d5e3e85
	//github.com/hyperledger/fabric v0.0.0
	//github.com/hyperledger/fabric-amcl v0.0.0 // indirect
	//github.com/hyperledger/fabric-chaincode-go v0.0.0
	//github.com/hyperledger/fabric-protos-go v0.0.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/jpillora/backoff v1.0.0
	github.com/klauspost/pgzip v1.2.5
	github.com/mholt/archiver/v3 v3.5.0
	github.com/micro/go-micro v1.18.0
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/opentracing-contrib/go-stdlib v1.0.0
	github.com/opentracing/opentracing-go v1.1.1-0.20190913142402-a7454ce5950e
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/sykesm/zap-logfmt v0.0.3 // indirect
	github.com/syndtr/goleveldb v1.0.0
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc
	github.com/uber/jaeger-client-go v2.24.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/ulikunitz/xz v0.5.10
	go.mongodb.org/mongo-driver v1.3.4
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.21.1
	moul.io/http2curl v1.0.0 // indirect
)

//replace (
//	//github.com/cloudflare/cfssl v0.0.0 => github.com/cloudflare/cfssl v1.4.1
//	github.com/hyperledger/fabric v0.0.0 => github.com/hyperledger/fabric v2.1.1+incompatible
//	github.com/hyperledger/fabric-amcl v0.0.0 => github.com/hyperledger/fabric-amcl v0.0.0-20200424173818-327c9e2cf77a
//	github.com/hyperledger/fabric-chaincode-go v0.0.0 => github.com/hyperledger/fabric-chaincode-go v0.0.0-20200511190512-bcfeb58dd83a
//	github.com/hyperledger/fabric-protos-go v0.0.0 => github.com/hyperledger/fabric-protos-go v0.0.0-20191121202242-f5500d5e3e85
//)
