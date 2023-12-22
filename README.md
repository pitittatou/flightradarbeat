# Flightradarbeat

Welcome to Flightradarbeat.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/pitittatou/flightradarbeat`

## Getting Started with Flightradarbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7


### Build

To build the binary for Flightradarbeat run the command below. This will generate a binary
in the same directory with the name flightradarbeat.

```
make
```


### Run

To run Flightradarbeat with debugging output enabled, run:

```
./flightradarbeat -c flightradarbeat.yml -e -d "*"
```


### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Flightradarbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/pitittatou/flightradarbeat
git clone https://github.com/pitittatou/flightradarbeat ${GOPATH}/src/github.com/pitittatou/flightradarbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
