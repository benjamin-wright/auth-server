{{- define "ph.users.full_name" -}}
{{ printf "%s-%s" .Values.prefix .Values.userService.name }}
{{- end -}}