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
