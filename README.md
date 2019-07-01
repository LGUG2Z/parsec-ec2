# Parsec-EC2
Easily start and stop Parsec-ready EC2 spot instances to make cloud gaming even cheaper.

## Background
Building on the excellent work done by [Larry Gadea](https://lg.io/), [Daniel Thomas](https://github.com/DanielThomas/ec2gaming),
[Josh McGhee](https://github.com/joshpmcghee/parsec-terraform), [Benjamin Malley](https://github.com/BenjaminMalley/ec2gaming),
and the [Parsec team](https://parsec.tv/), I started working on this project to allow for two very specific pieces of 
functionality I was looking for that I had not yet seen implemented anywhere else:

* To be able to easily switch between instance types without manually editing files
* To be able to arbitrarily specify a spot bid price per session relative to the current highest spot price in a given availability zone

This is very much a work in progress and my first attempt at writing a non-trivial application in Go. Improvements and pull
requests are very much encouraged and welcome.

## Requirements
* [Terraform](https://github.com/hashicorp/terraform)
* [aws-cli](https://github.com/aws/aws-cli)
* [Go](https://github.com/golang/go)
* [Parsec Client](https://parsec.tv/downloads)

This has been developed with MacOS in mind, but should also work on Linux.

Note that for now this project will not work with Terraform 0.11.0, as this version introduced some breaking changes in how the `terraform output` command works.
For now, to continue using `parsec-ec2`, use Terraform 0.10.0. Multiple versions of Terraform can be managed using the [tfenv](https://github.com/kamatama41/tfenv) project.

## Installation
The latest version of `parsec-ec2` can be installed using `go get`.

```
go get -u github.com/LGUG2Z/parsec-ec2
```

Make sure `$GOPATH` is set correctly that and that `$GOPATH/bin` is in your `$PATH`.

The `parsec-ec2` executable will be installed under the `$GOPATH/bin` directory.

Once installed, add `export PARSEC_EC2_SERVER_KEY=your_server_key` to your shell rc file.

You can find your server key by going to your [Parsec account page](https://parsec.tv/account) and looking at the
'Your Configuration Settings for Self Hosting in the Cloud' section at the bottom. You may first have to click the
'Generate Config Settings' button. You should eventually see something like the following:

```
network_server_start_port=8000
app_host=1
server_key=xxxxx
app_check_user_data=1
app_first_run=0
```

The `server_key` value is the one that you should assign to the `PARSEC_EC2_SERVER_KEY` environment variable.

## Usage
### init
After an initial installation or upgrade, all users should run `parsec-ec2 init`.

The init command will create the directory `$HOME/.parsec-ec2` and the required Terraform template and provisioning
userdata files. The command can safely be run multiple times.

### price
The `price` command looks for the current highest spot price for the requested instance type in the requested region.

The `--region` and `--instance-type` flags are required.

Examples:
```
$ parsec-ec2 price --region eu-west-1 --instance-type g2.2xlarge

>> The highest spot price in region eu-west-1 for g2.2xlarge instances is currently $0.87/hour.
```

### start
The `start` command makes a spot request for the requested EC2 instance type in the specified region. If
`PARSEC_EC2_SERVER_KEY` has not been exported in the shell rc file, it must be passed to the command using the 
`--server-key` flag.

The amount to bid above the current highest spot price for the instance is specified using the `--bid` flag, so if the
current highest spot price is $0.20, running the command with `--bid 0.10` will make a spot request with a bid price
of $0.30. Alternatively this flag can be left blank if you don't want to bid higher than the current highest bid price.

You also need to provide --volume-size parameter which is going to be used to determine the size of the root volume.
And --ami-name will be used to search AMI either from Parsec's AWS account or your own.

If the `--plan` flag is used, the spot request will not be sent and instead the `terraform plan` command will be run
which will output to the terminal the details of any AWS resources that will be created by running the `start` command.

Examples:
```
# With PARSEC_EC2_SERVER_KEY already set as an env variable
parsec-ec2 start \
--region eu-west-1 \
--instance-type g3.4xlarge \
--bid 0.10 \
--volume-size 45 \
--ami-name parsec-ksp-1

```
```
# With the server key passed using the --server-key flag
parsec-ec2 start \
--region eu-west-2 \
--instance-type g2.2xlarge \
--bid 0.10 \
--volume-size 45 \
--ami-name parsec-ksp-1 \
---server-key xxxxx
```
```
# With the --plan flag
parsec-ec2 start \
--region eu-central-1 \
--instance-type g2.2xlarge \
--bid 0.10 \
--volume-size 45 \
--ami-name parsec-ksp-1 \
--plan
```

### status
The `status` command queries the launched instance and gets the current initialisation status.

Once an instance is reporting a status of initialised, it may still take some time for the instance to show up in the
Parsec desktop application. This is because time is still required for the provisioning script to run on the instance, 
which is what will allow the Parsec application to launch and log in with the provided Parsec server key.

Example:
```
parsec-ec2 status
```

### stop
The `stop` command stops a Parsec EC2 instance created using the `start` command. Under the hood this command runs 
`terraform destroy`, with removes all AWS resources that are identified for creation in the terraform template.

This command depends on session information that is created by the `start` command and stored in `$HOME/.parsec-ec2/currentSession.json`,
so if this has been manually modified or removed after running the `start` command, the `stop` command will not execute. In
this situation it is still possible to manually run `terraform destroy` in the `$HOME/.parsec-ec2` directory. You will 
receive prompts for variable values, but these can all be left blank with the exception of the region variable, which
can be set to the region the instances were started in.

Example:
```
parsec-ec2 stop
```
