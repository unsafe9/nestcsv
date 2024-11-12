// Code generated by "nestcsv"; YOU CAN ONLY EDIT WITHIN THE TAGGED REGIONS!

#pragma once

#include "NestTableBase.h"
#include "NestTypes.h"

//NESTCSV:NESTTYPES_EXTRA_INCLUDE_START

//NESTCSV:NESTTYPES_EXTRA_INCLUDE_END

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

    virtual bool Load(const TSharedPtr<FJsonValue>& JsonValue) override
    {
        if (!JsonValue.IsValid()) return false;
        TMap<FString, FNestTypes> _Result;

        const TSharedPtr<FJsonObject>* RowsMap = nullptr;
        if (!JsonValue->TryGetObject(RowsMap)) return false;
        for (const auto& Row : (*RowsMap)->Values)
        {
            const TSharedPtr<FJsonObject> *RowValue = nullptr;
            if (!Row.Value->TryGetObject(RowValue)) return false;
            FNestTypes RowItem;
            if (!RowItem.Load(*RowValue)) return false;
            _Result.Add(Row.Key, RowItem);
        }

        Rows = MoveTemp(_Result);
        return true;
    }

    const FNestTypes* Find(int32 ID) const
    {
        return Rows.Find(FString::FromInt(ID));
    }
                        
    const FNestTypes& FindChecked(int32 ID) const
    {
        const FNestTypes* Row = Find(ID);
        check(Row != nullptr);
        return *Row;
    }

    //NESTCSV:NESTTYPES_EXTRA_BODY_START
    
    //NESTCSV:NESTTYPES_EXTRA_BODY_END
};