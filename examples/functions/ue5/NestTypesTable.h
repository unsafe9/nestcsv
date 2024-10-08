// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once

#include "NestTableBase.h"
#include "NestTypes.h"
#include "NestTypesTable.generated.h"

USTRUCT(BlueprintType)
struct FNestTypesTable : public FNestTableBase
{
    GENERATED_BODY()

    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    TMap<FString, FNestTypes> Rows;
    
    virtual FString GetSheetName() const override
    {
        return TEXT("types");
    }

    virtual void Load(const TSharedPtr<FJsonValue>& JsonValue) override
    {
        const TSharedPtr<FJsonObject>* RowsMap = nullptr;
        if (JsonValue->TryGetObject(RowsMap))
        {
            for (const auto& Row : (*RowsMap)->Values)
            {
                const TSharedPtr<FJsonObject> *RowValue = nullptr;
                if (Row.Value->TryGetObject(RowValue))
                {
                    FNestTypes RowItem;
                    RowItem.Load(*RowValue);
                    Rows.Add(Row.Key, RowItem);
                }
            }
        }
    }
};