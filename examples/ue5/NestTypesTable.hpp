#pragma once

#include "NestTableBase.h"
#include "NestTypes.hpp"

USTRUCT(BlueprintType)
struct FNestTypesTable : public FNestTableBase
{
    GENERATED_BODY()

    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TMap<FString, FNestTypes> Rows;

    void Load(const TSharedPtr<FJsonValue>& JsonValue) override
    {
        const TSharedPtr<FJsonObject>* RowsMap;
        if (JsonValue->TryGetObject(RowsMap))
        {
            for (const auto& Row : (*RowsMap)->Values)
            {
                FNestTypes RowItem;
                RowItem.Load(Row.Value);
                Rows.Add(Row.Key, RowItem);
            }
        }
    }
};