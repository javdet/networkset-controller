{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "networkset-controller.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "networkset-controller.fullname" -}}
{{-   if .Values.fullnameOverride -}}
{{-     .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{-   else -}}
{{-     $name := default .Chart.Name .Values.nameOverride -}}
{{-     if hasPrefix $name .Release.Name -}}
{{-       .Release.Name | trunc 63 | trimSuffix "-" -}}
{{-     else -}}
{{-       printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{-     end -}}
{{-   end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "networkset-controller.chart" -}}
{{-   printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Define the name of the secret containing the tokens
*/}}
{{- define "networkset-controller.secret" -}}
{{- default (include "networkset-controller.fullname" .) .Values.runners.secret | quote -}}
{{- end -}}

{{/*
Define the image, using .Chart.AppVersion and networkset controller image as a default value
*/}}
{{- define "networkset-controller.image" }}
{{- if kindIs "string" .Values.image -}}
{{- .Values.image }}
{{- else -}}
{{- $appVersion := ternary "bleeding" (print "v" .Chart.AppVersion) (eq .Chart.AppVersion "bleeding") -}}
{{- $appVersionImageTag := printf "alpine-%s" $appVersion -}}
{{- $imageTag := default $appVersionImageTag .Values.image.tag -}}
{{- printf "%s/%s:%s" .Values.image.registry .Values.image.image $imageTag }}
{{- end -}}
{{- end -}}


{{/*
Define the server session internal port, using 9000 as a default value
*/}}
{{- define "networkset-controller.server-session-external-port" }}
{{-   default 9000 .Values.sessionServer.externalPort }}
{{- end -}}

{{/*
Define the server session external port, using 8093 as a default value
*/}}
{{- define "networkset-controller.server-session-internal-port" }}
{{-   default 8093 .Values.sessionServer.internalPort }}
{{- end -}}

