import * as models from './models';

export interface go {
  "main": {
    "App": {
		Greet(arg1:string):Promise<string>
		OpenLinkInLocalBrowser(arg1:string):Promise<void>
		ScanNearbyNodes():Promise<string>
		TryConnect(arg1:Array<string>,arg2:string,arg3:string,arg4:boolean):Promise<Array<string>>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
