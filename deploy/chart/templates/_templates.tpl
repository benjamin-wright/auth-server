{{- define "ph.migrations.full_name" -}}
{{ printf "%s-%s" .Values.prefix .Values.users.migrations.name }}
{{- end -}}

{{- define "ph.users.full_name" -}}
{{ printf "%s-%s" .Values.prefix .Values.users.name }}
{{- end -}}

{{- define "ph.tokens.full_name" -}}
{{ printf "%s-%s" .Values.prefix .Values.tokens.name }}
{{- end -}}

{{- define "ph.verify.full_name" -}}
{{ printf "%s-%s" .Values.prefix .Values.verify.name }}
{{- end -}}

{{- define "ph.forms.full_name" -}}
{{ printf "%s-%s" .Values.prefix .Values.forms.name }}
{{- end -}}