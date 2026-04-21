# nestcsv
A Go CLI tool that analyzes CSV data based on a predefined schema and converts it into a nested JSON structure.

## Installation
```bash
go install github.com/unsafe9/nestcsv/cmd/nestcsv@latest
```

## Usage
Compose your configurations:
```yaml
# nestcsv.yaml
datasources:
  - spreadsheet_gas:
      url: <YOUR_GOOGLE_APPS_SCRIPT_WEB_APP_ENDPOINT>
      password: <YOUR_GOOGLE_APPS_SCRIPT_WEB_APP_PASSWORD>
      google_drive_folder_ids:
        - <YOUR_GOOGLE_DRIVE_FOLDER_ID>
      spreadsheet_file_ids:
        - <YOUR_GOOGLE_SPREADSHEET_FILE_ID>
      debug_save_dir: ./debug
  - excel:
      patterns:
        - ./datasource/*.xlsx
      debug_save_dir: ./debug
  - csv:
      patterns:
        - ./datasource/*.csv
        #- ./debug/*.csv

outputs:
  - tags: [server, client]
    json:
      root_dir: ./output
      indent: "  "

codegens:
  - tags: [server]
    go:
      root_dir: ./go
      package_name: table
      file_suffix: ".gen.go"  # optional, default ".go"
  - tags: [client]
    ue5:
      root_dir: ./ue5
      prefix: Nest
      file_suffix: ".gen.h"        # optional, default ".h"
  - tags: [client]
    unity:
      root_dir: ./unity
      namespace: MyGame.Tables
      singleton: true
      data_suffix: Data            # optional, e.g. FooData
      table_suffix: DB             # optional, default "Table" (e.g. FooDB)
      resource_folder: MetaData    # optional, enables {Foo}DB.inst() auto-load from Resources/MetaData/foo.json
      file_suffix: ".gen.cs"       # optional, default ".cs"
    
```

Run the following command:
```bash
nestcsv

# specify your config file
nestcsv -c ../config/config.yaml
```

## How to structure the schema
Every table (CSV sheet / spreadsheet tab) must have a 5-row header, followed by the data rows:

| Row | Purpose | Notes |
|-----|---------|-------|
| 0 | Metadata query | Placed in column 0 only. Query-string syntax (see below). Leave empty if no options are needed. |
| 1 | Tags | Comma-separated tags per column. Used by `outputs`/`codegens` to filter which fields to emit. |
| 2 | Field names | Supports `.` for struct nesting and a leading `[]` for multi-line arrays (see below). |
| 3 | Field types | One of `int`, `long`, `float`, `bool`, `string`, `time`, `json`. Prefix with `[]` for a cell-level array. |
| 4 | Description | Free-form comments. Ignored by the parser. |
| 5+ | Data | Actual rows. Column 0 is the row ID and must be `int`, `long`, or `string`. |

### Column / row drop rules
- Column 0 (the ID column) of a data row is empty or starts with `#` → the row is skipped.
- A field name (row 2) is empty or starts with `#` → the entire column is dropped.

### Metadata query (row 0, column 0)
Written as a URL-style query string. Available keys:

| Key | Value | Description |
|-----|-------|-------------|
| `as_map` | `true` \| `false` | Emit the table as a map keyed by ID instead of an array. Mutually exclusive with `sort_*_by`. |
| `sort_asc_by` | field name | Sort the output array by the given field (ascending). Cannot be a `json`, `bool`, or array field. |
| `sort_desc_by` | field name | Same as above, descending. |
| `struct` | `<fieldId>:<TypeName>` | Promote a nested object to a **named struct** that is emitted as its own type and can be shared across tables (see below). Wrap the id in `/.../` to match by regex. Repeatable. |

Example:
```
as_map=false&sort_asc_by=ID&struct=Rewards:Reward&struct=/.*SKU.*/:SKU
```

### Nesting & arrays
- **Struct nesting** — use `.` in the field name. `A.B.C` creates `{ "A": { "B": { "C": ... } } }`.
- **Cell array** — prefix the _type_ with `[]`. The cell value is split by `,` (e.g. type `[]int` with cell `1,2,3`).
- **Multi-line array** — prefix the _field name_ with `[]`. Rows that share the same ID are grouped, and the `[]`-prefixed field collects one element per row. Works with struct nesting (e.g. `[]Rewards.Type`). Nested multi-line arrays are not allowed.

### Anonymous vs. named structs
By default a `.`-nested object is emitted as an **anonymous struct** — an auto-named, per-table type (e.g. `Item_Rewards_ParamValue`). Two tables with the same shape still get two unrelated types.

Use `struct=<fieldId>:<TypeName>` in row 0 to promote it to a **named struct**: it is generated as its own top-level type, and tables that map a field to the **same `TypeName` share one type** (so you can write code that takes any `Reward`). Shapes must match across tables — otherwise codegen fails with `named struct "X" has different fields`. `/regex/` covers many fields at once (`struct=/.*SKU.*/:SKU`), and the field id drops the `[]` prefix and any enclosing named-struct path (see `table_metadata.go:14-20`).

See [examples/functions/csv](./examples/functions/csv) for a working demo and the JSON output it produces.

## Roadmap
### Docs
- [ ] Add an example of UE5 json file loading
### Datasource
- [ ] Implement Google OAuth2 authentication for Google Apps Script
- [ ] Integrate spreadsheet datasource using Sheets API
### Config
- [ ] Extract time format settings into the configuration file
### Output
- [ ] Generate SQL dump file
### Code generation
- [x] Generate Unity (C#) code (Unity 6, Newtonsoft.Json)
- [ ] Generate Protobuf schema
- [ ] Generate Rust code
- [ ] Generate Node.js code with type definitions
- [ ] Generate PostgreSQL DDL
- [ ] Generate MySQL DDL
