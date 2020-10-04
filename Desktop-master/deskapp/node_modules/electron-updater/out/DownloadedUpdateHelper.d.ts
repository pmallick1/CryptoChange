import { UpdateInfo } from "builder-util-runtime";
import { ResolvedUpdateFileInfo } from "./main";
/** @private **/
export declare class DownloadedUpdateHelper {
    private setupPath;
    private _packagePath;
    private versionInfo;
    private fileInfo;
    readonly file: string | null;
    readonly packagePath: string | null;
    getDownloadedFile(versionInfo: UpdateInfo, fileInfo: ResolvedUpdateFileInfo): string | null;
    setDownloadedFile(file: string, packagePath: string | null, versionInfo: UpdateInfo, fileInfo: ResolvedUpdateFileInfo): void;
    clear(): void;
}
