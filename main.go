package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type credential struct {
	protocol string
	host     string
	username string
	password string
}

func (c credential) String() string {
	return fmt.Sprintf("protocol=%s\nhost=%s\nusername=%s\npassword=%s\n",
		c.protocol, c.host, c.username, c.password)
}

func ExitError(e error) {
	fmt.Fprintln(os.Stderr, e)
	os.Exit(1)
}

func main() {

	paramName := flag.String("parameter", "", "Name or arn of the ssm parameter where credentials are stored")
	local := flag.Bool("default-auth", false, "Use AWS default credential chain, rather than EC2 metadata")
	flag.Parse()
	args := flag.Args()

	// Only handles "get" operation (get/store/erase)
	if len(args) > 0 && args[0] != "get" {
		os.Exit(0)
	}
	if *paramName == "" {
		flag.PrintDefaults()
		ExitError(fmt.Errorf("missing flag \"parameter\""))
	}

	cfg, err := config.LoadDefaultConfig(context.Background())

	var svc *ssm.Client
	if *local {
		if err != nil {
			ExitError(err)
		}
		svc = ssm.NewFromConfig(cfg)
	} else {
		provider := ec2rolecreds.New()
		svc = ssm.New(ssm.Options{
			Credentials: aws.NewCredentialsCache(provider),
			Region:      cfg.Region,
		})
	}

	out, err := svc.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           paramName,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		ExitError(err)
	}

	val := aws.ToString(out.Parameter.Value)

	cred, err := parseUrl(val)
	if err != nil {
		ExitError(err)
	}

	// output credential string in the expected format
	// https://git-scm.com/docs/git-credential#IOFMT
	fmt.Println(cred)

}

// parseUrl translates from url format to credential type
func parseUrl(s string) (credential, error) {

	url, err := url.Parse(s)

	pass, ok := url.User.Password()
	if !ok {
		return credential{}, fmt.Errorf("url is missing password")
	}

	return credential{
		host:     url.Host,
		protocol: url.Scheme,
		username: url.User.Username(),
		password: pass,
	}, err
}
