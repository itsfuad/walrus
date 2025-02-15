import * as path from "path";
import { workspace, ExtensionContext } from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  const serverExec = context.asAbsolutePath(
    path.join("bin", "walrus-lsp.exe") // Ensure the server binary is here
  );

  const serverOptions: ServerOptions = {
    run: { command: serverExec, transport: TransportKind.stdio },
    debug: { command: serverExec, transport: TransportKind.stdio },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "walrus" }],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher("**/*.{wal,walrus}"),
    },
    middleware: {
      handleDiagnostics: (uri, diagnostics, next) => {
        console.log("Received diagnostics:", uri.toString(), diagnostics);
        next(uri, diagnostics); // Pass diagnostics to VSCode
      },
    },
  };

  client = new LanguageClient(
    "walrusLanguageServer",
    "Walrus Language Server",
    serverOptions,
    clientOptions
  );

  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
