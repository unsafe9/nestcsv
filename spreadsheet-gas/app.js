function escapeCSVCell(v, tz) {
  if (v === undefined || v === null || v === "") {
    return "";
  }
  if (v instanceof Date) {
    v = Utilities.formatDate(v, tz, "yyyy-MM-dd HH:mm:ss");
  }
  if (/[,"\n\r\t]/g.test(v)) {
    v = '"' + v.replace(/"/g, '""') + '"';
  }
  return v;
}

function readSheetsAsCSV(res, fileId) {
  const spreadsheet = SpreadsheetApp.openById(fileId);
  const tz = spreadsheet.getSpreadsheetTimeZone();
  const sheets = SpreadsheetApp.openById(fileId).getSheets();
  
  for (const sheet of sheets) {
    if (sheet.getName().startsWith("#")) {
      continue;
    }
    const values = sheet.getDataRange().getValues();
    const csvLines = values.map(row => row.map(cell => escapeCSVCell(cell, tz)).join(","));
    res[sheet.getName()] = csvLines.join("\n");
  }
}

function readFolderSheets(res, folderId) {
  const parent = DriveApp.getFolderById(folderId);

  const files = parent.getFilesByType(MimeType.GOOGLE_SHEETS);
  while (files.hasNext()) {
    const file = files.next();
    readSheetsAsCSV(res, file.getId());
  }

  const folders = parent.getFolders();
  while (folders.hasNext()) {
    const subfolder = folders.next();
    readFolderSheets(res, subfolder.getId());
  }
}

function getParameters(e, name) {
  const value = e.parameters[name];
  if (Array.isArray(value)) {
    return value;
  }
  if (value) {
    return [value];
  }
  return [];
}

function doGet(e) {
  // Simply use a password to prevent unauthorized access.
  // You can use more secure methods such as OAuth2.
  const password = PropertiesService.getScriptProperties().getProperty("password");
  if (e.parameter.password !== password) {
    throw new Error("Forbidden");
  }

  const fileIds = getParameters(e, "fileIds");
  const folderIds = getParameters(e, "folderIds");
  let res = {};

  for (const fileId of fileIds) {
    readSheetsAsCSV(res, fileId);
  }
  for (const folderId of folderIds) {
    readFolderSheets(res, folderId);
  }
  
  const csvBlobs = [];
  for (const [name, csv] of Object.entries(res)) {
    csvBlobs.push(Utilities.newBlob(csv, "text/csv", name + ".csv"));
  }
  const zipBlob = Utilities.zip(csvBlobs, "cricket-tables.zip");
  const zipBase64 = Utilities.base64Encode(zipBlob.getBytes());

  return ContentService.createTextOutput(zipBase64).setMimeType(ContentService.MimeType.TEXT);
}