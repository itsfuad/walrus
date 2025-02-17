import * as path from "path";
import { workspace, ExtensionContext } from "vscode";
import * as net from "net";
import * as cp from "child_process";

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  StreamInfo,
} from "vscode-languageclient/node";

let client: LanguageClient;

function setup(childProcess: cp.ChildProcessWithoutNullStreams): Promise<StreamInfo> {
  return new Promise<StreamInfo>((resolve, reject) => {
    let portData: string = '';
    // Read the port number from stdout
    childProcess.stdout.on('data', (data: Buffer) => {
      portData += data.toString();
      const match = RegExp(/PORT:(\d+)/).exec(portData);
      if (match) {
        const portNumber = parseInt(match[1]);
        
        // Create TCP connection
        const socket = new net.Socket();
        socket.connect(portNumber, "127.0.0.1", () => {
          resolve({
            reader: socket,
            writer: socket
          });
        });
        
        socket.on('error', (err) => {
          reject(err);
        });
      }
    });
    
    childProcess.stderr.on('data', (data: Buffer) => {
      console.error(`LSP Server Error: ${data.toString()}`);
    });
    
    childProcess.on('error', (err) => {
      reject(err);
    });
  });
}

export function activate(context: ExtensionContext) {
  const serverModule = context.asAbsolutePath(
    path.join("bin", "walrus-lsp.exe")
  );

  const serverOptions: ServerOptions = async () => {
    const childProcess = cp.spawn(serverModule);
    return setup(childProcess);
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: 'walrus' }],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher('**/*.{wal,walrus}'),
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