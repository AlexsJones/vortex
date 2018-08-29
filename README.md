```
   :::     :::  ::::::::  ::::::::: ::::::::::: :::::::::: :::    :::
  :+:     :+: :+:    :+: :+:    :+:    :+:     :+:        :+:    :+:  
 +:+     +:+ +:+    +:+ +:+    +:+    +:+     +:+         +:+  +:+    
+#+     +:+ +#+    +:+ +#++:++#:     +#+     +#++:++#     +#++:+      
+#+   +#+  +#+    +#+ +#+    +#+    +#+     +#+         +#+  +#+      
#+#+#+#   #+#    #+# #+#    #+#    #+#     #+#        #+#    #+#      
 ###      ########  ###    ###    ###     ########## ###    ###       
```

---

[![Build Status](https://travis-ci.org/AlexsJones/vortex.svg?branch=master)](https://travis-ci.org/AlexsJones/vortex)
[![Maintainability](https://api.codeclimate.com/v1/badges/93b3be49a1b077adc0ba/maintainability)](https://codeclimate.com/github/AlexsJones/vortex/maintainability)

A simple template reader and variable injector

- Used for when you have a bunch of templates (e.g. kubernetes files) and want to inject a yaml file of variables
- Supports giving it a directory of nested templates and an output path (it will reproduce the directory structure)
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

Loading a variable from a running vault instance.
```
env:
    API_KEY: {{ vaultsecret "/secret/path/to/secret" "keyInDataMap" }}
```

For this to work, you will need to have:
- VAULT_ADDR exported in your shell to the running vault instance
- VAULT_TOKEN exported in your shell if "${HOME}/.vault-token" isn't present
