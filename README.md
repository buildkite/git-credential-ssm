# git-credential-ssm
git credential helper using AWS SSM parameter store.

AWS auth is fetched from EC2 IMDS by default. However you can specify `-default-auth` to use [AWS default credential chain](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-gosdk.html#specifying-credentials).

The credential helper currently only supports the `get` operation.

## Configuration

If the `git-credential-ssm` binary is in `$PATH` or `$GIT_EXEC_PATH`, it can be referenced as below:

```
git config credential.helper "ssm -parameter <name or arn>"
```

Otherwise

```
git config credential.helper "/path/to/git-credential-ssm -parameter <name or arn>"

```

## Local Testing

Run using default AWS credential chain:

```
go build
./git-credential-ssm -parameter /my/parameter -default-auth
```

Run using EC2 IMDS:

```
aws-vault exec <role-name> --ec2-server
./git-credential-ssm -parameter /my/parameter
```
