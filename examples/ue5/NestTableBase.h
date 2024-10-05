#pragma once

#include "Json.h"
#include "NestTableBase.generated.h"

USTRUCT(BlueprintType)
struct FNestTableBase
{
    GENERATED_BODY()

    FNestTableBase() {}
    virtual ~FNestTableBase() {}

    virtual FString GetSheetName() const;
    virtual void Load(const TSharedPtr<FJsonValue>& JsonValue);
};
