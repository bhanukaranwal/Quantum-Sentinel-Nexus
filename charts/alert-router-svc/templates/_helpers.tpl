{{/*
charts/alert-router-svc/templates/_helpers.tpl
This file contains common helper templates for the alert-router-svc chart.
*/}}

{{/*
Expand the name of the chart.
*/}}
{{- define "alert-router-svc.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this length.
*/}}
{{- define "alert-router-svc.fullname" -}}
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
{{- define "alert-router-svc.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels to be applied to all resources.
*/}}
{{- define "alert-router-svc.labels" -}}
helm.sh/chart: {{ include "alert-router-svc.chart" . }}
{{ include "alert-router-svc.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels used by Deployments and Services to identify matching pods.
*/}}
{{- define "alert-router-svc.selectorLabels" -}}
app.kubernetes.io/name: {{ include "alert-router-svc.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use.
*/}}
{{- define "alert-router-svc.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "alert-router-svc.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
