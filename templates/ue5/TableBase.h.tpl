// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once

#include "Json.h"
#include "{{ .Prefix }}TableBase.generated.h"

USTRUCT(BlueprintType)
struct F{{ .Prefix }}TableBase
{
    GENERATED_BODY()

    F{{ .Prefix }}TableBase() {}
    virtual ~F{{ .Prefix }}TableBase() {}

    virtual FString GetSheetName() const { return TEXT(""); }
    virtual bool Load(const TSharedPtr<FJsonValue>& JsonValue) { return false; }
    virtual bool Load(const FString& JsonString)
    {
        TSharedPtr<FJsonValue> JsonValue;
        TSharedRef<TJsonReader<TCHAR>> JsonReader = TJsonReaderFactory<TCHAR>::Create(JsonString);
        if (FJsonSerializer::Deserialize(JsonReader, JsonValue) && JsonValue.IsValid())
        {
            return Load(JsonValue);
        }
        return false;
    }
};
