// Code generated by "nestcsv"; YOU CAN ONLY EDIT WITHIN THE TAGGED REGIONS!

#pragma once

#include "NestTableDataBase.h"

//NESTCSV:NESTSKU_EXTRA_INCLUDE_START

//NESTCSV:NESTSKU_EXTRA_INCLUDE_END

#include "NestSKU.generated.h"

USTRUCT(BlueprintType)
struct FNestSKU : public FNestTableDataBase
{
    GENERATED_BODY()
    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    FString Type;
    UPROPERTY(VisibleAnywhere, BlueprintReadOnly)
    FString ID;

    virtual void Load(const TSharedPtr<FJsonObject>& JsonObject) override
    {
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("Type"), Type);
        JsonObject.ToSharedRef()->TryGetStringField(TEXT("ID"), ID);

        OnLoad();
    }

    //NESTCSV:NESTSKU_EXTRA_BODY_START
    
    //NESTCSV:NESTSKU_EXTRA_BODY_END
};
