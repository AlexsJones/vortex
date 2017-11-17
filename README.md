# vortex

A simple template reader and variable injector

- Used for when you have a bunch of templates (e.g. kubernetes files) and want to inject a yaml file of variables
- Annoyingly finding something simple like vortex is hard because everything out there seems to be overcomplicated. This may as well be a little script :shrug:


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
    - name: {{.name}}
      image: {{.image}}

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
vortex -template example/demo.tmpl -output test.txt -varpath example/vars.yaml

```

## Scripting

With the following setup it is very trivial to script a simple way of parsing many templates
```
environments/
           production.yaml
templates/
          api/
             kubernetes-deployment.yaml
          web/
             service.yaml
deployment/
```

```
rm -rf deployment || true

for d in ./templates/**/*; do
    filename=$(dirname $d)
    foldername=`echo basename $filename`
    folderaltered=$(echo $foldername | sed 's/templates/deployment/g')
    echo "Creating $folderaltered"
    mkdir -p deployment/$folderaltered

   newpath=$(echo $d | sed 's/templates/deployment/g')

   vortex --template $d --output $newpath -varpath ./environments/$1.yaml
done

```
```
./script.sh production
```

### Other examples
```
build:
  commands: docker build --no-cache=true -t {{ .name }}:{{ .version }} .
  docker:
    containerID: {{ .name }}:{{ .version }}
    buildArgs:
      url: {{ .repoistory }}/{{ .name }}:{{ .version }}
kubernetes:
  namespace: {{ .namespace }}
  service: |-
    kind: Service
    apiVersion: v1
    metadata:
      name: {{ .name }}
      namespace: {{ .namespace }}
    spec:
      type: NodePort
      selector:
        app: {{ .name }}
      ports:
        - protocol: TCP
          port: 9090
          name: openport
 ```
