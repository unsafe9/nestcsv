#pragma once

#include "Json.h"

USTRUCT(BlueprintType)
struct FNestComplexSKU
{
    GENERATED_BODY()
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString Type;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString ID;

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("Type"), Type);
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("ID"), ID);
    }
};

USTRUCT(BlueprintType)
struct FNestComplex
{
    GENERATED_BODY()
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FString> Tags;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FNestComplexSKU> SKU;

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
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
                const TSharedPtr<FJsonObject> *JsonObject;
                if (Item->TryGetObject(JsonObject))
                {
                    FNestComplexSKU FieldItem;
                    ObjItem.Load(JsonObject);
                    SKU.Add(FieldItem);
                }
            }
        }
    }
};
