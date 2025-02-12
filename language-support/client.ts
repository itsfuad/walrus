import { workspace, ExtensionContext } from 'vscode';
import { LanguageClient, LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
    // Server options - configure to run Go executable
    const serverOptions: ServerOptions = {
        command: 'lsp_server',
        args: []
    };

    // Client options
    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'walrus' }],
        synchronize: {
            fileEvents: workspace.createFileSystemWatcher('**/*.wal')
        }
    };

    // Create the language client
    client = new LanguageClient(
        'walrusLanguageServer',
        'Walrus Language Server',
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
