module github.com/multitheftauto/community

go 1.12

require (
	github.com/Masterminds/squirrel v1.1.0
	github.com/appleboy/gin-jwt v0.0.0-20190409072159-633d983b91f0
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/lib/pq v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/tidwall/gjson v1.2.1 // indirect
	github.com/tidwall/match v1.0.1 // indirect
	gocloud.dev v0.13.0
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	golang.org/x/sys v0.0.0-20190527104216-9cd6430ef91e // indirect
	golang.org/x/text v0.3.2 // indirect
)

// Hack from https://github.com/gin-gonic/gin/issues/1673
// Watch https://github.com/ugorji/go/issues/299
replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
