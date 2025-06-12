{{/*
Expand the name of the chart.
*/}}
{{- define "go-prometheus-test.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this length.
*/}}
{{- define "go-prometheus-test.fullname" -}}
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
Create chart labels.
*/}}
{{- define "go-prometheus-test.labels" -}}
helm.sh/chart: {{ include "go-prometheus-test.name" . }}-{{ .Chart.Version }}
{{ include "go-prometheus-test.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create selector labels.
*/}}
{{- define "go-prometheus-test.selectorLabels" -}}
app.kubernetes.io/name: {{ include "go-prometheus-test.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "go-prometheus-test.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "go-prometheus-test.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

