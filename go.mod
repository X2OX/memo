module github.com/x2ox/memo

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/yanyiwu/gojieba v1.1.2
	github.com/yuin/goldmark v1.3.9
	go.uber.org/zap v1.18.1
	go.x2ox.com/blackdatura v1.7.0
	go.x2ox.com/tea v1.0.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)

replace github.com/yanyiwu/gojieba v1.1.2 => github.com/Aoang/gojieba v1.1.3
