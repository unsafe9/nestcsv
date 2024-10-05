#pragma once

#include "Json.h"
#include "{{ .Prefix }}TableDataBase.generated.h"

USTRUCT(BlueprintType)
struct F{{ .Prefix }}TableDataBase
{
    GENERATED_BODY()

    virtual void Load(const TSharedPtr<FJsonObject>& JsonObject) = 0;
};
