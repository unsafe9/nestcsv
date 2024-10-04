#pragma once

#include "Json.h"


USTRUCT(BlueprintType)
struct FNestTypes
{
    GENERATED_USTRUCT_BODY()

public:
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    int32 Int;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    int64 Long;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    double Float;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FString String;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    FDateTime Time;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TSharedPtr<FJsonValue> Json;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<int32> IntArray;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<int64> LongArray;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<double> FloatArray;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FString> StringArray;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TArray<FDateTime> TimeArray;

    void Load(const TSharedPtr<FJsonObject>& JsonObject)
    {
        JsonObject.TryGetNumberField(TEXT("Int"), Int);
        JsonObject.TryGetNumberField(TEXT("Long"), Long);
        JsonObject.TryGetNumberField(TEXT("Float"), Float);
        JsonObject.TryGetStringField(TEXT("String"), String);
        FString TimeDtStr;
        if (JsonObject.TryGetStringField(TEXT("Time"), TimeDtStr))
        {
            Time = FDateTime::FromIso8601(TimeDtStr);
        }
        JsonObject.TryGetField(TEXT("Json"), Json);
        const TArray<TSharedPtr<FJsonValue>>* IntArrayArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("IntArray"), IntArrayArray))
        {
            for (const auto& Item : *IntArrayArray)
            {
                IntArray.Add(Item->AsNumber());
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* LongArrayArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("LongArray"), LongArrayArray))
        {
            for (const auto& Item : *LongArrayArray)
            {
                LongArray.Add(Item->AsNumber());
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* FloatArrayArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("FloatArray"), FloatArrayArray))
        {
            for (const auto& Item : *FloatArrayArray)
            {
                FloatArray.Add(Item->AsNumber());
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* StringArrayArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("StringArray"), StringArrayArray))
        {
            for (const auto& Item : *StringArrayArray)
            {
                StringArray.Add(Item->AsString());
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* TimeArrayArray = nullptr;
        if (JsonObject.TryGetArrayField(TEXT("TimeArray"), TimeArrayArray))
        {
            for (const auto& Item : *TimeArrayArray)
            {
                FString TimeArrayDtStr = Item->AsString();
                TimeArray.Add(FDateTime::FromIso8601(TimeArrayDtStr));
            }
        }
    }
};
