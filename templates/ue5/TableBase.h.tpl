#pragma once

#include "Json.h"
#include "{{ .Prefix }}TableBase.generated.h"

USTRUCT(BlueprintType)
struct F{{ .Prefix }}TableBase
{
    GENERATED_BODY()

    F{{ .Prefix }}TableBase() {}
    virtual ~F{{ .Prefix }}TableBase() {}

    virtual FString GetSheetName() const;
    virtual void Load(const TSharedPtr<FJsonValue>& JsonValue);
};
