#pragma once

#include "NestTableBase.h"
#include "NestComplex.hpp"

USTRUCT(BlueprintType)
struct FNestComplexTable : public FNestTableBase
{
    GENERATED_BODY()

    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FNestComplex> Rows;

    void Load(const TSharedPtr<FJsonValue>& JsonValue) override
    {
        const TArray<TSharedPtr<FJsonValue>>* RowsArray;
        if (JsonValue->TryGetArray(RowsArray))
        {
            for (const auto& Row : *RowsArray)
            {
                FNestComplex RowItem;
                RowItem.Load(Row);
                Rows.Add(RowItem);
            }
        }
    }
};