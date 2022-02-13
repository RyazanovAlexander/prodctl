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
      mage remove ...
      mage test ...
      mage publish ...
```

The [src](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/repo.engine/src) directory contains the microservice source code.

The [.deploy](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/repo.engine/.deploy) directory contains manifests for deploying this microservice with its infrastructure services. Depending on the type of environment, the manifest can be described using Helm, Terraform, Ansible, etc.

The [.pipelines](https://github.com/RyazanovAlexander/prodctl/tree/main/fakes/.repositories/repo.engine/.pipelines) directory contains pipelines, with the help of which we release artifacts for this microservice.

Pipelines, the prodctl tool, and the developer call the methods defined in the [Magefile](https://github.com/RyazanovAlexander/prodctl/blob/main/fakes/.repositories/repo.engine/Magefile.go), which is located at the root of the repository.
