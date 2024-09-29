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
datasource:
  spreadsheet_gas:
    url: <YOUR_GOOGLE_APPS_SCRIPT_WEB_APP_ENDPOINT>
    password: <YOUR_GOOGLE_APPS_SCRIPT_WEB_APP_PASSWORD>
    google_drive_folder_ids:
      - <YOUR_GOOGLE_DRIVE_FOLDER_ID>
    spreadsheet_file_ids:
      - <YOUR_GOOGLE_SPREADSHEET_FILE_ID>
    debug_save_dir: ./debug
  excel:
    directories:
      - <YOUR_EXCEL_FILE_DIRECTORY>
    extensions:
      - .xlsx
    files:
      - <YOUR_EXCEL_FILE>
    debug_save_dir: ./debug
  #local_file:
  #  root_dir: ./debug

output:
  indent: "  "
  root_dir: ./output
  as_map: false
  drop_id: false
```

Run the following command:
```bash
nestcsv

# specify the config file
nestcsv -c ../config/config.yaml
```

## How to structure the schema
See [examples](./examples)

## 0.0.1 Roadmap
- [ ] Add schema structure guide
- [ ] Add google oauth2 authentication for google apps scripts
- [ ] Add spreadsheet datasource using sheets api
- [ ] Extract time format into the config file
- [ ] Add code generation for the schema
- 