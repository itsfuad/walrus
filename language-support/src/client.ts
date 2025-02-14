import path from 'path';
import * as vscode from 'vscode';
import { LanguageClient, LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
	console.log("Activating Walrus Language Extension...");
  
  const serverExe = context.asAbsolutePath(path.join('bin', 'walrus-lsp.exe'));
	const serverOptions: ServerOptions = {
		run: { command: serverExe, args: [] },
		debug: { command: serverExe, args: [] }
	};

	const clientOptions: LanguageClientOptions = {
		documentSelector: [{ scheme: 'file', language: 'walrus' }],
		synchronize: {
			fileEvents: vscode.workspace.createFileSystemWatcher('**/*.{wal,walrus}')
		}
	};

	client = new LanguageClient(
		'walrusLanguageServer',
		'Walrus LSP', // changed from "Walrus Language Server"
		serverOptions,
		clientOptions
	);
	
	context.subscriptions.push({
		dispose: () => client.stop()
	});
	  
	client.start();
}

export function deactivate(): Thenable<void> | undefined {
	return client ? client.stop() : undefined;
}