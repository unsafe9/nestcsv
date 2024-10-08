// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once

#include "Json.h"
#include "{{ .Prefix }}TableDataBase.generated.h"

USTRUCT(BlueprintType)
struct F{{ .Prefix }}TableDataBase
{
    GENERATED_BODY()

    F{{ .Prefix }}TableDataBase() {}
    virtual ~F{{ .Prefix }}TableDataBase() {}

    virtual void Load(const TSharedPtr<FJsonObject>& JsonObject) {}
    virtual void Load(const FString& JsonString)
    {
        TSharedPtr<FJsonObject> JsonObject;
        TSharedRef<TJsonReader<TCHAR>> JsonReader = TJsonReaderFactory<TCHAR>::Create(JsonString);
        if (FJsonSerializer::Deserialize(JsonReader, JsonObject) && JsonObject.IsValid())
        {
            Load(JsonObject);
        }
    }
};
