# Data Models Service

[![GoDoc](https://godoc.org/github.com/chop-dbhi/data-models-service?status.svg)](https://godoc.org/github.com/chop-dbhi/data-models-service) [![Circle CI](https://circleci.com/gh/chop-dbhi/data-models-service.svg?style=svg)](https://circleci.com/gh/chop-dbhi/data-models-service)

Service for consuming files in the [data models format](https://github.com/chop-dbhi/data-models). The service is publicly hosted here: http://data-models.origins.link

## Build

**Dependencies**

- [Git](https://git-scm.com)
- [Go 1.3+](http://golang.org) ([test your installation](http://golang.org/doc/install#testing))

**Run**

```bash
make install build
```

This will put the `data-models` binary in your `$GOPATH/bin` directory. The examples below assumes `$GOPATH/bin` has been added to your `$PATH`.

## Usage

Once the binary is built, running the binary without any options will start the service. The service clones the data-models repository from GitHub into a `data-models` directory in your working directory. It is recommended to print the usage message to see the available options.

```bash
data-models -help
```

## Docker

Use the pre-built image on Docker Hub.

```bash
docker run -it -p 8123:8123 dbhi/data-models
```

Or build the image locally.

```bash
docker build -t dbhi/data-models .
```

### Top-level Resources

- Model Specifications - [/models](http://data-models.origins.link/models) - A specification of each data model version is available at a `/models/<data model>/<version>` endpoint (e.g., [/models/omop/5.0.0](http://data-models.origins.link/models/omop/5.0.0)). 

### Content Negotiation

The service supports representing each resource in various formats using simple content negotation. The supported formats are:

- HTML - `text/html`
- Markdown - `text/markdown`
- JSON - `application/json`

The desired format can be requested either by setting the `Accept` header to the corresponding mimetype or by adding a `format` parameter to the URL. For example, below is the OMOP v5 specification resource represented in each format:

- HTML - [/models/omop/5.0.0?format=html](http://data-models.origins.link/models/omop/5.0.0?format=html)
- Markdown - [/models/omop/5.0.0?format=md](http://data-models.origins.link/models/omop/5.0.0?format=md)
- JSON - [/models/omop/5.0.0?format=json](http://data-models.origins.link/models/omop/5.0.0?format=json)

Representations are tailored to the clients that are expected to use the resource, as described below. Note that some resources do not support all formats. The HTML format is the default format provided when neither method of content negotiation are used.

### Model Specification Resources

#### HTML

The HTML format (e.g., [OMOP v5](http://data-models.origins.link/models/omop/5.0.0?format=html)) is intended as a very simple proof of concept for displaying the data model specification in a web client for review by data model and/or data users. As such, it begins with the data model version id and a reference URL, followed by a list of tables (which serves as a linked table of contents). Each table section includes the table description and a list of fields (again, a linked table of contents). For each field, "refers to" information, if it exists, is followed by the description and any schema specifications. A table of mappings and a table of inbound references are also provided, if that information is found. This content represents an aggregation of information about the data model which we think would be useful for data model and/or data users.

#### Markdown

The Markdown format (e.g., [OMOP v5](http://data-models.origins.link/models/omop/5.0.0?format=md)) provides the same information as the HTML format. In fact, the HTML format is derived directly from the Markdown. The specific choices about header levels and organization can be seen at the actual endpoints linked above. This is intended as an API of sorts from which use-case-specific clients can retrieve, process, and display aggregated data model specification information as they wish.

#### JSON

The JSON format (e.g., [OMOP v5](http://data-models.origins.link/models/omop/5.0.0?format=json)), unlike the previously described formats, is intended for technical implementation clients and therefore presents a readily machine-processable and exhaustive representation of the data model specification. The top-level object contains the data model `name`, `version`, and reference `url` as well as an array of `tables`. Each object in the `tables` array contains the table `name` and `description`, an array of `fields`, and the `model` name and model `version`, to unambiguously identify the model to which the table belongs. Each object in the `fields` array contains the field `name`, `description`, `type`, and `required` status (as per governance), as well as the `default` (which defaults to `""`), `length`, `precision`, and `scale` (which all default to `0`). Each field object also contains the `table` name. This format should be useful in dynamically creating all sorts of data model operations, from schema creation to annotation to transformations.
