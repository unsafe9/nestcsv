// Code generated by "nestcsv"; YOU CAN ONLY EDIT WITHIN THE TAGGED REGIONS!

#pragma once

#include "NestTableDataBase.h"
#include "NestSKU.h"

//nestcsv:additional_include_start
#include "CustomInclude.h"
//nestcsv:additional_include_end

#include "NestComplex.generated.h"

USTRUCT(BlueprintType)
struct FNestComplex : public FNestTableDataBase
{
    GENERATED_BODY()
    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    TArray<FString> Tags;
    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    TArray<FNestSKU> SKU;

    virtual void Load(const TSharedPtr<FJsonObject>& JsonObject) override
    {
        const TArray<TSharedPtr<FJsonValue>>* TagsArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("Tags"), TagsArray))
        {
            for (const auto& Item : *TagsArray)
            {
                FString FieldItem;
                if (Item->TryGetString(FieldItem))
                {
                    Tags.Add(FieldItem);
                }
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* SKUArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("SKU"), SKUArray))
        {
            for (const auto& Item : *SKUArray)
            {
                const TSharedPtr<FJsonObject> *ObjPtr = nullptr;
                if (Item->TryGetObject(ObjPtr))
                {
                    FNestSKU FieldItem;
                    FieldItem.Load(*ObjPtr);
                    SKU.Add(FieldItem);
                }
            }
        }
    }

    //nestcsv:additional_struct_body_start
    void CustomFunction()
    {
        // Custom function body
    }
    //nestcsv:additional_struct_body_end
};
