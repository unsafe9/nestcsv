#pragma once

#include "Json.h"


USTRUCT(BlueprintType)
struct FNestRewardParamValue
{
    GENERATED_USTRUCT_BODY()

public:
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    int32 Int;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    double Float;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString Str;

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("Int"), Int);
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("Float"), Float);
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("Str"), Str);
    }
};

USTRUCT(BlueprintType)
struct FNestReward
{
    GENERATED_USTRUCT_BODY()

public:
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString Type;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString ParamType;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    NestRewardParamValue ParamValue;

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("Type"), Type);
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("ParamType"), ParamType);
        TSharedPtr<FJsonObject> ParamValueObject;
        if (JsonObject.ToSharedRef()->TryGetObjectField(TEXT("ParamValue"), ParamValueObject))
        {
            ParamValue.Load(ParamValueObject);
        }
    }
};
