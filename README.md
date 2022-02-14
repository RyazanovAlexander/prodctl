# prodctl

This project demonstrates an approach to product lifecycle management: deploying, updating and removing the environment and services with their resources.

All the artifacts needed to manage resources are provided through the bundle delivered via the Docker image.

![image](https://github.com/RyazanovAlexander/prodctl/blob/feature/base-implementation/diagrams/product-bundle.png)

The repository includes a "walking skeleton" of the prodctl utility and [sample repositories](https://github.com/RyazanovAlexander/prodctl/tree/feature/base-implementation/fakes/.repositories) with the resources needed for the utility.

The [product bundle](https://github.com/RyazanovAlexander/prodctl/tree/feature/base-implementation/fakes/.bundle) itself is first formed from the specified artifacts in the [release file](https://github.com/RyazanovAlexander/prodctl/blob/feature/base-implementation/fakes/.repositories/cfg.releases/values.yaml). Subsequently, it is packaged into a Docker image.

The product bundle is delivered as a Docker image. It includes all the utilities necessary to work with the resources included in it.

# Demonstration

Download the product of the required version:
```
docker run -it docker.io/aryazanov/product:0.0.1 /bin/bash
```

The following commands are executed inside the container:

`prodctl` - displays a list of available commands for the given context.

`prodctl deploy --namespace test` - creation of the environment + deployment of all services with their resources.

`prodctl environment deploy --namespace test` - creation of the environment.

`prodctl environment delete` - deleting the environment.

`prodctl release deploy --namespace test` - deploy all services with their resources to the current AKS cluster.

`prodctl release engine deploy --namespace test` - deploying the engine service with their resources to the current AKS cluster.

`prodctl release engine delete` - deleting the engine service with their resources from the current AKS cluster.

`prodctl release test --filter smoke` - running test.

`prodctl release delete` - deleting a release.


# Description of the project repositories
## Git branches and Versioning

When working with branches, it is recommended to follow the [Trunk-based](https://cloud.google.com/architecture/devops/devops-tech-trunk-based-development) development approach, as opposed to the [Gitflow Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) approach. Stick to [SemVer](https://semver.org/) when versioning artifacts, otherwise you lose some of Helm's features.

## Microservice repository

A specific example can be seen [here](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/repo.engine).
The microservice repository has the following structure:

```
repository (git clone https://github.com/engine)
├── .deploy/
│    ├── aws/
│    |   ├── helm/
|    |   └── terraform/
│    ├── azure/
│    |   ├── helm/
|    |   └── terraform/
│    ├── onPrem/
│    |   └── ansible/
|    └── ...
├── .pipelines/
│    ├── ci.yaml
│    ├── delete.yml
│    ├── deploy.yml
│    └── pr.yml
├── src/
|    └── ...
├── Dockerfile
├── context.yaml
└── Magefile.go (commands applied to engine)
      mage build ...
      mage bundle type={aws|azure|onPrem|...}
      mage deploy ...
      mage delete ...
      mage test ...
      mage publish ...
```

The [src](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/repo.engine/src) directory contains the microservice source code.

The [.deploy](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/repo.engine/.deploy) directory contains manifests for deploying this microservice with its infrastructure services. Depending on the type of environment, the manifest can be described using Helm, Terraform, Ansible, etc.

The [.pipelines](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/repo.engine/.pipelines) directory contains pipelines, with the help of which we release artifacts for this microservice.

Pipelines, the prodctl tool, and the developer call the methods defined in the [Magefile](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/repo.engine/Magefile.go), which is located at the root of the repository.

### Standardized deployment modules
In order not to duplicate helm manifests for each microservice, all repeating blocks are placed in separate Helm charts and included in the dependencies section in [Chart.yaml](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/repo.engine/.deploy/azure-dev/helm/Chart.yaml):
```yaml
apiVersion: v2
name: engine
description: Engine helm chart
type: application
version: 0.1.2
appVersion: 0.1.2

dependencies:
- name: service-template
  version: 1.0.1
  repository: "https://example.com/charts"
- name: azure-app-insights
  version: 0.1.0
  repository: "https://example.com/charts"
- name: azure-identity
  version: 0.1.0
  repository: "https://example.com/charts"
- name: azure-key-vault
  version: 0.1.0
  repository: "https://example.com/charts"
```

For a Helm chart [intended for AWS](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/repo.engine/.deploy/aws/helm/Chart.yaml), other manifests can be defined:
```yaml
...
dependencies:
- name: service-template
  version: 1.0.1
  repository: "https://example.com/charts"
- name: aws-some-resource
  version: 0.1.0
  repository: "https://example.com/charts"
```

We should follow a similar strategy for Terraform and other tools.

All dependencies for [Helm chart](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/lib.charts) and [Terraform modules](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/lib.terraform) are located in their respective repositories.

These repositories are separately versioned and have their own pipelines.

### Pipeline
Each repository with a microservice contains standardized pipelines for releasing artifacts. All these pipelines only refer to a common repository with pipelines, in which the real work is done.

For example, the [ci.yml](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/repo.engine/.pipelines/ci.yml) pipeline refers to [ci.yaml](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/lib.pipelines/ci.yml) from the [lib.pipelines](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/lib.pipelines) repository, which in turn uses standardized blocks:

```yaml
- template: blocks/build.yml
  parameters:
  param: ${{ parameters.param }}

- template: blocks/test.yml
  parameters:
  param: ${{ parameters.param }}

- template: blocks/scan.yml
  parameters:
  param: ${{ parameters.param }}

- template: blocks/bundle.yml
  parameters:
  param: ${{ parameters.param }}

- template: blocks/publish.yml
  parameters:
  param: ${{ parameters.param }}
```

Each of the [blocks](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/lib.pipelines/blocks/bundle.yml) calls the corresponding method from the Magefile.go file, which must be defined in it:
```yaml
parameters:
- name: param
  type: string

steps:
- script: mage bundle ${{ parameters.param }}
```

### Magefile
The Magefile contains methods for working with the current repository. For example, when calling the 'Deploy(envType string)' method, the resources specified in the './deploy/[envType]' directory will be deployed:
```Go
// Deploy deploys resources to the specified environment
// Params:
//   envType: environment type
func Deploy(envType string) error {
	...
	return nil
}
```

All common code for Magefiles is moved to a separate [repository](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/lib.mage/pkg).
A link to a specific version of this repository must be written in the go.mod file:
```Go
module github.com/RyazanovAlexander/engine/v1

go 1.17

require (
	github.com/RyazanovAlexander/lib.mage v1.12.1
)
```

### Context File
The [context file](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/repo.engine/context.yaml) contains variables for working with a specific environment. It is used for local development. So, for example, when calling the Deploy method in Magefile.go, the current value of the 'current-context' variable and the corresponding values for the given context are read from the context.yaml file, which will be used when executing the Deploy command.

```yaml
current-context: minikube

contexts:
  - name: minikube
    version: 1.46.6
    type: remote
    override:
      key: value

  - name: azure-dev
    version: 0.11.5
    type: remote
    override:
      cluster:
        name: dev
        resourceGroup: test
```

## Environments repository
The [environment repository](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/cfg.environments) contains the resources required to create the corresponding environment, as well as the environment variables. Environments ideologically know nothing about specific microservices or releases. Environment variables are passed to helm charts when they are set. So, for example, if the environment variable value "durable: false" is specified in the [environment variable file](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/cfg.environments/azure-dev/values.yaml), then the helm chart of a particular service may interpret this as the need to run in one replica, otherwise it will run in three replicas.

```
repository (git clone https://github.com/environments)
├── .pipelines/
│    ├── ci.yaml
│    ├── delete.yml
│    ├── deploy.yml
│    └── pr.yml
├── aws/
│    ├── terraform/
│    ├── context.yaml
│    └── values.yaml
├── onprem/
│    ├── ansible/
│    ├── context.yaml
│    └── values.yaml
├── ...
├── Dockerfile
└── Magefile.go
      mage bundle type={aws|azure|onPrem|...}
      mage deploy ...
      mage delete ...
      mage publish ...
```

## Releases repository
In the [releases repository](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/cfg.environments), the only resource of interest is the [values.yaml](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/cfg.releases/values.yaml) file. This file defines the composition of services, their versions, overridden variables for helm charts and the deployment order for a given version of the release. A new release is created by creating the appropriate "release/X.X.X" branch. The release ideologically has complete information about its microservices and can overwrite their variables, but knows nothing about the environment in which they will be delivered.
```
repository (git clone https://github.com/releases)
├── .pipelines/
│    ├── ci.yaml
│    ├── delete.yml
│    ├── deploy.yml
│    └── pr.yml
├── context.yaml
├── values.yaml
├── Dockerfile
└── Magefile.go
      mage build ...
      mage bundle ...
      mage deploy ...
      mage delete ...
      mage test ...
      mage publish ...
```

An example of the content of values.yaml file:
```yaml
name: White Rabbit

servies:
  - name: infrastructure
    version: 0.1.6
    order: 0
  - name: monitoring
    version: 1.1.0
    order: 1
  - name: engine
    version: 1.4.6
    order: 2
    parameters:
      someValue: override
      replicas: 5
  - name: test
    version: 0.5.5
    order: 2
  - name: ui
    version: 2.6.9
    order: 2
  - name: worker
    version: 5.1.0
    order: 2
```

## Product repository
The final artifact is assembled in the [product repository](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/cfg.product), consisting of an environment bundle, a release bundle, and the utilities necessary to manage them. All these artifacts are collected in one docker image. An example of a product bundle can be viewed [here](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.bundle). Optionally, the bundle can [include](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.bundle/bin) basic helm charts and all used docker images.
```
repository (git clone https://github.com/product)
├── .pipelines/
│    ├── ci.yaml
│    ├── delete.yml
│    ├── deploy.yml
│    └── pr.yml
├── Dockerfile
└── Magefile.go
      mage deploy ...
      mage delete ...
      mage ready ...
      mage binaries ...
```
