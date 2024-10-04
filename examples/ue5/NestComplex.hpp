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
        JsonObject.TryGetStringField(TEXT("Type"), Type);
        JsonObject.TryGetStringField(TEXT("ID"), ID);
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
    TArray<NestComplexSKU> SKU;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<NestReward> Rewards;

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        JsonObject.TryGetNumberField(TEXT("ID"), ID);
        const TArray<TSharedPtr<FJsonValue>>* TagsArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("Tags"), TagsArray))
        {
            for (const auto& Item : *TagsArray)
            {
                Tags.Add(Item->AsString());
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* SKUArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("SKU"), SKUArray))
        {
            for (const auto& Item : *SKUArray)
            {
                TSharedPtr<FJsonObject> SKUObject = Item->AsObject();
                if (SKUObject.IsValid())
                {
                    TArray<NestComplexSKU> SKUItem;
                    SKUItem.Load(SKUObject);
                    SKU.Add(SKUItem);
                }
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* RewardsArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("Rewards"), RewardsArray))
        {
            for (const auto& Item : *RewardsArray)
            {
                TSharedPtr<FJsonObject> RewardsObject = Item->AsObject();
                if (RewardsObject.IsValid())
                {
                    TArray<NestReward> RewardsItem;
                    RewardsItem.Load(RewardsObject);
                    Rewards.Add(RewardsItem);
                }
            }
        }
    }
};
