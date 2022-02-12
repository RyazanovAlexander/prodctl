This project demonstrates an approach to product lifecycle management: deploying, updating and removing the environment and services with their resources.

All the artifacts needed to manage resources are provided through the bundle delivered via the Docker image.

![image](https://github.com/RyazanovAlexander/prodctl/blob/feature/base-implementation/diagrams/product-bundle.png)



The repository includes a "walking skeleton" of the prodctl utility and [sample repositories](https://github.com/RyazanovAlexander/prodctl/tree/feature/base-implementation/fakes/.repositories) with the resources needed for the utility.

The [product bundle](https://github.com/RyazanovAlexander/prodctl/tree/feature/base-implementation/fakes/.bundle) itself is first formed from the specified artifacts in the [release file](https://github.com/RyazanovAlexander/prodctl/blob/feature/base-implementation/fakes/.repositories/cfg.releases/values.yaml). Subsequently, it is packaged into a Docker image.

The product bundle is delivered as a Docker image. It includes all the utilities necessary to work with the resources included in it.

To download the product, run the command:
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
