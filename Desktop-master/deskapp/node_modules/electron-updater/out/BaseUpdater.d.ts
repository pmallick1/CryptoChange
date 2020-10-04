import { AllPublishOptions, DownloadOptions } from "builder-util-runtime";
import { AppUpdater } from "./AppUpdater";
import { DownloadedUpdateHelper } from "./DownloadedUpdateHelper";
import { ResolvedUpdateFileInfo } from "./main";
export declare abstract class BaseUpdater extends AppUpdater {
    protected readonly downloadedUpdateHelper: DownloadedUpdateHelper;
    protected quitAndInstallCalled: boolean;
    private quitHandlerAdded;
    constructor(options?: AllPublishOptions | null, app?: any);
    quitAndInstall(isSilent?: boolean, isForceRunAfter?: boolean): void;
    protected executeDownload(downloadOptions: DownloadOptions, fileInfo: ResolvedUpdateFileInfo, task: (tempDir: string, destinationFile: string, removeTempDirIfAny: () => Promise<any>) => Promise<any>): Promise<void>;
    protected abstract doInstall(installerPath: string, isSilent: boolean, isRunAfter: boolean): boolean;
    protected install(isSilent: boolean, isRunAfter: boolean): boolean;
    protected addQuitHandler(): void;
}
