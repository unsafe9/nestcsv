// Code generated by "nestcsv"; YOU CAN ONLY EDIT WITHIN THE TAGGED REGIONS!

#pragma once

#include "NestTableDataBase.h"
#include "NestSKU.h"

//NESTCSV:NESTCOMPLEX_EXTRA_INCLUDE_START
#include "CustomInclude.h"
//NESTCSV:NESTCOMPLEX_EXTRA_INCLUDE_END

#include "NestComplex.generated.h"

USTRUCT(BlueprintType)
struct FNestComplex : public FNestTableDataBase
{
    GENERATED_BODY()
    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    TArray<FString> Tags;
    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    TArray<FNestSKU> SKU;

    virtual bool Load(const TSharedPtr<FJsonObject>& JsonObject) override
    {
        if (!JsonObject.IsValid()) return false;
        FNestComplex _Result;

        {
            const TArray<TSharedPtr<FJsonValue>>* TagsArray = nullptr;
            if (!JsonObject.ToSharedRef()->TryGetArrayField(TEXT("Tags"), TagsArray)) return false;
            for (const auto& Item : *TagsArray)
            {
                FString FieldItem;
                if (!Item->TryGetString(FieldItem)) return false;
                _Result.Tags.Add(FieldItem);
            }
        }
        {
            const TArray<TSharedPtr<FJsonValue>>* SKUArray = nullptr;
            if (!JsonObject.ToSharedRef()->TryGetArrayField(TEXT("SKU"), SKUArray)) return false;
            for (const auto& Item : *SKUArray)
            {
                const TSharedPtr<FJsonObject> *ObjPtr = nullptr;
                if (!Item->TryGetObject(ObjPtr)) return false;
                FNestSKU FieldItem;
                FieldItem.Load(*ObjPtr);
                _Result.SKU.Add(FieldItem);
            }
        }

        *this = MoveTemp(_Result);
        return true;
    }

    //NESTCSV:NESTCOMPLEX_EXTRA_BODY_START
    void CustomFunction()
    {
        // Custom function body
    }
    //NESTCSV:NESTCOMPLEX_EXTRA_BODY_END
};
