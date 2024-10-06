// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once

#include "Json.h"
#include "NestTableDataBase.generated.h"

USTRUCT(BlueprintType)
struct FNestTableDataBase
{
    GENERATED_BODY()

    FNestTableDataBase() {}
    virtual ~FNestTableDataBase() {}

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
