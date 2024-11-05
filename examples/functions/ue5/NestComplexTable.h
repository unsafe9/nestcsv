// Code generated by "nestcsv"; YOU CAN ONLY EDIT WITHIN THE TAGGED REGIONS!

#pragma once

#include "NestTableBase.h"
#include "NestComplex.h"

//NESTCSV:NESTCOMPLEX_EXTRA_INCLUDE_START
#include "CustomInclude.h"
//NESTCSV:NESTCOMPLEX_EXTRA_INCLUDE_END

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

    const FNestComplex* Find(int32 ID) const
    {
        return Rows.FindByPredicate([ID](const FNestComplex& Row) { return Row.ID == ID; });
    }
                        
    const FNestComplex& FindChecked(int32 ID) const
    {
        const FNestComplex* Row = Find(ID);
        check(Row != nullptr);
        return *Row;
    }

    //NESTCSV:NESTCOMPLEX_EXTRA_BODY_START
    void CustomFunction()
    {
        // Custom function body
    }
    //NESTCSV:NESTCOMPLEX_EXTRA_BODY_END
};