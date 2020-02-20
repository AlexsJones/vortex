This project is forked from github.com/AlexsJones/vortex

A simple template reader and variable injector

- Used for when you have a bunch of templates (e.g. kubernetes files) and want to inject a yaml file of variables
- Supports giving it a directory of nested templates and an output path (it will reproduce the directory structure)


## Install

`go get github.com/codelity/vortex`

_Or navigate to the releases page and install as a binary on the path_

### Run straight from docker on a local directory

` docker run -v /home/foo/myproject:/tmp tibbar/vortex:v1 -template /tmp/demo.yaml -output /tmp/deployment -varpath /tmp/vars.vortex`

Outputs on `/home/foo/myproject/deployment`

## Example

demo.tmpl
```
apiVersion: v1
kind: Pod
metadata:
  name: console
spec:
  restartPolicy: Always
  containers:
    - name: {{{.name}}}
      image: {{{.image}}}

```

vars.yaml
```
name: "test"
image: "us.gcr.io/test"

```

The result;
```
apiVersion: v1
kind: Pod
metadata:
  name: console
spec:
  restartPolicy: Always
  containers:
    - name: test
      image: us.gcr.io/test
````
## Usage

```
vortex -template example/demo.yaml -output deployment -varpath example/vars.vortex

```

## Recursive folder templating

Vortex also can recursively follow a template folder

e.g.
```
somefolders/
            foo/
                template.yaml
            bar/
                another.yaml

```

```
vortex -template somefolders -output anoutputfolder -var example/vars.yaml
```


### Other examples

The perfect Kubernetes companion...

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{{.info.namespace}}}-ingress
  namespace: {{{.info.namespace}}}
  annotations:
    kubernetes.io/ingress.class: "nginx"
    kubernetes.io/ingress.allow-http: "false"
spec:
  tls:
  {{{ range .ingress }}}
  - hosts:
    - {{{.ing.hostname}}}
    secretName: {{{.ing.tlssecretname}}}
  {{{end}}}
  rules:
  {{{ range .ingress }}}
  - host: {{{.ing.hostname}}}
    http:
      paths:
      - backend:
          serviceName: {{{.ing.servicename}}}
          servicePort: {{{.ing.serviceport}}}
  {{{end}}}

 ```
 ```
 apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubernetes-{{{.info.environment}}}-service-account
  namespace: frontier

---

kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kubernetes-{{{.info.environment}}}-service-account-binding
subjects:
  - kind: ServiceAccount
    name: kubernetes-{{{.info.environment}}}-service-account
    namespace: frontier
roleRef:
  kind: ClusterRole
  name: view
  apiGroup: ""
 ```

### Loading a variable from a connected vault instance.

```
env:
    API_KEY: {{{ vaultsecret "/secret/path/to/secret" "keyInDataMap" }}}
```

For this to work, you will need to have:
- VAULT_ADDR exported in your shell to the running vault instance
- VAULT_TOKEN exported in your shell if "${HOME}/.vault-token" isn't present

Using environment variables inside your templates:

```
annotations:
  UpdatedBy: {{{ getenv "USER" }}}
  SecretUsed: {{{ getenv "SECRET_TOKEN" }}}
```

This enables secrets to be loaded via environment variables rather than alternative methods such as using `sed` over
the template before processing.


#### Docker Build

```
go mod vendor
docker build .
```

