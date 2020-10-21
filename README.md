# m-aws-basic-infrastructure

Epiphany Module: AWS Basic Infrastructure

# Basic usage

## Build image

In main directory run:

```shell
make build
```

## Run module with provided example

### Prepare config file

Prepare your own variables in vars.mk file to use in the building process.
Sample file (examples/basic_flow/vars.mk.sample):

```shell
AWS_ACCESS_KEY_ID = "xxx"
AWS_ACCESS_KEY_SECRET = "xxx"
```

### Create an environment

```shell
cd examples/basic_flow
make all
```

or step-by-step:

```shell
cd examples/basic_flow
make init
make plan
make apply
```

### Delete environment

```shell
cd examples/basic_flow
make all-destroy
```

or step-by-step

```shell
cd examples/basic_flow
make destroy-plan
make destroy
```

## Release module

```shell
make release
```

or if you want to set a different version number:

```shell
make release VERSION=number_of_your_choice
```

## Integration tests execution

Prior to run integration tests on for AWS module specify variables on OS where you want to run tests:
- AWS_ACCESS_KEY_ID - this is your access key
- AWS_SECRET_ACCESS_KEY - this is your secret

and after that run shell command:

```shell
  make test
```

# Windows users

This module is designed for Linux/Unix development/usage only. If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then ```examples/basic_flow_devcontainer``` folder.


## Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AWS provider    | 3.7.0   | https://github.com/terraform-providers/terraform-provider-aws | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-aws/blob/master/LICENSE) |
| Make                      | 4.3     | https://www.gnu.org/software/make/                    | [ GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                        | 3.3.4   | https://github.com/mikefarah/yq/                      | [ MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
