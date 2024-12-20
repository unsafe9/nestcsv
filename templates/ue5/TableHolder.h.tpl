// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once

{{ range .Tables }}
#include "{{ $.Prefix }}{{ pascal .Name }}Table.h"
{{- end }}
#include "{{ .Prefix }}TableHolder.generated.h"

UCLASS(BlueprintType)
class U{{ .Prefix }}TableHolder : public UObject
{
    GENERATED_BODY()

public:
    {{- range .Tables }}
    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    F{{ $.Prefix }}{{ pascal .Name }}Table {{ pascal .Name }};
    {{- end }}

    TArray<F{{ .Prefix }}TableBase*> GetTables()
    {
        return {
            {{- range .Tables }}
            &{{ pascal .Name }},
            {{- end }}
        };
    }

    F{{ .Prefix }}TableBase* GetBySheetName(const FString& SheetName)
    {
        {{- range .Tables }}
        if (SheetName == {{ pascal .Name }}.GetSheetName()) return &{{ pascal .Name }};
        {{- end }}
        return nullptr;
    }

    template <class T = F{{ .Prefix }}TableBase>
    T* Get()
    {
        {{- range .Tables }}
        if constexpr (std::is_same_v<T, F{{ $.Prefix }}{{ pascal .Name }}Table>) return &{{ pascal .Name }};
        {{- end }}
        return nullptr;
    }
};
