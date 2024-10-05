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
                Tags.Add(Item->AsString());
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* SKUArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("SKU"), SKUArray))
        {
            for (const auto& Item : *SKUArray)
            {
                TSharedPtr<FJsonObject> Obj = Item->AsObject();
                if (Obj.IsValid())
                {
                    FNestComplexSKU ObjItem;
                    ObjItem.Load(Obj);
                    SKU.Add(ObjItem);
                }
            }
        }
    }
};
