package main

/**
unmodified copy from https://github.com/terraform-providers/terraform-provider-aws
This would not be included in any future pull request.
*/

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

func isAWSErr(err error, code string, message string) bool {
	if err, ok := err.(awserr.Error); ok {
		return err.Code() == code && strings.Contains(err.Message(), message)
	}
	return false
}
