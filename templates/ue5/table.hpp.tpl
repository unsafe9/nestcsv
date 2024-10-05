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
    
    virtual FString GetSheetName() const override
    {
        return TEXT("{{ .Name }}");
    }

    virtual void Load(const TSharedPtr<FJsonValue>& JsonValue) override
    {
        {{- if .IsMap }}
        const TSharedPtr<FJsonObject>* RowsMap = nullptr;
        if (JsonValue->TryGetObject(RowsMap))
        {
            for (const auto& Row : (*RowsMap)->Values)
            {
                const TSharedPtr<FJsonObject> *RowValue = nullptr;
                if (Row.Value->TryGetObject(RowValue))
                {
                    F{{ $.Prefix }}{{ pascal .Name }} RowItem;
                    RowItem.Load(*RowValue);
                    Rows.Add(Row.Key, RowItem);
                }
            }
        }
        {{- else }}
        const TArray<TSharedPtr<FJsonValue>>* RowsArray = nullptr;
        if (JsonValue->TryGetArray(RowsArray))
        {
            for (const auto& Row : *RowsArray)
            {
                const TSharedPtr<FJsonObject> *RowValue = nullptr;
                if (Row.Value->TryGetObject(RowValue))
                {
                    F{{ $.Prefix }}{{ pascal .Name }} RowItem;
                    RowItem.Load(*RowValue);
                    Rows.Add(RowItem);
                }
            }
        }
        {{- end }}
    }
};
{{- end -}}
