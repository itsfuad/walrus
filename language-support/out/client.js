"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode = __importStar(require("vscode"));
const node_1 = require("vscode-languageclient/node");
let client;
function activate(context) {
    console.log("Activating Walrus Language Extension...");
    const serverExe = context.asAbsolutePath('../lsp/walrus-lsp.exe');
    const serverOptions = {
        run: { command: serverExe, args: [] },
        debug: { command: serverExe, args: [] }
    };
    const clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'walrus' }],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.{wal,walrus}')
        }
    };
    client = new node_1.LanguageClient('walrusLanguageServer', 'Walrus Language Server', serverOptions, clientOptions);
    context.subscriptions.push(client.start());
}
exports.activate = activate;
function deactivate() {
    return client ? client.stop() : undefined;
}
exports.deactivate = deactivate;
//# sourceMappingURL=client.js.map