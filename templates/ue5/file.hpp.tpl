{{- with .File -}}
#pragma once

#include "Json.h"
{{- range .FileRefs }}
#include "{{ $.Prefix }}{{ pascal .Name }}.hpp"
{{- end }}
{{ range append .AnonymousStructs .Struct }}
USTRUCT(BlueprintType)
struct F{{ $.Prefix }}{{ pascal .Name }}
{
    GENERATED_BODY()

    {{- range .Fields }}
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    {{ fieldType . }} {{ .Name }};
    {{- end }}

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        {{- range .Fields }}
        {{- if .IsArray }}
        const TArray<TSharedPtr<FJsonValue>>* {{ .Name }}Array = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("{{ .Name }}"), {{ .Name }}Array))
        {
            for (const auto& Item : *{{ .Name }}Array)
            {
                {{- if eq .Type "time" }}
                FString DateTimeStr;
                if (Item->TryGetString(DateTimeStr))
                {
                    FDateTime DateTime;
                    if (FDateTime::ParseIso8601(DateTimeStr, DateTime))
                    {
                        {{ .Name }}.Add(DateTime);
                    }
                }
                {{- else if eq .Type "json" }}
                {{ .Name }}.Add(Item);
                {{- else if eq .Type "struct" }}
                const TSharedPtr<FJsonObject> *JsonObject;
                if (Item->TryGetObject(JsonObject))
                {
                    {{ fieldElemType . }} FieldItem;
                    ObjItem.Load(JsonObject);
                    {{ .Name }}.Add(FieldItem);
                }
                {{- else }}
                {{ fieldElemType . }} FieldItem;
                {{- if eq .Type "int" }}
                if (Item->TryGetNumber(FieldItem))
                {{- else if eq .Type "long" }}
                if (Item->TryGetNumber(FieldItem))
                {{- else if eq .Type "float" }}
                if (Item->TryGetNumber(FieldItem))
                {{- else if eq .Type "bool" }}
                if (Item->TryGetBool(FieldItem))
                {{- else if eq .Type "string" }}
                if (Item->TryGetString(FieldItem))
                {{- end }}
                {
                    {{ .Name }}.Add(FieldItem);
                }
                {{- end }}
            }
        }
        {{- else if eq .Type "int" }}
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "long" }}
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "float" }}
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "bool" }}
        JsonObject.ToSharedRef()->TryGetBoolField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "string" }}
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "time" }}
        FString {{ .Name }}DtStr;
        if (JsonObject.ToSharedRef()->TryGetStringField(TEXT("{{ .Name }}"), {{ .Name }}DtStr))
        {
            FDateTime::ParseIso8601(*{{ .Name }}DtStr, {{ .Name }});
        }
        {{- else if eq .Type "json" }}
        JsonObject.ToSharedRef()->TryGetField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "struct" }}
        const TSharedPtr<FJsonObject> *{{ .Name }}ObjPtr;
        if (JsonObject.ToSharedRef()->TryGetObjectField(TEXT("{{ .Name }}"), {{ .Name }}ObjPtr))
        {
            {{ .Name }}.Load(*{{ .Name }}ObjPtr);
        }
        {{- end }}
        {{- end }}
    }
};
{{ end }}
{{- end -}}
