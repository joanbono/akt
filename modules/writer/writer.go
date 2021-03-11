package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/youyo/awsprofile"
)

var yellow = color.New(color.FgYellow)
var red = color.New(color.FgRed)
var green = color.New(color.FgGreen)
var cyan = color.New(color.FgCyan)
var bold = color.New(color.FgHiWhite, color.Bold)

//Printer will print the keys to copy/paste them
//wherever the user wants. Colorized, of course
func Printer(profile, accessKey, secretKey string) {
	fmt.Printf("\n[%v]\n%v = %v\n%v = %v\n\n", green.Sprintf(profile), red.Sprintf("aws_access_key_id"), bold.Sprintf(accessKey), red.Sprintf("aws_secret_access_key"), bold.Sprintf(secretKey))
}

//Profiler will get the old credentials
//and send them to the update function
func Profiler(profile, newSecretKey, newAccessKey string) {
	awsProfile := awsprofile.New()
	err := awsProfile.Parse()
	CheckErr(err)

	keyid := awsProfile.GetCredentials()
	//this is an error in the lib
	oldSecretKey, err := keyid.GetAwsAccessKeyID(profile)
	oldAccessKey, err := keyid.GetAwsSecretAccessKey(profile)
	CheckErr(err)

	awsCredentials := Reader()
	UpdateCredentials(awsCredentials, oldAccessKey, newAccessKey)
	UpdateCredentials(awsCredentials, oldSecretKey, newSecretKey)
}

//Reader will get the location
//of the aws credentials file
func Reader() string {
	homedir, err := os.UserHomeDir()
	CheckErr(err)
	awsCredentials := fmt.Sprintf("%v%s.aws%scredentials", homedir, string(os.PathSeparator), string(os.PathSeparator))
	return awsCredentials
}

//UpdateCredentials will take the old
//keys and replace them with the new ones
func UpdateCredentials(awsCredentials, oldKey, newKey string) {
	read, err := ioutil.ReadFile(awsCredentials)
	CheckErr(err)

	newKeys := strings.Replace(string(read), oldKey, newKey, -1)

	err = ioutil.WriteFile(awsCredentials, []byte(newKeys), 0)
	CheckErr(err)
}

//CheckErr will check the errors
func CheckErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
