package client

// import (
// 	"fmt"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/awserr"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/secretsmanager"
// 	log "github.com/sirupsen/logrus"
// )

// type SecretClient struct {
// 	region  string
// 	session *session.Session
// }

// // get secret from AWS SecretManager
// func (c *SecretClient) GetSecret(key string) (string, error) {
// 	log.Debug("Request secret from AWS SecretManager")
// 	cfg := aws.NewConfig().WithRegion(c.region).WithCredentials(&credentials.Credentials{})
// 	sess, _ := session.NewSession(cfg)
// 	svc := secretsmanager.New(sess, cfg)
// 	output, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
// 		SecretId:     aws.String(key),
// 		VersionStage: aws.String("AWSCURRENT"),
// 	})
// 	if err != nil {
// 		if err, ok := err.(awserr.Error); ok {
// 			switch err.Code() {
// 			case secretsmanager.ErrCodeDecryptionFailure:
// 				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
// 				log.Error("Secrets Manager can't decrypt the protected secret text using the provided KMS key.")
// 			case secretsmanager.ErrCodeInternalServiceError:
// 				// An error occurred on the server side.
// 				log.Error("An error occurred on the server side.", err.Error())
// 			case secretsmanager.ErrCodeInvalidParameterException:
// 				// You provided an invalid value for a parameter.
// 				log.Error("You provided an invalid value for a parameter.", err.Error())
// 			case secretsmanager.ErrCodeInvalidRequestException:
// 				// You provided a parameter value that is not valid for the current state of the resource.
// 				log.Error("You provided a parameter value that is not valid for the current state of the resource.", err.Error())
// 			case secretsmanager.ErrCodeResourceNotFoundException:
// 				// We can't find the resource that you asked for.
// 				log.Error("We can't find the resource that you asked for.", err.Error())
// 			default:
// 				log.Error(err.Code(), err.Error())
// 			}
// 		} else {
// 			// Print the error, cast err to awserr.Error to get the Code and
// 			// Message from an error.
// 			log.Error("Unknown AWS Error", err.Error())
// 		}
// 		return "", err
// 	}
// 	if output.SecretString == nil {
// 		return "", fmt.Errorf("no secret string")
// 	}
// 	return *output.SecretString, nil
// }
