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
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("{{ .Name }}"), {{ .Name }}Array))
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
                FDateTime Dt;
                if (FDateTime::ParseIso8601(Item->AsString(), Dt))
                {
                    {{ .Name }}.Add(Dt);
                }
                {{- else if eq .Type "json" }}
                {{ .Name }}.Add(Item);
                {{- else if eq .Type "struct" }}
                TSharedPtr<FJsonObject> Obj = Item->AsObject();
                if (Obj.IsValid())
                {
                    {{ fieldElemType . }} ObjItem;
                    ObjItem.Load(Obj);
                    {{ .Name }}.Add(ObjItem);
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
