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
  - tags: [client]
    ue5:
      root_dir: ./ue5
      prefix: Nest
    
```

Run the following command:
```bash
nestcsv

# specify your config file
nestcsv -c ../config/config.yaml
```

## How to structure the schema
See [examples](./examples)

## Roadmap
### Docs
- [ ] Add a csv(datasource) schema structure guide
- [ ] Add an example of UE5 json file loading
### Datasource
- [ ] Implement Google OAuth2 authentication for Google Apps Script
- [ ] Integrate spreadsheet datasource using Sheets API
### Config
- [ ] Extract time format settings into the configuration file
### Output
- [ ] Generate SQL dump file
### Code generation
- [ ] Generate Protobuf schema
- [ ] Generate Node.js code with type definitions
- [ ] Generate PostgreSQL DDL
- [ ] Generate MySQL DDL
