{{- with .File -}}
#pragma once

#include "{{ $.Prefix }}TableBase.h"
#include "{{ $.Prefix }}{{ pascal .Name }}.hpp"

USTRUCT(BlueprintType)
struct F{{ $.Prefix }}{{ pascal .Name }}Table : public F{{ $.Prefix }}TableBase
{
    GENERATED_BODY()

    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    {{- if .IsMap }}
    TMap<FString, F{{ $.Prefix }}{{ pascal .Name }}> Rows;
    {{- else }}
    TArray<F{{ $.Prefix }}{{ pascal .Name }}> Rows;
    {{- end }}

    void Load(const TSharedPtr<FJsonValue>& JsonValue) override
    {
        {{- if .IsMap }}
        const TSharedPtr<FJsonObject>* RowsMap;
        if (JsonValue->TryGetObject(RowsMap))
        {
            for (const auto& Row : (*RowsMap)->Values)
            {
                F{{ $.Prefix }}{{ pascal .Name }} RowItem;
                RowItem.Load(Row.Value);
                Rows.Add(Row.Key, RowItem);
            }
        }
        {{- else }}
        const TArray<TSharedPtr<FJsonValue>>* RowsArray;
        if (JsonValue->TryGetArray(RowsArray))
        {
            for (const auto& Row : *RowsArray)
            {
                F{{ $.Prefix }}{{ pascal .Name }} RowItem;
                RowItem.Load(Row);
                Rows.Add(RowItem);
            }
        }
        {{- end }}
    }
};
{{- end -}}
