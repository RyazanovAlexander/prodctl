docker run -it docker.io/aryazanov/product:0.0.1 /bin/bash

prodctl
prodctl version
prodctl deploy --namespace test
prodctl environment deploy --namespace test
prodctl release deploy --namespace test
prodctl release engine deploy --namespace test
prodctl release test --filter smoke
prodctl release delete