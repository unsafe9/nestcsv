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
    GENERATED_USTRUCT_BODY()

public:
    {{- range .Fields }}
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    {{ fieldType . }} {{ .Name }};
    {{- end }}

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        {{- range .Fields }}
        {{- if .IsArray }}
        const TArray<TSharedPtr<FJsonValue>>* {{ .Name }}Array = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("{{ .Name }}"), {{ .Name }}Array))
        {
            for (const auto& Item : *{{ .Name }}Array)
            {
                {{- if eq .Type "int" }}
                {{ .Name }}.Add(Item->AsNumber());
                {{- else if eq .Type "long" }}
                {{ .Name }}.Add(Item->AsNumber());
                {{- else if eq .Type "float" }}
                {{ .Name }}.Add(Item->AsNumber());
                {{- else if eq .Type "bool" }}
                {{ .Name }}.Add(Item->AsBool());
                {{- else if eq .Type "string" }}
                {{ .Name }}.Add(Item->AsString());
                {{- else if eq .Type "time" }}
                FString {{ .Name }}DtStr = Item->AsString();
                {{ .Name }}.Add(FDateTime::FromIso8601({{ .Name }}DtStr));
                {{- else if eq .Type "json" }}
                {{ .Name }}.Add(Item);
                {{- else if eq .Type "struct" }}
                TSharedPtr<FJsonObject> {{ .Name }}Object = Item->AsObject();
                if ({{ .Name }}Object.IsValid())
                {
                    {{ fieldType . }} {{ .Name }}Item;
                    {{ .Name }}Item.Load({{ .Name }}Object);
                    {{ .Name }}.Add({{ .Name }}Item);
                }
                {{- end }}
            }
        }
        {{- else if eq .Type "int" }}
        JsonObject.TryGetNumberField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "long" }}
        JsonObject.TryGetNumberField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "float" }}
        JsonObject.TryGetNumberField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "bool" }}
        JsonObject.TryGetBoolField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "string" }}
        JsonObject.TryGetStringField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "time" }}
        FString {{ .Name }}DtStr;
        if (JsonObject.TryGetStringField(TEXT("{{ .Name }}"), {{ .Name }}DtStr))
        {
            {{ .Name }} = FDateTime::FromIso8601({{ .Name }}DtStr);
        }
        {{- else if eq .Type "json" }}
        JsonObject.TryGetField(TEXT("{{ .Name }}"), {{ .Name }});
        {{- else if eq .Type "struct" }}
        TSharedPtr<FJsonObject> {{ .Name }}Object;
        if (JsonObject->TryGetObjectField(TEXT("{{ .Name }}"), {{ .Name }}Object))
        {
            {{ .Name }}.Load({{ .Name }}Object);
        }
        {{- end }}
        {{- end }}
    }
};
{{ end }}
{{- end -}}
