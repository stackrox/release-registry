{{/*
Expand the name of the chart.
*/}}
{{- define "release-registry.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "release-registry.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "release-registry.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "release-registry.labels" -}}
helm.sh/chart: {{ include "release-registry.chart" . }}
{{ include "release-registry.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "release-registry.selectorLabels" -}}
app.kubernetes.io/name: {{ include "release-registry.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "release-registry.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "release-registry.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create secret to access docker registry
*/}}
{{- define "release-registry.imagePullSecret" }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.imageRegistry.name (printf "%s:%s" .Values.imageRegistry.username .Values.imageRegistry.password | b64enc) | b64enc }}
{{- end }}

{{/*
Create configuration secret
*/}}
{{- define "release-registry.configuration" }}
database:
  type: postgres
server:
  port: {{ .Values.service.port }}
  cert: /certs/tls.crt
  key: /certs/tls.key
  staticDir: ui
  docsDir: docs
  metrics:
    port: {{ .Values.server.metrics.port }}
    measureLatency: {{ .Values.server.metrics.measureLatency }}
tenant:
  emailDomain: {{ .Values.server.emailDomain }}
  password: {{ .Values.server.adminPassword }}
  oidcConfigFile: /config/oidc.yaml
{{- end }}
