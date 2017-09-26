# vortex

A simple template reader and variable injector

- Used for when you have a bunch of templates (e.g. kubernetes files) and want to inject a yaml file of variables

## Usage

```
vortex -template example/demo.tmpl -output test.txt -varpath example/vars.yaml

```
