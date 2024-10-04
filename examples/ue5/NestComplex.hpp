#pragma once

#include "Json.h"
#include "NestReward.hpp"


USTRUCT(BlueprintType)
struct FNestComplexSKU
{
    GENERATED_USTRUCT_BODY()

public:
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
    GENERATED_USTRUCT_BODY()

public:
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    int32 ID;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FString> Tags;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FNestComplexSKU> SKU;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FNestReward> Rewards;

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("ID"), ID);
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
        const TArray<TSharedPtr<FJsonValue>>* RewardsArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("Rewards"), RewardsArray))
        {
            for (const auto& Item : *RewardsArray)
            {
                TSharedPtr<FJsonObject> Obj = Item->AsObject();
                if (Obj.IsValid())
                {
                    FNestReward ObjItem;
                    ObjItem.Load(Obj);
                    Rewards.Add(ObjItem);
                }
            }
        }
    }
};
