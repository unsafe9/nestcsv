// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once

#include "NestTableDataBase.h"
#include "NestSKU.h"
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
};