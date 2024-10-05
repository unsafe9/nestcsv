#pragma once

#include "Json.h"
#include "NestTableBase.generated.h"

USTRUCT(BlueprintType)
struct FNestTableBase
{
    GENERATED_BODY()

    virtual void Load(const TSharedPtr<FJsonValue>& JsonValue) = 0;
};
