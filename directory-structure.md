.
├── Dockerfile
├── README.md
├── cmd
│   └── main.go
├── database
│   ├── definition.sql
│   └── seeds.sql
├── directory-structure.md
├── docker-compose.yml
├── documents
│   ├── erd
│   │   ├── OnlineStoreERD.png
│   │   └── erd_syntax.docx
│   └── postman
│       ├── Online Store API.postman_collection.json
│       ├── PostmanSS1.png
│       └── PostmanSS2.png
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── cache
│   │   │   └── cache.go
│   │   ├── config
│   │   │   ├── app.go
│   │   │   ├── gorm.go
│   │   │   ├── log.go
│   │   │   ├── redis.go
│   │   │   ├── response.go
│   │   │   ├── validator.go
│   │   │   └── viper.go
│   │   ├── constants
│   │   │   └── constant.go
│   │   ├── delivery
│   │   │   ├── handler
│   │   │   │   ├── cart.go
│   │   │   │   ├── customer.go
│   │   │   │   ├── order.go
│   │   │   │   └── product.go
│   │   │   └── router
│   │   │       └── router.go
│   │   ├── model
│   │   │   ├── cart.go
│   │   │   ├── cart_item.go
│   │   │   ├── customer.go
│   │   │   ├── order_item.go
│   │   │   ├── order_status.go
│   │   │   ├── orders.go
│   │   │   ├── product.go
│   │   │   └── product_category.go
│   │   ├── payloads
│   │   │   ├── request.go
│   │   │   └── response.go
│   │   ├── repository
│   │   │   ├── cart.go
│   │   │   ├── cart_item.go
│   │   │   ├── customer.go
│   │   │   ├── order.go
│   │   │   ├── order_item.go
│   │   │   └── product.go
│   │   └── service
│   │       ├── cart.go
│   │       ├── customer.go
│   │       ├── order.go
│   │       └── product.go
│   └── pkg
│       ├── builder
│       │   └── response_builder.go
│       ├── hash
│       │   └── hash.go
│       ├── jwt
│       │   └── jwt.go
│       └── pagination
│           └── pagination.go
└── params

24 directories, 52 files
