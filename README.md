# Operator Environment (WIP)

**Work In Progress! Do not use on you CI or Production environments!**

Simple commandline tool to spin Kubernetes or OpenShift using Docker in Docker concept.
Can simulate different type of issues like node failure etc.

## Requirements

Your host system should have Docker >= 1.12
Go >= 1.11

## Installation / Building from Source

```bash
go build -o op-env
```

## Usage

```text
./op-env up kubernetes -v
```

## Testing Features
- Failure of random Node
- Some cool stuff

## TODO
- Nodes creation
- Failure of random Node
 


