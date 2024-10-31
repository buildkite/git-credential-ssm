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

var version string = "HEAD"

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
	paramFlag := flag.String("parameter", "", "Name or arn of the ssm parameter where credentials are stored")
	localFlag := flag.Bool("default-credential-chain", false, "Use AWS default credential chain, rather than EC2 metadata endpoint")
	versionFlag := flag.Bool("version", false, "Print the version")
	flag.Parse()
	args := flag.Args()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if *paramFlag == "" {
		flag.PrintDefaults()
		ExitError(fmt.Errorf("missing flag \"parameter\""))
	}

	// Only handles "get" operation (get/store/erase)
	if len(args) > 0 && args[0] != "get" {
		ExitError(fmt.Errorf("unexpected argument. Only \"get\" is supported"))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		conditionalIMDSRegion(!*localFlag))
	if err != nil {
		ExitError(err)
	}

	var svc *ssm.Client
	if *localFlag {
		svc = ssm.NewFromConfig(cfg)
	} else {
		provider := ec2rolecreds.New()
		svc = ssm.New(ssm.Options{
			Credentials: aws.NewCredentialsCache(provider),
			Region:      cfg.Region,
		})
	}

	out, err := svc.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           paramFlag,
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

func conditionalIMDSRegion(ec2 bool) config.LoadOptionsFunc {
	if ec2 {
		return config.WithEC2IMDSRegion()
	}
	return func(l *config.LoadOptions) error { return nil }
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
