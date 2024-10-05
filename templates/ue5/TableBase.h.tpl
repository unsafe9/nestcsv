#pragma once

#include "Json.h"
#include "{{ .Prefix }}TableBase.generated.h"

USTRUCT(BlueprintType)
struct F{{ .Prefix }}TableBase
{
    GENERATED_BODY()

    virtual void Load(const TSharedPtr<FJsonValue>& JsonValue) = 0;
};
