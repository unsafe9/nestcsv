// Code generated by "nestcsv"; DO NOT EDIT.

#pragma once
#include "NestComplexTable.hpp"
#include "NestTypesTable.hpp"

UCLASS(BlueprintType)
class UNestTableHolder : public UObject
{
    GENERATED_BODY()

public:
    UNestTableHolder() :
        Complex(MakeUnique<FNestComplexTable>()),
        Types(MakeUnique<FNestTypesTable>())
        {}
    virtual ~UNestTableHolder() override {}

    
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TUniquePtr<FNestComplexTable> Complex;
    UPROPERTY(EditAnywhere, BlueprintReadWrite)
    TUniquePtr<FNestTypesTable> Types;

    UFUNCTION(BlueprintCallable)
    TArray<FNestTableBase*> GetTables() const
    {
        TArray<FNestTableBase*> Tables;
        Tables.Add(Complex.Get());
        Tables.Add(Types.Get());
        return Tables;
    }

    UFUNCTION(BlueprintCallable)
    FNestTableBase* GetBySheetName(FString SheetName) const
    {
        if (SheetName == TEXT("complex"))
        {
            return Complex.Get();
        }
        if (SheetName == TEXT("types"))
        {
            return Types.Get();
        }
        return nullptr;
    }
};
