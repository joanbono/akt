package rotate

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fatih/color"

	"github.com/aws/aws-sdk-go/service/iam"
)

// Defining colors
var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var green = color.New(color.FgGreen)
var cyan = color.New(color.FgCyan)
var bold = color.New(color.FgHiWhite, color.Bold)

//GetUsername is the same as
//aws iam get-user output
func GetUsername(profile string) string {
	var input *iam.GetUserInput
	session, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
	})

	if err != nil {
		fmt.Printf("%v Error: %v\n", red.Sprintf("[-]"), err)
	}

	service := iam.New(session)
	result, errIam := service.GetUser(input)
	CheckIAMErr(errIam, profile, *result.User.UserName)

	return *result.User.UserName
}

//Rotate is the main function
//which will rotate the keys
func Rotate(profile, username string) (string, string, string) {
	var accessKey, secretKey, user string
	sess, _ := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		//Config: aws.Config{
		//	Region: aws.String("us-west-2"),
		//},
	})

	service := iam.New(sess)

	activeKeys, inactiveKeys := ListUserKeys(service, profile, username)
	DeleteUserKeys(service, inactiveKeys, profile, username)

	if len(activeKeys) == 1 {
		accessKey, secretKey, user = GetNewPair(service, profile, username)
		DeleteUserKeys(service, activeKeys, profile, username)
	} else if len(activeKeys) > 1 {
		DeleteUserKeys(service, activeKeys[1:], profile, username)
		accessKey, secretKey, user = GetNewPair(service, profile, username)
		DeleteUserKeys(service, activeKeys[:1], profile, username)
	}
	return accessKey, secretKey, user
}

//DeleteUserKeys will delete
//the old keys for a user
func DeleteUserKeys(service *iam.IAM, keyArray []string, profile, username string) {
	var input *iam.DeleteAccessKeyInput

	if username != "" {
		input = &iam.DeleteAccessKeyInput{
			UserName: aws.String(username),
		}
	}

	for k := range keyArray {
		input = &iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(keyArray[k]),
		}
		_, err := service.DeleteAccessKey(input)
		CheckIAMErr(err, profile, username)
	}

}

//ListUserKeys will list
//the current keys for a user
func ListUserKeys(service *iam.IAM, profile, username string) ([]string, []string) {
	var input *iam.ListAccessKeysInput
	var keyActive, keyInactive []string

	if username != "" {
		input = &iam.ListAccessKeysInput{
			UserName: aws.String(username),
		}
	}
	listKeys, err := service.ListAccessKeys(input)
	CheckIAMErr(err, profile, username)

	for k := range listKeys.AccessKeyMetadata {
		if *listKeys.AccessKeyMetadata[k].Status == "Active" {
			keyActive = append(keyActive, *listKeys.AccessKeyMetadata[k].AccessKeyId)
		} else if *listKeys.AccessKeyMetadata[k].Status == "Inactive" {
			keyInactive = append(keyInactive, *listKeys.AccessKeyMetadata[k].AccessKeyId)
		}
	}
	return keyActive, keyInactive
}

//GetNewPair will generate
//new keys for a user
func GetNewPair(service *iam.IAM, profile, username string) (string, string, string) {
	var input *iam.CreateAccessKeyInput

	if username != "" {
		input = &iam.CreateAccessKeyInput{
			UserName: aws.String(username),
		}
	}

	accessKey, err := service.CreateAccessKey(input)
	CheckIAMErr(err, profile, username)
	return *accessKey.AccessKey.AccessKeyId, *accessKey.AccessKey.SecretAccessKey, *accessKey.AccessKey.UserName
}

// CheckIAMErr will handle errors
// on IAM requests for the entire program
func CheckIAMErr(err error, profile, username string) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				//fmt.Println(iam.ErrCodeNoSuchEntityException, aerr.Error())
				fmt.Printf("%v Username %s not found.\n", red.Sprintf("[-]"), bold.Sprintf(username))
			case iam.ErrCodeLimitExceededException:
				//fmt.Println("LIMIT EXCEED", iam.ErrCodeLimitExceededException, aerr.Error())
				fmt.Printf("%v API limit exceeded.\n", red.Sprintf("[-]"))
			case iam.ErrCodeServiceFailureException:
				//fmt.Println("FAILURE", iam.ErrCodeServiceFailureException, aerr.Error())
				fmt.Printf("%v Unknown error %v\n", red.Sprintf("[-]"), bold.Sprintf(`¯\_(ツ)_/¯`))
			default:
				if aerr.Error()[:21] == "NoCredentialProviders" {
					fmt.Printf("[-] Profile %s not found\n", profile)
				} else if aerr.Error()[:12] == "AccessDenied" {
					if username == "" {
						username = GetUsername(profile)
					}
					fmt.Printf("%v %v is not authorized to perform %v\n", red.Sprintf("[-]"), bold.Sprintf(username), bold.Sprintf("iam:CreateAccessKey"))
				} else {
					fmt.Printf("%v Generic error: %v\n", red.Sprintf("[-]"), aerr.Error())
				}
			}
		} else {
			fmt.Printf("%v Generic error: %v\n", red.Sprintf("[-]"), err.Error())
		}
		os.Exit(1)
	}

}
