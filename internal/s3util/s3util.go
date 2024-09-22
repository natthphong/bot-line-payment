package s3util

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/natthphong/bot-line-payment/config"
	"github.com/natthphong/bot-line-payment/internal/logz"
)

func OpenS3(config config.AwsS3Config) (*s3.S3, error) {
	logger := logz.NewLogger()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.DoSpaceRegion),
		Endpoint:    aws.String(config.DoSpaceEndpoint),
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
	})
	if err != nil {
		logger.Error("Error creating AWS session:")
		return nil, err
	}
	return s3.New(sess), err

}
