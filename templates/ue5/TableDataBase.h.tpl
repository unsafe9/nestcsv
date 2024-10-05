#pragma once

#include "Json.h"
#include "{{ .Prefix }}TableDataBase.generated.h"

USTRUCT(BlueprintType)
struct F{{ .Prefix }}TableDataBase
{
    GENERATED_BODY()

    F{{ .Prefix }}TableDataBase() {}
    virtual ~F{{ .Prefix }}TableDataBase() {}

    virtual void Load(const TSharedPtr<FJsonObject>& JsonObject);
};
