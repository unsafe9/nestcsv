// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once

#include "NestTableBase.h"
#include "NestComplex.h"
#include "NestComplexTable.generated.h"

USTRUCT(BlueprintType)
struct FNestComplexTable : public FNestTableBase
{
    GENERATED_BODY()

    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
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
                if (Row->TryGetObject(RowValue))
                {
                    FNestComplex RowItem;
                    RowItem.Load(*RowValue);
                    Rows.Add(RowItem);
                }
            }
        }
    }
};