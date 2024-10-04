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
        JsonObject.TryGetNumberField(TEXT("Int"), Int);
        JsonObject.TryGetNumberField(TEXT("Float"), Float);
        JsonObject.TryGetStringField(TEXT("Str"), Str);
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
        JsonObject.TryGetStringField(TEXT("Type"), Type);
        JsonObject.TryGetStringField(TEXT("ParamType"), ParamType);
        TSharedPtr<FJsonObject> ParamValueObject;
        if (JsonObject->TryGetObjectField(TEXT("ParamValue"), ParamValueObject))
        {
            ParamValue.Load(ParamValueObject);
        }
    }
};
