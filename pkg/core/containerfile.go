package core

import (
	"bytes"
	"html/template"

	"github.com/lcavajani/gojo/pkg/util"
)

var containerfile = `ARG FROM_IMAGE
{{- if .Spec.FromImageBuilder }}
ARG FROM_IMAGE_BUILDER
FROM ${FROM_IMAGE_BUILDER} AS builder

ARG VERSION

RUN apk add --no-cache git make curl gcc libc-dev ncurses

{{- if $source := (index .Spec.Sources 0).Provider.GitHub }}

RUN curl -OL "https://github.com/{{ $source.Owner }}/{{ $source.Repository }}/archive/v${VERSION}.tar.gz" && \\
    tar zxf "v${VERSION}.tar.gz" && cd "{{ $source.Repository }}-${VERSION}" && \\
    make && mv ./{{ $source.Repository }} /go/bin/{{ $source.Repository }}
{{- end }}
{{- end }}

FROM ${FROM_IMAGE}

{{- if not .Spec.FromImageBuilder }}
ARG VERSION
{{- end }}

LABEL maintainer="_me@spiarh.fr"

{{ if .Spec.FromImageBuilder }}
{{- $source := (index .Spec.Sources 0).Provider.GitHub }}
COPY --from=builder /go/bin/{{ $source.Repository }} /usr/local/bin/{{ $source.Repository }}
{{- end }}

RUN apk add --no-cache "{{ .Metadata.Name }}~=${VERSION}"

COPY entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

`

type Containerfile []byte

func (c Containerfile) String() string {
	return string(c)
}

func (c Containerfile) WriteToFile(path string) error {
	return util.WriteToFile(path, c, 0644)
}

func TemplateContainerfile(build *Build) (Containerfile, error) {
	tmpl, err := template.New(ContainerfileName).
		Option("missingkey=error").
		Parse(containerfile)
	if err != nil {
		return nil, err
	}

	var data bytes.Buffer
	if err := tmpl.Execute(&data, build); err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}
