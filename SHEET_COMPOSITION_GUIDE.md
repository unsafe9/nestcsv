# Sheet Composition Guide

## Table of Contents
- [Overview](#overview)
- [Sheet Structure](#sheet-structure)
- [Row Definitions](#row-definitions)
- [Field Types](#field-types)
- [Composition Patterns](#composition-patterns)
- [Metadata Options](#metadata-options)
- [Tag System](#tag-system)
- [Rules and Constraints](#rules-and-constraints)
- [Examples](#examples)

---

## Overview

**nestcsv** enables sophisticated data modeling through **sheet composition** - a system that converts structured CSV/Excel sheets into nested JSON and type-safe generated code. Sheet composition allows you to:

- Define complex nested data structures using dot notation
- Create arrays across multiple rows (Multi-Line Arrays)
- Generate reusable struct definitions
- Filter data by tags for different targets (client/server)
- Control output format through metadata

---

## Sheet Structure

Every sheet must follow a strict 5-row schema structure before data rows:

```
Row 0: Metadata (query string)
Row 1: Field Tags (comma-separated)
Row 2: Field Names (with dot notation for nesting)
Row 3: Field Types (with [] prefix for arrays)
Row 4: Comments (ignored, typically empty or comments)
Row 5+: Data rows
```

### Minimal Example

```csv
as_map=false
all,all,all
ID,Name,Level
int,string,int

1,Alice,10
2,Bob,15
```

### Visual Structure

```
┌──────────────────────────────────────────────────────┐
│ Row 0: as_map=false&sort_asc_by=ID                  │ ← Metadata
├──────────────────────────────────────────────────────┤
│ Row 1: all,client,server                            │ ← Tags
├──────────────────────────────────────────────────────┤
│ Row 2: ID,ClientField,ServerField                   │ ← Field Names
├──────────────────────────────────────────────────────┤
│ Row 3: int,string,string                            │ ← Field Types
├──────────────────────────────────────────────────────┤
│ Row 4: (empty or comments)                          │ ← Comments
├──────────────────────────────────────────────────────┤
│ Row 5: 1,ClientValue1,ServerValue1                  │ ← Data
│ Row 6: 2,ClientValue2,ServerValue2                  │ ← Data
│ ...                                                  │
└──────────────────────────────────────────────────────┘
```

---

## Row Definitions

### Row 0: Metadata Query String

A URL query-string format that controls output behavior:

**Syntax:** `key1=value1&key2=value2&key3=value3`

**Available Options:**

| Option | Type | Description | Example |
|--------|------|-------------|---------|
| `as_map` | bool | Return data as map (using ID as key) instead of array | `as_map=true` |
| `sort_asc_by` | string | Sort output array by field in ascending order | `sort_asc_by=ID` |
| `sort_desc_by` | string | Sort output array by field in descending order | `sort_desc_by=Level` |
| `struct` | map | Map field identifiers to named struct types (reusable) | `struct=Rewards:Reward` |

**Examples:**
```
as_map=false&sort_asc_by=ID
as_map=true
struct=Rewards:Reward&struct=/.*SKU.*/:SKU
sort_desc_by=CreatedAt
```

### Row 1: Field Tags

Comma-separated tags that control which fields are included in different outputs.

**Syntax:** `tag1,tag2,tag3` or empty

**Common Tags:**
- `all` - included in all outputs
- `client` - only included in client-tagged outputs
- `server` - only included in server-tagged outputs
- Custom tags (any identifier)

**Examples:**
```csv
all,client,server,all
ID,ClientData,ServerData,SharedData
```

**Tag Filtering:**
Outputs specify which tags to include:
```yaml
outputs:
  - tags: [server]      # Only includes columns tagged 'server' or 'all'
  - tags: [client]      # Only includes columns tagged 'client' or 'all'
  - tags: [client, server]  # Includes all columns
```

### Row 2: Field Names

Field names define the structure of your data using dot notation for nesting.

**Syntax Rules:**

1. **Simple Field:** `FieldName`
2. **Nested Field:** `Parent.Child.GrandChild`
3. **Multi-Line Array:** `[]ArrayField` or `[]Parent.Child`
4. **Comments:** Fields starting with `#` are ignored
5. **Empty Columns:** Empty field names are skipped

**Special Considerations:**

- **Column 0 (ID Column):** Must be a non-nested simple field (cannot contain dots)
- **Case Sensitive:** `ID`, `Id`, and `id` are different fields
- **No Whitespace:** Leading/trailing spaces are trimmed, internal spaces preserved

**Examples:**
```csv
ID                    → Simple field
User.Name             → Nested: User { Name }
User.Profile.Age      → Deeply nested: User { Profile { Age } }
[]Items.Type          → Multi-line array of structs
[]Tags                → Multi-line array of primitives
#DebugField           → Ignored column
```

### Row 3: Field Types

Field types define the data type and array behavior.

**Syntax:**
- Primitive: `typename`
- Cell Array: `[]typename`

**Available Types:**

| Type | Description | Zero Value | Example |
|------|-------------|------------|---------|
| `int` | 32-bit integer | `0` | `42`, `-10` |
| `long` | 64-bit integer | `0` | `9999999999` |
| `float` | 64-bit floating point | `0.0` | `3.14`, `-0.5` |
| `bool` | Boolean | `false` | `true`, `false` |
| `string` | Text | `""` (empty string) | `hello`, `multi word` |
| `time` | DateTime | `0001-01-01T00:00:00Z` | `2024-09-30 11:00:00` |
| `json` | Raw JSON | `null` | `{"key":"value"}` |
| `struct` | Nested object | (inferred from dot notation) | N/A - auto-detected |

**Array Types:**

Cell arrays use `[]` prefix: `[]int`, `[]string`, `[]float`, etc.

**Notes:**
- `json` type cannot be an array (`[]json` is invalid)
- `struct` type is automatically inferred when field has children via dot notation
- Time format: `YYYY-MM-DD HH:MM:SS` (parses to RFC3339/ISO8601)

### Row 4: Comments Row

This row is typically empty or contains comments for human readers. The parser completely ignores this row.

**Common Uses:**
```csv
,,comments!,,notes here,,
```

### Rows 5+: Data Rows

Data rows contain your actual data values.

**Special Behaviors:**

1. **Empty ID:** Rows with empty ID in column 0 are skipped
2. **Comment Rows:** Rows where ID starts with `#` are skipped
3. **Empty Cells:** Empty cells become zero values for their type
4. **Duplicate IDs:** Only allowed when Multi-Line Arrays exist (see composition patterns)

**Examples:**
```csv
1,Alice,25          ← Normal row
2,Bob,30            ← Normal row
#3,Charlie,35       ← Skipped (commented)
,Dave,40            ← Skipped (empty ID)
```

---

## Field Types

### Primitive Types

#### Integer (`int`)
- **Description:** 32-bit signed integer
- **Range:** -2,147,483,648 to 2,147,483,647
- **Zero Value:** `0`
- **Examples:** `1`, `42`, `-10`, `999999`

#### Long (`long`)
- **Description:** 64-bit signed integer
- **Range:** -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
- **Zero Value:** `0`
- **Examples:** `9999999999`, `-1000000000000`

#### Float (`float`)
- **Description:** 64-bit floating point
- **Zero Value:** `0.0`
- **Examples:** `3.14`, `-0.5`, `0.1`, `1.0e10`

#### Boolean (`bool`)
- **Description:** True or false value
- **Zero Value:** `false`
- **Valid Values:** `true`, `false`, `1`, `0`, `t`, `f`, `T`, `F`, `TRUE`, `FALSE`

#### String (`string`)
- **Description:** Text data
- **Zero Value:** `""` (empty string)
- **Examples:** `hello`, `multi word text`, `123` (as text)

#### Time (`time`)
- **Description:** DateTime value
- **Format:** `YYYY-MM-DD HH:MM:SS` (24-hour format)
- **Output:** RFC3339/ISO8601 format (`2024-09-30T11:00:00Z`)
- **Zero Value:** `0001-01-01T00:00:00Z`
- **Examples:** `2024-09-30 11:00:00`, `2023-12-25 15:30:45`

#### JSON (`json`)
- **Description:** Raw JSON data embedded in cell
- **Zero Value:** `null`
- **Format:** Valid JSON string (double quotes for strings)
- **Examples:** `{"key":"value"}`, `[1,2,3]`, `{"nested":{"deep":"value"}}`
- **Note:** Cannot be used as array type (`[]json` is invalid)

### Array Types

Arrays come in two forms: **Cell Arrays** and **Multi-Line Arrays**.

#### Cell Arrays

Arrays defined within a single cell using comma-separated values.

**Syntax:**
- Type: `[]typename` (e.g., `[]int`, `[]string`)
- Value: `item1,item2,item3`

**Examples:**

| Field Type | Cell Value | Parsed Result |
|------------|-----------|---------------|
| `[]int` | `1,2,3` | `[1, 2, 3]` |
| `[]string` | `apple,banana,cherry` | `["apple", "banana", "cherry"]` |
| `[]float` | `0.1,0.2,0.3` | `[0.1, 0.2, 0.3]` |
| `[]time` | `2024-01-01 10:00:00,2024-01-02 10:00:00` | `[time1, time2]` |
| `[]string` | (empty) | `[]` (empty array) |

**Limitations:**
- Cannot contain commas in individual values
- Cannot create arrays of structs (use Multi-Line Arrays instead)
- Cannot create arrays of JSON

---

## Composition Patterns

Composition patterns are the core power of nestcsv, enabling complex data structures from simple CSV layouts.

### 1. Nested Structures (Dot Notation)

Create hierarchical objects using dot notation in field names.

**Syntax:** `Parent.Child.GrandChild`

**Example:**

```csv
as_map=false
all,all,all,all
ID,User.Name,User.Profile.Age,User.Profile.City
int,string,int,string

1,Alice,25,NYC
2,Bob,30,LA
```

**Output:**
```json
[
  {
    "ID": 1,
    "User": {
      "Name": "Alice",
      "Profile": {
        "Age": 25,
        "City": "NYC"
      }
    }
  },
  {
    "ID": 2,
    "User": {
      "Name": "Bob",
      "Profile": {
        "Age": 30,
        "City": "LA"
      }
    }
  }
]
```

**How It Works:**
- Each dot creates a new nested level
- Fields with the same parent path are grouped together
- `User.Name` and `User.Profile.Age` both contribute to the same `User` object
- Parser builds a tree structure automatically

**Use Cases:**
- Grouping related fields
- Modeling entity relationships
- Creating clean, hierarchical data structures

---

### 2. Multi-Line Arrays (MLA)

Create arrays by repeating the same ID across multiple rows. This is the most powerful composition pattern.

**Syntax:** Prefix field name with `[]`
- `[]ArrayField`
- `[]Parent.Child`

**Key Concept:** Rows with the same ID value compose into a single record, with MLA fields creating array elements.

#### Example 1: Simple Multi-Line Array

```csv
as_map=false
all,all,all
ID,Name,[]Tags
int,string,string

1,Alice,gold
1,,platinum
2,Bob,silver
```

**Output:**
```json
[
  {
    "ID": 1,
    "Name": "Alice",
    "Tags": ["gold", "platinum"]
  },
  {
    "ID": 2,
    "Name": "Bob",
    "Tags": ["silver"]
  }
]
```

**How It Works:**
1. Parser groups rows by ID
2. For ID=1, there are 2 rows
3. `Tags` field has `[]` prefix, so creates array with 2 elements
4. `Name` field (non-MLA) only uses first occurrence (empty cells in subsequent rows are ignored)

#### Example 2: Array of Structs

```csv
as_map=false
all,all,all,all
ID,Name,[]Items.Type,[]Items.Quantity
int,string,string,int

1,Order1,Sword,2
1,,Shield,1
2,Order2,Potion,5
```

**Output:**
```json
[
  {
    "ID": 1,
    "Name": "Order1",
    "Items": [
      {"Type": "Sword", "Quantity": 2},
      {"Type": "Shield", "Quantity": 1}
    ]
  },
  {
    "ID": 2,
    "Name": "Order2",
    "Items": [
      {"Type": "Potion", "Quantity": 5}
    ]
  }
]
```

**How It Works:**
1. `[]Items.Type` and `[]Items.Quantity` share same MLA parent (`Items`)
2. Each row creates one element in the array
3. Both fields contribute to the same array element

#### Example 3: Multiple Independent Arrays

```csv
as_map=false
all,all,all,all,all
ID,[]SKU.Type,[]SKU.ID,[]Rewards.Type,[]Rewards.Amount
int,string,string,string,int

1,Google,IAP_1,Gold,100
1,,,Silver,50
2,Apple,IAP_2,Gold,200
2,Apple,IAP_3,Gold,300
```

**Output:**
```json
[
  {
    "ID": 1,
    "SKU": [
      {"Type": "Google", "ID": "IAP_1"}
    ],
    "Rewards": [
      {"Type": "Gold", "Amount": 100},
      {"Type": "Silver", "Amount": 50}
    ]
  },
  {
    "ID": 2,
    "SKU": [
      {"Type": "Apple", "ID": "IAP_2"},
      {"Type": "Apple", "ID": "IAP_3"}
    ],
    "Rewards": [
      {"Type": "Gold", "Amount": 200},
      {"Type": "Gold", "Amount": 300}
    ]
  }
]
```

**How It Works:**
1. `SKU` and `Rewards` are independent MLA fields
2. For ID=1: Only 1 SKU element (row 2 has all SKU fields empty), but 2 Rewards elements
3. Empty cell detection: If all cells in an MLA group are empty, that array element is skipped

#### Example 4: Real-World Complex Example

From `complex.csv`:

```csv
as_map=false&sort_asc_by=ID&struct=Rewards:Reward&struct=/.*SKU.*/:SKU
server,client,client,client,server,server,server,server,server
ID,Tags,[]SKU.Type,[]SKU.ID,[]Rewards.Type,[]Rewards.ParamValue.Str,[]Rewards.ParamType,[]Rewards.ParamValue.Int,[]Rewards.ParamValue.Float
int,[]string,string,string,string,string,string,int,float

1,"gold,package",Google,IAP_Google_1,Gold,,Int,10,
1,,,,Gear,Weapon,Str,,
2,dollar,Google,IAP_Google_2,Dollar,,Float,,0.5
2,,Apple,IAP_Apple_2,Dollar,,Float,,0.8
```

**Output (with client tag):**
```json
[
  {
    "ID": 1,
    "Tags": ["gold", "package"],
    "SKU": [
      {"Type": "Google", "ID": "IAP_Google_1"}
    ]
  },
  {
    "ID": 2,
    "Tags": ["dollar"],
    "SKU": [
      {"Type": "Google", "ID": "IAP_Google_2"},
      {"Type": "Apple", "ID": "IAP_Apple_2"}
    ]
  }
]
```

**Output (with server tag):**
```json
[
  {
    "ID": 1,
    "Rewards": [
      {
        "Type": "Gold",
        "ParamType": "Int",
        "ParamValue": {"Str": "", "Int": 10, "Float": 0}
      },
      {
        "Type": "Gear",
        "ParamType": "Str",
        "ParamValue": {"Str": "Weapon", "Int": 0, "Float": 0}
      }
    ]
  },
  {
    "ID": 2,
    "Rewards": [
      {
        "Type": "Dollar",
        "ParamType": "Float",
        "ParamValue": {"Str": "", "Int": 0, "Float": 0.5}
      },
      {
        "Type": "Dollar",
        "ParamType": "Float",
        "ParamValue": {"Str": "", "Int": 0, "Float": 0.8}
      }
    ]
  }
]
```

**Key Observations:**
1. Different outputs for different tags
2. Complex nested structures within arrays
3. Multiple MLA fields create independent arrays
4. Named struct mapping (see next section)

---

### 3. Named Struct Mapping

Map field identifiers to named struct types for code generation reusability.

**Problem:** Without named structs, each array of structs generates anonymous struct types in code, leading to duplication.

**Solution:** Use `struct=FieldId:StructName` metadata to create reusable struct definitions.

**Syntax:**
```
struct=FieldIdentifier:StructName
struct=/regex/:StructName
```

**Examples:**

#### Exact Match
```
struct=Rewards:Reward
```
Maps field `Rewards` to struct name `Reward`.

#### Regex Match
```
struct=/.*SKU.*/:SKU
```
Maps any field containing "SKU" to struct name `SKU`.

#### Multiple Mappings
```
struct=Rewards:Reward&struct=/.*SKU.*/:SKU
```

**Field Identifier Rules:**

1. Use the full path without parent structs
   ```
   Field: A.B.C.D
   Field Identifier: C.D (omit parent A.B if C is a named struct)
   ```

2. For MLA fields, omit the `[]` prefix
   ```
   Field: []Rewards.Type
   Identifier: Rewards
   ```

**Example:**

```csv
struct=Rewards:Reward
all,all,all
ID,[]Rewards.Type,[]Rewards.Amount
int,string,int

1,Gold,100
1,Silver,50
```

**Generated Code (Go):**
```go
// Reward - reusable named struct
type Reward struct {
    Type   string `json:"Type"`
    Amount int    `json:"Amount"`
}

// Table struct uses named type
type ComplexRow struct {
    ID      int      `json:"ID"`
    Rewards []Reward `json:"Rewards"`
}
```

**Without Named Struct:**
```go
type ComplexRow struct {
    ID      int `json:"ID"`
    Rewards []struct {  // Anonymous struct - not reusable
        Type   string `json:"Type"`
        Amount int    `json:"Amount"`
    } `json:"Rewards"`
}
```

**Benefits:**
- Code reusability across tables
- Better type checking
- Cleaner generated code
- Shared struct definitions

---

### 4. Combining Patterns

You can combine all composition patterns together.

**Example:**

```csv
as_map=false&struct=Items:Item
all,all,all,all,all,all
ID,User.Name,User.Level,[]Items.Type,[]Items.Stats.Attack,[]Items.Stats.Defense
int,string,int,string,int,int

1,Alice,10,Sword,50,10
1,,,Shield,0,50
2,Bob,15,Potion,0,0
```

**Output:**
```json
[
  {
    "ID": 1,
    "User": {
      "Name": "Alice",
      "Level": 10
    },
    "Items": [
      {
        "Type": "Sword",
        "Stats": {"Attack": 50, "Defense": 10}
      },
      {
        "Type": "Shield",
        "Stats": {"Attack": 0, "Defense": 50}
      }
    ]
  },
  {
    "ID": 2,
    "User": {
      "Name": "Bob",
      "Level": 15
    },
    "Items": [
      {
        "Type": "Potion",
        "Stats": {"Attack": 0, "Defense": 0}
      }
    ]
  }
]
```

**This combines:**
- Nested structs (`User.Name`, `User.Level`)
- Multi-line arrays (`[]Items`)
- Nested structs within arrays (`[]Items.Stats.Attack`)
- Named struct mapping (`Items:Item`)

---

## Metadata Options

Metadata in Row 0 controls output behavior using query string format.

### `as_map`

**Type:** `boolean`
**Default:** `false`

Controls whether output is an array or a map (keyed by ID).

**Syntax:** `as_map=true` or `as_map=false`

**Example:**

```csv
as_map=true
all,all,all
ID,Name,Level
int,string,int

1,Alice,10
2,Bob,15
```

**Output:**
```json
{
  "1": {
    "ID": 1,
    "Name": "Alice",
    "Level": 10
  },
  "2": {
    "ID": 2,
    "Name": "Bob",
    "Level": 15
  }
}
```

**Use Cases:**
- Fast lookups by ID
- Dictionary/map data structures
- When order doesn't matter

**Constraints:**
- Mutually exclusive with `sort_asc_by` and `sort_desc_by`

---

### `sort_asc_by` / `sort_desc_by`

**Type:** `string` (field name)
**Default:** (no sorting)

Sort output array by specified field.

**Syntax:**
- `sort_asc_by=FieldName` (ascending)
- `sort_desc_by=FieldName` (descending)

**Example:**

```csv
as_map=false&sort_desc_by=Level
all,all,all
ID,Name,Level
int,string,int

1,Alice,10
2,Bob,25
3,Charlie,15
```

**Output:**
```json
[
  {"ID": 2, "Name": "Bob", "Level": 25},
  {"ID": 3, "Name": "Charlie", "Level": 15},
  {"ID": 1, "Name": "Alice", "Level": 10}
]
```

**Valid Sort Types:**
- `int`, `long`, `float`, `string`, `time`

**Invalid Sort Types:**
- Arrays (cell or multi-line)
- `bool`
- `json`
- Nested fields with dots

**Constraints:**
- Cannot use both `sort_asc_by` and `sort_desc_by`
- Cannot use with `as_map=true`
- Field must exist in the table
- Field cannot be an array type

---

### `struct`

**Type:** `map[string]string`
**Format:** `struct=FieldId:StructName`

Map field identifiers to named struct types for code generation.

**Syntax:**
```
struct=FieldId:StructName
struct=/regex/:StructName
```

**Multiple Mappings:**
```
struct=Rewards:Reward&struct=/.*SKU.*/:SKU&struct=Items:Item
```

**Field Identifier Rules:**

1. **For MLA fields:** Use field name without `[]` prefix
   ```
   Field: []Rewards.Type
   Identifier: Rewards
   Mapping: struct=Rewards:Reward
   ```

2. **For nested fields:** Omit parent struct names if field is a named struct
   ```
   Field: Order.Items.Name
   If Items is named struct: struct=Items:Item
   Not: struct=Order.Items:Item
   ```

3. **Regex patterns:** Enclose in `/regex/`
   ```
   struct=/.*SKU.*/:SKU    → matches SKU, PlayerSKU, SKUID, etc.
   struct=/Reward.*/:Reward → matches RewardWin, RewardLose, etc.
   ```

**Example:**

```csv
struct=Rewards:Reward&struct=/.*Bonus/:BonusData
all,all,all,all,all
ID,[]Rewards.Type,[]Rewards.Amount,[]WinBonus.Type,[]LoseBonus.Type
int,string,int,string,string

1,Gold,100,XP,HP
1,Silver,50,MP,
```

**Generated Structs:**
```go
type Reward struct {
    Type   string
    Amount int
}

type BonusData struct {
    Type string
}

type TableRow struct {
    ID         int
    Rewards    []Reward      // Uses named struct
    WinBonus   []BonusData   // Uses named struct (regex match)
    LoseBonus  []BonusData   // Uses named struct (regex match)
}
```

---

## Tag System

Tags provide powerful filtering to generate different outputs from the same source data.

### How Tags Work

1. **Define tags in Row 1** for each column
2. **Specify tags in output/codegen config**
3. **Parser includes columns** where tag matches

**Tag Matching Rule:**
- Column is included if **ANY** tag in the column matches **ANY** tag in the configuration
- Empty tag cell means column has no tags (included only if config has no tags)

### Tag Examples

#### Example 1: Basic Tag Filtering

```csv
,
all,client,server
ID,ClientData,ServerData
int,string,string

1,ClientValue1,ServerValue1
2,ClientValue2,ServerValue2
```

**Config:**
```yaml
outputs:
  - tags: [client]
    json:
      root_dir: ./client
  - tags: [server]
    json:
      root_dir: ./server
```

**Client Output:**
```json
[
  {"ID": 1, "ClientData": "ClientValue1"},
  {"ID": 2, "ClientData": "ClientValue2"}
]
```

**Server Output:**
```json
[
  {"ID": 1, "ServerData": "ServerValue1"},
  {"ID": 2, "ServerData": "ServerValue2"}
]
```

#### Example 2: Shared Fields

```csv
,
all,client,server,all
ID,ClientData,ServerData,SharedData
int,string,string,string

1,C1,S1,Shared1
```

**Client Output:**
```json
[
  {"ID": 1, "ClientData": "C1", "SharedData": "Shared1"}
]
```

**Server Output:**
```json
[
  {"ID": 1, "ServerData": "S1", "SharedData": "Shared1"}
]
```

#### Example 3: Multiple Tags Per Column

```csv
,
all,"client,server",server
ID,SharedField,ServerOnlyField
int,string,string

1,Shared,ServerValue
```

**Explanation:**
- Column 1: Tagged with both `client` and `server`
- Included in both client and server outputs

#### Example 4: Custom Tags

```csv
,
all,premium,debug,premium
ID,PremiumFeature,DebugInfo,AnotherPremium
int,string,string,string

1,Gold,DebugData,VIP
```

**Config:**
```yaml
outputs:
  - tags: [premium]
    json:
      root_dir: ./premium
  - tags: [debug]
    json:
      root_dir: ./debug
```

**Premium Output:**
```json
[
  {"ID": 1, "PremiumFeature": "Gold", "AnotherPremium": "VIP"}
]
```

**Debug Output:**
```json
[
  {"ID": 1, "DebugInfo": "DebugData"}
]
```

### Tag Strategy

**Common Patterns:**

1. **Client/Server Separation**
   ```
   all,client,server
   ```
   Use for separating sensitive data (server) from client-exposed data.

2. **Environment Separation**
   ```
   all,dev,staging,production
   ```
   Different data for different environments.

3. **Feature Flags**
   ```
   all,premium,free
   ```
   Different data based on user tier.

4. **Platform Separation**
   ```
   all,ios,android,web
   ```
   Platform-specific configurations.

**Best Practices:**

- Use `all` tag for universally included fields
- Keep tag names lowercase
- Use descriptive tag names
- Document your tag strategy
- ID column should typically be tagged `all`

---

## Rules and Constraints

### Schema Rules

1. **Required Rows**
   - Must have at least 5 rows (metadata + schema)
   - Row 0: Metadata
   - Row 1: Tags
   - Row 2: Field names
   - Row 3: Field types
   - Row 4: Comments (can be empty)

2. **ID Column (Column 0)**
   - **Must exist** in every sheet
   - **Must be first column**
   - **Type:** Only `int`, `long`, or `string`
   - **Cannot be nested** (no dots in name)
   - **Cannot be empty** in data rows (empty rows skipped)

3. **Field Names**
   - **Case sensitive**
   - **Dot notation** creates nesting: `A.B.C`
   - **`[]` prefix** creates multi-line arrays: `[]Items`
   - **Cannot start with `#`** (reserved for comments)
   - **Cannot be empty** (empty columns skipped)
   - **No duplicate field names** at the same level

4. **Field Types**
   - **Must be valid type:** `int`, `long`, `float`, `bool`, `string`, `time`, `json`
   - **`[]` prefix for cell arrays:** `[]int`, `[]string`, etc.
   - **`struct` type** automatically inferred (don't specify)
   - **No `[]json`** (invalid)

### Composition Rules

1. **Multi-Line Arrays (MLA)**
   - **Prefix with `[]`:** `[]Items`, `[]Rewards.Type`
   - **Duplicate IDs required** for multiple elements
   - **Cannot nest MLAs:** `[][]Items` is invalid
   - **Empty row detection:** If all MLA fields in a row are empty, element is skipped
   - **Independent arrays:** Multiple MLA fields create separate arrays

2. **Nested Structures**
   - **Unlimited depth:** `A.B.C.D.E...` is valid
   - **Mixed types:** Can combine primitives and structs at any level
   - **Shared parents:** `A.B` and `A.C` share the same `A` parent

3. **Named Struct Mapping**
   - **Field identifier:** Use field name without `[]` prefix
   - **Regex support:** `/pattern/` for matching multiple fields
   - **Struct consistency:** All fields mapped to the same struct must have identical structure
   - **Omit parent names:** When identifying nested fields

### Data Rules

1. **Data Rows (Row 5+)**
   - **ID required:** Empty ID rows are skipped
   - **Comment rows:** ID starting with `#` skips row
   - **Duplicate IDs:** Only allowed with MLA fields
   - **Empty cells:** Become zero values for type

2. **Type Parsing**
   - **Int:** Empty → `0`, invalid → error
   - **Long:** Empty → `0`, invalid → error
   - **Float:** Empty → `0.0`, invalid → error
   - **Bool:** Empty → `false`, invalid values → error
   - **String:** Empty → `""`, all values valid
   - **Time:** Empty → `0001-01-01T00:00:00Z`, invalid format → error
   - **JSON:** Empty → `null`, invalid JSON → error

3. **Cell Arrays**
   - **Comma-separated:** `value1,value2,value3`
   - **No nested commas:** Cannot escape commas in values
   - **Empty cell:** Produces empty array `[]`
   - **All values must be valid** for the type

### Metadata Rules

1. **as_map**
   - **Mutually exclusive** with `sort_asc_by` and `sort_desc_by`
   - **Keys are string IDs:** Even if ID type is int

2. **Sorting**
   - **Cannot use both** `sort_asc_by` and `sort_desc_by`
   - **Field must exist** in table
   - **Field cannot be array** (MLA or cell array)
   - **Invalid types:** `bool`, `json`, nested fields
   - **Valid types:** `int`, `long`, `float`, `string`, `time`

3. **Struct Mapping**
   - **Syntax:** `struct=FieldId:StructName`
   - **Multiple mappings:** Separated by `&`
   - **Regex patterns:** Enclosed in `/regex/`
   - **Struct consistency:** Mapped fields must have identical structure

### Validation Errors

Common errors and their causes:

| Error | Cause | Fix |
|-------|-------|-----|
| `invalid table data` | Less than 5 rows | Add missing schema rows |
| `no columns in the csv file` | All columns are empty or start with `#` | Add valid field names |
| `invalid index field` | ID field has dots | Remove dots from ID field name |
| `invalid index field type` | ID type is not int/long/string | Change ID type |
| `nested multi-line array is not allowed` | Two `[]` prefixes in a field path | Remove one `[]` prefix |
| `there is no multi-line array field but id is duplicated` | Duplicate IDs without MLA fields | Add `[]` prefix or remove duplicate IDs |
| `as_map and sort_by are mutually exclusive` | Both as_map and sorting enabled | Remove one option |
| `sort_by: field is array` | Trying to sort by array field | Sort by non-array field |
| `sort_by: invalid field type` | Sorting by bool or json | Use valid sort type |

---

## Examples

### Example 1: Simple Flat Data

**CSV:**
```csv
as_map=false
all,all,all
ID,Name,Age
int,string,int

1,Alice,25
2,Bob,30
```

**Output:**
```json
[
  {"ID": 1, "Name": "Alice", "Age": 25},
  {"ID": 2, "Name": "Bob", "Age": 30}
]
```

---

### Example 2: Nested Structure

**CSV:**
```csv
as_map=false
all,all,all,all,all
ID,Name,Profile.Age,Profile.City,Profile.Country
int,string,int,string,string

1,Alice,25,NYC,USA
2,Bob,30,London,UK
```

**Output:**
```json
[
  {
    "ID": 1,
    "Name": "Alice",
    "Profile": {
      "Age": 25,
      "City": "NYC",
      "Country": "USA"
    }
  },
  {
    "ID": 2,
    "Name": "Bob",
    "Profile": {
      "Age": 30,
      "City": "London",
      "Country": "UK"
    }
  }
]
```

---

### Example 3: Cell Arrays

**CSV:**
```csv
as_map=false
all,all,all
ID,Name,Tags
int,string,[]string

1,Alice,"gold,premium,vip"
2,Bob,silver
3,Charlie,
```

**Output:**
```json
[
  {"ID": 1, "Name": "Alice", "Tags": ["gold", "premium", "vip"]},
  {"ID": 2, "Name": "Bob", "Tags": ["silver"]},
  {"ID": 3, "Name": "Charlie", "Tags": []}
]
```

---

### Example 4: Multi-Line Array

**CSV:**
```csv
as_map=false
all,all,all
ID,[]Item,[]Quantity
int,string,int

1,Sword,2
1,Shield,1
2,Potion,5
2,Elixir,3
```

**Output:**
```json
[
  {
    "ID": 1,
    "Item": ["Sword", "Shield"],
    "Quantity": [2, 1]
  },
  {
    "ID": 2,
    "Item": ["Potion", "Elixir"],
    "Quantity": [5, 3]
  }
]
```

---

### Example 5: Array of Structs (MLA)

**CSV:**
```csv
as_map=false&struct=Items:Item
all,all,all,all
ID,Name,[]Items.Type,[]Items.Quantity
int,string,string,int

1,Order1,Sword,2
1,,Shield,1
2,Order2,Potion,5
```

**Output:**
```json
[
  {
    "ID": 1,
    "Name": "Order1",
    "Items": [
      {"Type": "Sword", "Quantity": 2},
      {"Type": "Shield", "Quantity": 1}
    ]
  },
  {
    "ID": 2,
    "Name": "Order2",
    "Items": [
      {"Type": "Potion", "Quantity": 5}
    ]
  }
]
```

---

### Example 6: Deep Nesting with Arrays

**CSV:**
```csv
as_map=false&struct=Items:Item
all,all,all,all,all,all
ID,User.Name,User.Level,[]Items.Type,[]Items.Stats.HP,[]Items.Stats.MP
int,string,int,string,int,int

1,Alice,10,Sword,0,0
1,,,Potion,50,30
2,Bob,15,Shield,0,0
```

**Output:**
```json
[
  {
    "ID": 1,
    "User": {
      "Name": "Alice",
      "Level": 10
    },
    "Items": [
      {
        "Type": "Sword",
        "Stats": {"HP": 0, "MP": 0}
      },
      {
        "Type": "Potion",
        "Stats": {"HP": 50, "MP": 30}
      }
    ]
  },
  {
    "ID": 2,
    "User": {
      "Name": "Bob",
      "Level": 15
    },
    "Items": [
      {
        "Type": "Shield",
        "Stats": {"HP": 0, "MP": 0}
      }
    ]
  }
]
```

---

### Example 7: Map Output

**CSV:**
```csv
as_map=true
all,all,all
ID,Name,Level
int,string,int

1,Alice,10
2,Bob,15
```

**Output:**
```json
{
  "1": {"ID": 1, "Name": "Alice", "Level": 10},
  "2": {"ID": 2, "Name": "Bob", "Level": 15}
}
```

---

### Example 8: Sorted Output

**CSV:**
```csv
as_map=false&sort_desc_by=Score
all,all,all,all
ID,Name,Score,Rank
int,string,int,string

1,Alice,100,Gold
2,Bob,150,Platinum
3,Charlie,75,Silver
```

**Output:**
```json
[
  {"ID": 2, "Name": "Bob", "Score": 150, "Rank": "Platinum"},
  {"ID": 1, "Name": "Alice", "Score": 100, "Rank": "Gold"},
  {"ID": 3, "Name": "Charlie", "Score": 75, "Rank": "Silver"}
]
```

---

### Example 9: All Data Types

**CSV:**
```csv
as_map=false
all,all,all,all,all,all,all
ID,Int,Long,Float,Bool,String,Time,Json
int,int,long,float,bool,string,time,json

1,42,9999999999,3.14,true,hello,2024-09-30 11:00:00,"{""key"":""value""}"
2,0,0,0.0,false,"",0001-01-01 00:00:00,
```

**Output:**
```json
[
  {
    "ID": 1,
    "Int": 42,
    "Long": 9999999999,
    "Float": 3.14,
    "Bool": true,
    "String": "hello",
    "Time": "2024-09-30T11:00:00Z",
    "Json": {"key": "value"}
  },
  {
    "ID": 2,
    "Int": 0,
    "Long": 0,
    "Float": 0.0,
    "Bool": false,
    "String": "",
    "Time": "0001-01-01T00:00:00Z",
    "Json": null
  }
]
```

---

### Example 10: Tag Filtering

**CSV:**
```csv
as_map=false
all,client,server,all
ID,ClientData,ServerData,SharedData
int,string,string,string

1,ClientValue1,ServerValue1,SharedValue1
2,ClientValue2,ServerValue2,SharedValue2
```

**Client Output (tags: [client]):**
```json
[
  {"ID": 1, "ClientData": "ClientValue1", "SharedData": "SharedValue1"},
  {"ID": 2, "ClientData": "ClientValue2", "SharedData": "SharedValue2"}
]
```

**Server Output (tags: [server]):**
```json
[
  {"ID": 1, "ServerData": "ServerValue1", "SharedData": "SharedValue1"},
  {"ID": 2, "ServerData": "ServerValue2", "SharedData": "SharedValue2"}
]
```

---

### Example 11: Complex Real-World Example

**CSV:**
```csv
as_map=false&sort_asc_by=ID&struct=Rewards:Reward&struct=/.*SKU.*/:SKU
server,client,"client,server","client,server",server,server,server,server,server
ID,Tags,[]SKU.Type,[]SKU.ID,[]Rewards.Type,[]Rewards.ParamValue.Str,[]Rewards.ParamType,[]Rewards.ParamValue.Int,[]Rewards.ParamValue.Float
int,[]string,string,string,string,string,string,int,float

1,"gold,package",Google,IAP_Google_1,Gold,,Int,10,
1,,,,Gear,Weapon,Str,,
2,dollar,Google,IAP_Google_2,Dollar,,Float,,0.5
2,,Apple,IAP_Apple_2,Dollar,,Float,,0.8
2,,,,Dollar,,Float,,0.9
```

**Client Output (tags: [client]):**
```json
[
  {
    "ID": 1,
    "Tags": ["gold", "package"],
    "SKU": [
      {"Type": "Google", "ID": "IAP_Google_1"}
    ]
  },
  {
    "ID": 2,
    "Tags": ["dollar"],
    "SKU": [
      {"Type": "Google", "ID": "IAP_Google_2"},
      {"Type": "Apple", "ID": "IAP_Apple_2"}
    ]
  }
]
```

**Server Output (tags: [server]):**
```json
[
  {
    "ID": 1,
    "Rewards": [
      {
        "Type": "Gold",
        "ParamType": "Int",
        "ParamValue": {
          "Str": "",
          "Int": 10,
          "Float": 0
        }
      },
      {
        "Type": "Gear",
        "ParamType": "Str",
        "ParamValue": {
          "Str": "Weapon",
          "Int": 0,
          "Float": 0
        }
      }
    ]
  },
  {
    "ID": 2,
    "Rewards": [
      {
        "Type": "Dollar",
        "ParamType": "Float",
        "ParamValue": {
          "Str": "",
          "Int": 0,
          "Float": 0.5
        }
      },
      {
        "Type": "Dollar",
        "ParamType": "Float",
        "ParamValue": {
          "Str": "",
          "Int": 0,
          "Float": 0.8
        }
      },
      {
        "Type": "Dollar",
        "ParamType": "Float",
        "ParamValue": {
          "Str": "",
          "Int": 0,
          "Float": 0.9
        }
      }
    ]
  }
]
```

**Generated Go Code (partial):**
```go
// Named structs (reusable)
type Reward struct {
    Type       string       `json:"Type"`
    ParamType  string       `json:"ParamType"`
    ParamValue RewardParamValue `json:"ParamValue"`
}

type RewardParamValue struct {
    Str   string  `json:"Str"`
    Int   int     `json:"Int"`
    Float float64 `json:"Float"`
}

type SKU struct {
    Type string `json:"Type"`
    ID   string `json:"ID"`
}

// Server table
type ComplexRow struct {
    ID      int      `json:"ID"`
    Rewards []Reward `json:"Rewards"`
}

// Client table
type ComplexClientRow struct {
    ID   int      `json:"ID"`
    Tags []string `json:"Tags"`
    SKU  []SKU    `json:"SKU"`
}
```

---

## Summary

Sheet composition in nestcsv enables powerful data modeling through:

1. **Schema Structure:** 5 required rows (metadata, tags, names, types, comments) + data
2. **Field Types:** 8 primitive types (int, long, float, bool, string, time, json, struct)
3. **Nested Structures:** Dot notation creates hierarchies (`A.B.C`)
4. **Multi-Line Arrays:** `[]` prefix enables array composition across rows
5. **Cell Arrays:** `[]type` creates comma-separated arrays in single cell
6. **Named Structs:** Map field identifiers to reusable struct types
7. **Metadata Control:** Query string options (as_map, sorting, struct mapping)
8. **Tag Filtering:** Generate different outputs from same source data

**Key Principles:**
- ID column (column 0) is mandatory and must be simple type
- Duplicate IDs only allowed with multi-line arrays
- Empty cells become zero values
- Rows starting with `#` or empty ID are skipped
- Flexible composition patterns can be combined
- Type-safe code generation from schema

**For more examples, see:**
- `/examples/functions/` - Feature demonstrations
- `/examples/downstream/` - Real-world usage
- Source code: `table_parser.go`, `table_data.go`, `table_field.go`

---

*Generated for nestcsv - A Go CLI tool for converting structured CSV data into nested JSON and type-safe code.*
