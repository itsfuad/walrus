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
    path.join("bin", "walrus-lsp.exe")
  );
  // If the extension is launched in debug mode then the debug server options are used
  // Otherwise the run options are used
  const serverOptions: ServerOptions = {
    run: { command: serverExec, transport: TransportKind.stdio },
    debug: {
      command: serverExec,
      transport: TransportKind.stdio,
    },
  };

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: 'walrus' }],
    synchronize: {
      // Notify the server about file changes to .wal files contained in the workspace
      fileEvents: workspace.createFileSystemWatcher('**/*.{wal,walrus}'),
    },
  };

  // Create the language client and start the client.
  client = new LanguageClient(
    "walrusLanguageServer",
    "Walrus Language Server",
    serverOptions,
    clientOptions
  );

  // Start the client. This will also launch the server
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}