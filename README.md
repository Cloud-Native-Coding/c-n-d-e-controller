# c-n-d-e Controller

The *c-n-d-e Controller* serves as a client to the [c-n-d-e Dashboard](https://github.com/Cloud-Native-Coding/c-n-d-e-dashboard). 

Additionally it creates CRs and other necessary Kubernetes resources for the [c-n-d-e Operator](https://github.com/Cloud-Native-Coding/c-n-d-e-operator).

## Installation

- create *c-n-d-e Operator* and *c-n-d-e Dashboard* instances
- If you are using a Dashboard with mTLS, then add the certs to folder `config/api-client-cert` and uncomment the respective sections in `kustomization.yaml` and `controller.yaml`
- configure a new Cluster via the Dashboard and add _CNDE_CLUSTER_NAME_, _CNDE_API_KEY_, _CNDE_URL_ and _CNDE_KEYCLOAK_HOST_ to file `config/cnde-controller.properties`
- select the c-n-d-e Operator Namespace, cd to folder `config` and execute  `kubectl apply --dry-run=server -k .` and if does not show any errors, execute `kubectl apply -k .`
