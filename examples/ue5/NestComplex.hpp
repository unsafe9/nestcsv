#pragma once

#include "NestTableDataBase.h"

USTRUCT(BlueprintType)
struct FNestComplexSKU : public FNestTableDataBase
{
    GENERATED_BODY()
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString Type;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString ID;

    virtual void Load(const TSharedPtr<FJsonObject>& JsonObject) override
    {
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("Type"), Type);
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("ID"), ID);
    }
};

USTRUCT(BlueprintType)
struct FNestComplex : public FNestTableDataBase
{
    GENERATED_BODY()
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FString> Tags;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FNestComplexSKU> SKU;

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
                    FNestComplexSKU FieldItem;
                    FieldItem.Load(*ObjPtr);
                    SKU.Add(FieldItem);
                }
            }
        }
    }
};
