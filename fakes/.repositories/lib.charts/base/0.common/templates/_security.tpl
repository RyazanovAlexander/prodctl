{{/*
Create the name of the service account to use
*/}}
{{- define "security.serviceAccount.name" -}}
{{- default (include "common.fullname" .) .Values.serviceAccount.name }}
{{- end }}
