#pragma once

#include "Json.h"
#include "NestTableDataBase.generated.h"

USTRUCT(BlueprintType)
struct FNestTableDataBase
{
    GENERATED_BODY()

    FNestTableDataBase() {}
    virtual ~FNestTableDataBase() {}

    virtual void Load(const TSharedPtr<FJsonObject>& JsonObject);
};
