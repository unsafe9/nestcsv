#pragma once

#include "Json.h"

USTRUCT(BlueprintType)
struct FNestTypes
{
    GENERATED_BODY()
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
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("Int"), Int);
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("Long"), Long);
        JsonObject.ToSharedRef()->TryGetNumberField(TEXT("Float"), Float);
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("String"), String);
        FString TimeDtStr;
        if (JsonObject.ToSharedRef()->TryGetStringField(TEXT("Time"), TimeDtStr))
        {
            FDateTime::ParseIso8601(*TimeDtStr, Time);
        }
        JsonObject.ToSharedRef()->TryGetField(TEXT("Json"), Json);
        const TArray<TSharedPtr<FJsonValue>>* IntArrayArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("IntArray"), IntArrayArray))
        {
            for (const auto& Item : *IntArrayArray)
            {
                int32 FieldItem;
                if (Item->TryGetNumber(FieldItem))
                {
                    IntArray.Add(FieldItem);
                }
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* LongArrayArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("LongArray"), LongArrayArray))
        {
            for (const auto& Item : *LongArrayArray)
            {
                int64 FieldItem;
                if (Item->TryGetNumber(FieldItem))
                {
                    LongArray.Add(FieldItem);
                }
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* FloatArrayArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("FloatArray"), FloatArrayArray))
        {
            for (const auto& Item : *FloatArrayArray)
            {
                double FieldItem;
                if (Item->TryGetNumber(FieldItem))
                {
                    FloatArray.Add(FieldItem);
                }
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* StringArrayArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("StringArray"), StringArrayArray))
        {
            for (const auto& Item : *StringArrayArray)
            {
                FString FieldItem;
                if (Item->TryGetString(FieldItem))
                {
                    StringArray.Add(FieldItem);
                }
            }
        }
        const TArray<TSharedPtr<FJsonValue>>* TimeArrayArray = nullptr;
        if (JsonObject.ToSharedRef()->TryGetArrayField(TEXT("TimeArray"), TimeArrayArray))
        {
            for (const auto& Item : *TimeArrayArray)
            {
                FString DateTimeStr;
                if (Item->TryGetString(DateTimeStr))
                {
                    FDateTime DateTime;
                    if (FDateTime::ParseIso8601(DateTimeStr, DateTime))
                    {
                        TimeArray.Add(DateTime);
                    }
                }
            }
        }
    }
};
