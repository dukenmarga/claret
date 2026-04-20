// @ts-check
// Manual update to unblock frontend build

export function SelectServiceAccountFile() {
    return window['go']['main']['App']['SelectServiceAccountFile']();
}

export function ParseServiceAccount(arg1) {
    return window['go']['main']['App']['ParseServiceAccount'](arg1);
}

export function ListDatabases(arg1) {
    return window['go']['main']['App']['ListDatabases'](arg1);
}

export function GetConfig() {
    return window['go']['main']['App']['GetConfig']();
}

export function AddProject(arg1) {
    return window['go']['main']['App']['AddProject'](arg1);
}

export function RemoveProject(arg1) {
    return window['go']['main']['App']['RemoveProject'](arg1);
}

export function ConnectProject(arg1, arg2) {
    return window['go']['main']['App']['ConnectProject'](arg1, arg2);
}

export function GetCollections(arg1) {
    return window['go']['main']['App']['GetCollections'](arg1);
}

export function GetDocuments(arg1, arg2) {
    return window['go']['main']['App']['GetDocuments'](arg1, arg2);
}

export function UpdateDocument(arg1, arg2) {
    return window['go']['main']['App']['UpdateDocument'](arg1, arg2);
}

export function DeleteDocuments(arg1) {
    return window['go']['main']['App']['DeleteDocuments'](arg1);
}

export function GetUsers(arg1) {
    return window['go']['main']['App']['GetUsers'](arg1);
}

export function DeleteUser(arg1) {
    return window['go']['main']['App']['DeleteUser'](arg1);
}
