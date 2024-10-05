#pragma once

#include "NestTableBase.h"
#include "NestComplex.hpp"

USTRUCT(BlueprintType)
struct FNestComplexTable : public FNestTableBase
{
    GENERATED_BODY()

    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FNestComplex> Rows;
    
    virtual FString GetSheetName() const override
    {
        return TEXT("complex");
    }

    virtual void Load(const TSharedPtr<FJsonValue>& JsonValue) override
    {
        const TArray<TSharedPtr<FJsonValue>>* RowsArray = nullptr;
        if (JsonValue->TryGetArray(RowsArray))
        {
            for (const auto& Row : *RowsArray)
            {
                const TSharedPtr<FJsonObject> *RowValue = nullptr;
                if (Row.Value->TryGetObject(RowValue))
                {
                    FNestComplex RowItem;
                    RowItem.Load(*RowValue);
                    Rows.Add(RowItem);
                }
            }
        }
    }
};