module github.com/hogehoge-banana/sls-rtc-backend

require (
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.33.11
	github.com/google/uuid v1.1.1
	github.com/stretchr/testify v1.5.1
)

go 1.13

replace github.com/hogehoge-banana/sls-rtc-backend/ => ./
